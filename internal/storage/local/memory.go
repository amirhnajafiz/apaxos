package local

import "github.com/f24-cse535/apaxos/pkg/models"

// Memory is a local storage that is used
// for each node. It keeps the state of each node.
type Memory struct {
	sequenceNumber int64
	clients        map[string]int64
	ballotNumber   *models.BallotNumber
	acceptedNum    *models.BallotNumber
	acceptedVal    []*models.Block
	datastore      []*models.Transaction
}

// NewMemory returns an instance of the memory struct.
// It accepts the node_id, an initiali balance value, and the list of clients.
func NewMemory(nodeId string, initBalanceValue int64, clients ...string) *Memory {
	// this hashmap is used to store clients and their balance
	clientsHashmap := make(map[string]int64, len(clients))
	for _, client := range clients {
		clientsHashmap[client] = initBalanceValue
	}

	return &Memory{
		sequenceNumber: 0,
		clients:        clientsHashmap,
		acceptedVal:    make([]*models.Block, 0),
		datastore:      make([]*models.Transaction, 0),
		ballotNumber: &models.BallotNumber{ // initialize the ballot number for each node
			Number: 0,
			NodeId: nodeId,
		},
	}
}
