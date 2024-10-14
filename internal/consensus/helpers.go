package consensus

import "github.com/f24-cse535/apaxos/pkg/messages"

// forward to instance is a helper function to pass packets in apaxos instance channel
func (c Consensus) forwardToInstance(pkt *messages.Packet) {
	if c.instance != nil {
		c.instance.InChannel <- pkt
	}
}

// instance exists return true if the apaxos instance is started and running
func (c Consensus) instanceExists() bool {
	return c.instance == nil
}
