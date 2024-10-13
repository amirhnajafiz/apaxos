package consensus

import (
	"github.com/f24-cse535/apaxos/internal/storage/database"
	"github.com/f24-cse535/apaxos/internal/storage/local"
)

type Consensus struct {
	Memory   *local.Memory
	Database *database.Database
}
