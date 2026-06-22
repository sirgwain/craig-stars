CREATE TABLE
  users (
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

CREATE TABLE
  races (
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
    CONSTRAINT fk_users_races FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
  );

CREATE TABLE
  api_tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL DEFAULT '',
    token_prefix TEXT NOT NULL DEFAULT '',
    token_hash TEXT UNIQUE NOT NULL,
    scope TEXT NOT NULL DEFAULT '',
    expires_at TIMESTAMP NOT NULL,
    last_used_at TIMESTAMP,
    revoked_at TIMESTAMP,
    CONSTRAINT fk_users_api_tokens FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
  );

CREATE INDEX idx_api_tokens_token_hash ON api_tokens (token_hash);
CREATE INDEX idx_api_tokens_user_id ON api_tokens (user_id);

CREATE TABLE
  games (
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
    archived BOOLEAN NOT NULL DEFAULT 0,
    galaxy_clumping BOOLEAN NOT NULL DEFAULT 0
  );

CREATE TABLE
  players (
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
    battle_plans TEXT NOT NULL DEFAULT '{}',
    production_plans TEXT NOT NULL DEFAULT '{}',
    transport_plans TEXT NOT NULL DEFAULT '{}',
    relations TEXT NOT NULL DEFAULT '{}',
    cargo_transfers TEXT NOT NULL DEFAULT '{}',
    messages TEXT NOT NULL DEFAULT '{}',
    battle_records TEXT NOT NULL DEFAULT '{}',
    player_intels TEXT NOT NULL DEFAULT '{}',
    score_intels TEXT NOT NULL DEFAULT '{}',
    planet_intels TEXT NOT NULL DEFAULT '{}',
    fleet_intels TEXT NOT NULL DEFAULT '{}',
    ship_design_intels TEXT NOT NULL DEFAULT '{}',
    mineral_packet_intels TEXT NOT NULL DEFAULT '{}',
    minefield_intels TEXT NOT NULL DEFAULT '{}',
    wormhole_intels TEXT NOT NULL DEFAULT '{}',
    mystery_trader_intels TEXT NOT NULL DEFAULT '{}',
    salvage_intels TEXT NOT NULL DEFAULT '{}',
    race TEXT NOT NULL DEFAULT '{}',
    stats TEXT,
    score_history TEXT NOT NULL DEFAULT '{}',
    achieved_victory_conditions INTEGER NOT NULL DEFAULT 0,
    victor BOOLEAN NOT NULL DEFAULT 0,
    guest BOOLEAN NOT NULL DEFAULT 0,
    ai_difficulty TEXT DEFAULT '',
    acquired_techs TEXT NOT NULL DEFAULT '{}',
    archived BOOLEAN NOT NULL DEFAULT 0,
    UNIQUE (game_id, num),
    CONSTRAINT fk_games_players FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE CASCADE
  );

CREATE TABLE
  fleets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    game_id INTEGER NOT NULL DEFAULT 0,
    report_age INTEGER NOT NULL DEFAULT 0,
    battle_plan_num INTEGER NOT NULL DEFAULT 0,
    x REAL NOT NULL DEFAULT 0,
    y REAL NOT NULL DEFAULT 0,
    name TEXT NOT NULL DEFAULT '',
    num INTEGER NOT NULL DEFAULT 0,
    player_num INTEGER NOT NULL DEFAULT 0,
    tokens TEXT NOT NULL DEFAULT '{}',
    waypoints TEXT NOT NULL DEFAULT '{}',
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
    spec TEXT NOT NULL DEFAULT '{}',
    purpose TEXT NOT NULL DEFAULT '',
    tags TEXT NOT NULL DEFAULT '{}',
    CONSTRAINT fk_players_fleets FOREIGN KEY (game_id, player_num) REFERENCES players (game_id, num) ON DELETE CASCADE,
    CONSTRAINT fk_games_fleets FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE CASCADE
  );

CREATE TABLE
  ship_designs (
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
    slots TEXT NOT NULL DEFAULT '{}',
    purpose TEXT NOT NULL DEFAULT '',
    spec TEXT NOT NULL DEFAULT '{}',
    cannot_delete BOOLEAN NOT NULL DEFAULT 0,
    original_player_num INTEGER DEFAULT 0,
    mystery_trader BOOLEAN NOT NULL DEFAULT 0,
    UNIQUE (game_id, player_num, num),
    UNIQUE (game_id, player_num, name),
    CONSTRAINT fk_players_designs FOREIGN KEY (game_id, player_num) REFERENCES players (game_id, num) ON DELETE CASCADE
  );

CREATE TABLE
  planets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    game_id INTEGER NOT NULL DEFAULT 0,
    report_age INTEGER NOT NULL DEFAULT 0,
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
    spec TEXT NOT NULL DEFAULT '{}',
    tags TEXT NOT NULL DEFAULT '{}',
    random_artifact boolean NOT NULL DEFAULT 0,
    UNIQUE (game_id, num),
    CONSTRAINT fk_games_planets FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE CASCADE
  );

