package lua

// WorkerPool manages concurrent execution with bounded concurrency.
type WorkerPool struct {
	slots chan struct{}
}

// NewWorkerPool creates a new worker pool with the specified max concurrent workers.
// If maxConcurrent <= 0, it defaults to 10.
func NewWorkerPool(maxConcurrent int) *WorkerPool {
	if maxConcurrent <= 0 {
		maxConcurrent = 10
	}

	return &WorkerPool{
		slots: make(chan struct{}, maxConcurrent),
	}
}

// Acquire blocks until a worker slot is available.
func (p *WorkerPool) Acquire() {
	p.slots <- struct{}{}
}

// Release releases a worker slot.
func (p *WorkerPool) Release() {
	<-p.slots
}

// Close closes the worker pool.
func (p *WorkerPool) Close() {
	close(p.slots)
}
