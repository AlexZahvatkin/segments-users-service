-- name: CreateSegment :one
INSERT INTO segments (name, created_at, updated_at) 
VALUES ($1, now(), now())
RETURNING *;

-- name: DeleteSegment :exec
DELETE FROM segments 
WHERE name = $1;