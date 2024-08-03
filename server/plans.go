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

type battlePlanRequest struct {
	*cs.BattlePlan
}

func (req *battlePlanRequest) Bind(r *http.Request) error {
	return nil
}

type productionPlanRequest struct {
	*cs.ProductionPlan
}

func (req *productionPlanRequest) Bind(r *http.Request) error {
	return nil
}

type transportPlanRequest struct {
	*cs.TransportPlan
}

func (req *transportPlanRequest) Bind(r *http.Request) error {
	return nil
}

// context for /api/games/{id}/battle-plans/{num} calls that require a battlePlan
func (s *server) battlePlanCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		player := s.contextPlayer(r)

		num, err := s.intURLParam(r, "num")
		if num == nil || err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		var battlePlan *cs.BattlePlan
		for i := range player.BattlePlans {
			if player.BattlePlans[i].Num == *num {
				battlePlan = &player.BattlePlans[i]
				break
			}
		}

		if battlePlan == nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), keyBattlePlan, battlePlan)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// context for /api/games/{id}/production-plans/{num} calls that require a productionPlan
func (s *server) productionPlanCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		player := s.contextPlayer(r)

		num, err := s.intURLParam(r, "num")
		if num == nil || err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		var productionPlan *cs.ProductionPlan
		for i := range player.ProductionPlans {
			if player.ProductionPlans[i].Num == *num {
				productionPlan = &player.ProductionPlans[i]
				break
			}
		}

		if productionPlan == nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), keyProductionPlan, productionPlan)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// context for /api/games/{id}/transport-plans/{num} calls that require a transportPlan
func (s *server) transportPlanCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		player := s.contextPlayer(r)

		num, err := s.intURLParam(r, "num")
		if num == nil || err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		var transportPlan *cs.TransportPlan
		for i := range player.TransportPlans {
			if player.TransportPlans[i].Num == *num {
				transportPlan = &player.TransportPlans[i]
				break
			}
		}

		if transportPlan == nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), keyTransportPlan, transportPlan)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *server) contextBattlePlan(r *http.Request) *cs.BattlePlan {
	return r.Context().Value(keyBattlePlan).(*cs.BattlePlan)
}

func (s *server) contextProductionPlan(r *http.Request) *cs.ProductionPlan {
	return r.Context().Value(keyProductionPlan).(*cs.ProductionPlan)
}

func (s *server) contextTransportPlan(r *http.Request) *cs.TransportPlan {
	return r.Context().Value(keyTransportPlan).(*cs.TransportPlan)
}

