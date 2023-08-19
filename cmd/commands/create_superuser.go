package commands

import (
	"go-base-structure/models"
)

// createSuperUser is a command to create a superuser in database
func createSuperUser() {
	user := models.User{
		Email:    "admin@test.com",
		Password: "testPass",
		// consider a is_superuser field and set it to true value here...
	}
	result := commandsApp.DB.GormDB.Create(&user)
	commandsApp.Logger.InfoLog.Printf("superuser created successfully!, row affected: %d", result.RowsAffected)
}
