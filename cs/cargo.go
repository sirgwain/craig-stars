package cs

import (
	"fmt"
	"strings"
)

// Cargo represents minerals and colonists that are in cargo holds, salvage, mineral packets, or on planets.
// 1 kT of Cargo respresents 1 unit of minerals or 100 colonists.
type Cargo struct {
	Ironium   int `json:"ironium,omitempty"`
	Boranium  int `json:"boranium,omitempty"`
	Germanium int `json:"germanium,omitempty"`
	Colonists int `json:"colonists,omitempty"`
}

// Create a new Cargo struct from a Mineral struct and population amount.
// Assumes pop is already in kT.
func NewCargoFromMineral(mineral Mineral, pop int) Cargo {
	return Cargo{
		Ironium:   mineral.Ironium,
		Boranium:  mineral.Ironium,
		Germanium: mineral.Ironium,
		Colonists: pop,
	}
}

func NewCargoFromType(cargoType CargoType, amount int) Cargo {
	c := Cargo{}
	switch cargoType {
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

type CargoType = ResourceType

var CargoTypes = [4]CargoType{
	Ironium,
	Boranium,
	Germanium,
	Colonists,
}

func NewCargoFromMineralsAndPop(mineral Mineral, pop int) Cargo {
	return Cargo{
		Ironium:   mineral.Ironium,
		Boranium:  mineral.Boranium,
		Germanium: mineral.Germanium,
		Colonists: pop / 100,
	}
}

func NewCargoFromArray(values [4]int) Cargo {
	return Cargo{
		Ironium:   values[0],
		Boranium:  values[1],
		Germanium: values[2],
		Colonists: values[3],
	}
}

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
	if c.Ironium != 0 {
		texts = append(texts, fmt.Sprintf("%dkT ironium", c.Ironium))
	}
	if c.Boranium != 0 {
		texts = append(texts, fmt.Sprintf("%dkT boranium", c.Boranium))
	}
	if c.Germanium != 0 {
		texts = append(texts, fmt.Sprintf("%dkT germanium", c.Germanium))
	}
	if c.Colonists != 0 {
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

// HasNegative returns true if any cargo is negative
func (c Cargo) HasNegative() bool {
	return c.Ironium < 0 || c.Boranium < 0 || c.Germanium < 0 || c.Colonists < 0
}

// HasPositive returns true if any cargo is positive
func (c Cargo) HasPositive() bool {
	return c.Ironium > 0 || c.Boranium > 0 || c.Germanium > 0 || c.Colonists > 0
}

// return this cargo with a minimum of zero for each value
func (c Cargo) MinZero() Cargo {
	return Cargo{
		Ironium:   max(c.Ironium, 0),
		Boranium:  max(c.Boranium, 0),
		Germanium: max(c.Germanium, 0),
		Colonists: max(c.Colonists, 0),
	}
}

// NegativeOnly returns a cargo with only negative values
// used for identifying stealing cargo
func (c Cargo) NegativeOnly() Cargo {
	return Cargo{
		Ironium:   min(c.Ironium, 0),
		Boranium:  min(c.Boranium, 0),
		Germanium: min(c.Germanium, 0),
		Colonists: min(c.Colonists, 0),
	}
}

// PositiveOnly returns a cargo with only negative values
// used for identifying unloading cargo
func (c Cargo) PositiveOnly() Cargo {
	return Cargo{
		Ironium:   max(c.Ironium, 0),
		Boranium:  max(c.Boranium, 0),
		Germanium: max(c.Germanium, 0),
		Colonists: max(c.Colonists, 0),
	}
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

func (c Cargo) Subtract(other Cargo) Cargo {
	return Cargo{
		Ironium:   c.Ironium - other.Ironium,
		Boranium:  c.Boranium - other.Boranium,
		Germanium: c.Germanium - other.Germanium,
		Colonists: c.Colonists - other.Colonists,
	}
}

func (c Cargo) Multiply(product float64) Cargo {
	return Cargo{
		int(float64(c.Ironium) * product),
		int(float64(c.Boranium) * product),
		int(float64(c.Germanium) * product),
		int(float64(c.Colonists) * product),
	}
}

func (c Cargo) ToMineral() Mineral {
	return Mineral{
		Ironium:   c.Ironium,
		Boranium:  c.Boranium,
		Germanium: c.Germanium,
	}
}

func (c Cargo) ToCost() Cost {
	return Cost{
		Ironium:   c.Ironium,
		Boranium:  c.Boranium,
		Germanium: c.Germanium,
	}
}

func (c Cargo) Total() int {
	return c.Ironium + c.Boranium + c.Germanium + c.Colonists
}

func (c Cargo) absSum() int {
	return Abs(c.Ironium) + Abs(c.Boranium) + Abs(c.Germanium) + Abs(c.Colonists)
}

func (c Cargo) ToArray() [4]int {
	return [4]int{
		c.Ironium,
		c.Boranium,
		c.Germanium,
		c.Colonists,
	}
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

func (c Cargo) SubtractAmount(cargoType CargoType, transferAmount int) Cargo {
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

func (c Cargo) AddAmount(cargoType CargoType, transferAmount int) Cargo {
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

func (c Cargo) SetAmount(t CargoType, amount int) Cargo {
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

// TODO: Remove this in favor of simple assignment (this just seems dumb lol)
func (c Cargo) WithPopulation(amount int) Cargo {
	c.Colonists = amount / 100
	return c
}

// return the mineral with the highest amount
func (c Cargo) GreatestMineralType() CargoType {
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

// split a cargo into two cargos based on capacity
func (source Cargo) Split(sourceCapacity, capacity1, capacity2 int) (Cargo, Cargo, error) {
	sourceArray := source.ToArray()
	split1, split2, err := splitValues(sourceCapacity, capacity1, capacity2, (sourceArray[:])...)
	if err != nil {
		return Cargo{}, Cargo{}, err
	}

	return NewCargoFromArray([4]int(split1)), NewCargoFromArray([4]int(split2)), nil
}

// Set the mineral portions of this Cargo, leaving resources unaffected.
func (c *Cargo) SetMineral(mineral Mineral) {
	c.Ironium = mineral.Ironium
	c.Boranium = mineral.Boranium
	c.Germanium = mineral.Germanium
}
