package port_forward_manager

import (
	"context"
	"github.com/kurtosis-tech/kurtosis/cli/cli/kurtosis_gateway/port_utils"
	"github.com/kurtosis-tech/stacktrace"
	"github.com/sirupsen/logrus"
	"strconv"
)

const (
	localhostIpString = "127.0.0.1"
)

type PortForwardManager struct {
	serviceEnumerator *ServiceEnumerator
}

func NewPortForwardManager(serviceEnumerator *ServiceEnumerator) *PortForwardManager {
	return &PortForwardManager{
		serviceEnumerator: serviceEnumerator,
	}
}

func (manager *PortForwardManager) Ping(ctx context.Context) error {
	return manager.serviceEnumerator.checkHealth(ctx)
}

// CreateUserServicePortForward
// This can run in two manners:
// 1. requestedLocalPort is specified: this will target only one (enclaveId, serviceId, portId), so all must be specified
// 2. requestedLocalPort is unspecified (0): we will bind all services to ephemeral local ports.  The list of services depends
// upon what's specified:
//   - (enclaveId): finds all services and ports within the enclave and binds them
//   - (enclaveId, serviceId): finds all ports in the given service and binds them
//   - (enclaveId, serviceId, portId): finds a specific service/port and binds that (similar to case 1 but ephemeral)
func (manager *PortForwardManager) CreateUserServicePortForward(ctx context.Context, enclaveServicePort EnclaveServicePort, requestedLocalPort uint16) (map[EnclaveServicePort]uint16, error) {
	if err := validateCreateUserServicePortForwardArgs(enclaveServicePort, requestedLocalPort); err != nil {
		return nil, stacktrace.Propagate(err, "Validation failed for arguments")
	}

	if requestedLocalPort == 0 {
		ephemeralLocalPortSpec, err := port_utils.GetFreeTcpPort(localhostIpString)
		if err != nil {
			return nil, stacktrace.Propagate(err, "Could not allocate a local port for the tunnel")
		}

		requestedLocalPort = ephemeralLocalPortSpec.GetNumber()
	}

	portForward, err := manager.createAndOpenPortForwardToUserService(ctx, enclaveServicePort, requestedLocalPort)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Failed to set up port forward to (enclave, service, port), %v", enclaveServicePort)
	}

	return map[EnclaveServicePort]uint16{enclaveServicePort: portForward.localPortNumber}, nil
}

func (manager *PortForwardManager) RemoveUserServicePortForward(ctx context.Context, enclaveServicePort EnclaveServicePort) error {
	panic("implement me")
}

func (manager *PortForwardManager) createAndOpenPortForwardToUserService(ctx context.Context, enclaveServicePort EnclaveServicePort, localPortToBind uint16) (*PortForwardTunnel, error) {
	serviceInterfaceDetail, err := manager.serviceEnumerator.collectServiceInformation(ctx, enclaveServicePort)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Failed to enumerate service information for (enclave, service, port), %v", enclaveServicePort)
	}

	portForward := NewPortForwardTunnel(localPortToBind, serviceInterfaceDetail)
	logrus.Infof("Opening port forward session on local port %d, to remote service %v", portForward.localPortNumber, serviceInterfaceDetail)
	err = portForward.RunAsync()
	if err != nil {
		return nil, stacktrace.Propagate(err, "Failed to open a port forward tunnel to remote service %v", serviceInterfaceDetail)
	}

	return portForward, nil
}

// Check for two modes of operation:
// 1. where a local port is requested, we need all of (enclaveId, serviceId, portId) to be specified; this has to target one service
// 2. if no local port is requested, we need at least enclaveId, and will target as many services as possible within the given context
func validateCreateUserServicePortForwardArgs(enclaveServicePort EnclaveServicePort, requestedLocalPort uint16) error {
	if enclaveServicePort.EnclaveId() == "" {
		return stacktrace.NewError("EnclaveId is always required but we received an empty string")
	}

	if requestedLocalPort != 0 {
		if enclaveServicePort.ServiceId() == "" || enclaveServicePort.PortId() == "" {
			return stacktrace.NewError("A static port '%d' was requested, but enclaveId, serviceId, and portId were not all specified: %v", requestedLocalPort, enclaveServicePort)
		}
	}

	return nil
}

func getLocalChiselServerUri(localPortToChiselServer uint16) string {
	return "localhost:" + strconv.Itoa(int(localPortToChiselServer))
}
