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

type raceRequest struct {
	*cs.Race
}

func (req *raceRequest) Bind(r *http.Request) error {
	return nil
}

// context for /api/races/{id} calls
func (s *server) raceCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := s.contextUser(r)

		// load the race by id from the database
		id, err := s.int64URLParam(r, "id")
		if id == nil || err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		race, err := s.db.GetRace(*id)
		if err != nil {
			render.Render(w, r, ErrInternalServerError(err))
			return
		}

		if race == nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		if race.UserID != user.ID {
			render.Render(w, r, ErrForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), keyRace, race)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *server) contextRace(r *http.Request) *cs.Race {
	return r.Context().Value(keyRace).(*cs.Race)
}

func (s *server) races(w http.ResponseWriter, r *http.Request) {
	user := s.contextUser(r)

	races, err := s.db.GetRacesForUser(user.ID)
	if err != nil {
		log.Error().Err(err).Int64("UserID", user.ID).Msg("get races from database")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	rest.RenderJSON(w, races)
}

func (s *server) race(w http.ResponseWriter, r *http.Request) {
	race := s.contextRace(r)
	rest.RenderJSON(w, race)
}

// create a new race for a user
func (s *server) createRace(w http.ResponseWriter, r *http.Request) {
	user := s.contextUser(r)

	race := raceRequest{}
	if err := render.Bind(r, &race); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	race.UserID = user.ID
	if err := s.db.CreateRace(race.Race); err != nil {
		log.Error().Err(err).Int64("UserID", user.ID).Msg("create race")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	rest.RenderJSON(w, race)
}

// get points for a race
func (s *server) getRacePoints(w http.ResponseWriter, r *http.Request) {

	race := raceRequest{}
	if err := render.Bind(r, &race); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// compute points
	points := race.ComputeRacePoints(cs.NewRules().RaceStartingPoints)
	rest.RenderJSON(w, rest.JSON{"points": points})
}

func (s *server) updateRace(w http.ResponseWriter, r *http.Request) {
	race := raceRequest{}
	if err := render.Bind(r, &race); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// load in the existing race from the context
	existingRace := s.contextRace(r)

	// validate
	if race.ID != existingRace.ID || race.UserID != existingRace.UserID {
		log.Error().Int64("ID", race.ID).Msgf("race.ID %d != existingRace.ID %d or race.UserID  %d != existingRace.UserID %d", race.ID, existingRace.ID, race.UserID, existingRace.UserID)
		render.Render(w, r, ErrBadRequest(fmt.Errorf("race id/user id does not match existing race")))
		return
	}

	if err := s.db.UpdateRace(race.Race); err != nil {
		log.Error().Err(err).Int64("ID", race.ID).Msg("update race in database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, race)
}

func (s *server) deleteRace(w http.ResponseWriter, r *http.Request) {
	race := s.contextRace(r)

	if err := s.db.DeleteRace(race.ID); err != nil {
		log.Error().Err(err).Int64("ID", race.ID).Msg("delete race from database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}
}
