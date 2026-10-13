package ai

import (
	"log/slog"
	"math"

	"github.com/sirgwain/craig-stars/cs"
)

func (ai *aiPlayer) invade() error {
	// enemy worlds already being invaded
	targeted := map[int]bool{}

	// make sure our invasions are valid
	for _, fleet := range ai.Fleets {
		if fleet.Purpose != cs.FleetPurposeInvader {
			continue
		}

		if len(fleet.Waypoints) == 1 {
			// the invasion is over, this fleet can transport colonists again
			fleet.Purpose = cs.FleetPurposeColonistFreighter
			continue
		}

		target := ai.GetPlanetIntel(fleet.Waypoints[1].TargetNum)

		// if this planet is no longer owned by an enemy, or it suddenly has a starbase, or its pop has grown out
		// of the threshold where we would invade, return to the nearest starbase
		if target == nil || !ai.canInvade(target, fleet.Cargo.Colonists) {
			fleet.Purpose = cs.FleetPurposeNone
			closestStarbase := ai.getClosestStarbasePlanet(fleet)
			if closestStarbase != nil {
				warpSpeed := ai.getWarpSpeed(fleet, closestStarbase.ToTarget(), true)

				fleet.Waypoints[1] = cs.NewPlanetWaypoint(closestStarbase.Position, closestStarbase.Num, closestStarbase.Name, warpSpeed).
					WithTask(cs.WaypointTaskTransport).
					WithTransportTasks(cs.WaypointTransportTasks{Colonists: cs.WaypointTransportTask{Action: cs.TransportActionUnloadAll}})
				ai.client.UpdateFleetOrders(ai.Player, fleet, fleet.FleetOrders)
				ai.log.Debug("Fleet called off invasion, returning to starbase",
					slog.String("fleet", fleet.Name),
					slog.String("starbase", closestStarbase.Name),
					slog.Int("Invaders", fleet.Cargo.Colonists*100))
			}
			continue
		}
		targeted[target.Num] = true
	}

	// send idle colonist freighters after weakly held enemy worlds
	for _, fleet := range ai.Fleets {
		if fleet.GetTag(cs.TagPurpose) != string(cs.FleetPurposeColonistFreighter) || fleet.Spec.Colonizer || !fleet.Idle() || fleet.Cargo.Total() != 0 {
			continue
		}
		source := ai.getPlanet(fleet.OrbitingPlanetNum)
		if source == nil {
			continue
		}
		available := min(fleet.Spec.CargoCapacity, ai.exportableColonists(source, ai.config.colonistTransportDensity))

		// invade the closest world we can overwhelm
		var target *cs.Planet
		shortestDist := math.MaxInt
		for _, planet := range ai.PlanetIntels {
			if targeted[planet.Num] || !ai.canInvade(planet, available) {
				continue
			}
			dist := fleet.Position.DistanceSquaredTo(planet.Position)
			if ai.enabled("bomb-then-invade") && ai.bombing(planet) {
				// a world our bombers are already over is the one to land on
				dist -= 1 << 40
			}
			if dist < shortestDist || (dist == shortestDist && planet.Num < target.Num) {
				shortestDist = dist
				target = planet
			}
		}
		if target == nil {
			continue
		}

		if err := ai.client.TransferByHand(&ai.game.Rules, ai.Player, fleet, source, cs.CargoTransferRequest{Cargo: cs.Cargo{Colonists: available}}); err != nil {
			return err
		}
		index := fleet.AddWaypoint(ai.Player, cs.WaypointDest{MO: target.MapObject}, 0, true)
		if index == 0 {
			// we can't get there, put the colonists back
			if err := ai.client.TransferByHand(&ai.game.Rules, ai.Player, fleet, source, cs.CargoTransferRequest{Cargo: cs.Cargo{Colonists: -available}}); err != nil {
				return err
			}
			continue
		}
		fleet.Purpose = cs.FleetPurposeInvader
		fleet.Waypoints[index].Task = cs.WaypointTaskTransport
		fleet.Waypoints[index].TransportTasks = cs.WaypointTransportTasks{Colonists: cs.WaypointTransportTask{Action: cs.TransportActionUnloadAll}}
		ai.client.UpdateFleetOrders(ai.Player, fleet, fleet.FleetOrders)
		targeted[target.Num] = true

		ai.log.Debug("Fleet invading planet",
			slog.String("fleet", fleet.Name),
			slog.String("target", target.Name),
			slog.Int("Invaders", available*100),
			slog.Int("Defenders", target.GetPopulation()))
	}

	return nil
}

// the share of a planet's defense coverage that applies against invaders
const invasionDefenseFactor = .75

// bombing returns true if one of our bomber fleets is in orbit of this planet
func (ai *aiPlayer) bombing(planet *cs.Planet) bool {
	for _, fleet := range ai.fleetsByPlanetNum[planet.Num] {
		if fleet.Spec.Bomber {
			return true
		}
	}
	return false
}

// canInvade returns true if this many kT of colonists can overwhelm an enemy
// world, using only what we have discovered about it.
func (ai *aiPlayer) canInvade(planet *cs.Planet, colonists int) bool {
	if !planet.Explored() || !planet.Owned() || !ai.IsEnemy(planet.PlayerNum) || planet.Spec.HasStarbase {
		return false
	}
	defenders := planet.GetPopulation()
	invaders := float64(colonists * 100)
	if ai.enabled("bomb-then-invade") {
		// planetary defenses kill invaders on the way down
		invaders *= 1 - invasionDefenseFactor*planet.Spec.DefenseCoverage
	}
	return defenders > 0 && invaders >= ai.config.invasionFactor*float64(defenders)
}
