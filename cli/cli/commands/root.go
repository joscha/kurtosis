/*
 * Copyright (c) 2021 - present Kurtosis Technologies Inc.
 * All Rights Reserved.
 */

package commands

import (
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/kurtosis-tech/kurtosis-cli/cli/commands/clean"
	"github.com/kurtosis-tech/kurtosis-cli/cli/commands/enclave"
	"github.com/kurtosis-tech/kurtosis-cli/cli/commands/engine"
	"github.com/kurtosis-tech/kurtosis-cli/cli/commands/module"
	"github.com/kurtosis-tech/kurtosis-cli/cli/commands/repl"
	"github.com/kurtosis-tech/kurtosis-cli/cli/commands/sandbox"
	"github.com/kurtosis-tech/kurtosis-cli/cli/commands/service"
	"github.com/kurtosis-tech/kurtosis-cli/cli/commands/test"
	"github.com/kurtosis-tech/kurtosis-cli/cli/commands/version"
	"github.com/kurtosis-tech/kurtosis-cli/cli/helpers/host_machine_directories"
	"github.com/kurtosis-tech/kurtosis-cli/cli/helpers/logrus_log_levels"
	"github.com/kurtosis-tech/kurtosis-cli/cli/kurtosis_cli_version"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	logLevelStrArg = "cli-log-level"

	latestReleaseOnGitHubURL   = "https://api.github.com/repos/kurtosis-tech/kurtosis-cli-release-artifacts/releases/latest"
	acceptHttpHeaderKey        = "Accept"
	acceptHttpHeaderValue      = "application/json"
	contentTypeHttpHeaderKey   = "Content-Type"
	contentTypeHttpHeaderValue = "application/json"
	userAgentHttpHeaderKey     = "User-Agent"
	userAgentHttpHeaderValue   = "kurtosis-tech"

	upgradeCLIInstructionsDocsPageURL = "https://docs.kurtosistech.com/installation.html#upgrading-kurtosis-cli"

	latestCLIReleaseCacheFileContentSeparator               = ";"
	latestCLIReleaseCacheFilePermissionsForSaveData         = 0644
	latestCLIReleaseCacheFilePermissionsForOpenOrCreateFile = 0755
	latestCLIReleaseCacheFileExpirationHours                = 24
	latestCLIReleaseCacheFileContentColumnsAmount           = 2
	latestCLIReleaseCacheFileContentDateKey                 = 0
	latestCLIReleaseCacheFileContentVersionKey              = 1

	zeroBytes = 0
)

type GitHubReleaseReponse struct {
	TagName string `json:"tag_name"`
}

var logLevelStr string
var defaultLogLevelStr = logrus.InfoLevel.String()

var RootCmd = &cobra.Command{
	// Leaving out the "use" will auto-use os.Args[0]
	Use:   "",
	Short: "A CLI for interacting with the Kurtosis engine",

	// Cobra will print usage whenever _any_ error occurs, including ones we throw in Kurtosis
	// This doesn't make sense in 99% of the cases, so just turn them off entirely
	SilenceUsage:      true,
	PersistentPreRunE: globalSetup,
}

func init() {
	RootCmd.PersistentFlags().StringVar(
		&logLevelStr,
		logLevelStrArg,
		defaultLogLevelStr,
		"Sets the level that the CLI will log at ("+strings.Join(logrus_log_levels.GetAcceptableLogLevelStrs(), "|")+")",
	)

	RootCmd.AddCommand(sandbox.SandboxCmd)
	RootCmd.AddCommand(test.TestCmd)
	RootCmd.AddCommand(enclave.EnclaveCmd)
	RootCmd.AddCommand(service.ServiceCmd)
	RootCmd.AddCommand(module.ModuleCmd)
	RootCmd.AddCommand(repl.REPLCmd)
	RootCmd.AddCommand(engine.EngineCmd)
	RootCmd.AddCommand(version.VersionCmd)
	RootCmd.AddCommand(clean.CleanCmd)
}

// ====================================================================================================
//                                       Private Helper Functions
// ====================================================================================================
func globalSetup(cmd *cobra.Command, args []string) error {
	logLevel, err := logrus.ParseLevel(logLevelStr)
	if err != nil {
		return stacktrace.Propagate(err, "Could not parse log level string '%v'", logLevelStr)
	}
	logrus.SetOutput(cmd.OutOrStdout())
	logrus.SetLevel(logLevel)

	checkCLIVersion()

	return nil
}

