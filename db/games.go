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

func (c *client) GetGamesWithPlayers(ctx context.Context) ([]cs.GameWithPlayers, error) {
	return c.getGameWithPlayersStatus(ctx, generated.GetGamesWithPlayersParams{
		ID:     nil,
		HostId: nil,
		UserId: nil,
		State:  nil,
		Open:   nil,
		Public: nil,
		Hash:   nil,
	})
}

func (c *client) GetGamesForHost(ctx context.Context, userID int64) ([]cs.GameWithPlayers, error) {
	return c.getGameWithPlayersStatus(ctx, generated.GetGamesWithPlayersParams{
		ID:     nil,
		HostId: userID,
		UserId: nil,
		State:  nil,
		Open:   nil,
		Public: nil,
		Hash:   nil,
	})
}

func (c *client) GetGamesForUser(ctx context.Context, userID int64) ([]cs.GameWithPlayers, error) {
	return c.getGameWithPlayersStatus(ctx, generated.GetGamesWithPlayersParams{
		ID:     nil,
		HostId: userID,
		UserId: userID,
		State:  nil,
		Open:   nil,
		Public: nil,
		Hash:   nil,
	})
}

func (c *client) GetOpenGames(ctx context.Context) ([]cs.GameWithPlayers, error) {
	return c.getGameWithPlayersStatus(ctx, generated.GetGamesWithPlayersParams{
		ID:     nil,
		HostId: nil,
		UserId: nil,
		State:  string(cs.GameStateSetup),
		Open:   true,
		Public: true,
		Hash:   nil,
	})
}

func (c *client) GetOpenGamesByHash(ctx context.Context, hash string) ([]cs.GameWithPlayers, error) {
	return c.getGameWithPlayersStatus(ctx, generated.GetGamesWithPlayersParams{
		ID:     nil,
		HostId: nil,
		UserId: nil,
		Public: nil,
		State:  string(cs.GameStateSetup),
		Open:   true,
		Hash:   hash,
	})
}

// get a game by id
func (c *client) GetGame(ctx context.Context, id int64) (*cs.GameWithPlayers, error) {
	games, err := c.getGameWithPlayersStatus(ctx, generated.GetGamesWithPlayersParams{
		ID:     id,
		HostId: nil,
		UserId: nil,
		State:  nil,
		Open:   nil,
		Public: nil,
		Hash:   nil,
	})
	if err != nil {
		return nil, err
	}
	if len(games) == 0 {
		return nil, nil
	}

	return &games[0], nil
}

