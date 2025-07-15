--
-- Games
--
-- name: GetGame :one
SELECT
    *
FROM
    games
WHERE
    id = ?;

-- name: GetGames :many
SELECT
    *
FROM
    games;

-- name: GetGamesForHost :many
SELECT
    *
FROM
    games
WHERE
    host_id = ?;

-- name: GetGameWithPlayers :many
SELECT
    sqlc.embed(g),
    p.*
FROM
    games g
    LEFT JOIN game_players p ON g.id = p.game_id
WHERE
    -- game by id
    g.id = @id
    OR @hash IS NOT NULL
    AND g.hash = @hash;

-- name: GetGamesWithPlayers :many
SELECT
    sqlc.embed(g),
    p.*
FROM
    games g
    LEFT JOIN game_players p ON g.id = p.game_id
WHERE
    -- game state
    (
        @state IS NULL
        OR state = @state
    )
    -- open games
    AND (
        @open IS NULL
        OR (
            @open
            AND g.open_player_slots > 0
        )
    )
    -- public games
    AND (
        @public IS NULL
        OR g.public = @public
    );

-- name: GetGamesWithPlayersForUser :many
SELECT
    sqlc.embed(g),
    p.*
FROM
    games g
    LEFT JOIN game_players p ON g.id = p.game_id
WHERE
    -- host or player in game
    g.host_id = @user_id
    OR g.id IN (
        SELECT
            game_id
        FROM
            players p
        WHERE
            p.user_id = @user_id
    );

-- name: CreateGame :execlastid
INSERT INTO
    games (
        created_at,
        updated_at,
        host_id,
        name,
        state,
        public,
        hash,
        size,
        density,
        player_positions,
        random_events,
        computer_players_form_alliances,
        public_player_scores,
        max_minerals,
        start_mode,
        quick_start_turns,
        open_player_slots,
        num_players,
        victory_conditions_conditions,
        victory_conditions_num_criteria_required,
        victory_conditions_years_passed,
        victory_conditions_own_planets,
        victory_conditions_attain_tech_level,
        victory_conditions_attain_tech_level_num_fields,
        victory_conditions_exceeds_score,
        victory_conditions_exceeds_second_place_score,
        victory_conditions_production_capacity,
        victory_conditions_own_capital_ships,
        victory_conditions_highest_score_after_years,
        seed,
        area_x,
        area_y,
        year,
        victor_declared,
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
        ?
    );

-- name: UpdateGame :execrows
UPDATE games
SET
    updated_at = CURRENT_TIMESTAMP,
    host_id = ?,
    name = ?,
    state = ?,
    public = ?,
    hash = ?,
    size = ?,
    density = ?,
    player_positions = ?,
    random_events = ?,
    computer_players_form_alliances = ?,
    public_player_scores = ?,
    max_minerals = ?,
    start_mode = ?,
    quick_start_turns = ?,
    open_player_slots = ?,
    num_players = ?,
    victory_conditions_conditions = ?,
    victory_conditions_num_criteria_required = ?,
    victory_conditions_years_passed = ?,
    victory_conditions_own_planets = ?,
    victory_conditions_attain_tech_level = ?,
    victory_conditions_attain_tech_level_num_fields = ?,
    victory_conditions_exceeds_score = ?,
    victory_conditions_exceeds_second_place_score = ?,
    victory_conditions_production_capacity = ?,
    victory_conditions_own_capital_ships = ?,
    victory_conditions_highest_score_after_years = ?,
    seed = ?,
    area_x = ?,
    area_y = ?,
    year = ?,
    victor_declared = ?,
    archived = ?
WHERE
    id = ?;

-- name: UpdateGameState :execrows
UPDATE games
SET
    updated_at = CURRENT_TIMESTAMP,
    state = ?
WHERE
    id = ?;

-- name: UpdateGameHost :execrows
UPDATE games
SET
    updated_at = CURRENT_TIMESTAMP,
    host_id = ?
WHERE
    id = ?;

-- name: DeleteGame :execrows
DELETE FROM games
WHERE
    id = ?;

-- name: DeleteUserGames :execrows
DELETE FROM games
WHERE
    host_id = ?;