func checkCLIVersion() {
	isLatestVersion, latestVersion, err := isLatestCLIVersion()
	if err != nil {
		logrus.Warning("An error occurred trying to check if you are running the latest Kurtosis CLI version.")
		logrus.Debugf("Checking latest version error: %v", err)
		logrus.Warningf("Your current version is '%v'", kurtosis_cli_version.KurtosisCLIVersion)
		logrus.Warningf("You can manually upgrade the CLI tool following these instructions: %v", upgradeCLIInstructionsDocsPageURL)
		return
	}
	if !isLatestVersion {
		logrus.Warningf("You are running an old version of the Kurtosis CLI; we suggest you to update it to the latest version, '%v'", latestVersion)
		logrus.Warningf("You can manually upgrade the CLI tool following these instructions: %v", upgradeCLIInstructionsDocsPageURL)
	}
	return
}

func isLatestCLIVersion() (bool, string, error) {
	ownVersionStr := kurtosis_cli_version.KurtosisCLIVersion
	latestVersionStr, err := getLatestCLIReleaseVersion()
	if err != nil {
		return false, "", stacktrace.Propagate(err, "An error occurred getting the latest release version number from the GitHub public API")
	}

	ownSemver, err := semver.StrictNewVersion(ownVersionStr)
	if err != nil {
		return false, "", stacktrace.Propagate(err, "An error occurred parsing own version string '%v' to sem version", ownVersionStr)
	}

	latestSemver, err := semver.StrictNewVersion(latestVersionStr)
	if err != nil {
		return false, "", stacktrace.Propagate(err, "An error occurred parsing latest version string '%v' to sem version", latestVersionStr)
	}

	compareResult := ownSemver.Compare(latestSemver)

	//compareResult = 1  means that the own version is newer than the latest version, (e.g.: during a new release)
	if compareResult >= 0 {
		return true, latestVersionStr, nil
	}

	return false, latestVersionStr, nil
}

func getLatestCLIReleaseVersion() (string, error) {
	latestCLIVersion, err := getLatestCLIReleaseVersionFromCacheFile()
	if err != nil {
		logrus.Debugf("An error occurred getting latest released CLI version from cache file. Error: \n%v", err)
	}
	if latestCLIVersion != "" {
		logrus.Debugf("Getting the latest released CLI version '%v' from the cache file", latestCLIVersion)
		return latestCLIVersion, nil
	}

	latestCLIVersion, err = getLatestCLIReleaseVersionFromGitHub()
	if err != nil {
		return "", stacktrace.Propagate(err, "An error occurred getting latest released CLI version from GitHub")
	}

	logrus.Debugf("Getting the latest released CLI version '%v' from GitHub API", latestCLIVersion)
	return latestCLIVersion, nil
}

func getLatestCLIReleaseVersionFromGitHub() (string, error) {
	var (
		client         = &http.Client{}
		requestMethod  = "GET"
		requestBody    io.Reader
		responseObject GitHubReleaseReponse
	)

	request, err := http.NewRequest(requestMethod, latestReleaseOnGitHubURL, requestBody)
	if err != nil {
		return "", stacktrace.Propagate(err, "An error occurred creating new HTTP GET request to URL '%v' ", latestReleaseOnGitHubURL)
	}

	request.Header.Add(acceptHttpHeaderKey, acceptHttpHeaderValue)
	request.Header.Add(contentTypeHttpHeaderKey, contentTypeHttpHeaderValue)
	request.Header.Add(userAgentHttpHeaderKey, userAgentHttpHeaderValue)

	response, err := client.Do(request)
	if err != nil {
		return "", stacktrace.Propagate(err, "An error occurred executing HTTP GET request to URL '%v' ", latestReleaseOnGitHubURL)
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			logrus.Warnf("We tried to close the response body, but doing so threw an error:\n%v", err)
		}
	}()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", stacktrace.Propagate(err, "An error occurred reading the HTTP response body")
	}

	if err := json.Unmarshal(bodyBytes, &responseObject); err != nil {
		return "", stacktrace.Propagate(err, "An error occurred deserializing the latest release body response")
	}

	latestVersion := responseObject.TagName
	if latestVersion == "" {
		return "", stacktrace.Propagate(err, "The latest release version got from GitHub releases is empty")
	}

	if err := saveLatestCLIReleaseVersionInCacheFile(latestVersion); err != nil {
		logrus.Debugf("We tried to save the latest release version '%v' in the cache file, but doing so threw an error:\n%v", latestVersion, err)
	}

	return latestVersion, nil
}

