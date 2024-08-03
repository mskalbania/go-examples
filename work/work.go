package work

import (
	"fmt"
	"sync"
)

//Pool
/*
Allows concurrent processing of work by multiple go routines.
Backed by unbuffered channel.
Guarantees that work is scheduled to some worker, exchange (send/receive) must happen between goroutine.
This is a different approach than queue based pool backed by buffered channel.
*/
type Pool struct {
	work chan Worker
	wg   sync.WaitGroup
}

type Worker interface { //types that want to use worker pool need to implement it
	Task()
}

func New(maxWorkers int) *Pool {
	pool := &Pool{
		work: make(chan Worker), //using unbuffered channel, as mentioned above to make send/receive blocking
	}
	pool.wg.Add(maxWorkers)
	for i := 1; i <= maxWorkers; i++ {
		go doWork(i, pool) //this spawns requested amount of workers
	}
	return pool
}

func doWork(id int, pool *Pool) {
	for work := range pool.work { //this blocks until work is submitted
		work.Task()
		fmt.Printf("Worker [%d] compleated task\n", id)
	}
	pool.wg.Done() //this is reached only when channel is closed (for loop interrupted)
	fmt.Printf("Worker [%d] shutdown\n", id)
}

func (p *Pool) Submit(work Worker) {
	p.work <- work //this is blocking until there is available worker go routine to pick up the task
}

func (p *Pool) Shutdown() {
	close(p.work)
	p.wg.Wait()
}
