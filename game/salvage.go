package game

type Salvage struct {
	MapObject
	Cargo Cargo `json:"cargo,omitempty" gorm:"embedded;embeddedPrefix:cargo_"`
}

// create a new salvage object
func NewSalvage(playerNum int, position Vector, cargo Cargo) Salvage {
	return Salvage{
		MapObject: MapObject{
			Type:      MapObjectTypeSalvage,
			Position:  position,
			PlayerNum: playerNum,
		},
		Cargo: cargo,
	}
}
