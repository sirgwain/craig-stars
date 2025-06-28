package ai

import (
	"math"

	"github.com/sirgwain/craig-stars/cs"
)

// make sure our fleets are all going at the speed they should
func (ai *aiPlayer) updateFleetWarpSpeed() error {
	for _, fleet := range ai.Fleets {
		if len(fleet.Waypoints) == 2 {
			wp1 := fleet.Waypoints[1]
			// this is set to free speed, see if we can make it go faster
			if wp1.WarpSpeed == fleet.Spec.Engine.FreeSpeed {
				warpSpeed := ai.getWarpSpeed(fleet, wp1.MapObjectTarget, false)
				if warpSpeed > wp1.WarpSpeed {
					ai.log.Debug().
						Msgf("Fleet %s increasing warp from %d to %d", fleet.Name, wp1.WarpSpeed, warpSpeed)

					fleet.Waypoints[1].WarpSpeed = warpSpeed
					ai.client.UpdateFleetOrders(ai.Player, fleet, fleet.FleetOrders)
				}
			}
		}
	}

	return nil
}

// getWarpSpeed returns the optimal warp speed to a target
func (ai *aiPlayer) getWarpSpeed(fleet *cs.Fleet, target cs.MapObjectTarget, fastest bool) int {
	dist := math.Ceil(fleet.Position.DistanceTo(target.TargetPosition))

	var orbiting *cs.PlanetIntel
	var targetPlanet *cs.PlanetIntel
	if fleet.OrbitingPlanetNum != cs.None {
		orbiting = ai.GetPlanetIntel(fleet.OrbitingPlanetNum)
	}
	if target.TargetType == cs.MapObjectTypePlanet {
		targetPlanet = ai.GetPlanetIntel(target.TargetNum)
	}

	return fleet.GetWarpSpeed(ai.Player, dist, orbiting, targetPlanet, 0, fastest)
}
