package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

type Logger struct {
	*logrus.Logger
}

// NewLogger configs application logger
func NewLogger() *Logger {
	logger := logrus.New()

	logger.SetReportCaller(true)
	logger.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})

	file, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Fatal("failed to open log file:", err)
	}
	logger.SetOutput(file)

	return &Logger{
		Logger: logger,
	}
}
