package cmd

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/f24-cse535/apaxos/internal/config"
	"github.com/f24-cse535/apaxos/internal/grpc/client"
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"

	"go.uber.org/zap"
)

var (
	errInvalidCommand = errors.New("command not found")
	errNumberOfArgs   = errors.New("args input are not enough")
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
		3: c.printBalance,      // printbalance <client> <node>
		4: c.printLogs,         // printlogs <node>
		5: c.printDB,           // printdb <node>
		6: c.performance,       // performance
		7: c.aggrigatedBalance, // aggrigated balance <client>
		8: c.newTransaction,    // new transaction <sender> <receiver> <amount> <node>
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
		return fmt.Errorf("input command should be an int in range of 1 to %d", len(c.commands))
	}

	// set arguments
	c.args = append(c.args, parts[1:]...)

	// run the command
	if val, ok := c.commands[command]; ok {
		return val()
	} else {
		return errInvalidCommand
	}
}

func (c *Controller) testcase() error {
	if len(c.args) < 1 {
		return errNumberOfArgs
	}

	// set path
	path := c.args[0]

	// open the CSV file
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// create a CSV reader
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // allow variable number of fields per row

	var currentIndex string
	var currentList string
	var currentTuples []string

	// read through the CSV records
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		// handle the index if present
		if record[0] != "" {
			if currentIndex != "" {
				// start the test-set
				c.execSet(currentIndex, strings.Split(currentList, ","), currentTuples)
			}

			// start a new block with the new index
			currentIndex = record[0]
			currentTuples = []string{} // reset tuples
			currentList = ""           // reset list

			// if a list is present in the first row, capture it
			if len(record) > 2 {
				currentList = strings.Trim(record[2], "[]") // remove square brackets
			}

			// if a tuple is present in the first row, add it to the tuples
			if len(record) > 1 {
				currentTuples = append(currentTuples, strings.Trim(record[1], "()"))
			}
		} else {
			// this is a continuation of the current block
			if len(record) > 1 && record[1] != "" {
				currentTuples = append(currentTuples, strings.Trim(record[1], "()"))
			}
		}
	}

	return nil
}

func (c *Controller) execSet(index string, servers []string, transactions []string) {
	fmt.Printf("starting set: %s\n", index)

	for _, server := range servers {
		fmt.Printf("%s ", server)
	}

	for _, transaction := range transactions {
		fmt.Printf("%s\n", transaction)
	}
}

func (c *Controller) block(address string) {
	if err := c.LDialer.ChangeState(address, false); err != nil {
		fmt.Printf("%s returned error: %v\n", address, err)
	}
}

func (c *Controller) reset() error {
	for key, value := range c.Cfg.GetNodes() {
		if err := c.LDialer.ChangeState(value, true); err != nil {
			fmt.Printf("%s returned error: %v\n", key, err)
		}
	}

	return nil
}

func (c *Controller) printBalance() error {
	if len(c.args) < 2 {
		return errNumberOfArgs
	}

	// get the client name and the node address
	client := c.args[0]
	address := c.Cfg.GetNodes()[c.args[1]]

	// call RPC call for printBalance
	balance, err := c.TDialer.PrintBalance(address, client)
	if err != nil {
		return err
	}

	fmt.Println(balance)

	return nil
}

func (c *Controller) printLogs() error {
	if len(c.args) < 1 {
		return errNumberOfArgs
	}

	// get the node address
	address := c.Cfg.GetNodes()[c.args[0]]

	// call RPC call for printLogs
	logs, err := c.TDialer.PrintLogs(address)
	if err != nil {
		return err
	}

	for _, log := range logs {
		fmt.Printf(
			"ballot-number: <%d - %s>\n",
			log.GetMetadata().GetBallotNumber().GetNumber(),
			log.GetMetadata().GetBallotNumber().GetNodeId(),
		)

		fmt.Println("transactions:")

		for _, transaction := range log.Transactions {
			fmt.Printf(
				"%d. (%s, %s, %d)",
				transaction.GetSequenceNumber(),
				transaction.GetSender(),
				transaction.GetReciever(),
				transaction.GetAmount(),
			)
		}
	}

	return nil
}

func (c *Controller) printDB() error {
	if len(c.args) < 1 {
		return errNumberOfArgs
	}

	// get the node address
	address := c.Cfg.GetNodes()[c.args[0]]

	// call RPC call for printDB
	blocks, err := c.TDialer.PrintDB(address)
	if err != nil {
		return err
	}

	for _, log := range blocks {
		fmt.Printf(
			"ballot-number: <%d - %s>\n",
			log.GetMetadata().GetBallotNumber().GetNumber(),
			log.GetMetadata().GetBallotNumber().GetNodeId(),
		)

		fmt.Println("transactions:")

		for _, transaction := range log.Transactions {
			fmt.Printf(
				"%d. (%s, %s, %d)\n",
				transaction.GetSequenceNumber(),
				transaction.GetSender(),
				transaction.GetReciever(),
				transaction.GetAmount(),
			)
		}
	}

	return nil
}

func (c *Controller) performance() error {
	for key, value := range c.Cfg.GetNodes() {
		if resp, err := c.TDialer.Performance(value); err == nil {
			fmt.Printf("%s: %f TPS, %f ms\n", key, resp.GetThroughput(), resp.GetLatency())
		} else {
			fmt.Printf("%s: no response: %v\n", key, err)
		}
	}

	return nil
}

func (c *Controller) aggrigatedBalance() error {
	if len(c.args) < 1 {
		return errNumberOfArgs
	}

	// get the client name
	client := c.args[0]

	for key, value := range c.Cfg.GetNodes() {
		if balance, err := c.TDialer.PrintBalance(value, client); err == nil {
			fmt.Printf("%s: %d\n", key, balance)
		} else {
			fmt.Printf("%s: no response: %v\n", key, err)
		}
	}

	return nil
}

func (c *Controller) newTransaction() error {
	if len(c.args) < 4 {
		return errNumberOfArgs
	}

	// get the node address
	address := c.args[3]

	// parse the amount
	amount, _ := strconv.Atoi(c.args[2])

	// create a new transaction
	t := apaxos.Transaction{
		Sender:   c.args[0],
		Reciever: c.args[1],
		Amount:   int64(amount),
	}

	// call rpc on the node
	if code, text, err := c.TDialer.NewTransaction(address, &t); err == nil {
		fmt.Printf("got %d from server: %s\n", code, text)
	} else {
		fmt.Printf("%s: returned error: %v\n", address, err)
	}

	return nil
}
