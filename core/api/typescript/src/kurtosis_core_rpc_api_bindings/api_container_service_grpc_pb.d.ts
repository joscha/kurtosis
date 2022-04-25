// GENERATED CODE -- DO NOT EDIT!

// package: api_container_api
// file: api_container_service.proto

import * as api_container_service_pb from "./api_container_service_pb";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";
import * as grpc from "@grpc/grpc-js";

interface IApiContainerServiceService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
  loadModule: grpc.MethodDefinition<api_container_service_pb.LoadModuleArgs, api_container_service_pb.LoadModuleResponse>;
  unloadModule: grpc.MethodDefinition<api_container_service_pb.UnloadModuleArgs, google_protobuf_empty_pb.Empty>;
  executeModule: grpc.MethodDefinition<api_container_service_pb.ExecuteModuleArgs, api_container_service_pb.ExecuteModuleResponse>;
  getModuleInfo: grpc.MethodDefinition<api_container_service_pb.GetModuleInfoArgs, api_container_service_pb.GetModuleInfoResponse>;
  registerFilesArtifacts: grpc.MethodDefinition<api_container_service_pb.RegisterFilesArtifactsArgs, google_protobuf_empty_pb.Empty>;
  registerService: grpc.MethodDefinition<api_container_service_pb.RegisterServiceArgs, api_container_service_pb.RegisterServiceResponse>;
  startService: grpc.MethodDefinition<api_container_service_pb.StartServiceArgs, api_container_service_pb.StartServiceResponse>;
  getServiceInfo: grpc.MethodDefinition<api_container_service_pb.GetServiceInfoArgs, api_container_service_pb.GetServiceInfoResponse>;
  removeService: grpc.MethodDefinition<api_container_service_pb.RemoveServiceArgs, google_protobuf_empty_pb.Empty>;
  repartition: grpc.MethodDefinition<api_container_service_pb.RepartitionArgs, google_protobuf_empty_pb.Empty>;
  execCommand: grpc.MethodDefinition<api_container_service_pb.ExecCommandArgs, api_container_service_pb.ExecCommandResponse>;
  waitForHttpGetEndpointAvailability: grpc.MethodDefinition<api_container_service_pb.WaitForHttpGetEndpointAvailabilityArgs, google_protobuf_empty_pb.Empty>;
  waitForHttpPostEndpointAvailability: grpc.MethodDefinition<api_container_service_pb.WaitForHttpPostEndpointAvailabilityArgs, google_protobuf_empty_pb.Empty>;
  getServices: grpc.MethodDefinition<google_protobuf_empty_pb.Empty, api_container_service_pb.GetServicesResponse>;
  getModules: grpc.MethodDefinition<google_protobuf_empty_pb.Empty, api_container_service_pb.GetModulesResponse>;
  uploadFilesArtifact: grpc.MethodDefinition<api_container_service_pb.UploadFilesArtifactArgs, api_container_service_pb.UploadFilesArtifactResponse>;
  downloadFilesArtifact: grpc.MethodDefinition<api_container_service_pb.DownloadFilesArtifactArgs, api_container_service_pb.DownloadFilesArtifactResponse>;
}

export const ApiContainerServiceService: IApiContainerServiceService;

