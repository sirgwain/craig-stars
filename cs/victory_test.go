//go:build !wasi && !wasm

package cs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_victory_checkForVictor(t *testing.T) {
	// create a game with 2 planets
	s := SingleUnitScenario()
	s.Planets = append(s.Planets, ScenarioPlanet{Name: "Planet 2"})
	u := newTestUniverse(t, s)
	game := u.Game

	// we own one planet and one fleet
	player := u.Player(1)
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

func Test_victory_checkForVictorExceedsSecondPlaceScore(t *testing.T) {
	u := newTestUniverse(t, TestScenario{Players: []ScenarioPlayer{
		{
			Relations: []PlayerRelationship{{Relation: PlayerRelationFriend}, {Relation: PlayerRelationNeutral}, {Relation: PlayerRelationFriend}},
		},
		{
			Relations: []PlayerRelationship{{Relation: PlayerRelationNeutral}, {Relation: PlayerRelationFriend}, {Relation: PlayerRelationFriend}},
		},
		{
			Relations: []PlayerRelationship{{Relation: PlayerRelationNeutral}, {Relation: PlayerRelationFriend}, {Relation: PlayerRelationFriend}},
		},
	}, Planets: []ScenarioPlanet{{Name: "Planet 1", Owner: 1, Cargo: Cargo{Colonists: 2500}}}})
	game, player1, player2, player3 := u.Game, u.Player(1), u.Player(2), u.Player(3)

	// we own one planet and one fleet
	player1.ScoreHistory = []PlayerScore{{
		Score: 100,
	}}
	player2.ScoreHistory = []PlayerScore{{
		Score: 200,
	}}
	player3.ScoreHistory = []PlayerScore{{
		Score: 50,
	}}

	// add another planet to the game. we own 50% of them now
	// and only 1 victory condition is required so we win
	// also make it so only one year must pass for us to win
	game.VictoryConditions.ExceedsSecondPlaceScore = 100
	game.VictoryConditions.NumCriteriaRequired = 1
	game.VictoryConditions.YearsPassed = 1
	game.Year = 2401
	victory := newVictoryChecker(game)
	victory.checkForVictor(player1)
	victory.checkForVictor(player2)
	victory.checkForVictor(player3)
	assert.Equal(t, 0, player1.AchievedVictoryConditions.countBits())
	assert.False(t, player1.Victor)
	assert.True(t, player2.AchievedVictoryConditions&Bitmask(VictoryConditionExceedsSecondPlaceScore) > 0)
	assert.True(t, player2.Victor)
	assert.Equal(t, 0, player3.AchievedVictoryConditions.countBits())
	assert.False(t, player3.Victor)

}

func Test_victory_checkForVictorNoConditions(t *testing.T) {
	s := SingleUnitScenario()
	u := newTestUniverse(t, s)
	game := u.Game

	player := u.Player(1)
	player.ScoreHistory = []PlayerScore{{Planets: 1}}

	// a game with no victory conditions and no criteria required should never declare a victor
	game.VictoryConditions = VictoryConditions{}

	victory := newVictoryChecker(game)
	victory.checkForVictor(player)
	assert.False(t, player.Victor)
	assert.False(t, game.VictorDeclared)

	// conditions but no criteria required should also not declare a victor
	game.VictoryConditions.Conditions = Bitmask(VictoryConditionOwnPlanets)
	game.VictoryConditions.OwnPlanets = 1
	victory.checkForVictor(player)
	assert.False(t, player.Victor)
	assert.False(t, game.VictorDeclared)
}
