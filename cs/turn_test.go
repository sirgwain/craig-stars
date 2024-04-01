package cs

import (
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

// many functions require a copy of the current game's rules.
// for testing, create a standard rules var every test can use
var rules = NewRules()

type MockRand struct {
	int63Result int64
}

func (m MockRand) Seed(seed int64) {}

func (m MockRand) Int63() int64 {
	return m.int63Result
}

func createSingleUnitGame() *FullGame {
	client := NewGamer()
	game := client.CreateGame(1, *NewGameSettings())
	game.Rules.ResetSeed(0) // keep the same seed for tests
	player := client.NewPlayer(1, *NewRace(), &game.Rules).withSpec(&game.Rules)
	player.Num = 1
	player.Relations = []PlayerRelationship{{Relation: PlayerRelationFriend}} // friends with themselves
	player.PlayerIntels.PlayerIntels = player.defaultPlayerIntels([]*Player{player})

	planet := &Planet{
		MapObject: MapObject{Type: MapObjectTypePlanet, Name: "Planet 1", Num: 1, PlayerNum: player.Num},
		Hab:       Hab{50, 50, 50},
		BaseHab:   Hab{50, 50, 50},
		Cargo: Cargo{
			Colonists: 2500,
		},
	}
	planet.Spec = computePlanetSpec(&game.Rules, player, planet)

	// setup initial planet intels for this planet
	player.initDefaultPlanetIntels(&game.Rules, []*Planet{planet})

	fleet := testLongRangeScout(player)
	fleet.OrbitingPlanetNum = planet.Num
	fleet.Waypoints = []Waypoint{
		NewPlanetWaypoint(Vector{}, 1, "Planet 1", 5),
	}
	player.Designs = []*ShipDesign{
		fleet.Tokens[0].design,
	}

	players := []*Player{player}

	universe := NewUniverse(&game.Rules)
	universe.Planets = append(universe.Planets, planet)
	universe.Fleets = append(universe.Fleets, fleet)

	universe.buildMaps(players)

	return &FullGame{
		Game:      game,
		Universe:  &universe,
		TechStore: &StaticTechStore,
		Players:   players,
	}

}

func createTwoPlayerGame() *FullGame {
	client := NewGamer()
	game := client.CreateGame(1, *NewGameSettings())
	game.Rules.ResetSeed(0) // keep the same seed for tests
	player1 := client.NewPlayer(1, *NewRace(), &game.Rules).withSpec(&game.Rules)
	player1.Num = 1
	player1.Relations = []PlayerRelationship{{Relation: PlayerRelationFriend}, {Relation: PlayerRelationNeutral}}

	player2 := client.NewPlayer(1, *NewRace(), &game.Rules).withSpec(&game.Rules)
	player2.Num = 2
	player2.Relations = []PlayerRelationship{{Relation: PlayerRelationFriend}, {Relation: PlayerRelationNeutral}}

	player1.PlayerIntels.PlayerIntels = player1.defaultPlayerIntels([]*Player{player1, player2})
	player2.PlayerIntels.PlayerIntels = player2.defaultPlayerIntels([]*Player{player2, player2})

	// create homeworlds
	planet1 := &Planet{
		MapObject: MapObject{Type: MapObjectTypePlanet, Name: "Planet 1", Num: 1, PlayerNum: player1.Num},
		Hab:       Hab{50, 50, 50},
		BaseHab:   Hab{50, 50, 50},
		Cargo: Cargo{
			Colonists: 2500,
		},
	}
	planet1.Spec = computePlanetSpec(&game.Rules, player1, planet1)

	planet2 := &Planet{
		MapObject: MapObject{Type: MapObjectTypePlanet, Name: "Planet 2", Num: 2, PlayerNum: player2.Num, Position: Vector{100, 0}},
		Hab:       Hab{50, 50, 50},
		BaseHab:   Hab{50, 50, 50},
		Cargo: Cargo{
			Colonists: 2500,
		},
	}
	planet2.Spec = computePlanetSpec(&game.Rules, player2, planet2)

	// setup initial planet intels for this planet
	player1.initDefaultPlanetIntels(&game.Rules, []*Planet{planet1, planet2})
	player2.initDefaultPlanetIntels(&game.Rules, []*Planet{planet1, planet2})

	// give each player a scout on their homeworld
	fleet1 := testLongRangeScout(player1)
	fleet1.OrbitingPlanetNum = planet1.Num
	fleet1.Waypoints = []Waypoint{
		NewPlanetWaypoint(Vector{}, 1, "Planet 1", 5),
	}
	player1.Designs = []*ShipDesign{
		fleet1.Tokens[0].design,
	}

	// give each player a scout on their homeworld
	fleet2 := testLongRangeScout(player2)
	fleet2.OrbitingPlanetNum = planet2.Num
	fleet2.Waypoints = []Waypoint{
		NewPlanetWaypoint(Vector{}, 2, "Planet 2", 5),
	}
	player2.Designs = []*ShipDesign{
		fleet2.Tokens[0].design,
	}

	players := []*Player{player1, player2}

	universe := NewUniverse(&game.Rules)
	universe.Planets = append(universe.Planets, planet1, planet2)
	universe.Fleets = append(universe.Fleets, fleet1, fleet2)

	universe.buildMaps(players)

	return &FullGame{
		Game:      game,
		Universe:  &universe,
		TechStore: &StaticTechStore,
		Players:   players,
	}

}

func Test_generateTurn(t *testing.T) {
	client := NewGamer()
	game := client.CreateGame(1, *NewGameSettings())
	player := client.NewPlayer(1, *NewRace(), &game.Rules)
	players := []*Player{player}
	player.AIControlled = true
	player.Num = 1
	universe, _ := client.GenerateUniverse(game, players)

	// build a ship on the planet
	pmo := universe.GetPlayerMapObjects(player.Num)
	planet := pmo.Planets[0]
	planet.ProductionQueue = append([]ProductionQueueItem{{Type: QueueItemTypeShipToken, Quantity: 1, DesignNum: player.Designs[0].Num}}, planet.ProductionQueue...)

	startingFleets := len(universe.Fleets)

	client.GenerateTurn(game, universe, players)

	assert.Equal(t, 2401, game.Year)

	// should have intel about planets
	assert.Equal(t, len(universe.Planets), len(player.PlanetIntels))

	// should have built a new scout
	assert.Greater(t, len(universe.Fleets), startingFleets)

	// should have grown pop
	assert.Greater(t, universe.Planets[0].population(), player.Race.Spec.StartingPlanets[0].Population)
}

func Test_generateTurns(t *testing.T) {
	client := NewGamer()
	game := client.CreateGame(1, *NewGameSettings())
	player := client.NewPlayer(1, *NewRace(), &game.Rules)
	player.AIControlled = true
	player.Num = 1
	players := []*Player{player}
	universe, _ := client.GenerateUniverse(game, players)

	// generate many turns
	numTurns := 10
	for i := 0; i < numTurns; i++ {
		client.GenerateTurn(game, universe, players)

		// manually delete these
		fleets := make([]*Fleet, 0, len(universe.Fleets))
		for _, fleet := range universe.Fleets {
			if !fleet.Delete {
				fleets = append(fleets, fleet)
			}
		}
		universe.Fleets = fleets

		salvages := make([]*Salvage, 0, len(universe.Salvages))
		for _, salvage := range universe.Salvages {
			if !salvage.Delete {
				salvages = append(salvages, salvage)
			}
		}
		universe.Salvages = salvages

		mineFields := make([]*MineField, 0, len(universe.MineFields))
		for _, mineField := range universe.MineFields {
			if !mineField.Delete {
				mineFields = append(mineFields, mineField)
			}
		}
		universe.MineFields = mineFields

		mineralPackets := make([]*MineralPacket, 0, len(universe.MineralPackets))
		for _, mineralPacket := range universe.MineralPackets {
			if !mineralPacket.Delete {
				mineralPackets = append(mineralPackets, mineralPacket)
			}
		}
		universe.MineralPackets = mineralPackets

		mysteryTraders := make([]*MysteryTrader, 0, len(universe.MysteryTraders))
		for _, mysteryTrader := range universe.MysteryTraders {
			if !mysteryTrader.Delete {
				mysteryTraders = append(mysteryTraders, mysteryTrader)
			}
		}
		universe.MysteryTraders = mysteryTraders

		wormholes := make([]*Wormhole, 0, len(universe.Wormholes))
		for _, wormhole := range universe.Wormholes {
			if !wormhole.Delete {
				wormholes = append(wormholes, wormhole)
			}
		}
		universe.Wormholes = wormholes

	}

	assert.Equal(t, game.Rules.StartingYear+numTurns, game.Year)

	// should have fleets
	assert.True(t, len(universe.Fleets) > 0)

	// should have grown pop
	assert.Greater(t, universe.Planets[0].population(), player.Race.Spec.StartingPlanets[0].Population)

	// should have built factories
	assert.Greater(t, universe.Planets[0].Factories, game.Rules.StartingFactories)

	// no victor
	assert.False(t, player.Victor)
	assert.False(t, game.VictorDeclared)

}

func Test_turn_grow(t *testing.T) {
	game := createSingleUnitGame()

	player := game.Players[0]

	// add a second, third, and fourth planet
	// first two are hostile and will lose colonists
	// the final planet is for overcrowding
	game.Planets = append(game.Planets,
		&Planet{
			MapObject: MapObject{Type: MapObjectTypePlanet, Name: "Planet 2", Num: 2, PlayerNum: player.Num},
			Hab:       Hab{0, 0, 0}, // bad planet
			BaseHab:   Hab{0, 0, 0},
		},
		&Planet{
			MapObject: MapObject{Type: MapObjectTypePlanet, Name: "Planet 3", Num: 3, PlayerNum: player.Num},
			Hab:       Hab{0, 0, 0}, // bad planet
			BaseHab:   Hab{0, 0, 0},
		},
		&Planet{
			MapObject: MapObject{Type: MapObjectTypePlanet, Name: "Planet 3", Num: 3, PlayerNum: player.Num},
			Hab:       Hab{50, 50, 50}, // good planet
			BaseHab:   Hab{50, 50, 50},
		},
	)

	planet1 := game.Planets[0]
	planet2 := game.Planets[1]
	planet3 := game.Planets[2]
	planet4 := game.Planets[3]
	planet1.setPopulation(100_000)
	planet2.setPopulation(100_000)
	planet3.setPopulation(100)       // last turn with us
	planet4.setPopulation(2_400_000) // should lose 4%

	turn := turn{
		game: game,
	}
	turn.game.Universe.buildMaps(game.Players)

	turn.generateTurn()

	// one planet should grow, another should not, the other should die off completely
	assert.Equal(t, 115_000, planet1.population())
	assert.Equal(t, 95_500, planet2.population())
	assert.Equal(t, 0, planet3.population())
	assert.Equal(t, false, planet3.Owned())
	assert.Equal(t, 2_304_000, planet4.population())
}

func Test_turn_fleetRoute(t *testing.T) {
	game := createSingleUnitGame()

	// add a second planet target
	game.Planets = append(game.Planets, &Planet{
		MapObject: MapObject{Type: MapObjectTypePlanet, Num: 2, Position: Vector{10, 10}},
	})

	player := game.Players[0]
	planet := game.Planets[0]
	target := game.Planets[1]
	fleet := game.Fleets[0]

	planet.RouteTargetType = MapObjectTypePlanet
	planet.RouteTargetNum = 2

	fleet.Waypoints[0].Task = WaypointTaskRoute

	turn := turn{
		game: game,
	}
	turn.game.Universe.buildMaps(game.Players)

	// route to planet 2
	// move
	turn.fleetRoute()

	assert.Equal(t, 2, len(fleet.Waypoints))
	assert.Equal(t, target.Num, fleet.Waypoints[1].TargetNum)
	assert.Equal(t, target.Type, fleet.Waypoints[1].TargetType)
	assert.Equal(t, target.PlayerNum, fleet.Waypoints[1].TargetPlayerNum)
	assert.Equal(t, 1, len(player.Messages))
}

func Test_turn_fleetMove(t *testing.T) {
	game := createSingleUnitGame()

	planet := game.Planets[0]
	fleet := game.Fleets[0]

	fleet.Waypoints = append(fleet.Waypoints, NewPositionWaypoint(Vector{10, 10}, 5))

	turn := turn{
		game: game,
	}
	turn.game.Universe.buildMaps(game.Players)

	// move to place
	turn.fleetMove()

	// should have consumed that waypoint and moved to the space
	assert.Equal(t, 1, len(fleet.Waypoints))
	assert.Equal(t, Vector{10, 10}, fleet.Position)
	assert.Equal(t, 1, len(game.getMapObjectsAtPosition(planet.Position)))
	assert.Equal(t, 1, len(game.getMapObjectsAtPosition(fleet.Position)))
}

func Test_turn_fleetMoveRepeatOrders(t *testing.T) {
	game := createSingleUnitGame()
	player := game.Players[0]

	planet := game.Planets[0]

	planet.Cargo = Cargo{1000, 1000, 1000, 1000}

	// make a new freighter for transport
	fleet := testSmallFreighter(player)
	player.Designs[0] = fleet.Tokens[0].design
	game.Fleets[0] = fleet

	// set a waypoint 2 turns away, load ironium from planet and move
	fleet.RepeatOrders = true
	fleet.OrbitingPlanetNum = planet.Num
	fleet.Waypoints[0] = NewPlanetWaypoint(planet.Position, planet.Num, planet.Name, 5)
	fleet.Waypoints[0].Task = WaypointTaskTransport
	fleet.Waypoints[0].TransportTasks.Ironium.Action = TransportActionLoadAll

	// dump all ironium at this waypoint and return
	fleet.Waypoints = append(fleet.Waypoints, NewPositionWaypoint(Vector{50, 0}, 5))
	fleet.Waypoints[1].Task = WaypointTaskTransport
	fleet.Waypoints[1].TransportTasks.Ironium.Action = TransportActionUnloadAll

	turn := turn{
		game: game,
	}
	turn.game.Universe.buildMaps(game.Players)

	// move one year
	turn.generateTurn()

	// should have loaded, moved, but still have waypoints
	assert.Equal(t, 120, fleet.Cargo.Ironium)
	assert.Equal(t, 880, planet.Cargo.Ironium)
	assert.Equal(t, Vector{25, 0}, fleet.Position)
	assert.Equal(t, 3, len(fleet.Waypoints))

	// generate the second turn, should move to dest and unload
	turn.generateTurn()
	assert.Equal(t, 0, fleet.Cargo.Ironium)
	assert.Equal(t, Vector{50, 0}, fleet.Position)
	assert.Equal(t, 2, len(fleet.Waypoints))
	// should have created salvage with ironium drop
	salvage := game.Salvages[0]
	assert.Equal(t, 120, salvage.Cargo.Ironium)

	// generate the third turn, should move back towards planet
	turn.generateTurn()
	assert.Equal(t, Vector{25, 0}, fleet.Position)
	assert.Equal(t, 3, len(fleet.Waypoints))

	// generate a fourth turn, should arrive at planet and load ironium
	turn.generateTurn()
	assert.Equal(t, Vector{0, 0}, fleet.Position)
	assert.Equal(t, 120, fleet.Cargo.Ironium)
	assert.Equal(t, 760, planet.Cargo.Ironium)
	assert.Equal(t, 2, len(fleet.Waypoints))

	// generate a fifth turn, should move again towards dest
	turn.generateTurn()
	assert.Equal(t, Vector{25, 0}, fleet.Position)
	assert.Equal(t, 3, len(fleet.Waypoints))

}

func Test_turn_fleetMoveTransportRepeat(t *testing.T) {
	game := createSingleUnitGame()
	player := game.Players[0]

	planet1 := game.Planets[0]
	// make a second planet we transfer cargo to
	planet2 := &Planet{
		MapObject: MapObject{Type: MapObjectTypePlanet, Name: "Planet 2", Num: 2, Position: Vector{10, 0}, PlayerNum: player.Num},
		Hab:       Hab{50, 50, 50},
		BaseHab:   Hab{50, 50, 50},
	}
	planet2.Spec = computePlanetSpec(&game.Rules, player, planet2)
	game.Planets = []*Planet{planet1, planet2}
	player.initDefaultPlanetIntels(&game.Rules, []*Planet{planet1, planet2})

	// planet1 has pop, planet2 is a starer colony
	planet1.Cargo = Cargo{1000, 1000, 1000, 10000}
	planet2.Cargo = Cargo{Colonists: 25}

	// make a new freighter for transport
	fleet := testGalleon(player)
	player.Designs[0] = fleet.Tokens[0].design
	game.Fleets[0] = fleet

	// set a waypoint nearby for the transport to load colonists from planet1 and dump them on planet2
	// until planet2 has 25% capacity
	fleet.RepeatOrders = true
	fleet.OrbitingPlanetNum = planet1.Num
	fleet.Waypoints[0] = NewPlanetWaypoint(planet1.Position, planet1.Num, planet1.Name, 5)
	fleet.Waypoints[0].Task = WaypointTaskTransport
	fleet.Waypoints[0].TransportTasks.Colonists.Action = TransportActionLoadAll

	// dump all colonists at this waypoint and return
	fleet.Waypoints = append(fleet.Waypoints, NewPlanetWaypoint(planet2.Position, planet2.Num, planet2.Name, 5))
	fleet.Waypoints[1].Task = WaypointTaskTransport
	fleet.Waypoints[1].TransportTasks.Colonists.Action = TransportActionSetWaypointTo
	fleet.Waypoints[1].TransportTasks.Colonists.Amount = 2500

	turn := turn{
		game: game,
	}
	turn.game.Universe.buildMaps(game.Players)

	// move one year
	turn.generateTurn()

	// should have loaded, moved, dropped
	assert.Equal(t, 10000-1000+150, planet1.Cargo.Colonists) // planet1 loaded colonists on freighter, then grew
	assert.Equal(t, Vector{10, 0}, fleet.Position)
	assert.Equal(t, Vector{10, 0}, fleet.Waypoints[0].Position)
	assert.Equal(t, MapObjectTypePlanet, fleet.Waypoints[0].TargetType)
	assert.Equal(t, planet2.Num, fleet.Waypoints[0].TargetNum)
	assert.Equal(t, 25+4+1000, planet2.Cargo.Colonists)
	assert.Equal(t, 2, len(fleet.Waypoints))

	// generate the second turn, should move back to planet1
	turn.generateTurn()

	// should have arrived back at homeworld, loaded
	assert.Equal(t, 8288, planet1.Cargo.Colonists)
	assert.Equal(t, Vector{0, 0}, fleet.Position)
	assert.Equal(t, Vector{0, 0}, fleet.Waypoints[0].Position)
	assert.Equal(t, MapObjectTypePlanet, fleet.Waypoints[0].TargetType)
	assert.Equal(t, planet1.Num, fleet.Waypoints[0].TargetNum)
	assert.Equal(t, 1000, fleet.Cargo.Colonists)
	assert.Equal(t, MapObjectTypePlanet, fleet.Waypoints[1].TargetType)
	assert.Equal(t, planet2.Num, fleet.Waypoints[1].TargetNum)
	assert.Equal(t, 2, len(fleet.Waypoints))

	// generate the third turn, should move back to planet2 and unload
	turn.generateTurn()

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
	turn.generateTurn()
	// p1 -> p2
	turn.generateTurn()

	assert.Equal(t, 7956, planet1.Cargo.Colonists)
	assert.Equal(t, Vector{10, 0}, fleet.Position)
	assert.Equal(t, Vector{10, 0}, fleet.Waypoints[0].Position)
	assert.Equal(t, MapObjectTypePlanet, fleet.Waypoints[0].TargetType)
	assert.Equal(t, planet2.Num, fleet.Waypoints[0].TargetNum)
	assert.Equal(t, Cargo{Colonists: 1000}, fleet.Cargo) // we have leftover
	assert.Equal(t, 3121, planet2.Cargo.Colonists)       // planet is ready to go!
	assert.Equal(t, 2, len(fleet.Waypoints))

}

func Test_turn_fleetMoveTransportWaitForPercent(t *testing.T) {
	game := createSingleUnitGame()
	player := game.Players[0]

	planet1 := game.Planets[0]
	// make a second planet we transfer cargo to
	planet2 := &Planet{
		MapObject: MapObject{Type: MapObjectTypePlanet, Name: "Planet 2", Num: 2, Position: Vector{10, 0}, PlayerNum: player.Num},
		Hab:       Hab{50, 50, 50},
		BaseHab:   Hab{50, 50, 50},
	}
	planet2.Spec = computePlanetSpec(&game.Rules, player, planet2)
	game.Planets = []*Planet{planet1, planet2}
	player.initDefaultPlanetIntels(&game.Rules, []*Planet{planet1, planet2})

	// pull from planet1 to planet2
	planet1.MineralConcentration = Mineral{100, 100, 100}
	planet1.Mines = 300
	planet1.Cargo = Cargo{100, 100, 100, 10000} // start with cargo, mine the rest
	planet2.Cargo = Cargo{0, 0, 0, 1000}

	// make a new freighter for transport
	fleet := testGalleon(player)
	player.Designs[0] = fleet.Tokens[0].design
	game.Fleets[0] = fleet

	// set a waypoint nearby for the transport to wait until we have an even amount of cargo in the hold, then dump on planet2
	fleet.RepeatOrders = true
	fleet.OrbitingPlanetNum = planet1.Num
	fleet.Waypoints[0] = NewPlanetWaypoint(planet1.Position, planet1.Num, planet1.Name, 5)
	fleet.Waypoints[0].Task = WaypointTaskTransport
	fleet.Waypoints[0].TransportTasks.Ironium.Action = TransportActionWaitForPercent
	fleet.Waypoints[0].TransportTasks.Ironium.Amount = 33
	fleet.Waypoints[0].TransportTasks.Boranium.Action = TransportActionWaitForPercent
	fleet.Waypoints[0].TransportTasks.Boranium.Amount = 33
	fleet.Waypoints[0].TransportTasks.Germanium.Action = TransportActionWaitForPercent
	fleet.Waypoints[0].TransportTasks.Germanium.Amount = 34

	// dump all at this waypoint and return
	fleet.Waypoints = append(fleet.Waypoints, NewPlanetWaypoint(planet2.Position, planet2.Num, planet2.Name, 5))
	fleet.Waypoints[1].Task = WaypointTaskTransport
	fleet.Waypoints[1].TransportTasks.Ironium.Action = TransportActionUnloadAll
	fleet.Waypoints[1].TransportTasks.Boranium.Action = TransportActionUnloadAll
	fleet.Waypoints[1].TransportTasks.Germanium.Action = TransportActionUnloadAll

	turn := turn{
		game: game,
	}
	turn.game.Universe.buildMaps(game.Players)

	// load and grow and wait
	turn.generateTurn()

	// should have loaded all cargo, but waited for more to be generated
	assert.Equal(t, Vector{0, 0}, fleet.Position)
	assert.Equal(t, Vector{0, 0}, fleet.Waypoints[0].Position)
	assert.Equal(t, MapObjectTypePlanet, fleet.Waypoints[0].TargetType)
	assert.Equal(t, planet1.Num, fleet.Waypoints[0].TargetNum)
	assert.Equal(t, Cargo{100, 100, 100, 0}, fleet.Cargo)
	assert.Equal(t, 2, len(fleet.Waypoints))

	// we should load the rest and move
	turn.generateTurn()

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

}

func Test_turn_fleetMoveDestroyedByMineField(t *testing.T) {
	game := createSingleUnitGame()
	rules := &game.Rules

	// change the rules so going 4 warp over the limit guarantee's a hit
	stats := rules.MineFieldStatsByType[MineFieldTypeStandard]
	stats.MaxSpeed = 5
	stats.ChanceOfHit = .25
	stats.MinDecay = 0 // turn off decay
	rules.MineFieldStatsByType[MineFieldTypeStandard] = stats

	// create a new MineField 20ly away with 10ly radius
	radius := 10
	mineFieldPlayer := NewPlayer(2, NewRace().WithSpec(rules)).WithNum(2).withSpec(rules)
	mineFieldPlayer.Race.Spec.MineFieldBaseDecayRate = 0
	mineFieldPlayer.Race.Spec.MineFieldMinDecayFactor = 0
	mineFieldPlayer.Race.Spec.MineFieldMaxDecayRate = 0
	mineField := newMineField(mineFieldPlayer, MineFieldTypeStandard, radius*radius, 1, Vector{20, 0})
	mineField.Spec = computeMinefieldSpec(rules, mineFieldPlayer, mineField, 0)
	// setup initial planet intels so turn generation works
	mineFieldPlayer.initDefaultPlanetIntels(&game.Rules, game.Planets)

	// make sure our player doesn't gain any tech levels since we're checking messages after turn generation
	player := game.Players[0]
	player.TechLevels = TechLevel{26, 26, 26, 26, 26, 26}

	game.Players = append(game.Players, mineFieldPlayer)
	game.MineFields = append(game.MineFields, mineField)

	// move us straight through a minefield
	fleet := game.Fleets[0]
	fleet.Waypoints = append(fleet.Waypoints, NewPositionWaypoint(Vector{81, 0}, 9))

	turn := turn{
		game: game,
	}
	turn.game.Universe.buildMaps(game.Players)

	// let's go!!
	turn.generateTurn()

	// we should have struck the minefield and lost the ship
	assert.True(t, fleet.Delete)
	assert.Equal(t, 1, len(game.Players[0].Messages))
	assert.Equal(t, 2, len(game.Players[1].Messages))

	// the MineField should have lost some mines in the collision
	assert.Equal(t, 88, mineField.NumMines)
}

func Test_turn_permaform(t *testing.T) {
	game := createSingleUnitGame()

	player := game.Players[0]
	planet := game.Planets[0]
	planet.Hab = Hab{49, 49, 49}

	turn := turn{
		game: game,
	}
	turn.game.Universe.buildMaps(game.Players)

	// 100% chance to permaform
	player.Race.Spec.PermaformChance = 1
	player.Race.Spec.PermaformPopulation = 0

	// mock the random number generator to return temp as the hab to permaform
	mockRand := MockRand{}
	mockRand.int63Result = int64(Temp) << 32 // rand.Intn calls this int63 and >> 32 the result
	game.Rules.random = rand.New(mockRand)

	turn.permaform()

	// should have permaformed the planet temp in one direction
	assert.Equal(t, Hab{49, 50, 49}, planet.Hab)
}

func Test_turn_fleetRemoteMine(t *testing.T) {

	type fields struct {
		task                    WaypointTask
		planetPlayerNum         int
		orbitingPlanetNum       int
		miningRate              int
		canRemoteMineOwnPlanets bool
	}

	tests := []struct {
		name            string
		fields          fields
		wantCargo       Cargo
		wantMessageType PlayerMessageType
	}{
		{name: "no task, do nothing", fields: fields{}, wantCargo: Cargo{}, wantMessageType: PlayerMessageNone},
		{name: "no planet, invalid message", fields: fields{task: WaypointTaskRemoteMining}, wantCargo: Cargo{}, wantMessageType: PlayerMessageInvalid},
		{name: "owned planet, invalid message", fields: fields{task: WaypointTaskRemoteMining, planetPlayerNum: 2, orbitingPlanetNum: 2}, wantCargo: Cargo{}, wantMessageType: PlayerMessageInvalid},
		{name: "owned by us, invalid", fields: fields{task: WaypointTaskRemoteMining, planetPlayerNum: 1, orbitingPlanetNum: 2}, wantCargo: Cargo{}, wantMessageType: PlayerMessageInvalid},
		{name: "owned by us, but we can remote mine our own, should skip", fields: fields{task: WaypointTaskRemoteMining, planetPlayerNum: 1, orbitingPlanetNum: 2, canRemoteMineOwnPlanets: true}, wantCargo: Cargo{}, wantMessageType: PlayerMessageNone},
		{name: "no miners, invalid message", fields: fields{task: WaypointTaskRemoteMining, orbitingPlanetNum: 2}, wantCargo: Cargo{}, wantMessageType: PlayerMessageInvalid},
		{name: "should mine", fields: fields{task: WaypointTaskRemoteMining, orbitingPlanetNum: 2, miningRate: 10}, wantCargo: Cargo{10, 10, 10, 0}, wantMessageType: PlayerMessageRemoteMined},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// create a new test game
			game := createSingleUnitGame()
			player := game.Players[0]
			fleet := game.Fleets[0]

			// setup params
			player.Race.Spec.CanRemoteMineOwnPlanets = tt.fields.canRemoteMineOwnPlanets
			fleet.Waypoints[0].Task = tt.fields.task
			fleet.Spec.MiningRate = tt.fields.miningRate
			fleet.OrbitingPlanetNum = tt.fields.orbitingPlanetNum

			// add a new planet
			planet := &Planet{
				MapObject:            MapObject{Type: MapObjectTypePlanet, Name: "Planet 2", Num: 2, PlayerNum: tt.fields.planetPlayerNum},
				MineralConcentration: Mineral{100, 100, 100},
			}
			planet.Spec = computePlanetSpec(&game.Rules, player, planet)
			game.Planets = append(game.Planets, planet)

			turn := turn{
				game: game,
			}
			turn.game.Universe.buildMaps(game.Players)

			// try and remote the planet
			turn.fleetRemoteMine()

			if tt.wantMessageType != PlayerMessageNone {
				assert.Equal(t, 1, len(player.Messages))
				assert.Equal(t, tt.wantMessageType, player.Messages[0].Type)
			} else {
				assert.Equal(t, 0, len(player.Messages))
			}

			// make sure the cargo matches what we want
			assert.Equal(t, tt.wantCargo, planet.Cargo)
		})
	}

}

