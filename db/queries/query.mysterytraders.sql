--
-- MysteryTraders
--
-- name: GetMysteryTrader :one
SELECT
    *
FROM
    mysterytraders
WHERE
    id = ?;

-- name: GetMysteryTraderByNum :one
SELECT
    *
FROM
    mysterytraders
WHERE
    gameId = ?
    AND num = ?;

-- name: GetMysteryTraders :many
SELECT
    *
FROM
    mysterytraders;

-- name: GetMysteryTradersForGame :many
SELECT
    *
FROM
    mysterytraders
WHERE
    gameId = ?;

-- name: CreateMysteryTrader :one
INSERT INTO
    mysterytraders (
        createdAt,
        updatedAt,
        gameId,
        x,
        y,
        name,
        num,
        tags,
        headingX,
        headingY,
        warpSpeed,
        requestedBoon,
        destinationX,
        destinationY,
        rewardType,
        playersRewarded,
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
        ?,
        ?,
        ?,
        ?,
        ?
    ) RETURNING *;

-- name: UpdateMysteryTrader :one
UPDATE mysterytraders
SET
    updatedAt = CURRENT_TIMESTAMP,
    gameId = ?,
    x = ?,
    y = ?,
    name = ?,
    num = ?,
    tags = ?,
    headingX = ?,
    headingY = ?,
    warpSpeed = ?,
    requestedBoon = ?,
    destinationX = ?,
    destinationY = ?,
    rewardType = ?,
    playersRewarded = ?,
    spec = ?
WHERE
    id = ? RETURNING *;

-- name: DeleteMysteryTrader :exec
DELETE FROM mysterytraders
WHERE
    id = ?;