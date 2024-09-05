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

func normalDistribution(μ float64, σ float64, x float64) float64 {
	var eNum = -math.Pow(x-μ, 2)
	var eDen = 2 * math.Pow(σ, 2)

	return math.Exp(eNum/eDen) / math.Sqrt(2*math.Pi*math.Pow(σ, 2))
}

func randomMineralConcentrationPDF(μ float64, σ float64, x int) float64 {
	maxY := normalDistribution(μ, σ, μ)

	if x < 30 {
		return 0.6
	} else {
		return normalDistribution(μ, σ, float64(x)) / maxY
	}
}

func randomMineralConcentration(random rng) int {
	var x int
	var y float64

	for x, y = random.Intn(126), random.Float64(); randomMineralConcentrationPDF(80, 20, x) < y; x, y = random.Intn(126), random.Float64() {
	}

	return x
}
