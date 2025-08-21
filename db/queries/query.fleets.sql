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
    game_id = ?
    AND player_num = ?
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
    game_id = ?
ORDER BY
    player_num,
    num;

-- name: GetFleetsForPlayer :many
SELECT
    *
FROM
    fleets
WHERE
    game_id = ?
    AND player_num = ?
ORDER BY
    num;

-- name: GetFleetsByNums :many
SELECT
    *
FROM
    fleets
WHERE
    game_id = ?
    AND player_num = ?
    AND num IN (sqlc.slice ('nums'))
ORDER BY
    player_num,
    num;

-- name: CreateFleet :execlastid
INSERT INTO
    fleets (
        created_at,
        updated_at,
        game_id,
        intel_player_num,
        report_age,
        battle_plan_num,
        x,
        y,
        name,
        num,
        player_num,
        tags,
        tokens,
        waypoints,
        repeat_orders,
        planet_num,
        base_name,
        ironium,
        boranium,
        germanium,
        colonists,
        fuel,
        age,
        heading_x,
        heading_y,
        warp_speed,
        previous_position_x,
        previous_position_y,
        orbiting_planet_num,
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
        ?,
        ?,
        ?
    );

-- name: UpdateFleet :execrows
UPDATE fleets
SET
    updated_at = CURRENT_TIMESTAMP,
    game_id = ?,
    intel_player_num = ?,
    report_age = ?,
    battle_plan_num = ?,
    x = ?,
    y = ?,
    name = ?,
    num = ?,
    player_num = ?,
    tags = ?,
    tokens = ?,
    waypoints = ?,
    repeat_orders = ?,
    planet_num = ?,
    base_name = ?,
    ironium = ?,
    boranium = ?,
    germanium = ?,
    colonists = ?,
    fuel = ?,
    age = ?,
    heading_x = ?,
    heading_y = ?,
    warp_speed = ?,
    previous_position_x = ?,
    previous_position_y = ?,
    orbiting_planet_num = ?,
    starbase = ?,
    purpose = ?,
    spec = ?
WHERE
    id = ?;

-- name: DeleteFleet :exec
DELETE FROM fleets
WHERE
    id = ?;