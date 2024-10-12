package models

import "github.com/f24-cse535/apaxos/pkg/transactions"

// Block model is a struct that acts as
// a data-model for MongoDB database. It is also used
// for internal processing.
type Block struct {
	Uid            string         `bson:"uid"`
	NodeId         string         `bson:"node_id"`
	SequenceNumber int64          `bson:"sequence_number"`
	BallotNumber   *BallotNumber  `bson:"ballot_number"`
	Transactions   []*Transaction `bson:"transactions"`
}

// The following methods are being used to cast our data-model
// to proto-model which is used in RPC calls.
// Each model comes with two methods to create proto-model from
// the existing model, and a build a data-model from the given proto-model.

func (b *Block) ToProtoModel() *transactions.Block {
	list := make([]*transactions.Transaction, len(b.Transactions))
	for index, value := range b.Transactions {
		list[index] = value.ToProtoModel()
	}

	return &transactions.Block{
		Uid:            b.Uid,
		NodeId:         b.NodeId,
		SequenceNumber: b.SequenceNumber,
		BallotNumber:   b.BallotNumber.ToProtoModel(),
		Transactions:   list,
	}
}

func (b *Block) FromProtoModel(instance *transactions.Block) {
	list := make([]*Transaction, len(instance.Transactions))
	for index, value := range instance.Transactions {
		tmp := &Transaction{}
		tmp.FromProtoModel(value)

		list[index] = tmp
	}

	ballotNumber := &BallotNumber{}
	ballotNumber.FromProtoModel(instance.BallotNumber)

	b.Transactions = list
	b.NodeId = instance.GetNodeId()
	b.Uid = instance.GetUid()
	b.SequenceNumber = instance.GetSequenceNumber()
	b.BallotNumber = ballotNumber
}
