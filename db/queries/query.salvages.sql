--
-- Salvages
--
-- name: GetSalvage :one
SELECT
    *
FROM
    salvages
WHERE
    id = ?;

-- name: GetSalvageByNum :one
SELECT
    *
FROM
    salvages
WHERE
    gameId = ?
    AND num = ?;

-- name: GetSalvages :many
SELECT
    *
FROM
    salvages;

-- name: GetSalvagesForGame :many
SELECT
    *
FROM
    salvages
WHERE
    gameId = ?;

-- name: GetSalvagesForPlayer :many
SELECT
    *
FROM
    salvages
WHERE
    gameId = ?
    AND playerNum = ?;

-- name: CreateSalvage :one
INSERT INTO
    salvages (
        createdAt,
        updatedAt,
        gameId,
        x,
        y,
        name,
        num,
        playerNum,
        tags,
        ironium,
        boranium,
        germanium
    )
VALUES
    (
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?
    ) RETURNING *;

-- name: UpdateSalvage :one
UPDATE salvages
SET
    updatedAt = CURRENT_TIMESTAMP,
    gameId = ?,
    x = ?,
    y = ?,
    name = ?,
    num = ?,
    playerNum = ?,
    tags = ?,
    ironium = ?,
    boranium = ?,
    germanium = ?
WHERE
    id = ? RETURNING *;

-- name: DeleteSalvage :exec
DELETE FROM salvages
WHERE
    id = ?;