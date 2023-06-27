package cs

import (
	"fmt"
	"math"

	"github.com/rs/zerolog/log"
)

// producers perform planet production
type producer interface {
	produce() productionResult
}

// create a new planet production object
func newProducer(planet *Planet, player *Player) producer {
	return &production{
		planet: planet,
		player: player,
	}
}

type production struct {
	planet *Planet
	player *Player
}

type productionResult struct {
	tokens   []ShipToken
	packets  []Cargo
	scanner  bool
	starbase *ShipDesign
	alchemy  Mineral
}

type QueueItemCompletionEstimate struct {
	Skipped         bool    `json:"skipped,omitempty"`
	YearsToBuildOne int     `json:"yearsToBuildOne,omitempty"`
	YearsToBuildAll int     `json:"yearsToBuildAll,omitempty"`
	PercentComplete float64 `json:"percentComplete,omitempty"`
}

type ProductionQueueItem struct {
	QueueItemCompletionEstimate
	Type         QueueItemType `json:"type"`
	DesignNum    int           `json:"designNum,omitempty"`
	Quantity     int           `json:"quantity"`
	CostOfOne    Cost          `json:"costOfOne"`
	MaxBuildable int           `json:"maxBuildable"`
	Allocated    Cost          `json:"allocated"`
}

type QueueItemType string

const (
	QueueItemTypeIroniumMineralPacket   QueueItemType = "IroniumMineralPacket"
	QueueItemTypeBoraniumMineralPacket  QueueItemType = "BoraniumMineralPacket"
	QueueItemTypeGermaniumMineralPacket QueueItemType = "GermaniumMineralPacket"
	QueueItemTypeMixedMineralPacket     QueueItemType = "MixedMineralPacket"
	QueueItemTypeFactory                QueueItemType = "Factory"
	QueueItemTypeMine                   QueueItemType = "Mine"
	QueueItemTypeDefenses               QueueItemType = "Defenses"
	QueueItemTypeMineralAlchemy         QueueItemType = "MineralAlchemy"
	QueueItemTypeTerraformEnvironment   QueueItemType = "TerraformEnvironment"
	QueueItemTypeAutoMines              QueueItemType = "AutoMines"
	QueueItemTypeAutoFactories          QueueItemType = "AutoFactories"
	QueueItemTypeAutoDefenses           QueueItemType = "AutoDefenses"
	QueueItemTypeAutoMineralAlchemy     QueueItemType = "AutoMineralAlchemy"
	QueueItemTypeAutoMinTerraform       QueueItemType = "AutoMinTerraform"
	QueueItemTypeAutoMaxTerraform       QueueItemType = "AutoMaxTerraform"
	QueueItemTypeAutoMineralPacket      QueueItemType = "AutoMineralPacket"
	QueueItemTypeShipToken              QueueItemType = "ShipToken"
	QueueItemTypeStarbase               QueueItemType = "Starbase"
	QueueItemTypePlanetaryScanner       QueueItemType = "PlanetaryScanner"
)

// true if this is an auto type
func (t QueueItemType) IsAuto() bool {
	return t == QueueItemTypeAutoMines ||
		t == QueueItemTypeAutoFactories ||
		t == QueueItemTypeAutoDefenses ||
		t == QueueItemTypeAutoMineralAlchemy ||
		t == QueueItemTypeAutoMinTerraform ||
		t == QueueItemTypeAutoMaxTerraform ||
		t == QueueItemTypeAutoMineralPacket
}

// true if this is an auto type
func (t QueueItemType) IsPacket() bool {
	return t == QueueItemTypeAutoMineralPacket ||
		t == QueueItemTypeMixedMineralPacket ||
		t == QueueItemTypeIroniumMineralPacket ||
		t == QueueItemTypeBoraniumMineralPacket ||
		t == QueueItemTypeGermaniumMineralPacket
}

// return the concrete version of this auto type
func (t QueueItemType) concreteType() QueueItemType {
	switch t {
	case QueueItemTypeAutoMines:
		return QueueItemTypeMine
	case QueueItemTypeAutoFactories:
		return QueueItemTypeFactory
	case QueueItemTypeAutoDefenses:
		return QueueItemTypeDefenses
	case QueueItemTypeAutoMaxTerraform:
		return QueueItemTypeTerraformEnvironment
	case QueueItemTypeAutoMinTerraform:
		return QueueItemTypeTerraformEnvironment
	case QueueItemTypeAutoMineralAlchemy:
		return QueueItemTypeMineralAlchemy
	case QueueItemTypeAutoMineralPacket:
		return QueueItemTypeMixedMineralPacket
	}
	return t
}

