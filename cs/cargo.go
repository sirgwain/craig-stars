package cs

import (
	"fmt"
	"strings"
)

type Cargo struct {
	Ironium   int `json:"ironium,omitempty"`
	Boranium  int `json:"boranium,omitempty"`
	Germanium int `json:"germanium,omitempty"`
	Colonists int `json:"colonists,omitempty"`
}

type CargoType int

const (
	Ironium CargoType = iota
	Boranium
	Germanium
	Colonists
	Fuel
)

func (c CargoType) String() string {
	switch c {
	case Ironium:
		return "Ironium"
	case Boranium:
		return "Boranium"
	case Germanium:
		return "Germanium"
	case Colonists:
		return "Colonists"
	}
	return ""
}

func (c Cargo) PrettyString() string {
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
	if c.Colonists > 0 {
		texts = append(texts, fmt.Sprintf("%dkT colonists", c.Colonists))
	}
	return strings.Join(texts, ", ")
}

func (c Cargo) HasColonists() bool {
	return c.Colonists > 0
}

func (c Cargo) HasMinerals() bool {
	return (c.Ironium + c.Boranium + c.Germanium) > 0
}

func (c Cargo) Negative() Cargo {
	return Cargo{
		Ironium:   -c.Ironium,
		Boranium:  -c.Boranium,
		Germanium: -c.Germanium,
		Colonists: -c.Colonists,
	}
}

func (c Cargo) Add(other Cargo) Cargo {
	return Cargo{
		Ironium:   c.Ironium + other.Ironium,
		Boranium:  c.Boranium + other.Boranium,
		Germanium: c.Germanium + other.Germanium,
		Colonists: c.Colonists + other.Colonists,
	}
}

func (c Cargo) Subtract(other Cargo) Cargo {
	return Cargo{
		Ironium:   c.Ironium - other.Ironium,
		Boranium:  c.Boranium - other.Boranium,
		Germanium: c.Germanium - other.Germanium,
		Colonists: c.Colonists - other.Colonists,
	}
}

func (c Cargo) AddMineral(other Mineral) Cargo {
	return Cargo{
		Ironium:   c.Ironium + other.Ironium,
		Boranium:  c.Boranium + other.Boranium,
		Germanium: c.Germanium + other.Germanium,
		Colonists: c.Colonists,
	}
}

func (c Cargo) AddCostMinerals(other Cost) Cargo {
	return Cargo{
		Ironium:   c.Ironium + other.Ironium,
		Boranium:  c.Boranium + other.Boranium,
		Germanium: c.Germanium + other.Germanium,
		Colonists: c.Colonists,
	}
}

func (c Cargo) ToMineral() Mineral {
	return Mineral{
		Ironium:   c.Ironium,
		Boranium:  c.Boranium,
		Germanium: c.Germanium,
	}
}

func (c Cargo) Total() int {
	return c.Ironium + c.Boranium + c.Germanium + c.Colonists
}

// return true if this cargo can have transferAmount taken from it
func (c Cargo) CanTransfer(transferAmount Cargo) bool {
	return (c.Ironium >= transferAmount.Ironium &&
		c.Boranium >= transferAmount.Boranium &&
		c.Germanium >= transferAmount.Germanium &&
		c.Colonists >= transferAmount.Colonists)

}

func (c Cargo) CanTransferAmount(cargoType CargoType, transferAmount int) bool {
	switch cargoType {
	case Ironium:
		return c.Ironium >= transferAmount
	case Boranium:
		return c.Boranium >= transferAmount
	case Germanium:
		return c.Germanium >= transferAmount
	case Colonists:
		return c.Colonists >= transferAmount
	}
	return false

}

func (c *Cargo) SubtractAmount(cargoType CargoType, transferAmount int) *Cargo {
	switch cargoType {
	case Ironium:
		c.Ironium -= transferAmount
	case Boranium:
		c.Boranium -= transferAmount
	case Germanium:
		c.Germanium -= transferAmount
	case Colonists:
		c.Colonists -= transferAmount
	}
	return c
}

func (c *Cargo) AddAmount(cargoType CargoType, transferAmount int) *Cargo {
	switch cargoType {
	case Ironium:
		c.Ironium += transferAmount
	case Boranium:
		c.Boranium += transferAmount
	case Germanium:
		c.Germanium += transferAmount
	case Colonists:
		c.Colonists += transferAmount
	}
	return c
}

// get the amount for a type of cargo
func (c Cargo) GetAmount(t CargoType) int {
	switch t {
	case Ironium:
		return c.Ironium
	case Boranium:
		return c.Boranium
	case Germanium:
		return c.Germanium
	case Colonists:
		return c.Colonists
	}
	return 0
}

// get the amount for a type of cargo
func (c Cargo) WithCargo(t CargoType, amount int) Cargo {
	switch t {
	case Ironium:
		c.Ironium = amount
	case Boranium:
		c.Boranium = amount
	case Germanium:
		c.Germanium = amount
	case Colonists:
		c.Colonists = amount
	}
	return c
}

func (c Cargo) WithPopulation(amount int) Cargo {
	c.Colonists = amount / 100
	return c
}
