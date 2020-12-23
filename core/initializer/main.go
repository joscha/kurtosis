/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package main

import (
	"encoding/json"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/kurtosis-tech/kurtosis/commons/logrus_log_levels"
	"github.com/kurtosis-tech/kurtosis/initializer/auth/access_controller"
	"github.com/kurtosis-tech/kurtosis/initializer/auth/auth0_authorizers"
	"github.com/kurtosis-tech/kurtosis/initializer/auth/auth0_constants"
	"github.com/kurtosis-tech/kurtosis/initializer/auth/session_cache"
	"github.com/kurtosis-tech/kurtosis/initializer/docker_flag_parser"
	"github.com/kurtosis-tech/kurtosis/initializer/initializer_container_constants"
	"github.com/kurtosis-tech/kurtosis/initializer/test_suite_constants"
	"github.com/kurtosis-tech/kurtosis/initializer/test_suite_metadata_acquirer"
	"github.com/kurtosis-tech/kurtosis/initializer/test_suite_runner"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"sort"
	"strings"
)

const (

	successExitCode = 0
	failureExitCode = 1

	// We don't want to overwhelm slow machines, since it becomes not-obvious what's happening
	defaultParallelism = 2

	// Web link shown to users who do not authenticate.
	licenseWebUrl = "https://kurtosistech.com/"

	// The location on the INITIALIZER container where the suite execution volume will be mounted
	// A user MUST mount a volume here
	initializerContainerSuiteExVolMountDirpath = "/suite-execution"

	// The location on the INITIALIZER container where the Kurtosis storage directory (containing things like JWT
	//  tokens) will be bind-mounted from the host filesystem
	storageDirectoryBindMountDirpath = "/kurtosis"

	// Name of the file within the Kurtosis storage directory where the session cache will be stored
	sessionCacheFilename = "session-cache"
	sessionCacheFileMode os.FileMode = 0600

	// Can make these configurable if needed
	hostPortTrackerInterfaceIp = "127.0.0.1"
	hostPortTrackerStartRange = 8000
	hostPortTrackerEndRange = 10000

	defaultDebuggerPort = 2778
	debuggerPortProtocol = "tcp"


	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! IMPORTANT !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	//                  If you change the below, you need to update the Dockerfile!
	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! IMPORTANT !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	doListArg = "DO_LIST"
	testSuiteImageArg = "TEST_SUITE_IMAGE"
	testNamesArg = "TEST_NAMES"
	kurtosisLogLevelArg = "KURTOSIS_LOG_LEVEL"
	testSuiteLogLevelArg = "TEST_SUITE_LOG_LEVEL"
	clientIdArg = "CLIENT_ID"
	clientSecretArg = "CLIENT_SECRET"
	kurtosisApiImageArg = "KURTOSIS_API_IMAGE"
	parallelismArg = "PARALLELISM"
	customEnvVarsJsonArg = "CUSTOM_ENV_VARS_JSON"
	suiteExecutionVolumeArg = "SUITE_EXECUTION_VOLUME"
	testSuiteDebuggerPortArg = "DEBUGGER_PORT"
	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! IMPORTANT !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	//                     If you change the above, you need to update the Dockerfile!
	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! IMPORTANT !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
)


// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! IMPORTANT !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//          If you change default values below, you need to update the Dockerfile!
// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! IMPORTANT !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
var flagConfigs = map[string]docker_flag_parser.FlagConfig{
	clientIdArg: {
		Required: false,
		Default:  "",
		HelpText: fmt.Sprintf("An OAuth client ID which is needed for running Kurtosis in CI, and should be left empty when running Kurtosis on a local machine"),
		Type:     docker_flag_parser.StringFlagType,
	},
	clientSecretArg: {
		Required: false,
		Default:  "",
		HelpText: fmt.Sprintf("An OAuth client secret which is needed for running Kurtosis in CI, and should be left empty when running Kurtosis on a local machine"),
		Type:     docker_flag_parser.StringFlagType,
	},
	customEnvVarsJsonArg: {
		Required: false,
		Default:  "{}",
		HelpText: "JSON containing key-value mappings of custom environment variables that will be passed through to the test suite container, e.g. '{\"MY_VAR\":\\ \"/some/value\"}' (note the escaped spaces!)",
		Type:     docker_flag_parser.StringFlagType,
	},
	doListArg: {
		Required: false,
		Default:  false,
		HelpText: "Rather than running the tests, lists the tests available to run",
		Type:     docker_flag_parser.BoolFlagType,
	},
	kurtosisApiImageArg: {
		Required: true,
		Default:  "",
		HelpText: "The Docker image from the Kurtosis API image repo (https://hub.docker.com/repository/docker/kurtosistech/kurtosis-core_api) that will be used during operation",
		Type:     docker_flag_parser.StringFlagType,
	},
	kurtosisLogLevelArg: {
		Required: false,
		Default: "info",
		HelpText: fmt.Sprintf(
			"The log level that all output generated by the Kurtosis framework itself should log at (%v)",
			strings.Join(logrus_log_levels.AcceptableLogLevels, "|"),
		),
		Type: docker_flag_parser.StringFlagType,
	},
	parallelismArg: {
		Required: false,
		Default:  defaultParallelism,
		HelpText: "A positive integer telling Kurtosis how many tests to run concurrently (should be set no higher than the number of cores on your machine, else you'll slow down your tests and potentially hit test timeouts!)",
		Type:     docker_flag_parser.IntFlagType,
	},
	suiteExecutionVolumeArg: {
		Required: true,
		Default:  "",
		HelpText: "The name of the Docker volume that will contain all the data for the test suite execution (should be a new volume for each execution!)",
		Type:     docker_flag_parser.StringFlagType,
	},
	testNamesArg: {
		Required: false,
		Default:  "",
		HelpText: "List of test names to run, separated by '" + initializer_container_constants.TestNameArgSeparator + "' (default or empty: run all tests)",
		Type:     docker_flag_parser.StringFlagType,
	},
	testSuiteDebuggerPortArg: {
		Required: false,
		Default: defaultDebuggerPort,
		HelpText: "The port that debuggers running inside the testsuite should listen on, which Kurtosis will expose" +
			"to the host machine",
		Type:   docker_flag_parser.IntFlagType,
	},
	testSuiteImageArg: {
		Required: true,
		Default:  "",
		HelpText: "The name of the Docker image containing your test suite to run",
		Type:     docker_flag_parser.StringFlagType,
	},
	testSuiteLogLevelArg: {
		Required: false,
		Default:  "debug",
		HelpText: fmt.Sprintf("A string that will be passed as-is to the test suite container to indicate " +
			"what log level the test suite container should output at; this string should be meaningful to " +
			"the test suite container because Kurtosis won't know what logging framework the testsuite uses"),
		Type:     docker_flag_parser.StringFlagType,
	},
}
// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! IMPORTANT !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//             If you change default values above, you need to update the Dockerfile!
// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! IMPORTANT !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!


