package cs

import "math"

type Rect struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`
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
	newPoint := point.Subtract(center).ToFloat64()
	// rotate
	newPoint = VectorFloat64{float64(newPoint.X)*c - float64(newPoint.Y)*s, float64(newPoint.X)*s + float64(newPoint.Y)*c}
	// put origin back
	newPoint = newPoint.Add(center.ToFloat64())

	// check if our transformed point is in the rectangle, which is no longer
	// rotated relative to the point
	return newPoint.X >= float64(rect.X) && newPoint.X <= float64(rect.X+rect.Width) && newPoint.Y >= float64(rect.Y) && newPoint.Y <= float64(rect.Y+rect.Height)
}
