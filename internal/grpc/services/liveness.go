package services

import (
	"context"

	"github.com/f24-cse535/apaxos/internal/storage/local"
	"github.com/f24-cse535/apaxos/pkg/rpc/liveness"
)

// liveness server handles the running state of the gRPC server.
type Liveness struct {
	liveness.UnimplementedLivenessServer
	Memory *local.Memory
}

// Ping RPC is used to check if a server is alive and can process or not.
func (l *Liveness) Ping(ctx context.Context, input *liveness.LivePingMessage) (*liveness.LivePingMessage, error) {
	if l.Memory.GetServiceStatus() {
		return &liveness.LivePingMessage{
			Random: input.GetRandom(),
		}, nil
	}

	return &liveness.LivePingMessage{Random: -1}, nil
}

// ChangeStatus is used to update the liveness of the gRPC server.
func (l *Liveness) ChangeStatus(ctx context.Context, input *liveness.LiveChangeStatusMessage) (*liveness.LiveChangeStatusMessage, error) {
	l.Memory.SetServiceStatus(input.GetStatus())

	return &liveness.LiveChangeStatusMessage{Status: l.Memory.GetServiceStatus()}, nil
}
