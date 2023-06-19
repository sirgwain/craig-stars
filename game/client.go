package game

import (
	"log"
	"time"
)

type client struct {
}

type Client interface {
	CreateGame(hostID int64, settings GameSettings) Game
	NewPlayer(userID int64, race Race) *Player
	GenerateUniverse(game *Game, players []*Player) (*Universe, error)
	SubmitTurn(player *Player)
	GenerateTurn(game *FullGame) error
}

func NewClient() client {
	return client{}
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
	player := NewPlayer(userID, &race)
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
	turnGenerator := NewTurnGenerator(&FullGame{game, universe, players})
	return turnGenerator.GenerateTurn()
}
