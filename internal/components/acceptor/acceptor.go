package acceptor

type Acceptor struct {
	GRPCChannel chan bool
}

func (a Acceptor) Start() {
	// on start method, the acceptor waits for messages from the gRPC server.
	for {
		// wait on gRPC notify channel
		<-a.GRPCChannel
	}
}
