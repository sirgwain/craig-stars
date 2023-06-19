package cs

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

func clampFloat64(value, min, max float64) float64 {
	if value < min {
		return min
	} else {
		if value > max {
			return max
		}
	}
	return value
}

func maxInt(nums ...int) int {
	result := math.MinInt
	for _, value := range nums {
		if value > result {
			result = value
		}
	}

	return result
}

func minInt(nums ...int) int {
	result := math.MaxInt
	for _, value := range nums {
		if value < result {
			result = value
		}
	}

	return result
}

func minFloat64(nums ...float64) float64 {
	result := math.MaxFloat64
	for _, value := range nums {
		if value < result {
			result = value
		}
	}

	return result
}

func absInt(num int) int {
	if num < 0 {
		return -num
	}
	return num
}
