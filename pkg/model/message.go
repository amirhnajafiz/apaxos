package model

type Message struct {
	BallotNumber int           `json:"ballot_number"`
	AcceptNum    int           `json:"accept_number"`
	AcceptValue  []Transaction `json:"accept_value"`
}
