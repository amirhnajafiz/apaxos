package client

import (
	"context"
	"fmt"

	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ApaxosDialer is used to call RPCs for apaxos protocol.
type ApaxosDialer struct {
	Logger *zap.Logger
}

// connect should be called in the beginning of each method to establish a connection.
func (a *ApaxosDialer) connect(address string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("[grpc/client/apaxosDialer] failed to open connection to %s: %v", address, err)
	}

	return conn, nil
}

// Propose sends a prepare message to the given address.
func (a *ApaxosDialer) Propose(address string, message *apaxos.PrepareMessage) {
	// base connection
	conn, err := a.connect(address)
	if err != nil {
		a.Logger.Debug("failed to connect", zap.String("address", address), zap.Error(err))
		return
	}
	defer conn.Close()

	// call Propose RPC
	_, err = apaxos.NewApaxosClient(conn).Propose(context.Background(), message)
	if err != nil {
		a.Logger.Debug("failed to call propose rpc", zap.String("address", address), zap.Error(err))
	}
}

// Promise sends a promise message to the give address.
func (a *ApaxosDialer) Promise(address string, message *apaxos.PromiseMessage) {
	// base connection
	conn, err := a.connect(address)
	if err != nil {
		a.Logger.Debug("failed to connect", zap.String("address", address), zap.Error(err))
		return
	}
	defer conn.Close()

	// call Promise RPC
	_, err = apaxos.NewApaxosClient(conn).Promise(context.Background(), message)
	if err != nil {
		a.Logger.Debug("failed to call promise rpc", zap.String("address", address), zap.Error(err))
	}
}

// Accept sends an accept message to the given address.
func (a *ApaxosDialer) Accept(address string, message *apaxos.AcceptMessage) {
	// base connection
	conn, err := a.connect(address)
	if err != nil {
		a.Logger.Debug("failed to connect", zap.String("address", address), zap.Error(err))
		return
	}
	defer conn.Close()

	// call Accept RPC
	_, err = apaxos.NewApaxosClient(conn).Accept(context.Background(), message)
	if err != nil {
		a.Logger.Debug("failed to call accept rpc", zap.String("address", address), zap.Error(err))
	}
}

// Accepted just calls the accepted RPC on the given address.
func (a *ApaxosDialer) Accepted(address string) {
	// base connection
	conn, err := a.connect(address)
	if err != nil {
		a.Logger.Debug("failed to connect", zap.String("address", address), zap.Error(err))
		return
	}
	defer conn.Close()

	// call Accepted RPC
	_, err = apaxos.NewApaxosClient(conn).Accepted(context.Background(), &emptypb.Empty{})
	if err != nil {
		a.Logger.Debug("failed to call accepted rpc", zap.String("address", address), zap.Error(err))
	}
}

// Commit just calls the commit RPC on the given address.
func (a *ApaxosDialer) Commit(address string) {
	// base connection
	conn, err := a.connect(address)
	if err != nil {
		a.Logger.Debug("failed to connect", zap.String("address", address), zap.Error(err))
		return
	}
	defer conn.Close()

	// call Commit RPC
	_, err = apaxos.NewApaxosClient(conn).Commit(context.Background(), &emptypb.Empty{})
	if err != nil {
		a.Logger.Debug("failed to call commit rpc", zap.String("address", address), zap.Error(err))
	}
}

// Sync sends a sync message to the given address.
func (a *ApaxosDialer) Sync(address string, messages *apaxos.SyncMessage) {
	// base connection
	conn, err := a.connect(address)
	if err != nil {
		a.Logger.Debug("failed to connect", zap.String("address", address), zap.Error(err))
		return
	}
	defer conn.Close()

	// call Sync RPC
	_, err = apaxos.NewApaxosClient(conn).Sync(context.Background(), messages)
	if err != nil {
		a.Logger.Debug("failed to call sync rpc", zap.String("address", address), zap.Error(err))
	}
}
