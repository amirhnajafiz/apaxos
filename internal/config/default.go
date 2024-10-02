package config

import (
	"github.com/f24-cse535/apaxos/internal/config/grpc"
	"github.com/f24-cse535/apaxos/internal/config/storage"
)

// Default return default configuration.
func Default() Config {
	return Config{
		NodeID: "unique",
		Client: "unique",
		Nodes:  make([]string, 0),
		GRPC: grpc.GRPC{
			Host:           "127.0.0.1",
			Port:           8080,
			WaitingTimeout: 10, // in milliseconds
		},
		MongoDB: storage.MongoDB{
			URI: "your atlas connection string",
		},
		Redis: storage.Redis{
			Host:     "",
			Port:     0,
			Password: "",
			Database: 0,
		},
	}
}
