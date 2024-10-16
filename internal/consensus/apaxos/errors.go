package apaxos

import "errors"

var (
	ErrSlowNode         = errors.New("sync received, got a slow node case")
	ErrNotEnoughServers = errors.New("did not get majority responses")
	ErrRequestTimeout   = errors.New("did not get enough responses")
	ErrNotEnoughBalance = errors.New("not enough balance to process")
)
