package main

import (
	"io"
	"log"
	"os"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	Access  *log.Logger
)

func InitLogging(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer,
	accessHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Access = log.New(accessHandle,
		"ACCESS: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func NewLogger(logfile string, prefix string) (*log.Logger, error) {
	// log to stdout
	if logfile == "-" {
		logger := log.New(os.Stderr, prefix+": ", log.Ldate|log.Ltime|log.Lshortfile)
		return logger, nil
	}
	// log to file
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open logfile", logfile, ":", err)
		return nil, err
	}
	logger := log.New(file, prefix, log.Ldate|log.Ltime|log.Lshortfile)
	return logger, nil
}
