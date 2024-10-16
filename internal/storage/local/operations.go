package local

import (
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"
)

// SetServiceStatus sets a new status for service.
func (m *Memory) SetServiceStatus(status bool) {
	m.serviceStatus = status
}

// GetServiceStatus returns the current service status.
func (m *Memory) GetServiceStatus() bool {
	return m.serviceStatus
}

// GetSequenceNumber is used for labeling transactions
// inside each node to keep an order of them.
func (m *Memory) GetSequenceNumber() int64 {
	tmp := m.sequenceNumber

	// inc sequence number after each read
	m.IncSequenceNumber()

	return tmp
}

// IncSequenceNumber increaments the sequence_number by one.
func (m *Memory) IncSequenceNumber() {
	m.sequenceNumber++
}

// SetBalance is used to reset a balance value to the given balance.
func (m *Memory) SetBalance(client string, balance int64) {
	m.clients[client] = balance
}

// UpdateBalance is used to change the balance of a client by adding an amount to it.
func (m *Memory) UpdateBalance(client string, amount int64) {
	m.clients[client] = m.clients[client] + amount
}

// GetBalance is used to get the balance of a client.
func (m *Memory) GetBalance(client string) int64 {
	return m.clients[client]
}

// GetClients returns the list of clients with balances.
func (m *Memory) GetClients() map[string]int64 {
	return m.clients
}

// SetBallotNumber updates the current ballot-number.
func (m *Memory) SetBallotNumber(instance *apaxos.BallotNumber) {
	m.ballotNumber = instance
}

// GetBallotNumber is used to return the current ballot-number.
func (m *Memory) GetBallotNumber() *apaxos.BallotNumber {
	return m.ballotNumber
}

// SetAcceptedNum is used to update the current accepted_num.
func (m *Memory) SetAcceptedNum(instance *apaxos.BallotNumber) {
	m.acceptedNum = instance
}

// GetAcceptedNum is used to return the current accepted_num.
func (m *Memory) GetAcceptedNum() *apaxos.BallotNumber {
	return m.acceptedNum
}

// SetAcceptedVal is used to update the current accepted_val.
func (m *Memory) SetAcceptedVal(instance []*apaxos.Block) {
	m.acceptedVal = instance
}

// GetAcceptedVal is used to return the current accepted_val.
func (m *Memory) GetAcceptedVal() []*apaxos.Block {
	return m.acceptedVal
}

// AddTransactionToDatastore stores a transaction into datastore.
func (m *Memory) AddTransactionToDatastore(instance *apaxos.Transaction) {
	m.datastore.Transactions = append(m.datastore.Transactions, instance)
}

// ClearDatastore gets a block and removes the transactions from datastore
// that are inside that block.
func (m *Memory) ClearDatastore(instance *apaxos.Block) {
	// create a map to store elements of block for quick lookup
	hashMap := make(map[int64]bool)
	for _, transaction := range instance.Transactions {
		hashMap[transaction.SequenceNumber] = true
	}

	// create a new datastore
	datastore := make([]*apaxos.Transaction, 0)
	for _, transaction := range m.datastore.Transactions {
		// add transactions that are not in the given block
		if !hashMap[transaction.SequenceNumber] {
			datastore = append(datastore, transaction)
		}
	}

	// reset the datastore
	m.SetDatastore(datastore)
}

// GetDatastore returns the datastore block.
func (m *Memory) GetDatastore() *apaxos.Block {
	return m.datastore
}

// SetDatastore only updates the datastore transactions list.
func (m *Memory) SetDatastore(instance []*apaxos.Transaction) {
	m.datastore.Transactions = instance
}

// SetLastCommittedMessage updates last committed value.
func (m *Memory) SetLastCommittedMessage(instance *apaxos.BallotNumber) {
	m.lastCommitted = instance
}

// GetLastCommittedMessage returns the current last committed value.
func (m *Memory) GetLastCommittedMessage() *apaxos.BallotNumber {
	return m.lastCommitted
}
