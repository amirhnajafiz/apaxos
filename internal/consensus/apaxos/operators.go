package apaxos

import (
	"github.com/f24-cse535/apaxos/pkg/models"
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"

	"go.uber.org/zap"
)

// broadcase propose, calls propose RPC on each node.
func (a *Apaxos) broadcastPropose(b *models.BallotNumber) {
	for _, node := range a.Nodes {
		a.Logger.Debug("send prepare message", zap.String("to", node))

		a.Dialer.Propose(node, &apaxos.PrepareMessage{
			NodeId:              a.NodeId,
			BallotNumber:        b.ToProtoModel(),
			LastComittedMessage: a.Memory.GetLastCommittedMessage().ToProtoModel(),
		})
	}
}

// broadcast accept, calls accept RPC on each node.
func (a *Apaxos) broadcastAccept(b *models.BallotNumber, blocks []*models.Block) {
	// convert models block to apaxos.Block
	list := make([]*apaxos.Block, len(blocks))
	for index, block := range blocks {
		list[index] = block.ToProtoModel()
	}

	for _, node := range a.Nodes {
		a.Logger.Debug("send accept message", zap.String("to", node))

		a.Dialer.Accept(node, &apaxos.AcceptMessage{
			NodeId:       a.NodeId,
			BallotNumber: b.ToProtoModel(),
			Blocks:       list,
		})
	}
}

// broadcast commit, calls commit RPC on each node.
func (a *Apaxos) broadcastCommit() {
	for _, node := range a.Nodes {
		a.Logger.Debug("send commit message", zap.String("to", node))

		a.Dialer.Commit(node)
	}
}
