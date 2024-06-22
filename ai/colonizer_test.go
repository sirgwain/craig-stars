package ai

import (
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/test"
)

func Test_aiPlayer_getColonizerBuilderPlanets(t *testing.T) {

	tests := []struct {
		name             string
		planetsWithDocks []*cs.Planet
		want             []colonizerBuilderPlanet
	}{
		{"homeworld not ready",
			[]*cs.Planet{
				cs.NewPlanet().WithNum(1).WithCargo(cs.Cargo{Colonists: 350}),
			},
			nil,
		},
		{"homeworld at ~25%",
			[]*cs.Planet{
				cs.NewPlanet().WithNum(1).WithCargo(cs.Cargo{Colonists: 3100}),
			},
			[]colonizerBuilderPlanet{
				{
					Planet:             cs.NewPlanet().WithNum(1).WithCargo(cs.Cargo{Colonists: 3200}),
					availableColonists: 200,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ai := testAIPlayer()
			ai.planetsWithDocks = tt.planetsWithDocks

			if got := ai.getColonizerBuilderPlanets(); test.CompareAsJSON(t, got, tt.want) {
				t.Errorf("aiPlayer.getColonizerBuilderPlanets() = %v, want %v", got, tt.want)
			}
		})
	}
}
