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
