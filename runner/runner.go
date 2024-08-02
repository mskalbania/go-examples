package runner

import (
	"errors"
	"os"
	"os/signal"
	"time"
)

var (
	TimeoutError     = errors.New("received timeout")
	InterruptedError = errors.New("received interrupted signal")
)

//Runner
/*
Allows to schedule sequential tasks to run concurrently,
while also monitoring specified time boundary and OS interrupt signal.
*/
type Runner struct {
	interrupt chan os.Signal //reports signal from operating system
	complete  chan error     //reports processing is finished, returns nil (success) or error (failed)
	//reports the amount of time runner has to complete its task
	//arrow here forces the direction of data flow - it can only receive
	timeout time.Duration
	tasks   []func(int) //a set of functions runner will execute in indexed order
}

func New(timeout time.Duration) *Runner {
	return &Runner{
		//buffered channel, it is non-blocking - goroutine consumes it when ready.
		//also size 1 since it is enough to process 1 interrupt signal
		interrupt: make(chan os.Signal, 1),
		//unbuffered blocking channel when goroutine finish (either with error or success) it will wait for receiver
		complete: make(chan error),
		timeout:  timeout,
	}
}

func (r *Runner) Add(task ...func(int)) { //var arg Go operator becomes slice within a function
	r.tasks = append(r.tasks, task...) //unpacks slice to (1,2,3)
}

func (r *Runner) Start() error {
	//this creates & start timer that will push time value to channel once time limit expires
	timer := time.After(r.timeout)
	//registers that we want to listen for all interrupt signals on this channel
	signal.Notify(r.interrupt, os.Interrupt)
	go func() {
		err := r.run() //invokes tasks, pushes status to complete
		r.complete <- err
	}()
	select { //this select listens on both below channels, if any message arrives returns.
	case err := <-r.complete:
		return err
	case <-timer:
		return TimeoutError
	}
}

func (r *Runner) run() error {
	for id, task := range r.tasks {
		if r.isInterrupted() {
			return InterruptedError
		}
		task(id + 1) //executes task giving it id
	}
	return nil //success
}

func (r *Runner) isInterrupted() bool {
	//using select so call to receive from channel interrupt is non-blocking
	//it is buffered channel but calling receive on empty blocks coroutine
	select {
	case <-r.interrupt: //check if interrupt signal came in
		signal.Stop(r.interrupt) //this makes sure no more signals are received on interrupt channel
		return true
	default: //if not just continue
		return false
	}
}
