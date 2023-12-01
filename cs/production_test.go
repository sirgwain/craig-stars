package cs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_production_produce(t *testing.T) {
	player, planet := newTestPlayerPlanet()

	// build 1 mine
	planet.ProductionQueue = []ProductionQueueItem{{Type: QueueItemTypeMine, Quantity: 1}}
	planet.Cargo = Cargo{10, 20, 30, 2500}
	planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 100, MaxPossibleMines: 100}
	planet.Mines = 0

	// should build 1 mine, leaving empty queue
	producer := newProducer(planet, player)
	producer.produce()
	assert.Equal(t, 1, planet.Mines)
	assert.Equal(t, 0, len(planet.ProductionQueue))

	// build 5 auto mines, leaving them in the queue
	planet.ProductionQueue = []ProductionQueueItem{{Type: QueueItemTypeAutoMines, Quantity: 5}}
	planet.Cargo = Cargo{10, 20, 30, 2500}
	planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 100, MaxMines: 100}
	planet.Mines = 0
	player.Messages = []PlayerMessage{}

	// should build 5 mine, leaving the auto build in the queu
	producer = newProducer(planet, player)
	producer.produce()
	assert.Equal(t, 5, planet.Mines)
	assert.Equal(t, 1, len(planet.ProductionQueue))
	assert.Equal(t, QueueItemTypeAutoMines, planet.ProductionQueue[0].Type)
	assert.Equal(t, 5, planet.ProductionQueue[0].Quantity)
}

func Test_production_produce2(t *testing.T) {
	player, planet := newTestPlayerPlanet()

	// build 5 auto factories, leaving them in the queue
	planet.ProductionQueue = []ProductionQueueItem{{Type: QueueItemTypeAutoFactories, Quantity: 5}}
	planet.Cargo = Cargo{10, 20, 30, 2500}
	planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 100, MaxFactories: 100}
	planet.Factories = 0
	player.Messages = []PlayerMessage{}

	// should build 5 mine, leaving the auto build in the queu
	producer := newProducer(planet, player)
	producer.produce()
	assert.Equal(t, 5, planet.Factories)
	assert.Equal(t, 1, len(planet.ProductionQueue))
	assert.Equal(t, QueueItemTypeAutoFactories, planet.ProductionQueue[0].Type)
	assert.Equal(t, 5, planet.ProductionQueue[0].Quantity)

}

func Test_production_produce3(t *testing.T) {
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
	planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 100, MaxFactories: 100, MaxMines: 100}
	planet.Factories = 0
	player.Messages = []PlayerMessage{}

	// should build 2 factories and 5 mines, leaving the auto builds in the queue
	producer := newProducer(planet, player)
	producer.produce()
	assert.Equal(t, 2, planet.Factories)
	assert.Equal(t, 5, planet.Mines)
	assert.Equal(t, 2, len(planet.ProductionQueue))
	assert.Equal(t, QueueItemTypeAutoFactories, planet.ProductionQueue[0].Type)
	assert.Equal(t, 5, planet.ProductionQueue[0].Quantity)
	assert.Equal(t, QueueItemTypeAutoMines, planet.ProductionQueue[1].Type)
	assert.Equal(t, 5, planet.ProductionQueue[1].Quantity)
}

func Test_production_produce4(t *testing.T) {
	player, planet := newTestPlayerPlanet()

	// build 2/5 auto factories and one partial mine
	planet.ProductionQueue = []ProductionQueueItem{
		{Type: QueueItemTypeAutoFactories, Quantity: 5},
		{Type: QueueItemTypeAutoMines, Quantity: 5},
	}
	planet.Cargo = Cargo{0, 0, 8, 2500}
	planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 10*2 + 8, MaxFactories: 100, MaxMines: 100}
	planet.Factories = 0
	player.Messages = []PlayerMessage{}

	// should build 2 factories, 1 mine and one partial mine, leaving the auto builds in the queue
	producer := newProducer(planet, player)
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

func Test_production_produce5(t *testing.T) {
	player, planet := newTestPlayerPlanet()

	// build 100 auto defenses when we have 90 already and one partial in the queue
	planet.ProductionQueue = []ProductionQueueItem{
		{Type: QueueItemTypeDefenses, Quantity: 1, Allocated: Cost{5, 5, 5, 14}},
		{Type: QueueItemTypeAutoDefenses, Quantity: 100},
	}
	planet.Cargo = Cargo{5000, 5000, 5000, 1_000_000}
	planet.Defenses = 90
	planet.Spec = computePlanetSpec(&rules, player, planet)
	player.Messages = []PlayerMessage{}

	// should end up with 100 defenses and the auto defenses still in the queue
	producer := newProducer(planet, player)
	producer.produce()
	assert.Equal(t, 100, planet.Defenses)
	assert.Equal(t, 1, len(planet.ProductionQueue))
	assert.Equal(t, QueueItemTypeAutoDefenses, planet.ProductionQueue[0].Type)
	assert.Equal(t, 100, planet.ProductionQueue[0].Quantity)
}

