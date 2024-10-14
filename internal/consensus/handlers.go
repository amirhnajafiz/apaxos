package consensus

import (
	"github.com/f24-cse535/apaxos/internal/utils"
	"github.com/f24-cse535/apaxos/pkg/models"
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"
)

// prepareHandler get's a prepare message from the proposer.
func (c Consensus) prepareHandler(msg *apaxos.PrepareMessage) {
	// now we extract data of the input message
	ballotNumber := models.BallotNumber{}.FromProtoModel(msg.GetBallotNumber())
	lastCommitted := models.BallotNumber{}.FromProtoModel(msg.LastComittedMessage)

	// first we check the given last committed message to our own last committed message
	// if they matched, then we are probably having a synced proposer.
	// if they didn't, we check if we had that message committed or not. if we had committed that
	// message, it means that the proposer should get a sync request.
	if utils.CompareBallotNumbers(&lastCommitted, c.Memory.GetLastCommittedMessage()) == 0 {
		// if the last committed messages match, we follow the promise
		c.promiseHandler(c.Nodes[msg.NodeId], ballotNumber)
	} else {
		// if the last committed message was not the same, we check to see if we had that
		// message committed sometime in past or not.
		if exist, err := c.Database.IsBlockExists(&lastCommitted); err == nil && exist {
			// if the message exists, we send a sync message to update the node
			c.transmitSync(c.Nodes[msg.NodeId])
		} else {
			// if the message was not committed before, we follow the promise
			c.promiseHandler(c.Nodes[msg.NodeId], ballotNumber)
		}
	}
}

// promiseHandler will be called by the prepare handler, to send back a promise message.
func (c Consensus) promiseHandler(address string, ballotNumber models.BallotNumber) {
	// first we get our ballot-number
	savedBallotNumber := c.Memory.GetBallotNumber()

	// now we check the proposer's ballot-number with our own ballot-number.
	// the ballot-number should be absolute greater than ours.
	if utils.CompareBallotNumbers(&ballotNumber, savedBallotNumber) < 1 {
		return
	}

	// if the proposer's ballot-number was absolute greater, we update our ballot-number
	c.Memory.SetBallotNumber(&ballotNumber)

	// then we create a promise message with init fields of our node_id and our last_committed_message
	promiseMessage := &apaxos.PromiseMessage{
		NodeId:              c.NodeId,
		LastComittedMessage: c.Memory.GetLastCommittedMessage().ToProtoModel(),
	}

	// now we need to get our current accepted_num and accepted_val
	acceptedNum := c.Memory.GetAcceptedNum()
	acceptedVal := c.Memory.GetAcceptedVal()

	// then we check to see if we have accepted something from another proposer before, and it is not committed.
	if acceptedNum != nil {
		// first we create a blocklist of our accepted_vals
		blockList := make([]*apaxos.Block, len(acceptedVal))
		for index, item := range acceptedVal {
			blockList[index] = item.ToProtoModel()
		}

		// then we update the promise message
		promiseMessage.Blocks = blockList                        // set accepted_val as blocks
		promiseMessage.BallotNumber = acceptedNum.ToProtoModel() // set accepted_num as ballot-number
	} else { // if nothing was in accepted fields, we send our own log block
		// get the current datastore as a block and set the block metadata
		block := c.Memory.GetDatastoreAsBlock()
		block.Metadata = models.BlockMetadata{
			NodeId:       c.NodeId,     // set node-id to identify the node's block in future
			BallotNumber: ballotNumber, // ballot-number is the same as what the proposer said
		}

		// then update the promise message
		promiseMessage.Blocks = []*apaxos.Block{block.ToProtoModel()} // set node's block as blocks
		promiseMessage.BallotNumber = ballotNumber.ToProtoModel()     // set ballot-number as proposer's sent
	}

	// send the promise message
	c.Dialer.Promise(address, promiseMessage)
}

// acceptHandler get's a accept message from the proposer.
func (c Consensus) acceptHandler(msg *apaxos.AcceptMessage) {
	// first we extract data of the input message
	ballotNumber := models.BallotNumber{}.FromProtoModel(msg.GetBallotNumber())

	// then we get our ballot-number
	savedBallotNumber := c.Memory.GetBallotNumber()

	// now we check the proposer's ballot-number with our own ballot-number.
	// the ballot-number should be greater than ours.
	if utils.CompareBallotNumbers(&ballotNumber, savedBallotNumber) < 0 {
		return
	}

	// update accepted_num with proposer's ballot number
	c.Memory.SetAcceptedNum(&ballotNumber)

	// update accepted_val with proposer's give blocks
	blocks := make([]*models.Block, len(msg.Blocks))
	for index, block := range msg.Blocks {
		tmp := models.Block{}.FromProtoModel(block)
		blocks[index] = &tmp
	}

	c.Memory.SetAcceptedVal(blocks)

	// send accepted message
	c.Dialer.Accepted(c.Nodes[msg.NodeId])
}

// commitHandler get's a commit message and emptys the datastore by executing
// the transactions, and storing them inside database.
func (c Consensus) commitHandler() error {
	return nil
}

// transmitSync will be called by the prepare handler to update the proposer.
func (c Consensus) transmitSync(address string) {
	// get a clone of the clients
	clients := c.Memory.GetClients()

	// create an instance of sync message
	message := &apaxos.SyncMessage{
		LastComittedMessage: c.Memory.GetLastCommittedMessage().ToProtoModel(),
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
	c.Dialer.Sync(address, message)
}

// syncHandler get's a sync message and updates itself to catch up with others.
func (c Consensus) syncHandler() error {
	return nil
}
