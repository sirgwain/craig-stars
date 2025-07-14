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
    f.createdAt AS 'fleet.createdAt',
    f.updatedAt AS 'fleet.updatedAt',
    f.gameId AS 'fleet.gameId',
    f.battlePlanNum AS 'fleet.battlePlanNum',
    f.x AS 'fleet.x',
    f.y AS 'fleet.y',
    f.name AS 'fleet.name',
    f.num AS 'fleet.num',
    f.playerNum AS 'fleet.playerNum',
    f.tokens AS 'fleet.tokens',
    f.waypoints AS 'fleet.waypoints',
    f.repeatOrders AS 'fleet.repeatOrders',
    f.planetNum AS 'fleet.planetNum',
    f.baseName AS 'fleet.baseName',
    f.ironium AS 'fleet.ironium',
    f.boranium AS 'fleet.boranium',
    f.germanium AS 'fleet.germanium',
    f.colonists AS 'fleet.colonists',
    f.fuel AS 'fleet.fuel',
    f.age AS 'fleet.age',
    f.headingX AS 'fleet.headingX',
    f.headingY AS 'fleet.headingY',
    f.warpSpeed AS 'fleet.warpSpeed',
    f.previousPositionX AS 'fleet.previousPositionX',
    f.previousPositionY AS 'fleet.previousPositionY',
    f.orbitingPlanetNum AS 'fleet.orbitingPlanetNum',
    f.starbase AS 'fleet.starbase',
    f.spec AS 'fleet.spec',
    f.purpose AS 'fleet.purpose',
    f.tags AS 'fleet.tags'
FROM
    planets p
    LEFT JOIN fleets f ON p.gameId = f.gameId
    AND p.num = f.planetNum
WHERE
    p.gameId = ?
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
    gameId = ?;

-- name: GetPlanetsForPlayer :many
SELECT
    *
FROM
    planets
WHERE
    gameId = ?
    AND playerNum = ?;

-- name: CreatePlanet :one
INSERT INTO
    planets (
        createdAt,
        updatedAt,
        gameId,
        x,
        y,
        name,
        num,
        playerNum,
    grav,
        TEMP,
        rad,
        baseGrav,
        baseTemp,
        baseRad,
        terraformedAmountGrav,
        terraformedAmountTemp,
        terraformedAmountRad,
        mineralConcIronium,
        mineralConcBoranium,
        mineralConcGermanium,
        mineYearsIronium,
        mineYearsBoranium,
        mineYearsGermanium,
        ironium,
        boranium,
        germanium,
        colonists,
        partialPopulation,
        mines,
        factories,
        defenses,
        homeworld,
        contributesOnlyLeftoverToResearch,
        scanner,
        routeTargetType,
        routeTargetNum,
        routeTargetPlayerNum,
        packetTargetNum,
        packetSpeed,
        productionQueue,
        spec,
        tags,
        randomArtifact
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
    ) RETURNING id, createdAt, updatedAt;

-- name: UpdatePlanet :one
UPDATE planets
SET
    updatedAt = CURRENT_TIMESTAMP,
    gameId = ?,
    x = ?,
    y = ?,
    name = ?,
    num = ?,
    playerNum = ?,
    grav = ?,
    TEMP = ?,
    rad = ?,
    baseGrav = ?,
    baseTemp = ?,
    baseRad = ?,
    terraformedAmountGrav = ?,
    terraformedAmountTemp = ?,
    terraformedAmountRad = ?,
    mineralConcIronium = ?,
    mineralConcBoranium = ?,
    mineralConcGermanium = ?,
    mineYearsIronium = ?,
    mineYearsBoranium = ?,
    mineYearsGermanium = ?,
    ironium = ?,
    boranium = ?,
    germanium = ?,
    colonists = ?,
    partialPopulation = ?,
    mines = ?,
    factories = ?,
    defenses = ?,
    homeworld = ?,
    contributesOnlyLeftoverToResearch = ?,
    scanner = ?,
    routeTargetType = ?,
    routeTargetNum = ?,
    routeTargetPlayerNum = ?,
    packetTargetNum = ?,
    packetSpeed = ?,
    productionQueue = ?,
    spec = ?,
    tags = ?,
    randomArtifact = ?
WHERE
    id = ? RETURNING updatedAt;

-- name: UpdatePlanetSpec :one
UPDATE planets
SET
    spec = ?
WHERE
    id = ? RETURNING updatedAt;

-- name: DeletePlanet :exec
DELETE FROM planets
WHERE
    id = ?;