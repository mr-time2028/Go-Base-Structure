package models

import "go-base-structure/database"

// modelsApp is wide configuration instance belong to models package
var modelsApp *modelsConf

// modelsConf contains wide configuration settings we need in models package
type modelsConf struct {
	DB *database.DB
}

// NewModelsApp assign sent wide configuration instance to the modelsApp variable
func NewModelsApp(DB *database.DB) {
	modelsApp = &modelsConf{
		DB: DB,
	}
}
