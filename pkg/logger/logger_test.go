package logger

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// isFileEmpty check if a file is empty
func isFileEmpty(filename string) (bool, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return false, err
	}

	return fileInfo.Size() == 0, nil
}

func TestNewLogger(t *testing.T) {
	// test logging to the file
	logr, err := NewLogger(fileName)
	if err != nil {
		t.Errorf("unexpected error while initial logger: %s", err.Error())
	}

	logr.Info("log something to the file")
	isEmptyFile, err := isFileEmpty(fileName)
	if err != nil {
		logr.Fatal("failed to check that file is empty or not: ", err.Error())
	}

	if isEmptyFile {
		t.Errorf("we log something to the file but the file is empty")
	}

	// test logging to the console
	logr, err = NewLogger()
	if err != nil {
		t.Errorf("unexpected error while initial logger: %s", err.Error())
	}

	logOutput := &bytes.Buffer{}
	logr.Out = logOutput

	logr.Info("log something to the console")

	logOutputStr := logOutput.String()
	if logOutputStr == "" {
		t.Errorf("we log something to the console but the console is empty")
	}
}

func TestLogger_ServerError(t *testing.T) {
	logr, _ := NewLogger(fileName)
	rr := httptest.NewRecorder()

	logr.ServerError(rr, "unable convert int to string", errors.New("bad int, lol"))

	// check for status code (should be 500)
	expectedStatusCode := http.StatusInternalServerError
	if rr.Code != expectedStatusCode {
		t.Errorf("expected %d status code but got %d", expectedStatusCode, rr.Code)
	}
}
