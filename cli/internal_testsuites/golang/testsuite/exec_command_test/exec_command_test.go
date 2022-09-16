package exec_command_test

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-cli/golang_internal_testsuite/test_helpers"
	"github.com/kurtosis-tech/kurtosis-core-api-lib/api/golang/lib/services"
	"github.com/kurtosis-tech/stacktrace"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	testName = "exec-command"
	isPartitioningEnabled = false

	execCmdTestImage      = "alpine:3.12.4"
	inputForLogOutputTest = "hello"
	expectedLogOutput     = "hello\n"
	inputForAdvancedLogOutputTest = "hello && hello"
	expectedAdvancedLogOutput = "hello && hello\n"
	testServiceId         = "test"

	successExitCode int32 = 0

)

var execCommandThatShouldWork = []string{
	"true",
}

var execCommandThatShouldHaveLogOutput = []string{
	"echo",
	inputForLogOutputTest,
}

// This command tests to ensure that the commands the user is running get passed exactly as-is to the Docker
// container. If Kurtosis code is magically wrapping the code with "sh -c", this will fail.
var execCommandThatWillFailIfShWrapped = []string{
	"echo",
	inputForAdvancedLogOutputTest,
}

var execCommandThatShouldFail = []string{
	"false",
}

func TestExecCommand(t *testing.T) {
	ctx := context.Background()

	// ------------------------------------- ENGINE SETUP ----------------------------------------------
	enclaveCtx, stopEnclaveFunc, _, err := test_helpers.CreateEnclave(t, ctx, testName, isPartitioningEnabled)
	require.NoError(t, err, "An error occurred creating an enclave")
	defer stopEnclaveFunc()

	// ------------------------------------- TEST SETUP ----------------------------------------------
	containerConfig := getContainerConfig()

	testServiceContext, err := enclaveCtx.AddService(testServiceId, containerConfig)
	require.NoError(t, err, "An error occurred starting service '%v'", testServiceId)

	// ------------------------------------- TEST RUN ----------------------------------------------
	logrus.Infof("Running exec command '%+v' that should return a successful exit code...", execCommandThatShouldWork)
	shouldWorkExitCode, _, err := runExecCmd(testServiceContext, execCommandThatShouldWork)
	require.NoError(t, err, "An error occurred running exec command '%+v'", execCommandThatShouldWork)
	require.Equal(t, successExitCode, shouldWorkExitCode, "Exec command '%+v' should work, but got unsuccessful exit code %v", execCommandThatShouldWork, shouldWorkExitCode)
	logrus.Info("Exec command returned successful exit code as expected")

	logrus.Infof("Running exec command '%+v' that should return an error exit code...", execCommandThatShouldFail)
	shouldFailExitCode, _, err := runExecCmd(testServiceContext, execCommandThatShouldFail)
	require.NoError(t, err, "An error occurred running exec command '%+v'", execCommandThatShouldFail)
	require.NotEqual(t, successExitCode, shouldFailExitCode, "Exec command '%+v' should fail, but got successful exit code %v", execCommandThatShouldFail, successExitCode)
	logrus.Infof("Exec command returning an error exited with error")

	logrus.Infof("Running exec command '%+v' that should return log output...", execCommandThatShouldHaveLogOutput)
	shouldHaveLogOutputExitCode, logOutput, err := runExecCmd(testServiceContext, execCommandThatShouldHaveLogOutput)
	require.NoError(t, err, "An error occurred running exec command '%+v'", execCommandThatShouldHaveLogOutput)
	require.Equal(t, successExitCode, shouldHaveLogOutputExitCode, "Exec command '%+v' should work, but got unsuccessful exit code %v", execCommandThatShouldHaveLogOutput, shouldHaveLogOutputExitCode)
	require.Equal(t, expectedLogOutput, logOutput, "Exec command '%+v' should return '%v', but got '%v'.", execCommandThatShouldHaveLogOutput, expectedLogOutput, logOutput)
	logrus.Info("Exec command returning log output passed as expected")

	logrus.Infof("Running exec command '%+v' that will fail if Kurtosis is accidentally sh-wrapping the command...", execCommandThatWillFailIfShWrapped)
	shouldNotGetShWrappedExitCode, shouldNotGetShWrappedLogOutput, err := runExecCmd(testServiceContext, execCommandThatWillFailIfShWrapped)
	require.NoError(t, err, "An error occurred running exec command '%+v'", execCommandThatWillFailIfShWrapped)
	require.Equal(t, successExitCode, shouldNotGetShWrappedExitCode, "Exec command '%v' should work, but got unsuccessful exit code %v", execCommandThatWillFailIfShWrapped, shouldNotGetShWrappedExitCode)
	require.Equal(t, expectedAdvancedLogOutput, shouldNotGetShWrappedLogOutput, "Exec command '%+v' should return '%v', but got '%v'.", execCommandThatWillFailIfShWrapped, expectedAdvancedLogOutput, shouldNotGetShWrappedLogOutput)
	logrus.Info("Exec command that will fail if Kurtosis is accidentally sh-wrapping did not fail")
}

// ====================================================================================================
//                                       Private helper functions
// ====================================================================================================
func getContainerConfig() *services.ContainerConfig {
	// We sleep because the only function of this container is to test Docker executing a command while it's running
	// NOTE: We could just as easily combine this into a single array (rather than splitting between ENTRYPOINT and CMD
	// args), but this provides a nice little regression test of the ENTRYPOINT overriding
	entrypointArgs := []string{"sleep"}
	cmdArgs := []string{"30"}

	containerConfig := services.NewContainerConfigBuilder(
		execCmdTestImage,
	).WithEntrypointOverride(
		entrypointArgs,
	).WithCmdOverride(
		cmdArgs,
	).Build()
	return containerConfig
}

func runExecCmd(serviceContext *services.ServiceContext, command []string) (int32, string, error) {
	exitCode, logOutput, err := serviceContext.ExecCommand(command)
	if err != nil {
		return 0, "", stacktrace.Propagate(
			err,
			"An error occurred executing command '%v'", command)
	}
	return exitCode, logOutput, nil
}
