package models

import "github.com/f24-cse535/apaxos/pkg/transactions"

// BallotNumber model is a struct that presents
// PAXOS ballot_number.
type BallotNumber struct {
	Number int64  `bson:"number"`
	NodeId string `bson:"node_id"`
}

// The following methods are being used to cast our data-model
// to proto-model which is used in RPC calls.
// Each model comes with two methods to create proto-model from
// the existing model, and a build a data-model from the given proto-model.

func (b *BallotNumber) ToProtoModel() *transactions.BallotNumber {
	return &transactions.BallotNumber{
		Number: b.Number,
		NodeId: b.NodeId,
	}
}

func (b *BallotNumber) FromProtoModel(instance *transactions.BallotNumber) {
	b.Number = instance.GetNumber()
	b.NodeId = instance.GetNodeId()
}
