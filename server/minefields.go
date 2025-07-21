//go:build !wasi && !wasm

package server

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
)

type minefieldRequest struct {
	*cs.Minefield
}

func (req *minefieldRequest) Bind(r *http.Request) error {
	return nil
}

// context for /api/games/{id}/minefields/{num} calls that require a shipDesign
func (s *server) minefieldCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := s.contextDb(r)
		player := s.contextPlayer(r)

		num, err := s.intURLParam(r, "num")
		if num == nil || err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		minefield, err := db.GetMinefieldByNum(r.Context(), player.GameID, player.Num, *num)
		if err != nil {
			render.Render(w, r, ErrInternalServerError(err))
			return
		}

		if minefield == nil {
			log.Error().Int64("GameID", player.GameID).Msgf("unable to find minefield %d", num)
			render.Render(w, r, ErrNotFound)
			return
		}

		// only minefield owners can load this minefield
		if minefield.PlayerNum != player.Num {
			render.Render(w, r, ErrForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), keyMinefield, minefield)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *server) contextMinefield(r *http.Request) *cs.Minefield {
	return r.Context().Value(keyMinefield).(*cs.Minefield)
}

func (s *server) minefield(w http.ResponseWriter, r *http.Request) {
	minefield := s.contextMinefield(r)
	RenderJSON(w, minefield)
}

// Allow a user to update a minefield's orders
func (s *server) updateMinefieldOrders(w http.ResponseWriter, r *http.Request) {
	existingMinefield := s.contextMinefield(r)
	player := s.contextPlayer(r)
	game := s.contextGame(r)

	minefield := minefieldRequest{}
	if err := render.Bind(r, &minefield); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	orderer := cs.NewOrderer()
	if err := orderer.UpdateMinefieldOrders(player, existingMinefield, minefield.MinefieldOrders); err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Str("Minefield", existingMinefield.Name).Msg("update minefield orders")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// update this minefield and the player's spec in the database
	if err := s.db.WrapInTransaction(func(c db.Client) error {
		if err := c.SaveMinefield(r.Context(), existingMinefield); err != nil {
			log.Error().Err(err).Int64("ID", minefield.ID).Msg("update minefield in database")
			return err
		}

		return nil
	}); err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	RenderJSON(w, JSON{"minefield": existingMinefield, "player": player})
}
