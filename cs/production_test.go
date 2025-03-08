package cs

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func Test_production_produce(t *testing.T) {
	// TODO: Add tests for mineral alchemy & auto mineral alchemy
	t.Run("1 Mine", func(t *testing.T) {
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
		planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 100, MaxMines: 100, MaxPopulation: 1_000_000}
		planet.Mines = 0
		player.Messages = []PlayerMessage{}

		// should build 5 mines, leaving the auto build in the queue
		producer = newProducer(testLogger, &rules, planet, player)
		producer.produce()
		assert.Equal(t, 5, planet.Mines)
		assert.Equal(t, 1, len(planet.ProductionQueue))
		assert.Equal(t, QueueItemTypeAutoMines, planet.ProductionQueue[0].Type)
		assert.Equal(t, 5, planet.ProductionQueue[0].Quantity)
	})

	t.Run("Auto Factories", func(t *testing.T) {
		player, planet := newTestPlayerPlanet()

		// build 5 auto factories, leaving them in the queue
		planet.ProductionQueue = []ProductionQueueItem{{Type: QueueItemTypeAutoFactories, Quantity: 5}}
		planet.Cargo = Cargo{10, 20, 30, 2500}
		planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 100, MaxFactories: 100, MaxPopulation: 1_000_000}
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
	})

	t.Run("Factories + Mines, low germ", func(t *testing.T) {
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
		planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 100, MaxFactories: 100, MaxMines: 100, MaxPopulation: 1_000_000}
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
	})

	t.Run("Refund invalid items", func(t *testing.T) {
		player, planet := newTestPlayerPlanet()

		// make defenses an even 10 in all for cost
		rCopy := rules
		rCopy.DefenseCost = Cost{10, 10, 10, 10}
		rCopy.PlanetaryScannerCost = Cost{999, 999, 999, 999}

		// exactly enough to finish 10 defenses
		planet.Cargo = Cargo{100, 100, 100, 1000}
		planet.Defenses = 90
		planet.ContributesOnlyLeftoverToResearch = true

		player.Race = *player.Race.WithSpec(&rCopy)
		planet.Spec = computePlanetSpec(&rCopy, player, planet)
		player.Spec = computePlayerSpec(player, &rCopy, []*Planet{planet})

		// The 100 auto defenses should remove the auto defense from the queue
		planet.ProductionQueue = []ProductionQueueItem{
			{Type: QueueItemTypeAutoDefenses, Quantity: 100},
			{Type: QueueItemTypeDefenses, Quantity: 90, Allocated: Cost{5, 5, 5, 5}},
			// super expensive scanner to soak up leftover minerals
			{Type: QueueItemTypePlanetaryScanner, Quantity: 1},
		}
		planet.Spec = computePlanetSpec(&rCopy, player, planet)
		player.Messages = []PlayerMessage{}

		// should end up with 100 defenses, with auto defenses still in the queue;
		// scanner should soak up allocated resources from first defense item
		wantQueue := []ProductionQueueItem{
			{Type: QueueItemTypeAutoDefenses, Quantity: 100},
			{Type: QueueItemTypePlanetaryScanner, Quantity: 1, Allocated: Cost{5, 5, 5, 5}},
		}
		producer := newProducer(testLogger, &rCopy, planet, player)
		producer.produce()
		assert.Equal(t, 100, planet.Defenses)
		test.CompareAsJSON(t, planet.ProductionQueue, wantQueue)
	})

	t.Run("Auto Factories + Mines", func(t *testing.T) {
		player, planet := newTestPlayerPlanet()

		// auto build up to the max
		planet.ProductionQueue = []ProductionQueueItem{
			{Type: QueueItemTypeAutoFactories, Quantity: 10},
			{Type: QueueItemTypeAutoMines, Quantity: 10},
		}
		planet.Cargo = Cargo{1000, 1000, 1000, 100}
		planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 1000, MaxFactories: 10, MaxMines: 10, MaxPopulation: 1_000_000}
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

	})

	t.Run("Partial mines", func(t *testing.T) {
		player, planet := newTestPlayerPlanet()

		// auto build factories, but we have no mines or minerals, so they'll never build
		// until we build mines and mine a bit
		planet.ProductionQueue = []ProductionQueueItem{
			{Type: QueueItemTypeAutoFactories, Quantity: 100},
			{Type: QueueItemTypeAutoMines, Quantity: 100},
		}
		planet.Cargo = Cargo{0, 0, 0, 25}
		planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 2, MaxFactories: 10, MaxMines: 10, MaxPopulation: 1_000_000}
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

	})

	t.Run("Partial Factory + auto builds", func(t *testing.T) {
		player, planet := newTestPlayerPlanet()

		// build the half completed factory, keep building more factories but don't add a partial mine
		planet.ProductionQueue = []ProductionQueueItem{
			{Type: QueueItemTypeFactory, Quantity: 1, Allocated: Cost{Germanium: 2, Resources: 5}},
			{Type: QueueItemTypeAutoMinTerraform, Quantity: 1},
			{Type: QueueItemTypeAutoFactories, Quantity: 100},
			{Type: QueueItemTypeAutoMines, Quantity: 100},
		}
		planet.Cargo = Cargo{361, 382, 1173, 331} // 42 res
		planet.Mines = 11
		planet.Factories = 16
		planet.Spec = computePlanetSpec(&rules, player, planet)

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

	})

	t.Run("Partial factory blocks queue", func(t *testing.T) {
		player, planet := newTestPlayerPlanet()

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
		producer := newProducer(testLogger, &rules, planet, player)
		producer.produce()

		// We should consume 1kT germanium and allocate appropriate resources to match
		assert.Equal(t, Cargo{7, 2, 0, 37}, planet.Cargo)
		assert.Equal(t, Cost{Germanium: 3, Resources: 7}, planet.ProductionQueue[0].Allocated)

	})

	t.Run("Growth check", func(t *testing.T) {
		player, planet := newTestPlayerPlanet()

		// make mines/factories cheap so we can build them
		player.Race.MineCost = 1
		player.Race.FactoryCost = 1
		player.Race.Spec = computeRaceSpec(&player.Race, &rules)

		// auto build with future growth taken into account
		planet.ProductionQueue = []ProductionQueueItem{
			{Type: QueueItemTypeAutoMines, Quantity: 100},
			{Type: QueueItemTypeAutoFactories, Quantity: 100},
		}
		planet.Cargo = Cargo{1000, 1000, 1000, 1000}
		planet.Spec = computePlanetSpec(&rules, player, planet)

		// max mines for current setting
		planet.Mines = planet.Spec.MaxMines
		planet.Factories = planet.Spec.MaxFactories

		// should build nothing, but queue up a mine partially done
		producer := newProducer(testLogger, &rules, planet, player)
		producer.produce()

		// we should build mines/factories accounting for future growth
		assert.Equal(t, 115, planet.Mines)
		assert.Equal(t, 115, planet.Factories)

	})

	t.Run("Ships", func(t *testing.T) {
		player, planet := newTestPlayerPlanet()
		player.TechLevels = TechLevel{0, 0, 3, 4, 0, 2}
		player.Race.PRT = SD
		player.Race.LRTs = Bitmask(IFE) | Bitmask(ARM) | Bitmask(BET) | Bitmask(RS)
		player.Race.PopEfficiency = 9
		player.Race.FactoryOutput = 11
		player.Race.Spec = computeRaceSpec(&player.Race, &rules)
		player.Spec = computePlayerSpec(player, &rules, []*Planet{planet})

		// add two designs, a colony ship and medium freighter with fuel mizer
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
		planet.Spec = computePlanetSpec(&rules, player, planet)

		// should build nothing, but queue up a mine partially done
		producer := newProducer(testLogger, &rules, planet, player)
		result, err := producer.produce()

		assert.Nil(t, err)
		assert.Equal(t, 3, len(result.itemsBuilt))
		assert.Equal(t, 1, result.itemsBuilt[0].designNum)
		assert.Equal(t, 1, result.itemsBuilt[0].numBuilt)
		assert.Equal(t, 2, result.itemsBuilt[1].designNum)
		assert.Equal(t, 2, result.itemsBuilt[1].numBuilt)

	})

	t.Run("Starbase Upgrade", func(t *testing.T) {
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
		planet.Spec = computePlanetSpec(&rules, player, planet)

		// should build nothing, but queue up a mine partially done
		producer := newProducer(testLogger, &rules, planet, player)
		result, err := producer.produce()

		// we should have built a starbase
		assert.Nil(t, err)
		assert.Equal(t, len(result.itemsBuilt), 1)
		assert.Equal(t, result.itemsBuilt[0].queueItemType, QueueItemTypeStarbase)
		assert.Equal(t, result.starbase, starbaseDesign2)

	})

	t.Run("Terraform", func(t *testing.T) {
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

	})

	t.Run("Packets", func(t *testing.T) {
		player, planet := newTestPlayerPlanet()

		// build 5 auto factories, leaving them in the queue
		planet.ProductionQueue = []ProductionQueueItem{{Type: QueueItemTypeMixedMineralPacket, Quantity: 1}}
		planet.Cargo = Cargo{100, 100, 100, 2500}
		planet.PacketSpeed = 6
		planet.PacketTargetNum = 1
		planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 100, PlanetStarbaseSpec: PlanetStarbaseSpec{HasMassDriver: true, SafePacketSpeed: 6, BasePacketSpeed: 6}}

		// should build 5 mine, leaving the auto build in the queue
		producer := newProducer(testLogger, &rules, planet, player)
		result, err := producer.produce()
		assert.Nil(t, err)
		assert.Equal(t, Cargo{40, 40, 40, 0}, result.packets)
		assert.Equal(t, 0, len(planet.ProductionQueue))
	})

	t.Run("Scanner", func(t *testing.T) {
		player, planet := newTestPlayerPlanet()

		// build a scanner
		planet.ProductionQueue = []ProductionQueueItem{{Type: QueueItemTypePlanetaryScanner, Quantity: 1}}
		planet.Cargo = Cargo{1000, 1000, 1000, 100_000}
		planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 1000, MaxPopulation: 1_000_000}
		planet.Scanner = false

		// planet should have scanner and an empty queue
		producer := newProducer(testLogger, &rules, planet, player)
		producer.produce()
		assert.True(t, planet.Scanner)
		assert.Equal(t, 0, len(planet.ProductionQueue))
	})
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
