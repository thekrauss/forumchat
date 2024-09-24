package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Store interface {
	OpenDatabase() (*sql.DB, error)
	CloseDatabase(db *sql.DB) error
}

type DBStore struct{}

func (s *DBStore) OpenDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "database/forum.DB")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	log.Println("Ping to database")

	return db, nil
}

func (s *DBStore) CloseDatabase(db *sql.DB) error {
	if err := db.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}
	return nil
}
