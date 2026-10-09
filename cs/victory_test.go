//go:build !wasi && !wasm

package cs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_victory_checkForVictor(t *testing.T) {
	// create a game with 2 planets
	s := SingleUnitScenario()
	s.Players = append(s.Players, ScenarioPlayer{})
	s.Planets = append(s.Planets, ScenarioPlanet{Name: "Planet 2"})
	u := newTestUniverse(t, s)
	game := u.Game
	game.Year = game.Rules.StartingYear + game.Rules.ShowPublicScoresAfterYears
	game.VictoryConditions.Conditions = Bitmask(VictoryConditionOwnPlanets)
	game.VictoryConditions.OwnPlanets = 100

	// we own one planet and one fleet
	player := u.Player(1)
	player.ScoreHistory = []PlayerScore{{
		Planets:      1,
		UnarmedShips: 1,
	}}

	victory := newVictoryChecker(game)
	victory.checkForVictor(player)

	// Not enough planets for victory.
	assert.False(t, player.Victor)
	assert.Equal(t, 0, (player.AchievedVictoryConditions & game.VictoryConditions.Conditions).countBits())

	// add another planet to the game. we own 50% of them now
	// and only 1 victory condition is required so we win
	// allow victory as soon as condition evaluation starts
	game.VictoryConditions.OwnPlanets = 50
	game.VictoryConditions.NumCriteriaRequired = 1
	game.VictoryConditions.YearsPassed = 1
	game.Year = game.Rules.StartingYear + game.Rules.ShowPublicScoresAfterYears
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
		Score: 100, Planets: 1,
	}}
	player2.ScoreHistory = []PlayerScore{{
		Score: 200, Planets: 1,
	}}
	player3.ScoreHistory = []PlayerScore{{
		Score: 50, Planets: 1,
	}}

	// add another planet to the game. we own 50% of them now
	// and only 1 victory condition is required so we win
	// allow victory as soon as condition evaluation starts
	game.VictoryConditions.Conditions = Bitmask(VictoryConditionExceedsSecondPlaceScore)
	game.VictoryConditions.ExceedsSecondPlaceScore = 100
	game.VictoryConditions.NumCriteriaRequired = 1
	game.VictoryConditions.YearsPassed = 1
	game.Year = game.Rules.StartingYear + game.Rules.ShowPublicScoresAfterYears
	victory := newVictoryChecker(game)
	victory.checkForVictor(player1)
	victory.checkForVictor(player2)
	victory.checkForVictor(player3)
	assert.Equal(t, 0, (player1.AchievedVictoryConditions & game.VictoryConditions.Conditions).countBits())
	assert.False(t, player1.Victor)
	assert.True(t, player2.AchievedVictoryConditions&Bitmask(VictoryConditionExceedsSecondPlaceScore) > 0)
	assert.True(t, player2.Victor)
	assert.Equal(t, 0, (player3.AchievedVictoryConditions & game.VictoryConditions.Conditions).countBits())
	assert.False(t, player3.Victor)

}

func Test_victory_checkForVictorNoConditions(t *testing.T) {
	s := SingleUnitScenario()
	s.Players = append(s.Players, ScenarioPlayer{})
	u := newTestUniverse(t, s)
	game := u.Game
	game.Year = game.Rules.StartingYear + game.Rules.ShowPublicScoresAfterYears

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

func Test_victory_checkForVictorStartYear(t *testing.T) {
	const allConditions = Bitmask(VictoryConditionOwnPlanets | VictoryConditionAttainTechLevels | VictoryConditionExceedsScore | VictoryConditionExceedsSecondPlaceScore | VictoryConditionProductionCapacity | VictoryConditionOwnCapitalShips | VictoryConditionHighestScoreAfterYears)
	for _, tt := range []struct {
		name                    string
		startingYear, waitYears int
		yearsPassed             int
		want                    bool
	}{
		{name: "before default threshold", startingYear: 2400, waitYears: 20, yearsPassed: 19},
		{name: "at default threshold", startingYear: 2400, waitYears: 20, yearsPassed: 20, want: true},
		{name: "before custom threshold", startingYear: 2600, waitYears: 7, yearsPassed: 6},
		{name: "at custom threshold", startingYear: 2600, waitYears: 7, yearsPassed: 7, want: true},
		{name: "no waiting period", startingYear: 2400, waitYears: 0, yearsPassed: 0, want: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			u := newTestUniverse(t, TestScenario{Players: []ScenarioPlayer{{}, {}}})
			u.Game.Rules.StartingYear = tt.startingYear
			u.Game.Rules.ShowPublicScoresAfterYears = tt.waitYears
			u.Game.Year = tt.startingYear + tt.yearsPassed
			u.Game.VictoryConditions = VictoryConditions{
				Conditions: allConditions, NumCriteriaRequired: 1,
				ExceedsScore: 100, ExceedsSecondPlaceScore: 100,
			}
			player := u.Player(1)
			previous := PlayerScore{Score: 90, AchievedVictoryConditions: allConditions}
			current := PlayerScore{Score: 100, Planets: 1, AchievedVictoryConditions: allConditions}
			player.ScoreHistory = []PlayerScore{previous, current}
			player.AchievedVictoryConditions = allConditions
			u.Player(2).ScoreHistory = []PlayerScore{{Score: 40, Planets: 1}}

			checker := newVictoryChecker(u.Game)
			assert.NoError(t, checker.checkForVictor(player))
			wantConditions := Bitmask(VictoryConditionNone)
			if tt.want {
				wantConditions = allConditions
			}
			assert.Equal(t, wantConditions, player.AchievedVictoryConditions)
			current.AchievedVictoryConditions = wantConditions
			assert.Equal(t, current, player.GetScore())
			assert.Equal(t, previous, player.ScoreHistory[0])
			assert.Equal(t, tt.want, player.Victor)
		})
	}
}

