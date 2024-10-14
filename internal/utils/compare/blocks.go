package compare

import "github.com/f24-cse535/apaxos/pkg/models"

// CompareBlocks takes two blocks metadata and sorts them
// by comparing ballot-number and sequence number.
// Returns 1 if a > b.
// Returns -1 if b > a.
// Returns 0 if a = b.
func CompareBlocks(a, b *models.BlockMetadata) int {
	ballotNumbers := CompareBallotNumbers(&a.BallotNumber, &b.BallotNumber)
	if ballotNumbers != 0 {
		return ballotNumbers
	}

	if a.SequenceNumber > b.SequenceNumber {
		return 1
	} else if a.SequenceNumber < b.SequenceNumber {
		return -1
	}

	return 0
}

// SortBlocks return true if a < b, otherwise returns false
func SortBlocks(a, b *models.BlockMetadata) bool {
	return CompareBlocks(a, b) <= 0
}