func Test_turn_fleetRemoteMineAR(t *testing.T) {

	type fields struct {
		task                    WaypointTask
		planetPlayerNum         int
		orbitingPlanetNum       int
		miningRate              int
		canRemoteMineOwnPlanets bool
	}

	tests := []struct {
		name            string
		fields          fields
		wantCargo       Cargo
		wantMessageType PlayerMessageType
	}{
		{name: "no task, do nothing", fields: fields{}, wantCargo: Cargo{}, wantMessageType: PlayerMessageNone},
		{name: "no planet, invalid message", fields: fields{task: WaypointTaskRemoteMining}, wantCargo: Cargo{}, wantMessageType: PlayerMessageInvalid},
		{name: "no miners, invalid message", fields: fields{task: WaypointTaskRemoteMining, orbitingPlanetNum: 2, planetPlayerNum: 1, canRemoteMineOwnPlanets: true}, wantCargo: Cargo{}, wantMessageType: PlayerMessageInvalid},
		{name: "owned by us, but we can remote mine our own, should mine", fields: fields{task: WaypointTaskRemoteMining, orbitingPlanetNum: 2, planetPlayerNum: 1, canRemoteMineOwnPlanets: true, miningRate: 10}, wantCargo: Cargo{10, 10, 10, 0}, wantMessageType: PlayerMessageRemoteMined},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// create a new test game
			game := createSingleUnitGame()
			player := game.Players[0]
			fleet := game.Fleets[0]

			// setup params
			player.Race.Spec.CanRemoteMineOwnPlanets = tt.fields.canRemoteMineOwnPlanets
			fleet.Waypoints[0].Task = tt.fields.task
			fleet.Spec.MiningRate = tt.fields.miningRate
			fleet.OrbitingPlanetNum = tt.fields.orbitingPlanetNum

			// add a new planet
			planet := &Planet{
				MapObject:            MapObject{Type: MapObjectTypePlanet, Name: "Planet 2", Num: 2, PlayerNum: tt.fields.planetPlayerNum},
				MineralConcentration: Mineral{100, 100, 100},
			}
			planet.Spec = computePlanetSpec(&game.Rules, player, planet)
			game.Planets = append(game.Planets, planet)

			turn := turn{
				game: game,
			}
			turn.game.Universe.buildMaps(game.Players)

			// try and remote the planet
			turn.fleetRemoteMineAR()

			if tt.wantMessageType != PlayerMessageNone {
				assert.Equal(t, 1, len(player.Messages))
				assert.Equal(t, tt.wantMessageType, player.Messages[0].Type)
			} else {
				assert.Equal(t, 0, len(player.Messages))
			}

			// make sure the cargo matches what we want
			assert.Equal(t, tt.wantCargo, planet.Cargo)
		})
	}

}

