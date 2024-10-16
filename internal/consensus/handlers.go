package consensus

import (
	"github.com/f24-cse535/apaxos/internal/utils"
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"

	"go.uber.org/zap"
)

// promiseHandler will be called by the prepare to perform the promise logic.
func (c *Consensus) promiseHandler(address string, ballotNumber *apaxos.BallotNumber, synced bool) {
	// first we get our ballot-number
	savedBallotNumber := c.Memory.GetBallotNumber()

	// now we check the proposer's ballot-number with our own ballot-number
	if utils.CompareBallotNumbers(ballotNumber, savedBallotNumber) < 1 {
		// this means that the input ballot-number is < saved ballot-number
		return
	}

	// if the proposer's ballot-number was absolute greater, we update our ballot-number
	c.Memory.SetBallotNumber(ballotNumber)

	// then we create a promise message with init fields of our node_id and our last_committed_message
	promiseMessage := &apaxos.PromiseMessage{
		NodeId:              c.NodeId,
		LastComittedMessage: c.Memory.GetLastCommittedMessage(),
	}

	// now we need to get our current accepted_num and accepted_val to see what we should send in the promise message
	acceptedNum := c.Memory.GetAcceptedNum()
	if synced && acceptedNum != nil { // if we are syned and there is something we have from previous
		// then we update the promise message
		promiseMessage.Blocks = c.Memory.GetAcceptedVal() // set accepted_val as blocks
		promiseMessage.BallotNumber = acceptedNum         // set accepted_num as ballot-number
	} else { // if nothing was in accepted fields, we send our own log block
		// get the current datastore as a block and set the block ballot-number
		block := c.Memory.GetDatastore()
		block.Metadata.BallotNumber = ballotNumber     // ballot-number is the same as what the proposer said
		promiseMessage.Blocks = []*apaxos.Block{block} // set node's block as blocks
		promiseMessage.BallotNumber = ballotNumber     // set ballot-number as proposer's sent
	}

	// send the promise message
	c.Dialer.Promise(address, promiseMessage)
}

// transmitSync will be called by the prepare handler to update the proposer.
func (c *Consensus) transmitSync(address string) {
	// get a clone of the clients
	clients := c.Memory.GetClients()

	// create an instance of sync message
	message := &apaxos.SyncMessage{
		LastComittedMessage: c.Memory.GetLastCommittedMessage(),
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

// liveness handler checks to see if there are enough live servers or not.
func (c *Consensus) livenessHandler() bool {
	count := 0

	// send ping messages to servers
	for key, value := range c.Nodes {
		if c.Dialer.Ping(value) {
			c.Logger.Debug("found alive server", zap.String("node", key))
			count++
		}
	}

	return count >= c.Majority
}
