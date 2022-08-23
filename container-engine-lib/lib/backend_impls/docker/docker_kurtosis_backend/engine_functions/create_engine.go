package engine_functions

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_impls/docker/docker_kurtosis_backend/consts"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_impls/docker/docker_kurtosis_backend/engine_functions/logs_components"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_impls/docker/docker_kurtosis_backend/engine_functions/logs_components/fluentbit"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_impls/docker/docker_kurtosis_backend/engine_functions/logs_components/loki"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_impls/docker/docker_kurtosis_backend/shared_helpers"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_impls/docker/docker_manager"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_impls/docker/docker_manager/types"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_impls/docker/object_attributes_provider"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_interface/objects/engine"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_interface/objects/port_spec"
	"github.com/kurtosis-tech/container-engine-lib/lib/operation_parallelizer"
	"github.com/kurtosis-tech/container-engine-lib/lib/uuid_generator"
	"github.com/kurtosis-tech/stacktrace"
	"github.com/sirupsen/logrus"
	"reflect"
	"strings"
	"time"
)

const (
	maxWaitForEngineAvailabilityRetries         = 10
	timeBetweenWaitForEngineAvailabilityRetries = 1 * time.Second

	getAllEngineContainersOperationId operation_parallelizer.OperationID= "getAllEngineContainers"
	getAllLogsDatabaseContainersOperationId operation_parallelizer.OperationID= "getAllLogsDatabaseContainers"
	getAllLogsCollectorContainersOperationId operation_parallelizer.OperationID= "getAllLogsCollectorContainers"
)

