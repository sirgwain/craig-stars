--
-- Players
--
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
    userId = ?;

-- name: GetPlayersForGame :many
SELECT
    *
FROM
    players
WHERE
    gameId = ?;

-- name: GetPlayersWithDesignsForGame :many
SELECT
    sqlc.embed(p),
    d.id AS 'design.id',
    d.createdAt AS 'design.createdAt',
    d.updatedAt AS 'design.updatedAt',
    d.gameId AS 'design.gameId',
    d.num AS 'design.num',
    d.playerNum AS 'design.playerNum',
    d.name AS 'design.name',
    d.version AS 'design.version',
    d.hull AS 'design.hull',
    d.hullSetNumber AS 'design.hullSetNumber',
    d.canDelete AS 'design.canDelete',
    d.slots AS 'design.slots',
    d.purpose AS 'design.purpose',
    d.spec AS 'design.spec',
    d.cannotDelete AS 'design.cannotDelete',
    d.originalPlayerNum AS 'design.originalPlayerNum',
    d.mysteryTrader AS 'design.mysteryTrader'
FROM
    players p
    LEFT JOIN shipDesigns d ON p.gameId = d.gameId
    AND p.num = d.playerNum
WHERE
    p.gameId = ?;

-- name: GetPlayerForGame :many
SELECT
    sqlc.embed(p),
    d.id AS 'design.id',
    d.createdAt AS 'design.createdAt',
    d.updatedAt AS 'design.updatedAt',
    d.gameId AS 'design.gameId',
    d.num AS 'design.num',
    d.playerNum AS 'design.playerNum',
    d.name AS 'design.name',
    d.version AS 'design.version',
    d.hull AS 'design.hull',
    d.hullSetNumber AS 'design.hullSetNumber',
    d.canDelete AS 'design.canDelete',
    d.slots AS 'design.slots',
    d.purpose AS 'design.purpose',
    d.spec AS 'design.spec',
    d.cannotDelete AS 'design.cannotDelete',
    d.originalPlayerNum AS 'design.originalPlayerNum',
    d.mysteryTrader AS 'design.mysteryTrader'
FROM
    players p
    LEFT JOIN shipDesigns d ON p.gameId = d.gameId
    AND p.num = d.playerNum
WHERE
    p.gameId = @gameId
    --  playerNum
    AND (
        @playerNum IS NULL
        OR p.num = @playerNum
    )
    --  or userId
    AND (
        @userId IS NULL
        OR userId = @userId
    );

-- name: GetPlayersStatusForGame :many
SELECT
    id,
    createdAt,
    updatedAt,
    gameId,
    userId,
    name,
    num,
    ready,
    aiControlled,
    aiDifficulty,
    guest,
    submittedTurn,
    color
FROM
    players
WHERE
    gameId = ?
ORDER BY
    num;

-- name: GetPlayer :one
SELECT
    *
FROM
    players
WHERE
    id = ?;

-- name: GetLightPlayerForGame :one
SELECT
    id,
    createdAt,
    updatedAt,
    gameId,
    userId,
    name,
    num,
    ready,
    aiControlled,
    aiDifficulty,
    guest,
    submittedTurn,
    color,
    defaultHullSet,
    race,
    techLevelsEnergy,
    techLevelsWeapons,
    techLevelsPropulsion,
    techLevelsConstruction,
    techLevelsElectronics,
    techLevelsBiotechnology,
    techLevelsSpentEnergy,
    techLevelsSpentWeapons,
    techLevelsSpentPropulsion,
    techLevelsSpentConstruction,
    techLevelsSpentElectronics,
    techLevelsSpentBiotechnology,
    researchSpentLastYear,
    researchAmount,
    nextResearchField,
    researching,
    cargoTransfers,
    battlePlans,
    productionPlans,
    transportPlans,
    relations,
    stats,
    scoreHistory,
    acquiredTechs,
    achievedVictoryConditions,
    victor,
    archived,
    spec
FROM
    players
WHERE
    gameId = @gameId
    --  playerNum
    AND (
        @playerNum IS NULL
        OR num = @playerNum
    )
    --  or userId
    AND (
        @userId IS NULL
        OR userId = @userId
    );

-- name: GetPlayerNum :one
SELECT
    num
FROM
    players
WHERE
    gameId = ?
    AND userId = ?;

-- name: CreatePlayer :one
INSERT INTO
    players (
        createdAt,
        updatedAt,
        gameId,
        userId,
        name,
        num,
        ready,
        aiControlled,
        submittedTurn,
        color,
        defaultHullSet,
        techLevelsEnergy,
        techLevelsWeapons,
        techLevelsPropulsion,
        techLevelsConstruction,
        techLevelsElectronics,
        techLevelsBiotechnology,
        techLevelsSpentEnergy,
        techLevelsSpentWeapons,
        techLevelsSpentPropulsion,
        techLevelsSpentConstruction,
        techLevelsSpentElectronics,
        techLevelsSpentBiotechnology,
        researchAmount,
        researchSpentLastYear,
        nextResearchField,
        researching,
        battlePlans,
        productionPlans,
        transportPlans,
        relations,
        cargoTransfers,
        messages,
        battleRecords,
        playerIntels,
        scoreIntels,
        planetIntels,
        fleetIntels,
        shipDesignIntels,
        mineralPacketIntels,
        mineFieldIntels,
        wormholeIntels,
        mysteryTraderIntels,
        salvageIntels,
        race,
        stats,
        scoreHistory,
        achievedVictoryConditions,
        victor,
        spec,
        guest,
        aiDifficulty,
        acquiredTechs,
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
        ?,
        ?
    ) RETURNING *;

