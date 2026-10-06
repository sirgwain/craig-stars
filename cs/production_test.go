//go:build !wasi && !wasm

package cs

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_production_produceOneConcreteMine(t *testing.T) {
	player, planet := newTestPlayerPlanet()

	// build 1 mine
	planet.ProductionQueue = []ProductionQueueItem{{Type: QueueItemTypeMine, Quantity: 1}}
	planet.Cargo = Cargo{10, 20, 30, 2500}
	planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 100, MaxPossibleMines: 100, MaxPopulation: 1_000_000}
	planet.Mines = 0

	// should build 1 mine, leaving empty queue
	producer := newProducer(testLogger, &rules, planet, player)
	producer.produce()
	assert.Equal(t, 1, planet.Mines)
	assert.Equal(t, 0, len(planet.ProductionQueue))

	// build 5 auto mines, leaving them in the queue
	planet.ProductionQueue = []ProductionQueueItem{{Type: QueueItemTypeAutoMines, Quantity: 5}}
	planet.Cargo = Cargo{10, 20, 30, 2500}
	planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 100, MaxMines: 100, MaxPossibleMines: 100, MaxPopulation: 1_000_000}
	planet.Mines = 0
	player.Messages = []PlayerMessage{}

	// should build 5 mine, leaving the auto build in the queu
	producer = newProducer(testLogger, &rules, planet, player)
	producer.produce()
	assert.Equal(t, 5, planet.Mines)
	assert.Equal(t, 1, len(planet.ProductionQueue))
	assert.Equal(t, QueueItemTypeAutoMines, planet.ProductionQueue[0].Type)
	assert.Equal(t, 5, planet.ProductionQueue[0].Quantity)
}

func Test_production_produceAutoFactories(t *testing.T) {
	player, planet := newTestPlayerPlanet()

	// build 5 auto factories, leaving them in the queue
	planet.ProductionQueue = []ProductionQueueItem{{Type: QueueItemTypeAutoFactories, Quantity: 5}}
	planet.Cargo = Cargo{10, 20, 30, 2500}
	planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 100, MaxFactories: 100, MaxPossibleFactories: 100, MaxPopulation: 1_000_000}
	planet.Factories = 0
	player.Messages = []PlayerMessage{}

	// should build 5 mine, leaving the auto build in the queu
	producer := newProducer(testLogger, &rules, planet, player)
	result, err := producer.produce()
	assert.Nil(t, err)

	assert.Equal(t, 5, planet.Factories)
	assert.Equal(t, 1, len(planet.ProductionQueue))
	assert.Equal(t, QueueItemTypeAutoFactories, planet.ProductionQueue[0].Type)
	assert.Equal(t, 5, planet.ProductionQueue[0].Quantity)
	assert.Equal(t, 50, result.leftoverResources) // 50 resources leftover after building 5 mines
}

func Test_production_produceFactoriesThenMinesWithLowGerm(t *testing.T) {
	player, planet := newTestPlayerPlanet()

	// build 2/5 auto factories and 5 mines
	// we only have enough minerals on hand to build 2 factories so
	// we will skip after 2 and build mines
	planet.ProductionQueue = []ProductionQueueItem{
		{Type: QueueItemTypeAutoFactories, Quantity: 5},
		{Type: QueueItemTypeAutoMines, Quantity: 5},
	}
	// give a planet with enough germanium to build 2.5 factories
	// and enough resources to build all factories and all mines
	planet.Cargo = Cargo{0, 0, 10, 2500}
	planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 100, MaxFactories: 100, MaxMines: 100, MaxPossibleFactories: 100, MaxPossibleMines: 100, MaxPopulation: 1_000_000}
	planet.Factories = 0
	player.Messages = []PlayerMessage{}

	// should build 2 factories and 5 mines, leaving the auto builds in the queue
	producer := newProducer(testLogger, &rules, planet, player)
	producer.produce()
	assert.Equal(t, 2, planet.Factories)
	assert.Equal(t, 5, planet.Mines)
	assert.Equal(t, 2, len(planet.ProductionQueue))
	assert.Equal(t, QueueItemTypeAutoFactories, planet.ProductionQueue[0].Type)
	assert.Equal(t, 5, planet.ProductionQueue[0].Quantity)
	assert.Equal(t, QueueItemTypeAutoMines, planet.ProductionQueue[1].Type)
	assert.Equal(t, 5, planet.ProductionQueue[1].Quantity)
}

func Test_production_produceFactoriesToMaxAndPartialMine(t *testing.T) {
	player, planet := newTestPlayerPlanet()

	// build 2/5 auto factories and one partial mine
	planet.ProductionQueue = []ProductionQueueItem{
		{Type: QueueItemTypeAutoFactories, Quantity: 5},
		{Type: QueueItemTypeAutoMines, Quantity: 5},
	}
	planet.Cargo = Cargo{0, 0, 8, 2500}
	planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 10*2 + 8, MaxFactories: 100, MaxMines: 100, MaxPossibleFactories: 100, MaxPossibleMines: 100, MaxPopulation: 1_000_000}
	planet.Factories = 0
	player.Messages = []PlayerMessage{}

	// should build 2 factories, 1 mine and one partial mine, leaving the auto builds in the queue
	producer := newProducer(testLogger, &rules, planet, player)
	producer.produce()
	assert.Equal(t, 2, planet.Factories)
	assert.Equal(t, 1, planet.Mines)
	assert.Equal(t, 3, len(planet.ProductionQueue))
	assert.Equal(t, QueueItemTypeMine, planet.ProductionQueue[0].Type)
	assert.Equal(t, 1, planet.ProductionQueue[0].Quantity)
	assert.Equal(t, Cost{0, 0, 0, 3}, planet.ProductionQueue[0].Allocated)
	assert.Equal(t, QueueItemTypeAutoFactories, planet.ProductionQueue[1].Type)
	assert.Equal(t, 5, planet.ProductionQueue[1].Quantity)
	assert.Equal(t, QueueItemTypeAutoMines, planet.ProductionQueue[2].Type)
	assert.Equal(t, 5, planet.ProductionQueue[2].Quantity)

}

