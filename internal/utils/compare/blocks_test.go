package compare_test

import (
	"testing"

	"github.com/f24-cse535/apaxos/internal/utils/compare"
	"github.com/f24-cse535/apaxos/pkg/models"
)

func TestCompareBlocks(t *testing.T) {
	tests := []struct {
		name     string
		blockA   *models.BlockMetadata
		blockB   *models.BlockMetadata
		expected bool
	}{
		{
			name: "Higher BallotNumber.Number in blockA",
			blockA: &models.BlockMetadata{
				BallotNumber:   models.BallotNumber{Number: 2, NodeId: "1"},
				SequenceNumber: 1,
			},
			blockB: &models.BlockMetadata{
				BallotNumber:   models.BallotNumber{Number: 1, NodeId: "1"},
				SequenceNumber: 1,
			},
			expected: false,
		},
		{
			name: "Higher BallotNumber.NodeId in blockA",
			blockA: &models.BlockMetadata{
				BallotNumber:   models.BallotNumber{Number: 1, NodeId: "2"},
				SequenceNumber: 1,
			},
			blockB: &models.BlockMetadata{
				BallotNumber:   models.BallotNumber{Number: 1, NodeId: "1"},
				SequenceNumber: 1,
			},
			expected: false,
		},
		{
			name: "Higher SequenceNumber in blockA",
			blockA: &models.BlockMetadata{
				BallotNumber:   models.BallotNumber{Number: 1, NodeId: "1"},
				SequenceNumber: 2,
			},
			blockB: &models.BlockMetadata{
				BallotNumber:   models.BallotNumber{Number: 1, NodeId: "1"},
				SequenceNumber: 1,
			},
			expected: false,
		},
		{
			name: "Same BallotNumber and SequenceNumber",
			blockA: &models.BlockMetadata{
				BallotNumber:   models.BallotNumber{Number: 1, NodeId: "1"},
				SequenceNumber: 1,
			},
			blockB: &models.BlockMetadata{
				BallotNumber:   models.BallotNumber{Number: 1, NodeId: "1"},
				SequenceNumber: 1,
			},
			expected: true,
		},
		{
			name: "Lower BallotNumber.Number in blockA",
			blockA: &models.BlockMetadata{
				BallotNumber:   models.BallotNumber{Number: 1, NodeId: "1"},
				SequenceNumber: 1,
			},
			blockB: &models.BlockMetadata{
				BallotNumber:   models.BallotNumber{Number: 2, NodeId: "1"},
				SequenceNumber: 1,
			},
			expected: true,
		},
		{
			name: "Lower BallotNumber.NodeId in blockA",
			blockA: &models.BlockMetadata{
				BallotNumber:   models.BallotNumber{Number: 1, NodeId: "1"},
				SequenceNumber: 1,
			},
			blockB: &models.BlockMetadata{
				BallotNumber:   models.BallotNumber{Number: 1, NodeId: "2"},
				SequenceNumber: 1,
			},
			expected: true,
		},
		{
			name: "Lower SequenceNumber in blockA",
			blockA: &models.BlockMetadata{
				BallotNumber:   models.BallotNumber{Number: 1, NodeId: "1"},
				SequenceNumber: 1,
			},
			blockB: &models.BlockMetadata{
				BallotNumber:   models.BallotNumber{Number: 1, NodeId: "1"},
				SequenceNumber: 2,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compare.CompareBlocks(tt.blockA, tt.blockB)
			if result != tt.expected {
				t.Errorf("CompareBlocks() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
