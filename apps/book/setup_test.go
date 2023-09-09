package book

import (
	"github.com/sirupsen/logrus"
	"go-base-structure/cmd/settings"
	"go-base-structure/models"
	"go-base-structure/pkg/logger"
	"os"
	"testing"
)

// testApp is wide configuration for our book app tests
var testApp settings.Application

// TestMain is a base set up and configuration for tests
func TestMain(m *testing.M) {
	logger := &logger.Logger{Logger: logrus.New()}
	mdls := models.NewTestModels()

	testApp.Logger = logger
	testApp.Models = mdls

	bookApp = &testApp

	os.Exit(m.Run())
}
