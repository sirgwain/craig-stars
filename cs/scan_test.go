package cs

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/test"
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
			{RangeSquared: 150 * 150, RangePenSquared: 0},
		}},
		{"Single Long Range Scout", args{fleets: []*Fleet{testLongRangeScout(player).withPlayerNum(1)}}, []scanner{
			{RangeSquared: 66 * 66, RangePenSquared: 30 * 30},
		}},
		{"Planet and Scout same position", args{planets: []*Planet{NewPlanet().WithPlayerNum(1).WithScanner(true)}, fleets: []*Fleet{testLongRangeScout(player).withPlayerNum(1)}}, []scanner{
			{RangeSquared: 150 * 150, RangePenSquared: 30 * 30},
		}},
		{"Planet and Scout, diff position", args{planets: []*Planet{NewPlanet().WithPlayerNum(1).WithScanner(true)}, fleets: []*Fleet{testLongRangeScout(player).withPlayerNum(1).withPosition(Vector{1, 1})}}, []scanner{
			{RangeSquared: 150 * 150, RangePenSquared: 0},
			{RangeSquared: 66 * 66, RangePenSquared: 30 * 30, Position: Vector{1, 1}},
		}},
		{"Planet and two fleets, diff position", args{planets: []*Planet{NewPlanet().WithPlayerNum(1).WithScanner(true)}, fleets: []*Fleet{testLongRangeScout(player).withPlayerNum(1).withPosition(Vector{1, 1}), testSmallFreighter(player).withPlayerNum(1).withPosition(Vector{1, 1})}}, []scanner{
			{RangeSquared: 150 * 150, RangePenSquared: 0},
			{RangeSquared: 66 * 66, RangePenSquared: 30 * 30, Position: Vector{1, 1}},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scan := playerScan{&Universe{
				Planets:        tt.args.planets,
				Fleets:         tt.args.fleets,
				MineralPackets: tt.args.mineralPackets,
				MineFields:     tt.args.mineFields,
				rules:          &rules,
			}, &rules, player, []*Player{player}, make(map[int]bool), newDiscoverer(player)}
			if got := scan.getScanners(); !reflect.DeepEqual(got, tt.want) {
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
		enemyPlayer.initDefaultPlanetIntels(&game.Rules, game.Planets)

		// target the enemty fleet
		fleet.Waypoints = []Waypoint{
			NewPositionWaypoint(tt.args.fleetPosition, 5),
			NewFleetWaypoint(enemyFleet.Position, enemyFleet.Num, enemyFleet.PlayerNum, enemyFleet.Name, 5),
		}

		t.Run(tt.name, func(t *testing.T) {
			enemyFleet.Delete = tt.args.targetDestroyed
			scan := newPlayerScanner(game.Universe, game.Players, game.rules, player)
			scan.scan()
			// check the waypoints returned vs what we want
			if got := fleet.Waypoints; !test.CompareAsJSON(t, got, tt.want) {
				t.Errorf("updateFleetTargets() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}
