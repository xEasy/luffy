package workerpool

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type WorkersChan chan *Worker

type Pool struct {
	Name       string
	size       uint32
	wPool      WorkersChan
	jobQueue   chan *Job
	dispatcher *Dispatcher

	releaseLock    sync.Mutex
	released       bool
	releasedSignal chan struct{}

	started     bool
	startedLock sync.Mutex
	done        chan struct{}
}

func NewWorkPool(name string, poolSize uint32, maxJobSize uint32) *Pool {
	pool := &Pool{
		Name:     name,
		size:     poolSize,
		wPool:    make(WorkersChan, poolSize),
		jobQueue: make(chan *Job, maxJobSize),
		released: false,
		started:  false,
		done:     make(chan struct{}),
	}
	return pool
}

func (p *Pool) Done() <-chan struct{} {
	return p.done
}

func (p *Pool) Enqueue(job JobFunc, args ...any) {
	var jobID string
	// NOTE use NewRandom cos New will panic if fail
	jobUuid, err := uuid.NewRandom()
	if err != nil {
		jobID = "errJOBID"
	} else {
		jobID = jobUuid.String()
	}
	p.jobQueue <- &Job{
		ID:   jobID, // TODO random JobID
		Func: job,
		Args: args,
	}
}

// Start a worker pool
func (wp *Pool) Start() {
	wp.startedLock.Lock()
	if wp.started {
		wp.startedLock.Unlock()
		return
	} else {
		wp.started = true
	}
	wp.startedLock.Unlock()
	wp.startAndDispatch()

	fmt.Printf("[WorkPool] - %s  started and waiting for job ...\n", wp.Name)
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
	fmt.Printf("[WorkPool] - %s | releasing, waiting all worker's job done ... \n", p.Name)

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
		fmt.Printf("[WorkPool]- %s | free worker id: %d stoping\n", p.Name, worker.id)
		worker.stop()
		// waiting worker's stoped signal
		<-worker.stoped
		close(worker.stoped)
	}

	// 3 close workers chan
	fmt.Printf("[WorkPool] - %s | colse wPool \n", p.Name)
	close(p.wPool)
	close(p.done)
}

func (wp *Pool) addWorker(worker *Worker) {
	wp.wPool <- worker
}

func (p *Pool) pullJob() <-chan *Job {
	return p.jobQueue
}

func (p *Pool) pullFreeWorker() <-chan *Worker {
	return p.wPool
}