func CreateEngine(
	ctx context.Context,
	imageOrgAndRepo string,
	imageVersionTag string,
	grpcPortNum uint16,
	grpcProxyPortNum uint16,
	logsCollectorHttpPortNumber uint16,
	envVars map[string]string,
	dockerManager *docker_manager.DockerManager,
	objAttrsProvider object_attributes_provider.DockerObjectAttributesProvider,
) (
	*engine.Engine,
	error,
) {
	isThereEngineOrLogsComponentsContainersInTheCluster, existenceContainerIds,  err := isThereAnyOtherEngineOrLogsComponentsContainersInTheCluster(ctx, dockerManager)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred checking for engine containers and logs components containers existence")
	}
	canStartNewEngine := !isThereEngineOrLogsComponentsContainersInTheCluster
	if !canStartNewEngine {
		containerIdsStr := strings.Join(existenceContainerIds, ", ")
		return nil, stacktrace.NewError("No new engine won't be started because there exist an engine or logs component container in the cluster; the following containers with IDs '%v' should be removed before creating a new engine", containerIdsStr)
	}

	matchingNetworks, err := dockerManager.GetNetworksByName(ctx, consts.NameOfNetworkToStartEngineContainersIn)
	if err != nil {
		return nil, stacktrace.Propagate(
			err,
			"An error occurred getting networks matching the network we want to start the engine in, '%v'",
			consts.NameOfNetworkToStartEngineContainersIn,
		)
	}
	numMatchingNetworks := len(matchingNetworks)
	if numMatchingNetworks == 0 && numMatchingNetworks > 1 {
		return nil, stacktrace.NewError(
			"Expected exactly one network matching the name of the network that we want to start the engine in, '%v', but got %v",
			consts.NameOfNetworkToStartEngineContainersIn,
			numMatchingNetworks,
		)
	}
	targetNetwork := matchingNetworks[0]
	targetNetworkId := targetNetwork.GetId()

	engineGuidStr, err := uuid_generator.GenerateUUIDString()
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred generating a UUID string for the engine")
	}
	engineGuid := engine.EngineGUID(engineGuidStr)

	privateGrpcPortSpec, err := port_spec.NewPortSpec(grpcPortNum, consts.EnginePortProtocol)
	if err != nil {
		return nil, stacktrace.Propagate(
			err,
			"An error occurred creating the engine's private grpc port spec object using number '%v' and protocol '%v'",
			grpcPortNum,
			consts.EnginePortProtocol.String(),
		)
	}
	privateGrpcProxyPortSpec, err := port_spec.NewPortSpec(grpcProxyPortNum, consts.EnginePortProtocol)
	if err != nil {
		return nil, stacktrace.Propagate(
			err,
			"An error occurred creating the engine's private grpc proxy port spec object using number '%v' and protocol '%v'",
			grpcProxyPortNum,
			consts.EnginePortProtocol.String(),
		)
	}

	engineAttrs, err := objAttrsProvider.ForEngineServer(
		engineGuid,
		consts.KurtosisInternalContainerGrpcPortId,
		privateGrpcPortSpec,
		consts.KurtosisInternalContainerGrpcProxyPortId,
		privateGrpcProxyPortSpec,
	)
	if err != nil {
		return nil, stacktrace.Propagate(
			err,
			"An error occurred getting the engine server container attributes using GUID '%v', grpc port num '%v', and "+
				"grpc proxy port num '%v'",
			engineGuid,
			grpcPortNum,
			grpcProxyPortNum,
		)
	}

	privateGrpcDockerPort, err := shared_helpers.TransformPortSpecToDockerPort(privateGrpcPortSpec)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred transforming the private grpc port spec to a Docker port")
	}
	privateGrpcProxyDockerPort, err := shared_helpers.TransformPortSpecToDockerPort(privateGrpcProxyPortSpec)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred transforming the private grpc proxy port spec to a Docker port")
	}

	usedPorts := map[nat.Port]docker_manager.PortPublishSpec{
		privateGrpcDockerPort:      docker_manager.NewManualPublishingSpec(grpcPortNum),
		privateGrpcProxyDockerPort: docker_manager.NewManualPublishingSpec(grpcProxyPortNum),
	}

	bindMounts := map[string]string{
		// Necessary so that the engine server can interact with the Docker engine
		consts.DockerSocketFilepath: consts.DockerSocketFilepath,
	}

	containerImageAndTag := fmt.Sprintf(
		"%v:%v",
		imageOrgAndRepo,
		imageVersionTag,
	)

	labelStrs := map[string]string{}
	for labelKey, labelValue := range engineAttrs.GetLabels() {
		labelStrs[labelKey.GetString()] = labelValue.GetString()
	}

	createAndStartArgs := docker_manager.NewCreateAndStartContainerArgsBuilder(
		containerImageAndTag,
		engineAttrs.GetName().GetString(),
		targetNetworkId,
	).WithEnvironmentVariables(
		envVars,
	).WithBindMounts(
		bindMounts,
	).WithUsedPorts(
		usedPorts,
	).WithLabels(
		labelStrs,
	).Build()

	// Best-effort pull attempt
	if err = dockerManager.PullImage(ctx, containerImageAndTag); err != nil {
		logrus.Warnf("Failed to pull the latest version of engine server image '%v'; you may be running an out-of-date version", containerImageAndTag)
	}

	containerId, hostMachinePortBindings, err := dockerManager.CreateAndStartContainer(ctx, createAndStartArgs)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred starting the Kurtosis engine container")
	}
	shouldKillEngineContainer := true
	defer func() {
		if shouldKillEngineContainer {
			// NOTE: We use the background context here so that the kill will still go off even if the reason for
			// the failure was the original context being cancelled
			if err := dockerManager.KillContainer(context.Background(), containerId); err != nil {
				logrus.Errorf(
					"Launching the engine server with GUID '%v' and container ID '%v' didn't complete successfully so we "+
						"tried to kill the container we started, but doing so exited with an error:\n%v",
					engineGuid,
					containerId,
					err)
				logrus.Errorf("ACTION REQUIRED: You'll need to manually stop engine server with GUID '%v'!!!!!!", engineGuid)
			}
		}
	}()

	if err := shared_helpers.WaitForPortAvailabilityUsingNetstat(
		ctx,
		dockerManager,
		containerId,
		privateGrpcPortSpec,
		maxWaitForEngineAvailabilityRetries,
		timeBetweenWaitForEngineAvailabilityRetries,
	); err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred waiting for the engine server's grpc port to become available")
	}

	// TODO UNCOMMENT THIS ONCE WE HAVE GRPC-PROXY WIRED UP!!
	/*
		if err := waitForPortAvailabilityUsingNetstat(ctx, backend.dockerManager, containerId, grpcProxyPortNum); err != nil {
			return nil, stacktrace.Propagate(err, "An error occurred waiting for the engine server's grpc proxy port to become available")
		}
	*/

	result, err := getEngineObjectFromContainerInfo(containerId, labelStrs, types.ContainerStatus_Running, hostMachinePortBindings)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred creating an engine object from container with GUID '%v'", containerId)
	}

	killCentralizedLogsComponentsContainersAndVolumesFunc, err := createCentralizedLogsComponents(
		ctx,
		engineGuid,
		targetNetworkId,
		logsCollectorHttpPortNumber,
		objAttrsProvider,
		dockerManager,
	)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred creating the centralized logs components for the engine with GUID '%v' and network ID '%v'", engineGuid, targetNetworkId)
	}
	shouldKillCentralizedLogsComponentsContainers := true
	defer func() {
		if shouldKillCentralizedLogsComponentsContainers {
			killCentralizedLogsComponentsContainersAndVolumesFunc()
		}
	}()

	shouldKillEngineContainer = false
	shouldKillCentralizedLogsComponentsContainers = false
	return result, nil
}

