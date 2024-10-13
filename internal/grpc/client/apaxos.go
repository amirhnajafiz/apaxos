package client

import (
	"context"
	"fmt"
	"log"

	"github.com/f24-cse535/apaxos/pkg/transactions"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ApaxosDialer is used to call RPCs for apaxos protocol.
type ApaxosDialer struct{}

// connect should be called in the beginning of each method to establish a connection.
func (a *ApaxosDialer) connect(address string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption

	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, fmt.Errorf("[grpc/client/apaxosdialer] failed to make rpc call: %v", err)
	}

	return conn, nil
}

// Propose message is sent by a proposer.
func (a *ApaxosDialer) Propose(address string, message *transactions.PrepareMessage) {
	conn, err := a.connect(address)
	if err != nil {
		log.Printf("failed to call %s: %v\n", address, err)
		return
	}
	defer conn.Close()

	_, _ = transactions.NewApaxosClient(conn).Propose(context.Background(), message)
}

// Promise message is sent by an acceptor.
func (a *ApaxosDialer) Promise(address string, message *transactions.PromiseMessage) {
	conn, err := a.connect(address)
	if err != nil {
		log.Printf("failed to call %s: %v\n", address, err)
		return
	}
	defer conn.Close()

	_, _ = transactions.NewApaxosClient(conn).Promise(context.Background(), message)
}

// Accept message is sent by a proposer.
func (a *ApaxosDialer) Accept(address string, message *transactions.AcceptMessage) {
	conn, err := a.connect(address)
	if err != nil {
		log.Printf("failed to call %s: %v\n", address, err)
		return
	}
	defer conn.Close()

	_, _ = transactions.NewApaxosClient(conn).Accept(context.Background(), message)
}

// Accepted message is sent by an acceptor.
func (a *ApaxosDialer) Accepted(address string) {
	conn, err := a.connect(address)
	if err != nil {
		log.Printf("failed to call %s: %v\n", address, err)
		return
	}
	defer conn.Close()

	_, _ = transactions.NewApaxosClient(conn).Accepted(context.Background(), &emptypb.Empty{})
}

// Commit message is sent by a proposer.
func (a *ApaxosDialer) Commit(address string) {
	conn, err := a.connect(address)
	if err != nil {
		log.Printf("failed to call %s: %v\n", address, err)
		return
	}
	defer conn.Close()

	_, _ = transactions.NewApaxosClient(conn).Commit(context.Background(), &emptypb.Empty{})
}

// Sync message is sent by a proposer to a felt-behind acceptor
// or by an acceptor to a felt-behind proposer
func (a *ApaxosDialer) Sync(address string, messages []*transactions.SyncMessage) {
	conn, err := a.connect(address)
	if err != nil {
		log.Printf("failed to call %s: %v\n", address, err)
		return
	}
	defer conn.Close()

	stream, err := transactions.NewApaxosClient(conn).Sync(context.Background())
	if err != nil {
		log.Printf("failed to open an stream to %s: %v\n", address, err)
		return
	}

	for _, message := range messages {
		if err := stream.Send(message); err != nil {
			log.Printf("failed to send sync message to %s: %v\n", address, err)
			return
		}
	}

	_, _ = stream.CloseAndRecv()
}
