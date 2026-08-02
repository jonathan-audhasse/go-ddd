-- Migration: 001_create_users
-- Direction: DOWN
-- Description: Reverses the users table creation

BEGIN;

DROP TABLE IF EXISTS users;
 
COMMIT;