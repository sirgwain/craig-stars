//go:build !wasi && !wasm

package cs

import (
	"math"
	"os"
	"slices"
	"testing"

	"log/slog"

	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// many functions require a copy of the current game's rules.
// for testing, create a standard rules var every test can use
var rules = NewRulesWithSeed(0)
var testLogger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
	AddSource: true,
	Level:     slog.LevelDebug,
	ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			return slog.Attr{}
		}
		return a
	},
}))

type MockRand struct {
	int63Result int64
}

func (m MockRand) Seed(seed int64) {}

func (m MockRand) Int63() int64 {
	return m.int63Result
}

// newGeneratedGame generates a real universe like a new game does, with a fixed
// seed so the generated galaxy and every turn after it are the same each run.
func newGeneratedGame(t *testing.T) (Gamer, *Game, *Universe, *Player) {
	t.Helper()
	client := NewGamer()
	game := newSeededGame(*NewGameSettings())
	game.Area, _ = game.Rules.GetArea(game.Size)
	player := client.NewPlayer(1, *NewRace(), &game.Rules)
	player.AIControlled = true
	player.Num = 1
	universe, err := client.GenerateUniverse(game, []*Player{player})
	if err != nil {
		t.Fatal(err)
	}
	return client, game, universe, player
}

func Test_generateTurn(t *testing.T) {
	client, game, universe, player := newGeneratedGame(t)
	players := []*Player{player}

	// build a ship on the planet
	pmo := universe.GetPlayerMapObjects(player.Num)
	planet := pmo.Planets[0]
	planet.ProductionQueue = append([]ProductionQueueItem{{Type: QueueItemTypeShipToken, Quantity: 1, DesignNum: player.Designs[0].Num}}, planet.ProductionQueue...)

	startingFleets := len(universe.Fleets)

	if err := client.GenerateTurn(game, universe, players); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 2401, game.Year)

	// should have intel about planets
	assert.Equal(t, len(universe.Planets), len(player.PlanetIntels))

	// should have built a new scout
	assert.Greater(t, len(universe.Fleets), startingFleets)

	// should have grown pop
	pmo = universe.GetPlayerMapObjects(player.Num)
	assert.Greater(t, pmo.Planets[0].exactPopulation(), player.Race.Spec.StartingPlanets[0].Population)
}

func Test_generateTurns(t *testing.T) {
	client, game, universe, player := newGeneratedGame(t)
	players := []*Player{player}

	// generate many turns
	numTurns := 10
	for i := 0; i < numTurns; i++ {
		if err := client.GenerateTurn(game, universe, players); err != nil {
			t.Fatal(err)
		}
		removeDeleted(universe)
	}

	assert.Equal(t, game.Rules.StartingYear+numTurns, game.Year)

	// should have fleets
	assert.True(t, len(universe.Fleets) > 0)

	// should have grown pop
	pmo := universe.GetPlayerMapObjects(player.Num)
	assert.Greater(t, pmo.Planets[0].exactPopulation(), player.Race.Spec.StartingPlanets[0].Population)

	// should have built factories
	assert.Greater(t, pmo.Planets[0].Factories, player.Race.Spec.StartingPlanets[0].Factories)

	// no victor
	assert.False(t, player.Victor)
	assert.False(t, game.VictorDeclared)
}

func Test_turn_grow(t *testing.T) {
	u := newTestUniverse(t, TestScenario{
		Players: []ScenarioPlayer{{}},
		Planets: []ScenarioPlanet{
			{Name: "Planet 1", Owner: 1, Cargo: Cargo{Colonists: 1000}},
			{Name: "Planet 2", Owner: 1, Hab: ScenarioValue(Hab{}), Cargo: Cargo{Colonists: 1000}},
			{Name: "Planet 3", Owner: 1, Hab: ScenarioValue(Hab{}), Cargo: Cargo{Colonists: 1}},
			{Name: "Planet 4", Owner: 1, Cargo: Cargo{Colonists: 24000}},
		}})
	planet1, planet2 := u.Planet("Planet 1"), u.Planet("Planet 2")
	planet3, planet4 := u.Planet("Planet 3"), u.Planet("Planet 4")
	u.Run((*turnGenerator).planetGrow)

	// #1 grows, #2 dies, #3 dies but can't go lower, #4 overpops a tad bit
	assert.Equal(t, 115_000, planet1.exactPopulation())
	assert.Equal(t, 95_500, planet2.exactPopulation())
	assert.Equal(t, 100, planet3.exactPopulation())
	assert.Equal(t, 2_304_000, planet4.exactPopulation())
}

func Test_turn_fleetByHandUnloads(t *testing.T) {

	t.Run("jettison", func(t *testing.T) {
		u := newTestUniverse(t, TestScenario{
			Players: []ScenarioPlayer{
				{
					Designs: Designs(DesignTeamster),
					Fleets:  []ScenarioFleet{{Design: "Teamster", Position: Vector{10, 10}, Cargo: Cargo{Ironium: 50}}},
				},
			},
			Planets: []ScenarioPlanet{Homeworld("Planet 1", 1)},
		})
		game := u.Game
		u.TransferByHand(1, "Teamster #1", "", Cargo{Ironium: -50})

		u.GenerateTurn()

		assert.Equal(t, 1, len(game.Salvages))
		assert.Equal(t, Vector{10, 10}, game.Salvages[0].Position)
		assert.Equal(t, Cargo{Ironium: 40}, game.Salvages[0].Cargo)
		assert.Equal(t, Cargo{}, game.Fleets[0].Cargo)
	})

	t.Run("unload on our planet", func(t *testing.T) {
		u := newTestUniverse(t, TestScenario{
			Players: []ScenarioPlayer{
				{
					Designs: Designs(DesignTeamster),
					Fleets:  []ScenarioFleet{{Design: "Teamster", At: "Planet 1", Quantity: 2, Cargo: Cargo{Ironium: 50}}},
				},
			},
			Planets: []ScenarioPlanet{{Name: "Planet 1", Owner: 1, Cargo: Cargo{Colonists: 2500}}},
		})
		game := u.Game
		u.TransferByHand(1, "Teamster #1", "Planet 1", Cargo{Ironium: -50})

		u.GenerateTurn()

		assert.Equal(t, Cargo{Ironium: 50}.WithPopulation(game.Planets[0].GetPopulation()), game.Planets[0].Cargo)
		assert.Equal(t, Cargo{}, game.Fleets[0].Cargo)
	})

	t.Run("unload on unowned planet", func(t *testing.T) {
		u := newTestUniverse(t, TestScenario{
			Players: []ScenarioPlayer{
				{
					Designs: Designs(DesignTeamster),
					Fleets:  []ScenarioFleet{{Design: "Teamster", At: "Planet 2", Quantity: 2, Cargo: Cargo{Ironium: 50}}},
				},
			},
			Planets: []ScenarioPlanet{Homeworld("Planet 1", 1), {Name: "Planet 2", Position: Vector{10, 10}}},
		})
		game := u.Game
		u.TransferByHand(1, "Teamster #1", "Planet 2", Cargo{Ironium: -50})

		u.GenerateTurn()

		assert.Equal(t, Cargo{Ironium: 50}, game.Planets[1].Cargo)
		assert.Equal(t, Cargo{}, game.Fleets[0].Cargo)
	})
	t.Run("invade enemy planet", func(t *testing.T) {
		u := newTestUniverse(t, TestScenario{
			Players: []ScenarioPlayer{
				{
					Designs: Designs(DesignTeamster),
					Fleets:  []ScenarioFleet{{Design: "Teamster", At: "Planet 2", Quantity: 20, Cargo: Cargo{Colonists: 5000}}},
				},
				AIPlayer("Player 2"),
			},
			Planets: []ScenarioPlanet{
				{Name: "Planet 1", Owner: 1, Cargo: Cargo{Colonists: 2500}},
				{Name: "Planet 2", Owner: 2, Position: Vector{100, 0}, Cargo: Cargo{Colonists: 2500}},
			},
		})
		game := u.Game
		player := u.Player(1)
		u.TransferByHand(1, "Teamster #1", "Planet 2", Cargo{Colonists: -5000})

		u.GenerateTurn()

		assert.Equal(t, player.Num, game.Planets[1].PlayerNum)
		assert.Equal(t, Cargo{}, game.Fleets[0].Cargo)
		assert.True(t, slices.ContainsFunc(player.Messages, func(pm PlayerMessage) bool { return pm.Type == PlayerMessageFleetInvadedPlanet }))
		assert.True(t, slices.ContainsFunc(game.Players[1].Messages, func(pm PlayerMessage) bool { return pm.Type == PlayerMessagePlanetInvaded }))
	})
}

func Test_turn_fleetByHandLoads(t *testing.T) {

	t.Run("jettison load from another fleet's jettison", func(t *testing.T) {
		s := singleFleetScenario(DesignTeamster)
		s.Players[0].Fleets = append(s.Players[0].Fleets, s.Players[0].Fleets[0])
		s.Players[0].Fleets[0].Position = Vector{10, 10}
		s.Players[0].Fleets[0].Cargo = Cargo{Ironium: 50}
		s.Players[0].Fleets[1].Position = Vector{10, 10}
		u := newTestUniverse(t, s)
		game := u.Game
		u.TransferByHand(1, "Teamster #1", "", Cargo{Ironium: -50})
		u.TransferByHand(1, "Teamster #2", "", Cargo{Ironium: 10})

		u.GenerateTurn()

		// should end up with a salvage from the jettison, but fleet2 also gets some cargo it grabbed
		assert.Equal(t, 1, len(game.Salvages))
		assert.Equal(t, Vector{10, 10}, game.Salvages[0].Position)
		assert.Equal(t, Cargo{Ironium: 30}, game.Salvages[0].Cargo)
		assert.Equal(t, Cargo{}, game.Fleets[0].Cargo)
		assert.Equal(t, Cargo{Ironium: 10}, game.Fleets[1].Cargo)
	})

	t.Run("load from our planet", func(t *testing.T) {
		s := singleFleetScenario(DesignTeamster)
		s.Players[0].Fleets[0].At = "Planet 1"
		s.Planets[0].Cargo.Ironium = 50
		u := newTestUniverse(t, s)
		game := u.Game
		u.TransferByHand(1, "Teamster #1", "Planet 1", Cargo{Ironium: 50})

		u.GenerateTurn()

		assert.Equal(t, Cargo{Ironium: 0}.WithPopulation(game.Planets[0].GetPopulation()), game.Planets[0].Cargo)
		assert.Equal(t, Cargo{Ironium: 50}, game.Fleets[0].Cargo)
	})

	t.Run("load from owned mineral packet", func(t *testing.T) {
		s := singleFleetScenario(DesignTeamster)
		s.Players[0].Fleets[0].Position = Vector{50, 0}
		s.Players[0].MineralPackets = []ScenarioMineralPacket{{Name: "Packet", To: "Planet 1", Position: Vector{50, 0}, WarpSpeed: 5, SafeWarpSpeed: 5, Cargo: Cargo{Ironium: 100}}}
		u := newTestUniverse(t, s)
		game := u.Game
		u.TransferByHand(1, "Teamster #1", "Packet", Cargo{Ironium: 50})

		u.GenerateTurn()

		// mineral packet should have 50 left, fleet should have 50
		assert.Equal(t, Cargo{Ironium: 50}, game.MineralPackets[0].Cargo)
		assert.Equal(t, Cargo{Ironium: 50}, game.Fleets[0].Cargo)
	})

	t.Run("load from enemy mineral packet", func(t *testing.T) {
		s := twoPlayerFleetScenario(DesignTeamster)
		s.Players[0].Fleets[0].Position = Vector{50, 0}
		s.Players[1].MineralPackets = []ScenarioMineralPacket{{Name: "Packet", To: "Planet 1", Position: Vector{50, 0}, WarpSpeed: 5, SafeWarpSpeed: 5, Cargo: Cargo{Ironium: 100}}}
		u := newTestUniverse(t, s)
		game := u.Game
		newDiscoverer(testLogger, u.Player(1)).discoverMineralPacket(&game.Rules, game.MineralPackets[0], u.Player(2), u.Planet("Planet 1"))
		u.TransferByHand(1, "Teamster #1", "Packet", Cargo{Ironium: 50})

		u.GenerateTurn()

		// mineral packet should have 50 left, fleet should have 50
		assert.Equal(t, Cargo{Ironium: 50}, game.MineralPackets[0].Cargo)
		assert.Equal(t, Cargo{Ironium: 50}, game.Fleets[0].Cargo)
	})

	t.Run("load from salvage", func(t *testing.T) {
		s := singleFleetScenario(DesignTeamster)
		s.Players[0].Fleets[0].Position = Vector{50, 0}
		s.Players[0].Salvages = []Salvage{{MapObject: MapObject{Position: Vector{50, 0}}, Cargo: Cargo{Ironium: 100}}}
		u := newTestUniverse(t, s)
		game := u.Game
		newDiscoverer(testLogger, u.Player(1)).discoverSalvage(game.Salvages[0])
		u.TransferByHand(1, "Teamster #1", "Salvage #1", Cargo{Ironium: 50})

		u.GenerateTurn()

		// mineral packet should have 40 left (after decay), fleet should have 50
		assert.Equal(t, Cargo{Ironium: 40}, game.Salvages[0].Cargo)
		assert.Equal(t, Cargo{Ironium: 50}, game.Fleets[0].Cargo)
	})

	t.Run("steal from enemy planet", func(t *testing.T) {
		s := twoPlayerFleetScenario(DesignStealingFreighter)
		s.Players[0].Fleets[0].At = "Planet 2"
		s.Planets[1].Cargo.Ironium = 50
		u := newTestUniverse(t, s)
		game := u.Game
		u.Player(1).GetPlanetIntel(2).Cargo = u.Planet("Planet 2").Cargo
		u.TransferByHand(1, "Stealing Freighter #1", "Planet 2", Cargo{Ironium: 50})

		u.GenerateTurn()

		// should steal
		assert.Equal(t, Cargo{Ironium: 0}.WithPopulation(game.Planets[0].GetPopulation()), game.Planets[1].Cargo)
		assert.Equal(t, Cargo{Ironium: 50}, game.Fleets[0].Cargo)
	})

	t.Run("fail load from enemy planet", func(t *testing.T) {
		s := twoPlayerFleetScenario(DesignTeamster)
		s.Players[0].Fleets[0].At = "Planet 2"
		s.Planets[1].Cargo.Ironium = 50
		u := newTestUniverse(t, s)
		game := u.Game
		u.Player(1).GetPlanetIntel(2).Cargo = u.Planet("Planet 2").Cargo
		u.TransferByHand(1, "Teamster #1", "Planet 2", Cargo{Ironium: 50})
		player := u.Player(1)

		u.GenerateTurn()

		// should fail to steal
		assert.Equal(t, Cargo{Ironium: 50}.WithPopulation(game.Planets[0].GetPopulation()), game.Planets[1].Cargo)
		assert.Equal(t, Cargo{Ironium: 0}, game.Fleets[0].Cargo)
		assert.True(t, slices.ContainsFunc(player.Messages, func(pm PlayerMessage) bool { return pm.Type == PlayerMessageFleetByHandTransferIncomplete }))
	})

	t.Run("load from our fleet", func(t *testing.T) {
		s := singleFleetScenario(DesignTeamster)
		s.Players[0].Fleets = append(s.Players[0].Fleets, s.Players[0].Fleets[0])
		s.Players[0].Fleets[1].Cargo = Cargo{Ironium: 50}
		u := newTestUniverse(t, s)
		game := u.Game
		u.TransferByHand(1, "Teamster #1", "Teamster #2", Cargo{Ironium: 50})

		u.GenerateTurn()

		assert.Equal(t, Cargo{Ironium: 50}, game.Fleets[0].Cargo)
		assert.Equal(t, Cargo{Ironium: 0}, game.Fleets[1].Cargo)
	})

	t.Run("steal from enemy fleet", func(t *testing.T) {
		s := twoPlayerFleetScenario(DesignStealingFreighter)
		s.Players[0].Fleets[0].At = "Planet 2"
		s.Planets[1].Cargo.Ironium = 50
		u := newTestUniverse(t, s)
		game := u.Game
		u.Player(1).GetPlanetIntel(2).Cargo = u.Planet("Planet 2").Cargo
		u.TransferByHand(1, "Stealing Freighter #1", "Planet 2", Cargo{Ironium: 50})

		u.GenerateTurn()

		// should steal
		assert.Equal(t, Cargo{Ironium: 0}.WithPopulation(game.Planets[0].GetPopulation()), game.Planets[1].Cargo)
		assert.Equal(t, Cargo{Ironium: 50}, game.Fleets[0].Cargo)
	})
}

