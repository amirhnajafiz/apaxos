package models

// State holds the memory information of each node. It is used
// to keep a snapshot of the node inside MongoDB.
type State struct {
	Clients              map[string]int64 `bson:"clients"`
	LastCommittedMessage BallotNumber     `bson:"last_committed_message"`
	BallotNumber         BallotNumber     `bson:"ballot_number"`
	AcceptedNum          BallotNumber     `bson:"accepted_num"`
	Datastore            Block            `bson:"datastore"`
	AcceptedVal          []Block          `bson:"accepted_val"`
}
