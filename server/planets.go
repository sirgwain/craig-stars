//go:build !wasi && !wasm

package server

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-pkgz/rest"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
)

type planetRequest struct {
	*cs.Planet
}

type planetProductionEstimateRequest struct {
	Planet *cs.Planet `json:"planet,omitempty"`
	Rules  *cs.Rules  `json:"rules,omitempty"`
	Player *cs.Player `json:"player,omitempty"`
}

type starbaseUpgradeRequest struct {
	Design    *cs.ShipDesign `json:"design,omitempty"`
	NewDesign *cs.ShipDesign `json:"newDesign,omitempty"`
}

func (req *planetRequest) Bind(r *http.Request) error {
	return nil
}

func (req *planetProductionEstimateRequest) Bind(r *http.Request) error {
	return nil
}

func (req *starbaseUpgradeRequest) Bind(r *http.Request) error {
	return nil
}

// context for /api/games/{id}/planets/{num} calls that require a shipDesign
func (s *server) planetCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := s.contextDb(r)
		player := s.contextPlayer(r)

		num, err := s.intURLParam(r, "num")
		if num == nil || err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		planet, err := db.GetPlanetByNum(player.GameID, *num)
		if err != nil {
			render.Render(w, r, ErrInternalServerError(err))
			return
		}

		if planet == nil {
			log.Error().Int64("GameID", player.GameID).Msgf("unable to find planet %d", num)
			render.Render(w, r, ErrNotFound)
			return
		}

		// only planet owners can load this planet
		if planet.PlayerNum != player.Num {
			render.Render(w, r, ErrForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), keyPlanet, planet)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *server) contextPlanet(r *http.Request) *cs.Planet {
	return r.Context().Value(keyPlanet).(*cs.Planet)
}

func (s *server) planet(w http.ResponseWriter, r *http.Request) {
	planet := s.contextPlanet(r)
	rest.RenderJSON(w, planet)
}

// Allow a user to update a planet's orders
func (s *server) updatePlanetOrders(w http.ResponseWriter, r *http.Request) {
	dbClient := s.contextDb(r)
	game := s.contextGame(r)
	existingPlanet := s.contextPlanet(r)
	player := s.contextPlayer(r)

	planet := planetRequest{}
	if err := render.Bind(r, &planet); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// load the full player to update planet production estimates
	player, err := dbClient.GetPlayerWithDesignsForGame(game.ID, player.Num)
	if err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	// load all a player's planets so we can recompute research estimates
	planets, err := dbClient.GetPlanetsForPlayer(game.ID, player.Num)
	if err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	orderer := cs.NewOrderer()
	if err := orderer.UpdatePlanetOrders(&game.Rules, player, existingPlanet, planet.PlanetOrders, planets); err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Str("Planet", existingPlanet.Name).Msg("update planet orders")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	// update this planet and the player's spec in the database
	if err := s.db.WrapInTransaction(func(c db.Client) error {
		if err := c.UpdatePlanet(existingPlanet); err != nil {
			log.Error().Err(err).Int64("ID", planet.ID).Msg("update planet in database")
			return err
		}

		// update the player spec as well because changes in planet orders impact resources
		// available for research
		if err := c.UpdatePlayerSpec(player); err != nil {
			log.Error().Err(err).Int64("ID", planet.ID).Msg("update player spec in database")
			return err
		}
		return nil
	}); err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, rest.JSON{"planet": existingPlanet, "player": player})
}

// get an estimate for production completion based on a planet's production queue items
func (s *server) getPlanetProductionEstimate(w http.ResponseWriter, r *http.Request) {

	estimateRequest := planetProductionEstimateRequest{}
	if err := render.Bind(r, &estimateRequest); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	rules := cs.StandardRules
	if estimateRequest.Rules != nil {
		rules = *estimateRequest.Rules
	}

	planet := estimateRequest.Planet

	// populate the production queue estimates
	planet.PopulateProductionQueueEstimates(&rules, estimateRequest.Player)
	rest.RenderJSON(w, planet)
}

// get an estimate for production completion based on a planet's production queue items
func (s *server) getStarbaseUpgradeCost(w http.ResponseWriter, r *http.Request) {
	rules := s.contextGame(r).Game.Rules
	player := s.contextPlayer(r)
	upgradeRequest := starbaseUpgradeRequest{}
	if err := render.Bind(r, &upgradeRequest); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	calculator := cs.NewCostCalculator()
	cost, err := calculator.StarbaseUpgradeCost(&rules, player.TechLevels, player.Race.Spec, upgradeRequest.Design, upgradeRequest.NewDesign)
	if err != nil {
		render.Render(w, r, ErrBadRequest(err))
	}
	rest.RenderJSON(w, cost)
}
