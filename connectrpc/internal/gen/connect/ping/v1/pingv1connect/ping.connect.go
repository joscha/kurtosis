// Copyright 2021-2023 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// The canonical location for this file is
// https://github.com/connectrpc/connect-go/blob/main/internal/proto/connect/ping/v1/ping.proto.

// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: connect/ping/v1/ping.proto

// The connect.ping.v1 package contains an echo service designed to test the
// connect-go implementation.
package pingv1connect

import (
	connect "connectrpc.com/connect"
	v1 "connectrpc.com/connect/internal/gen/connect/ping/v1"
	context "context"
	errors "errors"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_7_0

const (
	// PingServiceName is the fully-qualified name of the PingService service.
	PingServiceName = "connect.ping.v1.PingService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// PingServicePingProcedure is the fully-qualified name of the PingService's Ping RPC.
	PingServicePingProcedure = "/connect.ping.v1.PingService/Ping"
	// PingServiceFailProcedure is the fully-qualified name of the PingService's Fail RPC.
	PingServiceFailProcedure = "/connect.ping.v1.PingService/Fail"
	// PingServiceSumProcedure is the fully-qualified name of the PingService's Sum RPC.
	PingServiceSumProcedure = "/connect.ping.v1.PingService/Sum"
	// PingServiceCountUpProcedure is the fully-qualified name of the PingService's CountUp RPC.
	PingServiceCountUpProcedure = "/connect.ping.v1.PingService/CountUp"
	// PingServiceCumSumProcedure is the fully-qualified name of the PingService's CumSum RPC.
	PingServiceCumSumProcedure = "/connect.ping.v1.PingService/CumSum"
)

// PingServiceClient is a client for the connect.ping.v1.PingService service.
type PingServiceClient interface {
	// Ping sends a ping to the server to determine if it's reachable.
	Ping(context.Context, *connect.Request[v1.PingRequest]) (*connect.Response[v1.PingResponse], error)
	// Fail always fails.
	Fail(context.Context, *connect.Request[v1.FailRequest]) (*connect.Response[v1.FailResponse], error)
	// Sum calculates the sum of the numbers sent on the stream.
	Sum(context.Context) *connect.ClientStreamForClient[v1.SumRequest, v1.SumResponse]
	// CountUp returns a stream of the numbers up to the given request.
	CountUp(context.Context, *connect.Request[v1.CountUpRequest]) (*connect.ServerStreamForClient[v1.CountUpResponse], error)
	// CumSum determines the cumulative sum of all the numbers sent on the stream.
	CumSum(context.Context) *connect.BidiStreamForClient[v1.CumSumRequest, v1.CumSumResponse]
}

// NewPingServiceClient constructs a client for the connect.ping.v1.PingService service. By default,
// it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and
// sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC()
// or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewPingServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) PingServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &pingServiceClient{
		ping: connect.NewClient[v1.PingRequest, v1.PingResponse](
			httpClient,
			baseURL+PingServicePingProcedure,
			connect.WithIdempotency(connect.IdempotencyNoSideEffects),
			connect.WithClientOptions(opts...),
		),
		fail: connect.NewClient[v1.FailRequest, v1.FailResponse](
			httpClient,
			baseURL+PingServiceFailProcedure,
			opts...,
		),
		sum: connect.NewClient[v1.SumRequest, v1.SumResponse](
			httpClient,
			baseURL+PingServiceSumProcedure,
			opts...,
		),
		countUp: connect.NewClient[v1.CountUpRequest, v1.CountUpResponse](
			httpClient,
			baseURL+PingServiceCountUpProcedure,
			opts...,
		),
		cumSum: connect.NewClient[v1.CumSumRequest, v1.CumSumResponse](
			httpClient,
			baseURL+PingServiceCumSumProcedure,
			opts...,
		),
	}
}

// pingServiceClient implements PingServiceClient.
type pingServiceClient struct {
	ping    *connect.Client[v1.PingRequest, v1.PingResponse]
	fail    *connect.Client[v1.FailRequest, v1.FailResponse]
	sum     *connect.Client[v1.SumRequest, v1.SumResponse]
	countUp *connect.Client[v1.CountUpRequest, v1.CountUpResponse]
	cumSum  *connect.Client[v1.CumSumRequest, v1.CumSumResponse]
}

// Ping calls connect.ping.v1.PingService.Ping.
func (c *pingServiceClient) Ping(ctx context.Context, req *connect.Request[v1.PingRequest]) (*connect.Response[v1.PingResponse], error) {
	return c.ping.CallUnary(ctx, req)
}

