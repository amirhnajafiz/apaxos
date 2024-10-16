package apaxos

import (
	"github.com/f24-cse535/apaxos/internal/utils"
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"
)

// createMessage is used to process all promise messages
// and prepare accepted_val and accepted_num values.
func (a *Apaxos) processPromiseMessages() {
	// create a list of blocks
	a.selectedBlocks = make([]*apaxos.Block, 0)

	// we go through the promissed messages to check everything needed
	for _, msg := range a.promisedMessage {
		// get their ballot-number and last committed message
		ballotNumber := msg.GetBallotNumber()
		lastCommitted := msg.GetLastComittedMessage()

		// first we check to see if they had a different ballot-number or not
		if utils.CompareBallotNumbers(a.selectedBallotNumber, ballotNumber) == 0 {
			// if they send the same ballot-number we save their blocks
			a.selectedBlocks = append(a.selectedBlocks, msg.GetBlocks()...)
		} else {
			// in this case, they have sent a ballot-number different than ours
			// first we check their last committed message to see if they are behind or not
			if utils.CompareBallotNumbers(lastCommitted, a.Memory.GetLastCommittedMessage()) == 0 {
				// if they are synced, we are going to check their ballot-number
				a.getAccepteds(msg)
			} else {
				// else we need to check their message existance to send them sync or to accept them
				if exist, err := a.Database.IsBlockExists(lastCommitted.GetNumber(), lastCommitted.GetNodeId()); err == nil && exist {
					// if we have that message then we try to sync them cause it is committed
					go a.transmitSync(msg.NodeId)
				} else {
					// else we take their message to store
					a.getAccepteds(msg)
				}
			}
		}
	}

	// check to see if we need to update accepted_num and accepted_val or no
	if a.acceptedNum == nil {
		a.acceptedNum = a.selectedBallotNumber
		a.acceptedVal = a.selectedBlocks
	}
}

// getAccepted extracts the accepted_num and accepted_val of a promise message, then it will update
// their values if the ballot-number was bigger.
func (a *Apaxos) getAccepteds(msg *apaxos.PromiseMessage) {
	if a.acceptedNum == nil { // empty accepted_num and accepted_val
		a.acceptedNum = msg.GetBallotNumber()
		a.acceptedVal = msg.GetBlocks()
	} else {
		// then we select the biggest ballot-number as our accepted_val
		if utils.CompareBallotNumbers(msg.GetBallotNumber(), a.acceptedNum) == 1 {
			a.acceptedNum = msg.GetBallotNumber()
			a.acceptedVal = msg.GetBlocks()
		}
	}
}
