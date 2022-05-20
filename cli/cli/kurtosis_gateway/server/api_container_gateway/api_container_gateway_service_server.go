package api_container_gateway

import (
	"context"
	"fmt"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_interface/objects/port_spec"
	"github.com/kurtosis-tech/kurtosis-cli/cli/kurtosis_gateway/connection"
	"github.com/kurtosis-tech/kurtosis-core-api-lib/api/golang/kurtosis_core_rpc_api_bindings"
	"github.com/kurtosis-tech/stacktrace"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
)

const (
	localHostIpStr = "127.0.0.1"
)

type ApiContainerGatewayServiceServer struct {
	// This embedding is required by gRPC
	kurtosis_core_rpc_api_bindings.UnimplementedApiContainerServiceServer

	// Id of enclave the API container is running in
	enclaveId string
	// Client for the api container we'll be connecting too
	remoteApiContainerClient kurtosis_core_rpc_api_bindings.ApiContainerServiceClient

	// Provides connections to Kurtosis objectis in cluster
	connectionProvider *connection.GatewayConnectionProvider

	// ServiceMap and mutex to protect it
	mutex                           *sync.Mutex
	userServiceToLocalConnectionMap map[string]*runningLocalServiceConnection
}

type runningLocalServiceConnection struct {
	localPublicServicePorts map[string]*kurtosis_core_rpc_api_bindings.Port
	localPublicIp           string
	kurtosisConnection      connection.GatewayConnectionToKurtosis
}

func NewEnclaveApiContainerGatewayServer(connectionProvider *connection.GatewayConnectionProvider, remoteApiContainerClient kurtosis_core_rpc_api_bindings.ApiContainerServiceClient, enclaveId string) (resultCoreGatewayServerService *ApiContainerGatewayServiceServer, resultGatewayCloseFunc func()) {
	// Start out with 0 connections to user services
	userServiceToLocalConnectionMap := map[string]*runningLocalServiceConnection{}

	closeGatewayFunc := func() {
		// Stop any port forwarding
		for _, runningLocalServiceConnection := range resultCoreGatewayServerService.userServiceToLocalConnectionMap {
			runningLocalServiceConnection.kurtosisConnection.Stop()
		}
	}

	return &ApiContainerGatewayServiceServer{
		remoteApiContainerClient:        remoteApiContainerClient,
		connectionProvider:              connectionProvider,
		mutex:                           &sync.Mutex{},
		userServiceToLocalConnectionMap: userServiceToLocalConnectionMap,
		enclaveId:                       enclaveId,
	}, closeGatewayFunc
}

