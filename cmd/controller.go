package cmd

import (
	"github.com/f24-cse535/apaxos/internal/config"

	"go.uber.org/zap"
)

type Controller struct {
	Cfg    config.Config
	Logger *zap.Logger
}

func (c Controller) Main() error {
	// todo: the following commands
	// 1. testcase <csv path>
	// 2. reset
	// 3. printbalance <client>
	// 4. printlogs <node>
	// 5. printdb <node>
	// 6. performance
	// 7. aggrigated balance <client>
	// 8. new transaction <sender> <receiver> <amount>

	return nil
}