// ====================================================================================================
// 									   Private helper methods
// ====================================================================================================
//TODO we can run it in parallel after the network creation, and we can wait before returning the EngineInfo object
func createCentralizedLogsComponents(
	ctx context.Context,
	engineGuid engine.EngineGUID,
	targetNetworkId string,
	logsCollectorHttpPortNumber uint16,
	objAttrsProvider object_attributes_provider.DockerObjectAttributesProvider,
	dockerManager *docker_manager.DockerManager,
) (func(), error) {

	logsDatabaseContainerConfigProvider := loki.CreateLokiContainerConfigProviderForKurtosis()

	logsDatabaseHost, logsDatabasePort, killLogsDatabaseContainerAndVolumeFunc, err := createLogsDatabaseContainer(
		ctx,
		engineGuid,
		targetNetworkId,
		objAttrsProvider,
		dockerManager,
		logsDatabaseContainerConfigProvider,
	)
	if err != nil {
		return nil, stacktrace.Propagate(
			err,
			"An error occurred creating the logs database container for engine with GUID '%v' in Docker network with ID '%v'",
			engineGuid,
			targetNetworkId,
		)
	}
	shouldKillLogsDatabaseContainerAndVolume := true
	defer func() {
		if shouldKillLogsDatabaseContainerAndVolume {
			killLogsDatabaseContainerAndVolumeFunc()
		}
	}()

	logsCollectorContainerConfigProvider := fluentbit.CreateFluentbitContainerConfigProviderForKurtosis(logsDatabaseHost, logsDatabasePort, logsCollectorHttpPortNumber)

	killLogsCollectorContainerAndVolumeFunc, err := createLogsCollectorContainer(
		ctx,
		engineGuid,
		targetNetworkId,
		objAttrsProvider,
		dockerManager,
		logsCollectorContainerConfigProvider,
	)
	if err != nil {
		return nil, stacktrace.Propagate(
			err,
			"An error occurred creating the logs collector container for engine with GUID '%v' in Docker network with ID '%v'",
			engineGuid,
			targetNetworkId,
		)
	}

	killCentralizedLogsComponentsContainersAndVolumesFunc := func() {
		killLogsDatabaseContainerAndVolumeFunc()
		killLogsCollectorContainerAndVolumeFunc()
	}

	shouldKillLogsDatabaseContainerAndVolume = false
	return killCentralizedLogsComponentsContainersAndVolumesFunc, nil
}

