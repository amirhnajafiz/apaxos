package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/f24-cse535/apaxos/cmd"
	"github.com/f24-cse535/apaxos/internal/config"
)

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

	// load configs
	cfg := config.New(*configPath)

	// create cmd instances and pass the config file path
	ctl := cmd.Controller{
		ConfigPath: *configPath,
	}
	node := cmd.Node{
		Cfg: cfg,
	}

	// command is the first argument variable
	command := argv[1]

	// then we check the command to run different programs based on the
	// user input.
	switch command {
	case ControllerCmdName:
		ctl.Main()
	case NodeCmdName:
		node.Main()
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
