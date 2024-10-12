package learner

import "github.com/f24-cse535/apaxos/pkg/messages"

type Learner struct {
	GRPCChannel chan *messages.Packet
}

func (l Learner) Start() {
	// on start method, the learner waits for messages from the gRPC server.
	for {
		// wait on gRPC notify channel
		<-l.GRPCChannel
	}
}