func (service *ApiContainerGatewayServiceServer) LoadModule(ctx context.Context, args *kurtosis_core_rpc_api_bindings.LoadModuleArgs) (*kurtosis_core_rpc_api_bindings.LoadModuleResponse, error) {
	remoteApiContainerResponse, err := service.remoteApiContainerClient.LoadModule(ctx, args)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Expected to be able to call the remote api container from the gateway, instead a non nil err was returned")
	}

	return remoteApiContainerResponse, nil
}
func (service *ApiContainerGatewayServiceServer) UnloadModule(ctx context.Context, args *kurtosis_core_rpc_api_bindings.UnloadModuleArgs) (*emptypb.Empty, error) {
	remoteApiContainerResponse, err := service.remoteApiContainerClient.UnloadModule(ctx, args)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Expected to be able to call the remote api container from the gateway, instead a non nil err was returned")
	}

	return remoteApiContainerResponse, nil
}
func (service *ApiContainerGatewayServiceServer) ExecuteModule(ctx context.Context, args *kurtosis_core_rpc_api_bindings.ExecuteModuleArgs) (*kurtosis_core_rpc_api_bindings.ExecuteModuleResponse, error) {
	remoteApiContainerResponse, err := service.remoteApiContainerClient.ExecuteModule(ctx, args)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Expected to be able to call the remote api container from the gateway, instead a non nil err was returned")
	}

	return remoteApiContainerResponse, nil
}
func (service *ApiContainerGatewayServiceServer) GetModuleInfo(ctx context.Context, args *kurtosis_core_rpc_api_bindings.GetModuleInfoArgs) (*kurtosis_core_rpc_api_bindings.GetModuleInfoResponse, error) {
	remoteApiContainerResponse, err := service.remoteApiContainerClient.GetModuleInfo(ctx, args)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Expected to be able to call the remote api container from the gateway, instead a non nil err was returned")
	}

	return remoteApiContainerResponse, nil
}
func (service *ApiContainerGatewayServiceServer) RegisterService(ctx context.Context, args *kurtosis_core_rpc_api_bindings.RegisterServiceArgs) (*kurtosis_core_rpc_api_bindings.RegisterServiceResponse, error) {
	remoteApiContainerResponse, err := service.remoteApiContainerClient.RegisterService(ctx, args)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Expected to be able to call the remote api container from the gateway, instead a non nil err was returned")
	}

	return remoteApiContainerResponse, nil

}
func (service *ApiContainerGatewayServiceServer) StartService(ctx context.Context, args *kurtosis_core_rpc_api_bindings.StartServiceArgs) (*kurtosis_core_rpc_api_bindings.StartServiceResponse, error) {
	remoteApiContainerResponse, err := service.remoteApiContainerClient.StartService(ctx, args)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Expected to be able to call the remote api container from the gateway, instead a non nil err was returned")
	}
	cleanUpService := true
	defer func() {
		if cleanUpService {
			destroyEnclaveArgs := &kurtosis_core_rpc_api_bindings.RemoveServiceArgs{ServiceId: args.GetServiceId()}
			if _, err := service.remoteApiContainerClient.RemoveService(ctx, destroyEnclaveArgs); err != nil {
				logrus.Error("Connecting to the service running in the remote cluster failed, expected to be able to cleanup the created service, but an error occurred calling the backend to remove the service we created:")
				fmt.Fprintln(logrus.StandardLogger().Out, err)
				logrus.Errorf("ACTION REQUIRED: You'll need to manually remote the service with id '%v'", args.GetServiceId())
			}
		}
	}()

	runningLocalServiceConnection, err := service.startRunningConnectionForKurtosisService(args.GetServiceId(), args.PrivatePorts)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Expected to be able to start a local connection to service '%v', instead a non-nil error was returned", args.GetServiceId())
	}

	// Overwrite PublicPorts and PublicIp fields
	remoteApiContainerResponse.PublicIpAddr = runningLocalServiceConnection.localPublicIp
	remoteApiContainerResponse.PublicPorts = runningLocalServiceConnection.localPublicServicePorts

	cleanUpService = false
	return remoteApiContainerResponse, nil
}

func (service *ApiContainerGatewayServiceServer) GetServices(ctx context.Context, args *emptypb.Empty) (*kurtosis_core_rpc_api_bindings.GetServicesResponse, error) {
	remoteApiContainerResponse, err := service.remoteApiContainerClient.GetServices(ctx, args)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Expected to be able to call the remote api container from the gateway, instead a non nil err was returned")
	}

	return remoteApiContainerResponse, nil
}