-- name: UpdateLightPlayer :one
UPDATE players
SET
    updatedAt = CURRENT_TIMESTAMP,
    name = ?,
    num = ?,
    ready = ?,
    aiControlled = ?,
    aiDifficulty = ?,
    guest = ?,
    submittedTurn = ?,
    color = ?,
    defaultHullSet = ?,
    researchAmount = ?,
    nextResearchField = ?,
    researching = ?,
    spec = ?
WHERE
    id = ? RETURNING updatedAt;

-- name: UpdatePlayerOrders :one
UPDATE players
SET
    updatedAt = CURRENT_TIMESTAMP,
    submittedTurn = ?,
    defaultHullSet = ?,
    researchAmount = ?,
    nextResearchField = ?,
    researching = ?,
    cargoTransfers = ?,
    battlePlans = ?,
    productionPlans = ?,
    transportPlans = ?,
    relations = ?,
    spec = ?
WHERE
    id = ? RETURNING updatedAt;

-- name: UpdatePlayerCargoTransfers :one
UPDATE players
SET
    updatedAt = CURRENT_TIMESTAMP,
    cargoTransfers = ?
WHERE
    id = ? RETURNING updatedAt;

-- name: UpdatePlayerRelations :one
UPDATE players
SET
    updatedAt = CURRENT_TIMESTAMP,
    relations = ?
WHERE
    id = ? RETURNING updatedAt;

-- name: SubmitPlayerTurn :one
UPDATE players
SET
    updatedAt = CURRENT_TIMESTAMP,
    submittedTurn = ?
WHERE
    gameId = ?
    AND num = ? RETURNING updatedAt;

-- name: ArchivePlayer :one
UPDATE players
SET
    updatedAt = CURRENT_TIMESTAMP,
    archived = ?
WHERE
    gameId = ?
    AND num = ? RETURNING updatedAt;

-- name: UpdatePlayerPlans :one
UPDATE players
SET
    updatedAt = CURRENT_TIMESTAMP,
    battlePlans = ?,
    productionPlans = ?,
    transportPlans = ?
WHERE
    id = ? RETURNING updatedAt;

-- name: UpdatePlayerSpec :one
UPDATE players
SET
    updatedAt = CURRENT_TIMESTAMP,
    spec = ?
WHERE
    id = ? RETURNING updatedAt;

-- name: UpdatePlayerPlanetIntels :one
UPDATE players
SET
    updatedAt = CURRENT_TIMESTAMP,
    planetIntels = ?
WHERE
    id = ? RETURNING updatedAt;

-- name: UpdatePlayerFleetIntels :one
UPDATE players
SET
    updatedAt = CURRENT_TIMESTAMP,
    fleetIntels = ?
WHERE
    id = ? RETURNING updatedAt;

-- name: UpdatePlayerSalvageIntels :one
UPDATE players
SET
    updatedAt = CURRENT_TIMESTAMP,
    salvageIntels = ?
WHERE
    id = ? RETURNING updatedAt;

-- name: UpdatePlayerMineralPacketIntels :one
UPDATE players
SET
    updatedAt = CURRENT_TIMESTAMP,
    mineralPacketIntels = ?
WHERE
    id = ? RETURNING updatedAt;

-- name: UpdatePlayerUserID :exec
UPDATE players
SET
    updatedAt = CURRENT_TIMESTAMP,
    userId = ?
WHERE
    id = ?;

-- name: UpdatePlayer :one
UPDATE players
SET
    updatedAt = CURRENT_TIMESTAMP,
    gameId = ?,
    userId = ?,
    name = ?,
    num = ?,
    ready = ?,
    aiControlled = ?,
    submittedTurn = ?,
    color = ?,
    defaultHullSet = ?,
    techLevelsEnergy = ?,
    techLevelsWeapons = ?,
    techLevelsPropulsion = ?,
    techLevelsConstruction = ?,
    techLevelsElectronics = ?,
    techLevelsBiotechnology = ?,
    techLevelsSpentEnergy = ?,
    techLevelsSpentWeapons = ?,
    techLevelsSpentPropulsion = ?,
    techLevelsSpentConstruction = ?,
    techLevelsSpentElectronics = ?,
    techLevelsSpentBiotechnology = ?,
    researchAmount = ?,
    researchSpentLastYear = ?,
    nextResearchField = ?,
    researching = ?,
    battlePlans = ?,
    productionPlans = ?,
    transportPlans = ?,
    relations = ?,
    cargoTransfers = ?,
    messages = ?,
    battleRecords = ?,
    playerIntels = ?,
    scoreIntels = ?,
    planetIntels = ?,
    fleetIntels = ?,
    shipDesignIntels = ?,
    mineralPacketIntels = ?,
    mineFieldIntels = ?,
    wormholeIntels = ?,
    mysteryTraderIntels = ?,
    salvageIntels = ?,
    race = ?,
    stats = ?,
    scoreHistory = ?,
    achievedVictoryConditions = ?,
    victor = ?,
    spec = ?,
    guest = ?,
    aiDifficulty = ?,
    acquiredTechs = ?,
    archived = ?
WHERE
    id = ? RETURNING *;

-- name: DeletePlayer :exec
DELETE FROM players
WHERE
    id = ?;