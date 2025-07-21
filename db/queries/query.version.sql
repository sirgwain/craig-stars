-- name: GetVersion :one
SELECT * FROM versions;

-- name: UpdateVersion :exec
UPDATE versions
SET
    updated_at = CURRENT_TIMESTAMP,
    current = ?
WHERE
    id = ?;