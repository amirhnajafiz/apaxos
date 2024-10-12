package messages

import "github.com/f24-cse535/apaxos/pkg/enum"

type Packet struct {
	Type    enum.PacketType `json:"type"`
	Payload interface{}     `json:"payload"`
}
