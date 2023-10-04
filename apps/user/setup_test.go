package user

import (
	"github.com/sirupsen/logrus"
	"go-base-structure/cmd/settings"
	"go-base-structure/models"
	"go-base-structure/pkg/auth"
	"go-base-structure/pkg/database"
	"go-base-structure/pkg/logger"
	"os"
	"testing"
)

// testApp is wide configuration for our user app tests
var testApp settings.Application

// addDefaultData add some record(s) to a specific table in the database (if needed)
func addDefaultData() {
	users := []*models.User{
		{ID: 1, Email: "John@test.com", FirstName: "John", LastName: "Smith", Password: "JohnPass"},
		{ID: 2, Email: "Benjamin@test.com", FirstName: "Benjamin", LastName: "Montgomery", Password: "BenjaminPass"},
		{ID: 3, Email: "David@test.com", FirstName: "David", LastName: "Park", Password: "DavidPass"},
		{ID: 4, Email: "admin@test.com", FirstName: "FAdmin", LastName: "LAdmin", Password: "FAdminPass"},
	}

	_, _, err := testApp.Models.User.InsertManyUsers(users)
	if err != nil {
		testApp.Logger.Fatal("setUpTest error while isert record(s) to the database: ", err.Error())
	}
}

func setUpTest() {
	// init logger
	logr := &logger.Logger{Logger: logrus.New()}
	testApp.Logger = logr

	// init test DB
	DB, err := database.ConnectTestSQL()
	if err != nil {
		logr.Fatal("setUpTest error while connect to the database: ", err.Error())
	}
	testApp.DB = DB

	// drop table if exists
	err = DB.DropAllTables()
	if err != nil {
		logr.Fatal("setUpTest error while drop all tables: ", err.Error())
	}

	// init models and migration
	models.NewModelsApp(DB)
	mdls := models.NewModels()
	err = mdls.AutoMigrateModels(DB.GormDB)
	if err != nil {
		logr.Fatal("setUpTest error while migrate tables: ", err.Error())
	}
	testApp.Models = mdls

	// create some mock record(s) in the database
	addDefaultData()

	// init JWT
	jAuth := auth.NewTestJWTAuth()
	testApp.Auth = jAuth

	userApp = &testApp
}

func tearDownTest() {
	err := testApp.DB.DropAllTables()
	if err != nil {
		testApp.Logger.Fatal("tearDownTest error while drop tables: ", err.Error())
	}
}

func TestMain(m *testing.M) {
	setUpTest()
	exitCode := m.Run()
	tearDownTest()
	os.Exit(exitCode)
}
