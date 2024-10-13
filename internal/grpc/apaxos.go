package grpc

import (
	"context"
	"io"

	"github.com/f24-cse535/apaxos/internal/consensus"
	"github.com/f24-cse535/apaxos/pkg/enum"
	"github.com/f24-cse535/apaxos/pkg/messages"
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"

	"google.golang.org/protobuf/types/known/emptypb"
)

// apaxos server handles internal RPC calls for apaxos protocol.
type apaxosServer struct {
	apaxos.UnimplementedApaxosServer
	Consensus *consensus.Consensus
}

func (a *apaxosServer) Propose(ctx context.Context, input *apaxos.PrepareMessage) (*emptypb.Empty, error) {
	a.Consensus.Signal(&messages.Packet{
		Type:    enum.PacketPrepare,
		Payload: input,
	})

	return &emptypb.Empty{}, nil
}

func (a *apaxosServer) Promise(ctx context.Context, input *apaxos.PromiseMessage) (*emptypb.Empty, error) {
	a.Consensus.Signal(&messages.Packet{
		Type:    enum.PacketPromise,
		Payload: input,
	})

	return &emptypb.Empty{}, nil
}

func (a *apaxosServer) Accept(ctx context.Context, input *apaxos.AcceptMessage) (*emptypb.Empty, error) {
	a.Consensus.Signal(&messages.Packet{
		Type:    enum.PacketAccept,
		Payload: input,
	})

	return &emptypb.Empty{}, nil
}

func (a *apaxosServer) Accepted(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	a.Consensus.Signal(&messages.Packet{
		Type: enum.PacketAccepted,
	})

	return &emptypb.Empty{}, nil
}

func (a *apaxosServer) Commit(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	a.Consensus.Signal(&messages.Packet{
		Type: enum.PacketCommit,
	})

	return &emptypb.Empty{}, nil
}

func (a *apaxosServer) Sync(stream apaxos.Apaxos_SyncServer) error {
	sync := make(map[string]int64)

	for {
		in, err := stream.Recv()
		if err != nil {
			if err == io.EOF { // send a response once the stream is closed
				a.Consensus.Signal(&messages.Packet{
					Type:    enum.PacketSync,
					Payload: sync,
				})

				return stream.SendAndClose(&emptypb.Empty{})
			}

			return err
		}

		sync[in.Client] = in.Balance
	}
}
