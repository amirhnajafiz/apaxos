package client

import "go.uber.org/zap"

// Client has all RPCs to communicate with the gRPC server.
type Client struct {
	ApaxosDialer
	LivenessDialer
	TransactionsDialer
}

// NewClient returns a new RPC client to make RPC to the gRPC server.
func NewClient(logr *zap.Logger) *Client {
	return &Client{
		ApaxosDialer: ApaxosDialer{
			Logger: logr.Named("apaxos"),
		},
		TransactionsDialer: TransactionsDialer{
			Logger: logr.Named("transactions"),
		},
		LivenessDialer: LivenessDialer{},
	}
}
