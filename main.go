package main

import (
	"github.com/f24-cse535/apaxos/internal/config"
	"github.com/f24-cse535/apaxos/internal/ports/http"
	"github.com/f24-cse535/apaxos/internal/ports/socket/binder"
)

func main() {
	// load configs
	cfg := config.New("config.yaml")

	// TODO: open db connection
	// TODO: open redis connection

	// bootstrap http server as a goroutine
	go http.Bootstrap(cfg.HTTP)

	// bootstrap socket interface
	go binder.Bootstrap(cfg.Socket)
}
