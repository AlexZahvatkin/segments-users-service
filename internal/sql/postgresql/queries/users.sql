-- name: CreateUser :one
INSERT INTO users(name, created_at, updated_at) 
VALUES ($1, now(), now())
RETURNING *;

-- name: DeleteUser :exec 
DELETE FROM users 
WHERE id = $1;

-- name: GetAllUsers :many
SELECT id
FROM users;