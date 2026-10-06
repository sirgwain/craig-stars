//go:build !wasi && !wasm

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
		minefields     []*Minefield
		mineralPackets []*MineralPacket
	}
	tests := []struct {
		name string
		args args
		want []scanner
	}{
		{"Single Planet", args{planets: []*Planet{NewPlanet().WithPlayerNum(1).WithScanner(true)}}, []scanner{
			{Range: 150, RangePen: 0, CloakReductionFactor: 1},
		}},
		{"Single Long Range Scout", args{fleets: []*Fleet{testLongRangeScout(player).withPlayerNum(1)}}, []scanner{
			{Range: 66, RangePen: 30, CloakReductionFactor: 1},
		}},
		{"Planet and Scout same position", args{
			planets: []*Planet{NewPlanet().WithPlayerNum(1).WithScanner(true)},
			fleets:  []*Fleet{testLongRangeScout(player).withPlayerNum(1)},
		}, []scanner{
			{Range: 150, RangePen: 30, CloakReductionFactor: 1},
		}},
		{"Planet and Scout, diff position", args{
			planets: []*Planet{NewPlanet().WithPlayerNum(1).WithScanner(true)},
			fleets:  []*Fleet{testLongRangeScout(player).withPlayerNum(1).withPosition(Vector{1, 1})},
		}, []scanner{
			{Range: 150, RangePen: 0, CloakReductionFactor: 1},
			{Range: 66, RangePen: 30, Position: Vector{1, 1}, CloakReductionFactor: 1},
		}},
		{"Planet and two fleets, diff position", args{
			planets: []*Planet{NewPlanet().WithPlayerNum(1).WithScanner(true)},
			fleets: []*Fleet{
				testLongRangeScout(player).withPlayerNum(1).withPosition(Vector{1, 1}),
				testSmallFreighter(player).withPlayerNum(1).withPosition(Vector{1, 1}),
			},
		}, []scanner{
			{Range: 150, RangePen: 0, CloakReductionFactor: 1},
			{Range: 66, RangePen: 30, Position: Vector{1, 1}, CloakReductionFactor: 1},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// planet scanners come from the spec
			for _, planet := range tt.args.planets {
				planet.Spec = ComputePlanetSpec(&rules, player, planet)
			}
			scan := playerScanner{&Universe{
				Planets:        tt.args.planets,
				Fleets:         tt.args.fleets,
				MineralPackets: tt.args.mineralPackets,
				Minefields:     tt.args.minefields,
			}, &rules, player, []*Player{player}, make(map[int]bool), newDiscoverer(testLogger, player)}
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
		{
			"Single Planet, starbase with no gate",
			args{planet: NewPlanet(), player: NewPlayer(1, NewRace().WithSpec(&rules))},
			[]scanner{},
		},
		{
			"Single Planet, starbase with 100/250 gate, not IT",
			args{planet: NewPlanet(), player: NewPlayer(1, NewRace().WithSpec(&rules)), stargate: &Stargate100_250},
			[]scanner{},
		},
		{
			"Single Planet, starbase with 100/250 gate, IT",
			args{planet: NewPlanet(), player: NewPlayer(1, NewRace().WithPRT(IT).WithSpec(&rules)), stargate: &Stargate100_250},
			[]scanner{{RangePen: 250, CloakReductionFactor: 1}},
		},
		{
			"Single Planet, starbase with 100/any gate, IT",
			args{planet: NewPlanet(), player: NewPlayer(1, NewRace().WithPRT(IT).WithSpec(&rules)), stargate: &Stargate100_Any},
			[]scanner{{RangePen: math.MaxInt16, CloakReductionFactor: 1}},
		},
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

			planet.Spec = ComputePlanetSpec(&rules, player, planet)

			scan := playerScanner{&Universe{
				Planets: []*Planet{planet},
			}, &rules, player, []*Player{player}, make(map[int]bool), newDiscoverer(testLogger, player)}
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
		{
			"fleet at 0, 0 in scan range with 0 range scanner",
			args{
				player,
				testLongRangeScout(player).withPosition(Vector{0, 0}),
				scanner{Range: 0, RangePen: NoScanner, CloakReductionFactor: 1},
			},
			true,
		},
		{
			"fleet at 30, 0 in scan range with 30 range scanner",
			args{
				player,
				testLongRangeScout(player).withPosition(Vector{30, 0}),
				scanner{Range: 30, RangePen: NoScanner, CloakReductionFactor: 1},
			},
			true,
		},
		{
			"fleet at 31, 0 not in scan range with 30 range scanner",
			args{
				player,
				testLongRangeScout(player).withPosition(Vector{31, 0}),
				scanner{Range: 30, RangePen: NoScanner, CloakReductionFactor: 1},
			},
			false,
		},
		{
			"35% cloaked fleet at 66, 0 not in scan range with 100 range scanner",
			args{
				player,
				testCloakedScout(player).withPosition(Vector{66, 0}),
				scanner{Range: 100, RangePen: NoScanner, CloakReductionFactor: 1},
			},
			false,
		},
		{
			"35% cloaked fleet at 65, 0 in scan range with 100 range scanner",
			args{
				player,
				testCloakedScout(player).withPosition(Vector{65, 0}),
				scanner{Range: 100, RangePen: NoScanner, CloakReductionFactor: 1},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scan := playerScanner{player: tt.args.player}
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
				NewFleetWaypoint(Vector{10, 0}, 1, 2, "Target", warpSpeed),
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

		t.Run(tt.name, func(t *testing.T) {
			s := SingleUnitScenario()
			s.Players = append(s.Players, ScenarioPlayer{
				Designs: Designs(DesignLongRangeScout),
				Fleets:  []ScenarioFleet{{Name: "Target", Design: "Long Range Scout", Position: tt.args.targetPosition}},
			})
			u := newTestUniverse(t, s)
			game, player, fleet, enemyFleet := u.Game, u.Player(1), u.Fleet("Long Range Scout #1"), u.Fleet("Target")

			// target the enemty fleet
			fleet.Waypoints = []Waypoint{
				NewPositionWaypoint(tt.args.fleetPosition, 5),
				NewFleetWaypoint(enemyFleet.Position, enemyFleet.Num, enemyFleet.PlayerNum, enemyFleet.Name, 5),
			}

			enemyFleet.Delete = tt.args.targetDestroyed
			game.buildMaps(game.Players)
			player.clearTransientIntel()
			scan := newPlayerScanner(game.Universe, game.Players, &game.Rules, player)
			scan.scan()
			// check the waypoints returned vs what we want
			got := fleet.Waypoints
			test.CompareAsJSON(t, got, tt.want)
		})
	}
}

func Test_scanPlanetWithStargates(t *testing.T) {
	gatedStation := Designs(DesignSpaceStation)[0] // copy, so the preset is untouched
	gatedStation.Slots = append(gatedStation.Slots, ShipDesignSlot{HullComponent: Stargate100_250.Name, HullSlotIndex: 1, Quantity: 1})
	u := newTestUniverse(t, TestScenario{
		Players: []ScenarioPlayer{{Player: NewPlayer(1, NewRace().WithPRT(IT)), Designs: Designs(gatedStation)}, {Designs: Designs(gatedStation)}},
		Planets: []ScenarioPlanet{
			{Name: "Planet 1", Owner: 1, Starbase: gatedStation.Name, Cargo: Cargo{Colonists: 2500}},
			{Name: "Planet 2", Owner: 2, Starbase: gatedStation.Name, Position: Vector{500, 500}, Cargo: Cargo{Colonists: 2500}},
		},
	})
	game, player1, planet2 := u.Game, u.Player(1), u.Planet("Planet 2")
	starbase2 := planet2.Starbase
	scan := playerScanner{game.Universe, &game.Rules, player1, game.Players, make(map[int]bool), newDiscoverer(testLogger, player1)}

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
		intel     []*Wormhole
	}
	type args struct {
		scanners []scanner
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*Wormhole
	}{
		{
			name:   "scan wormhole",
			fields: fields{wormholes: []*Wormhole{newWormhole(Vector{}, 1, WormholeStabilityStable)}},
			args:   args{[]scanner{{Range: 10, CloakReductionFactor: 1}}},
			want:   []*Wormhole{{MapObject: MapObject{Type: MapObjectTypeWormhole, Num: 1}, Stability: WormholeStabilityStable}},
		},
		{
			name: "forget deleted wormhole",
			fields: fields{
				wormholes: []*Wormhole{{MapObject: MapObject{Num: 1, Type: MapObjectTypeWormhole, Delete: true}}},
				intel:     []*Wormhole{{MapObject: MapObject{Type: MapObjectTypeWormhole, Num: 1}, Stability: WormholeStabilityStable}},
			},
			args: args{[]scanner{{Range: 10, CloakReductionFactor: 1}}},
			want: []*Wormhole{},
		},
		{
			name: "forget wormhole we scanned again that no longer exists in universe",
			fields: fields{
				intel: []*Wormhole{{MapObject: MapObject{Type: MapObjectTypeWormhole, Num: 1}, Stability: WormholeStabilityStable}},
			},
			args: args{[]scanner{{Range: 10, CloakReductionFactor: 1}}},
			want: []*Wormhole{},
		},
		{
			name:   "wormhole 75%cloaked, out of range",
			fields: fields{wormholes: []*Wormhole{newWormhole(Vector{int(math.Ceil(50 * .75)), 0}, 1, WormholeStabilityStable)}},
			args:   args{[]scanner{{Range: 50, CloakReductionFactor: 1}}}, // 50ly scanner
			want:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scenario := TestScenario{Players: []ScenarioPlayer{{}}}
			for _, wormhole := range tt.fields.wormholes {
				scenario.Wormholes = append(scenario.Wormholes, *wormhole)
			}
			u := newTestUniverse(t, scenario)
			player, players, universe := u.Player(1), u.Game.Players, u.Game.Universe
			player.WormholeIntels = tt.fields.intel

			// make a new scanner
			discoverer := newDiscoverer(testLogger, player)
			scan := playerScanner{universe, &rules, player, players, make(map[int]bool, len(player.Intels.PlayerIntels)), discoverer}
			scan.scanWormholes(tt.args.scanners)

			// check the waypoints returned vs what we want
			got := player.WormholeIntels
			test.CompareAsJSON(t, got, tt.want)

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
			args: args{fleetPosition: Vector{}, scanner: scanner{Position: Vector{}, Range: 0, CloakReductionFactor: 1}},
			want: true,
		},
		{
			name: "too far away",
			args: args{fleetPosition: Vector{1, 0}, scanner: scanner{Position: Vector{}, Range: 0, CloakReductionFactor: 1}},
			want: false,
		},
		{
			name: "far but in scan range",
			args: args{fleetPosition: Vector{10, 0}, scanner: scanner{Position: Vector{}, Range: 10, CloakReductionFactor: 1}},
			want: true,
		},
		{
			name: "far but in pen scan range",
			args: args{fleetPosition: Vector{10, 0}, scanner: scanner{Position: Vector{}, RangePen: 10, CloakReductionFactor: 1}},
			want: true,
		},
		{
			name: "cloaked",
			args: args{
				fleetPosition: Vector{10, 0},
				fleetCloak:    50,
				scanner:       scanner{Position: Vector{}, RangePen: 10, CloakReductionFactor: 1},
			},
			want: false,
		},
		{
			name: "cloaked but in range",
			args: args{
				fleetPosition: Vector{5, 0},
				fleetCloak:    50,
				scanner:       scanner{Position: Vector{}, RangePen: 10, CloakReductionFactor: 1},
			},
			want: true,
		},
		{
			name: "cloaked with tachyon scanner",
			// 1 tachyon + 55% cloak is 52.25% effective cloaking, 47 dist is just in range
			args: args{
				fleetPosition: Vector{47, 0},
				fleetCloak:    55,
				scanner:       scanner{Position: Vector{}, RangePen: 100, CloakReductionFactor: math.Pow(.95, math.Sqrt(1))},
			},
			want: true,
		},
		{
			name: "cloaked with tachyon scanner, JUST out of range",
			// 1 tachyon + 55% cloak is 52.25% effective cloaking, 48 dist is just out of range
			args: args{
				fleetPosition: Vector{48, 0},
				fleetCloak:    55,
				scanner:       scanner{Position: Vector{}, RangePen: 100, CloakReductionFactor: math.Pow(.95, math.Sqrt(1))},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scan := playerScanner{}
			player := testPlayer()
			fleet := NewFleet(player, 1, "fleet", []Waypoint{NewPositionWaypoint(Vector{}, 0)})
			fleet.Spec.CloakPercent = tt.args.fleetCloak
			fleet.Position = tt.args.fleetPosition

			if got := scan.fleetInScannerRange(fleet, tt.args.scanner); got != tt.want {
				t.Errorf("playerScan.fleetInScannerRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_scanMinefields(t *testing.T) {
	type fields struct {
		minefields []*Minefield
		intel      []*Minefield
	}
	type args struct {
		scanners []scanner
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*Minefield
	}{
		{
			name:   "scan minefield, 1 mine, 1ly radius",
			fields: fields{minefields: []*Minefield{newMinefield(testPlayer().WithNum(2), MinefieldTypeStandard, 1, 1, Vector{})}},
			args:   args{[]scanner{{Range: 10, CloakReductionFactor: 1}}},
			want: []*Minefield{{
				MapObject: MapObject{
					Type:      MapObjectTypeMinefield,
					PlayerNum: 2, Num: 1, Name: "Humanoids Standard Minefield #1"},
				MinefieldType: MinefieldTypeStandard, NumMines: 1},
			},
		},
		{
			name: "minefield 75% cloaked, out of range",
			// minefield is 13 away, but has a 10ly radius so the edge is only 3 away
			// it is not spotted with 75% cloaking (scanner range is 2.5 instead of 10)
			fields: fields{minefields: []*Minefield{newMinefield(testPlayer().WithNum(2), MinefieldTypeStandard, 100, 1, Vector{13, 0})}},
			args:   args{[]scanner{{Range: 10, CloakReductionFactor: 1}}},
			want:   nil,
		},
		{
			name: "minefield 75% cloaked, edge is in range",
			// minefield is 12 away, but has a 10ly radius so the edge is only 2 away
			// it is spotted even with 75% cloaking (scanner range is 2.5 instead of 10)
			fields: fields{minefields: []*Minefield{newMinefield(testPlayer().WithNum(2), MinefieldTypeStandard, 100, 1, Vector{12, 0})}},
			args:   args{[]scanner{{Range: 10, CloakReductionFactor: 1}}},
			want: []*Minefield{{
				MapObject: MapObject{
					Type:      MapObjectTypeMinefield,
					Position:  Vector{12, 0},
					PlayerNum: 2, Num: 1, Name: "Humanoids Standard Minefield #1"},
				MinefieldType: MinefieldTypeStandard, NumMines: 100},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scenario := TestScenario{Players: []ScenarioPlayer{{}, {}}}
			for _, minefield := range tt.fields.minefields {
				scenario.Players[1].Minefields = append(scenario.Players[1].Minefields, *minefield)
			}
			u := newTestUniverse(t, scenario)
			player, players, universe := u.Player(1), u.Game.Players, u.Game.Universe
			player.MinefieldIntels = tt.fields.intel

			// make a new scanner
			discoverer := newDiscoverer(testLogger, player)
			scan := playerScanner{universe, &rules, player, players, make(map[int]bool, len(player.Intels.PlayerIntels)), discoverer}
			scan.scanMinefields(tt.args.scanners)

			// check the waypoints returned vs what we want
			got := player.MinefieldIntels
			test.CompareAsJSON(t, got, tt.want)

		})
	}
}
