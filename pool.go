package main

// ThreadPool allows to create a max amount of threads.
// Users can query how many threads they are allowed to create.
type ThreadPool struct {
	pool chan struct{}
	max  int
}

// NewThreadPool creates a new Threadpool with
func NewThreadPool(maxThreads int) *ThreadPool {
	tp := ThreadPool{
		pool: make(chan struct{}, maxThreads),
		max:  maxThreads,
	}
	for i := 0; i < maxThreads; i++ {
		tp.pool <- struct{}{}
	}
	return &tp
}

// Get requests `tasks` amount of threads from the pool.
// If the pool is not used much, a caller can get up to 1/3 of the available threads.
// Otherwise the caller gets only a single thread (once available).
// If the pool has more than `tasks` threads available it will only return `tasks`.
// It uses len(chan) which is a bit racy but shouldn't matter to much.
func (t *ThreadPool) Get(tasks int) int {
	threads := 1
	if len(t.pool) > t.max/2 {
		threads = len(t.pool) / 3
	}
	if tasks > 0 && threads > tasks {
		threads = tasks
	}
	for i := 0; i < threads; i++ {
		<-t.pool
	}
	return threads
}

// Put returns n threads back to the pool.
func (t *ThreadPool) Put(threads int) {
	for i := 0; i < threads; i++ {
		t.pool <- struct{}{}
	}
}
