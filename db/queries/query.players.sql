--
-- Players
--
-- name: GetPlayer :one
SELECT
    *
FROM
    players
WHERE
    id = ?;

-- name: GetPlayerForGame :many
SELECT
    sqlc.embed(p),
    d.*
FROM
    players p
    LEFT JOIN ship_designs d ON p.game_id = d.game_id
    AND p.num = d.player_num
WHERE
    p.game_id = ?
    AND p.num = ?
ORDER BY
    d.num;

-- name: GetPlayerIntel :one
SELECT
    score_history,
    battle_records,
    player_intels,
    score_intels,
    planet_intels,
    fleet_intels,
    ship_design_intels,
    mineral_packet_intels,
    minefield_intels,
    wormhole_intels,
    mystery_trader_intels,
    salvage_intels
FROM
    players p
WHERE
    p.game_id = ?
    AND p.num = ?;

-- name: GetPlayers :many
SELECT
    *
FROM
    players;

-- name: GetPlayersForUser :many
SELECT
    *
FROM
    players
WHERE
    user_id = ?;

-- name: GetPlayersForGame :many
SELECT
    *
FROM
    players
WHERE
    game_id = ?
ORDER BY
    num;

-- name: GetPlayersWithDesignsForGame :many
SELECT
    sqlc.embed(p),
    d.*
FROM
    players p
    LEFT JOIN ship_designs d ON p.game_id = d.game_id
    AND p.num = d.player_num
WHERE
    p.game_id = ?
ORDER BY
    p.num,
    d.num;

-- name: GetPlayersStatusForGame :many
SELECT
    id,
    created_at,
    updated_at,
    game_id,
    user_id,
    name,
    num,
    ready,
    ai_controlled,
    ai_difficulty,
    guest,
    submitted_turn,
    color
FROM
    players
WHERE
    game_id = ?
ORDER BY
    num;

-- name: GetLightPlayerForGame :one
SELECT
    id,
    created_at,
    updated_at,
    game_id,
    user_id,
    name,
    num,
    ready,
    ai_controlled,
    submitted_turn,
    color,
    default_hull_set,
    tech_levels_energy,
    tech_levels_weapons,
    tech_levels_propulsion,
    tech_levels_construction,
    tech_levels_electronics,
    tech_levels_biotechnology,
    tech_levels_spent_energy,
    tech_levels_spent_weapons,
    tech_levels_spent_propulsion,
    tech_levels_spent_construction,
    tech_levels_spent_electronics,
    tech_levels_spent_biotechnology,
    research_amount,
    research_spent_last_year,
    next_research_field,
    researching,
    battle_plans,
    production_plans,
    transport_plans,
    relations,
    cargo_transfers,
    messages,
    race,
    stats,
    score_history,
    achieved_victory_conditions,
    victor,
    guest,
    ai_difficulty,
    acquired_techs,
    archived
FROM
    players
WHERE
    game_id = @game_id
    --  player_num
    AND (
        @player_num IS NULL
        OR num = @player_num
    )
    --  or user_id
    AND (
        @user_id IS NULL
        OR user_id = @user_id
    );

-- name: GetPlayerNum :one
SELECT
    num
FROM
    players
WHERE
    game_id = ?
    AND user_id = ?;

