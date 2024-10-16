package utils

import "github.com/f24-cse535/apaxos/pkg/rpc/apaxos"

// CompareBlocks return true if a < b, otherwise returns false
func CompareBlocks(a, b *apaxos.BlockMetaData) bool {
	return CompareBallotNumbers(a.GetBallotNumber(), b.GetBallotNumber()) <= 0
}
