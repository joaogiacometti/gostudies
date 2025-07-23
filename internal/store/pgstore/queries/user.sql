-- name: CreateUser :one
INSERT INTO users (user_name, email, password_hash)
VALUES ($1, $2, $3)
RETURNING id;

-- name: GetUserByID :one
SELECT id, user_name, email, created_at
FROM users
WHERE id = $1;