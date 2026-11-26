package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/am-kenny/amnova/internal/domain"
	"github.com/am-kenny/amnova/internal/store/sqlc"
)

// CreateApplication stores a validated application, creating its company
// if needed. Returns a new application id.
func (s *Store) CreateApplication(ctx context.Context, in domain.ApplicationInput, now time.Time) (int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	q := sqlc.New(tx)
	stamp := formatTime(now)

	companyID, err := q.UpsertCompany(ctx, sqlc.UpsertCompanyParams{
		Name:      in.Company,
		CreatedAt: stamp,
		UpdatedAt: stamp,
	})
	if err != nil {
		return 0, fmt.Errorf("upsert company: %w", err)
	}

	app, err := q.InsertApplication(ctx, sqlc.InsertApplicationParams{
		CompanyID: companyID,
		Role:      in.Role,
		Status:    string(in.Status),
		AppliedAt: nullTime(in.AppliedAt),
		CreatedAt: stamp,
		UpdatedAt: stamp,
	})
	if err != nil {
		return 0, fmt.Errorf("insert application: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit: %w", err)
	}

	return app.ID, nil
}

// ListApplications returns all applications, newest first.
func (s *Store) ListApplications(ctx context.Context) ([]domain.Application, error) {
	rows, err := sqlc.New(s.db).ListApplications(ctx)
	if err != nil {
		return nil, fmt.Errorf("list applications: %w", err)
	}

	apps := make([]domain.Application, 0, len(rows))
	for _, r := range rows {
		createdAt, err := parseTime(r.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("application %d created at: %w", r.ID, err)
		}

		appliedAt, err := parseNullTime(r.AppliedAt)
		if err != nil {
			return nil, fmt.Errorf("application %d applied at: %w", r.ID, err)
		}

		apps = append(apps, domain.Application{
			ID:        r.ID,
			Company:   r.CompanyName,
			Role:      r.Role,
			Status:    domain.Status(r.Status),
			AppliedAt: appliedAt,
			CreatedAt: createdAt,
		})
	}

	return apps, nil
}

// formatTime returns t as RFC 3339 in UTC.
func formatTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// nullTime returns t as a nullable database time.
func nullTime(t *time.Time) sql.NullString {
	if t == nil {
		return sql.NullString{}
	}

	return sql.NullString{String: formatTime(*t), Valid: true}
}

// parseTime reads a time stored by formatTime.
func parseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

// parseNullTime reads a nullable time stored by nullTime.
func parseNullTime(s sql.NullString) (*time.Time, error) {
	if !s.Valid {
		return nil, nil
	}

	t, err := parseTime(s.String)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
