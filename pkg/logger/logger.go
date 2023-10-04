package logger

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

// Logger contains logger configs
type Logger struct {
	*logrus.Logger
}

// NewLogger configs application logger
func NewLogger(fileName ...string) (*Logger, error) {
	logger := logrus.New()

	logger.SetReportCaller(true)

	if len(fileName) > 0 {
		logger.SetFormatter(&logrus.JSONFormatter{
			PrettyPrint: true,
		})

		file, err := os.OpenFile(fileName[0], os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		logger.SetOutput(file)
	}

	return &Logger{
		Logger: logger,
	}, nil
}

func (l *Logger) ServerError(w http.ResponseWriter, logMessage string, err error) {
	http.Error(w, "internal server error", http.StatusInternalServerError)
	l.Errorf("%s: %s", logMessage, err.Error())
}
