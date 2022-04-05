package metrics_reporting

import (
	"context"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_interface"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_interface/objects/api_container"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_interface/objects/enclave"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_interface/objects/engine"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_interface/objects/module"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_interface/objects/networking_sidecar"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_interface/objects/port_spec"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_interface/objects/service"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_interface/objects/shell"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_interface/objects/wait_for_availability_http_methods"
	"github.com/kurtosis-tech/stacktrace"
	"io"
	"net"
)

// TODO CALL THE METRICS LIBRARY EVENT-REGISTRATION FUNCTIONS HERE!!!!
type MetricsReportingKurtosisBackend struct {
	underlying backend_interface.KurtosisBackend
}

func NewMetricsReportingKurtosisBackend(underlying backend_interface.KurtosisBackend) *MetricsReportingKurtosisBackend {
	return &MetricsReportingKurtosisBackend{underlying: underlying}
}

func (backend *MetricsReportingKurtosisBackend) PullImage(image string) error {
	if err := backend.underlying.PullImage(image); err != nil {
		return stacktrace.Propagate(err, "An error occurred pulling image '%v'", image)
	}
	return nil
}

func (backend *MetricsReportingKurtosisBackend) CreateEngine(ctx context.Context, imageOrgAndRepo string, imageVersionTag string, grpcPortNum uint16, grpcProxyPortNum uint16, engineDataDirpathOnHostMachine string, envVars map[string]string) (*engine.Engine, error) {
	result, err := backend.underlying.CreateEngine(ctx, imageOrgAndRepo, imageVersionTag, grpcPortNum, grpcProxyPortNum, engineDataDirpathOnHostMachine, envVars)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred creating the engine using image '%v' with tag '%v'", imageOrgAndRepo, imageVersionTag)
	}
	return result, nil
}

// Gets point-in-time data about engines matching the given filters
func (backend *MetricsReportingKurtosisBackend) GetEngines(ctx context.Context, filters *engine.EngineFilters) (map[string]*engine.Engine, error) {
	engines, err := backend.underlying.GetEngines(ctx, filters)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred getting engines using filters: %+v", filters)
	}
	return engines, nil
}

func (backend *MetricsReportingKurtosisBackend) StopEngines(ctx context.Context, filters *engine.EngineFilters) (
	successfulIds map[string]bool,
	failedIds map[string]error,
	resultErr error,
) {
	successes, failures, err := backend.underlying.StopEngines(ctx, filters)
	if err != nil {
		return nil, nil, stacktrace.Propagate(err, "An error occurred stopping engines using filters: %+v", filters)
	}
	return successes, failures, nil
}

func (backend *MetricsReportingKurtosisBackend) DestroyEngines(ctx context.Context, filters *engine.EngineFilters) (
	successfulIds map[string]bool,
	failedIds map[string]error,
	resultErr error,
) {
	successes, failures, err := backend.underlying.DestroyEngines(ctx, filters)
	if err != nil {
		return nil, nil, stacktrace.Propagate(err, "An error occurred destroying engines using filters: %+v", filters)
	}
	return successes, failures, nil
}

func (backend *MetricsReportingKurtosisBackend) CreateEnclave(
	ctx context.Context,
	enclaveId enclave.EnclaveID,
	isPartitioningEnabled bool,
)(*enclave.Enclave, error) {
	result, err := backend.underlying.CreateEnclave(ctx, enclaveId, isPartitioningEnabled)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred creating enclave with ID '%v' and is-partitioning-enabled value '%v'", enclaveId, isPartitioningEnabled)
	}
	return result, nil
}

func (backend *MetricsReportingKurtosisBackend) GetEnclaves(
	ctx context.Context,
	filters *enclave.EnclaveFilters,
)(
	map[enclave.EnclaveID]*enclave.Enclave,
	error,
){
	results, err := backend.underlying.GetEnclaves(ctx, filters)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred getting enclaves using filters: %+v", filters)
	}
	return results, nil
}

func (backend *MetricsReportingKurtosisBackend) StopEnclaves(
	ctx context.Context,
	filters *enclave.EnclaveFilters,
)(
	successfulEnclaveIds map[enclave.EnclaveID]bool,
	erroredEnclaveIds map[enclave.EnclaveID]error,
	resultErr error,
){
	successes, failures, err := backend.underlying.StopEnclaves(ctx, filters)
	if err != nil {
		return nil, nil, stacktrace.Propagate(err, "An error occurred stopping enclaves using filters: %+v", filters)
	}
	return successes, failures, nil
}

func (backend *MetricsReportingKurtosisBackend) DumpEnclave(
	ctx context.Context,
	enclaveId enclave.EnclaveID,
	outputDirpath string,
) error {
	if err := backend.underlying.DumpEnclave(ctx, enclaveId, outputDirpath); err != nil {
		return stacktrace.Propagate(err, "An error occurred dumping enclave '%v' to path '%v'", enclaveId, outputDirpath)
	}
	return nil
}

