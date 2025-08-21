package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db/generated"
)

// Params for getting players either by userId or PlayerNum
type GetPlayerParams struct {
	UserID    int64
	PlayerNum int
}

func (c *client) GetPlayers(ctx context.Context) ([]*cs.Player, error) {
	items, err := c.reader.GetPlayers(ctx)
	if err == sql.ErrNoRows {
		return []*cs.Player{}, nil
	}
	if err != nil {

		return nil, err
	}

	return c.converter.ConvertPlayers(items), nil
}

func (c *client) GetPlayersForUser(ctx context.Context, userID int64) ([]*cs.Player, error) {
	items, err := c.reader.GetPlayersForUser(ctx, userID)
	if err == sql.ErrNoRows {
		return []*cs.Player{}, nil
	}
	if err != nil {

		return nil, err
	}

	return c.converter.ConvertPlayers(items), nil

}

// get all the players for a game, with data loaded
func (c *client) GetPlayersForGame(ctx context.Context, gameID int64) ([]*cs.Player, error) {

	items, err := c.reader.GetPlayersForGame(ctx, gameID)
	if err == sql.ErrNoRows {
		return []*cs.Player{}, nil
	}
	if err != nil {

		return nil, err
	}
	// // players := make([]*cs.Player, 0, len(items))
	players := c.converter.ConvertPlayers(items)

	return players, nil
}

// get all the players for a game, with data loaded
func (c *client) getPlayersWithDesignsForGame(ctx context.Context, gameID int64) ([]*cs.Player, error) {

	items, err := c.reader.GetPlayersWithDesignsForGame(ctx, gameID)
	if err == sql.ErrNoRows {
		return []*cs.Player{}, nil
	}
	if err != nil {

		return nil, err
	}

	players := []*cs.Player{}
	var item generated.Player
	var player *cs.Player
	for _, row := range items {

		if row.Player.ID != item.ID {
			// convert this row into a game
			item = row.Player
			p := c.converter.ConvertPlayer(item)
			players = append(players, &p)
			player = players[len(players)-1]
		}

		if row.ID.Valid {
			player.Designs = append(player.Designs, c.converter.ConvertShipDesign(c.converter.ConvertGetPlayersWithDesignsForGameRowToShipDesign(row)))
		}
	}

	return players, nil
}

// get all the players for a game, with data loaded
func (c *client) GetPlayersStatusForGame(ctx context.Context, gameID int64) ([]*cs.Player, error) {

	items, err := c.reader.GetPlayersStatusForGame(ctx, gameID)
	if err == sql.ErrNoRows {
		return []*cs.Player{}, nil
	}
	if err != nil {

		return nil, err
	}

	return c.converter.ConvertPlayerStatuses(items), nil
}

