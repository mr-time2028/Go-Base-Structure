package server

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"go-base-structure/apps/book"
	"go-base-structure/apps/user"
	"go-base-structure/cmd/commands"
	"go-base-structure/cmd/routes"
	"go-base-structure/cmd/settings"
	"go-base-structure/database"
	"go-base-structure/models"
	"go-base-structure/pkg/env"
	"go-base-structure/pkg/logger"
	"net/http"
	"os"
)

// Serve start our servers
func Serve() (*settings.Application, error) {
	app := newApplication()

	// run command (if user want to run a command) else run http server
	command := flag.Bool("command", false, "run specific command")
	flag.Parse()

	if *command {
		commands.NewCommandsApp(app)
		runCommands(app)
	} else {
		err := HTTPServer(app)
		if err != nil {
			app.Logger.Fatal("failed to start the HTTP server: ", err)
		}
	}

	return app, nil
}

// newApplication start our program
func newApplication() *settings.Application {
	var app settings.Application

	// create a new logger
	logger := logger.NewLogger()
	app.Logger = logger

	// load .env file
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("cannot loading .env file")
	}

	// connect to the database
	logger.Info("connecting to the database...")
	DB, err := database.ConnectSQL()
	if err != nil {
		logger.Fatal("connecting to the database failed! ", err)
	}
	logger.Info("connected to the database successfully!")
	app.DB = DB

	// initial model
	mdls := models.NewModels()
	app.Models = mdls

	// auto migrations models
	err = mdls.AutoMigrateModels(DB.GormDB)
	if err != nil {
		logger.Fatal("auto migration failed! ", err)
	}
	logger.Info("auto migration was successful!")

	// pass some config to models package
	models.NewModelsApp(DB)

	// register your apps
	book.NewBookApp(&app)
	user.NewUserApp(&app)

	// initial config
	HTTPPort := env.GetEnvOrDefaultString("HTTP_PORT", "8000")
	Domain := env.GetEnvOrDefaultString("DOMAIN", "localhost")
	var cfg = &settings.Config{
		HTTPPort: HTTPPort,
		Domain:   Domain,
	}
	app.Config = cfg

	return &app
}

// HTTPServer starts a HTTP server
func HTTPServer(app *settings.Application) error {
	app.Logger.Info("the HTTP server is running on port ", app.Config.HTTPPort)
	err := http.ListenAndServe(fmt.Sprintf(":%s", app.Config.HTTPPort), routes.Routes())
	if err != nil {
		return err
	}

	return nil
}

// runCommands runs any command that user determine in CLI using -command parameter
func runCommands(app *settings.Application) {
	if len(os.Args) < 3 {
		app.Logger.Fatal("it seems you want to run a command. Usage: go run main.go -command <YOUR COMMAND NAME>")
	}

	commandName := os.Args[2]
	os.Args = os.Args[2:]

	command, ok := commands.Commands[commandName]
	if !ok {
		app.Logger.Fatal("invalid command.")
	}
	command.Function()
}
