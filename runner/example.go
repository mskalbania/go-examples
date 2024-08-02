package runner

import (
	"errors"
	"fmt"
	"os"
	"time"
)

func DemonstrateRunner() {
	const timeout = 2 * time.Second
	runner := New(timeout)
	runner.Add(func(i int) {
		fmt.Printf("\nExecuting task with id - %v", i)
		time.Sleep(1 * time.Second)
	}, func(i int) {
		fmt.Printf("\nExecuting task with id - %v", i)
		time.Sleep(2 * time.Second)
	})
	if err := runner.Start(); err != nil {
		fmt.Printf("\nError - %v", err)
		switch {
		case errors.Is(err, TimeoutError):
			os.Exit(1)
		case errors.Is(err, InterruptedError):
			os.Exit(2)
		}
	}
}
