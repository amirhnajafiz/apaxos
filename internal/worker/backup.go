package worker

import "github.com/f24-cse535/apaxos/pkg/models"

// snapshot nodes state, get's the current state of node's memory, and stores it inside MongoDB.
func (w Worker) snapshotNodeState() error {
	// build a new state (snapshot)
	state := models.State{
		Clients:              w.Memory.GetClients(),
		LastCommittedMessage: *w.Memory.GetLastCommittedMessage(),
		BallotNumber:         *w.Memory.GetBallotNumber(),
		AcceptedNum:          *w.Memory.GetAcceptedNum(),
		Datastore:            *w.Memory.GetDatastore(),
	}

	// copy accepted_val and current datastore
	vals := w.Memory.GetAcceptedVal()

	// converting addresses to values
	state.AcceptedVal = make([]models.Block, len(vals))
	for index, item := range vals {
		state.AcceptedVal[index] = *item
	}

	// store the snapshot in MongoDB
	return w.Database.InsertState(&state)
}
