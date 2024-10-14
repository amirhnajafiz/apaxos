package consensus

import (
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
		c.Logger.Debug("apaxos started")

		// start apaxos
		if err := c.instance.Start(); err != nil {
			c.Logger.Error("apaxos failed", zap.Error(err))
		}

		// submit the transaction after apaxos
		c.submitTransaction(transaction)

		// reset the instance
		c.instance = nil
	}()
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
