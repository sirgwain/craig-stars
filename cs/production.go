package cs

import (
	"fmt"
	"math"

	"github.com/rs/zerolog"
)

// The producer interface performs planetary production
type producer interface {
	produce() (productionResult, error)
}

// create a new planet production object
func newProducer(log zerolog.Logger, rules *Rules, planet *Planet, player *Player) producer {
	producerLogger := log.With().
		Int("Num", planet.Num).
		Str("Name", planet.Name).
		Int("PlayerNum", player.Num).
		Str("PlayerName", player.Race.PluralName).
		Logger()
	return &production{
		log:       producerLogger,
		rules:     rules,
		planet:    planet,
		player:    player,
		estimator: NewCompletionEstimator(),
	}
}

type production struct {
	log       zerolog.Logger
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
	Allocated Cost          `json:"allocated"`
	Tags      Tags          `json:"tags"`
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
	packets           []Cargo
	scanner           bool
	reset             bool
	starbase          *ShipDesign
	alchemy           Mineral
	mines             int
	factories         int
	defenses          int
	terraformResults  []TerraformResult
	messages          []PlayerMessage
	completed         bool
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

// produce all items in the production queue
func (p *production) produce() (productionResult, error) {
	planet := p.planet
	costCalculator := NewCostCalculator()
	result := productionResult{}
	available := Cost{Resources: planet.Spec.ResourcesPerYearAvailable}.AddCargoMinerals(planet.Cargo)
	newQueue := []ProductionQueueItem{}
	for itemIndex := range planet.ProductionQueue {
		item := planet.ProductionQueue[itemIndex]
		maxBuildable := planet.maxBuildable(p.player, item.Type)
		var err error
		var cost Cost
		if item.Type == QueueItemTypeStarbase && planet.Spec.HasStarbase {
			cost, err = costCalculator.StarbaseUpgradeCost(p.rules, p.player.TechLevels, p.player.Race.Spec, planet.Starbase.Tokens[0].design, item.design)
			if err != nil {
				p.log.Error().
					Err(err).
					Int("DesignNum", item.design.Num).
					Int("OldDesignNum", planet.Starbase.Tokens[0].design.Num).
					Msgf("StarbaseUpgradeCost returned error: %v", err)
				return productionResult{}, fmt.Errorf("failed to compute starbase upgrade cost")
			}
		} else if item.Type == QueueItemTypeStarbase || item.Type == QueueItemTypeShipToken {
			cost, err = costCalculator.GetDesignCost(p.rules, p.player.TechLevels, p.player.Race.Spec, item.design)
			if err != nil {
				p.log.Error().
					Err(err).
					Int("DesigNum", item.design.Num).
					Msgf("GetDesignCost returned error: %v", err)
				return productionResult{}, fmt.Errorf("failed to get design cost")
			}
		} else {
			cost, err = costCalculator.CostOfOne(p.player, item)
			if err != nil {
				p.log.Error().
					Err(err).
					Int("DesignNum", item.design.Num).
					Str("ItemType", string(item.Type)).
					Int("ItemQuantity", item.Quantity).
					Msgf("CostOfOne returned error: %v", err)
				return productionResult{}, fmt.Errorf("failed to compute cost of %s", item.Type)
			}
		}

		// Infinite is the constant int of -1, but for our purposes we want a very large number
		if maxBuildable == Infinite {
			maxBuildable = math.MaxInt
		}

		// check for auto items we should skip
		// we skip auto items if we can't build any more because we don't have the minerals
		if item.Type.IsAuto() && (maxBuildable <= 0 || available.DivideByMineral(cost.ToMineral()) < 1) {
			result.itemsBuilt = append(result.itemsBuilt, itemBuilt{index: item.index, skipped: true})
			newQueue = append(newQueue, item)

			// we skipped the last item in the queue, so we're done
			if itemIndex == len(planet.ProductionQueue)-1 {
				result.completed = true
			}
			continue
		}

		// make sure this item is buildable, and if not, send the player a message and move on.
		if message, valid := p.validateItem(item, maxBuildable, planet); !valid {
			// skip this item and remove it from the queue
			result.messages = append(result.messages, message)
			result.itemsBuilt = append(result.itemsBuilt, itemBuilt{index: item.index, never: true})
			continue
		}

		// add in any previously allocated resources to this item into our pot
		available = available.Add(item.Allocated)
		item.Allocated = Cost{}

		// determine how many we can build
		numBuilt, spent := p.getNumBuilt(item, cost, available, maxBuildable)

		// deduct what was built from available
		available = available.Minus(spent)

		// we built something, record ship buildings/packets for later, add installations and terraform right now
		if numBuilt > 0 {
			// add mines and factories to the planet, terraform
			p.addPlanetaryInstallations(item, numBuilt)

			if item.Type.IsTerraform() {
				result.terraformResults = append(result.terraformResults, p.terraformPlanet(numBuilt)...)
			}

			p.updateProductionResult(item, numBuilt, cost, &result)

			// if we built mineral alchemy, add it back in to our available amount
			available = available.Add(result.alchemy.ToCost())

			result.itemsBuilt = append(result.itemsBuilt, itemBuilt{index: item.index, queueItemType: item.Type, designNum: item.DesignNum, numBuilt: numBuilt})

			// planets are ending up with negative minerals. Trying to figure out why...
			if available.MinZero() != available {
				p.log.Warn().
					Str("Cargo", fmt.Sprintf("%+v", planet.Cargo)).
					Str("ProductionQueue", fmt.Sprintf("%+v", planet.ProductionQueue)).
					Str("itemResult", fmt.Sprintf("%+v", result)).
					Msgf("available minerals and resources went negative - available: %+v", available)
				available = available.MinZero()
			}
		}

		if itemIndex == len(planet.ProductionQueue)-1 && (numBuilt >= item.Quantity || numBuilt >= maxBuildable) {
			// we built all of the last item in the queue, we're all done
			result.completed = true
			if item.Type.IsAuto() {
				// append the unfinished queue back to the end of our remaining items
				newQueue = append(newQueue, planet.ProductionQueue[itemIndex:]...)
			}
			break
		}

		if item.Type.IsAuto() {
			// auto items stay in the queue
			newQueue = append(newQueue, item)

			// if we are an auto item and we have resources left, move on to the next
			// auto items don't block the queue
			if available.Resources > 0 {
				// if we still have auto items to build, and
				// if we have enough minerals to complete an auto item, add a concrete one to the queue

				if !(numBuilt >= item.Quantity || numBuilt >= maxBuildable) && available.DivideByMineral(cost.ToMineral()) >= 1 {
					// add the partially completed concrete item to the top of the queue
					newQueue = append([]ProductionQueueItem{
						{
							Type:      item.Type.concreteType(),
							Quantity:  1,
							Allocated: Cost{Resources: available.Resources},
							index:     -1, // we don't track concrete auto items, we only care about the first fully built auto item
						}}, newQueue...)
					available.Resources -= newQueue[0].Allocated.Resources

					if itemIndex < len(planet.ProductionQueue)-1 {
						// if this isn't the last item, append the unfinished queue back to the end of our remaining items
						newQueue = append(newQueue, planet.ProductionQueue[itemIndex+1:]...)
					}
					break
				}
			} else {
				if itemIndex < len(planet.ProductionQueue)-1 {
					// if this isn't the last item, append the unfinished queue back to the end of our remaining items
					newQueue = append(newQueue, planet.ProductionQueue[itemIndex+1:]...)
				}
				break
			}
		} else {
			item.Quantity -= numBuilt
			if item.Quantity > 0 {
				// allocate remaining resources to this partially built item
				item.Allocated.Resources = MinInt(cost.Resources, available.Resources)
				available.Resources -= item.Allocated.Resources
				planet.ProductionQueue[itemIndex] = item

				// keep it in the queue for next time
				newQueue = append(newQueue, planet.ProductionQueue[itemIndex:]...)
				break
			}
		}

	}
	// replace the queue with what's leftover
	planet.ProductionQueue = newQueue
	planet.Cargo = Cargo{available.Ironium, available.Boranium, available.Germanium, planet.Cargo.Colonists}
	if planet.Cargo.MinZero() != planet.Cargo {
		p.log.Warn().
			Str("Cargo", fmt.Sprintf("%+v", planet.Cargo)).
			Str("productionResult", fmt.Sprintf("%+v", result)).
			Msgf("planet cargo was negative after production: %s", planet.Cargo.PrettyString())
		// planet.Cargo = planet.Cargo.MinZero()
	}

	// any leftover resources go back to the player for research
	result.leftoverResources = available.Resources
	return result, nil
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
	case QueueItemTypePlanetaryScanner:
		p.planet.Scanner = true
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
		return newPlanetMessage(PlayerMessagePlanetBuiltInvalidMineralPacketNoMassDriver, planet), false
	}
	if item.Type.IsPacket() && planet.PacketTargetNum == None {
		return newPlanetMessage(PlayerMessagePlanetBuiltInvalidMineralPacketNoTarget, planet), false
	}
	if !item.Type.IsAuto() && maxBuildable == 0 {
		// can't build this, skip it
		// it shouldn't have been ever added to the queue, but just in case of a bug
		return newPlanetMessage(PlayerMessagePlanetBuiltInvalidItem, planet).withSpec(PlayerMessageSpec{Name: planet.Name, QueueItemType: item.Type}), false
	}

	return PlayerMessage{}, true
}

// the result of processing an item in the queue
type processQueueItemResult struct {
	numBuilt int
	spent    Cost
}

// determine how many of this production item we can build
func (p *production) getNumBuilt(item ProductionQueueItem, cost, availableToSpend Cost, maxBuildable int) (numBuilt int, spent Cost) {

	// add in anything allocated in previous turns
	availableToSpend = availableToSpend.Add(item.Allocated)
	item.Allocated = Cost{}

	if cost == (Cost{}) {
		return MinInt(item.Quantity, maxBuildable), Cost{}
	}

	// figure out how many we can build
	// and make sure we only build up to the quantity, and we don't build more than the planet supports
	numBuilt = MaxInt(0, MinInt(item.Quantity, maxBuildable, availableToSpend.NumBuildable(cost)))
	spent = cost.MultiplyInt(numBuilt)

	return numBuilt, spent
}

// add built items to planet, build fleets, update player messages, etc
func (p *production) updateProductionResult(item ProductionQueueItem, numBuilt int, cost Cost, result *productionResult) {
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
		cargo := cost.MultiplyFloat64(1 / p.player.Race.Spec.PacketMineralCostFactor).MultiplyInt(numBuilt).ToCargo()
		result.packets = append(result.packets, cargo)
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
