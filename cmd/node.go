package cmd

import (
	"github.com/f24-cse535/apaxos/internal/config"
	"github.com/f24-cse535/apaxos/internal/consensus"
	"github.com/f24-cse535/apaxos/internal/grpc"
	"github.com/f24-cse535/apaxos/internal/grpc/client"
	"github.com/f24-cse535/apaxos/internal/storage/database"
	"github.com/f24-cse535/apaxos/internal/storage/local"
	"github.com/f24-cse535/apaxos/internal/worker"

	"go.uber.org/zap"
)

// Each node of our transaction system runs using this main function.
type Node struct {
	Cfg    config.Config
	Logger *zap.Logger
}

func (n Node) Main() error {
	// open database connection
	db, err := database.NewDatabase(n.Cfg.MongoDB, n.Cfg.NodeID)
	if err != nil {
		return err
	}

	// create a local storage (aka memory)
	mem := local.NewMemory(n.Cfg.NodeID, n.Cfg.GetBalances())

	// check for previous state (aka snapshot)
	ss, err := db.GetLastState()
	if err != nil && ss != nil { // if ss exists, read from the previous state
		mem.ReadFromState(ss)
	}

	// create a new consensus module
	instance := consensus.Consensus{
		Database:        db,
		Memory:          mem,
		Client:          n.Cfg.Client,
		NodeId:          n.Cfg.NodeID,
		Nodes:           n.Cfg.GetNodes(),
		RequestTimeout:  n.Cfg.GRPC.RequestTimeout,
		MajorityTimeout: n.Cfg.GRPC.MajorityTimeout,
		Majority:        n.Cfg.Majority,
		Logger:          n.Logger.Named("consensus"),
		LivenessDialer:  &client.LivenessDialer{},
		Dialer: &client.ApaxosDialer{
			Logger: n.Logger.Named("apaxos-dialer"),
		},
	}

	// create a worker instance and execute it in a new sub-process
	go worker.Worker{
		Memory:   mem,
		Database: db,
		Interval: n.Cfg.WorkersInterval,
		Logger:   n.Logger.Named("worker"),
	}.Start()

	// create a new gRPC bootstrap instance and execute the server by running the boot commands
	boot := grpc.Bootstrap{
		Port:      n.Cfg.GRPC.Port,
		Memory:    mem,
		Database:  db,
		Consensus: &instance,
		Logger:    n.Logger.Named("grpc"),
	}

	return boot.ListenAnsServer()
}
