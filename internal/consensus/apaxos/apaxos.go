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
// In order to communicate with this module, there are two channels.
// The InChannel which the data from consensus will be received. And, the OutChannel which the data will be sent to consensus.
type Apaxos struct {
	Logger *zap.Logger   // logger is needed for tracing
	Memory *local.Memory // memory will be used for reading states

	// Dialer and nodes are needed to make RPC calls
	Dialer *client.Client
	Nodes  map[string]string // list of nodes and their addresses is needed for RPC calls
	NodeId string

	// These parameters are used for apaxos protocol
	Majority        int
	Timeout         int
	MajorityTimeout int

	InChannel  chan *messages.Packet // in channel is used to get inputs from the consensus module
	OutChannel chan *messages.Packet // out channel is used to return response to the client

	// internal messages
	// promised messages is used to store all input promise messages
	promisedMessage []*apaxos.PromiseMessage

	// selected ballot-number and selected-blocks
	selectedBallotNumber *apaxos.BallotNumber // we set this in the beginning of the start method
	selectedBlocks       []*apaxos.Block      // this will be fulled in createMessage method

	// accepted values to submit to others as accept request
	acceptedNum *apaxos.BallotNumber
	acceptedVal []*apaxos.Block
}

// Start will trigger a new apaxos protocol.
// This protocl first sends propose messages, then it waits on promised messages,
// then it sends accept messages, then it waits on accepted messages, finally, it
// sends commit messages.
func (a *Apaxos) Start() error {
	defer func() {
		// set new ballot-number in memory for next attempts
		a.Memory.SetBallotNumber(a.selectedBallotNumber)
	}()

	// create a new promised messages list to get all promise messages
	a.promisedMessage = make([]*apaxos.PromiseMessage, 0)

	// get the current ballot-number from memory to increase ballot number on each attempt
	tmp := a.Memory.GetBallotNumber()
	a.selectedBallotNumber = &apaxos.BallotNumber{
		NodeId: tmp.GetNodeId(),
		Number: tmp.GetNumber() + 1,
	}

	a.Logger.Info(
		"sending prepare",
		zap.Int64("number", a.selectedBallotNumber.Number),
		zap.String("node", a.selectedBallotNumber.NodeId),
	)

	// send a propose message to all existing nodes
	a.broadcastPropose(a.selectedBallotNumber)

	// wait for promise messages (first on majority, then on a timeout)
	err := a.waitForPromise()
	if err != nil {
		a.Logger.Debug("wait for promise error", zap.Error(err))
		return err
	}

	a.Logger.Info("got majority promise", zap.Int("promises", len(a.promisedMessage)))

	// create accepted_num and accepted_val by checking the promised messages
	a.processPromiseMessages()

	a.Logger.Info(
		"sending promise",
		zap.Int64("number", a.selectedBallotNumber.Number),
		zap.String("node", a.selectedBallotNumber.NodeId),
		zap.Int("blocks", len(a.selectedBlocks)),
	)

	// send a broadcast accept message to all other servers
	a.broadcastAccept(a.selectedBallotNumber, a.selectedBlocks)

	// wait for accepted messages (first on majority, then on a timeout)
	if err := a.waitForAccepted(); err != nil {
		a.Logger.Debug("wait for accepted error", zap.Error(err))
		return err
	}

	a.Logger.Info("got majority accepted")

	// send commit message to all other servers
	a.broadcastCommit()

	a.Logger.Info("sent commit messages, waiting for my own commit")

	// wait for own commit packet in case of race condition
	return a.waitForOwnCommit()
}
