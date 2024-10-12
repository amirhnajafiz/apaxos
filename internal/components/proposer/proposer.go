package proposer

type Proposer struct {
	GRPCChannel chan bool
}

func (p Proposer) Start() {
	// on start method, the proposer waits for messages from the gRPC server.
	// if the gRPC server signals a APAXOS start, then the proposer starts APAXOS.
	for {
		// wait on gRPC notify channel
		<-p.GRPCChannel
	}
}
