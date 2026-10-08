package ai

import (
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/stretchr/testify/assert"
)

// A planet is a possible target when it is ahead of the fleet and close to its
// heading. Direction alone must not mark planets behind or beside it as threats.
func Test_aiPlayer_findPlanetTargets(t *testing.T) {
	tests := []struct {
		name    string
		heading cs.VectorFloat64
		target  cs.Vector
		want    bool
	}{
		// Axis-aligned approaches must match even when one heading component is zero.
		{"east", cs.VectorFloat64{X: 1}, cs.Vector{X: 20, Y: 10}, true},
		{"west", cs.VectorFloat64{X: -1}, cs.Vector{X: 0, Y: 10}, true},
		{"north", cs.VectorFloat64{Y: 1}, cs.Vector{X: 10, Y: 20}, true},
		{"south", cs.VectorFloat64{Y: -1}, cs.Vector{X: 10, Y: 0}, true},
		{"diagonal", cs.VectorFloat64{X: 1, Y: 1}, cs.Vector{X: 20, Y: 20}, true},
		{"opposite diagonal", cs.VectorFloat64{X: -1, Y: -1}, cs.Vector{X: 0, Y: 0}, true},
		// A planet behind the fleet, at its position, or off its path is not a target.
		{"behind", cs.VectorFloat64{X: 1}, cs.Vector{X: 0, Y: 10}, false},
		{"diagonal behind", cs.VectorFloat64{X: 1, Y: 1}, cs.Vector{X: 0, Y: 0}, false},
		{"at fleet", cs.VectorFloat64{X: 1}, cs.Vector{X: 10, Y: 10}, false},
		{"off horizontal path", cs.VectorFloat64{X: 1}, cs.Vector{X: 20, Y: 11}, false},
		{"off vertical path", cs.VectorFloat64{Y: 1}, cs.Vector{X: 11, Y: 20}, false},
		// Map coordinates are rounded: a miss within half a map unit still matches;
		// a larger miss does not. This tolerance must also work on diagonal paths.
		{"near diagonal path", cs.VectorFloat64{X: 2, Y: 1}, cs.Vector{X: 21, Y: 15}, true},
		{"off diagonal path", cs.VectorFloat64{X: 2, Y: 1}, cs.Vector{X: 20, Y: 16}, false},
		// Without a heading there is no direction from which to infer a target.
		{"zero heading", cs.VectorFloat64{}, cs.Vector{X: 20, Y: 20}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planet := cs.NewPlanet()
			planet.Position = tt.target
			ai := &aiPlayer{}
			targets := ai.findPlanetTargets(cs.Vector{X: 10, Y: 10}, tt.heading, []*cs.Planet{planet})
			if tt.want {
				assert.Equal(t, []*cs.Planet{planet}, targets)
			} else {
				assert.Empty(t, targets)
			}
		})
	}
}
