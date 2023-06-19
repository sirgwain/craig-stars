package game

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

func (c *Cost) Total() int {
	return c.Ironium + c.Boranium + c.Germanium + c.Resources
}

func (c *Cost) Add(other Cost) Cost {
	return Cost{
		Ironium:   c.Ironium + other.Ironium,
		Boranium:  c.Boranium + other.Boranium,
		Germanium: c.Germanium + other.Germanium,
		Resources: c.Resources + other.Resources,
	}
}
