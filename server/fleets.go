package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-pkgz/rest"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
)

type cargoTransferRequest struct {
	MO                 cs.MapObject            `json:"mo,omitempty"`
	TransferAmount     cs.CargoTransferRequest `json:"transferAmount,omitempty"`
	FuelTransferAmount int                     `json:"fuelTransferAmount,omitempty"`
}

func (req *cargoTransferRequest) Bind(r *http.Request) error {
	return nil
}

type fleetRequest struct {
	*cs.Fleet
}

func (req *fleetRequest) Bind(r *http.Request) error {
	return nil
}

type fleetRenameRequest struct {
	Name string `json:"name,omitempty"`
}

func (req *fleetRenameRequest) Bind(r *http.Request) error {
	return nil
}

type mergeFleetRequest struct {
	FleetNums []int `json:"fleetNums,omitempty"`
}

func (req *mergeFleetRequest) Bind(r *http.Request) error {
	return nil
}

type splitFleetRequest struct {
	// the source/dest fleet num
	SourceFleetNum int `json:"sourceFleetNum,omitempty"`
	DestFleetNum   int `json:"destFleetNum,omitempty"`

	// a matching slice of source and dest tokens that only differ in token.Quantity
	SourceTokens []cs.ShipToken `json:"sourceTokens,omitempty"`
	DestTokens   []cs.ShipToken `json:"destTokens,omitempty"`

	// the amount of cargo to transfer from the source fleet to the dest when splitting
	TransferAmount cs.CargoTransferRequest `json:"transferAmount,omitempty"`
}

type splitFleetResponse struct {
	Source *cs.Fleet `json:"source,omitempty"`
	Dest   *cs.Fleet `json:"dest,omitempty"`
}

func (req *splitFleetRequest) Bind(r *http.Request) error {
	return nil
}

// context for /api/games/{id}/fleets/{num} calls that require a shipDesign
func (s *server) fleetCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := s.contextDb(r)
		player := s.contextPlayer(r)

		num, err := s.intURLParam(r, "num")
		if num == nil || err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		fleet, err := db.GetFleetByNum(player.GameID, player.Num, *num)
		if err != nil {
			render.Render(w, r, ErrInternalServerError(err))
			return
		}

		if fleet == nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		// only fleet owners can load this fleet
		if fleet.PlayerNum != player.Num {
			render.Render(w, r, ErrForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), keyFleet, fleet)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *server) contextFleet(r *http.Request) *cs.Fleet {
	return r.Context().Value(keyFleet).(*cs.Fleet)
}

func (s *server) fleet(w http.ResponseWriter, r *http.Request) {
	fleet := s.contextFleet(r)
	rest.RenderJSON(w, fleet)
}

