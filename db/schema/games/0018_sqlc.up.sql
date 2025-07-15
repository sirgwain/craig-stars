-- Drop some unused tables. We'll create them later if we actually need them
DROP TABLE IF EXISTS rules;

DROP TABLE IF EXISTS techDefenses;

DROP TABLE IF EXISTS techEngines;

DROP TABLE IF EXISTS techHullComponents;

DROP TABLE IF EXISTS techHulls;

DROP TABLE IF EXISTS techPlanetaryScanners;

DROP TABLE IF EXISTS techStores;

-- Migrate data from old table format to new table format
CREATE TABLE
  games_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    host_id INTEGER NOT NULL DEFAULT 0,
    name TEXT NOT NULL DEFAULT '',
    state TEXT NOT NULL DEFAULT '',
    public BOOLEAN NOT NULL DEFAULT 0,
    hash TEXT NOT NULL DEFAULT '',
    size TEXT NOT NULL DEFAULT '',
    density TEXT NOT NULL DEFAULT '',
    player_positions TEXT NOT NULL DEFAULT '',
    random_events BOOLEAN NOT NULL DEFAULT 0,
    computer_players_form_alliances BOOLEAN NOT NULL DEFAULT 0,
    public_player_scores BOOLEAN NOT NULL DEFAULT 0,
    start_mode TEXT NOT NULL DEFAULT '',
    quick_start_turns INTEGER NOT NULL DEFAULT 0,
    open_player_slots INTEGER NOT NULL DEFAULT 0,
    num_players INTEGER NOT NULL DEFAULT 0,
    victory_conditions_conditions INTEGER NOT NULL DEFAULT 0,
    victory_conditions_num_criteria_required INTEGER NOT NULL DEFAULT 0,
    victory_conditions_years_passed INTEGER NOT NULL DEFAULT 0,
    victory_conditions_own_planets INTEGER NOT NULL DEFAULT 0,
    victory_conditions_attain_tech_level INTEGER NOT NULL DEFAULT 0,
    victory_conditions_attain_tech_level_num_fields INTEGER NOT NULL DEFAULT 0,
    victory_conditions_exceeds_score INTEGER NOT NULL DEFAULT 0,
    victory_conditions_exceeds_second_place_score INTEGER NOT NULL DEFAULT 0,
    victory_conditions_production_capacity INTEGER NOT NULL DEFAULT 0,
    victory_conditions_own_capital_ships INTEGER NOT NULL DEFAULT 0,
    victory_conditions_highest_score_after_years INTEGER NOT NULL DEFAULT 0,
    seed INTEGER NOT NULL DEFAULT 0,
    area_x REAL NOT NULL DEFAULT 0,
    area_y REAL NOT NULL DEFAULT 0,
    year INTEGER NOT NULL DEFAULT 0,
    victor_declared BOOLEAN NOT NULL DEFAULT 0,
    max_minerals BOOLEAN NOT NULL DEFAULT 0,
    archived BOOLEAN NOT NULL DEFAULT 0
  );

