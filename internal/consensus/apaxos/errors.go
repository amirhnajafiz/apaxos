package apaxos

import "errors"

var (
	ErrSlowNode         = errors.New("sync received, got a slow node case")
	ErrNotEnoughServers = errors.New("not enough servers to decide")
	ErrRequestTimeout   = errors.New("did not get enough responses")
	ErrNotEnoughBalance = errors.New("not enough balance to process")
	ErrCommitTimeout    = errors.New("waited to long for own commitment")
)
