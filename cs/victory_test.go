package cs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_victory_checkForVictor(t *testing.T) {
	// create a game with 2 planets
	game := createSingleUnitGame()
	game.Planets = append(game.Planets, NewPlanet().WithNum(2))

	// we own one planet and one fleet
	player := game.Players[0]
	player.ScoreHistory = []PlayerScore{{
		Planets:      1,
		UnarmedShips: 1,
	}}

	victory := newVictoryChecker(game)
	victory.checkForVictor(player)

	// no score, no victory
	assert.False(t, player.Victor)
	assert.Equal(t, 0, player.AchievedVictoryConditions.countBits())

	// add another planet to the game. we own 50% of them now
	// and only 1 victory condition is required so we win
	// also make it so only one year must pass for us to win
	game.VictoryConditions.OwnPlanets = 50
	game.VictoryConditions.NumCriteriaRequired = 1
	game.VictoryConditions.YearsPassed = 1
	game.Year = 2401
	player.ScoreHistory[0].Planets = 1
	victory.checkForVictor(player)
	assert.True(t, player.AchievedVictoryConditions&Bitmask(VictoryConditionOwnPlanets) > 0)
	assert.True(t, player.Victor)

}
