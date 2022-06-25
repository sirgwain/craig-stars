package game

import "fmt"

type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (v *Vector) String() string {
	return fmt.Sprintf("(%f, %f)", v.X, v.Y)
}

func (v *Vector) DistanceSquaredTo(to *Vector) float64 {
	return (v.X-to.X)*(v.X-to.X) + (v.Y-to.Y)*(v.Y-to.Y)
}
