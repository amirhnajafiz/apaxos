package consensus

import "github.com/f24-cse535/apaxos/pkg/messages"

type Consensus struct{}

func (c Consensus) Signal(pkt *messages.Packet) error {
	return nil
}
