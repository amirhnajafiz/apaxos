package proposer

import (
	"github.com/f24-cse535/apaxos/internal/storage/local"
	"github.com/f24-cse535/apaxos/pkg/enum"
	"github.com/f24-cse535/apaxos/pkg/messages"
)

// Proposer state-machine handles the cases for
// propose messages, promise messages, and accepted messages.
type Proposer struct {
	Memory  *local.Memory
	Channel chan *messages.Packet
}

// Signal method is used to send a message to this machine.
func (p Proposer) Signal(pkt *messages.Packet) {
	p.Channel <- pkt
}

// Start method, the learner waits for messages from the dispatcher.
func (p Proposer) Start() {
	for {
		// wait on its notify channel
		pkt := <-p.Channel

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
