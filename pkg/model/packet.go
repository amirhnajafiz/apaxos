package model

type Packet struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}
