package ai

import (
	"cmp"
	"slices"

	"github.com/sirgwain/craig-stars/cs"
)

const (
	haulMinimumLoad      = 50  // the smallest mineral load, in kT, worth a trip
	haulSourceReserve    = 500 // kT of each mineral a source keeps for its own ships and starbases
	haulFactoriesPerTrip = 100 // deliver germanium for at most this many factories at once
)

// Factories need germanium, and a young colony rarely mines enough of its
// own. A colony's first starbase or ships can wait years for ironium or
// boranium too. Transports with no colonists to move carry minerals from
// worlds with a surplus to the worlds that are waiting for them.
func (ai *aiPlayer) haulMinerals() error {
	if ai.Race.Spec.InnateResources {
		return nil
	}
	perFactory := ai.game.Rules.FactoryCostGermanium
	if ai.Race.FactoriesCostLess {
		perFactory--
	}

	// what each planet needs beyond what it has and will mine this year
	needs := func(planet *cs.Planet) cs.Mineral {
		need := cs.Mineral{Germanium: max(0, planet.Spec.MaxFactories-planet.Factories) * perFactory}
		if ai.enabled("haul-minerals") {
			// and the ship or starbase at the front of its queue
			for _, item := range planet.ProductionQueue {
				if item.Type != cs.QueueItemTypeShipToken && item.Type != cs.QueueItemTypeStarbase {
					continue
				}
				if design := ai.GetDesign(item.DesignNum); design != nil {
					cost := design.Spec.Cost
					need = need.Add(cs.Mineral{
						Ironium:   max(0, cost.Ironium*item.Quantity-item.Allocated.Ironium),
						Boranium:  max(0, cost.Boranium*item.Quantity-item.Allocated.Boranium),
						Germanium: max(0, cost.Germanium*item.Quantity-item.Allocated.Germanium),
					})
				}
				break
			}
		}
		return need
	}

	surplus := map[int]cs.Mineral{}
	wanted := map[int]cs.Mineral{}
	planets := slices.Clone(ai.Planets)
	slices.SortFunc(planets, func(a, b *cs.Planet) int { return cmp.Compare(a.Num, b.Num) })
	for _, planet := range planets {
		need := needs(planet)
		stock := planet.Cargo.ToMineral()
		spare := cs.Mineral{
			Ironium:   max(0, stock.Ironium-need.Ironium-haulSourceReserve),
			Boranium:  max(0, stock.Boranium-need.Boranium-haulSourceReserve),
			Germanium: max(0, stock.Germanium-need.Germanium-haulSourceReserve),
		}
		// germanium for factories arrives a trip at a time
		need.Germanium = min(need.Germanium, haulFactoriesPerTrip*perFactory+need.Germanium-max(0, planet.Spec.MaxFactories-planet.Factories)*perFactory)
		short := cs.Mineral{
			Ironium:   max(0, need.Ironium-stock.Ironium-planet.Spec.MiningOutput.Ironium),
			Boranium:  max(0, need.Boranium-stock.Boranium-planet.Spec.MiningOutput.Boranium),
			Germanium: max(0, need.Germanium-stock.Germanium),
		}
		if spare.Total() >= haulMinimumLoad {
			surplus[planet.Num] = spare
		}
		if short.Total() >= haulMinimumLoad {
			wanted[planet.Num] = short
		}
	}

	fleets := []*cs.Fleet{}
	for _, fleet := range ai.Fleets {
		if fleet.GetTag(cs.TagPurpose) != string(cs.FleetPurposeColonistFreighter) || fleet.Spec.Colonizer || fleet.Purpose == cs.FleetPurposeInvader {
			continue
		}
		if !fleet.Idle() {
			// minerals already on their way
			if cargo := fleet.Cargo.ToMineral(); cargo.Total() > 0 && len(fleet.Waypoints) > 1 {
				target := fleet.Waypoints[1].TargetNum
				wanted[target] = mineralsLess(wanted[target], cargo)
			}
			continue
		}
		if fleet.Cargo.Total() == 0 {
			fleets = append(fleets, fleet)
		}
	}
	if len(surplus) == 0 {
		return nil
	}

	// what a source can usefully send a world, mineral by mineral
	load := func(source, target int, capacity int) cs.Mineral {
		spare, short := surplus[source], wanted[target]
		cargo := cs.Mineral{}
		cargo.Germanium = min(spare.Germanium, short.Germanium, capacity)
		cargo.Ironium = min(spare.Ironium, short.Ironium, capacity-cargo.Total())
		cargo.Boranium = min(spare.Boranium, short.Boranium, capacity-cargo.Total())
		return cargo
	}

	for _, fleet := range fleets {
		source := ai.getPlanet(fleet.OrbitingPlanetNum)
		if source == nil || surplus[source.Num].Total() < haulMinimumLoad {
			// go and wait at the nearest world with minerals some world is waiting for
			sources := map[int]*cs.Planet{}
			for num := range surplus {
				for target := range wanted {
					if num != target && load(num, target, fleet.Spec.CargoCapacity).Total() >= haulMinimumLoad {
						sources[num] = ai.getPlanet(num)
						break
					}
				}
			}
			if nearest := ai.getClosestPlanet(fleet, sources); nearest != nil {
				if index := fleet.AddWaypoint(ai.Player, cs.WaypointDest{MO: nearest.MapObject}, 0, false); index > 0 {
					ai.client.UpdateFleetOrders(ai.Player, fleet, fleet.FleetOrders)
				}
			}
			continue
		}

		// deliver to the nearest world this source has a worthwhile load for
		waiting := map[int]*cs.Planet{}
		for num := range wanted {
			if num != source.Num && load(source.Num, num, fleet.Spec.CargoCapacity).Total() >= haulMinimumLoad {
				waiting[num] = ai.getPlanet(num)
			}
		}
		target := ai.getClosestPlanet(fleet, waiting)
		if target == nil {
			continue
		}
		minerals := load(source.Num, target.Num, fleet.Spec.CargoCapacity)
		cargo := minerals.ToCargo()
		if err := ai.client.TransferByHand(&ai.game.Rules, ai.Player, fleet, source, cs.CargoTransferRequest{Cargo: cargo}); err != nil {
			return err
		}
		index := fleet.AddWaypoint(ai.Player, cs.WaypointDest{MO: target.MapObject}, 0, false)
		if index == 0 {
			// we can't get there, put it back
			if err := ai.client.TransferByHand(&ai.game.Rules, ai.Player, fleet, source, cs.CargoTransferRequest{Cargo: cargo.Negative()}); err != nil {
				return err
			}
			continue
		}
		unload := cs.WaypointTransportTask{Action: cs.TransportActionUnloadAll}
		fleet.Waypoints[index].Task = cs.WaypointTaskTransport
		fleet.Waypoints[index].TransportTasks = cs.WaypointTransportTasks{Ironium: unload, Boranium: unload, Germanium: unload}
		ai.client.UpdateFleetOrders(ai.Player, fleet, fleet.FleetOrders)
		surplus[source.Num] = mineralsLess(surplus[source.Num], minerals)
		wanted[target.Num] = mineralsLess(wanted[target.Num], minerals)
	}
	return nil
}

// mineralsLess subtracts minerals, never going below zero
func mineralsLess(a, b cs.Mineral) cs.Mineral {
	return cs.Mineral{
		Ironium:   max(0, a.Ironium-b.Ironium),
		Boranium:  max(0, a.Boranium-b.Boranium),
		Germanium: max(0, a.Germanium-b.Germanium),
	}
}
