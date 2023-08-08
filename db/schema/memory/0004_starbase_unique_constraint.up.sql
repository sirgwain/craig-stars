-- only one starbase per planet --
CREATE UNIQUE INDEX IF NOT EXISTS fleetStarbasePlanet on fleets(gameId, playerNum, planetNum) WHERE starbase = 1;
