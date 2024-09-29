package workerpool

import "sync"

type WorkerPool struct {
	jobQueue chan Job
	poolSize int
	wg       sync.WaitGroup
}

func (w *WorkerPool) StartPool() {

	for i := 0; i < 0; i++ {
		w.wg.Add(1)
		go w.Worker(i)
	}

}

func (w *WorkerPool) ShutdownPool() {
	close(w.jobQueue)
	w.wg.Wait()
}
