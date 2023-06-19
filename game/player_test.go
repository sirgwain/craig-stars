package game

import (
	"testing"
)

func TestPlayer_HasTech(t *testing.T) {

	type args struct {
		player *Player
		tech   *Tech
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Everyone can learn the Viewer50", args{NewPlayer(1, NewRace()), &Viewer50.Tech}, true},
		{"Player is missing tech levels for the Viewer90", args{NewPlayer(1, NewRace()), &Viewer90.Tech}, false},
		{"Player is missing tech levels for the TransGalacticDrive", args{NewPlayer(1, NewRace()), &TransGalacticDrive.Tech}, false},
		{"Player has tech levels for the TransGalacticDrive", args{NewPlayer(1, NewRace()).WithTechLevels(TechLevel{Propulsion: 9}), &TransGalacticDrive.Tech}, true},
		{"Only players with IFE can learn the FuelMizer", args{NewPlayer(1, NewRace()), &FuelMizer.Tech}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.player.HasTech(tt.args.tech); got != tt.want {
				t.Errorf("Player.HasTech() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayer_CanLearnTech(t *testing.T) {

	type args struct {
		player *Player
		tech   *Tech
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Everyone can learn the Viewer50", args{NewPlayer(1, NewRace()), &Viewer50.Tech}, true},
		{"Everyone can learn the Viewer90", args{NewPlayer(1, NewRace()), &Viewer90.Tech}, true},
		{"Only players with IFE can learn the FuelMizer", args{NewPlayer(1, NewRace()), &FuelMizer.Tech}, false},
		{"Only players with IFE can learn the FuelMizer", args{NewPlayer(1, NewRace().WithLRT(IFE)), &FuelMizer.Tech}, true},
		{"Players with NRSE cannot learn the ramscoops", args{NewPlayer(1, NewRace().WithLRT(NRSE)), &FuelMizer.Tech}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.player.CanLearnTech(tt.args.tech); got != tt.want {
				t.Errorf("Player.CanLearnTech() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayer_getNextFleetNum(t *testing.T) {

	tests := []struct {
		name   string
		fleets []*Fleet
		want   int
	}{
		{"No Fleets", []*Fleet{}, 1},
		{"Simple fleet", []*Fleet{{MapObject: MapObject{Num: 1}}}, 2},
		{"Skipped num in fleets", []*Fleet{
			{MapObject: MapObject{Num: 1}},
			{MapObject: MapObject{Num: 3}},
		}, 2},
		{"Skipped num in fleets 2", []*Fleet{
			{MapObject: MapObject{Num: 1}},
			{MapObject: MapObject{Num: 2}},
			{MapObject: MapObject{Num: 3}},
			{MapObject: MapObject{Num: 5}},
		}, 4},
		{"Many fleets", []*Fleet{
			{MapObject: MapObject{Num: 1}},
			{MapObject: MapObject{Num: 2}},
			{MapObject: MapObject{Num: 3}},
			{MapObject: MapObject{Num: 4}},
		}, 5},
		{"Out of order", []*Fleet{
			{MapObject: MapObject{Num: 4}},
			{MapObject: MapObject{Num: 2}},
			{MapObject: MapObject{Num: 3}},
			{MapObject: MapObject{Num: 1}},
		}, 5},
		{"Out of order, missing number", []*Fleet{
			{MapObject: MapObject{Num: 4}},
			{MapObject: MapObject{Num: 2}},
			{MapObject: MapObject{Num: 1}},
		}, 3},
		{"Multiple fleet num 1 for starbases", []*Fleet{
			{MapObject: MapObject{Num: 1}},
			{MapObject: MapObject{Num: 1}},
			{MapObject: MapObject{Num: 1}},
			{MapObject: MapObject{Num: 2}},
		}, 3},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPlayer(1, NewRace())
			p.Fleets = tt.fleets
			if got := p.getNextFleetNum(); got != tt.want {
				t.Errorf("Player.getNextFleetNum() = %v, want %v", got, tt.want)
			}
		})
	}
}
