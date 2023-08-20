package logger

import (
	"log"
	"os"
)

type Logger struct {
	ErrorLog   *log.Logger
	InfoLog    *log.Logger
	WarningLog *log.Logger
}

// NewLogger customize our logger
func NewLogger() *Logger {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	warningLog := log.New(os.Stdout, "WARNING\t", log.Ldate|log.Ltime|log.Lshortfile)

	return &Logger{
		ErrorLog:   errorLog,
		InfoLog:    infoLog,
		WarningLog: warningLog,
	}
}
