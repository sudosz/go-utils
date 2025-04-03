package gopool

import (
	"fmt"
	"sync"
)

// Job represents a function to be executed by a worker.
type Job func()

// Pool manages a group of workers that execute jobs concurrently.
type Pool struct {
	wg     sync.WaitGroup
	jobCh  chan Job
	closed bool
}

// NewPool initializes a pool with the specified number of workers.
// Optimization: Pre-allocates buffered channel to reduce contention.
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

// Submit adds a job to the pool for execution unless the pool is closed.
// Optimization: Buffered channel reduces blocking under load.
func (p *Pool) Submit(job Job) {
	if p.closed {
		fmt.Println("Pool: pool is stopped, do not submit the job.")
		return
	}
	p.jobCh <- job
}

// Stop waits for all workers to complete their current jobs.
// Optimization: Uses WaitGroup for efficient synchronization.
func (p *Pool) Stop() {
	p.wg.Wait()
}

// Close shuts down the pool by closing the job channel.
// Optimization: Ensures no new jobs are accepted efficiently.
func (p *Pool) Close() {
	if !p.closed {
		close(p.jobCh)
		p.closed = true
	}
}

// worker runs in a goroutine, processing jobs from the channel until it closes.
// Optimization: Defers wg.Done to ensure cleanup even on panic.
func (p *Pool) worker() {
	defer p.wg.Done()
	for job := range p.jobCh {
		job()
	}
}

// IsClosed checks if the channel is closed by attempting a non-blocking receive.
// Optimization: Uses select for efficient closure detection.
func IsClosed[T any](ch <-chan T) bool {
	select {
	case <-ch:
		return true
	default:
	}
	return false
}
