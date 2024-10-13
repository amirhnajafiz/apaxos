package hashing

import (
	"fmt"

	"github.com/f24-cse535/apaxos/pkg/models"
)

func HashBlock(instance *models.Block) string {
	return fmt.Sprintf("%s-%d-%d", instance.Metadata.NodeId, instance.Metadata.SequenceNumber, len(instance.Transactions))
}