func (s *server) createBattlePlan(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	player := s.contextPlayer(r)

	battlePlan := battlePlanRequest{}
	if err := render.Bind(r, &battlePlan); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	if err := battlePlan.Validate(player); err != nil {
		log.Error().Err(err).Int64("PlayerID", player.ID).Str("PlanName", battlePlan.Name).Msg("validate new player BattlePlan")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	battlePlan.Num = player.GetNextBattlePlanNum()
	player.BattlePlans = append(player.BattlePlans, *battlePlan.BattlePlan)

	if err := db.UpdatePlayerPlans(player); err != nil {
		log.Error().Err(err).Int64("PlayerID", player.ID).Str("PlanName", battlePlan.Name).Msg("save new player BattlePlan")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, battlePlan)
}

func (s *server) updateBattlePlan(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	player := s.contextPlayer(r)
	existingBattlePlan := s.contextBattlePlan(r)

	battlePlan := battlePlanRequest{}
	if err := render.Bind(r, &battlePlan); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// edge case for bad request where the url num doesn't match the request payload
	if battlePlan.Num != existingBattlePlan.Num {
		log.Error().Int("Num", battlePlan.Num).Msgf("battlePlan.Num %d != existingBattlePlan.Num %d", battlePlan.Num, existingBattlePlan.Num)
		render.Render(w, r, ErrBadRequest(fmt.Errorf("BattlePlan num does not match existing BattlePlan num")))
		return
	}

	// validate
	if err := battlePlan.Validate(player); err != nil {
		log.Error().Err(err).Int64("PlayerID", player.ID).Str("PlanName", battlePlan.Name).Msg("validate BattlePlan")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	for i := range player.BattlePlans {
		if player.BattlePlans[i].Num == battlePlan.Num {
			player.BattlePlans[i] = *battlePlan.BattlePlan
			break
		}
	}

	if err := db.UpdatePlayerPlans(player); err != nil {
		log.Error().Err(err).Int64("PlayerID", player.ID).Str("PlanName", battlePlan.Name).Msg("save new player BattlePlan")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, battlePlan)
}

func (s *server) deleteBattlePlan(w http.ResponseWriter, r *http.Request) {
	readWriteClient := s.contextDb(r)
	game := s.contextGame(r)
	player := s.contextPlayer(r)
	battlePlan := s.contextBattlePlan(r)

	if battlePlan.Num == 0 {
		log.Error().Int64("PlayerID", player.ID).Str("PlanName", battlePlan.Name).Msg("delete default BattlePlan")
		render.Render(w, r, ErrBadRequest(fmt.Errorf("cannot delete default battle plan")))
		return
	}

	// delete the battle plan
	// set all fleets using this battle plan to use the default one
	playerFleets, err := readWriteClient.GetFleetsForPlayer(game.ID, player.Num)
	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("load fleets from database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	fleetsToUpdate := []*cs.Fleet{}
	for _, fleet := range playerFleets {
		if fleet.BattlePlanNum == battlePlan.Num {
			fleet.BattlePlanNum = 0 // reset to default
			fleetsToUpdate = append(fleetsToUpdate, fleet)
		}
	}

	battlePlans := make([]cs.BattlePlan, 0, len(player.BattlePlans)-1)
	for _, plan := range player.BattlePlans {
		if plan.Num != battlePlan.Num {
			battlePlans = append(battlePlans, plan)
		}
	}
	player.BattlePlans = battlePlans

	// save the updated fleets back to the database
	if err := s.db.WrapInTransaction(func(c db.Client) error {
		for _, fleet := range fleetsToUpdate {
			if err := c.UpdateFleet(fleet); err != nil {
				log.Error().Err(err).Msg("update fleet in database")
				return err
			}
			log.Info().Int64("GameID", game.ID).Int("PlayerNum", player.Num).Int("Num", battlePlan.Num).Msgf("updated fleet %s after deleting BattlePlan", fleet.Name)
		}

		// update fleets in one transaction
		if err := c.UpdatePlayerPlans(player); err != nil {
			log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("update player plans in database")
			return err
		}
		log.Info().Int64("GameID", game.ID).Int("PlayerNum", player.Num).Int("Num", battlePlan.Num).Msgf("deleted BattlePlan %s", battlePlan.Name)

		return nil
	}); err != nil {
		log.Error().Err(err).Msg("update game in database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	// load all the fleets again and return them to the user
	allFleets, err := readWriteClient.GetFleetsForPlayer(game.ID, player.Num)
	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Int("Num", battlePlan.Num).Msg("load fleets from database")
		render.Render(w, r, ErrInternalServerError(err))
		return
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
	rest.RenderJSON(w, rest.JSON{"player": player, "fleets": fleets, "starbases": starbases})
}

func (s *server) createProductionPlan(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	player := s.contextPlayer(r)

	productionPlan := productionPlanRequest{}
	if err := render.Bind(r, &productionPlan); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	if err := productionPlan.Validate(player); err != nil {
		log.Error().Err(err).Int64("PlayerID", player.ID).Str("PlanName", productionPlan.Name).Msg("validate new player ProductionPlan")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	productionPlan.Num = player.GetNextProductionPlanNum()
	player.ProductionPlans = append(player.ProductionPlans, *productionPlan.ProductionPlan)

	if err := db.UpdatePlayerPlans(player); err != nil {
		log.Error().Err(err).Int64("PlayerID", player.ID).Str("PlanName", productionPlan.Name).Msg("save new player ProductionPlan")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, productionPlan)
}

func (s *server) updateProductionPlan(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	player := s.contextPlayer(r)
	existingProductionPlan := s.contextProductionPlan(r)

	productionPlan := productionPlanRequest{}
	if err := render.Bind(r, &productionPlan); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// edge case for bad request where the url num doesn't match the request payload
	if productionPlan.Num != existingProductionPlan.Num {
		log.Error().Int("Num", productionPlan.Num).Msgf("productionPlan.Num %d != existingProductionPlan.Num %d", productionPlan.Num, existingProductionPlan.Num)
		render.Render(w, r, ErrBadRequest(fmt.Errorf("ProductionPlan num does not match existing ProductionPlan num")))
		return
	}

	// validate
	if err := productionPlan.Validate(player); err != nil {
		log.Error().Err(err).Int64("PlayerID", player.ID).Str("PlanName", productionPlan.Name).Msg("validate ProductionPlan")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	for i := range player.ProductionPlans {
		if player.ProductionPlans[i].Num == productionPlan.Num {
			player.ProductionPlans[i] = *productionPlan.ProductionPlan
			break
		}
	}

	if err := db.UpdatePlayerPlans(player); err != nil {
		log.Error().Err(err).Int64("PlayerID", player.ID).Str("PlanName", productionPlan.Name).Msg("save new player ProductionPlan")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, productionPlan)
}

func (s *server) deleteProductionPlan(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	game := s.contextGame(r)
	player := s.contextPlayer(r)
	productionPlan := s.contextProductionPlan(r)

	if productionPlan.Num == 0 {
		log.Error().Int64("PlayerID", player.ID).Str("PlanName", productionPlan.Name).Msg("delete default ProductionPlan")
		render.Render(w, r, ErrBadRequest(fmt.Errorf("cannot delete default production plan")))
		return
	}

	productionPlans := make([]cs.ProductionPlan, 0, len(player.ProductionPlans)-1)
	for _, plan := range player.ProductionPlans {
		if plan.Num != productionPlan.Num {
			productionPlans = append(productionPlans, plan)
		}
	}
	player.ProductionPlans = productionPlans

	// update fleets in one transaction
	if err := db.UpdatePlayerPlans(player); err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("update player plans in database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	// log what we did
	log.Info().Int64("GameID", game.ID).Int("PlayerNum", player.Num).Int("Num", productionPlan.Num).Msgf("deleted ProductionPlan %s", productionPlan.Name)

	rest.RenderJSON(w, player)
}

func (s *server) createTransportPlan(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	player := s.contextPlayer(r)

	transportPlan := transportPlanRequest{}
	if err := render.Bind(r, &transportPlan); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	if err := transportPlan.Validate(player); err != nil {
		log.Error().Err(err).Int64("PlayerID", player.ID).Str("PlanName", transportPlan.Name).Msg("validate new player TransportPlan")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	transportPlan.Num = player.GetNextTransportPlanNum()
	player.TransportPlans = append(player.TransportPlans, *transportPlan.TransportPlan)

	if err := db.UpdatePlayerPlans(player); err != nil {
		log.Error().Err(err).Int64("PlayerID", player.ID).Str("PlanName", transportPlan.Name).Msg("save new player TransportPlan")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, transportPlan)
}

func (s *server) updateTransportPlan(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	player := s.contextPlayer(r)
	existingTransportPlan := s.contextTransportPlan(r)

	transportPlan := transportPlanRequest{}
	if err := render.Bind(r, &transportPlan); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// edge case for bad request where the url num doesn't match the request payload
	if transportPlan.Num != existingTransportPlan.Num {
		log.Error().Int("Num", transportPlan.Num).Msgf("transportPlan.Num %d != existingTransportPlan.Num %d", transportPlan.Num, existingTransportPlan.Num)
		render.Render(w, r, ErrBadRequest(fmt.Errorf("TransportPlan num does not match existing TransportPlan num")))
		return
	}

	// validate
	if err := transportPlan.Validate(player); err != nil {
		log.Error().Err(err).Int64("PlayerID", player.ID).Str("PlanName", transportPlan.Name).Msg("validate TransportPlan")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	for i := range player.TransportPlans {
		if player.TransportPlans[i].Num == transportPlan.Num {
			player.TransportPlans[i] = *transportPlan.TransportPlan
			break
		}
	}

	if err := db.UpdatePlayerPlans(player); err != nil {
		log.Error().Err(err).Int64("PlayerID", player.ID).Str("PlanName", transportPlan.Name).Msg("save new player TransportPlan")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, transportPlan)
}

func (s *server) deleteTransportPlan(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	game := s.contextGame(r)
	player := s.contextPlayer(r)
	transportPlan := s.contextTransportPlan(r)

	if transportPlan.Num == 0 {
		log.Error().Int64("PlayerID", player.ID).Str("PlanName", transportPlan.Name).Msg("delete default TransportPlan")
		render.Render(w, r, ErrBadRequest(fmt.Errorf("cannot delete default transport plan")))
		return
	}

	transportPlans := make([]cs.TransportPlan, 0, len(player.TransportPlans)-1)
	for _, plan := range player.TransportPlans {
		if plan.Num != transportPlan.Num {
			transportPlans = append(transportPlans, plan)
		}
	}
	player.TransportPlans = transportPlans

	// update fleets in one transaction
	if err := db.UpdatePlayerPlans(player); err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("update player plans in database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	// log what we did
	log.Info().Int64("GameID", game.ID).Int("PlayerNum", player.Num).Int("Num", transportPlan.Num).Msgf("deleted TransportPlan %s", transportPlan.Name)

	rest.RenderJSON(w, player)
}
