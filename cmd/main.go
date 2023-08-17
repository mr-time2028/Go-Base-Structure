package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"go-base-structure/cmd/commands"
	"go-base-structure/cmd/config"
	"go-base-structure/cmd/routes"
	"go-base-structure/database"
	"go-base-structure/models"
	"net/http"
	"os"
)

func main() {
	// create a new logger
	config.NewLogger()

	// load .env file
	err := godotenv.Load()
	if err != nil {
		config.AppConfig.ErrorLog.Fatal("Error loading .env file")
	}

	// connect to the database
	gormDB, sqlDB := database.ConnectSQL()

	// initial model and auto migration models
	models.NewDB(gormDB, sqlDB)
	models.AutoMigrateModels()

	// run command (if user want to run a command) else run http server
	command := flag.Bool("command", false, "Run specific command")
	flag.Parse()

	if *command {
		RunCommands()
	} else {
		err = serveHTTP()
		if err != nil {
			config.AppConfig.ErrorLog.Fatal("Failed to start the server: ", err)
		}
	}
}

// serveHTTP starts http server
func serveHTTP() error {
	config.AppConfig.HTTPPort = 8000
	config.AppConfig.Domain = "localhost"

	config.AppConfig.InfoLog.Printf("The HTTP server is running on port %d", config.AppConfig.HTTPPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.AppConfig.HTTPPort), routes.Routes())
	if err != nil {
		return err
	}

	return nil
}

// RunCommands runs any command that user determine in CLI using -command parameter
func RunCommands() {
	if len(os.Args) < 3 {
		fmt.Println("It seems you want to run a command.")
		fmt.Println("Usage: go run main.go -command <YOUR COMMAND NAME>")
		return
	}

	commandName := os.Args[2]
	command, ok := commands.Commands[commandName]
	if !ok {
		fmt.Println("Invalid command.")
		return
	}
	command.Function()
}