// get a player by id
func (c *client) GetPlayer(ctx context.Context, id int64) (*cs.Player, error) {
	item, err := c.reader.GetPlayer(ctx, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	player := c.converter.ConvertPlayer(item)
	return &player, nil
}

func (c *client) GetPlayerForGame(ctx context.Context, gameID int64, playerNum int) (*cs.Player, error) {
	rows, err := c.reader.GetPlayerForGame(ctx, generated.GetPlayerForGameParams{
		GameID: gameID,
		Num:    int64(playerNum),
	})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	player := c.converter.ConvertPlayer(rows[0].Player)
	for _, row := range rows {
		if row.ID.Valid {
			player.Designs = append(player.Designs,
				c.converter.ConvertShipDesign(
					c.converter.ConvertGetPlayerForGameRowToShipDesign(row)))
		}
	}

	return &player, nil
}

func (c *client) GetLightPlayerForGame(ctx context.Context, gameID int64, params GetPlayerParams) (*cs.Player, error) {
	queryParams := generated.GetLightPlayerForGameParams{
		GameID:    gameID,
		UserID:    nil,
		PlayerNum: nil,
	}
	if params.UserID != 0 {
		queryParams.UserID = params.UserID
	}
	if params.PlayerNum != 0 {
		queryParams.PlayerNum = params.PlayerNum
	}
	item, err := c.reader.GetLightPlayerForGame(ctx, queryParams)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	player := c.converter.ConvertLightPlayer(item)
	return &player, nil
}

func (c *client) GetLightPlayerForGameWithDesigns(ctx context.Context, gameID int64, params GetPlayerParams) (*cs.Player, error) {
	queryParams := generated.GetLightPlayerForGameParams{
		GameID:    gameID,
		UserID:    nil,
		PlayerNum: nil,
	}
	if params.UserID != 0 {
		queryParams.UserID = params.UserID
	}
	if params.PlayerNum != 0 {
		queryParams.PlayerNum = params.PlayerNum
	}
	item, err := c.reader.GetLightPlayerForGame(ctx, queryParams)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	player := c.converter.ConvertLightPlayer(item)

	// load this player's designs
	designs, err := c.reader.GetShipDesignsForPlayer(ctx, generated.GetShipDesignsForPlayerParams{GameID: gameID, PlayerNum: int64(player.Num)})
	player.Designs = c.converter.ConvertShipDesigns(designs)

	return &player, nil
}

func (c *client) GetPlayerIntel(ctx context.Context, gameID int64, playerNum int) (*cs.Intels, error) {
	intels, err := c.reader.GetPlayerIntel(ctx, generated.GetPlayerIntelParams{GameID: gameID, Num: int64(playerNum)})
	if err != nil {
		return nil, err
	}

	return &cs.Intels{
		BattleRecords:       *intels.BattleRecords,
		PlayerIntels:        *intels.PlayerIntels,
		ScoreIntels:         *intels.ScoreIntels,
		PlanetIntels:        *intels.PlanetIntels,
		FleetIntels:         *intels.FleetIntels,
		ShipDesignIntels:    *intels.ShipDesignIntels,
		MineralPacketIntels: *intels.MineralPacketIntels,
		MinefieldIntels:     *intels.MinefieldIntels,
		WormholeIntels:      *intels.WormholeIntels,
		MysteryTraderIntels: *intels.MysteryTraderIntels,
		SalvageIntels:       *intels.SalvageIntels,
	}, nil

}

func (c *client) GetPlayerMapObjects(ctx context.Context, gameID int64, playerNum int) (*cs.PlayerMapObjects, error) {

	mapObjects := cs.PlayerMapObjects{}

	planets, err := c.GetPlanetsForPlayer(ctx, gameID, playerNum)
	if err != nil {
		return nil, fmt.Errorf("get player planets: %w", err)
	}
	mapObjects.Planets = planets

	minefields, err := c.GetMinefieldsForPlayer(ctx, gameID, playerNum)
	if err != nil {
		return nil, fmt.Errorf("get player minefields: %w", err)
	}
	mapObjects.Minefields = minefields

	mineralPackets, err := c.GetMineralPacketsForPlayer(ctx, gameID, playerNum)
	if err != nil {
		return nil, fmt.Errorf("get player mineralPackets: %w", err)
	}
	mapObjects.MineralPackets = mineralPackets

	fleets, err := c.GetFleetsForPlayer(ctx, gameID, playerNum)
	if err != nil {
		return nil, fmt.Errorf("get player fleets: %w", err)
	}
	mapObjects.Fleets = fleets

	return &mapObjects, nil
}

func (c *client) SavePlayer(ctx context.Context, player *cs.Player) error {
	if player.ID == 0 {
		result, err := c.writer.CreatePlayer(ctx, c.converter.ConvertGamePlayerToCreateParams(player))
		if err != nil {
			return err
		}
		player.ID = result
	} else {
		_, err := c.writer.UpdatePlayer(ctx, c.converter.ConvertGamePlayerToUpdateParams(player))
		if err != nil {
			return err
		}
	}

	return nil
}

// update an existing player's lightweight fields
func (c *client) UpdateLightPlayer(ctx context.Context, player *cs.Player) error {

	_, err := c.writer.UpdateLightPlayer(ctx, generated.UpdateLightPlayerParams{
		ID:                player.ID,
		Name:              player.Name,
		Num:               int64(player.Num),
		Ready:             player.Ready,
		AiControlled:      player.AIControlled,
		AiDifficulty:      &player.AIDifficulty,
		Guest:             player.Guest,
		SubmittedTurn:     player.SubmittedTurn,
		Color:             player.Color,
		DefaultHullSet:    int64(player.DefaultHullSet),
		ResearchAmount:    int64(player.ResearchAmount),
		NextResearchField: player.NextResearchField,
		Researching:       player.Researching,
	})
	return err
}

// update an existing player's lightweight fields
func (c *client) UpdatePlayerOrders(ctx context.Context, player *cs.Player) error {

	_, err := c.writer.UpdatePlayerOrders(ctx, generated.UpdatePlayerOrdersParams{
		ID:                player.ID,
		SubmittedTurn:     player.SubmittedTurn,
		DefaultHullSet:    int64(player.DefaultHullSet),
		ResearchAmount:    int64(player.ResearchAmount),
		NextResearchField: player.NextResearchField,
		Researching:       player.Researching,
		CargoTransfers:    (*generated.CargoTransfers)(&player.CargoTransfers),
		BattlePlans:       (*generated.BattlePlans)(&player.BattlePlans),
		ProductionPlans:   (*generated.ProductionPlans)(&player.ProductionPlans),
		TransportPlans:    (*generated.TransportPlans)(&player.TransportPlans),
		Relations:         (*generated.PlayerRelationships)(&player.Relations),
	})
	return err
}

// update an existing player's lightweight fields
func (c *client) UpdatePlayerCargoTransfers(ctx context.Context, player *cs.Player) error {
	_, err := c.writer.UpdatePlayerCargoTransfers(ctx, generated.UpdatePlayerCargoTransfersParams{
		ID:             player.ID,
		CargoTransfers: (*generated.CargoTransfers)(&player.CargoTransfers),
	})
	return err

}

// update an existing player's lightweight fields
func (c *client) UpdatePlayerRelations(ctx context.Context, player *cs.Player) error {
	_, err := c.writer.UpdatePlayerRelations(ctx, generated.UpdatePlayerRelationsParams{
		ID:        player.ID,
		Relations: (*generated.PlayerRelationships)(&player.Relations),
	})
	return err

}

// update an existing player's lightweight fields
func (c *client) SubmitPlayerTurn(ctx context.Context, gameID int64, num int, submittedTurn bool) error {
	_, err := c.writer.SubmitPlayerTurn(ctx, generated.SubmitPlayerTurnParams{
		GameID:        gameID,
		Num:           int64(num),
		SubmittedTurn: submittedTurn,
	})
	return err
}

// update an existing player's lightweight fields
func (c *client) ArchivePlayer(ctx context.Context, gameID int64, num int, archived bool) error {
	_, err := c.writer.ArchivePlayer(ctx, generated.ArchivePlayerParams{
		GameID:   gameID,
		Num:      int64(num),
		Archived: archived,
	})
	return err
}

// update an existing player's lightweight fields
func (c *client) UpdatePlayerPlans(ctx context.Context, player *cs.Player) error {
	_, err := c.writer.UpdatePlayerPlans(ctx, generated.UpdatePlayerPlansParams{
		ID:              player.ID,
		BattlePlans:     (*generated.BattlePlans)(&player.BattlePlans),
		ProductionPlans: (*generated.ProductionPlans)(&player.ProductionPlans),
		TransportPlans:  (*generated.TransportPlans)(&player.TransportPlans),
	})
	return err
}

// update a players planet intels (used after creating a new planet)
func (c *client) UpdatePlayerPlanetIntels(ctx context.Context, player *cs.Player) error {
	_, err := c.writer.UpdatePlayerPlanetIntels(ctx, generated.UpdatePlayerPlanetIntelsParams{
		ID:           player.ID,
		PlanetIntels: (*generated.PlanetIntels)(&player.PlanetIntels),
	})
	return err
}

// update a players fleet intels (used after creating a new fleet)
func (c *client) UpdatePlayerFleetIntels(ctx context.Context, player *cs.Player) error {
	_, err := c.writer.UpdatePlayerFleetIntels(ctx, generated.UpdatePlayerFleetIntelsParams{
		ID:          player.ID,
		FleetIntels: (*generated.FleetIntels)(&player.FleetIntels),
	})
	return err
}

// update a players salvage intels (used after creating a new salvage)
func (c *client) UpdatePlayerSalvageIntels(ctx context.Context, player *cs.Player) error {
	_, err := c.writer.UpdatePlayerSalvageIntels(ctx, generated.UpdatePlayerSalvageIntelsParams{
		ID:            player.ID,
		SalvageIntels: (*generated.SalvageIntels)(&player.SalvageIntels),
	})
	return err
}

// update a players mineralPacket intels (used after creating a new mineralPacket)
func (c *client) UpdatePlayerMineralPacketIntels(ctx context.Context, player *cs.Player) error {
	_, err := c.writer.UpdatePlayerMineralPacketIntels(ctx, generated.UpdatePlayerMineralPacketIntelsParams{
		ID:                  player.ID,
		MineralPacketIntels: (*generated.MineralPacketIntels)(&player.MineralPacketIntels),
	})
	return err
}

// helper to update a player using a transaction or DB
// update an existing player
func (c *client) UpdatePlayerUserID(ctx context.Context, player *cs.Player) error {
	_, err := c.writer.UpdatePlayerUserID(ctx, generated.UpdatePlayerUserIDParams{
		ID:     player.ID,
		UserID: player.UserID,
	})
	return err
}

// delete a player by id
func (c *client) DeletePlayer(ctx context.Context, id int64) error {
	_, err := c.writer.DeletePlayer(ctx, id)
	return err
}
