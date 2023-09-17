package cs

import (
	"fmt"
	"strings"
)

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

func (c Mineral) PrettyString() string {
	texts := make([]string, 0, 4)
	if c.Ironium > 0 {
		texts = append(texts, fmt.Sprintf("%dkT ironium", c.Ironium))
	}
	if c.Boranium > 0 {
		texts = append(texts, fmt.Sprintf("%dkT boranium", c.Boranium))
	}
	if c.Germanium > 0 {
		texts = append(texts, fmt.Sprintf("%dkT germanium", c.Germanium))
	}
	return strings.Join(texts, ", ")
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

// convert a mineral to a cost
func (m Mineral) ToCost() Cost {
	return Cost{
		Ironium:   m.Ironium,
		Boranium:  m.Boranium,
		Germanium: m.Germanium,
	}
}

// add two minerals
func (m Mineral) Add(m2 Mineral) Mineral {
	return Mineral{
		Ironium:   m.Ironium + m2.Ironium,
		Boranium:  m.Boranium + m2.Boranium,
		Germanium: m.Germanium + m2.Germanium,
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

// subtract two minerals
func (m Mineral) Minus(m2 Mineral) Mineral {
	return Mineral{
		Ironium:   m.Ironium - m2.Ironium,
		Boranium:  m.Boranium - m2.Boranium,
		Germanium: m.Germanium - m2.Germanium,
	}
}

// subtract the mineral components of a Cost
func (m Mineral) MinusCost(m2 Cost) Mineral {
	return Mineral{
		Ironium:   m.Ironium - m2.Ironium,
		Boranium:  m.Boranium - m2.Boranium,
		Germanium: m.Germanium - m2.Germanium,
	}
}

func (c Mineral) MultiplyFloat64(factor float64) Mineral {
	return Mineral{
		Ironium:   int(float64(c.Ironium) * factor),
		Boranium:  int(float64(c.Boranium) * factor),
		Germanium: int(float64(c.Germanium) * factor),
	}
}

func (m Mineral) Clamp(min, max int) Mineral {
	return Mineral{
		Ironium:   Clamp(m.Ironium, min, max),
		Boranium:  Clamp(m.Boranium, min, max),
		Germanium: Clamp(m.Germanium, min, max),
	}
}

func (c Mineral) GreatestType() CargoType {
	if c.Ironium >= c.Boranium && c.Ironium >= c.Germanium {
		return Ironium
	}

	if c.Boranium >= c.Ironium && c.Boranium >= c.Germanium {
		return Boranium
	}

	if c.Germanium >= c.Ironium && c.Germanium >= c.Boranium {
		return Germanium
	}

	return None
}
