// GENERATED CODE -- DO NOT EDIT!

// package: portal_daemon_api
// file: portal_daemon_service.proto

import * as portal_daemon_service_pb from "./portal_daemon_service_pb";
import * as grpc from "@grpc/grpc-js";

interface IKurtosisPortalDaemonService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
  ping: grpc.MethodDefinition<portal_daemon_service_pb.PortalPing, portal_daemon_service_pb.PortalPong>;
  forwardUserServicePort: grpc.MethodDefinition<portal_daemon_service_pb.ForwardUserServicePortArgs, portal_daemon_service_pb.ForwardUserServicePortResponse>;
}

export const KurtosisPortalDaemonService: IKurtosisPortalDaemonService;

export interface IKurtosisPortalDaemonServer extends grpc.UntypedServiceImplementation {
  ping: grpc.handleUnaryCall<portal_daemon_service_pb.PortalPing, portal_daemon_service_pb.PortalPong>;
  forwardUserServicePort: grpc.handleUnaryCall<portal_daemon_service_pb.ForwardUserServicePortArgs, portal_daemon_service_pb.ForwardUserServicePortResponse>;
}

export class KurtosisPortalDaemonClient extends grpc.Client {
  constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
  ping(argument: portal_daemon_service_pb.PortalPing, callback: grpc.requestCallback<portal_daemon_service_pb.PortalPong>): grpc.ClientUnaryCall;
  ping(argument: portal_daemon_service_pb.PortalPing, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<portal_daemon_service_pb.PortalPong>): grpc.ClientUnaryCall;
  ping(argument: portal_daemon_service_pb.PortalPing, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<portal_daemon_service_pb.PortalPong>): grpc.ClientUnaryCall;
  forwardUserServicePort(argument: portal_daemon_service_pb.ForwardUserServicePortArgs, callback: grpc.requestCallback<portal_daemon_service_pb.ForwardUserServicePortResponse>): grpc.ClientUnaryCall;
  forwardUserServicePort(argument: portal_daemon_service_pb.ForwardUserServicePortArgs, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<portal_daemon_service_pb.ForwardUserServicePortResponse>): grpc.ClientUnaryCall;
  forwardUserServicePort(argument: portal_daemon_service_pb.ForwardUserServicePortArgs, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<portal_daemon_service_pb.ForwardUserServicePortResponse>): grpc.ClientUnaryCall;
}