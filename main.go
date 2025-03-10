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
		panic("you did not provide enough arguments to run! (./main <command> <configpath>)")
	}

	// load configs into a config struct
	cfg := config.New(argv[2])

	// create a new zap logger instance
	logr := logger.NewLogger(cfg.LogLevel)

	// create cmd instances and pass needed parameters
	commands := map[string]cmd.CMD{
		ControllerCmdName: cmd.Controller{
			Cfg:    cfg,
			Logger: logr.Named("controller"),
		}, // controller command
		NodeCmdName: cmd.Node{
			Cfg:    cfg,
			Logger: logr.Named("node"),
		}, // node command
		MongoCmdName: cmd.MongoDB{
			Cfg: cfg,
		}, // mongodb command
	}

	// command is the first argument variable
	command := argv[1]

	// then we check the command to run different programs based on the user input.
	if callback, ok := commands[command]; ok {
		if err := callback.Main(); err != nil {
			logr.Panic("failed to run the command", zap.Error(err), zap.String("command", command))
		}

		logr.Info("successful run", zap.String("command", command))
	} else {
		panic(
			fmt.Sprintf(
				"your input command must be the first argument variable, and it should be `%s`, `%s`, or `%s`.",
				ControllerCmdName,
				NodeCmdName,
				MongoCmdName,
			),
		)
	}
}
