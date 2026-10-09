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
	length := heading.Length()
	if length == 0 {
		return targets
	}

	for _, planet := range planets {
		offset := planet.Position.Subtract(position).ToFloat64()
		if offset.Dot(heading) <= 0 {
			continue
		}
		// Allow half a map unit for rounding, measured perpendicular to the
		// heading so the tolerance is the same for every direction.
		distance := math.Abs(offset.X*heading.Y-offset.Y*heading.X) / length
		if distance <= 0.5 {
			targets = append(targets, planet)
		}
	}

	return targets
}
