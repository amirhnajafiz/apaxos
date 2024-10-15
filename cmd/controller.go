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

// list of errors that could happen during the controller run
var (
	errInvalidCommand = errors.New("command not found")
	errNumberOfArgs   = errors.New("args input are not enough")
	errEndOfSets      = errors.New("no test-set available")
)

// test-set is a holder for each testcase in CSV.
type testSet struct {
	index        string
	serverList   []string
	transactions []string
}

// Controller is used to communicate with our distributed system using gRPC calls.
type Controller struct {
	// the config and logger modules
	Cfg    config.Config
	Logger *zap.Logger

	// gRPC clients to send messages
	TDialer client.TransactionsDialer
	LDialer client.LivenessDialer

	// commands is a map of all commands that are bound to a function
	commands map[int]func() error
	args     []string // args is a list that the command functions use to read their inputs from

	// testcases read from a csv file
	tests []testSet
	index int
}

func (c *Controller) Main() error {
	// init dialers
	c.TDialer = client.TransactionsDialer{}
	c.LDialer = client.LivenessDialer{}

	// init commands
	c.commands = map[int]func() error{
		0: c.help,              // print
		1: c.testcase,          // testcase <csv path>
		2: c.resetServers,      // reset
		3: c.printBalance,      // printbalance <client> <node>
		4: c.printLogs,         // printlogs <node>
		5: c.printDB,           // printdb <node>
		6: c.performance,       // performance
		7: c.aggrigatedBalance, // aggrigated balance <client>
		8: c.newTransaction,    // new transaction <sender> <receiver> <amount> <node>
		9: c.next,              // next test-set
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

// parseInput get's the user input to run system commands.
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

// help command displays controller instructions.
func (c *Controller) help() error {
	fmt.Println(
		`
exit: close the controller
0: print help
1: testcase <csv path>
2: reset
3: printbalance <client> <node>
4: printlogs <node>
5: printdb <node>
6: performance
7: aggrigated balance <client>
8: new_transaction <sender> <receiver> <amount> <node>
9: next (runs the next test-set)
		`,
	)

	return nil
}

// testcase accepts a testcase file to execute test sets.
func (c *Controller) testcase() error {
	if len(c.args) < 1 {
		return errNumberOfArgs
	}

	// create a new test-set
	c.index = 0
	c.tests = make([]testSet, 0)

	// set CSV path
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
	var currentTransactions []string

	// read through the CSV records
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		// handle the index if present
		if record[0] != "" {
			if currentIndex != "" {
				// store a test-set to run in future
				c.tests = append(c.tests, testSet{
					index:        currentIndex,
					serverList:   strings.Split(currentList, ","),
					transactions: currentTransactions,
				})
			}

			// start a new block with the new index
			currentIndex = record[0]
			currentTransactions = []string{} // reset transactions
			currentList = ""                 // reset list

			// if a list is present in the first row, capture it
			if len(record) > 2 {
				currentList = strings.Trim(record[2], "[]") // remove square brackets
			}

			// if a tuple is present in the first row, add it to the tuples
			if len(record) > 1 {
				currentTransactions = append(currentTransactions, strings.Trim(record[1], "()"))
			}
		} else {
			// this is a continuation of the current block
			if len(record) > 1 && record[1] != "" {
				currentTransactions = append(currentTransactions, strings.Trim(record[1], "()"))
			}
		}
	}

	return nil
}

// next runs the next test-set until it reaches the end of sets.
func (c *Controller) next() error {
	if c.index >= len(c.tests) {
		return errEndOfSets
	}

	// select the current index and execSet
	set := c.tests[c.index]
	c.execSet(set.index, set.serverList, set.transactions)

	// go for the next test-set
	c.index++

	return nil
}

// execSet runs a testcase set.
func (c *Controller) execSet(index string, servers []string, transactions []string) {
	fmt.Printf("starting set: %s\n", index)

	// create a map of living servers
	hashMap := make(map[string]bool)
	for _, server := range servers {
		hashMap[server] = true
	}

	// block servers that are not in hashMap
	for key, value := range c.Cfg.GetNodes() {
		if !hashMap[key] {
			c.blockServer(value)
		}
	}

	// run transactions
	for _, transaction := range transactions {
		parts := strings.Split(transaction, ",")

		// set args base args
		c.args = make([]string, 0)
		c.args = append(c.args, parts...)

		// set the address based on the client shard
		address := c.Cfg.GetClientShards()[c.args[0]]
		c.args = append(c.args, address)

		// submit a new transaction
		if err := c.newTransaction(); err != nil {
			fmt.Println(err)
		}
	}

	// reset all servers to unblock nodes
	if err := c.resetServers(); err != nil {
		fmt.Println(err)
	}
}

// block server takes an address and block the server.
func (c *Controller) blockServer(address string) {
	if err := c.LDialer.ChangeState(address, false); err != nil {
		fmt.Printf("%s returned error: %v\n", address, err)
	}
}

// reset servers unblocks all of the servers.
func (c *Controller) resetServers() error {
	for key, value := range c.Cfg.GetNodes() {
		if err := c.LDialer.ChangeState(value, true); err != nil {
			fmt.Printf("%s returned error: %v\n", key, err)
		}
	}

	return nil
}

// print balance runs a rpc call to get user balance.
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

// print logs prints the logs of a server.
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

// print db gets database of a node.
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

// performance loops over servers to get metrics.
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

// aggrigated balance gets the client name and runs the aggrigated balance method.
func (c *Controller) aggrigatedBalance() error {
	if len(c.args) < 1 {
		return errNumberOfArgs
	}

	// get the client name
	client := c.args[0]

	// runs print balance over servers
	for key, value := range c.Cfg.GetNodes() {
		if balance, err := c.TDialer.PrintBalance(value, client); err == nil {
			fmt.Printf("%s: %d\n", key, balance)
		} else {
			fmt.Printf("%s: no response: %v\n", key, err)
		}
	}

	return nil
}

// new transaction runs a new transaction over the system.
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
