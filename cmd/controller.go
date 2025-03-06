package cmd

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	goclient "github.com/f24-cse535/apaxos/cmd/client"
	"github.com/f24-cse535/apaxos/internal/config"
	"github.com/f24-cse535/apaxos/internal/grpc/client"

	"go.uber.org/zap"
)

// list of errors that could happen during the controller run
var (
	errInvalidCommand = errors.New("command not found")
	errEndOfSets      = errors.New("no test-set available")

	// currentTest is a holder for executing test sets in testCase array
	currentTest int
	// testCase is a reference to an array of test sets
	testCase *[]testSet
)

// testSet is a type to store test sets.
type testSet struct {
	index        string
	serverList   []string
	transactions []map[string]interface{}
}

// Controller is used to communicate with our distributed system using gRPC calls.
type Controller struct {
	// the config and logger modules
	Cfg    config.Config
	Logger *zap.Logger

	// client module to make gRPC calls
	client *goclient.Client
}

func (c Controller) Main() error {
	// init client to make rpc calls
	c.client = &goclient.Client{
		Dialer: client.NewClient(c.Logger.Named("client")),
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

		// call parse input handler
		if err := c.parseInput(input); err != nil {
			fmt.Println(err.Error())
		}
	}
}

// parseInput get's the user input to run system commands.
func (c Controller) parseInput(input string) error {
	// split into parts
	parts := strings.Split(input, " ")

	// create an error holder
	var err error

	// take out the command by parsing the first part and switch on the input command and run functions
	switch parts[0] {
	case "exit":
		os.Exit(0)
	case "help":
		c.printHelp()
	case "tests":
		err = c.readTests(parts[1])
	case "next":
		if currentTest >= len(*testCase) {
			err = errEndOfSets
		} else {
			c.execSet(&(*testCase)[currentTest])
			currentTest++
		}
	case "ping":
		c.pintServer(c.Cfg.GetNodes()[parts[1]])
	case "reset":
		c.resetServers()
	case "block":
		err = c.client.UpdateServerStatus(c.Cfg.GetNodes()[parts[1]], false)
	case "unblock":
		err = c.client.UpdateServerStatus(c.Cfg.GetNodes()[parts[1]], true)
	case "printbalance":
		err = c.client.PrintBalance(parts[1], c.Cfg.GetNodes()[c.Cfg.GetClientShards()[parts[1]]])
	case "printlogs":
		err = c.client.PrintLogs(c.Cfg.GetNodes()[parts[1]])
	case "printdb":
		err = c.client.PrintDB(c.Cfg.GetNodes()[parts[1]])
	case "performance":
		err = c.client.Performance(c.Cfg.GetNodes())
	case "aggrigated":
		err = c.client.AggrigatedBalance(parts[1], c.Cfg.GetNodes())
	case "transaction":
		sender := parts[1]
		receiver := parts[2]
		amount, _ := strconv.Atoi(parts[3])

		c.client.Transaction(sender, receiver, amount, c.Cfg.GetNodes()[c.Cfg.GetClientShards()[sender]])
	default:
		return errInvalidCommand
	}

	return err
}

// help command displays controller instructions.
func (c Controller) printHelp() error {
	fmt.Println(
		`
exit: close the controller app
help | prints help instructions

tests  <csv path> | loads a csv test file
next | runs the next test-set

reset | reset all servers status to active
block <node> | get a node out of access
unblock <node> | restore a single node 
ping <node> | send a ping message to a node

printbalance <client> | print the balance of a client (based on shards)
printlogs <node> | print logs of a node
printdb <node> | print database of a node
performance | gets performance of all nodes
aggrigated <client> | gets a clients balance in all servers
transaction <sender> <receiver> <amount> | make a transaction for a client`,
	)

	return nil
}

// readTests reads a CSV file into current testcase array
func (c Controller) readTests(path string) error {
	// set testcases into an array
	tmp := make([]testSet, 0)
	testCase = &tmp

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
				// first trim the servers names and make them all upper-case
				servers := strings.Split(currentList, ", ")
				for index, value := range servers {
					servers[index] = strings.ToUpper(strings.TrimSpace(value))
				}

				// now modify transactions
				transactions := make([]map[string]interface{}, 0)
				for _, transaction := range currentTransactions {
					// split them by `,`
					parts := strings.Split(transaction, ", ")
					for index, part := range parts {
						parts[index] = strings.ToUpper(strings.TrimSpace(part)) // make them all upper-case
					}

					sender := parts[0]
					receiver := parts[1]
					amount, _ := strconv.Atoi(parts[2])
					address := c.Cfg.GetClientShards()[sender]

					// save the transaction
					transactions = append(transactions, map[string]interface{}{
						"sender":   sender,
						"receiver": receiver,
						"amount":   amount,
						"address":  address,
					})
				}

				// store a test-set
				tmp = append(tmp, testSet{
					index:        currentIndex,
					serverList:   servers,
					transactions: transactions,
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

	// first trim the servers names and make them all upper-case
	servers := strings.Split(currentList, ", ")
	for index, value := range servers {
		servers[index] = strings.ToUpper(strings.TrimSpace(value))
	}

	// now modify transactions
	transactions := make([]map[string]interface{}, 0)
	for _, transaction := range currentTransactions {
		// split them by `,`
		parts := strings.Split(transaction, ", ")
		for index, part := range parts {
			parts[index] = strings.ToUpper(strings.TrimSpace(part)) // make them all upper-case
		}

		sender := parts[0]
		receiver := parts[1]
		amount, _ := strconv.Atoi(parts[2])
		address := c.Cfg.GetClientShards()[sender]

		// save the transaction
		transactions = append(transactions, map[string]interface{}{
			"sender":   sender,
			"receiver": receiver,
			"amount":   amount,
			"address":  address,
		})
	}

	// store a test-set
	tmp = append(tmp, testSet{
		index:        currentIndex,
		serverList:   servers,
		transactions: transactions,
	})

	fmt.Printf("file %s, loading %d sets.\n", path, len(tmp))

	return nil
}

// execSet runs a testcase set.
func (c Controller) execSet(set *testSet) {
	fmt.Printf("starting set: %s\n", set.index)

	// reset all servers to unblock nodes
	c.resetServers()

	// create a map of living servers
	hashMap := make(map[string]bool)
	for _, server := range set.serverList {
		fmt.Printf("active server: %s\n", server)
		hashMap[server] = true
	}

	// block servers that are not in hashMap
	for key, value := range c.Cfg.GetNodes() {
		if !hashMap[key] {
			fmt.Printf("blocking server: %s\n", key)
			c.client.UpdateServerStatus(value, false)
		}
	}

	// run transactions
	for _, ts := range set.transactions {
		// submit a new transaction
		sender := ts["sender"].(string)
		receiver := ts["receiver"].(string)
		amount := ts["amount"].(int)
		node := ts["address"].(string)

		address := c.Cfg.GetNodes()[node]

		if err := c.client.Transaction(sender, receiver, amount, address); err != nil {
			fmt.Println(err)
		}
	}

	// reset all servers to unblock nodes
	c.resetServers()
}

// reset servers unblocks all of the servers.
func (c Controller) resetServers() {
	for _, value := range c.Cfg.GetNodes() {
		if err := c.client.UpdateServerStatus(value, true); err != nil {
			fmt.Println(err)
		}
	}
}

// ping server sends a ping request to a node to check it's availability.
func (c Controller) pintServer(address string) {
	fmt.Println(c.client.Dialer.Ping(address))
}
