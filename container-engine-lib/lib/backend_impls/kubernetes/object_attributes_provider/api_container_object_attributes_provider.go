package object_attributes_provider

import (
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_impls/kubernetes/object_attributes_provider/annotation_key_consts"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_impls/kubernetes/object_attributes_provider/kubernetes_annotation_key"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_impls/kubernetes/object_attributes_provider/kubernetes_annotation_value"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_impls/kubernetes/object_attributes_provider/kubernetes_label_key"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_impls/kubernetes/object_attributes_provider/kubernetes_label_value"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_impls/kubernetes/object_attributes_provider/kubernetes_object_name"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_impls/kubernetes/object_attributes_provider/label_key_consts"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_impls/kubernetes/object_attributes_provider/label_value_consts"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_impls/kubernetes/object_attributes_provider/port_spec_serializer"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_interface/objects/enclave"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_interface/objects/port_spec"
	"github.com/kurtosis-tech/stacktrace"
	"strings"
)

const (
	apiContainerNamePrefix                = "kurtosis-api"
)

type KubernetesApiContainerObjectAttributesProvider interface {
	ForApiContainerService(privateGrpcPortId string, privateGrpcPortSpec *port_spec.PortSpec) (KubernetesObjectAttributes, error)
	ForApiContainerNamespace() (KubernetesObjectAttributes, error)
}

// Private so it can't be instantiated
type kubernetesApiContainerObjectAttributesProviderImpl struct {
	enclaveId string
}

func GetKubernetesApiContainerObjectAttributesProvider(enclaveId enclave.EnclaveID) KubernetesApiContainerObjectAttributesProvider {
	return newKubernetesApiContainerObjectAttributesProviderImpl(enclaveId)
}

func newKubernetesApiContainerObjectAttributesProviderImpl(enclaveId enclave.EnclaveID) *kubernetesApiContainerObjectAttributesProviderImpl {
	return &kubernetesApiContainerObjectAttributesProviderImpl{
		enclaveId: string(enclaveId),
	}
}

func (provider *kubernetesApiContainerObjectAttributesProviderImpl) ForApiContainerService(grpcPortId string, grpcPortSpec *port_spec.PortSpec) (KubernetesObjectAttributes, error) {
	nameStr := provider.getApiContainerObjectNameString(serviceNameSuffix, []string{})
	name, err := kubernetes_object_name.CreateNewKubernetesObjectName(nameStr)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred creating a name for api container service")
	}

	enclaveIdLabelValue, err := kubernetes_label_value.CreateNewKubernetesLabelValue(provider.enclaveId)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred creating the enclave ID Kubernetes label from string '%v'", provider.enclaveId)
	}

	labels := map[*kubernetes_label_key.KubernetesLabelKey]*kubernetes_label_value.KubernetesLabelValue{
		label_key_consts.KurtosisResourceTypeLabelKey: label_value_consts.APIContainerContainerTypeLabelValue,
		label_key_consts.EnclaveIDLabelKey:            enclaveIdLabelValue,
	}

	usedPorts := map[string]*port_spec.PortSpec{
		grpcPortId:      grpcPortSpec,
	}
	serializedPortsSpec, err := port_spec_serializer.SerializePortSpecs(usedPorts)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred serializing the following api container server ports to a string for storing in the ports annotation: %+v", usedPorts)
	}

	// Store Kurtosis port_spec info in annotation
	annotations := map[*kubernetes_annotation_key.KubernetesAnnotationKey]*kubernetes_annotation_value.KubernetesAnnotationValue{
		annotation_key_consts.PortSpecsAnnotationKey: serializedPortsSpec,
	}

	objectAttributes, err := newKubernetesObjectAttributesImpl(name, labels, annotations)
	if err != nil {
		stacktrace.Propagate(err, "An error occurred while creating the Kubernetes object attributes with the name '%s' and labels '%+v', and annotations '%+v'", name, labels, annotations)
	}

	return objectAttributes, nil
}

func (provider *kubernetesApiContainerObjectAttributesProviderImpl) ForApiContainerNamespace() (KubernetesObjectAttributes, error) {
	nameStr := provider.getApiContainerObjectNameString(namespaceSuffix, []string{})
	name, err := kubernetes_object_name.CreateNewKubernetesObjectName(nameStr)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred creating a name for api container namespace")
	}

	enclaveIdLabelValue, err := kubernetes_label_value.CreateNewKubernetesLabelValue(provider.enclaveId)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred creating the enclave ID Kubernetes label from string '%v'", provider.enclaveId)
	}

	labels := map[*kubernetes_label_key.KubernetesLabelKey]*kubernetes_label_value.KubernetesLabelValue{
		label_key_consts.KurtosisResourceTypeLabelKey: label_value_consts.APIContainerContainerTypeLabelValue,
		label_key_consts.EnclaveIDLabelKey:            enclaveIdLabelValue,
	}

	// No custom annotations for api container namespace
	annotations := map[*kubernetes_annotation_key.KubernetesAnnotationKey]*kubernetes_annotation_value.KubernetesAnnotationValue{}

	objectAttributes, err := newKubernetesObjectAttributesImpl(name, labels, annotations)
	if err != nil {
		stacktrace.Propagate(err, "An error occurred while creating the Kubernetes object attributes with the name '%s' and labels '%+v', and annotations '%+v'", name, labels, annotations)
	}

	return objectAttributes, nil
}

func (provider *kubernetesApiContainerObjectAttributesProviderImpl) getApiContainerObjectNameString(suffix string, elems []string) string {
	toJoin := []string{
		provider.enclaveId,
		apiContainerNamePrefix,
	}
	if elems != nil {
		toJoin = append(toJoin, elems...)
	}
	toJoin = append(toJoin, suffix)
	nameStr := strings.Join(
		toJoin,
		objectNameElementSeparator,
	)
	return nameStr
}
