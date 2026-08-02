-- Migration: 002_add_posts
-- Direction: UP
-- Description: Adds a posts table linked to users
 
BEGIN;
 
CREATE TABLE IF NOT EXISTS posts (
    id         UUID        PRIMARY KEY DEFAULT (gen_random_uuid()),
    user_id    UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    title      TEXT        NOT NULL CHECK (char_length(title) BETWEEN 1 AND 255),
    body       TEXT        NOT NULL DEFAULT '',
    published  BOOLEAN     NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
 
CREATE INDEX IF NOT EXISTS idx_posts_user_id   ON posts (user_id);
CREATE INDEX IF NOT EXISTS idx_posts_published ON posts (published) WHERE published = true;
 
COMMIT;