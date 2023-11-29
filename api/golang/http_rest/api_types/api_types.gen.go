// Package api_types provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version (devel) DO NOT EDIT.
package api_types

import (
	"encoding/json"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
)

// Defines values for ApiContainerStatus.
const (
	ApiContainerStatusNONEXISTENT ApiContainerStatus = "NON_EXISTENT"
	ApiContainerStatusRUNNING     ApiContainerStatus = "RUNNING"
	ApiContainerStatusSTOPPED     ApiContainerStatus = "STOPPED"
)

// Defines values for Connect.
const (
	CONNECT   Connect = "CONNECT"
	NOCONNECT Connect = "NO_CONNECT"
)

// Defines values for ContainerStatus.
const (
	ContainerStatusRUNNING ContainerStatus = "RUNNING"
	ContainerStatusSTOPPED ContainerStatus = "STOPPED"
	ContainerStatusUNKNOWN ContainerStatus = "UNKNOWN"
)

// Defines values for EnclaveContainersStatus.
const (
	EnclaveContainersStatusEMPTY   EnclaveContainersStatus = "EMPTY"
	EnclaveContainersStatusRUNNING EnclaveContainersStatus = "RUNNING"
	EnclaveContainersStatusSTOPPED EnclaveContainersStatus = "STOPPED"
)

// Defines values for EnclaveMode.
const (
	PRODUCTION EnclaveMode = "PRODUCTION"
	TEST       EnclaveMode = "TEST"
)

// Defines values for EnclaveTargetStatus.
const (
	STOP EnclaveTargetStatus = "STOP"
)

// Defines values for HttpMethodAvailability.
const (
	GET  HttpMethodAvailability = "GET"
	POST HttpMethodAvailability = "POST"
)

// Defines values for ImageDownloadMode.
const (
	ImageDownloadModeALWAYS  ImageDownloadMode = "ALWAYS"
	ImageDownloadModeMISSING ImageDownloadMode = "MISSING"
)

// Defines values for KurtosisFeatureFlag.
const (
	NOINSTRUCTIONSCACHING KurtosisFeatureFlag = "NO_INSTRUCTIONS_CACHING"
)

// Defines values for LogLineOperator.
const (
	DOESCONTAINMATCHREGEX    LogLineOperator = "DOES_CONTAIN_MATCH_REGEX"
	DOESCONTAINTEXT          LogLineOperator = "DOES_CONTAIN_TEXT"
	DOESNOTCONTAINMATCHREGEX LogLineOperator = "DOES_NOT_CONTAIN_MATCH_REGEX"
	DOESNOTCONTAINTEXT       LogLineOperator = "DOES_NOT_CONTAIN_TEXT"
)

// Defines values for ResponseType.
const (
	ERROR   ResponseType = "ERROR"
	INFO    ResponseType = "INFO"
	WARNING ResponseType = "WARNING"
)

// Defines values for RestartPolicy.
const (
	RestartPolicyALWAYS RestartPolicy = "ALWAYS"
	RestartPolicyNEVER  RestartPolicy = "NEVER"
)

// Defines values for ServiceStatus.
const (
	ServiceStatusRUNNING ServiceStatus = "RUNNING"
	ServiceStatusSTOPPED ServiceStatus = "STOPPED"
	ServiceStatusUNKNOWN ServiceStatus = "UNKNOWN"
)

// Defines values for TransportProtocol.
const (
	SCTP TransportProtocol = "SCTP"
	TCP  TransportProtocol = "TCP"
	UDP  TransportProtocol = "UDP"
)

// ApiContainerStatus defines model for ApiContainerStatus.
type ApiContainerStatus string

// AsyncStarlarkExecutionLogs Use it to asynchronously retrieve the execution logs via Websockets or http streaming
type AsyncStarlarkExecutionLogs struct {
	// AsyncStarlarkExecutionLogs Execution UUID to asynchronously retrieve the execution logs
	AsyncStarlarkExecutionLogs struct {
		StarlarkExecutionUuid string `json:"starlark_execution_uuid"`
	} `json:"async_starlark_execution_logs"`
}

// Connect 0 - CONNECT // Best effort port forwarding
// 1 - NO_CONNECT // Port forwarding disabled
type Connect string

// Container defines model for Container.
type Container struct {
	CmdArgs        []string          `json:"cmd_args"`
	EntrypointArgs []string          `json:"entrypoint_args"`
	EnvVars        map[string]string `json:"env_vars"`
	ImageName      string            `json:"image_name"`

	// Status 0 - STOPPED
	// 1 - RUNNING
	// 2 - UNKNOWN
	Status ContainerStatus `json:"status"`
}

// ContainerStatus 0 - STOPPED
// 1 - RUNNING
// 2 - UNKNOWN
type ContainerStatus string

// CreateEnclave defines model for CreateEnclave.
type CreateEnclave struct {
	ApiContainerLogLevel   *string      `json:"api_container_log_level,omitempty"`
	ApiContainerVersionTag *string      `json:"api_container_version_tag,omitempty"`
	EnclaveName            *string      `json:"enclave_name,omitempty"`
	Mode                   *EnclaveMode `json:"mode,omitempty"`
}

