package game

import (
	"time"

	"gorm.io/gorm"
)

type TechStore struct {
	ID                uint                   `gorm:"primaryKey" json:"id" header:"Username"`
	CreatedAt         time.Time              `json:"createdAt"`
	UpdatedAt         time.Time              `json:"updatedat"`
	DeletedAt         gorm.DeletedAt         `gorm:"index" json:"deletedAt"`
	RulesID           uint                   `json:"rulesId"`
	Engines           []TechEngine           `json:"engines"`
	PlanetaryScanners []TechPlanetaryScanner `json:"planetaryScanners"`
}

// simple static tech store
var StaticTechStore = TechStore{
	Engines:           TechEngines(),
	PlanetaryScanners: TechPlanetaryScanners(),
}

type TechFinder interface {
	GetBestPlanetaryScanner(player *Player) *TechPlanetaryScanner
}

func NewTechStore() TechFinder {
	return &StaticTechStore
}

// get the best planetary scanner for a player
func (store *TechStore) GetBestPlanetaryScanner(player *Player) *TechPlanetaryScanner {
	bestTech := &store.PlanetaryScanners[0]
	for i := range store.PlanetaryScanners {
		tech := &store.PlanetaryScanners[i]
		if player.HasTech(&tech.Tech) {
			bestTech = tech
		}
	}
	return bestTech
}

// TechEngines
var SettlersDelight = TechEngine{
	TechHullComponent: TechHullComponent{Tech: NewTech("Settler's Delight", NewCost(1, 0, 1, 2), TechRequirements{PRTRequired: HE}, 10, TechCategoryEngine), Mass: 2},
	IdealSpeed:        6,
	FreeSpeed:         6,
	FuelUsage: [11]int{
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		150,
		275,
		480,
		576,
	},
}
var QuickJump5 = TechEngine{
	TechHullComponent: TechHullComponent{Tech: NewTech("Quick Jump 5", NewCost(3, 0, 1, 3), TechRequirements{}, 20, TechCategoryEngine), Mass: 4},
	IdealSpeed:        5,
	FuelUsage: [11]int{
		0,    // 0
		0,    // 1
		25,   // 2
		100,  // 3
		100,  // 4
		100,  // 5
		180,  // 6
		500,  // 7
		800,  // 8
		900,  // 9
		1080, // 10
	},
}
var LongHump6 = TechEngine{
	TechHullComponent: TechHullComponent{Tech: NewTech("Long Hump 6", NewCost(5, 0, 1, 6), TechRequirements{TechLevel: TechLevel{Propulsion: 3}}, 30, TechCategoryEngine), Mass: 9},
	IdealSpeed:        6,
	FuelUsage: [11]int{
		0,    // 0
		0,    // 1
		20,   // 2
		60,   // 3
		100,  // 4
		100,  // 5
		105,  // 6
		450,  // 7
		750,  // 8
		900,  // 9
		1080, // 10
	},
}
var FuelMizer = TechEngine{
	TechHullComponent: TechHullComponent{Tech: NewTech("Fuel Mizer", NewCost(8, 0, 0, 11), TechRequirements{TechLevel: TechLevel{Propulsion: 2}, LRTsRequired: IFE}, 40, TechCategoryEngine), Mass: 6},
	IdealSpeed:        6,
	FreeSpeed:         4,
	FuelUsage: [11]int{
		0,   // 0
		0,   // 1
		0,   // 2
		0,   // 3
		0,   // 4
		35,  // 5
		120, // 6
		175, // 7
		235, // 8
		360, // 9
		420, // 10
	},
}

