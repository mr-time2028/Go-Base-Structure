package user

import (
	"github.com/sirupsen/logrus"
	"go-base-structure/cmd/settings"
	"go-base-structure/models"
	"go-base-structure/pkg/auth"
	"go-base-structure/pkg/logger"
	"net/http"
	"os"
	"testing"
	"time"
)

// testApp is wide configuration for our book app tests
var testApp settings.Application

// TestMain is a base set up and configuration for tests
func TestMain(m *testing.M) {
	logger := &logger.Logger{Logger: logrus.New()}
	mdls := models.NewTestModels()
	jAuth := &auth.Auth{
		Issuer:        "example.com",
		Audience:      "example.com",
		Secret:        "highsecret",
		TokenExpiry:   5 * time.Minute,
		RefreshExpiry: 60 * time.Minute,
	}

	testApp.Logger = logger
	testApp.Models = mdls
	testApp.Auth = jAuth

	userApp = &testApp

	os.Exit(m.Run())
}

type myHandler struct{}

func (mh myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
