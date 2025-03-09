//go:build !wasi && !wasm

package server

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
	"github.com/go-pkgz/rest"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
)

type playerOrdersRequest struct {
	*cs.PlayerOrders
}

func (req *playerOrdersRequest) Bind(r *http.Request) error {
	return nil
}

type playerPlansRequest struct {
	*cs.PlayerPlans
}

func (req *playerPlansRequest) Bind(r *http.Request) error {
	return nil
}

type researchCostRequest struct {
	cs.TechLevel
}

func (req *researchCostRequest) Bind(r *http.Request) error {
	return nil
}

type playerRelationsRequest struct {
	*cs.Player
}

func (req *playerRelationsRequest) Bind(r *http.Request) error {
	return nil
}

// context for /api/games/{id} calls that require a player
func (s *server) playerCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := s.contextDb(r)
		user := s.contextUser(r)
		game := s.contextGame(r)

		player, err := db.GetLightPlayerForGame(game.ID, user.ID)
		if err != nil {
			render.Render(w, r, ErrInternalServerError(err))
			return
		}

		if player == nil {
			log.Error().Int64("GameID", game.ID).Int64("UserID", user.ID).Msg("player not found")
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), keyPlayer, player)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *server) contextPlayer(r *http.Request) *cs.Player {
	return r.Context().Value(keyPlayer).(*cs.Player)
}

func (s *server) player(w http.ResponseWriter, r *http.Request) {
	player := s.contextPlayer(r)
	rest.RenderJSON(w, player)
}

func (s *server) playerIntels(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	user := s.contextUser(r)
	game := s.contextGame(r)
	intels, err := db.GetPlayerIntelsForGame(game.ID, user.ID)

	if err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, intels)
}

