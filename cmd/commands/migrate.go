package commands

import (
	"go-base-structure/models"
	"go-base-structure/pkg/database"
	"log"
	"reflect"
)

// migrate is a simple command to create a superuser in database
func migrate() {
	db, err := database.ConnectSQL()
	if err != nil {
		log.Fatal("cannot connect to the database: ", err)
	}

	models := models.NewModels()
	modelsValue := reflect.ValueOf(*models)
	for i := 0; i < modelsValue.NumField(); i++ {
		field := modelsValue.Field(i)
		if field.Kind() == reflect.Interface {
			model := field.Interface()
			if err = db.GormDB.AutoMigrate(model); err != nil {
				commandsApp.Logger.Fatal("error while migration: ", err)
			}
		}
	}

	commandsApp.Logger.Println("migration was successful!")
}
