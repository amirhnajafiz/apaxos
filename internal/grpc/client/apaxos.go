package client

import (
	"context"
	"fmt"
	"log"

	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ApaxosDialer is used to call RPCs for apaxos protocol.
type ApaxosDialer struct {
	Logger *zap.Logger
}

// connect should be called in the beginning of each method to establish a connection.
func (a *ApaxosDialer) connect(address string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption

	conn, err := grpc.NewClient(address, opts...)
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
		log.Printf("failed to call %s: %v\n", address, err)
		return
	}
	defer conn.Close()

	// call Propose RPC
	_, _ = apaxos.NewApaxosClient(conn).Propose(context.Background(), message)
}

// Promise sends a promise message to the give address.
func (a *ApaxosDialer) Promise(address string, message *apaxos.PromiseMessage) {
	// base connection
	conn, err := a.connect(address)
	if err != nil {
		log.Printf("failed to call %s: %v\n", address, err)
		return
	}
	defer conn.Close()

	// call Promise RPC
	_, _ = apaxos.NewApaxosClient(conn).Promise(context.Background(), message)
}

// Accept sends an accept message to the given address.
func (a *ApaxosDialer) Accept(address string, message *apaxos.AcceptMessage) {
	// base connection
	conn, err := a.connect(address)
	if err != nil {
		log.Printf("failed to call %s: %v\n", address, err)
		return
	}
	defer conn.Close()

	// call Accept RPC
	_, _ = apaxos.NewApaxosClient(conn).Accept(context.Background(), message)
}

// Accepted just calls the accepted RPC on the given address.
func (a *ApaxosDialer) Accepted(address string) {
	// base connection
	conn, err := a.connect(address)
	if err != nil {
		log.Printf("failed to call %s: %v\n", address, err)
		return
	}
	defer conn.Close()

	// call Accepted RPC
	_, _ = apaxos.NewApaxosClient(conn).Accepted(context.Background(), &emptypb.Empty{})
}

// Commit just calls the commit RPC on the given address.
func (a *ApaxosDialer) Commit(address string) {
	// base connection
	conn, err := a.connect(address)
	if err != nil {
		log.Printf("failed to call %s: %v\n", address, err)
		return
	}
	defer conn.Close()

	// call Commit RPC
	_, _ = apaxos.NewApaxosClient(conn).Commit(context.Background(), &emptypb.Empty{})
}

// Sync sends a sync message to the given address.
func (a *ApaxosDialer) Sync(address string, messages *apaxos.SyncMessage) {
	// base connection
	conn, err := a.connect(address)
	if err != nil {
		log.Printf("failed to call %s: %v\n", address, err)
		return
	}
	defer conn.Close()

	// call Sync RPC
	_, _ = apaxos.NewApaxosClient(conn).Sync(context.Background(), messages)
}
