package network

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

func RunDialListenTcpIp() {
	go listen()
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go dial(&wg)
	}
	wg.Wait()
}

func dial(wg *sync.WaitGroup) {
	defer wg.Done()
	//this is connect timout
	c, err := net.DialTimeout("tcp", "localhost:8080", 1*time.Second)
	if err != nil {
		log.Printf("error dialing %v", err)
	}
	//this is read/write timeout
	err = c.SetDeadline(time.Now().Add(1 * time.Second))
	if err != nil {
		log.Printf("error setting deadline %v", err)
	}
	b, err := io.ReadAll(c)
	if err != nil {
		log.Printf("error reading %v", err)
	}
	fmt.Println(string(b))
}

func listen() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("error creating listener %v", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error accepting connection %v", err)
		}
		go func(c net.Conn) { //async handling, release accept for next connection, handle current connection async
			defer c.Close()
			log.Printf("handling connection from %v", c.RemoteAddr())
			time.Sleep(time.Duration(950+rand.Intn(150)) * time.Millisecond)
			_, err := fmt.Fprintf(c, "current time: %v", time.Now().Format("2006-01-02 15:04:05"))
			if err != nil {
				log.Printf("error writing %v", err)
			}
		}(conn)
	}
}
