#### Usage

```golang
  // make a new workerpool
  pool := workerpool.NewWorkPool("testBatch", 5, 10)

  // start dispatching job with no block
  pool.Start()

  // let pool to init complete
  time.Sleep(time.Microsecond * 100) 

  // an example job:
  // inc counter with step 2
  counter := 0
	for i := 0; i < 100; i++ {
  job := func(args ...any) {
    // run any func you want
    prime := args[0].(int32)
      atomic.CompareAndSwapInt32(&counter, counter, counter+prime)
    }

    pool.Enqueue(job, int32(2))
  }

  // call release if shutting down
  // release will block until all jobs done
  pool.Release()

  // no block release
  fmt.Println("Wating for workerPool stop & release.")
  go pool.Release()
  select {
    case: <- pool.Done()
      // this code will run until pool release done
      fmt.Println("All jobs done, workerPool has stoped & released.")
  }
```
