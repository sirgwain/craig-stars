package game

import (
	"log"
	"time"
)

type client struct {
	orders orderer
}

// external interface for creating/interacting with game objects
type Client interface {

	// game creation
	CreateGame(hostID int64, settings GameSettings) Game
	NewPlayer(userID int64, race Race, rules *Rules) *Player

	// player orders to planets, fleets, research, etc
	UpdatePlayerOrders(player *Player, playerPlanets []*Planet, orders PlayerOrders)
	UpdatePlanetOrders(player *Player, planet *Planet, orders PlanetOrders)
	UpdateFleetOrders(player *Player, fleet *Fleet, orders FleetOrders)
	UpdateMineFieldOrders(player *Player, minefield *MineField, orders MineFieldOrders)
	TransferFleetCargo(source *Fleet, dest *Fleet, transferAmount Cargo) error
	TransferPlanetCargo(source *Fleet, dest *Planet, transferAmount Cargo) error

	// universe/turn generation
	GenerateUniverse(game *Game, players []*Player) (*Universe, error)
	SubmitTurn(player *Player)
	GenerateTurn(game *Game, universe *Universe, players []*Player) error
	CheckAllPlayersSubmitted(players []*Player) bool
}

func NewClient() Client {
	return &client{orders: &orders{}}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func (c *client) CreateGame(hostID int64, settings GameSettings) Game {
	g := NewGame().WithSettings(settings)
	g.HostID = hostID

	return *g
}

// create a new player
func (c *client) NewPlayer(userID int64, race Race, rules *Rules) *Player {
	player := newPlayer(userID, &race)
	player.Race.Spec = computeRaceSpec(&player.Race, rules)

	return player
}

// Generate a new universe
func (c *client) GenerateUniverse(game *Game, players []*Player) (*Universe, error) {
	defer timeTrack(time.Now(), "GenerateUniverse")

	ug := NewUniverseGenerator(game.Size, game.Density, players, &game.Rules)
	universe, err := ug.Generate()

	if err != nil {
		return nil, err
	}

	// save our area back to the game object now that it's been generated
	game.Area = ug.Area()

	return universe, nil
}

func (c *client) SubmitTurn(player *Player) {
	// TODO: anything else to do on turn submit?
	player.SubmittedTurn = true
}

// check if all players have submitted their turn
func (c *client) CheckAllPlayersSubmitted(players []*Player) bool {
	for _, player := range players {
		if !player.SubmittedTurn {
			return false
		}
	}
	return true
}

// generate a new turn for this game
func (c *client) GenerateTurn(game *Game, universe *Universe, players []*Player) error {
	defer timeTrack(time.Now(), "GenerateTurn")
	turnGenerator := newTurnGenerator(&FullGame{game, universe, players})
	return turnGenerator.generateTurn()
}

func (c *client) UpdatePlayerOrders(player *Player, playerPlanets []*Planet, orders PlayerOrders) {
	c.orders.updatePlayerOrders(player, playerPlanets, orders)
}

func (c *client) UpdatePlanetOrders(player *Player, planet *Planet, orders PlanetOrders) {
	c.orders.updatePlanetOrders(player, planet, orders)
}

func (c *client) UpdateFleetOrders(player *Player, fleet *Fleet, orders FleetOrders) {
	c.orders.updateFleetOrders(player, fleet, orders)
}

func (c *client) UpdateMineFieldOrders(player *Player, minefield *MineField, orders MineFieldOrders) {
	c.orders.updateMineFieldOrders(player, minefield, orders)
}

func (c *client) TransferFleetCargo(source *Fleet, dest *Fleet, transferAmount Cargo) error {
	return c.orders.transferFleetCargo(source, dest, transferAmount)
}

func (c *client) TransferPlanetCargo(source *Fleet, dest *Planet, transferAmount Cargo) error {
	return c.orders.transferPlanetCargo(source, dest, transferAmount)
}
