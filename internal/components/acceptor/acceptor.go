package acceptor

import (
	"github.com/f24-cse535/apaxos/pkg/enum"
	"github.com/f24-cse535/apaxos/pkg/messages"
)

// Acceptor state-machine handles the cases for
// prepare messages and accept messages.
type Acceptor struct {
	Channel chan *messages.Packet
}

// Signal method is used to send a message to this machine.
func (a Acceptor) Signal(pkt *messages.Packet) {
	a.Channel <- pkt
}

// Start method, the machine waits for messages from the dispatcher.
func (a Acceptor) Start() {
	for {
		// wait on its input notify channel
		pkt := <-a.Channel

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