CREATE TABLE
  players_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    game_id INTEGER NOT NULL DEFAULT 0,
    user_id INTEGER NOT NULL DEFAULT 0,
    name TEXT NOT NULL DEFAULT '',
    num INTEGER NOT NULL DEFAULT 0,
    ready BOOLEAN NOT NULL DEFAULT 0,
    ai_controlled BOOLEAN NOT NULL DEFAULT 0,
    submitted_turn BOOLEAN NOT NULL DEFAULT 0,
    color TEXT NOT NULL DEFAULT '',
    default_hull_set INTEGER NOT NULL DEFAULT 0,
    tech_levels_energy INTEGER NOT NULL DEFAULT 0,
    tech_levels_weapons INTEGER NOT NULL DEFAULT 0,
    tech_levels_propulsion INTEGER NOT NULL DEFAULT 0,
    tech_levels_construction INTEGER NOT NULL DEFAULT 0,
    tech_levels_electronics INTEGER NOT NULL DEFAULT 0,
    tech_levels_biotechnology INTEGER NOT NULL DEFAULT 0,
    tech_levels_spent_energy INTEGER NOT NULL DEFAULT 0,
    tech_levels_spent_weapons INTEGER NOT NULL DEFAULT 0,
    tech_levels_spent_propulsion INTEGER NOT NULL DEFAULT 0,
    tech_levels_spent_construction INTEGER NOT NULL DEFAULT 0,
    tech_levels_spent_electronics INTEGER NOT NULL DEFAULT 0,
    tech_levels_spent_biotechnology INTEGER NOT NULL DEFAULT 0,
    research_amount INTEGER NOT NULL DEFAULT 0,
    research_spent_last_year INTEGER NOT NULL DEFAULT 0,
    next_research_field TEXT NOT NULL DEFAULT '',
    researching TEXT NOT NULL DEFAULT '',
    battle_plans TEXT NOT NULL DEFAULT '',
    production_plans TEXT NOT NULL DEFAULT '',
    transport_plans TEXT NOT NULL DEFAULT '',
    relations TEXT NOT NULL DEFAULT '',
    cargo_transfers TEXT NOT NULL DEFAULT '',
    messages TEXT NOT NULL DEFAULT '',
    battle_records TEXT NOT NULL DEFAULT '',
    player_intels TEXT NOT NULL DEFAULT '',
    score_intels TEXT NOT NULL DEFAULT '',
    planet_intels TEXT NOT NULL DEFAULT '',
    fleet_intels TEXT NOT NULL DEFAULT '',
    ship_design_intels TEXT NOT NULL DEFAULT '',
    mineral_packet_intels TEXT NOT NULL DEFAULT '',
    minefield_intels TEXT NOT NULL DEFAULT '',
    wormhole_intels TEXT NOT NULL DEFAULT '',
    mystery_trader_intels TEXT NOT NULL DEFAULT '',
    salvage_intels TEXT NOT NULL DEFAULT '',
    race TEXT NOT NULL DEFAULT '',
    stats TEXT,
    score_history TEXT NOT NULL DEFAULT '',
    achieved_victory_conditions INTEGER NOT NULL DEFAULT 0,
    victor BOOLEAN NOT NULL DEFAULT 0,
    spec TEXT NOT NULL DEFAULT '',
    guest BOOLEAN NOT NULL DEFAULT 0,
    ai_difficulty TEXT DEFAULT "",
    acquired_techs TEXT NOT NULL DEFAULT '',
    archived BOOLEAN NOT NULL DEFAULT 0,
    UNIQUE (game_id, num),
    CONSTRAINT fk_games_players FOREIGN KEY (game_id) REFERENCES games_new (id) ON DELETE CASCADE
  );

CREATE TABLE
  fleets_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    game_id INTEGER NOT NULL DEFAULT 0,
    battle_plan_num INTEGER NOT NULL DEFAULT 0,
    x REAL NOT NULL DEFAULT 0,
    y REAL NOT NULL DEFAULT 0,
    name TEXT NOT NULL DEFAULT '',
    num INTEGER NOT NULL DEFAULT 0,
    player_num INTEGER NOT NULL DEFAULT 0,
    tokens TEXT NOT NULL DEFAULT '',
    waypoints TEXT NOT NULL DEFAULT '',
    repeat_orders BOOLEAN NOT NULL DEFAULT 0,
    planet_num INTEGER NOT NULL DEFAULT 0,
    base_name TEXT NOT NULL DEFAULT '',
    ironium INTEGER NOT NULL DEFAULT 0,
    boranium INTEGER NOT NULL DEFAULT 0,
    germanium INTEGER NOT NULL DEFAULT 0,
    colonists INTEGER NOT NULL DEFAULT 0,
    fuel INTEGER NOT NULL DEFAULT 0,
    age INTEGER NOT NULL DEFAULT 0,
    heading_x REAL NOT NULL DEFAULT 0,
    heading_y REAL NOT NULL DEFAULT 0,
    warp_speed INTEGER NOT NULL DEFAULT 0,
    previous_position_x REAL,
    previous_position_y REAL,
    orbiting_planet_num INTEGER NOT NULL DEFAULT 0,
    starbase BOOLEAN NOT NULL DEFAULT 0,
    spec TEXT NOT NULL DEFAULT '',
    purpose TEXT NOT NULL DEFAULT '',
    tags TEXT NOT NULL DEFAULT '',
    CONSTRAINT fk_players_fleets FOREIGN KEY (game_id, player_num) REFERENCES players_new (game_id, num) ON DELETE CASCADE,
    CONSTRAINT fk_games_fleets FOREIGN KEY (game_id) REFERENCES games_new (id) ON DELETE CASCADE
  );

CREATE UNIQUE INDEX idx_fleets_game_id_player_num_num ON fleets_new (game_id, player_num, num)
WHERE
  starbase = 0;

CREATE UNIQUE INDEX idx_starbase_game_id_player_num_planet_num ON fleets_new (game_id, player_num, planet_num)
WHERE
  starbase = 1;

