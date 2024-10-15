package grpc

import (
	"context"
	"time"

	"github.com/f24-cse535/apaxos/internal/consensus"
	"github.com/f24-cse535/apaxos/internal/monitoring/metrics"
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
	Metrics   *metrics.Metrics
}

// observeMetrics is used in RPCs to set new metrics values.
func (s *transactionsServer) observeMetrics(start, end time.Time) {
	tmp := float64(end.Sub(start).Milliseconds())

	s.Metrics.ObserveLatency(tmp)                    // latency is the time spent for each transaction
	s.Metrics.ObserveThroughput(float64(1000 / tmp)) // throughput is the number of transactions per second
}

// NewTransaction is called for registering a new transaction.
// The handler sends a demand to consensus and waits for a response.
func (s *transactionsServer) NewTransaction(ctx context.Context, req *apaxos.Transaction) (*transactions.TransactionResponse, error) {
	s.Logger.Debug("rpc called NewTransaction")

	// to set system metrics
	start := time.Now()

	// send a message to the consensus module to process a new transaction
	channel, code, err := s.Consensus.Demand(&messages.Packet{
		Type:    enum.PacketTransaction,
		Payload: req,
	})

	// set the first end of response
	end := time.Now()

	s.Logger.Debug("consensus responed", zap.Int("code", code))

	// if the consensus didn't accept our transaction, we return an error
	if code != enum.ResponseAccepted {
		s.observeMetrics(start, end)

		return &transactions.TransactionResponse{
			Status: int64(code),
			Text:   err.Error(),
		}, err
	} else { // else, it means that the consensus module has accepted our transaction
		// now we have to see that we should wait for a response, or the transaction is submitted
		response := transactions.TransactionResponse{
			Status: int64(code),
		}

		// wait for consensus module response
		if code != enum.ResponseOK {
			resp := <-channel

			end = time.Now()
			response.Text = resp.Payload.(string)
		}

		s.observeMetrics(start, end)

		return &response, nil
	}
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

	// call get values on metrics module
	lt, th := s.Metrics.GetValues()

	return &transactions.PerformanceResponse{
		Throughput: th,
		Latency:    lt,
	}, nil
}
