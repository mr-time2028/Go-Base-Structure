package models

import (
	"github.com/sirupsen/logrus"
	"go-base-structure/pkg/database"
	"go-base-structure/pkg/logger"
	"os"
	"reflect"
	"testing"
)

var testModelsApp modelsConf
var logr = &logger.Logger{Logger: logrus.New()}

func interfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("interfaceSlice() called with a non-slice type")
	}

	result := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		result[i] = s.Index(i).Interface()
	}

	return result
}

// addDefaultData adds default data to the database using a batch size
func addDefaultData() error {
	var defaultBooks = []*Book{
		{Name: "Harry Potter"},
		{Name: "Pride and Prejudice"},
	}

	if err := testModelsApp.DB.GormDB.CreateInBatches(defaultBooks, len(defaultBooks)).Error; err != nil {
		logr.Fatal("addDefaultData error while adding default book(s) data to the database: ", err.Error())
	}

	var defaultUsers = []*User{
		{Email: "David@test.com", Password: "DavidPass"},
		{Email: "John@test.com", Password: "JohnPass"},
	}

	if err := testModelsApp.DB.GormDB.CreateInBatches(defaultUsers, len(defaultUsers)).Error; err != nil {
		logr.Fatal("addDefaultData error while adding default user(s) data to the database: ", err.Error())
	}

	return nil
}

// resetTestDB reset the database (drop all tables and migrate models again amd add default data to the database) : use for tests
func resetTestDB() error {
	// drop all tables (if exits any from previous tests)
	err := testModelsApp.DB.DropAllTables()
	if err != nil {
		return err
	}

	// migration models
	mdls := NewModels()
	err = mdls.AutoMigrateModels(testModelsApp.DB.GormDB)
	if err != nil {
		return err
	}

	// add default data to the database
	err = addDefaultData()
	if err != nil {
		return err
	}

	return nil
}

func setUpTest() {
	// connect to the database
	testDB, err := database.ConnectTestSQL()
	if err != nil {
		logr.Fatal("setUpTest error while connect to the database: ", err.Error())
	}
	testModelsApp.DB = testDB

	err = resetTestDB()
	if err != nil {
		logr.Fatal("setUpTest error while reset the database: ", err.Error())
	}

	modelsApp = &testModelsApp
}

func tearDownTest() {
	err := testModelsApp.DB.DropAllTables()
	if err != nil {
		logr.Fatal("setUpTest error while drop all tables: ", err.Error())
	}
}

func TestMain(m *testing.M) {
	setUpTest()
	exitCode := m.Run()
	tearDownTest()
	os.Exit(exitCode)
}
