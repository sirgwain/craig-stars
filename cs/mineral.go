package cs

import (
	"fmt"
	"strings"
)

// Minerals are mined from planets and turned into Cargo
type Mineral struct {
	Ironium   int `json:"ironium,omitempty"`
	Boranium  int `json:"boranium,omitempty"`
	Germanium int `json:"germanium,omitempty"`
}

type MineralType = ResourceType

var MineralTypes = [3]MineralType{
	Ironium,
	Boranium,
	Germanium,
}

func NewMineral(values [3]int) Mineral {
	return Mineral{
		Ironium:   values[0],
		Boranium:  values[1],
		Germanium: values[2],
	}
}

func (m Mineral) String() string {
	return fmt.Sprintf("Ironium: %d, Boranium: %d, Germanium: %d", m.Ironium, m.Boranium, m.Germanium)
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

func (h *Mineral) Set(mineralType MineralType, value int) *Mineral {
	switch mineralType {
	case Ironium:
		h.Ironium = value
	case Boranium:
		h.Boranium = value
	case Germanium:
		h.Germanium = value
	}
	return h
}

// return higher of 2 Mineral structs for all MineralTypes separately
func (m Mineral) Max(other Mineral) Mineral {
	return Mineral{
		Ironium:   MaxInt(m.GetAmount(Ironium), other.GetAmount(Ironium)),
		Boranium:  MaxInt(m.GetAmount(Boranium), other.GetAmount(Boranium)),
		Germanium: MaxInt(m.GetAmount(Germanium), other.GetAmount(Germanium)),
	}
}

func (m Mineral) GetAmount(mineralType MineralType) int {
	var amt int
	switch mineralType {
	case Ironium:
		amt = m.Ironium
	case Boranium:
		amt = m.Boranium
	case Germanium:
		amt = m.Germanium
	}
	return amt
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

func (m Mineral) HighestType() MineralType {
	if m.Ironium >= m.Boranium && m.Ironium >= m.Germanium {
		return Ironium
	}

	if m.Boranium >= m.Ironium && m.Boranium >= m.Germanium {
		return Boranium
	}

	if m.Germanium >= m.Ironium && m.Germanium >= m.Boranium {
		return Germanium
	}

	return None
}

// returns 2nd lowest/highest mineral type
func (m Mineral) MiddleType() MineralType {
	if Boranium != m.HighestType() && Boranium != m.LowestType() {
		return Boranium
	}

	if Germanium != m.HighestType() && Germanium != m.LowestType() {
		return Germanium
	}

	if Ironium != m.HighestType() && Ironium != m.LowestType() {
		return Ironium
	}

	return None
}

func (m Mineral) LowestType() MineralType {
	if m.Germanium <= m.Ironium && m.Germanium <= m.Boranium {
		return Germanium
	}

	if m.Boranium <= m.Ironium && m.Boranium <= m.Germanium {
		return Boranium
	}

	if m.Ironium <= m.Germanium && m.Ironium <= m.Boranium {
		return Ironium
	}

	return None
}
