package models

import (
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

func (models *ModelManager) AutoMigrateModels(gormDB *gorm.DB) error {
	modelsValue := reflect.ValueOf(*models)

	for i := 0; i < modelsValue.NumField(); i++ {
		field := modelsValue.Field(i)
		if field.Kind() == reflect.Interface {
			model := field.Interface()
			if err := gormDB.AutoMigrate(model); err != nil {
				return err
			}
		}
	}

	return nil
}
