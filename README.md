# APAXOS

We are going to run `apaxos` when a client sends a request that has an amount more than the current balance of the client.

## Components

- gRPC server
  - get's requests from both other nodes and clients
- proposer
  - send's propose request with a _ballot number_
  - waits for the majority + a timeout perioud
  - collects all logs to create a _major block_
- acceptor
  - get's propose requests and compare ballot number with _accepted num_
  - returns the promise with accepted num and _accepted val_
  - get's the commit request
  - clears the block list
- learner
  - does the transaction process to update the values
  - returns the client response

## Data types

- Block
  - List of transactions (array of transactions)
  - ID: server id of that block
  - UID: a unique Id to check the block
  - Seq number: provided by the server that has it
  - Ballot Number
- Transaction
  - Sender
  - Reciever
  - Amount
- Major Block (Block List)
  - List of Blocks (array of blocks) ordered by their Ballot Number
- Ballot Number
  - Contains a N number and ID of server

## Phases

1. Propose, Promise, and Sync
2. Accept, and Accepted
3. Commit

## Requests

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

## Requirements

- Programming language: `Golang 1.23`
- Communication: `gRPC v2`, `protoc3`
- Datastores: `MongoDB` and `Redis`

## Functions

- new transaction file
- new transaction
- print balance (X)
- print logs
- print db
- performance (latency, throughput)
- aggregated balance (X)

## Design

### Services

1. A centeralized process that:
   1. Sends function requests using gRPC
   2. Creates nodes
   3. Deletes nodes
2. A node process that:
   1. Listens on a gRPC server
   2. Get's input requests that should be send to:
      1. Proposer sub-process
      2. Acceptor sub-process
      3. Learner sub-process
   3. A dialar module to call other nodes using gRPC
3. A MongoDB server
   1. Servers have their own collections using their SID as prefix
      1. Collections are block list and history
4. A Redis server
   1. Servers have their own key-value pairs using their SID as prefix
      1. They store their accepted num, accepted val, and clients balances

### Schema

```
client   =>    centeralized process
                          ||
                          \/
                   Node S1 gRPC   =>   Dispatcher  => Proposer, Acceptor, Learner
                          ||
                          \/
                   Dialer gRPC
                          ||
                          \/
  Redis  <=  Other Nodes gRPC  => MongoDB
```
