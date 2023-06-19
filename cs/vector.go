package cs

import (
	"fmt"
	"math"
)

type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

var VectorRight Vector = Vector{1, 0}
var VectorLeft Vector = Vector{-1, 0}
var VectorUp Vector = Vector{0, 1}
var VectorDown Vector = Vector{0, -1}

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
// intersect
func segmentIntersectsCircle(segmentFrom, segmentTo, circlePosition Vector, circleRadius float64) float64 {
	lineVec := segmentTo.Subtract(segmentFrom)
	vecToLine := segmentFrom.Subtract(circlePosition)

	// Create a quadratic formula of the form ax^2 + bx + c = 0
	var a, b, c float64

	a = lineVec.Dot(lineVec)
	b = 2 * vecToLine.Dot(lineVec)
	c = vecToLine.Dot(vecToLine) - circleRadius*circleRadius

	// Solve for t.
	sqrtterm := b*b - 4*a*c

	// If the term we intend to square root is less than 0 then the answer won't be real,
	// so it definitely won't be t in the range 0 to 1.
	if sqrtterm < 0 {
		return -1
	}

	// If we can assume that the line segment starts outside the circle (e.g. for continuous time collision detection)
	// then the following can be skipped and we can just return the equivalent of res1.
	sqrtterm = math.Sqrt(sqrtterm)
	res1 := (-b - sqrtterm) / (2 * a)
	res2 := (-b + sqrtterm) / (2 * a)

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
