package consensus

import (
	"sort"

	"github.com/f24-cse535/apaxos/internal/utils"
	"github.com/f24-cse535/apaxos/pkg/enum"
	"github.com/f24-cse535/apaxos/pkg/messages"
	"github.com/f24-cse535/apaxos/pkg/models"
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"

	"go.uber.org/zap"
)

// Prepare get's a prepare message from a propose and decides to send a promise message or not.
func (c *Consensus) Prepare(msg *apaxos.PrepareMessage) {
	// first we extract some data out of the input message
	proposerBallotNumber := msg.GetBallotNumber()
	proposerLastCommittedMessage := msg.GetLastComittedMessage()

	// first we check the proposer last committed message to our own last committed message
	compareResult := utils.CompareBallotNumbers(proposerLastCommittedMessage, c.Memory.GetLastCommittedMessage())
	if compareResult >= 0 {
		c.Logger.Debug("compare result states a sync client", zap.Int("result", compareResult))

		// if they where same as us or greater, we continue the logic in promise handler
		c.promiseHandler(c.Nodes[msg.NodeId], proposerBallotNumber, compareResult == 0)
	} else {
		c.Logger.Warn(
			"compare last committed messages states a sync needed (out-of-sync node detected)",
			zap.Int("result", compareResult),
		)

		// proposer is a not syned with us, therefore, we try to sync it by transmitting a sync request
		c.transmitSync(c.Nodes[msg.NodeId])
	}
}

// Accept get's a accept message from the proposer.
func (c *Consensus) Accept(msg *apaxos.AcceptMessage) {
	// first we get the proposer's ballot number
	proposerBallotNumber := msg.GetBallotNumber()

	// then we get our saved ballot-number
	savedBallotNumber := c.Memory.GetBallotNumber()

	// now we check the proposer's ballot-number with our own ballot-number.
	if proposerBallotNumber.GetNumber() < savedBallotNumber.GetNumber() {
		c.Logger.Debug(
			"no greater or equal ballot number",
			zap.Int64("saved_number", savedBallotNumber.GetNumber()),
			zap.String("saved_node_id", savedBallotNumber.GetNodeId()),
			zap.Int64("proposer_number", proposerBallotNumber.GetNumber()),
			zap.String("proposer_node_id", proposerBallotNumber.GetNodeId()),
		)

		// this means the the proposer's ballot-number is <= our saved ballot-number
		return
	}

	// update accepted_num with proposer's ballot-number
	c.Memory.SetAcceptedNum(proposerBallotNumber)
	// update accepted_val with proposer's give blocks
	c.Memory.SetAcceptedVal(msg.Blocks)

	c.Logger.Info(
		"accept",
		zap.String("from", proposerBallotNumber.GetNodeId()),
		zap.Int64("number", proposerBallotNumber.GetNumber()),
	)

	c.Logger.Info("accepted sent", zap.String("to", msg.NodeId))

	// send accepted message
	c.Dialer.Accepted(c.Nodes[msg.NodeId])
}

// Commit emptys the datastore and stores the blocks inside database.
func (c *Consensus) Commit() {
	// get our accepteds
	acceptedNum := c.Memory.GetAcceptedNum()
	acceptedVal := c.Memory.GetAcceptedVal()

	// sort the blocks by their ballot-numbers
	sort.Slice(acceptedVal, func(i, j int) bool { // we switch the places of i and j to sort the list in decreasing order
		return utils.CompareBlocks(acceptedVal[j].GetMetadata(), acceptedVal[i].GetMetadata())
	})

	// now we should execute the transactions
	for _, block := range acceptedVal {
		// update our own blocks in memory, to remove previous transactions
		if block.Metadata.GetNodeId() == c.NodeId {
			c.Memory.ClearDatastore(block)
		} else {
			// get transactions and sort them by sequence number
			tlist := block.GetTransactions()
			sort.Slice(tlist, func(i, j int) bool { // transactions are sorted in increasing order
				return tlist[i].GetSequenceNumber() < tlist[j].GetSequenceNumber()
			})

			// loop in transactions and execute them
			for _, transaction := range tlist {
				c.submitTransaction(transaction, false)
			}
		}
	}

	// not to mention that we should update our last committed message
	c.Memory.SetLastCommittedMessage(acceptedNum)

	// now we store the blocks inside MongoDB
	blocks := make([]*models.Block, 0)
	for _, block := range acceptedVal {
		if len(block.Transactions) > 0 { // only save the blocks that have transactions
			// convert to models.block
			tmp := models.Block{}
			tmp.FromProtoModel(block)

			// append to blocks
			blocks = append(blocks, &tmp)
		}
	}

	// query execution on MongoDB
	if err := c.Database.InsertBlocks(blocks); err != nil {
		c.Logger.Warn("failed to store blocks inside MongoDB", zap.Error(err))
	}

	// finally, we clear our accepted_num and accepted_val
	c.Memory.SetAcceptedNum(nil)
	c.Memory.SetAcceptedVal(nil)

	// signal the apaxos instance for a commit message
	c.Signal(&messages.Packet{
		Type: enum.PacketCommit,
	})

	c.Logger.Info("committed")
}

// Sync get's a sync message and updates itself to catch up with others.
func (c *Consensus) Sync(msg *apaxos.SyncMessage) {
	// update clients' balances
	for _, item := range msg.GetPairs() {
		c.Memory.SetBalance(item.GetClient(), item.GetBalance())
	}

	// update our last committed message
	c.Memory.SetLastCommittedMessage(msg.GetLastComittedMessage())
	c.Memory.SetAcceptedNum(nil)
	c.Memory.SetAcceptedVal(nil)

	c.Logger.Info(
		"node is syncronized",
		zap.Int64("last_number", msg.GetLastComittedMessage().GetNumber()),
		zap.String("last_node_id", msg.GetLastComittedMessage().GetNodeId()),
	)
}
