package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

var (
	Trace *log.Logger
	Info  *log.Logger
	Error *log.Logger
)

func logging(directory string, initMessage string) {
	// Logging setup
	if err := os.MkdirAll(directory, 0777); err != nil {
		log.Fatal(err)
	}
	infoF, err := os.OpenFile(
		path.Join(directory, fmt.Sprint(path.Base(os.Args[0]), ".log")), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", infoF, ":", err)
	}
	logMulti := io.MultiWriter(infoF, os.Stdout)

	Trace = log.New(os.Stdout,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(logMulti,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(logMulti,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info.Println(initMessage)
}
