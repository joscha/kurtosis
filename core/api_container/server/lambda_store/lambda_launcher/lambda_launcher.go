/*
 * Copyright (c) 2021 - present Kurtosis Technologies Inc.
 * All Rights Reserved.
 */

package lambda_launcher

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/kurtosis-tech/kurtosis-client/golang/kurtosis_core_rpc_api_consts"
	"github.com/kurtosis-tech/kurtosis-lambda-api-lib/golang/kurtosis_lambda_docker_api"
	"github.com/kurtosis-tech/kurtosis-lambda-api-lib/golang/kurtosis_lambda_rpc_api_bindings"
	"github.com/kurtosis-tech/kurtosis-lambda-api-lib/golang/kurtosis_lambda_rpc_api_consts"
	"github.com/kurtosis-tech/kurtosis/api_container/server/lambda_store/lambda_store_types"
	"github.com/kurtosis-tech/kurtosis/commons"
	"github.com/kurtosis-tech/kurtosis/commons/docker_manager"
	"github.com/kurtosis-tech/kurtosis/commons/object_name_providers"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"strconv"
	"time"
)

const (
	waitForLambdaAvailabilityTimeout = 10 * time.Second
)

type LambdaLauncher struct {
	dockerManager *docker_manager.DockerManager

	// Lambdas have a connection to the API container, so the launcher must know about the API container's IP addr
	apiContainerIpAddr string

	enclaveObjNameProvider *object_name_providers.EnclaveObjectNameProvider

	freeIpAddrTracker *commons.FreeIpAddrTracker

	shouldPublishPorts bool

	dockerNetworkId string

	enclaveDataVolName string
}

func NewLambdaLauncher(dockerManager *docker_manager.DockerManager, apiContainerIpAddr string, enclaveObjNameProvider *object_name_providers.EnclaveObjectNameProvider, freeIpAddrTracker *commons.FreeIpAddrTracker, shouldPublishPorts bool, dockerNetworkId string, enclaveDataVolName string) *LambdaLauncher {
	return &LambdaLauncher{dockerManager: dockerManager, apiContainerIpAddr: apiContainerIpAddr, enclaveObjNameProvider: enclaveObjNameProvider, freeIpAddrTracker: freeIpAddrTracker, shouldPublishPorts: shouldPublishPorts, dockerNetworkId: dockerNetworkId, enclaveDataVolName: enclaveDataVolName}
}

func (launcher LambdaLauncher) Launch(
		ctx context.Context,
		lambdaId lambda_store_types.LambdaID,
		containerImage string,
		serializedParams string) (newContainerId string, newContainerIpAddr net.IP, client kurtosis_lambda_rpc_api_bindings.LambdaServiceClient, lambdaPortHostPortBindng *nat.PortBinding, resultErr error) {

	lambdaPortNumStr := strconv.Itoa(kurtosis_lambda_rpc_api_consts.ListenPort)
	lambdaPortObj, err := nat.NewPort(kurtosis_lambda_rpc_api_consts.ListenProtocol, lambdaPortNumStr)
	if err != nil {
		return "", nil, nil, nil, stacktrace.Propagate(
			err,
			"An error occurred creating port object for Lambda port %v/%v",
			kurtosis_lambda_rpc_api_consts.ListenProtocol,
			kurtosis_lambda_rpc_api_consts.ListenPort,
		)
	}
	usedPorts := map[nat.Port]bool {
		lambdaPortObj: true,
	}

	lambdaIpAddr, err := launcher.freeIpAddrTracker.GetFreeIpAddr()
	if err != nil {
		return "", nil, nil, nil, stacktrace.Propagate(err, "An error occurred getting a free IP address for new module")
	}

	apiContainerSocket := fmt.Sprintf("%v:%v", launcher.apiContainerIpAddr, kurtosis_core_rpc_api_consts.ListenPort)
	envVars := map[string]string{
		kurtosis_lambda_docker_api.ApiContainerSocketEnvVar: apiContainerSocket,
		kurtosis_lambda_docker_api.SerializedCustomParamsEnvVar: serializedParams,
	}

	volumeMounts := map[string]string{
		launcher.enclaveDataVolName: kurtosis_lambda_docker_api.ExecutionVolumeMountpoint,
	}

	containerId, allHostPortBindings, err := launcher.dockerManager.CreateAndStartContainer(
		ctx,
		containerImage,
		launcher.enclaveObjNameProvider.ForLambdaContainer(lambdaId),
		false,	// Lambda containers don't run in interactive mode
		launcher.dockerNetworkId,
		lambdaIpAddr,
		map[docker_manager.ContainerCapability]bool{}, // No extra capapbilities needed for modules
		docker_manager.DefaultNetworkMode,
		usedPorts,
		launcher.shouldPublishPorts,
		nil, // No ENTRYPOINT overrides; modules are configured using env vars
		nil, // No CMD overrides; modules are configured using env vars
		envVars,
		nil, // No bind mounts needed
		volumeMounts,
		false, // Lambdas shouldn't have access to the host machine, for security purposes!
	)
	if err != nil {
		return "", nil, nil, nil, stacktrace.Propagate(err, "An error occurred launching the module container")
	}
	shouldDestroyContainer := true
	defer func() {
		if shouldDestroyContainer {
			if err := launcher.dockerManager.KillContainer(context.Background(), containerId); err != nil {
				logrus.Error("Launching the lambda container failed, but an error occurred killing container we started:")
				fmt.Fprintln(logrus.StandardLogger().Out, err)
				logrus.Errorf("ACTION REQUIRED: You'll need to manually kill container with ID '%v'", containerId)
			}
		}
	}()

	var resultHostPortBinding *nat.PortBinding = nil
	hostPortBindingFromMap, found := allHostPortBindings[lambdaPortObj]
	if found {
		resultHostPortBinding = hostPortBindingFromMap
	}

	lambdaSocket := fmt.Sprintf("%v:%v", lambdaIpAddr, kurtosis_lambda_rpc_api_consts.ListenPort)
	conn, err := grpc.Dial(
		lambdaSocket,
		grpc.WithInsecure(), // TODO SECURITY: Use HTTPS to verify we're connecting to the correct lambda
	)
	if err != nil {
		return "", nil, nil, nil, stacktrace.Propagate(err, "Couldn't dial lambda container '%v' at %v", lambdaId, lambdaSocket)
	}
	lambdaClient := kurtosis_lambda_rpc_api_bindings.NewLambdaServiceClient(conn)

	logrus.Debugf("Waiting for lambda container to become available...")
	if err := waitUntilLambdaContainerIsAvailable(ctx, lambdaClient); err != nil {
		return "", nil, nil, nil, stacktrace.Propagate(err, "An error occurred while waiting for lambda container '%v' to become available", lambdaId)
	}
	logrus.Debugf("Lambda container '%v' became available", lambdaId)

	shouldDestroyContainer = false
	return containerId, lambdaIpAddr, lambdaClient, resultHostPortBinding, nil
}

func waitUntilLambdaContainerIsAvailable(ctx context.Context, client kurtosis_lambda_rpc_api_bindings.LambdaServiceClient) error {
	contextWithTimeout, cancelFunc := context.WithTimeout(ctx, waitForLambdaAvailabilityTimeout)
	defer cancelFunc()
	if _, err := client.IsAvailable(contextWithTimeout, &emptypb.Empty{}, grpc.WaitForReady(true)); err != nil {
		return stacktrace.Propagate(err, "An error occurred waiting for the Lambda container to become available")
	}
	return nil
}
