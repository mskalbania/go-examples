package concurrency

import (
	"fmt"
	"math/rand/v2"
	"sync"
)

func StartTennisMatch() {
	var wg sync.WaitGroup
	wg.Add(2)

	fmt.Println("Starting tennis game...")

	court := make(chan int) //channel represents the court where balls pass between players (go routines)

	go player(court, "John", &wg)
	go player(court, "Jo", &wg)

	court <- 1 //simulates ball serve by first player

	wg.Wait()
}

func player(court chan int, name string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		ball, open := <-court
		//no ball coming at player - win
		if !open {
			fmt.Printf("Player %v won the game", name)
			break
		}
		//ball received now try to swing and throw it back
		random := rand.IntN(100)
		if random%13 == 0 {
			fmt.Printf("%s missed\n", name) //miss - lose
			close(court)
			break
		}
		//ball successfully hit, game is still on
		fmt.Printf("Player %s hit %d\n", name, ball)
		ball++
		court <- ball
	}
}
