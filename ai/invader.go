package ai

import (
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
)

func (ai *aiPlayer) invade() error {

	// make ure our invasions are valid
	for _, fleet := range ai.Fleets {
		if fleet.Purpose != cs.FleetPurposeInvader {
			continue
		}

		if len(fleet.Waypoints) == 1 {
			continue
		}

		target := ai.getPlanetIntel(fleet.Waypoints[1].TargetNum)

		// if this planet is no longer owned by a player, or it suddenly has a starbase, or it's pop has grown out
		// of the threshold where we would invade, return to the nearest starbase
		if !target.Owned() || target.Spec.HasStarbase || target.Spec.Population > int(ai.config.invasionFactor*float64(fleet.Cargo.Colonists*100)) {
			fleet.Purpose = cs.FleetPurposeNone
			closestStarbase := ai.getClosestStarbasePlanet(fleet)
			if closestStarbase != nil {
				warpSpeed := ai.getWarpSpeed(fleet, closestStarbase.Position)

				fleet.Waypoints[1] = cs.NewPlanetWaypoint(closestStarbase.Position, closestStarbase.Num, closestStarbase.Name, warpSpeed).
					WithTask(cs.WaypointTaskTransport).
					WithTransportTasks(cs.WaypointTransportTasks{Colonists: cs.WaypointTransportTask{Action: cs.TransportActionUnloadAll}})
				log.Debug().
					Int64("GameID", ai.GameID).
					Int("PlayerNum", ai.Num).
					Int("Invaders", fleet.Cargo.Colonists*100).
					Int("Defenders", target.Spec.Population).
					Bool("HasStarbase", target.Spec.HasStarbase).
					Msgf("%s called off invasion of %s, returning to %s", fleet.Name, target.Name, closestStarbase.Name)
			}
		}
	}

	return nil
}
