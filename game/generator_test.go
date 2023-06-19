package game

import (
	"testing"

	"github.com/sirgwain/craig-stars/test"
)

func Test_getStartingStarbaseDesigns(t *testing.T) {
	player := newPlayer(1, NewRace())
	player.Race.PRT = JoaT
	rules := NewRules()
	player.Race.Spec = computeRaceSpec(&player.Race, &rules)

	starbaseDesign := NewShipDesign(player)
	starbaseDesign.Name = "Starbase"
	starbaseDesign.Hull = "Space Station"
	starbaseDesign.Purpose = ShipDesignPurposeStarbase
	starbaseDesign.Slots = []ShipDesignSlot{
		{"Laser", 2, 8},
		{"Mole-skin Shield", 3, 8},
		{"Laser", 4, 8},
		{"Mole-skin Shield", 6, 8},
		{"Laser", 8, 8},
		{"Laser", 10, 8},
	}

	type args struct {
		techStore *TechStore
		player    *Player
	}
	tests := []struct {
		name string
		args args
		want []ShipDesign
	}{
		{"Humanoid Designs", args{&StaticTechStore, player}, []ShipDesign{*starbaseDesign}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gu := universeGenerator{}
			got := gu.getStartingStarbaseDesigns(tt.args.techStore, tt.args.player)

			// uuids are random, so just make our want/got's the same
			for i := range got {
				tt.want[i].UUID = got[i].UUID
			}
			if !test.CompareAsJSON(t, got, tt.want) {
				t.Errorf("getStartingStarbaseDesigns() = %v, want %v", got, tt.want)
			}
		})
	}
}
