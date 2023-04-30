package server

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-pkgz/rest"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
)

type planetRequest struct {
	*cs.Planet
}

func (req *planetRequest) Bind(r *http.Request) error {
	return nil
}

// context for /api/games/{id}/planets/{num} calls that require a shipDesign
func (s *server) planetCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		player := s.contextPlayer(r)

		num, err := s.intURLParam(r, "num")
		if num == nil || err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		planet, err := s.db.GetPlanetByNum(player.GameID, *num)
		if err != nil {
			render.Render(w, r, ErrInternalServerError(err))
			return
		}

		if planet == nil {
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
	existingPlanet := s.contextPlanet(r)
	player := s.contextPlayer(r)

	planet := planetRequest{}
	if err := render.Bind(r, &planet); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	orderer := cs.NewOrderer()
	orderer.UpdatePlanetOrders(player, existingPlanet, planet.PlanetOrders)

	if err := s.db.UpdatePlanet(existingPlanet); err != nil {
		log.Error().Err(err).Int64("ID", planet.ID).Msg("update planet in database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, existingPlanet)
}
