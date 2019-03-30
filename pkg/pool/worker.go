package pool

import (
	"database/sql"
	"log"
	"time"

	"timescale/pkg/job"
)

// Worker represents a job execution unit
type Worker struct {
	ID   int
	in   <-chan job.Job
	out  chan<- int64
	done chan<- bool
}

// Exec performs sql query
func (w Worker) Exec(db *sql.Stmt) {
	for j := range w.in {
		start := time.Now()
		_, err := db.Exec(j.Hostname, j.StartTime, j.EndTime)
		elapsed := time.Since(start)
		if err != nil {
			log.Println(err)
			continue
		}
		w.out <- elapsed.Nanoseconds()
	}
	w.done <- true
}

func newWorker(id int, in <-chan job.Job, out chan<- int64, done chan<- bool) Worker {
	return Worker{id, in, out, done}
}
