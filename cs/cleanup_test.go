//go:build !wasi && !wasm

package cs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_cleanup_FixRepeatOrders(t *testing.T) {
	a := NewPositionWaypoint(Vector{0, 0}, 5)
	b := NewPositionWaypoint(Vector{10, 0}, 5)
	c := NewPositionWaypoint(Vector{20, 0}, 5)
	partial := NewPositionWaypoint(Vector{5, 0}, 5)
	partial.PartiallyComplete = true

	positions := func(waypoints []Waypoint) []Vector {
		result := []Vector{}
		for _, wp := range waypoints {
			result = append(result, wp.Position)
		}
		return result
	}

	tests := []struct {
		name        string
		waypoints   []Waypoint
		want        []Waypoint
		wantChanged bool
	}{
		{"single waypoint", []Waypoint{a}, []Waypoint{a}, false},
		{"A to B loops back to A", []Waypoint{a, b}, []Waypoint{a, b, a}, true},
		{"A to B to C loops back to A", []Waypoint{a, b, c}, []Waypoint{a, b, c, a}, true},
		{"A to B and back to A already", []Waypoint{a, b, a}, []Waypoint{a, b, a}, false},
		{"partway to B already has the loop", []Waypoint{partial, b, a}, []Waypoint{partial, b, a}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, changed := NewCleaner().FixRepeatOrders(tt.waypoints)
			assert.Equal(t, tt.wantChanged, changed)
			assert.Equal(t, positions(tt.want), positions(got))
		})
	}
}
