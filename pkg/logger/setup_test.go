package logger

import (
	"log"
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
	err := removeFileIfExists(fileName)
	if err != nil {
		log.Fatal("testTearDown error while remove log file if exists: ", err.Error())
	}
}

func TestMain(m *testing.M) {
	setUpTest()
	exitCode := m.Run()
	tearDownTest()
	os.Exit(exitCode)
}
