package ai

import (
	"strconv"
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// an idle transport fleet orbiting one of the test AI's planets
func newTestTransport(t *testing.T, ai *aiPlayer, num int, planet *cs.Planet) *cs.Fleet {
	t.Helper()
	design := ai.GetDesign(100)
	if design == nil {
		var err error
		design, err = cs.DesignShip(&ai.game.Rules, ai.Player, &cs.SmallFreighter, "Freighter", 100, 0, cs.ShipDesignPurposeColonistFreighter, cs.FleetPurposeColonistFreighter)
		require.NoError(t, err)
		ai.Designs = append(ai.Designs, design)
	}
	fleet := cs.NewFleet(ai.Player, num, "Freighter", []cs.Waypoint{cs.NewPlanetWaypoint(planet.Position, planet.Num, planet.Name, 5)})
	fleet.OrbitingPlanetNum = planet.Num
	fleet.Position = planet.Position
	fleet.Tokens = []cs.ShipToken{{DesignNum: design.Num, Quantity: 1}}
	fleet.InjectDesigns(ai.Designs)
	fleet.Spec = cs.ComputeFleetSpec(&ai.game.Rules, ai.Player, fleet)
	fleet.Fuel = fleet.Spec.FuelCapacity
	fleet.SetTag(cs.TagPurpose, string(cs.FleetPurposeColonistFreighter))
	ai.Fleets = append(ai.Fleets, fleet)
	return fleet
}

// spread the test AI's planets out so fleets can travel between them
func placeTestPlanets(ai *aiPlayer) {
	for i, planet := range ai.Planets {
		planet.Position = cs.Vector{X: 100 + 30*i, Y: 100}
		planet.Name = "Planet " + strconv.Itoa(planet.Num)
	}
	ai.PlanetIntels = nil
	for _, planet := range ai.Planets {
		intel := *planet
		ai.PlanetIntels = append(ai.PlanetIntels, &intel)
	}
}

func TestIdleYearsAreTrackedOnTheFleet(t *testing.T) {
	ai := newProductionTestAI()
	placeTestPlanets(ai)
	fleet := newTestTransport(t, ai, 1, ai.Planets[0])

	ai.game.Year = 2420
	ai.trackIdleFleets()
	ai.game.Year = 2426
	ai.trackIdleFleets()
	assert.Equal(t, 6, ai.idleYears(fleet))

	// a fleet with orders starts again from zero
	fleet.Waypoints = append(fleet.Waypoints, cs.NewPlanetWaypoint(ai.Planets[1].Position, ai.Planets[1].Num, ai.Planets[1].Name, 5))
	ai.trackIdleFleets()
	assert.Zero(t, ai.idleYears(fleet))
}

func TestSurplusTransportsAreScrapped(t *testing.T) {
	ai := newProductionTestAI()
	placeTestPlanets(ai)
	ai.game.Year = 2440
	fleets := []*cs.Fleet{}
	for i, since := range []int{2430, 2437, 2438} {
		fleet := newTestTransport(t, ai, i+1, ai.Planets[0])
		fleet.SetTag(tagIdleSince, strconv.Itoa(since))
		fleets = append(fleets, fleet)
	}

	ai.retireTransports()

	// two are surplus, but only one has been idle long enough to scrap
	assert.Equal(t, cs.WaypointTaskScrapFleet, fleets[0].Waypoints[0].Task)
	assert.Equal(t, cs.WaypointTaskNone, fleets[1].Waypoints[0].Task)
	assert.Equal(t, cs.WaypointTaskNone, fleets[2].Waypoints[0].Task)
}

func TestGermaniumIsHauledToWaitingFactories(t *testing.T) {
	ai := newProductionTestAI()
	placeTestPlanets(ai)
	source, colony := ai.Planets[0], ai.Planets[1]
	source.Cargo.Germanium = 5000
	colony.Spec.MaxFactories = 200
	fleet := newTestTransport(t, ai, 1, source)
	require.GreaterOrEqual(t, fleet.Spec.CargoCapacity, haulMinimumLoad)

	require.NoError(t, ai.haulMinerals())

	require.Len(t, fleet.Waypoints, 2)
	assert.Equal(t, colony.Num, fleet.Waypoints[1].TargetNum)
	assert.Equal(t, cs.TransportActionUnloadAll, fleet.Waypoints[1].TransportTasks.Germanium.Action)
	assert.Equal(t, fleet.Spec.CargoCapacity, fleet.Cargo.Germanium)
	assert.Equal(t, 5000-fleet.Spec.CargoCapacity, source.Cargo.Germanium)
}

func TestGermaniumIsNotHauledFromAWorldThatNeedsIt(t *testing.T) {
	ai := newProductionTestAI()
	placeTestPlanets(ai)
	source, colony := ai.Planets[0], ai.Planets[1]
	source.Cargo.Germanium = haulSourceReserve + haulMinimumLoad - 1
	colony.Spec.MaxFactories = 200
	fleet := newTestTransport(t, ai, 1, source)

	require.NoError(t, ai.haulMinerals())
	assert.Len(t, fleet.Waypoints, 1)
	assert.Zero(t, fleet.Cargo.Germanium)
}

func TestScoutLimitFollowsWhatIsLeftToExplore(t *testing.T) {
	ai := newProductionTestAI()
	for i := range 100 {
		intel := cs.NewPlanet()
		intel.Num = 100 + i
		intel.ReportAge = cs.ReportAgeUnexplored
		ai.PlanetIntels = append(ai.PlanetIntels, intel)
	}
	assert.Equal(t, 100/unexploredPerScout, ai.scoutLimit())

	for _, intel := range ai.PlanetIntels {
		intel.ReportAge = 0
	}
	assert.Equal(t, 2, ai.scoutLimit())
}

func TestWarFleetsWaitWhileOrdersExceedTheBudget(t *testing.T) {
	ai := newProductionTestAI()
	planet := ai.Planets[0]
	planet.Spec.ResourcesPerYear = 1000
	assert.Equal(t, 1, ai.warLimit())

	design := ai.Designs[0]
	design.Spec.Cost = cs.Cost{Resources: 100}
	ai.addShipToTopOfQueue(planet, cs.FleetPurposeBomber, design, 10)
	assert.Zero(t, ai.warLimit(), "1000 on order is more than three years of 30% of 1000")
}

func TestSpareShipyardsRaiseTheColonyLimit(t *testing.T) {
	ai := newProductionTestAI()
	for _, planet := range ai.Planets {
		planet.ProductionQueue = nil
		planet.Spec.MaxFactories, planet.Spec.MaxMines = 100, 100
		planet.Spec.ResourcesPerYearAvailable = 500
	}
	assert.Equal(t, ai.config.colonyPipeline, ai.colonyLimit(), "undeveloped shipyards have no capacity to spare")

	// one shipyard is developed and has nothing on order
	developed := ai.Planets[0]
	developed.Factories, developed.Mines = 100, 100
	assert.Equal(t, ai.config.colonyPipeline+spareColonyShips, ai.colonyLimit())

	// a year's worth of ships on order uses its capacity
	design := ai.Designs[0]
	design.Spec.Cost = cs.Cost{Resources: 100}
	ai.addShipToTopOfQueue(developed, cs.FleetPurposeScout, design, 5)
	assert.Equal(t, ai.config.colonyPipeline, ai.colonyLimit())
}

func TestSpentScoutsAreReplaced(t *testing.T) {
	ai := newProductionTestAI()
	placeTestPlanets(ai)
	for i := range 40 {
		intel := cs.NewPlanet()
		// intel is looked up by planet number
		intel.Num = len(ai.PlanetIntels) + 1
		intel.Position = cs.Vector{X: 500 + 10*i, Y: 500}
		intel.ReportAge = cs.ReportAgeUnexplored
		ai.PlanetIntels = append(ai.PlanetIntels, intel)
	}
	limit := ai.scoutLimit()
	for i := range limit {
		scout := newTestTransport(t, ai, i+1, ai.Planets[0])
		scout.SetTag(cs.TagPurpose, string(cs.FleetPurposeScout))
		// on its way somewhere, with a full tank
		scout.Waypoints = append(scout.Waypoints, cs.NewPlanetWaypoint(ai.PlanetIntels[3+i].Position, ai.PlanetIntels[3+i].Num, "", 5))
	}
	require.NoError(t, ai.scout())
	assert.Zero(t, ai.requests.fleetBuilds[cs.FleetPurposeScout], "at the limit")

	// a scout that has run dry no longer counts, so another is requested
	ai.Fleets[0].Fuel = 0
	ai.requests.fleetBuilds = map[cs.FleetPurpose]int{}
	require.NoError(t, ai.scout())
	assert.Equal(t, 1, ai.requests.fleetBuilds[cs.FleetPurposeScout])
}

func TestMineralsAreHauledForAStarbaseOrder(t *testing.T) {
	ai := newProductionTestAI()
	placeTestPlanets(ai)
	source, colony := ai.Planets[0], ai.Planets[1]
	source.Cargo = cs.Cargo{Ironium: 5000, Boranium: 5000, Germanium: 5000}
	base := cs.NewShipDesign(1, 50).WithPurpose(cs.ShipDesignPurposeFort)
	base.Spec.Cost = cs.Cost{Ironium: 40, Boranium: 30, Resources: 300}
	ai.Designs = append(ai.Designs, base)
	colony.ProductionQueue = []cs.ProductionQueueItem{{Type: cs.QueueItemTypeStarbase, Quantity: 1, DesignNum: base.Num}}
	fleet := newTestTransport(t, ai, 1, source)
	require.GreaterOrEqual(t, fleet.Spec.CargoCapacity, 70)

	require.NoError(t, ai.haulMinerals())

	require.Len(t, fleet.Waypoints, 2)
	assert.Equal(t, colony.Num, fleet.Waypoints[1].TargetNum)
	assert.Equal(t, 40, fleet.Cargo.Ironium)
	assert.Equal(t, 30, fleet.Cargo.Boranium)
	assert.Equal(t, cs.TransportActionUnloadAll, fleet.Waypoints[1].TransportTasks.Ironium.Action)
}
