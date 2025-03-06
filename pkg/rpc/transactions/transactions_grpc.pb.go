// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: transactions.proto

package transactions

import (
	apaxos "github.com/f24-cse535/apaxos/pkg/rpc/apaxos"
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Transactions_NewTransaction_FullMethodName = "/transactions.Transactions/NewTransaction"
	Transactions_PrintBalance_FullMethodName   = "/transactions.Transactions/PrintBalance"
	Transactions_PrintLogs_FullMethodName      = "/transactions.Transactions/PrintLogs"
	Transactions_PrintDB_FullMethodName        = "/transactions.Transactions/PrintDB"
	Transactions_Performance_FullMethodName    = "/transactions.Transactions/Performance"
)

// TransactionsClient is the client API for Transactions service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// creating rpc services for transactions and apaxos.
// the transactions service is for handling client-server calls.
type TransactionsClient interface {
	NewTransaction(ctx context.Context, in *apaxos.Transaction, opts ...grpc.CallOption) (*TransactionResponse, error)
	PrintBalance(ctx context.Context, in *PrintBalanceRequest, opts ...grpc.CallOption) (*PrintBalanceResponse, error)
	PrintLogs(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (grpc.ServerStreamingClient[apaxos.Block], error)
	PrintDB(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (grpc.ServerStreamingClient[apaxos.Block], error)
	Performance(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PerformanceResponse, error)
}

type transactionsClient struct {
	cc grpc.ClientConnInterface
}

func NewTransactionsClient(cc grpc.ClientConnInterface) TransactionsClient {
	return &transactionsClient{cc}
}

func (c *transactionsClient) NewTransaction(ctx context.Context, in *apaxos.Transaction, opts ...grpc.CallOption) (*TransactionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TransactionResponse)
	err := c.cc.Invoke(ctx, Transactions_NewTransaction_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transactionsClient) PrintBalance(ctx context.Context, in *PrintBalanceRequest, opts ...grpc.CallOption) (*PrintBalanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PrintBalanceResponse)
	err := c.cc.Invoke(ctx, Transactions_PrintBalance_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transactionsClient) PrintLogs(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (grpc.ServerStreamingClient[apaxos.Block], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Transactions_ServiceDesc.Streams[0], Transactions_PrintLogs_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[emptypb.Empty, apaxos.Block]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Transactions_PrintLogsClient = grpc.ServerStreamingClient[apaxos.Block]

func (c *transactionsClient) PrintDB(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (grpc.ServerStreamingClient[apaxos.Block], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Transactions_ServiceDesc.Streams[1], Transactions_PrintDB_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[emptypb.Empty, apaxos.Block]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Transactions_PrintDBClient = grpc.ServerStreamingClient[apaxos.Block]

func (c *transactionsClient) Performance(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PerformanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PerformanceResponse)
	err := c.cc.Invoke(ctx, Transactions_Performance_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TransactionsServer is the server API for Transactions service.
// All implementations must embed UnimplementedTransactionsServer
// for forward compatibility.
//
// creating rpc services for transactions and apaxos.
// the transactions service is for handling client-server calls.
type TransactionsServer interface {
	NewTransaction(context.Context, *apaxos.Transaction) (*TransactionResponse, error)
	PrintBalance(context.Context, *PrintBalanceRequest) (*PrintBalanceResponse, error)
	PrintLogs(*emptypb.Empty, grpc.ServerStreamingServer[apaxos.Block]) error
	PrintDB(*emptypb.Empty, grpc.ServerStreamingServer[apaxos.Block]) error
	Performance(context.Context, *emptypb.Empty) (*PerformanceResponse, error)
	mustEmbedUnimplementedTransactionsServer()
}

// UnimplementedTransactionsServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedTransactionsServer struct{}

func (UnimplementedTransactionsServer) NewTransaction(context.Context, *apaxos.Transaction) (*TransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewTransaction not implemented")
}
func (UnimplementedTransactionsServer) PrintBalance(context.Context, *PrintBalanceRequest) (*PrintBalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PrintBalance not implemented")
}
func (UnimplementedTransactionsServer) PrintLogs(*emptypb.Empty, grpc.ServerStreamingServer[apaxos.Block]) error {
	return status.Errorf(codes.Unimplemented, "method PrintLogs not implemented")
}
func (UnimplementedTransactionsServer) PrintDB(*emptypb.Empty, grpc.ServerStreamingServer[apaxos.Block]) error {
	return status.Errorf(codes.Unimplemented, "method PrintDB not implemented")
}
func (UnimplementedTransactionsServer) Performance(context.Context, *emptypb.Empty) (*PerformanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Performance not implemented")
}
func (UnimplementedTransactionsServer) mustEmbedUnimplementedTransactionsServer() {}
func (UnimplementedTransactionsServer) testEmbeddedByValue()                      {}

// UnsafeTransactionsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TransactionsServer will
// result in compilation errors.
type UnsafeTransactionsServer interface {
	mustEmbedUnimplementedTransactionsServer()
}

func RegisterTransactionsServer(s grpc.ServiceRegistrar, srv TransactionsServer) {
	// If the following call pancis, it indicates UnimplementedTransactionsServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Transactions_ServiceDesc, srv)
}

func _Transactions_NewTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(apaxos.Transaction)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionsServer).NewTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Transactions_NewTransaction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionsServer).NewTransaction(ctx, req.(*apaxos.Transaction))
	}
	return interceptor(ctx, in, info, handler)
}

func _Transactions_PrintBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PrintBalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionsServer).PrintBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Transactions_PrintBalance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionsServer).PrintBalance(ctx, req.(*PrintBalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Transactions_PrintLogs_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TransactionsServer).PrintLogs(m, &grpc.GenericServerStream[emptypb.Empty, apaxos.Block]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Transactions_PrintLogsServer = grpc.ServerStreamingServer[apaxos.Block]

func _Transactions_PrintDB_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TransactionsServer).PrintDB(m, &grpc.GenericServerStream[emptypb.Empty, apaxos.Block]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Transactions_PrintDBServer = grpc.ServerStreamingServer[apaxos.Block]

func _Transactions_Performance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionsServer).Performance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Transactions_Performance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionsServer).Performance(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Transactions_ServiceDesc is the grpc.ServiceDesc for Transactions service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Transactions_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "transactions.Transactions",
	HandlerType: (*TransactionsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NewTransaction",
			Handler:    _Transactions_NewTransaction_Handler,
		},
		{
			MethodName: "PrintBalance",
			Handler:    _Transactions_PrintBalance_Handler,
		},
		{
			MethodName: "Performance",
			Handler:    _Transactions_Performance_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PrintLogs",
			Handler:       _Transactions_PrintLogs_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "PrintDB",
			Handler:       _Transactions_PrintDB_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "transactions.proto",
}
