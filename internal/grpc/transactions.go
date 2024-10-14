package grpc

import (
	"context"

	"github.com/f24-cse535/apaxos/internal/consensus"
	"github.com/f24-cse535/apaxos/internal/storage/database"
	"github.com/f24-cse535/apaxos/internal/storage/local"
	"github.com/f24-cse535/apaxos/pkg/enum"
	"github.com/f24-cse535/apaxos/pkg/messages"
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"
	"github.com/f24-cse535/apaxos/pkg/rpc/transactions"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

// transactions server handles the clients RPC calls.
type transactionsServer struct {
	transactions.UnimplementedTransactionsServer
	Memory    *local.Memory
	Database  *database.Database
	Consensus *consensus.Consensus
	Logger    *zap.Logger
}

// NewTransaction is called for registering a new transaction.
// The handler sends a demand to consensus and waits for a response.
func (s *transactionsServer) NewTransaction(ctx context.Context, req *apaxos.Transaction) (*transactions.TransactionResponse, error) {
	s.Logger.Debug("rpc called NewTransaction")

	// create a response instance
	response := transactions.TransactionResponse{}

	// send a message to consensus to process a transaction
	channel, err := s.Consensus.Demand(&messages.Packet{
		Type:    enum.PacketTransaction,
		Payload: req,
	})
	if err != nil {
		return &response, err
	}

	if channel == nil {
		// if channel is nil, it means the transaction was successful and no need to wait
		// for consensus protocl
	} else {
		// wait on the consensus response
		pkt := <-channel
	}

	return &response, nil
}

// PrintBalance is a simple operation that reads the client balance from node's memory.
func (s *transactionsServer) PrintBalance(ctx context.Context, req *transactions.PrintBalanceRequest) (*transactions.PrintBalanceResponse, error) {
	s.Logger.Debug("rpc called PrintBalance")

	return &transactions.PrintBalanceResponse{
		Balance: s.Memory.GetBalance(req.Client),
	}, nil
}

// PrintLogs returns the node datastore and accepted val.
func (s *transactionsServer) PrintLogs(req *emptypb.Empty, stream transactions.Transactions_PrintLogsServer) error {
	s.Logger.Debug("rpc called PrintLogs")

	// first send the datastore block
	if err := stream.Send(s.Memory.GetDatastore().ToProtoModel()); err != nil {
		return err
	}

	// send accepted_val blocks
	for _, block := range s.Memory.GetAcceptedVal() {
		if err := stream.Send(block.ToProtoModel()); err != nil {
			return err
		}
	}

	return nil
}

// PrintDB get's blocks from MongoDB and sends them as proto blocks.
func (s *transactionsServer) PrintDB(req *emptypb.Empty, stream transactions.Transactions_PrintDBServer) error {
	s.Logger.Debug("rpc called PrintDB")

	// get all blocks from MongoDB
	blocks, err := s.Database.GetBlocks()
	if err != nil {
		return err
	}

	// send them one by one
	for _, block := range blocks {
		if err := stream.Send(block.ToProtoModel()); err != nil {
			return err
		}
	}

	return stream.RecvMsg(nil)
}

// Performance function returns the node's throughput and latency.
func (s *transactionsServer) Performance(ctx context.Context, req *emptypb.Empty) (*transactions.PerformanceResponse, error) {
	s.Logger.Debug("rpc called Performance")

	return &transactions.PerformanceResponse{
		Throughput: 1000.5,
		Latency:    20.3,
	}, nil
}
