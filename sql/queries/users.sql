-- name: CreateUser :one
INSERT INTO users (id, email, created_at, updated_at, hashed_password)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteAllUsers :exec
TRUNCATE TABLE users;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;