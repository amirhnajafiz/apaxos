package hashing

import (
	"fmt"

	"github.com/f24-cse535/apaxos/pkg/models"
)

func HashTransaction(instance *models.Transaction) string {
	return fmt.Sprintf("%s-%s-%d", instance.Sender, instance.Reciever, instance.SequenceNumber)
}
