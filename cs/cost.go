package cs

import (
	"fmt"
	"math"
	"slices"

	"golang.org/x/exp/constraints"
)

//tygo:emit
var _ = `export type cost<T extends number = number> = {
    ironium?: T;
    boranium?: T;
    germanium?: T;
    resources?: T;
};
`

// A Cost represents minerals and resources required to build something, like a mine, factory, or ship
// These are by default integers, but sometimes need to be treated as floats for applying
// discounts and miniaturization
type cost[T number] struct {
	Ironium   T `json:"ironium,omitempty"`
	Boranium  T `json:"boranium,omitempty"`
	Germanium T `json:"germanium,omitempty"`
	Resources T `json:"resources,omitempty"`
}

// An integer cost, used for most outwards-facing cost-related operations.
type Cost = cost[int]

// A floating point cost, used within internal calculations for determining unit rates.
type CostFloat64 = cost[float64]

// An integer or floating point value.
type number interface {
	constraints.Integer | constraints.Float
}

type CostType = ResourceType

var CostTypes = [4]CostType{
	Ironium,
	Boranium,
	Germanium,
	Resources,
}

// Create a new Cost struct with the given values.
func NewCost[T number](ironium, boranium, germanium, resources T) cost[T] {
	return cost[T]{
		Ironium:   ironium,
		Boranium:  boranium,
		Germanium: germanium,
		Resources: resources,
	}
}

// Create a Cost struct from a Mineral struct and a resources value.
func NewCostFromMineralAndResources(m Mineral, resources int) Cost {
	return Cost{
		Ironium:   m.Ironium,
		Boranium:  m.Boranium,
		Germanium: m.Germanium,
		Resources: resources,
	}
}

// HighestType returns the CostType and numerical value of the
// Nth highest value in a Cost struct.
// Negative indices count backwards from lowest value.
// (1 = highest, 2 = 2nd highest, -1 = lowest, etc etc).
//
// Ties are broken in order of precendence (I>B>G>R); tie order not affected by negative indices
//
// panics if ranking is 0 or if abs(ranking) is greater than 4
func (c cost[T]) HighestType(ranking int) (costType CostType, value T) {
	// Fun fact: this code is designed to work for any arbitrarily large struct
	// of similarly typed comparable values, only requiring changes to the
	// method signature, doc comment and error message
	a := c.ToSlice()
	if ranking == 0 || Abs(ranking) > len(a) {
		panic(fmt.Sprintf("Cost.HighestType() called with incorrect ranking %d; must be non-zero integer between -%d and %[2]d", ranking, len(a)))
	}

	slices.Sort(a[:])
	if ranking > 0 {
		value = a[len(a)-ranking] // Slice is ordered in ascending order, so biggest values will be at the end
	} else {
		value = a[-ranking-1] // negative indices count from the start (lowest first)
	}

	return c.GetTypeFromAmount(value), value
}

// Return the first valid CostType in a Cost struct with the given numerical value;
// panics if no CostType with the corresponding value exists
func (c cost[T]) GetTypeFromAmount(amt T) CostType {
	switch amt {
	case c.Ironium:
		return Ironium
	case c.Germanium:
		return Germanium
	case c.Boranium:
		return Boranium
	case c.Resources:
		return Resources
	}
	panic(fmt.Sprintf("GetTypeFromAmount called with value %v but no corresponding costType was found in cost struct; Struct values:\n%#v", amt, c))
}

func (c cost[T]) GetAmount(costType CostType) T {
	switch costType {
	case Ironium:
		return c.Ironium
	case Boranium:
		return c.Boranium
	case Germanium:
		return c.Germanium
	case Resources:
		return c.Resources
	default:
		panic(fmt.Sprintf("GetAmount called with invalid CostType %s", costType))
	}
}

// Set sets the value corresponding to costType to amt.
// Unlike all the other Cost functions, this _will_ mutate the original struct's values,
// and is best used for more complex cases not handled by other functions.
func (c *cost[T]) Set(costType CostType, amt T) {
	switch costType {
	case Ironium:
		c.Ironium = amt
	case Boranium:
		c.Boranium = amt
	case Germanium:
		c.Germanium = amt
	case Resources:
		c.Resources = amt
	default:
		panic(fmt.Sprintf("cost.Set called with invalid CostType %s", costType))
	}
}

// Return the Cargo equivalent of a Cost struct, truncating values as necessary.
func (c cost[T]) ToCargo() Cargo {
	return Cargo{
		Ironium:   int(c.Ironium),
		Boranium:  int(c.Boranium),
		Germanium: int(c.Germanium),
	}
}

// Return the Mineral equivalent of a Cost struct, truncating values as necessary.
func (c cost[T]) ToMineral() Mineral {
	return Mineral{
		Ironium:   int(c.Ironium),
		Boranium:  int(c.Boranium),
		Germanium: int(c.Germanium),
	}
}

func (c cost[T]) ToSlice() [4]T {
	return [4]T{
		c.Ironium,
		c.Boranium,
		c.Germanium,
		c.Resources,
	}
}

// Convert an integer cost into a floating point cost.
func (c cost[T]) ToCostFloat64() CostFloat64 {
	return CostFloat64{
		Ironium:   float64(c.Ironium),
		Boranium:  float64(c.Boranium),
		Germanium: float64(c.Germanium),
		Resources: float64(c.Resources),
	}
}