func Test_production_produceDefensesToMax(t *testing.T) {
	player, planet := newTestPlayerPlanet()

	// build 100 auto defenses when we have 90 already and one partial in the queue
	planet.ProductionQueue = []ProductionQueueItem{
		{Type: QueueItemTypeDefenses, Quantity: 1, Allocated: Cost{5, 5, 5, 14}},
		{Type: QueueItemTypeAutoDefenses, Quantity: 100},
	}
	planet.Cargo = Cargo{5000, 5000, 5000, 1_000_000}
	planet.Defenses = 90
	planet.Spec = ComputePlanetSpec(&rules, player, planet)
	player.Messages = []PlayerMessage{}

	// should end up with 100 defenses and the auto defenses still in the queue
	producer := newProducer(testLogger, &rules, planet, player)
	producer.produce()
	assert.Equal(t, 100, planet.Defenses)
	assert.Equal(t, 1, len(planet.ProductionQueue))
	assert.Equal(t, QueueItemTypeAutoDefenses, planet.ProductionQueue[0].Type)
	assert.Equal(t, 100, planet.ProductionQueue[0].Quantity)
}

func Test_production_produceFactoriesToMaxThenMines(t *testing.T) {
	player, planet := newTestPlayerPlanet()

	// auto build up to the max
	planet.ProductionQueue = []ProductionQueueItem{
		{Type: QueueItemTypeAutoFactories, Quantity: 10},
		{Type: QueueItemTypeAutoMines, Quantity: 10},
	}
	planet.Cargo = Cargo{1000, 1000, 1000, 100}
	planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 1000, MaxFactories: 10, MaxMines: 10, MaxPossibleFactories: 10, MaxPossibleMines: 10, MaxPopulation: 1_000_000}
	planet.Factories = 9
	planet.Mines = 0
	player.Messages = []PlayerMessage{}

	// should build 1 factories, 10 mines and have leftover for research
	producer := newProducer(testLogger, &rules, planet, player)
	result, err := producer.produce()
	assert.Nil(t, err)
	assert.Equal(t, 10, planet.Factories)
	assert.Equal(t, 10, planet.Mines)
	assert.Equal(t, 2, len(planet.ProductionQueue))
	assert.Equal(t, QueueItemTypeAutoFactories, planet.ProductionQueue[0].Type)
	assert.Equal(t, 10, planet.ProductionQueue[0].Quantity)
	assert.Equal(t, QueueItemTypeAutoMines, planet.ProductionQueue[1].Type)
	assert.Equal(t, 10, planet.ProductionQueue[1].Quantity)
	assert.Equal(t, 940, result.leftoverResources)

}

func Test_production_producePartialMine(t *testing.T) {
	player, planet := newTestPlayerPlanet()

	// auto build factories, but we have no mines or minerals, so they'll never build
	// until we build mines and mine a bit
	planet.ProductionQueue = []ProductionQueueItem{
		{Type: QueueItemTypeAutoFactories, Quantity: 100},
		{Type: QueueItemTypeAutoMines, Quantity: 100},
	}
	planet.Cargo = Cargo{0, 0, 0, 25}
	planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 2, MaxFactories: 10, MaxMines: 10, MaxPossibleFactories: 10, MaxPossibleMines: 10, MaxPopulation: 1_000_000}
	player.Messages = []PlayerMessage{}

	// should build nothing, but queue up a mine partially done
	producer := newProducer(testLogger, &rules, planet, player)
	producer.produce()

	// nothing built
	assert.Equal(t, 0, planet.Factories)
	assert.Equal(t, 0, planet.Mines)
	assert.Equal(t, 3, len(planet.ProductionQueue))

	// concrete mine partially completed with 2 resources
	assert.Equal(t, QueueItemTypeMine, planet.ProductionQueue[0].Type)
	assert.Equal(t, 1, planet.ProductionQueue[0].Quantity)
	assert.Equal(t, Cost{Resources: 2}, planet.ProductionQueue[0].Allocated)
	assert.Equal(t, QueueItemTypeAutoFactories, planet.ProductionQueue[1].Type)
	assert.Equal(t, 100, planet.ProductionQueue[1].Quantity)
	assert.Equal(t, Cost{}, planet.ProductionQueue[1].Allocated)
	assert.Equal(t, QueueItemTypeAutoMines, planet.ProductionQueue[2].Type)
	assert.Equal(t, 100, planet.ProductionQueue[2].Quantity)
	assert.Equal(t, Cost{}, planet.ProductionQueue[2].Allocated)

}

func Test_production_producePartialFactoryAndMoreAuto(t *testing.T) {
	player, planet := newTestPlayerPlanet()

	// build the half completed factory, keep building more factories but don't add a partial mine
	planet.ProductionQueue = []ProductionQueueItem{
		{Type: QueueItemTypeFactory, Quantity: 1, Allocated: Cost{Germanium: 2, Resources: 5}},
		{Type: QueueItemTypeAutoMinTerraform, Quantity: 1},
		{Type: QueueItemTypeAutoFactories, Quantity: 100},
		{Type: QueueItemTypeAutoMines, Quantity: 100},
	}
	planet.Cargo = Cargo{361, 382, 1173, 331}
	planet.Mines = 11
	planet.Factories = 16
	planet.Spec = ComputePlanetSpec(&rules, player, planet)

	producer := newProducer(testLogger, &rules, planet, player)
	producer.produce()

	// build partial factory, then more factories until out of resources
	// then allocate to a new concrete factory
	assert.Equal(t, 20, planet.Factories)
	assert.Equal(t, 11, planet.Mines)
	assert.Equal(t, 4, len(planet.ProductionQueue))
	assert.Equal(t, QueueItemTypeFactory, planet.ProductionQueue[0].Type)
	assert.Equal(t, 1, planet.ProductionQueue[0].Quantity)
	assert.Equal(t, Cost{Germanium: 2, Resources: 7}, planet.ProductionQueue[0].Allocated)

	// auto orders remain empty of allocated resources
	assert.Equal(t, QueueItemTypeAutoMinTerraform, planet.ProductionQueue[1].Type)
	assert.Equal(t, 1, planet.ProductionQueue[1].Quantity)
	assert.Equal(t, Cost{}, planet.ProductionQueue[1].Allocated)
	assert.Equal(t, QueueItemTypeAutoFactories, planet.ProductionQueue[2].Type)
	assert.Equal(t, 100, planet.ProductionQueue[2].Quantity)
	assert.Equal(t, Cost{}, planet.ProductionQueue[2].Allocated)

	assert.Equal(t, QueueItemTypeAutoMines, planet.ProductionQueue[3].Type)
	assert.Equal(t, 100, planet.ProductionQueue[3].Quantity)
	assert.Equal(t, Cost{}, planet.ProductionQueue[3].Allocated)

}

