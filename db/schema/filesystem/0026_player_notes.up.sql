CREATE TABLE
  player_notes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    game_id INTEGER NOT NULL DEFAULT 0,
    player_num INTEGER NOT NULL DEFAULT 0,
    target_type TEXT NOT NULL DEFAULT '',
    target_num INTEGER NOT NULL DEFAULT 0,
    target_player_num INTEGER NOT NULL DEFAULT 0,
    x REAL NOT NULL DEFAULT 0,
    y REAL NOT NULL DEFAULT 0,
    text TEXT NOT NULL DEFAULT '',
    show_on_map BOOLEAN NOT NULL DEFAULT 0,
    CONSTRAINT fk_players_player_notes FOREIGN KEY (game_id, player_num) REFERENCES players (game_id, num) ON DELETE CASCADE,
    CONSTRAINT fk_games_player_notes FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE CASCADE
  );

CREATE INDEX idx_player_notes_gameid_playernum ON player_notes (game_id, player_num);

-- a player has at most one note per map object, notes without a target are labels on the map
CREATE UNIQUE INDEX idx_player_notes_target_unique ON player_notes (
  game_id,
  player_num,
  target_type,
  target_num,
  target_player_num
)
WHERE
  target_type <> '';
