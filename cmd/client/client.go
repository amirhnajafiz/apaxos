package client

import (
	"fmt"

	"github.com/f24-cse535/apaxos/internal/grpc/client"
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"
)

// Client is a simple struct that makes RPC calls to servers.
type Client struct {
	Dialer *client.Client
}

// UpdateServerStatus is used to change the status of a server.
func (c Client) UpdateServerStatus(address string, status bool) {
	if err := c.Dialer.ChangeState(address, status); err != nil {
		fmt.Printf("%s returned error: %v\n", address, err)
	}
}

// PrintBalance runs a rpc call to get user balance.
func (c Client) PrintBalance(client, address string) error {
	// call RPC call for printBalance
	balance, err := c.Dialer.PrintBalance(address, client)
	if err != nil {
		return err
	}

	fmt.Println(balance)

	return nil
}

// PrintLogs prints the logs of a server.
func (c Client) PrintLogs(address string) error {
	// call RPC call for printLogs
	logs, err := c.Dialer.PrintLogs(address)
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

// PrintDB gets database of a node.
func (c Client) PrintDB(address string) error {
	// call RPC call for printDB
	blocks, err := c.Dialer.PrintDB(address)
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

// Performance loops over servers to get metrics.
func (c Client) Performance(addresses map[string]string) error {
	for key, address := range addresses {
		if resp, err := c.Dialer.Performance(address); err == nil {
			fmt.Printf("%s: %f TPS, %f microseconds\n", key, resp.GetThroughput(), resp.GetLatency())
		} else {
			fmt.Printf("%s: no response: %v\n", key, err)
		}
	}

	return nil
}

// AggrigatedBalance gets the client name and runs the aggrigated balance method.
func (c Client) AggrigatedBalance(client string, addresses map[string]string) error {
	// runs print balance over servers
	for key, value := range addresses {
		if balance, err := c.Dialer.PrintBalance(value, client); err == nil {
			fmt.Printf("%s: %d\n", key, balance)
		} else {
			fmt.Printf("%s: no response: %v\n", key, err)
		}
	}

	return nil
}

// Transaction runs a new transaction over the system.
func (c Client) Transaction(sender string, receiver string, amount int, address string) error {
	// create a new transaction
	t := &apaxos.Transaction{
		Sender:   sender,
		Reciever: receiver,
		Amount:   int64(amount),
	}

	fmt.Printf("sending (%s, %s, %d) to %s\n", t.GetSender(), t.GetReciever(), t.GetAmount(), address)

	// call rpc on the node
	if text, err := c.Dialer.NewTransaction(address, t); err == nil {
		fmt.Printf("server: %s\n", text)
	} else {
		fmt.Printf("%s: returned error: %v\n", address, err)
	}

	return nil
}
