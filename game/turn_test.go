package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_generateTurn(t *testing.T) {
	game := NewGame()
	game.AddPlayer(NewPlayer(1, NewRace()))
	player := &game.Players[0]
	player.AIControlled = true

	game.GenerateUniverse()

	startingFleets := len(player.Fleets)

	game.GenerateTurn()

	assert.Equal(t, 2401, game.Year)

	// should have intel about planets
	assert.Equal(t, len(game.Planets), len(player.PlanetIntels))

	// should have built a new scout
	assert.Greater(t, len(player.Fleets), startingFleets)

	// should have grown pop
	assert.Greater(t, player.Planets[0].Population(), player.Race.Spec.StartingPlanets[0].Population)
}

func Test_generateTurns(t *testing.T) {
	game := NewGame()
	game.AddPlayer(NewPlayer(1, NewRace()))
	player := &game.Players[0]
	player.AIControlled = true

	game.GenerateUniverse()

	planetMinerals := player.Planets[0].Cargo.ToMineral()

	// generate many turns
	for i := 0; i < 100; i++ {
		game.GenerateTurn()
	}

	assert.Equal(t, 2500, game.Year)
	// should have intel about planets
	assert.Equal(t, len(game.Planets), len(player.PlanetIntels))

	// should have fleets
	assert.True(t, len(player.Fleets) > 0)

	// should have grown pop
	assert.Greater(t, player.Planets[0].Population(), player.Race.Spec.StartingPlanets[0].Population)

	// should have built factories
	assert.Greater(t, player.Planets[0].Factories, game.Rules.StartingFactories)

	// should have mined minerals
	assert.Greater(t, player.Planets[0].Cargo.ToMineral().Total(), planetMinerals.Total())

}
