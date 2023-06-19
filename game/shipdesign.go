package game

import (
	"time"

	"gorm.io/gorm"
)

type ShipDesign struct {
	ID        uint              `gorm:"primaryKey" json:"id" header:"Username"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedat"`
	DeletedAt gorm.DeletedAt    `gorm:"index" json:"deletedAt"`
	GameID    uint              `json:"gameId"`
	PlayerID  uint              `json:"playerId"`
	PlayerNum int               `json:"playerNum"`
	Dirty     bool              `json:"-"`
	Name      string            `json:"name"`
	Version   int               `json:"version"`
	Hull      string            `json:"hull"`
	Armor     int               `json:"armor"`
	Slots     []ShipDesignSlot  `json:"slots" gorm:"serializer:json"`
	Purpose   ShipDesignPurpose `json:"purpose,omitempty"`
	Spec      *ShipDesignSpec   `json:"spec" gorm:"serializer:json"`
}

type ShipDesignSlot struct {
	HullComponent string `json:"hullComponent"`
	HullSlotIndex int    `json:"hullSlotIndex"`
	Quantity      int    `json:"quantity"`
}

type ShipDesignSpec struct {
	IdealSpeed               int                   `json:"idealSpeed,omitempty"`
	WeaponSlots              []ShipDesignSlot      `json:"weaponSlots"`
	Computed                 bool                  `json:"computed"`
	Engine                   string                `json:"engine"`
	NumEngines               int                   `json:"numEngines"`
	Cost                     Cost                  `json:"cost"`
	Mass                     int                   `json:"mass"`
	Armor                    int                   `json:"armor"`
	FuelCapacity             int                   `json:"fuelCapacity"`
	ScanRange                int                   `json:"scanRange"`
	ScanRangePen             int                   `json:"scanRangePen"`
	TorpedoInaccuracyFactor  float64               `json:"torpedoInaccuracyFactor"`
	Initiative               int                   `json:"initiative"`
	Movement                 int                   `json:"movement"`
	Bombs                    []Bomb                `json:"bombs"`
	SmartBombs               []Bomb                `json:"smartBombs"`
	RetroBombs               []Bomb                `json:"retroBombs"`
	Scanner                  bool                  `json:"scanner"`
	MineLayingRateByMineType map[MineFieldType]int `json:"mineLayingRateByMineType"`
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

func computeShipDesignSpec(rules *Rules, player *Player, design *ShipDesign) *ShipDesignSpec {
	hull := rules.Techs.GetHull(design.Hull)
	_ = hull
	spec := ShipDesignSpec{}

	return &spec
}

func designShip(techStore *TechStore, hull *TechHull, name string, player *Player, hullSetNumber int, purpose ShipDesignPurpose) *ShipDesign {

	design := ShipDesign{Name: name, PlayerID: player.ID, PlayerNum: player.Num, Hull: hull.Name, Dirty: true}

	engine := techStore.GetBestEngine(player)
	scanner := techStore.GetBestScanner(player)

	for i, hullSlot := range hull.Slots {
		slot := ShipDesignSlot{HullSlotIndex: i + 1}
		slot.Quantity = hullSlot.Capacity

		switch hullSlot.Type {
		case HullSlotTypeEngine:
			slot.HullComponent = engine.Name
		case HullSlotTypeScanner:
			slot.HullComponent = scanner.Name
		}

		// we filled it, add it
		if slot.HullComponent != "" {
			design.Slots = append(design.Slots, slot)
		}
	}

	return &design
}
