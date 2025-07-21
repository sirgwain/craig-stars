--
-- Planets
--
-- name: GetPlanet :one
SELECT
    *
FROM
    planets
WHERE
    id = ?;

-- name: GetPlanetByNum :many
SELECT
    sqlc.embed(p),
    f.id AS 'fleet.id',
    f.created_at AS 'fleet.created_at',
    f.updated_at AS 'fleet.updated_at',
    f.game_id AS 'fleet.game_id',
    COALESCE(f.battle_plan_num, 0) AS 'fleet.battle_plan_num',
    f.x AS 'fleet.x',
    f.y AS 'fleet.y',
    f.name AS 'fleet.name',
    COALESCE(f.num, 0) AS 'fleet.num',
    COALESCE(f.player_num, 0) AS 'fleet.player_num',
    f.tokens AS 'fleet.tokens',
    f.waypoints AS 'fleet.waypoints',
    f.repeat_orders AS 'fleet.repeat_orders',
    COALESCE(f.planet_num, 0) AS 'fleet.planet_num',
    f.base_name AS 'fleet.base_name',
    COALESCE(f.ironium, 0) AS 'fleet.ironium',
    COALESCE(f.boranium, 0) AS 'fleet.boranium',
    COALESCE(f.germanium, 0) AS 'fleet.germanium',
    COALESCE(f.colonists, 0) AS 'fleet.colonists',
    COALESCE(f.fuel, 0) AS 'fleet.fuel',
    COALESCE(f.age, 0) AS 'fleet.age',
    f.heading_x AS 'fleet.heading_x',
    f.heading_y AS 'fleet.heading_y',
    COALESCE(f.warp_speed, 0) AS 'fleet.warp_speed',
    f.previous_position_x AS 'fleet.previous_position_x',
    f.previous_position_y AS 'fleet.previous_position_y',
    COALESCE(f.orbiting_planet_num, 0) AS 'fleet.orbiting_planet_num',
    f.starbase AS 'fleet.starbase',
    f.spec AS 'fleet.spec',
    f.purpose AS 'fleet.purpose',
    f.tags AS 'fleet.tags'
FROM
    planets p
    LEFT JOIN fleets f ON p.game_id = f.game_id
    AND p.num = f.planet_num
WHERE
    p.game_id = ?
    AND p.num = ?;

-- name: GetPlanets :many
SELECT
    *
FROM
    planets;

-- name: GetPlanetsForGame :many
SELECT
    *
FROM
    planets
WHERE
    game_id = ?
ORDER BY
    num;

-- name: GetPlanetsForPlayer :many
SELECT
    *
FROM
    planets
WHERE
    game_id = ?
    AND player_num = ?
ORDER BY
    num;

-- name: CreatePlanet :execlastid
INSERT INTO
    planets (
        created_at,
        updated_at,
        game_id,
        x,
        y,
        name,
        num,
        player_num,
        grav,
        temp,
        rad,
        base_grav,
        base_temp,
        base_rad,
        terraformed_amount_grav,
        terraformed_amount_temp,
        terraformed_amount_rad,
        mineral_conc_ironium,
        mineral_conc_boranium,
        mineral_conc_germanium,
        mine_years_ironium,
        mine_years_boranium,
        mine_years_germanium,
        ironium,
        boranium,
        germanium,
        colonists,
        partial_population,
        mines,
        factories,
        defenses,
        homeworld,
        contributes_only_leftover_to_research,
        scanner,
        route_target_type,
        route_target_num,
        route_target_player_num,
        packet_target_num,
        packet_speed,
        production_queue,
        spec,
        tags,
        random_artifact
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

-- name: UpdatePlanet :execrows
UPDATE planets
SET
    updated_at = CURRENT_TIMESTAMP,
    game_id = ?,
    x = ?,
    y = ?,
    name = ?,
    num = ?,
    player_num = ?,
    grav = ?,
    temp = ?,
    rad = ?,
    base_grav = ?,
    base_temp = ?,
    base_rad = ?,
    terraformed_amount_grav = ?,
    terraformed_amount_temp = ?,
    terraformed_amount_rad = ?,
    mineral_conc_ironium = ?,
    mineral_conc_boranium = ?,
    mineral_conc_germanium = ?,
    mine_years_ironium = ?,
    mine_years_boranium = ?,
    mine_years_germanium = ?,
    ironium = ?,
    boranium = ?,
    germanium = ?,
    colonists = ?,
    partial_population = ?,
    mines = ?,
    factories = ?,
    defenses = ?,
    homeworld = ?,
    contributes_only_leftover_to_research = ?,
    scanner = ?,
    route_target_type = ?,
    route_target_num = ?,
    route_target_player_num = ?,
    packet_target_num = ?,
    packet_speed = ?,
    production_queue = ?,
    spec = ?,
    tags = ?,
    random_artifact = ?
WHERE
    id = ?;

-- name: UpdatePlanetSpec :execrows
UPDATE planets
SET
    spec = ?
WHERE
    id = ?;

-- name: DeletePlanet :exec
DELETE FROM planets
WHERE
    id = ?;