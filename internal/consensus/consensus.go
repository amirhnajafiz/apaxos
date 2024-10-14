package consensus

import (
	"errors"

	protocol "github.com/f24-cse535/apaxos/internal/consensus/apaxos"
	"github.com/f24-cse535/apaxos/internal/grpc/client"
	"github.com/f24-cse535/apaxos/internal/storage/database"
	"github.com/f24-cse535/apaxos/internal/storage/local"
	"github.com/f24-cse535/apaxos/pkg/enum"
	"github.com/f24-cse535/apaxos/pkg/messages"
	"github.com/f24-cse535/apaxos/pkg/models"
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"

	"go.uber.org/zap"
)

// Consensus module is the core module that runs consensus protocols
// by getting the gRPC level packets.
type Consensus struct {
	Memory   *local.Memory      // memory is needed to update the node state
	Database *database.Database // database is needed to store blocks

	Logger *zap.Logger          // logger is needed for tracing
	Dialer *client.ApaxosDialer // apaxos dialer is needed in handler methods

	Client string // client is needed to identify input transactions
	NodeId string // nodeId is needed for making RPC calls

	Nodes map[string]string // list of nodes and their addresses is needed for RPC calls

	// these parameters needed to use in apaxos
	Majority        int
	RequestTimeout  int
	MajorityTimeout int

	// each consensus should keep track of the protocol instance
	instance *protocol.Apaxos
}

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

// Signal is used by the upper layer (gRPC functions) to send their
// packets to the consensus module without getting any response.
func (c Consensus) Signal(pkt *messages.Packet) {
	// switch case on pkt types to see if you should handle them or
	// they should go into the channel of apaxos instance.
	switch pkt.Type {
	case enum.PacketPrepare: // on prepare packet arrival, call prepare handler
		c.prepareHandler(pkt.Payload.(*apaxos.PrepareMessage))
	case enum.PacketAccept:
		c.acceptHandler(pkt.Payload.(*apaxos.AcceptMessage))
	case enum.PacketCommit:
		c.commitHandler()
	case enum.PacketSync:
		c.syncHandler(pkt.Payload.(*apaxos.SyncMessage))
		c.forwardToInstance(pkt)
	default:
		c.forwardToInstance(pkt)
	}
}

// Demand is used by components to use the consensus logic to perform an
// operation. When calling demand, the caller waits for consensus to return something.
func (c Consensus) Demand(pkt *messages.Packet) (chan *messages.Packet, int, error) {
	// get the payload of input request
	transaction := pkt.Payload.(*apaxos.Transaction)

	// if channel is nil without any errors, it means that the transaction should not handle on this machine
	// this should be check on apaxos only. In multiplaxos, we don't care about this.
	if transaction.Sender != c.Client && transaction.Reciever != c.Client {
		return nil, enum.ResponseRequestFailed, errors.New("your client request does not belong to this machine")
	}

	// if the receiver is our client then no need to run consensus protocol
	if transaction.Reciever == c.Client {
		// save the transaction into datastore
		t := models.Transaction{}.FromProtoModel(transaction)
		t.SequenceNumber = c.Memory.GetSequenceNumber()

		// make changes into memory for client's balances
		c.Memory.UpdateBalance(t.Sender, t.Amount*-1)
		c.Memory.UpdateBalance(t.Reciever, t.Amount)

		// save it into datastore
		c.Memory.AddTransactionToDatastore(t)

		return nil, enum.ResponseOK, nil
	}

	// now we check to see if we can run the consensus protocol
	if c.instanceExists() {
		return nil, enum.ResponseServerFailed, errors.New("cannot run multiple consensus protocols at the same time on this machine")
	}

	// if no instances exist, we create a new apaxos instance
	c.instance = &protocol.Apaxos{
		Dialer:          c.Dialer,
		Nodes:           c.Nodes,
		Memory:          c.Memory,
		Majority:        c.Majority,
		MajorityTimeout: c.MajorityTimeout,
		Timeout:         c.RequestTimeout,
		InChannel:       make(chan *messages.Packet),
		OutChannel:      make(chan *messages.Packet),
	}

	// start a new go-routine for apaxos instance
	go func() {
		c.Logger.Debug("apaxos started")

		// start apaxos
		if err := c.instance.Start(); err != nil {
			c.Logger.Error("apaxos failed", zap.Error(err))
		}

		// reset the instance
		c.instance = nil
	}()

	// send an accepted response, so the client waits for a response
	return c.instance.OutChannel, enum.ResponseAccepted, nil
}
