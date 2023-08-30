package server

import (
	"context"
	"fmt"
	"math"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-pkgz/rest"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
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
	Planets        []interface{} `json:"planets,omitempty"`
	Fleets         []interface{} `json:"fleets,omitempty"`
	Starbases      []interface{} `json:"starbases,omitempty"`
	Wormholes      []interface{} `json:"wormholes,omitempty"`
	MineralPackets []interface{} `json:"mineralPackets,omitempty"`
	MineFields     []interface{} `json:"mineFields,omitempty"`
	MysteryTraders []interface{} `json:"mysteryTraders,omitempty"`
	Salvages       []interface{} `json:"salvages,omitempty"`
	Designs        []interface{} `json:"designs,omitempty"`
	Players        []interface{} `json:"players,omitempty"`
	Scores         []interface{} `json:"scores,omitempty"`
	Battles        []interface{} `json:"battles,omitempty"`
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

	universe := buildUniverse(player, designs, *pmos, *intels)

	rest.RenderJSON(w, universe)
}

// build a universe response
func buildUniverse(player *cs.Player, designs []*cs.ShipDesign, pmos cs.PlayerMapObjects, intels cs.PlayerIntels) playerUniverseResponse {
	numPlayerFleets := len(pmos.Fleets)
	numPlayerStarbases := len(pmos.Starbases)
	numPlayerMineralPackets := len(pmos.MineralPackets)
	numPlayerMineFields := len(pmos.MineFields)
	numPlayerDesigns := len(designs)

	universe := playerUniverseResponse{
		Planets:        make([]interface{}, len(intels.PlanetIntels)),
		Fleets:         make([]interface{}, len(intels.FleetIntels)+numPlayerFleets),
		Starbases:      make([]interface{}, len(intels.StarbaseIntels)+numPlayerStarbases),
		MineralPackets: make([]interface{}, len(intels.MineralPacketIntels)+numPlayerMineralPackets),
		MineFields:     make([]interface{}, len(intels.MineFieldIntels)+numPlayerMineFields),
		Salvages:       make([]interface{}, len(intels.SalvageIntels)),
		Wormholes:      make([]interface{}, len(intels.WormholeIntels)),
		MysteryTraders: make([]interface{}, len(intels.MysteryTraderIntels)),
		Designs:        make([]interface{}, len(intels.ShipDesignIntels)+numPlayerDesigns),
		Players:        make([]interface{}, len(intels.PlayerIntels)),
		Scores:         make([]interface{}, len(intels.ScoreIntels)),
		Battles:        make([]interface{}, len(intels.BattleRecords)),
	}

	// merge player and design intels into the Designs data
	for i, item := range designs {
		universe.Designs[i] = item
	}
	for i, item := range intels.ShipDesignIntels {
		universe.Designs[i+numPlayerDesigns] = item
	}

	for i, item := range intels.PlayerIntels {
		universe.Players[i] = item
	}

	for i, item := range intels.ScoreIntels {
		universe.Scores[i] = item.ScoreHistory
	}
	if player.Num < len(universe.Scores) {
		universe.Scores[player.Num-1] = player.ScoreHistory
	}

	for i, item := range intels.BattleRecords {
		universe.Battles[i] = item
	}

	// merge player planets and planet intels
	for i, item := range intels.PlanetIntels {
		universe.Planets[i] = item
	}
	// we overwrite planets by num
	for _, item := range pmos.Planets {
		universe.Planets[item.Num-1] = item
	}

	// start with player objects, then append intel objects
	for i, item := range pmos.Fleets {
		universe.Fleets[i] = item
	}
	for i, item := range intels.FleetIntels {
		universe.Fleets[i+numPlayerFleets] = item
	}

	for i, item := range pmos.Starbases {
		universe.Starbases[i] = item
	}
	for i, item := range intels.StarbaseIntels {
		universe.Starbases[i+numPlayerStarbases] = item
	}

	for i, item := range pmos.MineralPackets {
		universe.MineralPackets[i] = item
	}
	for i, item := range intels.MineralPacketIntels {
		universe.MineralPackets[i+numPlayerMineralPackets] = item
	}

	for i, item := range pmos.MineFields {
		universe.MineFields[i] = item
	}
	for i, item := range intels.MineFieldIntels {
		universe.MineFields[i+numPlayerMineFields] = item
	}

	for i, item := range intels.SalvageIntels {
		universe.Salvages[i] = item
	}

	for i, item := range intels.WormholeIntels {
		universe.Wormholes[i] = item
	}

	for i, item := range intels.MysteryTraderIntels {
		universe.MysteryTraders[i] = item
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

	// TODO: this should probably be a goroutine or something
	gr := s.newGameRunner()
	result, err := gr.CheckAndGenerateTurn(player.GameID)
	if err != nil {
		log.Error().Err(err).Int64("GameID", player.GameID).Msg("check and generate new turn")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	if result == TurnGenerated {
		s.sendNewTurnNotification(db, game.ID)
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

	universe := buildUniverse(&fullPlayer.Player, fullPlayer.Designs, fullPlayer.PlayerMapObjects, fullPlayer.PlayerIntels)

	// don't clutter our response
	// TODO: do this fetching more elegantly
	fullPlayer.Player.Designs = nil
	fullPlayer.Player.PlayerIntels = cs.PlayerIntels{}

	rest.RenderJSON(w, rest.JSON{"game": game, "player": fullPlayer.Player, "universe": universe})
}

// Submit a turn for the player
func (s *server) updatePlayerOrders(w http.ResponseWriter, r *http.Request) {
	db := s.contextDb(r)
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

	planets, err := db.GetPlanetsForPlayer(player.GameID, player.Num)
	if err != nil {
		log.Error().Err(err).Int64("ID", player.ID).Msg("loading player planets from database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	// load this player but with designs so the update works correctly
	player, err = db.GetPlayerWithDesignsForGame(game.ID, player.Num)
	if err != nil {
		log.Error().Err(err).Int64("ID", player.ID).Msg("loading player from database")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	orderer := cs.NewOrderer()
	orderer.UpdatePlayerOrders(player, planets, *orders.PlayerOrders, &game.Rules)

	// save the player to the database
	if err := db.UpdatePlayerOrders(player); err != nil {
		log.Error().Err(err).Int64("GameID", player.GameID).Int("PlayerNum", player.Num).Msg("update player")
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	for _, planet := range planets {
		if planet.Dirty {
			// TODO: only update the planet spec? that's all that changes
			// TODO: do this all in one transaction?
			if err := db.UpdatePlanet(planet); err != nil {
				log.Error().Err(err).Int64("ID", player.ID).Msg("updating player planet in database")
				render.Render(w, r, ErrInternalServerError(err))
				return
			}
		}
	}

	log.Info().Int64("GameID", player.GameID).Int("PlayerNum", player.Num).Msg("update orders")
	rest.RenderJSON(w, rest.JSON{"player": player, "planets": planets})
}

// Submit a turn for the player
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
