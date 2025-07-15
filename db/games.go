package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db/generated"
)

func (c *client) GetGames(ctx context.Context) ([]cs.Game, error) {

	items, err := c.reader.GetGames(ctx)
	if err == sql.ErrNoRows {
		return []cs.Game{}, nil
	}
	if err != nil {

		return nil, err
	}

	return c.converter.ConvertGames(items), nil
}

func (c *client) GetGamesForHost(ctx context.Context, userID int64) ([]cs.Game, error) {
	items, err := c.reader.GetGamesForHost(ctx, userID)
	if err == sql.ErrNoRows {
		return []cs.Game{}, nil
	}
	if err != nil {

		return nil, err
	}

	return c.converter.ConvertGames(items), nil

}

func (c *client) GetGamesWithPlayers(ctx context.Context) ([]cs.GameWithPlayers, error) {
	return c.getGamesWithPlayersStatus(ctx, generated.GetGamesWithPlayersParams{
		State:  nil,
		Open:   nil,
		Public: nil,
	})
}

func (c *client) GetGamesForUser(ctx context.Context, userID int64) ([]cs.GameWithPlayers, error) {
	items, err := c.reader.GetGamesWithPlayersForUser(ctx, userID)
	if err == sql.ErrNoRows {
		return []cs.GameWithPlayers{}, nil
	}
	if err != nil {
		return nil, err
	}

	games := []cs.GameWithPlayers{}

	var item generated.Game
	var game *cs.GameWithPlayers
	for _, row := range items {

		if row.Game.ID != item.ID {
			// convert this row into a game
			item = row.Game
			g := c.converter.ConvertGame(item)
			games = append(games, cs.GameWithPlayers{Game: g, Players: []cs.PlayerStatus{}})
			game = &games[len(games)-1]
		}

		if row.ID.Valid {
			game.Players = append(game.Players, c.converter.ConvertGetGamesWithPlayersForUserRowToPlayerStatus(row))
		}
	}

	return games, nil
}

func (c *client) GetOpenGames(ctx context.Context) ([]cs.GameWithPlayers, error) {
	return c.getGamesWithPlayersStatus(ctx, generated.GetGamesWithPlayersParams{
		State:  string(cs.GameStateSetup),
		Open:   true,
		Public: true,
	})
}

// get a game by id
func (c *client) GetGame(ctx context.Context, id int64) (*cs.GameWithPlayers, error) {
	return c.getGameWithPlayers(ctx, generated.GetGameWithPlayersParams{
		ID: id,
	})
}

func (c *client) GetGameByHash(ctx context.Context, hash string) (*cs.GameWithPlayers, error) {
	return c.getGameWithPlayers(ctx, generated.GetGameWithPlayersParams{
		Hash: hash,
	})
}

func (c *client) getGameWithPlayers(ctx context.Context, params generated.GetGameWithPlayersParams) (*cs.GameWithPlayers, error) {
	rows, err := c.reader.GetGameWithPlayers(ctx, params)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}

	game := &cs.GameWithPlayers{
		Game: c.converter.ConvertGame(rows[0].Game),
	}

	for _, row := range rows {
		if row.ID.Valid {
			game.Players = append(game.Players, c.converter.ConvertGetGameWithPlayersRowToPlayerStatus(row))
		}
	}

	// TODO: eventually allow rules overrides, but for now, always use standard rules
	game.Rules = cs.NewRules()
	return game, nil
}

func (c *client) getGamesWithPlayersStatus(ctx context.Context, params generated.GetGamesWithPlayersParams) ([]cs.GameWithPlayers, error) {

	items, err := c.reader.GetGamesWithPlayers(ctx, params)
	if err == sql.ErrNoRows {
		return []cs.GameWithPlayers{}, nil
	}
	if err != nil {
		return nil, err
	}

	games := []cs.GameWithPlayers{}

	var item generated.Game
	var game *cs.GameWithPlayers
	for _, row := range items {

		if row.Game.ID != item.ID {
			// convert this row into a game
			item = row.Game
			g := c.converter.ConvertGame(item)
			games = append(games, cs.GameWithPlayers{Game: g, Players: []cs.PlayerStatus{}})
			game = &games[len(games)-1]
		}

		if row.ID.Valid {
			game.Players = append(game.Players, c.converter.ConvertGetGamesWithPlayersRowToPlayerStatus(row))
		}
	}

	return games, nil
}

