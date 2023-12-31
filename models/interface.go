package models

import (
	"gorm.io/gorm"
	"reflect"
)

type BookInterface interface {
	GetAll() ([]*Book, error)
	InsertOneBook(*Book) (int, error)
	InsertManyBooks([]*Book) (int64, []int, error)
}

type UserInterface interface {
	CheckIfExistsUser(email string) (bool, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	InsertOneUser(user *User) (int, error)
	InsertManyUsers(user []*User) (int64, []int, error)
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
