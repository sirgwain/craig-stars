package cs

import (
	"fmt"
	"math"

	"github.com/rs/zerolog"
)

// The producer struct performs planetary production.
type producer struct {
	log       zerolog.Logger
	rules     *Rules
	planet    *Planet
	player    *Player
	estimator CompletionEstimator
}

// create a new planet production object
func newProducer(log zerolog.Logger, rules *Rules, planet *Planet, player *Player) producer {
	producerLogger := log.With().
		Int("Num", planet.Num).
		Str("Name", planet.Name).
		Int("PlayerNum", player.Num).
		Str("PlayerName", player.Race.PluralName).
		Logger()
	return producer{
		log:       producerLogger,
		rules:     rules,
		planet:    planet,
		player:    player,
		estimator: NewCompletionEstimator(),
	}
}

type ProductionQueueItem struct {
	QueueItemCompletionEstimate `tstype:",extends"`
	Type                        QueueItemType `json:"type"`
	DesignNum                   int           `json:"designNum,omitempty"`
	Quantity                    int           `json:"quantity"`
	Allocated                   Cost          `json:"allocated"`
	Tags                        Tags          `json:"tags"`
	index                       int           // used for holding a place in the queue while estimating
	design                      *ShipDesign
}

func (item *ProductionQueueItem) SetDesign(design *ShipDesign) {
	item.design = design
}

func NewProductionQueueItemShip(quantity int, design *ShipDesign) *ProductionQueueItem {
	return &ProductionQueueItem{Type: QueueItemTypeShipToken, Quantity: quantity, DesignNum: design.Num}
}

func (item *ProductionQueueItem) WithTag(key, value string) *ProductionQueueItem {
	if item.Tags == nil {
		item.Tags = make(Tags)
	}
	item.Tags[key] = value
	return item
}

func (item *ProductionQueueItem) GetTag(key string) string {
	return item.Tags[key]
}

type QueueItemCompletionEstimate struct {
	Canceled        bool `json:"canceled,omitempty"`
	YearsToBuildOne int  `json:"yearsToBuildOne,omitempty"`
	YearsToBuildAll int  `json:"yearsToBuildAll,omitempty"`
	YearsToSkipAuto int  `json:"yearsToSkipAuto,omitempty"`
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
	QueueItemTypeGenesisDevice          QueueItemType = "GenesisDevice"
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
	tokens            []builtShip
	packets           Cargo
	scanner           bool
	reset             bool
	starbase          *ShipDesign
	alchemy           int
	mines             int
	factories         int
	defenses          int
	terraformResults  []TerraformResult
	messages          []PlayerMessage
	completed         bool
}

// A record of a built queue item, used for logging & estimating
type itemBuilt struct {
	queueItemType QueueItemType
	designNum     int
	index         int
	numBuilt      int
	skipped       bool // whether an auto item is skipped
	canceled      bool // whether an invalid concrete item is canceled
}

type builtShip struct {
	ShipToken
	tags Tags
}

