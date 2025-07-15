--
-- Races
--
-- name: GetRace :one
SELECT
    *
FROM
    races
WHERE
    id = ?;

-- name: GetRaces :many
SELECT
    *
FROM
    races;

-- name: GetRacesForUser :many
SELECT
    *
FROM
    races
WHERE
    user_id = ?;

-- name: CreateRace :one
INSERT INTO
    races (
        created_at,
        updated_at,
        user_id,
        name,
        plural_name,
        spend_leftover_points_on,
        prt,
        lrts,
        hab_low_grav,
        hab_low_temp,
        hab_low_rad,
        hab_high_grav,
        hab_high_temp,
        hab_high_rad,
        growth_rate,
        pop_efficiency,
        factory_output,
        factory_cost,
        num_factories,
        factories_cost_less,
        immune_grav,
        immune_temp,
        immune_rad,
        mine_output,
        mine_cost,
        num_mines,
        research_cost_energy,
        research_cost_weapons,
        research_cost_propulsion,
        research_cost_construction,
        research_cost_electronics,
        research_cost_biotechnology,
        techs_start_high,
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
        ?,
        ?,
        ?
    ) RETURNING id,
    created_at,
    updated_at;

-- name: UpdateRace :one
UPDATE races
SET
    updated_at = CURRENT_TIMESTAMP,
    user_id = ?,
    name = ?,
    plural_name = ?,
    spend_leftover_points_on = ?,
    prt = ?,
    lrts = ?,
    hab_low_grav = ?,
    hab_low_temp = ?,
    hab_low_rad = ?,
    hab_high_grav = ?,
    hab_high_temp = ?,
    hab_high_rad = ?,
    growth_rate = ?,
    pop_efficiency = ?,
    factory_output = ?,
    factory_cost = ?,
    num_factories = ?,
    factories_cost_less = ?,
    immune_grav = ?,
    immune_temp = ?,
    immune_rad = ?,
    mine_output = ?,
    mine_cost = ?,
    num_mines = ?,
    research_cost_energy = ?,
    research_cost_weapons = ?,
    research_cost_propulsion = ?,
    research_cost_construction = ?,
    research_cost_electronics = ?,
    research_cost_biotechnology = ?,
    techs_start_high = ?,
    spec = ?
WHERE
    id = ? RETURNING updated_at;

-- name: DeleteRace :exec
DELETE FROM races
WHERE
    id = ?;

-- name: DeleteUserRaces :exec
DELETE FROM races
WHERE
    user_id = ?;