func Test_turn_fleetTransferCargoInvade1(t *testing.T) {
	s := TwoPlayerScenario()
	s.Players[0].Player = NewPlayer(1, NewRace().WithPluralName("Attackers"))
	s.Players[1].Player = AIPlayer("Player 2").Player
	s.Players[1].Player.Race.PluralName = "Defenders"
	s.Players[0].Fleets[0].At = "Planet 2"
	s.Players[0].Fleets[0].Cargo = Cargo{Colonists: 5000}
	s.Players[0].Fleets[0].Waypoints = []ScenarioWaypoint{
		{
			To:             "Planet 2",
			Warp:           5,
			Task:           WaypointTaskTransport,
			TransportTasks: WaypointTransportTasks{Colonists: WaypointTransportTask{Action: TransportActionUnloadAll}},
		},
	}

	u := newTestUniverse(t, s)
	player1, player2, fleet, planet := u.Player(1), u.Player(2), u.FleetFor(1, "Long Range Scout #1"), u.Planet("Planet 2")

	// transfer
	initialDefenders := planet.GetPopulation()
	u.GenerateTurn()
	for _, message := range append(u.Messages(1, PlayerMessageFleetInvadedPlanet), u.Messages(2, PlayerMessagePlanetInvaded)...) {
		assert.Equal(t, 500000, message.Spec.Invasion.Attackers)
		assert.Equal(t, initialDefenders, message.Spec.Invasion.Defenders)
	}

	// should have invaded the planet and taken it over
	assert.Equal(t, planet.PlayerNum, fleet.PlayerNum)
	assert.True(t, slices.ContainsFunc(player1.Messages, func(message PlayerMessage) bool {
		return message.Type == PlayerMessageFleetInvadedPlanet
	}))
	assert.True(t, slices.ContainsFunc(player2.Messages, func(message PlayerMessage) bool {
		return message.Type == PlayerMessagePlanetInvaded
	}))

	// make sure player one knows they lost their planet
	assert.True(t, slices.ContainsFunc(player2.PlanetIntels, func(p *Planet) bool {
		return p.Num == planet.Num && p.PlayerNum == player1.Num
	}))

}

func Test_turn_fleetTransferCargoInvadeStarbase(t *testing.T) {
	s := TwoPlayerScenario()
	s.Players[0].Player = NewPlayer(1, NewRace().WithPluralName("Attackers"))
	s.Players[1].Player = AIPlayer("Player 2").Player
	s.Players[1].Player.Race.PluralName = "Defenders"
	s.Players[0].Fleets[0].At = "Planet 2"
	s.Players[0].Fleets[0].Cargo = Cargo{Colonists: 5000}
	s.Players[0].Fleets[0].Waypoints = []ScenarioWaypoint{
		{
			To:             "Planet 2",
			Warp:           5,
			Task:           WaypointTaskTransport,
			TransportTasks: WaypointTransportTasks{Colonists: WaypointTransportTask{Action: TransportActionUnloadAll}},
		},
	}
	s.Players[1].Designs = append(s.Players[1].Designs, ShipDesign{Name: "Starbase", Hull: SpaceStation.Name})
	s.Planets[1].Starbase = "Starbase"

	u := newTestUniverse(t, s)
	player1, player2, fleet, planet := u.Player(1), u.Player(2), u.FleetFor(1, "Long Range Scout #1"), u.Planet("Planet 2")
	numInvaders := 5000

	// transfer
	u.GenerateTurn()

	// should not have invaded the planet because we have a starbase
	assert.Equal(t, planet.PlayerNum, player2.Num)
	assert.Equal(t, numInvaders, fleet.Cargo.Colonists)
	assert.False(t, slices.ContainsFunc(player1.Messages, func(message PlayerMessage) bool {
		return message.Type == PlayerMessageFleetInvadedPlanet
	}))
	assert.False(t, slices.ContainsFunc(player2.Messages, func(message PlayerMessage) bool {
		return message.Type == PlayerMessagePlanetInvaded
	}))
}

func Test_turn_fleetRoute(t *testing.T) {
	s := SingleUnitScenario()
	s.Planets = append(s.Planets, ScenarioPlanet{Name: "Planet 2", Position: Vector{10, 10}})
	s.Planets[0].RouteTo = "Planet 2"
	s.Players[0].Fleets[0].Waypoints = []ScenarioWaypoint{{To: "Planet 1", Task: WaypointTaskRoute}}
	u := newTestUniverse(t, s)
	player, target, fleet := u.Player(1), u.Planet("Planet 2"), u.Fleet("Long Range Scout #1")

	// route to planet 2
	// move
	u.Run((*turnGenerator).fleetRoute)

	assert.Equal(t, 2, len(fleet.Waypoints))
	assert.Equal(t, target.Num, fleet.Waypoints[1].TargetNum)
	assert.Equal(t, target.Type, fleet.Waypoints[1].TargetType)
	assert.Equal(t, target.PlayerNum, fleet.Waypoints[1].TargetPlayerNum)
	assert.Equal(t, 1, len(player.Messages))
}

