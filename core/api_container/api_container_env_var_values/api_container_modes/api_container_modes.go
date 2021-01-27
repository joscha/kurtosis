/*
 * Copyright (c) 2021 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package api_container_modes

type ApiContainerMode string

const (
	SuiteMetadataSerializingMode ApiContainerMode = "PRINT_SUITE_METADATA"
	TestExecutionMode            ApiContainerMode = "EXECUTE_TEST"
)
