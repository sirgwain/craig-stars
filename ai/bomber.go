package ai

import (
	"math"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
)

func (ai *aiPlayer) bomb() error {

	bombablePlanets := map[int]cs.PlanetIntel{}

	// find all the unexplored planets
	for _, planet := range ai.Player.PlanetIntels {
		if ai.IsEnemy(planet.PlayerNum) {
			bombablePlanets[planet.Num] = planet
		}
	}

	fleetMakeup := ai.fleetsByPurpose[cs.FleetPurposeBomber]

	// find all idle ships that can be part of our colonizer fleets
	// and merge them into a single fleet with purpose
	bomberFleets, err := ai.assembleFromIdleFleets(fleetMakeup)
	if err != nil {
		return err
	}

	// go through fleets in space that may have been misassigned
	for _, fleet := range fleetMakeup.getFleetsMatchingMakeup(ai, ai.Fleets) {
		if fleet.GetTag(cs.TagPurpose) == string(cs.FleetPurposeBomber) && fleet.Spec.Bomber {
			if len(fleet.Waypoints) > 1 {
				// this fleet is already targeting a planet, remove the target from the bombable planets list
				for _, wp := range fleet.Waypoints[1:] {
					if _, found := bombablePlanets[wp.TargetNum]; !found {
						// this fleet is targeting a planet that is no longer bombable, clear its waypoints so it can be used
						// to bomb something else
						target := ai.getPlanetIntel(wp.TargetNum)
						fleet.Waypoints = fleet.Waypoints[:1]
						bomberFleets = append(bomberFleets, fleet)
						log.Debug().
							Int64("GameID", ai.GameID).
							Int("PlayerNum", ai.Num).
							Msgf("Fleet %s was going to bomb %s, but it's no longer a bombable target", fleet.Name, target.Name)

					} else {
						delete(bombablePlanets, wp.TargetNum)
					}
				}
			} else {
				// if we are already orbiting a bombable planet, we're good
				if _, found := bombablePlanets[fleet.OrbitingPlanetNum]; found {
					delete(bombablePlanets, fleet.OrbitingPlanetNum)
				} else {
					// we aren't orbiting a bombable planet yet, go find one!
					bomberFleets = append(bomberFleets, fleet)
				}
			}
		}
	}

	// after colonizing, we may have idle fleets leftover
	idleFleets := len(bomberFleets)
	log.Debug().
		Int64("GameID", ai.GameID).
		Int("PlayerNum", ai.Num).
		Msgf("%d bomber, %d bombable planets", idleFleets, len(bombablePlanets))

	for _, fleet := range bomberFleets {
		bestPlanet := ai.getBestPlanetToBomb(fleet, bombablePlanets)
		if bestPlanet != nil {
			warpSpeed := ai.getMaxWarp(fleet.Position.DistanceTo(bestPlanet.Position), fleet)
			fleet.Waypoints = append(fleet.Waypoints, cs.NewPlanetWaypoint(bestPlanet.Position, bestPlanet.Num, bestPlanet.Name, warpSpeed))
			ai.client.UpdateFleetOrders(ai.Player, fleet, fleet.FleetOrders, fleet.Tags)
			delete(bombablePlanets, bestPlanet.Num)
			idleFleets--

			log.Debug().
				Int64("GameID", ai.GameID).
				Int("PlayerNum", ai.Num).
				Int("WarpSpeed", warpSpeed).
				Int("Population", bestPlanet.Spec.Population).
				Bool("HasStarbase", bestPlanet.Spec.HasStarbase).
				Msgf("Fleet %s targeting %s for bombing", fleet.Name, bestPlanet.Name)

		}
	}

	// build colonizer fleets where necessary
	if len(bombablePlanets)-idleFleets > 0 {
		ai.addFleetBuildRequest(cs.FleetPurposeBomber, len(bombablePlanets)-idleFleets)
	}
	return nil
}

// get the planet with the best distance to hab ratio
func (ai *aiPlayer) getBestPlanetToBomb(fleet *cs.Fleet, planets map[int]cs.PlanetIntel) *cs.PlanetIntel {
	var best *cs.PlanetIntel = nil

	// lowest weight wins
	bestWeight := 0.0
	yearlyTravelDistance := float64(fleet.Spec.Engine.IdealSpeed * fleet.Spec.Engine.IdealSpeed)

	for num := range planets {
		intel := planets[num]
		// go for low pop planets to clean them up quickly
		// like the snowball method of paying off debt
		pop := intel.Spec.Population

		// avoid starbases if there are better targets
		starbaseFactor := 1.0
		if intel.Spec.HasStarbase {
			starbaseFactor = 2.0
		}
		dist := intel.Position.DistanceTo(fleet.Position)
		yearsToTravel := math.Ceil(dist / yearlyTravelDistance)

		// weight is based on distance (in years) and pop
		// closer is better, lower pop is better, starbases are discouraged
		weight := (1 / float64(pop)) / (2 * yearsToTravel) / starbaseFactor

		// log.Debug().
		// 	Int64("GameID", ai.GameID).
		// 	Int("PlayerNum", ai.Num).
		// 	Float64("dist", dist).
		// 	Int("pop", pop).
		// 	Float64("yearsToTravel", yearsToTravel).
		// 	Bool("hasStarbase", intel.Spec.HasStarbase).
		// 	Float64("weight", weight*1000).
		// 	Msgf("getBestPlanetToBomb %s", intel.Name)

		if weight > bestWeight {
			bestWeight = weight
			best = &intel
		}
	}

	// best = ai.getClosestPlanetIntel(fleet.Position, colonizablePlanets)

	return best
}
