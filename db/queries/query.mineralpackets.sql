--
-- MineralPackets
--
-- name: GetMineralPacket :one
SELECT
    *
FROM
    mineral_packets
WHERE
    id = ?;

-- name: GetMineralPacketByNum :one
SELECT
    *
FROM
    mineral_packets
WHERE
    game_id = ?
    AND player_num = ?
    AND num = ?;

-- name: GetMineralPackets :many
SELECT
    *
FROM
    mineral_packets;

-- name: GetMineralPacketsForGame :many
SELECT
    *
FROM
    mineral_packets
WHERE
    game_id = ?
ORDER BY
    player_num,
    num;

-- name: GetMineralPacketsForPlayer :many
SELECT
    *
FROM
    mineral_packets
WHERE
    game_id = ?
    AND player_num = ?
ORDER BY
    num;

-- name: CreateMineralPacket :execlastid
INSERT INTO
    mineral_packets (
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
        target_planet_num,
        ironium,
        boranium,
        germanium,
        safe_warp_speed,
        warp_speed,
        scan_range,
        scan_range_pen,
        heading_x,
        heading_y
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
        ?
    );

-- name: UpdateMineralPacket :execrows
UPDATE mineral_packets
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
    target_planet_num = ?,
    ironium = ?,
    boranium = ?,
    germanium = ?,
    safe_warp_speed = ?,
    warp_speed = ?,
    scan_range = ?,
    scan_range_pen = ?,
    heading_x = ?,
    heading_y = ?
WHERE
    id = ?;

-- name: DeleteMineralPacket :execrows
DELETE FROM mineral_packets
WHERE
    id = ?;