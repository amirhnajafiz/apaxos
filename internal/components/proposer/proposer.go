package proposer

import "github.com/f24-cse535/apaxos/pkg/messages"

type Proposer struct {
	GRPCChannel chan *messages.Packet
}

func (p Proposer) Start() {
	// on start method, the proposer waits for messages from the gRPC server.
	// if the gRPC server signals a APAXOS start, then the proposer starts APAXOS.
	for {
		// wait on gRPC notify channel
		<-p.GRPCChannel
	}
}
