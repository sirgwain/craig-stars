package cs

import (
	"fmt"
	"math"
)

// Represents a habitability point for a planet. Values range from 1 to 99
type Hab struct {
	Grav int `json:"grav,omitempty"`
	Temp int `json:"temp,omitempty"`
	Rad  int `json:"rad,omitempty"`
}

type HabType int

const (
	Grav HabType = iota
	Temp
	Rad
)

var HabTypes = [3]HabType{
	Grav,
	Temp,
	Rad,
}

func (h HabType) String() string {
	switch h {
	case Grav:
		return "Gravity"
	case Temp:
		return "Temperature"
	case Rad:
		return "Radiation"
	default:
		return fmt.Sprintf("Unknown HabType (%d)", h)
	}
}

func (h Hab) String() string {
	return fmt.Sprintf("Grav: %d, Temp: %d, Rad: %d", h.Grav, h.Temp, h.Rad)
}

func (h Hab) Add(other Hab) Hab {
	return Hab{
		Grav: h.Grav + other.Grav,
		Temp: h.Temp + other.Temp,
		Rad:  h.Rad + other.Rad,
	}
}

func (h Hab) Subtract(other Hab) Hab {
	return Hab{
		Grav: h.Grav - other.Grav,
		Temp: h.Temp - other.Temp,
		Rad:  h.Rad - other.Rad,
	}
}

func (h Hab) Clamp(min, max int) Hab {
	return Hab{
		Grav: Clamp(h.Grav, min, max),
		Temp: Clamp(h.Temp, min, max),
		Rad:  Clamp(h.Rad, min, max),
	}
}

func (h Hab) Get(habType HabType) int {
	switch habType {
	case Grav:
		return h.Grav
	case Temp:
		return h.Temp
	case Rad:
		return h.Rad
	}
	return 0
}

func (h *Hab) Set(habType HabType, value int) *Hab {
	switch habType {
	case Grav:
		h.Grav = value
	case Temp:
		h.Temp = value
	case Rad:
		h.Rad = value
	}
	return h
}

func HabFromInts(hab [3]int) Hab {
	return Hab{
		Grav: hab[0],
		Temp: hab[1],
		Rad:  hab[2],
	}
}

func (h Hab) absSum() int {
	return AbsInt(h.Grav) + AbsInt(h.Temp) + AbsInt(h.Rad)
}

func gravString(grav int) string {
	result := 0
	tmp := int(math.Abs(float64(grav - 50)))
	if tmp <= 25 {
		result = (tmp + 25) * 4
	} else {
		result = tmp*24 - 400
	}
	if grav < 50 {
		result = 10000 / result
	}

	value := float64(result)/100 + float64(result%100)/100.0
	return fmt.Sprintf("%.2fg", value)
}

func tempString(temp int) string {
	return fmt.Sprintf("%d°C", (temp-50)*4)
}

func radString(rad int) string {
	return fmt.Sprintf("%dmR", rad)
}
