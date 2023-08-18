package tests

import (
	"go-base-structure/cmd/config"
	"go-base-structure/models"
	"os"
	"testing"
)

var testConfig config.Config

func TestMain(m *testing.M) {
	// create a new logger
	config.NewLogger()

	// initial test model
	models.NewTestModels()

	os.Exit(m.Run())
}
