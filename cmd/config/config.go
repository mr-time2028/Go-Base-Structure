package config

import (
	"go-base-structure/database"
	"go-base-structure/models"
	"go-base-structure/pkg/logging"
)

type Application struct {
	Config *Config
	Logger *logging.Logger
	DB     *database.DB
	Models *models.ModelManager
}

// Config is our wide configuration for the application
type Config struct {
	HTTPPort string
	Domain   string
}