// produce all items in the production queue
func (p *producer) produce() (result productionResult, err error) {
	planet := p.planet

	if len(planet.ProductionQueue) == 0 {
		// no queue, no dice
		fmt.Println("BABA IS YOU \n\nAAAAAAAAAAAAAAAAAAAAAAAA")
		result.messages = append(result.messages,
			newPlanetMessage(PlayerMessagePlanetProductionQueueComplete, planet))
		return productionResult{leftoverResources: planet.Spec.ResourcesPerYearAvailable, completed: true}, nil
	}
	// TODO: Fix auto mineral alchemy:
	// * Only has 1 max quantity
	// * Always blocks queue if final item;
	// * otherwise only blocks if subsequent item is out of minerals

	var (
		// tracker of resources/minerals available to spend
		available = Cost{Resources: planet.Spec.ResourcesPerYearAvailable}.AddMineral(planet.Cargo.ToMineral())
		// new queue of production items, used to recreate an updated queue
		newQueue         = []ProductionQueueItem{}
		hasBuiltAnything bool // if we've built anything yet
	)

	// check each item in the queue in order
	for itemIndex, item := range planet.ProductionQueue {
		itemCost, err := p.getItemCost(p.rules, p.player, p.planet, item)
		if err != nil {
			p.log.Error().
				Err(err).
				Any("item", item).
				Msgf("produce() obtained error when calculating item costs: %v", err)
			return productionResult{}, err
		}

		maxBuildable := planet.maxBuildable(p.player, item.Type)
		// Infinite is the constant int of -1, but for our purposes we want a very large number
		if maxBuildable == Infinite {
			maxBuildable = math.MaxInt
		}

		// If we haven't built anything yet and this is a concrete item,
		// clamp its build quantity down to maxBuildable.
		// We do this for everything ahead of us after each successful build;
		// this just ensures we don't forget to check before that happens.
		if !hasBuiltAnything && !item.Type.IsAuto() {
			oldQty := item.Quantity
			if item = p.clampItemQty(item, maxBuildable); item.Quantity <= 0 {
				// can't build any more of this item; mark as canceled & move on
				p.handleInvalidQty(item, &available, &result, &itemIndex)
				p.log.Debug().
					Any("Item", item).
					Int("PrevQty", oldQty).
					Int("ClampedQty", item.Quantity).
					Int("maxBuildable", maxBuildable).
					Msgf("cancelling queue item; can't build any more")
				continue
			}
		}

		// Add in any previously allocated resources for this item into our pot.
		available = available.Add(item.Allocated)
		item.Allocated = Cost{}

		// check for auto items we should skip due to not being buildable or lacking minerals
		if item.Type.IsAuto() && (maxBuildable <= 0 || available.DivideMineral(itemCost.ToMineral()) < 1) {
			result.itemsBuilt = append(result.itemsBuilt, itemBuilt{index: item.index, skipped: true})
			newQueue = append(newQueue, item)

			// if we skipped the last item in the queue, mark result as done
			if itemIndex == len(planet.ProductionQueue)-1 {
				result.completed = true
			}
			continue
		}

		// make sure this item is buildable, notifying the player if not.
		if message, valid := p.validateItem(item, planet); !valid {
			// cancel this item and remove it from the queue
			result.messages = append(result.messages, message)
			result.itemsBuilt = append(result.itemsBuilt, itemBuilt{index: item.index, canceled: true})
			continue
		}

		// After all that, we can finally try to build the dang thing.

		// determine how much to build
		numBuilt, spent := p.getNumBuilt(item, itemCost, available, maxBuildable)

		// if we can't build anything, skip to "cleanup" section
		if numBuilt <= 0 {
			goto checkDone
		}

		hasBuiltAnything = true
		available = available.Subtract(spent)
		// record info and add installations/terraforming
		p.addPlanetaryInstallations(item, numBuilt)

		if item.Type.IsTerraform() {
			result.terraformResults = append(result.terraformResults,
				p.terraformPlanet(numBuilt)...)
		}

		p.updateProductionResult(item, numBuilt, itemCost, &result)

		// if we built mineral alchemy, add it back in to our pot to use later
		available = available.AddToAllMineral(result.alchemy)

		result.itemsBuilt = append(result.itemsBuilt, itemBuilt{index: item.index, queueItemType: item.Type, designNum: item.DesignNum, numBuilt: numBuilt})

		// planets are ending up with negative minerals. Trying to figure out why...
		if available.MinZero() != available {
			p.log.Warn().
				Str("Cargo", fmt.Sprintf("%+v", planet.Cargo)).
				Str("ProductionQueue", fmt.Sprintf("%+v", planet.ProductionQueue)).
				Str("itemResult", fmt.Sprintf("%+v", result)).
				Msgf("planet minerals/resources went negative - available: %+v", available)
			available = available.MinZero()
		}

		// Elide any remaining concrete items within the queue that are over cap.
		copy(planet.ProductionQueue[itemIndex+1:], MapSlice(planet.ProductionQueue[itemIndex+1:], func(i ProductionQueueItem) (newItem ProductionQueueItem, keep bool) {
			if i.Type.IsAuto() {
				// leave auto items alone
				return i, true
			}

			oldQty := i.Quantity
			maxBuildable := planet.maxBuildable(p.player, i.Type)

			// clamp item quantity down to maxBuildable
			if i = p.clampItemQty(i, maxBuildable); i.Quantity <= 0 {
				// item hit 0 quantity; mark as canceled
				p.handleInvalidQty(i, &itemCost, &result, &itemIndex)
				p.log.Debug().
					Any("Item", i).
					Int("PrevQty", oldQty).
					Int("ClampedQty", i.Quantity).
					Int("maxBuildable", maxBuildable).
					Msgf("cancelling queue item; can't build any more")
			}

			// return the modified item with clamped qty
			return i, i.Quantity > 0
		}))

	checkDone:
		if itemIndex == len(planet.ProductionQueue)-1 && (numBuilt >= item.Quantity || numBuilt >= maxBuildable) {
			// we built the last item in the queue; we're all done
			result.completed = true
			if item.Type.IsAuto() {
				// append the unfinished queue back to the end of our remaining items
				newQueue = append(newQueue, planet.ProductionQueue[itemIndex:]...)
			}
			break
		}

		// TODO: Refactor this to make the control flow somewhat more obvious
		if item.Type.IsAuto() {
			// auto items stay in the queue after being built
			newQueue = append(newQueue, item)

			// if we have resources left, try and move on to the next item
			// (auto items don't block the queue)
			if available.Resources > 0 {
				if numBuilt >= item.Quantity || numBuilt >= maxBuildable || // We've built all that we can
					available.DivideMineral(itemCost.ToMineral()) < 1 {
					// We've built all that we can for this auto item; move on
					continue
				}

				// we still have auto items left to build and enough minerals
				// to complete one auto item; add a concrete one to the top of the queue
				newQueue = append([]ProductionQueueItem{
					{
						Type:      item.Type.concreteType(),
						Quantity:  1,
						Allocated: p.allocatePartialBuild(itemCost, available),
						index:     -1, // we don't track concrete auto items, we only care about the first fully built auto item
					}}, newQueue...)
				available = available.Subtract(newQueue[0].Allocated)

				if itemIndex < len(planet.ProductionQueue)-1 {
					// if this isn't the last item, tack the rest of the queue back on
					newQueue = append(newQueue, planet.ProductionQueue[itemIndex+1:]...)
				}
				break
			} else {
				// all resources spent; wrap up
				if itemIndex < len(planet.ProductionQueue)-1 {
					// if this isn't the last item, tack the rest of the queue back on
					newQueue = append(newQueue, planet.ProductionQueue[itemIndex+1:]...)
				}
				break
			}
		} else {
			// dock amount built from the concrete item's remaining quantity
			item.Quantity -= numBuilt
			if item.Quantity <= 0 {
				// Finished concrete build; move on
				continue
			}

			// could not finish concrete item; done building
			// allocate remaining resources to partially built item
			item.Allocated = p.allocatePartialBuild(itemCost, available)
			available = available.Subtract(item.Allocated)
			planet.ProductionQueue[itemIndex] = item

			// keep it in the queue and break out
			newQueue = append(newQueue, planet.ProductionQueue[itemIndex:]...)
			break
		}
	}

	// ping player about an empty queue
	if result.completed {
		result.messages = append(result.messages, newPlanetMessage(PlayerMessagePlanetProductionQueueComplete, planet))
	}

	// replace queue & cargo with leftovers post-production
	planet.ProductionQueue = newQueue
	planet.Cargo.SetMineral(available.ToMineral())
	if planet.Cargo.MinZero() != planet.Cargo {
		p.log.Warn().
			Str("Cargo", fmt.Sprintf("%+v", planet.Cargo)).
			Str("productionResult", fmt.Sprintf("%+v", result)).
			Msgf("planet cargo was negative after production: %s", planet.Cargo.PrettyString())
		// planet.Cargo = planet.Cargo.MinZero()
		return result, fmt.Errorf("planet cargo was negative after production")
	}

	// leftover resources go back to the player for research
	result.leftoverResources = available.Resources
	return result, nil
}

