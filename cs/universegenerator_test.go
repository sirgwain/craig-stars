package cs

import (
	"testing"

	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func Test_GenerateUniverse(t *testing.T) {
	client := NewGamer()
	game := client.CreateGame(1, *NewGameSettings())

	numPlanets, err := game.Rules.GetNumPlanets(game.Size, game.Density)
	if err != nil {
		t.Error(err)
	}
	player := client.NewPlayer(1, *NewRace(), &game.Rules)
	players := []*Player{player}
	player.AIControlled = true
	player.Num = 1
	universe, _ := client.GenerateUniverse(game, players)

	assert.Equal(t, len(universe.Planets), numPlanets)
	assert.Greater(t, len(universe.Fleets), 0)
	assert.Greater(t, len(player.Designs), 0)
	assert.Greater(t, len(universe.Wormholes), 0)

	pmo := universe.GetPlayerMapObjects(player.Num)

	assert.Equal(t, 1, len(pmo.Planets))
	homeworld := pmo.Planets[0]
	assert.Equal(t, 25_000, homeworld.population())
	assert.True(t, homeworld.Spec.HasStarbase)
}

func Test_getStartingStarbaseDesigns(t *testing.T) {
	player := NewPlayer(1, NewRace())
	player.Race.PRT = JoaT
	player.Race.Spec = computeRaceSpec(&player.Race, &rules)

	starbaseDesign := NewShipDesign(player, 1)
	starbaseDesign.Name = "Starbase"
	starbaseDesign.Hull = "Space Station"
	starbaseDesign.Purpose = ShipDesignPurposeStarbase
	starbaseDesign.Slots = []ShipDesignSlot{
		{Laser.Name, 2, 8},
		{MoleSkinShield.Name, 3, 8},
		{Laser.Name, 4, 8},
		{MoleSkinShield.Name, 6, 8},
		{Laser.Name, 8, 8},
		{Laser.Name, 10, 8},
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
			got := gu.getStartingStarbaseDesigns(tt.args.techStore, tt.args.player, 1)

			if !test.CompareAsJSON(t, got, tt.want) {
				t.Errorf("getStartingStarbaseDesigns() = %v, want %v", got, tt.want)
			}
		})
	}
}
