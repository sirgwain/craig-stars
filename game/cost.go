package game

import "math"

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

func FromMineral(m *Mineral) *Cost {
	return &Cost{
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

func (c *Cost) Total() int {
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
