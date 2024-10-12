package enum

// Message is a new type used in internal messaging.
type Message int

const (
	MessagePropose Message = iota + 1
	MessagePromise
	MessageAccept
	MessageAccepted
	MessageCommit
	MessageSync
)
