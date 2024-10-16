package main

import (
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
	MongoCmdName      = "mongodb"
)

func main() {
	// get argument variables
	argv := os.Args
	if len(argv) < 3 {
		panic("not enough arguments to run apaxos!")
	}

	// load configs into a config struct
	cfg := config.New(argv[2])

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
	db := cmd.MongoDB{
		Cfg: cfg,
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
	case MongoCmdName:
		db.Main()
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
