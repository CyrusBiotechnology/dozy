package main

import (
	"log"
	"os"
)

var (
	Trace *log.Logger
	Info  *log.Logger
	Error *log.Logger
)

func logging(initMessage string) {
	// Logging setup
	Trace = log.New(os.Stdout,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(os.Stderr,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info.Println(initMessage)
}
