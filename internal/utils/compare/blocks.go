package compare

import "github.com/f24-cse535/apaxos/pkg/models"

// CompareBlocks return true if a < b, otherwise returns false
func CompareBlocks(a, b *models.BlockMetadata) bool {
	return CompareBallotNumbers(&a.BallotNumber, &b.BallotNumber) <= 0
}