// DeletionSummary defines model for DeletionSummary.
type DeletionSummary struct {
	RemovedEnclaveNameAndUuids *[]EnclaveNameAndUuid `json:"removed_enclave_name_and_uuids,omitempty"`
}

// EnclaveAPIContainerHostMachineInfo defines model for EnclaveAPIContainerHostMachineInfo.
type EnclaveAPIContainerHostMachineInfo struct {
	GrpcPortOnHostMachine int    `json:"grpc_port_on_host_machine"`
	IpOnHostMachine       string `json:"ip_on_host_machine"`
}

// EnclaveAPIContainerInfo defines model for EnclaveAPIContainerInfo.
type EnclaveAPIContainerInfo struct {
	BridgeIpAddress       string `json:"bridge_ip_address"`
	ContainerId           string `json:"container_id"`
	GrpcPortInsideEnclave int    `json:"grpc_port_inside_enclave"`
	IpInsideEnclave       string `json:"ip_inside_enclave"`
}

// EnclaveContainersStatus defines model for EnclaveContainersStatus.
type EnclaveContainersStatus string

// EnclaveIdentifiers defines model for EnclaveIdentifiers.
type EnclaveIdentifiers struct {
	EnclaveUuid   string `json:"enclave_uuid"`
	Name          string `json:"name"`
	ShortenedUuid string `json:"shortened_uuid"`
}

// EnclaveInfo defines model for EnclaveInfo.
type EnclaveInfo struct {
	ApiContainerHostMachineInfo *EnclaveAPIContainerHostMachineInfo `json:"api_container_host_machine_info,omitempty"`
	ApiContainerInfo            *EnclaveAPIContainerInfo            `json:"api_container_info,omitempty"`
	ApiContainerStatus          ApiContainerStatus                  `json:"api_container_status"`
	ContainersStatus            EnclaveContainersStatus             `json:"containers_status"`
	CreationTime                Timestamp                           `json:"creation_time"`
	EnclaveUuid                 string                              `json:"enclave_uuid"`
	Mode                        EnclaveMode                         `json:"mode"`
	Name                        string                              `json:"name"`
	ShortenedUuid               string                              `json:"shortened_uuid"`
}

// EnclaveMode defines model for EnclaveMode.
type EnclaveMode string

// EnclaveNameAndUuid defines model for EnclaveNameAndUuid.
type EnclaveNameAndUuid struct {
	Name string `json:"name"`
	Uuid string `json:"uuid"`
}

// EnclaveTargetStatus defines model for EnclaveTargetStatus.
type EnclaveTargetStatus string

// EngineInfo defines model for EngineInfo.
type EngineInfo struct {
	EngineVersion string `json:"engine_version"`
}

// ExecCommand Exec Command
type ExecCommand struct {
	CommandArgs []string `json:"command_args"`
}

// ExecCommandResult defines model for ExecCommandResult.
type ExecCommandResult struct {
	ExitCode int32 `json:"exit_code"`

	// LogOutput Assumes UTF-8 encoding
	LogOutput string `json:"log_output"`
}

// FileArtifactDescription defines model for FileArtifactDescription.
type FileArtifactDescription struct {
	// Path Path relative to the file artifact
	Path string `json:"path"`

	// Size Size of the file, in bytes
	Size int64 `json:"size"`

	// TextPreview A bit of text content, if the file allows (similar to UNIX's 'head')
	TextPreview *string `json:"text_preview,omitempty"`
}

// FileArtifactReference Files Artifact identifier
type FileArtifactReference struct {
	// Name UUID of the files artifact, for use when referencing it in the future
	Name string `json:"name"`

	// Uuid UUID of the files artifact, for use when referencing it in the future
	Uuid string `json:"uuid"`
}

// HttpMethodAvailability defines model for HttpMethodAvailability.
type HttpMethodAvailability string

// ImageDownloadMode 0 - ALWAYS
// 1 - MISSING
type ImageDownloadMode string

// KurtosisFeatureFlag 0 - NO_INSTRUCTIONS_CACHING
type KurtosisFeatureFlag string

// LogLine defines model for LogLine.
type LogLine struct {
	Line      []string  `json:"line"`
	Timestamp Timestamp `json:"timestamp"`
}

// LogLineFilter defines model for LogLineFilter.
type LogLineFilter struct {
	Operator    LogLineOperator `json:"operator"`
	TextPattern string          `json:"text_pattern"`
}

// LogLineOperator defines model for LogLineOperator.
type LogLineOperator string

// Port Shared Objects (Used By Multiple Endpoints)
type Port struct {
	ApplicationProtocol *string `json:"application_protocol,omitempty"`
	Number              int32   `json:"number"`

	// TransportProtocol 0 - TCP
	// 1 - SCTP
	// 2 - UDP
	TransportProtocol TransportProtocol `json:"transport_protocol"`

	// WaitTimeout The wait timeout duration in string
	WaitTimeout *string `json:"wait_timeout,omitempty"`
}

