// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.3
// source: portal_daemon_service.proto

package kurtosis_portal_rpc_api_bindings

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	KurtosisPortalDaemon_Ping_FullMethodName                   = "/portal_daemon_api.KurtosisPortalDaemon/Ping"
	KurtosisPortalDaemon_ForwardUserServicePort_FullMethodName = "/portal_daemon_api.KurtosisPortalDaemon/ForwardUserServicePort"
)

// KurtosisPortalDaemonClient is the client API for KurtosisPortalDaemon service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KurtosisPortalDaemonClient interface {
	// To check availability
	Ping(ctx context.Context, in *PortalPing, opts ...grpc.CallOption) (*PortalPong, error)
	ForwardUserServicePort(ctx context.Context, in *ForwardUserServicePortArgs, opts ...grpc.CallOption) (*ForwardUserServicePortResponse, error)
}

type kurtosisPortalDaemonClient struct {
	cc grpc.ClientConnInterface
}

func NewKurtosisPortalDaemonClient(cc grpc.ClientConnInterface) KurtosisPortalDaemonClient {
	return &kurtosisPortalDaemonClient{cc}
}

func (c *kurtosisPortalDaemonClient) Ping(ctx context.Context, in *PortalPing, opts ...grpc.CallOption) (*PortalPong, error) {
	out := new(PortalPong)
	err := c.cc.Invoke(ctx, KurtosisPortalDaemon_Ping_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kurtosisPortalDaemonClient) ForwardUserServicePort(ctx context.Context, in *ForwardUserServicePortArgs, opts ...grpc.CallOption) (*ForwardUserServicePortResponse, error) {
	out := new(ForwardUserServicePortResponse)
	err := c.cc.Invoke(ctx, KurtosisPortalDaemon_ForwardUserServicePort_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KurtosisPortalDaemonServer is the server API for KurtosisPortalDaemon service.
// All implementations should embed UnimplementedKurtosisPortalDaemonServer
// for forward compatibility
type KurtosisPortalDaemonServer interface {
	// To check availability
	Ping(context.Context, *PortalPing) (*PortalPong, error)
	ForwardUserServicePort(context.Context, *ForwardUserServicePortArgs) (*ForwardUserServicePortResponse, error)
}

// UnimplementedKurtosisPortalDaemonServer should be embedded to have forward compatible implementations.
type UnimplementedKurtosisPortalDaemonServer struct {
}

func (UnimplementedKurtosisPortalDaemonServer) Ping(context.Context, *PortalPing) (*PortalPong, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedKurtosisPortalDaemonServer) ForwardUserServicePort(context.Context, *ForwardUserServicePortArgs) (*ForwardUserServicePortResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ForwardUserServicePort not implemented")
}

// UnsafeKurtosisPortalDaemonServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KurtosisPortalDaemonServer will
// result in compilation errors.
type UnsafeKurtosisPortalDaemonServer interface {
	mustEmbedUnimplementedKurtosisPortalDaemonServer()
}

func RegisterKurtosisPortalDaemonServer(s grpc.ServiceRegistrar, srv KurtosisPortalDaemonServer) {
	s.RegisterService(&KurtosisPortalDaemon_ServiceDesc, srv)
}

func _KurtosisPortalDaemon_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PortalPing)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KurtosisPortalDaemonServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KurtosisPortalDaemon_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KurtosisPortalDaemonServer).Ping(ctx, req.(*PortalPing))
	}
	return interceptor(ctx, in, info, handler)
}

func _KurtosisPortalDaemon_ForwardUserServicePort_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ForwardUserServicePortArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KurtosisPortalDaemonServer).ForwardUserServicePort(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KurtosisPortalDaemon_ForwardUserServicePort_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KurtosisPortalDaemonServer).ForwardUserServicePort(ctx, req.(*ForwardUserServicePortArgs))
	}
	return interceptor(ctx, in, info, handler)
}

// KurtosisPortalDaemon_ServiceDesc is the grpc.ServiceDesc for KurtosisPortalDaemon service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KurtosisPortalDaemon_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "portal_daemon_api.KurtosisPortalDaemon",
	HandlerType: (*KurtosisPortalDaemonServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _KurtosisPortalDaemon_Ping_Handler,
		},
		{
			MethodName: "ForwardUserServicePort",
			Handler:    _KurtosisPortalDaemon_ForwardUserServicePort_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "portal_daemon_service.proto",
}