package enum

// PacketType is a new type used in internal messaging.
type PacketType int

const (
	PacketPropose PacketType = iota + 1
	PacketPromise
	PacketAccept
	PacketAccepted
	PacketCommit
	PacketSync
)
