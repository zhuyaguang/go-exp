package main

import (
	"fmt"
	"sync"
)

func main() {
	pool := &Pool{
		Name:      "test",
		Size:      5,
		QueueSize: 20,
	}
	pool.Initialize()
	pool.Start()
	defer pool.Stop()

	for i := 1; i <= 100; i++ {
		job := &PrintJob{
			Index: i,
		}
		pool.Queue <- job
	}
}

// PrintJob ...
type PrintJob struct {
	Index int
}

func (pj *PrintJob) Start(worker *Worker) error {

	fmt.Printf("job %s - %d\n", worker.Name, pj.Index)
	return nil
}



// Pool ...
type Pool struct {
	Name string

	Size    int
	Workers []*Worker

	QueueSize int
	Queue     chan Job
}

// Initialize  ...
func (p *Pool) Initialize() {
	// maintain minimum 1 worker
	if p.Size < 1 {
		p.Size = 1
	}
	p.Workers = []*Worker{}
	for i := 1; i <= p.Size; i++ {
		worker := &Worker{
			ID:   i - 1,
			Name: fmt.Sprintf("%s-worker-%d", p.Name, i-1),
		}
		p.Workers = append(p.Workers, worker)
	}

	// maintain min queue size as 1
	if p.QueueSize < 1 {
		p.QueueSize = 1
	}
	p.Queue = make(chan Job, p.QueueSize)
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
		go func(w *Worker) {
			defer wg.Done()

			w.Stop()
		}(worker)
	}
	wg.Wait()
	fmt.Println("all workers stopped")
}



// Worker ...
type Worker struct {
	ID       int
	Name     string
	StopChan chan bool
}

// Start ...
func (w *Worker) Start(jobQueue chan Job) {
	w.StopChan = make(chan bool)
	successChan := make(chan bool)

	go func() {
		successChan <- true
		for {
			// take job
			job := <-jobQueue
			if job != nil {
				job.Start(w)
			} else {
				fmt.Printf("worker %s to be stopped\n", w.Name)
				w.StopChan <- true
				break
			}
		}
	}()

	// wait for the worker to start
	<-successChan
}

// Stop ...
func (w *Worker) Stop() {
	// wait for the worker to stop, blocking
	_ = <-w.StopChan
	fmt.Printf("worker %s stopped\n", w.Name)
}

// Job ...
type Job interface {
	Start(worker *Worker) error
}
