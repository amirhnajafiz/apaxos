package proposer

// initApaxos starts the state-machine for running
// the paxos protocol.
func (p Proposer) initApaxos() error {
	// send promise messages
	// call promis
	// call accept
	// call commit
	// return gRPC response
}

// propose handler first creates a new ballot_number,
// then sends it to other nodes to get the majority response.
func (p Proposer) propose() error {
	// create the ballot_number
	// wait for the majority
}

// accept handler first creates a big block of transactions,
// then sends an accept message to other nodes, and waits
// for the majority.
func (p Proposer) accept() error {
	// collects all logs to create a major block
	// send accept request
	// wait for the majority
}

// commit handler sends a commit message to other nodes
// and waits then terminates the state-machine.
func (p Proposer) commit() error {
	// send commit message
}
