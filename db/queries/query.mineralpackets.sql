--
-- MineralPackets
--
-- name: GetMineralPacket :one
SELECT
    *
FROM
    mineralpackets
WHERE
    id = ?;

-- name: GetMineralPacketByNum :one
SELECT
    *
FROM
    mineralpackets
WHERE
    gameId = ?
    AND playerNum = ?
    AND num = ?;

-- name: GetMineralPackets :many
SELECT
    *
FROM
    mineralpackets;

-- name: GetMineralPacketsForGame :many
SELECT
    *
FROM
    mineralpackets
WHERE
    gameId = ?;

-- name: GetMineralPacketsForPlayer :many
SELECT
    *
FROM
    mineralpackets
WHERE
    gameId = ?
    AND playerNum = ?;

-- name: CreateMineralPacket :one
INSERT INTO
    mineralpackets (
        createdAt,
        updatedAt,
        gameId,
        x,
        y,
        name,
        num,
        playerNum,
        tags,
        targetPlanetNum,
        ironium,
        boranium,
        germanium,
        safeWarpSpeed,
        warpSpeed,
        scanRange,
        scanRangePen,
        headingX,
        headingY
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
        ?,
        ?,
        ?
    ) RETURNING id, createdAt, updatedAt;

-- name: UpdateMineralPacket :one
UPDATE mineralpackets
SET
    updatedAt = CURRENT_TIMESTAMP,
    gameId = ?,
    x = ?,
    y = ?,
    name = ?,
    num = ?,
    playerNum = ?,
    tags = ?,
    targetPlanetNum = ?,
    ironium = ?,
    boranium = ?,
    germanium = ?,
    safeWarpSpeed = ?,
    warpSpeed = ?,
    scanRange = ?,
    scanRangePen = ?,
    headingX = ?,
    headingY = ?
WHERE
    id = ? RETURNING updatedAt;

-- name: DeleteMineralPacket :exec
DELETE FROM mineralpackets
WHERE
    id = ?;