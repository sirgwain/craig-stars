package ai

import "github.com/sirgwain/craig-stars/cs"

// lay mines
func (ai *aiPlayer) layMines() error {
	design := ai.Player.GetLatestDesign(cs.ShipDesignPurposeDamageMineLayer)

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
			if len(fleet.Waypoints) <= 1 && fleet.Waypoints[0].Task != cs.WaypointTaskLayMineField {
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
				fleet.Waypoints[0].Task = cs.WaypointTaskLayMineField
				fleet.Waypoints[0].LayMineFieldDuration = cs.Indefinite
				ai.client.UpdateFleetOrders(ai.Player, fleet, fleet.FleetOrders)
				delete(planetsToProtectByNum, closestPlanet.Num)
			} else {
				warpSpeed := ai.getWarpSpeed(fleet, closestPlanet.Position)
				fleet.Waypoints = append(fleet.Waypoints, cs.NewPlanetWaypoint(closestPlanet.Position, closestPlanet.Num, closestPlanet.Name, warpSpeed))
				ai.client.UpdateFleetOrders(ai.Player, fleet, fleet.FleetOrders)
				delete(planetsToProtectByNum, closestPlanet.Num)
			}
		}
	}

	// TODO: build new minelayer fleets when some conditions are true...

	return nil
}
