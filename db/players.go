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
	items, err := c.reader.GetPlayersForUser(ctx, sql.NullInt64{Valid: true, Int64: userID})
	if err == sql.ErrNoRows {
		return []*cs.Player{}, nil
	}
	if err != nil {

		return nil, err
	}

	return c.converter.ConvertPlayers(items), nil

}

// get all the players for a game, with data loaded
func (c *client) getPlayersForGame(ctx context.Context, gameID int64) ([]*cs.Player, error) {

	items, err := c.reader.GetPlayersForGame(ctx, gameID)
	if err == sql.ErrNoRows {
		return []*cs.Player{}, nil
	}
	if err != nil {

		return nil, err
	}
	// // players := make([]*cs.Player, 0, len(items))
	players := c.converter.ConvertPlayers(items)
	for _, player := range players {
		designs, err := c.GetShipDesignsForPlayer(ctx, gameID, player.Num)
		if err != nil {
			return nil, fmt.Errorf("failed to get designs for player: %d %w", player.Num, err)
		}
		player.Designs = designs
	}

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

		if row.DesignID.Valid {
			player.Designs = append(player.Designs, c.converter.ConvertShipDesign(generated.Shipdesign{
				ID:                row.DesignID.Int64,
				Createdat:         row.DesignCreatedat.Time,
				Updatedat:         row.DesignUpdatedat.Time,
				Gameid:            row.DesignGameid.Int64,
				Num:               row.DesignNum.Int64,
				Playernum:         row.DesignPlayernum.Int64,
				Name:              row.DesignName.String,
				Version:           row.DesignVersion,
				Hull:              row.DesignHull,
				Hullsetnumber:     row.DesignHullsetnumber,
				Candelete:         row.DesignCandelete,
				Slots:             row.DesignSlots,
				Purpose:           row.DesignPurpose,
				Spec:              row.DesignSpec,
				Cannotdelete:      row.DesignCannotdelete.Bool,
				Originalplayernum: row.DesignOriginalplayernum,
				Mysterytrader:     row.DesignMysterytrader,
			}))
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

// Get all player data except universe intel
func (c *client) GetPlayerForGame(ctx context.Context, gameID int64, params GetPlayerParams) (*cs.Player, error) {
	queryParams := generated.GetPlayerForGameParams{
		GameId:    gameID,
		UserId:    nil,
		PlayerNum: nil,
	}
	if params.UserID != 0 {
		queryParams.UserId = params.UserID
	}
	if params.PlayerNum != 0 {
		queryParams.PlayerNum = params.PlayerNum
	}
	rows, err := c.reader.GetPlayerForGame(ctx, queryParams)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	player := c.converter.ConvertPlayer(rows[0].Player)
	for _, row := range rows {
		if row.DesignID.Valid {
			player.Designs = append(player.Designs, c.converter.ConvertShipDesign(generated.Shipdesign{
				ID:                row.DesignID.Int64,
				Createdat:         row.DesignCreatedat.Time,
				Updatedat:         row.DesignUpdatedat.Time,
				Gameid:            row.DesignGameid.Int64,
				Num:               row.DesignNum.Int64,
				Playernum:         row.DesignPlayernum.Int64,
				Name:              row.DesignName.String,
				Version:           row.DesignVersion,
				Hull:              row.DesignHull,
				Hullsetnumber:     row.DesignHullsetnumber,
				Candelete:         row.DesignCandelete,
				Slots:             row.DesignSlots,
				Purpose:           row.DesignPurpose,
				Spec:              row.DesignSpec,
				Cannotdelete:      row.DesignCannotdelete.Bool,
				Originalplayernum: row.DesignOriginalplayernum,
				Mysterytrader:     row.DesignMysterytrader,
			}))
		}
	}

	return &player, nil
}

func (c *client) GetLightPlayerForGame(ctx context.Context, gameID int64, params GetPlayerParams) (*cs.Player, error) {
	queryParams := generated.GetLightPlayerForGameParams{
		GameId:    gameID,
		UserId:    nil,
		PlayerNum: nil,
	}
	if params.UserID != 0 {
		queryParams.UserId = params.UserID
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

// get a full player by id with all dependencies loaded
func (c *client) GetFullPlayerForGame(ctx context.Context, gameID int64, params GetPlayerParams) (*cs.FullPlayer, error) {

	p, err := c.GetPlayerForGame(ctx, gameID, params)
	if err != nil {
		return nil, err
	}

	player := cs.FullPlayer{
		Player: *p,
	}

	planets, err := c.GetPlanetsForPlayer(ctx, player.GameID, player.Num)
	if err != nil {
		return nil, fmt.Errorf("get player planets: %w", err)
	}
	player.Planets = planets

	mineFields, err := c.GetMineFieldsForPlayer(ctx, player.GameID, player.Num)
	if err != nil {
		return nil, fmt.Errorf("get player mineFields: %w", err)
	}
	player.MineFields = mineFields

	mineralPackets, err := c.GetMineralPacketsForPlayer(ctx, player.GameID, player.Num)
	if err != nil {
		return nil, fmt.Errorf("get player mineralPackets: %w", err)
	}
	player.MineralPackets = mineralPackets

	fleets, err := c.GetFleetsForPlayer(ctx, player.GameID, player.Num)
	if err != nil {
		return nil, fmt.Errorf("get player fleets: %w", err)
	}

	// pre-instantiate the fleets/starbases arrays (make it a little bigger than necessary)
	player.Fleets = make([]*cs.Fleet, 0, len(fleets))
	player.Starbases = make([]*cs.Fleet, 0, len(planets))
	for i := range fleets {
		fleet := fleets[i]
		if fleet.Starbase {
			player.Starbases = append(player.Starbases, fleet)
		} else {
			player.Fleets = append(player.Fleets, fleet)
		}
	}

	return &player, nil
}

func (c *client) GetPlayerMapObjects(ctx context.Context, gameID, userID int64) (*cs.PlayerMapObjects, error) {
	num, err := c.reader.GetPlayerNum(ctx, generated.GetPlayerNumParams{Gameid: gameID, Userid: sql.NullInt64{Valid: true, Int64: userID}})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	playerNum := int(num)
	mapObjects := cs.PlayerMapObjects{}

	planets, err := c.GetPlanetsForPlayer(ctx, gameID, playerNum)
	if err != nil {
		return nil, fmt.Errorf("get player planets: %w", err)
	}
	mapObjects.Planets = planets

	mineFields, err := c.GetMineFieldsForPlayer(ctx, gameID, playerNum)
	if err != nil {
		return nil, fmt.Errorf("get player mineFields: %w", err)
	}
	mapObjects.MineFields = mineFields

	mineralPackets, err := c.GetMineralPacketsForPlayer(ctx, gameID, playerNum)
	if err != nil {
		return nil, fmt.Errorf("get player mineralPackets: %w", err)
	}
	mapObjects.MineralPackets = mineralPackets

	fleets, err := c.GetFleetsForPlayer(ctx, gameID, playerNum)
	if err != nil {
		return nil, fmt.Errorf("get player fleets: %w", err)
	}
	// pre-instantiate the fleets/starbases arrays (make it a little bigger than necessary)
	mapObjects.Fleets = make([]*cs.Fleet, 0, len(fleets))
	mapObjects.Starbases = make([]*cs.Fleet, 0, len(planets))
	for i := range fleets {
		fleet := fleets[i]
		if fleet.Starbase {
			mapObjects.Starbases = append(mapObjects.Starbases, fleet)
		} else {
			mapObjects.Fleets = append(mapObjects.Fleets, fleet)
		}
	}

	return &mapObjects, nil
}

func (c *client) CreatePlayer(ctx context.Context, player *cs.Player) (*cs.Player, error) {

	result, err := c.writer.CreatePlayer(ctx, c.converter.ConvertGamePlayerToCreateParams(player))
	if err != nil {
		return nil, err
	}

	created := c.converter.ConvertPlayer(result)
	return &created, nil
}

// update an existing player's lightweight fields
func (c *client) UpdateLightPlayer(ctx context.Context, player *cs.Player) error {

	result, err := c.writer.UpdateLightPlayer(ctx, generated.UpdateLightPlayerParams{
		ID:                player.ID,
		Name:              player.Name,
		Num:               int64(player.Num),
		Ready:             sql.NullBool{Valid: true, Bool: player.Ready},
		Aicontrolled:      sql.NullBool{Valid: true, Bool: player.AIControlled},
		Aidifficulty:      &player.AIDifficulty,
		Guest:             player.Guest,
		Submittedturn:     sql.NullBool{Valid: true, Bool: player.SubmittedTurn},
		Color:             sql.NullString{Valid: true, String: player.Color},
		Defaulthullset:    sql.NullInt64{Valid: true, Int64: int64(player.DefaultHullSet)},
		Researchamount:    sql.NullInt64{Valid: true, Int64: int64(player.ResearchAmount)},
		Nextresearchfield: player.NextResearchField,
		Researching:       player.Researching,
		Spec:              (*generated.PlayerSpec)(&player.Spec),
	})
	if err != nil {
		return err
	}

	player.UpdatedAt = result
	return nil
}

// update an existing player's lightweight fields
func (c *client) UpdatePlayerOrders(ctx context.Context, player *cs.Player) error {

	result, err := c.writer.UpdatePlayerOrders(ctx, generated.UpdatePlayerOrdersParams{
		ID:                player.ID,
		Submittedturn:     sql.NullBool{Valid: true, Bool: player.SubmittedTurn},
		Defaulthullset:    sql.NullInt64{Valid: true, Int64: int64(player.DefaultHullSet)},
		Researchamount:    sql.NullInt64{Valid: true, Int64: int64(player.ResearchAmount)},
		Nextresearchfield: player.NextResearchField,
		Researching:       player.Researching,
		Cargotransfers:    (*generated.CargoTransfers)(&player.CargoTransfers),
		Battleplans:       (*generated.BattlePlans)(&player.BattlePlans),
		Productionplans:   (*generated.ProductionPlans)(&player.ProductionPlans),
		Transportplans:    (*generated.TransportPlans)(&player.TransportPlans),
		Relations:         (*generated.PlayerRelationships)(&player.Relations),
		Spec:              (*generated.PlayerSpec)(&player.Spec),
	})
	if err != nil {
		return err
	}

	player.UpdatedAt = result
	return nil
}

// update an existing player's lightweight fields
func (c *client) UpdatePlayerCargoTransfers(ctx context.Context, player *cs.Player) error {
	result, err := c.writer.UpdatePlayerCargoTransfers(ctx, generated.UpdatePlayerCargoTransfersParams{
		ID:             player.ID,
		Cargotransfers: (*generated.CargoTransfers)(&player.CargoTransfers),
	})
	if err != nil {
		return err
	}

	player.UpdatedAt = result
	return nil

}

// update an existing player's lightweight fields
func (c *client) UpdatePlayerRelations(ctx context.Context, player *cs.Player) error {
	result, err := c.writer.UpdatePlayerRelations(ctx, generated.UpdatePlayerRelationsParams{
		ID:        player.ID,
		Relations: (*generated.PlayerRelationships)(&player.Relations),
	})
	if err != nil {
		return err
	}

	player.UpdatedAt = result
	return nil

}

// update an existing player's lightweight fields
func (c *client) SubmitPlayerTurn(ctx context.Context, gameID int64, num int, submittedTurn bool) error {
	_, err := c.writer.SubmitPlayerTurn(ctx, generated.SubmitPlayerTurnParams{
		Gameid:        gameID,
		Num:           int64(num),
		Submittedturn: sql.NullBool{Valid: true, Bool: submittedTurn},
	})
	if err != nil {
		return err
	}

	return nil
}

// update an existing player's lightweight fields
func (c *client) ArchivePlayer(ctx context.Context, gameID int64, num int, archived bool) error {
	_, err := c.writer.ArchivePlayer(ctx, generated.ArchivePlayerParams{
		Gameid:   gameID,
		Num:      int64(num),
		Archived: archived,
	})
	if err != nil {
		return err
	}

	return nil
}

// update an existing player's lightweight fields
func (c *client) UpdatePlayerPlans(ctx context.Context, player *cs.Player) error {
	result, err := c.writer.UpdatePlayerPlans(ctx, generated.UpdatePlayerPlansParams{
		ID:              player.ID,
		Battleplans:     (*generated.BattlePlans)(&player.BattlePlans),
		Productionplans: (*generated.ProductionPlans)(&player.ProductionPlans),
		Transportplans:  (*generated.TransportPlans)(&player.TransportPlans),
	})
	if err != nil {
		return err
	}

	player.UpdatedAt = result
	return nil
}

// update a player's spec in the database
func (c *client) UpdatePlayerSpec(ctx context.Context, player *cs.Player) error {
	result, err := c.writer.UpdatePlayerSpec(ctx, generated.UpdatePlayerSpecParams{
		ID:   player.ID,
		Spec: (*generated.PlayerSpec)(&player.Spec),
	})
	if err != nil {
		return err
	}

	player.UpdatedAt = result
	return nil
}

// update a players planet intels (used after creating a new planet)
func (c *client) UpdatePlayerPlanetIntels(ctx context.Context, player *cs.Player) error {
	result, err := c.writer.UpdatePlayerPlanetIntels(ctx, generated.UpdatePlayerPlanetIntelsParams{
		ID:           player.ID,
		Planetintels: (*generated.PlanetIntels)(&player.PlanetIntels),
	})
	if err != nil {
		return err
	}

	player.UpdatedAt = result
	return nil
}

// update a players fleet intels (used after creating a new fleet)
func (c *client) UpdatePlayerFleetIntels(ctx context.Context, player *cs.Player) error {
	result, err := c.writer.UpdatePlayerFleetIntels(ctx, generated.UpdatePlayerFleetIntelsParams{
		ID:          player.ID,
		Fleetintels: (*generated.FleetIntels)(&player.FleetIntels),
	})
	if err != nil {
		return err
	}

	player.UpdatedAt = result
	return nil
}

// update a players salvage intels (used after creating a new salvage)
func (c *client) UpdatePlayerSalvageIntels(ctx context.Context, player *cs.Player) error {
	result, err := c.writer.UpdatePlayerSalvageIntels(ctx, generated.UpdatePlayerSalvageIntelsParams{
		ID:            player.ID,
		Salvageintels: (*generated.SalvageIntels)(&player.SalvageIntels),
	})
	if err != nil {
		return err
	}

	player.UpdatedAt = result
	return nil
}

// update a players mineralPacket intels (used after creating a new mineralPacket)
func (c *client) UpdatePlayerMineralPacketIntels(ctx context.Context, player *cs.Player) error {
	result, err := c.writer.UpdatePlayerMineralPacketIntels(ctx, generated.UpdatePlayerMineralPacketIntelsParams{
		ID:                  player.ID,
		Mineralpacketintels: (*generated.MineralPacketIntels)(&player.MineralPacketIntels),
	})
	if err != nil {
		return err
	}

	player.UpdatedAt = result
	return nil
}

// helper to update a player using a transaction or DB
// update an existing player
func (c *client) UpdatePlayer(ctx context.Context, player *cs.Player) error {
	result, err := c.writer.UpdatePlayer(ctx, c.converter.ConvertGamePlayerToUpdateParams(player))
	if err != nil {
		return err
	}

	player.UpdatedAt = result.Updatedat
	return nil
}

// helper to update a player using a transaction or DB
// update an existing player
func (c *client) UpdatePlayerUserId(ctx context.Context, player *cs.Player) error {
	return c.writer.UpdatePlayerUserID(ctx, generated.UpdatePlayerUserIDParams{
		ID:     player.ID,
		Userid: sql.NullInt64{Valid: true, Int64: player.UserID},
	})
}

// delete a player by id
func (c *client) DeletePlayer(ctx context.Context, id int64) error {
	return c.writer.DeletePlayer(ctx, id)
}
