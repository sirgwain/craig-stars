package cs

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

func HabFromInts(hab [3]int) Hab {
	return Hab{
		Grav: hab[0],
		Temp: hab[1],
		Rad:  hab[2],
	}
}
