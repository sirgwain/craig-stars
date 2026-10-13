package ai

import (
	"github.com/sirgwain/craig-stars/cs"
	"log/slog"
)

// lay mines
func (ai *aiPlayer) layMines() error {
	if err := ai.detonateMinefields(); err != nil {
		return err
	}

	design := ai.Player.GetLatestDesign(cs.ShipDesignPurposeDamageMineLayer)
	// TODO: Add support for remote detonating minefield perimeters X years out in space if playing SD

	// Maintain the starting layers on every difficulty; hard replenishes them.
	hard := ai.AIDifficulty == cs.AIDifficultyHard || ai.AIDifficulty == cs.AIDifficultyCheater
	maintain := ai.doctrine.mineLayers || (hard && (ai.Race.PRT == cs.SD || ai.Race.PRT == cs.IS))
	if design == nil && !maintain {
		return nil
	}

	// lay mines to protect all our planets
	planetsToProtect := make([]*cs.Planet, len(ai.Planets))
	copy(planetsToProtect, ai.Planets)
	// a minefield is worth its upkeep over a world with a starbase to defend
	planetsToProtectByNum := map[int]*cs.Planet{}
	for _, planet := range ai.Planets {
		if planet.Spec.HasStarbase {
			planetsToProtectByNum[planet.Num] = planet
		}
	}

	// find all idle fleets that have scanners
	mineLayerFleets := []*cs.Fleet{}
	for _, fleet := range ai.Fleets {
		if _, contains := fleet.Spec.Purposes[cs.ShipDesignPurposeDamageMineLayer]; contains && fleet.Spec.MineLayingRateByMineType != nil {
			if len(fleet.Waypoints) <= 1 && fleet.Waypoints[0].Task != cs.WaypointTaskLayMinefield {
				// this fleet can be sent to scan a planet
				mineLayerFleets = append(mineLayerFleets, fleet)
			} else {
				if len(fleet.Waypoints) == 1 && fleet.Waypoints[0].Task == cs.WaypointTaskLayMinefield {
					delete(planetsToProtectByNum, fleet.OrbitingPlanetNum)
				}
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

	if maintain {
		ai.addFleetBuildRequest(cs.FleetPurposeMineLayer, min(ai.supportLimit(), len(planetsToProtectByNum)))
	}

	return nil
}

const (
	detonateEnemyShips = 5 // armed enemy ships in a minefield before it is worth detonating
	detonateShipRatio  = 3 // and they must outnumber our own ships in it by this much
)

// A race that can detonate its standard minefields damages every ship
// inside, its own included, and loses a quarter of the mines each year. So
// only do it while a real enemy force is in the field and few of our own
// ships are, and stop as soon as that is no longer true.
func (ai *aiPlayer) detonateMinefields() error {
	if !ai.Race.Spec.CanDetonateMinefields || !ai.enabled("detonate-minefields") {
		return nil
	}
	for _, minefield := range ai.Minefields {
		if !minefield.MinefieldType.CanDetonate() {
			continue
		}
		radius := minefield.Radius()
		within := func(position cs.Vector) bool {
			return float64(position.DistanceSquaredTo(minefield.Position)) <= radius*radius
		}

		enemies := 0
		for _, fleet := range ai.FleetIntels {
			if ai.IsEnemy(fleet.PlayerNum) && within(fleet.Position) && ai.hasAttackShips([]*cs.Fleet{fleet}) {
				for _, token := range fleet.Tokens {
					enemies += token.Quantity
				}
			}
		}
		// mine layers are immune to their own side's detonation
		ours := 0
		for _, fleet := range ai.Fleets {
			if within(fleet.Position) && len(fleet.Spec.MineLayingRateByMineType) == 0 {
				ours += fleet.Spec.TotalShips
			}
		}

		detonate := enemies >= detonateEnemyShips && enemies >= detonateShipRatio*ours
		if detonate != minefield.Detonate {
			if err := ai.client.UpdateMinefieldOrders(ai.Player, minefield, cs.MinefieldOrders{Detonate: detonate}); err != nil {
				return err
			}
		}
	}
	return nil
}
