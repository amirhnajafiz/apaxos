package models

import "github.com/f24-cse535/apaxos/pkg/rpc/apaxos"

// Transaction model is a struct that acts as
// a data-model for MongoDB database.
type Transaction struct {
	Sender         string `bson:"sender"`
	Reciever       string `bson:"reciever"`
	Amount         int64  `bson:"amount"`
	SequenceNumber int64  `bson:"sequence_number"`
}

// The following methods are being used to cast our data-model
// to proto-model which is used in RPC calls.
// Each model comes with two methods to create proto-model from
// the existing model, and a build a data-model from the given proto-model.

func (t Transaction) ToProtoModel() *apaxos.Transaction {
	return &apaxos.Transaction{
		Sender:         t.Sender,
		Reciever:       t.Reciever,
		Amount:         t.Amount,
		SequenceNumber: t.SequenceNumber,
	}
}

func (t Transaction) FromProtoModel(instance *apaxos.Transaction) Transaction {
	t.Sender = instance.GetSender()
	t.Reciever = instance.GetReciever()
	t.Amount = instance.GetAmount()
	t.SequenceNumber = instance.GetSequenceNumber()

	return t
}
