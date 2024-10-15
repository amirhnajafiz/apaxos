package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/f24-cse535/apaxos/internal/config"
	"github.com/f24-cse535/apaxos/internal/grpc/client"

	"go.uber.org/zap"
)

// Controller is used to communicate with our distributed system
// using gRPC calls.
type Controller struct {
	Cfg    config.Config
	Logger *zap.Logger

	TDialer client.TransactionsDialer
	LDialer client.LivenessDialer

	commands map[int]func() error
	args     []string
}

func (c *Controller) Main() error {
	// init dialers
	c.TDialer = client.TransactionsDialer{}
	c.LDialer = client.LivenessDialer{}

	// init commands
	c.commands = map[int]func() error{
		1: c.testcase,          // testcase <csv path>
		2: c.reset,             // reset
		3: c.printBalance,      // printbalance <client>
		4: c.printLogs,         // printlogs <node>
		5: c.printDB,           // printdb <node>
		6: c.performance,       // performance
		7: c.aggrigatedBalance, // aggrigated balance <client>
		8: c.newTransaction,    // new transaction <sender> <receiver> <amount>
	}

	// in a for loop, read user commands
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")

		input, _ := reader.ReadString('\n') // read input until newline
		input = strings.TrimSpace(input)

		// no input
		if len(input) == 0 {
			continue
		}

		// exit command
		if input == "exit" {
			break
		}

		// call parse input handler
		if err := c.parseInput(input); err != nil {
			fmt.Println(err.Error())
		}
	}

	return nil
}

func (c *Controller) parseInput(input string) error {
	// reset input args
	c.args = make([]string, 0)

	// split into parts
	parts := strings.Split(input, " ")

	// take out the command
	command, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("input command should be in range of 1 to %d", len(c.commands))
	}

	// set arguments
	c.args = append(c.args, parts[1:]...)

	// run the command
	return c.commands[command]()
}

func (c *Controller) testcase() error {
	return nil
}

func (c *Controller) reset() error {
	return nil
}

func (c *Controller) printBalance() error {
	return nil
}

func (c *Controller) printLogs() error {
	return nil
}

func (c *Controller) printDB() error {
	return nil
}

func (c *Controller) performance() error {
	return nil
}

func (c *Controller) aggrigatedBalance() error {
	return nil
}

func (c *Controller) newTransaction() error {
	return nil
}
