package grpc

import "github.com/f24-cse535/apaxos/pkg/transactions"

type apaxosServer struct {
	transactions.UnimplementedApaxosServer
}