func createLogsDatabaseContainer(
	ctx context.Context,
	engineGuid engine.EngineGUID,
	targetNetworkId string,
	objAttrsProvider object_attributes_provider.DockerObjectAttributesProvider,
	dockerManager *docker_manager.DockerManager,
	logsDatabaseContainerConfigProvider logs_components.LogsDatabaseContainerConfigProvider,
) (
	resultLogsDatabasePrivateHost string,
	resultLogsDatabasePrivatePort uint16,
	resultKillLogsDatabaseContainerFunc func(),
	resultErr error,
) {

	//Create the volume first
	logsDbVolumeAttrs, err := objAttrsProvider.ForLogsDatabaseVolume()
	if err != nil {
		return "", 0, nil, stacktrace.Propagate(err, "An error occurred getting the logs database volume attributes for engine with GUID %v", engineGuid)
	}
	volumeName := logsDbVolumeAttrs.GetName().GetString()
	volumeLabelStrs := map[string]string{}
	for labelKey, labelValue := range logsDbVolumeAttrs.GetLabels() {
		volumeLabelStrs[labelKey.GetString()] = labelValue.GetString()
	}

	//This method will create the volume if it doesn't exist, or it will get it if it exists
	//From Docker docs: If you specify a volume name already in use on the current driver, Docker assumes you want to re-use the existing volume and does not return an error.
	//https://docs.docker.com/engine/reference/commandline/volume_create/
	if err := dockerManager.CreateVolume(ctx, volumeName, volumeLabelStrs); err != nil {
		return "", 0, nil, stacktrace.Propagate(
			err,
			"An error occurred creating logs database volume with name '%v' and labels '%+v'",
			volumeName,
			volumeLabelStrs,
		)
	}
	deleteVolumeFunc := func() {
		if err := dockerManager.RemoveVolume(ctx, volumeName); err != nil {
			logrus.Errorf(
				"Launching the logs database server for the engine with GUID '%v' didn't complete successfully so we "+
					"tried to remove the associated logs database volume '%v' we started, but doing so exited with an error:\n%v",
				engineGuid,
				volumeName,
				err)
			logrus.Errorf("ACTION REQUIRED: You'll need to manually remove the logs database volume '%v'!!!!!!", volumeName)
		}
	}
	shouldDeleteLogsDatabaseVolume := true
	defer func() {
		if shouldDeleteLogsDatabaseVolume {
			deleteVolumeFunc()
		}
	}()

	privateHttpPortSpec, err := logsDatabaseContainerConfigProvider.GetPrivateHttpPortSpec()
	if err != nil {
		return "", 0, nil, stacktrace.Propagate(err, "An error occurred getting the logs database container's private port spec")
	}

	logsDatabaseAttrs, err := objAttrsProvider.ForLogsDatabase(
		engineGuid,
		consts.LogsDatabaseHttpPortId,
		privateHttpPortSpec,
	)
	if err != nil {
		return "", 0, nil, stacktrace.Propagate(
			err,
			"An error occurred getting the logs database container attributes using GUID '%v' and the HTTP port spec '%+v'",
			engineGuid,
			privateHttpPortSpec,
		)
	}

	containerLabelStrs := map[string]string{}
	for labelKey, labelValue := range logsDatabaseAttrs.GetLabels() {
		containerLabelStrs[labelKey.GetString()] = labelValue.GetString()
	}

	containerName := logsDatabaseAttrs.GetName().GetString()

	createAndStartArgs, err := logsDatabaseContainerConfigProvider.GetContainerArgs(containerName, containerLabelStrs, volumeName, targetNetworkId)
	if err != nil {
		return "", 0, nil,
			stacktrace.Propagate(
				err,
				"An error occurred getting the logs database container args with container name '%v', labels '%+v', volume name '%v' and network ID '%v",
				containerName,
				containerLabelStrs,
				volumeName,
				targetNetworkId,
			)
	}

	containerId, _, err := dockerManager.CreateAndStartContainer(ctx, createAndStartArgs)
	if err != nil {
		return "", 0, nil, stacktrace.Propagate(err, "An error occurred starting the logs database container with these args '%+v'", createAndStartArgs)
	}
	killContainerFunc := func() {
		if err := dockerManager.KillContainer(context.Background(), containerId); err != nil {
			logrus.Errorf(
				"Launching the logs database server with GUID '%v' and container ID '%v' didn't complete successfully so we "+
					"tried to kill the container we started, but doing so exited with an error:\n%v",
				engineGuid,
				containerId,
				err)
			logrus.Errorf("ACTION REQUIRED: You'll need to manually stop the logs database server with GUID '%v' and Docker container ID '%v'!!!!!!", engineGuid, containerId)
		}
	}
	shouldKillLogsDbContainer := true
	defer func() {
		if shouldKillLogsDbContainer {
			killContainerFunc()
		}
	}()

	if err := shared_helpers.WaitForPortAvailabilityUsingNetstat(
		ctx,
		dockerManager,
		containerId,
		privateHttpPortSpec,
		maxWaitForEngineAvailabilityRetries,
		timeBetweenWaitForEngineAvailabilityRetries,
	); err != nil {
		return "", 0, nil, stacktrace.Propagate(err, "An error occurred waiting for the log database's HTTP port to become available")
	}

	logsDatabaseIP, err := dockerManager.GetContainerIP(ctx, consts.NameOfNetworkToStartEngineContainersIn, containerId)
	if err != nil {
		return "", 0, nil, stacktrace.Propagate(err, "An error occurred ")
	}

	killContainerAndDeleteVolumeFunc := func() {
		killContainerFunc()
		deleteVolumeFunc()
	}

	shouldDeleteLogsDatabaseVolume = false
	shouldKillLogsDbContainer = false
	return logsDatabaseIP, privateHttpPortSpec.GetNumber(), killContainerAndDeleteVolumeFunc, nil
}

