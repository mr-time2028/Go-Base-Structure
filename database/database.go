package database

import (
	"database/sql"
	"fmt"
	"go-base-structure/pkg/env"
	"go-base-structure/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

const maxOpenDBConn = 10
const maxIdleDBConn = 5
const maxDBLifetime = 5 * time.Minute

type DB struct {
	GormDB *gorm.DB
	SqlDB  *sql.DB
}

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

func ConnectSQL(logger *logger.Logger) (*gorm.DB, *sql.DB) {
	logger.InfoLog.Println("Connecting to database...")

	dsn := getDSN()

	db, err := openDB(dsn)
	if err != nil {
		logger.ErrorLog.Fatal(err)
	}

	sdb, err := db.DB()
	if err != nil {
		logger.ErrorLog.Fatal(err)
	}
	sdb.SetMaxOpenConns(maxOpenDBConn)
	sdb.SetMaxIdleConns(maxIdleDBConn)
	sdb.SetConnMaxLifetime(maxDBLifetime)

	logger.InfoLog.Println("Testing database connection...")
	err = testDB(sdb)
	if err != nil {
		logger.ErrorLog.Fatal(err)
	}

	logger.InfoLog.Println("Connected to database successfully!")

	return db, sdb
}
