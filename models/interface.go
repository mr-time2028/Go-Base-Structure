package models

import (
	"go-base-structure/pkg/logging"
	"gorm.io/gorm"
	"reflect"
)

type BookInterface interface {
	GetAll() ([]*Book, error)
}

type UserInterface interface {
	GetOne(email string) (*User, error)
}

type ModelManager struct {
	Book BookInterface
	User UserInterface
}

func NewModels() *ModelManager {
	return &ModelManager{
		Book: &Book{},
		User: &User{},
	}
}

func NewTestModels() *ModelManager {
	return &ModelManager{
		Book: &TestBook{},
		User: &TestUser{},
	}
}

func AutoMigrateModels(logger *logging.Logger, gormDB *gorm.DB, models *ModelManager) {
	modelsValue := reflect.ValueOf(*models)

	for i := 0; i < modelsValue.NumField(); i++ {
		field := modelsValue.Field(i)
		if field.Kind() == reflect.Interface {
			model := field.Interface()
			if err := gormDB.AutoMigrate(model); err != nil {
				logger.ErrorLog.Fatal(err)
			}
		}
	}
}
