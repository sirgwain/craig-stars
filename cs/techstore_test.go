package cs

import (
	"reflect"
	"testing"
)

func TestTechStore_GetBestScanner(t *testing.T) {
	type args struct {
		player *Player
	}
	tests := []struct {
		name string
		args args
		want *TechHullComponent
	}{
		{"Get Lowest Scanner", args{newPlayer(1, NewRace())}, &BatScanner},
		{"Get Nicer Scanner", args{newPlayer(1, NewRace()).WithTechLevels(TechLevel{Electronics: 5})}, &PossumScanner},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := StaticTechStore.GetBestScanner(tt.args.player); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TechStore.GetBestScanner() = %v, want %v", got, tt.want)
			}
		})
	}
}
