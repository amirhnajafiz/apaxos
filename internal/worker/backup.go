package worker

import "github.com/f24-cse535/apaxos/pkg/models"

// snapshot nodes state, get's the current state of node's memory, and stores it inside MongoDB.
func (w Worker) snapshotNodeState() error {
	// build a new state (snapshot)
	state := models.State{
		Clients:              w.Memory.GetClients(),
		LastCommittedMessage: models.BallotNumber{},
		BallotNumber:         models.BallotNumber{},
		AcceptedNum:          models.BallotNumber{},
		Datastore:            models.Block{},
	}

	// set the state values from memory reads
	state.LastCommittedMessage.FromProtoModel(w.Memory.GetLastCommittedMessage())
	state.BallotNumber.FromProtoModel(w.Memory.GetBallotNumber())
	state.AcceptedNum.FromProtoModel(w.Memory.GetAcceptedNum())
	state.Datastore.FromProtoModel(w.Memory.GetDatastore())

	// copy accepted_val and current datastore
	vals := w.Memory.GetAcceptedVal()

	// converting addresses to values
	state.AcceptedVal = make([]models.Block, len(vals))
	for index, item := range vals {
		state.AcceptedVal[index].FromProtoModel(item)
	}

	// store the snapshot in MongoDB
	return w.Database.InsertState(&state)
}
