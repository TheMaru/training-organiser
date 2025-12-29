-- name: CreateUser :one
INSERT INTO users (password_hash, user_name)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserByUserName :one
SELECT * FROM users
WHERE user_name = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUsers :many
SELECT * FROM users;
