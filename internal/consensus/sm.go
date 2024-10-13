package consensus

import (
	"github.com/f24-cse535/apaxos/pkg/enum"
	"github.com/f24-cse535/apaxos/pkg/messages"
)

// StateMachine is an interface that represents
// each of the components in the system.
type StateMachine interface {
	Start()
	Signal(pkt *messages.Packet) (enum.StateType, error)
}
