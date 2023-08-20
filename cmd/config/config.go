package config

import (
	"go-base-structure/database"
	"go-base-structure/models"
	"go-base-structure/pkg/logger"
)

type Application struct {
	Config *Config
	Logger *logger.Logger
	DB     *database.DB
	Models *models.ModelManager
}

// Config is our wide configuration for the application
type Config struct {
	HTTPPort string
	Domain   string
}
