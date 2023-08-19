package commands

import "go-base-structure/cmd/config"

var commandsApp *config.Application

// Command is a struct for our custom commands
type command struct {
	Name        string
	Description string
	Function    func()
}

func NewCommandsApp(app *config.Application) {
	commandsApp = app
}

// Commands contains custom commands information
var Commands = map[string]command{
	"createSuperUser": {
		Name:        "createSuperUser",
		Description: "This command simply create a superuser in database with given username and password as CLI parameters.",
		Function:    createSuperUser,
	},
}
