package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func NewDB(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established")
	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}

// InitializeTables creates all necessary tables based on the ERD
func (db *DB) InitializeTables() error {
	// Enable foreign key support
	_, err := db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Create tables will be handled by migrations
	log.Println("Database initialized successfully")
	return nil
}