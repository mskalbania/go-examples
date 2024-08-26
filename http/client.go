package http

import (
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

func RunClient() {
	client := &http.Client{
		//all requests timeout
		Timeout: 1 * time.Second,
	}

	//more granular timeout and possibility to cancel the request at any time
	//ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	//defer cancel()
	//rq, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080", strings.NewReader(`{"message": "client request"}`))

	rq, err := http.NewRequest("GET", "http://localhost:8080", strings.NewReader(`{"message": "client request"}`))
	if err != nil {
		log.Fatalf("error creating request - %v", err)
	}
	rq.Header.Set("Content-Type", "application/json")

	rs, err := client.Do(rq)
	if err != nil {
		log.Fatalf("error sending request - %v", err)
	}
	defer rs.Body.Close()
	//dump used for testing purposes
	rsDump, err := httputil.DumpResponse(rs, true)
	//normally - b, err := io.ReadAll(rs.Body) and string(b)
	if err != nil {
		log.Fatalf("error dumping response - %v", err)
	}
	log.Printf("response: %v\n", string(rsDump))
}
