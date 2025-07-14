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
    hostId = ?;

-- name: GetGameWithPlayers :many
SELECT
    sqlc.embed(g),
    p.*
FROM
    games g
    LEFT JOIN game_players p ON g.id = p.gameId
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
    LEFT JOIN game_players p ON g.id = p.gameId
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
            AND g.openPlayerSlots > 0
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
    LEFT JOIN game_players p ON g.id = p.gameId
WHERE
    -- host or player in game
    g.hostId = @userId
    OR g.id IN (
        SELECT
            gameId
        FROM
            players p
        WHERE
            p.userId = @userId
    );

-- name: CreateGame :one
INSERT INTO
    games (
        createdAt,
        updatedAt,
        hostId,
        name,
        state,
        public,
        hash,
        size,
        density,
        playerPositions,
        randomEvents,
        computerPlayersFormAlliances,
        publicPlayerScores,
        maxMinerals,
        startMode,
        quickStartTurns,
        openPlayerSlots,
        numPlayers,
        victoryConditionsConditions,
        victoryConditionsNumCriteriaRequired,
        victoryConditionsYearsPassed,
        victoryConditionsOwnPlanets,
        victoryConditionsAttainTechLevel,
        victoryConditionsAttainTechLevelNumFields,
        victoryConditionsExceedsScore,
        victoryConditionsExceedsSecondPlaceScore,
        victoryConditionsProductionCapacity,
        victoryConditionsOwnCapitalShips,
        victoryConditionsHighestScoreAfterYears,
        seed,
        rules,
        areaX,
        areaY,
        year,
        victorDeclared,
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
        ?
    ) RETURNING id,
    createdAt,
    updatedAt;

-- name: UpdateGame :one
UPDATE games
SET
    updatedAt = CURRENT_TIMESTAMP,
    hostId = ?,
    name = ?,
    state = ?,
    public = ?,
    hash = ?,
    size = ?,
    density = ?,
    playerPositions = ?,
    randomEvents = ?,
    computerPlayersFormAlliances = ?,
    publicPlayerScores = ?,
    maxMinerals = ?,
    startMode = ?,
    quickStartTurns = ?,
    openPlayerSlots = ?,
    numPlayers = ?,
    victoryConditionsConditions = ?,
    victoryConditionsNumCriteriaRequired = ?,
    victoryConditionsYearsPassed = ?,
    victoryConditionsOwnPlanets = ?,
    victoryConditionsAttainTechLevel = ?,
    victoryConditionsAttainTechLevelNumFields = ?,
    victoryConditionsExceedsScore = ?,
    victoryConditionsExceedsSecondPlaceScore = ?,
    victoryConditionsProductionCapacity = ?,
    victoryConditionsOwnCapitalShips = ?,
    victoryConditionsHighestScoreAfterYears = ?,
    seed = ?,
    rules = ?,
    areaX = ?,
    areaY = ?,
    year = ?,
    victorDeclared = ?,
    archived = ?
WHERE
    id = ? RETURNING updatedAt;

-- name: UpdateGameState :exec
UPDATE games
SET
    updatedAt = CURRENT_TIMESTAMP,
    state = ?
WHERE
    id = ?;

-- name: UpdateGameHost :exec
UPDATE games
SET
    updatedAt = CURRENT_TIMESTAMP,
    hostId = ?
WHERE
    id = ?;

-- name: DeleteGame :exec
DELETE FROM games
WHERE
    id = ?;

-- name: DeleteUserGames :exec
DELETE FROM games
WHERE
    hostId = ?;