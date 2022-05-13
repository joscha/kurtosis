/*
 * Copyright (c) 2021 - present Kurtosis Technologies Inc.
 * All Rights Reserved.
 */

package user_service_launcher

import (
	"context"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_interface"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_interface/objects/enclave"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_interface/objects/port_spec"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_interface/objects/service"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_interface/objects/user_service_registration"
	"github.com/kurtosis-tech/free-ip-addr-tracker-lib/lib"
	"github.com/kurtosis-tech/kurtosis-core/server/api_container/server/service_network/user_service_launcher/files_artifact_expander"
	"github.com/kurtosis-tech/stacktrace"
	"net"
)

/*
Convenience struct whose only purpose is launching user services
*/
type UserServiceLauncher struct {
	kurtosisBackend          backend_interface.KurtosisBackend
	filesArtifactExpander    *files_artifact_expander.FilesArtifactExpander
	freeIpAddrTracker        *lib.FreeIpAddrTracker
}

func NewUserServiceLauncher(kurtosisBackend backend_interface.KurtosisBackend, filesArtifactExpander *files_artifact_expander.FilesArtifactExpander, freeIpAddrTracker *lib.FreeIpAddrTracker) *UserServiceLauncher {
	return &UserServiceLauncher{kurtosisBackend: kurtosisBackend, filesArtifactExpander: filesArtifactExpander, freeIpAddrTracker: freeIpAddrTracker}
}

/**
Launches a testnet service with the given parameters

Returns:
	* The container ID of the newly-launched service
	* The mapping of used_port -> host_port_binding (if no host port is bound, then the value will be nil)
*/
func (launcher UserServiceLauncher) Launch(
	ctx context.Context,
	registrationGuid user_service_registration.UserServiceRegistrationGUID,
	enclaveId enclave.EnclaveID,
	ipAddr net.IP,
	imageName string,
	privatePorts map[string]*port_spec.PortSpec,
	entrypointArgs []string,
	cmdArgs []string,
	envVars map[string]string,
	// Mapping of UUIDs of previously-registered files artifacts -> mountpoints on the container
	// being launched
	filesArtifactUuidsToMountpoints map[service.FilesArtifactID]string,
) (
	resultUserService *service.Service,
	resultErr error,
) {
	usedArtifactUuidSet := map[service.FilesArtifactID]bool{}
	for artifactUuid := range filesArtifactUuidsToMountpoints {
		usedArtifactUuidSet[artifactUuid] = true
	}

	// First expand the files artifacts into volumes, so that any errors get caught early
	// NOTE: if users don't need to investigate the volume contents, we could keep track of the volumes we create
	//  and delete them at the end of the test to keep things cleaner
	artifactUuidsToVolumes, err := launcher.filesArtifactExpander.ExpandArtifactsIntoVolumes(ctx, serviceGUID, usedArtifactUuidSet)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred expanding the requested files artifacts into volumes")
	}


	artifactVolumeMounts := map[string]string{}
	for artifactUuid, mountpoint := range filesArtifactUuidsToMountpoints {
		artifactVolume, found := artifactUuidsToVolumes[artifactUuid]
		if !found {
			return nil, stacktrace.NewError(
				"Even though we declared that we need files artifact '%v' to be expanded, no volume containing the "+
					"expanded contents was found; this is a bug in Kurtosis",
				artifactUuid,
			)
		}
		artifactVolumeMounts[string(artifactVolume)] = mountpoint
	}

	launchedUserService, err := launcher.kurtosisBackend.CreateUserService(
		ctx,
		registrationGuid,
		imageName,
		enclaveId,
		ipAddr,
		privatePorts,
		entrypointArgs,
		cmdArgs,
		envVars,
		artifactVolumeMounts,
	)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred starting the container for user service in enclave '%v' with image '%v'", enclaveId, imageName)
	}
	return launchedUserService, nil
}