func Test_turn_fleetMove(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		s := SingleUnitScenario()
		s.Players[0].Fleets[0].Waypoints = []ScenarioWaypoint{{To: "Planet 1"}, {Position: Vector{10, 10}, Warp: 5}}
		u := newTestUniverse(t, s)
		game, planet, fleet := u.Game, u.Planet("Planet 1"), u.Fleet("Long Range Scout #1")

		// move to place
		u.Run((*turnGenerator).fleetMove)

		// should have consumed that waypoint and moved to the space
		assert.Equal(t, 1, len(fleet.Waypoints))
		assert.Equal(t, Vector{10, 10}, fleet.Position)
		assert.Equal(t, 1, len(game.getMapObjectsAtPosition(planet.Position)))
		assert.Equal(t, 1, len(game.getMapObjectsAtPosition(fleet.Position)))
	})

	t.Run("Repeat Orders", func(t *testing.T) {
		s := singleFleetScenario(DesignSmallFreighter)
		s.Planets[0].Cargo = Cargo{1000, 1000, 1000, 1000}
		s.Players[0].Fleets = []ScenarioFleet{{Design: "Small Freighter", At: "Planet 1", RepeatOrders: true, Waypoints: []ScenarioWaypoint{
			{
				To:             "Planet 1",
				Warp:           5,
				Task:           WaypointTaskTransport,
				TransportTasks: WaypointTransportTasks{Ironium: WaypointTransportTask{Action: TransportActionLoadAll}},
			},
			{
				Position:       Vector{50, 0},
				Warp:           5,
				Task:           WaypointTaskTransport,
				TransportTasks: WaypointTransportTasks{Ironium: WaypointTransportTask{Action: TransportActionUnloadAll}},
			},
		}}}
		u := newTestUniverse(t, s)
		game, planet, fleet := u.Game, u.Planet("Planet 1"), u.Fleet("Small Freighter #1")

		// move one year
		u.GenerateTurn()

		// should have loaded, moved, but still have waypoints
		assert.Equal(t, 120, fleet.Cargo.Ironium)
		assert.Equal(t, 880, planet.Cargo.Ironium)
		assert.Equal(t, Vector{25, 0}, fleet.Position)
		assert.Equal(t, 3, len(fleet.Waypoints))

		// generate the second turn, should move to dest and unload
		u.GenerateTurn()
		assert.Equal(t, 0, fleet.Cargo.Ironium)
		assert.Equal(t, Vector{50, 0}, fleet.Position)
		assert.Equal(t, 2, len(fleet.Waypoints))
		// should have created salvage with ironium drop
		salvage := game.Salvages[0]
		assert.Equal(t, 120, salvage.Cargo.Ironium)

		// generate the third turn, should move back towards planet
		u.GenerateTurn()
		assert.Equal(t, Vector{25, 0}, fleet.Position)
		assert.Equal(t, 3, len(fleet.Waypoints))

		// generate a fourth turn, should arrive at planet and load ironium
		u.GenerateTurn()
		assert.Equal(t, Vector{0, 0}, fleet.Position)
		assert.Equal(t, 120, fleet.Cargo.Ironium)
		assert.Equal(t, 760, planet.Cargo.Ironium)
		assert.Equal(t, 2, len(fleet.Waypoints))

		// generate a fifth turn, should move again towards dest
		u.GenerateTurn()
		assert.Equal(t, Vector{25, 0}, fleet.Position)
		assert.Equal(t, 3, len(fleet.Waypoints))

	})

	t.Run("TransportRepeat", func(t *testing.T) {
		s := singleFleetScenario(DesignGalleon)
		s.Planets = []ScenarioPlanet{
			{Name: "Planet 1", Owner: 1, Cargo: Cargo{1000, 1000, 1000, 10000}, Mines: 0},
			{
				Name:          "Planet 2",
				Owner:         1,
				Position:      Vector{10, 0},
				Cargo:         Cargo{Colonists: 25},
				Concentration: ScenarioValue(Mineral{}),
			},
		}
		s.Players[0].Fleets = []ScenarioFleet{{Design: "Galleon", At: "Planet 1", RepeatOrders: true, Waypoints: []ScenarioWaypoint{
			{
				To:             "Planet 1",
				Warp:           5,
				Task:           WaypointTaskTransport,
				TransportTasks: WaypointTransportTasks{Colonists: WaypointTransportTask{Action: TransportActionLoadAll}},
			},
			{
				To:             "Planet 2",
				Warp:           5,
				Task:           WaypointTaskTransport,
				TransportTasks: WaypointTransportTasks{Colonists: WaypointTransportTask{Action: TransportActionSetWaypointTo, Amount: 2500}},
			},
		}}}
		u := newTestUniverse(t, s)
		planet1, planet2, fleet := u.Planet("Planet 1"), u.Planet("Planet 2"), u.Fleet("Galleon #1")

		// move one year
		u.GenerateTurn()

		// should have loaded, moved & dropped pop
		assert.Equal(t, 9150, planet1.Cargo.Colonists) // 10000-1000+150
		assert.Equal(t, Vector{10, 0}, fleet.Position)
		assert.Equal(t, Vector{10, 0}, fleet.Waypoints[0].Position)
		assert.Equal(t, MapObjectTypePlanet, fleet.Waypoints[0].TargetType)
		assert.Equal(t, planet2.Num, fleet.Waypoints[0].TargetNum)
		assert.Equal(t, 1028, planet2.Cargo.Colonists)
		assert.Equal(t, 75, planet2.PartialPopulation)
		assert.Equal(t, 2, len(fleet.Waypoints))

		// generate the second turn, should move back to planet 1
		u.GenerateTurn()

		// should have arrived back at homeworld & loaded more pop
		assert.Equal(t, 8287, planet1.Cargo.Colonists)
		assert.Equal(t, Vector{0, 0}, fleet.Position)
		assert.Equal(t, Vector{0, 0}, fleet.Waypoints[0].Position)
		assert.Equal(t, MapObjectTypePlanet, fleet.Waypoints[0].TargetType)
		assert.Equal(t, planet1.Num, fleet.Waypoints[0].TargetNum)
		assert.Equal(t, 1000, fleet.Cargo.Colonists)
		assert.Equal(t, MapObjectTypePlanet, fleet.Waypoints[1].TargetType)
		assert.Equal(t, planet2.Num, fleet.Waypoints[1].TargetNum)
		assert.Equal(t, 2, len(fleet.Waypoints))

		// generate the third turn, should move back to planet2 and unload
		u.GenerateTurn()

		assert.Equal(t, 8499, planet1.Cargo.Colonists)
		assert.Equal(t, Vector{10, 0}, fleet.Position)
		assert.Equal(t, Vector{10, 0}, fleet.Waypoints[0].Position)
		assert.Equal(t, MapObjectTypePlanet, fleet.Waypoints[0].TargetType)
		assert.Equal(t, planet2.Num, fleet.Waypoints[0].TargetNum)
		assert.Equal(t, Cargo{}, fleet.Cargo)
		assert.Equal(t, 2360, planet2.Cargo.Colonists)
		assert.Equal(t, 2, len(fleet.Waypoints))

		// generate a couple more turns, we should eventually stop unloading cargo due to the SetAmountTo and growth
		// p2 -> p1
		u.GenerateTurn()
		// p1 -> p2
		u.GenerateTurn()

		assert.Equal(t, 7956, planet1.Cargo.Colonists)
		assert.Equal(t, Vector{10, 0}, fleet.Position)
		assert.Equal(t, Vector{10, 0}, fleet.Waypoints[0].Position)
		assert.Equal(t, MapObjectTypePlanet, fleet.Waypoints[0].TargetType)
		assert.Equal(t, planet2.Num, fleet.Waypoints[0].TargetNum)
		assert.Equal(t, Cargo{Colonists: 1000}, fleet.Cargo) // we have leftover
		assert.Equal(t, 3121, planet2.Cargo.Colonists)       // planet is ready to go!
		assert.Equal(t, 2, len(fleet.Waypoints))

	})

	t.Run("Wait for %", func(t *testing.T) {
		s := singleFleetScenario(DesignGalleon)
		s.Planets = []ScenarioPlanet{
			{Name: "Planet 1", Owner: 1, Cargo: Cargo{100, 100, 100, 10000}, Mines: 300},
			{
				Name:          "Planet 2",
				Owner:         1,
				Position:      Vector{10, 0},
				Cargo:         Cargo{Colonists: 1000},
				Concentration: ScenarioValue(Mineral{}),
			},
		}
		s.Players[0].Fleets = []ScenarioFleet{{Design: "Galleon", At: "Planet 1", RepeatOrders: true, Waypoints: []ScenarioWaypoint{
			{
				To:   "Planet 1",
				Warp: 5,
				Task: WaypointTaskTransport,
				TransportTasks: WaypointTransportTasks{
					Ironium:   WaypointTransportTask{Action: TransportActionWaitForPercent, Amount: 33},
					Boranium:  WaypointTransportTask{Action: TransportActionWaitForPercent, Amount: 33},
					Germanium: WaypointTransportTask{Action: TransportActionWaitForPercent, Amount: 34},
				},
			},
			{
				To:   "Planet 2",
				Warp: 5,
				Task: WaypointTaskTransport,
				TransportTasks: WaypointTransportTasks{
					Ironium:   WaypointTransportTask{Action: TransportActionUnloadAll},
					Boranium:  WaypointTransportTask{Action: TransportActionUnloadAll},
					Germanium: WaypointTransportTask{Action: TransportActionUnloadAll},
				},
			},
		}}}
		u := newTestUniverse(t, s)
		planet1, planet2, fleet := u.Planet("Planet 1"), u.Planet("Planet 2"), u.Fleet("Galleon #1")

		// load and grow and wait
		u.GenerateTurn()

		// should have loaded all cargo, waited, then loaded what was mined after production
		assert.Equal(t, Vector{0, 0}, fleet.Position)
		assert.Equal(t, Vector{0, 0}, fleet.Waypoints[0].Position)
		assert.Equal(t, MapObjectTypePlanet, fleet.Waypoints[0].TargetType)
		assert.Equal(t, planet1.Num, fleet.Waypoints[0].TargetNum)
		assert.Equal(t, Cargo{330, 330, 340, 0}, fleet.Cargo)
		assert.Equal(t, 2, len(fleet.Waypoints))

		// our hold is full, so we move
		u.GenerateTurn()

		// should have loaded all cargo, and moved to planet2 to dump
		assert.Equal(t, Vector{10, 0}, fleet.Position)
		assert.Equal(t, Vector{10, 0}, fleet.Waypoints[0].Position)
		assert.Equal(t, MapObjectTypePlanet, fleet.Waypoints[0].TargetType)
		assert.Equal(t, planet2.Num, fleet.Waypoints[0].TargetNum)
		assert.Equal(t, Cargo{0, 0, 0, 0}, fleet.Cargo)
		assert.Equal(t, Mineral{330, 330, 340}, planet2.Cargo.ToMineral())
		assert.Equal(t, 2, len(fleet.Waypoints))
		// go back and load again from p1
		assert.Equal(t, Vector{0, 0}, fleet.Waypoints[1].Position)
		assert.Equal(t, MapObjectTypePlanet, fleet.Waypoints[1].TargetType)
		assert.Equal(t, planet1.Num, fleet.Waypoints[1].TargetNum)

	})

	t.Run("Stopped by minefield", func(t *testing.T) {
		s := SingleUnitScenario()
		s.Players[0].Player = NewPlayer(1, NewRace()).WithTechLevels(TechLevel{26, 26, 26, 26, 26, 26})
		s.Players = append(s.Players, ScenarioPlayer{
			Minefields: []Minefield{{MapObject: MapObject{Position: Vector{20, 0}}, MinefieldType: MinefieldTypeStandard, NumMines: 100}},
		})
		s.Players[0].Fleets[0].Waypoints = []ScenarioWaypoint{{To: "Planet 1"}, {Position: Vector{36, 0}, Warp: 6}}
		u := newTestUniverse(t, s)
		game, fleet, minefield := u.Game, u.Fleet("Long Range Scout #1"), u.Game.Minefields[0]
		stats := game.Rules.MinefieldStatsByType[MinefieldTypeStandard]
		stats.MaxSpeed, stats.ChanceOfHit, stats.MinDecay = 5, 1, 0
		game.Rules.MinefieldStatsByType[MinefieldTypeStandard] = stats

		// let's go!!
		u.GenerateTurn()

		// we should have struck the minefield and lost the ship
		assert.True(t, fleet.Delete)
		assert.Len(t, u.Messages(1, PlayerMessageFleetMinefieldHit), 1)
		assert.Len(t, u.Messages(2, PlayerMessageFleetMinefieldHit), 1)

		// the Minefield should have lost some mines in the collision
		assert.Equal(t, 88, minefield.NumMines)
		assert.Equal(t, Vector{10, 0}, fleet.Position)
	})

	t.Run("Destroyed by minefield", func(t *testing.T) {
		s := SingleUnitScenario()
		s.Players[0].Player = NewPlayer(1, NewRace()).WithTechLevels(TechLevel{26, 26, 26, 26, 26, 26})
		s.Players = append(s.Players, ScenarioPlayer{
			Minefields: []Minefield{{MapObject: MapObject{Position: Vector{20, 0}}, MinefieldType: MinefieldTypeStandard, NumMines: 100}},
		})
		s.Players[0].Fleets[0].Waypoints = []ScenarioWaypoint{{To: "Planet 1"}, {Position: Vector{81, 0}, Warp: 9}}
		u := newTestUniverse(t, s)
		game, fleet, minefield := u.Game, u.Fleet("Long Range Scout #1"), u.Game.Minefields[0]
		stats := game.Rules.MinefieldStatsByType[MinefieldTypeStandard]
		stats.MaxSpeed, stats.ChanceOfHit, stats.MinDecay = 5, .25, 0
		game.Rules.MinefieldStatsByType[MinefieldTypeStandard] = stats

		// let's go!!
		u.GenerateTurn()

		// we should have struck the minefield and lost the ship
		assert.True(t, fleet.Delete)
		assert.Len(t, u.Messages(1, PlayerMessageFleetMinefieldHit), 1)
		assert.Len(t, u.Messages(2, PlayerMessageFleetMinefieldHit), 1)

		// the Minefield should have lost some mines in the collision
		assert.Equal(t, 88, minefield.NumMines)
	})
}

func Test_turn_permaform(t *testing.T) {
	s := SingleUnitScenario()
	s.Players[0].Player = NewPlayer(1, NewRace().WithPRT(CA))
	s.Planets[0].Hab = ScenarioValue(Hab{49, 49, 49})
	u := newTestUniverse(t, s)
	game, planet := u.Game, u.Planet("Planet 1")

	// mock the random number generator to return temp as the hab to permaform
	rng := testRandom{}
	rng.addFloats(.1) // permaform chance
	rng.addInts(1)    // permaform temp
	game.Rules.random = &rng

	u.Run((*turnGenerator).permaform)

	// should have permaformed the planet temp in one direction
	assert.Equal(t, Hab{49, 50, 49}, planet.Hab)
}

func Test_turn_permaformNone(t *testing.T) {
	s := SingleUnitScenario()
	s.Players[0].Player = NewPlayer(1, NewRace().WithPRT(CA))
	s.Planets[0].Hab = ScenarioValue(Hab{49, 49, 49})
	u := newTestUniverse(t, s)
	game, planet := u.Game, u.Planet("Planet 1")

	// mock the random number generator to return temp as the hab to permaform
	rng := testRandom{}
	rng.addFloats(.2) // no permaform
	game.Rules.random = &rng

	u.Run((*turnGenerator).permaform)

	// should have permaformed the planet temp in one direction
	assert.Equal(t, Hab{49, 49, 49}, planet.Hab)
}

func remoteMiningScenario(at string, owner int, design ShipDesign, task WaypointTask, prt PRT) TestScenario {
	s := singleFleetScenario(design)
	s.Players[0].Player = NewPlayer(1, NewRace().WithPRT(prt))
	s.Players = append(s.Players, AIPlayer("Player 2"))
	s.Planets = append(s.Planets, ScenarioPlanet{Name: "Mine", Owner: owner})
	s.Players[0].Fleets[0].At = at
	s.Players[0].Fleets[0].Quantity = 2
	s.Players[0].Fleets[0].Waypoints = []ScenarioWaypoint{{To: at, Task: task}}
	return s
}

