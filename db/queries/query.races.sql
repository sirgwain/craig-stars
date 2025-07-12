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
    userId = ?;

-- name: CreateRace :one
INSERT INTO
    races (
        createdAt,
        updatedAt,
        userId,
        name,
        pluralName,
        spendLeftoverPointsOn,
        prt,
        lrts,
        habLowGrav,
        habLowTemp,
        habLowRad,
        habHighGrav,
        habHighTemp,
        habHighRad,
        growthRate,
        popEfficiency,
        factoryOutput,
        factoryCost,
        numFactories,
        factoriesCostLess,
        immuneGrav,
        immuneTemp,
        immuneRad,
        mineOutput,
        mineCost,
        numMines,
        researchCostEnergy,
        researchCostWeapons,
        researchCostPropulsion,
        researchCostConstruction,
        researchCostElectronics,
        researchCostBiotechnology,
        techsStartHigh,
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
    ) RETURNING *;

-- name: UpdateRace :one
UPDATE races
SET
    updatedAt = CURRENT_TIMESTAMP,
    userId = ?,
    name = ?,
    pluralName = ?,
    spendLeftoverPointsOn = ?,
    prt = ?,
    lrts = ?,
    habLowGrav = ?,
    habLowTemp = ?,
    habLowRad = ?,
    habHighGrav = ?,
    habHighTemp = ?,
    habHighRad = ?,
    growthRate = ?,
    popEfficiency = ?,
    factoryOutput = ?,
    factoryCost = ?,
    numFactories = ?,
    factoriesCostLess = ?,
    immuneGrav = ?,
    immuneTemp = ?,
    immuneRad = ?,
    mineOutput = ?,
    mineCost = ?,
    numMines = ?,
    researchCostEnergy = ?,
    researchCostWeapons = ?,
    researchCostPropulsion = ?,
    researchCostConstruction = ?,
    researchCostElectronics = ?,
    researchCostBiotechnology = ?,
    techsStartHigh = ?,
    spec = ?
WHERE
    id = ? RETURNING *;

-- name: DeleteRace :exec
DELETE FROM races
WHERE
    id = ?;

-- name: DeleteUserRaces :exec
DELETE FROM races
WHERE
    userId = ?;