func Test_turn_fleetLayMines(t *testing.T) {
	game := createSingleUnitGame()
	player := game.Players[0]

	// make a new freighter for transport
	fleet := testMiniMineLayer(player)
	player.Designs[0] = fleet.Tokens[0].design
	game.Fleets[0] = fleet

	// set a waypoint 2 turns away, load ironium from planet and move
	fleet.Waypoints[0].Task = WaypointTaskLayMineField

	turn := turn{
		game: game,
	}
	turn.game.Universe.buildMaps(game.Players)

	// move one year
	turn.generateTurn()

	// should have loaded, moved, but still have waypoints
	assert.Equal(t, 1, len(game.MineFields))
	mineField := game.MineFields[0]
	assert.Equal(t, 320, mineField.NumMines)
	assert.Equal(t, math.Sqrt(320), mineField.Spec.Radius)
	assert.Equal(t, Vector{0, 0}, mineField.Position)

}

func Test_turn_fleetSweepMines(t *testing.T) {
	game := createSingleUnitGame()
	rules := &game.Rules

	// change the rules so we don't decay
	stats := rules.MineFieldStatsByType[MineFieldTypeStandard]
	stats.MinDecay = 0 // turn off decay
	rules.MineFieldStatsByType[MineFieldTypeStandard] = stats

	player := game.Players[0]

	// make a new destroyer to sweep mines
	fleet := testStalwartDefender(player)
	player.Designs[0] = fleet.Tokens[0].design
	game.Fleets[0] = fleet

	// create a new MineField with 100 mines
	// and make it not decay
	radius := 10
	mineFieldPlayer := NewPlayer(2, NewRace().WithSpec(rules)).WithNum(2).withSpec(rules)
	mineFieldPlayer.Race.Spec.MineFieldBaseDecayRate = 0
	mineFieldPlayer.Race.Spec.MineFieldMinDecayFactor = 0
	mineFieldPlayer.Race.Spec.MineFieldMaxDecayRate = 0
	mineField := newMineField(mineFieldPlayer, MineFieldTypeStandard, radius*radius, 1, Vector{0, 0})
	mineField.Spec = computeMinefieldSpec(rules, mineFieldPlayer, mineField, 0)
	// setup initial planet intels so turn generation works
	mineFieldPlayer.initDefaultPlanetIntels(&game.Rules, game.Planets)

	game.Players = append(game.Players, mineFieldPlayer)
	game.MineFields = append(game.MineFields, mineField)

	player.Relations = []PlayerRelationship{{Relation: PlayerRelationFriend}, {Relation: PlayerRelationNeutral}}
	mineFieldPlayer.Relations = []PlayerRelationship{{Relation: PlayerRelationNeutral}, {Relation: PlayerRelationFriend}}

	turn := turn{
		game: game,
	}
	turn.game.Universe.buildMaps(game.Players)

	// sweep mines
	turn.generateTurn()

	// we should clear out some mines
	assert.Equal(t, 78, mineField.NumMines)
	assert.False(t, mineField.Delete)

	// upgrade a mine sweeper weapon
	fleet.Tokens[0].design.Slots[1].HullComponent = GatlingNeutrinoCannon.Name
	fleet.Tokens[0].design.Spec = ComputeShipDesignSpec(rules, player.TechLevels, player.Race.Spec, fleet.Tokens[0].design)
	fleet.Spec = ComputeFleetSpec(rules, player, fleet)

	// sweep mines
	turn.generateTurn()

	// we should clear out some mines
	assert.Equal(t, 0, mineField.NumMines)
	assert.True(t, mineField.Delete)
}

