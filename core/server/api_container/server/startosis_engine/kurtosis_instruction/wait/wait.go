package wait

import (
	"context"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/kurtosis_core_rpc_api_bindings"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/binding_constructors"
	"github.com/kurtosis-tech/kurtosis/core/server/api_container/server/service_network"
	"github.com/kurtosis-tech/kurtosis/core/server/api_container/server/startosis_engine/kurtosis_instruction"
	"github.com/kurtosis-tech/kurtosis/core/server/api_container/server/startosis_engine/kurtosis_instruction/assert"
	"github.com/kurtosis-tech/kurtosis/core/server/api_container/server/startosis_engine/kurtosis_instruction/shared_helpers"
	"github.com/kurtosis-tech/kurtosis/core/server/api_container/server/startosis_engine/recipe"
	"github.com/kurtosis-tech/kurtosis/core/server/api_container/server/startosis_engine/runtime_value_store"
	"github.com/kurtosis-tech/kurtosis/core/server/api_container/server/startosis_engine/startosis_errors"
	"github.com/kurtosis-tech/kurtosis/core/server/api_container/server/startosis_engine/startosis_validator"
	"github.com/kurtosis-tech/stacktrace"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"time"
)

const (
	WaitBuiltinName = "wait"

	recipeArgName           = "recipe"
	valueFieldArgName       = "field"
	assertionArgName        = "assertion"
	targetArgName           = "target_value"
	optionalIntervalArgName = "interval?"
	optionalTimeoutArgName  = "timeout?"
	emptyOptionalField      = ""

	defaultInterval = 1 * time.Second
	defaultTimeout  = 15 * time.Minute
)

func GenerateWaitBuiltin(instructionsQueue *[]kurtosis_instruction.KurtosisInstruction, recipeExecutor *runtime_value_store.RuntimeValueStore, serviceNetwork service_network.ServiceNetwork) func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	// TODO: Force returning an InterpretationError rather than a normal error
	return func(thread *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		instructionPosition := shared_helpers.GetCallerPositionFromThread(thread)
		waitInstruction := newEmptyWaitInstructionInstruction(serviceNetwork, instructionPosition, recipeExecutor)
		if interpretationError := waitInstruction.parseStartosisArgs(builtin, args, kwargs); interpretationError != nil {
			return nil, interpretationError
		}
		resultUuid, err := recipeExecutor.CreateValue()
		if err != nil {
			return nil, startosis_errors.NewInterpretationError("An error occurred while generating uuid for future reference for %v instruction", WaitBuiltinName)
		}
		waitInstruction.resultUuid = resultUuid
		returnValue, interpretationErr := waitInstruction.recipe.CreateStarlarkReturnValue(waitInstruction.resultUuid)
		if interpretationErr != nil {
			return nil, startosis_errors.NewInterpretationError("An error occurred while creating return value for %v instruction", WaitBuiltinName)
		}
		*instructionsQueue = append(*instructionsQueue, waitInstruction)
		return returnValue, nil
	}
}

type WaitInstruction struct {
	serviceNetwork service_network.ServiceNetwork

	position       *kurtosis_instruction.InstructionPosition
	starlarkKwargs starlark.StringDict

	runtimeValueStore *runtime_value_store.RuntimeValueStore
	recipe            recipe.Recipe
	resultUuid        string
	targetKey         string
	assertion         string
	target            starlark.Comparable
	backoff           backoff.BackOff
	timeout           time.Duration
}

func newEmptyWaitInstructionInstruction(serviceNetwork service_network.ServiceNetwork, position *kurtosis_instruction.InstructionPosition, runtimeValueStore *runtime_value_store.RuntimeValueStore) *WaitInstruction {
	return &WaitInstruction{
		serviceNetwork:    serviceNetwork,
		position:          position,
		runtimeValueStore: runtimeValueStore,
		recipe:            nil,
		resultUuid:        "",
		starlarkKwargs:    nil,
		targetKey:         "",
		assertion:         "",
		target:            nil,
		backoff:           nil,
		timeout:           defaultTimeout,
	}
}

