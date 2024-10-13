package apaxos

import (
	"github.com/f24-cse535/apaxos/pkg/transactions"
)

func (a Apaxos) broadcastPropose() {
	for _, node := range a.Nodes {
		a.Dialer.Propose(node, &transactions.PrepareMessage{
			BallotNumber: a.Memory.GetBallotNumber().ToProtoModel(),
			// LastComittedMessage:
		})
	}
}

func (a Apaxos) broadcastAccept() {
	for _, node := range a.Nodes {
		a.Dialer.Accept(node, &transactions.AcceptMessage{
			BallotNumber: a.Memory.GetBallotNumber().ToProtoModel(),
			// Blocks:
		})
	}
}

func (a Apaxos) broadcastCommit() {
	for _, node := range a.Nodes {
		a.Dialer.Commit(node)
	}
}
