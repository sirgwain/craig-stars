package cs

import (
	"math"
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func Test_getScanners(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).WithTechLevels(TechLevel{3, 3, 3, 3, 3, 3}).withSpec(&rules)
	player.Num = 1

	type args struct {
		planets        []*Planet
		fleets         []*Fleet
		mineFields     []*MineField
		mineralPackets []*MineralPacket
	}
	tests := []struct {
		name string
		args args
		want []scanner
	}{
		{"Single Planet", args{planets: []*Planet{NewPlanet().WithPlayerNum(1).WithScanner(true)}}, []scanner{
			{RangeSquared: 150 * 150, RangePenSquared: 0, CloakReductionFactor: 1},
		}},
		{"Single Long Range Scout", args{fleets: []*Fleet{testLongRangeScout(player).withPlayerNum(1)}}, []scanner{
			{RangeSquared: 66 * 66, RangePenSquared: 30 * 30, CloakReductionFactor: 1},
		}},
		{"Planet and Scout same position", args{planets: []*Planet{NewPlanet().WithPlayerNum(1).WithScanner(true)}, fleets: []*Fleet{testLongRangeScout(player).withPlayerNum(1)}}, []scanner{
			{RangeSquared: 150 * 150, RangePenSquared: 30 * 30, CloakReductionFactor: 1},
		}},
		{"Planet and Scout, diff position", args{planets: []*Planet{NewPlanet().WithPlayerNum(1).WithScanner(true)}, fleets: []*Fleet{testLongRangeScout(player).withPlayerNum(1).withPosition(Vector{1, 1})}}, []scanner{
			{RangeSquared: 150 * 150, RangePenSquared: 0, CloakReductionFactor: 1},
			{RangeSquared: 66 * 66, RangePenSquared: 30 * 30, Position: Vector{1, 1}, CloakReductionFactor: 1},
		}},
		{"Planet and two fleets, diff position", args{planets: []*Planet{NewPlanet().WithPlayerNum(1).WithScanner(true)}, fleets: []*Fleet{testLongRangeScout(player).withPlayerNum(1).withPosition(Vector{1, 1}), testSmallFreighter(player).withPlayerNum(1).withPosition(Vector{1, 1})}}, []scanner{
			{RangeSquared: 150 * 150, RangePenSquared: 0, CloakReductionFactor: 1},
			{RangeSquared: 66 * 66, RangePenSquared: 30 * 30, Position: Vector{1, 1}, CloakReductionFactor: 1},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// planet scanners come from the spec
			for _, planet := range tt.args.planets {
				planet.Spec = computePlanetSpec(&rules, player, planet)
			}
			scan := playerScan{&Universe{
				Planets:        tt.args.planets,
				Fleets:         tt.args.fleets,
				MineralPackets: tt.args.mineralPackets,
				MineFields:     tt.args.mineFields,
			}, &rules, player, []*Player{player}, make(map[int]bool), newDiscoverer(player)}
			if got := scan.getScanners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getScanners() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func Test_getStargateScanners(t *testing.T) {

	// get stargate scanner for a single planet/player with stargate
	type args struct {
		planet   *Planet
		player   *Player
		stargate *TechHullComponent
	}
	tests := []struct {
		name string
		args args
		want []scanner
	}{
		{"Single Planet, starbase with no gate", args{planet: NewPlanet(), player: NewPlayer(1, NewRace().WithSpec(&rules))}, []scanner{}},
		{"Single Planet, starbase with 100/250 gate, not IT", args{planet: NewPlanet(), player: NewPlayer(1, NewRace().WithSpec(&rules)), stargate: &Stargate100_250}, []scanner{}},
		{"Single Planet, starbase with 100/250 gate, IT", args{planet: NewPlanet(), player: NewPlayer(1, NewRace().WithPRT(IT).WithSpec(&rules)), stargate: &Stargate100_250}, []scanner{{RangePenSquared: 250 * 250, CloakReductionFactor: 1}}},
		{"Single Planet, starbase with 100/any gate, IT", args{planet: NewPlanet(), player: NewPlayer(1, NewRace().WithPRT(IT).WithSpec(&rules)), stargate: &Stargate100_Any}, []scanner{{RangePenSquared: math.MaxInt16 * math.MaxInt16, CloakReductionFactor: 1}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// planet scanners come from the spec
			planet := tt.args.planet
			player := tt.args.player
			player.Num = 1
			planet.PlayerNum = player.Num

			starbase := testSpaceStation(player, planet)
			if tt.args.stargate != nil {
				design := starbase.Tokens[0].design
				design.Slots = append(design.Slots, ShipDesignSlot{HullComponent: tt.args.stargate.Name, HullSlotIndex: 1, Quantity: 1})
				design.Spec, _ = ComputeShipDesignSpec(&rules, player.TechLevels, player.Race.Spec, design)
				starbase.Spec = ComputeFleetSpec(&rules, player, starbase)
			}
			planet.Starbase = starbase

			planet.Spec = computePlanetSpec(&rules, player, planet)

			scan := playerScan{&Universe{
				Planets: []*Planet{planet},
			}, &rules, player, []*Player{player}, make(map[int]bool), newDiscoverer(player)}
			if got := scan.getStarGateScanners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getScanners() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func Test_fleetInScannerRange(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).WithTechLevels(TechLevel{3, 3, 3, 3, 3, 3}).withSpec(&rules)

	type args struct {
		player  *Player
		fleet   *Fleet
		scanner scanner
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"fleet at 0, 0 in scan range with 0 range scanner", args{player, testLongRangeScout(player).withPosition(Vector{0, 0}), scanner{RangeSquared: 0, RangePenSquared: NoScanner}}, true},
		{"fleet at 30, 0 in scan range with 30 range scanner", args{player, testLongRangeScout(player).withPosition(Vector{30, 0}), scanner{RangeSquared: 30 * 30, RangePenSquared: NoScanner}}, true},
		{"fleet at 31, 0 not in scan range with 30 range scanner", args{player, testLongRangeScout(player).withPosition(Vector{31, 0}), scanner{RangeSquared: 30 * 30, RangePenSquared: NoScanner}}, false},
		{"cloaked fleet at 30, 0 not in scan range with 30 range scanner", args{player, testCloakedScout(player).withPosition(Vector{30, 0}), scanner{RangeSquared: 30 * 30, RangePenSquared: NoScanner}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scan := playerScan{player: tt.args.player}
			if got := scan.fleetInScannerRange(tt.args.fleet, tt.args.scanner); got != tt.want {
				t.Errorf("fleetInScannerRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_updateFleetTargets(t *testing.T) {
	warpSpeed := 5
	type args struct {
		fleetPosition   Vector
		targetPosition  Vector
		targetDestroyed bool
	}
	tests := []struct {
		name string
		args args
		want []Waypoint
	}{
		{
			name: "fleet still in range, no change",
			args: args{fleetPosition: Vector{0, 0}, targetPosition: Vector{10, 0}},
			want: []Waypoint{
				NewPositionWaypoint(Vector{}, warpSpeed),
				NewFleetWaypoint(Vector{10, 0}, 1, 2, "", warpSpeed),
			},
		},
		{
			name: "fleet out of range, make position waypoint",
			args: args{fleetPosition: Vector{0, 0}, targetPosition: Vector{1000, 0}},
			want: []Waypoint{
				NewPositionWaypoint(Vector{}, warpSpeed),
				NewPositionWaypoint(Vector{1000, 0}, warpSpeed),
			},
		},
		{
			name: "fleet destroyed, but was at our location",
			args: args{fleetPosition: Vector{0, 0}, targetPosition: Vector{0, 0}, targetDestroyed: true},
			want: []Waypoint{
				NewPositionWaypoint(Vector{}, warpSpeed),
			},
		},
	}
	for _, tt := range tests {

		game := createSingleUnitGame()
		player := game.Players[0]
		fleet := game.Fleets[0]

		// create a new enemy player with a fleet at targetPosition
		enemyPlayer := NewPlayer(2, NewRace().WithSpec(&rules)).WithNum(2).withSpec(&rules)
		enemyFleet := testLongRangeScout(enemyPlayer)
		enemyFleet.Position = tt.args.targetPosition
		game.Players = append(game.Players, enemyPlayer)
		game.Fleets = append(game.Fleets, enemyFleet)

		player.Relations = []PlayerRelationship{{Relation: PlayerRelationFriend}, {Relation: PlayerRelationNeutral}}
		enemyPlayer.Relations = []PlayerRelationship{{Relation: PlayerRelationNeutral}, {Relation: PlayerRelationFriend}}
		player.PlayerIntels.PlayerIntels = player.defaultPlayerIntels([]*Player{player, enemyPlayer})
		enemyPlayer.PlayerIntels.PlayerIntels = player.defaultPlayerIntels([]*Player{player, enemyPlayer})
		// setup initial planet intels so turn generation works
		enemyPlayer.initDefaultPlanetIntels(game.Planets)

		// target the enemty fleet
		fleet.Waypoints = []Waypoint{
			NewPositionWaypoint(tt.args.fleetPosition, 5),
			NewFleetWaypoint(enemyFleet.Position, enemyFleet.Num, enemyFleet.PlayerNum, enemyFleet.Name, 5),
		}

		t.Run(tt.name, func(t *testing.T) {
			enemyFleet.Delete = tt.args.targetDestroyed
			scan := newPlayerScanner(game.Universe, game.Players, &game.Rules, player)
			scan.scan()
			// check the waypoints returned vs what we want
			if got := fleet.Waypoints; !test.CompareAsJSON(t, got, tt.want) {
				t.Errorf("updateFleetTargets() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func Test_scanPlanetWithStargates(t *testing.T) {
	game := createTwoPlayerGame()
	// setup a player1 and a  planet with a starbase with a scanner
	player1 := game.Players[0]
	planet1 := game.Planets[0]

	// make player1 an IT
	player1.Race.PRT = IT
	player1.Race.Spec = computeRaceSpec(&player1.Race, &rules)

	// add a stargate to the space station
	starbase1 := testSpaceStation(player1, planet1)
	design1 := starbase1.Tokens[0].design
	design1.Slots = append(design1.Slots, ShipDesignSlot{HullComponent: Stargate100_250.Name, HullSlotIndex: 1, Quantity: 1})
	design1.Spec, _ = ComputeShipDesignSpec(&rules, player1.TechLevels, player1.Race.Spec, design1)
	starbase1.Spec = ComputeFleetSpec(&rules, player1, starbase1)

	planet1.Starbase = starbase1
	planet1.Spec = computePlanetSpec(&rules, player1, planet1)

	// create a second player/planet
	player2 := game.Players[1]
	planet2 := game.Planets[1]
	// add a stargate to the second space station
	starbase2 := testSpaceStation(player2, planet2)
	design2 := starbase2.Tokens[0].design
	design2.Slots = append(design2.Slots, ShipDesignSlot{HullComponent: Stargate100_250.Name, HullSlotIndex: 1, Quantity: 1})
	design2.Spec, _ = ComputeShipDesignSpec(&rules, player2.TechLevels, player2.Race.Spec, design2)
	starbase2.Spec = ComputeFleetSpec(&rules, player2, starbase1)

	planet2.Starbase = starbase1
	planet2.Spec = computePlanetSpec(&rules, player2, planet2)

	scan := playerScan{game.Universe, &rules, player1, game.Players, make(map[int]bool), newDiscoverer(player1)}

	// first test a faraway planet
	planet2.Position = Vector{500, 500}
	starbase2.Position = planet2.Position
	scan.scanPlanets([]scanner{}, []scanner{}, scan.getStarGateScanners())

	// player1's starbase should scan player2's starbase
	assert.Equal(t, ReportAgeUnexplored, player1.PlanetIntels[1].ReportAge)

	// now test a close up planet
	planet2.Position = Vector{250, 0}
	starbase2.Position = planet2.Position
	scan.scanPlanets([]scanner{}, []scanner{}, scan.getStarGateScanners())

	// player1's starbase should scan player2's starbase and therefore planet
	assert.Equal(t, planet2.Hab, player1.PlanetIntels[1].Hab)
	assert.Equal(t, planet2.MineralConcentration, player1.PlanetIntels[1].MineralConcentration)

}

func Test_scanWormholes(t *testing.T) {
	type fields struct {
		wormholes []*Wormhole
		intel     []WormholeIntel
	}
	type args struct {
		scanners []scanner
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []WormholeIntel
	}{
		{
			name:   "scan wormhole",
			fields: fields{wormholes: []*Wormhole{newWormhole(Vector{}, 1, WormholeStabilityStable)}},
			args:   args{[]scanner{{RangeSquared: 100}}},
			want:   []WormholeIntel{{MapObjectIntel: MapObjectIntel{Type: MapObjectTypeWormhole, Intel: Intel{Num: 1}}, Stability: WormholeStabilityStable}},
		},
		{
			name: "forget deleted wormhole",
			fields: fields{
				wormholes: []*Wormhole{{MapObject: MapObject{Num: 1, Type: MapObjectTypeWormhole, Delete: true}}},
				intel:     []WormholeIntel{{MapObjectIntel: MapObjectIntel{Type: MapObjectTypeWormhole, Intel: Intel{Num: 1}}, Stability: WormholeStabilityStable}},
			},
			args: args{[]scanner{{RangeSquared: 100}}},
			want: []WormholeIntel{},
		},
		{
			name: "forget wormhole we scanned again that no longer exists in universe",
			fields: fields{
				intel: []WormholeIntel{{MapObjectIntel: MapObjectIntel{Type: MapObjectTypeWormhole, Intel: Intel{Num: 1}}, Stability: WormholeStabilityStable}},
			},
			args: args{[]scanner{{RangeSquared: 100}}},
			want: []WormholeIntel{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(1, NewRace().WithSpec(&rules))
			player.WormholeIntels = tt.fields.intel

			players := []*Player{player}

			universe := NewUniverse(&rules)
			universe.Wormholes = tt.fields.wormholes
			universe.buildMaps(players)

			// make a new scanner
			discoverer := newDiscoverer(player)
			scan := playerScan{&universe, &rules, player, players, make(map[int]bool, len(player.PlayerIntels.PlayerIntels)), discoverer}
			scan.scanWormholes(tt.args.scanners)

			// check the waypoints returned vs what we want
			if got := player.WormholeIntels; !test.CompareAsJSON(t, got, tt.want) {
				t.Errorf("scanWormholes() = \n%v, want \n%v", got, tt.want)
			}

		})
	}
}

func Test_playerScan_fleetInScannerRange(t *testing.T) {
	type args struct {
		fleetCloak    int
		fleetPosition Vector
		scanner       scanner
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "on top of each other",
			args: args{fleetPosition: Vector{}, scanner: scanner{Position: Vector{}, RangeSquared: 0, CloakReductionFactor: 1}},
			want: true,
		},
		{
			name: "too far away",
			args: args{fleetPosition: Vector{1, 0}, scanner: scanner{Position: Vector{}, RangeSquared: 0, CloakReductionFactor: 1}},
			want: false,
		},
		{
			name: "far but in scan range",
			args: args{fleetPosition: Vector{10, 0}, scanner: scanner{Position: Vector{}, RangeSquared: 100, CloakReductionFactor: 1}},
			want: true,
		},
		{
			name: "far but in pen scan range",
			args: args{fleetPosition: Vector{10, 0}, scanner: scanner{Position: Vector{}, RangePenSquared: 100, CloakReductionFactor: 1}},
			want: true,
		},
		{
			name: "cloaked",
			args: args{fleetPosition: Vector{10, 0}, fleetCloak: 50, scanner: scanner{Position: Vector{}, RangePenSquared: 100, CloakReductionFactor: 1}},
			want: false,
		},
		{
			name: "cloaked but in range",
			args: args{fleetPosition: Vector{5, 0}, fleetCloak: 50, scanner: scanner{Position: Vector{}, RangePenSquared: 100, CloakReductionFactor: 1}},
			want: true,
		},
		{
			name: "cloaked with tachyon scanner",
			// 1 tachyon + 55% cloak is 52.25% effective cloaking
			args: args{fleetPosition: Vector{math.Sqrt(52.2), 0}, fleetCloak: 55, scanner: scanner{Position: Vector{}, RangePenSquared: 100, CloakReductionFactor: math.Pow(.95, math.Sqrt(1))}},
			want: true,
		},
		{
			name: "cloaked with tachyon scanner, JUST out of range",
			// 1 tachyon + 55% cloak is 52.25% effective cloaking
			args: args{fleetPosition: Vector{math.Sqrt(52.3), 0}, fleetCloak: 55, scanner: scanner{Position: Vector{}, RangePenSquared: 100, CloakReductionFactor: math.Pow(.95, math.Sqrt(1))}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scan := playerScan{}
			player := testPlayer()
			fleet := newFleet(player, 1, "fleet", []Waypoint{NewPositionWaypoint(Vector{}, 0)})
			fleet.Spec.CloakPercent = tt.args.fleetCloak
			fleet.Position = tt.args.fleetPosition

			if got := scan.fleetInScannerRange(&fleet, tt.args.scanner); got != tt.want {
				t.Errorf("playerScan.fleetInScannerRange() = %v, want %v", got, tt.want)
			}
		})
	}
}
