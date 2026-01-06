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
SELECT * FROM users
ORDER BY created_at DESC;

-- name: UpdateUser :one
UPDATE users
SET 
  user_name = $2,
  updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: GrantAdminRole :exec
UPDATE users SET platform_role = 'admin', updated_at = NOW()
WHERE id = $1;
