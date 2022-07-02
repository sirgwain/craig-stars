package game

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Fleet struct {
	MapObject
	PlanetID     uint        `json:"-"` // for starbase fleets that are owned by a planet
	BaseName     string      `json:"baseName"`
	Fuel         int         `json:"fuel"`
	BattlePlanID uint        `json:"battlePlan"`
	DockCapacity int         `json:"dockCapacity"`
	Tokens       []ShipToken `json:"tokens" gorm:"constraint:OnDelete:CASCADE;"`
	Waypoints    []Waypoint  `json:"waypoints" gorm:"serializer:json"`
	Spec         *FleetSpec  `json:"spec" gorm:"serializer:json"`
}

type FleetSpec struct {
	Computed                 bool                  `json:"computed"`
	Purposes                 []string              `json:"purposes"`
	WeaponSlots              []ShipDesignSlot      `json:"weaponSlots"`
	TotalShips               int                   `json:"totalShips"`
	MassEmpty                int                   `json:"massEmpty"`
	BaseCloakedCargo         int                   `json:"baseCloakedCargo"`
	Engine                   string                `json:"engine"`
	NumEngines               int                   `json:"numEngines"`
	Cost                     Cost                  `json:"cost"`
	Mass                     int                   `json:"mass"`
	Armor                    int                   `json:"armor"`
	FuelCapacity             int                   `json:"fuelCapacity"`
	ScanRange                int                   `json:"scanRange"`
	ScanRangePen             int                   `json:"scanRangePen"`
	Bombs                    []Bomb                `json:"bombs"`
	SmartBombs               []Bomb                `json:"smartBombs"`
	RetroBombs               []Bomb                `json:"retroBombs"`
	Scanner                  bool                  `json:"scanner"`
	MineLayingRateByMineType map[MineFieldType]int `json:"mineLayingRateByMineType"`
}

type ShipToken struct {
	ID        uint           `gorm:"primaryKey" json:"id" header:"Username"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedat"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	FleetID   uint           `json:"gameId"`
	Design    *ShipDesign    `json:"-" gorm:"foreignKey:DesignID"`
	DesignID  uint           `json:"designId"`
	Quantity  int            `json:"quantity"`
}

type Waypoint struct {
	FleetID          uint   `json:"-"`
	TargetID         uint   `json:"targetId,omitempty"`
	Position         Vector `json:"position,omitempty" gorm:"embedded"`
	WarpFactor       int    `json:"warpFactor,omitempty"`
	TargetPlanetNum  *int   `json:"targetPlanetNum,omitempty"`
	TransferToPlayer *int   `json:"transferToPlayer,omitempty"`
	TargetName       string `json:"targetName,omitempty"`
}

func NewFleet(player *Player, design *ShipDesign, num int, name string, waypoints []Waypoint) Fleet {
	return Fleet{
		MapObject: MapObject{
			GameID:    player.GameID,
			PlayerID:  player.ID,
			PlayerNum: &player.Num,
			Dirty:     true,
			Num:       num,
			Name:      fmt.Sprintf("%s #%d", name, num),
			Position:  waypoints[0].Position,
		},
		BaseName: name,
		Tokens: []ShipToken{
			{Design: design, Quantity: 1},
		},
		Waypoints: waypoints,
	}
}

func NewPlanetWaypoint(planet *Planet, warpFactor int) Waypoint {
	return Waypoint{
		Position:        planet.Position,
		TargetPlanetNum: &planet.Num,
		TargetName:      planet.Name,
		WarpFactor:      warpFactor,
	}
}
