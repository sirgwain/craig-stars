package ai

import (
	"math"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
)

// find all colonizable planets and send colony ships to them
func (ai *aiPlayer) colonize() error {
	colonizablePlanets := map[int]cs.PlanetIntel{}
	terraformablePlanets := map[int]cs.PlanetIntel{}

	// find all the unexplored planets
	for _, planet := range ai.Player.PlanetIntels {
		if planet.Explored() && !planet.Owned() {
			if planet.Spec.Habitability > 0 {
				colonizablePlanets[planet.Num] = planet
			} else if planet.Spec.TerraformedHabitability > 0 {
				terraformablePlanets[planet.Num] = planet
			}
		}
	}

	// find all idle fleets that are colonizers
	colonizerFleets := []*cs.Fleet{}
	for _, fleet := range ai.Fleets {
		if _, contains := fleet.Spec.Purposes[cs.ShipDesignPurposeColonizer]; contains && fleet.Spec.Colonizer {
			if len(fleet.Waypoints) <= 1 {
				planet := ai.getPlanet(fleet.OrbitingPlanetNum)
				if planet != nil && planet.OwnedBy(ai.Player.Num) && planet.Spec.PopulationDensity > ai.config.colonizerPopulationDensity {
					// this fleet can be sent to colonize a planet
					colonizerFleets = append(colonizerFleets, fleet)
				}
			} else {
				// this fleet is already colonizing a planet, remove the target from the colonizable planets list
				for _, wp := range fleet.Waypoints[1:] {
					if wp.TargetNum != cs.None {
						delete(colonizablePlanets, wp.TargetNum)
					}

					target := ai.getPlanetIntel(wp.TargetNum)
					if target.Owned() {
						// this planet is owned, don't target it anymore and return this colonizer to the queue
						fleet.Waypoints = fleet.Waypoints[:1]
						colonizerFleets = append(colonizerFleets, fleet)
						continue
					}
				}
			}
		}
	}

	// TODO: we should sort by best planet and find the closest fleet to that planet to send a colonizer
	for _, fleet := range colonizerFleets {
		bestPlanet := ai.getBestPlanetToColonize(fleet, colonizablePlanets)
		if bestPlanet == nil {
			// TODO: send larger colonizers to terraformable planets, otherwise they die off
			// try terraformable planets when we run out of easy planets
			// bestPlanet = ai.getBestPlanetToColonize(fleet, terraformablePlanets)
		}
		if bestPlanet != nil {

			// if our colonizer is out in space or already has cargo, don't try and load more
			if fleet.OrbitingPlanetNum != cs.None && fleet.Cargo.Total() == 0 {

				// make sure we aren't orbiting another player's planet, somehow
				planet := ai.getPlanetIntel(fleet.OrbitingPlanetNum)
				if planet.PlayerNum != ai.Num {
					continue
				}

				// we are over our world, load colonists
				if err := ai.client.TransferPlanetCargo(&ai.game.Rules, ai.Player, fleet, ai.getPlanet(fleet.OrbitingPlanetNum), cs.Cargo{Colonists: fleet.Spec.CargoCapacity}); err != nil {
					// something went wrong, skipi this planet
					log.Error().Err(err).Msg("transferring colonists from planet, skipping")
					continue
				}

			}

			warpSpeed := ai.getWarpSpeed(fleet, bestPlanet.Position)
			fleet.Waypoints = append(fleet.Waypoints, cs.NewPlanetWaypoint(bestPlanet.Position, bestPlanet.Num, bestPlanet.Name, warpSpeed).WithTask(cs.WaypointTaskColonize))
			ai.client.UpdateFleetOrders(ai.Player, fleet, fleet.FleetOrders)
			delete(colonizablePlanets, bestPlanet.Num)
			delete(terraformablePlanets, bestPlanet.Num)
		}
	}

	// build scout ships where necessary
	if len(colonizablePlanets) > 0 || len(terraformablePlanets) > 0 {
		ai.addShipBuildRequest(cs.ShipDesignPurposeColonizer, len(colonizablePlanets)+len(terraformablePlanets))
	}
	return nil
}

// get the planet with the best distance to hab ratio
func (ai *aiPlayer) getBestPlanetToColonize(fleet *cs.Fleet, colonizablePlanets map[int]cs.PlanetIntel) *cs.PlanetIntel {
	var best *cs.PlanetIntel = nil

	// lowest weight wins
	bestWeight := math.MaxFloat64
	yearlyTravelDistance := float64(fleet.Spec.Engine.IdealSpeed * fleet.Spec.Engine.IdealSpeed)

	for num := range colonizablePlanets {
		intel := colonizablePlanets[num]
		habValue := ai.Player.Race.GetPlanetHabitability(intel.Hab)
		dist := intel.Position.DistanceTo(fleet.Position)
		yearsToTravel := math.Ceil(dist / yearlyTravelDistance)

		// weight is based on distance (in years) and hab value
		// as distance goes up, weight goes up. As habValue goes up, weight goes down
		// low weight wins
		// distance is weighted heavier than habValue. In the above calc, a 100% planet 10y away is ranked equally with a
		// 25% planet 5y away
		weight := (yearsToTravel * yearsToTravel) / float64(habValue)
		if weight < bestWeight {
			bestWeight = weight
			best = &intel
		}
	}

	// best = ai.getClosestPlanetIntel(fleet.Position, colonizablePlanets)

	return best
}