func (s *server) fullPlayer(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	user := s.contextUser(r)
	game := s.contextGame(r)

	player, err := db.GetPlayerForGame(game.ID, user.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	if player == nil {
		render.Render(w, r, ErrNotFound)
		return
	}

	rest.RenderJSON(w, player)
}

// get mapObjects for a player
func (s *server) mapObjects(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	user := s.contextUser(r)

	gameID, err := s.int64URLParam(r, "id")
	if gameID == nil || err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	mapObjects, err := db.GetPlayerMapObjects(*gameID, user.ID)
	if err != nil {
		log.Error().Err(err).Int64("GameID", *gameID).Int64("UserID", user.ID).Msg("load player map objects database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	if mapObjects == nil {
		render.Render(w, r, ErrNotFound)
		return
	}

	rest.RenderJSON(w, mapObjects)
}

// data about a universe (planets, fleets, designs, other players, etc) for a single player in the game
// this aggregates player objects (full planets/fleets/mineralPackets) and intel objects
type playerUniverseResponse struct {
	cs.PlayerIntels
	Planets        []*cs.Planet        `json:"planets,omitempty"`
	Fleets         []*cs.Fleet         `json:"fleets,omitempty"`
	Starbases      []*cs.Fleet         `json:"starbases,omitempty"`
	MineFields     []*cs.MineField     `json:"mineFields,omitempty"`
	MineralPackets []*cs.MineralPacket `json:"mineralPackets,omitempty"`
	Designs        []*cs.ShipDesign    `json:"designs,omitempty"`
}

// get mapObjects for a player
func (s *server) universe(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	user := s.contextUser(r)
	game := s.contextGame(r)

	player, err := db.GetPlayerForGame(game.ID, user.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	pmos, err := db.GetPlayerMapObjects(game.ID, user.ID)
	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Int64("UserID", user.ID).Msg("load player map objects database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	if pmos == nil {
		render.Render(w, r, ErrNotFound)
		return
	}

	intels, err := db.GetPlayerIntelsForGame(game.ID, user.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	if intels == nil {
		render.Render(w, r, ErrNotFound)
		return
	}

	designs, err := db.GetShipDesignsForPlayer(game.ID, player.Num)
	if err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	player.PlayerIntels = *intels
	player.Designs = designs

	universe := buildUniverse(&cs.FullPlayer{
		Player:           *player,
		PlayerMapObjects: *pmos,
	})

	rest.RenderJSON(w, universe)
}

// build a universe response
func buildUniverse(player *cs.FullPlayer) playerUniverseResponse {

	universe := playerUniverseResponse{
		PlayerIntels:   player.PlayerIntels,
		Planets:        make([]*cs.Planet, len(player.Planets)),
		Fleets:         make([]*cs.Fleet, len(player.Fleets)),
		Starbases:      make([]*cs.Fleet, len(player.Starbases)),
		MineFields:     make([]*cs.MineField, len(player.MineFields)),
		MineralPackets: make([]*cs.MineralPacket, len(player.MineralPackets)),
		Designs:        make([]*cs.ShipDesign, len(player.Designs)),
	}

	// merge player and design intels into the Designs data
	copy(universe.Planets, player.Planets)
	copy(universe.Fleets, player.Fleets)
	copy(universe.Starbases, player.Starbases)
	copy(universe.MineFields, player.MineFields)
	copy(universe.MineralPackets, player.MineralPackets)
	copy(universe.Designs, player.Designs)

	if player.Num <= len(universe.ScoreIntels) {
		universe.ScoreIntels[player.Num-1].ScoreHistory = player.ScoreHistory
	}

	return universe
}

// submit a player turn and return the newly generated turn if there is one
func (s *server) submitTurn(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	game := s.contextGame(r)
	player := s.contextPlayer(r)

	// submit the turn
	player.SubmittedTurn = true
	if err := db.SubmitPlayerTurn(player.GameID, player.Num, true); err != nil {
		log.Error().Err(err).Int64("GameID", player.GameID).Int("PlayerNum", player.Num).Msg("update player")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	// only allow one CheckAndGenerate to run at a time
	// TODO: handle this differently if you ever scale out beyond one instance. :)
	result, err, _ := s.sf.Do(strconv.FormatInt(game.ID, 10), func() (interface{}, error) {
		gr := s.newGameRunner()
		result, err := gr.CheckAndGenerateTurn(player.GameID)
		if err != nil {
			return nil, err
		}
		return result, nil
	})

	if err != nil {
		log.Error().Err(err).Int64("GameID", player.GameID).Msg("check and generate new turn")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	if result == TurnGenerated {
		s.sendNewTurnNotification(r, game.ID)
		s.renderFullPlayerGame(w, r, player.GameID, player.UserID)
		return
	}

	// return the game status
	game, err = db.GetGame(player.GameID)
	if err != nil {
		log.Error().Err(err).Int64("GameID", player.GameID).Msg("load game")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, rest.JSON{"game": game, "player": player})
}

// submit a player turn and return the newly generated turn if there is one
func (s *server) unSubmitTurn(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	player := s.contextPlayer(r)

	// submit the turn
	player.SubmittedTurn = false
	if err := db.SubmitPlayerTurn(player.GameID, player.Num, false); err != nil {
		log.Error().Err(err).Int64("GameID", player.GameID).Int("PlayerNum", player.Num).Msg("update player")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	rest.RenderJSON(w, rest.JSON{"player": player})
}

func (s *server) renderFullPlayerGame(w http.ResponseWriter, r *http.Request, gameID, userID int64) {
	// return a new turn
	gr := s.newGameRunner()
	game, fullPlayer, err := gr.LoadPlayerGame(gameID, userID)
	if err != nil {
		log.Error().Err(err).Int64("GameID", game.ID).Msg("load full game from database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	universe := buildUniverse(fullPlayer)

	rest.RenderJSON(w, rest.JSON{"game": game, "player": fullPlayer.Player, "universe": universe})
}

// Update a player's orders (research field, research amount)
func (s *server) updatePlayerOrders(w http.ResponseWriter, r *http.Request) {
	readWriteClient := s.contextDb(r)
	game := s.contextGame(r)
	player := s.contextPlayer(r)

	orders := playerOrdersRequest{}
	if err := render.Bind(r, &orders); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	if orders.ResearchAmount < 0 || orders.ResearchAmount > 100 {
		render.Render(w, r, ErrBadRequest(fmt.Errorf("research ammount must be between 0 and 100")))
		return
	}

	planets, err := readWriteClient.GetPlanetsForPlayer(player.GameID, player.Num)
	if err != nil {
		log.Error().Err(err).Int64("ID", player.ID).Msg("loading player planets from database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	// load this player but with designs so the update works correctly
	player, err = readWriteClient.GetPlayerWithDesignsForGame(game.ID, player.Num)
	if err != nil {
		log.Error().Err(err).Int64("ID", player.ID).Msg("loading player from database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	orderer := cs.NewOrderer()
	orders.CargoTransfers = player.CargoTransfers // don't let the player update cargoTransfers outside of a fleet transfer function
	orderer.UpdatePlayerOrders(player, planets, *orders.PlayerOrders, &game.Rules)

	// save the updated fleets back to the database
	if err := s.db.WrapInTransaction(func(c db.Client) error {
		// save the player to the database
		if err := c.UpdatePlayerOrders(player); err != nil {
			log.Error().Err(err).Int64("GameID", player.GameID).Int("PlayerNum", player.Num).Msg("update player")
			return err
		}

		for _, planet := range planets {
			if planet.Dirty {
				// TODO: only update the planet spec? that's all that changes
				if err := c.UpdatePlanet(planet); err != nil {
					log.Error().Err(err).Int64("ID", player.ID).Msg("updating player planet in database")
					return err
				}
			}
		}
		return nil
	}); err != nil {
		log.Error().Err(err).Msg("update game in database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	log.Info().Int64("GameID", player.GameID).Int("PlayerNum", player.Num).Msg("update orders")
	rest.RenderJSON(w, rest.JSON{"player": player, "planets": planets})
}

// update the player's relations with other players
func (s *server) updatePlayerRelations(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	player := s.contextPlayer(r)

	relationsRequest := playerRelationsRequest{}
	if err := render.Bind(r, &relationsRequest); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	if len(relationsRequest.Relations) != len(player.Relations) {
		render.Render(w, r, ErrBadRequest(fmt.Errorf("must include all player relations")))
		return
	}

	// save the player to the database
	player.Relations = relationsRequest.Relations
	if err := db.UpdatePlayerRelations(player); err != nil {
		log.Error().Err(err).Int64("GameID", player.GameID).Int("PlayerNum", player.Num).Msg("update player relations")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	log.Info().Int64("GameID", player.GameID).Int("PlayerNum", player.Num).Msg("update relations")
	rest.RenderJSON(w, player.Relations)
}

// Update a player's plans
func (s *server) updatePlayerPlans(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
	player := s.contextPlayer(r)

	plans := playerPlansRequest{}
	if err := render.Bind(r, &plans); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	if len(plans.BattlePlans) == 0 {
		render.Render(w, r, ErrBadRequest(fmt.Errorf("must have at least one battle plan")))
		return
	}

	if plans.BattlePlans[0].Num != 0 {
		render.Render(w, r, ErrBadRequest(fmt.Errorf("must have a default battle plan")))
		return
	}

	// TODO: validate?
	// TODO: convert creates into a separate POST?
	// TODO: update fleets with deleted battle plans to use default battleplan
	nextNum := 0
	for i := range plans.BattlePlans {
		nextNum = int(math.Max(float64(plans.BattlePlans[i].Num+1), float64(nextNum)))
	}

	for i := range plans.BattlePlans {
		if plans.BattlePlans[i].Num == -1 {
			plans.BattlePlans[i].Num = nextNum
			nextNum++
		}
	}

	player.PlayerPlans = *plans.PlayerPlans

	// save the player to the database
	if err := db.UpdatePlayerPlans(player); err != nil {
		log.Error().Err(err).Int64("GameID", player.GameID).Int("PlayerNum", player.Num).Msg("update player")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	log.Info().Int64("GameID", player.GameID).Int("PlayerNum", player.Num).Msg("update plans")
	rest.RenderJSON(w, player)
}

// get an estimate for production completion based on a planet's production queue items
func (s *server) getResearchCost(w http.ResponseWriter, r *http.Request) {
	game := s.contextGame(r)
	player := s.contextPlayer(r)

	researchCost := researchCostRequest{}
	if err := render.Bind(r, &researchCost); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	resources := player.GetResearchCost(&game.Rules, researchCost.TechLevel)
	rest.RenderJSON(w, rest.JSON{"resources": resources})
}
