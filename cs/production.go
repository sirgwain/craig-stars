package cs

import (
	"fmt"
	"math"
	"slices"

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

// An item in a production queue.
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

type QueueItemType string

const (
	QueueItemTypeNone                   QueueItemType = ""
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

// true if this is a packet type (concrete or auto)
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

// A record used by the production estimator to record unbuilt
// ProductionQueueItem completion times
type QueueItemCompletionEstimate struct {
	Canceled        bool `json:"canceled,omitempty"`        // Whether an item is canceled due to an invalid order
	YearsToBuildOne int  `json:"yearsToBuildOne,omitempty"` // Years to build (or skip) the first item of its type
	YearsToBuildAll int  `json:"yearsToBuildAll,omitempty"` // Years to build (or skip) the last item of its type
	YearsToSkipAuto int  `json:"yearsToSkipAuto,omitempty"` // Years to skip the first auto item in a queue
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
		// no queue, no life
		result.messages = append(result.messages,
			newPlanetMessage(PlayerMessagePlanetProductionQueueEmpty, planet))
		return productionResult{leftoverResources: planet.Spec.ResourcesPerYearAvailable, completed: true}, nil
	}

	// TODO: Fix auto mineral alchemy:
	// * Only has 1 max quantity
	// * Always blocks queue if final item,
	// * otherwise only blocks if subsequent item is out of minerals

	// TODO: Add support for multiple starbase designs in a queue

	var (
		// tracker of resources/minerals available to spend
		available = Cost{Resources: planet.Spec.ResourcesPerYearAvailable}.AddMineral(planet.Cargo.ToMineral())
		// new queue of production items, used to recreate an updated queue
		newQueue      = []ProductionQueueItem{}
		builtAnything bool
	)

	c := NewCostCalculator()

	// check each item in the queue in order
	for itemIndex, item := range planet.ProductionQueue {
		itemCost, err := c.GetItemCost(p.rules, p.player, planet, item)
		if err != nil {
			p.log.Error().
				Err(err).
				Any("item", item).
				Msgf("produce() obtained error when calculating item costs: %v", err)
			return productionResult{}, err
		}

		maxBuildable := planet.MaxBuildable(p.player, item.Type)
		if maxBuildable == Infinite {
			// Infinite is the constant int of -1, but for our purposes we want a very large number
			maxBuildable = math.MaxInt
		}

		// If this is a concrete item and we haven't built anything yet,
		// check to make sure we aren't trying to build over cap.
		// We do this for everything ahead of us upon building (or trying to build)
		// an item, but this ensures we don't forget to check before that happens.
		if !builtAnything && !item.Type.IsAuto() {
			overCap := item.Quantity - min(item.Quantity, maxBuildable)
			if overCap > 0 {
				p.log.Debug().
					Any("Item", item).
					Int("Qty", item.Quantity).
					Int("maxBuildable", maxBuildable).
					Int("New Quantity", item.Quantity-overCap).
					Msgf("clamping queue item quantity")
				item.Quantity -= overCap
			}

			if item.Quantity <= 0 {
				// quantity below 0; skip building item
				available = available.Add(item.Allocated) // refund previously allocated amount
				result.itemsBuilt = append(result.itemsBuilt,
					itemBuilt{index: item.index, canceled: true})
				p.updateCanceledMessage(&result, item, overCap, maxBuildable)
				continue
			}
		}

		// make sure our concrete packets are buildable, notifying the player and canceling if not.
		if msgType, valid := p.validatePacket(item, planet, result.starbase); !valid {
			if item.Type == QueueItemTypeAutoMineralPacket {
				// auto mineral packets are skipped instead of canceled
				maxBuildable = 0
			} else {
				p.log.Debug().
					Int("Index", itemIndex).
					Str("QueueItemType", string(item.Type)).
					Int("MessageType", int(msgType)).
					Msg("Canceling packet")

				available = available.Add(item.Allocated)
				p.updatePacketCanceledMessage(&result, item, msgType,
					itemCost.ToMineral().MultiplyFloat64(1/p.player.Race.Spec.PacketMineralCostFactor, math.Floor).Total())
				result.itemsBuilt = append(result.itemsBuilt, itemBuilt{index: item.index, canceled: true})
				continue
			}
		}

		// Dump in any previously allocated resources for this item into our pot.
		available = available.Add(item.Allocated)
		item.Allocated = Cost{}

		// check for auto items we should skip due to not being buildable or lacking minerals.
		// Stars! doesn't bother starting auto items unless we have enough minerals for 1 full batch
		// (likely to prevent accidental queue blockages)
		if item.Type.IsAuto() && (maxBuildable <= 0 || available.DivideMineral(itemCost.ToMineral()) < 1) {
			result.itemsBuilt = append(result.itemsBuilt, itemBuilt{index: item.index, skipped: true})
			newQueue = append(newQueue, item) // auto items stick around

			// if we skipped the last item in the queue, mark result as done
			if itemIndex == len(planet.ProductionQueue)-1 {
				result.completed = true
			}
			continue
		}

		// After all that, we can finally try to build the dang thing.

		// determine how many copies to build and dock cost
		numBuilt, spent := p.getNumBuilt(item, itemCost, available, maxBuildable)
		available = available.Subtract(spent)

		// if we can't build anything, go directly to cleanup.
		// Do NOT pass GO, do NOT collect $200.
		if numBuilt <= 0 {
			goto checkDone
		}

		// woot woot, we built a thing!
		builtAnything = true

		// spend money, record info and add installations/terraforming
		p.updateProductionResult(&result, item, numBuilt, itemCost)

		// if we built mineral alchemy, add it back into our pot to use later
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

		// Elide any concrete items ahead of us inside the queue that are over cap.
		if itemIndex < len(planet.ProductionQueue)-1 {
			modQueue := planet.ProductionQueue[:itemIndex+1] // everything before us can stay
			for _, item := range planet.ProductionQueue[itemIndex+1:] {
				if item.Type.IsAuto() {
					// leave auto items alone
					modQueue = append(modQueue, item)
					continue
				}

				maxBuildable := planet.MaxBuildable(p.player, item.Type)
				if maxBuildable == Infinite {
					// infinite maxBuildable
					modQueue = append(modQueue, item)
					continue
				}

				// clamp item quantity down to maxBuildable
				overCap := item.Quantity - maxBuildable
				if overCap > 0 {
					p.log.Debug().
						Any("Item", item).
						Int("Qty", item.Quantity).
						Int("maxBuildable", maxBuildable).
						Int("New Quantity", item.Quantity-overCap).
						Msgf("clamping queue item quantity down to maxBuildable")
					item.Quantity -= overCap // a-(a-b) = a-a+b = b
				}

				if item.Quantity > 0 {
					modQueue = append(modQueue, item)
				} else {
					// quantity <= 0; remove from new queue
					available = available.Add(item.Allocated) // refund previously allocated amount
					result.itemsBuilt = append(result.itemsBuilt,
						itemBuilt{index: item.index, queueItemType: item.Type, canceled: true})
					p.updateCanceledMessage(&result, item, overCap, maxBuildable)
				}
			}

			planet.ProductionQueue = planet.ProductionQueue[:len(modQueue)]
		}

	checkDone:
		if itemIndex == len(planet.ProductionQueue)-1 && (numBuilt >= item.Quantity || numBuilt >= maxBuildable) {
			// we built (or tried to) built the last item in the queue; we're all done
			result.completed = true
			if item.Type.IsAuto() {
				// tack on the auto item to the end of the queue
				newQueue = append(newQueue, item)
			}
			break
		}

		if !item.Type.IsAuto() {
			// concrete items never reset, so dock amount built from remaining quantity
			item.Quantity -= numBuilt
			if item.Quantity < 0 {
				// should never happen, but covering our bases
				p.log.Warn().
					Any("Item", item).
					Int("Qty", item.Quantity).
					Int("PrevQty", item.Quantity+numBuilt).
					Int("NumBuilt", numBuilt).
					Msgf("concrete item quantity went negative after building")
				continue
			} else if item.Quantity == 0 {
				// fully built item; move on
				continue
			} else {
				// couldn't finish entire concrete item; done building

				// allocate any remaining resources to partially built item
				item.Allocated = p.allocatePartialBuild(itemCost, available)
				available = available.Subtract(item.Allocated)
				planet.ProductionQueue[itemIndex] = item

				// keep it & everything in front and break out
				newQueue = append(newQueue, planet.ProductionQueue[itemIndex:]...)
				break
			}
		}

		// auto items stay in the queue after being built
		newQueue = append(newQueue, item)

		if available.Resources <= 0 {
			// all resources spent; wrap up
			if itemIndex < len(planet.ProductionQueue)-1 {
				// if this isn't the last item, tack the rest of the queue back on
				newQueue = append(newQueue, planet.ProductionQueue[itemIndex+1:]...)
			}
			break
		}

		if numBuilt >= item.Quantity || numBuilt >= maxBuildable ||
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
	}

	// ping player if we finish queue
	if result.completed {
		result.messages = append(result.messages,
			newPlanetMessage(PlayerMessagePlanetProductionQueueComplete, planet))
	}

	// replace queue & surface minerals with leftovers
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

// validate a packet in the production queue
func (p *producer) validatePacket(item ProductionQueueItem, planet *Planet, builtBase *ShipDesign) (msgType PlayerMessageType, valid bool) {
	if item.Type.IsPacket() && !planet.Spec.HasMassDriver &&
		(builtBase == nil || builtBase.Spec.SafePacketSpeed == 0) {
		// We have no mass driver ATM and our current queued starbase
		// either doesn't exist or can't fling packets
		return PlayerMessagePlanetBuiltInvalidMineralPacketNoMassDriver, false
	}

	if item.Type.IsPacket() && planet.PacketTargetNum == None {
		return PlayerMessagePlanetBuiltInvalidMineralPacketNoTarget, false
	}

	return PlayerMessageNone, true
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

// updateProductionResult updates a production result with information about a newly built item,
// as well as applying corresponding changes to the planet.
func (p *producer) updateProductionResult(result *productionResult, item ProductionQueueItem, numBuilt int, itemCost Cost) {
	switch item.Type {
	case QueueItemTypeAutoMineralAlchemy, QueueItemTypeMineralAlchemy:
		result.alchemy += numBuilt
	case QueueItemTypeAutoMines, QueueItemTypeMine:
		result.mines += numBuilt
		p.planet.Mines += numBuilt
	case QueueItemTypeAutoFactories, QueueItemTypeFactory:
		result.factories += numBuilt
		p.planet.Factories += numBuilt
	case QueueItemTypeAutoDefenses, QueueItemTypeDefenses:
		result.defenses += numBuilt
		p.planet.Defenses += numBuilt
	case QueueItemTypeAutoMinTerraform, QueueItemTypeAutoMaxTerraform, QueueItemTypeTerraformEnvironment:
		// terraform 1 by 1 to accurately simulate results
		result.terraformResults = append(result.terraformResults,
			p.terraformPlanet(numBuilt)...)
	case QueueItemTypeAutoMineralPacket, QueueItemTypeMixedMineralPacket, QueueItemTypeIroniumMineralPacket, QueueItemTypeBoraniumMineralPacket, QueueItemTypeGermaniumMineralPacket:
		// add this packet's minerals to the production result
		// so it can be added as packets to the universe later
		// Multiply by 1/PacketMineralCostFactor to simulate overhead (pay 132 kT; only get 120)
		minsPerPacket := MultiplyCost(itemCost, 1/p.player.Race.Spec.PacketMineralCostFactor)
		result.packets = result.packets.AddCostMinerals(MultiplyCost(minsPerPacket, numBuilt))
	case QueueItemTypeShipToken:
		result.tokens = append(result.tokens, builtShip{
			ShipToken: ShipToken{
				Quantity:  numBuilt,
				design:    item.design,
				DesignNum: item.DesignNum,
			},
			tags: item.Tags,
		})
	case QueueItemTypeStarbase:
		result.starbase = item.design
	case QueueItemTypePlanetaryScanner:
		result.scanner = true
		p.planet.Scanner = true
	case QueueItemTypeGenesisDevice:
		result.reset = true
	}
}

// terraform a planet during production and save the results for messages
func (p *producer) terraformPlanet(numSteps int) []TerraformResult {
	planet, player := p.planet, p.player
	terraformer := NewTerraformer()
	terraformResults := make([]TerraformResult, numSteps)

	for i := range numSteps {
		// terraform one step at a time to ensure the best things get terraformed first
		terraformResults[i] = terraformer.TerraformOneStep(planet, player, nil, false)
	}

	return terraformResults
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

func (p *producer) updateCanceledMessage(result *productionResult, item ProductionQueueItem, numCanceled, maxBuildable int) {
	// For normal items being removed, we record removals for each item type separately
	if index := slices.IndexFunc(result.messages, func(message PlayerMessage) bool {
		return message.Type == PlayerMessagePlanetBuiltInvalidItem && message.Spec.QueueItemType == item.Type
	}); index == -1 {
		// message for this type doesn't exist; add one
		result.messages = append(result.messages, newPlanetMessage(PlayerMessagePlanetBuiltInvalidItem, p.planet).
			withSpec(PlayerMessageSpec{
				Name:          p.planet.Name,
				Cost:          item.Allocated,
				QueueItemType: item.Type,
				Amount:        numCanceled,                 // Amount canceled
				Amount2:       maxBuildable,                // MaxBuildable
				PrevAmount:    item.Quantity + numCanceled, // Total amount
			}))
	} else {
		// update the previous message with this item's amounts
		result.messages[index].Spec.Amount += numCanceled
		result.messages[index].Spec.PrevAmount += item.Quantity + numCanceled
		result.messages[index].Spec.Cost = result.messages[index].Spec.Cost.Add(item.Allocated)
	}
}

func (p *producer) updatePacketCanceledMessage(result *productionResult, item ProductionQueueItem, msgType PlayerMessageType, weight int) {
	// For packets, we don't care what *type* of packet got canceled, only that *a* packet was canceled with this message.
	if index := slices.IndexFunc(result.messages, func(message PlayerMessage) bool {
		return message.Type == msgType
	}); index == -1 {
		// message for this type doesn't exist; add one
		result.messages = append(result.messages, newPlanetMessage(msgType, p.planet).
			withSpec(PlayerMessageSpec{
				Name:          p.planet.Name,
				QueueItemType: item.Type, // Keep track of item type if this is the _only_ packet in the queue
				Cost:          item.Allocated,
				Amount:        weight * item.Quantity, // total kT of packet cargo
				Amount2:       1,                      // number of orders canceled
				PrevAmount:    item.Quantity,          // Items built (equal to items canceled)
			}))
	} else {
		// update the previous message with this item's amounts
		if result.messages[index].Spec.QueueItemType != item.Type {
			// If this is a different kind of packet than earlier, remove QueueItemType
			// to indicate multiple packet types were canceled
			result.messages[index].Spec.QueueItemType = QueueItemTypeNone
		}

		result.messages[index].Spec.PrevAmount += item.Quantity
		result.messages[index].Spec.Amount += weight * item.Quantity
		result.messages[index].Spec.Amount2++
		result.messages[index].Spec.Cost = result.messages[index].Spec.Cost.Add(item.Allocated)
	}
}
