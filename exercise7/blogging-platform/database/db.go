package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	DB           *sql.DB
	user         = os.Getenv("MYSQL_USER")
	password     = os.Getenv("MYSQL_PASSWORD")
	rootPassword = os.Getenv("MYSQL_ROOT_PASSWORD")
	host         = os.Getenv("DB_HOST")
	port         = os.Getenv("DB_PORT")
	dbName       = os.Getenv("MYSQL_DATABASE")
	rootDsn      = fmt.Sprintf("%s:%s@tcp(%s:%s)/", "root", rootPassword, host, port)
)

const maxRetries = 5
const retryInterval = 5 * time.Second

func connectToDatabase() (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		slog.Info("Attempting to connect to the database...", "rootDsn", rootDsn)
		db, err = sql.Open("mysql", rootDsn)
		if err != nil {
			slog.Error("Failed to open connection with database", "err", err)
			return nil, err
		}
		if err == nil {
			pingErr := db.Ping()
			if pingErr == nil {
				slog.Info("Successfully connected to the database!")
				return db, nil
			} else {
				slog.Warn("Failed to connect to the database, retrying...", "retry", i+1, "error", pingErr)
				return nil, pingErr
			}
		}
		time.Sleep(retryInterval)
	}

	return nil, fmt.Errorf("could not connect to the database after %d attempts: %w", maxRetries, err)
}

func DatabaseInit() {

	slog.Info("Initializing database...")
	var err error
	DB, err = connectToDatabase()
	if err != nil {
		slog.Error("Exiting application due to database connection failure", "error", err)
		os.Exit(1)
	}

	slog.Info("verifying database connection...")
	if err := DB.Ping(); err != nil {
		slog.Error("failed to connect to MySQL: ", "error", err)
	}

	slog.Info("creating database")
	if err := createDatabase(dbName); err != nil {
		slog.Error("Failed to create database", "dbName", dbName, "error", err)
	}

	slog.Info("reconnecting to the database...")
	DB.Close()
	DB, err = sql.Open("mysql", rootDsn+dbName)
	if err != nil {
		slog.Error("Failed to connect to database", "dbName", dbName, "error", err)
	}
	slog.Info("database reconnected")

	slog.Info("verifying database connection...")
	if err := DB.Ping(); err != nil {
		slog.Error("Failed to connect to database", "dbName", dbName, "error", err)
	}

	slog.Info("creating table..")
	if err := createPostsTable(); err != nil {
		slog.Error("Failed to create posts table: ", "error", err)
	}
	slog.Info("table created")
}

func createDatabase(dbName string) error {
	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName)
	_, err := DB.Exec(query)
	return err
}

func createPostsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS posts (
		id CHAR(36) PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
	    category VARCHAR(255) NOT NULL,
	    tags JSON NOT NULL,
	    createdat VARCHAR(255) NOT NULL,
	    updatedat VARCHAR(255) NOT NULL
	)`
	_, err := DB.Exec(query)
	return err
}
