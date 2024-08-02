package concurrency

import (
	"fmt"
	"sync"
	"time"
)

func StartRallyRace() {
	var wg sync.WaitGroup
	wg.Add(1)
	baton := make(chan int) //simulates exchange of the baton between runners

	go run(baton, &wg) //first runner prepares

	baton <- 1 //first runner starts

	wg.Wait()
}

func run(baton chan int, wg *sync.WaitGroup) {
	var next int
	//runner receives the baton and starts running
	runner := <-baton
	fmt.Printf("Runner %d now running\n", runner)

	//next runner prepares (so there is another goroutine which waits on chanel exchange)
	if runner != 4 {
		next = runner + 1
		go run(baton, wg)
	}

	//time passes while running around
	time.Sleep(1 * time.Second)

	//if runner 4 reaches the line, the race is over
	if runner == 4 {
		fmt.Println("Race done 4th reach the end")
		wg.Done()
		return
	}
	//else runner reached the next one and passed the baton
	fmt.Printf("Passing baton between %d -> %d\n", runner, next)
	baton <- next
}
