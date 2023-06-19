package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_generateTurn(t *testing.T) {
	client := NewClient()
	game := client.CreateGame(1, *NewGameSettings())
	player := client.NewPlayer(1, *NewRace(), &game.Rules)
	players := []*Player{player}
	player.AIControlled = true
	player.Num = 1
	universe, _ := client.GenerateUniverse(&game, players)

	startingFleets := len(universe.Fleets)

	client.GenerateTurn(&game, universe, players)

	assert.Equal(t, 2401, game.Year)

	// should have intel about planets
	assert.Equal(t, len(universe.Planets), len(player.PlanetIntels))

	// should have built a new scout
	assert.Greater(t, len(universe.Fleets), startingFleets)

	// should have grown pop
	assert.Greater(t, universe.Planets[0].population(), player.Race.Spec.StartingPlanets[0].Population)
}

func Test_generateTurns(t *testing.T) {
	client := NewClient()
	game := client.CreateGame(1, *NewGameSettings())
	player := client.NewPlayer(1, *NewRace(), &game.Rules)
	player.AIControlled = true
	player.Num = 1
	players := []*Player{player}
	universe, _ := client.GenerateUniverse(&game, players)

	// generate many turns
	for i := 0; i < 100; i++ {
		client.GenerateTurn(&game, universe, players)
	}

	assert.Equal(t, 2500, game.Year)

	// should have fleets
	assert.True(t, len(universe.Fleets) > 0)

	// should have grown pop
	assert.Greater(t, universe.Planets[0].population(), player.Race.Spec.StartingPlanets[0].Population)

	// should have built factories
	assert.Greater(t, universe.Planets[0].Factories, game.Rules.StartingFactories)

}
