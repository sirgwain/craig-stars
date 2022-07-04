package game

type Cargo struct {
	Ironium   int `json:"ironium,omitempty"`
	Boranium  int `json:"boranium,omitempty"`
	Germanium int `json:"germanium,omitempty"`
	Colonists int `json:"colonists,omitempty"`
}

func (c *Cargo) Add(other Cargo) Cargo {
	return Cargo{
		Ironium:   c.Ironium + other.Ironium,
		Boranium:  c.Boranium + other.Boranium,
		Germanium: c.Germanium + other.Germanium,
		Colonists: c.Colonists + other.Colonists,
	}
}

func (c *Cargo) AddMineral(other Mineral) Cargo {
	return Cargo{
		Ironium:   c.Ironium + other.Ironium,
		Boranium:  c.Boranium + other.Boranium,
		Germanium: c.Germanium + other.Germanium,
		Colonists: c.Colonists,
	}
}

func (c *Cargo) Total() int {
	return c.Ironium + c.Boranium + c.Germanium + c.Colonists
}
