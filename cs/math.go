package cs

import (
	"math"

	"golang.org/x/exp/constraints"
)

// Round a value to a multiple of 100 using the specified rounding function
// and return the result as an integer.
//
// Population is often updated with floating point/integer math, but we typically have to convert
// it back to Colonist cargo values, which are stored in units of 100 colonists per 1kT
func roundTo100[T int | float64](value T, roundFunc func(float64) float64) int {
	return int(roundFunc(float64(value)/100) * 100)
}

// Round a float to the given precision value using math.Round()
func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

// Round a float to the nearest whole number, rounding halves towards 0.
//
// This is distinct from math.Round() which rounds numbers *away* from 0.
func roundHalfTowards0(x float64) float64 {
	// Implementation taken from a comment found in Golang's math.Round() source code.
	// Thanks, golang devs!
	t := math.Trunc(x)
	if Abs(x-t) > 0.5 {
		return t + math.Copysign(1, x)
	}
	return t
}

// Clamps value between minVal and maxVal and returns the result.
//
// Equivalent to
//
//	max(minVal, min(value, maxVal))
func Clamp[T constraints.Ordered](value, minVal, maxVal T) T {
	return max(minVal, min(value, maxVal))
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
func Abs[T constraints.Integer | constraints.Float](num T) T {
	if num < 0 {
		return -num
	}
	return num
}
