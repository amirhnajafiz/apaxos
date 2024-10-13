package enum

// PacketType is a new type used in internal messaging.
type PacketType int

const (
	PacketPropose PacketType = iota + 1
	PacketPrepare
	PacketPromise
	PacketAccept
	PacketAccepted
	PacketCommit
	PacketSync
	PacketTransaction
	PacketError
)
