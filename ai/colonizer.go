package ai

import (
	"math"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
	"golang.org/x/exp/maps"
)

type colonizablePlanet struct {
	planet cs.PlanetIntel
}

type colonizer struct {
	*cs.Fleet
	target cs.PlanetIntel
}

type colonizerBuilderPlanet struct {
	*cs.Planet
	availableColonists int
}

// find all colonizable planets and send colony ships to them
func (ai *aiPlayer) colonize() error {
	ai.colonizablePlanets = ai.getColonizablePlanets()

	log.Debug().
		Int64("GameID", ai.GameID).
		Int("PlayerNum", ai.Num).
		Msgf("%d colonizable planets", len(ai.colonizablePlanets))

	// make sure our in-flight colonizers are still targeting uninhabited planets
	ai.rerouteInFlightColonizers(maps.Values(ai.colonizersByPlanet))

	// route any targeted colonizers we just build
	for _, planet := range ai.planetsWithDocks {
		for t, target := range ai.targets {
			if ai.isTargetWaitingForProduction(t) {
				log.Debug().
					Int64("GameID", ai.GameID).
					Int("PlayerNum", ai.Num).
					Msgf("colonizer targeting %s waiting for %d more ships to build", t.TargetName, len(target.production))
				continue
			}

			// assemble the fleet for this target
			fleet, err := ai.assembleFleetForTarget(t, planet.Num, cs.FleetPurposeColonizer)
			if err != nil {
				return err
			}

			if fleet != nil {
				planet := ai.getPlanetIntel(t.TargetNum)
				ai.orderColonizer(fleet, &planet)
			}
		}
	}

	// find any colonizer fleets that were built previously that need to be ordered
	colonizerFleets, err := ai.getIdleColonizers()
	if err != nil {
		return err
	}

	for _, fleet := range colonizerFleets {
		bestPlanet := ai.getBestPlanetToColonize(fleet.Position, fleet.Spec.Engine.IdealSpeed, ai.colonizablePlanets)
		ai.orderColonizer(fleet.Fleet, bestPlanet)
	}

	colonizersNeeded := len(ai.colonizablePlanets)
	if colonizersNeeded <= 0 {
		// all done, no need to build more
		return nil
	}

	// get the warp speed and cargo capacity of the colonizer fleets
	colonizerFleetMakeup := ai.fleetsByPurpose[cs.FleetPurposeColonizer]
	warpSpeed := 0
	cargoCapacity := 0
	estimatedRangeFull := 0
	for _, ship := range colonizerFleetMakeup.ships {
		if ship.design != nil {
			warpSpeed = ship.design.Spec.Engine.IdealSpeed
			cargoCapacity += ship.design.Spec.CargoCapacity
			estimatedRangeFull = ship.design.Spec.EstimatedRangeFull
		}
	}

	// no planets that can build, we're all done
	for _, planet := range ai.getColonizerBuilderPlanets() {
		numBuildableColonizers := planet.availableColonists / cargoCapacity
		for i := 0; i < numBuildableColonizers; i++ {
			targetPlanet := ai.getBestPlanetToColonize(planet.Position, warpSpeed, ai.colonizablePlanets)
			if targetPlanet == nil {
				break
			}
			target := targetPlanet.Target()
			fleetMakeup := ai.fleetsByPurpose[cs.FleetPurposeColonizer]
			distance := planet.Position.DistanceTo(targetPlanet.Position)
			if distance > float64(estimatedRangeFull) {
				// add a fuel freighter
				fleetMakeup.ships = append(fleetMakeup.ships, fleetShip{purpose: cs.ShipDesignPurposeFuelFreighter, quantity: 1, design: ai.fuelTransportDesign})
			}

			log.Debug().
				Int64("GameID", ai.GameID).
				Int("PlayerNum", ai.Num).
				Msgf("%s requesting build of colonizer fleet with %d ships", planet.Name, len(fleetMakeup.ships))

			ai.fleetBuildRequests = append(ai.fleetBuildRequests, fleetBuildRequest{buildPlanet: planet.Planet, fleet: fleetMakeup, target: &target})
			delete(ai.colonizablePlanets, targetPlanet.Num)
		}
	}
	return nil
}

func (ai *aiPlayer) getColonizablePlanets() map[int]colonizablePlanet {
	colonizablePlanets := make(map[int]colonizablePlanet)

	// find all the unowned habitable planets
	for _, planet := range ai.Player.PlanetIntels {
		if planet.Explored() && !planet.Owned() {
			if planet.Spec.Habitability > 0 || planet.Spec.TerraformedHabitability > 0 {
				if _, ok := ai.colonizersByPlanet[planet.Num]; !ok {
					colonizablePlanets[planet.Num] = colonizablePlanet{planet: planet}
				}
			}
		}
	}

	return colonizablePlanets
}

