package cs

import "math"

type Rect struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width,omitempty"`
	Height float64 `json:"height,omitempty"`
}

func (rect Rect) Center() Vector {
	return Vector{
		rect.X + rect.Width/2,
		rect.Y + rect.Height/2,
	}
}

func (rect Rect) PointInRectangle(point Vector) bool {
	return point.X >= rect.X && point.X <= (rect.X+rect.Width) && point.Y >= rect.Y && point.Y <= (rect.Y+rect.Height)
}

func (rect Rect) PointInRotatedRectangle(point Vector, rectAngle float64) bool {
	// rotate around rectangle center by -rectAngle
	var s = math.Sin(-rectAngle)
	var c = math.Cos(-rectAngle)

	// set origin to rect center
	center := rect.Center()
	newPoint := point.Subtract(center)
	// rotate
	newPoint = Vector{newPoint.X*c - newPoint.Y*s, newPoint.X*s + newPoint.Y*c}
	// put origin back
	newPoint = newPoint.Add(center)

	// check if our transformed point is in the rectangle, which is no longer
	// rotated relative to the point
	return rect.PointInRectangle(newPoint)
}
