package apaxos

import (
	"github.com/f24-cse535/apaxos/internal/grpc/client"
	"github.com/f24-cse535/apaxos/internal/storage/local"
	"github.com/f24-cse535/apaxos/pkg/messages"
)

// Apaxos module is used by the consensus module.
// In consensus module we create instances of this Apaxos.
// In order to communicate with this module, there are two
// channels. The InChannel which the data from consensus will be received.
// And, the OutChannel which the data will be sent to consensus.
type Apaxos struct {
	Dialer client.ApaxosDialer
	Memory *local.Memory

	Nodes           map[string]string
	Majority        int
	Timeout         int
	MajorityTimeout int

	InChannel  chan messages.Packet
	OutChannel chan messages.Packet
}

// Start will trigger a new apaxos protocol.
func (a Apaxos) Start() error {
	// Send prepare messages
	a.broadcastPropose()
	// Wait for promise messages (first on majority, then on a timeout)
	a.waitForPromise()
	// Create a message
	// Send accept messages
	a.broadcastAccept()
	// Wait for accepted messages (first on majority, then on a timeout)
	a.waitForAccepted()
	// Send commit message
	a.broadcastCommit()
	// Clean the dataset

	return nil
}
