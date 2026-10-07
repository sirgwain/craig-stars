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