-- name: CreatePlayer :execlastid
INSERT INTO
    players (
        created_at,
        updated_at,
        game_id,
        user_id,
        name,
        num,
        ready,
        ai_controlled,
        submitted_turn,
        color,
        default_hull_set,
        tech_levels_energy,
        tech_levels_weapons,
        tech_levels_propulsion,
        tech_levels_construction,
        tech_levels_electronics,
        tech_levels_biotechnology,
        tech_levels_spent_energy,
        tech_levels_spent_weapons,
        tech_levels_spent_propulsion,
        tech_levels_spent_construction,
        tech_levels_spent_electronics,
        tech_levels_spent_biotechnology,
        research_amount,
        research_spent_last_year,
        next_research_field,
        researching,
        battle_plans,
        production_plans,
        transport_plans,
        relations,
        cargo_transfers,
        messages,
        battle_records,
        player_intels,
        score_intels,
        planet_intels,
        fleet_intels,
        ship_design_intels,
        mineral_packet_intels,
        minefield_intels,
        wormhole_intels,
        mystery_trader_intels,
        salvage_intels,
        race,
        stats,
        score_history,
        achieved_victory_conditions,
        victor,
        guest,
        ai_difficulty,
        acquired_techs,
        archived
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

-- name: UpdateLightPlayer :execrows
UPDATE players
SET
    updated_at = CURRENT_TIMESTAMP,
    name = ?,
    num = ?,
    ready = ?,
    ai_controlled = ?,
    ai_difficulty = ?,
    guest = ?,
    submitted_turn = ?,
    color = ?,
    default_hull_set = ?,
    research_amount = ?,
    next_research_field = ?,
    researching = ?
WHERE
    id = ?;

-- name: UpdatePlayerOrders :execrows
UPDATE players
SET
    updated_at = CURRENT_TIMESTAMP,
    submitted_turn = ?,
    default_hull_set = ?,
    research_amount = ?,
    next_research_field = ?,
    researching = ?,
    cargo_transfers = ?,
    battle_plans = ?,
    production_plans = ?,
    transport_plans = ?,
    relations = ?
WHERE
    id = ?;

-- name: UpdatePlayerCargoTransfers :execrows
UPDATE players
SET
    updated_at = CURRENT_TIMESTAMP,
    cargo_transfers = ?
WHERE
    id = ?;

-- name: UpdatePlayerRelations :execrows
UPDATE players
SET
    updated_at = CURRENT_TIMESTAMP,
    relations = ?
WHERE
    id = ?;

-- name: SubmitPlayerTurn :execrows
UPDATE players
SET
    updated_at = CURRENT_TIMESTAMP,
    submitted_turn = ?
WHERE
    game_id = ?
    AND num = ?;

-- name: ArchivePlayer :execrows
UPDATE players
SET
    updated_at = CURRENT_TIMESTAMP,
    archived = ?
WHERE
    game_id = ?
    AND num = ?;

-- name: UpdatePlayerPlans :execrows
UPDATE players
SET
    updated_at = CURRENT_TIMESTAMP,
    battle_plans = ?,
    production_plans = ?,
    transport_plans = ?
WHERE
    id = ?;

-- name: UpdatePlayerPlanetIntels :execrows
UPDATE players
SET
    updated_at = CURRENT_TIMESTAMP,
    planet_intels = ?
WHERE
    id = ?;

-- name: UpdatePlayerFleetIntels :execrows
UPDATE players
SET
    updated_at = CURRENT_TIMESTAMP,
    fleet_intels = ?
WHERE
    id = ?;

-- name: UpdatePlayerSalvageIntels :execrows
UPDATE players
SET
    updated_at = CURRENT_TIMESTAMP,
    salvage_intels = ?
WHERE
    id = ?;

-- name: UpdatePlayerMineralPacketIntels :execrows
UPDATE players
SET
    updated_at = CURRENT_TIMESTAMP,
    mineral_packet_intels = ?
WHERE
    id = ?;

-- name: UpdatePlayerUserID :execrows
UPDATE players
SET
    updated_at = CURRENT_TIMESTAMP,
    user_id = ?
WHERE
    id = ?;

-- name: UpdatePlayer :execrows
UPDATE players
SET
    updated_at = CURRENT_TIMESTAMP,
    game_id = ?,
    user_id = ?,
    name = ?,
    num = ?,
    ready = ?,
    ai_controlled = ?,
    submitted_turn = ?,
    color = ?,
    default_hull_set = ?,
    tech_levels_energy = ?,
    tech_levels_weapons = ?,
    tech_levels_propulsion = ?,
    tech_levels_construction = ?,
    tech_levels_electronics = ?,
    tech_levels_biotechnology = ?,
    tech_levels_spent_energy = ?,
    tech_levels_spent_weapons = ?,
    tech_levels_spent_propulsion = ?,
    tech_levels_spent_construction = ?,
    tech_levels_spent_electronics = ?,
    tech_levels_spent_biotechnology = ?,
    research_amount = ?,
    research_spent_last_year = ?,
    next_research_field = ?,
    researching = ?,
    battle_plans = ?,
    production_plans = ?,
    transport_plans = ?,
    relations = ?,
    cargo_transfers = ?,
    messages = ?,
    battle_records = ?,
    player_intels = ?,
    score_intels = ?,
    planet_intels = ?,
    fleet_intels = ?,
    ship_design_intels = ?,
    mineral_packet_intels = ?,
    minefield_intels = ?,
    wormhole_intels = ?,
    mystery_trader_intels = ?,
    salvage_intels = ?,
    race = ?,
    stats = ?,
    score_history = ?,
    achieved_victory_conditions = ?,
    victor = ?,
    guest = ?,
    ai_difficulty = ?,
    acquired_techs = ?,
    archived = ?
WHERE
    id = ?;

-- name: DeletePlayer :execrows
DELETE FROM players
WHERE
    id = ?;

-- name: DeleteTransientFleets :execrows
DELETE FROM fleets
WHERE
    game_id = ?
    AND intel_player_num = ?;

-- name: DeleteTransientMinefields :execrows
DELETE FROM minefields
WHERE
    game_id = ?
    AND intel_player_num = ?;

-- name: DeleteTransientSalvages :execrows
DELETE FROM salvages
WHERE
    game_id = ?
    AND intel_player_num = ?;

-- name: DeleteTransientMineralPackets :execrows
DELETE FROM mineral_packets
WHERE
    game_id = ?
    AND intel_player_num = ?;

-- name: DeleteTransientMysteryTraders :execrows
DELETE FROM mystery_traders
WHERE
    game_id = ?
    AND intel_player_num = ?;