package startosis_module_test

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-cli/golang_internal_testsuite/test_helpers"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
)

const (
	invalidCaseModFileTestName          = "invalid-module-invalid-mod-file"
	moduleWithInvalidKurtosisModRelPath = "../../../startosis/invalid-mod-file"
)

func TestStartosisModule_InvalidModFile(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// ------------------------------------- ENGINE SETUP ----------------------------------------------
	enclaveCtx, destroyEnclaveFunc, _, err := test_helpers.CreateEnclave(t, ctx, invalidCaseModFileTestName, isPartitioningEnabled)
	require.NoError(t, err, "An error occurred creating an enclave")
	defer destroyEnclaveFunc()

	currentWorkingDirectory, err := os.Getwd()
	require.Nil(t, err)
	moduleDirpath := path.Join(currentWorkingDirectory, moduleWithInvalidKurtosisModRelPath)

	// ------------------------------------- TEST RUN ----------------------------------------------
	logrus.Info("Executing Startosis Module...")

	logrus.Infof("Startosis module path: \n%v", moduleDirpath)

	expectedErrorContents := "Field module.name in kurtosis.yml needs to be set and cannot be empty"
	_, _, err = enclaveCtx.ExecuteKurtosisModule(ctx, moduleDirpath, emptyExecuteParams, defaultDryRun)
	require.NotNil(t, err, "Unexpected error executing startosis module")
	require.Contains(t, err.Error(), expectedErrorContents)
}
