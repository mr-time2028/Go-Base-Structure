package commands

import (
	"go-base-structure/cmd/settings"
)

// commandsApp is wide configuration instance belong to commands package
var commandsApp *settings.Application

// Command is a struct for our custom commands
type Command struct {
	Description string
	Function    func()
}

// NewCommandsApp assign sent wide configuration instance to the commandsApp variable
func NewCommandsApp(app *settings.Application) {
	commandsApp = app
}

// Commands contains custom commands information
var Commands = map[string]Command{
	"createSuperUser": {
		Description: "this command simply create a superuser in database with given username and password as CLI parameters.",
		Function:    createSuperUser,
	},
}
