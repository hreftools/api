-- name: CreateSession :one
INSERT INTO sessions (user_id, hash, description, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetSessionByHash :one
SELECT * FROM sessions
WHERE hash = $1
LIMIT 1;

-- name: UpdateSessionExpiresAt :one
UPDATE sessions
SET expires_at = $2
WHERE id = $1
RETURNING *;

-- name: DeleteSessionByHash :exec
DELETE FROM sessions
WHERE hash = $1;

-- name: DeleteSessionsByUserID :exec
DELETE FROM sessions
WHERE user_id = $1;
