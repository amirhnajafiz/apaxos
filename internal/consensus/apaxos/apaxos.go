package apaxos

import (
	"github.com/f24-cse535/apaxos/internal/grpc/client"
	"github.com/f24-cse535/apaxos/internal/storage/database"
	"github.com/f24-cse535/apaxos/internal/storage/local"
	"github.com/f24-cse535/apaxos/internal/utils"
	"github.com/f24-cse535/apaxos/pkg/messages"
	"github.com/f24-cse535/apaxos/pkg/models"
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"

	"go.uber.org/zap"
)

// Apaxos module is used by the consensus module.
// In consensus module we create instances of this Apaxos.
// In order to communicate with this module, there are two
// channels. The InChannel which the data from consensus will be received.
// And, the OutChannel which the data will be sent to consensus.
type Apaxos struct {
	Logger   *zap.Logger        // logger is needed for tracing
	Memory   *local.Memory      // memory will be used for reading states
	Database *database.Database // database is needed to check committed messages

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

	// selected ballot-number
	selectedBallotNumber *models.BallotNumber
	selectedBlocks       []*apaxos.Block

	// accepted values
	acceptedNum *apaxos.BallotNumber
	acceptedVal []*apaxos.Block
}

// Start will trigger a new apaxos protocol.
func (a *Apaxos) Start() error {
	// create a new promised messages list
	a.promisedMessage = make([]*apaxos.PromiseMessage, 0)

	// increase ballot number on each attempt
	a.selectedBallotNumber = a.Memory.GetBallotNumber()
	a.selectedBallotNumber.Number++

	a.Logger.Debug(
		"sending prepare",
		zap.Int64("number", a.selectedBallotNumber.Number),
		zap.String("node", a.selectedBallotNumber.NodeId),
	)

	// send propose message to all
	a.broadcastPropose(a.selectedBallotNumber)

	// set new ballot-number for retry
	a.Memory.SetBallotNumber(a.selectedBallotNumber)

	// wait for promise messages (first on majority, then on a timeout)
	err := a.waitForPromise()
	if err != nil {
		a.Logger.Debug("wait for promise error", zap.Error(err))

		return err
	}

	// create a message and then send the values as accept message
	a.createMessage()
	a.broadcastAccept(a.acceptedNum, a.acceptedVal)

	// wait for accepted messages (first on majority, then on a timeout)
	a.waitForAccepted()

	// send commit message
	a.broadcastCommit()

	return nil
}

func (a *Apaxos) createMessage() {
	// create a list of blocks
	a.selectedBlocks = make([]*apaxos.Block, 0)

	// we go through the promissed messages to check everything needed
	for _, msg := range a.promisedMessage {
		// get their ballot-number and last committed message
		ballotNumber := models.BallotNumber{}.FromProtoModel(msg.BallotNumber)
		lastCommitted := models.BallotNumber{}.FromProtoModel(msg.LastComittedMessage)

		// first we check to see if they had a different ballot-number or not
		if utils.CompareBallotNumbers(a.selectedBallotNumber, &ballotNumber) == 0 {
			// if they send the same ballot-number we save their blocks
			a.selectedBlocks = append(a.selectedBlocks, msg.GetBlocks()...)
		} else {
			// in this case, they have sent a ballot-number different than ours
			// first we check their last committed message to see if they are behind or not
			if utils.CompareBallotNumbers(&lastCommitted, a.Memory.GetLastCommittedMessage()) == 0 {
				// if they are synced, we are going to check their ballot-number
				a.getAccepted(msg)
			} else {
				// else we need to check their message existance to send them sync or to accept them
				if exist, err := a.Database.IsBlockExists(&lastCommitted); err == nil && exist {
					// if we have that message then we try to sync them
					a.transmitSync(msg.NodeId)
				} else {
					// else we take their message
					a.getAccepted(msg)
				}
			}
		}
	}

	// check to see if we need to update accepted_num and accepted_val
	if a.acceptedNum == nil {
		a.acceptedNum = a.selectedBallotNumber.ToProtoModel()
		a.acceptedVal = a.selectedBlocks
	}
}

// getAccepted extracts the accepted_num and accepted_val of a promise message, then it will update
// their values if the ballot-number was bigger.
func (a *Apaxos) getAccepted(msg *apaxos.PromiseMessage) {
	if a.acceptedNum == nil {
		a.acceptedNum = msg.GetBallotNumber()
		a.acceptedVal = msg.GetBlocks()
	} else {
		tmpA := models.BallotNumber{}.FromProtoModel(a.acceptedNum)
		tmpB := models.BallotNumber{}.FromProtoModel(msg.BallotNumber)

		// then we select the biggest ballot-number as our accepted_val
		if utils.CompareBallotNumbers(&tmpB, &tmpA) == 1 {
			a.acceptedNum = msg.GetBallotNumber()
			a.acceptedVal = msg.GetBlocks()
		}
	}
}

// transmitSync will be called by the prepare handler to update the proposer.
func (a *Apaxos) transmitSync(address string) {
	// get a clone of the clients
	clients := a.Memory.GetClients()

	// create an instance of sync message
	message := &apaxos.SyncMessage{
		LastComittedMessage: a.Memory.GetLastCommittedMessage().ToProtoModel(),
		Pairs:               make([]*apaxos.ClientBalancePair, len(clients)),
	}

	// add client and their balances
	index := 0
	for key, value := range clients {
		message.Pairs[index] = &apaxos.ClientBalancePair{
			Client:  key,
			Balance: value,
		}

		index++
	}

	// send the sync message
	a.Dialer.Sync(address, message)
}
