package ai

import (
	"strconv"
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// an idle, finished warship of a war fleet, waiting over one of the test AI's planets
func newTestWarship(t *testing.T, ai *aiPlayer, num int, planet *cs.Planet, bomber bool, quantity int) *cs.Fleet {
	t.Helper()
	fleet := newTestTransport(t, ai, num, planet)
	fleet.Tokens[0].Quantity = quantity
	fleet.Spec.TotalShips = quantity
	fleet.Spec.Bomber = bomber
	fleet.SetTag(cs.TagPurpose, string(cs.FleetPurposeBomber))
	fleet.SetTag(tagIdleSince, strconv.Itoa(ai.game.Year-musterIdleYears))
	return fleet
}

func TestScatteredWarshipsGatherWhereMostAlreadyAre(t *testing.T) {
	ai := newProductionTestAI()
	placeTestPlanets(ai)
	ai.game.Year = 2450
	ai.fleetsByPurpose[cs.FleetPurposeBomber] = fleet{purpose: cs.FleetPurposeBomber, ships: []fleetShip{{purpose: cs.ShipDesignPurposeBomber, quantity: 20}}}
	bombers := newTestWarship(t, ai, 1, ai.Planets[0], true, 6)
	escorts := newTestWarship(t, ai, 2, ai.Planets[1], false, 3)
	ai.buildMaps()

	require.NoError(t, ai.musterWarships())

	assert.Len(t, bombers.Waypoints, 1, "the larger group stays where it is")
	require.Len(t, escorts.Waypoints, 2)
	assert.Equal(t, ai.Planets[0].Num, escorts.Waypoints[1].TargetNum)
}

func TestWarFleetWithLittleFuelGoesHome(t *testing.T) {
	ai := newProductionTestAI()
	placeTestPlanets(ai)
	fleet := newTestWarship(t, ai, 1, ai.Planets[1], true, 5)
	fleet.Tokens = append(fleet.Tokens, fleet.Tokens[0]) // an assembled fleet
	ai.Planets[1].Spec.HasStarbase = false
	fleet.Fuel = 0
	ai.buildMaps()

	require.NoError(t, ai.bomb())

	require.Len(t, fleet.Waypoints, 2)
	assert.Equal(t, ai.Planets[0].Num, fleet.Waypoints[1].TargetNum, "the nearest of our worlds with a starbase")
}

func TestNormalOnlyAttacksAnEvenFight(t *testing.T) {
	ai := newProductionTestAI()
	ai.AIDifficulty = cs.AIDifficultyNormal
	warship := cs.NewShipDesign(1, 60)
	warship.Spec.PowerRating = 100
	ai.Designs = append(ai.Designs, warship)
	fleet := &cs.Fleet{Tokens: []cs.ShipToken{{DesignNum: 60, Quantity: 3}}}

	base := cs.NewShipDesign(2, 1)
	base.Spec.PowerRating = 250
	ai.ShipDesignIntels = append(ai.ShipDesignIntels, base)
	target := cs.NewPlanet()
	target.PlayerNum = 2
	target.Spec.HasStarbase = true
	target.Spec.StarbaseDesignNum = 1

	assert.True(t, ai.canAttackPlanet(fleet, target), "300 against 250")
	base.Spec.PowerRating = 350
	assert.False(t, ai.canAttackPlanet(fleet, target), "300 against 350")

	ai.AIDifficulty = cs.AIDifficultyEasy
	assert.True(t, ai.canAttackPlanet(fleet, target), "Easy doesn't check")
	ai.AIDifficulty = cs.AIDifficultyHard
	base.Spec.PowerRating = 250
	assert.False(t, ai.canAttackPlanet(fleet, target), "Hard wants an edge of a quarter")
}

func TestDefensesRaiseTheInvadersNeeded(t *testing.T) {
	ai := newProductionTestAI()
	ai.Relations = []cs.PlayerRelationship{{}, {Relation: cs.PlayerRelationEnemy}}
	target := cs.NewPlanet()
	target.PlayerNum = 2
	target.ReportAge = 0
	target.Cargo.Colonists = 100 // 10,000 defenders

	assert.True(t, ai.canInvade(target, 200), "twice the defenders")
	target.Spec.DefenseCoverage = .8
	assert.False(t, ai.canInvade(target, 200), "most of them would die on the way down")
}
