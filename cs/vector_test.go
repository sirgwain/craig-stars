//go:build !wasi && !wasm

package cs

import (
	"testing"
)

func Test_segmentIntersectsCircle(t *testing.T) {
	tests := []struct {
		name          string
		segmentFrom   VectorFloat64
		segmentTo     VectorFloat64
		circleCenter  VectorFloat64
		circleRadius  float64
		expectedValue float64
	}{
		{
			name:          "segment is half out of circle",
			segmentFrom:   VectorFloat64{X: 0, Y: 0},
			segmentTo:     VectorFloat64{X: 1, Y: 0},
			circleCenter:  VectorFloat64{X: 0, Y: 0},
			circleRadius:  .5,
			expectedValue: 0.5,
		},
		{
			name:          "segment fully in circle",
			segmentFrom:   VectorFloat64{X: -.5, Y: 0},
			segmentTo:     VectorFloat64{X: .5, Y: 0},
			circleCenter:  VectorFloat64{X: 0, Y: 0},
			circleRadius:  .5,
			expectedValue: 0,
		},
		{
			name:          "segment is half in circle",
			segmentFrom:   VectorFloat64{X: -1, Y: 0},
			segmentTo:     VectorFloat64{X: 0, Y: 0},
			circleCenter:  VectorFloat64{X: 0, Y: 0},
			circleRadius:  0.5,
			expectedValue: 0.5,
		},
		{
			name:          "segment is 1/4 in circle",
			segmentFrom:   VectorFloat64{X: -1.25, Y: 0},
			segmentTo:     VectorFloat64{X: -.25, Y: 0},
			circleCenter:  VectorFloat64{X: 0, Y: 0},
			circleRadius:  0.5,
			expectedValue: 0.75,
		},
		{
			name:          "segment is 3/4 out of circle",
			segmentFrom:   VectorFloat64{X: .25, Y: 0},
			segmentTo:     VectorFloat64{X: 1.25, Y: 0},
			circleCenter:  VectorFloat64{X: 0, Y: 0},
			circleRadius:  0.5,
			expectedValue: 0.25,
		},
		{
			name:          "Segment does not intersect circle",
			segmentFrom:   VectorFloat64{X: 0, Y: 0},
			segmentTo:     VectorFloat64{X: 1, Y: 0},
			circleCenter:  VectorFloat64{X: 2, Y: 2},
			circleRadius:  .5,
			expectedValue: -1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := segmentIntersectsCircle(test.segmentFrom, test.segmentTo, test.circleCenter, test.circleRadius)
			if result != test.expectedValue {
				t.Errorf("Test case %q failed: expected %f, got %f", test.name, test.expectedValue, result)
			}
		})
	}
}
