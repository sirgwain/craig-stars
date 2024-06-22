package ai

import "github.com/rs/zerolog/log"

// make sure our fleets are all going the speed they should
func (ai *aiPlayer) updateFleetWarpSpeed() error {

	for _, fleet := range ai.Fleets {
		if len(fleet.Waypoints) == 2 {
			wp1 := fleet.Waypoints[1]
			if wp1.WarpSpeed == fleet.Spec.Engine.FreeSpeed {
				warpSpeed := ai.getWarpSpeed(fleet, wp1.Position)
				if warpSpeed >= fleet.Spec.Engine.IdealSpeed {
					log.Debug().
						Int64("GameID", ai.GameID).
						Int("PlayerNum", ai.Num).
						Msgf("Fleet %s increasing warp from %d to %d", fleet.Name, wp1.WarpSpeed, warpSpeed)

					fleet.Waypoints[1].WarpSpeed = warpSpeed
					ai.client.UpdateFleetOrders(ai.Player, fleet, fleet.FleetOrders, fleet.Tags)
				}
			}
		}
	}

	return nil
}
