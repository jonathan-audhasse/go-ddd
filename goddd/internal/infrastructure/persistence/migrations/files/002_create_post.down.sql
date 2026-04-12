-- Migration: 002_add_posts
-- Direction: DOWN
-- Description: Reverses the posts table creation
 
BEGIN;
 
DROP TABLE IF EXISTS posts;
 
COMMIT;
 