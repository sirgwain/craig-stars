package ai

import (
	"github.com/rs/zerolog/log"
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
		if planet.Spec.PopulationDensity >= ai.config.colonizerPopulationDensity {
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

	for _, fleet := range fleetMakeup.getFleetsMatchingMakeup(ai, ai.Fleets) {
		if _, contains := fleet.Spec.Purposes[cs.ShipDesignPurposeColonistFreighter]; contains {
			if len(fleet.Waypoints) <= 1 {
				// this fleet is over a feeder, great!
				if _, found := feedersByNum[fleet.OrbitingPlanetNum]; found {
					fleets = append(fleets, fleet)
				} else {
					// find the nearest feeder and head that way
					closestFeeder := ai.getClosestPlanet(fleet, feedersByNum)
					if closestFeeder != nil {
						warpSpeed := ai.getWarpSpeed(fleet, closestFeeder.Position)
						fleet.Waypoints = append(fleet.Waypoints, cs.NewPlanetWaypoint(closestFeeder.Position, closestFeeder.Num, closestFeeder.Name, warpSpeed))
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

	for _, fleet := range fleets {
		closest := ai.getClosestPlanet(fleet, needersByNum)
		if closest != nil {
			if err := ai.loadColonistsAndTarget(fleet, closest); err != nil {
				// something went wrong, skip this fleet
				log.Error().Err(err).Msg("transferring colonists from planet")
				continue
			}
			delete(needersByNum, closest.Num)
		}
		if len(needersByNum) == 0 {
			break
		}
	}

	// for each planet remaining, we need a transport
	if len(needersByNum) > 0 {
		ai.addFleetBuildRequest(cs.FleetPurposeColonistFreighter, len(needersByNum))
	}

	return nil
}

func (ai *aiPlayer) loadColonistsAndTarget(fleet *cs.Fleet, planet *cs.Planet) error {

	// if our transport is out in space or already has cargo, don't try and load more
	if fleet.OrbitingPlanetNum != cs.None || fleet.Cargo.Total() == 0 {
		// make sure the planet we're orbiting still has enough pop to feed
		orbiting := ai.getPlanet(fleet.OrbitingPlanetNum)
		if orbiting != nil && orbiting.Spec.PopulationDensity >= ai.config.colonizerPopulationDensity {
			// we are over our world, load colonists
			if err := ai.client.TransferPlanetCargo(&ai.game.Rules, ai.Player, fleet, orbiting, cs.Cargo{Colonists: fleet.Spec.CargoCapacity}); err != nil {
				return err
			}
		}
	}

	// unload on this planet
	warpSpeed := ai.getWarpSpeed(fleet, planet.Position)
	fleet.Waypoints = append(fleet.Waypoints, cs.NewPlanetWaypoint(planet.Position, planet.Num, planet.Name, warpSpeed).
		WithTask(cs.WaypointTaskTransport).
		WithTransportTasks(cs.WaypointTransportTasks{Colonists: cs.WaypointTransportTask{Action: cs.TransportActionUnloadAll}}))
	ai.client.UpdateFleetOrders(ai.Player, fleet, fleet.FleetOrders)

	return nil
}
