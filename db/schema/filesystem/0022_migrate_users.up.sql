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

-- this is taken care of by upgrade.upgrade6() in db/upgrade.go
-- CREATE UNIQUE INDEX idx_users_discord_id_unique ON users (discord_id)
-- WHERE
--     discord_id <> '';