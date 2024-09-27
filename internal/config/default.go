package config

import (
	"github.com/f24-cse535/apaxos/internal/config/http"
	"github.com/f24-cse535/apaxos/internal/config/socket"
)

// Default return default configuration.
func Default() Config {
	return Config{
		Nodes:  make([]string, 0),
		Client: "",
		HTTP: http.Config{
			Port: 8080,
		},
		Socket: socket.Config{
			Port:    8081,
			Timeout: 10, // in seconds
		},
	}
}