var DaddyLongLegs7 = TechEngine{
	TechHullComponent: TechHullComponent{Tech: NewTech("Daddy Long Legs 7", NewCost(11, 0, 3, 12), TechRequirements{TechLevel: TechLevel{Propulsion: 5}}, 50, TechCategoryEngine), Mass: 13},
	IdealSpeed:        7,
	FuelUsage: [11]int{
		0,   // 0
		0,   // 1
		20,  // 2
		60,  // 3
		70,  // 4
		100, // 5
		100, // 6
		110, // 7
		600, // 8
		750, // 9
		900, // 10
	},
}
var AlphaDrive8 = TechEngine{
	TechHullComponent: TechHullComponent{Tech: NewTech("Alpha Drive 8", NewCost(16, 0, 3, 28), TechRequirements{TechLevel: TechLevel{Propulsion: 7}}, 60, TechCategoryEngine), Mass: 17},
	IdealSpeed:        8,
	FuelUsage: [11]int{
		0,
		0,
		15,
		50,
		60,
		70,
		100,
		100,
		115,
		700,
		840,
	},
}
var TransGalacticDrive = TechEngine{
	TechHullComponent: TechHullComponent{Tech: NewTech("Galactic Drive", NewCost(20, 20, 9, 50), TechRequirements{TechLevel: TechLevel{Propulsion: 9}}, 70, TechCategoryEngine), Mass: 25},
	IdealSpeed:        9,
	FuelUsage: [11]int{
		0,
		0,
		15,
		35,
		45,
		55,
		70,
		80,
		90,
		100,
		120,
	},
}
var Interspace10 = TechEngine{
	TechHullComponent: TechHullComponent{Tech: NewTech("Interspace-10", NewCost(18, 25, 10, 60), TechRequirements{TechLevel: TechLevel{Propulsion: 11}, LRTsRequired: NRSE}, 80, TechCategoryEngine), Mass: 25},
	IdealSpeed:        10,
	FuelUsage: [11]int{
		0,
		0,
		10,
		30,
		40,
		50,
		60,
		70,
		80,
		90,
		100,
	},
}
var TransStar10 = TechEngine{
	TechHullComponent: TechHullComponent{Tech: NewTech("Trans-Star 10", NewCost(3, 0, 3, 10), TechRequirements{TechLevel: TechLevel{Propulsion: 23}}, 90, TechCategoryEngine), Mass: 5},
	IdealSpeed:        10,
	FuelUsage: [11]int{
		0,
		0,
		5,
		15,
		20,
		25,
		30,
		35,
		40,
		45,
		50,
	},
}
var RadiatingHydroRamScoop = TechEngine{
	TechHullComponent: TechHullComponent{Tech: NewTech("Radiating Hydro-Ram Scoop", NewCost(3, 2, 9, 8), TechRequirements{TechLevel: TechLevel{Energy: 2, Propulsion: 6}, LRTsDenied: NRSE}, 100, TechCategoryEngine), Mass: 10, Radiating: true},
	IdealSpeed:        6,
	FreeSpeed:         6,
	FuelUsage: [11]int{
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		165,
		375,
		600,
		720,
	},
}
var SubGalacticFuelScoop = TechEngine{
	TechHullComponent: TechHullComponent{Tech: NewTech("Sub-Galactic Fuel Scoop", NewCost(4, 4, 7, 12), TechRequirements{TechLevel: TechLevel{Energy: 2, Propulsion: 8}, LRTsDenied: NRSE}, 110, TechCategoryEngine), Mass: 20},
	IdealSpeed:        7,
	FreeSpeed:         7,
	FuelUsage: [11]int{
		0,
		0,
		0,
		0,
		0,
		0,
		85,
		105,
		210,
		380,
		456,
	},
}
var TransGalacticFuelScoop = TechEngine{
	TechHullComponent: TechHullComponent{Tech: NewTech("Trans-Galactic Fuel Scoop", NewCost(5, 4, 12, 18), TechRequirements{TechLevel: TechLevel{Energy: 3, Propulsion: 9}, LRTsDenied: NRSE}, 120, TechCategoryEngine), Mass: 19},
	IdealSpeed:        8,
	FreeSpeed:         8,
	FuelUsage: [11]int{
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		88,
		100,
		145,
		174,
	},
}
var TransGalacticSuperScoop = TechEngine{
	TechHullComponent: TechHullComponent{Tech: NewTech("Trans-Galactic Super Scoop", NewCost(6, 4, 16, 24), TechRequirements{TechLevel: TechLevel{Energy: 4, Propulsion: 12}, LRTsDenied: NRSE}, 130, TechCategoryEngine), Mass: 18},
	IdealSpeed:        9,
	FreeSpeed:         9,
	FuelUsage: [11]int{
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		65,
		90,
		108,
	},
}
var TransGalacticMizerScoop = TechEngine{
	TechHullComponent: TechHullComponent{Tech: NewTech("Trans-Galactic Mizer Scoop", NewCost(5, 2, 13, 11), TechRequirements{TechLevel: TechLevel{Energy: 4, Propulsion: 16}, LRTsDenied: NRSE}, 140, TechCategoryEngine), Mass: 11},
	IdealSpeed:        10,
	FreeSpeed:         10,
	FuelUsage: [11]int{
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		70,
		84,
	},
}
var GalaxyScoop = TechEngine{
	TechHullComponent: TechHullComponent{Tech: NewTech("Galaxy Scoop", NewCost(4, 2, 9, 12), TechRequirements{TechLevel: TechLevel{Energy: 5, Propulsion: 20}, LRTsRequired: IFE, LRTsDenied: NRSE}, 150, TechCategoryEngine), Mass: 8},
	IdealSpeed:        10,
	FreeSpeed:         10,
	FuelUsage: [11]int{
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		60,
	},
}

