package work

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const (
	workers = 3
	tasks   = 20
)

func DemonstrateWork() {
	var wg sync.WaitGroup
	wg.Add(tasks)

	pool := New(workers)

	for i := 0; i < tasks; i++ {
		go func() {
			printer := newPrinter()
			pool.Submit(printer)
			wg.Done()
		}()
	}
	wg.Wait()
	pool.Shutdown()
}

type printer struct {
	id int32
}

var counter int32

func newPrinter() *printer {
	atomic.AddInt32(&counter, 1)
	return &printer{
		id: counter,
	}
}

func (p *printer) Task() {
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
}
