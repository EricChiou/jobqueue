package main

import (
	"fmt"
	"jobqueue/jobqueue"
	"time"
)

func main() {
	jobQueue := jobqueue.New(1024) // new a jobqueue
	jobQueue.SetWorker(worker)     // set worker function
	jobQueue.Run()                 // run jobqueue

	jobQueue.Add(1) // add a new job
	jobQueue.Add(2)

	jobQueue.Stop() // stop jobqueue accept a new job
	jobQueue.Add(3) // can not add a new job after stop the jobqueue

	jobQueue.Start() // start jobqueue again
	jobQueue.Add(4)

	jobQueue.Close() // close jobqueue
	jobQueue.Run()   // can not run jobqueue again after close it
	jobQueue.Add(5)  // can not add a new job after close the jobqueue

	fmt.Println("all jobs added")
	jobQueue.Wait() // wait for jobqueue finished
	fmt.Println("all jobs finished")
}

func worker(param interface{}) {
	fmt.Println("job", param, "start")
	time.Sleep(2 * time.Second)
	fmt.Println("job", param, "finished")
}