export interface IApiContainerServiceServer extends grpc.UntypedServiceImplementation {
  loadModule: grpc.handleUnaryCall<api_container_service_pb.LoadModuleArgs, api_container_service_pb.LoadModuleResponse>;
  unloadModule: grpc.handleUnaryCall<api_container_service_pb.UnloadModuleArgs, google_protobuf_empty_pb.Empty>;
  executeModule: grpc.handleUnaryCall<api_container_service_pb.ExecuteModuleArgs, api_container_service_pb.ExecuteModuleResponse>;
  getModuleInfo: grpc.handleUnaryCall<api_container_service_pb.GetModuleInfoArgs, api_container_service_pb.GetModuleInfoResponse>;
  registerFilesArtifacts: grpc.handleUnaryCall<api_container_service_pb.RegisterFilesArtifactsArgs, google_protobuf_empty_pb.Empty>;
  registerService: grpc.handleUnaryCall<api_container_service_pb.RegisterServiceArgs, api_container_service_pb.RegisterServiceResponse>;
  startService: grpc.handleUnaryCall<api_container_service_pb.StartServiceArgs, api_container_service_pb.StartServiceResponse>;
  getServiceInfo: grpc.handleUnaryCall<api_container_service_pb.GetServiceInfoArgs, api_container_service_pb.GetServiceInfoResponse>;
  removeService: grpc.handleUnaryCall<api_container_service_pb.RemoveServiceArgs, google_protobuf_empty_pb.Empty>;
  repartition: grpc.handleUnaryCall<api_container_service_pb.RepartitionArgs, google_protobuf_empty_pb.Empty>;
  execCommand: grpc.handleUnaryCall<api_container_service_pb.ExecCommandArgs, api_container_service_pb.ExecCommandResponse>;
  waitForHttpGetEndpointAvailability: grpc.handleUnaryCall<api_container_service_pb.WaitForHttpGetEndpointAvailabilityArgs, google_protobuf_empty_pb.Empty>;
  waitForHttpPostEndpointAvailability: grpc.handleUnaryCall<api_container_service_pb.WaitForHttpPostEndpointAvailabilityArgs, google_protobuf_empty_pb.Empty>;
  getServices: grpc.handleUnaryCall<google_protobuf_empty_pb.Empty, api_container_service_pb.GetServicesResponse>;
  getModules: grpc.handleUnaryCall<google_protobuf_empty_pb.Empty, api_container_service_pb.GetModulesResponse>;
  uploadFilesArtifact: grpc.handleUnaryCall<api_container_service_pb.UploadFilesArtifactArgs, api_container_service_pb.UploadFilesArtifactResponse>;
  downloadFilesArtifact: grpc.handleUnaryCall<api_container_service_pb.DownloadFilesArtifactArgs, api_container_service_pb.DownloadFilesArtifactResponse>;
}

