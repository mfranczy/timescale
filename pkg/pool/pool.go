package pool

import (
	"database/sql"

	"timescale/pkg/job"
)

const buffer = 10

// Pool structure represents workers and channels
type Pool struct {
	w       []Worker
	wNum    int
	wDoneCh chan bool
	inCh    []chan job.Job
	OutCh   chan int64
}

// New returns a new Pool
func New(wNum int) *Pool {
	w := make([]Worker, wNum)
	inCh := make([]chan job.Job, wNum)
	outCh := make(chan int64, buffer*wNum)
	wDoneCh := make(chan bool, wNum)

	for i := 0; i < wNum; i++ {
		inCh[i] = make(chan job.Job, buffer)
		w[i] = newWorker(i+1, inCh[i], outCh, wDoneCh)
	}
	return &Pool{w, wNum, wDoneCh, inCh, outCh}
}

// Run worker pool
func (p *Pool) Run(jobs chan job.Job, db *sql.Stmt) {
	p.workersStart(db)
	go p.workersWait()
	go p.dispatch(jobs)
}

func (p *Pool) workersStart(db *sql.Stmt) {
	for _, w := range p.w {
		go w.Exec(db)
	}
}

func (p *Pool) workersWait() {
	for i := 1; i <= p.wNum; i++ {
		<-p.wDoneCh
		if i == p.wNum {
			close(p.OutCh)
		}
	}
}

func (p *Pool) dispatch(jobs chan job.Job) {
	hostWorker := make(map[string]int)
	i := 0

	for j := range jobs {
		workerID := i
		if id, ok := hostWorker[j.Hostname]; ok {
			workerID = id
		} else {
			hostWorker[j.Hostname] = workerID
		}

		p.inCh[workerID] <- j
		i++
		if i >= p.wNum {
			i = 0
		}
	}
	for _, c := range p.inCh {
		close(c)
	}
}
