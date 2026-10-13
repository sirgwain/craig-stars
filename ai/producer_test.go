package ai

import (
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// These tests check ship orders, not planet-spec recalculation. Keep the fixture's
// shipyard capacities fixed when produce submits its orders.
type productionTestOrderer struct{ cs.Orderer }

func (o productionTestOrderer) UpdatePlanetOrders(_ *cs.Rules, _ *cs.Player, planet *cs.Planet, orders cs.PlanetOrders) error {
	planet.PlanetOrders = orders
	return nil
}

// Use three ready shipyards and two designs so tests can distinguish correct
// allocation from broadcasting an order to every planet. The freighter is
// heavier than the scout so a small dock can build only part of a fleet.
func newProductionTestAI() *aiPlayer {
	game := cs.NewGame()
	player := cs.NewPlayer(1, cs.NewRace().WithSpec(&game.Rules))
	player.Num = 1
	player.Designs = []*cs.ShipDesign{
		cs.NewShipDesign(player.Num, 1).WithPurpose(cs.ShipDesignPurposeScout),
		cs.NewShipDesign(player.Num, 2).WithPurpose(cs.ShipDesignPurposeColonistFreighter),
	}
	player.Designs[0].Spec.Mass = 10
	player.Designs[1].Spec.Mass = 100
	planets := make([]*cs.Planet, 3)
	for i := range planets {
		planet := cs.NewPlanet()
		planet.Num = i + 1
		planet.PlayerNum = player.Num
		planet.Scanner = true
		planet.Spec.HasStarbase = true
		planet.Spec.DockCapacity = cs.UnlimitedSpaceDock
		// An existing base order keeps infrastructure planning out of these tests.
		planet.ProductionQueue = []cs.ProductionQueueItem{{Type: cs.QueueItemTypeStarbase, Quantity: 1}}
		planets[i] = planet
	}
	ai := NewAIPlayer(game, &cs.StaticTechStore, player, cs.PlayerMapObjects{Planets: planets})
	ai.client = productionTestOrderer{cs.NewOrderer()}
	return ai
}

// Count ship quantities for this mission only. Infrastructure rows and ships
// assigned to other missions should not affect the production assertions.
func queuedScoutShips(planet *cs.Planet) map[int]int {
	counts := map[int]int{}
	for _, item := range planet.ProductionQueue {
		if item.Type == cs.QueueItemTypeShipToken && item.GetTag(cs.TagPurpose) == string(cs.FleetPurposeScout) {
			counts[item.DesignNum] += item.Quantity
		}
	}
	return counts
}

// One scout is one fleet. Pending production must satisfy the requested count
// across the empire; complete idle fleets must not be deducted again because
// the managers already accounted for them when requesting more fleets.
func Test_aiPlayer_produce_fleetRequestCounts(t *testing.T) {
	tests := []struct {
		name      string
		requested int
		queued    [3]int
		idle      [3]int
		want      [3]int // Total queued scouts at each shipyard after production.
	}{
		{name: "no demand queues nothing"},
		{name: "one request uses one shipyard", requested: 1, want: [3]int{1, 0, 0}},
		{name: "two requests use two shipyards", requested: 2, want: [3]int{1, 1, 0}},
		{name: "three requests use all shipyards", requested: 3, want: [3]int{1, 1, 1}},
		{name: "pending fleet at a later shipyard satisfies demand", requested: 1, queued: [3]int{0, 0, 1}, want: [3]int{0, 0, 1}},
		{name: "separate queue rows supply two fleets", requested: 2, queued: [3]int{0, 0, 2}, want: [3]int{0, 0, 2}},
		{name: "pending fleet leaves one additional fleet to build", requested: 2, queued: [3]int{0, 0, 1}, want: [3]int{0, 0, 2}},
		{name: "idle fleet does not satisfy new demand", requested: 1, idle: [3]int{1, 0, 0}, want: [3]int{1, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ai := newProductionTestAI()
			for i, planet := range ai.Planets {
				for range tt.queued[i] {
					ai.addShipToTopOfQueue(planet, cs.FleetPurposeScout, ai.Designs[0], 1)
				}
				if tt.idle[i] > 0 {
					f := cs.NewFleet(ai.Player, i+1, "Idle Scout", []cs.Waypoint{{}})
					f.SetTag(cs.TagPurpose, string(cs.FleetPurposeScout))
					f.Tokens = []cs.ShipToken{{DesignNum: 1, Quantity: tt.idle[i]}}
					ai.fleetsByPlanetNum[planet.Num] = []*cs.Fleet{f}
				}
			}
			ai.addFleetBuildRequest(cs.FleetPurposeScout, tt.requested)

			// Reconsidering the same demand must not duplicate pending orders.
			for range 2 {
				require.NoError(t, ai.produce())
				for i, planet := range ai.Planets {
					assert.Equal(t, tt.want[i], queuedScoutShips(planet)[1], "shipyard %d", i+1)
				}
			}
		})
	}
}

// Finish a partially supplied fleet at its existing shipyard. Only two more
// freighters are needed: idle and queued ships already supply everything else.
func Test_aiPlayer_produce_completesPartialFleet(t *testing.T) {
	ai := newProductionTestAI()
	ai.fleetsByPurpose[cs.FleetPurposeScout] = fleet{
		purpose: cs.FleetPurposeScout,
		ships: []fleetShip{
			{purpose: cs.ShipDesignPurposeScout, quantity: 2},
			{purpose: cs.ShipDesignPurposeColonistFreighter, quantity: 3},
		},
	}
	planet := ai.Planets[2]
	ai.addShipToTopOfQueue(planet, cs.FleetPurposeScout, ai.Designs[0], 1)
	// Queued ships assigned to another mission cannot satisfy this request.
	ai.addShipToTopOfQueue(planet, cs.FleetPurposeColonizer, ai.Designs[0], 20)
	f := cs.NewFleet(ai.Player, 1, "Idle", []cs.Waypoint{{}})
	f.SetTag(cs.TagPurpose, string(cs.FleetPurposeScout))
	f.Tokens = []cs.ShipToken{{DesignNum: 1, Quantity: 1}, {DesignNum: 2, Quantity: 1}}
	ai.fleetsByPlanetNum[planet.Num] = []*cs.Fleet{f}
	ai.addFleetBuildRequest(cs.FleetPurposeScout, 1)

	require.NoError(t, ai.produce())
	assert.Empty(t, queuedScoutShips(ai.Planets[0]))
	assert.Empty(t, queuedScoutShips(ai.Planets[1]))
	assert.Equal(t, map[int]int{1: 1, 2: 2}, queuedScoutShips(planet))
}

// A fleet needs both a scout and a freighter. Skip shipyards that cannot build
// both and queue the whole fleet at the next capable yard. If a design is
// unavailable everywhere, queue nothing rather than an unusable partial fleet.
func Test_aiPlayer_produce_skipsShipyardsThatCannotCompleteFleet(t *testing.T) {
	tests := []struct {
		name         string
		dockCapacity int
		hasStarbase  bool
		secondShip   cs.ShipDesignPurpose
		wantPlanet   int // Index of the shipyard receiving the fleet; -1 means none.
	}{
		{
			name:         "dock can build the scout but not the freighter",
			dockCapacity: 50,
			hasStarbase:  true,
			secondShip:   cs.ShipDesignPurposeColonistFreighter,
			wantPlanet:   1,
		},
		{
			name:         "planet has no starbase",
			dockCapacity: cs.UnlimitedSpaceDock,
			secondShip:   cs.ShipDesignPurposeColonistFreighter,
			wantPlanet:   1,
		},
		{
			name:         "required bomber design is unavailable",
			dockCapacity: cs.UnlimitedSpaceDock,
			hasStarbase:  true,
			secondShip:   cs.ShipDesignPurposeBomber,
			wantPlanet:   -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ai := newProductionTestAI()
			ai.fleetsByPurpose[cs.FleetPurposeScout] = fleet{
				purpose: cs.FleetPurposeScout,
				ships: []fleetShip{
					{purpose: cs.ShipDesignPurposeScout, quantity: 1},
					{purpose: tt.secondShip, quantity: 1},
				},
			}
			ai.Planets[0].Spec.DockCapacity = tt.dockCapacity
			ai.Planets[0].Spec.HasStarbase = tt.hasStarbase
			ai.addFleetBuildRequest(cs.FleetPurposeScout, 1)

			require.NoError(t, ai.produce())
			for i, planet := range ai.Planets {
				if i == tt.wantPlanet {
					assert.Equal(t, map[int]int{1: 1, 2: 1}, queuedScoutShips(planet))
				} else {
					assert.Empty(t, queuedScoutShips(planet), "shipyard %d", i+1)
				}
			}
		})
	}
}