CREATE TABLE
  ship_designs_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    game_id INTEGER NOT NULL DEFAULT 0,
    num INTEGER NOT NULL DEFAULT 0,
    player_num INTEGER NOT NULL DEFAULT 0,
    name TEXT NOT NULL DEFAULT '',
    version INTEGER NOT NULL DEFAULT 0,
    hull TEXT NOT NULL DEFAULT '',
    hull_set_number INTEGER NOT NULL DEFAULT 0,
    slots TEXT NOT NULL DEFAULT '',
    purpose TEXT NOT NULL DEFAULT '',
    spec TEXT NOT NULL DEFAULT '',
    cannot_delete BOOLEAN NOT NULL DEFAULT 0,
    original_player_num INTEGER DEFAULT 0,
    mystery_trader BOOLEAN NOT NULL DEFAULT 0,
    UNIQUE (game_id, player_num, num),
    UNIQUE (game_id, player_num, name),
    CONSTRAINT fk_players_designs FOREIGN KEY (game_id, player_num) REFERENCES players_new (game_id, num) ON DELETE CASCADE
  );

CREATE TABLE
  planets_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    game_id INTEGER NOT NULL DEFAULT 0,
    x REAL NOT NULL DEFAULT 0,
    y REAL NOT NULL DEFAULT 0,
    name TEXT NOT NULL NOT NULL DEFAULT '',
    num INTEGER NOT NULL DEFAULT 0,
    player_num INTEGER NOT NULL DEFAULT 0,
    grav INTEGER NOT NULL DEFAULT 0,
    TEMP INTEGER NOT NULL DEFAULT 0,
    rad INTEGER NOT NULL DEFAULT 0,
    base_grav INTEGER NOT NULL DEFAULT 0,
    base_temp INTEGER NOT NULL DEFAULT 0,
    base_rad INTEGER NOT NULL DEFAULT 0,
    terraformed_amount_grav INTEGER NOT NULL DEFAULT 0,
    terraformed_amount_temp INTEGER NOT NULL DEFAULT 0,
    terraformed_amount_rad INTEGER NOT NULL DEFAULT 0,
    mineral_conc_ironium INTEGER NOT NULL DEFAULT 0,
    mineral_conc_boranium INTEGER NOT NULL DEFAULT 0,
    mineral_conc_germanium INTEGER NOT NULL DEFAULT 0,
    mine_years_ironium INTEGER NOT NULL DEFAULT 0,
    mine_years_boranium INTEGER NOT NULL DEFAULT 0,
    mine_years_germanium INTEGER NOT NULL DEFAULT 0,
    ironium INTEGER NOT NULL DEFAULT 0,
    boranium INTEGER NOT NULL DEFAULT 0,
    germanium INTEGER NOT NULL DEFAULT 0,
    colonists INTEGER NOT NULL DEFAULT 0,
    partial_population INTEGER NOT NULL DEFAULT 0,
    mines INTEGER NOT NULL DEFAULT 0,
    factories INTEGER NOT NULL DEFAULT 0,
    defenses INTEGER NOT NULL DEFAULT 0,
    homeworld boolean NOT NULL DEFAULT 0,
    contributes_only_leftover_to_research boolean NOT NULL DEFAULT 0,
    scanner boolean NOT NULL DEFAULT 0,
    route_target_type TEXT NOT NULL DEFAULT '',
    route_target_num INTEGER NOT NULL DEFAULT 0,
    route_target_player_num INTEGER NOT NULL DEFAULT 0,
    packet_target_num INTEGER NOT NULL DEFAULT 0,
    packet_speed INTEGER NOT NULL DEFAULT 0,
    production_queue TEXT NOT NULL DEFAULT '',
    spec TEXT NOT NULL DEFAULT '',
    tags TEXT NOT NULL DEFAULT '',
    random_artifact boolean NOT NULL DEFAULT 0,
    UNIQUE (game_id, num),
    CONSTRAINT fk_games_planets FOREIGN KEY (game_id) REFERENCES games_new (id) ON DELETE CASCADE
  );

