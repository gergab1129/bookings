package driver

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// holds the database connection pool
type DB struct {
	SQL *sqlx.DB
}

var dbConn = &DB{}

const (
	maxOpenDBConn = 10
	maxIdleDbConn = 5
	maxDBLifetime = 5 * time.Minute
)

// ConnectSQL create database pool to postgres
func ConnectSQL(dsn string) (*DB, error) {
	db, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(maxOpenDBConn)
	db.SetMaxIdleConns(maxIdleDbConn)
	db.SetConnMaxLifetime(maxDBLifetime)

	dbConn.SQL = db

	err = testDB(db)
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

// testDB tries to ping the database
func testDB(db *sqlx.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}

	return nil
}

// NewDatabase creates a new database abstraction for the application
func NewDatabase(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