func Test_production_producePartialFactory(t *testing.T) {
	player, planet := newTestPlayerPlanet()

	// build the half completed factory, keep building more factories but don't add a partial mine
	planet.ProductionQueue = []ProductionQueueItem{
		{Type: QueueItemTypeFactory, Quantity: 1, Allocated: Cost{Germanium: 2, Resources: 5}},
		{Type: QueueItemTypeAutoMinTerraform, Quantity: 1},
		{Type: QueueItemTypeAutoFactories, Quantity: 100},
		{Type: QueueItemTypeAutoMines, Quantity: 100},
	}
	planet.Cargo = Cargo{7, 2, 1, 37}
	planet.Mines = 2
	planet.Factories = 1
	planet.Spec = ComputePlanetSpec(&rules, player, planet)

	// should build nothing, but queue up a mine partially done
	producer := newProducer(testLogger, &rules, planet, player)
	producer.produce()

	// We should consume 1kT germanium and allocate appropriate resources to match
	assert.Equal(t, Cargo{7, 2, 0, 37}, planet.Cargo)
	assert.Equal(t, Cost{Germanium: 3, Resources: 7}, planet.ProductionQueue[0].Allocated)

}

func Test_production_produceBuildToMinesFactoriesToMax(t *testing.T) {
	player, planet := newTestPlayerPlanet()

	// make mines/factories cheap so we can build them
	player.Race.MineCost = 1
	player.Race.FactoryCost = 1
	player.Race.Spec = ComputeRaceSpec(&player.Race, &rules)

	// auto build with future growth taken into account
	planet.ProductionQueue = []ProductionQueueItem{
		{Type: QueueItemTypeAutoMines, Quantity: 100},
		{Type: QueueItemTypeAutoFactories, Quantity: 100},
	}
	planet.Cargo = Cargo{1000, 1000, 1000, 1000}
	planet.Spec = ComputePlanetSpec(&rules, player, planet)

	// max mines for current setting
	planet.Mines = planet.Spec.MaxMines
	planet.Factories = planet.Spec.MaxFactories

	// should build nothing, but queue up a mine partially done
	producer := newProducer(testLogger, &rules, planet, player)
	producer.produce()

	// we should build mines/factories accounting for future growth
	assert.Equal(t, 115, planet.Mines)
	assert.Equal(t, 115, planet.Factories)

}

func Test_production_produceColonizerAndPartialFreighters(t *testing.T) {
	player, planet := newTestPlayerPlanet()
	player.TechLevels = TechLevel{0, 0, 3, 4, 0, 2}
	player.Race.PRT = SD
	player.Race.LRTs = Bitmask(IFE) | Bitmask(ARM) | Bitmask(BET) | Bitmask(RS)
	player.Race.PopEfficiency = 9
	player.Race.FactoryOutput = 11
	player.Race.Spec = ComputeRaceSpec(&player.Race, &rules)
	player.Spec = ComputePlayerSpec(player, &rules)

	// add two designs, a colony ship w/fuel mizer and medium freighter w/fuel mizer
	player.Designs = append(player.Designs, NewShipDesign(player.Num, 1).
		WithHull(ColonyShip.Name).
		WithSlots([]ShipDesignSlot{
			{HullComponent: FuelMizer.Name, HullSlotIndex: 1, Quantity: 1},
			{HullComponent: ColonizationModule.Name, HullSlotIndex: 2, Quantity: 1},
		}).
		WithSpec(&rules, player))

	player.Designs = append(player.Designs, NewShipDesign(player.Num, 2).
		WithHull(MediumFreighter.Name).
		WithSlots([]ShipDesignSlot{
			{HullComponent: FuelMizer.Name, HullSlotIndex: 1, Quantity: 1},
			{HullComponent: CargoPod.Name, HullSlotIndex: 2, Quantity: 1},
		}).
		WithSpec(&rules, player))

	// add designs plus some auto items
	planet.ProductionQueue = []ProductionQueueItem{
		{Type: QueueItemTypeShipToken, Quantity: 1, DesignNum: 1, design: player.Designs[0]},
		{Type: QueueItemTypeShipToken, Quantity: 2, DesignNum: 2, design: player.Designs[1], Allocated: Cost{Ironium: 40, Germanium: 29, Resources: 75}},
		{Type: QueueItemTypeAutoMines, Quantity: 250},
		{Type: QueueItemTypeAutoFactories, Quantity: 250},
		{Type: QueueItemTypeAutoMaxTerraform, Quantity: 10},
	}
	planet.Cargo = Cargo{1000, 1000, 77, 3166}

	// ships need a starbase to build them
	starbaseDesign := NewShipDesign(player.Num, 3).WithHull(SpaceStation.Name).WithSpec(&rules, player)
	player.Designs = append(player.Designs, starbaseDesign)
	starbase := newStarbase(player, planet, starbaseDesign, "Starbase")
	starbase.Spec = ComputeFleetSpec(&rules, player, &starbase)
	planet.Starbase = &starbase
	planet.Spec = ComputePlanetSpec(&rules, player, planet)

	// should build nothing, but queue up a mine partially done
	producer := newProducer(testLogger, &rules, planet, player)
	result, err := producer.produce()

	assert.Nil(t, err)
	assert.Equal(t, 3, len(result.itemsBuilt))
	assert.Equal(t, 1, result.itemsBuilt[0].designNum)
	assert.Equal(t, 1, result.itemsBuilt[0].numBuilt)
	assert.Equal(t, 2, result.itemsBuilt[1].designNum)
	assert.Equal(t, 2, result.itemsBuilt[1].numBuilt)

}

