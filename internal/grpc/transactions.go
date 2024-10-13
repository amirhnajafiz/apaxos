package grpc

import (
	"context"
	"log"
	"time"

	"github.com/f24-cse535/apaxos/internal/consensus"
	"github.com/f24-cse535/apaxos/internal/storage/database"
	"github.com/f24-cse535/apaxos/internal/storage/local"
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"
	"github.com/f24-cse535/apaxos/pkg/rpc/transactions"

	"google.golang.org/protobuf/types/known/emptypb"
)

// transactions server handles the clients RPC calls.
type transactionsServer struct {
	transactions.UnimplementedTransactionsServer
	Memory    *local.Memory
	Database  *database.Database
	Consensus *consensus.Consensus
}

func (s *transactionsServer) NewTransaction(ctx context.Context, req *apaxos.Transaction) (*transactions.TransactionResponse, error) {
	log.Printf("NewTransaction called with: %+v", req)

	// Placeholder logic - always return success
	return &transactions.TransactionResponse{Result: true}, nil
}

// Implement PrintBalance RPC
func (s *transactionsServer) PrintBalance(ctx context.Context, req *transactions.PrintBalanceRequest) (*transactions.PrintBalanceResponse, error) {
	log.Println("PrintBalance called")

	// Placeholder logic - return some balance
	return &transactions.PrintBalanceResponse{Balance: 1000}, nil
}

// Implement PrintLogs RPC - server-side streaming
func (s *transactionsServer) PrintLogs(req *emptypb.Empty, stream transactions.Transactions_PrintLogsServer) error {
	log.Println("PrintLogs called")

	// Placeholder logic - stream 3 sample blocks
	for i := 0; i < 3; i++ {
		block := &apaxos.Block{
			Metadata: &apaxos.BlockMetaData{
				Uid:            "block-uid",
				NodeId:         "node-123",
				SequenceNumber: int64(i),
				BallotNumber:   &apaxos.BallotNumber{Number: int64(i), NodeId: "node-123"},
			},
		}
		if err := stream.Send(block); err != nil {
			return err
		}
		time.Sleep(time.Second) // Simulate processing delay
	}
	return nil
}

// Implement PrintDB RPC - server-side streaming
func (s *transactionsServer) PrintDB(req *emptypb.Empty, stream transactions.Transactions_PrintDBServer) error {
	log.Println("PrintDB called")

	// Placeholder logic - stream 3 sample blocks
	for i := 0; i < 3; i++ {
		block := &apaxos.Block{
			Metadata: &apaxos.BlockMetaData{
				Uid:            "block-uid",
				NodeId:         "node-123",
				SequenceNumber: int64(i),
				BallotNumber:   &apaxos.BallotNumber{Number: int64(i), NodeId: "node-123"},
			},
		}
		if err := stream.Send(block); err != nil {
			return err
		}
		time.Sleep(time.Second) // Simulate processing delay
	}
	return nil
}

// Implement Performance RPC
func (s *transactionsServer) Performance(ctx context.Context, req *emptypb.Empty) (*transactions.PerformanceResponse, error) {
	log.Println("Performance called")

	// Placeholder logic - return some performance metrics
	return &transactions.PerformanceResponse{
		Throughput: 1000.5,
		Latency:    20.3,
	}, nil
}

// Implement AggregatedBalance RPC - client-side streaming
func (s *transactionsServer) AggregatedBalance(req *transactions.AggregatedBalanceRequest, stream transactions.Transactions_AggregatedBalanceServer) error {
	log.Printf("AggregatedBalance called for client: %s", req.Client)

	// Placeholder logic - stream 3 aggregated balance responses
	for i := 0; i < 3; i++ {
		balanceResponse := &transactions.AggregatedBalanceResponse{
			Client:  req.Client,
			NodeId:  "node-123",
			Balance: int64(1000 + i),
		}
		if err := stream.Send(balanceResponse); err != nil {
			return err
		}
		time.Sleep(time.Second) // Simulate processing delay
	}
	return nil
}