func main() {
	// NOTE: we'll want to chnage the ForceColors to false if we ever want structured logging
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	flagParser := docker_flag_parser.NewFlagParser(flagConfigs)
	parsedFlags, err := flagParser.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred parsing the initializer CLI flags: %v\n", err)
		os.Exit(failureExitCode)
	}

	kurtosisLevel, err := logrus.ParseLevel(parsedFlags.GetString(kurtosisLogLevelArg))
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred parsing the Kurtosis log level string: %v\n", err)
		os.Exit(failureExitCode)
	}
	logrus.SetLevel(kurtosisLevel)

	// TODO Break this into a private helper function
	clientId := parsedFlags.GetString(clientIdArg)
	clientSecret := parsedFlags.GetString(clientSecretArg)
	var accessController access_controller.AccessController
	if len(clientId) > 0 && len(clientSecret) > 0 {
		logrus.Debugf("Running CI machine-to-machine auth flow...")
		accessController = access_controller.NewClientAuthAccessController(
			auth0_constants.RsaPublicKeyCertsPem,
			auth0_authorizers.NewStandardClientCredentialsAuthorizer(),
			clientId,
			clientSecret)
	} else {
		logrus.Debugf("Running developer device auth flow...")
		sessionCacheFilepath := path.Join(storageDirectoryBindMountDirpath, sessionCacheFilename)
		sessionCache := session_cache.NewEncryptedSessionCache(
			sessionCacheFilepath,
			sessionCacheFileMode,
		)
		accessController = access_controller.NewDeviceAuthAccessController(
			auth0_constants.RsaPublicKeyCertsPem,
			sessionCache,
			auth0_authorizers.NewStandardDeviceAuthorizer(),
		)
	}
	if err := accessController.Authorize(); err != nil {
		logrus.Fatalf(
			"The following error occurred when authenticating and authorizing your Kurtosis license: %v",
			err)
		os.Exit(failureExitCode)
	}

	// TODO Break everything from here on down into a private helper function
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logrus.Errorf("An error occurred creating the Docker client: %v", err)
		os.Exit(failureExitCode)
	}

	// Parse environment variables
	customEnvVarsJson := parsedFlags.GetString(customEnvVarsJsonArg)
	var customEnvVars map[string]string
	if err := json.Unmarshal([]byte(customEnvVarsJson), &customEnvVars); err != nil {
		logrus.Errorf("An error occurred parsing the custom environment variables JSON: %v", err)
		os.Exit(failureExitCode)
	}

	freeHostPortBindingSupplier, err := test_suite_runner.NewFreeHostPortBindingSupplier(
		hostPortTrackerInterfaceIp,
		debuggerPortProtocol,
		hostPortTrackerStartRange,
		hostPortTrackerEndRange)
	if err != nil {
		logrus.Errorf("An error occurred creating the free host port binding supplier: %v", err)
		os.Exit(failureExitCode)
	}

	testsuiteLauncher, err := test_suite_constants.NewTestsuiteContainerLauncher(
		parsedFlags.GetString(testSuiteImageArg),
		parsedFlags.GetString(testSuiteLogLevelArg),
		customEnvVars,
		parsedFlags.GetInt(testSuiteDebuggerPortArg))
	if err != nil {
		logrus.Errorf("An error occurred creating the testsuite launcher: %v", err)
		os.Exit(failureExitCode)
	}

	metadataAcquisitionHostPortBinding, err := freeHostPortBindingSupplier.GetFreePortBinding()
	if err != nil {
		logrus.Errorf("An error occurred getting the test suite metadata: %v", err)
		os.Exit(failureExitCode)
	}

	suiteMetadata, err := test_suite_metadata_acquirer.GetTestSuiteMetadata(
		parsedFlags.GetString(suiteExecutionVolumeArg),
		initializerContainerSuiteExVolMountDirpath,
		dockerClient,
		testsuiteLauncher,
		metadataAcquisitionHostPortBinding)
	if err != nil {
		logrus.Errorf("An error occurred getting the test suite metadata: %v", err)
		os.Exit(failureExitCode)
	}

	// If any test names have our special test name arg separator, we won't be able to select the test so throw an
	//  error and loudly alert the user
	for testName, _ := range suiteMetadata.TestNames {
		if strings.Contains(testName, initializer_container_constants.TestNameArgSeparator) {
			logrus.Errorf(
				"Test '%v' contains illegal character '%v'; we use this character for delimiting when choosing which tests to run so test names cannot contain it!",
				testName,
				initializer_container_constants.TestNameArgSeparator)
			os.Exit(failureExitCode)
		}
	}

	if parsedFlags.GetBool(doListArg) {
		testNames := []string{}
		for name := range suiteMetadata.TestNames {
			testNames = append(testNames, name)
		}
		sort.Strings(testNames)

		fmt.Println("\nTests in test suite:")
		for _, name := range testNames {
			// We intentionally don't use Logrus here so that we always see the output, even with a misconfigured loglevel
			fmt.Println("- " + name)
		}
		os.Exit(successExitCode)
	}


	// Split user-input string into actual candidate test names
	testNamesArgStr := strings.TrimSpace(parsedFlags.GetString(testNamesArg))
	testNamesToRun := map[string]bool{}
	if len(testNamesArgStr) > 0 {
		testNamesList := strings.Split(testNamesArgStr, initializer_container_constants.TestNameArgSeparator)
		for _, name := range testNamesList {
			testNamesToRun[name] = true
		}
	}

	parallelismUint := uint(parsedFlags.GetInt(parallelismArg))
	allTestsPassed, err := test_suite_runner.RunTests(
		dockerClient,
		parsedFlags.GetString(suiteExecutionVolumeArg),
		initializerContainerSuiteExVolMountDirpath,
		*suiteMetadata,
		testNamesToRun,
		parallelismUint,
		parsedFlags.GetString(kurtosisApiImageArg),
		parsedFlags.GetString(kurtosisLogLevelArg),
		testsuiteLauncher,
		freeHostPortBindingSupplier)
	if err != nil {
		logrus.Errorf("An error occurred running the tests:")
		fmt.Fprintln(logrus.StandardLogger().Out, err)
		os.Exit(failureExitCode)
	}

	var exitCode int
	if allTestsPassed {
		exitCode = successExitCode
	} else {
		exitCode = failureExitCode
	}
	os.Exit(exitCode)
}
