package cs

import (
	"fmt"
	"math"
)

// A simple 2D vector with handy functions for moving things in space, calculating distance, etc
// Many of these functions were taken from Godot source, thanks Godot folks.
type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (v Vector) String() string {
	return fmt.Sprintf("(%0.0f, %0.0f)", v.X, v.Y)
}

func (v Vector) DistanceSquaredTo(to Vector) float64 {
	return (v.X-to.X)*(v.X-to.X) + (v.Y-to.Y)*(v.Y-to.Y)
}

func (v Vector) DistanceTo(to Vector) float64 {
	return math.Sqrt((v.X-to.X)*(v.X-to.X) + (v.Y-to.Y)*(v.Y-to.Y))
}

func (minuend Vector) Subtract(subtrahend Vector) Vector {
	return Vector{minuend.X - subtrahend.X, minuend.Y - subtrahend.Y}
}

func (v1 Vector) Add(v2 Vector) Vector {
	return Vector{v1.X + v2.X, v1.Y + v2.Y}
}

func (v Vector) Scale(scale float64) Vector {
	return Vector{v.X * scale, v.Y * scale}
}

func (v Vector) LengthSquared() float64 {
	return (v.X * v.X) + (v.Y * v.Y)
}

func (v Vector) Length() float64 {
	return math.Sqrt((v.X * v.X) + (v.Y * v.Y))
}

func (v Vector) Normalized() Vector {
	lengthsq := v.LengthSquared()

	if lengthsq == 0 {
		v.X = 0
		v.Y = 0
	} else {
		length := math.Sqrt(lengthsq)
		v.X /= length
		v.Y /= length
	}
	return v
}

func (v Vector) Dot(other Vector) float64 {
	return v.X*other.X + v.Y*other.Y
}

func (v Vector) Round() Vector {
	return Vector{math.Round(v.X), math.Round(v.Y)}
}

// SegmentIntersectsCircle checks whether a segment intersects a circle or not.
// This returns what percent of the segment is NOT in the circle, or -1 if it doesn't
// intersect.
// Who would have thought the godot developers would write a perfect function for determining
// if we collide with a minefield
// while moving through space?
// https://github.com/godotengine/godot/blob/4.1.2-stable/core/math/geometry_2d.h#L217
func segmentIntersectsCircle(segmentFrom, segmentTo, circlePosition Vector, circleRadius float64) (percentOutside float64) {
	lineVec := segmentTo.Subtract(segmentFrom)
	vecToLine := segmentFrom.Subtract(circlePosition)

	// Create a quadratic formula of the form ax^2 + bx + c = 0
	// Hope you remembered your high school algebra!
	var a, b, c float64

	a = lineVec.Dot(lineVec)
	b = 2 * vecToLine.Dot(lineVec)
	c = vecToLine.Dot(vecToLine) - circleRadius*circleRadius

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
	res1 := (-b - discriminant) / (2 * a)
	res2 := (-b + discriminant) / (2 * a)

	if res1 >= 0 && res1 <= 1 {
		return res1
	}
	if res2 >= 0 && res2 <= 1 {
		return res2
	}
	return -1
}

// Returns true if this point is in a circle
func isPointInCircle(point, circlePosition Vector, circleRadius float64) bool {
	return point.DistanceSquaredTo(circlePosition) <= circleRadius*circleRadius
}
