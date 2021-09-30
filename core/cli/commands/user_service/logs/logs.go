/*
 * Copyright (c) 2021 - present Kurtosis Technologies Inc.
 * All Rights Reserved.
 */

package logs

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/kurtosis-tech/container-engine-lib/lib/docker_manager"
	"github.com/kurtosis-tech/kurtosis/commons/enclave_object_labels"
	"github.com/kurtosis-tech/kurtosis/commons/logrus_log_levels"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"sort"
	"strings"
)

const (
	kurtosisLogLevelArg = "kurtosis-log-level"
	nameArg             = "name"
)

var kurtosisLogLevelStr string
var name string

var defaultKurtosisLogLevel = logrus.InfoLevel.String()

var LogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Show user service logs by its name",
	RunE:  run,
}

func init() {
	LogsCmd.Flags().StringVarP(
		&kurtosisLogLevelStr,
		kurtosisLogLevelArg,
		"l",
		defaultKurtosisLogLevel,
		fmt.Sprintf(
			"The log level that Kurtosis itself should log at (%v)",
			strings.Join(logrus_log_levels.GetAcceptableLogLevelStrs(), "|"),
		),
	)

	LogsCmd.Flags().StringVarP(
		&name,
		nameArg,
		"n",
		"",
		"The name of the user service from which the logs are to be displayed",
	)
}

func run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	kurtosisLogLevel, err := logrus.ParseLevel(kurtosisLogLevelStr)
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred parsing Kurtosis loglevel string '%v' to a log level object", kurtosisLogLevelStr)
	}
	logrus.SetLevel(kurtosisLogLevel)

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred creating the Docker client")
	}
	dockerManager := docker_manager.NewDockerManager(
		logrus.StandardLogger(),
		dockerClient,
	)

	labels := getLabelsForListEnclaveUserServices()

	containers, err := dockerManager.GetContainersByLabels(ctx, labels, true)
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred getting containers by labels: '%+v'", labels)
	}

	if containers != nil {
		containersNames := getContainersNames(containers)
		for _, containerNames := range containersNames {
			fmt.Println(containerNames)
		}
	}

	return nil
}

// ====================================================================================================
// 									   Private helper methods
// ====================================================================================================
func getLabelsForListEnclaveUserServices() map[string]string {
	labels := map[string]string{}
	labels[enclave_object_labels.ContainerTypeLabel] = enclave_object_labels.ContainerTypeUserServiceContainer
	labels[enclave_object_labels.EnclaveIDContainerLabel] = name
	return labels
}

func getContainersNames(containers []*docker_manager.Container) []string{
	containersSet := map[string]*docker_manager.Container{}
	for _, container := range containers {
		if container != nil {
			containerId := container.GetId()
			containersSet[containerId] = container
		}
	}

	containersNames := []string{}
	for _, container := range containersSet {
		containerNames := container.GetNames()
		sort.Strings(containerNames)
		containerNamesJoined := strings.Join(containerNames,", ")
		containersNames = append(containersNames, containerNamesJoined)
	}

	sort.Strings(containersNames)

	return containersNames
}