// ResponseInfo defines model for ResponseInfo.
type ResponseInfo struct {
	Code    uint32       `json:"code"`
	Message string       `json:"message"`
	Type    ResponseType `json:"type"`
}

// ResponseType defines model for ResponseType.
type ResponseType string

// RestartPolicy 0 - NEVER
// 1 - ALWAYS
type RestartPolicy string

// RunStarlarkPackage defines model for RunStarlarkPackage.
type RunStarlarkPackage struct {
	// ClonePackage Whether the package should be cloned or not.
	// If false, then the package will be pulled from the APIC local package store. If it's a local package then is must
	// have been uploaded using UploadStarlarkPackage prior to calling RunStarlarkPackage.
	// If true, then the package will be cloned from GitHub before execution starts
	ClonePackage *bool `json:"clone_package,omitempty"`

	// CloudInstanceId Defaults to empty
	CloudInstanceId *string `json:"cloud_instance_id,omitempty"`

	// CloudUserId Defaults to empty
	CloudUserId *string `json:"cloud_user_id,omitempty"`

	// DryRun Defaults to false
	DryRun               *bool                  `json:"dry_run,omitempty"`
	ExperimentalFeatures *[]KurtosisFeatureFlag `json:"experimental_features,omitempty"`

	// ImageDownloadMode 0 - ALWAYS
	// 1 - MISSING
	ImageDownloadMode *ImageDownloadMode `json:"image_download_mode,omitempty"`

	// Local the payload of the local module
	Local *[]byte `json:"local,omitempty"`

	// MainFunctionName The name of the main function, the default value is "run"
	MainFunctionName *string `json:"main_function_name,omitempty"`

	// Parallelism Defaults to 4
	Parallelism *int32 `json:"parallelism,omitempty"`

	// Params Parameters data for the Starlark package main function
	Params *map[string]interface{} `json:"params,omitempty"`

	// RelativePathToMainFile The relative main file filepath, the default value is the "main.star" file in the root of a package
	RelativePathToMainFile *string `json:"relative_path_to_main_file,omitempty"`

	// Remote just a flag to indicate the module must be cloned inside the API
	Remote *bool `json:"remote,omitempty"`
}

// RunStarlarkScript defines model for RunStarlarkScript.
type RunStarlarkScript struct {
	// CloudInstanceId Defaults to empty
	CloudInstanceId *string `json:"cloud_instance_id,omitempty"`

	// CloudUserId Defaults to empty
	CloudUserId *string `json:"cloud_user_id,omitempty"`

	// DryRun Defaults to false
	DryRun               *bool                  `json:"dry_run,omitempty"`
	ExperimentalFeatures *[]KurtosisFeatureFlag `json:"experimental_features,omitempty"`

	// ImageDownloadMode 0 - ALWAYS
	// 1 - MISSING
	ImageDownloadMode *ImageDownloadMode `json:"image_download_mode,omitempty"`

	// MainFunctionName The name of the main function, the default value is "run"
	MainFunctionName *string `json:"main_function_name,omitempty"`

	// Parallelism Defaults to 4
	Parallelism *int32 `json:"parallelism,omitempty"`

	// Params Parameters data for the Starlark package main function
	Params           *map[string]interface{} `json:"params,omitempty"`
	SerializedScript string                  `json:"serialized_script"`
}

// ServiceIdentifiers An service identifier is a collection of uuid, name and shortened uuid
type ServiceIdentifiers struct {
	// Name Name of the service
	Name string `json:"name"`

	// ServiceUuid UUID of the service
	ServiceUuid string `json:"service_uuid"`

	// ShortenedUuid The shortened uuid of the service
	ShortenedUuid string `json:"shortened_uuid"`
}

// ServiceInfo defines model for ServiceInfo.
type ServiceInfo struct {
	Container Container `json:"container"`

	// Name Name of the service
	Name string `json:"name"`

	// PrivateIpAddr The IP address of the service inside the enclave
	PrivateIpAddr string          `json:"private_ip_addr"`
	PrivatePorts  map[string]Port `json:"private_ports"`

	// PublicIpAddr Public IP address *outside* the enclave where the service is reachable
	// NOTE: Will be empty if the service isn't running, the service didn't define any ports, or the backend doesn't support reporting public service info
	PublicIpAddr *string          `json:"public_ip_addr,omitempty"`
	PublicPorts  *map[string]Port `json:"public_ports,omitempty"`

	// ServiceStatus 0 - STOPPED
	// 1 - RUNNING
	// 2 - UNKNOWN
	ServiceStatus ServiceStatus `json:"service_status"`

	// ServiceUuid UUID of the service
	ServiceUuid string `json:"service_uuid"`

	// ShortenedUuid Shortened uuid of the service
	ShortenedUuid string `json:"shortened_uuid"`
}

// ServiceLogs defines model for ServiceLogs.
type ServiceLogs struct {
	NotFoundServiceUuidSet   *[]string           `json:"not_found_service_uuid_set,omitempty"`
	ServiceLogsByServiceUuid *map[string]LogLine `json:"service_logs_by_service_uuid,omitempty"`
}

