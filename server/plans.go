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

type battlePlanRequest struct {
	*cs.BattlePlan
}

func (req *battlePlanRequest) Bind(r *http.Request) error {
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

func (s *server) contextBattlePlan(r *http.Request) *cs.BattlePlan {
	return r.Context().Value(keyBattlePlan).(*cs.BattlePlan)
}

func (s *server) createBattlePlan(w http.ResponseWriter, r *http.Request) {
	player := s.contextPlayer(r)

	battlePlan := battlePlanRequest{}
	if err := render.Bind(r, &battlePlan); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	if err := battlePlan.Validate(player); err != nil {
		log.Error().Err(err).Int64("PlayerID", player.ID).Str("PlanName", battlePlan.Name).Msg("validate new player battlePlan")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	battlePlan.Num = player.GetNextBattlePlanNum()
	player.BattlePlans = append(player.BattlePlans, *battlePlan.BattlePlan)

	if err := s.db.UpdatePlayerPlans(player); err != nil {
		log.Error().Err(err).Int64("PlayerID", player.ID).Str("PlanName", battlePlan.Name).Msg("save new player plan")
		render.Render(w, r, ErrInternalServerError(err))
	}

	rest.RenderJSON(w, battlePlan)
}

func (s *server) updateBattlePlan(w http.ResponseWriter, r *http.Request) {
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
		render.Render(w, r, ErrBadRequest(fmt.Errorf("battlePlan num does not match existing battlePlan num")))
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

	if err := s.db.UpdatePlayerPlans(player); err != nil {
		log.Error().Err(err).Int64("PlayerID", player.ID).Str("PlanName", battlePlan.Name).Msg("save new player plan")
		render.Render(w, r, ErrInternalServerError(err))
	}

	rest.RenderJSON(w, battlePlan)
}

func (s *server) deleteBattlePlan(w http.ResponseWriter, r *http.Request) {
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
	playerFleets, err := s.db.GetFleetsForPlayer(game.ID, player.Num)
	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("load fleets from database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	fleetsToUpdate := []*cs.Fleet{}
	for _, fleet := range playerFleets {
		if fleet.BattlePlanNum == battlePlan.Num {
			fleet.BattlePlanNum = 0 // reset to default
			fleet.Dirty = true
			fleetsToUpdate = append(fleetsToUpdate, fleet)
		}
	}

	// update fleets in one transaction
	if err := s.db.CreateUpdateOrDeleteFleets(game.ID, fleetsToUpdate); err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("update fleets in database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	battlePlans := make([]cs.BattlePlan, 0, len(player.BattlePlans)-1)
	for _, plan := range player.BattlePlans {
		if plan.Num != battlePlan.Num {
			battlePlans = append(battlePlans, plan)
		}
	}
	player.BattlePlans = battlePlans

	// update fleets in one transaction
	if err := s.db.UpdatePlayerPlans(player); err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Msg("update player plans in database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	// log what we did
	log.Info().Int64("GameID", game.ID).Int("PlayerNum", player.Num).Int("Num", battlePlan.Num).Msgf("deleted BattlePlan %s", battlePlan.Name)

	for _, fleet := range fleetsToUpdate {
		log.Info().Int64("GameID", game.ID).Int("PlayerNum", player.Num).Int("Num", battlePlan.Num).Msgf("updated fleet %s after deleting plan", fleet.Name)
	}

	allFleets, err := s.db.GetFleetsForPlayer(game.ID, player.Num)

	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Int("Num", battlePlan.Num).Msg("load fleets from database")
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
	rest.RenderJSON(w, rest.JSON{"player": player, "fleets": fleets, "starbases": starbases})
}
