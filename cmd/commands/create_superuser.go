package commands

import (
	"flag"
	"go-base-structure/models"
)

// createSuperUser is a simple command to create a superuser in database
func createSuperUser() {
	email := flag.String("email", "", "user email")
	password := flag.String("password", "", "user password")
	firstName := flag.String("firstName", "", "user first name")
	lastName := flag.String("lastName", "", "user last name")
	flag.Parse()

	// validations goes here...
	if *email == "" {
		commandsApp.Logger.Fatal("email is required.")
	}
	if *password == "" {
		commandsApp.Logger.Fatal("password is required.")
	}

	user := models.User{
		Email:     *email,
		Password:  *password,
		FirstName: *firstName,
		LastName:  *lastName,
		// consider an is_superuser field and set it to true here...
	}
	result := commandsApp.DB.GormDB.Create(&user)
	commandsApp.Logger.Info("superuser created successfully!, row affected: ", result.RowsAffected)
}