func Test_turn_instaform(t *testing.T) {
	game := createSingleUnitGame()
	player := game.Players[0]
	planet := game.Planets[0]

	// allow Grav3 terraform
	player.Race.PRT = CA
	player.TechLevels = TechLevel{Propulsion: 1, Biotechnology: 1}

	planet.Hab = Hab{45, 50, 50}
	planet.BaseHab = Hab{45, 50, 50}
	planet.TerraformedAmount = Hab{}

	turn := turn{
		game: game,
	}
	turn.game.Universe.buildMaps(game.Players)

	// instaform
	turn.generateTurn()

	// should terraform 3 grav points
	// TODO: sometimes this fails if we randomly permaform...
	assert.Equal(t, Hab{48, 50, 50}, planet.Hab)
}

func Test_turn_fleetRepair(t *testing.T) {
	game := createSingleUnitGame()
	player := game.Players[0]
	fleet := game.Fleets[0]
	planet := game.Planets[0]

	// create a new starbase
	starbaseDesign := NewShipDesign(player, 2).WithHull(SpaceStation.Name).WithSpec(&rules, player)
	starbase := newStarbase(player, planet,
		starbaseDesign,
		"Starbase",
	)
	player.Designs = append(player.Designs, starbaseDesign)
	starbase.Spec = ComputeFleetSpec(&rules, player, &starbase)
	starbase.Tokens[0].QuantityDamaged = 1
	starbase.Tokens[0].Damage = 100
	game.Starbases = append(game.Starbases, &starbase)
	planet.Starbase = &starbase

	turn := turn{
		game: game,
	}
	turn.game.Universe.buildMaps(game.Players)

	// in space with 10 damage
	fleet.OrbitingPlanetNum = None
	fleet.Tokens[0].QuantityDamaged = 1
	fleet.Tokens[0].Damage = 10

	// repair
	turn.generateTurn()

	// should repair fleet and starbase
	assert.Equal(t, 9.0, fleet.Tokens[0].Damage)
	assert.Equal(t, 50.0, starbase.Tokens[0].Damage)

}

