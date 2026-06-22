-- name: CreateAPIToken :one
INSERT INTO
    api_tokens (
        user_id,
        name,
        token_prefix,
        token_hash,
        scope,
        expires_at
    )
VALUES
    (?, ?, ?, ?, ?, ?) RETURNING *;

-- name: GetAPITokenByHash :one
SELECT
    *
FROM
    api_tokens
WHERE
    token_hash = ?;

-- name: GetAPITokensForUser :many
SELECT
    *
FROM
    api_tokens
WHERE
    user_id = ?
ORDER BY
    created_at DESC;

-- name: TouchAPIToken :exec
UPDATE api_tokens
SET
    updated_at = CURRENT_TIMESTAMP,
    last_used_at = CURRENT_TIMESTAMP
WHERE
    id = ?;

-- name: RevokeAPIToken :exec
UPDATE api_tokens
SET
    updated_at = CURRENT_TIMESTAMP,
    revoked_at = CURRENT_TIMESTAMP
WHERE
    id = ?
    AND user_id = ?;
