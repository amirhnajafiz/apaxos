package config

import (
	"github.com/f24-cse535/apaxos/internal/config/grpc"
	"github.com/f24-cse535/apaxos/internal/config/mongodb"
)

// Default return default configuration.
func Default() Config {
	return Config{
		NodeID:          "unique",
		Client:          "unique",
		Majority:        0,
		WorkersEnabled:  false,
		WorkersInterval: 10, // in seconds
		LogLevel:        "debug",
		Nodes:           make([]Pair, 0),
		Clients:         make([]Pair, 0),
		ClientsShards:   make([]Pair, 0),
		GRPC: grpc.Config{
			Host:            "127.0.0.1",
			Port:            8080,
			RequestTimeout:  10, // in milliseconds
			MajorityTimeout: 10, // in microseconds
		},
		MongoDB: mongodb.Config{
			URI:      "your atlas connection string",
			Database: "",
		},
	}
}