// get the planet with the best distance to hab ratio
func (ai *aiPlayer) getBestPlanetToColonize(position cs.Vector, warpSpeed int, colonizablePlanets map[int]colonizablePlanet) *cs.PlanetIntel {
	var best *cs.PlanetIntel = nil

	// lowest weight wins
	bestWeight := math.MaxFloat64
	yearlyTravelDistance := float64(warpSpeed * warpSpeed)

	for num := range colonizablePlanets {
		intel := colonizablePlanets[num].planet
		habValue := intel.Spec.Habitability
		terraformHabValue := intel.Spec.TerraformedHabitability
		dist := intel.Position.DistanceTo(position)
		yearsToTravel := math.Ceil(dist / yearlyTravelDistance)

		// weight is based on distance (in years) and hab value
		// as distance goes up, weight goes up. As habValue goes up, weight goes down
		// low weight wins
		// distance is weighted heavier than habValue. In the above calc, a 100% planet 10y away is
		// ranked equally with a 25% planet 5y away
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

// get all planets that
func (ai *aiPlayer) getColonizerBuilderPlanets() []colonizerBuilderPlanet {
	planets := make([]colonizerBuilderPlanet, 0, len(ai.planetsWithDocks))
	for _, planet := range ai.planetsWithDocks {
		if planet.Spec.PopulationDensity >= ai.config.colonistTransportDensity {

			builder := colonizerBuilderPlanet{
				Planet:             planet,
				availableColonists: int((planet.Spec.PopulationDensity - ai.config.colonistTransportDensity) * float64(planet.Spec.MaxPopulation/100)),
			}

			planets = append(planets, builder)
		}
	}
	return planets
}

func (ai *aiPlayer) getIdleColonizers() ([]colonizer, error) {
	// find all idle ships that can be part of our colonizer fleets
	// and merge them into a single fleet with purpose
	fleets, err := ai.assembleFromIdleFleets(ai.fleetsByPurpose[cs.FleetPurposeColonizer])
	if err != nil {
		return nil, err
	}
	colonizers := make([]colonizer, 0, len(fleets))

	for _, fleet := range fleets {
		planet := ai.getPlanet(fleet.OrbitingPlanetNum)
		if planet != nil && planet.OwnedBy(ai.Player.Num) && planet.Spec.PopulationDensity > ai.config.colonizerPopulationDensity {
			// this fleet can be sent to colonize a planet
			colonizers = append(colonizers, colonizer{Fleet: fleet})
		} else if fleet.Cargo.Colonists > 0 {
			// this fleet already has colonists but the planet its over is probably already
			// colonized, so send it to a new planet
			colonizers = append(colonizers, colonizer{Fleet: fleet})
		}
	}

	return colonizers, nil
}

// make sure our colonizers are still pointing at valid planets
func (ai *aiPlayer) rerouteInFlightColonizers(colonizers []colonizer) {
	for _, colonizer := range colonizers {
		target := colonizer.target
		if target.Owned() {
			// our target is owned by someone else, see if they are an enemy see if we can invade them
			if ai.IsEnemy(target.PlayerNum) && !target.Spec.HasStarbase && target.Spec.Population < int(float64(colonizer.Cargo.Colonists*100)/ai.config.invasionFactor) {
				log.Debug().
					Int64("GameID", ai.GameID).
					Int("PlayerNum", ai.Num).
					Int("Invaders", colonizer.Cargo.Colonists*100).
					Int("Defenders", target.Spec.Population).
					Bool("HasStarbase", target.Spec.HasStarbase).
					Msgf("Colonizer %s switched to invasion of %s", colonizer.Name, target.Name)

				colonizer.Purpose = cs.FleetPurposeInvader
				warpSpeed := colonizer.Waypoints[1].WarpSpeed
				colonizer.Waypoints[1] = cs.NewPlanetWaypoint(target.Position, target.Num, target.Name, warpSpeed).
					WithTask(cs.WaypointTaskTransport).
					WithTransportTasks(cs.WaypointTransportTasks{Colonists: cs.WaypointTransportTask{Action: cs.TransportActionUnloadAll}})
			} else {
				// this planet is owned by someone else and we don't want to invade
				// remove the target return this colonizer to the available queue
				colonizer.Waypoints = colonizer.Waypoints[:1]
				colonizer.RemoveTag(cs.TagTarget)
				ai.updateFleetOrders(colonizer.Fleet)
				log.Debug().
					Int64("GameID", ai.GameID).
					Int("PlayerNum", ai.Num).
					Msgf("Fleet %s was targeting %s for colonizing, but it is owned by player %d", colonizer.Name, target.Name, target.PlayerNum)

				// attempt to reroute it to a new planet
				bestPlanet := ai.getBestPlanetToColonize(colonizer.Position, colonizer.Spec.Engine.IdealSpeed, ai.colonizablePlanets)
				ai.orderColonizer(colonizer.Fleet, bestPlanet)
			}
			continue
		}

	}
}

// order a colonizer fleet to colonize the best planet
func (ai *aiPlayer) orderColonizer(fleet *cs.Fleet, bestPlanet *cs.PlanetIntel) bool {
	if bestPlanet == nil {
		return false
	}
	// if our colonizer is out in space or already has cargo, don't try and load more
	if fleet.OrbitingPlanetNum != cs.None && fleet.Cargo.Total() == 0 {

		// make sure we aren't orbiting another player's planet, somehow
		planet := ai.getPlanetIntel(fleet.OrbitingPlanetNum)
		if planet.PlayerNum != ai.Num {
			return false
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

			return false
		}
		if err := ai.client.TransferPlanetCargo(&ai.game.Rules, ai.Player, fleet, orbiting, cs.CargoTransferRequest{Cargo: cs.Cargo{Colonists: colonistsToLoad}}); err != nil {
			// something went wrong, skipi this planet
			log.Error().Err(err).Msg("transferring colonists from planet, skipping")
			return false
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
	ai.updateFleetOrders(fleet)
	delete(ai.colonizablePlanets, bestPlanet.Num)

	log.Debug().
		Int64("GameID", ai.GameID).
		Int("PlayerNum", ai.Num).
		Int("WarpSpeed", warpSpeed).
		Msgf("Fleet %s targeting %s for colonizing", fleet.Name, bestPlanet.Name)

	// we did it!
	return true
}
