package apaxos

import (
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"
)

func (a Apaxos) broadcastPropose() {
	for _, node := range a.Nodes {
		a.Dialer.Propose(node, &apaxos.PrepareMessage{
			BallotNumber: a.Memory.GetBallotNumber().ToProtoModel(),
			// LastComittedMessage:
		})
	}
}

func (a Apaxos) broadcastAccept() {
	for _, node := range a.Nodes {
		a.Dialer.Accept(node, &apaxos.AcceptMessage{
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
