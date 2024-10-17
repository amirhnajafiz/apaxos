package grpc

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// package name for the liveness service.
const livenessServicePrefix = "/liveness."

// logPerRPC logs request info per each RPC call.
func (b *Bootstrap) logPerRPCUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	b.Logger.Info("rpc called", zap.String("method", info.FullMethod))

	return handler(ctx, req)
}

// selectiveStatusCheck interceptor checks the status
// of a service before running the gRPC function.
func (b *Bootstrap) selectiveStatusCheckUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// if status is true, allow all services to proceed
	if b.Memory.GetServiceStatus() {
		return handler(ctx, req)
	}

	// Ii status is false, only allow services in the liveness package
	if len(info.FullMethod) > len(livenessServicePrefix) && info.FullMethod[:len(livenessServicePrefix)] == livenessServicePrefix {
		return handler(ctx, req) // allow liveness service to proceed
	}

	// block all other services
	return nil, status.Error(13, "service is not responding") // return an error for blocked services
}
