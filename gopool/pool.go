package gopool

import (
	"fmt"
	"sync"
)

// Job represents the job to be executed by a worker.
type Job func()

// Pool is a worker group that runs jobs concurrently.
type Pool struct {
	wg     sync.WaitGroup
	jobCh  chan Job
	closed bool
}

// NewPool creates a new pool with the given number of workers.
func NewPool(numWorkers int) *Pool {
	p := &Pool{
		jobCh: make(chan Job, numWorkers),
	}

	p.wg.Add(numWorkers)
	for range numWorkers {
		go p.worker()
	}

	return p
}

// Submit submits a job to the pool for execution.
func (p *Pool) Submit(job Job) {
	if p.closed {
		fmt.Println("Pool: pool is stopped, do not submit the job.")
		return
	}
	p.jobCh <- job
}

// Stop stops the pool and waits for all workers to finish their jobs.
func (p *Pool) Stop() {
	p.wg.Wait()
}

// Close cloese the pool
func (p *Pool) Close() {
	if !p.closed {
		close(p.jobCh)
		p.closed = true
	}
}

// worker is a long-running goroutine that executes jobs from the job channel.
func (p *Pool) worker() {
	defer p.wg.Done()

	for job := range p.jobCh {
		job()
	}
}

func IsClosed[T any](ch <-chan T) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}
