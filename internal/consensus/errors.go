package consensus

import "errors"

var (
	ErrMultipleInstances = errors.New("cannot run multiple consensus protocols at the same time on this machine")
)
