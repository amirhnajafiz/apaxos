package grpc

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
type apaxosServer struct {
	apaxos.UnimplementedApaxosServer
	Consensus *consensus.Consensus
	Logger    *zap.Logger
}

// Propose will be called by the proposer's consensus module and waits for a call on promise.
func (a *apaxosServer) Propose(ctx context.Context, input *apaxos.PrepareMessage) (*emptypb.Empty, error) {
	a.Logger.Debug("rpc called propose", zap.String("caller", input.NodeId))

	a.Consensus.Signal(&messages.Packet{
		Type:    enum.PacketPrepare,
		Payload: input,
	})

	return &emptypb.Empty{}, nil
}

// Promise will be called by the acceptor's consensus module.
func (a *apaxosServer) Promise(ctx context.Context, input *apaxos.PromiseMessage) (*emptypb.Empty, error) {
	a.Logger.Debug("rpc called promise", zap.String("caller", input.NodeId))

	a.Consensus.Signal(&messages.Packet{
		Type:    enum.PacketPromise,
		Payload: input,
	})

	return &emptypb.Empty{}, nil
}

// Accept will be called by the proposer's consensus module and waits for a call on accepted.
func (a *apaxosServer) Accept(ctx context.Context, input *apaxos.AcceptMessage) (*emptypb.Empty, error) {
	a.Logger.Debug("rpc called accept", zap.String("caller", input.NodeId))

	a.Consensus.Signal(&messages.Packet{
		Type:    enum.PacketAccept,
		Payload: input,
	})

	return &emptypb.Empty{}, nil
}

// Accepted will be called by the acceptor's consensus module.
func (a *apaxosServer) Accepted(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	a.Logger.Debug("rpc called accepted")

	a.Consensus.Signal(&messages.Packet{
		Type: enum.PacketAccepted,
	})

	return &emptypb.Empty{}, nil
}

// Commit will be called by the proposer's consensus module.
func (a *apaxosServer) Commit(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	a.Logger.Debug("rpc called commit")

	a.Consensus.Signal(&messages.Packet{
		Type: enum.PacketCommit,
	})

	return &emptypb.Empty{}, nil
}

// Sync will be called by the proposer's or acceptor's consensus module.
// If the proposer is slow, one acceptor will call this sync.
// If the acceptor is slow, the proposer will call this after getting a call on promise.
func (a *apaxosServer) Sync(ctx context.Context, input *apaxos.SyncMessage) (*emptypb.Empty, error) {
	a.Logger.Debug("rpc called sync")

	a.Consensus.Signal(&messages.Packet{
		Type:    enum.PacketSync,
		Payload: input,
	})

	return &emptypb.Empty{}, nil
}
