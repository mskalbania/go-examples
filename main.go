package main

import (
	_ "embed"
	"go-examples/rest"
	"log"
	"regexp"
)

// injected with ldflags
// go build -ldflags "-X main.version=$VERSION"
var version string

// example of including some file in binary
//
//go:embed go.mod
var embedFile string

func main() {
	log.Printf("Application version=%v", version)
	match := regexp.MustCompile(`go (\d+\.\d+)`).FindStringSubmatch(embedFile)
	if match != nil && len(match) > 1 {
		log.Printf("Go version: %v", match[1])
	}
	//concurrency.StartTennisMatch()
	//concurrency.StartRallyRace()
	//concurrency.ProcessConcurrently()

	//runner.DemonstrateRunner()
	//pool.DemonstratePool()
	//work.DemonstrateWork()
	//network.RunDialListenTcpIp()
	//network.RunHttpExample()
	//network.RunWebsocketExample()
	rest.StartRestAPIExample()
}
