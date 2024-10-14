package cmd

import (
	"github.com/f24-cse535/apaxos/internal/config"

	"go.uber.org/zap"
)

type Controller struct {
	Cfg     config.Config
	Logger  *zap.Logger
	CSVPath string
}

func (c Controller) Main() {}
