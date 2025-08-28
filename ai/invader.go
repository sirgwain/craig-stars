package ai

import (
	"github.com/sirgwain/craig-stars/cs"
	"log/slog"
)

func (ai *aiPlayer) invade() error {

	// make sure our invasions are valid
	for _, fleet := range ai.Fleets {
		if fleet.Purpose != cs.FleetPurposeInvader {
			continue
		}

		if len(fleet.Waypoints) == 1 {
			continue
		}

		target := ai.GetPlanetIntel(fleet.Waypoints[1].TargetNum)

		// if this planet is no longer owned by a player, or it suddenly has a starbase, or its pop has grown out
		// of the threshold where we would invade, return to the nearest starbase
		if !target.Owned() || target.Spec.HasStarbase || target.GetPopulation() > int(ai.config.invasionFactor*float64(fleet.Cargo.Colonists*100)) {
			fleet.Purpose = cs.FleetPurposeNone
			closestStarbase := ai.getClosestStarbasePlanet(fleet)
			if closestStarbase != nil {
				warpSpeed := ai.getWarpSpeed(fleet, closestStarbase.ToTarget(), true)

				fleet.Waypoints[1] = cs.NewPlanetWaypoint(closestStarbase.Position, closestStarbase.Num, closestStarbase.Name, warpSpeed).
					WithTask(cs.WaypointTaskTransport).
					WithTransportTasks(cs.WaypointTransportTasks{Colonists: cs.WaypointTransportTask{Action: cs.TransportActionUnloadAll}})
				ai.log.Debug("Fleet called off invasion, returning to starbase",
					slog.String("fleet", fleet.Name),
					slog.String("target", target.Name),
					slog.String("starbase", closestStarbase.Name),
					slog.Int("Invaders", fleet.Cargo.Colonists*100),
					slog.Int("Defenders", target.GetPopulation()),
					slog.Bool("HasStarbase", target.Spec.HasStarbase))
			}
		}
	}

	return nil
}