// ServiceStatus 0 - STOPPED
// 1 - RUNNING
// 2 - UNKNOWN
type ServiceStatus string

// StarlarkDescription defines model for StarlarkDescription.
type StarlarkDescription struct {
	ExperimentalFeatures   []KurtosisFeatureFlag `json:"experimental_features"`
	MainFunctionName       string                `json:"main_function_name"`
	PackageId              string                `json:"package_id"`
	Parallelism            int32                 `json:"parallelism"`
	RelativePathToMainFile string                `json:"relative_path_to_main_file"`

	// RestartPolicy 0 - NEVER
	// 1 - ALWAYS
	RestartPolicy    RestartPolicy `json:"restart_policy"`
	SerializedParams string        `json:"serialized_params"`
	SerializedScript string        `json:"serialized_script"`
}

// StarlarkError defines model for StarlarkError.
type StarlarkError struct {
	Error StarlarkError_Error `json:"error"`
}

// StarlarkError_Error defines model for StarlarkError.Error.
type StarlarkError_Error struct {
	union json.RawMessage
}

// StarlarkExecutionError defines model for StarlarkExecutionError.
type StarlarkExecutionError struct {
	ExecutionError struct {
		ErrorMessage string `json:"error_message"`
	} `json:"execution_error"`
}

// StarlarkInfo defines model for StarlarkInfo.
type StarlarkInfo struct {
	Info struct {
		Instruction struct {
			InfoMessage string `json:"info_message"`
		} `json:"instruction"`
	} `json:"info"`
}

// StarlarkInstruction defines model for StarlarkInstruction.
type StarlarkInstruction struct {
	Arguments             []StarlarkInstructionArgument `json:"arguments"`
	ExecutableInstruction string                        `json:"executable_instruction"`
	InstructionName       string                        `json:"instruction_name"`
	IsSkipped             bool                          `json:"is_skipped"`
	Position              StarlarkInstructionPosition   `json:"position"`
}

// StarlarkInstructionArgument defines model for StarlarkInstructionArgument.
type StarlarkInstructionArgument struct {
	ArgName            *string `json:"arg_name,omitempty"`
	IsRepresentative   bool    `json:"is_representative"`
	SerializedArgValue string  `json:"serialized_arg_value"`
}

// StarlarkInstructionPosition defines model for StarlarkInstructionPosition.
type StarlarkInstructionPosition struct {
	Column   int32  `json:"column"`
	Filename string `json:"filename"`
	Line     int32  `json:"line"`
}

// StarlarkInstructionResult defines model for StarlarkInstructionResult.
type StarlarkInstructionResult struct {
	InstructionResult struct {
		SerializedInstructionResult string `json:"serialized_instruction_result"`
	} `json:"instruction_result"`
}

// StarlarkInterpretationError defines model for StarlarkInterpretationError.
type StarlarkInterpretationError struct {
	InterpretationError struct {
		ErrorMessage string `json:"error_message"`
	} `json:"interpretation_error"`
}

// StarlarkRunFinishedEvent defines model for StarlarkRunFinishedEvent.
type StarlarkRunFinishedEvent struct {
	RunFinishedEvent struct {
		IsRunSuccessful  bool    `json:"is_run_successful"`
		SerializedOutput *string `json:"serialized_output,omitempty"`
	} `json:"run_finished_event"`
}

// StarlarkRunLogs Starlark Execution Logs
type StarlarkRunLogs = []StarlarkRunResponseLine

// StarlarkRunProgress defines model for StarlarkRunProgress.
type StarlarkRunProgress struct {
	ProgressInfo struct {
		CurrentStepInfo   []string `json:"current_step_info"`
		CurrentStepNumber int32    `json:"current_step_number"`
		TotalSteps        int32    `json:"total_steps"`
	} `json:"progress_info"`
}

// StarlarkRunResponse defines model for StarlarkRunResponse.
type StarlarkRunResponse struct {
	union json.RawMessage
}

// StarlarkRunResponseLine Starlark Execution Response
type StarlarkRunResponseLine struct {
	union json.RawMessage
}

// StarlarkValidationError defines model for StarlarkValidationError.
type StarlarkValidationError struct {
	ValidationError struct {
		ErrorMessage string `json:"error_message"`
	} `json:"validation_error"`
}

// StarlarkWarning defines model for StarlarkWarning.
type StarlarkWarning struct {
	Warning struct {
		WarningMessage string `json:"warning_message"`
	} `json:"warning"`
}

// StoreFilesArtifactFromService defines model for StoreFilesArtifactFromService.
type StoreFilesArtifactFromService struct {
	// Name The name of the files artifact
	Name string `json:"name"`

	// SourcePath The absolute source path where the source files will be copied from
	SourcePath string `json:"source_path"`
}

// StoreWebFilesArtifact Store Web Files Artifact
type StoreWebFilesArtifact struct {
	// Name The name of the files artifact
	Name string `json:"name"`

	// Url URL to download the artifact from
	Url string `json:"url"`
}

