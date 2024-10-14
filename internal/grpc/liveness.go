package grpc

import (
	"context"
	"errors"

	"github.com/f24-cse535/apaxos/pkg/rpc/liveness"
)

// liveness server handles the running state of the gRPC server.
type livenessServer struct {
	liveness.UnimplementedLivenessServer
	state bool
}

// Ping RPC is used to check if a server is alive and can process or not.
func (l *livenessServer) Ping(ctx context.Context, input *liveness.LivePingMessage) (*liveness.LivePingMessage, error) {
	if l.state {
		return &liveness.LivePingMessage{
			Random: input.GetRandom(),
		}, nil
	}

	return &liveness.LivePingMessage{Random: -1}, errors.New("server is blocked")
}

// ChangeStatus is used to update the liveness of the gRPC server.
func (l *livenessServer) ChangeStatus(ctx context.Context, input *liveness.LiveChangeStatusMessage) (*liveness.LiveChangeStatusMessage, error) {
	l.state = input.GetStatus()

	return &liveness.LiveChangeStatusMessage{Status: l.state}, nil
}
