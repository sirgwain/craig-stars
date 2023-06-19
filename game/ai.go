package game

import (
	"math"

	"golang.org/x/exp/slices"
)

func processTurn(player *Player) {
	scout(player)
}

func scout(player *Player) {
	design := player.GetLatestDesign(ShipDesignPurposeScout)
	unknownPlanetsByNum := map[int]PlanetIntel{}
	buildablePlanets := player.GetBuildablePlanets(design.Spec.Mass)

	for _, planet := range player.PlanetIntels {
		if planet.Unexplored() {
			unknownPlanetsByNum[planet.Num] = planet
		}
	}

	scannerFleets := []*Fleet{}

	for _, fleet := range player.Fleets {
		if _, contains := fleet.Spec.Purposes[ShipDesignPurposeScout]; contains && fleet.Spec.Scanner {
			if len(fleet.Waypoints) <= 1 {
				// this fleet can be sent to scan a planet
				scannerFleets = append(scannerFleets, fleet)
			} else {
				// this fleet is already scanning a planet, remove the target from the unknown planets list
				for _, wp := range fleet.Waypoints[1:] {
					if wp.TargetNum != nil {
						delete(unknownPlanetsByNum, *wp.TargetNum)
					}
				}
			}
		}
	}

	for _, fleet := range scannerFleets {
		closestPlanet := getClosestPlanet(fleet, unknownPlanetsByNum)
		if closestPlanet != nil {
			// TODO: this isn't working once the player is loaded from the DB...
			fleet.Waypoints = append(fleet.Waypoints, NewPlanetWaypoint(closestPlanet.Position, closestPlanet.Num, closestPlanet.Name, fleet.Spec.IdealSpeed))
			fleet.Dirty = true
			delete(unknownPlanetsByNum, closestPlanet.Num)
		}
	}

	for _, planet := range buildablePlanets {
		existingQueueItemIndex := slices.IndexFunc(planet.ProductionQueue, func(item ProductionQueueItem) bool { return item.DesignName == design.Name })
		if existingQueueItemIndex == -1 {
			// put a new scout at the front of the queue
			// TODO: this isn't working once the player is loaded from the DB...
			planet.ProductionQueue = append([]ProductionQueueItem{{Type: QueueItemTypeShipToken, Quantity: 1, DesignName: design.Name}}, planet.ProductionQueue...)
			planet.Dirty = true
		}
	}

}

// get the closest planet to this fleet from a list of unknown planets
func getClosestPlanet(fleet *Fleet, unknownPlanetsByNum map[int]PlanetIntel) *PlanetIntel {
	shortestDist := math.MaxFloat64
	var closest *PlanetIntel = nil

	for num := range unknownPlanetsByNum {
		intel := unknownPlanetsByNum[num]
		distSquared := fleet.Position.DistanceSquaredTo(intel.Position)
		if shortestDist > distSquared {
			shortestDist = distSquared
			closest = &intel
		}
	}

	return closest
}
