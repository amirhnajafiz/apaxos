package utils

import (
	"testing"

	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"
)

func TestCompareBlocks(t *testing.T) {
	tests := []struct {
		name     string
		blockA   *apaxos.BlockMetaData
		blockB   *apaxos.BlockMetaData
		expected bool
	}{
		{
			name: "Higher BallotNumber.Number in blockA",
			blockA: &apaxos.BlockMetaData{
				BallotNumber: &apaxos.BallotNumber{Number: 2, NodeId: "1"},
			},
			blockB: &apaxos.BlockMetaData{
				BallotNumber: &apaxos.BallotNumber{Number: 1, NodeId: "1"},
			},
			expected: false,
		},
		{
			name: "Higher BallotNumber.NodeId in blockA",
			blockA: &apaxos.BlockMetaData{
				BallotNumber: &apaxos.BallotNumber{Number: 1, NodeId: "2"},
			},
			blockB: &apaxos.BlockMetaData{
				BallotNumber: &apaxos.BallotNumber{Number: 1, NodeId: "1"},
			},
			expected: false,
		},
		{
			name: "Same BallotNumber",
			blockA: &apaxos.BlockMetaData{
				BallotNumber: &apaxos.BallotNumber{Number: 1, NodeId: "1"},
			},
			blockB: &apaxos.BlockMetaData{
				BallotNumber: &apaxos.BallotNumber{Number: 1, NodeId: "1"},
			},
			expected: true,
		},
		{
			name: "Lower BallotNumber.Number in blockA",
			blockA: &apaxos.BlockMetaData{
				BallotNumber: &apaxos.BallotNumber{Number: 1, NodeId: "1"},
			},
			blockB: &apaxos.BlockMetaData{
				BallotNumber: &apaxos.BallotNumber{Number: 2, NodeId: "1"},
			},
			expected: true,
		},
		{
			name: "Lower BallotNumber.NodeId in blockA",
			blockA: &apaxos.BlockMetaData{
				BallotNumber: &apaxos.BallotNumber{Number: 1, NodeId: "1"},
			},
			blockB: &apaxos.BlockMetaData{
				BallotNumber: &apaxos.BallotNumber{Number: 1, NodeId: "2"},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareBlocks(tt.blockA, tt.blockB)
			if result != tt.expected {
				t.Errorf("CompareBlocks() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
