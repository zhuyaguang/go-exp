package demo1

import (
	"fmt"
)

func demo1() {
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








