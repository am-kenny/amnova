package store

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// Migrate applies every pending migration in order.
func (s *Store) Migrate(ctx context.Context) error {
	dir, err := fs.Sub(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}

	// Do not close the provider: that closes the store's database.
	p, err := goose.NewProvider(goose.DialectSQLite3, s.db, dir,
		goose.WithVerbose(true),
		goose.WithSlog(slog.Default()),
	)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}

	if _, err := p.Up(ctx); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}

	return nil
}
