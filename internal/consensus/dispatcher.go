package consensus

import (
	"github.com/f24-cse535/apaxos/pkg/messages"
)

// Dispatcher is the core component of our distributed system.
// the role of this component is receive messages Apaxos messages from the gRPC
// server on the node and sends it to other components.
type Dispatcher struct {
	machins     []StateMachine
	GRPCChannel chan *messages.Packet
}

// NewDispatcher first initialized the components, then it will return a
// dispatcher instance.
func NewDispatcher() *Dispatcher {
	// create a new dispatcher instance
	instance := Dispatcher{}

	// creating each of the three acceptor, learner, and proposer components.

	return &instance
}

func (d Dispatcher) ListenAndServe() {
	// the first job of the dispatcher is to start 3 state-machines
	for _, machine := range d.machins {
		go machine.Start()
	}

	// then it waits for input events from gRPC server
	for {
		pkt := <-d.GRPCChannel

		// after that it will broadcast the packet to all state-machines
		for _, machine := range d.machins {
			machine.Signal(pkt)
		}
	}
}
