package consensus

import (
	"errors"

	protocol "github.com/f24-cse535/apaxos/internal/consensus/apaxos"
	"github.com/f24-cse535/apaxos/pkg/messages"
	"github.com/f24-cse535/apaxos/pkg/models"
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"

	"go.uber.org/zap"
)

// submitTransaction is used to store a transaction into datastore, and
// update balances.
func (c Consensus) submitTransaction(transaction *apaxos.Transaction) {
	c.Logger.Info(
		"transaction submitted",
		zap.Int64("amount", transaction.Amount),
		zap.Int64("seq_number", transaction.SequenceNumber),
		zap.String("sender", transaction.Sender),
		zap.String("receiver", transaction.Reciever),
	)

	// save the transaction into datastore
	t := models.Transaction{}.FromProtoModel(transaction)
	t.SequenceNumber = c.Memory.GetSequenceNumber()

	// make changes into memory for client's balances
	c.Memory.UpdateBalance(t.Sender, t.Amount*-1)
	c.Memory.UpdateBalance(t.Reciever, t.Amount)

	// save it into datastore
	c.Memory.AddTransactionToDatastore(t)
}

// beginConsensus builds and starts a new apaxos instance.
func (c Consensus) beginConsensus(transaction *apaxos.Transaction) {
	// first we create a new instance for the protocol
	c.instance = &protocol.Apaxos{
		NodeId:          c.NodeId,
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
		// after this go-routine finished, clear the protocol instance
		defer func() {
			// terminate consensus by resetting the instance
			c.instance = nil
			c.Logger.Debug("apaxos terminated")
		}()

		c.Logger.Debug("apaxos started")

		// in a while loop, try to make consensus
		for {
			// start apaxos protocol
			err := c.instance.Start()
			if err != nil {
				c.Logger.Error("apaxos failed", zap.Error(err))

				// check the error for a proper handling mechanism
				switch {
				case errors.Is(err, protocol.ErrRequestTimeout):
					// in this case, we first check to see if we have enough servers up and running or not
					if !c.livenessHandler() {
						c.failedConsensus(protocol.ErrNotEnoughServers)
						return
					}
				case errors.Is(err, protocol.ErrSlowNode):
					// if we hit slow-node, it means that we got synced, so we should check the status of balance before
					// rerunning the consensus protocol.
					c.Logger.Info("slow server detected")

					// if the client balance is now enough, then we submit the transaction
					if c.recheckBalance(transaction) {
						c.submitTransaction(transaction)
						break
					}
				default:
					c.Logger.Error("consensus error", zap.Error(err))
				}
			} else { // a successful consensus
				// now we check to see if the client balance is enough or not
				if c.recheckBalance(transaction) {
					c.successfulConsensus(transaction)
				} else {
					c.failedConsensus(protocol.ErrNotEnoughBalance)
				}

				return
			}
		}
	}()
}

// failed consensus tells the client about the consensus failure and
func (c Consensus) failedConsensus(err error) {
	c.instance.OutChannel <- &messages.Packet{
		Payload: err.Error(),
	}
}

// recheck balance checks to see if the client balance has changed
// after the sync or not.
func (c Consensus) recheckBalance(t *apaxos.Transaction) bool {
	return c.Memory.GetBalance(c.Client) > t.Amount
}

// successful consensus submits the transaction and notifys the client.
func (c Consensus) successfulConsensus(t *apaxos.Transaction) {
	c.submitTransaction(t)

	c.instance.OutChannel <- &messages.Packet{
		Payload: "transaction submitted",
	}
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
