--
-- Fleets
--
-- name: GetFleet :one
SELECT
    *
FROM
    fleets
WHERE
    id = ?;

-- name: GetFleetByNum :one
SELECT
    *
FROM
    fleets
WHERE
    gameId = ?
    AND playerNum = ?
    AND num = ?;

-- name: GetFleets :many
SELECT
    *
FROM
    fleets;

-- name: GetFleetsForGame :many
SELECT
    *
FROM
    fleets
WHERE
    gameId = ?
ORDER BY
    playerNum,
    num;

-- name: GetFleetsForPlayer :many
SELECT
    *
FROM
    fleets
WHERE
    gameId = ?
    AND playerNum = ?
ORDER BY
    num;

-- name: GetFleetsByNums :many
SELECT
    *
FROM
    fleets
WHERE
    gameId = ?
    AND playerNum = ?
    AND num IN (sqlc.slice ('nums'))
ORDER BY
    playerNum,
    num;

-- name: CreateFleet :one
INSERT INTO
    fleets (
        createdAt,
        updatedAt,
        gameId,
        battlePlanNum,
        x,
        y,
        name,
        num,
        playerNum,
        tags,
        tokens,
        waypoints,
        repeatOrders,
        planetNum,
        baseName,
        ironium,
        boranium,
        germanium,
        colonists,
        fuel,
        age,
        headingX,
        headingY,
        warpSpeed,
        previousPositionX,
        previousPositionY,
        orbitingPlanetNum,
        starbase,
        purpose,
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
    ) RETURNING id,
    createdAt,
    updatedAt;

-- name: UpdateFleet :one
UPDATE fleets
SET
    updatedAt = CURRENT_TIMESTAMP,
    gameId = ?,
    battlePlanNum = ?,
    x = ?,
    y = ?,
    name = ?,
    num = ?,
    playerNum = ?,
    tags = ?,
    tokens = ?,
    waypoints = ?,
    repeatOrders = ?,
    planetNum = ?,
    baseName = ?,
    ironium = ?,
    boranium = ?,
    germanium = ?,
    colonists = ?,
    fuel = ?,
    age = ?,
    headingX = ?,
    headingY = ?,
    warpSpeed = ?,
    previousPositionX = ?,
    previousPositionY = ?,
    orbitingPlanetNum = ?,
    starbase = ?,
    purpose = ?,
    spec = ?
WHERE
    id = ? RETURNING updatedAt;

-- name: DeleteFleet :exec
DELETE FROM fleets
WHERE
    id = ?;