CREATE TABLE
  mineral_packets_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    game_id INTEGER NOT NULL DEFAULT 0,
    x REAL NOT NULL DEFAULT 0,
    y REAL NOT NULL DEFAULT 0,
    name TEXT NOT NULL DEFAULT '',
    num INTEGER NOT NULL DEFAULT 0,
    player_num INTEGER NOT NULL DEFAULT 0,
    target_planet_num INTEGER NOT NULL DEFAULT 0,
    ironium INTEGER NOT NULL DEFAULT 0,
    boranium INTEGER NOT NULL DEFAULT 0,
    germanium INTEGER NOT NULL DEFAULT 0,
    safe_warp_speed INTEGER NOT NULL DEFAULT 0,
    warp_speed INTEGER NOT NULL DEFAULT 0,
    scan_range INTEGER NOT NULL DEFAULT 0,
    scan_range_pen INTEGER NOT NULL DEFAULT 0,
    heading_x REAL NOT NULL DEFAULT 0,
    heading_y REAL NOT NULL DEFAULT 0,
    tags TEXT NOT NULL DEFAULT '',
    UNIQUE (game_id, player_num, num),
    CONSTRAINT fk_players_mineral_packets FOREIGN KEY (game_id, player_num) REFERENCES players_new (game_id, num),
    CONSTRAINT fk_games_mineral_packets FOREIGN KEY (game_id) REFERENCES games_new (id) ON DELETE CASCADE
  );

CREATE TABLE
  salvages_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    game_id INTEGER NOT NULL DEFAULT 0,
    x REAL NOT NULL DEFAULT 0,
    y REAL NOT NULL DEFAULT 0,
    name TEXT NOT NULL DEFAULT '',
    num INTEGER NOT NULL DEFAULT 0,
    player_num INTEGER NOT NULL DEFAULT 0,
    ironium INTEGER NOT NULL DEFAULT 0,
    boranium INTEGER NOT NULL DEFAULT 0,
    germanium INTEGER NOT NULL DEFAULT 0,
    tags TEXT NOT NULL DEFAULT '',
    UNIQUE (game_id, num),
    CONSTRAINT fk_players_salvages FOREIGN KEY (game_id, player_num) REFERENCES players_new (game_id, num),
    CONSTRAINT fk_games_salvages FOREIGN KEY (game_id) REFERENCES games_new (id) ON DELETE CASCADE
  );

CREATE TABLE
  wormholes_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    game_id INTEGER NOT NULL DEFAULT 0,
    x REAL NOT NULL DEFAULT 0,
    y REAL NOT NULL DEFAULT 0,
    name TEXT NOT NULL DEFAULT '',
    num INTEGER NOT NULL DEFAULT 0,
    destination_num INTEGER NOT NULL DEFAULT 0,
    stability TEXT NOT NULL DEFAULT '',
    years_at_stability INTEGER NOT NULL DEFAULT 0,
    spec TEXT NOT NULL DEFAULT '',
    tags TEXT NOT NULL DEFAULT '',
    UNIQUE (game_id, num),
    CONSTRAINT fk_games_wormholes FOREIGN KEY (game_id) REFERENCES games_new (id) ON DELETE CASCADE
  );

CREATE TABLE
  mystery_traders_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    game_id INTEGER NOT NULL DEFAULT 0,
    x REAL NOT NULL DEFAULT 0,
    y REAL NOT NULL DEFAULT 0,
    name TEXT NOT NULL DEFAULT '',
    num INTEGER NOT NULL DEFAULT 0,
    heading_x REAL NOT NULL DEFAULT 0,
    heading_y REAL NOT NULL DEFAULT 0,
    warp_speed INTEGER NOT NULL DEFAULT 0,
    spec TEXT NOT NULL DEFAULT '',
    tags TEXT NOT NULL DEFAULT '',
    requested_boon INTEGER NOT NULL DEFAULT 0,
    destination_x REAL NOT NULL DEFAULT 0,
    destination_y REAL NOT NULL DEFAULT 0,
    reward_type TEXT NOT NULL DEFAULT '',
    players_rewarded TEXT NOT NULL DEFAULT '',
    UNIQUE (game_id, num),
    CONSTRAINT fk_games_mystery_traders FOREIGN KEY (game_id) REFERENCES games_new (id) ON DELETE CASCADE
  );

CREATE TABLE
  minefields_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    game_id INTEGER NOT NULL DEFAULT 0,
    x REAL NOT NULL DEFAULT 0,
    y REAL NOT NULL DEFAULT 0,
    name TEXT NOT NULL DEFAULT '',
    num INTEGER NOT NULL DEFAULT 0,
    player_num INTEGER NOT NULL DEFAULT 0,
    num_mines INTEGER NOT NULL DEFAULT 0,
    detonate BOOLEAN NOT NULL DEFAULT 0,
    minefield_type TEXT NOT NULL DEFAULT '',
    spec TEXT NOT NULL DEFAULT '',
    tags TEXT NOT NULL DEFAULT '',
    UNIQUE (game_id, player_num, num),
    CONSTRAINT fk_players_minefields FOREIGN KEY (game_id, player_num) REFERENCES players_new (game_id, num),
    CONSTRAINT fk_games_minefields FOREIGN KEY (game_id) REFERENCES games_new (id) ON DELETE CASCADE
  );

