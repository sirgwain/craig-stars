package ai

import (
	"github.com/sirgwain/craig-stars/cs"
	"math"
)

// Colonizers need a destination and passengers, not a fixed escort recipe.
// This also lets the starting colony ships participate immediately.
func (ai *aiPlayer) colonize() error {
	targets := map[int]*cs.Planet{}
	// worlds we can only live on after terraforming
	terraformTargets := map[int]*cs.Planet{}
	for _, planet := range ai.PlanetIntels {
		if !planet.Explored() || planet.Owned() || ai.hasAttackShips(ai.enemyShipsAbovePlanet(planet)) {
			continue
		}
		if planet.Spec.Habitability > 0 || (ai.Race.Spec.Instaforming && planet.Spec.TerraformedHabitability > 0) {
			targets[planet.Num] = planet
		} else if planet.Spec.TerraformedHabitability > 0 && ai.AIDifficulty != cs.AIDifficultyEasy {
			terraformTargets[planet.Num] = planet
		}
	}
	idle := []*cs.Fleet{}
	active := 0
	for _, f := range ai.Fleets {
		if !f.Spec.Colonizer || f.Purpose == cs.FleetPurposeInvader {
			continue
		}
		f.Purpose = cs.FleetPurposeColonizer
		f.SetTag(cs.TagPurpose, string(cs.FleetPurposeColonizer))
		if len(f.Waypoints) > 1 {
			wp := f.Waypoints[1]
			target := ai.GetPlanetIntel(wp.TargetNum)
			if wp.TargetType == cs.MapObjectTypePlanet && target != nil && !target.Owned() {
				delete(targets, target.Num)
				delete(terraformTargets, target.Num)
				active++
				continue
			}
			// A colony claimed in transit needs a new destination, including when
			// another fleet of ours got there first.
			f.Waypoints = f.Waypoints[:1]
		}
		idle = append(idle, f)
	}
	// Colonists die on a hostile world until it is terraformed, so only settle
	// one when there is nowhere better to go.
	if len(targets) == 0 {
		targets = terraformTargets
	}
	for _, f := range idle {
		target := ai.getBestPlanetToColonize(f, targets)
		if target == nil {
			// Nowhere left to settle. Colonists do more good on one of our
			// worlds than waiting aboard, so take them to the nearest.
			if f.Cargo.Colonists > 0 {
				if err := ai.returnColonists(f); err != nil {
					return err
				}
			}
			continue
		}
		if f.Cargo.Colonists == 0 {
			source := ai.getPlanet(f.OrbitingPlanetNum)
			if source == nil {
				if home := ai.getClosestStarbasePlanet(f); home != nil {
					f.AddWaypoint(ai.Player, cs.WaypointDest{MO: home.MapObject}, 0, false)
					ai.client.UpdateFleetOrders(ai.Player, f, f.FleetOrders)
				}
				continue
			}
			desired := max(minColonistLoad, int(float64(target.Spec.MaxPopulation)*.10)/100)
			amount := min(desired, f.Spec.CargoCapacity-f.Cargo.Total(), ai.exportableColonists(source, ai.config.colonizerPopulationDensity))
			if seed := min(minColonistLoad, f.Spec.CargoCapacity-f.Cargo.Total()); amount < seed && ai.colonyHabitability(target) >= ai.seedHabitability() && source.Cargo.Colonists >= seed*seedColonySourceFactor {
				// A ship we already have claims a good world with a small load
				// rather than waiting years for the source to fill up.
				amount = seed
			}
			if amount < min(minColonistLoad, f.Spec.CargoCapacity) || amount <= 0 {
				continue
			}
			if err := ai.client.TransferByHand(&ai.game.Rules, ai.Player, f, source, cs.CargoTransferRequest{Cargo: cs.Cargo{Colonists: amount}}); err != nil {
				return err
			}
		}
		if index := f.AddWaypoint(ai.Player, cs.WaypointDest{MO: target.MapObject}, 0, true); index > 0 {
			f.Waypoints[index].Task = cs.WaypointTaskColonize
			ai.client.UpdateFleetOrders(ai.Player, f, f.FleetOrders)
			delete(targets, target.Num)
			active++
		}
	}
	// Bound commitments by both known destinations and a manageable pipeline.
	// Count every idle colony ship, including incomplete legacy escort fleets.
	idleShips := 0
	for _, f := range ai.Fleets {
		if f.Spec.Colonizer && len(f.Waypoints) <= 1 {
			for _, token := range f.Tokens {
				if d := ai.GetDesign(token.DesignNum); d != nil && d.Spec.Colonizer {
					idleShips += token.Quantity
				}
			}
		}
	}
	demand := max(0, min(len(targets), ai.colonyLimit()-active)-idleShips)
	ai.addFleetBuildRequest(cs.FleetPurposeColonizer, demand)
	return nil
}

