package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/f24-cse535/apaxos/cmd"
	"github.com/f24-cse535/apaxos/internal/config"
	"github.com/f24-cse535/apaxos/internal/monitoring/logger"

	"go.uber.org/zap"
)

// here is the list of current system commands
const (
	ControllerCmdName = "controller"
	NodeCmdName       = "node"
)

func main() {
	// get argument variables
	argv := os.Args
	if len(argv) < 2 {
		panic("not enough arguments to run apaxos!")
	}

	// config file path is a flag
	configPath := flag.String("config", "config.yaml", "this is the config file path.")

	// parse flags
	flag.Parse()

	// load configs module
	cfg := config.New(*configPath)

	// create a new zap logger instance
	logr := logger.NewLogger(cfg.LogLevel)

	// create cmd instances and pass needed parameters
	ctl := cmd.Controller{
		Cfg:    cfg,
		Logger: logr.Named("controller"),
	}
	node := cmd.Node{
		Cfg:    cfg,
		Logger: logr.Named("node"),
	}

	// command is the first argument variable
	command := argv[1]

	// then we check the command to run different programs based on the
	// user input.
	switch command {
	case ControllerCmdName:
		if err := ctl.Main(); err != nil {
			logr.Panic("failed run controller", zap.Error(err))
		}
	case NodeCmdName:
		if err := node.Main(); err != nil {
			logr.Panic("failed to run node", zap.Error(err))
		}
	default:
		panic(
			fmt.Sprintf(
				"your input command must be the first argument variable, and it should be `%s` or `%s`.",
				ControllerCmdName,
				NodeCmdName,
			),
		)
	}
}
