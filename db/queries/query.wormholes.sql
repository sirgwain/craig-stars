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
    game_id = ?
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
    game_id = ?
ORDER BY
    num;

-- name: CreateWormhole :execlastid
INSERT INTO
    wormholes (
        created_at,
        updated_at,
        intel_player_num,
        report_age,
        game_id,
        x,
        y,
        name,
        num,
        tags,
        destination_num,
        stability,
        years_at_stability,
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
        ?
    );

-- name: UpdateWormhole :execrows
UPDATE wormholes
SET
    updated_at = CURRENT_TIMESTAMP,
    game_id = ?,
    intel_player_num = ?,
    report_age = ?,
    x = ?,
    y = ?,
    name = ?,
    num = ?,
    tags = ?,
    destination_num = ?,
    stability = ?,
    years_at_stability = ?,
    spec = ?
WHERE
    id = ?;

-- name: DeleteWormhole :execrows
DELETE FROM wormholes
WHERE
    id = ?;