func (service *ApiContainerGatewayServiceServer) GetServiceInfo(ctx context.Context, args *kurtosis_core_rpc_api_bindings.GetServiceInfoArgs) (*kurtosis_core_rpc_api_bindings.GetServiceInfoResponse, error) {
	remoteApiContainerResponse, err := service.remoteApiContainerClient.GetServiceInfo(ctx, args)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Expected to be able to call the remote api container from the gateway, instead a non nil err was returned")
	}

	// Get the running connection if it's available, start one if there is no running connection
	serviceId := args.GetServiceId()
	var runningLocalConnection *runningLocalServiceConnection
	runningLocalConnection, isFound := service.userServiceToLocalConnectionMap[serviceId]
	if !isFound {
		runningLocalConnection, err = service.startRunningConnectionForKurtosisService(serviceId, remoteApiContainerResponse.PrivatePorts)
		if err != nil {
			return nil, stacktrace.Propagate(err, "Expected to be able to start a local connection to kurtosis service '%v', instead a non-nil error was returned", serviceId)
		}
	}
	// Overwrite PublicPorts and PublicIp fields
	remoteApiContainerResponse.PublicPorts = runningLocalConnection.localPublicServicePorts
	remoteApiContainerResponse.PublicIpAddr = runningLocalConnection.localPublicIp

	return remoteApiContainerResponse, nil
}
func (service *ApiContainerGatewayServiceServer) RemoveService(ctx context.Context, args *kurtosis_core_rpc_api_bindings.RemoveServiceArgs) (*emptypb.Empty, error) {
	remoteApiContainerResponse, err := service.remoteApiContainerClient.RemoveService(ctx, args)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Expected to be able to call the remote api container from the gateway, instead a non nil err was returned")
	}
	// Kill the connection if we can
	service.idempotentKillRunningConnectionForServiceGuid(args.GetServiceId())

	return remoteApiContainerResponse, nil
}
func (service *ApiContainerGatewayServiceServer) Repartition(ctx context.Context, args *kurtosis_core_rpc_api_bindings.RepartitionArgs) (*emptypb.Empty, error) {
	remoteApiContainerResponse, err := service.remoteApiContainerClient.Repartition(ctx, args)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Expected to be able to call the remote api container from the gateway, instead a non nil err was returned")
	}

	return remoteApiContainerResponse, nil
}
func (service *ApiContainerGatewayServiceServer) ExecCommand(ctx context.Context, args *kurtosis_core_rpc_api_bindings.ExecCommandArgs) (*kurtosis_core_rpc_api_bindings.ExecCommandResponse, error) {
	remoteApiContainerResponse, err := service.remoteApiContainerClient.ExecCommand(ctx, args)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Expected to be able to call the remote api container from the gateway, instead a non nil err was returned")
	}

	return remoteApiContainerResponse, nil
}
func (service *ApiContainerGatewayServiceServer) WaitForHttpGetEndpointAvailability(ctx context.Context, args *kurtosis_core_rpc_api_bindings.WaitForHttpGetEndpointAvailabilityArgs) (*emptypb.Empty, error) {
	remoteApiContainerResponse, err := service.remoteApiContainerClient.WaitForHttpGetEndpointAvailability(ctx, args)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Expected to be able to call the remote api container from the gateway, instead a non nil err was returned")
	}

	return remoteApiContainerResponse, nil
}
func (service *ApiContainerGatewayServiceServer) WaitForHttpPostEndpointAvailability(ctx context.Context, args *kurtosis_core_rpc_api_bindings.WaitForHttpPostEndpointAvailabilityArgs) (*emptypb.Empty, error) {
	remoteApiContainerResponse, err := service.remoteApiContainerClient.WaitForHttpPostEndpointAvailability(ctx, args)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Expected to be able to call the remote api container from the gateway, instead a non nil err was returned")
	}

	return remoteApiContainerResponse, nil
}

func (service *ApiContainerGatewayServiceServer) GetModules(ctx context.Context, args *emptypb.Empty) (*kurtosis_core_rpc_api_bindings.GetModulesResponse, error) {
	remoteApiContainerResponse, err := service.remoteApiContainerClient.GetModules(ctx, args)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Expected to be able to call the remote api container from the gateway, instead a non nil err was returned")
	}

	return remoteApiContainerResponse, nil
}
func (service *ApiContainerGatewayServiceServer) UploadFilesArtifact(ctx context.Context, args *kurtosis_core_rpc_api_bindings.UploadFilesArtifactArgs) (*kurtosis_core_rpc_api_bindings.UploadFilesArtifactResponse, error) {
	remoteApiContainerResponse, err := service.remoteApiContainerClient.UploadFilesArtifact(ctx, args)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Expected to be able to call the remote api container from the gateway, instead a non nil err was returned")
	}

	return remoteApiContainerResponse, nil
}
func (service *ApiContainerGatewayServiceServer) StoreWebFilesArtifact(ctx context.Context, args *kurtosis_core_rpc_api_bindings.StoreWebFilesArtifactArgs) (*kurtosis_core_rpc_api_bindings.StoreWebFilesArtifactResponse, error) {
	remoteApiContainerResponse, err := service.remoteApiContainerClient.StoreWebFilesArtifact(ctx, args)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Expected to be able to call the remote api container from the gateway, instead a non nil err was returned")
	}

	return remoteApiContainerResponse, nil
}
func (service *ApiContainerGatewayServiceServer) StoreFilesArtifactFromService(ctx context.Context, args *kurtosis_core_rpc_api_bindings.StoreFilesArtifactFromServiceArgs) (*kurtosis_core_rpc_api_bindings.StoreFilesArtifactFromServiceResponse, error) {
	remoteApiContainerResponse, err := service.remoteApiContainerClient.StoreFilesArtifactFromService(ctx, args)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Expected to be able to call the remote api container from the gateway, instead a non nil err was returned")
	}

	return remoteApiContainerResponse, nil
}

