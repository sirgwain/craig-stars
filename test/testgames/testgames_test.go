//go:build !wasi && !wasm

package testgames

import (
	"github.com/sirgwain/craig-stars/cs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTestGames(t *testing.T) {
	names := map[string]bool{}
	for _, scenario := range TestGames {
		t.Run(scenario.Name, func(t *testing.T) {
			assert.False(t, names[scenario.Name], "fixture names must be unique for browser lookups")
			names[scenario.Name] = true
			game := cs.BuildScenario(scenario)
			assert.Equal(t, scenario.Name, game.Name)
			assert.Len(t, game.Players, len(scenario.Players))
			assert.Len(t, game.Planets, len(scenario.Planets))
			for _, planet := range game.Planets {
				assert.Equal(t, planet.Hab, planet.BaseHab)
				if planet.Starbase != nil {
					assert.True(t, planet.Spec.HasStarbase)
					assert.Equal(t, planet.Position, planet.Starbase.Position)
				}
			}
		})
	}
}
