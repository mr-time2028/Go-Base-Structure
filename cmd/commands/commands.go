package commands

// Command is a struct for our custom commands
type Command struct {
	Name        string
	Description string
	Function    func()
}

// Commands contains custom commands information
var Commands = map[string]Command{
	"createSuperUser": {
		Name:        "createSuperUser",
		Description: "This command simply create a superuser in database with given username and password as CLI parameters.",
		Function:    createSuperUser,
	},
}
