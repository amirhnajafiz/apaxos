package compare_test

import (
	"testing"

	"github.com/f24-cse535/apaxos/internal/utils/compare"
	"github.com/f24-cse535/apaxos/pkg/models"
)

func TestCompareBallotNumbers(t *testing.T) {
	tests := []struct {
		name     string
		a        *models.BallotNumber
		b        *models.BallotNumber
		expected int
	}{
		{
			name:     "a.Number > b.Number",
			a:        &models.BallotNumber{Number: 3, NodeId: "1"},
			b:        &models.BallotNumber{Number: 2, NodeId: "1"},
			expected: 1,
		},
		{
			name:     "a.Number < b.Number",
			a:        &models.BallotNumber{Number: 2, NodeId: "1"},
			b:        &models.BallotNumber{Number: 3, NodeId: "1"},
			expected: -1,
		},
		{
			name:     "a.Number == b.Number, a.NodeId > b.NodeId",
			a:        &models.BallotNumber{Number: 2, NodeId: "2"},
			b:        &models.BallotNumber{Number: 2, NodeId: "1"},
			expected: 1,
		},
		{
			name:     "a.Number == b.Number, a.NodeId < b.NodeId",
			a:        &models.BallotNumber{Number: 2, NodeId: "1"},
			b:        &models.BallotNumber{Number: 2, NodeId: "2"},
			expected: -1,
		},
		{
			name:     "a.Number == b.Number, a.NodeId == b.NodeId",
			a:        &models.BallotNumber{Number: 2, NodeId: "1"},
			b:        &models.BallotNumber{Number: 2, NodeId: "1"},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compare.CompareBallotNumbers(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}
