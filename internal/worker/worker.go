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
	w.Logger.Info("worker process started.", zap.Int("interval", w.Interval))

	// block worker if not enabled
	if !enable {
		return
	}

	for {
		// wait in a period
		time.Sleep(time.Duration(w.Interval) * time.Second)

		// then run the worker tasks
		if err := w.snapshotNodeState(); err != nil {
			w.Logger.Error("backup worker failed", zap.Error(err))
		}
	}
}
