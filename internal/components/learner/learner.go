package learner

import (
	"github.com/f24-cse535/apaxos/pkg/enum"
	"github.com/f24-cse535/apaxos/pkg/messages"
)

// Learner state-machine handles the cases for
// sync messages and commit messages.
type Learner struct {
	Channel chan *messages.Packet
}

// Signal method is used to send a message to this machine.
func (l Learner) Signal(pkt *messages.Packet) {
	l.Channel <- pkt
}

// Start method, the learner waits for messages from the dispatcher.
func (l Learner) Start() {
	for {
		// wait on its notify channel
		pkt := <-l.Channel

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
