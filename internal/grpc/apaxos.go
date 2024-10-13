package grpc

import (
	"context"

	"github.com/f24-cse535/apaxos/pkg/transactions"

	"google.golang.org/protobuf/types/known/emptypb"
)

// apaxos server handles internal RPC calls for apaxos protocol.
type apaxosServer struct {
	transactions.UnimplementedApaxosServer
}

func (a *apaxosServer) Propose(ctx context.Context, input *transactions.BallotNumber) (emptypb.Empty, error) {

}

func (a *apaxosServer) Promise(ctx context.Context, input *transactions.PromiseMessage) (emptypb.Empty, error) {

}

func (a *apaxosServer) Accept(ctx context.Context, input *transactions.AcceptMessage) (emptypb.Empty, error) {

}

func (a *apaxosServer) Accepted(ctx context.Context, _ *emptypb.Empty) (emptypb.Empty, error) {

}

func (a *apaxosServer) Commit(ctx context.Context, _ *emptypb.Empty) (emptypb.Empty, error) {

}

func (a *apaxosServer) Sync(stream *transactions.SyncMessage) (*emptypb.Empty, error) {

}
