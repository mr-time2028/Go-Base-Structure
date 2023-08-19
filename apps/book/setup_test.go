package book

import (
	"go-base-structure/cmd/config"
	"go-base-structure/models"
	"go-base-structure/pkg/logging"
	"os"
	"testing"
)

var testApp config.Application

func TestMain(m *testing.M) {
	// create a new logging
	logger := logging.NewLogger()
	mdls := models.NewTestModels()

	testApp.Logger = logger
	testApp.Models = mdls

	bookApp = &testApp

	os.Exit(m.Run())
}
