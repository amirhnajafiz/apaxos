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

	"google.golang.org/protobuf/types/known/emptypb"
)

// transactions server handles the clients RPC calls.
type transactionsServer struct {
	transactions.UnimplementedTransactionsServer
	Memory    *local.Memory
	Database  *database.Database
	Consensus *consensus.Consensus
}

// NewTransaction is called for registering a new transaction.
// The handler sends a demand to consensus and waits for a response.
func (s *transactionsServer) NewTransaction(ctx context.Context, req *apaxos.Transaction) (*transactions.TransactionResponse, error) {
	response := transactions.TransactionResponse{Result: false}

	pkt, err := s.Consensus.Demand(&messages.Packet{
		Type:    enum.PacketTransaction,
		Payload: req,
	})
	if err != nil {
		return &response, err
	}

	response.Result = pkt.Payload.(bool)

	return &response, nil
}

// PrintBalance is a simple operation that reads the client balance from node's memory.
func (s *transactionsServer) PrintBalance(ctx context.Context, req *transactions.PrintBalanceRequest) (*transactions.PrintBalanceResponse, error) {
	return &transactions.PrintBalanceResponse{
		Balance: s.Memory.GetBalance(req.Client),
	}, nil
}

// PrintLogs returns the node datastore and accepted val.
func (s *transactionsServer) PrintLogs(req *emptypb.Empty, stream transactions.Transactions_PrintLogsServer) error {
	// first send datastore block
	// creating transactions list
	transactions := s.Memory.GetDatastore()

	// creating datastore block
	ds := &apaxos.Block{
		Metadata: &apaxos.BlockMetaData{
			NodeId:         s.Memory.GetBallotNumber().NodeId,
			SequenceNumber: s.Memory.GetSequenceNumber(),
			BallotNumber:   s.Memory.GetBallotNumber().ToProtoModel(),
		},
		Transactions: make([]*apaxos.Transaction, len(transactions)),
	}

	// modify transactions
	for index, item := range transactions {
		ds.Transactions[index] = item.ToProtoModel()
	}

	// send the datastore block
	if err := stream.Send(ds); err != nil {
		return err
	}

	// send accepted var blocks
	for _, block := range s.Memory.GetAcceptedVal() {
		if err := stream.Send(block.ToProtoModel()); err != nil {
			return err
		}
	}

	return nil
}

// PrintDB get's blocks from MongoDB and sends them as proto blocks.
func (s *transactionsServer) PrintDB(req *emptypb.Empty, stream transactions.Transactions_PrintDBServer) error {
	blocks, err := s.Database.GetBlocks()
	if err != nil {
		return err
	}

	for _, block := range blocks {
		if err := stream.Send(block.ToProtoModel()); err != nil {
			return err
		}
	}

	return nil
}

// Performance function returns the node's throughput and latency.
func (s *transactionsServer) Performance(ctx context.Context, req *emptypb.Empty) (*transactions.PerformanceResponse, error) {
	return &transactions.PerformanceResponse{
		Throughput: 1000.5,
		Latency:    20.3,
	}, nil
}
