package compare

import "github.com/f24-cse535/apaxos/pkg/models"

// CompareBlocks takes two blocks metadata and sorts them
// by comparing ballot-number and sequence number.
func CompareBlocks(a *models.BlockMetadata, b *models.BlockMetadata) bool {
	if a.BallotNumber.Number > b.BallotNumber.Number {
		return false
	} else if a.BallotNumber.NodeId > b.BallotNumber.NodeId {
		return false
	} else if a.SequenceNumber > b.SequenceNumber {
		return false
	} else {
		return true
	}
}