func (backend *MetricsReportingKurtosisBackend) DestroyEnclaves(
	ctx context.Context,
	filters *enclave.EnclaveFilters,
)(
	successfulEnclaveIds map[enclave.EnclaveID]bool,
	erroredEnclaveIds map[enclave.EnclaveID]error,
	resultErr error,
) {
	successes, failures, err := backend.underlying.DestroyEnclaves(ctx, filters)
	if err != nil {
		return nil, nil, stacktrace.Propagate(err, "An error occurred destroying enclaves using filters: %+v", filters)
	}
	return successes, failures, nil
}

func (backend *MetricsReportingKurtosisBackend) CreateAPIContainer(
	ctx context.Context,
	image string,
	enclaveId enclave.EnclaveID,
	ipAddr net.IP,
	grpcPortNum uint16,
	grpcProxyPortNum uint16,
	enclaveDataDirpathOnHostMachine string,
	envVars map[string]string,
) (*api_container.APIContainer, error) {
	result, err := backend.underlying.CreateAPIContainer(
		ctx,
		image,
		enclaveId,
		ipAddr,
		grpcPortNum,
		grpcProxyPortNum,
		enclaveDataDirpathOnHostMachine,
		envVars,
	)
	if err != nil {
		return nil, stacktrace.Propagate(
			err,
			"An error occurred creating an API container from image '%v' with envvars: %+v",
			image,
			envVars,
		)
	}
	return result, nil
}

func (backend *MetricsReportingKurtosisBackend) GetAPIContainers(ctx context.Context, filters *api_container.APIContainerFilters) (map[enclave.EnclaveID]*api_container.APIContainer, error) {
	results, err := backend.underlying.GetAPIContainers(ctx, filters)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred getting API containers matching filters: %+v", filters)
	}
	return results, nil
}

func (backend *MetricsReportingKurtosisBackend) StopAPIContainers(ctx context.Context, filters *api_container.APIContainerFilters) (successfulApiContainerIds map[enclave.EnclaveID]bool, erroredApiContainerIds map[enclave.EnclaveID]error, resultErr error) {
	successes, failures, err := backend.underlying.StopAPIContainers(ctx, filters)
	if err != nil {
		return nil, nil, stacktrace.Propagate(err, "An error occurred stopping API containers using filters: %+v", filters)
	}
	return successes, failures, nil
}

func (backend *MetricsReportingKurtosisBackend) DestroyAPIContainers(ctx context.Context, filters *api_container.APIContainerFilters) (successfulApiContainerIds map[enclave.EnclaveID]bool, erroredApiContainerIds map[enclave.EnclaveID]error, resultErr error) {
	successes, failures, err := backend.underlying.DestroyAPIContainers(ctx, filters)
	if err != nil {
		return nil, nil, stacktrace.Propagate(err, "An error occurred destroying API containers using filters: %+v", filters)
	}
	return successes, failures, nil
}


func (backend *MetricsReportingKurtosisBackend) CreateModule(
	ctx context.Context,
	image string,
	enclaveId enclave.EnclaveID,
	id module.ModuleID,
	guid module.ModuleGUID,
	ipAddr net.IP, // TODO REMOVE THIS ONCE WE FIX THE STATIC IP PROBLEM!!
	grpcPortNum uint16,
	enclaveDataDirpathOnHostMachine string,
	envVars map[string]string,
)(
	newModule *module.Module,
	resultErr error,
) {
	newModule, err := backend.underlying.CreateModule(
		ctx,
		image,
		enclaveId,
		id,
		guid,
		ipAddr,
		grpcPortNum,
		enclaveDataDirpathOnHostMachine,
		envVars,
	)
	if err != nil {
		return nil,
		stacktrace.Propagate(
			err,
			"An error occurred creating module with ID '%v', GUID '%v', and image '%v'",
			id,
			guid,
			image,
		)
	}

	return newModule, nil
}

func (backend *MetricsReportingKurtosisBackend) GetModules(
	ctx context.Context,
	filters *module.ModuleFilters,
)(
	map[module.ModuleGUID]*module.Module,
	error,
) {
	modules, err := backend.underlying.GetModules(ctx, filters)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred getting modules using filters: %+v", filters)
	}
	return modules, nil
}

func (backend *MetricsReportingKurtosisBackend) DestroyModules(
	ctx context.Context,
	filters *module.ModuleFilters,
)(
	successfulModuleIds map[module.ModuleGUID]bool,
	erroredModuleIds map[module.ModuleGUID]error,
	resultErr error,
) {
	successes, failures, err := backend.underlying.DestroyModules(ctx, filters)
	if err != nil {
		return nil, nil, stacktrace.Propagate(err, "An error occurred destroying modules using filters: %+v", filters)
	}
	return successes, failures, nil
}

