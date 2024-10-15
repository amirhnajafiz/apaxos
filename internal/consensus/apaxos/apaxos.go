package apaxos

import (
	"github.com/f24-cse535/apaxos/internal/grpc/client"
	"github.com/f24-cse535/apaxos/internal/storage/local"
	"github.com/f24-cse535/apaxos/pkg/messages"
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"

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

	// internal messages
	promisedMessage []*apaxos.PromiseMessage
}

// Start will trigger a new apaxos protocol.
func (a *Apaxos) Start() error {
	// create a new promised messages list
	a.promisedMessage = make([]*apaxos.PromiseMessage, 0)

	// in a for loop send prepare messages to get the majority or sync
	for {
		// increase ballot number on each attempt
		ballotNumber := a.Memory.GetBallotNumber()
		ballotNumber.Number++

		// send propose message to all
		a.broadcastPropose(ballotNumber)

		// set new ballot-number for retry
		a.Memory.SetBallotNumber(ballotNumber)

		// wait for promise messages (first on majority, then on a timeout)
		a.waitForPromise()
	}

	// Create a message
	// Send accept messages
	a.broadcastAccept(nil, nil)
	// Wait for accepted messages (first on majority, then on a timeout)
	a.waitForAccepted()
	// Send commit message
	a.broadcastCommit()
	// Clean the dataset

	return nil
}
