--
-- MysteryTraders
--
-- name: GetMysteryTrader :one
SELECT
    *
FROM
    mystery_traders
WHERE
    id = ?;

-- name: GetMysteryTraderByNum :one
SELECT
    *
FROM
    mystery_traders
WHERE
    game_id = ?
    AND num = ?;

-- name: GetMysteryTraders :many
SELECT
    *
FROM
    mystery_traders;

-- name: GetMysteryTradersForGame :many
SELECT
    *
FROM
    mystery_traders
WHERE
    game_id = ?
ORDER BY
    num;

-- name: CreateMysteryTrader :execlastid
INSERT INTO
    mystery_traders (
        created_at,
        updated_at,
        game_id,
        intel_player_num,
        report_age,
        x,
        y,
        name,
        num,
        tags,
        heading_x,
        heading_y,
        warp_speed,
        requested_boon,
        destination_x,
        destination_y,
        reward_type,
        players_rewarded,
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
        ?
    );

-- name: UpdateMysteryTrader :execrows
UPDATE mystery_traders
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
    heading_x = ?,
    heading_y = ?,
    warp_speed = ?,
    requested_boon = ?,
    destination_x = ?,
    destination_y = ?,
    reward_type = ?,
    players_rewarded = ?,
    spec = ?
WHERE
    id = ?;

-- name: DeleteMysteryTrader :execrows
DELETE FROM mystery_traders
WHERE
    id = ?;