// Allow a user to update a fleet's orders
func (s *server) renameFleet(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	fleet := s.contextFleet(r)
	player := s.contextPlayer(r)

	rename := fleetRenameRequest{}
	if err := render.Bind(r, &rename); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	if fleet.PlayerNum != player.Num {
		log.Error().Int64("ID", fleet.ID).Int("Num", fleet.Num).Int("PlayerNum", fleet.PlayerNum).Msg("rename fleet forbidden")
		render.Render(w, r, ErrForbidden)
		return
	}

	// update fleet name in db
	fleet.BaseName = rename.Name
	fleet.Name = fmt.Sprintf("%s #%d", rename.Name, fleet.Num)
	if err := db.UpdateFleet(fleet); err != nil {
		log.Error().Err(err).Int64("ID", fleet.ID).Msg("update fleet in database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, fleet)
}

// Allow a user to update a fleet's orders
func (s *server) updateFleetOrders(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	existingFleet := s.contextFleet(r)
	game := s.contextGame(r)
	player := s.contextPlayer(r)
	var err error

	fleet := fleetRequest{}
	if err := render.Bind(r, &fleet); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	player.Designs, err = db.GetShipDesignsForPlayer(game.ID, player.Num)
	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("get fleets for player")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	existingFleet.InjectDesigns(player.Designs)

	orderer := cs.NewOrderer()
	orderer.UpdateFleetOrders(player, existingFleet, fleet.FleetOrders)

	if err := db.UpdateFleet(existingFleet); err != nil {
		log.Error().Err(err).Int64("ID", fleet.ID).Msg("update fleet in database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, existingFleet)
}

// split a fleet into 2 fleets
func (s *server) split(w http.ResponseWriter, r *http.Request) {
	readWriteClient := s.contextDb(r)
	sourceFleet := s.contextFleet(r)
	game := s.contextGame(r)
	player := s.contextPlayer(r)

	splitFleet := splitFleetRequest{}
	if err := render.Bind(r, &splitFleet); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// load the dest fleet, designs, and player fleets from the db
	var err error
	var destFleet *cs.Fleet
	if splitFleet.DestFleetNum != 0 {
		destFleet, err = readWriteClient.GetFleetByNum(game.ID, player.Num, splitFleet.DestFleetNum)
		if err != nil {
			log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Int("Num", splitFleet.DestFleetNum).Msg("get dest fleet for split")
			render.Render(w, r, ErrInternalServerError(err))
			return
		}
	}

	player.Designs, err = readWriteClient.GetShipDesignsForPlayer(game.ID, player.Num)
	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("get fleets for player")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	playerFleets, err := readWriteClient.GetFleetsForPlayer(game.ID, player.Num)
	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("get fleets for player")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	// execute the order
	orderer := cs.NewOrderer()
	splitFleetRequest := cs.SplitFleetRequest{
		Source:         sourceFleet,
		Dest:           destFleet,
		SourceTokens:   splitFleet.SourceTokens,
		DestTokens:     splitFleet.DestTokens,
		TransferAmount: splitFleet.TransferAmount,
	}

	source, dest, err := orderer.SplitFleet(&game.Rules, player, playerFleets, splitFleetRequest)
	if err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// save the updated fleets back to the database
	if err := s.db.WrapInTransaction(func(c db.Client) error {
		if err := c.UpdateFleet(source); err != nil {
			log.Error().Err(err).Msg("update fleet in database")
			return err
		}

		if dest.ID == 0 {
			dest.GameID = game.ID
			if err := c.CreateFleet(dest); err != nil {
				log.Error().Err(err).Msg("update fleet in database")
				return err
			}
		} else {
			if err := c.UpdateFleet(dest); err != nil {
				log.Error().Err(err).Msg("update fleet in database")
				return err
			}
		}
		return nil
	}); err != nil {
		log.Error().Err(err).Msg("update game in database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, splitFleetResponse{Source: source, Dest: dest})
}

// split all a fleet's tokens into separate fleets
func (s *server) splitAll(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	fleet := s.contextFleet(r)
	game := s.contextGame(r)
	player := s.contextPlayer(r)

	fleets, err := db.GetFleetsForPlayer(game.ID, player.Num)
	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("get fleets for player")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	player.Designs, err = db.GetShipDesignsForPlayer(game.ID, player.Num)
	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("get fleets for player")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	orderer := cs.NewOrderer()
	newFleets, err := orderer.SplitAll(&game.Rules, player, fleets, fleet)
	if err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// save all the fleets
	newFleets = append(newFleets, fleet)
	if err := db.CreateUpdateOrDeleteFleets(game.ID, newFleets); err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("get fleets for player")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, newFleets)
}

// merge target fleets into this one
func (s *server) merge(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	fleet := s.contextFleet(r)
	game := s.contextGame(r)
	player := s.contextPlayer(r)

	mergeFleets := mergeFleetRequest{}
	if err := render.Bind(r, &mergeFleets); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	for _, num := range mergeFleets.FleetNums {
		if fleet.Num == num {
			log.Error().Int64("GameID", game.ID).Int("PlayerNum", player.Num).Int("FleetNum", fleet.Num).Msg("include source fleet Num in merge fleets request")
			render.Render(w, r, ErrBadRequest(fmt.Errorf("invalid merge fleet request")))
			return
		}
	}

	fleets, err := db.GetFleetsByNums(game.ID, player.Num, mergeFleets.FleetNums)
	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("get fleets for merge")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	player.Designs, err = db.GetShipDesignsForPlayer(game.ID, player.Num)
	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("get fleets for player")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	orderer := cs.NewOrderer()
	fleets = append([]*cs.Fleet{fleet}, fleets...)

	updatedFleet, err := orderer.Merge(&game.Rules, player, fleets)
	if err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// save all the fleets
	if err := db.CreateUpdateOrDeleteFleets(game.ID, fleets); err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("get fleets for player")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, updatedFleet)
}

// Transfer cargo from a player's fleet to/from a fleet or planet the player controls
func (s *server) transferCargo(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	game := s.contextGame(r)
	player := s.contextPlayer(r)
	fleet := s.contextFleet(r)
	var err error

	// figure out what type of object we have
	transfer := cargoTransferRequest{}
	if err := render.Bind(r, &transfer); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// the fleet needs designs to compute its spec after
	// transfering cargo
	player.Designs, err = db.GetShipDesignsForPlayer(game.ID, player.Num)
	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("get fleets for player")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	fleet.InjectDesigns(player.Designs)

	switch transfer.MO.Type {
	case cs.MapObjectTypePlanet:
		s.transferCargoFleetPlanet(w, r, &game.Game, player, fleet, transfer.MO.Num, transfer.TransferAmount)
	case cs.MapObjectTypeFleet:
		s.transferCargoFleetFleet(w, r, &game.Game, player, fleet, transfer.MO.PlayerNum, transfer.MO.Num, transfer.TransferAmount, transfer.FuelTransferAmount)
	case cs.MapObjectTypeSalvage:
		s.transferCargoFleetSalvage(w, r, &game.Game, player, fleet, transfer.MO.Num, transfer.TransferAmount)
	default:
		render.Render(w, r, ErrBadRequest(fmt.Errorf("unable to transfer cargo from fleet to %s", transfer.MO.Type)))
		return
	}

}

// transfer cargo from a fleet to/from a planet
func (s *server) transferCargoFleetPlanet(w http.ResponseWriter, r *http.Request, game *cs.Game, player *cs.Player, fleet *cs.Fleet, num int, transferAmount cs.CargoTransferRequest) {
	readClient := s.contextDb(r)
	// find the planet planet by id so we can perform the transfer
	planet, err := readClient.GetPlanetByNum(game.ID, num)
	if err != nil {
		log.Error().Err(err).Msg("get planet from database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	if planet == nil {
		log.Error().Int64("GameID", fleet.GameID).Int("Num", num).Msg("dest planet not found")
		render.Render(w, r, ErrNotFound)
		return
	}

	// don't allow cargo transfers from contested planets
	if !planet.Owned() {
		fleetsInOrbit, err := readClient.GetFleetsOrbitingPlanet(fleet.GameID, planet.Num)
		if err != nil {
			log.Error().Err(err).Msg("get fleets in orbit of planet from database")
			render.Render(w, r, ErrInternalServerError(err))
			return
		}

		// check if any of these fleets are freighters and are owned by someone other than the player
		for _, f := range fleetsInOrbit {
			if f.Spec.CargoCapacity > 0 && f.PlayerNum != fleet.PlayerNum {
				log.Error().Int64("GameID", fleet.GameID).Int("Num", num).Int("PlayerNum", planet.PlayerNum).Msg("dest planet is contested")
				render.Render(w, r, ErrForbidden)
				return
			}
		}
	}

	if planet.Owned() && !planet.OwnedBy(player.Num) {
		log.Error().Int64("GameID", fleet.GameID).Int("Num", num).Int("PlayerNum", planet.PlayerNum).Msg("dest planet not owned by player")
		render.Render(w, r, ErrForbidden)
		return
	}

	orderer := cs.NewOrderer()
	if err := orderer.TransferPlanetCargo(&game.Rules, player, fleet, planet, transferAmount); err != nil {
		log.Error().Err(err).Msg("transfer cargo")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	if err := s.db.WrapInTransaction(func(c db.Client) error {

		if err := c.UpdatePlanet(planet); err != nil {
			return err
		}

		if err := c.UpdateFleet(fleet); err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Error().Err(err).Int64("PlanetID", planet.ID).Int64("FleetID", fleet.ID).Msg("update planet and fleet in database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	log.Info().
		Int64("GameID", fleet.GameID).
		Int("Player", fleet.PlayerNum).
		Str("Fleet", fleet.Name).
		Str("Planet", planet.Name).
		Str("TransferAmount", fmt.Sprintf("%v", transferAmount)).
		Msgf("%s transfered %v to/from Planet %s", fleet.Name, transferAmount, planet.Name)

	// success
	// only return an updated mapobject if we own it
	if planet.PlayerNum == player.Num {
		rest.RenderJSON(w, rest.JSON{"fleet": fleet, "dest": planet})
	} else {
		rest.RenderJSON(w, rest.JSON{"fleet": fleet})
	}
}

// transfer cargo from a fleet to/from a planet
func (s *server) transferCargoFleetSalvage(w http.ResponseWriter, r *http.Request, game *cs.Game, player *cs.Player, fleet *cs.Fleet, num int, transferAmount cs.CargoTransferRequest) {
	db := s.contextDb(r)
	// find the salvage salvage by id so we can perform the transfer
	salvage, err := db.GetSalvageByNum(game.ID, num)
	if err != nil {
		log.Error().Err(err).Msg("get salvage from database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	salvages, err := db.GetSalvagesForGame(game.ID)
	if err != nil {
		log.Error().Err(err).Msg("get salvages from database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}
	nextSalvageNum := 1
	if len(salvages) > 0 {
		nextSalvageNum = salvages[len(salvages)-1].Num + 1
	}

	fullPlayer, err := db.GetPlayer(player.ID)
	if err != nil {
		log.Error().Err(err).Msg("get player from database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	orderer := cs.NewOrderer()
	salvage, err = orderer.TransferSalvageCargo(&game.Rules, fullPlayer, fleet, salvage, nextSalvageNum, transferAmount)
	if err != nil {
		log.Error().Err(err).Msg("transfer cargo")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	if salvage.ID == 0 {
		salvage.GameID = game.ID
		if err := db.CreateSalvage(salvage); err != nil {
			log.Error().Err(err).Int64("ID", salvage.ID).Msg("create salvage in database")
			render.Render(w, r, ErrInternalServerError(err))
			return
		}
	} else {
		if err := db.UpdateSalvage(salvage); err != nil {
			log.Error().Err(err).Int64("ID", salvage.ID).Msg("update salvage in database")
			render.Render(w, r, ErrInternalServerError(err))
			return
		}
	}

	if err := db.UpdateFleet(fleet); err != nil {
		log.Error().Err(err).Msg("update fleet in database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	if err := db.UpdatePlayerSalvageIntels(fullPlayer); err != nil {
		log.Error().Err(err).Msg("update player in database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	log.Info().
		Int64("GameID", fleet.GameID).
		Int("Player", fleet.PlayerNum).
		Str("Fleet", fleet.Name).
		Str("Salvage", salvage.Name).
		Str("TransferAmount", fmt.Sprintf("%v", transferAmount)).
		Msgf("%s transfered %v to/from Salvage %s", fleet.Name, transferAmount, salvage.Name)

	// success
	rest.RenderJSON(w, rest.JSON{"fleet": fleet, "dest": salvage, "salvages": fullPlayer.SalvageIntels})
}

// transfer cargo from a fleet to/from a fleet
func (s *server) transferCargoFleetFleet(w http.ResponseWriter, r *http.Request, game *cs.Game, player *cs.Player, fleet *cs.Fleet, destPlayerNum int, destNum int, transferAmount cs.CargoTransferRequest, fuelTransferAmount int) {
	readWriteClient := s.contextDb(r)
	// find the dest dest by id so we can perform the transfer
	dest, err := readWriteClient.GetFleetByNum(game.ID, destPlayerNum, destNum)
	if err != nil {
		log.Error().Err(err).Msg("get dest fleet from database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	if dest == nil {
		log.Error().Int64("GameID", fleet.GameID).Int("PlayerNum", destPlayerNum).Int("Num", destNum).Msg("dest fleet not found")
		render.Render(w, r, ErrNotFound)
		return
	}

	if dest.Owned() && !dest.OwnedBy(player.Num) {
		log.Error().Int64("GameID", fleet.GameID).Int("Num", fleet.Num).Int("PlayerNum", fleet.PlayerNum).Msg("dest fleet not owned by player")
		render.Render(w, r, ErrForbidden)
		return
	}

	// if we are transferring cargo to another player, load them from the DB
	destPlayer := player
	if dest.PlayerNum != player.Num {
		destPlayer, err = readWriteClient.GetPlayerByNum(game.ID, dest.PlayerNum)
		if err != nil {
			log.Error().Err(err).Msg("get dest player from database")
			render.Render(w, r, ErrInternalServerError(err))
			return
		}

		destPlayer.Designs, err = readWriteClient.GetShipDesignsForPlayer(game.ID, destPlayer.Num)
		if err != nil {
			log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", destPlayer.Num).Msg("get fleets for player")
			render.Render(w, r, ErrInternalServerError(err))
			return
		}

		dest.InjectDesigns(destPlayer.Designs)
	} else {
		dest.InjectDesigns(player.Designs)
	}

	orderer := cs.NewOrderer()
	if err := orderer.TransferFleetCargo(&game.Rules, player, destPlayer, fleet, dest, transferAmount); err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	if err := s.db.WrapInTransaction(func(c db.Client) error {
		if err := c.UpdateFleet(dest); err != nil {
			log.Error().Err(err).Msg("update fleet in database")
			return err
		}

		if err := c.UpdateFleet(fleet); err != nil {
			log.Error().Err(err).Msg("update fleet in database")
			return err
		}
		return nil
	}); err != nil {
		log.Error().Err(err).Msg("update game in database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	log.Info().
		Int64("GameID", fleet.GameID).
		Int("Player", fleet.PlayerNum).
		Str("Fleet", fleet.Name).
		Str("Planet", dest.Name).
		Str("TransferAmount", fmt.Sprintf("%v", transferAmount)).
		Msgf("%s transfered %v to/from Planet %s", fleet.Name, transferAmount, dest.Name)

	// success
	// only return an updated mapobject if we own it
	if dest.PlayerNum == player.Num {
		rest.RenderJSON(w, rest.JSON{"fleet": fleet, "dest": dest})
	} else {
		rest.RenderJSON(w, rest.JSON{"fleet": fleet})
	}
}