func (p *producer) getItemCost(rules *Rules, player *Player, planet *Planet, item ProductionQueueItem) (cost Cost, err error) {
	costCalculator := NewCostCalculator()
	switch {
	case item.Type == QueueItemTypeStarbase && planet.Spec.HasStarbase:
		// upgrade existing starbase
		cost, err = costCalculator.StarbaseUpgradeCost(rules, player.TechLevels, player.Race.Spec, planet.Starbase.Tokens[0].design, item.design)
		if err != nil {
			return Cost{}, fmt.Errorf("failed to compute starbase upgrade cost: %w", err)
		}
	case item.Type == QueueItemTypeStarbase || item.Type == QueueItemTypeShipToken:
		// make new ship/starbase
		cost, err = costCalculator.GetDesignCost(rules, player.TechLevels, player.Race.Spec, item.design)
		if err != nil {
			return Cost{}, fmt.Errorf("failed to get design cost: %w", err)
		}
	default:
		// regular old item
		cost, err = costCalculator.CostOfOne(player, item)
		if err != nil {
			return Cost{}, fmt.Errorf("failed to compute cost of %s: %w", item.Type, err)
		}
	}
	return cost, nil
}

// Clamp a ProductionQueueItem's quantity down to however much we can actually build,
// returning the modified queue item.
func (p *producer) clampItemQty(item ProductionQueueItem, maxBuildable int) ProductionQueueItem {
	if maxBuildable != Infinite && item.Quantity > maxBuildable {
		item.Quantity = maxBuildable
	}
	return item
}

