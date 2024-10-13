package enum

// StateType is a new type used in internal apaxos' state-machines.
type StateType int

const (
	StateInit StateType = iota + 1
	StateInWaitForMajority
	StateInTimeoutWait
	StateFreeToAccept
	StateBlocked
)
