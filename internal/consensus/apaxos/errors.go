package apaxos

import "errors"

var (
	errSlowNode       = errors.New("sync received, got a slow node case")
	errRequestTimeout = errors.New("did not get enough responses")
)