// Timestamp defines model for Timestamp.
type Timestamp = time.Time

// TransportProtocol 0 - TCP
// 1 - SCTP
// 2 - UDP
type TransportProtocol string

// ArtifactIdentifier defines model for artifact_identifier.
type ArtifactIdentifier = string

// ConjunctiveFilters defines model for conjunctive_filters.
type ConjunctiveFilters = []LogLineFilter

// EnclaveIdentifier defines model for enclave_identifier.
type EnclaveIdentifier = string

// ExpectedResponse defines model for expected_response.
type ExpectedResponse = string

// FollowLogs defines model for follow_logs.
type FollowLogs = bool

// HttpMethod defines model for http_method.
type HttpMethod = HttpMethodAvailability

// InitialDelayMilliseconds defines model for initial_delay_milliseconds.
type InitialDelayMilliseconds = int32

// NumLogLines defines model for num_log_lines.
type NumLogLines = int

// PackageId defines model for package_id.
type PackageId = string

// Path defines model for path.
type Path = string

// PortNumber defines model for port_number.
type PortNumber = int32

// RemoveAll defines model for remove_all.
type RemoveAll = bool

// RequestBody defines model for request_body.
type RequestBody = string

// Retries defines model for retries.
type Retries = int32

// RetriesDelayMilliseconds defines model for retries_delay_milliseconds.
type RetriesDelayMilliseconds = int32

// RetrieveLogsAsync defines model for retrieve_logs_async.
type RetrieveLogsAsync = bool

// ReturnAllLogs defines model for return_all_logs.
type ReturnAllLogs = bool

// ServiceIdentifier defines model for service_identifier.
type ServiceIdentifier = string

// ServiceUuidSet defines model for service_uuid_set.
type ServiceUuidSet = []string

// StarlarkExecutionUuid defines model for starlark_execution_uuid.
type StarlarkExecutionUuid = string

// NotOk defines model for NotOk.
type NotOk = ResponseInfo

// DeleteEnclavesParams defines parameters for DeleteEnclaves.
type DeleteEnclavesParams struct {
	// RemoveAll If true, remove all enclaves. Default is false
	RemoveAll *RemoveAll `form:"remove_all,omitempty" json:"remove_all,omitempty"`
}

// PostEnclavesEnclaveIdentifierArtifactsLocalFileMultipartBody defines parameters for PostEnclavesEnclaveIdentifierArtifactsLocalFile.
type PostEnclavesEnclaveIdentifierArtifactsLocalFileMultipartBody = openapi_types.File

// GetEnclavesEnclaveIdentifierLogsParams defines parameters for GetEnclavesEnclaveIdentifierLogs.
type GetEnclavesEnclaveIdentifierLogsParams struct {
	ServiceUuidSet     ServiceUuidSet      `form:"service_uuid_set" json:"service_uuid_set"`
	FollowLogs         *FollowLogs         `form:"follow_logs,omitempty" json:"follow_logs,omitempty"`
	ConjunctiveFilters *ConjunctiveFilters `form:"conjunctive_filters,omitempty" json:"conjunctive_filters,omitempty"`
	ReturnAllLogs      *ReturnAllLogs      `form:"return_all_logs,omitempty" json:"return_all_logs,omitempty"`
	NumLogLines        *NumLogLines        `form:"num_log_lines,omitempty" json:"num_log_lines,omitempty"`
}

// GetEnclavesEnclaveIdentifierServicesParams defines parameters for GetEnclavesEnclaveIdentifierServices.
type GetEnclavesEnclaveIdentifierServicesParams struct {
	// Services Select services to get information
	Services *[]string `form:"services,omitempty" json:"services,omitempty"`
}

// GetEnclavesEnclaveIdentifierServicesServiceIdentifierEndpointsPortNumberAvailabilityParams defines parameters for GetEnclavesEnclaveIdentifierServicesServiceIdentifierEndpointsPortNumberAvailability.
type GetEnclavesEnclaveIdentifierServicesServiceIdentifierEndpointsPortNumberAvailabilityParams struct {
	// HttpMethod The HTTP method used to check availability. Default is GET.
	HttpMethod *HttpMethod `form:"http_method,omitempty" json:"http_method,omitempty"`

	// Path The path of the service to check. It mustn't start with the first slash. For instance `service/health`
	Path *Path `form:"path,omitempty" json:"path,omitempty"`

	// InitialDelayMilliseconds The number of milliseconds to wait until executing the first HTTP call
	InitialDelayMilliseconds *InitialDelayMilliseconds `form:"initial_delay_milliseconds,omitempty" json:"initial_delay_milliseconds,omitempty"`

	// Retries Max number of HTTP call attempts that this will execute until giving up and returning an error
	Retries *Retries `form:"retries,omitempty" json:"retries,omitempty"`

	// RetriesDelayMilliseconds Number of milliseconds to wait between retries
	RetriesDelayMilliseconds *RetriesDelayMilliseconds `form:"retries_delay_milliseconds,omitempty" json:"retries_delay_milliseconds,omitempty"`

	// ExpectedResponse If the endpoint returns this value, the service will be marked as available (e.g. Hello World).
	ExpectedResponse *ExpectedResponse `form:"expected_response,omitempty" json:"expected_response,omitempty"`

	// RequestBody If the http_method is set to POST, this value will be send as the body of the availability request.
	RequestBody *RequestBody `form:"request_body,omitempty" json:"request_body,omitempty"`
}

