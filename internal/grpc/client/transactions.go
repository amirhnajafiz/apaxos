package client

import (
	"context"
	"fmt"
	"io"

	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"
	"github.com/f24-cse535/apaxos/pkg/rpc/transactions"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

// TransactionsDialer is used to call RPCs for transactions service.
type TransactionsDialer struct {
	Logger *zap.Logger
}

// connect should be called in the beginning of each method to establish a connection.
func (t *TransactionsDialer) connect(address string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("[grpc/client/transactionsdialer] failed to open connectio to %s: %v", address, err)
	}

	return conn, nil
}

// NewTransaction is used by the clients to submit a new transaction.
func (t *TransactionsDialer) NewTransaction(address string, instance *apaxos.Transaction) (string, error) {
	// base connection
	conn, err := t.connect(address)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// call NewTransaction RPC and get the response
	resp, err := transactions.NewTransactionsClient(conn).NewTransaction(context.Background(), instance)
	if err != nil {
		return "", fmt.Errorf("failed transaction: %v", err)
	}

	return resp.GetText(), nil
}

// PrintBalance is used for getting a client balance inside a specific node.
func (t *TransactionsDialer) PrintBalance(address string, client string) (int64, error) {
	// base connection
	conn, err := t.connect(address)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	// call PrintBalance with the client in the request, and wait for response
	resp, err := transactions.NewTransactionsClient(conn).PrintBalance(context.Background(), &transactions.PrintBalanceRequest{
		Client: client,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to process printBalance: %v", err)
	}

	return resp.GetBalance(), nil
}

// PrintLogs is used to get a specific node datastore.
func (t *TransactionsDialer) PrintLogs(address string) ([]*apaxos.Block, error) {
	// base connection
	conn, err := t.connect(address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// open an stream on PrintLogs RPC to get blocks
	stream, err := transactions.NewTransactionsClient(conn).PrintLogs(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("failed to process printLogs: %v", err)
	}

	// create a list to store blocks
	list := make([]*apaxos.Block, 0)

	for {
		// get items one by one
		in, err := stream.Recv()
		if err != nil {
			if err == io.EOF { // send a response once the stream is closed
				return list, nil
			}

			return nil, fmt.Errorf("failed to receive blocks: %v", err)
		}

		// append to the list of blocks
		list = append(list, in)
	}
}

// PrintDB is used to get stored blocks of a specific node.
func (t *TransactionsDialer) PrintDB(address string) ([]*apaxos.Block, error) {
	// base connection
	conn, err := t.connect(address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// open a stream on PrintDB to get blocks
	stream, err := transactions.NewTransactionsClient(conn).PrintDB(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("failed to process printDB: %v", err)
	}

	// create a list to store blocks
	list := make([]*apaxos.Block, 0)

	for {
		// read blocks one by one
		in, err := stream.Recv()
		if err != nil {
			if err == io.EOF { // send a response once the stream is closed
				return list, nil
			}

			return nil, fmt.Errorf("failed to receive log: %v", err)
		}

		// append to the list of blocks
		list = append(list, in)
	}
}

// Performance is used to get the performance of a specific node.
func (t *TransactionsDialer) Performance(address string) (*transactions.PerformanceResponse, error) {
	// base connection
	conn, err := t.connect(address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// call Performance RPC call on the target to get a performance response
	resp, err := transactions.NewTransactionsClient(conn).Performance(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("failed to process performance: %v", err)
	}

	return resp, nil
}
