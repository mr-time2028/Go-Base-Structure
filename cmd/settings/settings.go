package settings

import (
	"go-base-structure/cmd/server/config"
	"go-base-structure/models"
	"go-base-structure/pkg/auth"
	"go-base-structure/pkg/database"
	"go-base-structure/pkg/logger"
)

// Application is our wide configuration for the application
type Application struct {
	Config *config.Config
	Logger *logger.Logger
	DB     *database.DB
	Models *models.ModelManager
	Auth   *auth.Auth
}
