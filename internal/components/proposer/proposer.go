package proposer

import (
	"github.com/f24-cse535/apaxos/internal/storage/local"
	"github.com/f24-cse535/apaxos/pkg/enum"
	"github.com/f24-cse535/apaxos/pkg/messages"
)

type Proposer struct {
	Memory      *local.Memory
	GRPCChannel chan *messages.Packet
}

func (p Proposer) Start() {
	// on start method, the proposer waits for messages from the gRPC server.
	// if the gRPC server signals a APAXOS start, then the proposer starts APAXOS.
	for {
		// wait on gRPC notify channel
		pkt := <-p.GRPCChannel

		// a switch case for pkt type
		// the proposer only get's propose, promise, and accepted.
		switch pkt.Type {
		case enum.PacketPropose: // sent by the gRPC server
			return
		case enum.PacketPromise: // sent by the acceptors
			return
		case enum.PacketAccepted: // sent by the accptors
			return
		default: // drop the message if none
			return
		}
	}
}