// A scout queued at one planet and a freighter queued at another do not supply
// a local fleet. Expect a freighter at the first planet to complete one fleet,
// with no additional orders elsewhere.
func Test_aiPlayer_produce_partialFleetsAtDifferentPlanets(t *testing.T) {
	ai := newProductionTestAI()
	ai.fleetsByPurpose[cs.FleetPurposeScout] = fleet{
		purpose: cs.FleetPurposeScout,
		ships: []fleetShip{
			{purpose: cs.ShipDesignPurposeScout, quantity: 1},
			{purpose: cs.ShipDesignPurposeColonistFreighter, quantity: 1},
		},
	}
	ai.addShipToTopOfQueue(ai.Planets[0], cs.FleetPurposeScout, ai.Designs[0], 1)
	ai.addShipToTopOfQueue(ai.Planets[1], cs.FleetPurposeScout, ai.Designs[1], 1)
	ai.addFleetBuildRequest(cs.FleetPurposeScout, 1)

	require.NoError(t, ai.produce())
	assert.Equal(t, map[int]int{1: 1, 2: 1}, queuedScoutShips(ai.Planets[0]))
	assert.Equal(t, map[int]int{2: 1}, queuedScoutShips(ai.Planets[1]))
	assert.Empty(t, queuedScoutShips(ai.Planets[2]))
}

