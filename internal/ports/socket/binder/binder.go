package binder

import (
	"fmt"
	"log"
	"net"

	"github.com/f24-cse535/apaxos/internal/config/socket"
)

const (
	SERVER_TYPE = "tcp"
	SERVER_HOST = "127.0.0.1"
)

func Bootstrap(cfg socket.Config) error {
	// open a new interface
	server, err := net.Listen(SERVER_TYPE, fmt.Sprintf("%s:%d", SERVER_HOST, cfg.Port))
	if err != nil {
		return fmt.Errorf("failed to open node's socket interface on port %d: %v", cfg.Port, err)
	}

	// close after this method is done
	defer server.Close()

	// interface loop
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Println(fmt.Errorf("failed to accept client: %v", err))
			continue
		}

		// send conn to a state machine
		conn.Close()
	}
}
