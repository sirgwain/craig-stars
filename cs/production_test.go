package cs

import (
	"math"
	"testing"

	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func Test_production_produce(t *testing.T) {
	// TODO: Add tests for auto mineral alchemy
	t.Run("Concrete mine removed from queue", func(t *testing.T) {
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
	})

	t.Run("Auto items stay in queue", func(t *testing.T) {
		t.Run("Factories", func(t *testing.T) {
			player, planet := newTestPlayerPlanet()

			// build 5 auto factories, leaving them in the queue
			planet.ProductionQueue = []ProductionQueueItem{{Type: QueueItemTypeAutoFactories, Quantity: 5}}
			planet.Cargo = Cargo{10, 20, 30, 2500}
			planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 100, MaxFactories: 100, MaxPopulation: 1_000_000}
			planet.Factories = 0
			player.Messages = []PlayerMessage{}

			// should leave the auto build in the queue
			producer := newProducer(testLogger, &rules, planet, player)
			result, err := producer.produce()
			assert.Nil(t, err)

			assert.Equal(t, 5, planet.Factories)
			assert.Equal(t, 1, len(planet.ProductionQueue))
			assert.Equal(t, QueueItemTypeAutoFactories, planet.ProductionQueue[0].Type)
			assert.Equal(t, 5, planet.ProductionQueue[0].Quantity)
			assert.Equal(t, 50, result.leftoverResources) // 50 resources leftover
		})

		t.Run("Mines", func(t *testing.T) {
			player, planet := newTestPlayerPlanet()

			// build 5 auto mines, leaving them in the queue
			planet.ProductionQueue = []ProductionQueueItem{{Type: QueueItemTypeAutoMines, Quantity: 5}}
			planet.Cargo = Cargo{10, 20, 30, 2500}
			planet.Spec = PlanetSpec{ResourcesPerYearAvailable: 100, MaxMines: 100, MaxPopulation: 1_000_000}
			planet.Mines = 0
			player.Messages = []PlayerMessage{}

			// should build 5 mines, leaving the auto build in the queue
			producer := newProducer(testLogger, &rules, planet, player)
			producer.produce()
			assert.Equal(t, 5, planet.Mines)
			assert.Equal(t, 1, len(planet.ProductionQueue))
			assert.Equal(t, QueueItemTypeAutoMines, planet.ProductionQueue[0].Type)
			assert.Equal(t, 5, planet.ProductionQueue[0].Quantity)
		})
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

	t.Run("Refund invalid packets", func(t *testing.T) {
		player, planet := newTestPlayerPlanet()
		planet.Name = t.Name()

		// planet has target, but no driver
		planet.Cargo = Cargo{100, 100, 100, 1000}
		planet.Spec.ResourcesPerYearAvailable = 1
		planet.PacketTargetNum = 23
		planet.Spec.PlanetStarbaseSpec = PlanetStarbaseSpec{
			HasMassDriver: false, // oops, no driver!
		}

		// make a bunch of packets, along with some stuff to block queue
		planet.ProductionQueue = []ProductionQueueItem{
			{Type: QueueItemTypeMixedMineralPacket, Quantity: 1, Allocated: Cost{1, 1, 1, 3}},
			{Type: QueueItemTypeMixedMineralPacket, Quantity: 2},
			{Type: QueueItemTypeMine, Quantity: 1},
			{Type: QueueItemTypeBoraniumMineralPacket, Quantity: 2, Allocated: Cost{2, 2, 2, 6}},
			{Type: QueueItemTypeGermaniumMineralPacket, Quantity: 2},
			{Type: QueueItemTypeAutoMineralPacket, Quantity: 1},
			{Type: QueueItemTypeMixedMineralPacket, Quantity: 5, Allocated: Cost{2, 2, 2, 6}},
			{Type: QueueItemTypeDefenses, Quantity: 100},
			{Type: QueueItemTypeIroniumMineralPacket, Quantity: 3},
		}

		// build once
		player.Messages = []PlayerMessage{}
		producer := newProducer(testLogger, &rules, planet, player)
		result, err := producer.produce()
		assert.NoError(t, err)

		// Should cancel the first 2 packets before stopping to build the mine
		wantQueue := []ProductionQueueItem{
			{Type: QueueItemTypeMine, Quantity: 1, Allocated: Cost{0, 0, 0, 4}},
			{Type: QueueItemTypeBoraniumMineralPacket, Quantity: 2, Allocated: Cost{2, 2, 2, 6}},
			{Type: QueueItemTypeGermaniumMineralPacket, Quantity: 2},
			{Type: QueueItemTypeAutoMineralPacket, Quantity: 1},
			{Type: QueueItemTypeMixedMineralPacket, Quantity: 5, Allocated: Cost{2, 2, 2, 6}},
			{Type: QueueItemTypeDefenses, Quantity: 100},
			{Type: QueueItemTypeIroniumMineralPacket, Quantity: 3},
		}

		// We canceled 2 similar packet types, so QueueItemType is still set
		wantMessages := []PlayerMessage{{
			Target: PlayerMessageTarget{
				TargetType: TargetPlanet,
				TargetName: planet.Name,
				TargetNum:  planet.Num,
			},
			Type: PlayerMessagePlanetBuiltInvalidMineralPacketNoMassDriver,
			Spec: PlayerMessageSpec{
				Name:          planet.Name,
				QueueItemType: QueueItemTypeMixedMineralPacket,
				Cost:          Cost{1, 1, 1, 3},
				Amount:        360, // weight of canceled packets
				Amount2:       2,   // no. of canceled orders
				PrevAmount:    3,   // no. of canceled items
			},
		}}

		test.CompareAsJSON(t, planet.ProductionQueue, wantQueue)
		test.CompareAsJSON(t, result.messages, wantMessages)

		// give planet a starbase with a driver
		starbaseDesign := NewShipDesign(player.Num, 3).WithHull(SpaceStation.Name).
			WithSlots([]ShipDesignSlot{
				{HullComponent: MassDriver5.Name, HullSlotIndex: 1, Quantity: 1},
			}).WithSpec(&rules, player)
		starbaseFleet := newStarbase(player, planet, starbaseDesign, t.Name())
		starbaseFleet.Spec = ComputeFleetSpec(&rules, player, &starbaseFleet)
		planet.Starbase = &starbaseFleet
		planet.Spec.PlanetStarbaseSpec = computePlanetStarbaseSpec(planet)

		planet.PacketTargetNum = None // Oops, no target!

		// Reset messages and produce again
		player.Messages = []PlayerMessage{}
		producer = newProducer(testLogger, &rules, planet, player)
		result, err = producer.produce()
		assert.NoError(t, err)

		// Mine should be finished and 2nd batch of concrete packets should've been canceled;
		// last iron packet stays due to not being reached yet
		wantQueue = []ProductionQueueItem{
			{Type: QueueItemTypeAutoMineralPacket, Quantity: 1}, // gets skipped instead of canceled
			{Type: QueueItemTypeDefenses, Quantity: 100, Allocated: Cost{4, 4, 4, 12}},
			{Type: QueueItemTypeIroniumMineralPacket, Quantity: 3},
		}

		wantMessages = []PlayerMessage{{
			Target: PlayerMessageTarget{
				TargetType: TargetPlanet,
				TargetName: planet.Name,
				TargetNum:  planet.Num,
			},
			Type: PlayerMessagePlanetBuiltInvalidMineralPacketNoTarget,
			Spec: PlayerMessageSpec{
				Name:          planet.Name,
				QueueItemType: QueueItemTypeNone, // No QueueItemType due to multiple packet types being canceled
				Cost:          Cost{4, 4, 4, 12},
				Amount:        1000, // weight of canceled packets
				Amount2:       3,    // no. of canceled orders
				PrevAmount:    9,    // amount of canceled items
			},
		}}

		test.CompareAsJSON(t, planet.ProductionQueue, wantQueue)
		test.CompareAsJSON(t, result.messages, wantMessages)
	})

	t.Run("Refund invalid items", func(t *testing.T) {
		player, planet := newTestPlayerPlanet()

		// make defenses an even 10 in all for cost,
		// and make scanners abhorrently expensive
		rCopy := rules
		rCopy.DefenseCost = Cost{10, 10, 10, 10}
		rCopy.PlanetaryScannerCost = Cost{999, 999, 999, 999}

		planet.Name = t.Name()
		// exactly enough to finish 10 defenses
		planet.Cargo = Cargo{100, 100, 100, 1000}
		planet.Defenses = 90
		planet.ContributesOnlyLeftoverToResearch = true

		// Queue up 100 auto defenses, with some concrete defenses in the back of the queue
		planet.ProductionQueue = []ProductionQueueItem{
			{Type: QueueItemTypeAutoDefenses, Quantity: 100},
			{Type: QueueItemTypeDefenses, Quantity: 90, Allocated: Cost{5, 5, 5, 5}},
			{Type: QueueItemTypeDefenses, Quantity: 11, Allocated: Cost{5, 5, 5, 5}},
			{Type: QueueItemTypeDefenses, Quantity: 10, Allocated: Cost{5, 5, 5, 5}},
			// scanner to soak up leftover allocated stuff
			{Type: QueueItemTypePlanetaryScanner, Quantity: 1},
		}

		player.Race = *player.Race.WithSpec(&rCopy)
		planet.Spec = computePlanetSpec(&rCopy, player, planet)
		player.Spec = computePlayerSpec(player, &rCopy, []*Planet{planet})
		player.Messages = []PlayerMessage{}

		// should end up with 100 defenses, with auto defenses still in the queue;
		// scanner should soak up allocated cost of canceled items
		wantQueue := []ProductionQueueItem{
			{Type: QueueItemTypeAutoDefenses, Quantity: 100},
			{Type: QueueItemTypePlanetaryScanner, Quantity: 1, Allocated: Cost{15, 15, 15, 15}},
		}

		wantMessages := []PlayerMessage{{
			Target: PlayerMessageTarget{
				TargetType: TargetPlanet,
				TargetName: planet.Name,
				TargetNum:  planet.Num,
			},
			Type: PlayerMessagePlanetBuiltInvalidItem,
			Spec: PlayerMessageSpec{
				Name:          planet.Name,
				Cost:          Cost{15, 15, 15, 15},
				QueueItemType: QueueItemTypeDefenses,
				Amount:        111, // amount canceled
				Amount2:       0,   // MaxBuildable at moment of first build
				PrevAmount:    111, // Prior amount in queue
			},
		}}

		// build time!
		producer := newProducer(testLogger, &rCopy, planet, player)
		result, err := producer.produce()
		assert.NoError(t, err)
		assert.Equal(t, 100, planet.Defenses)
		assert.Equal(t, Mineral{}, planet.Cargo.ToMineral()) // make sure we actually paid for it
		test.CompareAsJSON(t, planet.ProductionQueue, wantQueue)
		test.CompareAsJSON(t, result.messages, wantMessages)
	})

	t.Run("Clamp autos over 5K", func(t *testing.T) {
		player, planet := newTestPlayerPlanet()
		planet.Cargo = Cargo{1000, 1000, 1000, 2000}

		// many many auto items
		planet.ProductionQueue = []ProductionQueueItem{
			{Type: QueueItemTypeAutoFactories, Quantity: 42069},
			{Type: QueueItemTypeAutoDefenses, Quantity: 1337},
			{Type: QueueItemTypeAutoMaxTerraform, Quantity: 69420},
			{Type: QueueItemTypeAutoMinTerraform, Quantity: math.MaxInt},
		}

		player.Race = *player.Race.WithSpec(&rules)
		planet.Spec = computePlanetSpec(&rules, player, planet)
		player.Spec = computePlayerSpec(player, &rules, []*Planet{planet})
		player.Messages = []PlayerMessage{}

		// scanner got built, and all auto items over cap got clamped
		wantQueue := []ProductionQueueItem{
			{Type: QueueItemTypeAutoFactories, Quantity: MaxBuildableCap},
			{Type: QueueItemTypeAutoDefenses, Quantity: 1337},
			{Type: QueueItemTypeAutoMaxTerraform, Quantity: MaxBuildableCap},
			{Type: QueueItemTypeAutoMinTerraform, Quantity: MaxBuildableCap},
		}

		producer := newProducer(testLogger, &rules, planet, player)
		result, err := producer.produce()
		assert.NoError(t, err)
		// should've built stuff and clamped the auto items;
		// no message produced due to nothing being explicitly "canceled".
		assert.Greater(t, planet.Factories, 0)
		test.CompareAsJSON(t, planet.ProductionQueue, wantQueue)
		assert.Equal(t, 0, len(result.messages))
	})

	t.Run("Don't refund invalid items if nothing built", func(t *testing.T) {
		player, planet := newTestPlayerPlanet()

		rCopy := rules
		rCopy.DefenseCost = Cost{0, 0, 0, 0} // defenses finish instantly if we get to them

		planet.Cargo = Cargo{1, 1, 7, 100}
		planet.Defenses = 100
		planet.ContributesOnlyLeftoverToResearch = true

		// 10 concrete defenses over cap
		planet.ProductionQueue = []ProductionQueueItem{
			// super expensive scanner to soak up leftover minerals
			{Type: QueueItemTypePlanetaryScanner, Quantity: 1},
			{Type: QueueItemTypeDefenses, Quantity: 10, Allocated: Cost{5, 5, 5, 5}},
		}

		player.Race = *player.Race.WithSpec(&rCopy)
		planet.Spec = computePlanetSpec(&rCopy, player, planet)
		player.Spec = computePlayerSpec(player, &rCopy, []*Planet{planet})
		player.Messages = []PlayerMessage{}

		// defenses stay in queue since nothing got built
		wantQueue := []ProductionQueueItem{
			{Type: QueueItemTypePlanetaryScanner, Quantity: 1, Allocated: Cost{1, 1, 7, 10}},
			{Type: QueueItemTypeDefenses, Quantity: 10, Allocated: Cost{5, 5, 5, 5}},
		}

		producer := newProducer(testLogger, &rCopy, planet, player)
		producer.produce()
		assert.Equal(t, 100, planet.Defenses)
		test.CompareAsJSON(t, planet.ProductionQueue, wantQueue)
	})

	t.Run("Autos cap at max", func(t *testing.T) {
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

	t.Run("Partial factory + auto builds", func(t *testing.T) {
		player, planet := newTestPlayerPlanet()

		// build the half completed factory, keep building more factories but don't add a partial mine
		planet.ProductionQueue = []ProductionQueueItem{
			{Type: QueueItemTypeFactory, Quantity: 1, Allocated: Cost{Germanium: 2, Resources: 5}},
			{Type: QueueItemTypeAutoMinTerraform, Quantity: 1}, // skipped
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

		// try and make 1 concrete factory when out of minerals
		planet.ProductionQueue = []ProductionQueueItem{
			{Type: QueueItemTypeFactory, Quantity: 1, Allocated: Cost{Germanium: 2, Resources: 5}},
			{Type: QueueItemTypeAutoMinTerraform, Quantity: 1},
			{Type: QueueItemTypeAutoFactories, Quantity: 100},
			{Type: QueueItemTypeAutoMines, Quantity: 100},
		}
		planet.Cargo = Cargo{7, 2, 1, 1000}
		planet.Spec.ResourcesPerYearAvailable = 5

		// build nothing due to partial factory blocking the queue
		wantQueue := []ProductionQueueItem{
			{Type: QueueItemTypeFactory, Quantity: 1, Allocated: Cost{Germanium: 3, Resources: 7}},
			{Type: QueueItemTypeAutoMinTerraform, Quantity: 1},
			{Type: QueueItemTypeAutoFactories, Quantity: 100},
			{Type: QueueItemTypeAutoMines, Quantity: 100},
		}

		producer := newProducer(testLogger, &rules, planet, player)
		result, err := producer.produce()
		assert.NoError(t, err)

		// should consume 1kT germanium, 2 res and allocate appropriate resources to match
		// Remaining resources go into research (we have nothing else to do with em)
		assert.Equal(t, Cargo{7, 2, 0, 1000}, planet.Cargo)
		test.CompareAsJSON(t, planet.ProductionQueue, wantQueue)
		assert.Equal(t, 0, result.factories)
		assert.Equal(t, 3, result.leftoverResources)

	})

	t.Run("Auto accounts for growth", func(t *testing.T) {
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
		planet.Cargo = Cargo{1000, 1000, 1000, 1000} // 100K pop; grows to 115K next year
		planet.Spec = computePlanetSpec(&rules, player, planet)

		// max out installations for current pop level
		planet.Mines = planet.Spec.MaxMines
		planet.Factories = planet.Spec.MaxFactories

		producer := newProducer(testLogger, &rules, planet, player)
		producer.produce()

		// should've built mines/factories up to the amount we can operate for next yr
		assert.Equal(t, planet.MaxBuildable(player, QueueItemTypeAutoFactories), 0)
		assert.Equal(t, planet.MaxBuildable(player, QueueItemTypeAutoMines), 0)

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

		// create 2 designs for the old and new bases
		starbaseDesign1 := NewShipDesign(player.Num, 2).WithHull(SpaceStation.Name).WithSpec(&rules, player)
		starbaseFleet := newStarbase(player, planet, starbaseDesign1, "Old Base")
		starbaseDesign2 := NewShipDesign(player.Num, 3).WithHull(SpaceStation.Name).
			WithSlots([]ShipDesignSlot{
				{HullComponent: MassDriver5.Name, HullSlotIndex: 1, Quantity: 1},
			}).WithSpec(&rules, player)

		player.Designs = append(player.Designs, starbaseDesign1, starbaseDesign2)
		starbaseFleet.Spec = ComputeFleetSpec(&rules, player, &starbaseFleet)
		planet.Starbase = &starbaseFleet

		planet.ProductionQueue = []ProductionQueueItem{
			{Type: QueueItemTypeStarbase, Quantity: 1, DesignNum: 3, design: starbaseDesign2},
		}
		planet.Cargo = Cargo{1000, 1000, 1000, 10000}
		planet.Spec = computePlanetSpec(&rules, player, planet)

		producer := newProducer(testLogger, &rules, planet, player)
		result, err := producer.produce()

		// should have built a starbase to replace old one
		assert.NoError(t, err)
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

	t.Run("Starbase, then packet", func(t *testing.T) {
		player, planet := newTestPlayerPlanet()
		starbaseDesign := NewShipDesign(player.Num, 3).WithHull(SpaceStation.Name).
			WithSlots([]ShipDesignSlot{
				{HullComponent: MassDriver5.Name, HullSlotIndex: 1, Quantity: 1},
			}).WithSpec(&rules, player)
		planet.PacketTargetNum = 33

		player.Designs = append(player.Designs, starbaseDesign)

		// build a new starbase with a mass driver
		planet.ProductionQueue = []ProductionQueueItem{
			{Type: QueueItemTypeStarbase, DesignNum: 3, Quantity: 1},
			{Type: QueueItemTypeMixedMineralPacket, Quantity: 1},
		}
		planet.Cargo = Cargo{500, 500, 500, 25000}
		planet.Spec.ResourcesPerYearAvailable = 2000

		// should build starbase, then packet
		planet.PopulateProductionQueueDesigns(player)
		producer := newProducer(testLogger, &rules, planet, player)
		result, err := producer.produce()
		assert.Nil(t, err)
		assert.NotNil(t, result.starbase)
		assert.Equal(t, Cargo{40, 40, 40, 0}, result.packets) // false if item is canceled
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
	// TODO: Add test for multiple starbases in 1 turn
}

func Test_production_allocatePartialBuild(t *testing.T) {
	player, planet := newTestPlayerPlanet()
	production := newProducer(testLogger, &rules, planet, player)

	tests := []struct {
		name        string
		costPerItem Cost
		available   Cost
		want        Cost
	}{
		{
			name:        "more than enough to build",
			costPerItem: Cost{10, 10, 10, 10},
			available:   Cost{100, 100, 100, 100},
			want:        Cost{10, 10, 10, 10},
		},
		{
			name:        "less than half a factory",
			costPerItem: Cost{0, 0, 4, 10},
			available:   Cost{100, 100, 2, 4},
			want:        Cost{0, 0, 1, 4},
		},
		{
			name:        "Up to 50% due to ironium shortage",
			costPerItem: Cost{10, 20, 30, 40},
			available:   Cost{5, 100, 100, 100},
			want:        Cost{5, 10, 15, 20},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := production.allocatePartialBuild(tt.costPerItem, tt.available); got != tt.want {
				t.Errorf("Planet.allocatePartialBuild() = %v, want %v", got, tt.want)
			}
		})
	}
}