// Perform necessary cleanup for concrete items with invalid quantities.
func (p *producer) handleInvalidQty(item ProductionQueueItem, costTally *Cost, result *productionResult, itemIndex *int) {
	// Since we only call this on the current queue item pre-production
	// or things ahead of us in the queue, none of them will have any
	// pre-existing records in itemsBuilt (meaning we're A-OK to tack em on).
	result.itemsBuilt = append(result.itemsBuilt, itemBuilt{index: item.index, canceled: true})
	result.messages = append(result.messages,
		newPlanetMessage(PlayerMessagePlanetBuiltInvalidItem, p.planet).
			withSpec(PlayerMessageSpec{Name: p.planet.Name, QueueItemType: item.Type}))

	// add any allocated resources back to our tally
	// (they should never have been spent to begin with)
	*costTally = costTally.Add(item.Allocated) // no need to clear item allocated since we skip it afterwards
	*itemIndex--                               // decrement itemIndex so we keep the loop in sync
}

// validate an item in the production queue
func (p *producer) validateItem(item ProductionQueueItem, planet *Planet) (PlayerMessage, bool) {
	if item.Type.IsPacket() && !planet.Spec.HasMassDriver {
		return newPlanetMessage(PlayerMessagePlanetBuiltInvalidMineralPacketNoMassDriver, planet), false
	}
	if item.Type.IsPacket() && planet.PacketTargetNum == None {
		return newPlanetMessage(PlayerMessagePlanetBuiltInvalidMineralPacketNoTarget, planet), false
	}

	return PlayerMessage{}, true
}

// getNumBuilt returns how many of a given production queue item we can build
// and how how much to spend on it.
func (p *producer) getNumBuilt(item ProductionQueueItem, cost, availableToSpend Cost, maxBuildable int) (numBuilt int, spent Cost) {
	if cost == (Cost{}) {
		return min(item.Quantity, maxBuildable), Cost{}
	}

	// figure out how many we can build;
	// make sure we only build up to the quantity required
	// and we don't build more than the planet supports
	numBuilt = max(0, min(item.Quantity, maxBuildable,
		int(availableToSpend.DivideCost(cost))))
	spent = MultiplyCost(cost, numBuilt)

	return numBuilt, spent
}

// add any planetary installations built during production to the planet building them.
func (p *producer) addPlanetaryInstallations(item ProductionQueueItem, numBuilt int) {
	switch item.Type {
	case QueueItemTypeAutoMines, QueueItemTypeMine:
		p.planet.Mines += numBuilt
	case QueueItemTypeAutoFactories, QueueItemTypeFactory:
		p.planet.Factories += numBuilt
	case QueueItemTypeAutoDefenses, QueueItemTypeDefenses:
		p.planet.Defenses += numBuilt
	case QueueItemTypePlanetaryScanner:
		p.planet.Scanner = true
	}
}

