package game

import (
	"time"

	"gorm.io/gorm"
)

type ShipDesign struct {
	ID        uint             `gorm:"primaryKey" json:"id" header:"Username"`
	CreatedAt time.Time        `json:"createdAt"`
	UpdatedAt time.Time        `json:"updatedat"`
	DeletedAt gorm.DeletedAt   `gorm:"index" json:"deletedAt"`
	GameID    uint             `json:"gameId"`
	PlayerID  uint             `json:"playerId"`
	Dirty     bool             `json:"-"`
	Name      string           `json:"name"`
	Version   int              `json:"version"`
	PlayerNum int              `json:"playerNum"`
	Hull      string           `json:"hull"`
	Slots     []ShipDesignSlot `json:"slots"`
	Spec      *ShipDesignSpec  `json:"spec" gorm:"-"`
	Armor     int              `json:"armor"`
}

type ShipDesignSlot struct {
	ID            uint           `gorm:"primaryKey" json:"id" header:"Username"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedat"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	ShipDesignID  uint           `json:"-"`
	HullComponent string         `json:"hullComponent"`
	HullSlotIndex int            `json:"hullSlotIndex"`
	Quantity      int            `json:"quantity"`
}

type ShipDesignSpec struct {
	WeaponSlots              []ShipDesignSlot                     `json:"weaponSlots"`
	Computed                 bool                                 `json:"computed"`
	Engine                   string                               `json:"engine"`
	NumEngines               int                                  `json:"numEngines"`
	Cost                     Cost                                 `json:"cost"`
	Mass                     int                                  `json:"mass"`
	Armor                    int                                  `json:"armor"`
	FuelCapacity             int                                  `json:"fuelCapacity"`
	ScanRange                int                                  `json:"scanRange"`
	ScanRangePen             int                                  `json:"scanRangePen"`
	TorpedoInaccuracyFactor  float64                              `json:"torpedoInaccuracyFactor"`
	Initiative               int                                  `json:"initiative"`
	Movement                 int                                  `json:"movement"`
	Bombs                    []Bomb                               `json:"bombs"`
	SmartBombs               []Bomb                               `json:"smartBombs"`
	RetroBombs               []Bomb                               `json:"retroBombs"`
	Scanner                  bool                                 `json:"scanner"`
	MineLayingRateByMineType map[MineLayingRateByMineType]float64 `json:"mineLayingRateByMineType"`
}

type MineLayingRateByMineType struct {
}

type ShipDesignPurpose string

const (
	ShipDesignPurposeScout             ShipDesignPurpose = "Scout"
	ShipDesignPurposeArmedScout        ShipDesignPurpose = "ArmedScout"
	ShipDesignPurposeColonizer         ShipDesignPurpose = "Colonizer"
	ShipDesignPurposeBomber            ShipDesignPurpose = "Bomber"
	ShipDesignPurposeFighter           ShipDesignPurpose = "Fighter"
	ShipDesignPurposeFighterScout      ShipDesignPurpose = "FighterScout"
	ShipDesignPurposeCapitalShip       ShipDesignPurpose = "CapitalShip"
	ShipDesignPurposeFreighter         ShipDesignPurpose = "Freighter"
	ShipDesignPurposeColonistFreighter ShipDesignPurpose = "ColonistFreighter"
	ShipDesignPurposeFuelFreighter     ShipDesignPurpose = "FuelFreighter"
	ShipDesignPurposeArmedFreighter    ShipDesignPurpose = "ArmedFreighter"
	ShipDesignPurposeMiner             ShipDesignPurpose = "Miner"
	ShipDesignPurposeTerraformer       ShipDesignPurpose = "Terraformer"
	ShipDesignPurposeDamageMineLayer   ShipDesignPurpose = "DamageMineLayer"
	ShipDesignPurposeSpeedMineLayer    ShipDesignPurpose = "SpeedMineLayer"
	ShipDesignPurposeStarbase          ShipDesignPurpose = "Starbase"
	ShipDesignPurposeFort              ShipDesignPurpose = "Fort"
	ShipDesignPurposeStarterColony     ShipDesignPurpose = "StarterColony"
)

func NewShipDesign(player *Player) *ShipDesign {
	return &ShipDesign{GameID: player.GameID, PlayerID: player.ID, PlayerNum: player.Num, Dirty: true}
}
