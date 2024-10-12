package local

import "github.com/f24-cse535/apaxos/pkg/models"

// Memory is a local storage that is used
// for each node.
type Memory struct {
	sequenceNumber int64
	ballotNumber   *models.BallotNumber
	acceptedNum    *models.BallotNumber
	acceptedVal    []*models.Block
	datastore      []*models.Transaction
}

// NewMemory returns an instance of the memory struct.
func NewMemory() *Memory {
	return &Memory{
		sequenceNumber: 0,
		acceptedVal:    make([]*models.Block, 0),
		datastore:      make([]*models.Transaction, 0),
	}
}
