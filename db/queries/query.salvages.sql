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
    game_id = ?
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
    game_id = ?
ORDER BY
    num;

-- name: GetSalvagesForPlayer :many
SELECT
    *
FROM
    salvages
WHERE
    game_id = ?
    AND player_num = ?
ORDER BY
    num;

-- name: CreateSalvage :execlastid
INSERT INTO
    salvages (
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
        ?,
        ?,
        ?
    );

-- name: UpdateSalvage :execrows
UPDATE salvages
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
    ironium = ?,
    boranium = ?,
    germanium = ?
WHERE
    id = ?;

-- name: DeleteSalvage :execrows
DELETE FROM salvages
WHERE
    id = ?;