package models

import "github.com/f24-cse535/apaxos/pkg/transactions"

// Transaction model is a struct that acts as
// a data-model for MongoDB database.
type Transaction struct {
	Uid      string `bson:"uid"`
	Sender   string `bson:"sender"`
	Reciever string `bson:"reciever"`
	Amount   int64  `bson:"amount"`
}

// The following methods are being used to cast our data-model
// to proto-model which is used in RPC calls.
// Each model comes with two methods to create proto-model from
// the existing model, and a build a data-model from the given proto-model.

func (t Transaction) ToProtoModel() *transactions.Transaction {
	return &transactions.Transaction{
		Uid:      t.Uid,
		Sender:   t.Sender,
		Reciever: t.Reciever,
		Amount:   t.Amount,
	}
}

func (t Transaction) FromProtoModel(instance *transactions.Transaction) Transaction {
	t.Uid = instance.GetUid()
	t.Sender = instance.GetSender()
	t.Reciever = instance.GetReciever()
	t.Amount = instance.GetAmount()

	return t
}
