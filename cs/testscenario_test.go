//go:build !wasi && !wasm

package cs

import (
	"github.com/stretchr/testify/assert"
	"slices"
	"testing"
)

func newTestPlayerPlanet() (player *Player, planet *Planet) {
	game := BuildScenario(TestScenario{Players: []ScenarioPlayer{{}}, Planets: []ScenarioPlanet{{Name: "Planet 1", Owner: 1, Cargo: Cargo{Colonists: 2500}}}})
	return game.Players[0], game.Planets[0]
}

// create a new long range scout fleet for testing with pre-computed specs.
func testLongRangeScout(player *Player) *Fleet {
	return testLongRangeScoutWithQuantity(player, 1)
}

// create a new long range scout ship design for testing.
// Does NOT come with precomputed specs; those will have to be done manually.
func testLongRangeScoutDesign(playerNum int) *ShipDesign {
	return NewShipDesign(playerNum, 1).
		WithName("Long Range Scout").
		WithHull(Scout.Name).
		WithSlots(slices.Clone(LongRangeScoutTestSlots))
}

func testLongRangeScoutWithQuantity(player *Player, quantity int) *Fleet {
	fleet := &Fleet{
		MapObject: MapObject{Type: MapObjectTypeFleet, Num: 1, PlayerNum: player.Num},
		BaseName:  "Long Range Scout",
		Tokens: []ShipToken{
			{
				Quantity:  quantity,
				DesignNum: 1,
				design: NewShipDesign(player.Num, 1).
					WithName("Long Range Scout").
					WithHull(Scout.Name).
					WithSlots(slices.Clone(LongRangeScoutTestSlots)).
					WithSpec(&rules, player),
			},
		},
		battlePlan:        &player.BattlePlans[0],
		OrbitingPlanetNum: None,
		FleetOrders: FleetOrders{
			Waypoints: []Waypoint{
				NewPositionWaypoint(Vector{}, 5),
			},
		},
	}
	fleet.Spec = ComputeFleetSpec(&rules, player, fleet)
	fleet.Fuel = fleet.Spec.FuelCapacity
	player.Designs = append(player.Designs, fleet.Tokens[0].design)
	return fleet
}

// create a new medium freighter fleet for cargo testing.
func testTeamster(player *Player) *Fleet {
	fleet := &Fleet{
		MapObject: MapObject{
			Type:      MapObjectTypeFleet,
			PlayerNum: player.Num,
		},
		BaseName: "Teamster",
		Tokens: []ShipToken{
			{
				Quantity:  1,
				DesignNum: 1,
				design: NewShipDesign(player.Num, 1).
					WithHull(MediumFreighter.Name).
					WithSlots(slices.Clone(TeamsterTestSlots)).
					WithSpec(&rules, player)},
		},
		battlePlan:        &player.BattlePlans[0],
		OrbitingPlanetNum: None,
		FleetOrders: FleetOrders{
			Waypoints: []Waypoint{
				NewPositionWaypoint(Vector{}, 5),
			},
		},
	}

	fleet.Spec = ComputeFleetSpec(&rules, player, fleet)
	fleet.Fuel = fleet.Spec.FuelCapacity
	return fleet
}

func TestBuildScenario_nameReferences(t *testing.T) {
	s := TestScenario{
		Players: []ScenarioPlayer{
			{
				Designs: Designs(DesignTeamster, DesignLongRangeScout, ShipDesign{Name: "Station", Hull: SpaceStation.Name}),
				Fleets: []ScenarioFleet{
					{
						Name:      "Transport",
						Design:    "Teamster",
						At:        "Home",
						Waypoints: []ScenarioWaypoint{{To: "Home", Task: WaypointTaskTransport}, {To: "Colony", Warp: 7}},
					},
				},
			},
		},
		Planets: []ScenarioPlanet{
			Homeworld("Home", 1),
			{
				Name:            "Colony",
				Position:        Vector{49, 0},
				Owner:           1,
				Starbase:        "Station",
				ProductionQueue: []ScenarioProductionQueueItem{{Type: QueueItemTypeShipToken, Design: "Long Range Scout", Quantity: 2}},
			},
		},
	}
	// Reordering both lists must not redirect any reference.
	slices.Reverse(s.Players[0].Designs)
	slices.Reverse(s.Planets)
	u := newTestUniverse(t, s)
	fleet, home, colony := u.Fleet("Transport"), u.Planet("Home"), u.Planet("Colony")
	assert.Equal(t, home.Num, fleet.OrbitingPlanetNum)
	assert.Equal(t, "Teamster", fleet.Tokens[0].design.Name)
	assert.Equal(t, home.Num, fleet.Waypoints[0].TargetNum)
	assert.Equal(t, colony.Num, fleet.Waypoints[1].TargetNum)
	assert.Equal(t, 7, fleet.Waypoints[1].WarpSpeed)
	assert.Equal(t, "Long Range Scout", colony.ProductionQueue[0].design.Name)
	assert.Equal(t, "Station", colony.Starbase.Tokens[0].design.Name)
	assert.True(t, colony.Spec.HasStarbase)
	assert.Equal(t, colony.Position, colony.Starbase.Position)
	assert.Equal(t, home.Hab, home.BaseHab)
	assert.Equal(t, fleet.Spec.FuelCapacity, fleet.Fuel)
}

