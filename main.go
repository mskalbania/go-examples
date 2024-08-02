package main

import (
	"go-examples/concurrency"
	"go-examples/runner"
)

func main() {
	concurrency.StartTennisMatch()
	concurrency.StartRallyRace()
	concurrency.ProcessConcurrently()

	runner.DemonstrateRunner()

}
