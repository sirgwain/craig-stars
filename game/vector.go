package game

import (
	"fmt"
	"math"
)

type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (v Vector) String() string {
	return fmt.Sprintf("(%f, %f)", v.X, v.Y)
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

func (v Vector) Round() Vector {
	return Vector{math.Round(v.X), math.Round(v.Y)}
}
