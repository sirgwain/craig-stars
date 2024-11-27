package cs

import (
	"fmt"
	"math"
	"slices"
)

// A Cost represents minerals and resources required to build something, i.e. a mine, factory, or ship
type Cost struct {
	Ironium   int `json:"ironium,omitempty"`
	Boranium  int `json:"boranium,omitempty"`
	Germanium int `json:"germanium,omitempty"`
	Resources int `json:"resources,omitempty"`
}

type CostType = ResourceType

var CostTypes = [4]CostType{
	Ironium,
	Boranium,
	Germanium,
	Resources,
}

func NewCost(ironium, boranium, germanium, resources int) Cost {
	return Cost{ironium, boranium, germanium, resources}
}

// return the CostType with the Nth highest numerical value in a Cost struct (1 = highest, 2 = 2nd highest, etc etc).
// Negative indices count backwards from lowest value
//
// Ties are broken in order of precendence (I>B>G>R); tie order not affected by negative indices
func (c Cost) HighestType(ranking int) CostType {
	return c.GetTypeFromAmount(c.HighestAmount(ranking))
}

// return the numerical value of the Nth highest CostType in a Cost struct (1 = highest, 2 = 2nd highest, etc etc).
// Negative indices count backwards from lowest value
//
// Ties are broken in order of precendence (I>B>G>R); tie order not affected by negative indices
func (c Cost) HighestAmount(ranking int) int {
	a := c.ToSlice()
	slice := slices.Clone(a[:])
	slices.Sort(slice)
	if ranking < 0 {
		slices.SortStableFunc(slice, func(a, b int) int { return b - a })
		ranking = -ranking
	}
	return slice[len(slice)-ranking]
}

// return the first valid CostType in a Cost struct with the given numerical value;
// panics if no CostType with the corresponding value exists
func (c Cost) GetTypeFromAmount(amt int) CostType {
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

func (c Cost) GetAmount(costType CostType) int {
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

func (c Cost) Set(costType CostType, amt int) Cost {
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

func FromMineralAndResources(m Mineral, resources int) Cost {
	return Cost{
		Ironium:   m.Ironium,
		Boranium:  m.Boranium,
		Germanium: m.Germanium,
		Resources: resources,
	}
}

func FromMineral(c Mineral) Cost {
	return Cost{
		Ironium:   c.Ironium,
		Boranium:  c.Boranium,
		Germanium: c.Germanium,
	}
}

func (c Cost) ToCargo() Cargo {
	return Cargo{
		Ironium:   c.Ironium,
		Boranium:  c.Boranium,
		Germanium: c.Germanium,
	}
}

func (c Cost) ToMineral() Mineral {
	return Mineral{
		Ironium:   c.Ironium,
		Boranium:  c.Boranium,
		Germanium: c.Germanium,
	}
}

func (c Cost) ToSlice() [4]int {
	return [4]int{
		c.Ironium,
		c.Boranium,
		c.Germanium,
		c.Resources,
	}
}

// convert an int cost into a costFloat64 struct for use in calculations
func (c Cost) ToCostFloat64() costFloat64 {
	return costFloat64{
		ironium:   float64(c.Ironium),
		boranium:  float64(c.Boranium),
		germanium: float64(c.Germanium),
		resources: float64(c.Resources),
	}
}

func (c Cost) Total() int {
	return c.Ironium + c.Boranium + c.Germanium + c.Resources
}

func (c Cost) Add(other Cost) Cost {
	return Cost{
		Ironium:   c.Ironium + other.Ironium,
		Boranium:  c.Boranium + other.Boranium,
		Germanium: c.Germanium + other.Germanium,
		Resources: c.Resources + other.Resources,
	}
}

func (c Cost) AddInt(costType CostType, amount int) Cost {
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

func (c Cost) AddMineral(other Mineral) Cost {
	return Cost{
		Ironium:   c.Ironium + other.Ironium,
		Boranium:  c.Boranium + other.Boranium,
		Germanium: c.Germanium + other.Germanium,
		Resources: c.Resources,
	}
}

func (c Cost) Subtract(other Cost) Cost {
	return Cost{
		Ironium:   c.Ironium - other.Ironium,
		Boranium:  c.Boranium - other.Boranium,
		Germanium: c.Germanium - other.Germanium,
		Resources: c.Resources - other.Resources,
	}
}

func (c Cost) SubtractMineral(other Mineral) Cost {
	return Cost{
		Ironium:   c.Ironium - other.Ironium,
		Boranium:  c.Boranium - other.Boranium,
		Germanium: c.Germanium - other.Germanium,
		Resources: c.Resources,
	}
}

func (c Cost) MultiplyInt(factor int) Cost {
	return Cost{
		Ironium:   c.Ironium * factor,
		Boranium:  c.Boranium * factor,
		Germanium: c.Germanium * factor,
		Resources: c.Resources * factor,
	}
}

func (c Cost) MultiplyFloat64(factor float64) Cost {
	return Cost{
		Ironium:   int(float64(c.Ironium) * factor),
		Boranium:  int(float64(c.Boranium) * factor),
		Germanium: int(float64(c.Germanium) * factor),
		Resources: int(float64(c.Resources) * factor),
	}
}

// divide a cost by another cost
// and return how many times divisor can go into dividend
//
// This functionally replaces cost.NumBuildable; 
// the latter can be e as the latter was essentially just "divide but int"
func (dividend Cost) Divide(divisor Cost) float64 {
	divisorFloat := divisor.ToCostFloat64()
	dividendFloat := dividend.ToCostFloat64()
	quotient := costFloat64{}
	for _, ct := range CostTypes {
		if divisorFloat.getAmount(ct) == 0 {
			quotient = quotient.set(ct, math.Inf(1))
		} else {
			quotient = quotient.set(ct, dividendFloat.getAmount(ct)/divisorFloat.getAmount(ct))
		}
	}
	return quotient.minAmount()
}

// divide a cost by a mineral and return how many times divisor can go into dividend.
//
// This will tell us if we have enough minerals to build some item
// (and how many we can make)
func (dividend Cost) DivideMineral(divisor Mineral) float64 {
	dc := divisor.ToCost()
	return dividend.Divide(dc)
}

// Return greater of 2 cost structs for all CostTypes separately
func (c Cost) Max(other Cost) Cost {
	return Cost{
		Ironium:   MaxInt(c.Ironium, other.Ironium),
		Boranium:  MaxInt(c.Boranium, other.Boranium),
		Germanium: MaxInt(c.Germanium, other.Germanium),
		Resources: MaxInt(c.Resources, other.Resources),
	}
}

func (c Cost) Negate() Cost {
	return Cost{
		Ironium:   -c.Ironium,
		Boranium:  -c.Boranium,
		Germanium: -c.Germanium,
		Resources: -c.Resources,
	}
}

// return this cost with a minimum of zero for each value
func (c Cost) MinZero() Cost {
	return Cost{
		Ironium:   MaxInt(c.Ironium, 0),
		Boranium:  MaxInt(c.Boranium, 0),
		Germanium: MaxInt(c.Germanium, 0),
		Resources: MaxInt(c.Resources, 0),
	}
}