func Test_turn_fleetReproduce(t *testing.T) {
	game := createSingleUnitGame()

	// make an IS race for reproducing
	player := game.Players[0]
	player.Race.PRT = IS
	player.Race.Spec = computeRaceSpec(&player.Race, &rules)

	// make a new freighter with some colonists
	fleet := testSmallFreighter(player)
	fleet.Cargo.Colonists = 50 // 5000 colonists
	player.Designs[0] = fleet.Tokens[0].design
	game.Fleets[0] = fleet

	// orbit a planet of ours
	planet := game.Planets[0]
	planet.PlayerNum = player.Num
	planet.Cargo.Colonists = 2500
	fleet.Waypoints[0] = NewPlanetWaypoint(planet.Position, planet.Num, planet.Name, 5)
	fleet.OrbitingPlanetNum = planet.Num

	turn := turn{
		game: game,
	}
	turn.game.Universe.buildMaps(game.Players)

	// don't generate a full turn, the planet will grow
	turn.fleetReproduce()

	// should have grown on freighter
	assert.Equal(t, 53, fleet.Cargo.Colonists)
	assert.Equal(t, 2500, planet.Cargo.Colonists)

	// fill it up, should overflow onto planet
	fleet.Cargo.Colonists = fleet.Spec.CargoCapacity

	// reproduce again
	turn.fleetReproduce()

	// should have grown on freighter and beamed down to planet
	assert.Equal(t, fleet.Spec.CargoCapacity, fleet.Cargo.Colonists)
	assert.Equal(t, 2509, planet.Cargo.Colonists) // 120kT * 7.5% = 900 colonists beamed to planet

}

