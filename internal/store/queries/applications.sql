-- name: InsertApplication :one
INSERT INTO applications (company_id, role, status, applied_at, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: ListApplications :many
SELECT a.id, c.name AS company_name, a.role, a.status, a.applied_at, a.created_at, a.updated_at
FROM applications a
JOIN companies c ON c.id = a.company_id
ORDER BY a.created_at DESC, a.id DESC;