func Test_turn_fleetRemoteMine(t *testing.T) {
	tests := []struct {
		name            string
		scenario        TestScenario
		wantCargo       Cargo
		wantMessageType PlayerMessageType
	}{
		{
			"no task, do nothing",
			remoteMiningScenario("", 0, DesignLongRangeScout, WaypointTaskNone, JoaT),
			Cargo{},
			PlayerMessageNone,
		},
		{
			"no planet, invalid message",
			remoteMiningScenario("", 0, DesignLongRangeScout, WaypointTaskRemoteMining, JoaT),
			Cargo{},
			PlayerMessageFleetRemoteMineInvalidDeepSpace,
		},
		{
			"owned planet, invalid message",
			remoteMiningScenario("Mine", 2, DesignLongRangeScout, WaypointTaskRemoteMining, JoaT),
			Cargo{},
			PlayerMessageFleetRemoteMineInvalidInhabited,
		},
		{
			"owned by us, invalid",
			remoteMiningScenario("Mine", 1, DesignLongRangeScout, WaypointTaskRemoteMining, JoaT),
			Cargo{},
			PlayerMessageFleetRemoteMineInvalidInhabited,
		},
		{
			"owned by us, but we can remote mine our own, should skip",
			remoteMiningScenario("Mine", 1, DesignLongRangeScout, WaypointTaskRemoteMining, AR),
			Cargo{},
			PlayerMessageNone,
		},
		{
			"no miners, invalid message",
			remoteMiningScenario("Mine", 0, DesignLongRangeScout, WaypointTaskRemoteMining, JoaT),
			Cargo{},
			PlayerMessageFleetRemoteMineInvalidNoMiners,
		},
		{
			"should mine",
			remoteMiningScenario("Mine", 0, DesignPotatoBug, WaypointTaskRemoteMining, JoaT),
			Cargo{10, 10, 10, 0},
			PlayerMessageFleetRemoteMined,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := newTestUniverse(t, tt.scenario)
			u.Run((*turnGenerator).fleetRemoteMine)
			if tt.wantMessageType != PlayerMessageNone {
				assert.Len(t, u.Player(1).Messages, 1)
				assert.Len(t, u.Messages(1, tt.wantMessageType), 1)
			} else {
				assert.Empty(t, u.Player(1).Messages)
			}
			assert.Equal(t, tt.wantCargo, u.Planet("Mine").Cargo)
		})
	}
}

func Test_turn_fleetRemoteMineAR(t *testing.T) {
	tests := []struct {
		name            string
		scenario        TestScenario
		wantCargo       Cargo
		wantMessageType PlayerMessageType
	}{
		{
			"no task, do nothing",
			remoteMiningScenario("", 0, DesignLongRangeScout, WaypointTaskNone, AR),
			Cargo{},
			PlayerMessageNone,
		},
		{
			"no planet, invalid message",
			remoteMiningScenario("", 0, DesignLongRangeScout, WaypointTaskRemoteMining, AR),
			Cargo{},
			PlayerMessageFleetRemoteMineInvalidDeepSpace,
		},
		{
			"no miners, invalid message",
			remoteMiningScenario("Mine", 1, DesignLongRangeScout, WaypointTaskRemoteMining, AR),
			Cargo{},
			PlayerMessageFleetRemoteMineInvalidNoMiners,
		},
		{
			"owned by us, but we can remote mine our own, should mine",
			remoteMiningScenario("Mine", 1, DesignPotatoBug, WaypointTaskRemoteMining, AR),
			Cargo{10, 10, 10, 0},
			PlayerMessageFleetRemoteMined,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := newTestUniverse(t, tt.scenario)
			u.Run((*turnGenerator).fleetRemoteMineAR)
			if tt.wantMessageType != PlayerMessageNone {
				assert.Len(t, u.Player(1).Messages, 1)
				assert.Len(t, u.Messages(1, tt.wantMessageType), 1)
			} else {
				assert.Empty(t, u.Player(1).Messages)
			}
			assert.Equal(t, tt.wantCargo, u.Planet("Mine").Cargo)
		})
	}
}

func Test_turn_fleetLayMines(t *testing.T) {
	s := singleFleetScenario(DesignMiniMineLayer)
	s.Players[0].Fleets[0].Waypoints = []ScenarioWaypoint{{Task: WaypointTaskLayMinefield, Warp: 5}}
	u := newTestUniverse(t, s)
	game := u.Game

	// run for one year; should have created newminefield
	u.GenerateTurn()

	assert.Equal(t, 1, len(game.Minefields))
	minefield := game.Minefields[0]
	assert.Equal(t, 320, minefield.NumMines)
	assert.Equal(t, math.Sqrt(320), minefield.Radius())
	assert.Equal(t, Vector{0, 0}, minefield.Position)

}

func Test_turn_fleetSweepMines(t *testing.T) {
	s := singleFleetScenario(DesignStalwartDefender)
	s.Players = append(s.Players, ScenarioPlayer{Minefields: []Minefield{{MinefieldType: MinefieldTypeStandard, NumMines: 100}}})
	u := newTestUniverse(t, s)
	game, fleet, minefield := u.Game, u.Fleet("Stalwart Defender #1"), u.Game.Minefields[0]
	rules := &game.Rules
	stats := rules.MinefieldStatsByType[MinefieldTypeStandard]
	stats.MinDecay = 0
	rules.MinefieldStatsByType[MinefieldTypeStandard] = stats

	// sweep mines
	u.GenerateTurn()

	// we should clear out some mines
	assert.Equal(t, 78, minefield.NumMines)
	assert.False(t, minefield.Delete)

	// upgrade a mine sweeper weapon
	fleet.Tokens[0].design.Slots[1].HullComponent = GatlingNeutrinoCannon.Name
	u.Recompute()

	// sweep mines
	u.GenerateTurn()

	// we should clear out some mines
	assert.Equal(t, 0, minefield.NumMines)
	assert.True(t, minefield.Delete)
}

func Test_turn_instaform(t *testing.T) {
	s := SingleUnitScenario()
	s.Players[0].Player = NewPlayer(1, NewRace().WithPRT(CA)).WithTechLevels(TechLevel{Propulsion: 1, Biotechnology: 1})
	s.Planets[0].Hab = ScenarioValue(Hab{45, 50, 50})
	s.Planets[0].BaseHab = ScenarioValue(Hab{45, 50, 50})
	u := newTestUniverse(t, s)
	planet := u.Planet("Planet 1")

	// instaform
	u.GenerateTurn()

	// should terraform 3 grav points
	// TODO: sometimes this fails if we randomly permaform...
	assert.Equal(t, Hab{48, 50, 50}, planet.Hab)
}

func Test_turn_instaformTakenPlanet(t *testing.T) {
	s := SingleUnitScenario()
	s.Players[0].Player = NewPlayer(1, NewRace().WithPRT(CA)).WithTechLevels(TechLevel{Propulsion: 1, Biotechnology: 1})
	s.Planets[0].Hab = ScenarioValue(Hab{40, 50, 50})
	s.Planets[0].BaseHab = ScenarioValue(Hab{45, 50, 50})
	u := newTestUniverse(t, s)
	planet := u.Planet("Planet 1")

	// instaform
	u.GenerateTurn()

	// should terraform 3 grav points
	// TODO: sometimes this fails if we randomly permaform...
	assert.Equal(t, Hab{48, 50, 50}, planet.Hab)
}

func Test_turn_fleetRepair(t *testing.T) {
	s := SingleUnitScenario()
	s.Players[0].Designs = Designs(DesignLongRangeScout, ShipDesign{Name: "Starbase", Hull: SpaceStation.Name})
	s.Planets[0].Starbase = "Starbase"
	s.Planets[0].StarbaseDamage = 400
	s.Players[0].Fleets = []ScenarioFleet{{Tokens: []ScenarioShipToken{{Design: "Long Range Scout", QuantityDamaged: 1, Damage: 10}}, Name: "Scout", Fuel: 300}}

	u := newTestUniverse(t, s)
	fleet := u.Game.Fleets[0]
	starbase := u.Planet("Planet 1").Starbase

	// in space with 10 damage

	// repair
	u.GenerateTurn()

	// should repair fleet and starbase
	assert.Equal(t, 9.0, fleet.Tokens[0].Damage)
	assert.Equal(t, 1, fleet.Tokens[0].QuantityDamaged)
	assert.Equal(t, 350.0, starbase.Tokens[0].Damage)
	assert.Equal(t, 1, starbase.Tokens[0].QuantityDamaged)

}

func Test_turn_colonizeAfterTravel(t *testing.T) {
	tests := []struct {
		name       string
		scenario   TestScenario
		population int
		starbase   bool
	}{
		{"normal colonizer", ScenarioColonizerTest(), 2500, false},
		{"AR colonizer", ScenarioColonizerTestAR(), 2400, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := newTestUniverse(t, tt.scenario)
			fleet := u.Fleet("Santa Maria #1")
			planet := u.Planet("Planet 2")
			u.TransferByHand(1, fleet.Name, "Planet 1", Cargo{Colonists: 25})
			wp := NewPlanetWaypoint(planet.Position, planet.Num, planet.Name, 5)
			wp.Task = WaypointTaskColonize
			fleet.Waypoints = append(fleet.Waypoints, wp)

			u.GenerateTurn()

			// AR warp dieoff happens before arrival colonization: 3% of 25kT
			// rounds to 1kT lost. The new colony doesn't grow until next turn.
			assert.Equal(t, 1, planet.PlayerNum)
			assert.Equal(t, tt.population, planet.GetPopulation())
			assert.True(t, fleet.Delete)
			assert.Equal(t, tt.starbase, planet.Spec.HasStarbase)
			if tt.starbase {
				assert.Equal(t, "Starter Colony", planet.Spec.StarbaseDesignName)
			}
		})
	}
}

func Test_turn_fleetReproduce(t *testing.T) {
	u := newTestUniverse(t, TestScenario{Players: []ScenarioPlayer{
		{
			Player:  NewPlayer(1, NewRace().WithPRT(IS)),
			Designs: Designs(DesignSmallFreighter),
			Fleets: []ScenarioFleet{
				{
					Design:    "Small Freighter",
					At:        "Planet 1",
					Cargo:     Cargo{Colonists: 50},
					Waypoints: []ScenarioWaypoint{{To: "Planet 1", Warp: 5}},
				},
			},
		},
		{
			Player:  NewPlayer(1, NewRace().WithPRT(AR)),
			Designs: Designs(DesignGalleon),
			Fleets:  []ScenarioFleet{{Design: "Galleon", Cargo: Cargo{Colonists: 50}, Waypoints: []ScenarioWaypoint{{Warp: 5}}}},
		},
	}, Planets: []ScenarioPlanet{
		{Name: "Planet 1", Owner: 1, Cargo: Cargo{Colonists: 2500}},
		{Name: "Planet 2", Owner: 2, Position: Vector{100, 0}, Cargo: Cargo{Colonists: 2500}},
	}})
	isPlayer, arPlayer := u.Player(1), u.Player(2)
	isFleet, arFleet, isPlanet := u.Fleet("Small Freighter #1"), u.Fleet("Galleon #1"), u.Planet("Planet 1")

	// AR colonists only die off when the fleet made a warp move
	arFleet.warped = true

	// don't generate a full turn, the planet will grow
	u.Run((*turnGenerator).fleetReproduce)

	// IS freighter should have grown; AR freighter should have lost pop slightly
	assert.Equal(t, 53, isFleet.Cargo.Colonists)
	assert.Equal(t, 2500, isPlanet.Cargo.Colonists)
	assert.Equal(t, 48, arFleet.Cargo.Colonists) // 3% of 50 is 1.5, rounded to 2 (49 in base game)

	// fill IS freighter up fully to overflow onto planet;
	// set AR freighter to 10K (we lose 3% or 300)
	isFleet.Cargo.Colonists = isFleet.Spec.CargoCapacity
	arFleet.Cargo.Colonists = 100

	// reproduce again
	u.Run((*turnGenerator).fleetReproduce)

	// IS should have grown on freighter and beamed down to planet
	assert.Equal(t, isFleet.Spec.CargoCapacity, isFleet.Cargo.Colonists)
	assert.Equal(t, 2509, isPlanet.Cargo.Colonists) // 12000 * 0.15 * 0.5 = 900 colonists beamed to planet
	assert.Equal(t, 97, arFleet.Cargo.Colonists)

	// Disable pop growth on both players & check for reproduction again;
	// IS should halt reproduction while AR should continue losing pop
	isPlayer.Race.GrowthRate = 0
	arPlayer.Race.GrowthRate = 0
	u.Run((*turnGenerator).fleetReproduce)
	assert.Equal(t, isFleet.Spec.CargoCapacity, isFleet.Cargo.Colonists)
	assert.Equal(t, 2509, isPlanet.Cargo.Colonists)
	assert.Equal(t, 94, arFleet.Cargo.Colonists)

	// AR colonists don't die off if the fleet didn't move
	arFleet.warped = false
	u.Run((*turnGenerator).fleetReproduce)
	assert.Equal(t, 94, arFleet.Cargo.Colonists)

	// or with 1000 colonists or fewer aboard
	arFleet.warped = true
	arFleet.Cargo.Colonists = 10
	u.Run((*turnGenerator).fleetReproduce)
	assert.Equal(t, 10, arFleet.Cargo.Colonists)
}