func (backend *MetricsReportingKurtosisBackend) CreateUserService(
	ctx context.Context,
	id service.ServiceID,
	guid service.ServiceGUID,
	containerImageName string,
	enclaveId enclave.EnclaveID,
	ipAddr net.IP,
	privatePorts map[string]*port_spec.PortSpec,
	entrypointArgs []string,
	cmdArgs []string,
	envVars map[string]string,
	enclaveDataDirpathOnHostMachine string,
	filesArtifactMountDirpaths map[string]string,
)(
	newUserService *service.Service,
	resultErr error,
) {
	userService, err := backend.underlying.CreateUserService(
		ctx,
		id,
		guid,
		containerImageName,
		enclaveId,
		ipAddr,
		privatePorts,
		entrypointArgs,
		cmdArgs,
		envVars,
		enclaveDataDirpathOnHostMachine,
		filesArtifactMountDirpaths,
		)
	if err != nil {
		return nil,
		stacktrace.Propagate(
			err,
			"An error occurred creating the user service with ID '%v' and GUID '%v' using image '%v' with private ports '%+v' with entry point args '%+v', command args '%+v', environment vars '%+v', enclave data mount dirpath '%v' and file artifacts mount dirpath '%v'",
			id,
			guid,
			containerImageName,
			privatePorts,
			entrypointArgs,
			cmdArgs,
			envVars,
			enclaveDataDirpathOnHostMachine,
			filesArtifactMountDirpaths,
			)
	}
	return userService, nil
}

func (backend *MetricsReportingKurtosisBackend) GetUserServices(
	ctx context.Context,
	filters *service.ServiceFilters,
)(
	map[service.ServiceGUID]*service.Service,
	error,
){
	services,  err := backend.underlying.GetUserServices(ctx, filters)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred getting user services using filters '%+v'", filters)
	}
	return services, nil
}

func (backend *MetricsReportingKurtosisBackend) GetUserServiceLogs(
	ctx context.Context,
	filters *service.ServiceFilters,
	shouldFollowLogs bool,
)(
	map[service.ServiceGUID]io.ReadCloser,
	map[service.ServiceGUID]error,
	error,
) {
	userServiceLogs, erroredUserServices, err := backend.underlying.GetUserServiceLogs(ctx, filters, shouldFollowLogs)
	if err != nil {
		return nil, nil,  stacktrace.Propagate(err, "An error occurred getting user service logs using filters '%+v'", filters)
	}
	return userServiceLogs, erroredUserServices, nil
}

func (backend *MetricsReportingKurtosisBackend) RunUserServiceExecCommands (
	ctx context.Context,
	enclaveId enclave.EnclaveID,
	userServiceCommands map[service.ServiceGUID][]string,
)(
	succesfulUserServiceGuids map[service.ServiceGUID]bool,
	erroredUserServiceGuids map[service.ServiceGUID]error,
	resultErr error,
) {
	exitCode, output, err := backend.underlying.RunUserServiceExecCommands(ctx, enclaveId, userServiceCommands)
	if err != nil {
		return nil, nil, stacktrace.Propagate(
			err,
			"An error occurred running user service exec commands '%+v' on enclave '%v'",
			userServiceCommands,
			enclaveId,
		)
	}
	return exitCode, output, nil
}

func (backend *MetricsReportingKurtosisBackend) WaitForUserServiceHttpEndpointAvailability(
	ctx context.Context,
	enclaveId enclave.EnclaveID,
	serviceGUID service.ServiceGUID,
	httpMethod wait_for_availability_http_methods.WaitForAvailabilityHttpMethod,
	port uint32,
	path string,
	requestBody string,
	expectedResponseBody string,
	initialDelayMilliseconds uint32,
	retries uint32,
	retriesDelayMilliseconds uint32,
)(
	resultErr error,
) {
	if err := backend.underlying.WaitForUserServiceHttpEndpointAvailability(
		ctx,
		enclaveId,
		serviceGUID,
		httpMethod,
		port,
		path,
		requestBody,
		expectedResponseBody,
		initialDelayMilliseconds,
		retries,
		retriesDelayMilliseconds,
		); err != nil {
		return stacktrace.Propagate(
			err,
			"An error occurred waiting for http endpoint with path '%v', port '%v', request body '%v', expected response body '%v' from service with GUID '%v' in enclave with ID '%v' to become available after '%v' retries and '%v' milliseconds between retries,",
			path,
			port,
			requestBody,
			expectedResponseBody,
			serviceGUID,
			enclaveId,
			retries,
			retriesDelayMilliseconds,
		)
	}
	return nil
}