func Test_turn_fleetDieoff(t *testing.T) {
	game := createSingleUnitGame()

	// make an AR race for dieoff
	player := game.Players[0]
	player.Race.PRT = AR
	player.Race.Spec = computeRaceSpec(&player.Race, &rules)

	// make a new freighter with some colonists
	fleet := testSmallFreighter(player)
	fleet.Cargo.Colonists = 50 // 5000 colonists
	player.Designs[0] = fleet.Tokens[0].design
	game.Fleets[0] = fleet

	turn := turn{
		game: game,
	}
	turn.game.Universe.buildMaps(game.Players)

	turn.fleetDieoff()

	// should have lost a min of 1kt on freighter
	assert.Equal(t, 49, fleet.Cargo.Colonists)

	// set to 10000, so we lose 300 or 3%
	fleet.Cargo.Colonists = 100
	turn.fleetDieoff()

	// should have lost a min of 1kt on freighter
	assert.Equal(t, 97, fleet.Cargo.Colonists)
}

func Test_turn_fleetRadiatingEngineDieoff(t *testing.T) {
	game := createSingleUnitGame()

	// make an IS race for reproducing
	player := game.Players[0]

	// make a new freighter with some colonists
	fleet := testSmallFreighter(player)
	fleet.Cargo.Colonists = 50 // 5000 colonists
	design := fleet.Tokens[0].design
	player.Designs[0] = design
	game.Fleets[0] = fleet

	turn := turn{
		game: game,
	}
	turn.game.Universe.buildMaps(game.Players)

	// generate turn to simulate die off
	turn.generateTurn()

	// should not die
	assert.Equal(t, 50, fleet.Cargo.Colonists)

	// add a radiating hydro ramscoop
	design.Slots[0].HullComponent = RadiatingHydroRamScoop.Name
	design.Spec = ComputeShipDesignSpec(&rules, player.TechLevels, player.Race.Spec, design)
	fleet.Spec = ComputeFleetSpec(&rules, player, fleet)

	// generate turn to simulate die off
	turn.generateTurn()

	// should not die
	assert.Equal(t, 41, fleet.Cargo.Colonists)

	// make the player a high rad race
	player.Race.HabHigh.Rad = 100
	player.Race.HabLow.Rad = 80
	player.Race.Spec = computeRaceSpec(&player.Race, &rules)
	turn.generateTurn()

	// should not die
	assert.Equal(t, 41, fleet.Cargo.Colonists)

}

func Test_turn_detonateMines(t *testing.T) {

	mineFieldPlayer := NewPlayer(1, NewRace().WithPRT(SD).WithSpec(&rules)).WithNum(1).withSpec(&rules)
	otherPlayer := NewPlayer(2, NewRace().WithSpec(&rules)).WithNum(2).withSpec(&rules)

	type args struct {
		mineField *MineField
		fleet     *Fleet
		players   []*Player
	}

	tests := []struct {
		name string
		args args
		want []ShipToken
	}{
		{
			name: "no op",
			args: args{
				mineField: newMineField(mineFieldPlayer, MineFieldTypeStandard, 10*10, 1, Vector{}),
				fleet:     testLongRangeScout(otherPlayer),
				players:   []*Player{mineFieldPlayer, otherPlayer},
			},
			want: []ShipToken{{Quantity: 1, QuantityDamaged: 0, Damage: 0}},
		},
		{
			name: "detonate, destroy ship",
			args: args{
				mineField: newMineField(mineFieldPlayer, MineFieldTypeStandard, 10*10, 1, Vector{}).
					withOrders(MineFieldOrders{Detonate: true}),
				fleet:   testLongRangeScout(otherPlayer),
				players: []*Player{mineFieldPlayer, otherPlayer},
			},
			want: []ShipToken{{Quantity: 0, QuantityDamaged: 0, Damage: 500}},
		},
		{
			name: "detonate, don't destroy mini mine layers",
			args: args{
				mineField: newMineField(mineFieldPlayer, MineFieldTypeStandard, 10*10, 1, Vector{}).
					withOrders(MineFieldOrders{Detonate: true}),
				fleet:   testMiniMineLayer(mineFieldPlayer),
				players: []*Player{mineFieldPlayer},
			},
			want: []ShipToken{{Quantity: 1, QuantityDamaged: 0, Damage: 0}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// create a new test game
			client := NewGamer()
			game := client.CreateGame(1, *NewGameSettings())

			for _, player := range tt.args.players {
				player.Relations = player.defaultRelationships(tt.args.players, false)
				player.PlayerIntels.PlayerIntels = player.defaultPlayerIntels(tt.args.players)
			}

			universe := NewUniverse(&game.Rules)
			universe.Fleets = []*Fleet{tt.args.fleet}
			universe.MineFields = []*MineField{tt.args.mineField}

			fg := FullGame{
				Game:      game,
				Universe:  &universe,
				TechStore: &StaticTechStore,
				Players:   tt.args.players,
			}

			// make sure the player knows about the designs of the fleet
			fleetPlayer := fg.getPlayer(tt.args.fleet.PlayerNum)
			for _, token := range tt.args.fleet.Tokens {
				fleetPlayer.Designs = append(fleetPlayer.Designs, token.design)
			}

			turn := turn{
				game: &fg,
			}
			turn.game.Universe.buildMaps(fg.Players)

			// try and remote the planet
			turn.generateTurn()

			for i := range tt.args.fleet.Tokens {
				token := tt.args.fleet.Tokens[i]
				if token.Quantity != tt.want[i].Quantity {
					t.Errorf("Fleet.detonateMines() token %d gotQuantity = %v, wantQuantity %v", i, token.Quantity, tt.want[i].Quantity)
				}
				if token.Damage != tt.want[i].Damage {
					t.Errorf("Fleet.detonateMines() token %d gotDamage = %v, wantDamage %v", i, token.Damage, tt.want[i].Damage)
				}
				if token.QuantityDamaged != tt.want[i].QuantityDamaged {
					t.Errorf("Fleet.detonateMines() token %d gotQuantityDamaged = %v, wantQuantityDamaged %v", i, token.QuantityDamaged, tt.want[i].QuantityDamaged)
				}
			}

		})
	}

}

