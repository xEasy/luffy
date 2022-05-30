package workerpool

import (
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

var counter int32

func setup() {
	counter = 0
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)
}

func TestNewPool(t *testing.T) {
	pool := NewWorkPool(5, 10)
	pool.Start()
	time.Sleep(time.Microsecond * 100) // let pool complete init

	// inc counter with step 2
	for i := 0; i < 10; i++ {
		job := func(args ...any) {
			prime := args[0].(int32)
			atomic.CompareAndSwapInt32(&counter, counter, counter+prime)
		}

		pool.Enqueue(job, int32(2))
	}

	// release should wait job done
	pool.Release()

	// time.Sleep(time.Second * 2)
	if counter != 20 {
		t.Fatal("pool work fail, counter: ", counter)
	}
}

func TestNewWorker(t *testing.T) {
	pool := NewWorkPool(100, 100)
	worker := newWorker(pool, 100)
	worker.start()

	// worker should regiest to pool's workers channel
	worker = <-pool.wPool
	if worker == nil {
		t.Fatal("worker regiest to pool's wPool channel fail")
	}

	// worker should call
	callResult := false
	wait := make(chan bool)

	callback := func(args ...interface{}) {
		callResult = true
		wait <- true
	}

	job := Job{
		ID:   "hello",
		Func: callback,
	}
	worker.process(&job)

	// wait for job done
	<-wait

	if !callResult {
		t.Fatalf("worker process fail!")
	}
}
