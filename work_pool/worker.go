package workerpool

import (
	"fmt"
)

type Worker struct {
	id       int
	workPool *Pool
	jobChan  chan *Job
	stoped   chan bool
}

func newWorker(pool *Pool, id int) *Worker {
	return &Worker{
		id:       id,
		workPool: pool,
		jobChan:  make(chan *Job),
		stoped:   make(chan bool),
	}
}

func (w *Worker) start() {
	go func() {
		for {
			// put worker back to pool after finishing job
			w.workPool.addWorker(w)
			select {
			case job := <-w.jobChan:
				fmt.Printf("[Worker] WID: %d Get job ID: %s \n", w.id, job.ID)
				job.Func(job.Args...)
			case <-w.stoped:
				close(w.jobChan)
				w.stoped <- true
				fmt.Printf("[Worker] WID: %d Stoped \n", w.id)
				return
			}
		}
	}()
}

func (w *Worker) process(job *Job) {
	w.jobChan <- job
}

func (w *Worker) stop() {
	w.stoped <- true
}