func TestBuildScenario_fleetDefaultsAndMixedTokens(t *testing.T) {
	u := newTestUniverse(t, TestScenario{Players: []ScenarioPlayer{{Designs: Designs(DesignTeamster, DesignLongRangeScout), Fleets: []ScenarioFleet{
		{Design: "Teamster"},
		{
			Tokens:    []ScenarioShipToken{{Design: "Teamster", Quantity: 2, QuantityDamaged: 1, Damage: 10}, {Design: "Long Range Scout"}},
			EmptyFuel: true,
		},
	}}}})
	assert.Equal(t, "Teamster #1", u.Game.Fleets[0].Name)
	assert.Equal(t, 1, u.Game.Fleets[0].Tokens[0].Quantity)
	mixed := u.Fleet("Teamster #2")
	assert.Equal(t, 3, mixed.Spec.TotalShips)
	assert.Equal(t, 1, mixed.Tokens[0].QuantityDamaged)
	assert.Equal(t, 10.0, mixed.Tokens[0].Damage)
	assert.Zero(t, mixed.Fuel)
}

func TestBuildScenario_explicitZeroEnvironment(t *testing.T) {
	u := newTestUniverse(t, TestScenario{
		Planets: []ScenarioPlanet{{Name: "Default"}, {Name: "Hostile", Hab: ScenarioValue(Hab{}), Concentration: ScenarioValue(Mineral{})}},
	})
	assert.Equal(t, Hab{50, 50, 50}, u.Planet("Default").Hab)
	assert.Equal(t, Mineral{100, 100, 100}, u.Planet("Default").MineralConcentration)
	assert.Equal(t, Hab{}, u.Planet("Hostile").Hab)
	assert.Equal(t, Hab{}, u.Planet("Hostile").BaseHab)
	assert.Equal(t, Mineral{}, u.Planet("Hostile").MineralConcentration)
}

func TestBuildScenario_reusableTemplates(t *testing.T) {
	s := ScenarioScoutTest()
	s.Players[0].Player = NewPlayer(1, NewRace())
	s.Players[0].Player.AcquiredTechs = map[string]bool{"example": true}
	s.MysteryTraders = []MysteryTrader{{Destination: Vector{100, 0}, PlayersRewarded: map[int]bool{1: true}}}
	first := BuildScenario(s)
	first.Players[0].Designs[0].Slots[0].Quantity = 99
	first.Players[0].AcquiredTechs["changed"] = true
	first.MysteryTraders[0].PlayersRewarded[2] = true
	first.Planets[0].Cargo.Colonists = 1
	second := BuildScenario(s)
	assert.Equal(t, 1, second.Players[0].Designs[0].Slots[0].Quantity)
	assert.NotContains(t, second.Players[0].AcquiredTechs, "changed")
	assert.NotContains(t, second.MysteryTraders[0].PlayersRewarded, 2)
	assert.Equal(t, 2500, second.Planets[0].Cargo.Colonists)
	assert.Equal(t, 1, DesignLongRangeScout.Slots[0].Quantity)
	assert.Equal(t, 0, second.Players[0].TechLevels.Energy)
	assert.Equal(t, first.Seed, second.Seed)
}

