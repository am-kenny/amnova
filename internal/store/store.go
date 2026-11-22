package store

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const fileName = "amnova.db"

// dsnParams are applied to every connection.
var dsnParams = url.Values{
	"_journal_mode": {"WAL"},
	"_foreign_keys": {"1"},
	"_busy_timeout": {"5000"},
}.Encode()

type Store struct {
	db *sql.DB
}

// Open creates a directory if it does not exist and opens the database inside it.
func Open(ctx context.Context, dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return nil, fmt.Errorf("create data dir %q: %w", dir, err)
	}

	db, err := sql.Open("sqlite", filepath.Join(dir, fileName)+"?"+dsnParams)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// SQLite allows one writer at a time.
	db.SetMaxOpenConns(1)

	s := &Store{db: db}
	if err := s.Ping(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return s, nil
}

// Ping reports whether the database answers.
func (s *Store) Ping(ctx context.Context) error {
	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	return nil
}

// Close closes the database.
func (s *Store) Close() error {
	return s.db.Close()
}