func (c *client) getGameWithPlayersStatus(ctx context.Context, params generated.GetGamesWithPlayersParams) ([]cs.GameWithPlayers, error) {

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

		if row.PlayerNum.Valid {
			game.Players = append(game.Players, cs.PlayerStatus{
				UpdatedAt:     &row.PlayerUpdatedat.Time,
				UserID:        row.PlayerUserid.Int64,
				Name:          row.PlayerName.String,
				Num:           int(row.PlayerNum.Int64),
				Ready:         row.PlayerReady.Bool,
				AIControlled:  row.PlayerAicontrolled.Bool,
				AIDifficulty:  *row.PlayerAidifficulty,
				Guest:         row.PlayerGuest.Bool,
				SubmittedTurn: row.PlayerSubmittedturn.Bool,
				Color:         row.PlayerColor.String,
				Victor:        row.PlayerVictor.Bool,
				Archived:      row.PlayerArchived.Bool,
			})
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

	players, err := c.getPlayersForGame(ctx, game.ID)
	if err != nil {
		return nil, fmt.Errorf("load players for game: %w", err)
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

	// TODO: allow rules overrides, but for now, always use standard rules
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

	created := c.converter.ConvertGame(result)
	return &created, nil
}

func (c *client) UpdateGameState(ctx context.Context, gameID int64, state cs.GameState) error {
	return c.writer.UpdateGameState(ctx, generated.UpdateGameStateParams{ID: gameID, State: state})
}

// update an existing game
func (c *client) UpdateGame(ctx context.Context, game *cs.Game) error {
	result, err := c.writer.UpdateGame(ctx, c.converter.ConvertGameGameToUpdateParams(game))
	if err != nil {
		return err
	}

	game.UpdatedAt = result.Updatedat
	return nil

}

// Save an entire game in the database. This should always be wrapped in a transaction
// TODO: move this into gameRunner so it's clear it should be wrapped in a transaction?
func (c *client) UpdateFullGame(ctx context.Context, fullGame *cs.FullGame) error {

	if err := c.UpdateGame(ctx, fullGame.Game); err != nil {
		return fmt.Errorf("update game: %w", err)
	}

	for i, player := range fullGame.Players {
		if player.ID == 0 {
			player.GameID = fullGame.ID
			created, err := c.CreatePlayer(ctx, player)
			if err != nil {
				return fmt.Errorf("create player: %w", err)
			}
			// copy designs over, they are part of a different table
			// TODO: move designs out of player
			created.Designs = player.Designs
			fullGame.Players[i] = created
			player = created
		}
		if err := c.updateFullPlayer(ctx, player); err != nil {
			return fmt.Errorf("update player: %w", err)
		}
	}

	for i, planet := range fullGame.Planets {
		if planet.ID == 0 {
			planet.GameID = fullGame.ID
			created, err := c.CreatePlanet(ctx, planet)
			if err != nil {
				return fmt.Errorf("create planet: %w", err)
			}
			// copy the startbase over, it's part of a different table
			// TODO: this is messy, what if we create the starbase below???
			created.Starbase = planet.Starbase
			fullGame.Planets[i] = created
			// log.Debug().Int64("GameID", planet.GameID).Int64("ID", planet.ID).Msgf("Created planet %s", planet.Name)
		} else if planet.Dirty {
			if err := c.UpdatePlanet(ctx, planet); err != nil {
				return fmt.Errorf("update planet: %w", err)
			}
			// log.Debug().Int64("GameID", planet.GameID).Int64("ID", planet.ID).Msgf("Updated planet %s", planet.Name)
		}
	}

	// save fleets and starbases
	remainingFleets := make([]*cs.Fleet, 0, len(fullGame.Fleets))

	// first delete fleets. This way if we end up creating a new fleet
	// with an in use unique index, we'll delete the old one first
	for _, fleet := range append(fullGame.Fleets, fullGame.Starbases...) {
		if fleet.Delete {
			if err := c.DeleteFleet(ctx, fleet.ID); err != nil {
				return fmt.Errorf("delete fleet: %w", err)
			}
			// log.Debug().Int64("GameID", fleet.GameID).Int64("ID", fleet.ID).Msgf("Deleted fleet %s", fleet.Name)
		}
	}

	for _, fleet := range append(fullGame.Fleets, fullGame.Starbases...) {
		if fleet.ID == 0 && !fleet.Delete {
			fleet.GameID = fullGame.ID
			fleet, err := c.CreateFleet(ctx, fleet)
			if err != nil {
				return fmt.Errorf("create fleet: %w", err)
			}
			remainingFleets = append(remainingFleets, fleet)
			// log.Debug().Int64("GameID", fleet.GameID).Int64("ID", fleet.ID).Msgf("Created fleet %s", fleet.Name)
		} else if !fleet.Delete {
			if err := c.UpdateFleet(ctx, fleet); err != nil {
				return fmt.Errorf("update fleet: %w", err)
			}
			remainingFleets = append(remainingFleets, fleet)
			// log.Debug().Int64("GameID", fleet.GameID).Int64("ID", fleet.ID).Msgf("Updated fleet %s", fleet.Name)
		}
	}
	fullGame.Fleets = remainingFleets
	for _, f := range fullGame.Fleets {
		if f.PlanetNum != cs.None {
			fullGame.Planets[f.PlanetNum-1].Starbase = f
		}
	}

	// save wormholes
	for i, wormhole := range fullGame.Wormholes {
		if wormhole.ID == 0 {
			wormhole.GameID = fullGame.ID
			wormhole, err := c.CreateWormhole(ctx, wormhole)
			if err != nil {
				return fmt.Errorf("create wormhole: %w", err)
			}
			fullGame.Wormholes[i] = wormhole
			// log.Debug().Int64("GameID", wormhole.GameID).Int64("ID", wormhole.ID).Msgf("Created wormhole %v", wormhole)
		} else if wormhole.Delete {
			if err := c.DeleteWormhole(ctx, wormhole.ID); err != nil {
				return fmt.Errorf("delete wormhole: %w", err)
			}
			// log.Debug().Int64("GameID", wormhole.GameID).Int64("ID", wormhole.ID).Msgf("Deleted wormhole %s", wormhole.Name)
		} else {
			if err := c.UpdateWormhole(ctx, wormhole); err != nil {
				return fmt.Errorf("update wormhole: %w", err)
			}
			// log.Debug().Int64("GameID", wormhole.GameID).Int64("ID", wormhole.ID).Msgf("Updated wormhole %v", wormhole)
		}
	}

	// save salvages
	for i, salvage := range fullGame.Salvages {
		if salvage.ID == 0 {
			salvage.GameID = fullGame.ID
			salvage, err := c.CreateSalvage(ctx, salvage)
			if err != nil {
				return fmt.Errorf("create wormhole: %w", err)
			}
			fullGame.Salvages[i] = salvage
			// log.Debug().Int64("GameID", salvage.GameID).Int64("ID", salvage.ID).Msgf("Created salvage %s", salvage.Name)
		} else if salvage.Delete {
			if err := c.DeleteSalvage(ctx, salvage.ID); err != nil {
				return fmt.Errorf("delete salvage: %w", err)
			}
			// log.Debug().Int64("GameID", salvage.GameID).Int64("ID", salvage.ID).Msgf("Deleted salvage %s", salvage.Name)
		} else {
			if err := c.UpdateSalvage(ctx, salvage); err != nil {
				return fmt.Errorf("update salvage: %w", err)
			}
			// log.Debug().Int64("GameID", salvage.GameID).Int64("ID", salvage.ID).Msgf("Updated salvage %s", salvage.Name)
		}
	}

	// save mineFields
	for i, mineField := range fullGame.MineFields {
		if mineField.ID == 0 {
			mineField.GameID = fullGame.ID
			mineField, err := c.CreateMineField(ctx, mineField)
			if err != nil {
				return fmt.Errorf("create mineField: %w", err)
			}
			fullGame.MineFields[i] = mineField
			// log.Debug().Int64("GameID", mineField.GameID).Int64("ID", mineField.ID).Msgf("Created mineField %s", mineField.Name)
		} else if mineField.Delete {
			if err := c.DeleteMineField(ctx, mineField.ID); err != nil {
				return fmt.Errorf("delete mineField: %w", err)
			}
			// log.Debug().Int64("GameID", mineField.GameID).Int64("ID", mineField.ID).Msgf("Deleted mineField %s", mineField.Name)
		} else {
			if err := c.UpdateMineField(ctx, mineField); err != nil {
				return fmt.Errorf("update mineField: %w", err)
			}
			// log.Debug().Int64("GameID", mineField.GameID).Int64("ID", mineField.ID).Msgf("Updated mineField %s", mineField.Name)
		}
	}

	// save mineralPackets
	for i, mineralPacket := range fullGame.MineralPackets {
		if mineralPacket.ID == 0 {
			mineralPacket.GameID = fullGame.ID
			mineralPacket, err := c.CreateMineralPacket(ctx, mineralPacket)
			if err != nil {
				return fmt.Errorf("create mineralPacket: %w", err)
			}
			fullGame.MineralPackets[i] = mineralPacket
			// log.Debug().Int64("GameID", mineralPacket.GameID).Int64("ID", mineralPacket.ID).Msgf("Created mineralPacket %s", mineralPacket.Name)
		} else if mineralPacket.Delete {
			if err := c.DeleteMineralPacket(ctx, mineralPacket.ID); err != nil {
				return fmt.Errorf("delete mineralPacket: %w", err)
			}
			// log.Debug().Int64("GameID", mineralPacket.GameID).Int64("ID", mineralPacket.ID).Msgf("Deleted mineralPacket %s", mineralPacket.Name)
		} else {
			if err := c.UpdateMineralPacket(ctx, mineralPacket); err != nil {
				return fmt.Errorf("update mineralPacket: %w", err)
			}
			// log.Debug().Int64("GameID", mineralPacket.GameID).Int64("ID", mineralPacket.ID).Msgf("Updated mineralPacket %s", mineralPacket.Name)
		}
	}

	// save mysteryTraders
	for i, mysteryTrader := range fullGame.MysteryTraders {
		if mysteryTrader.ID == 0 {
			mysteryTrader.GameID = fullGame.ID
			mysteryTrader, err := c.CreateMysteryTrader(ctx, mysteryTrader)
			if err != nil {
				return fmt.Errorf("create mysteryTrader: %w", err)
			}
			fullGame.MysteryTraders[i] = mysteryTrader
			// log.Debug().Int64("GameID", mysteryTrader.GameID).Int64("ID", mysteryTrader.ID).Msgf("Created mysteryTrader %s", mysteryTrader.Name)
		} else if mysteryTrader.Delete {
			if err := c.DeleteMysteryTrader(ctx, mysteryTrader.ID); err != nil {
				return fmt.Errorf("delete mysteryTrader: %w", err)
			}
			// log.Debug().Int64("GameID", mysteryTrader.GameID).Int64("ID", mysteryTrader.ID).Msgf("Deleted mysteryTrader %s", mysteryTrader.Name)
		} else {
			if err := c.UpdateMysteryTrader(ctx, mysteryTrader); err != nil {
				return fmt.Errorf("update mysteryTrader: %w", err)
			}
			// log.Debug().Int64("GameID", mysteryTrader.GameID).Int64("ID", mysteryTrader.ID).Msgf("Updated mysteryTrader %s", mysteryTrader.Name)
		}
	}
	return nil

}

// update a player and their designs
func (c *client) updateFullPlayer(ctx context.Context, player *cs.Player) error {

	if err := c.UpdatePlayer(ctx, player); err != nil {
		return fmt.Errorf("update player: %w", err)
	}

	for i := range player.Designs {
		design := player.Designs[i]
		if design.ID == 0 && !design.Delete {
			design.GameID = player.GameID
			created, err := c.CreateShipDesign(ctx, design)
			if err != nil {
				return fmt.Errorf("create design: %w", err)
			}
			player.Designs[i] = created
		} else if !design.Delete {
			if err := c.UpdateShipDesign(ctx, design); err != nil {
				return fmt.Errorf("update design: %w", err)
			}
		}
	}

	return nil
}

func (c *client) UpdateGameHost(ctx context.Context, gameID int64, hostId int64) error {
	return c.writer.UpdateGameHost(ctx, generated.UpdateGameHostParams{
		ID:     gameID,
		Hostid: sql.NullInt64{Valid: true, Int64: hostId},
	})
}

// delete a game by id
func (c *client) DeleteGame(ctx context.Context, id int64) error {
	return c.writer.DeleteGame(ctx, id)
}

// delete all games for user
func (c *client) DeleteUserGames(ctx context.Context, hostID int64) error {
	return c.writer.DeleteUserGames(ctx, sql.NullInt64{Valid: true, Int64: hostID})
}
