package service_pause_test

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-cli/golang_internal_testsuite/test_helpers"
	"github.com/kurtosis-tech/kurtosis-core-api-lib/api/golang/lib/services"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	testName                  = "pause-unpause"
	isPartitioningEnabled     = false
	pauseUnpauseTestImageName = "alpine:3.12.4"
	testServiceId             = "test"
)

func TestPauseUnpause(t *testing.T) {
	ctx := context.Background()
	// ------------------------------------- ENGINE SETUP ----------------------------------------------
	enclaveCtx, stopEnclaveFunc, _, err := test_helpers.CreateEnclave(t, ctx, testName, isPartitioningEnabled)
	require.NoError(t, err, "An error occurred creating an enclave")
	defer stopEnclaveFunc()

	// ------------------------------------- TEST SETUP ----------------------------------------------
	containerConfigSupplier := getContainerConfigSupplier()

	serviceCtx, err := enclaveCtx.AddService(testServiceId, containerConfigSupplier)
	require.NoError(t, err, "An error occurred adding the file server service")

	time.Sleep(10 * time.Second)
	// ------------------------------------- TEST RUN ----------------------------------------------
	// pause/unpause using servicectx
	err = enclaveCtx.PauseService(serviceCtx.GetServiceID())
	logrus.Infof("Paused service!")
	require.NoError(t, err, "An error occurred unpausing")
	time.Sleep(10 * time.Second)
	err = enclaveCtx.UnpauseService(serviceCtx.GetServiceID())
	require.NoError(t, err, "An error occurred unpausing")
	logrus.Infof("Unpaused service!")
	time.Sleep(10 * time.Second)
}

// ====================================================================================================
//                                       Private helper functions
// ====================================================================================================
func getContainerConfigSupplier() func(ipAddr string) (*services.ContainerConfig, error) {
	containerConfigSupplier := func(ipAddr string) (*services.ContainerConfig, error) {

		// We spam timestamps so that we can measure pausing processes (no more log output) and unpausing (log output resumes)
		entrypointArgs := []string{"/bin/sh", "-c"}
		cmdArgs := []string{"while sleep 1; do ts=$(date +\"%s\") ; echo \"Time: $ts\" ; done"}

		containerConfig := services.NewContainerConfigBuilder(
			pauseUnpauseTestImageName,
		).WithEntrypointOverride(
			entrypointArgs,
		).WithCmdOverride(
			cmdArgs,
		).Build()
		return containerConfig, nil
	}
	return containerConfigSupplier
}
