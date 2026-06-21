//go:build !wasi && !wasm

package cs

import "slices"

func newTestPlayerPlanet() (player *Player, planet *Planet) {
	game := BuildScenario(SingleUnitScenario())
	player = game.Players[0]
	planet = game.Planets[0]
	player.Designs = nil
	game.Fleets = nil
	return player, planet
}

func buildSingleTeamsterScenario() *FullGame {
	game := BuildScenario(SingleUnitScenario())
	player := game.Players[0]
	fleet := testTeamster(player).withNum(1)
	player.Designs[0] = fleet.Tokens[0].design
	game.Fleets[0] = fleet
	return game
}

func buildTwoTeamsterScenario() *FullGame {
	game := buildSingleTeamsterScenario()
	player := game.Players[0]
	fleet := testTeamster(player).withNum(2)
	game.Fleets = append(game.Fleets, fleet)
	return game
}

func buildEnemyPlanetTeamsterScenario() *FullGame {
	game := BuildScenario(TwoPlayerScenario())
	player := game.Players[0]
	fleet := testTeamster(player).withNum(1)
	player.Designs[0] = fleet.Tokens[0].design
	game.Fleets[0] = fleet
	return game
}

func buildEnemyPlanetStealingFreighterScenario() *FullGame {
	game := BuildScenario(TwoPlayerScenario())
	player := game.Players[0]
	fleet := testStealingFreighter(player, 1).withNum(1)
	player.Designs[0] = fleet.Tokens[0].design
	game.Fleets[0] = fleet
	return game
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