func Test_production_produceStarbaseUpgrade(t *testing.T) {
	player, planet := newTestPlayerPlanet()

	// create a new starbase
	starbaseDesign1 := NewShipDesign(player.Num, 2).WithHull(SpaceStation.Name).WithSpec(&rules, player)
	starbase1 := newStarbase(player, planet,
		starbaseDesign1,
		"Starbase",
	)
	starbaseDesign2 := NewShipDesign(player.Num, 3).WithHull(SpaceStation.Name).
		WithSlots([]ShipDesignSlot{
			{HullComponent: MassDriver5.Name, HullSlotIndex: 1, Quantity: 1},
		}).WithSpec(&rules, player)
	starbase2 := newStarbase(player, planet,
		starbaseDesign2,
		"Starbase2",
	)
	player.Designs = append(player.Designs, starbaseDesign1, starbaseDesign2)
	starbase1.Spec = ComputeFleetSpec(&rules, player, &starbase1)
	starbase2.Spec = ComputeFleetSpec(&rules, player, &starbase2)
	starbase1.Tokens[0].QuantityDamaged = 1
	starbase1.Tokens[0].Damage = 100
	planet.Starbase = &starbase1

	// build the half completed factory, keep building more factories but don't add a partial mine
	planet.ProductionQueue = []ProductionQueueItem{
		{Type: QueueItemTypeStarbase, Quantity: 1, DesignNum: 3, design: starbaseDesign2},
	}
	planet.Cargo = Cargo{1000, 1000, 1000, 10000}
	planet.Spec = ComputePlanetSpec(&rules, player, planet)

	// should build nothing, but queue up a mine partially done
	producer := newProducer(testLogger, &rules, planet, player)
	result, err := producer.produce()

	// we should have built a starbase
	assert.Nil(t, err)
	assert.Equal(t, len(result.itemsBuilt), 1)
	assert.Equal(t, result.itemsBuilt[0].queueItemType, QueueItemTypeStarbase)
	require.Len(t, result.starbases, 1)
	assert.Equal(t, starbaseDesign2, result.starbases[0].Tokens[0].design)
	assert.Equal(t, result.starbases[0], planet.Starbase)
	assert.Equal(t, &starbase1, result.replacedStarbase)

}

func Test_production_produceTerraform(t *testing.T) {
	player, planet := newTestPlayerPlanet()

	// build 5 terraform steps
	planet.ProductionQueue = []ProductionQueueItem{
		{Type: QueueItemTypeTerraformEnvironment, Quantity: 5},
	}
	planet.Cargo = Cargo{0, 0, 8, 2500}
	planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 1000, TerraformAmount: Hab{2, 2, 1}}
	planet.BaseHab = Hab{40, 40, 40}
	planet.Hab = Hab{40, 40, 40}

	player.Spec.Terraform[TerraformHabTypeAll] = &TotalTerraform3
	player.Messages = []PlayerMessage{}

	// build 5 terraform steps, make planet better, log some messages
	producer := newProducer(testLogger, &rules, planet, player)
	producer.produce()
	assert.Equal(t, Hab{42, 42, 41}, planet.Hab)
	assert.Equal(t, 0, len(planet.ProductionQueue))

}

func Test_production_produceMineralPackets(t *testing.T) {
	player, planet := newTestPlayerPlanet()

	// build 5 auto factories, leaving them in the queue
	planet.ProductionQueue = []ProductionQueueItem{{Type: QueueItemTypeMixedMineralPacket, Quantity: 1}}
	planet.Cargo = Cargo{100, 100, 100, 2500}
	planet.PacketSpeed = 6
	planet.PacketTargetNum = 1
	planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 100, PlanetStarbaseSpec: PlanetStarbaseSpec{HasMassDriver: true, SafePacketSpeed: 6, BasePacketSpeed: 6}}

	// should build 5 mine, leaving the auto build in the queu
	producer := newProducer(testLogger, &rules, planet, player)
	result, err := producer.produce()
	assert.Nil(t, err)
	assert.Equal(t, Cargo{40, 40, 40, 0}, result.packets)
	assert.Equal(t, 0, len(planet.ProductionQueue))
}

func Test_production_produceScanner(t *testing.T) {
	player, planet := newTestPlayerPlanet()

	// build a scanner
	planet.ProductionQueue = []ProductionQueueItem{{Type: QueueItemTypePlanetaryScanner, Quantity: 1}}
	planet.Cargo = Cargo{1000, 1000, 1000, 100_000}
	planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 1000, MaxPopulation: 1_000_000}
	planet.Scanner = false

	// should build 1 mine, leaving empty queue
	producer := newProducer(testLogger, &rules, planet, player)
	producer.produce()
	assert.True(t, planet.Scanner)
	assert.Equal(t, 0, len(planet.ProductionQueue))
}

func Test_production_allocatePartialBuild(t *testing.T) {
	_, planet := newTestPlayerPlanet()

	type args struct {
		costPerItem Cost
		allocated   Cost
	}
	tests := []struct {
		name string
		args args
		want Cost
	}{
		{"Factory 2kT Germ left", args{costPerItem: Cost{0, 0, 4, 10}, allocated: Cost{100, 100, 2, 100}}, Cost{0, 0, 2, 5}},
		{"Factory 5 resources left", args{costPerItem: Cost{0, 0, 4, 10}, allocated: Cost{100, 100, 100, 5}}, Cost{0, 0, 2, 5}},
		{"Up to 50% due to ironium shortage", args{costPerItem: Cost{10, 20, 30, 40}, allocated: Cost{5, 100, 100, 100}}, Cost{5, 10, 15, 20}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			production := producer{planet: planet}

			if got := production.allocatePartialBuild(tt.args.costPerItem, tt.args.allocated); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Planet.allocatePartialBuild() = %v, want %v", got, tt.want)
			}
		})
	}
}

// a planet with no installations or minerals, room for 100 of each installation
func newProductionTestPlanet(resources int, minerals Mineral) (*Player, *Planet) {
	player, planet := newTestPlayerPlanet()
	planet.Mines, planet.Factories, planet.Defenses = 0, 0, 0
	planet.Cargo = Cargo{Colonists: 1000}.AddMineral(minerals)
	planet.Spec = PlanetSpec{
		ResourcesPerYearAvailable: resources,
		MaxPossibleMines:          100,
		MaxPossibleFactories:      100,
		MaxDefenses:               100,
		MaxPopulation:             1_000_000,
	}
	return player, planet
}

// what a year of production did to the planet
type productionTestResult struct {
	alchemy   int // kT of each mineral made by alchemy
	mines     int
	factories int
	defenses  int
	ships     int
	minerals  Mineral
	leftover  int // resources left over for research
	queue     []ProductionQueueItem
	messages  []PlayerMessageType
}