// the habitability a world needs before an idle colony ship seeds it early
func (ai *aiPlayer) seedHabitability() int {
	if ai.doctrine.seedHabitability > 0 {
		return ai.doctrine.seedHabitability
	}
	return seedColonyHabitability
}

// returnColonists unloads a colony ship at the world it orbits, or sends it
// to our nearest world to unload there
func (ai *aiPlayer) returnColonists(f *cs.Fleet) error {
	if planet := ai.getPlanet(f.OrbitingPlanetNum); planet != nil {
		return ai.client.TransferByHand(&ai.game.Rules, ai.Player, f, planet, cs.CargoTransferRequest{Cargo: cs.Cargo{Colonists: -f.Cargo.Colonists}})
	}
	owned := map[int]*cs.Planet{}
	for _, planet := range ai.Planets {
		owned[planet.Num] = planet
	}
	if nearest := ai.getClosestPlanet(f, owned); nearest != nil {
		if index := f.AddWaypoint(ai.Player, cs.WaypointDest{MO: nearest.MapObject}, 0, false); index > 0 {
			f.Waypoints[index].Task = cs.WaypointTaskTransport
			f.Waypoints[index].TransportTasks = cs.WaypointTransportTasks{Colonists: cs.WaypointTransportTask{Action: cs.TransportActionUnloadAll}}
			ai.client.UpdateFleetOrders(ai.Player, f, f.FleetOrders)
		}
	}
	return nil
}

// the smallest colonist load, in kT, worth sending a ship for
const minColonistLoad = 25

// An idle colony ship may take a minimum load from a source below its reserve
// when the world is at least this habitable and the source has this many times
// the load.
const (
	seedColonyHabitability = 70
	seedColonySourceFactor = 10
)

// Cargo is measured in hundreds of colonists. Preserve the breeder population
// using people present now; next year's births cannot be loaded this year.
func (ai *aiPlayer) exportableColonists(planet *cs.Planet, density float64) int {
	reserve := int(math.Ceil(float64(planet.Spec.MaxPopulation) * density / 100))
	return max(0, planet.Cargo.Colonists-reserve)
}

// get the planet with the best distance to hab ratio
func (ai *aiPlayer) getBestPlanetToColonize(fleet *cs.Fleet, colonizablePlanets map[int]*cs.Planet) *cs.Planet {
	var best *cs.Planet = nil

	// lowest weight wins
	bestWeight := math.MaxFloat64
	yearlyTravelDistance := float64(fleet.Spec.Engine.IdealSpeed * fleet.Spec.Engine.IdealSpeed)

	for _, intel := range colonizablePlanets {
		habValue := ai.colonyHabitability(intel)
		terraformHabValue := intel.Spec.TerraformedHabitability
		dist := intel.Position.DistanceTo(fleet.Position)
		yearsToTravel := math.Ceil(dist / yearlyTravelDistance)

		// weight is based on distance (in years) and hab value
		// as distance goes up, weight goes up. As habValue goes up, weight goes down
		// low weight wins
		// distance is weighted heavier than habValue. In the above calc, a 100% planet 10y away is ranked equally with a
		// 25% planet 5y away
		weight := math.MaxFloat64
		if habValue > 0 {
			weight = (yearsToTravel * yearsToTravel) / float64(habValue)
		} else if terraformHabValue > 0 {
			// only colonize terraformable planets if they're really close
			weight = 2 * (yearsToTravel * yearsToTravel) / float64(terraformHabValue)
		}

		if ai.AIDifficulty == cs.AIDifficultyHard || ai.AIDifficulty == cs.AIDifficultyCheater {
			// Rich colonies start factories sooner; use only minerals in explored intel.
			weight /= 1 + float64(intel.MineralConcentration.Germanium+intel.MineralConcentration.Ironium)/200
		}
		if weight < bestWeight || (weight == bestWeight && best != nil && intel.Num < best.Num) {
			bestWeight = weight
			best = intel
		}
	}

	// best = ai.getClosestPlanetIntel(fleet.Position, colonizablePlanets)

	return best
}
