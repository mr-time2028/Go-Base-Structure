package logger

import (
	"os"
	"testing"
)

var fileName = "logFile.log"

func removeFileIfExists(filename string) error {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if err = os.Remove(filename); err != nil {
		return err
	}

	return nil
}

func setUpTest() {}

func tearDownTest() {
	_ = removeFileIfExists(fileName)
}

func TestMain(m *testing.M) {
	setUpTest()
	exitCode := m.Run()
	tearDownTest()
	os.Exit(exitCode)
}