// run one year of production and return what happened
func runProduction(t *testing.T, rules *Rules, player *Player, planet *Planet) productionTestResult {
	t.Helper()
	producer := newProducer(testLogger, rules, planet, player)
	result, err := producer.produce()
	require.NoError(t, err)

	got := productionTestResult{
		alchemy:   result.alchemy,
		mines:     planet.Mines,
		factories: planet.Factories,
		defenses:  planet.Defenses,
		minerals:  planet.Cargo.ToMineral(),
		leftover:  result.leftoverResources,
	}
	if len(planet.ProductionQueue) > 0 {
		got.queue = planet.ProductionQueue
	}
	for _, token := range result.tokens {
		got.ships += token.Quantity
	}
	for _, message := range result.messages {
		got.messages = append(got.messages, message.Type)
	}
	return got
}

// Costs with the default rules and race:
//
//	alchemy:  100 resources for 1kT of each mineral
//	factory:  10 resources, 4 germanium
//	mine:     5 resources
//	defenses: 15 resources, 5 of each mineral
func Test_production_alchemy(t *testing.T) {
	alchemy := ProductionQueueItem{Type: QueueItemTypeMineralAlchemy, Quantity: 1}
	autoAlchemy := ProductionQueueItem{Type: QueueItemTypeAutoMineralAlchemy, Quantity: 1}
	factory := ProductionQueueItem{Type: QueueItemTypeFactory, Quantity: 1}
	autoFactory := ProductionQueueItem{Type: QueueItemTypeAutoFactories, Quantity: 1}
	mine := ProductionQueueItem{Type: QueueItemTypeMine, Quantity: 1}

	// a factory started by auto factories, or alchemy started by auto alchemy at the end of the queue
	partialFactory := func(allocated Cost) ProductionQueueItem {
		return ProductionQueueItem{Type: QueueItemTypeFactory, Quantity: 1, Allocated: allocated}
	}
	partialAutoAlchemy := func(resources int) ProductionQueueItem {
		return ProductionQueueItem{Type: QueueItemTypeMineralAlchemy, Quantity: 1, Allocated: Cost{Resources: resources}}
	}
	// leftover resources from converting minerals for an order, put toward alchemy next year
	partialAlchemy := func(resources int) ProductionQueueItem {
		return ProductionQueueItem{Type: QueueItemTypeMineralAlchemy, Quantity: 1, Allocated: Cost{Resources: resources}, index: -1}
	}

	tests := []struct {
		name      string
		resources int
		minerals  Mineral
		factories int // factories already built
		queue     []ProductionQueueItem
		want      productionTestResult
	}{
		// concrete alchemy
		{
			name:      "alchemy makes minerals for the next order",
			resources: 105, // alchemy, mine
			queue:     []ProductionQueueItem{alchemy, mine},
			want:      productionTestResult{alchemy: 1, mines: 1, minerals: Mineral{1, 1, 1}},
		},
		{
			name:      "each alchemy order makes its own minerals",
			resources: 200, // 2 alchemy
			queue:     []ProductionQueueItem{alchemy, alchemy},
			want:      productionTestResult{alchemy: 2, minerals: Mineral{2, 2, 2}},
		},

		// auto alchemy before another order
		{
			name:      "auto alchemy makes the germanium a factory needs",
			resources: 410, // 4 alchemy, factory
			queue:     []ProductionQueueItem{autoAlchemy, factory},
			want:      productionTestResult{alchemy: 4, factories: 1, minerals: Mineral{4, 4, 0}},
		},
		{
			name:      "auto alchemy does nothing when we have the minerals",
			resources: 100,
			minerals:  Mineral{Germanium: 4},
			queue:     []ProductionQueueItem{autoAlchemy, factory},
			want:      productionTestResult{factories: 1, leftover: 90},
		},
		{
			name:      "auto alchemy only makes what a partly built factory still needs",
			resources: 205, // 2 alchemy, the rest of the factory
			queue: []ProductionQueueItem{autoAlchemy,
				{Type: QueueItemTypeFactory, Quantity: 1, Allocated: Cost{Germanium: 2, Resources: 5}}},
			want: productionTestResult{alchemy: 2, factories: 1, minerals: Mineral{2, 2, 0}},
		},
		{
			name:      "auto alchemy makes up every mineral we're short on",
			resources: 515, // 5 alchemy, defense
			minerals:  Mineral{0, 3, 1},
			queue:     []ProductionQueueItem{autoAlchemy, {Type: QueueItemTypeDefenses, Quantity: 1}},
			want:      productionTestResult{alchemy: 5, defenses: 1, minerals: Mineral{0, 3, 1}},
		},
		{
			name:      "auto alchemy only applies to the next order",
			resources: 115, // factory, mine, 100 left over
			minerals:  Mineral{Germanium: 4},
			queue:     []ProductionQueueItem{autoAlchemy, factory, autoFactory, mine},
			want:      productionTestResult{factories: 1, mines: 1, leftover: 100, queue: []ProductionQueueItem{autoFactory}},
		},
		{
			name:      "auto alchemy stays in the queue with an auto order",
			resources: 825, // 8 alchemy, 2 factories, mine
			queue:     []ProductionQueueItem{autoAlchemy, {Type: QueueItemTypeAutoFactories, Quantity: 2}, mine},
			want: productionTestResult{alchemy: 8, factories: 2, mines: 1, minerals: Mineral{8, 8, 0},
				queue: []ProductionQueueItem{autoAlchemy, {Type: QueueItemTypeAutoFactories, Quantity: 2}}},
		},
		{
			name:      "auto alchemy doesn't convert when we're short on resources",
			resources: 5,
			minerals:  Mineral{Germanium: 4},
			queue:     []ProductionQueueItem{autoAlchemy, factory},
			want: productionTestResult{minerals: Mineral{Germanium: 2},
				queue: []ProductionQueueItem{autoAlchemy, {Type: QueueItemTypeFactory, Quantity: 1, Allocated: Cost{Germanium: 2, Resources: 5}}}},
		},
		{
			name:      "auto factory short on resources starts a partial factory",
			resources: 5,
			minerals:  Mineral{Germanium: 4},
			queue:     []ProductionQueueItem{autoAlchemy, autoFactory, mine},
			want: productionTestResult{minerals: Mineral{Germanium: 2},
				queue: []ProductionQueueItem{partialFactory(Cost{Germanium: 2, Resources: 5}), autoAlchemy, autoFactory, mine}},
		},
		{
			name:      "not enough resources to finish an alchemy saves it for next year",
			resources: 150, // 1 alchemy, half of another
			queue:     []ProductionQueueItem{autoAlchemy, factory, mine},
			want: productionTestResult{alchemy: 1, minerals: Mineral{1, 1, 1},
				queue: []ProductionQueueItem{partialAlchemy(50), autoAlchemy, factory, mine}},
		},
		{
			name:      "next year we finish the alchemy, then the factory and mine",
			resources: 265, // the rest of the alchemy, 2 alchemy, factory, mine
			minerals:  Mineral{1, 1, 1},
			queue:     []ProductionQueueItem{partialAlchemy(50), autoAlchemy, factory, mine},
			want:      productionTestResult{alchemy: 3, factories: 1, mines: 1, minerals: Mineral{4, 4, 0}},
		},
		{
			name:      "not enough resources to finish an alchemy for an auto factory saves it for next year",
			resources: 150, // 1 alchemy, half of another
			queue:     []ProductionQueueItem{autoAlchemy, autoFactory, mine},
			want: productionTestResult{alchemy: 1, minerals: Mineral{1, 1, 1},
				queue: []ProductionQueueItem{partialAlchemy(50), autoAlchemy, autoFactory, mine}},
		},
		{
			name:      "next year we finish the alchemy, then the auto factory and mine",
			resources: 265, // the rest of the alchemy, 2 alchemy, factory, mine
			minerals:  Mineral{1, 1, 1},
			queue:     []ProductionQueueItem{partialAlchemy(50), autoAlchemy, autoFactory, mine},
			want: productionTestResult{alchemy: 3, factories: 1, mines: 1, minerals: Mineral{4, 4, 0},
				queue: []ProductionQueueItem{autoAlchemy, autoFactory}},
		},
		{
			name:      "a factory with no room is removed along with its auto alchemy",
			resources: 105,
			factories: 100,
			queue:     []ProductionQueueItem{autoAlchemy, factory, mine},
			want: productionTestResult{factories: 100, mines: 1, leftover: 100,
				messages: []PlayerMessageType{PlayerMessagePlanetBuiltBeyondMaximum}},
		},
		{
			name:      "an auto factory we can't build is skipped and keeps its auto alchemy",
			resources: 105,
			factories: 100,
			queue:     []ProductionQueueItem{autoAlchemy, autoFactory, mine},
			want: productionTestResult{factories: 100, mines: 1, leftover: 100,
				queue: []ProductionQueueItem{autoAlchemy, autoFactory}},
		},
		{
			name:      "auto factory short on any germanium waits, even when resources are shorter",
			resources: 2,
			minerals:  Mineral{Germanium: 3},
			queue:     []ProductionQueueItem{autoFactory, mine},
			want: productionTestResult{minerals: Mineral{Germanium: 3},
				queue: []ProductionQueueItem{autoFactory, {Type: QueueItemTypeMine, Quantity: 1, Allocated: Cost{Resources: 2}}}},
		},
		{
			name:      "auto factory with auto alchemy puts resources toward alchemy, even when resources are shorter",
			resources: 5,
			minerals:  Mineral{Germanium: 3},
			queue:     []ProductionQueueItem{autoAlchemy, autoFactory, mine},
			want: productionTestResult{minerals: Mineral{Germanium: 3},
				queue: []ProductionQueueItem{partialAlchemy(5), autoAlchemy, autoFactory, mine}},
		},
		{
			name:  "auto alchemy does nothing without resources",
			queue: []ProductionQueueItem{autoAlchemy, factory},
			want:  productionTestResult{queue: []ProductionQueueItem{autoAlchemy, factory}},
		},

		// auto alchemy at the end of the queue
		{
			name:      "auto alchemy at the end of the queue converts all our resources",
			resources: 350, // 3 alchemy, half of another
			queue:     []ProductionQueueItem{autoAlchemy},
			want: productionTestResult{alchemy: 3, minerals: Mineral{3, 3, 3},
				queue: []ProductionQueueItem{partialAutoAlchemy(50), autoAlchemy}},
		},
		{
			name:      "next year we finish the alchemy",
			resources: 50,
			queue:     []ProductionQueueItem{partialAutoAlchemy(50), autoAlchemy},
			want:      productionTestResult{alchemy: 1, minerals: Mineral{1, 1, 1}, queue: []ProductionQueueItem{autoAlchemy}},
		},
		{
			name:  "auto alchemy at the end of the queue does nothing without resources",
			queue: []ProductionQueueItem{autoAlchemy},
			want:  productionTestResult{queue: []ProductionQueueItem{autoAlchemy}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player, planet := newProductionTestPlanet(tt.resources, tt.minerals)
			planet.Factories = tt.factories
			planet.ProductionQueue = tt.queue

			assert.Equal(t, tt.want, runProduction(t, &rules, player, planet))
		})
	}
}

