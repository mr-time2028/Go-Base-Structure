package models

import (
	"database/sql"
	"go-base-structure/cmd/config"
	"gorm.io/gorm"
	"reflect"
)

var GormDB *gorm.DB
var SqlDB *sql.DB
var Models ModelManager

type ModelManager struct {
	Book BookInterface
	User UserInterface
}

func NewDB(GD *gorm.DB, SD *sql.DB) {
	GormDB = GD
	SqlDB = SD

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
			if err := GormDB.AutoMigrate(model); err != nil {
				config.AppConfig.ErrorLog.Fatal(err)
			}
		}
	}
}
