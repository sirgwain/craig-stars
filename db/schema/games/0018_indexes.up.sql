CREATE VIEW
  game_players AS
SELECT
  players.gameId,
  players.id,
  players.updatedAt,
  players.userId,
  players.name,
  players.num,
  players.ready,
  players.aiControlled,
  players.aiDifficulty,
  players.submittedTurn,
  players.color,
  players.victor,
  players.archived,
  players.guest
FROM
  players;

CREATE INDEX idx_games_hostid_gameid ON games(hostId);
CREATE INDEX idx_players_userid_gameid ON players(userId, gameId);
CREATE INDEX idx_players_gameid ON players(gameId);
CREATE INDEX idx_players_gameid_num ON players(gameId, num);
CREATE INDEX idx_planets_gameid ON planets(gameId);
CREATE INDEX idx_planets_gameid_player_num ON planets(gameId, playerNum);
CREATE INDEX idx_ship_designs_gameid ON shipDesigns(gameId);
CREATE INDEX idx_ship_designs_gameid_playernum ON shipDesigns(gameId, playerNum);
CREATE INDEX idx_fleets_gameid ON fleets(gameId);
CREATE INDEX idx_fleets_gameid_player_num ON fleets(gameId, playerNum);
CREATE INDEX idx_fleets_planet_num ON fleets(gameId, planetNum);
CREATE INDEX idx_salvages_gameid ON salvages(gameId);
CREATE INDEX idx_wormholes_gameid ON wormholes(gameId);
CREATE INDEX idx_minefields_gameid ON mineFields(gameId);
CREATE INDEX idx_mineral_packets_gameid ON mineralPackets(gameId);
CREATE INDEX idx_mystery_traders_gameid ON mysteryTraders(gameId);
