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
			if wp1.WarpSpeed == fleet.Spec.Engine.FreeSpeed {
				warpSpeed := ai.getWarpSpeed(fleet, wp1.Position)
				if warpSpeed >= fleet.Spec.Engine.IdealSpeed {
					ai.log.Debug().
						Int64("GameID", ai.GameID).
						Int("PlayerNum", ai.Num).
						Msgf("Fleet %s increasing warp from %d to %d", fleet.Name, wp1.WarpSpeed, warpSpeed)

					fleet.Waypoints[1].WarpSpeed = warpSpeed
					ai.client.UpdateFleetOrders(ai.Player, fleet, fleet.FleetOrders)
				}
			}
		}
	}

	return nil
}

func (ai *aiPlayer) getWarpSpeed(fleet *cs.Fleet, position cs.Vector) int {
	dist := fleet.Position.DistanceTo(position)
	return cs.Clamp(ai.getMaxWarp(dist, fleet), 1, 10)
}

// get the maximum warp we can travel to reach the destination
// in the minimal number of years, within our fuel constraints
func (ai *aiPlayer) getMaxWarp(dist float64, fleet *cs.Fleet) int {
	freeSpeed := fleet.Spec.Engine.FreeSpeed

	// start at freespeed+1 and move up until we run out of fuel
	var speed int
	for speed = freeSpeed + 1; speed < fleet.Spec.Engine.MaxSafeSpeed; speed++ {
		fuelUsed := fleet.GetFuelCost(ai.Player, speed, dist)

		// we are using too much fuel, go to the previous speed
		if fuelUsed > fleet.Fuel {
			speed--
			break
		}
	}

	idealSpeed := fleet.Spec.Engine.IdealSpeed

	// if we are using a ramscoop, make sure we at least go the ideal
	// speed of the engine. If we run out, oh well, it'll drop to
	// the free speed
	if freeSpeed > 1 && speed < idealSpeed {
		speed = idealSpeed
	}

	// don't go faster than we need
	return ai.getMinimalWarp(dist, speed, fleet)
}

// get the minimal warp starting at an idealSpeed and working downward
// if we can travel the same amount of time at a lower speed, do it
// TODO: Add fuel requirements?
func (ai *aiPlayer) getMinimalWarp(dist float64, idealSpeed int, fleet *cs.Fleet) int {
	speed := idealSpeed

	freeSpeed := fleet.Spec.Engine.FreeSpeed

	// travelling 49 ly at warp 7 takes one year
	yearsAtIdealSpeed := int(math.Ceil(dist / float64(idealSpeed*idealSpeed)))
	for i := idealSpeed; i > freeSpeed; i-- {
		yearsAtSpeed := int(math.Ceil(dist / float64(i*i)))

		// It takes the same time to go slower, so go slower
		if yearsAtIdealSpeed == yearsAtSpeed {
			speed = i
		}
	}

	return cs.Clamp(speed, 1, 10)
}
