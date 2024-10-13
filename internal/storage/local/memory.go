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

// ReadFromState is used to load a backup state into the current memory.
func (m *Memory) ReadFromState(state *models.State) {
	m.clients = state.Clients
	m.ballotNumber = &state.BallotNumber
	m.acceptedNum = &state.AcceptedNum

	m.acceptedVal = make([]*models.Block, len(state.AcceptedVal))
	for index, item := range state.AcceptedVal {
		m.acceptedVal[index] = &item
	}

	m.datastore = make([]*models.Transaction, len(state.Datastore))
	for index, item := range state.Datastore {
		m.datastore[index] = &item
	}
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
