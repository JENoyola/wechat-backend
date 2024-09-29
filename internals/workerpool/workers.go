package workerpool

type Job func()

// Worker will be in charge to complete an assigned task
func (w *WorkerPool) Worker(ID int) {
	defer w.wg.Done()

	for job := range w.jobQueue {
		job()
	}
}

// AssignJobToWorker assigns a job to a worker
func (w *WorkerPool) AssignJobToWorker(job Job) {
	w.jobQueue <- job
}
