package client

import (
	"fmt"

	"google.golang.org/grpc"
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

func (t *TransactionsDialer) NewTransaction() {}

func (t *TransactionsDialer) PrintBalance() {}

func (t *TransactionsDialer) PrintLogs() {}

func (t *TransactionsDialer) PrintDB() {}

func (t *TransactionsDialer) Performance() {}

func (t *TransactionsDialer) AggregatedBalance() {}
