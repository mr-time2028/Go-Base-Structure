package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"go-base-structure/apps/book"
	"go-base-structure/apps/user"
	"go-base-structure/cmd/commands"
	"go-base-structure/cmd/config"
	"go-base-structure/cmd/routes"
	"go-base-structure/database"
	"go-base-structure/models"
	"go-base-structure/pkg/env"
	"go-base-structure/pkg/logging"
	"net/http"
	"os"
)

func main() {
	var app config.Application

	// create a new logging
	logger := logging.NewLogger()
	app.Logger = logger

	// load .env file
	err := godotenv.Load()
	if err != nil {
		logger.ErrorLog.Fatal("Error loading .env file")
	}

	// connect to the database
	gormDB, sqlDB := database.ConnectSQL(logger)
	var DB = &database.DB{
		GormDB: gormDB,
		SqlDB:  sqlDB,
	}
	defer sqlDB.Close()
	app.DB = DB

	// initial model and auto migration repository
	mdls := models.NewModels()
	app.Models = mdls

	models.AutoMigrateModels(logger, gormDB, mdls)

	models.NewModelsApp(logger, DB)

	// register your apps
	book.NewBookApp(&app)
	user.NewUserApp(&app)

	// run command (if user want to run a command) else run http server
	command := flag.Bool("command", false, "Run specific command")
	flag.Parse()

	if *command {
		commands.NewCommandsApp(&app)
		RunCommands(logger)
	} else {
		HTTPPort := env.GetEnvOrDefaultString("HTTP_PORT", "8000")
		Domain := env.GetEnvOrDefaultString("DOMAIN", "localhost")
		var cfg = &config.Config{
			HTTPPort: HTTPPort,
			Domain:   Domain,
		}
		app.Config = cfg
		err = serveHTTP(&app)
		if err != nil {
			logger.ErrorLog.Fatal("Failed to start the server: ", err)
		}
	}
}

// serveHTTP starts http server
func serveHTTP(app *config.Application) error {
	app.Logger.InfoLog.Println("The HTTP server is running on port", app.Config.HTTPPort)
	err := http.ListenAndServe(fmt.Sprintf(":%s", app.Config.HTTPPort), routes.Routes())
	if err != nil {
		return err
	}

	return nil
}

// RunCommands runs any command that user determine in CLI using -command parameter
func RunCommands(logger *logging.Logger) {
	if len(os.Args) < 3 {
		logger.ErrorLog.Println("It seems you want to run a command.")
		logger.ErrorLog.Println("Usage: go run main.go -command <YOUR COMMAND NAME>")
		return
	}

	commandName := os.Args[2]
	command, ok := commands.Commands[commandName]
	if !ok {
		logger.ErrorLog.Println("Invalid command.")
		return
	}
	command.Function()
}
