package database

import (
	"database/sql"
	"fmt"
	"go-base-structure/cmd/config"
	"go-base-structure/pkg/env"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var GormDB *gorm.DB
var SqlDB *sql.DB

func getDSN() string {
	dbName := env.GetEnvOrDefaultString("DB_NAME", "")
	dbUser := env.GetEnvOrDefaultString("DB_USER", "")
	dbPass := env.GetEnvOrDefaultString("DB_PASS", "")
	dbHost := env.GetEnvOrDefaultString("DB_HOST", "localhost")
	dbPort := env.GetEnvOrDefaultString("DB_PORT", "5432")
	dbSSL := env.GetEnvOrDefaultString("DB_SSL", "disable")
	dbZone := env.GetEnvOrDefaultString("DB_ZONE", "Asia/Tehran")

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

func ConnectSQL() {
	config.AppConfig.InfoLog.Println("Connecting to database...")

	dsn := getDSN()

	db, err := openDB(dsn)
	if err != nil {
		config.AppConfig.ErrorLog.Fatal(err)
	}

	sdb, err := db.DB()
	if err != nil {
		config.AppConfig.ErrorLog.Fatal(err)
	}

	config.AppConfig.InfoLog.Println("Testing database connection...")
	err = testDB(sdb)
	if err != nil {
		config.AppConfig.ErrorLog.Fatal(err)
	}

	config.AppConfig.InfoLog.Println("Connected to database successfully!")

	GormDB = db
	SqlDB = sdb
}
