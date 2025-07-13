--
-- Wormholes
--
-- name: GetWormhole :one
SELECT
    *
FROM
    wormholes
WHERE
    id = ?;

-- name: GetWormholeByNum :one
SELECT
    *
FROM
    wormholes
WHERE
    gameId = ?
    AND num = ?;

-- name: GetWormholes :many
SELECT
    *
FROM
    wormholes;

-- name: GetWormholesForGame :many
SELECT
    *
FROM
    wormholes
WHERE
    gameId = ?;

-- name: CreateWormhole :one
INSERT INTO
    wormholes (
        createdAt,
        updatedAt,
        gameId,
        x,
        y,
        name,
        num,
        tags,
        destinationNum,
        stability,
        yearsAtStability,
        spec
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
    ) RETURNING id, createdAt, updatedAt;

-- name: UpdateWormhole :one
UPDATE wormholes
SET
    updatedAt = CURRENT_TIMESTAMP,
    gameId = ?,
    x = ?,
    y = ?,
    name = ?,
    num = ?,
    tags = ?,
    destinationNum = ?,
    stability = ?,
    yearsAtStability = ?,
    spec = ?
WHERE
    id = ? RETURNING updatedAt;

-- name: DeleteWormhole :exec
DELETE FROM wormholes
WHERE
    id = ?;