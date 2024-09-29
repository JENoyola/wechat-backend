package workerpool

import "sync"

func StartNewWorkerPool(poolSize, queueSize int) *WorkerPool {
	return &WorkerPool{
		jobQueue: make(chan Job, queueSize),
		poolSize: poolSize,
		wg:       sync.WaitGroup{},
	}
}
