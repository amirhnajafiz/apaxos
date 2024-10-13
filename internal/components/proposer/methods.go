package proposer

// propose handler first creates a new ballot_number,
// then sends it to other nodes to get the majority response.
func (p Proposer) propose() error {
	// create the ballot_number
	// wait for the majority

	return nil
}

// accept handler first creates a big block of transactions,
// then sends an accept message to other nodes, and waits
// for the majority.
func (p Proposer) accept() error {
	// collects all logs to create a major block
	// send accept request
	// wait for the majority

	return nil
}

// commit handler sends a commit message to other nodes
// and waits then terminates the state-machine.
func (p Proposer) commit() error {
	// send commit message

	return nil
}