func createLogsCollectorContainer(
	ctx context.Context,
	engineGuid engine.EngineGUID,
	targetNetworkId string,
	objAttrsProvider object_attributes_provider.DockerObjectAttributesProvider,
	dockerManager *docker_manager.DockerManager,
	logsCollectorContainerConfigProvider logs_components.LogsCollectorContainerConfigProvider,
) (func(), error) {

	//Create the volume first
	logsDbVolumeAttrs, err := objAttrsProvider.ForLogsCollectorVolume()
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred getting the logs collector volume attributes for engine with GUID %v", engineGuid)
	}
	volumeName := logsDbVolumeAttrs.GetName().GetString()
	volumeLabelStrs := map[string]string{}
	for labelKey, labelValue := range logsDbVolumeAttrs.GetLabels() {
		volumeLabelStrs[labelKey.GetString()] = labelValue.GetString()
	}

	//This method will create the volume if it doesn't exist, or it will get it if it exists
	//From Docker docs: If you specify a volume name already in use on the current driver, Docker assumes you want to re-use the existing volume and does not return an error.
	//https://docs.docker.com/engine/reference/commandline/volume_create/
	if err := dockerManager.CreateVolume(ctx, volumeName, volumeLabelStrs); err != nil {
		return nil, stacktrace.Propagate(
			err,
			"An error occurred creating logs collector volume with name '%v' and labels '%+v'",
			volumeName,
			volumeLabelStrs,
		)
	}
	deleteVolumeFunc := func() {
		if err := dockerManager.RemoveVolume(ctx, volumeName); err != nil {
			logrus.Errorf(
				"Launching the logs collector server for the engine with GUID '%v' didn't complete successfully so we "+
					"tried to remove the associated logs collector volume '%v' we started, but doing so exited with an error:\n%v",
				engineGuid,
				volumeName,
				err)
			logrus.Errorf("ACTION REQUIRED: You'll need to manually remove the logs collector volume '%v'!!!!!!", volumeName)
		}
	}
	shouldDeleteLogsCollectorVolume := true
	defer func() {
		if shouldDeleteLogsCollectorVolume {
			deleteVolumeFunc()
		}
	}()

	privateTcpPortSpec, err := logsCollectorContainerConfigProvider.GetPrivateTcpPortSpec()
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred getting the logs collector private TCP port spec")
	}

	privateHttpPortSpec, err := logsCollectorContainerConfigProvider.GetPrivateHttpPortSpec()
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred getting the logs collector private HTTP port spec")
	}

	logsCollectorAttrs, err := objAttrsProvider.ForLogsCollector(
		engineGuid,
		consts.LogsCollectorTcpPortId,
		privateTcpPortSpec,
		consts.LogsCollectorHttpPortId,
		privateHttpPortSpec,
	)
	if err != nil {
		return nil, stacktrace.Propagate(
			err,
			"An error occurred getting the logs collector container attributes using GUID '%v' with TCP port spec '%+v' and HTTP port spec '%+v'",
			engineGuid,
			privateTcpPortSpec,
			privateHttpPortSpec,
		)
	}

	containerName := logsCollectorAttrs.GetName().GetString()
	labelStrs := map[string]string{}
	for labelKey, labelValue := range logsCollectorAttrs.GetLabels() {
		labelStrs[labelKey.GetString()] = labelValue.GetString()
	}

	createAndStartArgs, err := logsCollectorContainerConfigProvider.GetContainerArgs(containerName, labelStrs, volumeName, targetNetworkId, dockerManager)
	if err != nil {
		return nil,
			stacktrace.Propagate(
				err,
				"An error occurred getting the logs-collector-container-args with container name '%v', labels '%+v', and network ID '%v",
				containerName,
				labelStrs,
				targetNetworkId,
			)
	}

	containerId, _, err := dockerManager.CreateAndStartContainer(ctx, createAndStartArgs)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred starting the logs controller container with these args '%+v'", createAndStartArgs)
	}
	killContainerFunc := func() {
		if err := dockerManager.KillContainer(context.Background(), containerId); err != nil {
			logrus.Errorf(
				"Launching the logs controller server for engine with GUID '%v' and container ID '%v' didn't complete successfully so we "+
					"tried to kill the container we started, but doing so exited with an error:\n%v",
				engineGuid,
				containerId,
				err)
			logrus.Errorf("ACTION REQUIRED: You'll need to manually stop the logs controller server for engine with GUID '%v' and Docker container ID '%v'!!!!!!", engineGuid, containerId)
		}
	}
	shouldKillLogsCollectorContainer := true
	defer func() {
		if shouldKillLogsCollectorContainer {
			killContainerFunc()
		}
	}()


	logsCollectorAvailabilityChecker := fluentbit.NewFluentbitAvailabilityChecker(privateHttpPortSpec.GetNumber())

	if err := logsCollectorAvailabilityChecker.WaitForAvailability(); err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred waiting for the log collector's to become available")
	}

	killContainerAndDeleteVolumeFunc := func() {
		killContainerFunc()
		deleteVolumeFunc()
	}

	shouldDeleteLogsCollectorVolume = false
	shouldKillLogsCollectorContainer = false
	return killContainerAndDeleteVolumeFunc, nil
}

