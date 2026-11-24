-- name: UpsertCompany :one
-- Returns the id of the company with this name, creating it if needed.
INSERT INTO companies (name, created_at, updated_at)
VALUES (?, ?, ?)
ON CONFLICT (name) DO UPDATE SET name = name
RETURNING id;
