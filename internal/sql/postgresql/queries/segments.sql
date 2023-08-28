-- name: AddSegment :one
INSERT INTO segments (name, created_at, updated_at, description) 
VALUES ($1, now(), now(), $2)
RETURNING *;

-- name: DeleteSegment :exec
DELETE FROM segments 
WHERE name = $1;