func saveLatestCLIReleaseVersionInCacheFile(latestReleaseVersion string) error {

	filepath, err := host_machine_directories.GetLatestCLIReleaseVersionCacheFilepath()
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred getting the latest release version cache filepath")
	}
	logrus.Debugf("Cache filepath: '%v'", filepath)

	cacheFile, err := os.OpenFile(filepath, os.O_WRONLY, latestCLIReleaseCacheFilePermissionsForSaveData)
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred opening the '%v' file", filepath)
	}
	defer func() {
		if err := cacheFile.Close(); err != nil {
			logrus.Warnf("We tried to close the latest release CLI version cache file, but doing so threw an error:\n%v", err)
		}
	}()

	now := time.Now()
	cacheCreationDateString := now.Format(time.RFC3339)
	content := strings.Join([]string{cacheCreationDateString, latestReleaseVersion}, latestCLIReleaseCacheFileContentSeparator)

	logrus.Debugf("Saving content '%v' in cache file...", content)
	if _, err := fmt.Fprintf(cacheFile, "%s", content); err != nil {
		return stacktrace.Propagate(err, "An error occurred saving content '%v' in latest release version cache file", content)
	}
	logrus.Debugf("Content successfully saved in cache file")
	return nil
}

func getLatestCLIReleaseVersionFromCacheFile() (string, error) {

	filepath, err := host_machine_directories.GetLatestCLIReleaseVersionCacheFilepath()
	if err != nil {
		return "", stacktrace.Propagate(err, "An error occurred getting the latest release version cache filepath")
	}
	logrus.Debugf("Cache filepath: '%v'", filepath)

	logrus.Debugf("Getting cache file content...")
	cacheFile, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, latestCLIReleaseCacheFilePermissionsForOpenOrCreateFile)
	if err != nil {
		return "", stacktrace.Propagate(err, "An error occurred opening the '%v' file", filepath)
	}
	defer func() {
		if err := cacheFile.Close(); err != nil {
			logrus.Warnf("We tried to close the latest release CLI version cache file, but doing so threw an error:\n%v", err)
		}
	}()

	bufferSize := 10
	buffer := make([]byte, bufferSize)
	fileContent := ""

	for {
		numberOfBytesRead, err := cacheFile.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", stacktrace.Propagate(err, "An error occurred reading cache file")
		}
		if numberOfBytesRead > zeroBytes {
			newString := string(buffer[:numberOfBytesRead])
			fileContent = fileContent + newString
		}
	}

	if fileContent == "" {
		logrus.Debug("The cache file is empty, skipping getting the latest released version from it")
		return "", nil
	}

	content := strings.Split(fileContent, latestCLIReleaseCacheFileContentSeparator)
	if len(content) != latestCLIReleaseCacheFileContentColumnsAmount {
		return "", stacktrace.NewError("The cache file content '%+v' is not valid", content)
	}
	dateString := content[latestCLIReleaseCacheFileContentDateKey]
	latestReleaseVersion := content[latestCLIReleaseCacheFileContentVersionKey]

	cacheCreationDate, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		return "", stacktrace.Propagate(err, "An error occurred parsing date string '%v' from cache file", dateString)
	}

	cacheExpirationDate := cacheCreationDate.Add(latestCLIReleaseCacheFileExpirationHours * time.Hour)
	logrus.Debugf("Cache Date expiration date '%v'", cacheExpirationDate)

	now := time.Now()
	logrus.Debugf("Now '%v'", now)

	if now.After(cacheExpirationDate) {
		logrus.Debugf("The latest release version cache file content is out-of-date, it was generated in '%v' and it expired in '%v'", cacheCreationDate, cacheExpirationDate)
		return "", nil
	}
	logrus.Debugf("Cache file content '%+v' successfully got", content)

	return latestReleaseVersion, nil
}
