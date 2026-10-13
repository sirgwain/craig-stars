CREATE TABLE
  player_mail (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    game_id INTEGER NOT NULL DEFAULT 0,
    year INTEGER NOT NULL DEFAULT 0,
    from_player_num INTEGER NOT NULL DEFAULT 0,
    -- 0 is mail to everybody
    to_player_num INTEGER NOT NULL DEFAULT 0,
    text TEXT NOT NULL DEFAULT '',
    CONSTRAINT fk_players_player_mail FOREIGN KEY (game_id, from_player_num) REFERENCES players (game_id, num) ON DELETE CASCADE,
    CONSTRAINT fk_games_player_mail FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE CASCADE
  );

CREATE INDEX idx_player_mail_gameid_year ON player_mail (game_id, year);