func Test_turn_fleetRadiatingEngineDieoff(t *testing.T) {
	s := singleFleetScenario(DesignSmallFreighter)
	s.Players[0].Fleets[0].Cargo = Cargo{Colonists: 50}
	u := newTestUniverse(t, s)
	player, fleet := u.Player(1), u.Fleet("Small Freighter #1")
	design := player.Designs[0]

	// generate turn to simulate die off; should not lose pop
	u.GenerateTurn()

	assert.Equal(t, 50, fleet.Cargo.Colonists)

	// add a radiating hydro ramscoop
	design.Slots[0].HullComponent = RadiatingHydroRamScoop.Name
	u.Recompute()

	// the engines didn't run, no radiation
	u.GenerateTurn()
	assert.Equal(t, 50, fleet.Cargo.Colonists)

	// move the fleet; should lose int((86 - 50) / 2) = 18% of the pop
	fleet.Waypoints = append(fleet.Waypoints, NewPositionWaypoint(Vector{1000, 0}, 5))
	u.GenerateTurn()
	assert.Equal(t, 41, fleet.Cargo.Colonists)

	// make the player a high rad race to prevent radiation damage
	player.Race.HabHigh.Rad = 100
	player.Race.HabLow.Rad = 80 // midpoint: 90mR
	u.Recompute()
	u.GenerateTurn()
	assert.Equal(t, 41, fleet.Cargo.Colonists)

}

func Test_turn_detonateMines(t *testing.T) {
	tests := []struct {
		name     string
		detonate bool
		friendly bool
		design   ShipDesign
		want     ShipToken
	}{
		{"no op", false, false, DesignLongRangeScout, ShipToken{Quantity: 1}},
		{"detonate, destroy ship", true, false, DesignLongRangeScout, ShipToken{Damage: 500}},
		{"detonate, don't destroy mini mine layers", true, true, DesignMiniMineLayer, ShipToken{Quantity: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := TestScenario{
				Players: []ScenarioPlayer{
					{
						Player:     NewPlayer(1, NewRace().WithPRT(SD)),
						Minefields: []Minefield{{MinefieldType: MinefieldTypeStandard, NumMines: 100, MinefieldOrders: MinefieldOrders{Detonate: tt.detonate}}},
					},
				},
			}
			fleetPlayer := 0
			if !tt.friendly {
				s.Players = append(s.Players, ScenarioPlayer{})
				fleetPlayer = 1
			}
			s.Players[fleetPlayer].Designs = Designs(tt.design)
			s.Players[fleetPlayer].Fleets = []ScenarioFleet{{Design: tt.design.Name, Waypoints: []ScenarioWaypoint{{Warp: 5}}}}
			u := newTestUniverse(t, s)
			fleet := u.Game.Fleets[0]
			u.GenerateTurn()
			if tt.want.Quantity == 0 {
				assert.True(t, fleet.Delete)
				assert.Empty(t, fleet.Tokens)
				return
			}
			token := fleet.Tokens[0]
			assert.Equal(t, tt.want.Quantity, token.Quantity)
			assert.Equal(t, tt.want.Damage, token.Damage)
			assert.Equal(t, tt.want.QuantityDamaged, token.QuantityDamaged)
		})
	}
}

func Test_turn_testPacketMoveHitPlanet(t *testing.T) {
	s := SingleUnitScenario()
	s.Players = append(s.Players, ScenarioPlayer{
		MineralPackets: []ScenarioMineralPacket{{To: "Planet 1", Position: Vector{25, 0}, WarpSpeed: 7, SafeWarpSpeed: 7, Cargo: Cargo{Ironium: 10}}},
	})

	u := newTestUniverse(t, s)
	player, planet := u.Player(1), u.Planet("Planet 1")

	// move packet, wipe out planet
	u.GenerateTurn()

	// packet hits, but planet is fine and we recover 1/3rd of the cargo
	assert.NotEqual(t, 0, planet.exactPopulation())
	assert.Equal(t, player.Num, planet.PlayerNum)
	assert.Equal(t, 10/3, planet.Cargo.Ironium)
}

func Test_turn_testPacketMoveDeleteStarbase(t *testing.T) {
	s := SingleUnitScenario()
	s.Players = append(s.Players, ScenarioPlayer{
		MineralPackets: []ScenarioMineralPacket{{To: "Planet 1", Position: Vector{25, 0}, WarpSpeed: 13, SafeWarpSpeed: 13, Cargo: Cargo{Ironium: 1000}}},
	})
	s.Planets[0].Starbase = "Starbase"
	s.Players[0].Designs = append(s.Players[0].Designs, ShipDesign{Name: "Starbase", Hull: SpaceStation.Name})

	u := newTestUniverse(t, s)
	planet := u.Planet("Planet 1")
	starbase := planet.Starbase

	// move packet, wipe out planet
	u.GenerateTurn()

	// no pop, no starbase
	assert.Equal(t, 0, planet.exactPopulation())
	assert.Equal(t, None, planet.PlayerNum)
	assert.Equal(t, true, starbase.Delete)
	assert.Nil(t, nil, planet.Starbase)

}

func Test_turn_decayPackets(t *testing.T) {
	s := SingleUnitScenario()
	s.Players[0].MineralPackets = []ScenarioMineralPacket{
		{Name: "Safe", To: "Planet 1", Position: Vector{200, 0}, WarpSpeed: 5, SafeWarpSpeed: 5, Cargo: Cargo{100, 100, 100, 0}},
		{
			Name:          "Too fast",
			To:            "Planet 1",
			Position:      Vector{0, 200},
			WarpSpeed:     10,
			SafeWarpSpeed: 7,
			Cargo:         Cargo{500, 500, 500, 0},
		},
		{Name: "New", To: "Planet 1", Position: Vector{200, 0}, WarpSpeed: 8, SafeWarpSpeed: 5, Cargo: Cargo{100, 100, 100, 0}},
	}
	u := newTestUniverse(t, s)
	packetSafe, packetTooFast, packetNewlyBuilt := u.Game.MineralPackets[0], u.Game.MineralPackets[1], u.Game.MineralPackets[2]
	packetNewlyBuilt.builtThisTurn = true

	// move and decay
	u.Run(func(turn *turnGenerator) { turn.packetMove(false) })
	u.Run(func(turn *turnGenerator) { turn.packetMove(true) })
	u.Run(func(turn *turnGenerator) { turn.decayPackets(false) })
	u.Run(func(turn *turnGenerator) { turn.decayPackets(true) })

	// no decay, 50% decay, and half of 50% decay for a newly built overfast packet
	assert.Equal(t, packetSafe.Cargo, Cargo{100, 100, 100, 0})
	assert.Equal(t, packetTooFast.Cargo, Cargo{250, 250, 250, 0})
	assert.Equal(t, packetNewlyBuilt.Cargo, Cargo{75, 75, 75, 0})
}

func Test_turn_randomCometStrikeOwnedPlanet(t *testing.T) {
	s := SingleUnitScenario()

	u := newTestUniverse(t, s)
	game, planet := u.Game, u.Planet("Planet 1")
	game.Year += game.Rules.RandomCometMinYear + game.Rules.RandomCometMinYearPlayerWorld

	startingPop := planet.exactPopulation()
	startingIronium := planet.Cargo.Ironium
	// strike planet
	u.Game.Rules.random = newFloat64Random(0) // 100% chance to strike planet
	u.Run((*turnGenerator).randomCometStrike)

	// comet hits, pop killed, minerals added
	assert.NotEqual(t, startingPop, planet.exactPopulation())
	assert.Greater(t, planet.Cargo.Ironium, startingIronium)

}

func Test_turn_randomCometStrikeOwnedPlanetAR(t *testing.T) {
	s := SingleUnitScenario()
	s.Players[0].Player = NewPlayer(1, NewRace().WithPRT(AR))

	u := newTestUniverse(t, s)
	game, planet := u.Game, u.Planet("Planet 1")
	game.Year += game.Rules.RandomCometMinYear + game.Rules.RandomCometMinYearPlayerWorld

	startingPop := planet.exactPopulation()
	startingIronium := planet.Cargo.Ironium
	// strike planet
	u.Game.Rules.random = newFloat64Random(0) // 100% chance to strike planet
	u.Run((*turnGenerator).randomCometStrike)

	// comet hits, pop not killed, minerals added
	assert.Equal(t, startingPop, planet.exactPopulation())
	assert.Greater(t, planet.Cargo.Ironium, startingIronium)
}

func Test_turn_fleetPatrol(t *testing.T) {
	s := singleFleetScenario(DesignStalwartDefender)
	s.Players[0].Fleets[0].Waypoints = []ScenarioWaypoint{{Warp: 5, Task: WaypointTaskPatrol, PatrolRange: 50}}
	s.Players = append(s.Players, ScenarioPlayer{
		Designs: Designs(DesignLongRangeScout),
		Fleets: []ScenarioFleet{
			{
				Name:      "Distant scout",
				Design:    "Long Range Scout",
				Position:  Vector{60, 0},
				Waypoints: []ScenarioWaypoint{{Position: Vector{60, 0}, Warp: 5}},
			},
			{
				Name:      "Close scout",
				Design:    "Long Range Scout",
				Position:  Vector{100, 0},
				Waypoints: []ScenarioWaypoint{{Position: Vector{100, 0}, Warp: 5}},
			},
		},
	})
	u := newTestUniverse(t, s)
	fleet, enemyFleet1, enemyFleet2 := u.Fleet("Stalwart Defender #1"), u.Fleet("Distant scout"), u.Fleet("Close scout")

	// no patrol target
	u.GenerateTurn()

	// should not attack
	assert.Equal(t, 1, len(fleet.Waypoints))

	// make a second enemy fleet closer
	enemyFleet2.Position = Vector{30, 0}

	// move first fleet within range as well
	enemyFleet1.Position = Vector{50, 0}

	// generate
	u.GenerateTurn()

	// should attack fleet 2
	assert.Equal(t, len(fleet.Waypoints), 2)
	assert.Equal(t, MapObjectTypeFleet, fleet.Waypoints[1].TargetType)
	assert.Equal(t, enemyFleet2.PlayerNum, fleet.Waypoints[1].TargetPlayerNum)
	assert.Equal(t, enemyFleet2.Num, fleet.Waypoints[1].TargetNum)
	assert.Equal(t, 6, fleet.Waypoints[1].WarpSpeed)

}

func Test_turn_fleetRemoteTerraform(t *testing.T) {
	u := newTestUniverse(t, TestScenario{Players: []ScenarioPlayer{
		{
			Player:    NewPlayer(1, NewRace().WithLRT(TT)),
			Relations: []PlayerRelationship{{Relation: PlayerRelationFriend}, {Relation: PlayerRelationNeutral}, {Relation: PlayerRelationFriend}},
			Designs:   Designs(DesignRemoteTerraformer),
			Fleets: []ScenarioFleet{
				{Name: "Enemy terraformer", Design: "Remote Terraformer", At: "Planet 1", Waypoints: []ScenarioWaypoint{{Warp: 5}}},
				{Name: "Friendly terraformer", Design: "Remote Terraformer", At: "Planet 2", Waypoints: []ScenarioWaypoint{{Warp: 5}}},
				{Name: "Gifted terraformer", Design: "Remote Terraformer", At: "Planet 3", Waypoints: []ScenarioWaypoint{{To: "Planet 3", Warp: 5, Task: WaypointTaskTransferFleet, TransferToPlayer: 3}}},
			},
		},
		{
			Relations: []PlayerRelationship{{Relation: PlayerRelationNeutral}, {Relation: PlayerRelationFriend}, {Relation: PlayerRelationNeutral}},
		},
		{
			Player:    NewPlayer(1, NewRace()).WithTechLevels(TechLevel{Propulsion: 1, Biotechnology: 1}),
			Relations: []PlayerRelationship{{Relation: PlayerRelationFriend}, {Relation: PlayerRelationNeutral}, {Relation: PlayerRelationFriend}},
		},
	}, Planets: []ScenarioPlanet{
		{Name: "Planet 1", Owner: 2, Cargo: Cargo{Colonists: 2500}},
		{Name: "Planet 2", Owner: 3, Hab: ScenarioValue(Hab{48, 50, 50}), Cargo: Cargo{Colonists: 2500}},
		{Name: "Planet 3", Owner: 3, Hab: ScenarioValue(Hab{48, 50, 50}), Cargo: Cargo{Colonists: 2500}},
	}})
	planet1, planet2, planet3 := u.Planet("Planet 1"), u.Planet("Planet 2"), u.Planet("Planet 3")
	// A gift retains the donor's terraformer specs even with the recipient's lower tech.
	u.Run((*turnGenerator).fleetTransferOwner)

	u.GenerateTurn()

	// should deterraform planet1 2 points
	assert.Equal(t, Hab{52, 50, 50}, planet1.Hab)

	// should terraform planet2 2 points
	assert.Equal(t, Hab{50, 50, 50}, planet2.Hab)

	// should terraform planet3 2 points
	assert.Equal(t, Hab{50, 50, 50}, planet3.Hab)

	// Both foreign owners receive a single report per fleet; our own world
	// receives one report, rather than duplicate sender/recipient reports.
	assert.Len(t, u.Messages(1, PlayerMessagePlanetRemoteTerraform), 2)
	assert.Len(t, u.Messages(2, PlayerMessagePlanetRemoteTerraform), 1)
	assert.Len(t, u.Messages(3, PlayerMessagePlanetRemoteTerraform), 2)
	enemyReport := u.Messages(2, PlayerMessagePlanetRemoteTerraform)[0]
	assert.Equal(t, planet1.Num, enemyReport.TargetNum)
	assert.Equal(t, 1, enemyReport.Spec.SourcePlayerNum)
	assert.Equal(t, 100, enemyReport.Spec.PrevAmount)
	assert.Equal(t, 99, enemyReport.Spec.Amount)
	assert.Equal(t, -1, enemyReport.Spec.Amount2)
	assert.Equal(t, Hab{Grav: 2}, enemyReport.Spec.TerraformAmount)
	assert.NotEqual(t, "Enemy terraformer", enemyReport.Spec.TargetName)

	// Optimal planets and neutral relationships produce no notification.
	for _, player := range u.Game.Players {
		player.Messages = nil
	}
	u.Player(1).Relations[1].Relation = PlayerRelationNeutral
	u.Fleet("Enemy terraformer").battlePlan.AttackWho = BattleAttackWhoEnemies
	u.Run((*turnGenerator).fleetRemoteTerraform)
	for _, player := range u.Game.Players {
		assert.Empty(t, u.Messages(player.Num, PlayerMessagePlanetRemoteTerraform))
	}
}

