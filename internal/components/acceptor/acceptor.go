package acceptor

import "github.com/f24-cse535/apaxos/pkg/messages"

type Acceptor struct {
	GRPCChannel chan *messages.Packet
}

func (a Acceptor) Start() {
	// on start method, the acceptor waits for messages from the gRPC server.
	for {
		// wait on gRPC notify channel
		<-a.GRPCChannel
	}
}
