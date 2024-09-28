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

// Keep a map of the nodes and their sockets
// Use a channel to communicate with http handler
// Accept new transactions (a block of transactions)
// Implement APAXOS for transactions.

// Open a socket listener to handle input requests from others
//     Their handlers are inside RPC
// Create a new process to communicate with HTTP handler

func Bootstrap(cfg socket.Config, client string, nodes ...string) error {
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
