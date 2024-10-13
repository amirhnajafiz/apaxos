package hashing_test

import (
	"testing"

	"github.com/f24-cse535/apaxos/internal/utils/hashing"
	"github.com/f24-cse535/apaxos/pkg/models"
)

func TestHashTransaction(t *testing.T) {
	tests := []struct {
		name        string
		transaction *models.Transaction
		expected    string
	}{
		{
			name: "Basic transaction",
			transaction: &models.Transaction{
				Sender:         "Alice",
				Reciever:       "Bob",
				SequenceNumber: 1,
			},
			expected: "Alice-Bob-1",
		},
		{
			name: "Transaction with larger sequence number",
			transaction: &models.Transaction{
				Sender:         "Charlie",
				Reciever:       "Dave",
				SequenceNumber: 100,
			},
			expected: "Charlie-Dave-100",
		},
		{
			name: "Transaction with empty fields",
			transaction: &models.Transaction{
				Sender:         "",
				Reciever:       "",
				SequenceNumber: 0,
			},
			expected: "--0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hashing.HashTransaction(tt.transaction)
			if result != tt.expected {
				t.Errorf("HashTransaction() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