// TechPlanetaryScanners

var Viewer50 = TechPlanetaryScanner{Tech: NewTech("Viewer 50", NewCost(10, 10, 70, 100), TechRequirements{PRTDenied: AR}, 0, TechCategoryPlanetaryScanner),

	ScanRange:    50,
	ScanRangePen: 0,
}
var Viewer90 = TechPlanetaryScanner{Tech: NewTech("Viewer 90", NewCost(10, 10, 70, 100), TechRequirements{TechLevel: TechLevel{Electronics: 1}, PRTDenied: AR}, 1, TechCategoryPlanetaryScanner),

	ScanRange:    90,
	ScanRangePen: 0,
}
var Scoper150 = TechPlanetaryScanner{Tech: NewTech("Scoper 150", NewCost(10, 10, 70, 100), TechRequirements{TechLevel: TechLevel{Electronics: 3}, PRTDenied: AR}, 30, TechCategoryPlanetaryScanner),

	ScanRange:    150,
	ScanRangePen: 0,
}
var Scoper220 = TechPlanetaryScanner{Tech: NewTech("Scoper 220", NewCost(10, 10, 70, 100), TechRequirements{TechLevel: TechLevel{Electronics: 6}, PRTDenied: AR}, 40, TechCategoryPlanetaryScanner),

	ScanRange:    220,
	ScanRangePen: 0,
}
var Scoper280 = TechPlanetaryScanner{Tech: NewTech("Scoper 280", NewCost(10, 10, 70, 100), TechRequirements{TechLevel: TechLevel{Electronics: 8}, PRTDenied: AR}, 50, TechCategoryPlanetaryScanner),

	ScanRange:    280,
	ScanRangePen: 0,
}
var Snooper320X = TechPlanetaryScanner{Tech: NewTech("Snooper 320X", NewCost(10, 10, 70, 100), TechRequirements{TechLevel: TechLevel{Energy: 3, Electronics: 10, Biotechnology: 3}, PRTDenied: AR, LRTsDenied: NAS}, 60, TechCategoryPlanetaryScanner),

	ScanRange:    320,
	ScanRangePen: 160,
}
var Snooper400X = TechPlanetaryScanner{Tech: NewTech("Snooper 400X", NewCost(10, 10, 70, 100), TechRequirements{TechLevel: TechLevel{Energy: 4, Electronics: 13, Biotechnology: 6}, PRTDenied: AR, LRTsDenied: NAS}, 70, TechCategoryPlanetaryScanner),

	ScanRange:    400,
	ScanRangePen: 200,
}
var Snooper500X = TechPlanetaryScanner{Tech: NewTech("Snooper 500X", NewCost(10, 10, 70, 100), TechRequirements{TechLevel: TechLevel{Energy: 5, Electronics: 16, Biotechnology: 7}, PRTDenied: AR, LRTsDenied: NAS}, 80, TechCategoryPlanetaryScanner),

	ScanRange:    500,
	ScanRangePen: 250,
}
var Snooper620X = TechPlanetaryScanner{Tech: NewTech("Snooper 620X", NewCost(10, 10, 70, 100), TechRequirements{TechLevel: TechLevel{Energy: 7, Electronics: 23, Biotechnology: 9}, PRTDenied: AR, LRTsDenied: NAS}, 90, TechCategoryPlanetaryScanner),

	ScanRange:    620,
	ScanRangePen: 310,
}

func TechEngines() []TechEngine {
	return []TechEngine{
		SettlersDelight,
		QuickJump5,
		LongHump6,
		FuelMizer,
		DaddyLongLegs7,
		AlphaDrive8,
		TransGalacticDrive,
		Interspace10,
		TransStar10,
		RadiatingHydroRamScoop,
		SubGalacticFuelScoop,
		TransGalacticFuelScoop,
		TransGalacticSuperScoop,
		TransGalacticMizerScoop,
		GalaxyScoop,
	}
}

func TechPlanetaryScanners() []TechPlanetaryScanner {

	return []TechPlanetaryScanner{
		Viewer50,
		Viewer90,
		Scoper150,
		Scoper220,
		Scoper280,
		Snooper320X,
		Snooper400X,
		Snooper500X,
		Snooper620X,
	}
}
