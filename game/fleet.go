package game

import (
	"time"

	"gorm.io/gorm"
)

type Fleet struct {
	MapObject
	PlanetID     uint        `json:"-"`        // for starbase fleets that are owned by a planet
	BaseName     string      `json:"baseName"`
	Fuel         int         `json:"fuel"`
	BattlePlanID uint        `json:"battlePlan"`
	Waypoints    []Waypoint  `json:"waypoints" gorm:"constraint:OnDelete:CASCADE;"`
	Tokens       []ShipToken `json:"tokens" gorm:"constraint:OnDelete:CASCADE;"`
	DockCapacity int         `json:"dockCapacity"`
	Spec         *FleetSpec  `json:"spec" gorm:"-"`
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
	DesignID  uint           `json:"designId"`
	Quantity  int            `json:"quantity"`
}

type Waypoint struct {
	FleetID          uint   `json:"-"`
	TargetID         uint   `json:"targetId"`
	Position         Vector `json:"position" gorm:"embedded"`
	WarpFactor       int    `json:"warpFactor"`
	TransferToPlayer int    `json:"transferToPlayer"`
	TargetName       string `json:"targetName"`
}
