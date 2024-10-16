package apaxos

import (
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"

	"go.uber.org/zap"
)

// broadcase propose, calls propose RPC on each node.
func (a *Apaxos) broadcastPropose(b *apaxos.BallotNumber) {
	for _, node := range a.Nodes {
		a.Logger.Debug("send prepare message", zap.String("to", node))

		a.Dialer.Propose(node, &apaxos.PrepareMessage{
			NodeId:              a.NodeId,
			BallotNumber:        b,
			LastComittedMessage: a.Memory.GetLastCommittedMessage(),
		})
	}
}

// broadcast accept, calls accept RPC on each node.
func (a *Apaxos) broadcastAccept(b *apaxos.BallotNumber, blocks []*apaxos.Block) {
	for _, node := range a.Nodes {
		a.Logger.Debug("send accept message", zap.String("to", node))

		a.Dialer.Accept(node, &apaxos.AcceptMessage{
			NodeId:       a.NodeId,
			BallotNumber: b,
			Blocks:       blocks,
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

// transmitSync will be called by the proposer handler to update the acceptor.
func (a *Apaxos) transmitSync(address string) {
	// get a clone of the clients
	clients := a.Memory.GetClients()

	// create an instance of sync message
	message := &apaxos.SyncMessage{
		LastComittedMessage: a.Memory.GetLastCommittedMessage(),
		Pairs:               make([]*apaxos.ClientBalancePair, len(clients)),
	}

	// add client and their balances
	index := 0
	for key, value := range clients {
		message.Pairs[index] = &apaxos.ClientBalancePair{
			Client:  key,
			Balance: value,
		}

		index++
	}

	// send the sync message
	a.Dialer.Sync(address, message)
}
