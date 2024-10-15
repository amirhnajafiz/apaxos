package consensus

import (
	"errors"

	protocol "github.com/f24-cse535/apaxos/internal/consensus/apaxos"
	"github.com/f24-cse535/apaxos/pkg/messages"
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"
	"go.uber.org/zap"
)

// newInstance builds and starts a new apaxos instance.
func (c Consensus) newInstance(transaction *apaxos.Transaction) {
	// first we create a new instance for the protocol
	c.instance = &protocol.Apaxos{
		NodeId:          c.NodeId,
		Dialer:          c.Dialer,
		Nodes:           c.Nodes,
		Memory:          c.Memory,
		Database:        c.Database,
		Majority:        c.Majority,
		MajorityTimeout: c.MajorityTimeout,
		Timeout:         c.RequestTimeout,
		InChannel:       make(chan *messages.Packet),
		OutChannel:      make(chan *messages.Packet),
	}

	// start a new go-routine for apaxos instance
	go func() {
		// after this go-routine finished, clear the protocol instance
		defer func() {
			// terminate consensus by resetting the instance
			c.instance = nil
			c.Logger.Debug("apaxos terminated")
		}()

		c.Logger.Debug("apaxos started")

		// in a while loop, try to make consensus
		for {
			// start apaxos protocol
			err := c.instance.Start()
			if err != nil {
				c.Logger.Error("apaxos failed", zap.Error(err))

				// check the error for a proper handling mechanism
				switch {
				case errors.Is(err, protocol.ErrRequestTimeout):
					// in this case, we first check to see if we have enough servers up and running or not
					if !c.livenessHandler() {
						c.failedConsensus(protocol.ErrNotEnoughServers)
						return
					}
				case errors.Is(err, protocol.ErrSlowNode):
					// if we hit slow-node, it means that we got synced, so we should check the status of balance before
					// rerunning the consensus protocol.
					c.Logger.Info("slow server detected")

					// if the client balance is now enough, then we submit the transaction
					if c.recheckBalance(transaction) {
						c.submitTransaction(transaction)
						break
					}
				default:
					c.Logger.Error("consensus error", zap.Error(err))
				}
			} else { // a successful consensus
				// now we check to see if the client balance is enough or not
				if c.recheckBalance(transaction) {
					c.successfulConsensus(transaction)
				} else {
					c.failedConsensus(protocol.ErrNotEnoughBalance)
				}

				return
			}
		}
	}()
}
