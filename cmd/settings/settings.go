package settings

import (
	"go-base-structure/database"
	"go-base-structure/models"
	"go-base-structure/pkg/auth"
	"go-base-structure/pkg/logger"
)

// Application is our wide configuration for the application
type Application struct {
	Config *Config
	Logger *logger.Logger
	DB     *database.DB
	Models *models.ModelManager
	Auth   *auth.Auth
}

// Config is our configuration about the server
type Config struct {
	HTTPPort string
	Domain   string
}
