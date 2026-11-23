-- +goose Up
CREATE TABLE companies (
  id          INTEGER PRIMARY KEY,
  name        TEXT NOT NULL UNIQUE COLLATE NOCASE,
  created_at  TEXT NOT NULL,
  updated_at  TEXT NOT NULL
);

CREATE TABLE applications (
  id                  INTEGER PRIMARY KEY,
  company_id          INTEGER NOT NULL REFERENCES companies(id),
  role                TEXT NOT NULL,
  status              TEXT NOT NULL CHECK (status IN ('saved','applied','interviewing','offer','accepted','rejected','no_response')),
  applied_at          TEXT,
  created_at          TEXT NOT NULL,
  updated_at          TEXT NOT NULL
);
CREATE INDEX applications_company_id ON applications(company_id);
