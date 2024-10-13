package consensus

// prepareHandler get's a prepare message and compares the ballot-number
// with it's own ballot-number.
// If everything was ok, it emits a promise message back to the proposer.
// It also checks the last committed block to see
// if it should sync the node that proposed a value.
func (c Consensus) prepareHandler() error {
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

// aggregatedBalanceHandler get's a request of client id and runs a RPC call
// to get an answer from each node.
func (c Consensus) aggregatedBalanceHandler() error {
	return nil
}
