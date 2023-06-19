package cs

type Salvage struct {
	MapObject
	Cargo Cargo `json:"cargo,omitempty"`
}

// create a new salvage object
func newSalvage(position Vector, num int, playerNum int, cargo Cargo) *Salvage {
	return &Salvage{
		MapObject: MapObject{
			Type:      MapObjectTypeSalvage,
			Position:  position,
			Num:       num,
			PlayerNum: playerNum,
			Dirty:     true,
		},
		Cargo: cargo,
	}
}
