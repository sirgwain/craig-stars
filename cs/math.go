package cs

import (
	"math"

	"golang.org/x/exp/constraints"
)

// Round a number (int or float) to the nearest multiple of 100 and return the resulting integer.
//
// Typically used to convert floating-point population values back into colonist Cargo values, 
// which are stored in discrete units of 100 colonists/1kT.
func roundToNearest100[T int | float64](value T) int {
	return int(math.Round(float64(value)/100) * 100)
}

// round a float to the given precision
func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

// round a float to the nearest whole number, rounding halves down
func roundHalfDown(x float64) float64 {
	if x > 0 {
		return math.Floor(x + 0.5)
	}
	return math.Ceil(x - 0.5)
}

// Clamps the passed in value between min and max by ensuring. 
func Clamp[T constraints.Ordered](value, min, max T) T {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}

// Returns the highest among a collection of similarly typed ordered values.
// Panics if given no arguments.
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

// Returns the lowest among a collection of similarly typed ordered values.
// Panics if given no arguments.
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

// Raise an integer to the power of another integer and return the result.
//
// Does not support negative exponents (we *are* dealing with integers here after all)
func PowInt[T constraints.Integer](base, exponent T) T {
	var result T = 1
	// According to internet, this is the fastest way to do int exponentiation
	for exponent != 0 {
		if exponent&1 == 1 {
			result *= base
		}
		exponent >>= 1
		base *= base
	}

	return result
}

// Returns the absolute value (unsigned portion) of a given number.
func Abs[T number](num T) T {
	if num < 0 {
		return -num
	}
	return num
}
