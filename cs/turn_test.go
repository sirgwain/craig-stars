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
	player := client.NewPlayer(1, *NewRace(), &game.Rules).withSpec(&game.Rules)
	player.Num = 1
	player.Relations = []PlayerRelationship{{Relation: PlayerRelationFriend}} // friends with themselves
	player.PlayerIntels.PlayerIntels = player.defaultPlayerIntels([]*Player{player})

	planet := &Planet{
		MapObject: MapObject{Type: MapObjectTypePlanet, Name: "Planet 1", Num: 1, PlayerNum: player.Num},
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

	return &FullGame{
		Game:     game,
		Players:  players,
		Universe: &universe,
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

	// no victor
	assert.False(t, player.Victor)
	assert.False(t, game.VictorDeclared)

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

func Test_turn_fleetMoveDestroyedByMineField(t *testing.T) {
	game := createSingleUnitGame()
	rules := game.rules

	// change the rules so going 4 warp over the limit guarantee's a hit
	stats := game.rules.MineFieldStatsByType[MineFieldTypeStandard]
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
	assert.Equal(t, 1, len(game.Players[1].Messages))

	// the MineField should have lost some mines in the collision
	assert.Equal(t, 90, mineField.NumMines)
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
	rules := game.rules

	// change the rules so we don't decay
	stats := game.rules.MineFieldStatsByType[MineFieldTypeStandard]
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
	assert.Equal(t, 84, mineField.NumMines)
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
	player.Race.Spec = computeRaceSpec(&player.Race, &rules)

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
	planet.starbase = &starbase

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
				player.Relations = player.defaultRelationships(tt.args.players)
				player.PlayerIntels.PlayerIntels = player.defaultPlayerIntels(tt.args.players)
			}

			universe := NewUniverse(&game.Rules)
			universe.Fleets = []*Fleet{tt.args.fleet}
			universe.MineFields = []*MineField{tt.args.mineField}

			fg := FullGame{
				Game:     game,
				Players:  tt.args.players,
				Universe: &universe,
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
	packetTooFast := newMineralPacket(player, 1, 8, 5, Cargo{100, 100, 100, 0}, Vector{200, 0}, 1)
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
	assert.Equal(t, packetTooFast.Cargo, Cargo{50, 50, 50, 0})
	assert.Equal(t, packetNewlyBuilt.Cargo, Cargo{75, 75, 75, 0})
}

func Test_turn_fleetPatrol(t *testing.T) {
	game := createSingleUnitGame()
	rules := game.rules

	// make a new destroyer to patrol within 50ly
	player := game.Players[0]
	fleet := testStalwartDefender(player)
	fleet.Waypoints[0].Task = WaypointTaskPatrol
	fleet.Waypoints[0].PatrolRange = 50
	player.Designs[0] = fleet.Tokens[0].design
	game.Fleets[0] = fleet

	// create a new enemy player with a fleet 100ly away
	enemyPlayer := NewPlayer(2, NewRace().WithSpec(rules)).WithNum(2).withSpec(rules)
	enemyFleet := testLongRangeScout(enemyPlayer)
	enemyFleet.Position = Vector{60, 0}
	enemyPlayer.Designs = []*ShipDesign{enemyFleet.Tokens[0].design}
	game.Players = append(game.Players, enemyPlayer)
	game.Fleets = append(game.Fleets, enemyFleet)

	player.Relations = []PlayerRelationship{{Relation: PlayerRelationFriend}, {Relation: PlayerRelationNeutral}}
	enemyPlayer.Relations = []PlayerRelationship{{Relation: PlayerRelationNeutral}, {Relation: PlayerRelationFriend}}
	player.PlayerIntels.PlayerIntels = player.defaultPlayerIntels([]*Player{player, enemyPlayer})
	enemyPlayer.PlayerIntels.PlayerIntels = player.defaultPlayerIntels([]*Player{player, enemyPlayer})

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
	assert.Equal(t, fleet.Waypoints[1].TargetType, MapObjectTypeFleet)
	assert.Equal(t, fleet.Waypoints[1].TargetPlayerNum, enemyFleet.PlayerNum)
	assert.Equal(t, fleet.Waypoints[1].TargetNum, enemyFleet.Num)
	assert.Equal(t, fleet.Waypoints[1].WarpSpeed, 6)

}
