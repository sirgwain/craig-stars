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

// context for /api/games/{id}/fleets/{num} calls that require a shipDesign
func (s *server) fleetCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		player := s.contextPlayer(r)

		num, err := s.intURLParam(r, "num")
		if num == nil || err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		fleet, err := s.db.GetFleetByNum(player.GameID, *num)
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
	player := s.contextPlayer(r)

	fleet := fleetRequest{}
	if err := render.Bind(r, &fleet); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	orderer := cs.NewOrderer()
	orderer.UpdateFleetOrders(player, existingFleet, fleet.FleetOrders)

	if err := s.db.UpdateFleet(existingFleet); err != nil {
		log.Error().Err(err).Int64("ID", fleet.ID).Msg("update fleet in database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, existingFleet)
}

// Transfer cargo from a player's fleet to/from a fleet or planet the player controls
func (s *server) transferCargo(w http.ResponseWriter, r *http.Request) {
	fleet := s.contextFleet(r)

	// figure out what type of object we have
	transfer := cargoTransferRequest{}
	if err := render.Bind(r, &transfer); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	switch transfer.MO.Type {
	case cs.MapObjectTypePlanet:
		s.transferCargoFleetPlanet(w, r, fleet, transfer.MO.ID, transfer.TransferAmount)
	case cs.MapObjectTypeFleet:
		s.transferCargoFleetFleet(w, r, fleet, transfer.MO.ID, transfer.TransferAmount)
	default:
		render.Render(w, r, ErrBadRequest(fmt.Errorf("unable to transfer cargo from fleet to %s", transfer.MO.Type)))
		return
	}

}

// transfer cargo from a fleet to/from a planet
func (s *server) transferCargoFleetPlanet(w http.ResponseWriter, r *http.Request, fleet *cs.Fleet, destID int64, transferAmount cs.Cargo) {
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
	if err := orderer.TransferPlanetCargo(fleet, planet, transferAmount); err != nil {
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
	rest.RenderJSON(w, fleet)
}

// transfer cargo from a fleet to/from a fleet
func (s *server) transferCargoFleetFleet(w http.ResponseWriter, r *http.Request, fleet *cs.Fleet, destID int64, transferAmount cs.Cargo) {
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
	}

	orderer := cs.NewOrderer()
	if err := orderer.TransferFleetCargo(fleet, dest, transferAmount); err != nil {
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
	rest.RenderJSON(w, fleet)
}
