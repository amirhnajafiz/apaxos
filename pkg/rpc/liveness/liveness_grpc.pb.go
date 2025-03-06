// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: liveness.proto

package liveness

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Liveness_Ping_FullMethodName         = "/liveness.Liveness/Ping"
	Liveness_ChangeStatus_FullMethodName = "/liveness.Liveness/ChangeStatus"
)

// LivenessClient is the client API for Liveness service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// creating rpc server for node's liveness status
type LivenessClient interface {
	Ping(ctx context.Context, in *LivePingMessage, opts ...grpc.CallOption) (*LivePingMessage, error)
	ChangeStatus(ctx context.Context, in *LiveChangeStatusMessage, opts ...grpc.CallOption) (*LiveChangeStatusMessage, error)
}

type livenessClient struct {
	cc grpc.ClientConnInterface
}

func NewLivenessClient(cc grpc.ClientConnInterface) LivenessClient {
	return &livenessClient{cc}
}

func (c *livenessClient) Ping(ctx context.Context, in *LivePingMessage, opts ...grpc.CallOption) (*LivePingMessage, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LivePingMessage)
	err := c.cc.Invoke(ctx, Liveness_Ping_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *livenessClient) ChangeStatus(ctx context.Context, in *LiveChangeStatusMessage, opts ...grpc.CallOption) (*LiveChangeStatusMessage, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LiveChangeStatusMessage)
	err := c.cc.Invoke(ctx, Liveness_ChangeStatus_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LivenessServer is the server API for Liveness service.
// All implementations must embed UnimplementedLivenessServer
// for forward compatibility.
//
// creating rpc server for node's liveness status
type LivenessServer interface {
	Ping(context.Context, *LivePingMessage) (*LivePingMessage, error)
	ChangeStatus(context.Context, *LiveChangeStatusMessage) (*LiveChangeStatusMessage, error)
	mustEmbedUnimplementedLivenessServer()
}

// UnimplementedLivenessServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedLivenessServer struct{}

func (UnimplementedLivenessServer) Ping(context.Context, *LivePingMessage) (*LivePingMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedLivenessServer) ChangeStatus(context.Context, *LiveChangeStatusMessage) (*LiveChangeStatusMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeStatus not implemented")
}
func (UnimplementedLivenessServer) mustEmbedUnimplementedLivenessServer() {}
func (UnimplementedLivenessServer) testEmbeddedByValue()                  {}

// UnsafeLivenessServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LivenessServer will
// result in compilation errors.
type UnsafeLivenessServer interface {
	mustEmbedUnimplementedLivenessServer()
}

func RegisterLivenessServer(s grpc.ServiceRegistrar, srv LivenessServer) {
	// If the following call pancis, it indicates UnimplementedLivenessServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Liveness_ServiceDesc, srv)
}

func _Liveness_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LivePingMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LivenessServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Liveness_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LivenessServer).Ping(ctx, req.(*LivePingMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Liveness_ChangeStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LiveChangeStatusMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LivenessServer).ChangeStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Liveness_ChangeStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LivenessServer).ChangeStatus(ctx, req.(*LiveChangeStatusMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// Liveness_ServiceDesc is the grpc.ServiceDesc for Liveness service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Liveness_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "liveness.Liveness",
	HandlerType: (*LivenessServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Liveness_Ping_Handler,
		},
		{
			MethodName: "ChangeStatus",
			Handler:    _Liveness_ChangeStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "liveness.proto",
}
