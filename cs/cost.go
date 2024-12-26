package cs

import (
	"cmp"
	"fmt"
	"math"
	"slices"
)

// Costs are by default ints, but sometimes we need to treat them as floats for applying
// discounts and miniaturization
type Cost = cost[int]
type CostFloat64 = cost[float64]

// Costs can be ints or float64... for now
type costTypeConstraints interface {
	int | float64
}

// A Cost represents minerals and resources required to build something, i.e. a mine, factory, or ship
type cost[T costTypeConstraints] struct {
	Ironium   T `json:"ironium,omitempty"`
	Boranium  T `json:"boranium,omitempty"`
	Germanium T `json:"germanium,omitempty"`
	Resources T `json:"resources,omitempty"`
}

type CostType = ResourceType

var CostTypes = [4]CostType{
	Ironium,
	Boranium,
	Germanium,
	Resources,
}

func NewCost[T costTypeConstraints](ironium, boranium, germanium, resources T) cost[T] {
	return cost[T]{ironium, boranium, germanium, resources}
}

func FromMineralAndResources(m Mineral, resources int) Cost {
	return Cost{
		Ironium:   m.Ironium,
		Boranium:  m.Boranium,
		Germanium: m.Germanium,
		Resources: resources,
	}
}

func FromMineral[T costTypeConstraints](c Mineral) cost[T] {
	return cost[T]{
		Ironium:   T(c.Ironium),
		Boranium:  T(c.Boranium),
		Germanium: T(c.Germanium),
	}
}

func MultiplyCost[T costTypeConstraints, F int | float64](c cost[T], factor F) cost[T] {
	return cost[T]{
		Ironium:   T(float64(c.Ironium) * float64(factor)),
		Boranium:  T(float64(c.Boranium) * float64(factor)),
		Germanium: T(float64(c.Germanium) * float64(factor)),
		Resources: T(float64(c.Resources) * float64(factor)),
	}
}

// return the CostType with the Nth highest numerical value in a Cost struct (1 = highest, 2 = 2nd highest, etc etc).
// Negative indices count backwards from lowest value
//
// Ties are broken in order of precendence (I>B>G>R); tie order not affected by negative indices
func (c cost[T]) HighestType(ranking int) CostType {
	return c.GetTypeFromAmount(c.HighestAmount(ranking))
}

// return the numerical value of the Nth highest CostType in a Cost struct (1 = highest, 2 = 2nd highest, etc etc).
// Negative indices count backwards from lowest value
//
// Ties are broken in order of precendence (I>B>G>R); tie order not affected by negative indices
func (c cost[T]) HighestAmount(ranking int) T {
	a := c.ToSlice()
	slice := slices.Clone(a[:])
	slices.SortStableFunc(slice, func(a, b T) int { return cmp.Compare(a, b) })
	if ranking < 0 {
		slices.Sort(slice)
		ranking = -ranking
	}
	return slice[len(slice)-ranking]
}

// return the first valid CostType in a Cost struct with the given numerical value;
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
	panic(fmt.Sprintf("GetTypeFromAmount called with value %v but no corresponding costType was found in cost struct; \nStruct values:\nIronium: %v\nBoranium: %v\nGermanium: %v\nResources: %v",
		amt, c.Ironium, c.Boranium, c.Germanium, c.Resources))
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
	}
	panic(fmt.Sprintf("GetAmount called with invalid CostType %s", costType))
}

func (c cost[T]) Set(costType CostType, amt T) cost[T] {
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
		panic(fmt.Sprintf("SetAmount called with invalid CostType %s", costType))
	}
	return c
}

func (c cost[T]) ToCargo() Cargo {
	return Cargo{
		Ironium:   int(c.Ironium),
		Boranium:  int(c.Boranium),
		Germanium: int(c.Germanium),
	}
}

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

// convert an int cost into a costFloat64 struct for use in calculations
func (c cost[T]) ToCostFloat64() CostFloat64 {
	return CostFloat64{
		Ironium:   float64(c.Ironium),
		Boranium:  float64(c.Boranium),
		Germanium: float64(c.Germanium),
		Resources: float64(c.Resources),
	}
}

func (c cost[T]) ToCost() Cost {
	return Cost{
		Ironium:   int(c.Ironium),
		Boranium:  int(c.Boranium),
		Germanium: int(c.Germanium),
		Resources: int(c.Resources),
	}
}

func (c cost[T]) Total() T {
	return c.Ironium + c.Boranium + c.Germanium + c.Resources
}

func (c cost[T]) Add(other cost[T]) cost[T] {
	return cost[T]{
		Ironium:   c.Ironium + other.Ironium,
		Boranium:  c.Boranium + other.Boranium,
		Germanium: c.Germanium + other.Germanium,
		Resources: c.Resources + other.Resources,
	}
}

func (c cost[T]) AddInt(costType CostType, amount T) cost[T] {
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
		panic(fmt.Sprintf("AddInt called with invalid CostType %s", costType))
	}
	return c
}

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

// divide a cost by another cost
// and return how many times divisor can go into dividend
//
// This functionally replaces cost.NumBuildable;
// the latter can be e as the latter was essentially just "divide but int"
func (dividend cost[T]) Divide(divisor CostFloat64) float64 {
	quotient := CostFloat64{}
	for _, ct := range CostTypes {
		if divisor.GetAmount(ct) == 0 {
			quotient = quotient.Set(ct, float64(math.Inf(1)))
		} else {
			quotient = quotient.Set(ct, float64(dividend.GetAmount(ct))/float64(divisor.GetAmount(ct)))
		}
	}
	return quotient.MinAmount()
}

// divide a cost by a mineral and return how many times divisor can go into dividend.
//
// This will tell us if we have enough minerals to build some item
// (and how many we can make)
func (dividend cost[T]) DivideMineral(divisor Mineral) float64 {
	dc := divisor.ToCost().ToCostFloat64()
	return dividend.Divide(dc)
}

// Return greater of 2 cost structs for all CostTypes separately
func (c cost[T]) Max(other cost[T]) cost[T] {
	return cost[T]{
		Ironium:   Max(c.Ironium, other.Ironium),
		Boranium:  Max(c.Boranium, other.Boranium),
		Germanium: Max(c.Germanium, other.Germanium),
		Resources: Max(c.Resources, other.Resources),
	}
}

func (c cost[T]) Negate() cost[T] {
	return cost[T]{
		Ironium:   -c.Ironium,
		Boranium:  -c.Boranium,
		Germanium: -c.Germanium,
		Resources: -c.Resources,
	}
}

// return this cost with a minimum of zero for each value
func (c cost[T]) MinZero() cost[T] {
	return cost[T]{
		Ironium:   Max(c.Ironium, 0),
		Boranium:  Max(c.Boranium, 0),
		Germanium: Max(c.Germanium, 0),
		Resources: Max(c.Resources, 0),
	}
}

func (c cost[T]) MinAmount() T {
	return Min(c.Ironium, c.Boranium, c.Germanium, c.Resources)
}

// Round a cost struct's values with passed in function
func (c cost[T]) Round(roundFunc func(float64) float64) cost[T] {
	return cost[T]{
		Ironium:   T(roundFunc(float64(c.Ironium))),
		Boranium:  T(roundFunc(float64(c.Boranium))),
		Germanium: T(roundFunc(float64(c.Germanium))),
		Resources: T(roundFunc(float64(c.Resources))),
	}
}
