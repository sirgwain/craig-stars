package cs

import "math"

// A Cost represents minerals and resources required to build something, i.e. a mine, factory, or ship
type Cost struct {
	Ironium   int `json:"ironium,omitempty"`
	Boranium  int `json:"boranium,omitempty"`
	Germanium int `json:"germanium,omitempty"`
	Resources int `json:"resources,omitempty"`
}

func NewCost(
	ironium int,
	boranium int,
	germanium int,
	resources int,
) Cost {
	return Cost{ironium, boranium, germanium, resources}
}

func (c Cost) GetAmount(resourceType int) int {
	switch resourceType {
	case 0:
		return c.Ironium
	case 1:
		return c.Boranium
	case 2:
		return c.Germanium
	case 3:
		return c.Resources
	}
	return 0
}

func (c Cost) AddInt(resourceType int, amount int) Cost {
	switch resourceType {
	case 0:
		c.Ironium += amount
	case 1:
		c.Boranium += amount
	case 2:
		c.Germanium += amount
	case 3:
		c.Resources += amount
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

func FromMineral(m Mineral) Cost {
	return Cost{
		Ironium:   m.Ironium,
		Boranium:  m.Boranium,
		Germanium: m.Germanium,
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
		var I, B, G, R int
		if c.Ironium%divisor > 0 {
			I = 1
		}
		if c.Boranium%divisor > 0 {
			B = 1
		}
		if c.Germanium%divisor > 0 {
			G = 1
		}
		if c.Resources%divisor > 0 {
			R = 1
		}
		return Cost{
			Ironium:   c.Ironium/divisor + I,
			Boranium:  c.Boranium/divisor + B,
			Germanium: c.Germanium/divisor + G,
			Resources: c.Resources/divisor + R,
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

// Return greater of 2 cost structs for the specified ResourceType;
// anything outside of range (0-3) evaluates it for all ResourceTypes separately
func (c Cost) Max(other Cost, resourceType int) Cost {
	a := Cost{}
	if 3 >= resourceType && resourceType >= 0 {
		switch resourceType {
		case 0:
			a.Ironium = MaxInt(c.Ironium, other.Ironium)
		case 1:
			a.Boranium = MaxInt(c.Boranium, other.Boranium)
		case 2:
			a.Germanium = MaxInt(c.Germanium, other.Germanium)
		case 3:
			a.Resources = MaxInt(c.Resources, other.Resources)
		} 
	} else {
		a.Ironium = MaxInt(c.Ironium, other.Ironium)
		a.Boranium = MaxInt(c.Boranium, other.Boranium)
		a.Germanium = MaxInt(c.Germanium, other.Germanium)
		a.Resources = MaxInt(c.Resources, other.Resources)
	}
	return a
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