func (service *ApiContainerGatewayServiceServer) PauseService(ctx context.Context, args *kurtosis_core_rpc_api_bindings.PauseServiceArgs) (*emptypb.Empty, error) {
	remoteApiContainerResponse, err := service.remoteApiContainerClient.PauseService(ctx, args)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Expected to be able to call the remote api container method from the gateway, instead a non nil err was returned")
	}

	return remoteApiContainerResponse, nil
}
func (service *ApiContainerGatewayServiceServer) UnpauseService(ctx context.Context, args *kurtosis_core_rpc_api_bindings.UnpauseServiceArgs) (*emptypb.Empty, error) {
	remoteApiContainerResponse, err := service.remoteApiContainerClient.UnpauseService(ctx, args)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Expected to be able to call the remote api container from the gateway, instead a non nil err was returned")
	}

	return remoteApiContainerResponse, nil
}

// Private functions for managing our running enclave api container gateways
func (service *ApiContainerGatewayServiceServer) startRunningConnectionForKurtosisService(serviceId string, servicePrivatePorts map[string]*kurtosis_core_rpc_api_bindings.Port) (*runningLocalServiceConnection, error) {
	service.mutex.Lock()
	defer service.mutex.Unlock()
	remotePrivatePortSpecs := map[string]*port_spec.PortSpec{}
	for portSpecId, coreApiPort := range servicePrivatePorts {
		if coreApiPort.GetProtocol() != kurtosis_core_rpc_api_bindings.Port_TCP {
			logrus.Warnf("Will not be able to forward service port with id '%v' in service '%v' for enclave '%v'. The protocol of this port is '%v', but only TCP protocol is support", portSpecId, serviceId, service.enclaveId, coreApiPort.GetProtocol())
			continue
		}
		portNumberUint16 := uint16(coreApiPort.GetNumber())
		remotePortSpec, err := port_spec.NewPortSpec(portNumberUint16, port_spec.PortProtocol_TCP)
		if err != nil {
			return nil, stacktrace.Propagate(err, "Expected to be able to create port spec describing remote port '%v', instead a non-nil error was returned", portSpecId)
		}
		remotePrivatePortSpecs[portSpecId] = remotePortSpec
	}

	// Start listening
	serviceConnection, err := service.connectionProvider.ForUserService(service.enclaveId, serviceId, remotePrivatePortSpecs)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Expected to be able to start a local connection service '%v' in enclave '%v', instead a non-nil error was returned", serviceId, service.enclaveId)
	}
	cleanUpConnection := true
	defer func() {
		if cleanUpConnection {
			serviceConnection.Stop()
		}
	}()
	// Locally forward ports described as expected by GRPC bindings
	localPublicApiPorts := map[string]*kurtosis_core_rpc_api_bindings.Port{}

	runingLocalServiceConnection := &runningLocalServiceConnection{
		kurtosisConnection:      serviceConnection,
		localPublicServicePorts: localPublicApiPorts,
		localPublicIp:           localHostIpStr,
	}

	// Store information about our running gateway
	service.userServiceToLocalConnectionMap[serviceId] = runingLocalServiceConnection
	cleanUpMapEntry := true
	defer func() {
		if cleanUpMapEntry {
			delete(service.userServiceToLocalConnectionMap, serviceId)
		}
	}()

	cleanUpMapEntry = false
	cleanUpConnection = false
	return runingLocalServiceConnection, nil
}

func (service *ApiContainerGatewayServiceServer) idempotentKillRunningConnectionForServiceGuid(serviceId string) {
	service.mutex.Lock()
	defer service.mutex.Unlock()
	runningLocalConnection, isRunning := service.userServiceToLocalConnectionMap[serviceId]
	// Nothing running, nothing to kill
	if !isRunning {
		return
	}

	// Close up the connection
	runningLocalConnection.kurtosisConnection.Stop()
	// delete the entry for the serve
	delete(service.userServiceToLocalConnectionMap, serviceId)
}