// terraform a planet during production and save the results for messages
func (p *producer) terraformPlanet(numSteps int) []TerraformResult {
	planet, player := p.planet, p.player
	terraformer := NewTerraformer()
	terraformResults := make([]TerraformResult, numSteps)

	for i := range numSteps {
		// terraform one step at a time to ensure the best things get terraformed
		terraformResults[i] = terraformer.TerraformOneStep(planet, player, nil, false)
	}

	return terraformResults
}

// updateProductionResult updates a production result with information about a newly built item.
func (p *producer) updateProductionResult(item ProductionQueueItem, numBuilt int, cost Cost, result *productionResult) {
	switch item.Type {
	case QueueItemTypeAutoMineralAlchemy, QueueItemTypeMineralAlchemy:
		result.alchemy += numBuilt
	case QueueItemTypeAutoMines, QueueItemTypeMine:
		result.mines += numBuilt
	case QueueItemTypeAutoFactories, QueueItemTypeFactory:
		result.factories += numBuilt
	case QueueItemTypeAutoDefenses, QueueItemTypeDefenses:
		result.defenses += numBuilt
	case QueueItemTypeAutoMineralPacket, QueueItemTypeMixedMineralPacket, QueueItemTypeIroniumMineralPacket, QueueItemTypeBoraniumMineralPacket, QueueItemTypeGermaniumMineralPacket:
		// add this packet cargo to the production result
		// so it can be added as packets to the universe later
		// Multiply by 1/PacketMineralCostFactor to simulate overhead (pay 132 kT; only get 120)
		cargo := MultiplyCost(MultiplyCost(cost, 1/p.player.Race.Spec.PacketMineralCostFactor), numBuilt).ToCargo()
		result.packets = result.packets.Add(cargo)
	case QueueItemTypeShipToken:
		result.tokens = append(result.tokens, builtShip{ShipToken: ShipToken{Quantity: numBuilt, design: item.design, DesignNum: item.DesignNum}, tags: item.Tags})
	case QueueItemTypeStarbase:
		result.starbase = item.design
	case QueueItemTypePlanetaryScanner:
		result.scanner = true
	case QueueItemTypeGenesisDevice:
		result.reset = true
	}
}

// allocatePartialBuild returns the amount of minerals and resources to allocate
// to a partially built production queue item, based on its cost and available
// minerals/resources.
func (p *producer) allocatePartialBuild(costPerItem Cost, available Cost) (allocated Cost) {
	// Costs are allocated by lowest percentage (except resources), i.e. if we require
	// Cost(10, 10, 10, 100) and we only have Cost(1, 10, 10, 100)
	// we allocate Cost(1, 1, 1, 100).
	// The min amount we have is 10% for ironium, so we
	// apply 10% to each cost amount
	ironiumPerc := 1.0
	if costPerItem.Ironium > 0 {
		ironiumPerc = min(1, float64(available.Ironium)/float64(costPerItem.Ironium))
	}
	boraniumPerc := 1.0
	if costPerItem.Boranium > 0 {
		boraniumPerc = min(1, float64(available.Boranium)/float64(costPerItem.Boranium))
	}
	germaniumPerc := 1.0
	if costPerItem.Germanium > 0 {
		germaniumPerc = min(1, float64(available.Germanium)/float64(costPerItem.Germanium))
	}
	resourcesPerc := 1.0
	if costPerItem.Resources > 0 {
		resourcesPerc = min(1, float64(available.Resources)/float64(costPerItem.Resources))
	}

	// figure out the lowest percentage
	minPerc := min(ironiumPerc, boraniumPerc, germaniumPerc, resourcesPerc)

	// allocate the lowest percentage of each cost, rounded down
	allocated = Cost{
		int(float64(costPerItem.Ironium) * minPerc),
		int(float64(costPerItem.Boranium) * minPerc),
		int(float64(costPerItem.Germanium) * minPerc),
		int(float64(costPerItem.Resources) * minPerc),
	}

	return allocated
}
