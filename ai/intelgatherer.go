package ai

import (
	"log/slog"
	"math"

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
		if !ai.hasAttackShips([]*cs.Fleet{fleet}) {
			continue
		}

		// check if these hostile fleets are headed towards one of our planets
		targets := ai.findPlanetTargets(fleet.Position, fleet.Heading, ai.Planets)
		for _, target := range targets {
			ai.log.Debug("Planet is being targetted by enemy fleet",
				slog.String("planet", target.Name),
				slog.Int("playerNum", fleet.PlayerNum),
				slog.String("fleet", fleet.Name))

			ai.targetedPlanets[target.Num] = append(ai.targetedPlanets[target.Num], fleet)
		}

	}
}

// find any planets that are possible targets of a fleet
func (ai *aiPlayer) findPlanetTargets(position cs.Vector, heading cs.VectorFloat64, planets []*cs.Planet) []*cs.Planet {

	targets := []*cs.Planet{}
	// y = mx + b
	// slope is the heading
	m := heading.Y / heading.X
	b := float64(position.Y) - (m * float64(position.X))

	for _, planet := range planets {
		// if the equation is true, and the planet is further along the line, we have a hit
		if int(math.Round(float64(planet.Position.Y))) == int(math.Round(m*float64(planet.Position.X)+b)) &&
			((heading.X > 0 && planet.Position.X > position.X) || (heading.X < 0 && planet.Position.X < position.X)) &&
			((heading.Y > 0 && planet.Position.Y > position.Y) || (heading.Y < 0 && planet.Position.Y < position.Y)) {
			targets = append(targets, planet)
		}
	}

	return targets
}
