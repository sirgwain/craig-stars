package game

import "math"

func roundToNearest100f(value float64) int {
	return int(math.Round(value/100) * 100)
}

func roundToNearest100(value int) int {
	return int(math.Round(float64(value)/100) * 100)
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	} else {
		if value > max {
			return max
		}
	}
	return value
}

func clampf(value, min, max float64) float64 {
	if value < min {
		return min
	} else {
		if value > max {
			return max
		}
	}
	return value
}
