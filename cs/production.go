package cs

import (
	"fmt"
	"log/slog"
	"math"
	"slices"
)

// The producer interface performs planetary production

// create a new planet production object
func newProducer(log *slog.Logger, rules *Rules, planet *Planet, player *Player) producer {
	producerLogger := log.With(
		slog.Int("Num", planet.Num),
		slog.String("Name", planet.Name),
		slog.Int("PlayerNum", player.Num),
		slog.String("Player", player.Name),
	)
	return producer{
		log:       producerLogger,
		rules:     rules,
		planet:    planet,
		player:    player,
		estimator: NewCompletionEstimator(),
	}
}

type producer struct {
	log       *slog.Logger
	rules     *Rules
	planet    *Planet
	player    *Player
	estimator CompletionEstimator
}

type QueueItemCompletionEstimate struct {
	Skipped         bool `json:"skipped,omitempty"`
	YearsToBuildOne int  `json:"yearsToBuildOne,omitempty"`
	YearsToBuildAll int  `json:"yearsToBuildAll,omitempty"`
	YearsToSkipAuto int  `json:"yearsToSkipAuto,omitempty"`
}

type ProductionQueueItem struct {
	QueueItemCompletionEstimate
	Type      QueueItemType `json:"type"`
	DesignNum int           `json:"designNum,omitempty"`
	Quantity  int           `json:"quantity"`
	Allocated Cost          `json:"allocated,omitzero"`
	Tags      Tags          `json:"tags,omitzero"`
	index     int           // used for holding a place in the queue while estimating
	design    *ShipDesign
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

type QueueItemType string

const (
	QueueItemTypeUnspecified            QueueItemType = ""
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

type builtShip struct {
	ShipToken
	tags Tags
}

// productionStatus is the outcome of building one queue item (mdProdStat in the original)
type productionStatus int

const (
	// a concrete item finished, remove it from the queue
	productionStatusComplete productionStatus = iota
	// an auto item built its quota for the year
	productionStatusCompleteAuto
	// an auto item had nothing to build (i.e. the planet is at capacity)
	productionStatusSkippedAuto
	// an auto item is waiting on minerals, move on to the next item
	productionStatusMineralBlockedAuto
	// out of resources (or minerals for a concrete item), stop production for the year
	productionStatusBlocked
)

// the result of building a single queue item
type buildItemResult struct {
	status    productionStatus
	numBuilt  int
	alchemy   int                  // kT of each mineral created by alchemy
	allocated Cost                 // resources/minerals spent toward the next unfinished item
	partial   *ProductionQueueItem // a partially built item to add to the front of the queue
}

// produce all items in the production queue
//
// This follows Produce() in the original. Each item is built in order until we run out of resources.
// Concrete items are removed from the queue when complete, auto items stay in the queue.
// Auto alchemy in front of another item converts resources to minerals as needed to build that item.
// Auto alchemy at the end of the queue converts all leftover resources into minerals.
func (p *producer) produce() (productionResult, error) {
	planet := p.planet
	result := productionResult{}
	available := Cost{Resources: planet.Spec.ResourcesPerYearAvailable}.AddMineral(planet.Cargo.ToMineral())

	// work on a copy of the queue, removing items as they complete
	queue := slices.Clone(planet.ProductionQueue)
	alchemy := false // true if the previous item is auto alchemy
	for i := 0; i < len(queue); {
		item := queue[i]
		if item.Quantity <= 0 {
			queue, i = removeQueueItem(queue, i, alchemy)
			alchemy = false
			continue
		}

		maxBuildable := p.maxBuildable(item.Type)

		// make sure this item is buildable, and if not, send the player a message and remove it
		if message, valid := p.validateItem(item, maxBuildable, planet); !valid {
			result.messages = append(result.messages, message)
			result.itemsBuilt = append(result.itemsBuilt, itemBuilt{index: item.index, never: true})
			queue, i = removeQueueItem(queue, i, alchemy)
			alchemy = false
			continue
		}

		// concrete orders beyond the planet's capacity are reduced to the max
		if !item.Type.IsAuto() && item.Quantity > maxBuildable {
			result.messages = append(result.messages, newPlanetMessage(PlayerMessagePlanetBuiltBeyondMaximum, planet).
				withSpec(PlayerMessageSpec{QueueItemType: item.Type, Amount: maxBuildable}))
			item.Quantity = maxBuildable
			queue[i] = item
		}

		// auto alchemy in front of another item makes minerals for that item
		if item.Type == QueueItemTypeAutoMineralAlchemy && i < len(queue)-1 {
			alchemy = true
			i++
			continue
		}

		cost, err := p.getItemCost(p.rules, p.player, planet, item)
		if err != nil {
			p.log.Error("produce() returned error when calculating costs",
				slog.Any("err", err),
				slog.Any("item", item),
			)
			return productionResult{}, err
		}

		quantity := min(item.Quantity, maxBuildable)
		if item.Type == QueueItemTypeAutoMineralAlchemy {
			// auto alchemy at the end of the queue converts all our leftover resources
			quantity = maxBuildable
		}

		built, err := p.buildItem(item, cost, quantity, &available, alchemy)
		if err != nil {
			return productionResult{}, err
		}

		result.alchemy += built.alchemy
		if built.numBuilt > 0 {
			// add mines and factories to the planet, terraform
			p.addPlanetaryInstallations(item, built.numBuilt)
			if item.Type.IsTerraform() {
				result.terraformResults = append(result.terraformResults, p.terraformPlanet(built.numBuilt)...)
				p.refreshPlanetSpec()
			}

			// record ship buildings/packets for later
			p.updateProductionResult(item, built.numBuilt, cost, &result)
			result.itemsBuilt = append(result.itemsBuilt, itemBuilt{index: item.index, queueItemType: item.Type, designNum: item.DesignNum, numBuilt: built.numBuilt})
		} else if built.status == productionStatusSkippedAuto || built.status == productionStatusMineralBlockedAuto {
			result.itemsBuilt = append(result.itemsBuilt, itemBuilt{index: item.index, skipped: true})
		}

		if built.status == productionStatusBlocked {
			// save our progress and stop production for the year
			if !item.Type.IsAuto() {
				queue[i].Quantity -= built.numBuilt
				queue[i].Allocated = built.allocated
			}
			if built.partial != nil {
				queue = slices.Insert(queue, 0, *built.partial)
			}
			break
		}

		if built.status == productionStatusComplete {
			queue, i = removeQueueItem(queue, i, alchemy)
		} else {
			// auto items (and any auto alchemy in front of them) stay in the queue
			i++
		}
		alchemy = false
	}

	planet.ProductionQueue = queue
	planet.Cargo = Cargo{available.Ironium, available.Boranium, available.Germanium, planet.Cargo.Colonists}
	if planet.Cargo.MinZero() != planet.Cargo {
		p.log.Warn("planet cargo was negative after production",
			slog.String("Cargo", planet.Cargo.PrettyString()),
			slog.Any("productionResult", result),
		)
		return result, fmt.Errorf("planet cargo was negative after production")
	}

	// any leftover resources go back to the player for research
	result.leftoverResources = available.Resources
	return result, nil
}

// refresh the planet spec after terraforming so later items use the new terraform amount,
// installation caps, and growth. This year's resources were already budgeted, so they stay the same.
func (p *producer) refreshPlanetSpec() {
	previous := p.planet.Spec
	p.planet.Spec = ComputePlanetSpec(p.rules, p.player, p.planet)
	p.planet.Spec.ResourcesPerYear = previous.ResourcesPerYear
	p.planet.Spec.ResourcesPerYearAvailable = previous.ResourcesPerYearAvailable
	p.planet.Spec.ResourcesPerYearResearch = previous.ResourcesPerYearResearch
	p.planet.Spec.ResourcesPerYearResearchEstimatedLeftover = previous.ResourcesPerYearResearchEstimatedLeftover
}

// remove the item at index i, along with the auto alchemy in front of it
// returns the new queue and the index of the next item
func removeQueueItem(queue []ProductionQueueItem, i int, alchemy bool) ([]ProductionQueueItem, int) {
	start := i
	if alchemy {
		start--
	}
	return slices.Delete(queue, start, i+1), start
}

// the most we can build of an item this year
func (p *producer) maxBuildable(itemType QueueItemType) int {
	if itemType == QueueItemTypeAutoMineralPacket && (!p.planet.Spec.HasMassDriver || p.planet.PacketTargetNum == None) {
		// auto packets wait for a mass driver and a target
		return 0
	}
	maxBuildable := p.planet.MaxBuildable(p.player, itemType)
	// Infinite is the constant int of -1, but for our purposes we want a very large number
	if maxBuildable == Infinite {
		return math.MaxInt
	}
	return maxBuildable
}

// buildItem builds up to quantity of an item, spending from available.
// If alchemy is true, resources are converted to minerals whenever minerals block progress.
//
// This follows CBuildProdItem() in the original
func (p *producer) buildItem(item ProductionQueueItem, cost Cost, quantity int, available *Cost, alchemy bool) (buildItemResult, error) {
	result := buildItemResult{}
	auto := item.Type.IsAuto()

	if cost == (Cost{}) {
		result.numBuilt = quantity
		quantity = 0
	}

	var alchemyCost int
	if alchemy {
		cost, err := p.getItemCost(p.rules, p.player, p.planet, ProductionQueueItem{Type: QueueItemTypeMineralAlchemy})
		if err != nil {
			return result, err
		}
		alchemyCost = cost.Resources
	}

	// auto items don't carry allocations, their partial builds are concrete items
	paid := item.Allocated
	remaining := quantity
	for remaining > 0 {
		// build as many as we can afford, including what we've already paid toward the first one
		funds := available.Add(paid)
		if numBuilt := min(remaining, int(funds.DivideCost(cost))); numBuilt > 0 {
			*available = funds.Subtract(MultiplyCost(cost, numBuilt))
			paid = Cost{}
			result.numBuilt += numBuilt
			remaining -= numBuilt
			continue
		}

		limit, mineralBlocked := limitingCost(cost, funds)
		if auto && funds.DivideMineral(cost.ToMineral()) < 1 {
			// auto items short on any mineral wait for minerals without spending anything,
			// unless they have alchemy, then they go straight to converting
			if !alchemy {
				result.status = productionStatusMineralBlockedAuto
				return result, nil
			}
		} else {
			// put what we have toward the next one
			paid = p.allocatePartialBuild(cost, funds)
			*available = funds.Subtract(paid)
			if !mineralBlocked || !alchemy {
				break
			}
		}

		// convert resources into what we're shortest on, alchemy makes some of every mineral
		// Like the original, an auto item short on resources ends up here too, and its
		// leftover resources go toward alchemy next year
		needed := cost.GetAmount(limit) - paid.GetAmount(limit) - available.GetAmount(limit)
		converted := min(needed, available.Resources/alchemyCost)
		available.Resources -= converted * alchemyCost
		*available = available.AddToAllMineral(converted)
		result.alchemy += converted
		if converted < needed {
			// out of resources, put the rest toward alchemy next year
			if available.Resources > 0 {
				result.partial = &ProductionQueueItem{Type: QueueItemTypeMineralAlchemy, Quantity: 1, Allocated: Cost{Resources: available.Resources}, index: -1}
				available.Resources = 0
			}
			break
		}
	}

	// alchemy items make minerals
	if item.Type == QueueItemTypeMineralAlchemy || item.Type == QueueItemTypeAutoMineralAlchemy {
		result.alchemy += result.numBuilt
		*available = available.AddToAllMineral(result.numBuilt)
	}

	switch {
	case remaining > 0:
		result.status = productionStatusBlocked
		result.allocated = paid
		if auto && paid != (Cost{}) && result.partial == nil {
			result.partial = &ProductionQueueItem{
				Type:      item.Type.concreteType(),
				Quantity:  1,
				Allocated: paid,
				index:     -1, // we don't track concrete auto items, we only care about the first fully built auto item
			}
		}
	case !auto:
		result.status = productionStatusComplete
	case result.numBuilt == 0:
		result.status = productionStatusSkippedAuto
	default:
		result.status = productionStatusCompleteAuto
	}
	return result, nil
}

// limitingCost returns the cost type we're shortest on to build one more of an item,
// and true if it's a mineral. Minerals win ties.
func limitingCost(cost, funds Cost) (CostType, bool) {
	lowest := math.Inf(1)
	limit := Resources
	for _, costType := range CostTypes {
		if amount := cost.GetAmount(costType); amount > 0 {
			if ratio := float64(funds.GetAmount(costType)) / float64(amount); ratio < lowest {
				lowest, limit = ratio, costType
			}
		}
	}
	return limit, limit != Resources
}

func (p *producer) getItemCost(rules *Rules, player *Player, planet *Planet, item ProductionQueueItem) (Cost, error) {
	costCalculator := NewCostCalculator()
	var err error
	var cost Cost
	if item.Type == QueueItemTypeStarbase && planet.Spec.HasStarbase {
		cost, err = costCalculator.StarbaseUpgradeCost(rules, player.TechLevels, player.Race.Spec, planet.Starbase.Tokens[0].design, item.design)
		if err != nil {
			return Cost{}, fmt.Errorf("failed to compute starbase upgrade cost: %w", err)
		}
	} else if item.Type == QueueItemTypeStarbase || item.Type == QueueItemTypeShipToken {
		cost, err = costCalculator.GetDesignCost(rules, player.TechLevels, player.Race.Spec, item.design)
		if err != nil {
			return Cost{}, fmt.Errorf("failed to get design cost: %w", err)
		}
	} else {
		cost, err = costCalculator.CostOfOne(player, item)
		if err != nil {
			return Cost{}, fmt.Errorf("failed to compute cost of %s: %w", item.Type, err)
		}
	}
	return cost, nil
}

// Allocate minerals and resources to the top item on this production queue
// and return the leftover resources
//
// Costs are allocated by lowest percentage (except resources), i.e. if we require
// Cost(10, 10, 10, 100) and we only have Cost(1, 10, 10, 100)
// we allocate Cost(1, 1, 1, 100)
//
// The min amount we have is 10 percent of the ironium, so we
// apply 10 percent to each cost amount
func (p *producer) allocatePartialBuild(costPerItem Cost, allocated Cost) Cost {
	ironiumPerc := 1.0
	if costPerItem.Ironium > 0 {
		ironiumPerc = min(1, float64(allocated.Ironium)/float64(costPerItem.Ironium))
	}
	boraniumPerc := 1.0
	if costPerItem.Boranium > 0 {
		boraniumPerc = min(1, float64(allocated.Boranium)/float64(costPerItem.Boranium))
	}
	germaniumPerc := 1.0
	if costPerItem.Germanium > 0 {
		germaniumPerc = min(1, float64(allocated.Germanium)/float64(costPerItem.Germanium))
	}
	resourcesPerc := 1.0
	if costPerItem.Resources > 0 {
		resourcesPerc = min(1, float64(allocated.Resources)/float64(costPerItem.Resources))
	}

	// figure out the lowest percentage
	minPerc := min(ironiumPerc, boraniumPerc, germaniumPerc, resourcesPerc)

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

// for things that are built on the planet (mines, factories, etc) add them
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

// terraform the planet and save the results for messages
func (p *producer) terraformPlanet(numSteps int) []TerraformResult {
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
func (p *producer) validateItem(item ProductionQueueItem, maxBuildable int, planet *Planet) (PlayerMessage, bool) {
	if item.Type.IsAuto() {
		// auto items are skipped rather than removed when they can't be built
		return PlayerMessage{}, true
	}
	if item.Type.IsPacket() && !planet.Spec.HasMassDriver {
		return newPlanetMessage(PlayerMessagePlanetBuiltInvalidMineralPacketNoMassDriver, planet), false
	}
	if item.Type.IsPacket() && planet.PacketTargetNum == None {
		return newPlanetMessage(PlayerMessagePlanetBuiltInvalidMineralPacketNoTarget, planet), false
	}
	if item.Type == QueueItemTypeShipToken && (item.design == nil || !planet.CanBuild(item.design.Spec.Mass)) {
		// the starbase was destroyed or downgraded since this was queued
		spec := PlayerMessageSpec{Name: planet.Name, QueueItemType: item.Type, Amount2: planet.Spec.DockCapacity}
		if item.design != nil {
			spec.Name, spec.Amount = item.design.Name, item.design.Spec.Mass
		}
		return newPlanetMessage(PlayerMessagePlanetBuiltInvalidShip, planet).withSpec(spec), false
	}
	if maxBuildable == 0 && (item.Type == QueueItemTypeMine || item.Type == QueueItemTypeFactory ||
		item.Type == QueueItemTypeDefenses || item.Type == QueueItemTypeTerraformEnvironment) {
		// no room left for these
		return newPlanetMessage(PlayerMessagePlanetBuiltBeyondMaximum, planet).withSpec(PlayerMessageSpec{QueueItemType: item.Type}), false
	}
	if maxBuildable == 0 {
		// can't build this, skip it
		// it shouldn't have been ever added to the queue, but just in case of a bug
		return newPlanetMessage(PlayerMessagePlanetBuiltInvalidItem, planet).withSpec(PlayerMessageSpec{Name: planet.Name, QueueItemType: item.Type}), false
	}

	return PlayerMessage{}, true
}

// add built items to planet, build fleets, update player messages, etc
func (p *producer) updateProductionResult(item ProductionQueueItem, numBuilt int, cost Cost, result *productionResult) {
	switch item.Type {
	case QueueItemTypeAutoMines, QueueItemTypeMine:
		result.mines += numBuilt
	case QueueItemTypeAutoFactories, QueueItemTypeFactory:
		result.factories += numBuilt
	case QueueItemTypeAutoDefenses, QueueItemTypeDefenses:
		result.defenses += numBuilt
	case QueueItemTypeAutoMineralPacket, QueueItemTypeMixedMineralPacket, QueueItemTypeIroniumMineralPacket, QueueItemTypeBoraniumMineralPacket, QueueItemTypeGermaniumMineralPacket:
		// add this packet cargo to the production result
		// so it can be added as packets to the universe later
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
