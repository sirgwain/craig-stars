package ai

import (
	"github.com/sirgwain/craig-stars/cs"
)

// transport resources and people from planet to planet
func (ai *aiPlayer) transport() error {
	ai.transportColonists()

	return nil
}

func (ai *aiPlayer) transportColonists() error {
	feedersByNum := map[int]*cs.Planet{}
	needersByNum := map[int]*cs.Planet{}

	// transport from planets that are fully terraformed and above the growth rate
	for _, planet := range ai.Planets {
		if planet.Spec.PopulationDensity >= ai.config.colonistTransportDensity {
			if !planet.Spec.CanTerraform {
				feedersByNum[planet.Num] = planet
			}
		} else {
			needersByNum[planet.Num] = planet
		}
	}

	// nothing to transport
	if len(feedersByNum) == 0 || len(needersByNum) == 0 {
		return nil
	}

	fleetMakeup := ai.fleetsByPurpose[cs.FleetPurposeColonistFreighter]
	// find all idle ships that can be part of our transport fleets
	// and merge them into a single fleet with purpose
	fleets, err := ai.assembleFromIdleFleets(fleetMakeup)
	if err != nil {
		return err
	}

	ai.log.Debug().
		Msgf("%d colonist transport fleets assembled from idle fleets", len(fleets))

	for _, fleet := range fleetMakeup.getFleetsMatchingMakeup(ai, ai.Fleets) {
		// don't use colonizer fleets as transports
		if !fleet.Spec.Colonizer {
			if len(fleet.Waypoints) <= 1 {
				// this fleet is over a feeder, great!
				if _, found := feedersByNum[fleet.OrbitingPlanetNum]; found {
					fleets = append(fleets, fleet)

					orbiting := ai.getPlanet(fleet.OrbitingPlanetNum)
					ai.log.Debug().
						Msgf("%s will load colonists from %s for transport to a needy world", fleet.Name, orbiting.Name)

				} else {
					// find the nearest feeder and head that way
					closestFeeder := ai.getClosestPlanet(fleet, feedersByNum)
					if closestFeeder != nil {
						newWpIndex := fleet.AddWaypoint(ai.Player, cs.WaypointDest{MO: closestFeeder.MapObject}, len(fleet.Waypoints)-1, false)
						if newWpIndex == 0 {
							ai.log.Warn().
								Msgf("Fleet %s tried to target %s for feeding but did not add the waypoint", fleet.Name, closestFeeder.Name)
							continue
						}

						// TODO: only remove a feeder if we have too many targets?
						delete(feedersByNum, closestFeeder.Num)

						ai.log.Debug().
							Msgf("%s is heading to %s to load colonists for transport to a needy world", fleet.Name, closestFeeder.Name)
					}
				}
			} else {
				// this fleet is already transporting to a needer, remove the target from the needers list
				for _, wp := range fleet.Waypoints[1:] {
					if wp.TargetType == cs.MapObjectTypePlanet && wp.TargetNum != cs.None {
						delete(needersByNum, wp.TargetNum)
					}
				}
			}
		}

	}

	idleFleets := len(fleets)
	ai.log.Debug().
		Msgf("%d transport, %d needy planets", idleFleets, len(needersByNum))

	for _, fleet := range fleets {
		planet := ai.getClosestPlanet(fleet, needersByNum)
		if planet != nil {

			// if our transport is out in space or already has cargo, don't try and load more
			if fleet.OrbitingPlanetNum != cs.None || fleet.Cargo.Total() == 0 {
				// make sure the planet we're orbiting still has enough pop to feed
				orbiting := ai.getPlanet(fleet.OrbitingPlanetNum)
				if orbiting != nil {
					// don't load more than 25% of the target planet
					// it will grow slower after 25%
					colonistsToLoad := min(int(float64(planet.Spec.MaxPopulation)*ai.config.colonistTransportDensity), fleet.Spec.CargoCapacity)

					// load colonists but only if taking these colonists doesn't reduce our pop too much
					// take into account how much we're going to grow
					orbiting := ai.getPlanet(fleet.OrbitingPlanetNum)
					popNextYear := orbiting.PopNextYear()
					newDensity := float64(popNextYear-colonistsToLoad*100) / float64(orbiting.Spec.MaxPopulation)
					if newDensity < ai.config.colonistTransportDensity {
						ai.log.Debug().
							Int64("GameID", ai.GameID).
							Int("PlayerNum", ai.Num).
							Int("ColonistsAvailable", popNextYear).
							Int("ColonistsNeeded", colonistsToLoad*100).
							Int("DensityAfterLoad", int(newDensity)).
							Msgf("Fleet %s cannot load colonists from %s", fleet.Name, orbiting.Name)

						continue
					}
					if err := ai.client.TransferByHand(&ai.game.Rules, ai.Player, fleet, orbiting, cs.CargoTransferRequest{Cargo: cs.Cargo{Colonists: colonistsToLoad}}); err != nil {
						// something went wrong, skip this planet
						ai.log.Error().Err(err).Msg("transferring colonists from planet returned error, skipping")
						continue
					}
				}
			}

			// unload on this planet
			newWpIndex := fleet.AddWaypoint(ai.Player, cs.WaypointDest{MO: planet.MapObject}, len(fleet.Waypoints)-1, false)
			if newWpIndex == 0 {
				ai.log.Warn().
					Msgf("Fleet %s tried to target %s for unloading but did not add the waypoint", fleet.Name, planet.Name)
				continue
			}

			fleet.Waypoints[newWpIndex].Task = cs.WaypointTaskTransport
			fleet.Waypoints[newWpIndex].TransportTasks = cs.WaypointTransportTasks{Colonists: cs.WaypointTransportTask{Action: cs.TransportActionUnloadAll}}
			ai.client.UpdateFleetOrders(ai.Player, fleet, fleet.FleetOrders)

			idleFleets--
			delete(needersByNum, planet.Num)

			ai.log.Debug().
				Int64("GameID", ai.GameID).
				Int("PlayerNum", ai.Num).
				Msgf("fleet %s transporting %d colonists to %s", fleet.Name, fleet.Cargo.Colonists*100, planet.Name)

		}
		if len(needersByNum) == 0 {
			break
		}
	}

	// for each planet remaining, we need a transport
	if len(needersByNum)-idleFleets > 0 {
		ai.addFleetBuildRequest(cs.FleetPurposeColonistFreighter, len(needersByNum)-idleFleets)
	}

	return nil
}

// TODO: implement this
func (ai *aiPlayer) loadColonistsAndTarget(fleet *cs.Fleet, planet *cs.Planet) error {

	return nil
}
