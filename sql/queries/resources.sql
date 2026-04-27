-- name: GetResource :one
SELECT r.*, c.title AS collection_title
FROM resources r
    LEFT JOIN collections c ON r.collection_id = c.id
WHERE r.id = $1 AND r.user_id = $2
LIMIT 1;

-- name: ListResources :many
SELECT r.*, c.title AS collection_title
FROM resources r
    LEFT JOIN collections c ON r.collection_id = c.id
WHERE r.user_id = $1
ORDER BY r.created_at;

-- name: CreateResource :one
INSERT INTO resources (
    user_id, title, description, url, collection_id
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateResource :one
UPDATE resources
SET
    title = $3,
    description = $4,
    url = $5,
    collection_id = $6
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteResource :one
DELETE FROM resources
WHERE id = $1 AND user_id = $2
RETURNING *;
