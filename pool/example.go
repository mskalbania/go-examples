package pool

import (
	"fmt"
	"io"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const (
	maxGoroutines = 50
	poolSize      = 2
)

type dbConnection struct {
	ID int32
}

func DemonstratePool() {
	var wg sync.WaitGroup
	wg.Add(maxGoroutines)

	pool, _ := New(poolSize, createDbConnection)

	for i := 0; i < maxGoroutines; i++ {
		go func() {
			//don't query at the same time, differentiate it a bit
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			performQueries(pool)
			wg.Done()
		}()
	}

	wg.Wait()
	pool.Close()
	fmt.Printf("\n"+
		"--------------------------------------------------------\n"+
		"Concurrency [%v], Pool size [%v], Connections opened [%v]"+
		"\n--------------------------------------------------------\n", maxGoroutines, poolSize, counter)
}

func (c *dbConnection) Close() error {
	return nil
}

var counter int32

func createDbConnection() (io.Closer, error) {
	atomic.AddInt32(&counter, 1)
	return &dbConnection{ID: counter}, nil
}

func performQueries(pool *Pool) {
	conn, err := pool.Acquire()
	if err != nil {
		fmt.Println(err)
	}
	defer pool.Release(conn)
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	fmt.Printf("Quering DB Using ConnID[%d]\n", conn.(*dbConnection).ID)
}
