package workerpool

import (
	"fmt"
	"time"
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
				fmt.Printf("[WorkPool] - %s, WID: %d processing job ID: %s \n", w.workPool.Name, w.id, job.ID)
				now := time.Now()
				job.Func(job.Args...)

				duration := time.Now().Sub(now)
				fmt.Printf("[WorkPool] - %s, WID: %d job ID: %s done in %s \n", w.workPool.Name, w.id, job.ID, duration)
			case <-w.stoped:
				close(w.jobChan)
				w.stoped <- true
				fmt.Printf("[WorkPool] - %s,  WID: %d Stoped \n", w.workPool.Name, w.id)
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
