package local

import "github.com/f24-cse535/apaxos/pkg/models"

// Memory operations for sequence_number.
func (m *Memory) GetSequenceNumber() int64 {
	tmp := m.sequenceNumber

	// inc sequence number after each read
	m.IncSequenceNumber()

	return tmp
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
	// process the transaction before adding it to the datastore
	m.clients[instance.Sender] = m.clients[instance.Sender] - instance.Amount
	m.clients[instance.Reciever] = m.clients[instance.Reciever] + instance.Amount

	m.datastore = append(m.datastore, instance)
}

func (m *Memory) GetDatastore() []*models.Transaction {
	return m.datastore
}

func (m *Memory) SetDatastore(instance []*models.Transaction) {
	m.datastore = instance
}
