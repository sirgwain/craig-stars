package ai

import (
	"github.com/sirgwain/craig-stars/cs"
	"log/slog"
)

// lay mines
func (ai *aiPlayer) layMines() error {
	design := ai.Player.GetLatestDesign(cs.ShipDesignPurposeDamageMineLayer)
	// TODO: Add support for remote detonating minefield perimeters X years out in space if playing SD

	// no mine layers
	if design == nil {
		return nil
	}

	// lay mines to protect all our planets
	planetsToProtect := make([]*cs.Planet, len(ai.Planets))
	copy(planetsToProtect, ai.Planets)
	planetsToProtectByNum := map[int]*cs.Planet{}
	for _, planet := range ai.Planets {
		planetsToProtectByNum[planet.Num] = planet
	}

	// find all idle fleets that have scanners
	mineLayerFleets := []*cs.Fleet{}
	for _, fleet := range ai.Fleets {
		if _, contains := fleet.Spec.Purposes[cs.ShipDesignPurposeDamageMineLayer]; contains && fleet.Spec.MineLayingRateByMineType != nil {
			if len(fleet.Waypoints) <= 1 && fleet.Waypoints[0].Task != cs.WaypointTaskLayMinefield {
				// this fleet can be sent to scan a planet
				mineLayerFleets = append(mineLayerFleets, fleet)
			} else {
				// this fleet is already scanning a planet, remove the target from the unknown planets list
				for _, wp := range fleet.Waypoints[1:] {
					if wp.TargetNum != cs.None {
						delete(planetsToProtectByNum, wp.TargetNum)
					}
				}
			}
		}
	}

	for _, fleet := range mineLayerFleets {
		closestPlanet := ai.getClosestPlanet(fleet, planetsToProtectByNum)
		if closestPlanet != nil {
			if fleet.Position == closestPlanet.Position {
				fleet.Waypoints[0].Task = cs.WaypointTaskLayMinefield
				fleet.Waypoints[0].LayMinefieldDuration = cs.Indefinite
				ai.client.UpdateFleetOrders(ai.Player, fleet, fleet.FleetOrders)
				delete(planetsToProtectByNum, closestPlanet.Num)
			} else {
				newWpIndex := fleet.AddWaypoint(ai.Player, cs.WaypointDest{MO: closestPlanet.MapObject}, len(fleet.Waypoints)-1, false)
				if newWpIndex == 0 {
					ai.log.Warn("Fleet tried to target planet for mine laying but did not add the waypoint",
						slog.String("fleet", fleet.Name),
						slog.String("planet", closestPlanet.Name))
					continue
				}

				ai.client.UpdateFleetOrders(ai.Player, fleet, fleet.FleetOrders)
				delete(planetsToProtectByNum, closestPlanet.Num)
			}
		}
	}

	// TODO: build new minelayer fleets when some conditions are true...

	return nil
}
