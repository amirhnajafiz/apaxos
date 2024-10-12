package proposer

func (p Proposer) initApaxos() error {
	// send promise messages
	// call promis
	// call accept
	// call commit
	// return gRPC response
}

func (p Proposer) propose() error {
	// create the ballot_number
	// wait for the majority
}

func (p Proposer) accept() error {
	// collects all logs to create a major block
	// send accept request
	// wait for the majority
}

func (p Proposer) commit() error {
	// send commit message
}
