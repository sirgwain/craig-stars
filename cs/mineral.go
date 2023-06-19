package cs

type Mineral struct {
	Ironium   int `json:"ironium,omitempty"`
	Boranium  int `json:"boranium,omitempty"`
	Germanium int `json:"germanium,omitempty"`
}

func NewMineral(values [3]int) Mineral {
	return Mineral{
		Ironium:   values[0],
		Boranium:  values[1],
		Germanium: values[2],
	}

}

func (m Mineral) Total() int {
	return m.Ironium + m.Boranium + m.Germanium
}

func (m Mineral) ToSlice() [3]int {
	return [3]int{
		m.Ironium,
		m.Boranium,
		m.Germanium,
	}
}

// convert a mineral to a cargo
func (m Mineral) ToCargo() Cargo {
	return Cargo{
		Ironium:   m.Ironium,
		Boranium:  m.Boranium,
		Germanium: m.Germanium,
	}
}

// add an int to all components of the mineral
func (m Mineral) AddInt(num int) Mineral {
	return Mineral{
		Ironium:   m.Ironium + num,
		Boranium:  m.Boranium + num,
		Germanium: m.Germanium + num,
	}
}