func Test_production_alchemyCost(t *testing.T) {
	tests := []struct {
		name        string
		alchemyCost int
		lrts        Bitmask
		resources   int
		want        productionTestResult
	}{
		{
			name:        "custom alchemy cost",
			alchemyCost: 115,
			resources:   470, // 4 alchemy at 115, factory
			want:        productionTestResult{alchemy: 4, factories: 1, minerals: Mineral{4, 4, 0}},
		},
		{
			name:        "custom alchemy cost with MA",
			alchemyCost: 115,
			lrts:        Bitmask(MA),
			resources:   170, // 4 alchemy at 40 (MA takes 75 off), factory
			want:        productionTestResult{alchemy: 4, factories: 1, minerals: Mineral{4, 4, 0}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules := rules
			rules.MineralAlchemyCost = tt.alchemyCost
			player, planet := newProductionTestPlanet(tt.resources, Mineral{})
			player.Race.LRTs |= tt.lrts
			player.Race.Spec = ComputeRaceSpec(&player.Race, &rules)
			planet.ProductionQueue = []ProductionQueueItem{
				{Type: QueueItemTypeAutoMineralAlchemy, Quantity: 1},
				{Type: QueueItemTypeFactory, Quantity: 1},
			}

			assert.Equal(t, tt.want, runProduction(t, &rules, player, planet))
		})
	}
}

