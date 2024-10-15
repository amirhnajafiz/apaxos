package consensus

import (
	"github.com/f24-cse535/apaxos/pkg/messages"
	"github.com/f24-cse535/apaxos/pkg/models"
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"

	"go.uber.org/zap"
)

// submitTransaction is used to store a transaction into datastore, and update balances.
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

// failed consensus tells the client about the consensus failure.
func (c Consensus) failedConsensus(err error) {
	c.instance.OutChannel <- &messages.Packet{
		Payload: err.Error(),
	}
}

// successful consensus submits the transaction and notifys the client.
func (c Consensus) successfulConsensus(t *apaxos.Transaction) {
	c.submitTransaction(t)

	c.instance.OutChannel <- &messages.Packet{
		Payload: "transaction submitted",
	}
}

// recheck balance checks to see if the client balance has changed or not.
func (c Consensus) recheckBalance(t *apaxos.Transaction) bool {
	return c.Memory.GetBalance(c.Client) > t.Amount
}

// forward to instance is a helper function to pass packets in apaxos instance channel.
func (c Consensus) forwardToInstance(pkt *messages.Packet) {
	if c.instance != nil {
		c.instance.InChannel <- pkt
	}
}

// instance exists return true if the apaxos instance is started and running
func (c Consensus) instanceExists() bool {
	return c.instance == nil
}