export class ApiContainerServiceClient extends grpc.Client {
  constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
  loadModule(argument: api_container_service_pb.LoadModuleArgs, callback: grpc.requestCallback<api_container_service_pb.LoadModuleResponse>): grpc.ClientUnaryCall;
  loadModule(argument: api_container_service_pb.LoadModuleArgs, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<api_container_service_pb.LoadModuleResponse>): grpc.ClientUnaryCall;
  loadModule(argument: api_container_service_pb.LoadModuleArgs, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<api_container_service_pb.LoadModuleResponse>): grpc.ClientUnaryCall;
  unloadModule(argument: api_container_service_pb.UnloadModuleArgs, callback: grpc.requestCallback<google_protobuf_empty_pb.Empty>): grpc.ClientUnaryCall;
  unloadModule(argument: api_container_service_pb.UnloadModuleArgs, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<google_protobuf_empty_pb.Empty>): grpc.ClientUnaryCall;
  unloadModule(argument: api_container_service_pb.UnloadModuleArgs, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<google_protobuf_empty_pb.Empty>): grpc.ClientUnaryCall;
  executeModule(argument: api_container_service_pb.ExecuteModuleArgs, callback: grpc.requestCallback<api_container_service_pb.ExecuteModuleResponse>): grpc.ClientUnaryCall;
  executeModule(argument: api_container_service_pb.ExecuteModuleArgs, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<api_container_service_pb.ExecuteModuleResponse>): grpc.ClientUnaryCall;
  executeModule(argument: api_container_service_pb.ExecuteModuleArgs, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<api_container_service_pb.ExecuteModuleResponse>): grpc.ClientUnaryCall;
  getModuleInfo(argument: api_container_service_pb.GetModuleInfoArgs, callback: grpc.requestCallback<api_container_service_pb.GetModuleInfoResponse>): grpc.ClientUnaryCall;
  getModuleInfo(argument: api_container_service_pb.GetModuleInfoArgs, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<api_container_service_pb.GetModuleInfoResponse>): grpc.ClientUnaryCall;
  getModuleInfo(argument: api_container_service_pb.GetModuleInfoArgs, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<api_container_service_pb.GetModuleInfoResponse>): grpc.ClientUnaryCall;
  registerFilesArtifacts(argument: api_container_service_pb.RegisterFilesArtifactsArgs, callback: grpc.requestCallback<google_protobuf_empty_pb.Empty>): grpc.ClientUnaryCall;
  registerFilesArtifacts(argument: api_container_service_pb.RegisterFilesArtifactsArgs, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<google_protobuf_empty_pb.Empty>): grpc.ClientUnaryCall;
  registerFilesArtifacts(argument: api_container_service_pb.RegisterFilesArtifactsArgs, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<google_protobuf_empty_pb.Empty>): grpc.ClientUnaryCall;
  registerService(argument: api_container_service_pb.RegisterServiceArgs, callback: grpc.requestCallback<api_container_service_pb.RegisterServiceResponse>): grpc.ClientUnaryCall;
  registerService(argument: api_container_service_pb.RegisterServiceArgs, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<api_container_service_pb.RegisterServiceResponse>): grpc.ClientUnaryCall;
  registerService(argument: api_container_service_pb.RegisterServiceArgs, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<api_container_service_pb.RegisterServiceResponse>): grpc.ClientUnaryCall;
  startService(argument: api_container_service_pb.StartServiceArgs, callback: grpc.requestCallback<api_container_service_pb.StartServiceResponse>): grpc.ClientUnaryCall;
  startService(argument: api_container_service_pb.StartServiceArgs, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<api_container_service_pb.StartServiceResponse>): grpc.ClientUnaryCall;
  startService(argument: api_container_service_pb.StartServiceArgs, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<api_container_service_pb.StartServiceResponse>): grpc.ClientUnaryCall;
  getServiceInfo(argument: api_container_service_pb.GetServiceInfoArgs, callback: grpc.requestCallback<api_container_service_pb.GetServiceInfoResponse>): grpc.ClientUnaryCall;
  getServiceInfo(argument: api_container_service_pb.GetServiceInfoArgs, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<api_container_service_pb.GetServiceInfoResponse>): grpc.ClientUnaryCall;
  getServiceInfo(argument: api_container_service_pb.GetServiceInfoArgs, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<api_container_service_pb.GetServiceInfoResponse>): grpc.ClientUnaryCall;
  removeService(argument: api_container_service_pb.RemoveServiceArgs, callback: grpc.requestCallback<google_protobuf_empty_pb.Empty>): grpc.ClientUnaryCall;
  removeService(argument: api_container_service_pb.RemoveServiceArgs, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<google_protobuf_empty_pb.Empty>): grpc.ClientUnaryCall;
  removeService(argument: api_container_service_pb.RemoveServiceArgs, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<google_protobuf_empty_pb.Empty>): grpc.ClientUnaryCall;
  repartition(argument: api_container_service_pb.RepartitionArgs, callback: grpc.requestCallback<google_protobuf_empty_pb.Empty>): grpc.ClientUnaryCall;
  repartition(argument: api_container_service_pb.RepartitionArgs, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<google_protobuf_empty_pb.Empty>): grpc.ClientUnaryCall;
  repartition(argument: api_container_service_pb.RepartitionArgs, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<google_protobuf_empty_pb.Empty>): grpc.ClientUnaryCall;
  execCommand(argument: api_container_service_pb.ExecCommandArgs, callback: grpc.requestCallback<api_container_service_pb.ExecCommandResponse>): grpc.ClientUnaryCall;
  execCommand(argument: api_container_service_pb.ExecCommandArgs, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<api_container_service_pb.ExecCommandResponse>): grpc.ClientUnaryCall;
  execCommand(argument: api_container_service_pb.ExecCommandArgs, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<api_container_service_pb.ExecCommandResponse>): grpc.ClientUnaryCall;
  waitForHttpGetEndpointAvailability(argument: api_container_service_pb.WaitForHttpGetEndpointAvailabilityArgs, callback: grpc.requestCallback<google_protobuf_empty_pb.Empty>): grpc.ClientUnaryCall;
  waitForHttpGetEndpointAvailability(argument: api_container_service_pb.WaitForHttpGetEndpointAvailabilityArgs, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<google_protobuf_empty_pb.Empty>): grpc.ClientUnaryCall;
  waitForHttpGetEndpointAvailability(argument: api_container_service_pb.WaitForHttpGetEndpointAvailabilityArgs, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<google_protobuf_empty_pb.Empty>): grpc.ClientUnaryCall;
  waitForHttpPostEndpointAvailability(argument: api_container_service_pb.WaitForHttpPostEndpointAvailabilityArgs, callback: grpc.requestCallback<google_protobuf_empty_pb.Empty>): grpc.ClientUnaryCall;
  waitForHttpPostEndpointAvailability(argument: api_container_service_pb.WaitForHttpPostEndpointAvailabilityArgs, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<google_protobuf_empty_pb.Empty>): grpc.ClientUnaryCall;
  waitForHttpPostEndpointAvailability(argument: api_container_service_pb.WaitForHttpPostEndpointAvailabilityArgs, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<google_protobuf_empty_pb.Empty>): grpc.ClientUnaryCall;
  getServices(argument: google_protobuf_empty_pb.Empty, callback: grpc.requestCallback<api_container_service_pb.GetServicesResponse>): grpc.ClientUnaryCall;
  getServices(argument: google_protobuf_empty_pb.Empty, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<api_container_service_pb.GetServicesResponse>): grpc.ClientUnaryCall;
  getServices(argument: google_protobuf_empty_pb.Empty, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<api_container_service_pb.GetServicesResponse>): grpc.ClientUnaryCall;
  getModules(argument: google_protobuf_empty_pb.Empty, callback: grpc.requestCallback<api_container_service_pb.GetModulesResponse>): grpc.ClientUnaryCall;
  getModules(argument: google_protobuf_empty_pb.Empty, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<api_container_service_pb.GetModulesResponse>): grpc.ClientUnaryCall;
  getModules(argument: google_protobuf_empty_pb.Empty, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<api_container_service_pb.GetModulesResponse>): grpc.ClientUnaryCall;
  uploadFilesArtifact(argument: api_container_service_pb.UploadFilesArtifactArgs, callback: grpc.requestCallback<api_container_service_pb.UploadFilesArtifactResponse>): grpc.ClientUnaryCall;
  uploadFilesArtifact(argument: api_container_service_pb.UploadFilesArtifactArgs, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<api_container_service_pb.UploadFilesArtifactResponse>): grpc.ClientUnaryCall;
  uploadFilesArtifact(argument: api_container_service_pb.UploadFilesArtifactArgs, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<api_container_service_pb.UploadFilesArtifactResponse>): grpc.ClientUnaryCall;
  downloadFilesArtifact(argument: api_container_service_pb.DownloadFilesArtifactArgs, callback: grpc.requestCallback<api_container_service_pb.DownloadFilesArtifactResponse>): grpc.ClientUnaryCall;
  downloadFilesArtifact(argument: api_container_service_pb.DownloadFilesArtifactArgs, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<api_container_service_pb.DownloadFilesArtifactResponse>): grpc.ClientUnaryCall;
  downloadFilesArtifact(argument: api_container_service_pb.DownloadFilesArtifactArgs, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<api_container_service_pb.DownloadFilesArtifactResponse>): grpc.ClientUnaryCall;
}
