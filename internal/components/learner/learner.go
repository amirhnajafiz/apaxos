package learner

type Learner struct {
	GRPCChannel chan bool
}

func (l Learner) Start() {
	// on start method, the learner waits for messages from the gRPC server.
	for {
		// wait on gRPC notify channel
		<-l.GRPCChannel
	}
}
