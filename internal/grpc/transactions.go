package grpc

import "github.com/f24-cse535/apaxos/pkg/transactions"

// transactions server handles the clients RPC calls.
type transactionsServer struct {
	transactions.UnimplementedTransactionsServer
}
