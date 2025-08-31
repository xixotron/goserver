-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES($1, NOW(), NOW(), $2, NOW() + '60 days', NULL)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT user_id FROM refresh_tokens
WHERE revoked_at IS NULL
  AND token = $1
  AND expires_at > NOW();

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens SET
  updated_at = NOW(),
  revoked_at = NOW()
WHERE token = $1;

