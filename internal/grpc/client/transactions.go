package client

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"
	"github.com/f24-cse535/apaxos/pkg/rpc/transactions"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// TransactionsDialer is used to call RPCs for transactions service.
type TransactionsDialer struct{}

// connect should be called in the beginning of each method to establish a connection.
func (t *TransactionsDialer) connect(address string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption

	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, fmt.Errorf("[grpc/client/transactionsdialer] failed to make rpc call: %v", err)
	}

	return conn, nil
}

// NewTransaction is used by the clients to submit a new transaction.
func (t *TransactionsDialer) NewTransaction(address string, instance *apaxos.Transaction) error {
	conn, err := t.connect(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	resp, err := transactions.NewTransactionsClient(conn).NewTransaction(context.Background(), instance)
	if err != nil {
		return fmt.Errorf("failed to process transaction: %v", err)
	}

	if !resp.GetResult() {
		return fmt.Errorf("cannot perform this transactions: %t", resp.GetResult())
	}

	return nil
}

// PrintBalance is used for getting a client balance inside a specific node.
func (t *TransactionsDialer) PrintBalance(address string, client string) (int64, error) {
	conn, err := t.connect(address)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	resp, err := transactions.NewTransactionsClient(conn).PrintBalance(context.Background(), &transactions.PrintBalanceRequest{
		Client: client,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to process printbalance: %v", err)
	}

	return resp.GetBalance(), nil
}

// PrintLogs is used to get a specific node datastore.
func (t *TransactionsDialer) PrintLogs(address string) ([]*apaxos.Block, error) {
	conn, err := t.connect(address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	stream, err := transactions.NewTransactionsClient(conn).PrintLogs(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("failed to process printlogs: %v", err)
	}

	list := make([]*apaxos.Block, 0)

	for {
		in, err := stream.Recv()
		if err != nil {
			if err == io.EOF { // send a response once the stream is closed
				return list, nil
			}

			return nil, fmt.Errorf("failed to receive log: %v", err)
		}

		list = append(list, in)
	}
}

// PrintDB is used to get stored blocks of a specific node.
func (t *TransactionsDialer) PrintDB(address string) ([]*apaxos.Block, error) {
	conn, err := t.connect(address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	stream, err := transactions.NewTransactionsClient(conn).PrintDB(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("failed to process printdb: %v", err)
	}

	list := make([]*apaxos.Block, 0)

	for {
		in, err := stream.Recv()
		if err != nil {
			if err == io.EOF { // send a response once the stream is closed
				return list, nil
			}

			return nil, fmt.Errorf("failed to receive log: %v", err)
		}

		list = append(list, in)
	}
}

// Performance is used to get the performance of a specific node.
func (t *TransactionsDialer) Performance(address string) (*transactions.PerformanceResponse, error) {
	conn, err := t.connect(address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	resp, err := transactions.NewTransactionsClient(conn).Performance(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("failed to process performance: %v", err)
	}

	return resp, nil
}

// AggregatedBalance will run printbalance function on the given addresses.
func (t *TransactionsDialer) AggregatedBalance(client string, addresses ...string) map[string]int64 {
	balances := make(map[string]int64)

	for _, address := range addresses {
		balance, err := t.PrintBalance(address, client)
		if err != nil {
			log.Println(fmt.Errorf("failed to get the balance from %s: %v", address, err))
		}

		balances[address] = balance
	}

	return balances
}