func Test_turn_decayPackets(t *testing.T) {
	game := createSingleUnitGame()
	player := game.Players[0]

	// create some Packets with 100kT of each mineral
	packetSafe := newMineralPacket(player, 1, 5, 5, Cargo{100, 100, 100, 0}, Vector{200, 0}, 1)
	packetTooFast := newMineralPacket(player, 1, 10, 7, Cargo{500, 500, 500, 0}, Vector{0, 200}, 1)
	packetNewlyBuilt := newMineralPacket(player, 1, 8, 5, Cargo{100, 100, 100, 0}, Vector{200, 0}, 1)
	packetNewlyBuilt.builtThisTurn = true

	game.MineralPackets = append(game.MineralPackets, packetSafe, packetTooFast, packetNewlyBuilt)

	turn := turn{
		game: game,
	}
	turn.game.Universe.buildMaps(game.Players)

	// move and decay
	turn.packetMove(false)
	turn.packetMove(true)
	turn.decayPackets()

	// no decay, 50% decay, and half of 50% decay for a newly built overfast packet
	assert.Equal(t, packetSafe.Cargo, Cargo{100, 100, 100, 0})
	assert.Equal(t, packetTooFast.Cargo, Cargo{250, 250, 250, 0})
	assert.Equal(t, packetNewlyBuilt.Cargo, Cargo{75, 75, 75, 0})
}

func Test_turn_fleetPatrol(t *testing.T) {
	game := createSingleUnitGame()
	rules := &game.Rules

	// make a new destroyer to patrol within 50ly
	player := game.Players[0]
	fleet := testStalwartDefender(player)
	fleet.Waypoints[0].Task = WaypointTaskPatrol
	fleet.Waypoints[0].PatrolRange = 50
	game.Fleets[0] = fleet

	// create a new enemy player with a fleet 100ly away
	enemyPlayer := NewPlayer(2, NewRace().WithSpec(rules)).WithNum(2).withSpec(rules)
	enemyFleet := testLongRangeScout(enemyPlayer)
	enemyFleet.Position = Vector{60, 0}
	game.Players = append(game.Players, enemyPlayer)
	game.Fleets = append(game.Fleets, enemyFleet)

	player.Relations = []PlayerRelationship{{Relation: PlayerRelationFriend}, {Relation: PlayerRelationNeutral}}
	enemyPlayer.Relations = []PlayerRelationship{{Relation: PlayerRelationNeutral}, {Relation: PlayerRelationFriend}}
	player.PlayerIntels.PlayerIntels = player.defaultPlayerIntels([]*Player{player, enemyPlayer})
	enemyPlayer.PlayerIntels.PlayerIntels = player.defaultPlayerIntels([]*Player{player, enemyPlayer})
	// setup initial planet intels so turn generation works
	enemyPlayer.initDefaultPlanetIntels(&game.Rules, game.Planets)

	turn := turn{
		game: game,
	}
	turn.game.Universe.buildMaps(game.Players)

	// no patrol target
	err := turn.generateTurn()
	if err != nil {
		t.Error("failed to generate turn", err)
	}

	// should not attack
	assert.Equal(t, 1, len(fleet.Waypoints))

	// move closer, should target
	enemyFleet.Position = Vector{50, 0}
	turn.generateTurn()

	// should not attack
	assert.Equal(t, len(fleet.Waypoints), 2)
	assert.Equal(t, MapObjectTypeFleet, fleet.Waypoints[1].TargetType)
	assert.Equal(t, enemyFleet.PlayerNum, fleet.Waypoints[1].TargetPlayerNum)
	assert.Equal(t, enemyFleet.Num, fleet.Waypoints[1].TargetNum)
	assert.Equal(t, 6, fleet.Waypoints[1].WarpSpeed)

}

func Test_turn_fleetRemoteTerraform(t *testing.T) {
	game := createSingleUnitGame()
	rules := &game.Rules

	// give player TT so it can terraform these other worlds
	player := game.Players[0]
	player.Race.LRTs |= Bitmask(TT)
	player.Race.Spec = computeRaceSpec(&player.Race, rules)

	// make two new remote terraformers, one over each planet to test deterraforming and terraforming
	fleet1 := testRemoteTerraformer(player)
	fleet2 := testRemoteTerraformer(player)
	player.Designs[0] = fleet1.Tokens[0].design
	game.Fleets = []*Fleet{fleet1, fleet2}

	// create a new enemy player and friendly player to own planets
	enemyPlayer := NewPlayer(2, NewRace().WithSpec(rules)).WithNum(2).withSpec(rules)
	friendlyPlayer := NewPlayer(3, NewRace().WithSpec(rules)).WithNum(3).withSpec(rules)
	game.Players = append(game.Players, enemyPlayer, friendlyPlayer)

	player.Relations = []PlayerRelationship{{Relation: PlayerRelationFriend}, {Relation: PlayerRelationNeutral}, {Relation: PlayerRelationFriend}}
	enemyPlayer.Relations = []PlayerRelationship{{Relation: PlayerRelationNeutral}, {Relation: PlayerRelationFriend}, {Relation: PlayerRelationNeutral}}
	friendlyPlayer.Relations = []PlayerRelationship{{Relation: PlayerRelationFriend}, {Relation: PlayerRelationNeutral}, {Relation: PlayerRelationFriend}}
	player.PlayerIntels.PlayerIntels = player.defaultPlayerIntels([]*Player{player, enemyPlayer, friendlyPlayer})
	enemyPlayer.PlayerIntels.PlayerIntels = player.defaultPlayerIntels([]*Player{player, enemyPlayer, friendlyPlayer})
	friendlyPlayer.PlayerIntels.PlayerIntels = player.defaultPlayerIntels([]*Player{player, enemyPlayer, friendlyPlayer})

	// give planet1 to the enemy and orbit it with fleet1
	planet1 := &Planet{
		MapObject: MapObject{Type: MapObjectTypePlanet, Name: "Planet 1", Num: 1, PlayerNum: enemyPlayer.Num},
		Cargo: Cargo{
			Colonists: 2500,
		},
		Hab:     Hab{50, 50, 50},
		BaseHab: Hab{50, 50, 50},
	}
	planet1.Spec = computePlanetSpec(&game.Rules, player, planet1)
	fleet1.OrbitingPlanetNum = planet1.Num

	// give planet2 to the friend and orbit it with fleet2
	planet2 := &Planet{
		MapObject: MapObject{Type: MapObjectTypePlanet, Name: "Planet 2", Num: 2, PlayerNum: friendlyPlayer.Num},
		Cargo: Cargo{
			Colonists: 2500,
		},
		Hab:     Hab{48, 50, 50},
		BaseHab: Hab{48, 50, 50},
	}
	planet2.Spec = computePlanetSpec(&game.Rules, player, planet2)
	fleet2.OrbitingPlanetNum = planet2.Num

	game.Planets = []*Planet{planet1, planet2}
	player.initDefaultPlanetIntels(&game.Rules, []*Planet{planet1, planet2})
	enemyPlayer.initDefaultPlanetIntels(&game.Rules, []*Planet{planet1, planet2})
	friendlyPlayer.initDefaultPlanetIntels(&game.Rules, []*Planet{planet1, planet2})

	turn := turn{
		game: game,
	}
	turn.game.Universe.buildMaps(game.Players)

	err := turn.generateTurn()
	if err != nil {
		t.Error("failed to generate turn", err)
	}

	// should deterraform planet1 2 points
	assert.Equal(t, Hab{52, 50, 50}, planet1.Hab)

	// should terraform planet2 2 points
	assert.Equal(t, Hab{50, 50, 50}, planet2.Hab)

}

