-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;


-- name: UpdateRefreshToken :exec
UPDATE refresh_tokens
SET updated_at = $2, revoked_at = $3
WHERE token = $1;