package add_service

import (
	"context"
	"fmt"
	"github.com/kurtosis-tech/kurtosis/container-engine-lib/lib/backend_interface/objects/service"
	"github.com/kurtosis-tech/kurtosis/core/server/api_container/server/service_network"
	"github.com/kurtosis-tech/kurtosis/core/server/api_container/server/startosis_engine/enclave_structure"
	"github.com/kurtosis-tech/kurtosis/core/server/api_container/server/startosis_engine/kurtosis_starlark_framework"
	"github.com/kurtosis-tech/kurtosis/core/server/api_container/server/startosis_engine/kurtosis_starlark_framework/builtin_argument"
	"github.com/kurtosis-tech/kurtosis/core/server/api_container/server/startosis_engine/kurtosis_starlark_framework/kurtosis_plan_instruction"
	"github.com/kurtosis-tech/kurtosis/core/server/api_container/server/startosis_engine/kurtosis_types/service_config"
	"github.com/kurtosis-tech/kurtosis/core/server/api_container/server/startosis_engine/runtime_value_store"
	"github.com/kurtosis-tech/kurtosis/core/server/api_container/server/startosis_engine/startosis_errors"
	"github.com/kurtosis-tech/kurtosis/core/server/api_container/server/startosis_engine/startosis_validator"
	"github.com/kurtosis-tech/stacktrace"
	"go.starlark.net/starlark"
	"reflect"
)

const (
	AddServiceBuiltinName = "add_service"

	ServiceNameArgName   = "name"
	ServiceConfigArgName = "config"
)

func NewAddService(serviceNetwork service_network.ServiceNetwork, runtimeValueStore *runtime_value_store.RuntimeValueStore) *kurtosis_plan_instruction.KurtosisPlanInstruction {
	return &kurtosis_plan_instruction.KurtosisPlanInstruction{
		KurtosisBaseBuiltin: &kurtosis_starlark_framework.KurtosisBaseBuiltin{
			Name: AddServiceBuiltinName,

			Arguments: []*builtin_argument.BuiltinArgument{
				{
					Name:              ServiceNameArgName,
					IsOptional:        false,
					ZeroValueProvider: builtin_argument.ZeroValueProvider[starlark.String],
					Validator: func(value starlark.Value) *startosis_errors.InterpretationError {
						return builtin_argument.NonEmptyString(value, ServiceNameArgName)
					},
				},
				{
					Name:              ServiceConfigArgName,
					IsOptional:        false,
					ZeroValueProvider: builtin_argument.ZeroValueProvider[*service_config.ServiceConfig],
					Validator: func(value starlark.Value) *startosis_errors.InterpretationError {
						// we just try to convert the configs here to validate their shape, to avoid code duplication
						// with Interpret
						if _, _, err := validateAndConvertConfigAndReadyCondition(serviceNetwork, value); err != nil {
							return err
						}
						return nil
					},
				},
			},
		},

		Capabilities: func() kurtosis_plan_instruction.KurtosisPlanInstructionCapabilities {
			return &AddServiceCapabilities{
				serviceNetwork:    serviceNetwork,
				runtimeValueStore: runtimeValueStore,

				serviceName:   "",  // populated at interpretation time
				serviceConfig: nil, // populated at interpretation time

				resultUuid:     "",  // populated at interpretation time
				readyCondition: nil, // populated at interpretation time
			}
		},

		DefaultDisplayArguments: map[string]bool{
			ServiceNameArgName:   true,
			ServiceConfigArgName: true,
		},
	}
}

type AddServiceCapabilities struct {
	serviceNetwork    service_network.ServiceNetwork
	runtimeValueStore *runtime_value_store.RuntimeValueStore

	serviceName    service.ServiceName
	serviceConfig  *service.ServiceConfig
	readyCondition *service_config.ReadyCondition

	resultUuid string
}

//TODO lo único que se me ocurre es que haya un método que retorne las storableCapabilities en un formato que se puedan guardar
//TODO y que luego haya otro metodo que tome ese formato y lo devuelva en el formato que queremos comparar

func (builtin *AddServiceCapabilities) Interpret(_ string, arguments *builtin_argument.ArgumentValuesSet) (starlark.Value, *startosis_errors.InterpretationError) {
	serviceName, err := builtin_argument.ExtractArgumentValue[starlark.String](arguments, ServiceNameArgName)
	if err != nil {
		return nil, startosis_errors.WrapWithInterpretationError(err, "Unable to extract value for '%s' argument", ServiceNameArgName)
	}

	serviceConfig, err := builtin_argument.ExtractArgumentValue[*service_config.ServiceConfig](arguments, ServiceConfigArgName)
	if err != nil {
		return nil, startosis_errors.WrapWithInterpretationError(err, "Unable to extract value for '%s' argument", ServiceConfigArgName)
	}
	apiServiceConfig, readyCondition, interpretationErr := validateAndConvertConfigAndReadyCondition(builtin.serviceNetwork, serviceConfig)
	if interpretationErr != nil {
		return nil, interpretationErr
	}

	builtin.serviceName = service.ServiceName(serviceName.GoString())
	builtin.serviceConfig = apiServiceConfig
	builtin.readyCondition = readyCondition
	builtin.resultUuid, err = builtin.runtimeValueStore.GetOrCreateValueAssociatedWithService(builtin.serviceName)
	if err != nil {
		return nil, startosis_errors.WrapWithInterpretationError(err, "Unable to create runtime value to hold '%v' command return values", AddServiceBuiltinName)
	}

	returnValue, interpretationErr := makeAddServiceInterpretationReturnValue(serviceName, builtin.serviceConfig, builtin.resultUuid)
	if interpretationErr != nil {
		return nil, interpretationErr
	}
	return returnValue, nil
}