func Test_production_produce6(t *testing.T) {
	player, planet := newTestPlayerPlanet()

	// auto build up to the max
	planet.ProductionQueue = []ProductionQueueItem{
		{Type: QueueItemTypeAutoFactories, Quantity: 10},
		{Type: QueueItemTypeAutoMines, Quantity: 10},
	}
	planet.Cargo = Cargo{1000, 1000, 1000, 2500}
	planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 1000, MaxFactories: 10, MaxMines: 10}
	planet.Factories = 9
	planet.Mines = 0
	player.Messages = []PlayerMessage{}

	// should build 1 factories, 10 mines and have leftover for research
	producer := newProducer(planet, player)
	result := producer.produce()
	assert.Equal(t, 10, planet.Factories)
	assert.Equal(t, 10, planet.Mines)
	assert.Equal(t, 2, len(planet.ProductionQueue))
	assert.Equal(t, QueueItemTypeAutoFactories, planet.ProductionQueue[0].Type)
	assert.Equal(t, 10, planet.ProductionQueue[0].Quantity)
	assert.Equal(t, QueueItemTypeAutoMines, planet.ProductionQueue[1].Type)
	assert.Equal(t, 10, planet.ProductionQueue[1].Quantity)
	assert.Equal(t, 940, result.leftoverResources)

}

func Test_production_produce7(t *testing.T) {
	player, planet := newTestPlayerPlanet()

	// auto build factories, but we have no mines or minerals, so they'll never build
	// until we build mines and mine a bit
	planet.ProductionQueue = []ProductionQueueItem{
		{Type: QueueItemTypeAutoFactories, Quantity: 100},
		{Type: QueueItemTypeAutoMines, Quantity: 100},
	}
	planet.Cargo = Cargo{0, 0, 0, 25}
	planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 2, MaxFactories: 10, MaxMines: 10}
	player.Messages = []PlayerMessage{}

	// should build nothing, but queue up a mine partially done
	producer := newProducer(planet, player)
	producer.produce()
	assert.Equal(t, 0, planet.Factories)
	assert.Equal(t, 0, planet.Mines)
	assert.Equal(t, 3, len(planet.ProductionQueue))
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

func Test_production_produce8(t *testing.T) {
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
	planet.Spec = computePlanetSpec(&rules, player, planet)

	// should build nothing, but queue up a mine partially done
	producer := newProducer(planet, player)
	producer.produce()
	assert.Equal(t, 20, planet.Factories)
	assert.Equal(t, 11, planet.Mines)
	assert.Equal(t, 4, len(planet.ProductionQueue))
	assert.Equal(t, QueueItemTypeFactory, planet.ProductionQueue[0].Type)
	assert.Equal(t, 1, planet.ProductionQueue[0].Quantity)
	assert.Equal(t, Cost{Resources: 7}, planet.ProductionQueue[0].Allocated)
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

func Test_production_produce9(t *testing.T) {
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
	planet.Spec = computePlanetSpec(&rules, player, planet)

	// should build nothing, but queue up a mine partially done
	producer := newProducer(planet, player)
	producer.produce()

	// don't go negative
	assert.Equal(t, planet.Cargo, planet.Cargo.MinZero())

}

func Test_production_produceStarbaseUpgrade(t *testing.T) {
	player, planet := newTestPlayerPlanet()

	// create a new starbase
	starbaseDesign1 := NewShipDesign(player, 2).WithHull(SpaceStation.Name).WithSpec(&rules, player)
	starbase1 := newStarbase(player, planet,
		starbaseDesign1,
		"Starbase",
	)
	starbaseDesign2 := NewShipDesign(player, 3).WithHull(SpaceStation.Name).
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
	planet.Spec = computePlanetSpec(&rules, player, planet)

	// should build nothing, but queue up a mine partially done
	producer := newProducer(planet, player)
	result := producer.produce()

	// we should have built a starbase
	assert.Equal(t, len(result.itemsBuilt), 1)
	assert.Equal(t, result.itemsBuilt[0].queueItemType, QueueItemTypeStarbase)
	assert.Equal(t, result.starbase, starbaseDesign2)

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
	producer := newProducer(planet, player)
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
	producer := newProducer(planet, player)
	result := producer.produce()
	assert.Equal(t, 1, len(result.packets))
	assert.Equal(t, Cargo{40, 40, 40, 0}, result.packets[0])
	assert.Equal(t, 0, len(planet.ProductionQueue))
}
