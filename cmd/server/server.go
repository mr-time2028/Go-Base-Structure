package server

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"go-base-structure/apps/book"
	"go-base-structure/apps/user"
	"go-base-structure/cmd/commands"
	"go-base-structure/cmd/routes"
	"go-base-structure/cmd/server/config"
	"go-base-structure/cmd/settings"
	"go-base-structure/models"
	"go-base-structure/pkg/auth"
	"go-base-structure/pkg/database"
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
	logr := logger.NewLogger()
	app.Logger = logr

	// load .env file
	err := godotenv.Load()
	if err != nil {
		logr.Fatal("cannot loading .env file")
	}

	// initial config
	cfg := config.NewConfig()
	app.Config = cfg

	// JWT settings
	jAuth := auth.NewJWTAuth()
	app.Auth = jAuth

	// connect to the database
	logr.Info("connecting to the database...")
	DB, err := database.ConnectSQL()
	if err != nil {
		logr.Fatal("connecting to the database failed! ", err)
	}
	logr.Info("connected to the database successfully!")
	app.DB = DB

	// initial model
	mdls := models.NewModels()
	app.Models = mdls

	// auto migrations models
	err = mdls.AutoMigrateModels(DB.GormDB)
	if err != nil {
		logr.Fatal("auto migration failed! ", err)
	}
	logr.Info("auto migration was successful!")

	// pass some config to models package
	models.NewModelsApp(DB)

	// register your apps
	book.NewBookApp(&app)
	user.NewUserApp(&app)

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