// Convert a floating point cost into an integer cost by truncating its values.
func (c cost[T]) ToCost() Cost {
	return Cost{
		Ironium:   int(c.Ironium),
		Boranium:  int(c.Boranium),
		Germanium: int(c.Germanium),
		Resources: int(c.Resources),
	}
}

// Returns the total sum of all resources in this Cost.
func (c cost[T]) Total() T {
	return c.Ironium + c.Boranium + c.Germanium + c.Resources
}

// Add 2 cost structs together and return the result.
func (c cost[T]) Add(other cost[T]) cost[T] {
	return cost[T]{
		Ironium:   c.Ironium + other.Ironium,
		Boranium:  c.Boranium + other.Boranium,
		Germanium: c.Germanium + other.Germanium,
		Resources: c.Resources + other.Resources,
	}
}

// Add a number to any singular component of a Cost struct.
func (c cost[T]) AddNum(costType CostType, amount T) cost[T] {
	switch costType {
	case Ironium:
		c.Ironium += amount
	case Boranium:
		c.Boranium += amount
	case Germanium:
		c.Germanium += amount
	case Resources:
		c.Resources += amount
	default:
		panic(fmt.Sprintf("AddNum called with invalid CostType %q", costType))
	}
	return c
}

// Add a Mineral to a cost struct and return the result.
func (c cost[T]) AddMineral(other Mineral) cost[T] {
	return cost[T]{
		Ironium:   c.Ironium + T(other.Ironium),
		Boranium:  c.Boranium + T(other.Boranium),
		Germanium: c.Germanium + T(other.Germanium),
		Resources: c.Resources,
	}
}

func (c cost[T]) Subtract(other cost[T]) cost[T] {
	return cost[T]{
		Ironium:   c.Ironium - other.Ironium,
		Boranium:  c.Boranium - other.Boranium,
		Germanium: c.Germanium - other.Germanium,
		Resources: c.Resources - other.Resources,
	}
}

func (c cost[T]) SubtractMineral(other Mineral) cost[T] {
	return cost[T]{
		Ironium:   c.Ironium - T(other.Ironium),
		Boranium:  c.Boranium - T(other.Boranium),
		Germanium: c.Germanium - T(other.Germanium),
		Resources: c.Resources,
	}
}

// Multiply a cost by an int or float and return the result.
func MultiplyCost[T number, F int | float64](c cost[T], factor F) cost[T] {
	return cost[T]{
		Ironium:   T(float64(c.Ironium) * float64(factor)),
		Boranium:  T(float64(c.Boranium) * float64(factor)),
		Germanium: T(float64(c.Germanium) * float64(factor)),
		Resources: T(float64(c.Resources) * float64(factor)),
	}
}

// Multiply a cost by another cost and return the resulting Cost struct.
//
// For multiplying a cost by an integer, use [MultiplyCost] instead
func MultiplyByCost[T, F number](c cost[T], other cost[F]) (result cost[T]) {
	return cost[T]{
		Ironium:   T(float64(c.Ironium) * float64(other.Ironium)),
		Boranium:  T(float64(c.Boranium) * float64(other.Boranium)),
		Germanium: T(float64(c.Germanium) * float64(other.Germanium)),
		Resources: T(float64(c.Resources) * float64(other.Resources)),
	}
}

// Divide a cost by another cost
// and return how many times divisor can go into dividend
// as a float64
func (dividend cost[T]) DivideCost(divisor cost[T]) float64 {
	quotient := CostFloat64{}
	for _, ct := range CostTypes {
		if divisor.GetAmount(ct) == 0 {
			quotient.Set(ct, float64(math.Inf(1)))
		} else {
			quotient.Set(ct, float64(dividend.GetAmount(ct))/float64(divisor.GetAmount(ct)))
		}
	}

	return quotient.MinAmount()
}

// Divide a cost by a mineral and return how many times divisor can go into dividend.
//
// This will tell us if we have enough minerals to build some item
// (and how many we can make)
func (dividend cost[T]) DivideMineral(divisor Mineral) float64 {
	dc := divisor.ToCost().ToCostFloat64()
	return dividend.ToCostFloat64().DivideCost(dc)
}

// Return greater of 2 Cost structs for all CostTypes separately
func (c cost[T]) Max(other cost[T]) cost[T] {
	return cost[T]{
		Ironium:   Max(c.Ironium, other.Ironium),
		Boranium:  Max(c.Boranium, other.Boranium),
		Germanium: Max(c.Germanium, other.Germanium),
		Resources: Max(c.Resources, other.Resources),
	}
}

// Return this Cost with a minimum of zero for each value
func (c cost[T]) MinZero() cost[T] {
	return cost[T]{
		Ironium:   Max(c.Ironium, 0),
		Boranium:  Max(c.Boranium, 0),
		Germanium: Max(c.Germanium, 0),
		Resources: Max(c.Resources, 0),
	}
}

// Return the lowest numerical value in a Cost struct
func (c cost[T]) MinAmount() T {
	return Min(c.Ironium, c.Boranium, c.Germanium, c.Resources)
}

// Round a cost struct's values by calling roundFunc on each of its values in turn.
func (c cost[T]) Round(roundFunc func(T) T) cost[T] {
	return cost[T]{
		Ironium:   roundFunc(c.Ironium),
		Boranium:  roundFunc(c.Boranium),
		Germanium: roundFunc(c.Germanium),
		Resources: roundFunc(c.Resources),
	}
}
