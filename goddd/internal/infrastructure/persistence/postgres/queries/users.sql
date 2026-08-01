-- name: GetUserByID :one
SELECT id, email, username, created_at, updated_at
FROM users
WHERE id = $1;
 
-- name: ListUsersPaged :many
SELECT id, email, username, created_at, updated_at
FROM users
WHERE (
    $1::uuid = '00000000-0000-0000-0000-000000000000'::uuid
    OR created_at < (SELECT created_at FROM users WHERE id = $1::uuid)
    OR (created_at = (SELECT created_at FROM users WHERE id = $1::uuid) AND id > $1::uuid)
)
ORDER BY created_at DESC, id ASC
LIMIT $2;
 
-- name: CountUsers :one
SELECT COUNT(*) FROM users;
 
-- name: CreateUser :one
INSERT INTO users (email, username) 
VALUES ($1, $2)
RETURNING *;

-- name: BulkCreateUsers :many
INSERT INTO users (email, username)
SELECT unnest(@emails::text[]), unnest(@usernames::text[])
RETURNING *;


-- name: UpdateUser :one
UPDATE users
SET
    email         = $1,
    username      = $2,
    updated_at    = now()
WHERE id = $3
RETURNING *;
 
-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;