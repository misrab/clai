package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Store handles all database operations
type Store struct {
	db *sqlx.DB
}

// NewStore creates a new Store instance and initializes the database
func NewStore() (*Store, error) {
	dbPath, err := GetDBPath()
	if err != nil {
		return nil, fmt.Errorf("get db path: %w", err)
	}

	fmt.Printf("Initializing database at: %s\n", dbPath)

	// Open database with pragmas for better performance and safety
	// _foreign_keys=on enables foreign key constraints
	// _journal_mode=WAL enables Write-Ahead Logging for better concurrency
	db, err := sqlx.Open("sqlite3", dbPath+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	// Set connection pool settings
	// SQLite only supports one writer at a time, so we limit connections
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	// Run migrations
	if err := runMigrations(db.DB); err != nil {
		db.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	return &Store{db: db}, nil
}

// Close closes the database connection
func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// DB returns the underlying sqlx database connection
func (s *Store) DB() *sqlx.DB {
	return s.db
}
