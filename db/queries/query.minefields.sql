--
-- MineFields
--
-- name: GetMineField :one
SELECT
    *
FROM
    minefields
WHERE
    id = ?;

-- name: GetMineFieldByNum :one
SELECT
    *
FROM
    minefields
WHERE
    gameId = ?
    AND playerNum = ?
    AND num = ?;

-- name: GetMineFields :many
SELECT
    *
FROM
    minefields;

-- name: GetMineFieldsForGame :many
SELECT
    *
FROM
    minefields
WHERE
    gameId = ?;

-- name: GetMineFieldsForPlayer :many
SELECT
    *
FROM
    minefields
WHERE
    gameId = ?
    AND playerNum = ?;

-- name: CreateMineField :one
INSERT INTO
    minefields (
        createdAt,
        updatedAt,
        gameId,
        x,
        y,
        name,
        num,
        playerNum,
        tags,
        mineFieldType,
        numMines,
        detonate,
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
        ?,
        ?
    ) RETURNING *;

-- name: UpdateMineField :one
UPDATE minefields
SET
    updatedAt = CURRENT_TIMESTAMP,
    gameId = ?,
    x = ?,
    y = ?,
    name = ?,
    num = ?,
    playerNum = ?,
    tags = ?,
    mineFieldType = ?,
    numMines = ?,
    detonate = ?,
    spec = ?
WHERE
    id = ? RETURNING *;

-- name: DeleteMineField :exec
DELETE FROM minefields
WHERE
    id = ?;