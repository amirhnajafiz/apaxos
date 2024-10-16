# APAXOS

This project is a distributed transaction management service. The goal is to replicate transactions over various number of nodes, keep the data consistant, and tolerate node failures. To achive these goals, we are going to use `PAXOS` consensus protocl. However, we are going to modify this protocol to handle multiple values in a single instance.

## System's flow

Clients send their transactions to servers (there is user-level sharding that map each client's account to one single machine). Servers store them inside a local storage as a block of transactions. After that, they run each transaction.

They are going to run `apaxos` consensus protocol when a client sends a request that has an amount more than the current balance of that client. By running this protocol, a leader node collects all other node's transactions and forms a huge block of transactions (or a list of transaction blocks). Then it replicates them over other servers.

Finally, as the majority of servers get that major block, each server starts running these transactions and stores them in a persistante storage (in our case its `MongoDB`).

## System components

In this section, there is a list of the system components:

- `gRPC` server
  - gets requests from both other nodes and clients
  - contains a `apaxos` server, `transactions` server, and `liveness` server
  - it uses a consensus module to run `apaxos` instances
- `consensus`
  - running an apaxos instance:
    1. sends prepare requests with a _ballot number_
    2. waits for the majority/a timeout perioud
    3. collects all transactions to create a _major block_
    4. sends accept requests with a _major block_
    5. waits for the majority/a timeout perioud
    6. sends a commit message
  - handling input requests:
    1. gets propose requests and compare ballot number with its own ballot number
    2. returns the promise with accepted num and _accepted val_ or its own transactions
    3. gets the commit request, clears the block list and executes the transactions
- `database`
  - connects to a `MongoDB` cluster
- `memory`
  - uses local memory of the node to keep data
- `worker`
  - runs backup tasks to keep track of node's states

### System functions

- new transaction file
- new transaction
- print balance (X)
- print logs
- print db
- performance (latency, throughput)
- aggregated balance (X)

### Data types

- Block
  - List of transactions (array of transactions)
  - Ballot Number
- Transaction
  - Sender
  - Reciever
  - Amount
  - Sequence number
- Major Block (Block List)
  - Ballot Number
  - List of Blocks (array of blocks) ordered by their Ballot Number
- Ballot Number
  - Contains a `N` number and `ID` of server

## Phases

1. Propose, Promise, and Sync
2. Accept, and Accepted
3. Commit

### Propose

A proposer sends this request with a ballot number. Contains it's last committed Block info (UID).

### Promise

An acceptor checks the propose ballot number with it's accepted num. If ballot number is greater than accepted num,
it will return a promise response with it's accepted val or block list.

### Accept

After collecting all logs from acceptors (waits for the majority). The proposer, creates a Major Block, or selects an accepted val with highest ballot number. Then it sends an accept request with its own ballot number.

### Accepted

Each acceptor checks its accepted num and accepted val with the given accept request. If it is ok it will update it's
accepted val and accepted num.

### Commit

Finally, the proposer waits for the majority. If it get's enough accepted responses, it will send a commit message.
After getting the commit message, each node clears it's block list by comparing it to the accepted val. It will also
store the accepted val and removes it to accept future messages.

### Sync

During the propose process, or the accept process. If a server sends an old commit block in return, the propose
sends a sync request and the list of blocks that where stored after that block. So, the node will be synced.

## User/System diagram

![system diagram](.github/assets/diagram.svg)

## Requirements

- Programming language: `Golang 1.23`
- Communication: `gRPC v2`, `protoc3`
- Database: `MongoDB`
- Logging: `zap logger`
- Configs: `koanf`
