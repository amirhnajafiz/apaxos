package apaxos

import (
	"github.com/f24-cse535/apaxos/internal/grpc/client"
	"github.com/f24-cse535/apaxos/internal/storage/local"
	"github.com/f24-cse535/apaxos/pkg/messages"

	"go.uber.org/zap"
)

// Apaxos module is used by the consensus module.
// In consensus module we create instances of this Apaxos.
// In order to communicate with this module, there are two
// channels. The InChannel which the data from consensus will be received.
// And, the OutChannel which the data will be sent to consensus.
type Apaxos struct {
	Logger *zap.Logger   // logger is needed for tracing
	Memory *local.Memory // memory will be used for reading states

	// Dialer and nodes are needed to make RPC calls
	Dialer *client.ApaxosDialer
	Nodes  map[string]string // list of nodes and their addresses is needed for RPC calls
	NodeId string

	// These parameters are used for apaxos protocol
	Majority        int
	Timeout         int
	MajorityTimeout int

	InChannel  chan *messages.Packet // in channel is used to get inputs from the consensus module
	OutChannel chan *messages.Packet // out channel is used to return response to the client
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
