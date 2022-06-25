package game

import (
	"testing"
)

func Test_generatePlanets(t *testing.T) {
	type args struct {
		g    *Game
		area Vector
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Generate Planets for empty game", args{NewGame(), Vector{800, 800}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := generatePlanets(tt.args.g, tt.args.area); (err != nil) != tt.wantErr {
				t.Errorf("generatePlanets() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
