package concurrency

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

func ProcessConcurrently() {
	concurrencyLevel := 3
	workLoad := 10

	var wg sync.WaitGroup
	wg.Add(concurrencyLevel)
	workQueue := make(chan string, workLoad)

	for i := 1; i <= concurrencyLevel; i++ { //spawn worker group
		go process(fmt.Sprintf("worker-%d", i), workQueue, &wg)
	}

	for i := 1; i <= workLoad; i++ { //filling in the workQueue (buffered channel)
		workQueue <- fmt.Sprintf("task-%d", i)
	}
	close(workQueue) //closing channel here - no more send but receive still possible

	wg.Wait()
}

func process(workerId string, workQueue chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for { //infinite processing loop
		item := <-workQueue
		//we can either check for zero value or use "closed" flag - item, closed :=...
		if item == "" {
			//no more items to process can exit
			fmt.Printf("Worker %s shutting down.\n", workerId)
			break
		}
		fmt.Printf("Processing [%s]\n", item)
		time.Sleep(time.Duration(rand.IntN(500)) * time.Millisecond)
	}
}
