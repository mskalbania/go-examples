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
		Timeout: 1 * time.Second,
	}

	rq, err := http.NewRequest("GET", "http://localhost:8080", strings.NewReader(`{"message": "client request"}`))
	if err != nil {
		log.Printf("error creating request - %v", err)
	}
	rq.Header.Set("Content-Type", "application/json")

	rs, err := client.Do(rq)
	if err != nil {
		log.Printf("error sending request - %v", err)
	}
	defer rs.Body.Close()
	rsDump, err := httputil.DumpResponse(rs, true)
	if err != nil {
		log.Printf("error dumping response - %v", err)
	}
	log.Printf("response: %v\n", string(rsDump))
}