// concrete orders beyond what the planet can hold are reduced to the max
func Test_production_ordersBeyondMaximum(t *testing.T) {
	alchemy := ProductionQueueItem{Type: QueueItemTypeMineralAlchemy, Quantity: 1}
	beyondMaximum := []PlayerMessageType{PlayerMessagePlanetBuiltBeyondMaximum}

	tests := []struct {
		name      string
		resources int
		minerals  Mineral
		queue     []ProductionQueueItem
		want      productionTestResult
	}{
		{
			name:      "5 mines with room for 1, then alchemy",
			resources: 105, // mine, alchemy
			queue:     []ProductionQueueItem{{Type: QueueItemTypeMine, Quantity: 5}, alchemy},
			want:      productionTestResult{mines: 1, alchemy: 1, minerals: Mineral{1, 1, 1}, messages: beyondMaximum},
		},
		{
			name:      "5 factories with room for 1, then alchemy",
			resources: 110, // factory, alchemy
			minerals:  Mineral{Germanium: 4},
			queue:     []ProductionQueueItem{{Type: QueueItemTypeFactory, Quantity: 5}, alchemy},
			want:      productionTestResult{factories: 1, alchemy: 1, minerals: Mineral{1, 1, 1}, messages: beyondMaximum},
		},
		{
			name:      "5 defenses with room for 1, then alchemy",
			resources: 115, // defense, alchemy
			minerals:  Mineral{5, 5, 5},
			queue:     []ProductionQueueItem{{Type: QueueItemTypeDefenses, Quantity: 5}, alchemy},
			want:      productionTestResult{defenses: 1, alchemy: 1, minerals: Mineral{1, 1, 1}, messages: beyondMaximum},
		},
		{
			name:      "5 terraforms with room for 1, then alchemy",
			resources: 200, // terraform, alchemy
			queue:     []ProductionQueueItem{{Type: QueueItemTypeTerraformEnvironment, Quantity: 5}, alchemy},
			want:      productionTestResult{alchemy: 1, minerals: Mineral{1, 1, 1}, messages: beyondMaximum},
		},
		{
			name:      "a second terraform is canceled once the first leaves nothing to terraform",
			resources: 200, // terraform, 100 left over
			queue: []ProductionQueueItem{
				{Type: QueueItemTypeTerraformEnvironment, Quantity: 1},
				{Type: QueueItemTypeTerraformEnvironment, Quantity: 1},
			},
			want: productionTestResult{leftover: 100, messages: beyondMaximum},
		},
		{
			name:      "a reduced order keeps what we already spent on it",
			resources: 3,
			queue:     []ProductionQueueItem{{Type: QueueItemTypeMine, Quantity: 5, Allocated: Cost{Resources: 1}}, alchemy},
			want: productionTestResult{messages: beyondMaximum,
				queue: []ProductionQueueItem{{Type: QueueItemTypeMine, Quantity: 1, Allocated: Cost{Resources: 4}}, alchemy}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player, planet := newProductionTestPlanet(tt.resources, tt.minerals)
			// room for one of everything
			planet.Spec.MaxPossibleMines, planet.Spec.MaxPossibleFactories, planet.Spec.MaxDefenses = 1, 1, 1
			planet.Spec.TerraformAmount = Hab{Grav: 1}
			planet.Hab, planet.BaseHab = Hab{49, 50, 50}, Hab{49, 50, 50}
			player.Spec.Terraform[TerraformHabTypeAll] = &TotalTerraform3
			planet.ProductionQueue = tt.queue

			assert.Equal(t, tt.want, runProduction(t, &rules, player, planet))
		})
	}
}

// ships are only built at a starbase with a dock big enough for them
func Test_production_ships(t *testing.T) {
	invalidShip := []PlayerMessageType{PlayerMessagePlanetBuiltInvalidShip}

	tests := []struct {
		name         string
		starbaseHull string // no starbase if empty
		want         productionTestResult
	}{
		{
			name:         "a space station builds a scout",
			starbaseHull: SpaceStation.Name,
			want:         productionTestResult{ships: 1, minerals: Mineral{96, 98, 96}, leftover: 91}, // a scout costs 4, 2, 4, 9
		},
		{
			name: "a scout with no starbase is canceled",
			want: productionTestResult{minerals: Mineral{100, 100, 100}, leftover: 100, messages: invalidShip},
		},
		{
			name:         "a scout too big for an orbital fort's dock is canceled",
			starbaseHull: OrbitalFort.Name,
			want:         productionTestResult{minerals: Mineral{100, 100, 100}, leftover: 100, messages: invalidShip},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player, planet := newProductionTestPlanet(100, Mineral{100, 100, 100})
			scout := NewShipDesign(player.Num, 1).WithHull(Scout.Name).WithSpec(&rules, player)
			player.Designs = append(player.Designs, scout)
			if tt.starbaseHull != "" {
				starbaseDesign := NewShipDesign(player.Num, 2).WithHull(tt.starbaseHull).WithSpec(&rules, player)
				player.Designs = append(player.Designs, starbaseDesign)
				starbase := newStarbase(player, planet, starbaseDesign, "Starbase")
				starbase.Spec = ComputeFleetSpec(&rules, player, &starbase)
				planet.Starbase = &starbase
				planet.Spec.PlanetStarbaseSpec = computePlanetStarbaseSpec(planet)
			}
			planet.ProductionQueue = []ProductionQueueItem{{Type: QueueItemTypeShipToken, Quantity: 1, DesignNum: scout.Num, design: scout}}

			assert.Equal(t, tt.want, runProduction(t, &rules, player, planet))
		})
	}
}

