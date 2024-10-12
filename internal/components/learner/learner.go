package learner

import (
	"github.com/f24-cse535/apaxos/pkg/enum"
	"github.com/f24-cse535/apaxos/pkg/messages"
)

type Learner struct {
	GRPCChannel chan *messages.Packet
}

func (l Learner) Start() {
	// on start method, the learner waits for messages from the gRPC server.
	for {
		// wait on gRPC notify channel
		pkt := <-l.GRPCChannel

		// a switch case for pkt type
		// the learner only get's sync and commit.
		switch pkt.Type {
		case enum.PacketSync: // sent by the proposer
			return
		case enum.PacketCommit: // sent by the proposer
			return
		default: // drop the message if none
			return
		}
	}
}
