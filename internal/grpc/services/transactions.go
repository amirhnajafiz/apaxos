package services

import (
	"context"
	"time"

	"github.com/f24-cse535/apaxos/internal/consensus"
	"github.com/f24-cse535/apaxos/internal/monitoring/metrics"
	"github.com/f24-cse535/apaxos/internal/storage/database"
	"github.com/f24-cse535/apaxos/internal/storage/local"
	"github.com/f24-cse535/apaxos/pkg/messages"
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"
	"github.com/f24-cse535/apaxos/pkg/rpc/transactions"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

// transactions server handles the clients RPC calls.
type Transactions struct {
	transactions.UnimplementedTransactionsServer
	Memory    *local.Memory
	Database  *database.Database
	Consensus *consensus.Consensus
	Logger    *zap.Logger
	Metrics   *metrics.Metrics
}

// observeMetrics is used in RPCs to set new metrics values.
func (s *Transactions) observeMetrics(duration time.Duration) {
	tmp := duration.Microseconds()

	if tmp == 0 {
		s.Metrics.ObserveLatency(0)
		s.Metrics.ObserveThroughput(1000000)
	} else {
		s.Metrics.ObserveLatency(float64(tmp))              // latency is the time spent for each transaction
		s.Metrics.ObserveThroughput(float64(1000000 / tmp)) // throughput is the number of transactions per second
	}
}

// NewTransaction is called for registering a new transaction.
// The handler sends a demand to consensus and waits for a response.
func (s *Transactions) NewTransaction(ctx context.Context, req *apaxos.Transaction) (*transactions.TransactionResponse, error) {
	s.Logger.Debug(
		"rpc called NewTransaction",
		zap.String("sender", req.GetSender()),
		zap.String("receiver", req.GetReciever()),
		zap.Int64("amount", req.GetAmount()),
	)

	// to set system metrics
	start := time.Now()

	// send a message to the consensus module to process a new transaction
	channel, err := s.Consensus.Checkout(&messages.Packet{
		Payload: req,
	})
	if err != nil {
		return nil, err
	}

	elapsed := time.Since(start) // calculate the elapsed time

	// create a response instance
	response := transactions.TransactionResponse{
		Text: "transaction submitted",
	}

	// if we got a channel, we should wait for consensus module response over it
	if channel != nil {
		resp := <-channel

		elapsed = time.Since(start) // calculate the elapsed time

		response.Text = resp.Payload.(string)
	}

	s.observeMetrics(elapsed)

	return &response, nil
}

// PrintBalance is a simple operation that reads the client balance from node's memory.
func (s *Transactions) PrintBalance(ctx context.Context, req *transactions.PrintBalanceRequest) (*transactions.PrintBalanceResponse, error) {
	s.Logger.Debug("rpc called PrintBalance")

	return &transactions.PrintBalanceResponse{
		Balance: s.Memory.GetBalance(req.Client),
	}, nil
}

// PrintLogs returns the node datastore and accepted val.
func (s *Transactions) PrintLogs(req *emptypb.Empty, stream transactions.Transactions_PrintLogsServer) error {
	s.Logger.Debug("rpc called PrintLogs")

	// first send the datastore block
	if err := stream.Send(s.Memory.GetDatastore()); err != nil {
		return err
	}

	// send accepted_val blocks
	for _, block := range s.Memory.GetAcceptedVal() {
		if err := stream.Send(block); err != nil {
			return err
		}
	}

	return nil
}

// PrintDB get's blocks from MongoDB and sends them as proto blocks.
func (s *Transactions) PrintDB(req *emptypb.Empty, stream transactions.Transactions_PrintDBServer) error {
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

	return nil
}

// Performance function returns the node's throughput and latency.
func (s *Transactions) Performance(ctx context.Context, req *emptypb.Empty) (*transactions.PerformanceResponse, error) {
	s.Logger.Debug("rpc called Performance")

	// call get values on metrics module
	lt, th := s.Metrics.GetValues()

	return &transactions.PerformanceResponse{
		Throughput: th,
		Latency:    lt,
	}, nil
}
