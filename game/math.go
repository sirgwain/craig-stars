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

func MaxInt(nums ...int) int {
	result := math.MinInt
	for _, value := range nums {
		if value > result {
			result = value
		}
	}

	return result
}

func MinInt(nums ...int) int {
	result := math.MaxInt
	for _, value := range nums {
		if value < result {
			result = value
		}
	}

	return result
}

func MinFloat64(nums ...float64) float64 {
	result := math.MaxFloat64
	for _, value := range nums {
		if value < result {
			result = value
		}
	}

	return result
}