CREATE TABLE
  versions_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CURRENT INTEGER NOT NULL DEFAULT 0
  );

--
-- Copy data
-- starting with games
--
INSERT INTO
  games_new (
    id,
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
    max_minerals,
    archived
  )
SELECT
  id,
  COALESCE(createdAt, CURRENT_TIMESTAMP),
  COALESCE(updatedAt, CURRENT_TIMESTAMP),
  COALESCE(hostId, 0),
  COALESCE(name, ''),
  COALESCE(state, ''),
  COALESCE(public, 0),
  COALESCE(hash, ''),
  COALESCE(size, ''),
  COALESCE(density, ''),
  COALESCE(playerPositions, ''),
  COALESCE(randomEvents, 0),
  COALESCE(computerPlayersFormAlliances, 0),
  COALESCE(publicPlayerScores, 0),
  COALESCE(startMode, ''),
  COALESCE(quickStartTurns, 0),
  COALESCE(openPlayerSlots, 0),
  COALESCE(numPlayers, 0),
  COALESCE(victoryConditionsConditions, 0),
  COALESCE(victoryConditionsNumCriteriaRequired, 0),
  COALESCE(victoryConditionsYearsPassed, 0),
  COALESCE(victoryConditionsOwnPlanets, 0),
  COALESCE(victoryConditionsAttainTechLevel, 0),
  COALESCE(victoryConditionsAttainTechLevelNumFields, 0),
  COALESCE(victoryConditionsExceedsScore, 0),
  COALESCE(victoryConditionsExceedsSecondPlaceScore, 0),
  COALESCE(victoryConditionsProductionCapacity, 0),
  COALESCE(victoryConditionsOwnCapitalShips, 0),
  COALESCE(victoryConditionsHighestScoreAfterYears, 0),
  COALESCE(seed, 0),
  COALESCE(areaX, 0),
  COALESCE(areaY, 0),
  COALESCE(year, 0),
  COALESCE(victorDeclared, 0),
  COALESCE(maxMinerals, 0),
  COALESCE(archived, 0)
FROM
  games;

DROP TABLE games;

ALTER TABLE games_new
RENAME TO games;

--
-- Copy players
--
INSERT INTO
  players_new (
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
    spec,
    guest,
    ai_difficulty,
    acquired_techs,
    archived
  )
SELECT
  id,
  COALESCE(createdAt, CURRENT_TIMESTAMP),
  COALESCE(updatedAt, CURRENT_TIMESTAMP),
  COALESCE(gameId, 0),
  COALESCE(userId, 0),
  COALESCE(name, ''),
  COALESCE(num, 0),
  COALESCE(ready, 0),
  COALESCE(aiControlled, 0),
  COALESCE(submittedTurn, 0),
  COALESCE(color, ''),
  COALESCE(defaultHullSet, 0),
  COALESCE(techLevelsEnergy, 0),
  COALESCE(techLevelsWeapons, 0),
  COALESCE(techLevelsPropulsion, 0),
  COALESCE(techLevelsConstruction, 0),
  COALESCE(techLevelsElectronics, 0),
  COALESCE(techLevelsBiotechnology, 0),
  COALESCE(techLevelsSpentEnergy, 0),
  COALESCE(techLevelsSpentWeapons, 0),
  COALESCE(techLevelsSpentPropulsion, 0),
  COALESCE(techLevelsSpentConstruction, 0),
  COALESCE(techLevelsSpentElectronics, 0),
  COALESCE(techLevelsSpentBiotechnology, 0),
  COALESCE(researchAmount, 0),
  COALESCE(researchSpentLastYear, 0),
  COALESCE(nextResearchField, ''),
  COALESCE(researching, ''),
  COALESCE(battlePlans, ''),
  COALESCE(productionPlans, ''),
  COALESCE(transportPlans, ''),
  COALESCE(relations, ''),
  COALESCE(cargoTransfers, ''),
  COALESCE(messages, ''),
  COALESCE(battleRecords, ''),
  COALESCE(playerIntels, ''),
  COALESCE(scoreIntels, ''),
  COALESCE(planetIntels, ''),
  COALESCE(fleetIntels, ''),
  COALESCE(shipDesignIntels, ''),
  COALESCE(mineralPacketIntels, ''),
  COALESCE(minefieldIntels, ''),
  COALESCE(wormholeIntels, ''),
  COALESCE(mysteryTraderIntels, ''),
  COALESCE(salvageIntels, ''),
  COALESCE(race, ''),
  stats,
  COALESCE(scoreHistory, ''),
  COALESCE(achievedVictoryConditions, 0),
  COALESCE(victor, 0),
  COALESCE(spec, 0),
  COALESCE(guest, 0),
  COALESCE(aiDifficulty, ''),
  COALESCE(acquiredTechs, ''),
  COALESCE(archived, 0)
