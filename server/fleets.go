package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-pkgz/rest"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
)

type cargoTransferRequest struct {
	MO             cs.MapObject `json:"mo,omitempty"`
	TransferAmount cs.Cargo     `json:"transferAmount,omitempty"`
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

type mergeFleetRequest struct {
	FleetNums []int `json:"fleetNums,omitempty"`
}

func (req *mergeFleetRequest) Bind(r *http.Request) error {
	return nil
}

// context for /api/games/{id}/fleets/{num} calls that require a shipDesign
func (s *server) fleetCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		player := s.contextPlayer(r)

		num, err := s.intURLParam(r, "num")
		if num == nil || err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		fleet, err := s.db.GetFleetByNum(player.GameID, player.Num, *num)
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
func (s *server) updateFleetOrders(w http.ResponseWriter, r *http.Request) {
	existingFleet := s.contextFleet(r)
	game := s.contextGame(r)
	player := s.contextPlayer(r)
	var err error

	fleet := fleetRequest{}
	if err := render.Bind(r, &fleet); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	player.Designs, err = s.db.GetShipDesignsForPlayer(game.ID, player.Num)
	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("get fleets for player")
		render.Render(w, r, ErrInternalServerError(err))
	}

	existingFleet.InjectDesigns(player.Designs)

	orderer := cs.NewOrderer()
	orderer.UpdateFleetOrders(player, existingFleet, fleet.FleetOrders)

	if err := s.db.UpdateFleet(existingFleet); err != nil {
		log.Error().Err(err).Int64("ID", fleet.ID).Msg("update fleet in database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, existingFleet)
}

// split all a fleet's tokens into separate fleets
func (s *server) splitAll(w http.ResponseWriter, r *http.Request) {
	fleet := s.contextFleet(r)
	game := s.contextGame(r)
	player := s.contextPlayer(r)

	fleets, err := s.db.GetFleetsForPlayer(game.ID, player.Num)
	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("get fleets for player")
		render.Render(w, r, ErrInternalServerError(err))
	}

	player.Designs, err = s.db.GetShipDesignsForPlayer(game.ID, player.Num)
	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("get fleets for player")
		render.Render(w, r, ErrInternalServerError(err))
	}

	orderer := cs.NewOrderer()
	newFleets, err := orderer.SplitAll(&game.Rules, player, fleets, fleet)
	if err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// save all the fleets
	newFleets = append(newFleets, fleet)
	if err := s.db.CreateUpdateOrDeleteFleets(game.ID, newFleets); err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("get fleets for player")
		render.Render(w, r, ErrInternalServerError(err))
	}

	rest.RenderJSON(w, newFleets)
}

// merge target fleets into this one
func (s *server) merge(w http.ResponseWriter, r *http.Request) {
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

	fleets, err := s.db.GetFleetsByNums(game.ID, player.Num, mergeFleets.FleetNums)
	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("get fleets for merge")
		render.Render(w, r, ErrInternalServerError(err))
	}

	player.Designs, err = s.db.GetShipDesignsForPlayer(game.ID, player.Num)
	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("get fleets for player")
		render.Render(w, r, ErrInternalServerError(err))
	}

	orderer := cs.NewOrderer()
	fleets = append([]*cs.Fleet{fleet}, fleets...)

	updatedFleet, err := orderer.Merge(&game.Rules, player, fleets)
	if err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// save all the fleets
	if err := s.db.CreateUpdateOrDeleteFleets(game.ID, fleets); err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("get fleets for player")
		render.Render(w, r, ErrInternalServerError(err))
	}

	rest.RenderJSON(w, updatedFleet)
}

// Transfer cargo from a player's fleet to/from a fleet or planet the player controls
func (s *server) transferCargo(w http.ResponseWriter, r *http.Request) {
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
	player.Designs, err = s.db.GetShipDesignsForPlayer(game.ID, player.Num)
	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("get fleets for player")
		render.Render(w, r, ErrInternalServerError(err))
	}

	fleet.InjectDesigns(player.Designs)

	switch transfer.MO.Type {
	case cs.MapObjectTypePlanet:
		s.transferCargoFleetPlanet(w, r, &game.Game, player, fleet, transfer.MO.ID, transfer.TransferAmount)
	case cs.MapObjectTypeFleet:
		s.transferCargoFleetFleet(w, r, &game.Game, player, fleet, transfer.MO.ID, transfer.TransferAmount)
	default:
		render.Render(w, r, ErrBadRequest(fmt.Errorf("unable to transfer cargo from fleet to %s", transfer.MO.Type)))
		return
	}

}

// transfer cargo from a fleet to/from a planet
func (s *server) transferCargoFleetPlanet(w http.ResponseWriter, r *http.Request, game *cs.Game, player *cs.Player, fleet *cs.Fleet, destID int64, transferAmount cs.Cargo) {
	// find the planet planet by id so we can perform the transfer
	planet, err := s.db.GetPlanet(destID)
	if err != nil {
		log.Error().Err(err).Msg("get planet from database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	if planet == nil {
		log.Error().Int64("GameID", fleet.GameID).Int64("DestID", destID).Msg("dest planet not found")
		render.Render(w, r, ErrNotFound)
		return
	}

	orderer := cs.NewOrderer()
	if err := orderer.TransferPlanetCargo(&game.Rules, player, fleet, planet, transferAmount); err != nil {
		log.Error().Err(err).Msg("transfer cargo")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	if err := s.db.UpdatePlanet(planet); err != nil {
		log.Error().Err(err).Int64("ID", planet.ID).Msg("update planet in database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	if err := s.db.UpdateFleet(fleet); err != nil {
		log.Error().Err(err).Msg("update fleet in database")
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

// transfer cargo from a fleet to/from a fleet
func (s *server) transferCargoFleetFleet(w http.ResponseWriter, r *http.Request, game *cs.Game, player *cs.Player, fleet *cs.Fleet, destID int64, transferAmount cs.Cargo) {
	// find the dest dest by id so we can perform the transfer
	dest, err := s.db.GetFleet(destID)
	if err != nil {
		log.Error().Err(err).Msg("get dest fleet from database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	if dest == nil {
		log.Error().Int64("GameID", fleet.GameID).Int64("DestID", destID).Msg("dest fleet not found")
		render.Render(w, r, ErrNotFound)
		return
	}

	// if we are transferring cargo to another player, load them from the DB
	destPlayer := player
	if dest.PlayerNum != player.Num {
		destPlayer, err = s.db.GetPlayerByNum(game.ID, dest.PlayerNum)
		if err != nil {
			log.Error().Err(err).Msg("get dest player from database")
			render.Render(w, r, ErrInternalServerError(err))
			return
		}

		destPlayer.Designs, err = s.db.GetShipDesignsForPlayer(game.ID, destPlayer.Num)
		if err != nil {
			log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", destPlayer.Num).Msg("get fleets for player")
			render.Render(w, r, ErrInternalServerError(err))
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

	if err := s.db.UpdateFleet(dest); err != nil {
		log.Error().Err(err).Msg("update fleet in database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	if err := s.db.UpdateFleet(fleet); err != nil {
		log.Error().Err(err).Msg("update fleet in database")
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
