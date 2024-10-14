package local

import (
	"time"

	"github.com/f24-cse535/apaxos/pkg/models"
)

// Memory is a local storage that is used for each node. It keeps the state of each node.
type Memory struct {
	sequenceNumber int64 // is used for ordering transactions within each block

	clients map[string]int64 // a map of all clients with their balances

	ballotNumber  *models.BallotNumber // the ballot number of each node
	lastCommitted *models.BallotNumber // apaxos last committed message for sync
	acceptedNum   *models.BallotNumber // apaxos accepted num

	datastore   *models.Block   // local transactions datastore for each node
	acceptedVal []*models.Block // apaxos accepted var
}

// ReadFromState is used to load a backup state into the current memory.
func (m *Memory) ReadFromState(state *models.State) {
	m.clients = state.Clients

	m.lastCommitted = &state.LastCommittedMessage
	m.ballotNumber = &state.BallotNumber
	m.acceptedNum = &state.AcceptedNum

	m.acceptedVal = make([]*models.Block, len(state.AcceptedVal))
	for index, item := range state.AcceptedVal {
		m.acceptedVal[index] = &item
	}

	m.datastore = &state.Datastore
}

// NewMemory returns an instance of the memory struct.
// It accepts the node_id, an initiali balance value within the list of clients.
func NewMemory(nodeId string, balances map[string]int64) *Memory {
	return &Memory{
		sequenceNumber: time.Now().UnixMilli(), // initalized by the system booting timestamp
		clients:        balances,

		acceptedNum: nil,

		acceptedVal: make([]*models.Block, 0),
		datastore: &models.Block{
			Metadata: models.BlockMetadata{
				NodeId: nodeId,
			},
			Transactions: make([]models.Transaction, 0),
		},

		lastCommitted: &models.BallotNumber{
			Number: 0,
			NodeId: "",
		},

		ballotNumber: &models.BallotNumber{ // initialize the ballot number for each node
			Number: 0,
			NodeId: nodeId,
		},
	}
}
