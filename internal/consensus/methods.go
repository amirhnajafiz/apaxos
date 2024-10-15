package consensus

import (
	"sort"

	"github.com/f24-cse535/apaxos/internal/utils"
	"github.com/f24-cse535/apaxos/pkg/models"
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"

	"go.uber.org/zap"
)

// Prepare get's a prepare message from a propose and decides to send a promise message or not.
func (c Consensus) Prepare(msg *apaxos.PrepareMessage) {
	// first we extract some data out of the input message
	proposerBallotNumber := models.BallotNumber{}.FromProtoModel(msg.GetBallotNumber())
	proposerLastCommittedMessage := models.BallotNumber{}.FromProtoModel(msg.GetLastComittedMessage())

	// first we check the proposer last committed message to our own last committed message
	if utils.CompareBallotNumbers(&proposerLastCommittedMessage, c.Memory.GetLastCommittedMessage()) == 0 {
		// if they where same as each, we continue the logic in promise handler
		c.promiseHandler(c.Nodes[msg.NodeId], proposerBallotNumber)
	} else {
		// we check to see if we had that message committed sometime in past or not
		if exist, err := c.Database.IsBlockExists(&proposerBallotNumber); err == nil && exist {
			c.transmitSync(c.Nodes[msg.NodeId]) // the message exists but is old
		} else {
			c.promiseHandler(c.Nodes[msg.NodeId], proposerBallotNumber) // the message is new
		}
	}
}

// Accept get's a accept message from the proposer.
func (c Consensus) Accept(msg *apaxos.AcceptMessage) {
	// first we get the proposer's ballot number
	proposerBallotNumber := models.BallotNumber{}.FromProtoModel(msg.GetBallotNumber())

	// then we get our saved ballot-number
	savedBallotNumber := c.Memory.GetBallotNumber()

	// now we check the proposer's ballot-number with our own ballot-number.
	if utils.CompareBallotNumbers(&proposerBallotNumber, savedBallotNumber) < 0 {
		// this means the the proposer's ballot-number is <= our saved ballot-number
		return
	}

	// update accepted_num with proposer's ballot-number
	c.Memory.SetAcceptedNum(&proposerBallotNumber)

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

// Commit emptys the datastore and stores the blocks inside database.
func (c Consensus) Commit() {
	// get our accepted_val
	acceptedVal := c.Memory.GetAcceptedVal()

	// sort the blocks by their ballot-numbers
	sort.Slice(acceptedVal, func(i, j int) bool { // we switch the places of i and j to sort the list in decreasing order
		return utils.CompareBlocks(&acceptedVal[j].Metadata, &acceptedVal[i].Metadata)
	})

	// now we should execute the transactions
	for _, block := range acceptedVal {
		// update our own blocks in memory, to remove previous transactions
		if block.Metadata.NodeId == c.NodeId {
			c.Memory.ClearDatastore(block)
		} else {
			// get transactions and sort them by sequence number
			tlist := block.Transactions
			sort.Slice(tlist, func(i, j int) bool { // transactions are sorted in increasing order
				return tlist[i].SequenceNumber < tlist[j].SequenceNumber
			})

			// loop in transactions and execute them
			for _, transaction := range tlist {
				c.Memory.UpdateBalance(transaction.Sender, transaction.Amount*-1)
				c.Memory.UpdateBalance(transaction.Reciever, transaction.Amount)
			}
		}
	}

	// now we store the blocks inside MongoDB
	for {
		if err := c.Database.InsertBlocks(acceptedVal); err != nil {
			c.Logger.Error("failed to store blocks inside MongoDB", zap.Error(err))
		} else {
			break
		}
	}

	// finally, we clear our accepted_num and accepted_val
	c.Memory.SetAcceptedNum(nil)
	c.Memory.SetAcceptedVal(nil)
}

// Sync get's a sync message and updates itself to catch up with others.
func (c Consensus) Sync(msg *apaxos.SyncMessage) {
	// update clients' balances
	for _, item := range msg.GetPairs() {
		c.Memory.SetBalance(item.GetClient(), item.GetBalance())
	}

	// update our last committed message
	ballotNumber := models.BallotNumber{}.FromProtoModel(msg.GetLastComittedMessage())
	c.Memory.SetLastCommittedMessage(&ballotNumber)
}
