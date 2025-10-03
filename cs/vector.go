package cs

import (
	"fmt"
	"math"
)

// A simple 2D vector with handy functions for moving things in space, calculating distance, etc
// Many of these functions were taken from Godot source, thanks Godot folks.
type VectorGeneric[T number] struct {
	X T `json:"x"`
	Y T `json:"y"`
}

type Vector = VectorGeneric[int]
type VectorFloat64 = VectorGeneric[float64]

func (v VectorGeneric[T]) String() string {
	return fmt.Sprintf("(%v, %v)", v.X, v.Y)
}

// compare int vectors by their x coord
func VectorCompareX(a, b Vector) int {
	if a.X == b.X {
		return a.Y - b.Y
	}
	return a.X - b.X
}

func (v VectorGeneric[T]) ToFloat64() VectorFloat64 {
	return VectorFloat64{X: float64(v.X), Y: float64(v.Y)}
}

func (v VectorGeneric[T]) ToInt(round bool) Vector {
	if round {
		return Vector{X: int(math.Round(float64(v.X))), Y: int(math.Round(float64(v.Y)))}
	}
	return Vector{X: int(v.X), Y: int(v.Y)}
}

func (v VectorGeneric[T]) DistanceSquaredTo(to VectorGeneric[T]) T {
	dx := v.X - to.X
	dy := v.Y - to.Y
	return dx*dx + dy*dy
}

func (v VectorGeneric[T]) DistanceTo(to VectorGeneric[T]) float64 {
	return math.Sqrt(float64(v.DistanceSquaredTo(to)))
}

func (minuend VectorGeneric[T]) Subtract(subtrahend VectorGeneric[T]) VectorGeneric[T] {
	return VectorGeneric[T]{minuend.X - subtrahend.X, minuend.Y - subtrahend.Y}
}

func (v1 VectorGeneric[T]) Add(v2 VectorGeneric[T]) VectorGeneric[T] {
	return VectorGeneric[T]{v1.X + v2.X, v1.Y + v2.Y}
}

func (v VectorGeneric[T]) Scale(scale T) VectorGeneric[T] {
	return VectorGeneric[T]{v.X * scale, v.Y * scale}
}

func (v VectorGeneric[T]) LengthSquared() T {
	return (v.X * v.X) + (v.Y * v.Y)
}

func (v VectorGeneric[T]) Length() float64 {
	return math.Sqrt(float64((v.X * v.X) + (v.Y * v.Y)))
}

func (v VectorGeneric[T]) Normalized() VectorFloat64 {
	vf := v.ToFloat64()
	lengthsq := vf.ToFloat64().LengthSquared()

	if lengthsq == 0 {
		vf.X = 0
		vf.Y = 0
	} else {
		length := math.Sqrt(float64(lengthsq))
		vf.X = vf.X / length
		vf.Y = vf.Y / length
	}
	return vf
}

func (v VectorGeneric[T]) Dot(other VectorGeneric[T]) T {
	return v.X*other.X + v.Y*other.Y
}

func (v VectorGeneric[T]) Round() VectorGeneric[T] {
	return VectorGeneric[T]{T(math.Round(float64(v.X))), T(math.Round(float64(v.Y)))}
}

// SegmentIntersectsCircle checks whether a segment intersects a circle or not.
// This returns what percent of the segment is NOT in the circle, or -1 if it doesn't
// intersect.
// Who would have thought the godot developers would write a perfect function for determining
// if we collide with a minefield
// while moving through space?
// https://github.com/godotengine/godot/blob/4.1.2-stable/core/math/geometry_2d.h#L217
func segmentIntersectsCircle[T number](segmentFrom, segmentTo, circlePosition VectorGeneric[T], circleRadius float64) (percentOutside float64) {
	lineVec := segmentTo.Subtract(segmentFrom)
	vecToLine := segmentFrom.Subtract(circlePosition)

	// Create a quadratic formula of the form ax^2 + bx + c = 0
	// Hope you remembered your high school algebra!
	var a, b, c float64

	a = float64(lineVec.Dot(lineVec))
	b = 2 * float64(vecToLine.Dot(lineVec))
	c = float64(vecToLine.Dot(vecToLine)) - circleRadius*circleRadius

	// Calculate the discriminant - b^2 - 4ac
	var discriminant = b*b - 4*a*c

	// A discriminant below 0 implies a non-real value,
	// so it definitely won't be in the range of 0 to 1.
	if discriminant < 0 {
		return -1
	}

	// If we can assume that the line segment starts outside the circle
	// (e.g. for continuous time collision detection), the following can be
	// skipped and we can just return the equivalent of res1.
	discriminant = math.Sqrt(discriminant)
	root1 := (-b - discriminant) / (2 * a)
	root2 := (-b + discriminant) / (2 * a)

	if root1 >= 0 && root1 <= 1 {
		return root1
	}
	if root2 >= 0 && root2 <= 1 {
		return root2
	}
	return -1
}

// Returns true if this point is in a circle
func isPointInCircle[T number](point, circlePosition VectorGeneric[T], circleRadius float64) bool {
	return float64(point.DistanceSquaredTo(circlePosition)) <= circleRadius*circleRadius
}

func (v1 VectorGeneric[T]) chebyshevDistance(v2 VectorGeneric[T]) T {
	return max(Abs(v1.X-v2.X), Abs(v1.Y-v2.Y))
}

func (v VectorGeneric[T]) scaleInt(scale int) Vector {
	return Vector{int(v.X) * scale, int(v.Y) * scale}
}