// produce one turns worth of items from the production queue
func (p *production) produce() productionResult {
	player, planet := p.player, p.planet
	planet.PopulateProductionQueueCosts(player)
	result := productionResult{}
	available := Cost{Resources: planet.Spec.ResourcesPerYearAvailable}.AddCargoMinerals(planet.Cargo)
	newQueue := []ProductionQueueItem{}
	for itemIndex, item := range planet.ProductionQueue {
		maxBuildable := planet.maxBuildable(item.Type)
		if maxBuildable == Infinite {
			maxBuildable = math.MaxInt
		}

		// skip auto items that will never complete, or that we don't need
		// this way we can put auto terraforming items up top
		// and skip them to build factories, then mines
		if item.Skipped || (item.Type.IsAuto() && (item.YearsToBuildOne > 100 || item.YearsToBuildOne == Infinite)) {
			newQueue = append(newQueue, item)
			continue
		}

		if item.Type.IsPacket() && !planet.Spec.HasMassDriver {
			messager.buildMineralPacketNoMassDriver(player, planet)
			continue
		}
		if item.Type.IsPacket() && planet.PacketTargetNum == None {
			messager.buildMineralPacketNoTarget(player, planet)
			continue
		}

		// add in anything allocated in previous turns
		available = available.Add(item.Allocated)
		item.Allocated = Cost{}

		// get the cost of the current item
		cost := item.CostOfOne

		if (cost != Cost{}) {
			// figure out how many we can build
			// and make sure we only build up to the quantity, and we don't build more than the planet supports
			numBuilt := maxInt(0, available.NumBuildable(cost))
			numBuilt = minInt(numBuilt, item.Quantity)
			numBuilt = minInt(numBuilt, maxBuildable)

			if numBuilt > 0 {
				// build the items on the planet and remove from our available
				p.buildItems(item, numBuilt, &result)
				available = available.Minus(cost.MultiplyInt(numBuilt))
				available = available.Add(result.alchemy.ToCost())
			}

			if numBuilt < item.Quantity {
				// allocate to this item
				item.Allocated = p.allocatePartialBuild(cost, available)
				available = available.Minus(item.Allocated)
			}

			if item.Type.IsAuto() {
				if available.Resources == 0 {
					// we are out of resources, create a partial item end production
					if (item.Allocated != Cost{}) && numBuilt < item.Quantity {
						// we partially built an auto items, create a partial concrete item
						// we have some leftover to allocate so create a concrete item
						concreteItem := ProductionQueueItem{Type: item.Type.concreteType(), Quantity: 1, Allocated: item.Allocated}
						item.Allocated = Cost{}

						// add the concreate item to the top of the queue
						newQueue = append([]ProductionQueueItem{concreteItem}, newQueue...)
					}
					// auto items stay in the list
					newQueue = append(newQueue, item)

					if available.Resources == 0 {
						// we are out of resources, so we are done building
						if itemIndex < len(planet.ProductionQueue)-1 {
							// append the unfinished queue back to the end of our remaining items
							newQueue = append(newQueue, planet.ProductionQueue[itemIndex+1:]...)
						}
						break
					}
				} else {
					// auto items stay in the list
					// and we have resources leftover so move on
					item.Allocated = Cost{}
					newQueue = append(newQueue, item)
				}
			} else {
				item.Quantity -= numBuilt
				if item.Quantity != 0 {
					// we didn't finish, add the item back onto the remaining list
					newQueue = append(newQueue, item)
					if itemIndex < len(planet.ProductionQueue)-1 {
						// append the unfinished queue back to the end of our remaining items
						newQueue = append(newQueue, planet.ProductionQueue[itemIndex+1:]...)
					}
					// we finished, break out
					break
				}
			}
		}
	}
	// replace the queue with what's leftover
	planet.ProductionQueue = newQueue
	player.leftoverResources += available.Resources
	planet.Cargo = Cargo{available.Ironium, available.Boranium, available.Germanium, planet.Cargo.Colonists}

	return result
}

