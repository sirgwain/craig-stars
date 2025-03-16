package ai

import (
	"math"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
)

// search through our intel and build lists of important
// things like which of our planets are being threatened by fleets
func (ai *aiPlayer) gatherIntel() {

	for _, fleet := range ai.FleetIntels {
		// skip idle fleets
		if fleet.WarpSpeed == 0 {
			continue
		}
		// skip friendly fleets
		if !ai.IsEnemy(fleet.PlayerNum) {
			continue
		}
		// skip transports and colonizers
		if !ai.hasAttackShips([]*cs.FleetIntel{&fleet}) {
			continue
		}

		// check if these hostile fleets are headed towards one of our planets
		targets := ai.findPlanetTargets(fleet.Position, fleet.Heading, ai.Planets)
		for _, target := range targets {
			log.Debug().
				Int64("GameID", ai.GameID).
				Int("PlayerNum", ai.Num).
				Msgf("Planet %s is being targetted by player %d fleet %s", target.Name, fleet.PlayerNum, fleet.Name)

			ai.targetedPlanets[target.Num] = append(ai.targetedPlanets[target.Num], &fleet)
		}

	}
}

// find any planets that are possible targets of a fleet
func (ai *aiPlayer) findPlanetTargets(position, heading cs.Vector, planets []*cs.Planet) []*cs.Planet {

	targets := []*cs.Planet{}
	// y = mx + b
	// slope is the heading
	m := heading.Y / heading.X
	b := position.Y - (m * position.X)

	for _, planet := range planets {
		// if the equation is true, and the planet is further along the line, we have a hit
		if int(math.Round(planet.Position.Y)) == int(math.Round(m*planet.Position.X+b)) &&
			((heading.X > 0 && planet.Position.X > position.X) || (heading.X < 0 && planet.Position.X < position.X)) &&
			((heading.Y > 0 && planet.Position.Y > position.Y) || (heading.Y < 0 && planet.Position.Y < position.Y)) {
			targets = append(targets, planet)
		}
	}

	return targets
}