CREATE TABLE
  mineral_packets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    game_id INTEGER NOT NULL DEFAULT 0,
    report_age INTEGER NOT NULL DEFAULT 0,
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
    tags TEXT NOT NULL DEFAULT '{}',
    UNIQUE (game_id, player_num, num),
    CONSTRAINT fk_players_mineral_packets FOREIGN KEY (game_id, player_num) REFERENCES players (game_id, num),
    CONSTRAINT fk_games_mineral_packets FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE CASCADE
  );

CREATE TABLE
  salvages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    game_id INTEGER NOT NULL DEFAULT 0,
    report_age INTEGER NOT NULL DEFAULT 0,
    x REAL NOT NULL DEFAULT 0,
    y REAL NOT NULL DEFAULT 0,
    name TEXT NOT NULL DEFAULT '',
    num INTEGER NOT NULL DEFAULT 0,
    player_num INTEGER NOT NULL DEFAULT 0,
    ironium INTEGER NOT NULL DEFAULT 0,
    boranium INTEGER NOT NULL DEFAULT 0,
    germanium INTEGER NOT NULL DEFAULT 0,
    tags TEXT NOT NULL DEFAULT '{}',
    UNIQUE (game_id, num),
    CONSTRAINT fk_players_salvages FOREIGN KEY (game_id, player_num) REFERENCES players (game_id, num),
    CONSTRAINT fk_games_salvages FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE CASCADE
  );

CREATE TABLE
  wormholes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    game_id INTEGER NOT NULL DEFAULT 0,
    report_age INTEGER NOT NULL DEFAULT 0,
    x REAL NOT NULL DEFAULT 0,
    y REAL NOT NULL DEFAULT 0,
    name TEXT NOT NULL DEFAULT '',
    num INTEGER NOT NULL DEFAULT 0,
    destination_num INTEGER NOT NULL DEFAULT 0,
    stability TEXT NOT NULL DEFAULT '',
    years_at_stability INTEGER NOT NULL DEFAULT 0,
    spec TEXT NOT NULL DEFAULT '{}',
    tags TEXT NOT NULL DEFAULT '{}',
    UNIQUE (game_id, num),
    CONSTRAINT fk_games_wormholes FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE CASCADE
  );

CREATE TABLE
  mystery_traders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    game_id INTEGER NOT NULL DEFAULT 0,
    report_age INTEGER NOT NULL DEFAULT 0,
    x REAL NOT NULL DEFAULT 0,
    y REAL NOT NULL DEFAULT 0,
    name TEXT NOT NULL DEFAULT '',
    num INTEGER NOT NULL DEFAULT 0,
    heading_x REAL NOT NULL DEFAULT 0,
    heading_y REAL NOT NULL DEFAULT 0,
    warp_speed INTEGER NOT NULL DEFAULT 0,
    spec TEXT NOT NULL DEFAULT '{}',
    tags TEXT NOT NULL DEFAULT '{}',
    requested_boon INTEGER NOT NULL DEFAULT 0,
    destination_x REAL NOT NULL DEFAULT 0,
    destination_y REAL NOT NULL DEFAULT 0,
    reward_type TEXT NOT NULL DEFAULT '',
    players_rewarded TEXT NOT NULL DEFAULT '{}',
    UNIQUE (game_id, num),
    CONSTRAINT fk_games_mystery_traders FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE CASCADE
  );

CREATE TABLE
  minefields (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    game_id INTEGER NOT NULL DEFAULT 0,
    report_age INTEGER NOT NULL DEFAULT 0,
    x REAL NOT NULL DEFAULT 0,
    y REAL NOT NULL DEFAULT 0,
    name TEXT NOT NULL DEFAULT '',
    num INTEGER NOT NULL DEFAULT 0,
    player_num INTEGER NOT NULL DEFAULT 0,
    num_mines INTEGER NOT NULL DEFAULT 0,
    detonate BOOLEAN NOT NULL DEFAULT 0,
    minefield_type TEXT NOT NULL DEFAULT '',
    tags TEXT NOT NULL DEFAULT '{}',
    UNIQUE (game_id, player_num, num),
    CONSTRAINT fk_players_minefields FOREIGN KEY (game_id, player_num) REFERENCES players (game_id, num),
    CONSTRAINT fk_games_minefields FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE CASCADE
  );

CREATE TABLE
  versions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CURRENT INTEGER NOT NULL DEFAULT 0
  );

CREATE UNIQUE INDEX idx_fleets_game_id_player_num_num ON fleets (game_id, player_num, num)
WHERE
  starbase = 0;

CREATE UNIQUE INDEX idx_starbase_game_id_player_num_planet_num ON fleets (game_id, player_num, planet_num)
WHERE
  starbase = 1;

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

CREATE UNIQUE INDEX idx_users_discord_id_unique ON users (discord_id)
WHERE
  discord_id <> '';

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

INSERT INTO
  versions (CURRENT)
VALUES
  (0);
