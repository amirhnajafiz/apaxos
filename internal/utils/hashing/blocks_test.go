package hashing_test

import (
	"testing"

	"github.com/f24-cse535/apaxos/internal/utils/hashing"
	"github.com/f24-cse535/apaxos/pkg/models"
)

func TestHashBlock(t *testing.T) {
	tests := []struct {
		name     string
		block    *models.Block
		expected string
	}{
		{
			name: "Block with one transaction",
			block: &models.Block{
				Metadata: models.BlockMetadata{
					NodeId:         "node-1",
					SequenceNumber: 10,
				},
				Transactions: []models.Transaction{
					{Uid: "tx-1"},
				},
			},
			expected: "node-1-10-1",
		},
		{
			name: "Block with multiple transactions",
			block: &models.Block{
				Metadata: models.BlockMetadata{
					NodeId:         "node-2",
					SequenceNumber: 20,
				},
				Transactions: []models.Transaction{
					{Uid: "tx-1"},
					{Uid: "tx-2"},
				},
			},
			expected: "node-2-20-2",
		},
		{
			name: "Block with no transactions",
			block: &models.Block{
				Metadata: models.BlockMetadata{
					NodeId:         "node-3",
					SequenceNumber: 30,
				},
				Transactions: []models.Transaction{},
			},
			expected: "node-3-30-0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hashing.HashBlock(tt.block)
			if result != tt.expected {
				t.Errorf("HashBlock() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
