/*
 * Copyright (c) 2021 - present Kurtosis Technologies Inc.
 * All Rights Reserved.
 */

package enclave_manager

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/kurtosis-tech/kurtosis/commons/docker_manager"
	"github.com/kurtosis-tech/kurtosis/commons/docker_network_allocator"
	"github.com/kurtosis-tech/kurtosis/commons/enclave_manager/enclave_context"
	"github.com/kurtosis-tech/kurtosis/commons/object_name_providers"
	"github.com/kurtosis-tech/kurtosis/initializer/api_container_launcher"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
	"net"
	"time"
)

const (
	// The API container is responsible for disconnecting/stopping everything in its network when stopped, so we need
	//  to give it some time to do so
	apiContainerStopTimeout = 10 * time.Second
)

// Manages Kurtosis enclaves, and creates new ones in response to running tasks
type EnclaveManager struct {
	// Will be wrapped in the DockerManager that logs to the proper location
	dockerClient *client.Client

	// TODO Hide this all inside this class, rather than taking it in as a constructor param????
	dockerNetworkAllocator *docker_network_allocator.DockerNetworkAllocator

	apiContainerLauncher *api_container_launcher.ApiContainerLauncher
}

// TODO Constructor

// NOTE: thisContainerId is the ID of the container in which this code is executing, so that it can be mounted inside
//  the new enclave such that it can communicate with the API container
func (manager *EnclaveManager) CreateEnclave(
		setupCtx context.Context,
		log *logrus.Logger,
		thisContainerId string,
		enclaveId string,
		isPartitioningEnabled bool) (*enclave_context.EnclaveContext, error) {
	dockerManager := docker_manager.NewDockerManager(log, manager.dockerClient)

	matchingNetworks, err := dockerManager.GetNetworkIdsByName(setupCtx, enclaveId)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred finding enclaves with name '%v', which is necessary to ensure that our enclave doesn't exist yet", enclaveId)
	}
	if len(matchingNetworks) > 0 {
		return nil, stacktrace.NewError("Cannot create enclave '%v' because an enclave with that name already exists", enclaveId)
	}

	enclaveObjNameProvider := object_name_providers.NewEnclaveObjectNameProvider(enclaveId)

	teardownCtx := context.Background()  // Separate context for tearing stuff down in case the input context is cancelled

	log.Debugf("Creating Docker network for enclave '%v'...", enclaveId)
	networkId, networkIpAndMask, gatewayIp, freeIpAddrTracker, err := manager.dockerNetworkAllocator.CreateNewNetwork(
		setupCtx,
		dockerManager,
		log,
		enclaveId,
	)
	if err != nil {
		// TODO If the user Ctrl-C's while the CreateNetwork call is ongoing then the CreateNetwork will error saying
		//  that the Context was cancelled as expected, but *the Docker engine will still create the network*!!! We'll
		//  need to parse the log message for the string "context canceled" and, if found, do another search for
		//  networks with our network name and delete them
		return nil, stacktrace.Propagate(err, "An error occurred allocating a new network for enclave '%v'", enclaveId)
	}
	shouldDeleteNetwork := true
	defer func() {
		if shouldDeleteNetwork {
			if err := dockerManager.RemoveNetwork(teardownCtx, networkId); err != nil {
				log.Errorf("Creating the enclave didn't complete successfully, so we tried to delete network '%v' that we created but an error was thrown:")
				fmt.Fprintln(log.Out, err)
				log.Errorf("ACTION REQUIRED: You'll need to manually remove network with ID '%v'!!!!!!!", networkId)
			}
		}
	}()
	log.Debugf("Docker network '%v' created successfully with ID '%v' and subnet CIDR '%v'", enclaveId, networkId, networkIpAndMask.String())

	log.Debugf("Connecting this container to the enclave network so that it can interact with the containers in the enclave...")
	mountIp, err := freeIpAddrTracker.GetFreeIpAddr()
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred getting a free IP for mounting this container inside the enclave")
	}
	if err := dockerManager.ConnectContainerToNetwork(setupCtx, networkId, thisContainerId, mountIp); err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred connecting this container to the enclave network")
	}
	shouldDisconnectThisContainer := true
	defer func() {
		if shouldDisconnectThisContainer {
			if err := dockerManager.DisconnectContainerFromNetwork(teardownCtx, thisContainerId, networkId); err != nil {
				log.Errorf("Creating the enclave didn't complete successfully, so we tried to disconnect this container from enclave network but an error was thrown:")
				fmt.Fprintln(log.Out, err)
				log.Errorf("ACTION REQUIRED: You'll need to manually disconnect container with ID '%v' from network with ID '%v'!!!!!!!", thisContainerId, networkId)
			}
		}
	}()
	log.Debugf("Successfully connected this container to the enclave network so that it can interact with the containers in the enclave")

	// TODO use hostnames rather than IPs, which makes things nicer and which we'll need for Docker swarm support
	// We need to create the IP addresses for BOTH containers because the testsuite needs to know the IP of the API
	//  container which will only be started after the testsuite container
	apiContainerIpAddr, err := freeIpAddrTracker.GetFreeIpAddr()
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred getting an IP for the Kurtosis API container")
	}

	if err := dockerManager.CreateVolume(setupCtx, enclaveId); err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred creating enclave volume '%v'", enclaveId)
	}
	// NOTE: We could defer a deletion of this volume unless the function completes successfully - right now, Kurtosis
	//  doesn't do any volume deletion

	apiContainerName := enclaveObjNameProvider.ForApiContainer()
	apiContainerId, err := manager.apiContainerLauncher.Launch(
		setupCtx,
		log,
		dockerManager,
		apiContainerName,
		enclaveId,
		networkId,
		networkIpAndMask.String(),
		gatewayIp,
		apiContainerIpAddr,
		[]net.IP{},  // TODO Add the other containers that we mount in here
		isPartitioningEnabled,
		map[string]bool{thisContainerId: true},
	)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred launching the API container")
	}
	// The API container is started successfully and it will disconnect/stop everything in its network when it shuts down,
	//  so it takes over the responsibility of disconnecting this container
	shouldDisconnectThisContainer = false
	shouldStopApiContainer := true
	defer func() {
		if shouldStopApiContainer {
			if err := dockerManager.StopContainer(teardownCtx, apiContainerId, apiContainerStopTimeout); err != nil {
				log.Errorf("Creating the enclave didn't complete successfully, so we tried to stop the API container but an error was thrown:")
				fmt.Fprintln(log.Out, err)
				log.Errorf("ACTION REQUIRED: You'll need to manually stop API container with ID '%v'", apiContainerId)
			}
		}
	}()

	result := enclave_context.NewEnclaveContext(
		networkId,
		enclaveId,
		networkIpAndMask,
		apiContainerIpAddr,
	)

	// Everything started successfully, so the responsibility of deleting the network is now transferred to the caller
	shouldDeleteNetwork = false
	shouldStopApiContainer = false
	return result, nil
}

func DestroyEnclave(ctx context.Context, enclaveId string) {

}