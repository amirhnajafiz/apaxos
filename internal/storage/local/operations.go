package local

import "github.com/f24-cse535/apaxos/pkg/models"

// Memory operations for sequence_number.
func (m *Memory) GetSequenceNumber() int64 {
	return m.sequenceNumber
}

func (m *Memory) IncSequenceNumber() {
	m.sequenceNumber++
}

// Memory operations for clients.
func (m *Memory) SetBalance(client string, balance int64) {
	m.clients[client] = balance
}

func (m *Memory) GetBalance(client string) int64 {
	return m.clients[client]
}

func (m *Memory) GetClients() map[string]int64 {
	return m.clients
}

// Memory operations for ballot_number.
func (m *Memory) SetBallotNumber(instance *models.BallotNumber) {
	m.ballotNumber = instance
}

func (m *Memory) GetBallotNumber() *models.BallotNumber {
	return m.ballotNumber
}

// Memory operations for accepted_num.
func (m *Memory) SetAcceptedNum(instance *models.BallotNumber) {
	m.acceptedNum = instance
}

func (m *Memory) GetAcceptedNum() *models.BallotNumber {
	return m.acceptedNum
}

// Memory operations for accepted_val.
func (m *Memory) SetAcceptedVal(instance []*models.Block) {
	m.acceptedVal = instance
}

func (m *Memory) GetAcceptedVal() []*models.Block {
	return m.acceptedVal
}

// Memory operations for datastore.
func (m *Memory) AddTransactionToDatastore(instance *models.Transaction) {
	m.datastore = append(m.datastore, instance)
}

func (m *Memory) GetDatastore() []*models.Transaction {
	return m.datastore
}

func (m *Memory) ResetDatastore() {
	m.datastore = make([]*models.Transaction, 0)
}
