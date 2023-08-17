package config

import (
	"log"
	"os"
)

// NewLogger customize our logger
func NewLogger() {
	infoLog := log.New(os.Stdout, "Info\t", log.Ldate|log.Ltime)
	AppConfig.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	AppConfig.ErrorLog = errorLog

	warningLog := log.New(os.Stdout, "WARNING\t", log.Ldate|log.Ltime|log.Lshortfile)
	AppConfig.WarningLog = warningLog
}
