package game

import (
	"reflect"
	"testing"
)

func Test_getScanners(t *testing.T) {
	rules := NewRules()
	player := NewPlayer(1, NewRace().WithSpec(&rules)).WithTechLevels(TechLevel{3, 3, 3, 3, 3, 3}).WithSpec(&rules)

	type args struct {
		planets          []*Planet
		fleets           []*Fleet
		mineFields       []*MineField
		mineralPackets   []*MineralPacket
		planetaryScanner TechPlanetaryScanner
	}
	tests := []struct {
		name string
		args args
		want []scanner
	}{
		{"Single Planet", args{planets: []*Planet{NewPlanet().WithScanner(true)}, planetaryScanner: Scoper150}, []scanner{
			{RangeSquared: 150 * 150, RangePenSquared: 0},
		}},
		{"Single Long Range Scout", args{fleets: []*Fleet{testLongRangeScout(player, &rules)}, planetaryScanner: Scoper150}, []scanner{
			{RangeSquared: 66 * 66, RangePenSquared: 30 * 30},
		}},
		{"Planet and Scout same position", args{planets: []*Planet{NewPlanet().WithScanner(true)}, fleets: []*Fleet{testLongRangeScout(player, &rules)}, planetaryScanner: Scoper150}, []scanner{
			{RangeSquared: 150 * 150, RangePenSquared: 30 * 30},
		}},
		{"Planet and Scout, diff position", args{planets: []*Planet{NewPlanet().WithScanner(true)}, fleets: []*Fleet{testLongRangeScout(player, &rules).WithPosition(Vector{1, 1})}, planetaryScanner: Scoper150}, []scanner{
			{RangeSquared: 150 * 150, RangePenSquared: 0},
			{RangeSquared: 66 * 66, RangePenSquared: 30 * 30, Position: Vector{1, 1}},
		}},
		{"Planet and two fleets, diff position", args{planets: []*Planet{NewPlanet().WithScanner(true)}, fleets: []*Fleet{testLongRangeScout(player, &rules).WithPosition(Vector{1, 1}), testSmallFreighter(player, &rules).WithPosition(Vector{1, 1})}, planetaryScanner: Scoper150}, []scanner{
			{RangeSquared: 150 * 150, RangePenSquared: 0},
			{RangeSquared: 66 * 66, RangePenSquared: 30 * 30, Position: Vector{1, 1}},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getScanners(tt.args.planets, tt.args.fleets, tt.args.mineralPackets, tt.args.mineFields, tt.args.planetaryScanner); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getScanners() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func Test_fleetInScannerRange(t *testing.T) {
	rules := NewRules()
	player := NewPlayer(1, NewRace().WithSpec(&rules)).WithTechLevels(TechLevel{3, 3, 3, 3, 3, 3}).WithSpec(&rules)

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
		{"fleet at 0, 0 in scan range with 0 range scanner", args{player, testLongRangeScout(player, &rules).WithPosition(Vector{0, 0}), scanner{RangeSquared: 0, RangePenSquared: NoScanner}}, true},
		{"fleet at 30, 0 in scan range with 30 range scanner", args{player, testLongRangeScout(player, &rules).WithPosition(Vector{30, 0}), scanner{RangeSquared: 30 * 30, RangePenSquared: NoScanner}}, true},
		{"fleet at 31, 0 not in scan range with 30 range scanner", args{player, testLongRangeScout(player, &rules).WithPosition(Vector{31, 0}), scanner{RangeSquared: 30 * 30, RangePenSquared: NoScanner}}, false},
		{"cloaked fleet at 30, 0 not in scan range with 30 range scanner", args{player, testCloakedScout(player, &rules).WithPosition(Vector{30, 0}), scanner{RangeSquared: 30 * 30, RangePenSquared: NoScanner}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fleetInScannerRange(tt.args.player, tt.args.fleet, tt.args.scanner); got != tt.want {
				t.Errorf("fleetInScannerRange() = %v, want %v", got, tt.want)
			}
		})
	}
}