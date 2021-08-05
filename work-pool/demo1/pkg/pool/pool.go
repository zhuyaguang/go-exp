package pool

import (
	"demo1/pkg/work"
	"fmt"
	"sync"
)

// Pool ...
type Pool struct {
	Name string

	Size    int
	Workers []*work.Worker

	QueueSize int
	Queue     chan work.Job
}

// Initialize  ...
func (p *Pool) Initialize() {
	// maintain minimum 1 worker
	if p.Size < 1 {
		p.Size = 1
	}
	p.Workers = []*work.Worker{}
	for i := 1; i <= p.Size; i++ {
		worker := &work.Worker{
			ID:   i - 1,
			Name: fmt.Sprintf("%s-worker-%d", p.Name, i-1),
		}
		p.Workers = append(p.Workers, worker)
	}

	// maintain min queue size as 1
	if p.QueueSize < 1 {
		p.QueueSize = 1
	}
	p.Queue = make(chan work.Job, p.QueueSize)
}

// Start ...
func (p *Pool) Start() {
	for _, worker := range p.Workers {
		worker.Start(p.Queue)
	}
	fmt.Println("all workers started")
}

// Stop ...
func (p *Pool) Stop() {
	close(p.Queue) // close the queue channel

	var wg sync.WaitGroup
	for _, worker := range p.Workers {
		wg.Add(1)
		go func(w *work.Worker) {
			defer wg.Done()

			w.Stop()
		}(worker)
	}
	wg.Wait()
	fmt.Println("all workers stopped")
}
