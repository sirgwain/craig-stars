CREATE INDEX idx_players_gameid ON players (gameId);

CREATE INDEX idx_planets_gameid ON planets (gameId);

CREATE INDEX idx_fleets_gameid ON fleets (gameId);

CREATE INDEX idx_minefields_gameid ON mineFields (gameId);

CREATE INDEX idx_mineral_packets_gameid ON mineralPackets (gameId);

CREATE INDEX idx_mystery_traders_gameid ON mysteryTraders (gameId);

CREATE INDEX idx_salvages_gameid ON salvages (gameId);

CREATE INDEX idx_ship_designs_gameid ON shipDesigns (gameId);

CREATE INDEX idx_ship_designs_gameid_playernum ON shipDesigns (gameId, playerNum);

CREATE INDEX idx_fleets_gameid_player_num ON fleets (gameId, playerNum);