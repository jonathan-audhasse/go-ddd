-- Migration: 001_create_users
-- Direction: UP
-- Description: Creates the users table with basic auth fields

BEGIN;

CREATE TABLE IF NOT EXISTS users (
    id          UUID primary key DEFAULT (gen_random_uuid()),
    email       TEXT NOT NULL UNIQUE,
    username    TEXT NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

COMMIT;