FROM
  players;

DROP TABLE players;

ALTER TABLE players_new
RENAME TO players;

--
-- Copy fleets
--
INSERT INTO
  fleets_new (
    id,
    created_at,
    updated_at,
    game_id,
    battle_plan_num,
    x,
    y,
    name,
    num,
    player_num,
    tokens,
    waypoints,
    repeat_orders,
    planet_num,
    base_name,
    ironium,
    boranium,
    germanium,
    colonists,
    fuel,
    age,
    heading_x,
    heading_y,
    warp_speed,
    previous_position_x,
    previous_position_y,
    orbiting_planet_num,
    starbase,
    spec,
    purpose,
    tags
  )
SELECT
  id,
  COALESCE(createdAt, CURRENT_TIMESTAMP),
  COALESCE(updatedAt, CURRENT_TIMESTAMP),
  COALESCE(gameId, 0),
  COALESCE(battlePlanNum, 0),
  COALESCE(x, 0),
  COALESCE(y, 0),
  COALESCE(name, ''),
  COALESCE(num, 0),
  COALESCE(playerNum, 0),
  COALESCE(tokens, ''),
  COALESCE(waypoints, ''),
  COALESCE(repeatOrders, 0),
  COALESCE(planetNum, 0),
  COALESCE(baseName, ''),
  COALESCE(ironium, 0),
  COALESCE(boranium, 0),
  COALESCE(germanium, 0),
  COALESCE(colonists, 0),
  COALESCE(fuel, 0),
  COALESCE(age, 0),
  COALESCE(headingX, 0),
  COALESCE(headingY, 0),
  COALESCE(warpSpeed, 0),
  previousPositionX,
  previousPositionY,
  COALESCE(orbitingPlanetNum, 0),
  COALESCE(starbase, 0),
  COALESCE(spec, ''),
  COALESCE(purpose, ''),
  COALESCE(tags, '')
FROM
  fleets;

DROP TABLE fleets;

ALTER TABLE fleets_new
RENAME TO fleets;

--
-- Copy ship_designs
--
INSERT INTO
  ship_designs_new (
    id,
    created_at,
    updated_at,
    game_id,
    num,
    player_num,
    name,
    version,
    hull,
    hull_set_number,
    slots,
    purpose,
    spec,
    cannot_delete,
    original_player_num,
    mystery_trader
  )
SELECT
  id,
  COALESCE(createdAt, CURRENT_TIMESTAMP),
  COALESCE(updatedAt, CURRENT_TIMESTAMP),
  COALESCE(gameId, 0),
  COALESCE(num, 0),
  COALESCE(playerNum, 0),
  COALESCE(name, ''),
  COALESCE(version, ''),
  COALESCE(hull, ''),
  COALESCE(hullSetNumber, 0),
  COALESCE(slots, ''),
  COALESCE(purpose, ''),
  COALESCE(spec, ''),
  COALESCE(cannotDelete, 0),
  COALESCE(originalPlayerNum, 0),
  COALESCE(mysteryTrader, 0)
FROM
  shipDesigns;

DROP TABLE shipDesigns;

ALTER TABLE ship_designs_new
RENAME TO ship_designs;

