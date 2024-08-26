package network

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

func RunHttpExample() {
	go func() {
		time.Sleep(2 * time.Second)
		client()
	}()
	server()
}

func client() {
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

func server() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", defaultHandler)

	server := &http.Server{
		//":" is required here - listen on all available network interfaces (external IPs, localhost)
		Addr:         ":8080",
		Handler:      mux,
		IdleTimeout:  10 * time.Second, //time to wait for the next request when keep-alive is enabled
		ReadTimeout:  time.Second,      //specifies the maximum duration allowed to read the entire client request
		WriteTimeout: time.Second,      //specifies the maximum duration before timing out write of the response
	}

	log.Printf("server starting on %v", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("error starting server - %v", err)
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintf(w, `{"message": "server response"}`)
	if err != nil {
		log.Printf("error writing response - %v", err)
	}
	rq, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Printf("error dumping request - %v", err)
	}
	log.Printf("incoming request %v\n", string(rq))
}
