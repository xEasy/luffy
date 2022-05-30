package workerpool

import (
	"fmt"
	"sync"
)

type WorkersChan chan *Worker

type Pool struct {
	size       uint32
	wPool      WorkersChan
	jobQueue   chan Job
	dispatcher *Dispatcher

	releaseLock    sync.Mutex
	released       bool
	releasedSignal chan struct{}

	started     bool
	startedLock sync.Mutex
}

func NewWorkPool(poolSize uint32, maxJobSize uint32) *Pool {
	pool := &Pool{
		size:     poolSize,
		wPool:    make(WorkersChan, poolSize),
		jobQueue: make(chan Job, maxJobSize),
		released: false,
		started:  false,
	}
	return pool
}

func (p *Pool) Enqueue(job JobFunc, args ...any) {
	p.jobQueue <- Job{
		ID:   "jobId", // TODO random JobID
		Func: job,
		Args: args,
	}
}

// Start a worker pool
func (wp *Pool) Start() {
	wp.startedLock.Lock()
	if wp.started {
		return
	} else {
		wp.started = true
	}
	wp.startedLock.Unlock()
	wp.startAndDispatch()
}

func (p *Pool) startAndDispatch() {

	// init pool workers by poolSize
	for i := 0; i < int(p.size); i++ {
		worker := newWorker(p, i)
		worker.start()
	}

	// init dispatcher and run
	dispatcher := newDispatcher(p)
	go dispatcher.run()
}

func (p *Pool) Release() {
	fmt.Println("WorkPool releasing, waiting all worker's job done ...")

	// 0 set release flag
	p.releaseLock.Lock()
	if p.released {
		return
	} else {
		p.released = true
	}
	p.releaseLock.Unlock()

	// 1 close jobQueue to prevent new job processing
	close(p.jobQueue)
	p.dispatcher.stop()
	<-p.dispatcher.stoped

	// 2 stop all worker
	for i := 0; i < cap(p.wPool); i++ {
		worker := <-p.pullFreeWorker()
		fmt.Printf("[WorkPool] free worker id: %d stoping\n", worker.id)
		worker.stop()
		// waiting worker's stoped signal
		<-worker.stoped
		close(worker.stoped)
	}

	// 3 close workers chan
	fmt.Println("[WorkerPool] colse wPool")
	close(p.wPool)
}

func (wp *Pool) addWorker(worker *Worker) {
	wp.wPool <- worker
}

func (p *Pool) pullJob() <-chan Job {
	return p.jobQueue
}

func (p *Pool) pullFreeWorker() <-chan *Worker {
	return p.wPool
}
