package compare

import "github.com/f24-cse535/apaxos/pkg/models"

// CompareBallotNumbers get's two a and b numbers to compare them.
// Returns 1 if a > b.
// Returns -1 if b > a.
// Returns 0 is a = b.
func CompareBallotNumbers(a *models.BallotNumber, b *models.BallotNumber) int {
	if a.Number > b.Number {
		return 1
	} else if a.Number < b.Number {
		return -1
	} else {
		if a.NodeId > b.NodeId {
			return 1
		} else if a.NodeId < b.NodeId {
			return -1
		} else {
			return 0
		}
	}
}
