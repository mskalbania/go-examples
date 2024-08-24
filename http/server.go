package http

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

func RunServer() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", defaultHandler)

	server := &http.Server{
		//":" is required here - listen on all available network interfaces (external IPs, localhost)
		Addr:         ":8080",
		Handler:      mux,
		IdleTimeout:  10 * time.Second, //time to wait for the next request when keep-alive is enabled
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
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
