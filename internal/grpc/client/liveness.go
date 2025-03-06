package client

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/f24-cse535/apaxos/pkg/rpc/liveness"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// LivenessDialer is used to call RPC methods for checking a server
// liveness status.
type LivenessDialer struct {
	Logger *zap.Logger
}

// connect should be called in the beginning of each method to establish a connection.
func (l *LivenessDialer) connect(address string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("[grpc/client/livenessDialer] failed to open connection to %s: %v", address, err)
	}

	return conn, nil
}

// Ping is used to send a ping request to a server. If the server is available, it returns true.
func (l *LivenessDialer) Ping(address string) bool {
	// base connection
	conn, err := l.connect(address)
	if err != nil {
		l.Logger.Debug("failed to connect", zap.String("address", address), zap.Error(err))

		return false
	}
	defer conn.Close()

	// call RPC of ping
	resp, err := liveness.NewLivenessClient(conn).Ping(context.Background(), &liveness.LivePingMessage{
		Random: rand.Int64(), // a non-negative number
	})
	if err != nil {
		l.Logger.Debug("failed to call ping RPC", zap.String("address", address), zap.Error(err))

		return false
	}

	// check the response code
	if resp.Random == -1 {
		return false
	}

	// server is ok
	return true
}

// ChangeState is used to modify the state of a gRPC server.
// if the state is true, then the server is alive, else the server will be blocked.
func (l *LivenessDialer) ChangeState(address string, state bool) error {
	// base connection
	conn, err := l.connect(address)
	if err != nil {
		return fmt.Errorf("failed to call %s: %v", address, err)
	}
	defer conn.Close()

	// call RPC of change status
	resp, err := liveness.NewLivenessClient(conn).ChangeStatus(context.Background(), &liveness.LiveChangeStatusMessage{
		Status: state,
	})
	if err != nil {
		return fmt.Errorf("failed to change server %s status: %v", address, err)
	}

	// check the response for changes
	if resp.GetStatus() != state {
		return fmt.Errorf("server status is not changed to %t, it is %t", state, resp.GetStatus())
	}

	return nil
}
