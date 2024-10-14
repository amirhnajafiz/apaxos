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
			c.sendSyncHandler(c.Nodes[msg.NodeId])
		} else {
			// if the message was not committed before, we follow the promise
			c.promiseHandler(c.Nodes[msg.NodeId], ballotNumber)
		}
	}
}

// promiseHandler will be called by the prepare handler, to send back a promise message.
func (c Consensus) promiseHandler(address string, ballotNumber models.BallotNumber) {
	// first we get our accepted_num, accepted_val
	acceptedNum := c.Memory.GetAcceptedNum()
	acceptedVal := c.Memory.GetAcceptedVal()

	// then we create a promise message with init fields of our node_id and our last_committed_message
	promiseMessage := &apaxos.PromiseMessage{
		NodeId:              c.NodeId,
		LastComittedMessage: c.Memory.GetLastCommittedMessage().ToProtoModel(),
	}

	// then we check to see if we have accepted something from another proposer before, and it is not committed.
	// if we had something, first we should check the ballot-numbers.
	if acceptedNum != nil {
		// the ballot-number should be absolute greater than accepted_num
		if utils.CompareBallotNumbers(&ballotNumber, acceptedNum) == 1 {
			// first we create a blocklist of our accepted_vals
			blockList := make([]*apaxos.Block, len(acceptedVal))
			for index, item := range acceptedVal {
				blockList[index] = item.ToProtoModel()
			}

			// then we update the promise message
			promiseMessage.Blocks = blockList                        // set accepted_val as blocks
			promiseMessage.BallotNumber = acceptedNum.ToProtoModel() // set accepted_num as ballot-number
		}
	} else {
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

// acceptHandler get's an accept message and compares the ballot-number
// with it's own ballot-number.
// If everything was ok, it updates its accepted num and accepted var, and
// emits an accepted message.
func (c Consensus) acceptHandler() {
}

// commitHandler get's a commit message and emptys the datastore by executing
// the transactions, and storing them inside database.
func (c Consensus) commitHandler() error {
	return nil
}

// sendSyncHandler will be called by the prepare handler to update the proposer.
func (c Consensus) sendSyncHandler(address string) {
	// get a clone of the clients
	clients := c.Memory.GetClients()

	// create sync messages
	messages := make([]*apaxos.SyncMessage, len(clients))
	for key, value := range clients {
		messages = append(messages, &apaxos.SyncMessage{
			Client:  key,
			Balance: value,
		})
	}

	// send the sync message
	c.Dialer.Sync(address, messages)
}

// receiveSyncHandler get's a sync message and updates itself to catch up with others.
func (c Consensus) receiveSyncHandler() error {
	return nil
}
