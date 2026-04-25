-- name: ListTags :many
SELECT t.*
FROM tags t
    LEFT JOIN resource_tags rt ON t.id = rt.tag_id
WHERE t.user_id = $1
GROUP BY t.id
ORDER BY COUNT(rt.resource_id) DESC, t.name;

-- name: GetTag :one
SELECT * FROM tags
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: UpsertTag :one
INSERT INTO tags (user_id, name)
VALUES ($1, $2)
ON CONFLICT (user_id, name) DO NOTHING
RETURNING *;

-- name: GetTagByName :one
SELECT * FROM tags
WHERE user_id = $1 AND name = $2
LIMIT 1;

-- name: UpdateTag :one
UPDATE tags
SET name = $3
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteTag :one
DELETE FROM tags
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: GetTagsForResource :many
SELECT t.name
FROM tags t
    JOIN resource_tags rt ON t.id = rt.tag_id
WHERE rt.resource_id = $1;

-- name: GetTagsForResources :many
SELECT rt.resource_id, t.name
FROM resource_tags rt
    JOIN tags t ON t.id = rt.tag_id
WHERE rt.resource_id = ANY($1::uuid []);

-- name: DeleteResourceTags :exec
DELETE FROM resource_tags
WHERE resource_id = $1;

-- name: CreateResourceTag :exec
INSERT INTO resource_tags (resource_id, tag_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;
