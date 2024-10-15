package utils

import "github.com/f24-cse535/apaxos/pkg/models"

// CompareBallotNumbers compares two ballot numbers a and b.
// Returns 1 if a > b, -1 if b > a, and 0 if a == b.
func CompareBallotNumbers(a, b *models.BallotNumber) int {
	if a.Number != b.Number {
		if a.Number > b.Number {
			return 1
		}

		return -1
	}

	if a.NodeId > b.NodeId {
		return 1
	} else if a.NodeId < b.NodeId {
		return -1
	}

	return 0
}
