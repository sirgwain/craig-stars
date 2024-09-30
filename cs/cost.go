package cs

import (
	"fmt"
	"math"
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

func NewCost(
	ironium int,
	boranium int,
	germanium int,
	resources int,
) Cost {
	return Cost{ironium, boranium, germanium, resources}
}


// return the CostType with the Nth highest numerical value in a Cost struct (1 = highest, 2 = 2nd highest, etc etc)
//
// Ties are broken by REVERSE order of precendence (I/B/G/R)
func (c Cost) HighestType(ranking int) CostType {
	copy := c // make copy of cost struct so we can zero out values without affecting the original
	var highestType CostType
	for i := 0; i < MinInt(ranking, 4); i++ {
		// get the highest type in the cost struct
		highestType = copy.GetTypeFromAmount(MaxInt(copy.Ironium, copy.Boranium, copy.Germanium, copy.Resources))
		// For the record, this will never cause GetTypeFromAmount to panic because we are literally
		// comparing the cost struct's own values against themselves
		// set it to 0 
		copy.Set(highestType, 0)
	}
	return highestType
}

// return the first valid CostType in a Cost struct with the given numerical value
// returns an error if no CostType with the corresponding value exists
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
func FromMineralAndResources(c Mineral, resources int) Cost {
	return Cost{
		Ironium:   c.Ironium,
		Boranium:  c.Boranium,
		Germanium: c.Germanium,
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

func (c Cost) AddCargoMinerals(other Cargo) Cost {
	return Cost{
		Ironium:   c.Ironium + other.Ironium,
		Boranium:  c.Boranium + other.Boranium,
		Germanium: c.Germanium + other.Germanium,
		Resources: c.Resources,
	}
}

func (c Cost) Minus(other Cost) Cost {
	return Cost{
		Ironium:   c.Ironium - other.Ironium,
		Boranium:  c.Boranium - other.Boranium,
		Germanium: c.Germanium - other.Germanium,
		Resources: c.Resources - other.Resources,
	}
}

func (c Cost) MinusMineral(other Mineral) Cost {
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

func (a Cost) Divide(b Cost) float64 {
	var newIronium float64
	if b.Ironium == 0 {
		newIronium = math.Inf(1)
	} else {
		newIronium = float64(a.Ironium) / float64(b.Ironium)
	}

	var newBoranium float64
	if b.Boranium == 0 {
		newBoranium = math.Inf(1)
	} else {
		newBoranium = float64(a.Boranium) / float64(b.Boranium)
	}

	var newGermanium float64
	if b.Germanium == 0 {
		newGermanium = math.Inf(1)
	} else {
		newGermanium = float64(a.Germanium) / float64(b.Germanium)
	}

	var newResources float64
	if b.Resources == 0 {
		newResources = math.Inf(1)
	} else {
		newResources = float64(a.Resources) / float64(b.Resources)
	}

	return math.Min(newResources, math.Min(newIronium, math.Min(newBoranium, newGermanium)))
}

// divide a cost by a mineral. This will tell us if we have enough minerals to build some item
func (a Cost) DivideByMineral(b Mineral) float64 {
	var newIronium float64
	if b.Ironium == 0 {
		newIronium = math.Inf(1)
	} else {
		newIronium = float64(a.Ironium) / float64(b.Ironium)
	}

	var newBoranium float64
	if b.Boranium == 0 {
		newBoranium = math.Inf(1)
	} else {
		newBoranium = float64(a.Boranium) / float64(b.Boranium)
	}

	var newGermanium float64
	if b.Germanium == 0 {
		newGermanium = math.Inf(1)
	} else {
		newGermanium = float64(a.Germanium) / float64(b.Germanium)
	}

	return math.Min(newIronium, math.Min(newBoranium, newGermanium))
}

// Divide cost by an integer, either truncating or rounding up the result
func (c Cost) DivideByInt(divisor int, roundUp bool) Cost {
	if divisor == 0 {
		return Cost{int(math.Inf(1)), int(math.Inf(1)), int(math.Inf(1)), int(math.Inf(1))}
	}

	if roundUp {
		var ironium, boranium, germanium, resources int
		if c.Ironium%divisor > 0 {
			ironium = 1
		}
		if c.Boranium%divisor > 0 {
			boranium = 1
		}
		if c.Germanium%divisor > 0 {
			germanium = 1
		}
		if c.Resources%divisor > 0 {
			resources = 1
		}
		return Cost{
			Ironium:   c.Ironium/divisor + ironium,
			Boranium:  c.Boranium/divisor + boranium,
			Germanium: c.Germanium/divisor + germanium,
			Resources: c.Resources/divisor + resources,
		}
	} else {
		return Cost{
			Ironium:   c.Ironium / divisor,
			Boranium:  c.Boranium / divisor,
			Germanium: c.Germanium / divisor,
			Resources: c.Resources / divisor,
		}
	}
}

// Return greater of 2 cost structs for all ResourceTypes separately
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

// determine how many times and item costing cost can be built by available resources
func (available Cost) NumBuildable(cost Cost) int {
	buildable := Cost{math.MaxInt, math.MaxInt, math.MaxInt, math.MaxInt}

	if cost.Ironium > 0 {
		buildable.Ironium = available.Ironium / cost.Ironium
	}
	if cost.Boranium > 0 {
		buildable.Boranium = available.Boranium / cost.Boranium
	}
	if cost.Germanium > 0 {
		buildable.Germanium = available.Germanium / cost.Germanium
	}
	if cost.Resources > 0 {
		buildable.Resources = available.Resources / cost.Resources
	}

	return MinInt(
		buildable.Ironium,
		buildable.Boranium,
		buildable.Germanium,
		buildable.Resources,
	)
}
