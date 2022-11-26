CREATE TABLE users (
  id INTEGER PRIMARY KEY,
  createdAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  updatedAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  username TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL,
  role TEXT NOT NULL
);
CREATE TABLE races (
  id INTEGER PRIMARY KEY,
  createdAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  updatedAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  userId INTEGER NOT NULL,
  playerId INTEGER,
  name TEXT NOT NULL,
  pluralName TEXT NOT NULL,
  prt TEXT,
  lrts INTEGER,
  habLowGrav INTEGER,
  habLowTemp INTEGER,
  habLowRad INTEGER,
  habHighGrav INTEGER,
  habHighTemp INTEGER,
  habHighRad INTEGER,
  growthRate INTEGER,
  popEfficiency INTEGER,
  factoryOutput INTEGER,
  factoryCost INTEGER,
  numFactories INTEGER,
  factoriesCostLess NUMERIC,
  immuneGrav NUMERIC,
  immuneTemp NUMERIC,
  immuneRad NUMERIC,
  mineOutput INTEGER,
  mineCost INTEGER,
  numMines INTEGER,
  researchCostEnergy TEXT,
  researchCostWeapons TEXT,
  researchCostPropulsion TEXT,
  researchCostConstruction TEXT,
  researchCostElectronics TEXT,
  researchCostBiotechnology TEXT,
  techsStartHigh NUMERIC,
  spec TEXT,
  CONSTRAINT fkUsersRaces FOREIGN KEY (userId) REFERENCES users (id) ON DELETE
  SET NULL,
    CONSTRAINT fkPlayersRaces FOREIGN KEY (playerId) REFERENCES players (id) ON DELETE
  SET NULL
);
CREATE TABLE games (
  id INTEGER PRIMARY KEY,
  createdAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  updatedAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  name TEXT NOT NULL,
  hostId INTEGER,
  quickStartTurns INTEGER,
  size TEXT,
  density TEXT,
  playerPositions TEXT,
  randomEvents NUMERIC,
  computerPlayersFormAlliances NUMERIC,
  publicPlayerScores NUMERIC,
  startMode TEXT,
  year INTEGER,
  state TEXT,
  openPlayerSlots INTEGER,
  numPlayers INTEGER,
  victoryConditionsConditions TEXT,
  victoryConditionsNumCriteriaRequired INTEGER,
  victoryConditionsYearsPassed INTEGER,
  victoryConditionsOwnPlanets INTEGER,
  victoryConditionsAttainTechLevel INTEGER,
  victoryConditionsAttainTechLevelNumFields INTEGER,
  victoryConditionsExceedsScore INTEGER,
  victoryConditionsExceedsSecondPlaceScore INTEGER,
  victoryConditionsProductionCapacity INTEGER,
  victoryConditionsOwnCapitalShips INTEGER,
  victoryConditionsHighestScoreAfterYears INTEGER,
  victorDeclared NUMERIC,
  seed INTEGER,
  rules TEXT,
  areaX REAL,
  areaY REAL
);
CREATE TABLE rules (
  id INTEGER PRIMARY KEY,
  createdAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  updatedAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  gameId INTEGER NOT NULL,
  seed INTEGER,
  tachyonCloakReduction INTEGER,
  maxPopulation INTEGER,
  fleetsScanWhileMoving NUMERIC,
  populationScannerError REAL,
  smartDefenseCoverageFactor REAL,
  invasionDefenseCoverageFactor REAL,
  numBattleRounds INTEGER,
  movesToRunAway INTEGER,
  beamRangeDropoff REAL,
  torpedoSplashDamage REAL,
  salvageDecayRate REAL,
  salvageDecayMin INTEGER,
  mineFieldCloak INTEGER,
  stargateMaxRangeFactor INTEGER,
  stargateMaxHullMassFactor INTEGER,
  randomEventChances TEXT,
  randomMineralDepositBonusRange TEXT,
  wormholeCloak INTEGER,
  wormholeMinDistance INTEGER,
  wormholeStatsByStability TEXT,
  wormholePairsForSize TEXT,
  mineFieldStatsByType TEXT,
  repairRates TEXT,
  maxPlayers INTEGER,
  startingYear INTEGER,
  showPublicScoresAfterYears INTEGER,
  planetMinDistance INTEGER,
  maxExtraWorldDistance INTEGER,
  minExtraWorldDistance INTEGER,
  minHomeworldMineralConcentration INTEGER,
  minExtraPlanetMineralConcentration INTEGER,
  minMineralConcentration INTEGER,
  minStartingMineralConcentration INTEGER,
  maxStartingMineralConcentration INTEGER,
  highRadGermaniumBonus INTEGER,
  highRadGermaniumBonusThreshold INTEGER,
  maxStartingMineralSurface INTEGER,
  minStartingMineralSurface INTEGER,
  mineralDecayFactor INTEGER,
  startingMines INTEGER,
  startingFactories INTEGER,
  startingDefenses INTEGER,
  raceStartingPoints INTEGER,
  scrapMineralAmount REAL,
  scrapResourceAmount REAL,
  factoryCostGermanium INTEGER,
  defenseCost TEXT,
  mineralAlchemyCost INTEGER,
  terraformCost TEXT,
  starbaseComponentCostFactor REAL,
  packetDecayRate TEXT,
  maxTechLevel INTEGER,
  techBaseCost TEXT,
  prtSpecs TEXT,
  lrtSpecs TEXT,
  techsId INTEGER
);
CREATE TABLE players (
  id INTEGER PRIMARY KEY,
  createdAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  updatedAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  gameId INTEGER NOT NULL,
  userId INTEGER NOT NULL,
  name TEXT NOT NULL,
  num INTEGER,
  ready NUMERIC,
  aiControlled NUMERIC,
  submittedTurn NUMERIC,
  color TEXT,
  defaultHullSet INTEGER,
  techLevelsEnergy INTEGER,
  techLevelsWeapons INTEGER,
  techLevelsPropulsion INTEGER,
  techLevelsConstruction INTEGER,
  techLevelsElectronics INTEGER,
  techLevelsBiotechnology INTEGER,
  techLevelsSpentEnergy INTEGER,
  techLevelsSpentWeapons INTEGER,
  techLevelsSpentPropulsion INTEGER,
  techLevelsSpentConstruction INTEGER,
  techLevelsSpentElectronics INTEGER,
  techLevelsSpentBiotechnology INTEGER,
  researchAmount INTEGER,
  researchSpentLastYear INTEGER,
  nextResearchField TEXT,
  researching TEXT,
  battlePlans TEXT,
  productionPlans TEXT,
  transportPlans TEXT,
  messages TEXT,
  planetIntels TEXT,
  fleetIntels TEXT,
  shipDesignIntels TEXT,
  mineralPacketIntels TEXT,
  mineFieldIntels TEXT,
  race TEXT,
  stats TEXT,
  spec TEXT,
  CONSTRAINT fkGamesPlayers FOREIGN KEY (gameId) REFERENCES games (id),
  CONSTRAINT fkUsersPlayers FOREIGN KEY (userId) REFERENCES users (id)
);
CREATE TABLE fleets (
  id INTEGER PRIMARY KEY,
  createdAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  updatedAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  gameId INTEGER NOT NULL,
  playerId INTEGER NOT NULL,
  battlePlanName TEXT NOT NULL,
  x REAL,
  y REAL,
  name TEXT NOT NULL,
  num INTEGER,
  playerNum INTEGER,
  waypoints TEXT,
  repeatOrders NUMERIC,
  planetId INTEGER,
  baseName TEXT NOT NULL,
  ironium INTEGER,
  boranium INTEGER,
  germanium INTEGER,
  colonists INTEGER,
  fuel INTEGER,
  damage INTEGER,
  headingX REAL,
  headingY REAL,
  warpSpeed INTEGER,
  previousPositionX REAL,
  previousPositionY REAL,
  orbitingPlanetNum INTEGER,
  starbase NUMERIC,
  spec TEXT
);
CREATE TABLE shipDesigns (
  id INTEGER PRIMARY KEY,
  createdAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  updatedAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  playerId INTEGER NOT NULL,
  playerNum INTEGER,
  uuid TEXT,
  name TEXT NOT NULL,
  version INTEGER,
  hull TEXT,
  hullSetNumber INTEGER,
  canDelete NUMERIC,
  slots TEXT,
  purpose TEXT,
  spec TEXT,
  CONSTRAINT fkPlayersDesigns FOREIGN KEY (playerId) REFERENCES players(id)
);
CREATE TABLE shipTokens (
  id INTEGER PRIMARY KEY,
  createdAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  updatedAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  fleetId INTEGER NOT NULL,
  designUuid TEXT NOT NULL,
  quantity INTEGER,
  damage REAL,
  quantityDamaged INTEGER,
  CONSTRAINT fkFleetsTokens FOREIGN KEY (fleetId) REFERENCES fleets(id)
);
CREATE TABLE planets (
  id INTEGER PRIMARY KEY,
  gameId INTEGER NOT NULL,
  createdAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  updatedAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  playerId INTEGER,
  x REAL,
  y REAL,
  name TEXT NOT NULL,
  num INTEGER,
  playerNum INTEGER,
  grav INTEGER,
  temp INTEGER,
  rad INTEGER,
  baseGrav INTEGER,
  baseTemp INTEGER,
  baseRad INTEGER,
  terraformedAmountGrav INTEGER,
  terraformedAmountTemp INTEGER,
  terraformedAmountRad INTEGER,
  mineralConcIronium INTEGER,
  mineralConcBoranium INTEGER,
  mineralConcGermanium INTEGER,
  mineYearsIronium INTEGER,
  mineYearsBoranium INTEGER,
  mineYearsGermanium INTEGER,
  ironium INTEGER,
  boranium INTEGER,
  germanium INTEGER,
  colonists INTEGER,
  mines INTEGER,
  factories INTEGER,
  defenses INTEGER,
  homeworld NUMERIC,
  contributesOnlyLeftoverToResearch NUMERIC,
  scanner NUMERIC,
  packetSpeed INTEGER,
  productionQueue TEXT,
  spec TEXT,
  CONSTRAINT fkGamesPlanets FOREIGN KEY (gameId) REFERENCES games (id)
);
CREATE TABLE mineralPackets (
  id INTEGER PRIMARY KEY,
  gameId INTEGER NOT NULL,
  createdAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  updatedAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  playerId INTEGER NOT NULL,
  x REAL,
  y REAL,
  name TEXT NOT NULL,
  num INTEGER,
  playerNum INTEGER,
  targetPlanetNum INTEGER,
  cargoIronium INTEGER,
  cargoBoranium INTEGER,
  cargoGermanium INTEGER,
  cargoColonists INTEGER,
  safeWarpSpeed INTEGER,
  warpFactor INTEGER,
  distanceTravelled REAL,
  headingX REAL,
  headingY REAL,
  CONSTRAINT fkGamesMineralPackets FOREIGN KEY (gameId) REFERENCES games (id),
  CONSTRAINT fkPlayersMineralPackets FOREIGN KEY (playerId) REFERENCES players (id)
);
CREATE TABLE salvages (
  id INTEGER PRIMARY KEY,
  gameId INTEGER NOT NULL,
  createdAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  updatedAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  playerId INTEGER,
  x REAL,
  y REAL,
  name TEXT NOT NULL,
  num INTEGER,
  playerNum INTEGER,
  cargoIronium INTEGER,
  cargoBoranium INTEGER,
  cargoGermanium INTEGER,
  cargoColonists INTEGER
);
CREATE TABLE wormoholes (
  id INTEGER PRIMARY KEY,
  gameId INTEGER NOT NULL,
  createdAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  updatedAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  x REAL,
  y REAL,
  name TEXT NOT NULL,
  num INTEGER,
  destinationNum INTEGER,
  stability TEXT,
  yearsAtStability INTEGER
);
CREATE TABLE mineFields (
  id INTEGER PRIMARY KEY,
  gameId INTEGER NOT NULL,
  createdAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  updatedAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  playerId INTEGER NOT NULL,
  x REAL,
  y REAL,
  name TEXT NOT NULL,
  num INTEGER,
  playerNum INTEGER,
  numMines INTEGER,
  detonate NUMERIC
);
CREATE TABLE techStores (
  id INTEGER PRIMARY KEY,
  createdAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  updatedAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  rulesId INTEGER
);
CREATE TABLE techEngines (
  id INTEGER PRIMARY KEY,
  createdAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  updatedAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  techStoreId INTEGER NOT NULL,
  name TEXT NOT NULL,
  costIronium INTEGER,
  costBoranium INTEGER,
  costGermanium INTEGER,
  costResources INTEGER,
  requirementsEnergy INTEGER,
  requirementsWeapons INTEGER,
  requirementsPropulsion INTEGER,
  requirementsConstruction INTEGER,
  requirementsElectronics INTEGER,
  requirementsBiotechnology INTEGER,
  requirementsPrtDenied TEXT,
  requirementsLrTsRequired INTEGER,
  requirementsLrTsDenied INTEGER,
  requirementsPrtRequired TEXT,
  ranking INTEGER,
  category TEXT,
  hullSlotType INTEGER,
  mass INTEGER,
  scanRange INTEGER,
  scanRangePen INTEGER,
  safeHullMass INTEGER,
  safeRange INTEGER,
  maxHullMass INTEGER,
  maxRange INTEGER,
  radiating NUMERIC,
  packetSpeed INTEGER,
  cloakUnits INTEGER,
  terraformRate INTEGER,
  miningRate INTEGER,
  killRate REAL,
  minKillRate INTEGER,
  structureDestroyRate REAL,
  unterraformRate INTEGER,
  smart NUMERIC,
  canStealFleetCargo NUMERIC,
  canStealPlanetCargo NUMERIC,
  armor INTEGER,
  shield INTEGER,
  torpedoBonus REAL,
  initiativeBonus INTEGER,
  beamBonus REAL,
  reduceMovement INTEGER,
  torpedoJamming REAL,
  reduceCloaking NUMERIC,
  cloakUnarmedOnly NUMERIC,
  mineFieldType TEXT,
  mineLayingRate INTEGER,
  beamDefense INTEGER,
  cargoBonus INTEGER,
  colonizationModule NUMERIC,
  fuelBonus INTEGER,
  movementBonus INTEGER,
  orbitalConstructionModule NUMERIC,
  power INTEGER,
  range INTEGER,
  initiative INTEGER,
  gattling NUMERIC,
  hitsAllTargets NUMERIC,
  damageShieldsOnly NUMERIC,
  fuelRegenerationRate INTEGER,
  accuracy INTEGER,
  capitalShipMissile NUMERIC,
  idealSpeed INTEGER,
  freeSpeed INTEGER,
  fuelUsage TEXT,
  CONSTRAINT fkTechStoresEngines FOREIGN KEY (techStoreId) REFERENCES techStores(id)
);
CREATE TABLE techPlanetaryScanners (
  id INTEGER PRIMARY KEY,
  createdAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  updatedAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  techStoreId INTEGER NOT NULL,
  name TEXT NOT NULL,
  costIronium INTEGER,
  costBoranium INTEGER,
  costGermanium INTEGER,
  costResources INTEGER,
  requirementsEnergy INTEGER,
  requirementsWeapons INTEGER,
  requirementsPropulsion INTEGER,
  requirementsConstruction INTEGER,
  requirementsElectronics INTEGER,
  requirementsBiotechnology INTEGER,
  requirementsPrtDenied TEXT,
  requirementsLrTsRequired INTEGER,
  requirementsLrTsDenied INTEGER,
  requirementsPrtRequired TEXT,
  ranking INTEGER,
  category TEXT,
  scanRange INTEGER,
  scanRangePen INTEGER,
  CONSTRAINT fkTechStoresPlanetaryScanners FOREIGN KEY (techStoreId) REFERENCES techStores(id)
);
CREATE TABLE techDefenses (
  id INTEGER PRIMARY KEY,
  createdAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  updatedAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  techStoreId INTEGER NOT NULL,
  name TEXT NOT NULL,
  costIronium INTEGER,
  costBoranium INTEGER,
  costGermanium INTEGER,
  costResources INTEGER,
  requirementsEnergy INTEGER,
  requirementsWeapons INTEGER,
  requirementsPropulsion INTEGER,
  requirementsConstruction INTEGER,
  requirementsElectronics INTEGER,
  requirementsBiotechnology INTEGER,
  requirementsPrtDenied TEXT,
  requirementsLrTsRequired INTEGER,
  requirementsLrTsDenied INTEGER,
  requirementsPrtRequired TEXT,
  ranking INTEGER,
  category TEXT,
  defenseCoverage REAL,
  CONSTRAINT fkTechStoresDefenses FOREIGN KEY (techStoreId) REFERENCES techStores(id)
);
CREATE TABLE techHullComponents (
  id INTEGER PRIMARY KEY,
  createdAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  updatedAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  techStoreId INTEGER NOT NULL,
  name TEXT NOT NULL,
  costIronium INTEGER,
  costBoranium INTEGER,
  costGermanium INTEGER,
  costResources INTEGER,
  requirementsEnergy INTEGER,
  requirementsWeapons INTEGER,
  requirementsPropulsion INTEGER,
  requirementsConstruction INTEGER,
  requirementsElectronics INTEGER,
  requirementsBiotechnology INTEGER,
  requirementsPrtDenied TEXT,
  requirementsLrTsRequired INTEGER,
  requirementsLrTsDenied INTEGER,
  requirementsPrtRequired TEXT,
  ranking INTEGER,
  category TEXT,
  hullSlotType INTEGER,
  mass INTEGER,
  scanRange INTEGER,
  scanRangePen INTEGER,
  safeHullMass INTEGER,
  safeRange INTEGER,
  maxHullMass INTEGER,
  maxRange INTEGER,
  radiating NUMERIC,
  packetSpeed INTEGER,
  cloakUnits INTEGER,
  terraformRate INTEGER,
  miningRate INTEGER,
  killRate REAL,
  minKillRate INTEGER,
  structureDestroyRate REAL,
  unterraformRate INTEGER,
  smart NUMERIC,
  canStealFleetCargo NUMERIC,
  canStealPlanetCargo NUMERIC,
  armor INTEGER,
  shield INTEGER,
  torpedoBonus REAL,
  initiativeBonus INTEGER,
  beamBonus REAL,
  reduceMovement INTEGER,
  torpedoJamming REAL,
  reduceCloaking NUMERIC,
  cloakUnarmedOnly NUMERIC,
  mineFieldType TEXT,
  mineLayingRate INTEGER,
  beamDefense INTEGER,
  cargoBonus INTEGER,
  colonizationModule NUMERIC,
  fuelBonus INTEGER,
  movementBonus INTEGER,
  orbitalConstructionModule NUMERIC,
  power INTEGER,
  range INTEGER,
  initiative INTEGER,
  gattling NUMERIC,
  hitsAllTargets NUMERIC,
  damageShieldsOnly NUMERIC,
  fuelRegenerationRate INTEGER,
  accuracy INTEGER,
  capitalShipMissile NUMERIC,
  CONSTRAINT fkTechStoresHullComponents FOREIGN KEY (techStoreId) REFERENCES techStores(id)
);
CREATE TABLE techHulls (
  id INTEGER PRIMARY KEY,
  createdAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  updatedAt TIMESTAMP NOT NULL DEFAULT CURRENTTIMESTAMP,
  techStoreId INTEGER NOT NULL,
  name TEXT NOT NULL,
  costIronium INTEGER,
  costBoranium INTEGER,
  costGermanium INTEGER,
  costResources INTEGER,
  requirementsEnergy INTEGER,
  requirementsWeapons INTEGER,
  requirementsPropulsion INTEGER,
  requirementsConstruction INTEGER,
  requirementsElectronics INTEGER,
  requirementsBiotechnology INTEGER,
  requirementsPrtDenied TEXT,
  requirementsLrTsRequired INTEGER,
  requirementsLrTsDenied INTEGER,
  requirementsPrtRequired TEXT,
  ranking INTEGER,
  category TEXT,
  type TEXT,
  mass INTEGER,
  armor INTEGER,
  fuelCapacity INTEGER,
  cargoCapacity INTEGER,
  spaceDock INTEGER,
  mineLayingFactor INTEGER,
  builtInScanner NUMERIC,
  initiative INTEGER,
  repairBonus REAL,
  immuneToOwnDetonation NUMERIC,
  rangeBonus INTEGER,
  starbase NUMERIC,
  orbitalConstructionHull NUMERIC,
  doubleMineEfficiency NUMERIC,
  innateScanRangePenFactor REAL,
  slots TEXT,
  CONSTRAINT fkTechStoresHulls FOREIGN KEY (techStoreId) REFERENCES techStores(id)
);