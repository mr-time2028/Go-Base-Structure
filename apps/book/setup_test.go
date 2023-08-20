package book

import (
	"go-base-structure/cmd/config"
	"go-base-structure/models"
	"go-base-structure/pkg/logger"
	"os"
	"testing"
)

var testApp config.Application

func TestMain(m *testing.M) {
	// create a new logger
	logger := logger.NewLogger()
	mdls := models.NewTestModels()

	testApp.Logger = logger
	testApp.Models = mdls

	bookApp = &testApp

	os.Exit(m.Run())
}
