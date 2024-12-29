package concurrency

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

/*
Client is a simple HTTP client that can be used to send multiple requests concurrently.
Both time limit and number of requests can be specified.
Finishes when either time limit is reached or all requests are sent.

	client := NewClient(
		2,
		0,
		10,
		func() *http.Request {
			req, _ := http.NewRequest("GET", "https://webhook.site/xxx-xxx", nil)
			return req
		},
	)
	client.Run()

This will send 10 GET requests to the specified URL with a concurrency of 2.
Result:

	Requests: 10, RPS: 0.00, Success: 0, Failed: 10
	Errors encountered (10 total):
		[10] unexpected status code: 429
*/
type Client struct {
	concurrency    int
	timeLimit      time.Duration
	numRequests    int
	client         *http.Client
	requestFactory func() *http.Request
	errors         struct {
		sync.Mutex
		list []error
	}
}

func NewClient(concurrency int, timeLimit time.Duration, numRequests int, requestFactory func() *http.Request) *Client {
	return &Client{
		concurrency:    concurrency,
		timeLimit:      timeLimit,
		numRequests:    numRequests,
		client:         &http.Client{},
		requestFactory: requestFactory,
	}
}

func (c *Client) Run() {
	var (
		workersWg       sync.WaitGroup
		reporterWg      sync.WaitGroup
		requestsCounter atomic.Int64
		successCount    atomic.Int64
		start           = time.Now()
		ctx, cancel     = context.WithCancel(context.Background())
	)

	if c.timeLimit > 0 {
		ctx, cancel = context.WithTimeout(ctx, c.timeLimit)
	}
	defer cancel()

	for i := 0; i < c.concurrency; i++ {
		workersWg.Add(1)
		go func() {
			defer workersWg.Done()
			for {
				if c.numRequests > 0 && requestsCounter.Load() >= int64(c.numRequests) {
					return
				}
				select {
				case <-ctx.Done():
					return
				default:
					requestsCounter.Add(1)
					if err := c.makeRequest(); err != nil {
						c.errors.Lock()
						c.errors.list = append(c.errors.list, err)
						c.errors.Unlock()
					} else {
						successCount.Add(1)
					}
				}
			}
		}()
	}

	reporterWg.Add(1)
	go trackProgress(&reporterWg, ctx, start, &successCount, &requestsCounter)

	workersWg.Wait()
	cancel() //isn't automatically cancelled when non timeout context
	reporterWg.Wait()

	c.printErrorSummary()
}

func (c *Client) makeRequest() error {
	resp, err := c.client.Do(c.requestFactory())
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) printErrorSummary() {
	if len(c.errors.list) > 0 {
		fmt.Printf("\nErrors encountered (%d total):\n", len(c.errors.list))
		errorCounts := make(map[string]int)
		for _, err := range c.errors.list {
			errorCounts[err.Error()]++
		}
		for errMsg, count := range errorCounts {
			fmt.Printf("  [%d] %s\n", count, errMsg)
		}
	}
}

func trackProgress(
	reporterWg *sync.WaitGroup,
	ctx context.Context,
	start time.Time,
	successCount *atomic.Int64,
	counter *atomic.Int64) {

	defer reporterWg.Done()
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			printProgress(start, successCount, counter)
			return
		case <-ticker.C:
			printProgress(start, successCount, counter)
		}
	}
}

func printProgress(start time.Time, successCount, counter *atomic.Int64) {
	elapsed := time.Since(start).Seconds()
	successReqs := successCount.Load()
	rps := float64(successReqs) / elapsed
	failedReqs := counter.Load() - successReqs
	fmt.Printf("\rRequests: %d, RPS: %.2f, Success: %d, Failed: %d", counter.Load(), rps, successReqs, failedReqs)
}
