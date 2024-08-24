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

type shipDesignRequest struct {
	*cs.ShipDesign
}

func (req *shipDesignRequest) Bind(r *http.Request) error {
	return nil
}

// context for /api/games/{id}/designs/{num} calls that require a shipDesign
func (s *server) shipdDesignCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := s.contextDb(r)
		player := s.contextPlayer(r)

		num, err := s.intURLParam(r, "num")
		if num == nil || err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		design, err := db.GetShipDesignByNum(player.GameID, player.Num, *num)
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
	db := s.contextDb(r)
	player := s.contextPlayer(r)

	shipDesigns, err := db.GetShipDesignsForPlayer(player.GameID, player.Num)
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
	db := s.contextDb(r)
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

	designs, err := db.GetShipDesignsForPlayer(game.ID, player.Num)
	if err != nil {
		log.Error().Err(err).Int64("ID", player.ID).Msg("load designs for player")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	design.PlayerNum = player.Num
	design.GameID = player.GameID
	design.Num = player.GetNextDesignNum(designs)
	design.Spec = cs.ComputeShipDesignSpec(&game.Rules, player.TechLevels, player.Race.Spec, design.ShipDesign)

	if err := db.CreateShipDesign(design.ShipDesign); err != nil {
		log.Error().Err(err).Int64("ID", player.ID).Str("DesignName", design.Name).Msg("save new player design")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	log.Info().Int64("GameID", design.GameID).Int("PlayerNum", player.Num).Str("DesignName", design.Name).Msg("created player design")

	rest.RenderJSON(w, design)
}

func (s *server) updateShipDesign(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	game := s.contextGame(r)
	player := s.contextPlayer(r)

	design := shipDesignRequest{}
	if err := render.Bind(r, &design); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// load in the existing shipDesign from the context
	existingDesign := s.contextShipDesign(r)

	if existingDesign.Spec.NumInstances > 0 {
		log.Error().
			Int64("GameID", existingDesign.GameID).
			Int64("ID", existingDesign.ID).
			Int("PlayerNum", existingDesign.PlayerNum).
			Str("DesignName", design.Name).
			Msgf("design in use (%d instances), cannot update", existingDesign.Spec.NumInstances)
		render.Render(w, r, ErrBadRequest(fmt.Errorf("design is in use, cannot update")))
		return
	}

	if existingDesign.OriginalPlayerNum != cs.None {
		log.Error().
			Int64("GameID", existingDesign.GameID).
			Int64("ID", existingDesign.ID).
			Int("PlayerNum", existingDesign.PlayerNum).
			Str("DesignName", design.Name).
			Msg("design is transferred from another player, cannot update")
		render.Render(w, r, ErrBadRequest(fmt.Errorf("design is transferred, cannot update")))
		return
	}

	// edge case for bad request where the url num doesn't match the request payload
	if design.Num != existingDesign.Num {
		log.Error().Int64("ID", design.ID).Msgf("design.Num %d != existingDesign.Num %d", design.ID, existingDesign.ID)
		render.Render(w, r, ErrBadRequest(fmt.Errorf("shipDesign id/user id does not match existing shipDesign")))
		return
	}

	// validate
	design.PlayerNum = player.Num
	design.GameID = player.GameID
	design.Spec = cs.ComputeShipDesignSpec(&game.Rules, player.TechLevels, player.Race.Spec, design.ShipDesign)

	if err := design.Validate(&game.Rules, player); err != nil {
		log.Error().Err(err).Int64("ID", player.ID).Str("DesignName", design.Name).Msg("validate player design")
		render.Render(w, r, ErrBadRequest(fmt.Errorf("shipDesign id/user id does not match existing shipDesign")))
		return
	}

	if err := db.UpdateShipDesign(design.ShipDesign); err != nil {
		log.Error().Err(err).Int64("ID", design.ID).Msg("update shipDesign in database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	log.Info().Int64("GameID", design.GameID).Int("PlayerNum", player.Num).Str("DesignName", design.Name).Msg("updated player design")

	rest.RenderJSON(w, design)
}

func (s *server) deleteShipDesign(w http.ResponseWriter, r *http.Request) {
	readWriteClient := s.contextDb(r)
	game := s.contextGame(r)
	player := s.contextPlayer(r)
	design := s.contextShipDesign(r)

	// validate

	if design.CannotDelete {
		log.Error().Int64("ID", player.ID).Str("DesignName", design.Name).Msg("delete design with CannotDelete = true")
		render.Render(w, r, ErrBadRequest(fmt.Errorf("shipDesign cannot be deleted")))
		return
	}

	// to delete ship designs we have to remove any tokens using the design from fleets
	playerFleets, err := readWriteClient.GetFleetsForPlayer(game.ID, player.Num)
	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("load fleets from database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	playerDesigns, err := readWriteClient.GetShipDesignsForPlayer(game.ID, player.Num)
	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("get designs for player")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	playerPlanets, err := readWriteClient.GetPlanetsForPlayer(game.ID, player.Num)
	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("load planets from database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	fleetsToDelete := []*cs.Fleet{}
	fleetsToUpdate := []*cs.Fleet{}
	leftoverPlayerFleets := []*cs.Fleet{}
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
				fleet.InjectDesigns(playerDesigns)
				fleet.Spec = cs.ComputeFleetSpec(&game.Rules, player, fleet)
				fleetsToUpdate = append(fleetsToUpdate, fleet)
			}
			leftoverPlayerFleets = append(leftoverPlayerFleets, fleet)
		}
	}

	// remove this design from any production queues
	planetsToUpdate := []*cs.Planet{}
	for _, planet := range playerPlanets {
		newQueue := make([]cs.ProductionQueueItem, 0, len(planet.ProductionQueue))
		for _, item := range planet.ProductionQueue {
			if item.DesignNum == design.Num {
				// remove this design from the queue and mark this planet for update
				planetsToUpdate = append(planetsToUpdate, planet)
				continue
			}
			newQueue = append(newQueue, item)
		}
		planet.ProductionQueue = newQueue
	}

	// update fleets to remove tokens, delete fleets without tokens, update planets with production queue changes
	if err := s.db.WrapInTransaction(func(c db.Client) error {

		for _, fleet := range fleetsToUpdate {
			if err := c.UpdateFleet(fleet); err != nil {
				return fmt.Errorf("update fleet in database %w", err)
			}
			log.Info().Int64("GameID", game.ID).Int("PlayerNum", player.Num).Int("Num", design.Num).Msgf("updated fleet %s after deleting design", fleet.Name)
		}

		for _, fleet := range fleetsToDelete {
			if err := c.DeleteFleet(fleet.ID); err != nil {
				return fmt.Errorf("delete fleet from database %w", err)
			}
			log.Info().Int64("GameID", game.ID).Int("PlayerNum", player.Num).Int("Num", design.Num).Msgf("deleted fleet %s after deleting design", fleet.Name)
		}

		for _, planet := range planetsToUpdate {
			if err := c.UpdatePlanet(planet); err != nil {
				return fmt.Errorf("update planet in database %w", err)
			}
			log.Info().Int64("GameID", game.ID).Int("PlayerNum", player.Num).Int("Num", design.Num).Msgf("updated planet %s after deleting design", planet.Name)

		}

		if err := c.DeleteShipDesign(design.ID); err != nil {
			return fmt.Errorf("delete design from database %w", err)
		}
		log.Info().Int64("GameID", game.ID).Int("PlayerNum", player.Num).Int("Num", design.Num).Msgf("deleted design %s", design.Name)

		return nil
	}); err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("delete design from database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	
	// split the player fleets into fleets and starbases
	fleets := make([]*cs.Fleet, 0, len(leftoverPlayerFleets))
	starbases := make([]*cs.Fleet, 0)
	for i := range leftoverPlayerFleets {
		fleet := leftoverPlayerFleets[i]
		if fleet.Starbase {
			starbases = append(starbases, fleet)
		} else {
			fleets = append(fleets, fleet)
		}
	}
	rest.RenderJSON(w, rest.JSON{"fleets": fleets, "starbases": starbases, "planets": playerPlanets})
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