func (backend *MetricsReportingKurtosisBackend) GetShellOnUserService(
	ctx context.Context,
	enclaveId enclave.EnclaveID,
	serviceGUID service.ServiceGUID,
)(
	resultShell *shell.Shell,
	resultErr error,
) {
	newShell, err := backend.underlying.GetShellOnUserService(ctx, enclaveId, serviceGUID)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred getting shell on user service with GUID '%v'", serviceGUID)
	}
	return newShell, nil
}

func (backend *MetricsReportingKurtosisBackend) StopUserServices(
	ctx context.Context,
	filters *service.ServiceFilters,
)(
	successfulUserServiceGuids map[service.ServiceGUID]bool,
	erroredUserServiceGuids map[service.ServiceGUID]error,
	resultErr error,
) {
	successes, failures, err := backend.underlying.StopUserServices(ctx, filters)
	if err != nil {
		return nil, nil, stacktrace.Propagate(err, "An error occurred stopping user services using filters: %+v", filters)
	}
	return successes, failures, nil
}

func (backend *MetricsReportingKurtosisBackend) DestroyUserServices(
	ctx context.Context,
	filters *service.ServiceFilters,
)(
	successfulUserServiceGuids map[service.ServiceGUID]bool,
	erroredUserServiceGuids map[service.ServiceGUID]error,
	resultErr error,
) {
	successes, failures, err := backend.underlying.DestroyUserServices(ctx, filters)
	if err != nil {
		return nil, nil, stacktrace.Propagate(err, "An error occurred destroying user services using filters: %+v", filters)
	}
	return successes, failures, nil
}

func (backend *MetricsReportingKurtosisBackend) CreateNetworkingSidecar(
	ctx context.Context,
	enclaveId enclave.EnclaveID,
	serviceGuid service.ServiceGUID,
	ipAddr net.IP, // TODO REMOVE THIS ONCE WE FIX THE STATIC IP PROBLEM!!
)(
	*networking_sidecar.NetworkingSidecar,
	error,
){
	networkingSidecar, err := backend.underlying.CreateNetworkingSidecar(ctx, enclaveId, serviceGuid, ipAddr)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred creating networking sidecar for user service with GUID '%v' in enclave with ID '%v'", serviceGuid, enclaveId)
	}
	return networkingSidecar, nil
}

func (backend *MetricsReportingKurtosisBackend) GetNetworkingSidecars(
	ctx context.Context,
	filters *networking_sidecar.NetworkingSidecarFilters,
)(
	map[service.ServiceGUID]*networking_sidecar.NetworkingSidecar,
	error,
) {
	networkingSidecars, err := backend.underlying.GetNetworkingSidecars(ctx, filters)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred getting networking sidecars using filters '%+v'", filters)
	}
	return networkingSidecars, nil
}

func (backend *MetricsReportingKurtosisBackend) RunNetworkingSidecarExecCommands(
	ctx context.Context,
	enclaveId enclave.EnclaveID,
	networkingSidecarsCommands map[service.ServiceGUID][]string,
)(
	map[service.ServiceGUID]bool,
	map[service.ServiceGUID]error,
	error,
){
	successfulUserServiceGuids, erroredUserServiceGuids, err := backend.underlying.RunNetworkingSidecarExecCommands(ctx, enclaveId, networkingSidecarsCommands)
	if err != nil {
		return nil, nil, stacktrace.Propagate(err, "An error occurred running networking sidecar exec commands '%+v' in enclave with ID '%v'", networkingSidecarsCommands, enclaveId)
	}
	return successfulUserServiceGuids, erroredUserServiceGuids, nil
}

func (backend *MetricsReportingKurtosisBackend) StopNetworkingSidecars(
	ctx context.Context,
	filters *networking_sidecar.NetworkingSidecarFilters,
)(
	map[service.ServiceGUID]bool,
	map[service.ServiceGUID]error,
	error,
){
	successfulUserServiceGuids, erroredUserServiceGuids, err := backend.underlying.StopNetworkingSidecars(ctx, filters)
	if err != nil {
		return nil, nil, stacktrace.Propagate(err, "An error occurred stopping networking sidecars using filters '%+v'", filters)
	}
	return successfulUserServiceGuids, erroredUserServiceGuids, nil
}

func (backend *MetricsReportingKurtosisBackend) DestroyNetworkingSidecars(
	ctx context.Context,
	filters *networking_sidecar.NetworkingSidecarFilters,
)(
	map[service.ServiceGUID]bool,
	map[service.ServiceGUID]error,
	error,
){
	successfulUserServiceGuids, erroredUserServiceGuids, err := backend.underlying.DestroyNetworkingSidecars(ctx, filters)
	if err != nil {
		return nil, nil, stacktrace.Propagate(err, "An error occurred destroying networking sidecars using filters '%+v'", filters)
	}
	return successfulUserServiceGuids, erroredUserServiceGuids, nil
}
