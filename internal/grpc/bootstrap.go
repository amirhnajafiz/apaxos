package grpc

import (
	"fmt"
	"net"

	"github.com/f24-cse535/apaxos/internal/consensus"
	"github.com/f24-cse535/apaxos/internal/storage/database"
	"github.com/f24-cse535/apaxos/internal/storage/local"
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"
	"github.com/f24-cse535/apaxos/pkg/rpc/transactions"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Bootstrap is a wrapper that holds
// every required thing for the gRPC server starting.
type Bootstrap struct {
	Port      int
	Memory    *local.Memory
	Database  *database.Database
	Consensus *consensus.Consensus
	Logger    *zap.Logger
}

// ListenAnsServer creates a new gRPC instance
// and registers both apaxos and transactions servers.
func (b Bootstrap) ListenAnsServer() error {
	// on the local network, listen to a port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", b.Port))
	if err != nil {
		return fmt.Errorf("[grcp] failed to start the listener server: %v", err)
	}

	// create a new grpc instance
	server := grpc.NewServer()

	// register both gRPC services
	apaxos.RegisterApaxosServer(server, &apaxosServer{
		Consensus: b.Consensus,
		Logger:    b.Logger.Named("apaxos"),
	})
	transactions.RegisterTransactionsServer(server, &transactionsServer{
		Consensus: b.Consensus,
		Memory:    b.Memory,
		Database:  b.Database,
		Logger:    b.Logger.Named("transactions"),
	})

	// starting the server
	b.Logger.Info("grpc server started", zap.Int("port", b.Port))
	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("[grpc] failed to start the server: %v", err)
	}

	return nil
}
