package local

import (
	"time"

	"github.com/f24-cse535/apaxos/pkg/models"
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"
)

// Memory is a local storage that is used for each node. It keeps the state of each node.
type Memory struct {
	sequenceNumber int64 // is used for ordering transactions within each block

	clients map[string]int64 // a map of all clients with their balances

	ballotNumber  *apaxos.BallotNumber // the ballot number of each node
	lastCommitted *apaxos.BallotNumber // apaxos last committed message for sync
	acceptedNum   *apaxos.BallotNumber // apaxos accepted num

	datastore   *apaxos.Block   // local transactions datastore for each node
	acceptedVal []*apaxos.Block // apaxos accepted var
}

// ReadFromState is used to load a backup state into the current memory.
func (m *Memory) ReadFromState(state *models.State) {
	m.clients = state.Clients

	m.lastCommitted = state.LastCommittedMessage.ToProtoModel()
	m.ballotNumber = state.BallotNumber.ToProtoModel()
	m.acceptedNum = state.AcceptedNum.ToProtoModel()
	m.datastore = state.Datastore.ToProtoModel()

	m.acceptedVal = make([]*apaxos.Block, len(state.AcceptedVal))
	for index, item := range state.AcceptedVal {
		m.acceptedVal[index] = item.ToProtoModel()
	}
}

// NewMemory returns an instance of the memory struct.
// It accepts the node_id, an initiali balance value within the list of clients.
func NewMemory(nodeId string, balances map[string]int64) *Memory {
	return &Memory{
		sequenceNumber: time.Now().UnixMilli(), // initalized by the system booting timestamp
		clients:        balances,

		acceptedNum: nil,

		acceptedVal: make([]*apaxos.Block, 0),
		datastore: &apaxos.Block{
			Metadata: &apaxos.BlockMetaData{
				NodeId: nodeId,
			},
			Transactions: make([]*apaxos.Transaction, 0),
		},

		lastCommitted: &apaxos.BallotNumber{ // initialize the last committed as zero
			Number: 0,
			NodeId: "",
		},

		ballotNumber: &apaxos.BallotNumber{ // initialize the ballot number for each node
			Number: 0,
			NodeId: nodeId,
		},
	}
}
