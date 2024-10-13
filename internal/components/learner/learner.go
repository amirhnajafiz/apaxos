package learner

import (
	"github.com/f24-cse535/apaxos/internal/storage/local"
	"github.com/f24-cse535/apaxos/pkg/enum"
	"github.com/f24-cse535/apaxos/pkg/messages"
)

// Learner state-machine handles the cases for
// sync messages and commit messages.
type Learner struct {
	Memory  *local.Memory
	Channel chan *messages.Packet
	State   enum.StateType
}

// Signal method is used to send a message to this machine.
func (l Learner) Signal(pkt *messages.Packet) (enum.StateType, error) {
	l.Channel <- pkt

	return l.State, nil
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
			l.sync()
		case enum.PacketCommit: // sent by the proposer
			l.commit()
		default: // drop the message if none
			continue
		}
	}
}