func Test_victory_conditions(t *testing.T) {
	tests := []struct {
		name          string
		condition     VictoryCondition
		settings      VictoryConditions
		score, other  PlayerScore
		planets, year int
		want          bool
	}{
		{name: "exact score threshold", condition: VictoryConditionExceedsScore, settings: VictoryConditions{ExceedsScore: 100}, score: PlayerScore{Score: 100}, want: true},
		{name: "below score threshold", condition: VictoryConditionExceedsScore, settings: VictoryConditions{ExceedsScore: 100}, score: PlayerScore{Score: 99}, want: false},
		{name: "round planet count down", condition: VictoryConditionOwnPlanets, settings: VictoryConditions{OwnPlanets: 40}, score: PlayerScore{Planets: 4}, planets: 11, want: true},
		{name: "round planet count up", condition: VictoryConditionOwnPlanets, settings: VictoryConditions{OwnPlanets: 50}, score: PlayerScore{Planets: 4}, planets: 9, want: false},
		{name: "second place fractional threshold", condition: VictoryConditionExceedsSecondPlaceScore, settings: VictoryConditions{ExceedsSecondPlaceScore: 50}, score: PlayerScore{Score: 151}, other: PlayerScore{Score: 101}, want: true},
		{name: "below second place threshold", condition: VictoryConditionExceedsSecondPlaceScore, settings: VictoryConditions{ExceedsSecondPlaceScore: 50}, score: PlayerScore{Score: 150}, other: PlayerScore{Score: 101}, want: false},
		{name: "second place zero", condition: VictoryConditionExceedsSecondPlaceScore, settings: VictoryConditions{ExceedsSecondPlaceScore: 100}, score: PlayerScore{Score: 100}, other: PlayerScore{}, want: true},
		{name: "not the leader", condition: VictoryConditionExceedsSecondPlaceScore, settings: VictoryConditions{ExceedsSecondPlaceScore: 0}, score: PlayerScore{Score: 100}, other: PlayerScore{Score: 200}, want: false},
		{name: "highest score before year", condition: VictoryConditionHighestScoreAfterYears, settings: VictoryConditions{HighestScoreAfterYears: 80}, score: PlayerScore{Score: 200}, other: PlayerScore{Score: 100}, year: 2479, want: false},
		{name: "highest score at year", condition: VictoryConditionHighestScoreAfterYears, settings: VictoryConditions{HighestScoreAfterYears: 80}, score: PlayerScore{Score: 200}, other: PlayerScore{Score: 100}, year: 2480, want: true},
		{name: "highest score tied", condition: VictoryConditionHighestScoreAfterYears, settings: VictoryConditions{HighestScoreAfterYears: 80}, score: PlayerScore{Score: 200}, other: PlayerScore{Score: 200}, year: 2480, want: false},
		{name: "highest score lower", condition: VictoryConditionHighestScoreAfterYears, settings: VictoryConditions{HighestScoreAfterYears: 80}, score: PlayerScore{Score: 100}, other: PlayerScore{Score: 200}, year: 2480, want: false},
		{name: "production uses score", condition: VictoryConditionProductionCapacity, settings: VictoryConditions{ProductionCapacity: 1}, score: PlayerScore{Resources: 1000}, want: true},
		{name: "production below threshold", condition: VictoryConditionProductionCapacity, settings: VictoryConditions{ProductionCapacity: 1}, score: PlayerScore{Resources: 999}, want: false},
		{name: "capital count equality", condition: VictoryConditionOwnCapitalShips, settings: VictoryConditions{OwnCapitalShips: 10}, score: PlayerScore{CapitalShips: 10}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := newTestUniverse(t, TestScenario{Players: []ScenarioPlayer{{}, {}}})
			u.Game.VictoryConditions = tt.settings
			u.Game.VictoryConditions.Conditions = Bitmask(tt.condition)
			u.Game.Year = tt.year
			if tt.year == 0 {
				u.Game.Year = u.Game.Rules.StartingYear + u.Game.Rules.ShowPublicScoresAfterYears
			}
			for range tt.planets {
				u.Game.Planets = append(u.Game.Planets, NewPlanet())
			}
			u.Player(1).ScoreHistory = []PlayerScore{tt.score}
			u.Player(2).ScoreHistory = []PlayerScore{tt.other}
			v := newVictoryChecker(u.Game)
			assert.NoError(t, v.checkForVictor(u.Player(1)))
			assert.Equal(t, tt.want, u.Player(1).AchievedVictoryConditions&Bitmask(tt.condition) != 0)
			assert.Equal(t, u.Player(1).AchievedVictoryConditions, u.Player(1).GetScore().AchievedVictoryConditions)
		})
	}
}
