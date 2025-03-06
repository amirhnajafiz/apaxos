package utils

import "github.com/f24-cse535/apaxos/pkg/rpc/apaxos"

// CompareBallotNumbers compares two ballot numbers a and b.
// Returns 1 if a > b, -1 if b > a, and 0 if a == b.
func CompareBallotNumbers(a, b *apaxos.BallotNumber) int {
	if a.GetNumber() != b.GetNumber() {
		if a.GetNumber() > b.GetNumber() {
			return 1
		}

		return -1
	}

	if a.GetNodeId() > b.GetNodeId() {
		return 1
	} else if a.GetNodeId() < b.GetNodeId() {
		return -1
	}

	return 0
}
