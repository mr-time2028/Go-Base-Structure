package models

import (
	"go-base-structure/database"
	"go-base-structure/pkg/logger"
)

var modelsApp *models

type models struct {
	DB     *database.DB
	Logger *logger.Logger
}

func NewModelsApp(logger *logger.Logger, DB *database.DB) {
	modelsApp = &models{
		DB:     DB,
		Logger: logger,
	}
}
