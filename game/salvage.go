package game

type Salvage struct {
	MapObject
	Cargo Cargo `json:"cargo,omitempty"`
}

// create a new salvage object
func newSalvage(playerNum int, position Vector, cargo Cargo) Salvage {
	return Salvage{
		MapObject: MapObject{
			Type:      MapObjectTypeSalvage,
			Position:  position,
			PlayerNum: playerNum,
		},
		Cargo: cargo,
	}
}
