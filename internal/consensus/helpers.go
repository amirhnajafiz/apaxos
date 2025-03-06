package consensus

import (
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"

	"go.uber.org/zap"
)

// submitTransaction is used to store a transaction into datastore, and update balances.
func (c *Consensus) submitTransaction(transaction *apaxos.Transaction, store bool) {
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
	if store {
		c.Memory.AddTransactionToDatastore(transaction)
	}
}

// checkBalance checks to see if the client balance has changed or not.
func (c *Consensus) checkBalance(t *apaxos.Transaction) bool {
	return c.Memory.GetBalance(c.Client) >= t.GetAmount()
}

// instance exists return true if the apaxos instance is started and running.
func (c *Consensus) instanceExists() bool {
	return c.instance != nil
}