// Fail calls connect.ping.v1.PingService.Fail.
func (c *pingServiceClient) Fail(ctx context.Context, req *connect.Request[v1.FailRequest]) (*connect.Response[v1.FailResponse], error) {
	return c.fail.CallUnary(ctx, req)
}

// Sum calls connect.ping.v1.PingService.Sum.
func (c *pingServiceClient) Sum(ctx context.Context) *connect.ClientStreamForClient[v1.SumRequest, v1.SumResponse] {
	return c.sum.CallClientStream(ctx)
}

// CountUp calls connect.ping.v1.PingService.CountUp.
func (c *pingServiceClient) CountUp(ctx context.Context, req *connect.Request[v1.CountUpRequest]) (*connect.ServerStreamForClient[v1.CountUpResponse], error) {
	return c.countUp.CallServerStream(ctx, req)
}

// CumSum calls connect.ping.v1.PingService.CumSum.
func (c *pingServiceClient) CumSum(ctx context.Context) *connect.BidiStreamForClient[v1.CumSumRequest, v1.CumSumResponse] {
	return c.cumSum.CallBidiStream(ctx)
}

// PingServiceHandler is an implementation of the connect.ping.v1.PingService service.
type PingServiceHandler interface {
	// Ping sends a ping to the server to determine if it's reachable.
	Ping(context.Context, *connect.Request[v1.PingRequest]) (*connect.Response[v1.PingResponse], error)
	// Fail always fails.
	Fail(context.Context, *connect.Request[v1.FailRequest]) (*connect.Response[v1.FailResponse], error)
	// Sum calculates the sum of the numbers sent on the stream.
	Sum(context.Context, *connect.ClientStream[v1.SumRequest]) (*connect.Response[v1.SumResponse], error)
	// CountUp returns a stream of the numbers up to the given request.
	CountUp(context.Context, *connect.Request[v1.CountUpRequest], *connect.ServerStream[v1.CountUpResponse]) error
	// CumSum determines the cumulative sum of all the numbers sent on the stream.
	CumSum(context.Context, *connect.BidiStream[v1.CumSumRequest, v1.CumSumResponse]) error
}

// NewPingServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewPingServiceHandler(svc PingServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	pingServicePingHandler := connect.NewUnaryHandler(
		PingServicePingProcedure,
		svc.Ping,
		connect.WithIdempotency(connect.IdempotencyNoSideEffects),
		connect.WithHandlerOptions(opts...),
	)
	pingServiceFailHandler := connect.NewUnaryHandler(
		PingServiceFailProcedure,
		svc.Fail,
		opts...,
	)
	pingServiceSumHandler := connect.NewClientStreamHandler(
		PingServiceSumProcedure,
		svc.Sum,
		opts...,
	)
	pingServiceCountUpHandler := connect.NewServerStreamHandler(
		PingServiceCountUpProcedure,
		svc.CountUp,
		opts...,
	)
	pingServiceCumSumHandler := connect.NewBidiStreamHandler(
		PingServiceCumSumProcedure,
		svc.CumSum,
		opts...,
	)
	return "/connect.ping.v1.PingService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case PingServicePingProcedure:
			pingServicePingHandler.ServeHTTP(w, r)
		case PingServiceFailProcedure:
			pingServiceFailHandler.ServeHTTP(w, r)
		case PingServiceSumProcedure:
			pingServiceSumHandler.ServeHTTP(w, r)
		case PingServiceCountUpProcedure:
			pingServiceCountUpHandler.ServeHTTP(w, r)
		case PingServiceCumSumProcedure:
			pingServiceCumSumHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedPingServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedPingServiceHandler struct{}

func (UnimplementedPingServiceHandler) Ping(context.Context, *connect.Request[v1.PingRequest]) (*connect.Response[v1.PingResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("connect.ping.v1.PingService.Ping is not implemented"))
}

func (UnimplementedPingServiceHandler) Fail(context.Context, *connect.Request[v1.FailRequest]) (*connect.Response[v1.FailResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("connect.ping.v1.PingService.Fail is not implemented"))
}

func (UnimplementedPingServiceHandler) Sum(context.Context, *connect.ClientStream[v1.SumRequest]) (*connect.Response[v1.SumResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("connect.ping.v1.PingService.Sum is not implemented"))
}

func (UnimplementedPingServiceHandler) CountUp(context.Context, *connect.Request[v1.CountUpRequest], *connect.ServerStream[v1.CountUpResponse]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("connect.ping.v1.PingService.CountUp is not implemented"))
}

func (UnimplementedPingServiceHandler) CumSum(context.Context, *connect.BidiStream[v1.CumSumRequest, v1.CumSumResponse]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("connect.ping.v1.PingService.CumSum is not implemented"))
}