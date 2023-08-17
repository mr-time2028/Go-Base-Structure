package database

import (
	"database/sql"
	"fmt"
	"go-base-structure/cmd/config"
	"go-base-structure/helpers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getDSN() string {
	dbName := helpers.GetEnvOrDefaultString("DB_NAME", "")
	dbUser := helpers.GetEnvOrDefaultString("DB_USER", "")
	dbPass := helpers.GetEnvOrDefaultString("DB_PASS", "")
	dbHost := helpers.GetEnvOrDefaultString("DB_HOST", "localhost")
	dbPort := helpers.GetEnvOrDefaultString("DB_PORT", "5432")
	dbSSL := helpers.GetEnvOrDefaultString("DB_SSL", "disable")
	dbZone := helpers.GetEnvOrDefaultString("DB_ZONE", "Asia/Tehran")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbHost, dbUser, dbPass, dbName, dbPort, dbSSL, dbZone)

	return dsn
}

func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

func openDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, err
}

func ConnectSQL() (*gorm.DB, *sql.DB) {
	config.AppConfig.InfoLog.Println("Connecting to database...")

	dsn := getDSN()

	db, err := openDB(dsn)
	if err != nil {
		config.AppConfig.ErrorLog.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		config.AppConfig.ErrorLog.Fatal(err)
	}

	config.AppConfig.InfoLog.Println("Testing database connection...")
	err = testDB(sqlDB)
	if err != nil {
		config.AppConfig.ErrorLog.Fatal(err)
	}

	config.AppConfig.InfoLog.Println("Connected to database successfully!")

	return db, sqlDB
}