func Test_turn_fleetRefuel(t *testing.T) {
	s := SingleUnitScenario()
	s.Players[0].Designs = Designs(DesignLongRangeScout, ShipDesign{Name: "Starbase", Hull: SpaceStation.Name})
	s.Planets[0].Starbase = "Starbase"
	s.Planets[0].StarbaseDamage = 100
	s.Players[0].Fleets[0].EmptyFuel = true
	s.Players[0].Fleets[0].Fuel = 0

	u := newTestUniverse(t, s)
	fleet := u.Game.Fleets[0]

	// Refuel the empty tank at the starbase.
	u.GenerateTurn()

	// should refuel at starbase
	assert.Equal(t, fleet.Spec.FuelCapacity, fleet.Fuel)

}

func Test_turn_playerResearch(t *testing.T) {
	s := SingleUnitScenario()
	s.Players[0].Player = NewPlayer(1, NewRace())
	s.Players[0].Player.Researching = Energy
	s.Players[0].Player.ResearchAmount = 100
	s.Players[0].Player.NextResearchField = NextResearchFieldEnergy
	s.Planets[0].Cargo = Cargo{Colonists: 5000}
	s.Planets[0].Factories = 500
	u := newTestUniverse(t, s)
	player := u.Player(1)

	// let's go!!
	u.GenerateTurn()

	// costs go up per level so
	// 1. 50,
	// 2. 80 + 10,
	// 3. 130 + 20,
	// 4. 210 + 30,
	// 5. 340 + 40 (910 resources for level 5)
	// 6. 550 (we'll bleed over 90, leaving 360 for this next level)

	assert.Equal(t, TechLevel{Energy: 5}, player.TechLevels, "should have raised 5 energy levels")
	assert.Equal(t, TechLevel{Energy: 90}, player.TechLevelsSpent, "should leave 90 spent on energy level 6")
	assert.Equal(t, 1000, player.ResearchSpentLastYear, " spent 1000 on research last year")
	assert.True(t, len(player.TechsJustGained) > 0)
}

func Test_turn_buildStarbase(t *testing.T) {
	s := SingleUnitScenario()
	s.Players[0].Designs = Designs(DesignLongRangeScout, ShipDesign{Name: "Sad empty base", Hull: SpaceStation.Name}, ShipDesign{
		Name: "LASER BASE!!!!!!",
		Hull: SpaceStation.Name,
		Slots: []ShipDesignSlot{
			{HullComponent: Laser.Name, HullSlotIndex: 2, Quantity: 1},
			{HullComponent: Superlatanium.Name, HullSlotIndex: 4, Quantity: 1},
		},
	})
	s.Planets[0].Cargo = Cargo{1000, 1000, 1000, 10000}
	s.Planets[0].Factories = 1000
	s.Planets[0].ProductionQueue = []ScenarioProductionQueueItem{{Type: QueueItemTypeStarbase, Quantity: 1, Design: "Sad empty base"}}
	u := newTestUniverse(t, s)
	game, player, planet := u.Game, u.Player(1), u.Planet("Planet 1")
	game.Rules.RepairRates[RepairRateStarbase] = 0
	emptyBaseDesign, starbaseDesignUpgrade := player.Designs[1], player.Designs[2]

	// should have no starbase
	assert.Nil(t, planet.Starbase)
	assert.Equal(t, 0, len(game.Starbases))

	// generate a turn to build the starbase
	u.GenerateTurn()

	emptyBaseDesign.Spec.NumBuilt++ // increment numBuilt so json doesn't error

	// should have the old base on the planet
	assert.NotNil(t, planet.Starbase)
	assert.Equal(t, 1, len(game.Starbases))
	test.CompareAsJSON(t, planet.Starbase.Tokens[0].design, emptyBaseDesign)

	// give new base 90% damage
	planet.Starbase.Tokens[0].Damage = 450

	// upgrade the starbase with A LASER!
	// Also some superlat to check armor dmg stats
	planet.ProductionQueue = append(planet.ProductionQueue, ProductionQueueItem{
		Type:      QueueItemTypeStarbase,
		Quantity:  1,
		DesignNum: starbaseDesignUpgrade.Num,
	})

	// generate a turn to upgrade the starbase
	u.GenerateTurn()

	starbaseDesignUpgrade.Spec.NumBuilt++

	// should have an upgraded starbase at the planet,
	// with the other base marked for deletion
	// New base inherits 90% damage from the original due to no repairs
	assert.Equal(t, 2, len(game.Starbases))
	test.CompareAsJSON(t, planet.Starbase.Tokens[0].design, starbaseDesignUpgrade)
	assert.InDelta(t, 0.9, planet.Starbase.Tokens[0].Damage/
		float64(planet.Starbase.Tokens[0].design.Spec.Armor), 0.005)
	assert.True(t, game.Starbases[0].Delete)
	assert.False(t, game.Starbases[1].Delete)

	// give player RS and swap back to the old base
	player.Race = *player.Race.WithLRT(RS)
	planet.Starbase.Tokens[0].Damage = 1125 // 90% of our new 1250 max dp

	planet.ProductionQueue = append(planet.ProductionQueue, ProductionQueueItem{Type: QueueItemTypeStarbase, Quantity: 1, DesignNum: emptyBaseDesign.Num})

	u.GenerateTurn()

	// should have the old base again, with the laser base marked for deletion
	// still has roughly 90% damage
	// the first base was removed before this turn, like the server does on save
	assert.Equal(t, 2, len(game.Starbases))
	test.CompareAsJSON(t, planet.Starbase.Tokens[0].design, emptyBaseDesign)
	assert.InDelta(t, 0.9, planet.Starbase.Tokens[0].Damage/
		float64(planet.Starbase.Tokens[0].design.Spec.Armor), 0.005)
	assert.True(t, game.Starbases[0].Delete)
	assert.False(t, game.Starbases[1].Delete)

}

func Test_turn_fleetTransferOwner(t *testing.T) {
	s := TwoPlayerScenario()
	s.Players[0].Player = NewPlayer(1, NewRace().WithPluralName("Rabbitoids"))
	s.Players[1].Relations = []PlayerRelationship{{Relation: PlayerRelationFriend}, {Relation: PlayerRelationNeutral}}
	s.Players[0].Fleets[0].Waypoints = []ScenarioWaypoint{{To: "Planet 1", Task: WaypointTaskTransferFleet, TransferToPlayer: 2}}
	u := newTestUniverse(t, s)
	player1, player2, fleet := u.Player(1), u.Player(2), u.FleetFor(1, "Long Range Scout #1")

	// transfer
	u.GenerateTurn()

	// should have transferred the fleet, updated the name and the design
	assert.Equal(t, player2.Num, fleet.PlayerNum)
	assert.Equal(t, "Rabbitoids Long Range Scout #2", fleet.Name)
	assert.Equal(t, player1.Num, fleet.Tokens[0].design.OriginalPlayerNum)
	assert.Equal(t, "Rabbitoids Long Range Scout", fleet.Tokens[0].design.Name)
	assert.Equal(t, 2, len(player2.Designs))
	assert.Equal(t, 1, len(fleet.Waypoints))
	assert.Equal(t, None, fleet.Waypoints[0].TransferToPlayer)
	assert.Equal(t, WaypointTaskNone, fleet.Waypoints[0].Task)

}

func Test_turn_fleetBattle(t *testing.T) {
	s := TwoPlayerScenario()
	s.Players[0].Relations = []PlayerRelationship{{Relation: PlayerRelationFriend}, {Relation: PlayerRelationEnemy}}
	s.Players[1].Relations = []PlayerRelationship{{Relation: PlayerRelationEnemy}, {Relation: PlayerRelationFriend}}
	s.Players[1].Player.Race.PluralName = "Attackers"
	s.Players[1].Designs = Designs(DesignStalwartDefender)
	s.Players[0].Fleets = []ScenarioFleet{
		{
			Name:      "Scout",
			Design:    "Long Range Scout",
			Position:  Vector{100, 100},
			Waypoints: []ScenarioWaypoint{{Position: Vector{100, 100}, Warp: 5}},
		},
	}
	s.Players[1].Fleets = []ScenarioFleet{
		{
			Name:      "Destroyer",
			Design:    "Stalwart Defender",
			Position:  Vector{100, 100},
			Waypoints: []ScenarioWaypoint{{Position: Vector{100, 100}, Warp: 5}},
		},
	}
	u := newTestUniverse(t, s)
	game, player1, player2 := u.Game, u.Player(1), u.Player(2)
	fleet1 := u.Fleet("Scout")
	design1, design2 := player1.Designs[0], player2.Designs[0]

	// make sure our battle is always uses the same random seed
	u.Run((*turnGenerator).fleetBattle)

	// should have a battle record
	assert.Equal(t, 1, len(player1.BattleRecords))
	assert.Equal(t, 1, len(player2.BattleRecords))

	// ensure players were discovered
	assert.Equal(t, player1.Race.PluralName, player2.Intels.PlayerIntels[0].RacePluralName)
	assert.Equal(t, player2.Race.PluralName, player1.Intels.PlayerIntels[1].RacePluralName)

	// ensure designs were discovered
	assert.Equal(t, 1, len(player1.ShipDesignIntels))
	assert.Equal(t, design2.Slots, player1.ShipDesignIntels[0].Slots)
	assert.Equal(t, 1, len(player2.ShipDesignIntels))
	assert.Equal(t, design1.Slots, player2.ShipDesignIntels[0].Slots)

	// scout was destroyed, salvage created
	assert.Equal(t, true, fleet1.Delete)
	assert.Equal(t, 1, len(game.Salvages))
}

func Test_turn_fleetBattle3Players(t *testing.T) {
	u := newTestUniverse(t, TestScenario{Players: []ScenarioPlayer{
		{
			Player:    NewPlayer(1, NewRace().WithPluralName("Player1s")),
			Relations: []PlayerRelationship{{Relation: PlayerRelationFriend}, {Relation: PlayerRelationEnemy}, {Relation: PlayerRelationEnemy}},
			Designs:   Designs(DesignLongRangeScout),
			Fleets:    []ScenarioFleet{{Name: "Scout 1", Design: "Long Range Scout", Waypoints: []ScenarioWaypoint{{Warp: 5}}}},
		},
		{
			Player:    NewPlayer(1, NewRace().WithPluralName("Player2s")),
			Relations: []PlayerRelationship{{Relation: PlayerRelationEnemy}, {Relation: PlayerRelationFriend}, {Relation: PlayerRelationEnemy}},
			Designs:   Designs(DesignLongRangeScout),
			Fleets:    []ScenarioFleet{{Name: "Scout 2", Design: "Long Range Scout", Waypoints: []ScenarioWaypoint{{Warp: 5}}}},
		},
		{
			Player:    NewPlayer(1, NewRace().WithPluralName("Player3s")),
			Relations: []PlayerRelationship{{Relation: PlayerRelationEnemy}, {Relation: PlayerRelationEnemy}, {Relation: PlayerRelationFriend}},
			Designs:   Designs(DesignStalwartDefender),
			Fleets:    []ScenarioFleet{{Name: "Warships", Design: "Stalwart Defender", Quantity: 5, Waypoints: []ScenarioWaypoint{{Warp: 5}}}},
		},
	}})
	player1, player2, player3 := u.Player(1), u.Player(2), u.Player(3)
	fleet2 := u.Fleet("Scout 2")
	fg := u.Game

	// make sure our battle is always uses the same random seed
	u.Run((*turnGenerator).fleetBattle)

	// should have a battle record
	assert.Equal(t, 1, len(player1.BattleRecords))
	assert.Equal(t, 1, len(player2.BattleRecords))
	assert.Equal(t, 1, len(player3.BattleRecords))

	// ensure players were discovered
	assert.Equal(t, player1.Race.PluralName, player2.Intels.PlayerIntels[0].RacePluralName)
	assert.Equal(t, player1.Race.PluralName, player3.Intels.PlayerIntels[0].RacePluralName)
	assert.Equal(t, player2.Race.PluralName, player1.Intels.PlayerIntels[1].RacePluralName)
	assert.Equal(t, player2.Race.PluralName, player3.Intels.PlayerIntels[1].RacePluralName)
	assert.Equal(t, player3.Race.PluralName, player1.Intels.PlayerIntels[2].RacePluralName)
	assert.Equal(t, player3.Race.PluralName, player2.Intels.PlayerIntels[2].RacePluralName)

	// ensure designs were discovered
	assert.Equal(t, 2, len(player1.ShipDesignIntels))
	assert.True(t, len(player1.ShipDesignIntels[0].Slots) > 0)
	assert.True(t, len(player1.ShipDesignIntels[1].Slots) > 0)
	assert.Equal(t, 2, len(player2.ShipDesignIntels))
	assert.True(t, len(player2.ShipDesignIntels[0].Slots) > 0)
	assert.True(t, len(player2.ShipDesignIntels[1].Slots) > 0)
	assert.Equal(t, 2, len(player3.ShipDesignIntels))
	assert.True(t, len(player3.ShipDesignIntels[0].Slots) > 0)
	assert.True(t, len(player3.ShipDesignIntels[1].Slots) > 0)

	// scout was destroyed, salvage created
	assert.Equal(t, true, fleet2.Delete)
	assert.Equal(t, 1, len(fg.Salvages))
}

