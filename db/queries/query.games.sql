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

-- name: GetGamesWithPlayers :many
SELECT
    sqlc.embed(g),
    p.updatedAt AS 'player.UpdatedAt',
    p.userId AS 'player.UserId',
    p.name AS 'player.Name',
    p.num AS 'player.Num',
    p.ready AS 'player.Ready',
    p.aiControlled AS 'player.AIControlled',
    p.aiDifficulty AS 'player.AIDifficulty',
    p.submittedTurn AS 'player.SubmittedTurn',
    p.color AS 'player.Color',
    p.victor AS 'player.Victor',
    p.archived AS 'player.Archived',
    p.guest AS 'player.Guest'
FROM
    games g
    LEFT JOIN players p ON g.id = p.gameId
WHERE
    -- game by id
    (
        @id IS NULL
        OR g.id = @id
    )
    -- host or player in game
    AND (
        (
            @hostId IS NULL
            AND @userId IS NULL
        )
        OR (
            g.hostId = @hostId
            OR g.id IN (
                SELECT
                    gameId
                FROM
                    players p
                WHERE
                    p.userId = @userId
            )
        )
    )
    --  host of game
    AND (
        @hostId IS NULL
        OR @userID IS NULL -- we handle userId above
        OR hostId = @hostId
    )
    -- game state
    AND (
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
    )
    -- game by hash
    AND (
        @hash IS NULL
        OR g.hash = @hash
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
    ) RETURNING *;

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