func newWaitInstructionInstructionForTest(serviceNetwork service_network.ServiceNetwork, position *kurtosis_instruction.InstructionPosition, runtimeValueStore *runtime_value_store.RuntimeValueStore, recipe recipe.Recipe, resultUuid string, targetKey string, assertion string, target starlark.Comparable, starlarkKwargs starlark.StringDict) *WaitInstruction {
	return &WaitInstruction{
		serviceNetwork:    serviceNetwork,
		position:          position,
		runtimeValueStore: runtimeValueStore,
		recipe:            recipe,
		resultUuid:        resultUuid,
		starlarkKwargs:    starlarkKwargs,
		targetKey:         targetKey,
		assertion:         assertion,
		target:            target,
		backoff:           nil,
		timeout:           defaultTimeout,
	}
}

func (instruction *WaitInstruction) GetPositionInOriginalScript() *kurtosis_instruction.InstructionPosition {
	return instruction.position
}

func (instruction *WaitInstruction) GetCanonicalInstruction() *kurtosis_core_rpc_api_bindings.StarlarkInstruction {
	args := []*kurtosis_core_rpc_api_bindings.StarlarkInstructionArg{
		binding_constructors.NewStarlarkInstructionKwarg(shared_helpers.CanonicalizeArgValue(instruction.starlarkKwargs[recipeArgName]), recipeArgName, kurtosis_instruction.Representative),
		binding_constructors.NewStarlarkInstructionKwarg(shared_helpers.CanonicalizeArgValue(instruction.starlarkKwargs[valueFieldArgName]), valueFieldArgName, kurtosis_instruction.Representative),
		binding_constructors.NewStarlarkInstructionKwarg(shared_helpers.CanonicalizeArgValue(instruction.starlarkKwargs[assertionArgName]), assertionArgName, kurtosis_instruction.Representative),
		binding_constructors.NewStarlarkInstructionKwarg(shared_helpers.CanonicalizeArgValue(instruction.starlarkKwargs[targetArgName]), targetArgName, kurtosis_instruction.Representative),
		binding_constructors.NewStarlarkInstructionKwarg(shared_helpers.CanonicalizeArgValue(instruction.starlarkKwargs[optionalIntervalArgName]), optionalIntervalArgName, kurtosis_instruction.NotRepresentative),
		binding_constructors.NewStarlarkInstructionKwarg(shared_helpers.CanonicalizeArgValue(instruction.starlarkKwargs[optionalTimeoutArgName]), optionalTimeoutArgName, kurtosis_instruction.NotRepresentative),
	}
	return binding_constructors.NewStarlarkInstruction(instruction.position.ToAPIType(), WaitBuiltinName, instruction.String(), args)
}

func (instruction *WaitInstruction) Execute(ctx context.Context) (*string, error) {
	var (
		requestErr error
		assertErr  error
	)
	tries := 0
	timedOut := false
	lastResult := map[string]starlark.Comparable{}
	startTime := time.Now()
	for {
		tries += 1
		backoffDuration := instruction.backoff.NextBackOff()
		if backoffDuration == backoff.Stop || time.Since(startTime) > instruction.timeout {
			timedOut = true
			break
		}
		lastResult, requestErr = instruction.recipe.Execute(ctx, instruction.serviceNetwork, instruction.runtimeValueStore)
		if requestErr != nil {
			time.Sleep(backoffDuration)
			continue
		}
		instruction.runtimeValueStore.SetValue(instruction.resultUuid, lastResult)
		value, found := lastResult[instruction.targetKey]
		if !found {
			return nil, stacktrace.NewError("Error grabbing value from key '%v'", instruction.targetKey)
		}
		assertErr = assert.Assert(value, instruction.assertion, instruction.target)
		if assertErr != nil {
			time.Sleep(backoffDuration)
			continue
		}
		break
	}
	if timedOut {
		return nil, stacktrace.NewError("Wait timed-out waiting for the assertion to become valid. Waited for '%v'. Last assertion error was: \n%v", time.Since(startTime), assertErr)
	}
	if requestErr != nil {
		return nil, stacktrace.Propagate(requestErr, "Error executing HTTP recipe")
	}
	if assertErr != nil {
		return nil, stacktrace.Propagate(assertErr, "Error asserting HTTP recipe on '%v'", WaitBuiltinName)
	}
	instructionResult := fmt.Sprintf("Wait took %d tries (%v in total). Assertion passed with following:\n%s", tries, time.Since(startTime), instruction.recipe.ResultMapToString(lastResult))
	return &instructionResult, nil
}

func (instruction *WaitInstruction) String() string {
	return shared_helpers.CanonicalizeInstruction(WaitBuiltinName, kurtosis_instruction.NoArgs, instruction.starlarkKwargs)
}

