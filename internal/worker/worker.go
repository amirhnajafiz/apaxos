package worker

import (
	"log"
	"time"

	"github.com/f24-cse535/apaxos/internal/storage/database"
	"github.com/f24-cse535/apaxos/internal/storage/local"
)

type Worker struct {
	Memory   *local.Memory
	Database *database.Database

	Interval int
}

func (w Worker) Start() {
	for {
		// wait in a period
		time.Sleep(time.Duration(w.Interval) * time.Second)

		// then run the worker tasks
		if err := w.snapshotNodeState(); err != nil {
			log.Printf("backup worker failed: %v\n", err)
		}
	}
}
