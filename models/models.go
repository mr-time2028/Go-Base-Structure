package models

import (
	"go-base-structure/database"
	"go-base-structure/pkg/logging"
)

var modelsApp *models

type models struct {
	DB     *database.DB
	Logger *logging.Logger
}

func NewModelsApp(logger *logging.Logger, DB *database.DB) {
	modelsApp = &models{
		DB:     DB,
		Logger: logger,
	}
}
