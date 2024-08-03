package pool

import (
	"errors"
	"io"
	"sync"
)

var ErrPoolClosed = errors.New("pool closed")

//Pool
/*
Manages a resources that can be shared by multiple goroutines (like db connections).
Resource must implement io.Closer interface.
Doesn't restrict the size but aims to maintain the specified targetSize.
*/
type Pool struct {
	lock      sync.Mutex                //used to wrap critical sections
	resources chan io.Closer            //represents actual resources being shared
	factory   func() (io.Closer, error) //used to create new resource when required
	closed    bool
	poolSize  uint
}

func New(targetSize uint, factory func() (io.Closer, error)) (*Pool, error) {
	if targetSize == 0 {
		return nil, errors.New("cannot create empty pool")
	}
	return &Pool{
		factory:   factory,
		resources: make(chan io.Closer, targetSize),
		poolSize:  targetSize,
	}, nil
}

func (p *Pool) Acquire() (io.Closer, error) {
	select {
	//this tries to return shared resource if ones available in chanel
	case r, ok := <-p.resources: //ok can be retrieved to check if channel is still opened
		if !ok {
			return nil, ErrPoolClosed
		}
		return r, nil
	//then if it fails creates a new one
	default:
		return p.factory()
	}
}

func (p *Pool) Close() {
	//wraps entire method as critical section
	//prevents race condition with Release function - both access shared bool & chan resource values
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.closed { //just returns if called multiple times by different go routines
		return
	}
	p.closed = true                     //marks as closed and
	close(p.resources)                  //closing channel before draining it, if omitted below "for" would create deadlock
	for resource := range p.resources { //this blocks when empty, hence above
		resource.Close()
	}
}

func (p *Pool) Release(resource io.Closer) {
	p.lock.Lock()
	defer p.lock.Unlock()
	//if pool marked as closed, we can't return it to channel since it was already drained and closed.
	//closing leftover resource here
	if p.closed {
		resource.Close()
		return
	}
	select {
	case p.resources <- resource: //it tries to put the resource back to chanel
	default:
		//when it fails it means buffer is full and resource should be discarded to remain at target pool size
		resource.Close()
	}
}
