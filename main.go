package main

import (
	"go-examples/cmd/statistics"
	"log"
)

// injected with ldflags
// go build -ldflags "-X main.version=$VERSION"
var version string

func main() {
	log.Printf("Starting application version=%v", version)
	//concurrency.StartTennisMatch()
	//concurrency.StartRallyRace()
	//concurrency.ProcessConcurrently()

	//runner.DemonstrateRunner()
	//pool.DemonstratePool()
	//work.DemonstrateWork()
	statistics.RunReadExample()
}
