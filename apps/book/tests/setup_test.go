package tests

import (
	"go-base-structure/cmd/config"
	"os"
	"testing"
)

var cfg config.Config

func TestMain(m *testing.M) {
	cfg.Port = 8000

	os.Exit(m.Run())
}
