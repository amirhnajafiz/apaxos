package apaxos

import (
	"github.com/f24-cse535/apaxos/internal/utils"
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"
)

// processPromiseMessages is used to process all promise messages
// and prepare accepted_val and accepted_num values.
func (a *Apaxos) processPromiseMessages() {
	// create a list of blocks
	a.selectedBlocks = make([]*apaxos.Block, 0)

	// we go through the promissed messages to check everything needed
	for _, msg := range a.promisedMessage {
		// check to see if we the node is synced or not
		if utils.CompareBallotNumbers(msg.GetLastComittedMessage(), a.Memory.GetLastCommittedMessage()) != 0 {
			a.Logger.Info("out-of-sync node detected")

			// sync the old acceptor's log
			a.transmitSync(msg.NodeId)
		}

		// process their promise messages, check to see if they have a different ballot-number
		if utils.CompareBallotNumbers(msg.GetBallotNumber(), a.selectedBallotNumber) != 0 {
			if a.acceptedNum == nil { // empty accepted_num and accepted_val
				a.acceptedNum = msg.GetBallotNumber()
				a.acceptedVal = msg.GetBlocks()
			} else if utils.CompareBallotNumbers(msg.GetBallotNumber(), a.acceptedNum) == 1 {
				a.acceptedNum = msg.GetBallotNumber()
				a.acceptedVal = msg.GetBlocks()
			}
		} else {
			// if they send the same ballot-number, we save their blocks
			a.selectedBlocks = append(a.selectedBlocks, msg.GetBlocks()...)
		}
	}

	// check to see if we need anything from previous time
	if a.acceptedNum != nil {
		a.selectedBlocks = a.acceptedVal
	}
}
