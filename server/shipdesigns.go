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

type shipDesignRequest struct {
	*cs.ShipDesign
}

func (req *shipDesignRequest) Bind(r *http.Request) error {
	return nil
}

// context for /api/games/{id}/designs/{num} calls that require a shipDesign
func (s *server) shipdDesignCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		player := s.contextPlayer(r)

		num, err := s.intURLParam(r, "num")
		if num == nil || err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		design, err := s.db.GetShipDesignByNum(player.GameID, player.Num, *num)
		if err != nil {
			render.Render(w, r, ErrInternalServerError(err))
			return
		}

		if design == nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), keyShipDesign, design)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *server) contextShipDesign(r *http.Request) *cs.ShipDesign {
	return r.Context().Value(keyShipDesign).(*cs.ShipDesign)
}

// CRUD for ship designs
func (s *server) shipDesigns(w http.ResponseWriter, r *http.Request) {
	player := s.contextPlayer(r)

	shipDesigns, err := s.db.GetShipDesignsForPlayer(player.GameID, player.Num)
	if err != nil {
		log.Error().Err(err).Int64("GameID", player.GameID).Int("PlayerNum", player.Num).Msg("get shipDesigns from database")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	rest.RenderJSON(w, shipDesigns)
}

func (s *server) shipDesign(w http.ResponseWriter, r *http.Request) {
	shipDesign := s.contextShipDesign(r)
	rest.RenderJSON(w, shipDesign)
}

func (s *server) createShipDesign(w http.ResponseWriter, r *http.Request) {
	game := s.contextGame(r)
	player := s.contextPlayer(r)

	design := shipDesignRequest{}
	if err := render.Bind(r, &design); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	if err := design.Validate(&game.Rules, player); err != nil {
		log.Error().Err(err).Int64("ID", player.ID).Str("DesignName", design.Name).Msg("validate new player design")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	designs, err := s.db.GetShipDesignsForPlayer(game.ID, player.Num)
	if err != nil {
		log.Error().Err(err).Int64("ID", player.ID).Msg("load designs for player")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	design.PlayerNum = player.Num
	design.GameID = player.GameID
	design.Num = player.GetNextDesignNum(designs)
	design.Spec = cs.ComputeShipDesignSpec(&game.Rules, player.TechLevels, player.Race.Spec, design.ShipDesign)

	if err := s.db.CreateShipDesign(design.ShipDesign); err != nil {
		log.Error().Err(err).Int64("ID", player.ID).Str("DesignName", design.Name).Msg("save new player design")
		render.Render(w, r, ErrInternalServerError(err))
	}

	rest.RenderJSON(w, design)
}

func (s *server) updateShipDesign(w http.ResponseWriter, r *http.Request) {
	game := s.contextGame(r)
	player := s.contextPlayer(r)

	design := shipDesignRequest{}
	if err := render.Bind(r, &design); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// load in the existing shipDesign from the context
	existingDesign := s.contextShipDesign(r)

	// validate

	design.PlayerNum = player.Num
	design.GameID = player.GameID
	design.Spec = cs.ComputeShipDesignSpec(&game.Rules, player.TechLevels, player.Race.Spec, design.ShipDesign)

	// edge case for bad request where the url num doesn't match the request payload
	if design.Num != existingDesign.Num {
		log.Error().Int64("ID", design.ID).Msgf("design.Num %d != existingDesign.Num %d", design.ID, existingDesign.ID)
		render.Render(w, r, ErrBadRequest(fmt.Errorf("shipDesign id/user id does not match existing shipDesign")))
		return
	}

	if err := design.Validate(&game.Rules, player); err != nil {
		log.Error().Err(err).Int64("ID", player.ID).Str("DesignName", design.Name).Msg("validate player design")
		render.Render(w, r, ErrBadRequest(fmt.Errorf("shipDesign id/user id does not match existing shipDesign")))
		return
	}

	if err := s.db.UpdateShipDesign(design.ShipDesign); err != nil {
		log.Error().Err(err).Int64("ID", design.ID).Msg("update shipDesign in database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, design)
}

func (s *server) deleteShipDesign(w http.ResponseWriter, r *http.Request) {
	game := s.contextGame(r)
	player := s.contextPlayer(r)
	design := s.contextShipDesign(r)

	// validate

	if !design.CanDelete {
		log.Error().Int64("ID", player.ID).Str("DesignName", design.Name).Msg("delete design with CanDelete = false")
		render.Render(w, r, ErrBadRequest(fmt.Errorf("shipDesign cannot be deleted")))
		return
	}

	// delete the ship design

	playerFleets, err := s.db.GetFleetsForPlayer(game.ID, player.Num)
	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("load fleets from database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	fleetsToDelete := []*cs.Fleet{}
	fleetsToUpdate := []*cs.Fleet{}
	for _, fleet := range playerFleets {
		// find any tokens using this design
		updatedTokens := make([]cs.ShipToken, 0, len(fleet.Tokens))
		for _, token := range fleet.Tokens {
			if token.DesignNum != design.Num {
				updatedTokens = append(updatedTokens, token)
			}
		}
		// if we have no tokens left, delete the fleet
		if len(updatedTokens) == 0 {
			fleetsToDelete = append(fleetsToDelete, fleet)
		} else {
			// if we have a different number of tokens than we
			// had before, update this fleet
			if len(updatedTokens) != len(fleet.Tokens) {
				fleet.Tokens = updatedTokens
				fleet.Spec = cs.ComputeFleetSpec(&game.Rules, player, fleet)
				fleetsToUpdate = append(fleetsToUpdate, fleet)
			}
		}

	}

	// delete a ship design and update fleets in one transaction
	if err := s.db.DeleteShipDesignWithFleets(design.ID, fleetsToUpdate, fleetsToDelete); err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("load fleets from database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	// log what we did
	log.Info().Int64("GameID", game.ID).Int("PlayerNum", player.Num).Int("Num", design.Num).Msgf("deleted design %s", design.Name)

	for _, fleet := range fleetsToUpdate {
		log.Info().Int64("GameID", game.ID).Int("PlayerNum", player.Num).Int("Num", design.Num).Msgf("updated fleet %s after deleting design", fleet.Name)
	}

	for _, fleet := range fleetsToDelete {
		log.Info().Int64("GameID", game.ID).Int("PlayerNum", player.Num).Int("Num", design.Num).Msgf("deleted fleet %s after deleting design", fleet.Name)
	}

	allFleets, err := s.db.GetFleetsForPlayer(game.ID, player.Num)

	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Int("Num", design.Num).Msg("load fleets from database")
		render.Render(w, r, ErrInternalServerError(err))
	}

	// split the player fleets into fleets and starbases
	fleets := make([]*cs.Fleet, 0, len(allFleets))
	starbases := make([]*cs.Fleet, 0)
	for i := range allFleets {
		fleet := allFleets[i]
		if fleet.Starbase {
			starbases = append(starbases, fleet)
		} else {
			fleets = append(fleets, fleet)
		}
	}
	rest.RenderJSON(w, rest.JSON{"fleets": fleets, "starbases": starbases})
}

func (s *server) computeShipDesignSpec(w http.ResponseWriter, r *http.Request) {
	game := s.contextGame(r)
	player := s.contextPlayer(r)

	design := shipDesignRequest{}
	if err := render.Bind(r, &design); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	design.Spec = cs.ComputeShipDesignSpec(&game.Rules, player.TechLevels, player.Race.Spec, design.ShipDesign)

	rest.RenderJSON(w, design.Spec)
}
