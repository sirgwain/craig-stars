package cs

import (
	"math"

	"golang.org/x/exp/constraints"
)

// Round a value to the nearest 100 using the specified rounding function
// and return the result
//
// Population is often updated with floating point math, but we typically have to convert
// it back to Colonist cargo values, which are stored in units of 100 colonists per 1kT
func roundTo100[T int | float64](value T, roundFunc func(float64) float64) int {
	return int(roundFunc(float64(value)/100) * 100)
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

// round a float to nearest whole number, rounding halves down
func roundHalfDown(x float64) float64 {
	if x > 0 {
		return math.Floor(x + 0.5)
	}
	return math.Ceil(x - 0.5)
}

func Clamp[T constraints.Ordered](value, min, max T) T {
	if value < min {
		return min
	} else {
		if value > max {
			return max
		}
	}
	return value
}

func Max[T constraints.Ordered](nums ...T) T {
	if len(nums) == 0 {
		panic("Max called with no arguments")
	}

	result := nums[0]
	for _, value := range nums[1:] {
		if value > result {
			result = value
		}
	}

	return result
}

func Min[T constraints.Ordered](nums ...T) T {
	if len(nums) == 0 {
		panic("Min called with no arguments")
	}

	result := nums[0]
	for _, value := range nums[1:] {
		if value < result {
			result = value
		}
	}

	return result
}

func Abs[T int](num T) T {
	if num < 0 {
		return -num
	}
	return num
}
