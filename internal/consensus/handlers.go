package consensus

import (
	"github.com/f24-cse535/apaxos/internal/utils/compare"
	"github.com/f24-cse535/apaxos/pkg/models"
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"
)

// prepareHandler get's a prepare message and compares the ballot-number
// with it's own ballot-number.
// If everything was ok, it emits a promise message back to the proposer.
// It also checks the last committed block to see
// if it should sync the node that proposed a value.
func (c Consensus) prepareHandler(msg *apaxos.PrepareMessage) error {
	acceptedNum := c.Memory.GetAcceptedNum()
	ballotNumber := models.BallotNumber{}.FromProtoModel(msg.GetBallotNumber())

	// the given accepted num should be absolute bigger than the current ballot-number
	if acceptedNum != nil && compare.CompareBallotNumbers(&ballotNumber, acceptedNum) != 1 {
		return nil
	} else {
		if exist, err := c.Database.IsBlockExists(&ballotNumber); err == nil && exist {
			// send back sync request
			// c.Dialer.Sync(msg.NodeId, )
		} else {
			// send promise message
			// c.Dialer.Promise(c.Nodes[msg.NodeId], &apaxos.PromiseMessage{

			// })
		}
	}

	return nil
}

// acceptHandler get's an accept message and compares the ballot-number
// with it's own ballot-number.
// If everything was ok, it updates its accepted num and accepted var, and
// emits an accepted message.
func (c Consensus) acceptHandler() error {
	return nil
}

// commitHandler get's a commit message and emptys the datastore by executing
// the transactions, and storing them inside database.
func (c Consensus) commitHandler() error {
	return nil
}

// syncHandler get's a sync message and updates itself to catch up with others.
func (c Consensus) syncHandler() error {
	return nil
}
