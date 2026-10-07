package ai

import (
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/stretchr/testify/assert"
)

// a colonist freighter that is already full while orbiting a feeder planet should head to a needy
// planet instead of trying (and failing) to load more colonists
func Test_aiPlayer_transportColonists_fullFreighter(t *testing.T) {
	gamer := cs.NewGamer()
	game := gamer.CreateGame(0, *cs.NewGameSettings().WithAIPlayer(cs.AIDifficultyEasy, 0))
	player := gamer.NewPlayer(0, *cs.NewRace(), &game.Rules)
	player.Num = 1
	universe, err := gamer.GenerateUniverse(game, []*cs.Player{player})
	if err != nil {
		t.Fatalf("gamer.generateUniverse() failed: \n%v", err)
	}

	ai := NewAIPlayer(game, &cs.StaticTechStore, player, universe.GetPlayerMapObjects(player.Num))

	// one full planet to feed colonists from, one empty planet that needs them
	feeder := ai.Planets[0]
	feeder.Spec.PopulationDensity = 1
	feeder.Spec.CanTerraform = false

	needer := cs.NewPlanet()
	needer.Num = feeder.Num + 1
	needer.Name = "Needer"
	needer.PlayerNum = player.Num
	needer.Position = feeder.Position.Add(cs.Vector{X: 50})
	needer.Spec.MaxPopulation = 1_000_000
	ai.Planets = []*cs.Planet{feeder, needer}

	colonistFreighter := cs.NewShipDesign(player.Num, 100).WithHull(cs.SmallFreighter.Name).WithPurpose(cs.ShipDesignPurposeColonistFreighter)
	fuelFreighter := cs.NewShipDesign(player.Num, 101).WithHull(cs.FuelTransport.Name).WithPurpose(cs.ShipDesignPurposeFuelFreighter)
	for _, design := range []*cs.ShipDesign{colonistFreighter, fuelFreighter} {
		spec, err := cs.ComputeShipDesignSpec(&game.Rules, player.TechLevels, player.Race.Spec, design)
		if err != nil {
			t.Fatalf("ComputeShipDesignSpec() failed: \n%v", err)
		}
		design.Spec = spec
	}
	player.Designs = append(player.Designs, colonistFreighter, fuelFreighter)

	fleet := cs.NewFleet(player, 100, "Colonist Freighter", []cs.Waypoint{cs.NewPlanetWaypoint(feeder.Position, feeder.Num, feeder.Name, 5)})
	fleet.Tokens = []cs.ShipToken{
		{DesignNum: colonistFreighter.Num, Quantity: 2},
		{DesignNum: fuelFreighter.Num, Quantity: 1},
	}
	fleet.InjectDesigns(player.Designs)
	fleet.Spec = cs.ComputeFleetSpec(&game.Rules, player, fleet)
	fleet.SetTag(cs.TagPurpose, string(cs.FleetPurposeColonistFreighter))
	fleet.OrbitingPlanetNum = feeder.Num
	fleet.Cargo = cs.Cargo{Colonists: fleet.Spec.CargoCapacity}
	ai.Fleets = []*cs.Fleet{fleet}
	ai.buildMaps()

	if err := ai.transportColonists(); err != nil {
		t.Fatalf("transportColonists() failed: \n%v", err)
	}

	// no colonists loaded, and the fleet is headed to the needy planet to unload
	assert.Equal(t, fleet.Spec.CargoCapacity, fleet.Cargo.Colonists)
	if assert.Len(t, fleet.Waypoints, 2) {
		assert.Equal(t, needer.Num, fleet.Waypoints[1].TargetNum)
		assert.Equal(t, cs.TransportActionUnloadAll, fleet.Waypoints[1].TransportTasks.Colonists.Action)
	}
}
