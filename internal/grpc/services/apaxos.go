package services

import (
	"context"

	"github.com/f24-cse535/apaxos/internal/consensus"
	"github.com/f24-cse535/apaxos/pkg/enum"
	"github.com/f24-cse535/apaxos/pkg/messages"
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

// apaxos server handles internal RPC calls for apaxos protocol.
type Apaxos struct {
	apaxos.UnimplementedApaxosServer
	Consensus *consensus.Consensus
	Logger    *zap.Logger
}

// Propose will be called by the proposer's consensus module and waits for a call on promise.
func (a *Apaxos) Propose(ctx context.Context, input *apaxos.PrepareMessage) (*emptypb.Empty, error) {
	a.Logger.Debug("rpc called propose", zap.String("caller", input.NodeId))

	// call prepare method of the consensus module
	go a.Consensus.Prepare(input)

	return &emptypb.Empty{}, nil
}

// Promise will be called by the acceptor's consensus module.
func (a *Apaxos) Promise(ctx context.Context, input *apaxos.PromiseMessage) (*emptypb.Empty, error) {
	a.Logger.Debug("rpc called promise", zap.String("caller", input.NodeId))

	// send a signal to consensus module
	a.Consensus.Signal(&messages.Packet{
		Type:    enum.PacketPromise,
		Payload: input,
	})

	return &emptypb.Empty{}, nil
}

// Accept will be called by the proposer's consensus module and waits for a call on accepted.
func (a *Apaxos) Accept(ctx context.Context, input *apaxos.AcceptMessage) (*emptypb.Empty, error) {
	a.Logger.Debug("rpc called accept", zap.String("caller", input.NodeId))

	// call accept method of the consensus module
	go a.Consensus.Accept(input)

	return &emptypb.Empty{}, nil
}

// Accepted will be called by the acceptor's consensus module.
func (a *Apaxos) Accepted(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	a.Logger.Debug("rpc called accepted")

	// send an accepted signal to the consensus module
	a.Consensus.Signal(&messages.Packet{
		Type: enum.PacketAccepted,
	})

	return &emptypb.Empty{}, nil
}

// Commit will be called by the proposer's consensus module.
func (a *Apaxos) Commit(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	a.Logger.Debug("rpc called commit")

	// call commit method of the consensus module
	go a.Consensus.Commit()

	return &emptypb.Empty{}, nil
}

// Sync will be called by the proposer's or acceptor's consensus module.
func (a *Apaxos) Sync(ctx context.Context, input *apaxos.SyncMessage) (*emptypb.Empty, error) {
	a.Logger.Debug("rpc called sync")

	// call sync method of the consensus module
	go a.Consensus.Sync(input)

	return &emptypb.Empty{}, nil
}
