package local

import (
	"time"

	"github.com/f24-cse535/apaxos/pkg/models"
)

// Memory is a local storage that is used for each node. It keeps the state of each node.
type Memory struct {
	sequenceNumber int64                // is used for ordering transactions within each block
	clients        map[string]int64     // a map of all clients with their balances
	ballotNumber   *models.BallotNumber // the ballot number of each node

	acceptedNum *models.BallotNumber // apaxos accepted num
	acceptedVal []*models.Block      // apaxos accepted var

	datastore []*models.Transaction // local transactions datastore for each node
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
		sequenceNumber: time.Now().UnixMilli(), // initalized by the system booting timestamp
		clients:        clientsHashmap,

		acceptedNum: nil,
		acceptedVal: make([]*models.Block, 0),

		datastore: make([]*models.Transaction, 0),

		ballotNumber: &models.BallotNumber{ // initialize the ballot number for each node
			Number: 0,
			NodeId: nodeId,
		},
	}
}
