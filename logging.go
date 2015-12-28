package main

import (
	"os"
	"io"
	"log"
	"path"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Error   *log.Logger
)

func logging(directory string, initMessage string) {
	// Logging setup
	if err := os.MkdirAll(directory, 0777); err != nil {
		log.Fatal(err)
	}
	infoF, err := os.OpenFile(path.Join(directory, "compressor.info.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", infoF, ":", err)
	}
	errorF, err := os.OpenFile(path.Join(directory, "compressor.error.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", errorF, ":", err)
	}
	infoMulti := io.MultiWriter(infoF, os.Stdout)
	errMulti := io.MultiWriter(errorF, os.Stdout)

	Trace = log.New(os.Stdout,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoMulti,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errMulti,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info.Println(initMessage)
}
