package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateUniverse(t *testing.T) {

	type args struct {
		g       *Game
		players []Player
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Generate Universe", args{NewGame(), []Player{*NewPlayer(1, NewRace())}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := tt.args.g
			for i := range tt.args.players {
				game.AddPlayer(&tt.args.players[i])
			}
			if err := game.GenerateUniverse(); (err != nil) != tt.wantErr {
				t.Errorf("GenerateUniverse() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.True(t, len(game.Planets) >= 24, "GenerateUniverse() did not generate planets")
			assert.True(t, len(game.GetOwnedPlanets()) > 0, "GenerateUniverse() did not create player planets")
			assert.True(t, len(game.Players[0].PlanetIntels) > 0, "GenerateUniverse() did not discover planets for player")
		})
	}
}
