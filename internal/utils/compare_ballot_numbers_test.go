package utils

import (
	"testing"

	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"
)

func TestCompareBallotNumbers(t *testing.T) {
	tests := []struct {
		name     string
		a        *apaxos.BallotNumber
		b        *apaxos.BallotNumber
		expected int
	}{
		{
			name:     "a.Number > b.Number",
			a:        &apaxos.BallotNumber{Number: 3, NodeId: "1"},
			b:        &apaxos.BallotNumber{Number: 2, NodeId: "1"},
			expected: 1,
		},
		{
			name:     "a.Number < b.Number",
			a:        &apaxos.BallotNumber{Number: 2, NodeId: "1"},
			b:        &apaxos.BallotNumber{Number: 3, NodeId: "1"},
			expected: -1,
		},
		{
			name:     "a.Number == b.Number, a.NodeId > b.NodeId",
			a:        &apaxos.BallotNumber{Number: 2, NodeId: "2"},
			b:        &apaxos.BallotNumber{Number: 2, NodeId: "1"},
			expected: 1,
		},
		{
			name:     "a.Number == b.Number, a.NodeId < b.NodeId",
			a:        &apaxos.BallotNumber{Number: 2, NodeId: "1"},
			b:        &apaxos.BallotNumber{Number: 2, NodeId: "2"},
			expected: -1,
		},
		{
			name:     "a.Number == b.Number, a.NodeId == b.NodeId",
			a:        &apaxos.BallotNumber{Number: 2, NodeId: "S1"},
			b:        &apaxos.BallotNumber{Number: 2, NodeId: "S1"},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareBallotNumbers(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}