// test a fleet with repeat patrol orders
// it should intercept and kill one fleet, then
// return to base and target another
func Test_turn_fleetPatrolBattleRepeat(t *testing.T) {
	s := TwoPlayerScenario()
	s.Players[0].Relations = []PlayerRelationship{{Relation: PlayerRelationFriend}, {Relation: PlayerRelationEnemy}}
	s.Players[1].Relations = []PlayerRelationship{{Relation: PlayerRelationEnemy}, {Relation: PlayerRelationFriend}}
	s.Players[0].Designs = Designs(DesignJihadCruiser, ShipDesign{Name: "Starbase", Hull: SpaceStation.Name})
	s.Planets[0].Starbase = "Starbase"
	s.Players[0].Fleets = []ScenarioFleet{
		{
			Design:       "Jihad Cruiser",
			At:           "Planet 1",
			RepeatOrders: true,
			Waypoints:    []ScenarioWaypoint{{To: "Planet 1", Warp: 5, Task: WaypointTaskPatrol}},
		},
	}
	s.Players[1].Designs = Designs(DesignLongRangeScout)
	s.Players[1].Fleets = []ScenarioFleet{
		{
			Name:      "Scout 2",
			Design:    "Long Range Scout",
			Position:  Vector{25, 0},
			Waypoints: []ScenarioWaypoint{{Position: Vector{25, 0}, Warp: 5}},
		},
		{
			Name:      "Scout 3",
			Design:    "Long Range Scout",
			Position:  Vector{-30, 0},
			Waypoints: []ScenarioWaypoint{{Position: Vector{-30, 0}, Warp: 5}},
		},
	}
	u := newTestUniverse(t, s)
	game, player1, player2 := u.Game, u.Player(1), u.Player(2)
	fleet1, fleet2, fleet3 := u.Fleet("Jihad Cruiser #1"), u.Fleet("Scout 2"), u.Fleet("Scout 3")
	planet := u.Planet("Planet 1")

	// generate a turn to setup the patrol target
	u.GenerateTurn()

	// fleet1 should target fleet2
	assert.Equal(t, 2, len(fleet1.Waypoints))
	assert.Equal(t, fleet2.PlayerNum, fleet1.Waypoints[1].TargetPlayerNum)
	assert.Equal(t, fleet2.Num, fleet1.Waypoints[1].TargetNum)

	// generate a turn to allow the player1 fleet to attack player2's fleet
	u.GenerateTurn()

	// fleet moves to intercept
	assert.Equal(t, fleet1.Position, fleet2.Position)

	// should have a battle record
	assert.Equal(t, 1, len(player1.BattleRecords))
	assert.Equal(t, 1, len(player2.BattleRecords))

	// scout was destroyed, salvage created
	assert.Equal(t, true, fleet2.Delete)
	assert.Equal(t, 1, len(game.Salvages))

	// fleet1 should be returning to base after killing fleet2
	assert.Equal(t, 2, len(fleet1.Waypoints))
	assert.Equal(t, planet.Num, fleet1.Waypoints[1].TargetNum)
	assert.Equal(t, MapObjectTypePlanet, fleet1.Waypoints[1].TargetType)

	// generate a turn to allow the fleet1 to return home
	// and target fleet 3
	u.GenerateTurn()

	// fleet1 should target fleet3
	assert.Equal(t, planet.Position, fleet1.Position)
	assert.Equal(t, planet.Num, fleet1.OrbitingPlanetNum)
	assert.Equal(t, 2, len(fleet1.Waypoints))
	assert.Equal(t, fleet3.PlayerNum, fleet1.Waypoints[1].TargetPlayerNum)
	assert.Equal(t, fleet3.Num, fleet1.Waypoints[1].TargetNum)

}

// test a fleet with open ended patrol orders
// it should intercept and kill one fleet, then target
// another
func Test_turn_fleetPatrolKillPatrolAgain(t *testing.T) {
	s := TwoPlayerScenario()
	s.Players[0].Relations = []PlayerRelationship{{Relation: PlayerRelationFriend}, {Relation: PlayerRelationEnemy}}
	s.Players[1].Relations = []PlayerRelationship{{Relation: PlayerRelationEnemy}, {Relation: PlayerRelationFriend}}
	s.Players[0].Designs = Designs(DesignJihadCruiser, ShipDesign{Name: "Starbase", Hull: SpaceStation.Name})
	s.Planets[0].Starbase = "Starbase"
	s.Players[0].Fleets = []ScenarioFleet{
		{
			Design:       "Jihad Cruiser",
			At:           "Planet 1",
			RepeatOrders: false,
			Waypoints:    []ScenarioWaypoint{{To: "Planet 1", Warp: 5, Task: WaypointTaskPatrol}},
		},
	}
	s.Players[1].Designs = Designs(DesignLongRangeScout)
	s.Players[1].Fleets = []ScenarioFleet{
		{
			Name:      "Scout 2",
			Design:    "Long Range Scout",
			Position:  Vector{25, 0},
			Waypoints: []ScenarioWaypoint{{Position: Vector{25, 0}, Warp: 5}},
		},
		{
			Name:      "Scout 3",
			Design:    "Long Range Scout",
			Position:  Vector{-30, 0},
			Waypoints: []ScenarioWaypoint{{Position: Vector{-30, 0}, Warp: 5}},
		},
	}
	u := newTestUniverse(t, s)
	game, player1, player2 := u.Game, u.Player(1), u.Player(2)
	fleet1, fleet2, fleet3 := u.Fleet("Jihad Cruiser #1"), u.Fleet("Scout 2"), u.Fleet("Scout 3")

	// generate a turn to setup the patrol target
	u.GenerateTurn()

	// fleet1 should target fleet2
	assert.Equal(t, 2, len(fleet1.Waypoints))
	assert.Equal(t, fleet2.PlayerNum, fleet1.Waypoints[1].TargetPlayerNum)
	assert.Equal(t, fleet2.Num, fleet1.Waypoints[1].TargetNum)

	// generate a turn to allow the player1 fleet to attack player2's fleet
	u.GenerateTurn()

	// fleet moves to intercept
	assert.Equal(t, fleet1.Position, fleet2.Position)

	// should have a battle record
	assert.Equal(t, 1, len(player1.BattleRecords))
	assert.Equal(t, 1, len(player2.BattleRecords))

	// scout was destroyed, salvage created
	assert.Equal(t, true, fleet2.Delete)
	assert.Equal(t, 1, len(game.Salvages))

	// fleet1 should now target fleet3
	assert.Equal(t, 2, len(fleet1.Waypoints))
	assert.Equal(t, fleet3.PlayerNum, fleet1.Waypoints[1].TargetPlayerNum)
	assert.Equal(t, fleet3.Num, fleet1.Waypoints[1].TargetNum)
}

func Test_turn_mysteryTraderSpawn(t *testing.T) {
	u := newTestUniverse(t, SingleUnitScenario())
	game := u.Game
	game.RandomEvents = true
	game.Rules.random = newIntRandom()
	game.Year += game.Rules.MysteryTraderRules.MinYear

	// move to place
	u.Run((*turnGenerator).mysteryTraderSpawn)

	// should have consumed that waypoint and moved to the space
	assert.Equal(t, 1, len(game.MysteryTraders))
	mt := game.MysteryTraders[0]
	assert.Equal(t, 1, len(game.getMapObjectsAtPosition(mt.Position)))
}

func Test_turn_mysteryTraderMove(t *testing.T) {
	s := SingleUnitScenario()
	s.MysteryTraders = []MysteryTrader{
		{
			MapObject:     MapObject{Position: Vector{}},
			WarpSpeed:     7,
			Destination:   Vector{100, 0},
			RequestedBoon: 5000,
			RewardType:    MysteryTraderRewardResearch,
		},
	}
	u := newTestUniverse(t, s)
	game := u.Game
	game.RandomEvents = true
	game.Rules.random = newIntRandom(1)
	mt := game.MysteryTraders[0]

	// move to place
	u.Run((*turnGenerator).mysteryTraderMove)

	// should have moved to our location
	assert.Equal(t, Vector{49, 0}, mt.Position)
	assert.Equal(t, 1, len(game.getMapObjectsAtPosition(mt.Position)))
}

func Test_turn_mysteryTraderMoveChangeCourse(t *testing.T) {
	s := SingleUnitScenario()
	s.MysteryTraders = []MysteryTrader{
		{
			MapObject:     MapObject{Position: Vector{}},
			WarpSpeed:     7,
			Destination:   Vector{100, 0},
			RequestedBoon: 5000,
			RewardType:    MysteryTraderRewardResearch,
		},
	}
	u := newTestUniverse(t, s)
	game := u.Game
	game.RandomEvents = true
	game.Rules.random = newIntRandom(0)
	mt := game.MysteryTraders[0]

	// move to place
	u.Run((*turnGenerator).mysteryTraderMove)

	// should have changed course, so we won't move along the previous path
	assert.NotEqual(t, Vector{49, 0}, mt.Position)
	player := game.Players[0]
	assert.Equal(t, 1, len(player.Messages))
	assert.Equal(t, PlayerMessageMysteryTraderChangedCourse, player.Messages[0].Type)
}

func Test_turn_mysteryTraderFinished(t *testing.T) {
	s := SingleUnitScenario()
	s.MysteryTraders = []MysteryTrader{
		{
			MapObject:     MapObject{Position: Vector{}},
			WarpSpeed:     7,
			Destination:   Vector{49, 0},
			RequestedBoon: 5000,
			RewardType:    MysteryTraderRewardResearch,
		},
	}
	u := newTestUniverse(t, s)
	game := u.Game
	game.RandomEvents = true
	game.Rules.random = newIntRandom(1, 1)
	mt := game.MysteryTraders[0]

	// move to place
	u.Run((*turnGenerator).mysteryTraderMove)

	// should have changed course, so we won't move along the previous path
	assert.Equal(t, Vector{49, 0}, mt.Position)
	assert.True(t, mt.Delete)
}

func Test_turn_mysteryTraderAgain(t *testing.T) {
	s := SingleUnitScenario()
	s.MysteryTraders = []MysteryTrader{
		{
			MapObject:     MapObject{Position: Vector{}},
			WarpSpeed:     7,
			Destination:   Vector{49, 0},
			RequestedBoon: 5000,
			RewardType:    MysteryTraderRewardResearch,
		},
	}
	u := newTestUniverse(t, s)
	game := u.Game
	game.RandomEvents = true
	game.Rules.random = newIntRandom(1, 0)
	mt := game.MysteryTraders[0]

	// move to place
	u.Run((*turnGenerator).mysteryTraderMove)

	// should have changed course, so we won't move along the previous path
	assert.Equal(t, Vector{49, 0}, mt.Position)
	player := game.Players[0]
	assert.Equal(t, 1, len(player.Messages))
	assert.Equal(t, PlayerMessageMysteryTraderAgain, player.Messages[0].Type)
	assert.NotEqual(t, Vector{49, 0}, mt.Destination)
}

func Test_turn_mysteryTraderMeetNoReward(t *testing.T) {
	s := SingleUnitScenario()
	s.MysteryTraders = []MysteryTrader{
		{
			MapObject:     MapObject{Position: Vector{}},
			WarpSpeed:     7,
			Destination:   Vector{100, 0},
			RequestedBoon: 5000,
			RewardType:    MysteryTraderRewardResearch,
		},
	}
	s.Players[0].Fleets[0].At = ""
	s.Players[0].Fleets[0].Waypoints = []ScenarioWaypoint{{To: "Mystery Trader #1", Warp: 5}}
	s.Players[0].Player = NewPlayer(1, NewRace())
	u := newTestUniverse(t, s)
	game := u.Game
	game.RandomEvents = true
	game.Rules.random = &testRandom{}
	fleet, player := u.Fleet("Long Range Scout #1"), u.Player(1)

	// meet mystery trader
	u.RunE((*turnGenerator).mysteryTraderMeet)

	// fleet should be left alone
	assert.Equal(t, false, fleet.Delete)
	assert.Equal(t, PlayerMessageMysteryTraderMetWithoutReward, player.Messages[0].Type)
	assert.Equal(t, TechLevel{}, player.TechLevels)
}

