package apaxos

import "time"

func (a Apaxos) waitForPromise() {
	receivedEvents := 0
	timeoutDuration := time.Duration(a.Timeout) * time.Millisecond
	majorityTimeoutDuration := time.Duration(a.MajorityTimeout) * time.Microsecond

	timer := time.NewTimer(timeoutDuration)
	defer timer.Stop()

	for {
		select {
		case _ = <-a.InChannel:
			receivedEvents++

			if receivedEvents == a.Majority {
				if !timer.Stop() {
					<-timer.C // drain the channel if needed
				}

				timer.Reset(majorityTimeoutDuration)
			}
		case <-timer.C:
			if receivedEvents < a.Majority {
				// return error
			}
		}
	}
}

func (a Apaxos) waitForAccepted() {
	receivedEvents := 0
	timeoutDuration := time.Duration(a.Timeout) * time.Millisecond
	majorityTimeoutDuration := time.Duration(a.MajorityTimeout) * time.Microsecond

	timer := time.NewTimer(timeoutDuration)
	defer timer.Stop()

	for {
		select {
		case _ = <-a.InChannel:
			receivedEvents++

			if receivedEvents == a.Majority {
				if !timer.Stop() {
					<-timer.C // drain the channel if needed
				}

				timer.Reset(majorityTimeoutDuration)
			}
		case <-timer.C:
			if receivedEvents < a.Majority {
				// return error
			}
		}
	}
}
