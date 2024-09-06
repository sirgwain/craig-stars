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

type MineralType int

var MineralTypes = [3]MineralType{
	MineralType(Ironium),
	MineralType(Boranium),
	MineralType(Germanium),
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
	case MineralType(Ironium):
		h.Ironium = value
	case MineralType(Boranium):
		h.Boranium = value
	case MineralType(Germanium):
		h.Germanium = value
	}
	return h
}

func (m Mineral) GetAmount(mineralType MineralType) int {
	var amt int
	switch mineralType {
	case MineralType(Ironium):
		amt = m.Ironium
	case MineralType(Boranium):
		amt = m.Boranium
	case MineralType(Germanium):
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

func (m Mineral) GreatestType() MineralType {
	if m.Ironium >= m.Boranium && m.Ironium >= m.Germanium {
		return MineralType(Ironium)
	}

	if m.Boranium >= m.Ironium && m.Boranium >= m.Germanium {
		return MineralType(Boranium)
	}

	if m.Germanium >= m.Ironium && m.Germanium >= m.Boranium {
		return MineralType(Germanium)
	}

	return None
}

// returns 2nd lowest/highest mineral type
func (m Mineral) MiddleType() MineralType {
	if MineralType(Boranium) != m.GreatestType() && MineralType(Boranium) != m.LowestType() {
		return MineralType(Boranium)
	}	
	
	if MineralType(Germanium) != m.GreatestType() && MineralType(Germanium) != m.LowestType() {
		return MineralType(Germanium)
	}

	if MineralType(Ironium) != m.GreatestType() && MineralType(Ironium) != m.LowestType() {
		return MineralType(Ironium)
	}

	return None
}

func (m Mineral) LowestType() MineralType {
	if m.Germanium <= m.Ironium && m.Germanium <= m.Boranium {
		return MineralType(Germanium)
	}

	if m.Boranium <= m.Ironium && m.Boranium <= m.Germanium {
		return MineralType(Boranium)
	}

	if m.Ironium <= m.Germanium && m.Ironium <= m.Boranium {
		return MineralType(Ironium)
	}

	return None
}
