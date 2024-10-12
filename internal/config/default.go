package config

import (
	"github.com/f24-cse535/apaxos/internal/config/grpc"
	"github.com/f24-cse535/apaxos/internal/config/mongodb"
)

// Default return default configuration.
func Default() Config {
	return Config{
		NodeID: "unique",
		Client: "unique",
		Nodes:  make([]string, 0),
		GRPC: grpc.Config{
			Host:           "127.0.0.1",
			Port:           8080,
			WaitingTimeout: 10, // in milliseconds
		},
		MongoDB: mongodb.Config{
			URI:      "your atlas connection string",
			Database: "",
		},
	}
}
