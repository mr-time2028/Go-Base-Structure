package book

import (
	"github.com/sirupsen/logrus"
	"go-base-structure/cmd/settings"
	"go-base-structure/models"
	"go-base-structure/pkg/database"
	"go-base-structure/pkg/logger"
	"os"
	"testing"
)

// testBookApp is wide configuration for our book app tests
var testBookApp settings.Application

// addDefaultData add some record(s) to a specific table in the database (if needed)
func addDefaultData() {
	books := []*models.Book{
		{ID: 1, Name: "The Lord of the Rings"},
		{ID: 2, Name: "A Song of Ice and Fire"},
		{ID: 3, Name: "Harry Potter"},
		{ID: 4, Name: "Pride and Prejudice"},
	}

	_, _, err := testBookApp.Models.Book.InsertManyBooks(books)
	if err != nil {
		testBookApp.Logger.Fatal("setUpTest error while insert record(s) to the database: ", err.Error())
	}
}

// tearDownTest tear down tests
func tearDownTest() {
	err := testBookApp.DB.DropAllTables()
	if err != nil {
		testBookApp.Logger.Fatal("tearDownTest error while drop tables: ", err.Error())
	}
}

// setUpTest set up tests
func setUpTest() {
	logr := &logger.Logger{Logger: logrus.New()}
	testBookApp.Logger = logr

	// connect to test DB
	DB, err := database.ConnectTestSQL()
	if err != nil {
		logr.Fatal("setUpTest error while connect to the database: ", err.Error())
	}
	testBookApp.DB = DB

	// drop all tables if exists from previous tests
	err = DB.DropAllTables()

	// init models and migration
	models.NewModelsApp(DB)
	mdls := models.NewModels()
	err = mdls.AutoMigrateModels(DB.GormDB)
	if err != nil {
		logr.Fatal("setUpTest error while migrate tables: ", err.Error())
	}
	testBookApp.Models = mdls

	// create some mock record(s) in the database (if needed)
	addDefaultData()

	// init bookApp to use in tests
	bookApp = &testBookApp
}

// TestMain is the configuration for the tests (setUp and tearDown)
func TestMain(m *testing.M) {
	setUpTest()
	exitCode := m.Run()
	tearDownTest()
	os.Exit(exitCode)
}
