-- name: CreateCollection :one
INSERT INTO collections (user_id, title, description, public)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetCollection :one
SELECT * FROM collections
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: ListCollections :many
SELECT * FROM collections
WHERE user_id = $1
ORDER BY created_at;

-- name: UpdateCollection :one
UPDATE collections
SET title = $3, description = $4, public = $5
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteCollection :one
DELETE FROM collections
WHERE id = $1 AND user_id = $2
RETURNING *;
