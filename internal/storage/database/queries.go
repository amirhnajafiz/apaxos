package database

import (
	"fmt"

	"github.com/f24-cse535/apaxos/pkg/model"
)

func (d Database) InsertTransactions(transactions []*model.Transaction) error {
	return d.db.CreateInBatches(transactions, len(transactions)).Error
}

func (d Database) GetTransactions() ([]*model.Transaction, error) {
	var transactions []*model.Transaction

	if err := d.db.Model(&model.Transaction{}).Find(transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to query transactions: %v", err)
	}

	return transactions, nil
}
