-- name: AddUser :one
INSERT INTO users(name, created_at, updated_at) 
VALUES ($1, now(), now())
RETURNING *;

-- name: DeleteUser :exec 
DELETE FROM users 
WHERE id = $1;

-- name: GetAllUsersId :many
SELECT id
FROM users;

-- name: GetUserById :one 
SELECT *
FROM users
WHERE id = $1;