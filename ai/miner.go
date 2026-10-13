package ai

import "github.com/sirgwain/craig-stars/cs"

// AR can mine its own worlds remotely. Give every colony a miner rather than
// leaving the race permanently dependent on its small innate mining output.
func (ai *aiPlayer) mine() error {
	if !ai.Race.Spec.CanRemoteMineOwnPlanets {
		return nil
	}
	needing := map[int]*cs.Planet{}
	for _, planet := range ai.Planets {
		needing[planet.Num] = planet
	}
	idle := []*cs.Fleet{}
	for _, f := range ai.Fleets {
		if f.Spec.MiningRate <= 0 {
			continue
		}
		if len(f.Waypoints) > 1 {
			delete(needing, f.Waypoints[1].TargetNum)
			continue
		}
		if needing[f.OrbitingPlanetNum] != nil {
			if f.Waypoints[0].Task != cs.WaypointTaskRemoteMining {
				f.Waypoints[0].Task = cs.WaypointTaskRemoteMining
				ai.client.UpdateFleetOrders(ai.Player, f, f.FleetOrders)
			}
			delete(needing, f.OrbitingPlanetNum)
		} else {
			idle = append(idle, f)
		}
	}
	for _, f := range idle {
		if target := ai.getClosestPlanet(f, needing); target != nil {
			if i := f.AddWaypoint(ai.Player, cs.WaypointDest{MO: target.MapObject}, 0, false); i > 0 {
				f.Waypoints[i].Task = cs.WaypointTaskRemoteMining
				ai.client.UpdateFleetOrders(ai.Player, f, f.FleetOrders)
				delete(needing, target.Num)
			}
		}
	}
	ai.addFleetBuildRequest(cs.FleetPurposeMiner, min(ai.supportLimit(), len(needing)))
	return nil
}