// add built items to planet, build fleets, update player messages, etc
func (p *production) buildItems(item ProductionQueueItem, numBuilt int, result *productionResult) {

	player, planet := p.player, p.planet

	switch item.Type {
	case QueueItemTypeAutoMineralAlchemy:
		fallthrough
	case QueueItemTypeMineralAlchemy:
		result.alchemy = Mineral{
			numBuilt,
			numBuilt,
			numBuilt,
		}
		messager.mineralAlchemyBuilt(player, planet, numBuilt)
	case QueueItemTypeAutoMines:
		fallthrough
	case QueueItemTypeMine:
		planet.Mines += numBuilt
		messager.minesBuilt(player, planet, numBuilt)
	case QueueItemTypeAutoFactories:
		fallthrough
	case QueueItemTypeFactory:
		planet.Factories += numBuilt
		messager.factoriesBuilt(player, planet, numBuilt)
	case QueueItemTypeAutoDefenses:
		fallthrough
	case QueueItemTypeDefenses:
		planet.Defenses += numBuilt
		messager.defensesBuilt(player, planet, numBuilt)
	case QueueItemTypeAutoMineralPacket:
		fallthrough
	case QueueItemTypeMixedMineralPacket:
		fallthrough
	case QueueItemTypeIroniumMineralPacket:
		fallthrough
	case QueueItemTypeBoraniumMineralPacket:
		fallthrough
	case QueueItemTypeGermaniumMineralPacket:
		// add this packet cargo to the production result
		// so it can be added as packets to the universe later
		cargo := player.Race.Spec.Costs[item.Type].MultiplyInt(numBuilt).ToCargo()
		result.packets = append(result.packets, cargo)
	case QueueItemTypeAutoMaxTerraform:
		fallthrough
	case QueueItemTypeAutoMinTerraform:
		fallthrough
	case QueueItemTypeTerraformEnvironment:
		terraformer := NewTerraformer()
		for i := 0; i < numBuilt; i++ {
			// terraform one at a time to ensure the best things get terraformed
			result := terraformer.TerraformOneStep(planet, player, nil, false)
			messager.terraform(player, planet, result.Type, result.Direction)
		}
	case QueueItemTypeShipToken:
		design := player.GetDesign(item.DesignNum)
		if design == nil {
			err := fmt.Errorf("tried to build a ship from design %d, but the design was not found", item.DesignNum)
			messager.error(player, err)
			log.Error().
				Int("Player", player.Num).
				Str("Planet", planet.Name).
				Str("Item", string(item.Type)).
				Int("DesignNum", item.DesignNum).
				Int("NumBuilt", numBuilt).
				Err(err).
				Msg("failed to build ShipToken")
			break
		}
		design.Spec.NumBuilt += numBuilt
		design.Spec.NumInstances += numBuilt
		design.MarkDirty()
		result.tokens = append(result.tokens, ShipToken{Quantity: numBuilt, design: design, DesignNum: design.Num})
	case QueueItemTypeStarbase:
		design := player.GetDesign(item.DesignNum)
		if design == nil {
			err := fmt.Errorf("tried to build a starbase from design %d, but the design was not found", item.DesignNum)
			messager.error(player, err)
			log.Error().
				Int("Player", player.Num).
				Str("Planet", planet.Name).
				Str("Item", string(item.Type)).
				Int("DesignNum", item.DesignNum).
				Int("NumBuilt", numBuilt).
				Err(err).
				Msg("failed to build Starbase")
			break
		}
		result.starbase = design
	case QueueItemTypePlanetaryScanner:
		result.scanner = true
	}

	log.Debug().
		Int("Player", player.Num).
		Str("Planet", planet.Name).
		Str("Item", string(item.Type)).
		Int("DesignNum", item.DesignNum).
		Int("NumBuilt", numBuilt).
		Msgf("built item")

}

// Allocate resources to the top item on this production queue
// and return the leftover resources
//
// Costs are allocated by lowest percentage, i.e. if (we require
// Cost(10, 10, 10, 100) and we only have Cost(1, 10, 10, 100)
// we allocate Cost(1, 1, 1, 10)
//
// The min amount we have is 10 percent of the ironium, so we
// apply 10 percent to each cost amount
func (p *production) allocatePartialBuild(costPerItem Cost, allocated Cost) Cost {
	ironiumPerc := 100.0
	if costPerItem.Ironium > 0 {
		ironiumPerc = float64(allocated.Ironium) / float64(costPerItem.Ironium)
	}
	boraniumPerc := 100.0
	if costPerItem.Boranium > 0 {
		boraniumPerc = float64(allocated.Boranium) / float64(costPerItem.Boranium)
	}
	germaniumPerc := 100.0
	if costPerItem.Germanium > 0 {
		germaniumPerc = float64(allocated.Germanium) / float64(costPerItem.Germanium)
	}
	resourcesPerc := 100.0
	if costPerItem.Resources > 0 {
		resourcesPerc = float64(allocated.Resources) / float64(costPerItem.Resources)
	}

	// figure out the lowest percentage
	minPerc := minFloat64(ironiumPerc, boraniumPerc, germaniumPerc, resourcesPerc)

	// allocate the lowest percentage of each cost
	newAllocated := Cost{
		int(float64(costPerItem.Ironium) * minPerc),
		int(float64(costPerItem.Boranium) * minPerc),
		int(float64(costPerItem.Germanium) * minPerc),
		int(float64(costPerItem.Resources) * minPerc),
	}

	// return the amount we allocate to the top queued item
	return newAllocated
}
