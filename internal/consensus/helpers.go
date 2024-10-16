package consensus

import (
	"github.com/f24-cse535/apaxos/pkg/messages"
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"

	"go.uber.org/zap"
)

// submitTransaction is used to store a transaction into datastore, and update balances.
func (c Consensus) submitTransaction(transaction *apaxos.Transaction) {
	// set a unique sequence number for transaction
	transaction.SequenceNumber = c.Memory.GetSequenceNumber()

	c.Logger.Info(
		"transaction submitted",
		zap.Int64("amount", transaction.GetAmount()),
		zap.Int64("seq_number", transaction.GetSequenceNumber()),
		zap.String("sender", transaction.GetSender()),
		zap.String("receiver", transaction.GetReciever()),
	)

	// make changes into memory for client's balances
	c.Memory.UpdateBalance(transaction.GetSender(), transaction.GetAmount()*-1)
	c.Memory.UpdateBalance(transaction.GetReciever(), transaction.GetAmount())

	// save it into datastore
	c.Memory.AddTransactionToDatastore(transaction)
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
	return c.Memory.GetBalance(c.Client) > t.GetAmount()
}

// instance exists return true if the apaxos instance is started and running
func (c Consensus) instanceExists() bool {
	return c.instance == nil
}
