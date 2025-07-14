--
-- ShipDesigns
--
-- name: GetShipDesign :one
SELECT
    *
FROM
    shipdesigns
WHERE
    id = ?;

-- name: GetShipDesignByNum :one
SELECT
    *
FROM
    shipDesigns
WHERE
    gameId = ?
    AND playerNum = ?
    AND num = ?;

-- name: GetShipDesigns :many
SELECT
    *
FROM
    shipdesigns;

-- name: GetShipDesignsForGame :many
SELECT
    *
FROM
    shipDesigns
WHERE
    gameId = ?;

-- name: GetShipDesignsForPlayer :many
SELECT
    *
FROM
    shipDesigns
WHERE
    gameId = ?
    AND playerNum = ?;

-- name: CreateShipDesign :one
INSERT INTO
    shipDesigns (
        createdAt,
        updatedAt,
        gameId,
        num,
        playerNum,
        originalPlayerNum,
        name,
        version,
        hull,
        hullSetNumber,
        cannotDelete,
        slots,
        purpose,
        canDelete,
        mysteryTrader,
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
        ?
    ) RETURNING id, createdAt, updatedAt;

-- name: UpdateShipDesign :one
UPDATE shipdesigns
SET
    updatedAt = CURRENT_TIMESTAMP,
    gameId = ?,
    num = ?,
    playerNum = ?,
    originalPlayerNum = ?,
    name = ?,
    version = ?,
    hull = ?,
    hullSetNumber = ?,
    cannotDelete = ?,
    slots = ?,
    purpose = ?,
    canDelete = ?,
    mysteryTrader = ?,
    spec = ?
WHERE
    id = ? RETURNING updatedAt;

-- name: DeleteShipDesign :exec
DELETE FROM shipdesigns
WHERE
    id = ?;