func Test_turn_mysteryTraderMeetReward(t *testing.T) {
	s := SingleUnitScenario()
	s.MysteryTraders = []MysteryTrader{
		{
			MapObject:     MapObject{Position: Vector{}},
			WarpSpeed:     7,
			Destination:   Vector{100, 0},
			RequestedBoon: 5000,
			RewardType:    MysteryTraderRewardResearch,
		},
	}
	s.Players[0].Fleets[0].At = ""
	s.Players[0].Fleets[0].Waypoints = []ScenarioWaypoint{{To: "Mystery Trader #1", Warp: 5}}
	s.Players[0].Fleets[0].Cargo = Cargo{Ironium: 5000}
	s.Players[0].Player = NewPlayer(1, NewRace())
	u := newTestUniverse(t, s)
	game := u.Game
	game.RandomEvents = true
	game.Rules.random = &testRandom{}
	fleet, player := u.Fleet("Long Range Scout #1"), u.Player(1)

	// meet mystery trader
	u.RunE((*turnGenerator).mysteryTraderMeet)

	// fleet should be deleted, player gained tech
	assert.Equal(t, true, fleet.Delete)
	assert.Equal(t, PlayerMessageMysteryTraderMetWithReward, player.Messages[0].Type)
	assert.Equal(t, TechLevel{Energy: 6}, player.TechLevels)
}

func Test_turn_mysteryTraderMeetRewardTech(t *testing.T) {
	s := SingleUnitScenario()
	s.MysteryTraders = []MysteryTrader{
		{
			MapObject:     MapObject{Position: Vector{}},
			WarpSpeed:     7,
			Destination:   Vector{100, 0},
			RequestedBoon: 5000,
			RewardType:    MysteryTraderRewardTorpedo,
		},
	}
	s.Players[0].Fleets[0].At = ""
	s.Players[0].Fleets[0].Waypoints = []ScenarioWaypoint{{To: "Mystery Trader #1", Warp: 5}}
	s.Players[0].Fleets[0].Cargo = Cargo{Ironium: 5000}
	u := newTestUniverse(t, s)
	game := u.Game
	game.RandomEvents = true
	game.Rules.random = &testRandom{}
	fleet, player := u.Fleet("Long Range Scout #1"), u.Player(1)

	// meet mystery trader
	u.RunE((*turnGenerator).mysteryTraderMeet)

	// fleet should be deleted, player gained tech
	assert.Equal(t, true, fleet.Delete)
	assert.Equal(t, PlayerMessageMysteryTraderMetWithReward, player.Messages[0].Type)
	assert.True(t, player.HasAcquiredTech(&AntiMatterTorpedo.Tech))
}

func Test_turn_mysteryTraderMeetRewardTechAlreadyAcquired(t *testing.T) {
	s := SingleUnitScenario()
	s.MysteryTraders = []MysteryTrader{
		{
			MapObject:     MapObject{Position: Vector{}},
			WarpSpeed:     7,
			Destination:   Vector{100, 0},
			RequestedBoon: 5000,
			RewardType:    MysteryTraderRewardTorpedo,
		},
	}
	s.Players[0].Fleets[0].At = ""
	s.Players[0].Fleets[0].Waypoints = []ScenarioWaypoint{{To: "Mystery Trader #1", Warp: 5}}
	s.Players[0].Fleets[0].Cargo = Cargo{Ironium: 5000}
	s.Players[0].Player = NewPlayer(1, NewRace())
	s.Players[0].Player.AcquiredTechs = map[string]bool{AntiMatterTorpedo.Name: true}
	u := newTestUniverse(t, s)
	game := u.Game
	game.RandomEvents = true
	game.Rules.random = &testRandom{}
	fleet, player := u.Fleet("Long Range Scout #1"), u.Player(1)

	// meet mystery trader
	u.RunE((*turnGenerator).mysteryTraderMeet)

	// fleet should be deleted, player gained an additional tech
	assert.Equal(t, true, fleet.Delete)
	assert.Equal(t, PlayerMessageMysteryTraderMetWithReward, player.Messages[0].Type)
	assert.Equal(t, 2, len(player.AcquiredTechs), "player should have acquired an additional tech. message is: %v", player.Messages[0])
}

func Test_turn_mysteryTraderMeetRewardShip(t *testing.T) {
	s := SingleUnitScenario()
	s.MysteryTraders = []MysteryTrader{
		{
			MapObject:     MapObject{Position: Vector{}},
			WarpSpeed:     7,
			Destination:   Vector{100, 0},
			RequestedBoon: 5000,
			RewardType:    MysteryTraderRewardLifeboat,
		},
	}
	s.Players[0].Fleets[0].At = ""
	s.Players[0].Fleets[0].Waypoints = []ScenarioWaypoint{{To: "Mystery Trader #1", Warp: 5}}
	s.Players[0].Fleets[0].Cargo = Cargo{Ironium: 5000}
	s.Players[0].Player = NewPlayer(1, NewRace())
	u := newTestUniverse(t, s)
	game := u.Game
	game.RandomEvents = true
	game.Rules.random = newIntRandom()
	fleet, player := u.Fleet("Long Range Scout #1"), u.Player(1)

	// meet mystery trader
	u.RunE((*turnGenerator).mysteryTraderMeet)

	// fleet should be deleted, player gained tech
	assert.Equal(t, true, fleet.Delete)
	assert.Equal(t, PlayerMessageMysteryTraderMetWithReward, player.Messages[0].Type)
	assert.True(t, player.Messages[0].Spec.MysteryTrader.Ship.Name != "")
	assert.True(t, player.Messages[0].Spec.MysteryTrader.ShipCount != 0)

	// old fleet gone, new fleet should be a mystery trader design
	assert.Equal(t, 2, len(game.Fleets))
	assert.Equal(t, 2, len(player.Designs))
	rewardFleet := game.Fleets[1]
	design := player.Designs[1]
	assert.True(t, rewardFleet.Tokens[0].design.MysteryTrader)
	assert.True(t, rewardFleet.Tokens[0].Quantity > 0)
	assert.True(t, design.MysteryTrader)

}

func Test_turn_mysteryTraderMeetAlreadyRewarded(t *testing.T) {
	s := SingleUnitScenario()
	s.MysteryTraders = []MysteryTrader{
		{
			MapObject:     MapObject{Position: Vector{}},
			WarpSpeed:     7,
			Destination:   Vector{100, 0},
			RequestedBoon: 5000,
			RewardType:    MysteryTraderRewardTorpedo,
		},
	}
	s.Players[0].Fleets[0].At = ""
	s.Players[0].Fleets[0].Waypoints = []ScenarioWaypoint{{To: "Mystery Trader #1", Warp: 5}}
	s.Players[0].Fleets[0].Cargo = Cargo{Ironium: 5000}
	s.MysteryTraders[0].PlayersRewarded = map[int]bool{1: true}
	u := newTestUniverse(t, s)
	game := u.Game
	game.RandomEvents = true
	game.Rules.random = &testRandom{}
	fleet, player := u.Fleet("Long Range Scout #1"), u.Player(1)

	// meet mystery trader
	u.RunE((*turnGenerator).mysteryTraderMeet)

	// fleet should be deleted, player gained tech
	assert.Equal(t, PlayerMessageMysteryTraderAlreadyRewarded, player.Messages[0].Type)
	assert.False(t, fleet.Delete)
	assert.False(t, player.HasAcquiredTech(&AntiMatterTorpedo.Tech))
}

func Test_turn_buildMysteryTraderGenesisDevice(t *testing.T) {
	s := SingleUnitScenario()
	s.Planets[0].Cargo = Cargo{1000, 1000, 1000, 10000}
	s.Planets[0].Defenses, s.Planets[0].Mines, s.Planets[0].Factories = 100, 1000, 1000
	s.Planets[0].ProductionQueue = []ScenarioProductionQueueItem{{Type: QueueItemTypeGenesisDevice, Quantity: 1}}
	u := newTestUniverse(t, s)
	u.Game.Rules.MysteryTraderRules.GenesisDeviceCost = Cost{0, 0, 0, 100}

	// generate a turn to build a starbase
	u.GenerateTurn()

	// should have an upgraded starbase at the planet, and the other should be
	// marked for deletion
	// TODO: not sure how to test this. Random number gen I guess...
}

func Test_turn_productionQueueMessages(t *testing.T) {
	for _, tt := range []struct {
		name      string
		queue     []ScenarioProductionQueueItem
		resources int
		want      PlayerMessageType
	}{
		{name: "empty", resources: 100, want: PlayerMessagePlanetProductionQueueEmpty},
		{name: "completed", queue: []ScenarioProductionQueueItem{{Type: QueueItemTypeMine, Quantity: 1}}, resources: 100, want: PlayerMessagePlanetProductionQueueComplete},
		{name: "blocked", queue: []ScenarioProductionQueueItem{{Type: QueueItemTypeMine, Quantity: 1}}, resources: 0},
		{name: "automatic", queue: []ScenarioProductionQueueItem{{Type: QueueItemTypeAutoMines, Quantity: 1}}, resources: 100},
		{name: "canceled", queue: []ScenarioProductionQueueItem{{Type: QueueItemTypeShipToken, Design: "Long Range Scout", Quantity: 1}}, resources: 100, want: PlayerMessagePlanetProductionQueueEmpty},
		{name: "completed and canceled", queue: []ScenarioProductionQueueItem{{Type: QueueItemTypeMine, Quantity: 1}, {Type: QueueItemTypeShipToken, Design: "Long Range Scout", Quantity: 1}}, resources: 100, want: PlayerMessagePlanetProductionQueueEmpty},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s := SingleUnitScenario()
			s.Planets[0].ProductionQueue = tt.queue
			u := newTestUniverse(t, s)
			planet := u.Game.Planets[0]
			planet.Spec.ResourcesPerYearAvailable = tt.resources
			planet.Cargo.Ironium, planet.Cargo.Boranium, planet.Cargo.Germanium = 1000, 1000, 1000
			u.RunE((*turnGenerator).planetProduction)
			empty := u.Messages(1, PlayerMessagePlanetProductionQueueEmpty)
			complete := u.Messages(1, PlayerMessagePlanetProductionQueueComplete)
			if tt.want == PlayerMessageNone {
				assert.Empty(t, empty)
				assert.Empty(t, complete)
			} else {
				messages := append(empty, complete...)
				require.Len(t, messages, 1)
				assert.Equal(t, tt.want, messages[0].Type)
				assert.Equal(t, planet.Num, messages[0].TargetNum)
			}
		})
	}
}

func Test_turn_fleetUnloadColonistsOnPlanetThatDiedThisTurn(t *testing.T) {
	newUniverse := func(owner int) *testUniverse {
		s := TwoPlayerScenario()
		s.Players[1].Player = AIPlayer("Player 2").Player
		s.Players[0].Fleets[0].At = "Planet 2"
		s.Players[0].Fleets[0].Cargo = Cargo{Colonists: 10}
		s.Players[0].Fleets[0].Waypoints = []ScenarioWaypoint{{
			To:             "Planet 2",
			Task:           WaypointTaskTransport,
			TransportTasks: WaypointTransportTasks{Colonists: WaypointTransportTask{Action: TransportActionUnloadAll}},
		}}
		s.Planets[1].Owner = owner
		return newTestUniverse(t, s)
	}

	t.Run("take over a planet whose population died this turn", func(t *testing.T) {
		u := newUniverse(2)
		planet, fleet := u.Planet("Planet 2"), u.FleetFor(1, "Long Range Scout #1")
		u.turn.planetInit()
		// bombed out before we unload
		planet.emptyPlanet()

		u.turn.fleetUnload()

		assert.Equal(t, 1, planet.PlayerNum)
		assert.Equal(t, 1000, planet.GetPopulation())
		assert.Equal(t, 0, fleet.Cargo.Colonists)
		assert.Len(t, u.Messages(2, PlayerMessagePlanetInvaded), 1)
	})

	t.Run("retake our own planet that died this turn", func(t *testing.T) {
		u := newUniverse(1)
		planet := u.Planet("Planet 2")
		u.turn.planetInit()
		planet.emptyPlanet()

		u.turn.fleetUnload()

		assert.Equal(t, 1, planet.PlayerNum)
		assert.Equal(t, 1000, planet.GetPopulation())
	})

	t.Run("planets empty at the start of the turn can't be taken", func(t *testing.T) {
		u := newUniverse(2)
		planet, fleet := u.Planet("Planet 2"), u.FleetFor(1, "Long Range Scout #1")
		planet.emptyPlanet()
		u.turn.planetInit()

		u.turn.fleetUnload()

		assert.False(t, planet.Owned())
		assert.Equal(t, 10, fleet.Cargo.Colonists)
		assert.Len(t, u.Messages(1, PlayerMessageFleetTransportInvalid), 1)
	})
}