func Test_turn_fleetRefuel(t *testing.T) {
	game := createSingleUnitGame()
	player := game.Players[0]
	fleet := game.Fleets[0]
	planet := game.Planets[0]

	// create a new starbase
	starbaseDesign := NewShipDesign(player, 2).WithHull(SpaceStation.Name).WithSpec(&rules, player)
	starbase := newStarbase(player, planet,
		starbaseDesign,
		"Starbase",
	)
	player.Designs = append(player.Designs, starbaseDesign)
	starbase.Spec = ComputeFleetSpec(&rules, player, &starbase)
	starbase.Tokens[0].QuantityDamaged = 1
	starbase.Tokens[0].Damage = 100
	game.Starbases = append(game.Starbases, &starbase)
	planet.Starbase = &starbase

	turn := turn{
		game: game,
	}
	turn.game.Universe.buildMaps(game.Players)

	// orbit and need fuel
	fleet.Fuel = 0

	// repair
	turn.generateTurn()

	// should refuel at starbase
	assert.Equal(t, fleet.Spec.FuelCapacity, fleet.Fuel)

}

func Test_turn_playerResearch(t *testing.T) {
	game := createSingleUnitGame()
	planet := game.Planets[0]
	player := game.Players[0]

	// research energy and keep going
	player.TechLevels = TechLevel{}
	player.TechLevelsSpent = TechLevel{}
	player.NextResearchField = NextResearchFieldEnergy
	player.Researching = Energy
	player.ResearchAmount = 100

	// give us 1000 resources, 500 from pop, 500 from factories
	planet.setPopulation(500_000)
	planet.Factories = 500

	turn := turn{
		game: game,
	}
	turn.game.Universe.buildMaps(game.Players)

	// let's go!!
	turn.generateTurn()

	// costs go up per level so
	// 1. 50,
	// 2. 80 + 10,
	// 3. 130 + 20,
	// 4. 210 + 30,
	// 5. 340 + 40 (910 resources for level 5)
	// 6. 550 (we'll bleed over 90, leaving 360 for this next level)

	assert.Equal(t, TechLevel{Energy: 5}, player.TechLevels, "should have raised 5 energy levels")
	assert.Equal(t, TechLevel{Energy: 90}, player.TechLevelsSpent, "should leave 90 spent on energy level 6")
	assert.Equal(t, 550+(5*10), player.Spec.CurrentResearchCost, "total cost of energy level 6 if we have 5 total tech levels")
	assert.Equal(t, 1000, player.ResearchSpentLastYear, " spent 1000 on research last year")
	assert.Equal(t, 1045, player.Spec.ResourcesPerYear, " planet grew, spend more next year")
	assert.Equal(t, 1045, player.Spec.ResourcesPerYearResearch)
}

func Test_turn_buildStarbase(t *testing.T) {
	game := createSingleUnitGame()
	player := game.Players[0]
	planet := game.Planets[0]

	// create a new starbase design
	starbaseDesign := NewShipDesign(player, 2).WithName("Sad empty base").WithHull(SpaceStation.Name).WithSpec(&rules, player)
	player.Designs = append(player.Designs, starbaseDesign)

	// build a starbase
	planet.ProductionQueue = append(planet.ProductionQueue, ProductionQueueItem{Type: QueueItemTypeStarbase, Quantity: 1, DesignNum: starbaseDesign.Num})
	planet.Cargo = Mineral{1000, 1000, 1000}.ToCargo()
	planet.setPopulation(1_000_000)
	planet.Factories = 1000

	turn := turn{
		game: game,
	}
	turn.game.Universe.buildMaps(game.Players)

	// should have no starbase
	assert.Nil(t, planet.Starbase)
	assert.Equal(t, 0, len(game.Starbases))

	// generate a turn to build a starbase
	turn.generateTurn()

	// should have a starbase at the planet
	assert.NotNil(t, planet.Starbase)
	assert.Equal(t, 1, len(game.Starbases))

	// upgrade the starbase with A LASER!
	starbaseDesignUpgrade := NewShipDesign(player, 3).WithName("LASER BASE!").WithHull(SpaceStation.Name).WithSlots(
		[]ShipDesignSlot{
			{HullComponent: Laser.Name, HullSlotIndex: 2, Quantity: 1},
		},
	).WithSpec(&rules, player)
	player.Designs = append(player.Designs, starbaseDesignUpgrade)

	planet.ProductionQueue = append(planet.ProductionQueue, ProductionQueueItem{Type: QueueItemTypeStarbase, Quantity: 1, DesignNum: starbaseDesignUpgrade.Num})

	// generate a turn to upgrade the starbase
	turn.generateTurn()

	// should have an upgraded starbase at the planet, and the other should be
	// marked for deletion
	assert.Equal(t, "LASER BASE!", planet.Starbase.Tokens[0].design.Name)
	assert.Equal(t, 2, len(game.Starbases))
	assert.True(t, game.Starbases[0].Delete)
}

func Test_turn_fleetTransferOwner(t *testing.T) {
	game := createTwoPlayerGame()
	player1 := game.Players[0]
	player2 := game.Players[1]
	fleet := game.Fleets[0]

	player1.Race.Name = "Rabbitoid"
	player1.Race.PluralName = "Rabbitoids"

	turn := turn{
		game: game,
	}
	turn.game.Universe.buildMaps(game.Players)

	fleet.Waypoints[0].Task = WaypointTaskTransferFleet
	fleet.Waypoints[0].TransferToPlayer = player2.Num

	// transfer
	turn.generateTurn()

	// should have transferred the fleet, updated the name and the design
	assert.Equal(t, player2.Num, fleet.PlayerNum)
	assert.Equal(t, "Rabbitoids Long Range Scout #2", fleet.Name)
	assert.Equal(t, player1.Num, fleet.Tokens[0].design.OriginalPlayerNum)
	assert.Equal(t, "Rabbitoids Long Range Scout", fleet.Tokens[0].design.Name)
	assert.Equal(t, 2, len(player2.Designs))
	assert.Equal(t, 1, len(fleet.Waypoints))
	assert.Equal(t, None, fleet.Waypoints[0].TransferToPlayer)
	assert.Equal(t, WaypointTaskNone, string(fleet.Waypoints[0].Task))

}
