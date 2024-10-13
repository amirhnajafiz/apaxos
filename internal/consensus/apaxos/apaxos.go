package apaxos

import "github.com/f24-cse535/apaxos/pkg/messages"

type Apaxos struct {
	Channel chan messages.Packet
}

func (a Apaxos) Start() error {
	// Send prepare messages
	// Wait for promise messages (first on majority, then on a timeout)
	// Create a message
	// Send accept messages
	// Wait for accepted messages (first on majority, then on a timeout)
	// Send commit message
	// Clean the dataset

	return nil
}
