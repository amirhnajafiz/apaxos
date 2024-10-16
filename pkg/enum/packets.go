package enum

// PacketType is a new type used in internal messaging.
type PacketType int

const (
	PacketPromise PacketType = iota + 1
	PacketAccepted
	PacketSync
)
