-- Migrate users table
CREATE TABLE
  users_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    username TEXT UNIQUE NOT NULL DEFAULT '',
    password TEXT NOT NULL DEFAULT '',
    email TEXT NOT NULL DEFAULT '',
    verified BOOLEAN NOT NULL DEFAULT 0,
    banned BOOLEAN NOT NULL DEFAULT 0,
    role TEXT NOT NULL DEFAULT '',
    last_login TIMESTAMP,
    discord_id TEXT NOT NULL DEFAULT '',
    discord_avatar TEXT NOT NULL DEFAULT '',
    discord_webhook_url TEXT NOT NULL DEFAULT '',
    game_id INTEGER NOT NULL DEFAULT 0,
    player_num INT NOT NULL DEFAULT 0
  );

INSERT INTO
  users_new (
    id,
    created_at,
    updated_at,
    username,
    password,
    email,
    verified,
    banned,
    role,
    last_login,
    discord_id,
    discord_avatar,
    discord_webhook_url,
    game_id,
    player_num
  )
SELECT
  id,
  COALESCE(createdAt, CURRENT_TIMESTAMP),
  COALESCE(updatedAt, CURRENT_TIMESTAMP),
  COALESCE(username, ''),
  COALESCE(password, ''),
  COALESCE(email, ''),
  COALESCE(verified, 0),
  COALESCE(banned, 0),
  COALESCE(role, ''),
  lastLogin,
  COALESCE(discordId, ''),
  COALESCE(discordAvatar, ''),
  COALESCE(discordWebhookUrl, ''),
  COALESCE(gameId, 0),
  COALESCE(playerNum, 0)
FROM
  users;

DROP TABLE users;

ALTER TABLE users_new
RENAME TO users;

-- migrate races table
CREATE TABLE
  races_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    plural_name TEXT NOT NULL,
    spend_leftover_points_on TEXT NOT NULL,
    prt TEXT NOT NULL DEFAULT '',
    lrts INTEGER NOT NULL DEFAULT 0,
    hab_low_grav INTEGER NOT NULL DEFAULT 0,
    hab_low_temp INTEGER NOT NULL DEFAULT 0,
    hab_low_rad INTEGER NOT NULL DEFAULT 0,
    hab_high_grav INTEGER NOT NULL DEFAULT 0,
    hab_high_temp INTEGER NOT NULL DEFAULT 0,
    hab_high_rad INTEGER NOT NULL DEFAULT 0,
    growth_rate INTEGER NOT NULL DEFAULT 0,
    pop_efficiency INTEGER NOT NULL DEFAULT 0,
    factory_output INTEGER NOT NULL DEFAULT 0,
    factory_cost INTEGER NOT NULL DEFAULT 0,
    num_factories INTEGER NOT NULL DEFAULT 0,
    factories_cost_less BOOLEAN NOT NULL DEFAULT 0,
    immune_grav BOOLEAN NOT NULL DEFAULT 0,
    immune_temp BOOLEAN NOT NULL DEFAULT 0,
    immune_rad BOOLEAN NOT NULL DEFAULT 0,
    mine_output INTEGER NOT NULL DEFAULT 0,
    mine_cost INTEGER NOT NULL DEFAULT 0,
    num_mines INTEGER NOT NULL DEFAULT 0,
    research_cost_energy TEXT NOT NULL DEFAULT '',
    research_cost_weapons TEXT NOT NULL DEFAULT '',
    research_cost_propulsion TEXT NOT NULL DEFAULT '',
    research_cost_construction TEXT NOT NULL DEFAULT '',
    research_cost_electronics TEXT NOT NULL DEFAULT '',
    research_cost_biotechnology TEXT NOT NULL DEFAULT '',
    techs_start_high BOOLEAN NOT NULL DEFAULT 0,
    spec TEXT NOT NULL DEFAULT '{}',
    CONSTRAINT fk_users_races FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
  );

INSERT INTO
  races_new (
    id,
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
SELECT
  id,
  COALESCE(createdAt, CURRENT_TIMESTAMP),
  COALESCE(updatedAt, CURRENT_TIMESTAMP),
  userId,
  name,
  pluralName,
  spendLeftoverPointsOn,
  COALESCE(prt, ''),
  COALESCE(lrts, 0),
  COALESCE(habLowGrav, 0),
  COALESCE(habLowTemp, 0),
  COALESCE(habLowRad, 0),
  COALESCE(habHighGrav, 0),
  COALESCE(habHighTemp, 0),
  COALESCE(habHighRad, 0),
  COALESCE(growthRate, 0),
  COALESCE(popEfficiency, 0),
  COALESCE(factoryOutput, 0),
  COALESCE(factoryCost, 0),
  COALESCE(numFactories, 0),
  COALESCE(factoriesCostLess, 0),
  COALESCE(immuneGrav, 0),
  COALESCE(immuneTemp, 0),
  COALESCE(immuneRad, 0),
  COALESCE(mineOutput, 0),
  COALESCE(mineCost, 0),
  COALESCE(numMines, 0),
  COALESCE(researchCostEnergy, ''),
  COALESCE(researchCostWeapons, ''),
  COALESCE(researchCostPropulsion, ''),
  COALESCE(researchCostConstruction, ''),
  COALESCE(researchCostElectronics, ''),
  COALESCE(researchCostBiotechnology, ''),
  COALESCE(techsStartHigh, 0),
  REPLACE(COALESCE(spec, '{}'), "ineField", "inefield")
FROM
  races;

DROP TABLE races;

ALTER TABLE races_new
RENAME TO races;