package cs

import (
	"math"
)

// the rng rules all
type rng interface {
	// Float64 returns, as a float64, a pseudo-random number in [0.0,1.0)
	// from the default Source.
	Float64() float64

	// Intn returns, as an int, a non-negative pseudo-random number in [0,n).
	// It panics if n <= 0.
	Intn(n int) int

	// Shuffle pseudo-randomizes the order of elements using the default Source.
	// n is the number of elements. Shuffle panics if n < 0.
	// swap swaps the elements with indexes i and j.
	Shuffle(n int, swap func(i, j int))
}

var modifiedNormalDistribution [126]float64
var initialized = false

func normalDistribution(μ float64, σ float64, x float64) float64 {
	var eNum = -math.Pow(x-μ, 2)
	var eDen = 2 * math.Pow(σ, 2)

	return math.Exp(eNum/eDen) / math.Sqrt(2*math.Pi*math.Pow(σ, 2))
}

func initRandom() {
	maxY := 0.0

	for i := range modifiedNormalDistribution {
		modifiedNormalDistribution[i] = normalDistribution(80, 20, float64(i))

		if modifiedNormalDistribution[i] > maxY {
			maxY = modifiedNormalDistribution[i]
		}
	}

	for i, v := range modifiedNormalDistribution {
		modifiedNormalDistribution[i] = v / maxY
	}

	for i := 0; i < 30; i++ {
		modifiedNormalDistribution[i] = 0.6
	}
}

func modifiedNormalRandom(random rng) int {
	if !initialized {
		initRandom()
		initialized = true
	}

	var x int
	var y float64

	for x, y = random.Intn(126), random.Float64(); modifiedNormalDistribution[x] < y; x, y = random.Intn(126), random.Float64() {
	}

	return x
}
