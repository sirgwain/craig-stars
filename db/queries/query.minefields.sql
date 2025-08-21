--
-- Minefields
--
-- name: GetMinefield :one
SELECT
    *
FROM
    minefields
WHERE
    id = ?;

-- name: GetMinefieldByNum :one
SELECT
    *
FROM
    minefields
WHERE
    game_id = ?
    AND player_num = ?
    AND num = ?;

-- name: GetMinefields :many
SELECT
    *
FROM
    minefields;

-- name: GetMinefieldsForGame :many
SELECT
    *
FROM
    minefields
WHERE
    game_id = ?
ORDER BY
    player_num,
    num;

-- name: GetMinefieldsForPlayer :many
SELECT
    *
FROM
    minefields
WHERE
    game_id = ?
    AND player_num = ?
ORDER BY
    num;

-- name: CreateMinefield :execlastid
INSERT INTO
    minefields (
        created_at,
        updated_at,
        game_id,
        intel_player_num,
        report_age,
        x,
        y,
        name,
        num,
        player_num,
        tags,
        minefield_type,
        num_mines,
        detonate
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

-- name: UpdateMinefield :execrows
UPDATE minefields
SET
    updated_at = CURRENT_TIMESTAMP,
    game_id = ?,
    intel_player_num = ?,
    report_age = ?,
    x = ?,
    y = ?,
    name = ?,
    num = ?,
    player_num = ?,
    tags = ?,
    minefield_type = ?,
    num_mines = ?,
    detonate = ?
WHERE
    id = ?;

-- name: DeleteMinefield :execrows
DELETE FROM minefields
WHERE
    id = ?;