package consensus

import (
	"github.com/f24-cse535/apaxos/internal/storage/database"
	"github.com/f24-cse535/apaxos/internal/storage/local"
	"github.com/f24-cse535/apaxos/pkg/messages"
)

// Consensus module is the core module that runs consensus protocols
// by getting the gRPC level packets.
type Consensus struct {
	Memory   *local.Memory
	Database *database.Database

	Clients map[string]string
	Nodes   map[string]string

	RequestTimeout  int `koanf:"request_timeout"`
	MajorityTimeout int `koanf:"majority_timeout"`
}

// Signal is used by the upper layer (gRPC functions) to send their
// packets to the consensus module.
func (c Consensus) Signal(pkt *messages.Packet) error {
	return nil
}