// get a game by id
func (c *client) GetFullGame(ctx context.Context, id int64) (*cs.FullGame, error) {
	item, err := c.reader.GetGame(ctx, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	game := c.converter.ConvertGame(item)

	players, err := c.GetPlayersForGame(ctx, game.ID)
	if err != nil {
		return nil, fmt.Errorf("load players for game: %w", err)
	}

	designs, err := c.GetShipDesignsForGame(ctx, game.ID)
	if err != nil {
		return nil, fmt.Errorf("load designs for game: %w", err)
	}

	for _, design := range designs {
		player := players[design.PlayerNum-1]
		player.Designs = append(player.Designs, design)
	}

	universeLogger := log.With().Int64("GameID", game.ID).Str("GameName", game.Name).Logger()
	universe := cs.NewUniverse(universeLogger, &game.Rules)

	planets, err := c.GetPlanetsForGame(ctx, game.ID)
	if err != nil {
		return nil, fmt.Errorf("load planets for game: %w", err)
	}
	universe.Planets = planets

	// load fleets and starbases
	fleets, err := c.GetFleetsForGame(ctx, game.ID)
	if err != nil {
		return nil, fmt.Errorf("load fleets for game: %w", err)
	}
	// pre-instantiate the fleets/starbases arrays (make it a little bigger than necessary)
	universe.Fleets = make([]*cs.Fleet, 0, len(fleets))
	universe.Starbases = make([]*cs.Fleet, 0, len(planets))
	for i := range fleets {
		fleet := fleets[i]
		if fleet.Starbase {
			universe.Starbases = append(universe.Starbases, fleet)
		} else {
			universe.Fleets = append(universe.Fleets, fleet)
		}
	}

	wormholes, err := c.GetWormholesForGame(ctx, game.ID)
	if err != nil {
		return nil, fmt.Errorf("load wormholes for game: %w", err)
	}
	universe.Wormholes = wormholes

	salvages, err := c.GetSalvagesForGame(ctx, game.ID)
	if err != nil {
		return nil, fmt.Errorf("load salvages for game: %w", err)
	}
	universe.Salvages = salvages

	mineFields, err := c.getMineFieldsForGame(ctx, game.ID)
	if err != nil {
		return nil, fmt.Errorf("load mineFields for game: %w", err)
	}
	universe.MineFields = mineFields

	mineralPackets, err := c.getMineralPacketsForGame(ctx, game.ID)
	if err != nil {
		return nil, fmt.Errorf("load mineralPackets for game: %w", err)
	}
	universe.MineralPackets = mineralPackets

	mysteryTraders, err := c.GetMysteryTradersForGame(ctx, game.ID)
	if err != nil {
		return nil, fmt.Errorf("load mysteryTraders for game: %w", err)
	}
	universe.MysteryTraders = mysteryTraders

	// TODO: eventually allow rules overrides, but for now, always use standard rules
	game.Rules = cs.NewRules()

	// load a tech store if this game has a separate one
	techStore := &cs.StaticTechStore
	if game.Rules.TechsID != 0 {
		techStore, err = c.GetTechStore(ctx, game.Rules.TechsID)
		if err != nil {
			return nil, err
		}
	}

	game.Rules.SetTechStore(techStore)

	fg := cs.FullGame{
		Game:      &game,
		Universe:  &universe,
		TechStore: techStore,
		Players:   players,
	}

	return &fg, nil
}

// create a new game
func (c *client) CreateGame(ctx context.Context, game *cs.Game) (*cs.Game, error) {

	result, err := c.writer.CreateGame(ctx, c.converter.ConvertGameGameToCreateParams(game))

	if err != nil {
		return nil, err
	}

	game.ID = result
	return game, nil
}

func (c *client) UpdateGameState(ctx context.Context, gameID int64, state cs.GameState) error {
	_, err := c.writer.UpdateGameState(ctx, generated.UpdateGameStateParams{ID: gameID, State: state})
	return err
}

// update an existing game
func (c *client) UpdateGame(ctx context.Context, game *cs.Game) error {
	_, err := c.writer.UpdateGame(ctx, c.converter.ConvertGameGameToUpdateParams(game))
	if err != nil {
		return err
	}

	return nil

}

// Save an entire game in the database. This should always be wrapped in a transaction
// TODO: move this into gameRunner so it's clear it should be wrapped in a transaction?
func (c *client) UpdateFullGame(ctx context.Context, fullGame *cs.FullGame) error {

	if err := c.UpdateGame(ctx, fullGame.Game); err != nil {
		return fmt.Errorf("update game: %w", err)
	}

	for _, player := range fullGame.Players {
		// in case this player is new, set the gameId
		player.GameID = fullGame.ID
		if err := c.SavePlayer(ctx, player); err != nil {
			return fmt.Errorf("update player: %w", err)
		}

		// delete designs
		remainingDesigns := make([]*cs.ShipDesign, 0, len(player.Designs))
		for _, design := range player.Designs {
			if !design.Delete {
				remainingDesigns = append(remainingDesigns, design)
				continue
			}

			// possible an AI created and deleted a design in a single turn so check
			// before we delete a design with no ID
			if design.ID != 0 {
				if err := c.DeleteShipDesign(ctx, design.ID); err != nil {
					return fmt.Errorf("update design: %w", err)
				}
			}
		}
		player.Designs = remainingDesigns

		// save designs
		for _, design := range player.Designs {
			design.GameID = player.GameID
			if err := c.SaveShipDesign(ctx, design); err != nil {
				return fmt.Errorf("update design: %w", err)
			}
		}
	}

	for _, planet := range fullGame.Planets {
		if planet.ID == 0 || planet.Dirty {
			planet.GameID = fullGame.ID
			if err := c.SavePlanet(ctx, planet); err != nil {
				return fmt.Errorf("create planet: %w", err)
			}
		}
	}

	// save fleets
	remainingFleets := make([]*cs.Fleet, 0, len(fullGame.Fleets))

	// first delete fleets. This way if we end up creating a new fleet
	// with an in use unique index, we'll delete the old one first
	for _, fleet := range fullGame.Fleets {
		if !fleet.Delete {
			remainingFleets = append(remainingFleets, fleet)
			continue
		}
		// possible a fleet was created and destroyed in one turn
		if fleet.ID != 0 {
			if err := c.DeleteFleet(ctx, fleet.ID); err != nil {
				return fmt.Errorf("delete fleet: %w", err)
			}
		}
	}
	fullGame.Fleets = remainingFleets

	for _, fleet := range fullGame.Fleets {
		fleet.GameID = fullGame.ID
		if err := c.SaveFleet(ctx, fleet); err != nil {
			return fmt.Errorf("save fleet: %w", err)
		}
	}

	// save fleets
	remainingStarbases := make([]*cs.Fleet, 0, len(fullGame.Starbases))

	// first delete fleets. This way if we end up creating a new fleet
	// with an in use unique index, we'll delete the old one first
	for _, starbase := range fullGame.Starbases {
		if !starbase.Delete {
			remainingStarbases = append(remainingStarbases, starbase)
			continue
		}
		// possible a fleet was created and destroyed in one turn
		if starbase.ID != 0 {
			if err := c.DeleteFleet(ctx, starbase.ID); err != nil {
				return fmt.Errorf("delete fleet: %w", err)
			}
		}
	}
	fullGame.Starbases = remainingStarbases

	for _, starbase := range fullGame.Starbases {
		starbase.GameID = fullGame.ID
		if err := c.SaveFleet(ctx, starbase); err != nil {
			return fmt.Errorf("update starbase: %w", err)
		}
	}

	// lastly update each planet with the starbase it's associated with
	for _, f := range fullGame.Starbases {
		fullGame.Planets[f.PlanetNum-1].Starbase = f
	}

	// save wormholes
	remainingWormholes := make([]*cs.Wormhole, 0, len(fullGame.Wormholes))
	for _, wormhole := range fullGame.Wormholes {
		if !wormhole.Delete {
			remainingWormholes = append(remainingWormholes, wormhole)
			continue
		}
		// possible a wormhole was created and destroyed in one turn
		if wormhole.ID != 0 {
			if err := c.DeleteWormhole(ctx, wormhole.ID); err != nil {
				return fmt.Errorf("delete wormhole: %w", err)
			}
		}
	}

	fullGame.Wormholes = remainingWormholes
	for _, wormhole := range fullGame.Wormholes {
		wormhole.GameID = fullGame.ID
		if err := c.SaveWormhole(ctx, wormhole); err != nil {
			return fmt.Errorf("update wormhole: %w", err)
		}
	}

	// save salvages
	remainingSalvages := make([]*cs.Salvage, 0, len(fullGame.Salvages))
	for _, salvage := range fullGame.Salvages {
		if !salvage.Delete {
			remainingSalvages = append(remainingSalvages, salvage)
			continue
		}
		// possible a salvage was created and destroyed in one turn
		if salvage.ID != 0 {
			if err := c.DeleteSalvage(ctx, salvage.ID); err != nil {
				return fmt.Errorf("delete salvage: %w", err)
			}
		}
	}

	fullGame.Salvages = remainingSalvages
	for _, salvage := range fullGame.Salvages {
		salvage.GameID = fullGame.ID
		if err := c.SaveSalvage(ctx, salvage); err != nil {
			return fmt.Errorf("update salvage: %w", err)
		}
	}

	// save mineFields
	remainingMineFields := make([]*cs.MineField, 0, len(fullGame.MineFields))
	for _, mineField := range fullGame.MineFields {
		if !mineField.Delete {
			remainingMineFields = append(remainingMineFields, mineField)
			continue
		}
		// possible a minefield was created and destroyed in one turn
		if mineField.ID != 0 {
			if err := c.DeleteMineField(ctx, mineField.ID); err != nil {
				return fmt.Errorf("delete minefield: %w", err)
			}
		}
	}

	fullGame.MineFields = remainingMineFields
	for _, mineField := range fullGame.MineFields {
		mineField.GameID = fullGame.ID
		if err := c.SaveMinefield(ctx, mineField); err != nil {
			return fmt.Errorf("update minefield: %w", err)
		}
	}

	// save mineralPackets
	remainingMineralPackets := make([]*cs.MineralPacket, 0, len(fullGame.MineralPackets))
	for _, mineralPacket := range fullGame.MineralPackets {
		if !mineralPacket.Delete {
			remainingMineralPackets = append(remainingMineralPackets, mineralPacket)
			continue
		}
		// possible a mineralPacket was created and destroyed in one turn
		if mineralPacket.ID != 0 {
			if err := c.DeleteMineralPacket(ctx, mineralPacket.ID); err != nil {
				return fmt.Errorf("delete mineralpacket: %w", err)
			}
		}
	}

	fullGame.MineralPackets = remainingMineralPackets
	for _, mineralPacket := range fullGame.MineralPackets {
		mineralPacket.GameID = fullGame.ID
		if err := c.SaveMineralPacket(ctx, mineralPacket); err != nil {
			return fmt.Errorf("update mineralpacket: %w", err)
		}
	}

	// save mysteryTraders
	remainingMysteryTraders := make([]*cs.MysteryTrader, 0, len(fullGame.MysteryTraders))
	for _, mysteryTrader := range fullGame.MysteryTraders {
		if !mysteryTrader.Delete {
			remainingMysteryTraders = append(remainingMysteryTraders, mysteryTrader)
			continue
		}
		// possible a mysteryTrader was created and destroyed in one turn
		if mysteryTrader.ID != 0 {
			if err := c.DeleteMysteryTrader(ctx, mysteryTrader.ID); err != nil {
				return fmt.Errorf("delete mysterytrader: %w", err)
			}
		}
	}

	fullGame.MysteryTraders = remainingMysteryTraders
	for _, mysteryTrader := range fullGame.MysteryTraders {
		mysteryTrader.GameID = fullGame.ID
		if err := c.SaveMysteryTrader(ctx, mysteryTrader); err != nil {
			return fmt.Errorf("update mysterytrader: %w", err)
		}
	}
	return nil

}

func (c *client) UpdateGameHost(ctx context.Context, gameID, hostID int64) error {
	_, err := c.writer.UpdateGameHost(ctx, generated.UpdateGameHostParams{
		ID:     gameID,
		HostID: hostID,
	})
	return err
}

// delete a game by id
func (c *client) DeleteGame(ctx context.Context, id int64) error {
	_, err := c.writer.DeleteGame(ctx, id)
	return err
}

// delete all games for user
func (c *client) DeleteUserGames(ctx context.Context, hostID int64) error {
	_, err := c.writer.DeleteUserGames(ctx, hostID)
	return err
}
