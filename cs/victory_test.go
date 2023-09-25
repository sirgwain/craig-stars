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

func Test_victory_checkForVictorExceedsSecondPlaceScore(t *testing.T) {
	// create a game with 2 planets
	game := createSingleUnitGame()
	player1 := game.Players[0]

	// create a new player with a lower score
	player2 := NewPlayer(2, NewRace().WithSpec(&rules)).WithNum(2).withSpec(&rules)
	game.Players = append(game.Players, player2)
	player3 := NewPlayer(3, NewRace().WithSpec(&rules)).WithNum(3).withSpec(&rules)
	game.Players = append(game.Players, player3)

	player1.Relations = []PlayerRelationship{{Relation: PlayerRelationFriend}, {Relation: PlayerRelationNeutral}, {Relation: PlayerRelationFriend}}
	player2.Relations = []PlayerRelationship{{Relation: PlayerRelationNeutral}, {Relation: PlayerRelationFriend}, {Relation: PlayerRelationFriend}}
	player3.Relations = []PlayerRelationship{{Relation: PlayerRelationNeutral}, {Relation: PlayerRelationFriend}, {Relation: PlayerRelationFriend}}
	player1.PlayerIntels.PlayerIntels = player1.defaultPlayerIntels([]*Player{player1, player2, player3})
	player2.PlayerIntels.PlayerIntels = player2.defaultPlayerIntels([]*Player{player1, player2, player3})
	player3.PlayerIntels.PlayerIntels = player3.defaultPlayerIntels([]*Player{player1, player2, player3})

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
