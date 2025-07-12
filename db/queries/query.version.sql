-- name: GetVersion :one
SELECT * FROM versions;

-- name: UpdateVersion :exec
UPDATE versions
SET
    updatedAt = CURRENT_TIMESTAMP,
    current = CAST(@current AS INTEGER)
WHERE
    id = ?;