// a new starbase goes on the planet as soon as it's built, so later items use it
func Test_production_starbaseAtItsQueuePosition(t *testing.T) {
	player, planet := newProductionTestPlanet(10_000, Mineral{10_000, 10_000, 10_000})
	newBase := func(num int, slots ...ShipDesignSlot) *ShipDesign {
		design := NewShipDesign(player.Num, num).WithHull(SpaceStation.Name).WithSlots(slots).WithSpec(&rules, player)
		player.Designs = append(player.Designs, design)
		return design
	}
	emptyBase := newBase(1)
	laserBase := newBase(2, ShipDesignSlot{HullComponent: Laser.Name, HullSlotIndex: 2, Quantity: 1})
	armoredLaserBase := newBase(3,
		ShipDesignSlot{HullComponent: Laser.Name, HullSlotIndex: 2, Quantity: 1},
		ShipDesignSlot{HullComponent: Tritanium.Name, HullSlotIndex: 4, Quantity: 1},
	)
	originalStarbase := newStarbase(player, planet, emptyBase, "Starbase")
	originalStarbase.Spec = ComputeFleetSpec(&rules, player, &originalStarbase)
	planet.Starbase = &originalStarbase
	planet.Spec.PlanetStarbaseSpec = computePlanetStarbaseSpec(planet)

	// upgrade twice in one year
	planet.ProductionQueue = []ProductionQueueItem{
		{Type: QueueItemTypeStarbase, Quantity: 1, DesignNum: laserBase.Num, design: laserBase},
		{Type: QueueItemTypeStarbase, Quantity: 1, DesignNum: armoredLaserBase.Num, design: armoredLaserBase},
	}
	producer := newProducer(testLogger, &rules, planet, player)
	result, err := producer.produce()
	require.NoError(t, err)

	// the second upgrade is priced against the laser base, not the empty base we started with
	costCalculator := NewCostCalculator()
	toLaserBase, err := costCalculator.StarbaseUpgradeCost(&rules, player.TechLevels, player.Race.Spec, emptyBase, laserBase)
	require.NoError(t, err)
	toArmoredLaserBase, err := costCalculator.StarbaseUpgradeCost(&rules, player.TechLevels, player.Race.Spec, laserBase, armoredLaserBase)
	require.NoError(t, err)
	assert.Equal(t, 10_000-toLaserBase.Resources-toArmoredLaserBase.Resources, result.leftoverResources)

	require.Len(t, result.starbases, 2)
	assert.Equal(t, armoredLaserBase, planet.Starbase.Tokens[0].design)
	assert.Equal(t, result.starbases[1], planet.Starbase)
	assert.Equal(t, &originalStarbase, result.replacedStarbase)
}

func Test_production_packetsAfterMassDriverBase(t *testing.T) {
	player, planet := newProductionTestPlanet(10_000, Mineral{10_000, 10_000, 10_000})
	player.Race.PRT = PP
	player.Race.Spec = ComputeRaceSpec(&player.Race, &rules)
	player.TechLevels = TechLevel{Energy: 4}
	massDriverBase := NewShipDesign(player.Num, 1).WithHull(SpaceStation.Name).
		WithSlots([]ShipDesignSlot{{HullComponent: MassDriver5.Name, HullSlotIndex: 1, Quantity: 1}}).
		WithSpec(&rules, player)
	player.Designs = append(player.Designs, massDriverBase)
	planet.PacketTargetNum = 1

	// no starbase yet, the packet is built by the base we build first
	planet.ProductionQueue = []ProductionQueueItem{
		{Type: QueueItemTypeStarbase, Quantity: 1, DesignNum: massDriverBase.Num, design: massDriverBase},
		{Type: QueueItemTypeMixedMineralPacket, Quantity: 1},
	}
	producer := newProducer(testLogger, &rules, planet, player)
	result, err := producer.produce()
	require.NoError(t, err)

	assert.Empty(t, result.messages)
	assert.Equal(t, Cargo{25, 25, 25, 0}, result.packets) // PP mixed packets carry 25kT of each
}

// packets carry the race's packet payload, regardless of what they cost to launch
func Test_production_packetPayloads(t *testing.T) {
	tests := []struct {
		name     string
		prt      PRT
		itemType QueueItemType
		want     Cargo
	}{
		{"mixed packet", JoaT, QueueItemTypeMixedMineralPacket, Cargo{40, 40, 40, 0}},
		{"ironium packet", JoaT, QueueItemTypeIroniumMineralPacket, Cargo{Ironium: 100}},
		{"PP mixed packet", PP, QueueItemTypeMixedMineralPacket, Cargo{25, 25, 25, 0}},
		{"PP germanium packet", PP, QueueItemTypeGermaniumMineralPacket, Cargo{Germanium: 70}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player, planet := newProductionTestPlanet(1000, Mineral{1000, 1000, 1000})
			player.Race.PRT = tt.prt
			player.Race.Spec = ComputeRaceSpec(&player.Race, &rules)
			planet.Spec.PlanetStarbaseSpec = PlanetStarbaseSpec{HasMassDriver: true, SafePacketSpeed: 6, BasePacketSpeed: 6}
			planet.PacketTargetNum = 1
			planet.ProductionQueue = []ProductionQueueItem{{Type: tt.itemType, Quantity: 1}}

			producer := newProducer(testLogger, &rules, planet, player)
			result, err := producer.produce()
			require.NoError(t, err)
			assert.Equal(t, tt.want, result.packets)
		})
	}
}

// a genesis device destroys everything on the planet and rerolls it right away
func Test_production_genesisDevice(t *testing.T) {
	tests := []struct {
		name       string
		estimating bool
	}{
		{"rerolls the planet", false},
		{"estimates don't know the reroll, so keep the environment", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules := NewRulesWithSeed(1)
			player, planet := newProductionTestPlanet(5000, Mineral{100, 100, 100})
			planet.Mines, planet.Factories, planet.Defenses, planet.Scanner = 10, 10, 10, true
			planet.Hab, planet.BaseHab, planet.TerraformedAmount = Hab{48, 50, 50}, Hab{45, 50, 50}, Hab{3, 0, 0}
			planet.ProductionQueue = []ProductionQueueItem{{Type: QueueItemTypeGenesisDevice, Quantity: 1}}

			producer := newProducer(testLogger, &rules, planet, player)
			producer.estimating = tt.estimating
			result, err := producer.produce()
			require.NoError(t, err)

			assert.True(t, result.reset)
			assert.Equal(t, []int{0, 0, 0}, []int{planet.Mines, planet.Factories, planet.Defenses})
			assert.False(t, planet.Scanner)
			assert.Equal(t, Mineral{100, 100, 100}, planet.Cargo.ToMineral()) // surface minerals are kept
			if tt.estimating {
				assert.Equal(t, Hab{48, 50, 50}, planet.Hab)
			} else {
				assert.NotEqual(t, Hab{48, 50, 50}, planet.Hab)
				assert.Equal(t, planet.Hab, planet.BaseHab)
				assert.Equal(t, Hab{}, planet.TerraformedAmount)
			}
		})
	}
}
