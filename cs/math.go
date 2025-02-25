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

// Round a float to the given precision value using math.Round()
func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

// Round a float to the nearest whole number, rounding halves towards 0.
// (This is distinct from math.Round() which rounds numbers *away* from 0.)
func roundHalfTowards0(x float64) float64 {
	// implementation taken from a comment found in Golang's math.Round() source code. Thanks, golang devs!
	t := math.Trunc(x)
	if Abs(x-t) > 0.5 {
		return t + math.Copysign(1, x)
	}
	return t
}

// Clamps value between min and max and returns the result.
// Equivalent to
//
//	Min(min, Max(value, max))
func Clamp[T constraints.Ordered](value, min, max T) T {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}

// Max returns the largest among a collection of similarly typed ordered values.
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

// Min returns the smallest among a collection of similarly typed ordered values.
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

// AbsMin returns the absolutely lowest (closest to 0)
// among a collection of similarly typed signed values.
// Panics if given no arguments.
//
// In the event one or more arguments have the same absolute value,
// the last one passed will take precedence.
func AbsMin[S constraints.Signed | constraints.Float](nums ...S) S {
	if len(nums) == 0 {
		panic("AbsMin called with no arguments")
	}

	result := nums[0]
	for _, value := range nums[1:] {
		if Abs(value) <= Abs(result) {
			result = value
		}
	}

	return result
}

// Raise an integer to the power of another integer and return the result.
//
// Does not support negative exponents (we *are* dealing with integers here after all)
func PowInt[I constraints.Integer](base, exponent I) I {
	var result I = 1
	// According to internet, this is the fastest way to do int exponentiation - by squaring
	for exponent != 0 {
		if exponent&1 == 1 {
			result *= base
		}
		exponent >>= 1
		base *= base
	}

	return result
}

// Abs returns the absolute value (unsigned portion) of a given number.
//
// Special cases:
//
//	Abs(±Inf) = +Inf
//	Abs(NaN) = NaN
func Abs[S constraints.Signed | constraints.Float](num S) S {
	if num < 0 {
		return -num
	}
	return num
}
