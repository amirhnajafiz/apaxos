package grpc

import (
	"context"

	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"

	"google.golang.org/protobuf/types/known/emptypb"
)

// apaxos server handles internal RPC calls for apaxos protocol.
type apaxosServer struct {
	apaxos.UnimplementedApaxosServer
}

func (a *apaxosServer) Propose(ctx context.Context, input *apaxos.PrepareMessage) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (a *apaxosServer) Promise(ctx context.Context, input *apaxos.PromiseMessage) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (a *apaxosServer) Accept(ctx context.Context, input *apaxos.AcceptMessage) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (a *apaxosServer) Accepted(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (a *apaxosServer) Commit(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (a *apaxosServer) Sync(stream apaxos.Apaxos_SyncServer) error {
	return nil
}
