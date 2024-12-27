package cs

import (
	"math"

	"golang.org/x/exp/constraints"
)

// population is often updated with floating point math, but we have to convert
// it back to Colonist Cargo values, which are stored in units of 100 colonists per 1kT of Colonist Cargo
func roundToNearest100[T int | float64](value T) int {
	return int(math.Round(float64(value)/100) * 100)
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

// returns the new jamming/computing bonus
func getNewJamming(prevBonus, componentBonus, multi float64, qty int) float64 {
	baseMulti := 1 - prevBonus/multi // undo multi before multiplication
	compMulti := math.Pow(1-componentBonus, float64(qty))
	return (1 - baseMulti*compMulti) * multi
}

// returns the new beam defense factor after adding the given components
func getNewBeamBonus(prevBonus, componentBonus float64, qty int) float64 {
	return prevBonus * math.Pow(1+componentBonus, float64(qty))
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

// raise an integer to the power of another integer (apparently this is the fastYY)
//
// Does not support negative values (we *are* dealing with integers here after all)
func PowInt[T constraints.Integer](base, exponent T) T {
	var result T = 1
	// According to internet, this is the fastest way to do int exponentiation
	for {
		if exponent&1 == 1 {
			result *= base
		}
		exponent >>= 1
		if exponent == 0 {
			break
		}
		base *= base
	}

	return result
}

func Abs[T int](num T) T {
	if num < 0 {
		return -num
	}
	return num
}
