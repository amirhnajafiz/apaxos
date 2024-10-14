package models

import "github.com/f24-cse535/apaxos/pkg/rpc/apaxos"

// Block metadata model is a struct that acts as
// a data-model for MongoDB database. It is also used
// for internal processing.
type BlockMetadata struct {
	NodeId       string       `bson:"node_id"`
	BallotNumber BallotNumber `bson:"ballot_number"`
}

// Block model is a struct that acts as
// a data-model for MongoDB database. It is also used
// for internal processing.
type Block struct {
	Metadata     BlockMetadata
	Transactions []Transaction `bson:"transactions"`
}

// The following methods are being used to cast our data-model
// to proto-model which is used in RPC calls.
// Each model comes with two methods to create proto-model from
// the existing model, and a build a data-model from the given proto-model.

func (b BlockMetadata) ToProtoModel() *apaxos.BlockMetaData {
	return &apaxos.BlockMetaData{
		NodeId:       b.NodeId,
		BallotNumber: b.BallotNumber.ToProtoModel(),
	}
}

func (b Block) ToProtoModel() *apaxos.Block {
	list := make([]*apaxos.Transaction, len(b.Transactions))
	for index, value := range b.Transactions {
		list[index] = value.ToProtoModel()
	}

	return &apaxos.Block{
		Metadata:     b.Metadata.ToProtoModel(),
		Transactions: list,
	}
}

func (b BlockMetadata) FromProtoModel(instance *apaxos.BlockMetaData) BlockMetadata {
	b.NodeId = instance.GetNodeId()
	b.BallotNumber = BallotNumber{}.FromProtoModel(instance.BallotNumber)

	return b
}

func (b Block) FromProtoModel(instance *apaxos.Block) Block {
	list := make([]Transaction, len(instance.Transactions))

	for index, value := range instance.Transactions {
		list[index] = Transaction{}.FromProtoModel(value)
	}

	b.Transactions = list
	b.Metadata = BlockMetadata{}.FromProtoModel(instance.GetMetadata())

	return b
}