func TestBuildScenario_relationsBeforeInitialScan(t *testing.T) {
	s := TestScenario{Players: []ScenarioPlayer{
		{
			Relations: []PlayerRelationship{{Relation: PlayerRelationFriend, ShareMap: true}, {Relation: PlayerRelationFriend, ShareMap: true}},
		},
		{
			Relations: []PlayerRelationship{{Relation: PlayerRelationFriend, ShareMap: true}, {Relation: PlayerRelationFriend, ShareMap: true}},
		},
	}, Planets: []ScenarioPlanet{{Name: "Ally", Owner: 2, Position: Vector{1000, 1000}, Cargo: Cargo{Ironium: 50, Colonists: 2500}}}}
	game := BuildScenario(s)
	assert.Equal(t, PlayerRelationFriend, game.Players[0].Relations[1].Relation)
	assert.Equal(t, game.Planets[0].Cargo.ToMineral(), game.Players[0].GetPlanetIntel(1).Cargo.ToMineral())
	assert.Equal(t, 2, game.Players[0].GetPlanetIntel(1).PlayerNum)
	assert.Equal(t, game.Planets[0].Hab, game.Players[0].GetPlanetIntel(1).Hab)
}

func TestBuildScenario_invalidReferences(t *testing.T) {
	tests := []struct {
		name    string
		change  func(*TestScenario)
		message string
	}{
		{
			"design typo",
			func(s *TestScenario) { s.Players[0].Fleets[0].Design = "Typo" },
			`fleet "Typo #1": unknown design "Typo" for player 1`,
		},
		{
			"planet typo",
			func(s *TestScenario) { s.Players[0].Fleets[0].At = "Typo" },
			`fleet "Long Range Scout #1": unknown planet "Typo"`,
		},
		{
			"waypoint typo",
			func(s *TestScenario) { s.Players[0].Fleets[0].Waypoints = []ScenarioWaypoint{{To: "Typo"}} },
			`fleet "Long Range Scout #1": unknown planet or trader "Typo"`,
		},
		{
			"duplicate planet",
			func(s *TestScenario) { s.Planets = append(s.Planets, s.Planets[0]) },
			`duplicate planet "Planet 1"`,
		},
		{
			"duplicate design",
			func(s *TestScenario) { s.Players[0].Designs = append(s.Players[0].Designs, s.Players[0].Designs[0]) },
			`player 1: duplicate design "Long Range Scout"`,
		},
		{"owner typo", func(s *TestScenario) { s.Planets[0].Owner = 2 }, `planet "Planet 1": unknown owner 2`},
		{
			"starbase typo",
			func(s *TestScenario) { s.Planets[0].Starbase = "Typo" },
			`planet "Planet 1": unknown design "Typo" for player 1`,
		},
		{"queue typo", func(s *TestScenario) {
			s.Planets[0].ProductionQueue = []ScenarioProductionQueueItem{{Type: QueueItemTypeShipToken, Design: "Typo", Quantity: 1}}
		}, `planet "Planet 1": unknown design "Typo" for player 1`},
		{
			"relation count",
			func(s *TestScenario) { s.Players[0].Relations = []PlayerRelationship{} },
			`player 1: relations must contain 1 entries`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SingleUnitScenario()
			tt.change(&s)
			assert.PanicsWithValue(t, tt.message, func() { BuildScenario(s) })
		})
	}
}

func TestSharedScenario_cargoInvasion(t *testing.T) {
	u := newTestUniverse(t, ScenarioCargoTransferInvasion())
	u.TransferByHand(1, "Teamster #1", "Planet 1", Cargo{Colonists: -210})
	u.GenerateTurn()
	assert.Equal(t, 1, u.Planet("Planet 1").PlayerNum)
	assert.NotEmpty(t, u.Messages(1, PlayerMessageFleetInvadedPlanet))
	assert.NotEmpty(t, u.Messages(2, PlayerMessagePlanetInvaded))
}

func TestBuildScenario_fleetBaseName(t *testing.T) {
	u := newTestUniverse(t, TestScenario{Players: []ScenarioPlayer{{Designs: Designs(DesignTeamster), Fleets: []ScenarioFleet{
		{Design: "Teamster", Name: "Hauler #7"},
		{Design: "Teamster", Name: "Hauler"},
		{Design: "Teamster", Name: "Route #A"},
		{Design: "Teamster"},
	}}}})
	assert.Equal(t, "Hauler", u.Fleet("Hauler #7").BaseName)
	assert.Equal(t, "Hauler", u.Fleet("Hauler").BaseName)
	assert.Equal(t, "Route #A", u.Fleet("Route #A").BaseName)
	assert.Equal(t, "Teamster", u.Fleet("Teamster #4").BaseName)
}
