package models

import (
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"
)

// Block metadata model is a struct that acts as
// a data-model for MongoDB database. It is also used
// for internal processing.
type BlockMetaData struct {
	NodeId       string       `bson:"node_id"`
	BallotNumber BallotNumber `bson:"ballot_number"`
}

// Block model is a struct that acts as
// a data-model for MongoDB database. It is also used
// for internal processing.
type Block struct {
	Metadata     BlockMetaData
	Transactions []Transaction `bson:"transactions"`
}

// The following methods are being used to cast our data-model
// to proto-model which is used in RPC calls.
// Each model comes with two methods to create proto-model from
// the existing model, and a build a data-model from the given proto-model.

func (b *BlockMetaData) ToProtoModel() *apaxos.BlockMetaData {
	return &apaxos.BlockMetaData{
		NodeId:       b.NodeId,
		BallotNumber: b.BallotNumber.ToProtoModel(),
	}
}

func (b *Block) ToProtoModel() *apaxos.Block {
	list := make([]*apaxos.Transaction, len(b.Transactions))
	for index, value := range b.Transactions {
		list[index] = value.ToProtoModel()
	}

	return &apaxos.Block{
		Metadata:     b.Metadata.ToProtoModel(),
		Transactions: list,
	}
}

func (b *BlockMetaData) FromProtoModel(instance *apaxos.BlockMetaData) {
	b.NodeId = instance.GetNodeId()

	b.BallotNumber = BallotNumber{}
	b.BallotNumber.FromProtoModel(instance.BallotNumber)
}

func (b *Block) FromProtoModel(instance *apaxos.Block) {
	// initialize the Transactions slice
	list := make([]Transaction, len(instance.Transactions))

	for index, value := range instance.Transactions {
		// properly initialize each Transaction before calling FromProtoModel
		t := Transaction{}
		t.FromProtoModel(value)
		list[index] = t
	}

	b.Transactions = list

	// ensure that b.Metadata is initialized before calling FromProtoModel
	b.Metadata = BlockMetaData{} // assuming Metadata is a pointer type
	b.Metadata.FromProtoModel(instance.GetMetadata())
}
