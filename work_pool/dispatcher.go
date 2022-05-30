package workerpool

import "fmt"

type Dispatcher struct {
	pool   *Pool
	stoped chan bool
}

func newDispatcher(pool *Pool) *Dispatcher {
	dp := &Dispatcher{
		pool:   pool,
		stoped: make(chan bool),
	}
	pool.dispatcher = dp
	return dp
}

func (dp *Dispatcher) run() {
	for {
		select {
		case job, ok := <-dp.pool.pullJob():
			if !ok {
				continue
			}
			fmt.Printf("[WorkPool] - %s | get job id: %s \n", dp.pool.Name, job.ID)
			worker := <-dp.pool.pullFreeWorker()
			worker.process(job)
		case <-dp.stoped:
			// handler existing job in queue
			fmt.Printf("[WorkPool] - %s | job queue residue job size = %d, waiting all jobs done ... \n", dp.pool.Name, len(dp.pool.jobQueue))
			lastCount := len(dp.pool.jobQueue)
			if lastCount > 0 {
				for i := 0; i < lastCount; i++ {
					job := <-dp.pool.pullJob()
					fmt.Println("[WorkPool] get job id: ", job.ID)
					worker := <-dp.pool.pullFreeWorker()
					worker.process(job)
				}
			}
			dp.stoped <- true
			return
		}
	}
}

func (dp *Dispatcher) stop() {
	dp.stoped <- true
}