func (builtin *AddServiceCapabilities) Validate(_ *builtin_argument.ArgumentValuesSet, validatorEnvironment *startosis_validator.ValidatorEnvironment) *startosis_errors.ValidationError {
	if validationErr := validateSingleService(validatorEnvironment, builtin.serviceName, builtin.serviceConfig); validationErr != nil {
		return validationErr
	}
	return nil
}

func (builtin *AddServiceCapabilities) Execute(ctx context.Context, _ *builtin_argument.ArgumentValuesSet) (string, error) {
	replacedServiceName, replacedServiceConfig, err := replaceMagicStrings(builtin.runtimeValueStore, builtin.serviceName, builtin.serviceConfig)
	if err != nil {
		return "", stacktrace.Propagate(err, "An error occurred replace a magic string in '%s' instruction arguments for service '%s'. Execution cannot proceed", AddServiceBuiltinName, builtin.serviceName)
	}
	var startedService *service.Service
	exist, err := builtin.serviceNetwork.ExistServiceRegistration(builtin.serviceName)
	if err != nil {
		return "", stacktrace.Propagate(err, "An error occurred getting service registration for service '%s'", builtin.serviceName)
	}
	if exist {
		startedService, err = builtin.serviceNetwork.UpdateService(ctx, replacedServiceName, replacedServiceConfig)
	} else {
		startedService, err = builtin.serviceNetwork.AddService(ctx, replacedServiceName, replacedServiceConfig)
	}
	if err != nil {
		return "", stacktrace.Propagate(err, "Unexpected error occurred starting service '%s'", replacedServiceName)
	}

	if err := runServiceReadinessCheck(
		ctx,
		builtin.serviceNetwork,
		builtin.runtimeValueStore,
		replacedServiceName,
		builtin.readyCondition,
	); err != nil {
		return "", stacktrace.Propagate(err, "An error occurred while checking if service '%v' is ready", replacedServiceName)
	}

	fillAddServiceReturnValueWithRuntimeValues(startedService, builtin.resultUuid, builtin.runtimeValueStore)
	instructionResult := fmt.Sprintf("Service '%s' added with service UUID '%s'", replacedServiceName, startedService.GetRegistration().GetUUID())
	return instructionResult, nil
}

func (builtin *AddServiceCapabilities) TryResolveWith(instructionsAreEqual bool, other kurtosis_plan_instruction.KurtosisPlanInstructionCapabilities, enclaveComponents *enclave_structure.EnclaveComponents) enclave_structure.InstructionResolutionStatus {
	// if other instruction is nil or other instruction is not an add_service instruction, status is unknown
	if other == nil {
		enclaveComponents.AddService(builtin.serviceName, enclave_structure.ComponentIsNew)
		return enclave_structure.InstructionIsUnknown
	}
	otherAddServiceCapabilities, ok := other.(*AddServiceCapabilities)
	if !ok {
		enclaveComponents.AddService(builtin.serviceName, enclave_structure.ComponentIsNew)
		return enclave_structure.InstructionIsUnknown
	}

	// if service names don't match, status is unknown, instructions can't be resolved together
	if otherAddServiceCapabilities.serviceName != builtin.serviceName {
		enclaveComponents.AddService(builtin.serviceName, enclave_structure.ComponentIsNew)
		return enclave_structure.InstructionIsUnknown
	}

	// if service names are equal but the instructions are not equal, it means the service config has been updated.
	// The instruction should be rerun
	if !instructionsAreEqual {
		enclaveComponents.AddService(builtin.serviceName, enclave_structure.ComponentIsUpdated)
		return enclave_structure.InstructionIsUpdate
	}

	// From here instructions are equal
	// We check if there has been some updates to the files it's mounting. If that's the case, it should be rerun
	filesArtifactsExpansion := builtin.serviceConfig.GetFilesArtifactsExpansion()
	if filesArtifactsExpansion != nil {
		for _, filesArtifactName := range filesArtifactsExpansion.ServiceDirpathsToArtifactIdentifiers {
			if enclaveComponents.HasFilesArtifactBeenUpdated(filesArtifactName) {
				enclaveComponents.AddService(builtin.serviceName, enclave_structure.ComponentIsUpdated)
				return enclave_structure.InstructionIsUpdate
			}
		}
	}

	enclaveComponents.AddService(builtin.serviceName, enclave_structure.ComponentWasLeftIntact)
	return enclave_structure.InstructionIsEqual
}

func validateAndConvertConfigAndReadyCondition(
	serviceNetwork service_network.ServiceNetwork,
	rawConfig starlark.Value,
) (*service.ServiceConfig, *service_config.ReadyCondition, *startosis_errors.InterpretationError) {
	config, ok := rawConfig.(*service_config.ServiceConfig)
	if !ok {
		return nil, nil, startosis_errors.NewInterpretationError("The '%s' argument is not a ServiceConfig (was '%s').", ConfigsArgName, reflect.TypeOf(rawConfig))
	}
	apiServiceConfig, interpretationErr := config.ToKurtosisType(serviceNetwork)
	if interpretationErr != nil {
		return nil, nil, interpretationErr
	}

	readyCondition, interpretationErr := config.GetReadyCondition()
	if interpretationErr != nil {
		return nil, nil, interpretationErr
	}

	return apiServiceConfig, readyCondition, nil
}
