package jobqueue

import (
	"log"
	"sync"
)

type Queue struct {
	channel chan interface{}
	wg      sync.WaitGroup
	worker  func(param interface{})
	blocked bool
	running bool
	closed  bool
	run     func(jobQueue *Queue)
}

// Set jobqueue worker function
func (queue *Queue) SetWorker(worker func(param interface{})) {
	queue.worker = worker
}

// Set if jobqueue can be blocked. Default is false
func (queue *Queue) SetBlocked(block bool) {
	queue.blocked = block
}

// Run jobqueue
func (queue *Queue) Run() bool {
	if queue.closed {
		queue.running = false
		log.Println("jobqueue is closed.")
		return false
	}

	if queue.worker == nil {
		log.Println("jobqueue is not set worker function yet.")
		return false
	}

	queue.running = true
	if queue.run == nil {
		queue.run = run
		go queue.run(queue)
	}
	return true
}

// Add a new job to jobqueue
func (queue *Queue) Add(param interface{}) bool {
	if queue.closed {
		log.Println("jobqueue is closed.")
		return false
	}

	if queue.running {
		if queue.blocked {
			queue.channel <- param
			queue.wg.Add(1)

		} else {
			select {
			case queue.channel <- param:
				queue.wg.Add(1)
				return true
			default:
				log.Println("jobqueue is full.")
				return false
			}
		}
	}

	log.Println("jobqueue is not start.")
	return false
}

// Wait for jobqueue finished
func (queue *Queue) Wait() {
	queue.wg.Wait()
}

// Start jobqueue
func (queue *Queue) Start() {
	if queue.run != nil {
		queue.running = true
	} else {
		log.Println("jobqueue is not running.")
	}
}

// Stop jobqueue accept a new job
func (queue *Queue) Stop() {
	queue.running = false
}

// Close jobqueue
func (queue *Queue) Close() {
	queue.running = false
	queue.closed = true
	close(queue.channel)
}

// New a job queue
func New(size int) Queue {
	return Queue{
		channel: make(chan interface{}, size),
		worker:  nil,
		blocked: false,
		running: false,
		closed:  false,
	}
}

func run(queue *Queue) {
	for job := range queue.channel {
		queue.worker(job)
		queue.wg.Done()
	}
}