func isThereAnyOtherEngineOrLogsComponentsContainersInTheCluster(
	ctx context.Context,
	dockerManager *docker_manager.DockerManager,
) (bool, []string, error){

	existentContainerIds := []string{}

	getAllEngineContainersOperation := func() (interface{}, error) {
		getAllEngineFilters := &engine.EngineFilters{}
		allEngineContainers, err := getMatchingEngines(ctx, getAllEngineFilters, dockerManager)
		if err != nil {
			return false, stacktrace.Propagate(err, "An error occurred getting all the engines using filters '%+v'", getAllEngineFilters)
		}

		allEngineContainerIDs := map[string]bool{}

		for containerId := range allEngineContainers{
			allEngineContainerIDs[containerId] = true
		}

		return allEngineContainers, nil
	}

	getAllLogsDatabaseContainersOperation := func() (interface{}, error) {
		allLogsDatabaseContainers, err := getAllLogsDatabaseContainers(ctx, dockerManager)
		if err != nil {
			return false, stacktrace.Propagate(err, "An error occurred getting all logs database containers")
		}

		allLogsDatabaseContainerIDs := map[string]bool{}

		for _, container := range allLogsDatabaseContainers{
			allLogsDatabaseContainerIDs[container.GetId()] = true
		}

		return allLogsDatabaseContainerIDs, nil
	}

	getAllLogsCollectorContainersOperation := func() (interface{}, error) {
		allLogsCollectorContainers, err := shared_helpers.GetAllLogsCollectorContainers(ctx, dockerManager)
		if err != nil {
			return false, stacktrace.Propagate(err, "An error occurred getting all logs collector containers")
		}

		allLogsCollectorContainerIDs := map[string]bool{}

		for _, container := range allLogsCollectorContainers{
			allLogsCollectorContainerIDs[container.GetId()] = true
		}

		return allLogsCollectorContainerIDs, nil
	}

	allOperations := map[operation_parallelizer.OperationID]operation_parallelizer.Operation{
		getAllEngineContainersOperationId: getAllEngineContainersOperation,
		getAllLogsDatabaseContainersOperationId: getAllLogsDatabaseContainersOperation,
		getAllLogsCollectorContainersOperationId: getAllLogsCollectorContainersOperation,
	}

	successfullOperations, erroredOperations := operation_parallelizer.RunOperationsInParallel(allOperations)
	if len(erroredOperations) > 0 {
		return false, nil, stacktrace.NewError("An error occurred running these operations '%+v' in parallel\n Operations with errors: %+v", allOperations, erroredOperations)
	}

	for _, uncastedContainerIds := range successfullOperations {
		containerIdsValue := reflect.ValueOf(uncastedContainerIds)
		for _, containerIdValue := range containerIdsValue.MapKeys() {
			existentContainerIds = append(existentContainerIds, containerIdValue.String())
		}
	}

	if len(existentContainerIds) > 0 {
		return true, existentContainerIds, nil
	}
	return false, existentContainerIds, nil
}