func TestProductionDropsUnneededColonyOrders(t *testing.T) {
	ai := newProductionTestAI()
	design := cs.NewShipDesign(1, 3).WithPurpose(cs.ShipDesignPurposeColonizer)
	design.Spec.Colonizer = true
	ai.Designs = append(ai.Designs, design)
	ai.addShipToTopOfQueue(ai.Planets[0], cs.FleetPurposeColonizer, design, 12)
	ai.Planets[0].ProductionQueue[0].Allocated = cs.Cost{Resources: 5}
	ai.addShipToTopOfQueue(ai.Planets[1], cs.FleetPurposeColonizer, design, 12)
	ai.addFleetBuildRequest(cs.FleetPurposeColonizer, 0)
	assert.NoError(t, ai.produce())
	assert.Equal(t, 1, ai.Planets[0].ProductionQueue[0].Quantity)
	assert.False(t, ai.isFleetInQueue(ai.Planets[1], cs.FleetPurposeColonizer))
}

func TestInfrastructureEstimateUsesMineralsInStock(t *testing.T) {
	ai := newProductionTestAI()
	planet := ai.Planets[0]
	planet.Spec.MiningOutput = cs.Mineral{}
	planet.Spec.ResourcesPerYearAvailable = 1000
	planet.Cargo = cs.Cargo{Ironium: 1000, Boranium: 1000, Germanium: 1000}
	years, err := ai.getYearsToBuild(planet)
	assert.NoError(t, err)
	assert.Equal(t, 1, years)
}

func TestMinesGoAheadOfAnOrderShortOfMinerals(t *testing.T) {
	ai := newProductionTestAI()
	planet := ai.Planets[0]
	planet.Spec.MaxMines = 100
	planet.Spec.ResourcesPerYearAvailable = 200
	design := ai.Designs[0]
	design.Spec.Cost = cs.Cost{Ironium: 50, Resources: 20}
	planet.ProductionQueue = nil
	ai.addShipToTopOfQueue(planet, cs.FleetPurposeScout, design, 1)

	// enough ironium in stock and from this year's mining: leave the queue alone
	planet.Cargo.Ironium = 30
	planet.Spec.MiningOutput = cs.Mineral{Ironium: 20}
	assert.NoError(t, ai.reserveEconomyProduction(planet))
	assert.Len(t, planet.ProductionQueue, 1)

	// short of ironium: the ship would stop the queue, so mines go first
	planet.Cargo.Ironium = 10
	assert.NoError(t, ai.reserveEconomyProduction(planet))
	assert.Len(t, planet.ProductionQueue, 2)
	assert.Equal(t, cs.QueueItemTypeAutoMines, planet.ProductionQueue[0].Type)
	assert.Equal(t, 10, planet.ProductionQueue[0].Quantity)

	// running again replaces the reservation rather than stacking another
	assert.NoError(t, ai.reserveEconomyProduction(planet))
	assert.Len(t, planet.ProductionQueue, 2)

	// and it is removed once the order can be paid for
	planet.Cargo.Ironium = 100
	assert.NoError(t, ai.reserveEconomyProduction(planet))
	assert.Len(t, planet.ProductionQueue, 1)
	assert.Equal(t, cs.QueueItemTypeShipToken, planet.ProductionQueue[0].Type)
}

func TestPacketsGoAheadOfShips(t *testing.T) {
	ai := newProductionTestAI()
	planet := ai.Planets[0]
	planet.ProductionQueue = []cs.ProductionQueueItem{{Type: cs.QueueItemTypeIroniumMineralPacket, Quantity: 1}}
	ai.addShipToTopOfQueue(planet, cs.FleetPurposeScout, ai.Designs[0], 1)
	ai.addShipToTopOfQueue(planet, cs.FleetPurposeColonistFreighter, ai.Designs[1], 2)

	assert.NoError(t, ai.produce())
	assert.Equal(t, cs.QueueItemTypeIroniumMineralPacket, planet.ProductionQueue[0].Type)
}

func TestTerraformingStartsOnceAWorldIsHalfDeveloped(t *testing.T) {
	ai := newProductionTestAI()
	planet := ai.Planets[0]
	planet.ProductionQueue = nil
	planet.Spec.CanTerraform = true
	planet.Spec.MaxFactories, planet.Spec.MaxMines = 100, 100
	planet.Spec.ResourcesPerYearAvailable = 800

	planet.Factories, planet.Mines = 40, 40
	assert.NoError(t, ai.reserveEconomyProduction(planet))
	assert.Empty(t, planet.ProductionQueue, "factories and mines come first on a young world")

	planet.Factories, planet.Mines = 60, 60
	assert.NoError(t, ai.reserveEconomyProduction(planet))
	if assert.Len(t, planet.ProductionQueue, 1) {
		assert.Equal(t, cs.QueueItemTypeAutoMaxTerraform, planet.ProductionQueue[0].Type)
		assert.Equal(t, 800/4/ai.game.Rules.TerraformCost.Resources, planet.ProductionQueue[0].Quantity)
	}

	// a finished world terraforms through its normal plan
	planet.Factories, planet.Mines = 100, 100
	assert.NoError(t, ai.reserveEconomyProduction(planet))
	assert.Empty(t, planet.ProductionQueue)
}
