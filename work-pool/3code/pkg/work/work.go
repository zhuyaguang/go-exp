package work

import "fmt"

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