// GetEnclavesEnclaveIdentifierServicesServiceIdentifierLogsParams defines parameters for GetEnclavesEnclaveIdentifierServicesServiceIdentifierLogs.
type GetEnclavesEnclaveIdentifierServicesServiceIdentifierLogsParams struct {
	FollowLogs         *FollowLogs         `form:"follow_logs,omitempty" json:"follow_logs,omitempty"`
	ConjunctiveFilters *ConjunctiveFilters `form:"conjunctive_filters,omitempty" json:"conjunctive_filters,omitempty"`
	ReturnAllLogs      *ReturnAllLogs      `form:"return_all_logs,omitempty" json:"return_all_logs,omitempty"`
	NumLogLines        *NumLogLines        `form:"num_log_lines,omitempty" json:"num_log_lines,omitempty"`
}

// PostEnclavesEnclaveIdentifierStarlarkPackagesMultipartBody defines parameters for PostEnclavesEnclaveIdentifierStarlarkPackages.
type PostEnclavesEnclaveIdentifierStarlarkPackagesMultipartBody = openapi_types.File

// PostEnclavesEnclaveIdentifierStarlarkPackagesPackageIdParams defines parameters for PostEnclavesEnclaveIdentifierStarlarkPackagesPackageId.
type PostEnclavesEnclaveIdentifierStarlarkPackagesPackageIdParams struct {
	// RetrieveLogsAsync If false, block http response until all logs are available. Default is true
	RetrieveLogsAsync *RetrieveLogsAsync `form:"retrieve_logs_async,omitempty" json:"retrieve_logs_async,omitempty"`
}

// PostEnclavesEnclaveIdentifierStarlarkScriptsParams defines parameters for PostEnclavesEnclaveIdentifierStarlarkScripts.
type PostEnclavesEnclaveIdentifierStarlarkScriptsParams struct {
	// RetrieveLogsAsync If false, block http response until all logs are available. Default is true
	RetrieveLogsAsync *RetrieveLogsAsync `form:"retrieve_logs_async,omitempty" json:"retrieve_logs_async,omitempty"`
}

// PostEnclavesJSONRequestBody defines body for PostEnclaves for application/json ContentType.
type PostEnclavesJSONRequestBody = CreateEnclave

// PostEnclavesEnclaveIdentifierArtifactsLocalFileMultipartRequestBody defines body for PostEnclavesEnclaveIdentifierArtifactsLocalFile for multipart/form-data ContentType.
type PostEnclavesEnclaveIdentifierArtifactsLocalFileMultipartRequestBody = PostEnclavesEnclaveIdentifierArtifactsLocalFileMultipartBody

// PostEnclavesEnclaveIdentifierArtifactsRemoteFileJSONRequestBody defines body for PostEnclavesEnclaveIdentifierArtifactsRemoteFile for application/json ContentType.
type PostEnclavesEnclaveIdentifierArtifactsRemoteFileJSONRequestBody = StoreWebFilesArtifact

// PostEnclavesEnclaveIdentifierArtifactsServicesServiceIdentifierJSONRequestBody defines body for PostEnclavesEnclaveIdentifierArtifactsServicesServiceIdentifier for application/json ContentType.
type PostEnclavesEnclaveIdentifierArtifactsServicesServiceIdentifierJSONRequestBody = StoreFilesArtifactFromService

// PostEnclavesEnclaveIdentifierServicesConnectionJSONRequestBody defines body for PostEnclavesEnclaveIdentifierServicesConnection for application/json ContentType.
type PostEnclavesEnclaveIdentifierServicesConnectionJSONRequestBody = Connect

// PostEnclavesEnclaveIdentifierServicesServiceIdentifierCommandJSONRequestBody defines body for PostEnclavesEnclaveIdentifierServicesServiceIdentifierCommand for application/json ContentType.
type PostEnclavesEnclaveIdentifierServicesServiceIdentifierCommandJSONRequestBody = ExecCommand

// PostEnclavesEnclaveIdentifierStarlarkPackagesMultipartRequestBody defines body for PostEnclavesEnclaveIdentifierStarlarkPackages for multipart/form-data ContentType.
type PostEnclavesEnclaveIdentifierStarlarkPackagesMultipartRequestBody = PostEnclavesEnclaveIdentifierStarlarkPackagesMultipartBody

// PostEnclavesEnclaveIdentifierStarlarkPackagesPackageIdJSONRequestBody defines body for PostEnclavesEnclaveIdentifierStarlarkPackagesPackageId for application/json ContentType.
type PostEnclavesEnclaveIdentifierStarlarkPackagesPackageIdJSONRequestBody = RunStarlarkPackage

