package logging

import (
	"io"
	"log"
	"os"
)

var (
	Trace *log.Logger
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
)

func init() {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error creating/opening log file - %v", err)
	}

	//discard here could be used to disable logging for that level
	Trace = log.New(io.Discard, "[TRACE] ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Warn = log.New(os.Stdout, "[WARN] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(os.Stderr, file), "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func logging() {
	Info.Println("info") //standard log call
	//log and return code 1 (os.Exit(1)), non-recoverable immediately stops execution
	Error.Fatalf("error")
	//log and call panic, recoverable (when recover used in deferred function
	Error.Panicf("error")
}
