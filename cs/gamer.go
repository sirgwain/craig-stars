package cs

import (
	"log"
	"time"
)

type client struct {
}

// external interface for creating/interacting with games
type Gamer interface {

	// game creation
	CreateGame(hostID int64, settings GameSettings) *Game
	NewPlayer(userID int64, race Race, rules *Rules) *Player

	// universe/turn generation
	GenerateUniverse(game *Game, players []*Player) (*Universe, error)
	SubmitTurn(player *Player)
	GenerateTurn(game *Game, universe *Universe, players []*Player) error
	CheckAllPlayersSubmitted(players []*Player) bool
}

func NewGamer() Gamer {
	return &client{}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func (c *client) CreateGame(hostID int64, settings GameSettings) *Game {
	game := NewGame().WithSettings(settings)
	game.HostID = hostID

	return game
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

	ug := NewUniverseGenerator(game, players)
	universe, err := ug.Generate()

	if err != nil {
		return nil, err
	}

	// save our area back to the game object now that it's been generated
	game.Area = ug.Area()

	return universe, nil
}

func (c *client) SubmitTurn(player *Player) {
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
	turnGenerator := newTurnGenerator(&FullGame{game, universe, game.Rules.techs, players})
	return turnGenerator.generateTurn()
}