// PostEnclavesEnclaveIdentifierStarlarkScriptsJSONRequestBody defines body for PostEnclavesEnclaveIdentifierStarlarkScripts for application/json ContentType.
type PostEnclavesEnclaveIdentifierStarlarkScriptsJSONRequestBody = RunStarlarkScript

// PostEnclavesEnclaveIdentifierStatusJSONRequestBody defines body for PostEnclavesEnclaveIdentifierStatus for application/json ContentType.
type PostEnclavesEnclaveIdentifierStatusJSONRequestBody = EnclaveTargetStatus

// AsStarlarkInterpretationError returns the union data inside the StarlarkError_Error as a StarlarkInterpretationError
func (t StarlarkError_Error) AsStarlarkInterpretationError() (StarlarkInterpretationError, error) {
	var body StarlarkInterpretationError
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromStarlarkInterpretationError overwrites any union data inside the StarlarkError_Error as the provided StarlarkInterpretationError
func (t *StarlarkError_Error) FromStarlarkInterpretationError(v StarlarkInterpretationError) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeStarlarkInterpretationError performs a merge with any union data inside the StarlarkError_Error, using the provided StarlarkInterpretationError
func (t *StarlarkError_Error) MergeStarlarkInterpretationError(v StarlarkInterpretationError) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(b, t.union)
	t.union = merged
	return err
}

// AsStarlarkValidationError returns the union data inside the StarlarkError_Error as a StarlarkValidationError
func (t StarlarkError_Error) AsStarlarkValidationError() (StarlarkValidationError, error) {
	var body StarlarkValidationError
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromStarlarkValidationError overwrites any union data inside the StarlarkError_Error as the provided StarlarkValidationError
func (t *StarlarkError_Error) FromStarlarkValidationError(v StarlarkValidationError) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeStarlarkValidationError performs a merge with any union data inside the StarlarkError_Error, using the provided StarlarkValidationError
func (t *StarlarkError_Error) MergeStarlarkValidationError(v StarlarkValidationError) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(b, t.union)
	t.union = merged
	return err
}

// AsStarlarkExecutionError returns the union data inside the StarlarkError_Error as a StarlarkExecutionError
func (t StarlarkError_Error) AsStarlarkExecutionError() (StarlarkExecutionError, error) {
	var body StarlarkExecutionError
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromStarlarkExecutionError overwrites any union data inside the StarlarkError_Error as the provided StarlarkExecutionError
func (t *StarlarkError_Error) FromStarlarkExecutionError(v StarlarkExecutionError) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeStarlarkExecutionError performs a merge with any union data inside the StarlarkError_Error, using the provided StarlarkExecutionError
func (t *StarlarkError_Error) MergeStarlarkExecutionError(v StarlarkExecutionError) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(b, t.union)
	t.union = merged
	return err
}

func (t StarlarkError_Error) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *StarlarkError_Error) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}

// AsAsyncStarlarkExecutionLogs returns the union data inside the StarlarkRunResponse as a AsyncStarlarkExecutionLogs
func (t StarlarkRunResponse) AsAsyncStarlarkExecutionLogs() (AsyncStarlarkExecutionLogs, error) {
	var body AsyncStarlarkExecutionLogs
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromAsyncStarlarkExecutionLogs overwrites any union data inside the StarlarkRunResponse as the provided AsyncStarlarkExecutionLogs
func (t *StarlarkRunResponse) FromAsyncStarlarkExecutionLogs(v AsyncStarlarkExecutionLogs) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeAsyncStarlarkExecutionLogs performs a merge with any union data inside the StarlarkRunResponse, using the provided AsyncStarlarkExecutionLogs
func (t *StarlarkRunResponse) MergeAsyncStarlarkExecutionLogs(v AsyncStarlarkExecutionLogs) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(b, t.union)
	t.union = merged
	return err
}

// AsStarlarkRunLogs returns the union data inside the StarlarkRunResponse as a StarlarkRunLogs
func (t StarlarkRunResponse) AsStarlarkRunLogs() (StarlarkRunLogs, error) {
	var body StarlarkRunLogs
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromStarlarkRunLogs overwrites any union data inside the StarlarkRunResponse as the provided StarlarkRunLogs
func (t *StarlarkRunResponse) FromStarlarkRunLogs(v StarlarkRunLogs) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeStarlarkRunLogs performs a merge with any union data inside the StarlarkRunResponse, using the provided StarlarkRunLogs
func (t *StarlarkRunResponse) MergeStarlarkRunLogs(v StarlarkRunLogs) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(b, t.union)
	t.union = merged
	return err
}

func (t StarlarkRunResponse) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *StarlarkRunResponse) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}

// AsStarlarkInstruction returns the union data inside the StarlarkRunResponseLine as a StarlarkInstruction
func (t StarlarkRunResponseLine) AsStarlarkInstruction() (StarlarkInstruction, error) {
	var body StarlarkInstruction
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromStarlarkInstruction overwrites any union data inside the StarlarkRunResponseLine as the provided StarlarkInstruction
func (t *StarlarkRunResponseLine) FromStarlarkInstruction(v StarlarkInstruction) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeStarlarkInstruction performs a merge with any union data inside the StarlarkRunResponseLine, using the provided StarlarkInstruction
func (t *StarlarkRunResponseLine) MergeStarlarkInstruction(v StarlarkInstruction) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(b, t.union)
	t.union = merged
	return err
}

