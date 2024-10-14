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
)

// Consensus module is the core module that runs consensus protocols
// by getting the gRPC level packets.
type Consensus struct {
	Memory   *local.Memory
	Database *database.Database
	Dialer   client.ApaxosDialer

	Client  string
	NodeId  string
	Clients map[string]string
	Nodes   map[string]string

	Majority        int
	RequestTimeout  int
	MajorityTimeout int

	channel chan *messages.Packet
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
		c.acceptHandler()
	case enum.PacketCommit:
		c.commitHandler()
	case enum.PacketSync:
		c.receiveSyncHandler()
	default:
		if c.channel != nil {
			c.channel <- pkt
		}
	}
}

// Demand is used by components to use the consensus logic to perform an
// operation. When calling demand, the caller waits for consensus to return something.
func (c Consensus) Demand(pkt *messages.Packet) (chan *messages.Packet, error) {
	// get the payload
	transaction := pkt.Payload.(*apaxos.Transaction)

	// if channel is nil without any errors, it means that the transaction should not handle on this machine
	// this should be check on apaxos only. In multiplaxos, we don't care about this.
	if transaction.Sender != c.Client && transaction.Reciever != c.Client {
		return nil, errors.New("your client request does not belong to this machine")
	}

	// if the receiver is our client then no need to run consensus protocol
	if transaction.Reciever == c.Client {
		// create a new transaction
		t := models.Transaction{
			SequenceNumber: c.Memory.GetSequenceNumber(),
		}.FromProtoModel(transaction)

		// save it into datastore
		c.Memory.AddTransactionToDatastore(&t)

		return nil, nil
	}

	// now we check to see if we can run the consensus protocol
	if c.channel != nil {
		return nil, errors.New("cannot run multiple consensus protocols at the same time on this machine")
	}

	// now running consensus by creating a new apaxos instance
	c.channel = make(chan *messages.Packet)
	channel := make(chan *messages.Packet)

	// start a new go-routine
	go func() {
		// creating a new instance of apaxos protocol
		protocol.Apaxos{
			Dialer:          c.Dialer,
			Memory:          c.Memory,
			Nodes:           c.Nodes,
			Majority:        c.Majority,
			MajorityTimeout: c.MajorityTimeout,
			Timeout:         c.RequestTimeout,
			InChannel:       c.channel,
			OutChannel:      channel,
		}.Start()

		// close both channels
		close(channel)
		close(c.channel)

		c.channel = nil
	}()

	return channel, nil
}
