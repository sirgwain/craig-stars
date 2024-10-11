package ai

import (
	"math"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
)

// find all colonizable planets and send colony ships to them
func (ai *aiPlayer) colonize() error {
	colonizablePlanets := map[int]cs.PlanetIntel{}

	// find all the unexplored planets
	for _, planet := range ai.Player.PlanetIntels {
		if planet.Explored() && !planet.Owned() {
			if planet.Spec.Habitability > 0 || planet.Spec.TerraformedHabitability > 0 {
				colonizablePlanets[planet.Num] = planet
			}
		}
	}

	fleetMakeup := ai.fleetsByPurpose[cs.FleetPurposeColonizer]

	// find all idle ships that can be part of our colonizer fleets
	// and merge them into a single fleet with purpose
	colonizerFleets, err := ai.assembleFromIdleFleets(fleetMakeup)
	if err != nil {
		return err
	}

	// log.Debug().
	// 	Int64("GameID", ai.GameID).
	// 	Int("PlayerNum", ai.Num).
	// 	Msgf("%d colonizer fleets assembled from idle fleets", len(colonizerFleets))

	// go through fleets in space that may have been misassigned
	for _, fleet := range fleetMakeup.getFleetsMatchingMakeup(ai, ai.Fleets) {
		if fleet.GetTag(cs.TagPurpose) == string(cs.FleetPurposeColonizer) && fleet.Spec.Colonizer {
			if len(fleet.Waypoints) <= 1 {
				planet := ai.getPlanet(fleet.OrbitingPlanetNum)
				if planet != nil && planet.OwnedBy(ai.Player.Num) && planet.Spec.PopulationDensity > ai.config.colonizerPopulationDensity {
					// this fleet can be sent to colonize a planet
					colonizerFleets = append(colonizerFleets, fleet)
				} else if fleet.Cargo.Colonists > 0 {
					// this fleet already has colonists but the planet its over is probably already
					// colonized, so send it to a new planet
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
						// our target is owned by someone else, see if they are an enemy and if we can invade them
						if ai.IsEnemy(target.PlayerNum) && !target.Spec.HasStarbase && target.Spec.Population < int(float64(fleet.Cargo.Colonists*100)/ai.config.invasionFactor) {
							log.Debug().
								Int64("GameID", ai.GameID).
								Int("PlayerNum", ai.Num).
								Int("Invaders", fleet.Cargo.Colonists*100).
								Int("Defenders", target.Spec.Population).
								Bool("HasStarbase", target.Spec.HasStarbase).
								Msgf("Colonizer %s switched to invasion of %s", fleet.Name, target.Name)

							fleet.Purpose = cs.FleetPurposeInvader
							warpSpeed := fleet.Waypoints[1].WarpSpeed
							fleet.Waypoints[1] = cs.NewPlanetWaypoint(target.Position, target.Num, target.Name, warpSpeed).
								WithTask(cs.WaypointTaskTransport).
								WithTransportTasks(cs.WaypointTransportTasks{Colonists: cs.WaypointTransportTask{Action: cs.TransportActionUnloadAll}})
						} else {
							// this planet is owned by someone else and we don't want to invade
							// remove the target return this colonizer to the available queue
							fleet.Waypoints = fleet.Waypoints[:1]
							colonizerFleets = append(colonizerFleets, fleet)
							log.Debug().
								Int64("GameID", ai.GameID).
								Int("PlayerNum", ai.Num).
								Msgf("Fleet %s was targeting %s for colonizing, but it is owned by player %d", fleet.Name, target.Name, target.PlayerNum)

						}
						continue
					}
				}
			}
		}
	}

	// after colonizing, we may have idle fleets leftover
	idleFleets := len(colonizerFleets)
	log.Debug().
		Int64("GameID", ai.GameID).
		Int("PlayerNum", ai.Num).
		Msgf("%d colonizerFleets, %d colonizable planets", idleFleets, len(colonizablePlanets))

	for _, fleet := range colonizerFleets {
		bestPlanet := ai.getBestPlanetToColonize(fleet, colonizablePlanets)
		if bestPlanet != nil {
			// if our colonizer is out in space or already has cargo, don't try and load more
			if fleet.OrbitingPlanetNum != cs.None && fleet.Cargo.Total() == 0 {

				// make sure we aren't orbiting another player's planet, somehow
				planet := ai.getPlanetIntel(fleet.OrbitingPlanetNum)
				if planet.PlayerNum != ai.Num {
					continue
				}

				// don't load more than 100% of the planet cap
				colonistsToLoad := cs.MinInt(planet.Spec.MaxPopulation, fleet.Spec.CargoCapacity)

				// we are over our world, load colonists
				// but only if taking  these colonists doesn't reduce our pop too much
				// take into account how much we're going to grow
				orbiting := ai.getPlanet(fleet.OrbitingPlanetNum)
				growth := orbiting.Spec.GrowthAmount
				newDensity := float64(((orbiting.Cargo.Colonists-colonistsToLoad)*100)+growth) / float64(orbiting.Spec.MaxPopulation)
				if newDensity < ai.config.colonizerPopulationDensity {
					log.Debug().
						Int64("GameID", ai.GameID).
						Int("PlayerNum", ai.Num).
						Int("ColonistsAvailable", orbiting.Cargo.Colonists*100).
						Int("ColonistsNeeded", colonistsToLoad*100).
						Int("DensityAfterLoad", int(newDensity)).
						Msgf("Fleet %s cannot load colonists from %s", fleet.Name, orbiting.Name)

					continue
				}
				if err := ai.client.TransferPlanetCargo(&ai.game.Rules, ai.Player, fleet, orbiting, cs.CargoTransferRequest{Cargo: cs.Cargo{Colonists: colonistsToLoad}}, ai.Planets); err != nil {
					// something went wrong, skipi this planet
					log.Error().Err(err).Msg("transferring colonists from planet, skipping")
					continue
				}

				log.Debug().
					Int64("GameID", ai.GameID).
					Int("PlayerNum", ai.Num).
					Int("ColonistsAvailable", orbiting.Cargo.Colonists*100).
					Int("ColonistsNeeded", colonistsToLoad*100).
					Int("DensityAfterLoad", int(newDensity)).
					Msgf("Fleet %s loaded %d colonists from %s", fleet.Name, colonistsToLoad*100, orbiting.Name)

			}

			warpSpeed := ai.getWarpSpeed(fleet, bestPlanet.Position)
			fleet.Waypoints = append(fleet.Waypoints, cs.NewPlanetWaypoint(bestPlanet.Position, bestPlanet.Num, bestPlanet.Name, warpSpeed).WithTask(cs.WaypointTaskColonize))
			ai.client.UpdateFleetOrders(ai.Player, fleet, fleet.FleetOrders)
			delete(colonizablePlanets, bestPlanet.Num)
			idleFleets--

			log.Debug().
				Int64("GameID", ai.GameID).
				Int("PlayerNum", ai.Num).
				Int("WarpSpeed", warpSpeed).
				Msgf("Fleet %s targeting %s for colonizing", fleet.Name, bestPlanet.Name)

		}
	}

	// build colonizer fleets where necessary
	if len(colonizablePlanets)-idleFleets > 0 {
		ai.addFleetBuildRequest(cs.FleetPurposeColonizer, len(colonizablePlanets)-idleFleets)
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
		habValue := intel.Spec.Habitability
		terraformHabValue := intel.Spec.TerraformedHabitability
		dist := intel.Position.DistanceTo(fleet.Position)
		yearsToTravel := math.Ceil(dist / yearlyTravelDistance)

		// weight is based on distance (in years) and hab value
		// as distance goes up, weight goes up. As habValue goes up, weight goes down
		// low weight wins
		// distance is weighted heavier than habValue. In the above calc, a 100% planet 10y away is ranked equally with a
		// 25% planet 5y away
		weight := math.MaxFloat64
		if habValue > 0 {
			weight = (yearsToTravel * yearsToTravel) / float64(habValue)
		} else if terraformHabValue > 0 {
			// only colonize terraformable planets if they're really close
			weight = 2 * (yearsToTravel * yearsToTravel) / float64(terraformHabValue)
		}
		if weight < bestWeight {
			bestWeight = weight
			best = &intel
		}
	}

	// best = ai.getClosestPlanetIntel(fleet.Position, colonizablePlanets)

	return best
}
