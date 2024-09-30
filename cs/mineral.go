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

func (m Mineral) PrettyString() string {
	texts := make([]string, 0, 4)
	if m.Ironium > 0 {
		texts = append(texts, fmt.Sprintf("%dkT ironium", m.Ironium))
	}
	if m.Boranium > 0 {
		texts = append(texts, fmt.Sprintf("%dkT boranium", m.Boranium))
	}
	if m.Germanium > 0 {
		texts = append(texts, fmt.Sprintf("%dkT germanium", m.Germanium))
	}
	return strings.Join(texts, ", ")
}

func (m *Mineral) Set(mineralType MineralType, value int) *Mineral {
	switch mineralType {
	case Ironium:
		m.Ironium = value
	case Boranium:
		m.Boranium = value
	case Germanium:
		m.Germanium = value
	}
	return m
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


// return the Nth highest MineralType in a Mineral struct 
// (1 = highest, 2 = middle, 3 = lowest)
// 
// Ties are broken in order of precendence (I/B/G) 
func (m Mineral) HighestType(ranking int) MineralType {
	copy := m // make copy of struct so we can zero out values without affecting the original
	var highestType MineralType
	for i := 0; i < ranking; i++ {
		// get the highest type in the cost struct
		highestType = copy.GetTypeFromAmount(MaxInt(copy.Ironium, copy.Boranium, copy.Germanium))
		// For the record, this will never cause GetTypeFromAmount to panic because we are  
		// comparing the struct's own values against themselves
		copy.Set(highestType, 0)
	}
	return highestType
}

// return the first valid MineralType in a Mineral struct with the given numerical value
// returns an error if no MineralType with the corresponding value exists
func (m Mineral) GetTypeFromAmount(amt int) MineralType {
	switch amt {
	case m.Ironium:
		return Ironium
	case m.Boranium:
		return Boranium
	case m.Germanium:
		return Germanium
	}
	panic(fmt.Sprintf("GetTypeFromAmount called with value %v but no corresponding MineralType was found in mineral struct; \nStruct values:\nIronium: %v\nBoranium: %v\nGermanium: %v",
		amt, m.Ironium, m.Boranium, m.Germanium))
}