package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	*sql.DB
}

func New(connectionString string) (*DB, error) {
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

func MustNew(connectionString string) *DB {
	db, err := New(connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal database error: %v\n", err)
		os.Exit(1)
	}
	return db
}
