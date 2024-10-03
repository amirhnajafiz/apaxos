package models

import "github.com/f24-cse535/apaxos/pkg/transactions"

// Block model is a struct that acts as
// a data-model for MongoDB database.
type Block struct {
	Transactions   []*Transaction `bson:"transactions"`
	NodeID         string         `bson:"node_id"`
	UID            string         `bson:"uid"`
	SequenceNumber int64          `bson:"sequence_number"`
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
		Transactions:   list,
		NodeId:         b.NodeID,
		Uid:            b.UID,
		SequenceNumber: b.SequenceNumber,
	}
}

func (b *Block) FromProtoModel(instance *transactions.Block) {
	list := make([]*Transaction, len(instance.Transactions))

	for index, value := range instance.Transactions {
		tmp := &Transaction{}
		tmp.FromProtoModel(value)

		list[index] = tmp
	}

	b.Transactions = list
	b.NodeID = instance.GetNodeId()
	b.UID = instance.GetUid()
	b.SequenceNumber = instance.GetSequenceNumber()
}
