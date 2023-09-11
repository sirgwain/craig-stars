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
		planet:    planet,
		player:    player,
		estimator: NewCompletionEstimator(),
	}
}

type production struct {
	planet    *Planet
	player    *Player
	estimator CompletionEstimator
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
	index        int           // used for holding a place in the queue while estimating
	design       *ShipDesign
}

func (item *ProductionQueueItem) SetDesign(design *ShipDesign) {
	item.design = design
}

// get the percent this build item has been completed
func (item ProductionQueueItem) percentComplete() float64 {
	if item.Allocated.Total() == 0 {
		return 0
	}

	// update the percent complete based on how much we've allocated vs the total cost of all items
	costOfAll := item.CostOfOne.MultiplyInt(item.Quantity)
	if !(item.Allocated == Cost{}) {
		return clampFloat64(item.Allocated.Divide(costOfAll), 0, 1)
	}
	return 0
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

// true if this is a terraform type
func (t QueueItemType) IsTerraform() bool {
	return t == QueueItemTypeAutoMaxTerraform ||
		t == QueueItemTypeAutoMinTerraform ||
		t == QueueItemTypeTerraformEnvironment
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

type productionResult struct {
	itemsBuilt        []itemBuilt
	leftoverResources int
	tokens            []ShipToken
	packets           []Cargo
	scanner           bool
	starbase          *ShipDesign
	alchemy           Mineral
	mines             int
	factories         int
	defenses          int
	terraformResults  []TerraformResult
	messages          []PlayerMessage
}

// for logging and for estimating, keep track of each item built
type itemBuilt struct {
	queueItemType QueueItemType
	designNum     int
	index         int
	numBuilt      int
	skipped       bool
	never         bool
}

// produce all items in the production queue
func (p *production) produce() productionResult {
	planet := p.planet

	productionResult := productionResult{}
	available := Cost{Resources: planet.Spec.ResourcesPerYearAvailable}.AddCargoMinerals(planet.Cargo)
	newQueue := []ProductionQueueItem{}
	for itemIndex := range planet.ProductionQueue {
		item := planet.ProductionQueue[itemIndex]
		maxBuildable := planet.maxBuildable(item.Type)
		if maxBuildable == Infinite {
			maxBuildable = math.MaxInt
		}

		// make sure this item is buildable, and if not, send the player a message and move on.
		if message, valid := p.validateItem(item, maxBuildable, planet); !valid {
			// skip this item and remove it from the queue
			productionResult.messages = append(productionResult.messages, message)
			productionResult.itemsBuilt = append(productionResult.itemsBuilt, itemBuilt{index: item.index, never: true})
			continue
		}

		// we can't build anymore of this auto item, skip it but leave it in the queue
		if maxBuildable == 0 && item.Type.IsAuto() {
			productionResult.itemsBuilt = append(productionResult.itemsBuilt, itemBuilt{index: item.index, skipped: true})
			newQueue = append(newQueue, item)
			continue
		}

		// build this item
		itemResult := p.processQueueItem(item, available, maxBuildable)

		// this one is an auto that will never be completed
		// leave it in the queue but move on to the next
		if itemResult.never {
			productionResult.itemsBuilt = append(productionResult.itemsBuilt, itemBuilt{index: item.index, never: itemResult.never})
			newQueue = append(newQueue, item)
			continue
		}

		// we built something, record ship buildings/packets for later, add installations and terraform right now
		if itemResult.numBuilt > 0 {
			// build the items on the planet and remove from our available
			p.addPlanetaryInstallations(item, itemResult.numBuilt)
			if item.Type.IsTerraform() {
				productionResult.terraformResults = append(productionResult.terraformResults, p.terraformPlanet(itemResult.numBuilt)...)
			}

			p.updateProductionResult(item, itemResult.numBuilt, &productionResult)

			// remove how much we spent from our available amount
			available = itemResult.leftover

			// if we built mineral alchemy, add it back in to our available amount
			available = available.Add(productionResult.alchemy.ToCost())

			productionResult.itemsBuilt = append(productionResult.itemsBuilt, itemBuilt{index: item.index, queueItemType: item.Type, designNum: item.DesignNum, numBuilt: itemResult.numBuilt})

			// planets are ending up with negative minerals. Trying to figure out why...
			if available.MinZero() != available {
				log.Warn().
					Int64("GameID", planet.GameID).
					Int64("ID", planet.ID).
					Str("Name", planet.Name).
					Str("Cargo", fmt.Sprintf("%+v", planet.Cargo)).
					Str("ProductionQueue", fmt.Sprintf("%+v", planet.ProductionQueue)).
					Str("itemResult", fmt.Sprintf("%+v", productionResult)).
					Msgf("available minerals and resources went negative - available: %+v", available)
				available = available.MinZero()
			}
		}

		// once we partially allocate, we're done
		allocatedToPartialItem := itemResult.allocated.Total() > 0
		if allocatedToPartialItem {
			if item.Type.IsAuto() && itemResult.leftover.Resources == 0 {
				// add the partially completed concrete item to the top of the queue
				newQueue = append([]ProductionQueueItem{
					{
						Type:      item.Type.concreteType(),
						Quantity:  1,
						Allocated: itemResult.allocated,
						CostOfOne: item.CostOfOne,
						index:     item.index, // for build estimates, keep track of the index of the auto item we're building
					}}, newQueue...)
				// add the rest of the items back to the queue and quit
				newQueue = append(newQueue, planet.ProductionQueue[itemIndex:]...)
				break
			} else {
				item.Allocated = itemResult.allocated
			}
		}

		if item.Type.IsAuto() {
			if available.Resources == 0 {
				// we are out of resources, so we are done building
				// append the unfinished queue back to the end of our remaining items
				newQueue = append(newQueue, planet.ProductionQueue[itemIndex:]...)
				break
			}

			// auto items stay in the queue
			newQueue = append(newQueue, item)
		} else {
			item.Quantity -= itemResult.numBuilt
			if item.Quantity > 0 {
				// keep it in the queue for next time
				newQueue = append(newQueue, item)
				if itemIndex < len(planet.ProductionQueue)-1 {
					// append the unfinished queue back to the end of our remaining items
					newQueue = append(newQueue, planet.ProductionQueue[itemIndex+1:]...)
				}
				break
			}
		}
	}

	// replace the queue with what's leftover
	planet.ProductionQueue = newQueue
	planet.Cargo = Cargo{available.Ironium, available.Boranium, available.Germanium, planet.Cargo.Colonists}
	if planet.Cargo.MinZero() != planet.Cargo {
		log.Warn().
			Int64("GameID", planet.GameID).
			Int64("ID", planet.ID).
			Str("Name", planet.Name).
			Str("Cargo", fmt.Sprintf("%+v", planet.Cargo)).
			Str("productionResult", fmt.Sprintf("%+v", productionResult)).
			Msgf("planet cargo was negative after production: %s", planet.Cargo.PrettyString())
		// planet.Cargo = planet.Cargo.MinZero()
	}

	// any leftover resources go back to the player for research
	productionResult.leftoverResources = available.Resources
	return productionResult
}

// for things that are built on the planet (mines, factories, etc) add them
func (p *production) addPlanetaryInstallations(item ProductionQueueItem, numBuilt int) {
	switch item.Type {
	case QueueItemTypeAutoMines:
		fallthrough
	case QueueItemTypeMine:
		p.planet.Mines += numBuilt
	case QueueItemTypeAutoFactories:
		fallthrough
	case QueueItemTypeFactory:
		p.planet.Factories += numBuilt
	case QueueItemTypeAutoDefenses:
		fallthrough
	case QueueItemTypeDefenses:
		p.planet.Defenses += numBuilt
	}
}

// terraform the planet and save the results for messages
func (p *production) terraformPlanet(numSteps int) []TerraformResult {
	planet, player := p.planet, p.player
	terraformer := NewTerraformer()
	terraformResults := make([]TerraformResult, numSteps)

	for i := 0; i < numSteps; i++ {
		// terraform one at a time to ensure the best things get terraformed
		terraformResults[i] = terraformer.TerraformOneStep(planet, player, nil, false)
	}

	return terraformResults
}

// validate an item in the production queue
func (p *production) validateItem(item ProductionQueueItem, maxBuildable int, planet *Planet) (PlayerMessage, bool) {
	if item.Type.IsPacket() && !planet.Spec.HasMassDriver {
		return newPlanetMessage(PlayerMessageBuildMineralPacketNoMassDriver, planet), false
	}
	if item.Type.IsPacket() && planet.PacketTargetNum == None {
		return newPlanetMessage(PlayerMessageBuildMineralPacketNoTarget, planet), false
	}
	if !item.Type.IsAuto() && maxBuildable == 0 {
		// can't build this, skip it
		// it shouldn't have been ever added to the queue, but just in case of a bug
		return newPlanetMessage(PlayerMessageBuildInvalidItem, planet).withSpec(PlayerMessageSpec{QueueItemType: item.Type}), false
	}

	return PlayerMessage{}, true
}

// the result of processing an item in the queue
type processQueueItemResult struct {
	never     bool
	numBuilt  int
	spent     Cost
	leftover  Cost
	allocated Cost
}

// build a single item in the queue, returning a result of what was built
func (p *production) processQueueItem(item ProductionQueueItem, availableToSpend Cost, maxBuildable int) processQueueItemResult {

	// skip auto items that will never complete, or that we don't need
	// this way we can put auto terraforming items up top
	// and skip them to build factories, then mines
	yearlyAvailableToSpend := FromMineralAndResources(p.planet.Spec.MiningOutput, p.planet.Spec.ResourcesPerYearAvailable)
	yearsToBuildOne := p.estimator.GetYearsToBuildOne(item, availableToSpend.ToMineral(), yearlyAvailableToSpend)
	if item.Type.IsAuto() && (yearsToBuildOne > 100 || yearsToBuildOne == Infinite) {
		return processQueueItemResult{never: true}
	}

	result := processQueueItemResult{}

	// add in anything allocated in previous turns
	availableToSpend = availableToSpend.Add(item.Allocated)
	item.Allocated = Cost{}

	// get the cost of the current item
	cost := item.CostOfOne

	if (cost != Cost{}) {
		// figure out how many we can build
		// and make sure we only build up to the quantity, and we don't build more than the planet supports
		numBuilt := maxInt(0, availableToSpend.NumBuildable(cost))
		numBuilt = minInt(numBuilt, item.Quantity)
		numBuilt = minInt(numBuilt, maxBuildable)

		if numBuilt > 0 {
			result.numBuilt = numBuilt
			result.spent = cost.MultiplyInt(numBuilt)
		}
		result.leftover = availableToSpend.Minus(result.spent)

		// If we didn't finish building all the items and we can still build more, allocate leftover resources to this item
		if numBuilt < item.Quantity && (maxBuildable == Infinite || numBuilt < maxBuildable) {
			result.allocated = p.allocatePartialBuild(cost, result.leftover)
			// don't allocate resources to an auto item if we're totally out minerals. We might never build it
			// and should move on to the next one
			// if item.Type.IsAuto() && item.CostOfOne.ToMineral().Total() > 0 && result.allocated.ToMineral().Total() == 0 {
			// 	result.allocated.Resources = 0
			// }
			result.leftover = result.leftover.Minus(result.allocated)
		}
	}

	return result
}

// add built items to planet, build fleets, update player messages, etc
func (p *production) updateProductionResult(item ProductionQueueItem, numBuilt int, result *productionResult) {
	switch item.Type {
	case QueueItemTypeAutoMineralAlchemy:
		fallthrough
	case QueueItemTypeMineralAlchemy:
		result.alchemy = Mineral{
			numBuilt,
			numBuilt,
			numBuilt,
		}
	case QueueItemTypeAutoMines:
		fallthrough
	case QueueItemTypeMine:
		result.mines += numBuilt
	case QueueItemTypeAutoFactories:
		fallthrough
	case QueueItemTypeFactory:
		result.factories += numBuilt
	case QueueItemTypeAutoDefenses:
		fallthrough
	case QueueItemTypeDefenses:
		result.defenses += numBuilt
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
		cargo := item.CostOfOne.MultiplyInt(numBuilt).ToCargo()
		result.packets = append(result.packets, cargo)
	case QueueItemTypeShipToken:
		result.tokens = append(result.tokens, ShipToken{Quantity: numBuilt, design: item.design, DesignNum: item.DesignNum})
	case QueueItemTypeStarbase:
		result.starbase = item.design
	case QueueItemTypePlanetaryScanner:
		result.scanner = true
	}
}

// Allocate resources to the top item on this production queue
// and return the leftover resources
//
// Costs are allocated by lowest percentage (except resources), i.e. if we require
// Cost(10, 10, 10, 100) and we only have Cost(1, 10, 10, 100)
// we allocate Cost(1, 1, 1, 100)
//
// The min amount we have is 10 percent of the ironium, so we
// apply 10 percent to each cost amount
func (p *production) allocatePartialBuild(costPerItem Cost, allocated Cost) Cost {
	ironiumPerc := 1.0
	if costPerItem.Ironium > 0 {
		ironiumPerc = math.Min(1, float64(allocated.Ironium)/float64(costPerItem.Ironium))
	}
	boraniumPerc := 1.0
	if costPerItem.Boranium > 0 {
		boraniumPerc = math.Min(1, float64(allocated.Boranium)/float64(costPerItem.Boranium))
	}
	germaniumPerc := 1.0
	if costPerItem.Germanium > 0 {
		germaniumPerc = math.Min(1, float64(allocated.Germanium)/float64(costPerItem.Germanium))
	}
	resourcesPerc := 1.0
	if costPerItem.Resources > 0 {
		resourcesPerc = math.Min(1, float64(allocated.Resources)/float64(costPerItem.Resources))
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
