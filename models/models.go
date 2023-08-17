package models

import (
	"go-base-structure/cmd/config"
	"go-base-structure/database"
	"reflect"
)

var Models ModelManager

type ModelManager struct {
	Book BookInterface
	User UserInterface
}

func NewModels() {
	Models = ModelManager{
		Book: &Book{},
		User: &User{},
	}
}

func AutoMigrateModels() {
	modelsValue := reflect.ValueOf(Models)

	for i := 0; i < modelsValue.NumField(); i++ {
		field := modelsValue.Field(i)
		if field.Kind() == reflect.Interface {
			model := field.Interface()
			if err := database.GormDB.AutoMigrate(model); err != nil {
				config.AppConfig.ErrorLog.Fatal(err)
			}
		}
	}
}
