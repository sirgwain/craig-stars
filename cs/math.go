package cs

import (
	"fmt"
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

// Round a float to the given precision value using math.Round().
func roundFloat(val float64, precision uint) float64 {
	ratio := float64(PowInt(10, precision))
	return math.Round(val*ratio) / ratio
}

// Round a float to the nearest whole number, rounding halves towards 0.
//
// This is distinct from [math.Round] which rounds numbers *away* from 0.
func roundHalfTowards0(x float64) float64 {
	// Implementation taken from a comment found in Golang's math.Round() source code.
	// Thanks, golang devs!
	t := math.Trunc(x)
	if Abs(x-t) > 0.5 {
		return t + math.Copysign(1, x)
	}
	return t
}

func splitValues(sourceCapacity, destCapacity1, destCapacity2 int, values ...int) ([]int, []int, error) {
	if destCapacity1+destCapacity2 != sourceCapacity {
		return nil, nil, fmt.Errorf("bucket sizes must sum to %d", sourceCapacity)
	}

	bucket1 := make([]int, len(values))
	bucket2 := make([]int, len(values))

	remaining1 := destCapacity1
	remaining2 := destCapacity2

	for i, count := range values {
		// Distribute proportionally
		split1 := (count*destCapacity1 + sourceCapacity/2) / sourceCapacity // Round to nearest
		split2 := count - split1

		// Ensure we do not exceed the remaining capacity
		if split2 > remaining2 {
			split2 = remaining2
			split1 = count - split2
		}
		if split1 > remaining1 {
			split1 = remaining1
			split2 = count - split1
		}

		bucket1[i] = split1
		bucket2[i] = split2

		remaining1 -= split1
		remaining2 -= split2
	}

	return bucket1, bucket2, nil
}

// Clamps value between minVal and maxVal and returns the result.
//
// Equivalent to
//
//	max(minVal, min(value, maxVal))
func Clamp[T constraints.Ordered](value, minVal, maxVal T) T {
	return max(minVal, min(value, maxVal))
}

// AbsMax returns the absolutely highest (furthest from 0)
// among a collection of similarly typed signed values.
// Panics if given no arguments.
//
// In the event one or more arguments have the same absolute value,
// the last one passed will take precedence.
func AbsMax[S constraints.Signed | constraints.Float](nums ...S) S {
	if len(nums) == 0 {
		panic("AbsMin called with no arguments")
	}

	result := nums[0]
	for _, value := range nums[1:] {
		if Abs(value) >= Abs(result) {
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

// Divide dividend by divisor and round the result away from 0.
//
// Effectively the opposite of truncation (which rounds towards 0 by chopping off decimals)
func divideRoundAway0[I constraints.Integer](dividend, divisor I) (quotient I) {
	if (dividend > 0) == (divisor > 0) { // Same sign
		return (dividend + divisor - 1) / divisor
	}
	return (dividend - divisor + 1) / divisor
}

// LogBase computes the logarithm of a number with the given base.
// Returns NaN if base <= 0 or base == 1, or if x <= 0.
func LogBase[F constraints.Float](base, x F) F {
	if base <= 0 || base == 1 || x <= 0 {
		return F(math.NaN())
	}
	// logB(x) = logA(x) / logA(B) for any A, B > 0 and != 1
	return F(math.Log2(float64(x)) / math.Log2(float64(base)))
}
