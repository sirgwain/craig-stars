package ai

import (
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/stretchr/testify/assert"
)

func TestBomberSetsOutForTheNearestTargetWhenNoneIsInRange(t *testing.T) {
	ai := newProductionTestAI()
	placeTestPlanets(ai)
	fleet := newTestWarship(t, ai, 1, ai.Planets[0], true, 5)

	// enemy worlds far beyond the fleet's reach in the allowed years
	targets := map[int]*cs.Planet{}
	for i, x := range []int{3000, 2000} {
		target := ai.PlanetIntels[i+1]
		target.PlayerNum = 2
		target.Position = cs.Vector{X: x, Y: 100}
		targets[target.Num] = target
	}

	best := ai.getBestPlanetToBomb(fleet, targets)

	if assert.NotNil(t, best, "a fleet with nothing in range still sets out") {
		assert.Equal(t, ai.PlanetIntels[2].Num, best.Num, "the nearer world")
	}
}
