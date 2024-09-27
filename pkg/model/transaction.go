package model

import "gorm.io/gorm"

// Transaction represents the system's data.
// Each transaction is signed by their node provider. They have
// have a sender, reciever, and an amount.
// Other data is used for storing them in database.
type Transaction struct {
	gorm.Model
	Sign     string
	Sender   string
	Reciever string
	Amount   int
}