--
-- Copy planets
--
INSERT INTO
  planets_new (
    id,
    created_at,
    updated_at,
    game_id,
    x,
    y,
    name,
    num,
    player_num,
    grav,
    TEMP,
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
SELECT
  id,
  COALESCE(createdAt, CURRENT_TIMESTAMP),
  COALESCE(updatedAt, CURRENT_TIMESTAMP),
  COALESCE(gameId, 0),
  COALESCE(x, 0),
  COALESCE(y, 0),
  COALESCE(name, ''),
  COALESCE(num, 0),
  COALESCE(playerNum, 0),
  COALESCE(grav, 0),
  COALESCE(TEMP, 0),
  COALESCE(rad, 0),
  COALESCE(baseGrav, 0),
  COALESCE(baseTemp, 0),
  COALESCE(baseRad, 0),
  COALESCE(terraformedAmountGrav, 0),
  COALESCE(terraformedAmountTemp, 0),
  COALESCE(terraformedAmountRad, 0),
  COALESCE(mineralConcIronium, 0),
  COALESCE(mineralConcBoranium, 0),
  COALESCE(mineralConcGermanium, 0),
  COALESCE(mineYearsIronium, 0),
  COALESCE(mineYearsBoranium, 0),
  COALESCE(mineYearsGermanium, 0),
  COALESCE(ironium, 0),
  COALESCE(boranium, 0),
  COALESCE(germanium, 0),
  COALESCE(colonists, 0),
  COALESCE(partialPopulation, 0),
  COALESCE(mines, 0),
  COALESCE(factories, 0),
  COALESCE(defenses, 0),
  COALESCE(homeworld, 0),
  COALESCE(contributesOnlyLeftoverToResearch, 0),
  COALESCE(scanner, 0),
  COALESCE(routeTargetType, ''),
  COALESCE(routeTargetNum, 0),
  COALESCE(routeTargetPlayerNum, 0),
  COALESCE(packetTargetNum, 0),
  COALESCE(packetSpeed, 0),
  COALESCE(productionQueue, ''),
  COALESCE(spec, ''),
  COALESCE(tags, ''),
  COALESCE(randomArtifact, 0)
FROM
  planets;

DROP TABLE planets;

ALTER TABLE planets_new
RENAME TO planets;

--
-- Copy mineral_packets
--
INSERT INTO
  mineral_packets_new (
    id,
    created_at,
    updated_at,
    game_id,
    x,
    y,
    name,
    num,
    player_num,
    target_planet_num,
    ironium,
    boranium,
    germanium,
    safe_warp_speed,
    warp_speed,
    scan_range,
    scan_range_pen,
    heading_x,
    heading_y,
    tags
  )
SELECT
  id,
  COALESCE(createdAt, CURRENT_TIMESTAMP),
  COALESCE(updatedAt, CURRENT_TIMESTAMP),
  COALESCE(gameId, 0),
  COALESCE(x, 0),
  COALESCE(y, 0),
  COALESCE(name, ''),
  COALESCE(num, 0),
  COALESCE(playerNum, 0),
  COALESCE(targetPlanetNum, 0),
  COALESCE(ironium, 0),
  COALESCE(boranium, 0),
  COALESCE(germanium, 0),
  COALESCE(safeWarpSpeed, 0),
  COALESCE(warpSpeed, 0),
  COALESCE(scanRange, 0),
  COALESCE(scanRangePen, 0),
  COALESCE(headingX, 0),
  COALESCE(headingY, 0),
  COALESCE(tags, '')
FROM
  mineralPackets;

DROP TABLE mineralPackets;

ALTER TABLE mineral_packets_new
RENAME TO mineral_packets;

--
-- Copy salvages
--
INSERT INTO
  salvages_new (
    id,
    created_at,
    updated_at,
    game_id,
    x,
    y,
    name,
    num,
    player_num,
    ironium,
    boranium,
    germanium,
    tags
  )
SELECT
  id,
  COALESCE(createdAt, CURRENT_TIMESTAMP),
  COALESCE(updatedAt, CURRENT_TIMESTAMP),
  COALESCE(gameId, 0),
  COALESCE(x, 0),
  COALESCE(y, 0),
  COALESCE(name, ''),
  COALESCE(num, 0),
  COALESCE(playerNum, 0),
  COALESCE(ironium, 0),
  COALESCE(boranium, 0),
  COALESCE(germanium, 0),
  COALESCE(tags, '')
FROM
  salvages;

DROP TABLE salvages;

ALTER TABLE salvages_new
RENAME TO salvages;

--
-- Copy wormholes
--
INSERT INTO
  wormholes_new (
    id,
    created_at,
    updated_at,
    game_id,
    x,
    y,
    name,
    num,
    destination_num,
    stability,
    years_at_stability,
    spec,
    tags
  )
SELECT
  id,
  COALESCE(createdAt, CURRENT_TIMESTAMP),
  COALESCE(updatedAt, CURRENT_TIMESTAMP),
  COALESCE(gameId, 0),
  COALESCE(x, 0),
  COALESCE(y, 0),
  COALESCE(name, ''),
  COALESCE(num, 0),
  COALESCE(destinationNum, 0),
  COALESCE(stability, ''),
  COALESCE(yearsAtStability, 0),
  COALESCE(spec, ''),
  COALESCE(tags, '')
FROM
  wormholes;

DROP TABLE wormholes;

ALTER TABLE wormholes_new
RENAME TO wormholes;

--
-- Copy mystery_traders
--
INSERT INTO
  mystery_traders_new (
    id,
    created_at,
    updated_at,
    game_id,
    x,
    y,
    name,
    num,
    heading_x,
    heading_y,
    warp_speed,
    spec,
    tags,
    requested_boon,
    destination_x,
    destination_y,
    reward_type,
    players_rewarded
  )
SELECT
  id,
  COALESCE(createdAt, CURRENT_TIMESTAMP),
  COALESCE(updatedAt, CURRENT_TIMESTAMP),
  COALESCE(gameId, 0),
  COALESCE(x, 0),
  COALESCE(y, 0),
  COALESCE(name, ''),
  COALESCE(num, 0),
  COALESCE(headingX, 0),
  COALESCE(headingY, 0),
  COALESCE(warpSpeed, 0),
  COALESCE(spec, ''),
  COALESCE(tags, ''),
  COALESCE(requestedBoon, 0),
  COALESCE(destinationX, 0),
  COALESCE(destinationY, 0),
  COALESCE(rewardType, ''),
  COALESCE(playersRewarded, '')
FROM
  mysteryTraders;

DROP TABLE mysteryTraders;

ALTER TABLE mystery_traders_new
RENAME TO mystery_traders;

--
-- Copy minefields
--
INSERT INTO
  minefields_new (
    id,
    created_at,
    updated_at,
    game_id,
    x,
    y,
    name,
    num,
    player_num,
    num_mines,
    detonate,
    minefield_type,
    spec,
    tags
  )
SELECT
  id,
  COALESCE(createdAt, CURRENT_TIMESTAMP),
  COALESCE(updatedAt, CURRENT_TIMESTAMP),
  COALESCE(gameId, 0),
  COALESCE(x, 0),
  COALESCE(y, 0),
  COALESCE(name, ''),
  COALESCE(num, 0),
  COALESCE(playerNum, 0),
  COALESCE(numMines, 0),
  COALESCE(detonate, 0),
  COALESCE(minefieldType, ''),
  COALESCE(spec, ''),
  COALESCE(tags, '')
FROM
  mineFields;

DROP TABLE mineFields;

ALTER TABLE minefields_new
RENAME TO minefields;

--
-- Copy versions
--
INSERT INTO
  versions_new (id, created_at, updated_at, CURRENT)
SELECT
  id,
  COALESCE(createdAt, CURRENT_TIMESTAMP),
  COALESCE(updatedAt, CURRENT_TIMESTAMP),
  CURRENT
FROM
  versions;

DROP TABLE versions;

ALTER TABLE versions_new
RENAME TO versions;

-- after migration of all the above tables, create some views and indexes we are using now
CREATE VIEW
  game_players AS
SELECT
  players.game_id,
  players.id,
  players.updated_at,
  players.user_id,
  players.name,
  players.num,
  players.ready,
  players.ai_controlled,
  players.ai_difficulty,
  players.submitted_turn,
  players.color,
  players.victor,
  players.archived,
  players.guest
FROM
  players;

CREATE INDEX idx_games_hostid_gameid ON games (host_id);

CREATE INDEX idx_players_userid_gameid ON players (user_id, game_id);

CREATE INDEX idx_players_gameid ON players (game_id);

CREATE INDEX idx_players_gameid_num ON players (game_id, num);

CREATE INDEX idx_planets_gameid ON planets (game_id);

CREATE INDEX idx_planets_gameid_player_num ON planets (game_id, player_num);

CREATE INDEX idx_ship_designs_gameid ON ship_designs (game_id);

CREATE INDEX idx_ship_designs_gameid_playernum ON ship_designs (game_id, player_num);

CREATE INDEX idx_fleets_gameid ON fleets (game_id);

CREATE INDEX idx_fleets_gameid_player_num ON fleets (game_id, player_num);

CREATE INDEX idx_fleets_planet_num ON fleets (game_id, player_num);

CREATE INDEX idx_salvages_gameid ON salvages (game_id);

CREATE INDEX idx_wormholes_gameid ON wormholes (game_id);

CREATE INDEX idx_minefields_gameid ON minefields (game_id);

CREATE INDEX idx_mineral_packets_gameid ON mineral_packets (game_id);

CREATE INDEX idx_mystery_traders_gameid ON mystery_traders (game_id);