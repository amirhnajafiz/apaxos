package apaxos

import (
	"time"

	"github.com/f24-cse535/apaxos/pkg/enum"
	"github.com/f24-cse535/apaxos/pkg/rpc/apaxos"
)

// waitForPromise state goes into a for loop to get a majority of promise messages.
func (a *Apaxos) waitForPromise() error {
	// start with 0 events
	receivedEvents := 0

	// set a request timeout and a majority timeout
	timeoutDuration := time.Duration(a.Timeout) * time.Millisecond
	majorityTimeoutDuration := time.Duration(a.MajorityTimeout) * time.Microsecond

	// set the timer to request timeout first
	timer := time.NewTimer(timeoutDuration)
	defer timer.Stop()

	// go inside a loop
	for {
		select {
		case pkt := <-a.InChannel:
			// if received a sync packet, we should halt
			if pkt.Type != enum.PacketSync {
				return ErrSlowNode
			}

			// if the packet is not promised then don't count
			if pkt.Type != enum.PacketPromise {
				continue
			}

			// got one promise message
			receivedEvents++
			a.promisedMessage = append(a.promisedMessage, pkt.Payload.(*apaxos.PromiseMessage))

			// check to see if we got the majority
			if receivedEvents == a.Majority {
				if !timer.Stop() {
					<-timer.C // drain the channel if needed
				}
				// reset the timer for the rest of the nodes
				timer.Reset(majorityTimeoutDuration)
			}
		case <-timer.C:
			if receivedEvents < a.Majority {
				// return error of not getting the majority
				return ErrNotEnoughServers
			} else {
				return nil
			}
		}
	}
}

// waitForAccepted state goes into a for loop to get the majority of accepted messages.
func (a *Apaxos) waitForAccepted() error {
	// start with 0 events
	receivedEvents := 0

	// set a request timeout and a majority timeout
	timeoutDuration := time.Duration(a.Timeout) * time.Millisecond
	majorityTimeoutDuration := time.Duration(a.MajorityTimeout) * time.Microsecond

	// set the timer to request timeout first
	timer := time.NewTimer(timeoutDuration)
	defer timer.Stop()

	// go inside a loop
	for {
		select {
		case pkt := <-a.InChannel:
			// if the packet is not promised then don't count
			if pkt.Type != enum.PacketPromise {
				continue
			}

			// got one accepted message
			receivedEvents++

			// check to see if we got the majority
			if receivedEvents == a.Majority {
				if !timer.Stop() {
					<-timer.C // drain the channel if needed
				}
				// reset the timer for the rest of the nodes
				timer.Reset(majorityTimeoutDuration)
			}
		case <-timer.C:
			if receivedEvents < a.Majority {
				// return error of not getting the majority
				return ErrNotEnoughServers
			} else {
				return nil
			}
		}
	}
}
