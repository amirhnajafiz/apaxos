package grpc

import "github.com/f24-cse535/apaxos/pkg/transactions"

// apaxos server handles internal RPC calls for apaxos protocol.
type apaxosServer struct {
	transactions.UnimplementedApaxosServer
}
