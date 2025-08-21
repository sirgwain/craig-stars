--
-- ShipDesigns
--
-- name: GetShipDesign :one
SELECT
    *
FROM
    ship_designs
WHERE
    id = ?;

-- name: GetShipDesignByNum :one
SELECT
    *
FROM
    ship_designs
WHERE
    game_id = ?
    AND player_num = ?
    AND num = ?;

-- name: GetShipDesigns :many
SELECT
    *
FROM
    ship_designs;

-- name: GetShipDesignsForGame :many
SELECT
    *
FROM
    ship_designs
WHERE
    game_id = ?;

-- name: GetShipDesignsForPlayer :many
SELECT
    *
FROM
    ship_designs
WHERE
    game_id = ?
    AND player_num = ?;

-- name: CreateShipDesign :execlastid
INSERT INTO
    ship_designs (
        created_at,
        updated_at,
        game_id,
        intel_player_num,
        num,
        player_num,
        original_player_num,
        name,
        version,
        hull,
        hull_set_number,
        cannot_delete,
        slots,
        purpose,
        mystery_trader,
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
    );

-- name: UpdateShipDesign :execrows
UPDATE ship_designs
SET
    updated_at = CURRENT_TIMESTAMP,
    game_id = ?,
    intel_player_num = ?,
    num = ?,
    player_num = ?,
    original_player_num = ?,
    name = ?,
    version = ?,
    hull = ?,
    hull_set_number = ?,
    cannot_delete = ?,
    slots = ?,
    purpose = ?,
    mystery_trader = ?,
    spec = ?
WHERE
    id = ?;

-- name: DeleteShipDesign :execrows
DELETE FROM ship_designs
WHERE
    id = ?;