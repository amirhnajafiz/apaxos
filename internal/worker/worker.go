package worker

import (
	"time"

	"github.com/f24-cse535/apaxos/internal/storage/database"
	"github.com/f24-cse535/apaxos/internal/storage/local"

	"go.uber.org/zap"
)

// Worker is a component that runs node jobs like backup.
type Worker struct {
	Memory   *local.Memory
	Database *database.Database

	Logger *zap.Logger

	Interval int
}

// Start triggers a process for running the worker jobs.
func (w Worker) Start(enable bool) {
	// block worker if not enabled
	if !enable {
		return
	}

	w.Logger.Info("worker process started.", zap.Int("interval", w.Interval))

	for {
		// wait in a period
		time.Sleep(time.Duration(w.Interval) * time.Second)

		// then run the worker tasks
		if err := w.snapshotNodeState(); err != nil {
			w.Logger.Info("backup worker failed", zap.Error(err))
		}
	}
}
