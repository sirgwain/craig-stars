package ai

import (
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func colonizationFixture(t *testing.T) (*aiPlayer, *cs.Planet, *cs.Planet, *cs.Fleet) {
	t.Helper()
	gamer := cs.NewGamer()
	game := cs.NewGame()
	player := gamer.NewPlayer(0, *cs.NewRace(), &game.Rules)
	player.Num = 1
	u, err := gamer.GenerateUniverse(game, []*cs.Player{player})
	require.NoError(t, err)
	ai := NewAIPlayer(game, &cs.StaticTechStore, player, u.GetPlayerMapObjects(1))
	source := ai.Planets[0]
	var target *cs.Planet
	for _, p := range ai.PlanetIntels {
		if p.Num != source.Num {
			target = p
			break
		}
	}
	require.NotNil(t, target)
	for _, p := range ai.PlanetIntels {
		if p.Num != source.Num {
			p.ReportAge = cs.ReportAgeUnexplored
		}
	}
	target.ReportAge = 0
	target.PlayerNum = cs.Unowned
	target.Spec.Habitability = 80
	target.Spec.MaxPopulation = 800_000
	player.TechLevels.Construction = 4
	design, err := cs.DesignShip(&game.Rules, player, &cs.Privateer, "Colony", 100, 0, cs.ShipDesignPurposeColonizer, cs.FleetPurposeColonizer)
	require.NoError(t, err)
	player.Designs = append(player.Designs, design)
	f := cs.NewFleet(player, 100, "Colony", []cs.Waypoint{cs.NewPlanetWaypoint(source.Position, source.Num, source.Name, 5)})
	f.OrbitingPlanetNum = source.Num
	f.Tokens = []cs.ShipToken{{DesignNum: design.Num, Quantity: 1}}
	f.InjectDesigns(player.Designs)
	f.Spec = cs.ComputeFleetSpec(&game.Rules, player, f)
	f.Fuel = f.Spec.FuelCapacity
	f.SetTag(cs.TagPurpose, string(cs.FleetPurposeColonizer))
	ai.Fleets = []*cs.Fleet{f}
	ai.buildMaps()
	return ai, source, target, f
}

func TestColonizerDispatchesLoadedFleet(t *testing.T) {
	ai, source, target, f := colonizationFixture(t)
	source.Cargo.Colonists = 100 // far below the breeder reserve
	f.Cargo.Colonists = 25
	f.OrbitingPlanetNum = cs.None // also works when stranded in deep space
	require.NoError(t, ai.colonize())
	require.Len(t, f.Waypoints, 2)
	assert.Equal(t, target.Num, f.Waypoints[1].TargetNum)
	assert.Equal(t, cs.WaypointTaskColonize, f.Waypoints[1].Task)
	assert.Equal(t, 25, f.Cargo.Colonists)
}

func TestColonizerLoadsAvailableSurplus(t *testing.T) {
	ai, source, target, f := colonizationFixture(t)
	reserve := int(float64(source.Spec.MaxPopulation)*ai.config.colonizerPopulationDensity) / 100
	source.Cargo.Colonists = reserve + 50
	source.Spec.GrowthAmount = 10_000 // future births must not count as cargo
	require.Greater(t, f.Spec.CargoCapacity, 50)
	require.NoError(t, ai.colonize())
	require.Len(t, f.Waypoints, 2)
	assert.Equal(t, target.Num, f.Waypoints[1].TargetNum)
	assert.Equal(t, 50, f.Cargo.Colonists)
	assert.Equal(t, reserve, source.Cargo.Colonists)
}

func TestColonizerDoesNotQueueWithoutDestinations(t *testing.T) {
	ai, _, target, _ := colonizationFixture(t)
	target.PlayerNum = 2
	require.NoError(t, ai.colonize())
	assert.Zero(t, ai.requests.fleetBuilds[cs.FleetPurposeColonizer])
}

func TestColonizerRetargetsClaimedWorld(t *testing.T) {
	ai, _, target, f := colonizationFixture(t)
	var claimed *cs.Planet
	for _, intel := range ai.PlanetIntels {
		if intel.Num != target.Num && intel.Num != f.OrbitingPlanetNum {
			claimed = intel
			break
		}
	}
	require.NotNil(t, claimed)
	claimed.PlayerNum = 2
	f.Cargo.Colonists = 25
	f.Waypoints = append(f.Waypoints, cs.NewPlanetWaypoint(claimed.Position, claimed.Num, claimed.Name, 5))
	require.NoError(t, ai.colonize())
	require.Len(t, f.Waypoints, 2)
	assert.Equal(t, target.Num, f.Waypoints[1].TargetNum)
}

func TestColonizerSeedsGoodWorldBelowReserve(t *testing.T) {
	ai, source, target, f := colonizationFixture(t)
	source.Cargo.Colonists = minColonistLoad * seedColonySourceFactor // far below the breeder reserve
	target.Spec.Habitability = seedColonyHabitability
	require.NoError(t, ai.colonize())
	require.Len(t, f.Waypoints, 2)
	assert.Equal(t, minColonistLoad, f.Cargo.Colonists)
}

func TestColonizerDoesNotSeedPoorWorld(t *testing.T) {
	ai, source, target, f := colonizationFixture(t)
	source.Cargo.Colonists = minColonistLoad * seedColonySourceFactor
	target.Spec.Habitability = seedColonyHabitability - 1
	require.NoError(t, ai.colonize())
	assert.Len(t, f.Waypoints, 1)
	assert.Zero(t, f.Cargo.Colonists)
}

func TestColonizerWithNowhereToGoUnloads(t *testing.T) {
	ai, source, target, f := colonizationFixture(t)
	target.PlayerNum = 2 // the only destination was claimed
	f.Cargo.Colonists = 25
	before := source.Cargo.Colonists
	require.NoError(t, ai.colonize())
	assert.Zero(t, f.Cargo.Colonists)
	assert.Equal(t, before+25, source.Cargo.Colonists)
}