func (instruction *WaitInstruction) ValidateAndUpdateEnvironment(environment *startosis_validator.ValidatorEnvironment) error {
	// TODO(vcolombo): Add validation step here
	return nil
}

func (instruction *WaitInstruction) parseStartosisArgs(b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) *startosis_errors.InterpretationError {
	var (
		recipeConfigArg  starlark.Value
		valueFieldArg    starlark.String
		assertionArg     starlark.String
		targetArg        starlark.Comparable
		optionalInterval starlark.String = emptyOptionalField
		optionalTimeout  starlark.String = emptyOptionalField
	)

	if err := starlark.UnpackArgs(b.Name(), args, kwargs, recipeArgName, &recipeConfigArg, valueFieldArgName, &valueFieldArg, assertionArgName, &assertionArg, targetArgName, &targetArg, optionalIntervalArgName, &optionalInterval, optionalTimeoutArgName, &optionalTimeout); err != nil {
		return startosis_errors.NewInterpretationError(err.Error())
	}

	instruction.starlarkKwargs = starlark.StringDict{
		recipeArgName:           recipeConfigArg,
		valueFieldArgName:       valueFieldArg,
		assertionArgName:        assertionArg,
		targetArgName:           targetArg,
		optionalIntervalArgName: optionalInterval,
		optionalTimeoutArgName:  optionalTimeout,
	}
	instruction.starlarkKwargs.Freeze()

	var ok bool
	instruction.recipe, ok = recipeConfigArg.(*recipe.HttpRequestRecipe)
	if !ok {
		instruction.recipe, ok = recipeConfigArg.(*recipe.ExecRecipe)
		if !ok {
			// TODO: remove this after 2 or 3 weeks?
			err := ensureBackwardCompatibleForWait(instruction, recipeConfigArg)
			if err != nil {
				return startosis_errors.NewInterpretationError("There was no valid recipe found for '%v' "+
					"(Expected ExecRecipe or PostHttpRequestRecipe or GetHttpRequestRecipe type)", recipeConfigArg)
			}
		}
	}

	instruction.assertion = string(assertionArg)
	instruction.target = targetArg
	instruction.targetKey = string(valueFieldArg)
	if optionalInterval != emptyOptionalField {
		interval, parseErr := time.ParseDuration(optionalInterval.GoString())
		if parseErr != nil {
			return startosis_errors.WrapWithInterpretationError(parseErr, "An error occurred when parsing interval '%v'", optionalInterval.GoString())
		}
		instruction.backoff = backoff.NewConstantBackOff(interval)
	} else {
		instruction.backoff = backoff.NewConstantBackOff(defaultInterval)
	}

	if optionalTimeout != emptyOptionalField {
		timeout, parseErr := time.ParseDuration(optionalTimeout.GoString())
		if parseErr != nil {
			return startosis_errors.NewInterpretationError("An error occurred when parsing timeout '%v'", optionalTimeout)
		}
		instruction.timeout = timeout
	} else {
		instruction.timeout = defaultTimeout
	}

	if _, found := assert.StringTokenToComparisonStarlarkToken[instruction.assertion]; !found && instruction.assertion != "IN" && instruction.assertion != "NOT_IN" {
		return startosis_errors.NewInterpretationError("'%v' is not a valid assertion", assertionArg)
	}
	if _, ok := instruction.target.(*starlark.List); (instruction.assertion == "IN" || instruction.assertion == "NOT_IN") && !ok {
		return startosis_errors.NewInterpretationError("'%v' assertion requires list, got '%v'", assertionArg, targetArg.Type())
	}
	return nil
}

// TODO: remove this code after we stop using this -- maybe in 2 or 3 weeks?
func ensureBackwardCompatibleForWait(instruction *WaitInstruction, recipeConfigArg starlark.Value) *startosis_errors.InterpretationError {
	var err *startosis_errors.InterpretationError
	recipeConfigStruct, ok := recipeConfigArg.(*starlarkstruct.Struct)
	if !ok {
		return startosis_errors.NewInterpretationError("Error occurred while parsing starlark value to starlark struct %v", recipeConfigArg)
	}

	instruction.recipe, err = kurtosis_instruction.ParseHttpRequestRecipe(recipeConfigStruct)
	if err != nil {
		instruction.recipe, err = kurtosis_instruction.ParseExecRecipe(recipeConfigStruct)
	}

	return err
}
