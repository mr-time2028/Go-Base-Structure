package database

import (
	"database/sql"
	"fmt"
	"go-base-structure/pkg/env"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

// DB contains sql and nosql dbs
type DB struct {
	GormDB *gorm.DB
	SqlDB  *sql.DB
}

const maxOpenDBConn = 10
const maxIdleDBConn = 5
const maxDBLifetime = 5 * time.Minute

// getDSN return dsn string for connection to the database
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

// testDB ping to the database to ensure database is open
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

// openDB open the database with dsn from getDSN
func openDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	return db, err
}

// ConnectSQL get dsn and open db and return DB instance
func ConnectSQL() (*DB, error) {
	dsn := getDSN()

	db, err := openDB(dsn)
	if err != nil {
		return nil, err
	}

	sdb, err := db.DB()
	if err != nil {
		return nil, err
	}
	sdb.SetMaxOpenConns(maxOpenDBConn)
	sdb.SetMaxIdleConns(maxIdleDBConn)
	sdb.SetConnMaxLifetime(maxDBLifetime)

	err = testDB(sdb)
	if err != nil {
		return nil, err
	}

	return &DB{
		GormDB: db,
		SqlDB:  sdb,
	}, nil
}

// ConnectTestSQL used when you want set up the database for the tests
func ConnectTestSQL() (*DB, error) {
	dsn := `host=localhost user=postgres password=password dbname=testDB port=5432 sslmode=disable TimeZone=Asia/Tehran`
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	sdb, _ := db.DB()

	return &DB{
		GormDB: db,
		SqlDB:  sdb,
	}, nil
}

// GetAllTables gather all tables name that exists in the database
func (db *DB) GetAllTables() ([]string, error) {
	tables, err := db.GormDB.Migrator().GetTables()
	if err != nil {
		return nil, err
	}

	return tables, nil
}

// DropAllTables drop all tables in the database
func (db *DB) DropAllTables() error {
	tables, err := db.GetAllTables()
	if err != nil {
		return err
	}

	for _, table := range tables {
		if err = db.GormDB.Migrator().DropTable(table); err != nil {
			return err
		}
	}

	return nil
}
