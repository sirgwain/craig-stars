package cs

import (
	"fmt"
	"math"
)

const doesNotIntersect = -1

// A 2D vector (direction or point), containing
// handy functions for moving things in space, calculating distance, etc.
// Many of these functions were taken from Godot source, thanks Godot folks!
type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (v Vector) String() string {
	return fmt.Sprintf("(%0.0f, %0.0f)", v.X, v.Y)
}

func (v Vector) DistanceSquaredTo(to Vector) float64 {
	return math.Pow(v.X-to.X, 2) + math.Pow(v.Y-to.Y, 2)
}

// Return the distance from one vector to another using the Pythagorean theorem.
func (v Vector) DistanceTo(to Vector) float64 {
	return math.Sqrt(v.DistanceSquaredTo(to))
}

func (addend Vector) Add(augend Vector) Vector {
	return Vector{addend.X + augend.X, addend.Y + augend.Y}
}

// Subtract 2 vectors and return the result.
//
// Typically used to generate a direction vector from 2 point vectors.
func (minuend Vector) Subtract(subtrahend Vector) Vector {
	return Vector{minuend.X - subtrahend.X, minuend.Y - subtrahend.Y}
}

// Multiply a vector by the given scale factor and return the result.
func (v Vector) Multiply(scale float64) Vector {
	return Vector{v.X * scale, v.Y * scale}
}

func (v Vector) LengthSquared() float64 {
	return math.Pow(v.X, 2) + math.Pow(v.Y, 2)
}

func (v Vector) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

// return this vector with length normalized to equal 1.
func (v Vector) Normalized() Vector {
	if v == (Vector{}) {
		return v
	}
	return Vector{
		X: v.X / v.Length(),
		Y: v.Y / v.Length(),
	}
}

// Return the dot product of 2 vectors.
//
//	a.Dot(b) = a.x*b.x + a.y*b.y
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
		return doesNotIntersect
	}

	// If we can assume that the line segment starts outside the circle
	// (e.g. for continuous time collision detection), the following can be
	// skipped and we can just return the equivalent of root1.
	discriminant = math.Sqrt(discriminant)
	root1 := (-b - discriminant) / (2 * a)
	root2 := (-b + discriminant) / (2 * a)

	if root1 >= 0 && root1 <= 1 {
		return root1
	}
	if root2 >= 0 && root2 <= 1 {
		return root2
	}
	return doesNotIntersect
}

// Returns true if this point is in a circle
func isPointInCircle(point, circlePosition Vector, circleRadius float64) bool {
	return point.DistanceSquaredTo(circlePosition) <= circleRadius*circleRadius
}
