package logger

import (
	"log"
	"os"
)

var (
	InfoLog  *log.Logger
	ErrorLog *log.Logger
)

func init() {
	newLogger()
}

func newLogger() {
	InfoLog = log.New(os.Stdout, "Info:\t", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	ErrorLog = log.New(os.Stderr, "Error:\t", log.Ldate|log.Lmicroseconds|log.Lshortfile)
}
