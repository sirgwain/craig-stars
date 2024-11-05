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

type mineFieldRequest struct {
	*cs.MineField
}

func (req *mineFieldRequest) Bind(r *http.Request) error {
	return nil
}

// context for /api/games/{id}/mineFields/{num} calls that require a shipDesign
func (s *server) mineFieldCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := s.contextDb(r)
		player := s.contextPlayer(r)

		num, err := s.intURLParam(r, "num")
		if num == nil || err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}

		mineField, err := db.GetMineFieldByNum(player.GameID, player.Num, *num)
		if err != nil {
			render.Render(w, r, ErrInternalServerError(err))
			return
		}

		if mineField == nil {
			log.Error().Int64("GameID", player.GameID).Msgf("unable to find mineField %d", num)
			render.Render(w, r, ErrNotFound)
			return
		}

		// only mineField owners can load this mineField
		if mineField.PlayerNum != player.Num {
			render.Render(w, r, ErrForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), keyMineField, mineField)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *server) contextMineField(r *http.Request) *cs.MineField {
	return r.Context().Value(keyMineField).(*cs.MineField)
}

func (s *server) mineField(w http.ResponseWriter, r *http.Request) {
	mineField := s.contextMineField(r)
	rest.RenderJSON(w, mineField)
}

// Allow a user to update a mineField's orders
func (s *server) updateMineFieldOrders(w http.ResponseWriter, r *http.Request) {
	existingMineField := s.contextMineField(r)
	player := s.contextPlayer(r)
	game := s.contextGame(r)

	mineField := mineFieldRequest{}
	if err := render.Bind(r, &mineField); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	orderer := cs.NewOrderer()
	if err := orderer.UpdateMineFieldOrders(player, existingMineField, mineField.MineFieldOrders); err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int("PlayerNum", player.Num).Str("MineField", existingMineField.Name).Msg("update mineField orders")
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// update this mineField and the player's spec in the database
	if err := s.db.WrapInTransaction(func(c db.Client) error {
		if err := c.UpdateMineField(existingMineField); err != nil {
			log.Error().Err(err).Int64("ID", mineField.ID).Msg("update mineField in database")
			return err
		}

		return nil
	}); err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, rest.JSON{"mineField": existingMineField, "player": player})
}
