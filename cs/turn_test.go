package cs

import (
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
	player := client.NewPlayer(1, *NewRace(), &game.Rules).withSpec(&game.Rules)
	player.Num = 1
	player.Relations = []PlayerRelationship{{Relation: PlayerRelationFriend}} // friends with themselves

	planet := &Planet{
		MapObject: MapObject{Type: MapObjectTypePlanet, Name: "Planet 1", Num: 1, PlayerNum: player.Num},
		Cargo: Cargo{
			Colonists: 2500,
		},
	}
	planet.Spec = computePlanetSpec(&game.Rules, player, planet)

	fleet := testLongRangeScout(player)
	fleet.OrbitingPlanetNum = planet.Num
	fleet.Waypoints = []Waypoint{
		NewPlanetWaypoint(Vector{}, 1, "Planet 1", 5),
	}
	player.Designs = []*ShipDesign{
		fleet.Tokens[0].design,
	}

	players := []*Player{player}

	return &FullGame{
		Game:    game,
		Players: players,
		Universe: &Universe{
			Planets: []*Planet{planet},
			Fleets:  []*Fleet{fleet},
			rules:   &game.Rules,
		},
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
	planet.ProductionQueue = append([]ProductionQueueItem{{Type: QueueItemTypeShipToken, Quantity: 1, DesignName: player.GetDesign("Long Range Scout").Name}}, planet.ProductionQueue...)

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
	for i := 0; i < 100; i++ {
		client.GenerateTurn(game, universe, players)
	}

	assert.Equal(t, 2500, game.Year)

	// should have fleets
	assert.True(t, len(universe.Fleets) > 0)

	// should have grown pop
	assert.Greater(t, universe.Planets[0].population(), player.Race.Spec.StartingPlanets[0].Population)

	// should have built factories
	assert.Greater(t, universe.Planets[0].Factories, game.Rules.StartingFactories)

	// should have researched
	assert.NotEqual(t, player.TechLevels, TechLevel{3, 3, 3, 3, 3, 3})

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

	planet.TargetType = MapObjectTypePlanet
	planet.TargetNum = 2

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
