package acceptor

import (
	"github.com/f24-cse535/apaxos/pkg/enum"
	"github.com/f24-cse535/apaxos/pkg/messages"
)

type Acceptor struct {
	GRPCChannel chan *messages.Packet
}

func (a Acceptor) Start() {
	// on start method, the acceptor waits for messages from the gRPC server.
	for {
		// wait on gRPC notify channel
		pkt := <-a.GRPCChannel

		// a switch case for pkt type
		// the acceptor only get's prepare and accept.
		switch pkt.Type {
		case enum.PacketPrepare: // sent by the proposer
			return
		case enum.PacketAccept: // sent by the proposer
			return
		default: // drop the message if none
			return
		}
	}
}
