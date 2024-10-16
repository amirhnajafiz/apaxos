package consensus

import (
	"errors"

	protocol "github.com/f24-cse535/apaxos/internal/consensus/apaxos"
	"github.com/f24-cse535/apaxos/internal/grpc/client"
	"github.com/f24-cse535/apaxos/internal/storage/database"
	"github.com/f24-cse535/apaxos/internal/storage/local"
	"github.com/f24-cse535/apaxos/pkg/messages"
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"

	"go.uber.org/zap"
)

// Consensus module is the core module that runs consensus protocols
// by getting the gRPC level packets.
type Consensus struct {
	Memory   *local.Memory      // memory is needed to update the node state
	Database *database.Database // database is needed to store blocks

	Logger *zap.Logger    // logger is needed for tracing
	Dialer *client.Client // dialer is used to make RPC calls

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

// Signal is used by the upper layer (gRPC functions) to send their
// packets to the consensus's instance module without waiting for any response.
func (c *Consensus) Signal(pkt *messages.Packet) {
	if c.instanceExists() {
		c.instance.InChannel <- pkt
	}
}

// notify is used to send a notification over user channel.
func (c *Consensus) notify(err error) {
	if err != nil {
		c.instance.OutChannel <- &messages.Packet{
			Payload: err.Error(),
		}
	} else {
		c.instance.OutChannel <- &messages.Packet{
			Payload: "transaction submitted",
		}
	}
}

// Checkout is the main method of our consensus module to submit or decline a new transaction.
func (c *Consensus) Checkout(pkt *messages.Packet) (chan *messages.Packet, error) {
	// get the payload of input request
	transaction := pkt.Payload.(*apaxos.Transaction)

	// if the receiver is our client then no need to run consensus protocol
	if transaction.Reciever == c.Client || c.Memory.GetBalance(c.Client) > transaction.GetAmount() {
		c.submitTransaction(transaction)

		return nil, nil
	}

	// now we check to see if we can run the consensus protocol
	if c.instanceExists() {
		return nil, ErrMultipleInstances
	}

	// if no instances exist, we create a new apaxos instance by running begin consensus
	c.newInstance(transaction)

	// send an accepted response, so the client waits for a response over the instance out channel
	return c.instance.OutChannel, nil
}

// newInstance builds and starts a new apaxos instance.
func (c *Consensus) newInstance(transaction *apaxos.Transaction) {
	// first we create a new instance for the protocol
	c.instance = &protocol.Apaxos{
		NodeId:          c.NodeId,
		Dialer:          c.Dialer,
		Nodes:           c.Nodes,
		Memory:          c.Memory,
		Majority:        c.Majority,
		MajorityTimeout: c.MajorityTimeout,
		Timeout:         c.RequestTimeout,
		Logger:          c.Logger.Named("apaxos-instance"),
		InChannel:       make(chan *messages.Packet),
		OutChannel:      make(chan *messages.Packet),
	}

	// start a new go-routine for apaxos instance
	go func() {
		// after this go-routine finished, clear the protocol instance
		defer func() {
			c.instance = nil
			c.Logger.Debug("apaxos terminated")
		}()

		c.Logger.Debug("apaxos started")

		// in a while loop, try to make consensus
		for {
			c.Logger.Debug("a round of apaxos")

			// start apaxos protocol
			err := c.instance.Start()
			if err != nil {
				c.Logger.Debug("apaxos failed", zap.Error(err))

				// check the error for a proper handling mechanism
				switch {
				case errors.Is(err, protocol.ErrRequestTimeout):
					// in this case, we first check to see if we have enough servers up and running or not
					if !c.livenessHandler() {
						c.notify(protocol.ErrNotEnoughBalance)
						return
					}
				case errors.Is(err, protocol.ErrSlowNode):
					// if we hit slow-node, it means that we got synced, so we should check the status of balance before
					// rerunning the consensus protocol.
					c.Logger.Info("slow server detected")

					// if the client balance is now enough, then we submit the transaction
					if c.recheckBalance(transaction) {
						c.submitTransaction(transaction)
						c.notify(nil)

						return
					}
				default:
					c.Logger.Error("consensus error", zap.Error(err))
				}
			} else {
				// now we check to see if the client balance is enough or not
				if c.recheckBalance(transaction) {
					c.submitTransaction(transaction)
					c.notify(nil)
				} else {
					c.notify(protocol.ErrNotEnoughBalance)
				}

				return
			}
		}
	}()
}