// AsStarlarkError returns the union data inside the StarlarkRunResponseLine as a StarlarkError
func (t StarlarkRunResponseLine) AsStarlarkError() (StarlarkError, error) {
	var body StarlarkError
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromStarlarkError overwrites any union data inside the StarlarkRunResponseLine as the provided StarlarkError
func (t *StarlarkRunResponseLine) FromStarlarkError(v StarlarkError) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeStarlarkError performs a merge with any union data inside the StarlarkRunResponseLine, using the provided StarlarkError
func (t *StarlarkRunResponseLine) MergeStarlarkError(v StarlarkError) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(b, t.union)
	t.union = merged
	return err
}

// AsStarlarkRunProgress returns the union data inside the StarlarkRunResponseLine as a StarlarkRunProgress
func (t StarlarkRunResponseLine) AsStarlarkRunProgress() (StarlarkRunProgress, error) {
	var body StarlarkRunProgress
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromStarlarkRunProgress overwrites any union data inside the StarlarkRunResponseLine as the provided StarlarkRunProgress
func (t *StarlarkRunResponseLine) FromStarlarkRunProgress(v StarlarkRunProgress) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeStarlarkRunProgress performs a merge with any union data inside the StarlarkRunResponseLine, using the provided StarlarkRunProgress
func (t *StarlarkRunResponseLine) MergeStarlarkRunProgress(v StarlarkRunProgress) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(b, t.union)
	t.union = merged
	return err
}

// AsStarlarkInstructionResult returns the union data inside the StarlarkRunResponseLine as a StarlarkInstructionResult
func (t StarlarkRunResponseLine) AsStarlarkInstructionResult() (StarlarkInstructionResult, error) {
	var body StarlarkInstructionResult
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromStarlarkInstructionResult overwrites any union data inside the StarlarkRunResponseLine as the provided StarlarkInstructionResult
func (t *StarlarkRunResponseLine) FromStarlarkInstructionResult(v StarlarkInstructionResult) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeStarlarkInstructionResult performs a merge with any union data inside the StarlarkRunResponseLine, using the provided StarlarkInstructionResult
func (t *StarlarkRunResponseLine) MergeStarlarkInstructionResult(v StarlarkInstructionResult) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(b, t.union)
	t.union = merged
	return err
}

// AsStarlarkRunFinishedEvent returns the union data inside the StarlarkRunResponseLine as a StarlarkRunFinishedEvent
func (t StarlarkRunResponseLine) AsStarlarkRunFinishedEvent() (StarlarkRunFinishedEvent, error) {
	var body StarlarkRunFinishedEvent
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromStarlarkRunFinishedEvent overwrites any union data inside the StarlarkRunResponseLine as the provided StarlarkRunFinishedEvent
func (t *StarlarkRunResponseLine) FromStarlarkRunFinishedEvent(v StarlarkRunFinishedEvent) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeStarlarkRunFinishedEvent performs a merge with any union data inside the StarlarkRunResponseLine, using the provided StarlarkRunFinishedEvent
func (t *StarlarkRunResponseLine) MergeStarlarkRunFinishedEvent(v StarlarkRunFinishedEvent) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(b, t.union)
	t.union = merged
	return err
}

// AsStarlarkWarning returns the union data inside the StarlarkRunResponseLine as a StarlarkWarning
func (t StarlarkRunResponseLine) AsStarlarkWarning() (StarlarkWarning, error) {
	var body StarlarkWarning
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromStarlarkWarning overwrites any union data inside the StarlarkRunResponseLine as the provided StarlarkWarning
func (t *StarlarkRunResponseLine) FromStarlarkWarning(v StarlarkWarning) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeStarlarkWarning performs a merge with any union data inside the StarlarkRunResponseLine, using the provided StarlarkWarning
func (t *StarlarkRunResponseLine) MergeStarlarkWarning(v StarlarkWarning) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(b, t.union)
	t.union = merged
	return err
}

// AsStarlarkInfo returns the union data inside the StarlarkRunResponseLine as a StarlarkInfo
func (t StarlarkRunResponseLine) AsStarlarkInfo() (StarlarkInfo, error) {
	var body StarlarkInfo
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromStarlarkInfo overwrites any union data inside the StarlarkRunResponseLine as the provided StarlarkInfo
func (t *StarlarkRunResponseLine) FromStarlarkInfo(v StarlarkInfo) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeStarlarkInfo performs a merge with any union data inside the StarlarkRunResponseLine, using the provided StarlarkInfo
func (t *StarlarkRunResponseLine) MergeStarlarkInfo(v StarlarkInfo) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(b, t.union)
	t.union = merged
	return err
}

func (t StarlarkRunResponseLine) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *StarlarkRunResponseLine) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}
