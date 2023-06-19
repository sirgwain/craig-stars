package cs

import (
	"math"
	"time"
)

const ScanWithZeroRange = 1
const UnlimitedSpaceDock = -1
const NoScanner = -1
const NoGate = -1
const InfinteGate = math.MaxInt32

type TechStore struct {
	ID                       int64                                 `json:"id"`
	CreatedAt                time.Time                             `json:"createdAt"`
	UpdatedAt                time.Time                             `json:"updatedAt"`
	RulesID                  int64                                 `json:"rulesId"`
	Engines                  []TechEngine                          `json:"engines"`
	PlanetaryScanners        []TechPlanetaryScanner                `json:"planetaryScanners"`
	Terraforms               []TechTerraform                       `json:"terraforms"`
	Defenses                 []TechDefense                         `json:"defenses"`
	HullComponents           []TechHullComponent                   `json:"hullComponents"`
	Hulls                    []TechHull                            `json:"hulls,omitempty"`
	techsByName              map[string]interface{}                `json:"-"`
	hullComponentsByName     map[string]*TechHullComponent         `json:"-"`
	hullsByName              map[string]*TechHull                  `json:"-"`
	enginesByName            map[string]*TechEngine                `json:"-"`
	hullComponetnsByCategory map[TechCategory][]*TechHullComponent `json:"-"`
}

// simple static tech store
var StaticTechStore = TechStore{
	Engines:           TechEngines(),
	PlanetaryScanners: TechPlanetaryScanners(),
	Terraforms:        TechTerraforms(),
	Defenses:          TechDefenses(),
	HullComponents:    TechHullComponents(),
	Hulls:             TechHulls(),
}

func init() {
	StaticTechStore.Init()
}

type TechFinder interface {
	GetBestPlanetaryScanner(player *Player) *TechPlanetaryScanner
	GetBestDefense(player *Player) *TechDefense
	GetBestTerraform(player *Player, terraformHabType TerraformHabType) *TechTerraform
	GetBestScanner(player *Player) *TechHullComponent
	GetBestEngine(player *Player) *TechEngine
	GetEngine(name string) *TechEngine
	GetTech(name string) interface{}
	GetHull(name string) *TechHull
	GetHullComponent(name string) *TechHullComponent
	GetHullComponentsByCategory(category TechCategory) []TechHullComponent
}

func NewTechStore() TechFinder {
	store := &StaticTechStore

	store.Init()
	return store
}

func (store *TechStore) Init() {
	store.techsByName = make(map[string]interface{},
		len(store.Engines)+
			len(store.PlanetaryScanners)+
			len(store.Terraforms)+
			len(store.Defenses)+
			len(store.HullComponents)+
			len(store.Hulls))
	store.hullsByName = make(map[string]*TechHull, len(store.Hulls))
	store.enginesByName = make(map[string]*TechEngine, len(store.Engines))
	store.hullComponentsByName = make(map[string]*TechHullComponent, len(store.Engines)+len(store.HullComponents))
	store.hullComponetnsByCategory = map[TechCategory][]*TechHullComponent{}

	for i := range store.Hulls {
		tech := &store.Hulls[i]
		store.techsByName[tech.Name] = tech
		store.hullsByName[tech.Name] = tech
	}

	for i := range store.Engines {
		tech := &store.Engines[i]
		store.techsByName[tech.Name] = tech
		store.enginesByName[tech.Name] = tech
		store.hullComponentsByName[tech.Name] = &tech.TechHullComponent
	}

	for i := range store.HullComponents {
		tech := &store.HullComponents[i]
		store.techsByName[tech.Name] = tech
		store.hullComponentsByName[tech.Name] = tech

		if _, ok := store.hullComponetnsByCategory[tech.Category]; !ok {
			store.hullComponetnsByCategory[tech.Category] = []*TechHullComponent{}
		}
		store.hullComponetnsByCategory[tech.Category] = append(store.hullComponetnsByCategory[tech.Category], tech)
	}

	for i := range store.PlanetaryScanners {
		tech := &store.PlanetaryScanners[i]
		store.techsByName[tech.Name] = tech
	}

	for i := range store.Terraforms {
		tech := &store.Terraforms[i]
		store.techsByName[tech.Name] = tech
	}

	for i := range store.Defenses {
		tech := &store.Defenses[i]
		store.techsByName[tech.Name] = tech
	}
}

func (store *TechStore) GetTech(name string) interface{} {
	return store.techsByName[name]
}

func (store *TechStore) GetEngine(name string) *TechEngine {
	return store.enginesByName[name]
}

func (store *TechStore) GetHull(name string) *TechHull {
	return store.hullsByName[name]
}

func (store *TechStore) GetHullComponent(name string) *TechHullComponent {
	return store.hullComponentsByName[name]
}

// get all techs by category
func (store *TechStore) GetHullComponentsByCategory(category TechCategory) []TechHullComponent {
	techs := []TechHullComponent{}
	for _, tech := range store.HullComponents {
		if tech.Category == category {
			techs = append(techs, tech)
		}
	}

	return techs
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

// get the best defense for a player
func (store *TechStore) GetBestDefense(player *Player) *TechDefense {
	bestTech := &store.Defenses[0]
	for i := range store.Defenses {
		tech := &store.Defenses[i]
		if player.HasTech(&tech.Tech) {
			bestTech = tech
		}
	}
	return bestTech
}

// get the best terraform for a player
func (store *TechStore) GetBestTerraform(player *Player, terraformHabType TerraformHabType) (bestTech *TechTerraform) {
	for i := range store.Terraforms {
		tech := &store.Terraforms[i]
		if tech.HabType == terraformHabType && player.HasTech(&tech.Tech) {
			bestTech = tech
		}
	}
	return bestTech
}

// get the best engine for a player
func (store *TechStore) GetBestEngine(player *Player) *TechEngine {
	bestTech := &store.Engines[0]
	for i := range store.Engines {
		tech := &store.Engines[i]
		if player.HasTech(&tech.Tech) {
			// techs are sorted by rank, so the latest is the best
			bestTech = tech
		}
	}
	return bestTech
}

// get the best scanner for a player
func (store *TechStore) GetBestScanner(player *Player) *TechHullComponent {
	var bestTech *TechHullComponent
	for i := range store.HullComponents {
		tech := &store.HullComponents[i]
		if (tech.ScanRange > 0 || tech.ScanRangePen > 0) && player.HasTech(&tech.Tech) {
			// techs are sorted by rank, so the latest is the best
			bestTech = tech
		}
	}
	return bestTech
}

// get the best fuel tank for a player
func (store *TechStore) GetBestFuelTank(player *Player) *TechHullComponent {
	var bestTech *TechHullComponent
	for i := range store.HullComponents {
		tech := &store.HullComponents[i]
		if (tech.FuelBonus > 0) && player.HasTech(&tech.Tech) {
			// techs are sorted by rank, so the latest is the best
			bestTech = tech
		}
	}
	return bestTech
}

// get the best fuel tank for a player
func (store *TechStore) GetBestColonizationModule(player *Player) *TechHullComponent {
	var bestTech *TechHullComponent
	for i := range store.HullComponents {
		tech := &store.HullComponents[i]
		if tech.ColonizationModule && player.HasTech(&tech.Tech) {
			// techs are sorted by rank, so the latest is the best
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
	TechHullComponent: TechHullComponent{Tech: NewTech("Long Hump 6", NewCost(5, 0, 1, 6), TechRequirements{TechLevel: TechLevel{Propulsion: 3}}, 40, TechCategoryEngine), Mass: 9},
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
	TechHullComponent: TechHullComponent{Tech: NewTech("Fuel Mizer", NewCost(8, 0, 0, 11), TechRequirements{TechLevel: TechLevel{Propulsion: 2}, LRTsRequired: IFE}, 30, TechCategoryEngine), Mass: 6},
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
	TechHullComponent: TechHullComponent{Tech: NewTech("Trans-Galactic Drive", NewCost(20, 20, 9, 50), TechRequirements{TechLevel: TechLevel{Propulsion: 9}}, 70, TechCategoryEngine), Mass: 25},
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

// TechTerraforms

var TotalTerraform3 = TechTerraform{Tech: NewTech("Total Terraform ±3", NewCost(0, 0, 0, 70), TechRequirements{TechLevel: TechLevel{}, LRTsRequired: TT}, 0, TechCategoryTerraforming),
	Ability: 3,
	HabType: TerraformHabTypeAll,
}
var TotalTerraform5 = TechTerraform{Tech: NewTech("Total Terraform ±5", NewCost(0, 0, 0, 70), TechRequirements{TechLevel: TechLevel{Biotechnology: 3}, LRTsRequired: TT}, 10, TechCategoryTerraforming),
	Ability: 5,
	HabType: TerraformHabTypeAll,
}
var TotalTerraform7 = TechTerraform{Tech: NewTech("Total Terraform ±7", NewCost(0, 0, 0, 70), TechRequirements{TechLevel: TechLevel{Biotechnology: 6}, LRTsRequired: TT}, 20, TechCategoryTerraforming),
	Ability: 7,
	HabType: TerraformHabTypeAll,
}
var TotalTerraform10 = TechTerraform{Tech: NewTech("Total Terraform ±10", NewCost(0, 0, 0, 70), TechRequirements{TechLevel: TechLevel{Biotechnology: 9}, LRTsRequired: TT}, 30, TechCategoryTerraforming),
	Ability: 10,
	HabType: TerraformHabTypeAll,
}
var TotalTerraform15 = TechTerraform{Tech: NewTech("Total Terraform ±15", NewCost(0, 0, 0, 70), TechRequirements{TechLevel: TechLevel{Biotechnology: 13}, LRTsRequired: TT}, 40, TechCategoryTerraforming),
	Ability: 15,
	HabType: TerraformHabTypeAll,
}
var TotalTerraform20 = TechTerraform{Tech: NewTech("Total Terraform ±20", NewCost(0, 0, 0, 70), TechRequirements{TechLevel: TechLevel{Biotechnology: 17}, LRTsRequired: TT}, 50, TechCategoryTerraforming),
	Ability: 20,
	HabType: TerraformHabTypeAll,
}
var TotalTerraform25 = TechTerraform{Tech: NewTech("Total Terraform ±25", NewCost(0, 0, 0, 70), TechRequirements{TechLevel: TechLevel{Biotechnology: 22}, LRTsRequired: TT}, 60, TechCategoryTerraforming),
	Ability: 25,
	HabType: TerraformHabTypeAll,
}
var TotalTerraform30 = TechTerraform{Tech: NewTech("Total Terraform ±30", NewCost(0, 0, 0, 70), TechRequirements{TechLevel: TechLevel{Biotechnology: 25}, LRTsRequired: TT}, 70, TechCategoryTerraforming),
	Ability: 30,
	HabType: TerraformHabTypeAll,
}
var GravityTerraform3 = TechTerraform{Tech: NewTech("Gravity Terraform ±3", NewCost(0, 0, 0, 100), TechRequirements{TechLevel: TechLevel{Propulsion: 1, Biotechnology: 1}}, 80, TechCategoryTerraforming),
	Ability: 3,
	HabType: TerraformHabTypeGrav,
}
var GravityTerraform7 = TechTerraform{Tech: NewTech("Gravity Terraform ±7", NewCost(0, 0, 0, 100), TechRequirements{TechLevel: TechLevel{Propulsion: 5, Biotechnology: 2}}, 90, TechCategoryTerraforming),
	Ability: 7,
	HabType: TerraformHabTypeGrav,
}
var GravityTerraform11 = TechTerraform{Tech: NewTech("Gravity Terraform ±11", NewCost(0, 0, 0, 100), TechRequirements{TechLevel: TechLevel{Propulsion: 10, Biotechnology: 3}}, 100, TechCategoryTerraforming),
	Ability: 11,
	HabType: TerraformHabTypeGrav,
}
var GravityTerraform15 = TechTerraform{Tech: NewTech("Gravity Terraform ±15", NewCost(0, 0, 0, 100), TechRequirements{TechLevel: TechLevel{Propulsion: 16, Biotechnology: 4}}, 110, TechCategoryTerraforming),
	Ability: 15,
	HabType: TerraformHabTypeGrav,
}
var TempTerraform3 = TechTerraform{Tech: NewTech("Temp Terraform ±3", NewCost(0, 0, 0, 100), TechRequirements{TechLevel: TechLevel{Energy: 1, Biotechnology: 1}}, 120, TechCategoryTerraforming),
	Ability: 3,
	HabType: TerraformHabTypeTemp,
}
var TempTerraform7 = TechTerraform{Tech: NewTech("Temp Terraform ±7", NewCost(0, 0, 0, 100), TechRequirements{TechLevel: TechLevel{Energy: 5, Biotechnology: 2}}, 130, TechCategoryTerraforming),
	Ability: 7,
	HabType: TerraformHabTypeTemp,
}
var TempTerraform11 = TechTerraform{Tech: NewTech("Temp Terraform ±11", NewCost(0, 0, 0, 100), TechRequirements{TechLevel: TechLevel{Energy: 10, Biotechnology: 3}}, 140, TechCategoryTerraforming),
	Ability: 11,
	HabType: TerraformHabTypeTemp,
}
var TempTerraform15 = TechTerraform{Tech: NewTech("Temp Terraform ±15", NewCost(0, 0, 0, 100), TechRequirements{TechLevel: TechLevel{Energy: 16, Biotechnology: 4}}, 150, TechCategoryTerraforming),
	Ability: 15,
	HabType: TerraformHabTypeTemp,
}
var RadiationTerraform3 = TechTerraform{Tech: NewTech("Radiation Terraform ±3", NewCost(0, 0, 0, 100), TechRequirements{TechLevel: TechLevel{Weapons: 1, Biotechnology: 1}}, 160, TechCategoryTerraforming),
	Ability: 3,
	HabType: TerraformHabTypeRad,
}
var RadiationTerraform7 = TechTerraform{Tech: NewTech("Radiation Terraform ±7", NewCost(0, 0, 0, 100), TechRequirements{TechLevel: TechLevel{Weapons: 5, Biotechnology: 2}}, 170, TechCategoryTerraforming),
	Ability: 7,
	HabType: TerraformHabTypeRad,
}
var RadiationTerraform11 = TechTerraform{Tech: NewTech("Radiation Terraform ±11", NewCost(0, 0, 0, 100), TechRequirements{TechLevel: TechLevel{Weapons: 10, Biotechnology: 3}}, 180, TechCategoryTerraforming),
	Ability: 11,
	HabType: TerraformHabTypeRad,
}
var RadiationTerraform15 = TechTerraform{Tech: NewTech("Radiation Terraform ±15", NewCost(0, 0, 0, 100), TechRequirements{TechLevel: TechLevel{Weapons: 16, Biotechnology: 4}}, 190, TechCategoryTerraforming),
	Ability: 15,
	HabType: TerraformHabTypeRad,
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

// TechDefenses

var SDI = TechDefense{Tech: NewTech("SDI", NewCost(5, 5, 5, 15), TechRequirements{PRTDenied: AR}, 0, TechCategoryPlanetaryDefense),
	DefenseCoverage: .99,
}
var MissileBattery = TechDefense{Tech: NewTech("Missile Battery", NewCost(5, 5, 5, 15), TechRequirements{TechLevel: TechLevel{Energy: 5}, PRTDenied: AR}, 10, TechCategoryPlanetaryDefense),
	DefenseCoverage: 1.99,
}
var LaserBattery = TechDefense{Tech: NewTech("Laser Battery", NewCost(5, 5, 5, 15), TechRequirements{TechLevel: TechLevel{Energy: 10}, PRTDenied: AR}, 20, TechCategoryPlanetaryDefense),
	DefenseCoverage: 2.39,
}
var PlanetaryShield = TechDefense{Tech: NewTech("Planetary Shield", NewCost(5, 5, 5, 15), TechRequirements{TechLevel: TechLevel{Energy: 16}, PRTDenied: AR}, 30, TechCategoryPlanetaryDefense),
	DefenseCoverage: 2.99,
}
var NeutronShield = TechDefense{Tech: NewTech("Neutron Shield", NewCost(5, 5, 5, 15), TechRequirements{TechLevel: TechLevel{Energy: 23}, PRTDenied: AR}, 40, TechCategoryPlanetaryDefense),
	DefenseCoverage: 3.79,
}

// TechHullComponents

var Stargate100_250 = TechHullComponent{Tech: NewTech("Stargate 100-250", NewCost(50, 20, 20, 200), TechRequirements{TechLevel: TechLevel{Propulsion: 5, Construction: 5}, PRTDenied: HE}, 0, TechCategoryOrbital),

	Mass:         0,
	SafeHullMass: 100,
	MaxHullMass:  500,
	SafeRange:    250,
	MaxRange:     1250,
	HullSlotType: HullSlotTypeOrbital,
}

var StargateAny_300 = TechHullComponent{Tech: NewTech("Stargate any-300", NewCost(50, 20, 20, 250), TechRequirements{TechLevel: TechLevel{Propulsion: 6, Construction: 10}, PRTRequired: IT, PRTDenied: HE}, 10, TechCategoryOrbital),

	Mass:         0,
	SafeHullMass: InfinteGate,
	MaxHullMass:  InfinteGate,
	SafeRange:    300,
	MaxRange:     1500,
	HullSlotType: HullSlotTypeOrbital,
}

var Stargate150_600 = TechHullComponent{Tech: NewTech("Stargate 150-600", NewCost(50, 20, 20, 500), TechRequirements{TechLevel: TechLevel{Propulsion: 11, Construction: 7}, PRTDenied: HE}, 20, TechCategoryOrbital),

	Mass:         0,
	SafeHullMass: 150,
	MaxHullMass:  750,
	SafeRange:    600,
	MaxRange:     3000,
	HullSlotType: HullSlotTypeOrbital,
}
var Stargate300_500 = TechHullComponent{Tech: NewTech("Stargate 300-500", NewCost(50, 20, 20, 600), TechRequirements{TechLevel: TechLevel{Propulsion: 9, Construction: 13}, PRTDenied: HE}, 30, TechCategoryOrbital),

	Mass:         0,
	SafeHullMass: 300,
	MaxHullMass:  1500,
	SafeRange:    500,
	MaxRange:     2500,
	HullSlotType: HullSlotTypeOrbital,
}
var Stargate100_Any = TechHullComponent{Tech: NewTech("Stargate 100-any", NewCost(50, 20, 20, 700), TechRequirements{TechLevel: TechLevel{Propulsion: 16, Construction: 12}, PRTRequired: IT, PRTDenied: HE}, 40, TechCategoryOrbital),

	Mass:         0,
	SafeHullMass: 100,
	MaxHullMass:  500,
	SafeRange:    InfinteGate,
	MaxRange:     InfinteGate,
	HullSlotType: HullSlotTypeOrbital,
}
var StargateAny_800 = TechHullComponent{Tech: NewTech("Stargate any-800", NewCost(50, 20, 20, 700), TechRequirements{TechLevel: TechLevel{Propulsion: 12, Construction: 18}, PRTRequired: IT, PRTDenied: HE}, 50, TechCategoryOrbital),

	Mass:         0,
	SafeHullMass: InfinteGate,
	MaxHullMass:  InfinteGate,
	SafeRange:    800,
	MaxRange:     4000,
	HullSlotType: HullSlotTypeOrbital,
}
var StargateAny_Any = TechHullComponent{Tech: NewTech("Stargate any-any", NewCost(50, 20, 20, 800), TechRequirements{TechLevel: TechLevel{Propulsion: 19, Construction: 24}, PRTRequired: IT, PRTDenied: HE}, 60, TechCategoryOrbital),

	Mass:         0,
	SafeHullMass: InfinteGate,
	MaxHullMass:  InfinteGate,
	SafeRange:    InfinteGate,
	MaxRange:     InfinteGate,
	HullSlotType: HullSlotTypeOrbital,
}
var MassDriver5 = TechHullComponent{Tech: NewTech("Mass Driver 5", NewCost(24, 20, 20, 70), TechRequirements{TechLevel: TechLevel{Energy: 4}, PRTRequired: PP}, 70, TechCategoryOrbital),

	Mass:         0,
	PacketSpeed:  5,
	HullSlotType: HullSlotTypeOrbital,
}
var MassDriver6 = TechHullComponent{Tech: NewTech("Mass Driver 6", NewCost(24, 20, 20, 144), TechRequirements{TechLevel: TechLevel{Energy: 7}, PRTRequired: PP}, 80, TechCategoryOrbital),

	Mass:         0,
	PacketSpeed:  6,
	HullSlotType: HullSlotTypeOrbital,
}
var MassDriver7 = TechHullComponent{Tech: NewTech("Mass Driver 7", NewCost(100, 100, 100, 512), TechRequirements{TechLevel: TechLevel{Energy: 9}}, 90, TechCategoryOrbital),

	Mass:         0,
	PacketSpeed:  7,
	HullSlotType: HullSlotTypeOrbital,
}
var SuperDriver8 = TechHullComponent{Tech: NewTech("Super Driver 8", NewCost(24, 20, 20, 256), TechRequirements{TechLevel: TechLevel{Energy: 11}, PRTRequired: PP}, 100, TechCategoryOrbital),

	Mass:         0,
	PacketSpeed:  8,
	HullSlotType: HullSlotTypeOrbital,
}
var SuperDriver9 = TechHullComponent{Tech: NewTech("Super Driver 9", NewCost(24, 20, 20, 324), TechRequirements{TechLevel: TechLevel{Energy: 13}, PRTRequired: PP}, 110, TechCategoryOrbital),

	Mass:         0,
	PacketSpeed:  9,
	HullSlotType: HullSlotTypeOrbital,
}
var UltraDriver10 = TechHullComponent{Tech: NewTech("Ultra Driver 10", NewCost(100, 100, 100, 968), TechRequirements{TechLevel: TechLevel{Energy: 15}}, 120, TechCategoryOrbital),

	Mass:         0,
	PacketSpeed:  10,
	HullSlotType: HullSlotTypeOrbital,
}
var UltraDriver11 = TechHullComponent{Tech: NewTech("Ultra Driver 11", NewCost(24, 20, 20, 484), TechRequirements{TechLevel: TechLevel{Energy: 17}, PRTRequired: PP}, 130, TechCategoryOrbital),

	Mass:         0,
	PacketSpeed:  11,
	HullSlotType: HullSlotTypeOrbital,
}
var UltraDriver12 = TechHullComponent{Tech: NewTech("Ultra Driver 12", NewCost(24, 20, 20, 576), TechRequirements{TechLevel: TechLevel{Energy: 20}, PRTRequired: PP}, 140, TechCategoryOrbital),

	Mass:         0,
	PacketSpeed:  12,
	HullSlotType: HullSlotTypeOrbital,
}
var UltraDriver13 = TechHullComponent{Tech: NewTech("Ultra Driver 13", NewCost(24, 20, 20, 676), TechRequirements{TechLevel: TechLevel{Energy: 24}, PRTRequired: PP}, 150, TechCategoryOrbital),

	Mass:         0,
	PacketSpeed:  13,
	HullSlotType: HullSlotTypeOrbital,
}
var RoboMiner = TechHullComponent{Tech: NewTech("Robo-Miner", NewCost(30, 0, 7, 100), TechRequirements{TechLevel: TechLevel{Construction: 4, Electronics: 2}, LRTsDenied: OBRM}, 0, TechCategoryMineRobot),

	Mass:         240,
	MiningRate:   12,
	HullSlotType: HullSlotTypeMining,
}
var RoboMaxiMiner = TechHullComponent{Tech: NewTech("Robo-Maxi-Miner", NewCost(30, 0, 7, 100), TechRequirements{TechLevel: TechLevel{Construction: 7, Electronics: 4}, LRTsDenied: OBRM}, 10, TechCategoryMineRobot),

	Mass:         240,
	MiningRate:   18,
	HullSlotType: HullSlotTypeMining,
}
var RoboMidgetMiner = TechHullComponent{Tech: NewTech("Robo-Midget-Miner", NewCost(12, 0, 4, 44), TechRequirements{TechLevel: TechLevel{}, LRTsRequired: ARM}, 20, TechCategoryMineRobot),

	Mass:         80,
	MiningRate:   5,
	HullSlotType: HullSlotTypeMining,
}
var RoboMiniMiner = TechHullComponent{Tech: NewTech("Robo-Mini-Miner", NewCost(29, 0, 7, 96), TechRequirements{TechLevel: TechLevel{Construction: 2, Electronics: 1}}, 30, TechCategoryMineRobot),

	Mass:         240,
	MiningRate:   4,
	HullSlotType: HullSlotTypeMining,
}
var RoboSuperMiner = TechHullComponent{Tech: NewTech("Robo-Super-Miner", NewCost(30, 0, 7, 100), TechRequirements{TechLevel: TechLevel{Construction: 12, Electronics: 6}, LRTsDenied: OBRM}, 40, TechCategoryMineRobot),

	Mass:         240,
	MiningRate:   27,
	HullSlotType: HullSlotTypeMining,
}
var RoboUltraMiner = TechHullComponent{Tech: NewTech("Robo-Ultra-Miner", NewCost(14, 0, 4, 100), TechRequirements{TechLevel: TechLevel{Construction: 15, Electronics: 8}, LRTsRequired: ARM, LRTsDenied: OBRM}, 50, TechCategoryMineRobot),

	Mass:         80,
	MiningRate:   25,
	HullSlotType: HullSlotTypeMining,
}
var OrbitalAdjuster = TechHullComponent{Tech: NewTech("Orbital Adjuster", NewCost(25, 25, 25, 50), TechRequirements{TechLevel: TechLevel{Biotechnology: 6}, PRTRequired: CA}, 60, TechCategoryMineRobot),

	Mass:          80,
	CloakUnits:    25,
	TerraformRate: 1,
	HullSlotType:  HullSlotTypeMining,
}
var LadyFingerBomb = TechHullComponent{Tech: NewTech("Lady Finger Bomb", NewCost(1, 19, 0, 5), TechRequirements{TechLevel: TechLevel{Weapons: 2}}, 0, TechCategoryBomb),

	Mass:                 40,
	MinKillRate:          300,
	StructureDestroyRate: .2,
	KillRate:             .6,
	HullSlotType:         HullSlotTypeBomb,
}
var BlackCatBomb = TechHullComponent{Tech: NewTech("Black Cat Bomb", NewCost(1, 22, 0, 7), TechRequirements{TechLevel: TechLevel{Weapons: 5}}, 10, TechCategoryBomb),

	Mass:                 45,
	MinKillRate:          300,
	StructureDestroyRate: .4,
	KillRate:             .9,
	HullSlotType:         HullSlotTypeBomb,
}
var M70Bomb = TechHullComponent{Tech: NewTech("M-70 Bomb", NewCost(1, 24, 0, 9), TechRequirements{TechLevel: TechLevel{Weapons: 8}}, 20, TechCategoryBomb),

	Mass:                 50,
	MinKillRate:          300,
	StructureDestroyRate: .6,
	KillRate:             1.2,
	HullSlotType:         HullSlotTypeBomb,
}
var M80Bomb = TechHullComponent{Tech: NewTech("M-80 Bomb", NewCost(1, 25, 0, 12), TechRequirements{TechLevel: TechLevel{Weapons: 11}}, 30, TechCategoryBomb),

	Mass:                 55,
	MinKillRate:          300,
	StructureDestroyRate: .7,
	KillRate:             1.7,
	HullSlotType:         HullSlotTypeBomb,
}
var CherryBomb = TechHullComponent{Tech: NewTech("Cherry Bomb", NewCost(1, 25, 0, 11), TechRequirements{TechLevel: TechLevel{Weapons: 14}}, 40, TechCategoryBomb),

	Mass:                 52,
	MinKillRate:          300,
	StructureDestroyRate: 1.0,
	KillRate:             2.5,
	HullSlotType:         HullSlotTypeBomb,
}
var LBU17Bomb = TechHullComponent{Tech: NewTech("LBU-17 Bomb", NewCost(1, 15, 15, 7), TechRequirements{TechLevel: TechLevel{Weapons: 5, Electronics: 8}}, 50, TechCategoryBomb),

	Mass:                 30,
	StructureDestroyRate: 1.6,
	KillRate:             .2,
	HullSlotType:         HullSlotTypeBomb,
}
var LBU32Bomb = TechHullComponent{Tech: NewTech("LBU-32 Bomb", NewCost(1, 24, 15, 10), TechRequirements{TechLevel: TechLevel{Weapons: 10, Electronics: 10}}, 60, TechCategoryBomb),

	Mass:                 35,
	StructureDestroyRate: 2.8,
	KillRate:             .3,
	HullSlotType:         HullSlotTypeBomb,
}
var LBU74Bomb = TechHullComponent{Tech: NewTech("LBU-74 Bomb", NewCost(1, 33, 12, 14), TechRequirements{TechLevel: TechLevel{Weapons: 15, Electronics: 12}}, 70, TechCategoryBomb),

	Mass:                 45,
	StructureDestroyRate: 4.5,
	KillRate:             .4,
	HullSlotType:         HullSlotTypeBomb,
}
var RetroBomb = TechHullComponent{Tech: NewTech("Retro Bomb", NewCost(15, 15, 10, 50), TechRequirements{TechLevel: TechLevel{Weapons: 10, Biotechnology: 12}, PRTRequired: CA}, 80, TechCategoryBomb),

	Mass:            45,
	UnterraformRate: 1,
	HullSlotType:    HullSlotTypeBomb,
}
var SmartBomb = TechHullComponent{Tech: NewTech("Smart Bomb", NewCost(1, 22, 0, 27), TechRequirements{TechLevel: TechLevel{Weapons: 5, Biotechnology: 7}}, 90, TechCategoryBomb),

	Mass:         50,
	Smart:        true,
	KillRate:     1.3,
	HullSlotType: HullSlotTypeBomb,
}
var NeutronBomb = TechHullComponent{Tech: NewTech("Neutron Bomb", NewCost(1, 30, 0, 30), TechRequirements{TechLevel: TechLevel{Weapons: 10, Biotechnology: 10}}, 110, TechCategoryBomb),

	Mass:         57,
	Smart:        true,
	KillRate:     2.2,
	HullSlotType: HullSlotTypeBomb,
}
var EnrichedNeutronBomb = TechHullComponent{Tech: NewTech("Enriched Neutron Bomb", NewCost(1, 36, 0, 25), TechRequirements{TechLevel: TechLevel{Weapons: 15, Biotechnology: 12}}, 120, TechCategoryBomb),

	Mass:         64,
	Smart:        true,
	KillRate:     3.5,
	HullSlotType: HullSlotTypeBomb,
}
var PeerlessBomb = TechHullComponent{Tech: NewTech("Peerless Bomb", NewCost(1, 33, 0, 32), TechRequirements{TechLevel: TechLevel{Weapons: 22, Biotechnology: 15}}, 130, TechCategoryBomb),

	Mass:         55,
	Smart:        true,
	KillRate:     5.0,
	HullSlotType: HullSlotTypeBomb,
}
var AnnihilatorBomb = TechHullComponent{Tech: NewTech("Annihilator Bomb", NewCost(1, 30, 0, 28), TechRequirements{TechLevel: TechLevel{Weapons: 26, Biotechnology: 17}}, 140, TechCategoryBomb),

	Mass:         50,
	Smart:        true,
	KillRate:     7.0,
	HullSlotType: HullSlotTypeBomb,
}
var BatScanner = TechHullComponent{Tech: NewTech("Bat Scanner", NewCost(1, 0, 1, 1), TechRequirements{TechLevel: TechLevel{}}, 10, TechCategoryScanner),

	HullSlotType: HullSlotTypeScanner,
	Mass:         2,
	ScanRange:    ScanWithZeroRange,
}
var RhinoScanner = TechHullComponent{Tech: NewTech("Rhino Scanner", NewCost(3, 0, 2, 3), TechRequirements{TechLevel: TechLevel{Electronics: 1}}, 20, TechCategoryScanner),

	HullSlotType: HullSlotTypeScanner,
	Mass:         5,
	ScanRange:    50,
}
var MoleScanner = TechHullComponent{Tech: NewTech("Mole Scanner", NewCost(2, 0, 2, 9), TechRequirements{TechLevel: TechLevel{Electronics: 4}}, 30, TechCategoryScanner),

	HullSlotType: HullSlotTypeScanner,

	Mass:      2,
	ScanRange: 100,
}
var DNAScanner = TechHullComponent{Tech: NewTech("DNA Scanner", NewCost(1, 1, 1, 5), TechRequirements{TechLevel: TechLevel{Propulsion: 3, Biotechnology: 6}}, 40, TechCategoryScanner),

	HullSlotType: HullSlotTypeScanner,
	Mass:         2,
	ScanRange:    125,
}
var PossumScanner = TechHullComponent{Tech: NewTech("Possum Scanner", NewCost(3, 0, 3, 18), TechRequirements{TechLevel: TechLevel{Electronics: 5}}, 50, TechCategoryScanner),

	HullSlotType: HullSlotTypeScanner,
	Mass:         3,
	ScanRange:    150,
}
var PickPocketScanner = TechHullComponent{Tech: NewTech("Pick Pocket Scanner", NewCost(8, 10, 6, 35), TechRequirements{TechLevel: TechLevel{Energy: 4, Electronics: 4, Biotechnology: 4}, PRTRequired: SS}, 60, TechCategoryScanner),

	HullSlotType:       HullSlotTypeScanner,
	Mass:               15,
	CanStealFleetCargo: true,
	ScanRange:          80,
}
var ChameleonScanner = TechHullComponent{Tech: NewTech("Chameleon Scanner", NewCost(4, 6, 4, 25), TechRequirements{TechLevel: TechLevel{Energy: 3, Electronics: 6}, PRTRequired: SS}, 70, TechCategoryScanner),

	HullSlotType: HullSlotTypeScanner,
	Mass:         6,
	ScanRange:    160,
	CloakUnits:   40,
	ScanRangePen: 45,
}
var FerretScanner = TechHullComponent{Tech: NewTech("Ferret Scanner", NewCost(2, 0, 8, 36), TechRequirements{TechLevel: TechLevel{Energy: 3, Electronics: 7, Biotechnology: 2}, LRTsDenied: NAS}, 80, TechCategoryScanner),

	HullSlotType: HullSlotTypeScanner,

	Mass:         6,
	ScanRange:    185,
	ScanRangePen: 50,
}
var DolphinScanner = TechHullComponent{Tech: NewTech("Dolphin Scanner", NewCost(5, 5, 10, 40), TechRequirements{TechLevel: TechLevel{Energy: 5, Electronics: 10, Biotechnology: 4}, LRTsDenied: NAS}, 90, TechCategoryScanner),

	HullSlotType: HullSlotTypeScanner,
	Mass:         4,
	ScanRange:    220,
	ScanRangePen: 100,
}
var GazelleScanner = TechHullComponent{Tech: NewTech("Gazelle Scanner", NewCost(4, 0, 5, 24), TechRequirements{TechLevel: TechLevel{Energy: 4, Electronics: 8}}, 100, TechCategoryScanner),

	HullSlotType: HullSlotTypeScanner,
	Mass:         5,
	ScanRange:    225,
}
var RNAScanner = TechHullComponent{Tech: NewTech("RNA Scanner", NewCost(1, 1, 2, 20), TechRequirements{TechLevel: TechLevel{Propulsion: 5, Biotechnology: 10}}, 110, TechCategoryScanner),

	HullSlotType: HullSlotTypeScanner,

	Mass:      2,
	ScanRange: 230,
}
var CheetahScanner = TechHullComponent{Tech: NewTech("Cheetah Scanner", NewCost(3, 1, 13, 50), TechRequirements{TechLevel: TechLevel{Energy: 5, Electronics: 11}}, 120, TechCategoryScanner),

	HullSlotType: HullSlotTypeScanner,
	Mass:         4,
	ScanRange:    275,
}
var ElephantScanner = TechHullComponent{Tech: NewTech("Elephant Scanner", NewCost(8, 5, 14, 70), TechRequirements{TechLevel: TechLevel{Energy: 6, Electronics: 16, Biotechnology: 7}, LRTsDenied: NAS}, 130, TechCategoryScanner),

	HullSlotType: HullSlotTypeScanner,
	Mass:         6,
	ScanRange:    300,
	ScanRangePen: 200,
}
var EagleEyeScanner = TechHullComponent{Tech: NewTech("Eagle Eye Scanner", NewCost(3, 2, 21, 64), TechRequirements{TechLevel: TechLevel{Energy: 6, Electronics: 14}}, 140, TechCategoryScanner),

	HullSlotType: HullSlotTypeScanner,
	Mass:         3,
	ScanRange:    335,
}
var RobberBaronScanner = TechHullComponent{Tech: NewTech("Robber Baron Scanner", NewCost(10, 10, 10, 90), TechRequirements{TechLevel: TechLevel{Energy: 10, Electronics: 15, Biotechnology: 10}, PRTRequired: SS}, 150, TechCategoryScanner),

	HullSlotType:        HullSlotTypeScanner,
	Mass:                20,
	CanStealFleetCargo:  true,
	CanStealPlanetCargo: true,
	ScanRange:           220,
	ScanRangePen:        120,
}
var PeerlessScanner = TechHullComponent{Tech: NewTech("Peerless Scanner", NewCost(3, 2, 30, 90), TechRequirements{TechLevel: TechLevel{Energy: 7, Electronics: 24}}, 160, TechCategoryScanner),

	HullSlotType: HullSlotTypeScanner,

	Mass:      4,
	ScanRange: 500,
}
var Tritanium = TechHullComponent{Tech: NewTech("Tritanium", NewCost(5, 0, 0, 9), TechRequirements{TechLevel: TechLevel{}}, 10, TechCategoryArmor),

	Mass:         60,
	Armor:        50,
	HullSlotType: HullSlotTypeArmor,
}
var Crobmnium = TechHullComponent{Tech: NewTech("Crobmnium", NewCost(6, 0, 0, 13), TechRequirements{TechLevel: TechLevel{Construction: 3}}, 20, TechCategoryArmor),

	Mass:         56,
	Armor:        75,
	HullSlotType: HullSlotTypeArmor,
}
var Carbonic = TechHullComponent{Tech: NewTech("Carbonic Armor", NewCost(5, 0, 0, 15), TechRequirements{TechLevel: TechLevel{Biotechnology: 4}}, 30, TechCategoryArmor),

	Mass:         25,
	Armor:        100,
	HullSlotType: HullSlotTypeArmor,
}
var Strobnium = TechHullComponent{Tech: NewTech("Strobnium", NewCost(8, 0, 0, 18), TechRequirements{TechLevel: TechLevel{Construction: 6}}, 40, TechCategoryArmor),

	Mass:         54,
	Armor:        120,
	HullSlotType: HullSlotTypeArmor,
}
var Organic = TechHullComponent{Tech: NewTech("Organic Armor", NewCost(0, 0, 6, 20), TechRequirements{TechLevel: TechLevel{Biotechnology: 7}}, 50, TechCategoryArmor),

	Mass:         15,
	Armor:        175,
	HullSlotType: HullSlotTypeArmor,
}
var Kelarium = TechHullComponent{Tech: NewTech("Kelarium", NewCost(9, 1, 0, 25), TechRequirements{TechLevel: TechLevel{Construction: 9}}, 60, TechCategoryArmor),

	Mass:         50,
	Armor:        180,
	HullSlotType: HullSlotTypeArmor,
}
var FieldedKelarium = TechHullComponent{Tech: NewTech("Fielded Kelarium", NewCost(10, 0, 2, 28), TechRequirements{TechLevel: TechLevel{Energy: 4, Construction: 10}, PRTRequired: IS}, 70, TechCategoryArmor),

	Mass:         50,
	Shield:       50,
	Armor:        175,
	HullSlotType: HullSlotTypeArmor,
}
var DepletedNeutronium = TechHullComponent{Tech: NewTech("Depleted Neutronium", NewCost(10, 0, 2, 28), TechRequirements{TechLevel: TechLevel{Construction: 10, Electronics: 3}, PRTRequired: SS}, 80, TechCategoryArmor),

	Mass:         50,
	Armor:        200,
	CloakUnits:   50,
	HullSlotType: HullSlotTypeArmor,
}
var Neutronium = TechHullComponent{Tech: NewTech("Neutronium", NewCost(11, 2, 1, 30), TechRequirements{TechLevel: TechLevel{Construction: 12}}, 90, TechCategoryArmor),

	Mass:         45,
	Armor:        275,
	HullSlotType: HullSlotTypeArmor,
}
var Valanium = TechHullComponent{Tech: NewTech("Valanium", NewCost(15, 0, 0, 50), TechRequirements{TechLevel: TechLevel{Construction: 16}}, 100, TechCategoryArmor),

	Mass:         40,
	Armor:        500,
	HullSlotType: HullSlotTypeArmor,
}
var Superlatanium = TechHullComponent{Tech: NewTech("Superlatanium", NewCost(25, 0, 0, 100), TechRequirements{TechLevel: TechLevel{Construction: 24}}, 110, TechCategoryArmor),

	Mass:         30,
	Armor:        1500,
	HullSlotType: HullSlotTypeArmor,
}
var TransportCloaking = TechHullComponent{Tech: NewTech("Transport Cloaking", NewCost(2, 0, 2, 3), TechRequirements{TechLevel: TechLevel{}, PRTRequired: SS}, 0, TechCategoryElectrical),

	Mass:             1,
	CloakUnarmedOnly: true,
	CloakUnits:       300,
	HullSlotType:     HullSlotTypeElectrical,
}
var StealthCloak = TechHullComponent{Tech: NewTech("Stealth Cloak", NewCost(2, 0, 2, 5), TechRequirements{TechLevel: TechLevel{Energy: 2, Electronics: 5}}, 10, TechCategoryElectrical),

	Mass:         2,
	CloakUnits:   70,
	HullSlotType: HullSlotTypeElectrical,
}
var SuperStealthCloak = TechHullComponent{Tech: NewTech("Super-Stealth Cloak", NewCost(8, 0, 8, 15), TechRequirements{TechLevel: TechLevel{Energy: 4, Electronics: 10}}, 20, TechCategoryElectrical),

	Mass:         3,
	CloakUnits:   140,
	HullSlotType: HullSlotTypeElectrical,
}
var UltraStealthCloak = TechHullComponent{Tech: NewTech("Ultra-Stealth Cloak", NewCost(10, 0, 10, 25), TechRequirements{TechLevel: TechLevel{Energy: 10, Electronics: 12}, PRTRequired: SS}, 30, TechCategoryElectrical),

	Mass:         5,
	CloakUnits:   540,
	HullSlotType: HullSlotTypeElectrical,
}
var BattleComputer = TechHullComponent{Tech: NewTech("Battle Computer", NewCost(0, 0, 13, 5), TechRequirements{TechLevel: TechLevel{}}, 40, TechCategoryElectrical),

	Mass:            1,
	InitiativeBonus: 1,
	TorpedoBonus:    .2,
	HullSlotType:    HullSlotTypeElectrical,
}
var BattleSuperComputer = TechHullComponent{Tech: NewTech("Battle Super Computer", NewCost(0, 0, 25, 14), TechRequirements{TechLevel: TechLevel{Energy: 5, Electronics: 11}}, 50, TechCategoryElectrical),

	Mass:            1,
	InitiativeBonus: 2,
	TorpedoBonus:    .3,
	HullSlotType:    HullSlotTypeElectrical,
}
var BattleNexus = TechHullComponent{Tech: NewTech("Battle Nexus", NewCost(0, 0, 30, 15), TechRequirements{TechLevel: TechLevel{Energy: 10, Electronics: 19}}, 60, TechCategoryElectrical),

	Mass:            1,
	InitiativeBonus: 3,
	TorpedoBonus:    .5,
	HullSlotType:    HullSlotTypeElectrical,
}
var Jammer10 = TechHullComponent{Tech: NewTech("Jammer 10", NewCost(0, 0, 2, 6), TechRequirements{TechLevel: TechLevel{Energy: 2, Electronics: 6}, PRTRequired: IS}, 70, TechCategoryElectrical),

	Mass:           1,
	TorpedoJamming: .1,
	HullSlotType:   HullSlotTypeElectrical,
}
var Jammer20 = TechHullComponent{Tech: NewTech("Jammer 20", NewCost(1, 0, 5, 20), TechRequirements{TechLevel: TechLevel{Energy: 4, Electronics: 10}}, 80, TechCategoryElectrical),

	Mass:           1,
	TorpedoJamming: .2,
	HullSlotType:   HullSlotTypeElectrical,
}
var Jammer30 = TechHullComponent{Tech: NewTech("Jammer 30", NewCost(1, 0, 6, 20), TechRequirements{TechLevel: TechLevel{Energy: 8, Electronics: 16}}, 90, TechCategoryElectrical),

	Mass:           1,
	TorpedoJamming: .3,
	HullSlotType:   HullSlotTypeElectrical,
}
var Jammer50 = TechHullComponent{Tech: NewTech("Jammer 50", NewCost(2, 0, 7, 20), TechRequirements{TechLevel: TechLevel{Energy: 16, Electronics: 22}}, 100, TechCategoryElectrical),

	Mass:           1,
	TorpedoJamming: .5,
	HullSlotType:   HullSlotTypeElectrical,
}
var EnergyCapacitor = TechHullComponent{Tech: NewTech("Energy Capacitor", NewCost(0, 0, 8, 5), TechRequirements{TechLevel: TechLevel{Energy: 7, Electronics: 4}}, 110, TechCategoryElectrical),

	Mass:         1,
	BeamBonus:    .1,
	HullSlotType: HullSlotTypeElectrical,
}
var FluxCapacitor = TechHullComponent{Tech: NewTech("Flux Capacitor", NewCost(0, 0, 8, 5), TechRequirements{TechLevel: TechLevel{Energy: 14, Electronics: 8}, PRTRequired: HE}, 120, TechCategoryElectrical),

	Mass:         1,
	BeamBonus:    .1,
	HullSlotType: HullSlotTypeElectrical,
}
var EnergyDampener = TechHullComponent{Tech: NewTech("Energy Dampener", NewCost(5, 10, 0, 50), TechRequirements{TechLevel: TechLevel{Energy: 14, Propulsion: 8}, PRTRequired: SD}, 130, TechCategoryElectrical),

	Mass:           2,
	ReduceMovement: 1,
	HullSlotType:   HullSlotTypeElectrical,
}
var TachyonDetector = TechHullComponent{Tech: NewTech("Tachyon Detector", NewCost(1, 5, 0, 70), TechRequirements{TechLevel: TechLevel{Energy: 8, Electronics: 14}, PRTRequired: IS}, 140, TechCategoryElectrical),

	Mass:           1,
	ReduceCloaking: true,
	HullSlotType:   HullSlotTypeElectrical,
}
var AntiMatterGenerator = TechHullComponent{Tech: NewTech("Anti-Matter Generator", NewCost(8, 3, 3, 10), TechRequirements{TechLevel: TechLevel{Weapons: 12, Biotechnology: 7}, PRTRequired: IT}, 150, TechCategoryElectrical),

	Mass:                 10,
	FuelRegenerationRate: 50,
	FuelBonus:            200,
	HullSlotType:         HullSlotTypeElectrical,
}
var MineDispenser40 = TechHullComponent{Tech: NewTech("Mine Dispenser 40", NewCost(2, 9, 7, 40), TechRequirements{TechLevel: TechLevel{}, PRTRequired: SD}, 0, TechCategoryMineLayer),

	Mass:           25,
	MineFieldType:  MineFieldTypeStandard,
	MineLayingRate: 40,
	HullSlotType:   HullSlotTypeMineLayer,
}
var MineDispenser50 = TechHullComponent{Tech: NewTech("Mine Dispenser 50", NewCost(2, 12, 10, 55), TechRequirements{TechLevel: TechLevel{Energy: 2, Biotechnology: 4}}, 10, TechCategoryMineLayer),

	Mass:           30,
	MineFieldType:  MineFieldTypeStandard,
	MineLayingRate: 50,
	HullSlotType:   HullSlotTypeMineLayer,
}
var MineDispenser80 = TechHullComponent{Tech: NewTech("Mine Dispenser 80", NewCost(2, 12, 10, 65), TechRequirements{TechLevel: TechLevel{Energy: 3, Biotechnology: 7}, PRTRequired: SD}, 20, TechCategoryMineLayer),

	Mass:           30,
	MineFieldType:  MineFieldTypeStandard,
	MineLayingRate: 80,
	HullSlotType:   HullSlotTypeMineLayer,
}
var MineDispenser130 = TechHullComponent{Tech: NewTech("Mine Dispenser 130", NewCost(2, 18, 10, 80), TechRequirements{TechLevel: TechLevel{Energy: 6, Biotechnology: 12}, PRTRequired: SD}, 30, TechCategoryMineLayer),

	Mass:           30,
	MineFieldType:  MineFieldTypeStandard,
	MineLayingRate: 130,
	HullSlotType:   HullSlotTypeMineLayer,
}
var HeavyDispenser50 = TechHullComponent{Tech: NewTech("Heavy Dispenser 50", NewCost(2, 20, 5, 50), TechRequirements{TechLevel: TechLevel{Energy: 5, Biotechnology: 3}, PRTRequired: SD}, 40, TechCategoryMineLayer),

	Mass:           10,
	MineFieldType:  MineFieldTypeHeavy,
	MineLayingRate: 50,
	HullSlotType:   HullSlotTypeMineLayer,
}
var HeavyDispenser110 = TechHullComponent{Tech: NewTech("Heavy Dispenser 110", NewCost(2, 20, 5, 50), TechRequirements{TechLevel: TechLevel{Energy: 9, Biotechnology: 5}, PRTRequired: SD}, 50, TechCategoryMineLayer),

	Mass:           15,
	MineFieldType:  MineFieldTypeHeavy,
	MineLayingRate: 110,
	HullSlotType:   HullSlotTypeMineLayer,
}
var HeavyDispenser200 = TechHullComponent{Tech: NewTech("Heavy Dispenser 200", NewCost(2, 45, 5, 90), TechRequirements{TechLevel: TechLevel{Energy: 14, Biotechnology: 7}, PRTRequired: SD}, 60, TechCategoryMineLayer),

	Mass:           20,
	MineFieldType:  MineFieldTypeHeavy,
	MineLayingRate: 200,
	HullSlotType:   HullSlotTypeMineLayer,
}
var SpeedTrap20 = TechHullComponent{Tech: NewTech("Speed Trap 20", NewCost(29, 0, 12, 58), TechRequirements{TechLevel: TechLevel{Propulsion: 2, Biotechnology: 2}, PRTRequired: SD}, 70, TechCategoryMineLayer),

	Mass:           100,
	MineFieldType:  MineFieldTypeSpeedBump,
	MineLayingRate: 20,
	HullSlotType:   HullSlotTypeMineLayer,
}
var SpeedTrap30 = TechHullComponent{Tech: NewTech("Speed Trap 30", NewCost(32, 0, 14, 72), TechRequirements{TechLevel: TechLevel{Propulsion: 3, Biotechnology: 6}, PRTRequired: SD}, 80, TechCategoryMineLayer),

	Mass:           135,
	MineFieldType:  MineFieldTypeSpeedBump,
	MineLayingRate: 30,
	HullSlotType:   HullSlotTypeMineLayer,
}
var SpeedTrap50 = TechHullComponent{Tech: NewTech("Speed Trap 50", NewCost(40, 0, 15, 80), TechRequirements{TechLevel: TechLevel{Propulsion: 5, Biotechnology: 11}, PRTRequired: SD}, 90, TechCategoryMineLayer),

	Mass:           140,
	MineFieldType:  MineFieldTypeSpeedBump,
	MineLayingRate: 50,
	HullSlotType:   HullSlotTypeMineLayer,
}
var ColonizationModule = TechHullComponent{Tech: NewTech("Colonization Module", NewCost(11, 9, 9, 9), TechRequirements{TechLevel: TechLevel{}}, 0, TechCategoryMechanical),

	Mass:               32,
	ColonizationModule: true,
	HullSlotType:       HullSlotTypeMechanical,
}
var OrbitalConstructionModule = TechHullComponent{Tech: NewTech("Orbital Construction Module", NewCost(18, 13, 13, 18), TechRequirements{TechLevel: TechLevel{}, PRTRequired: AR}, 10, TechCategoryMechanical),

	Mass:                      50,
	MinKillRate:               2000,
	OrbitalConstructionModule: true,
	HullSlotType:              HullSlotTypeArmor,
}
var CargoPod = TechHullComponent{Tech: NewTech("Cargo Pod", NewCost(5, 0, 2, 10), TechRequirements{TechLevel: TechLevel{Construction: 3}}, 20, TechCategoryMechanical),

	Mass:         5,
	CargoBonus:   50,
	HullSlotType: HullSlotTypeMechanical,
}
var SuperCargoPod = TechHullComponent{Tech: NewTech("Super Cargo Pod", NewCost(8, 0, 2, 15), TechRequirements{TechLevel: TechLevel{Energy: 3, Construction: 8}}, 30, TechCategoryMechanical),

	Mass:         7,
	CargoBonus:   100,
	HullSlotType: HullSlotTypeMechanical,
}
var FuelTank = TechHullComponent{Tech: NewTech("Fuel Tank", NewCost(5, 0, 0, 4), TechRequirements{TechLevel: TechLevel{}}, 40, TechCategoryMechanical),

	Mass:         3,
	FuelBonus:    250,
	HullSlotType: HullSlotTypeMechanical,
}
var SuperFuelTank = TechHullComponent{Tech: NewTech("Super Fuel Tank", NewCost(8, 0, 0, 8), TechRequirements{TechLevel: TechLevel{Energy: 6, Propulsion: 4, Construction: 14}}, 50, TechCategoryMechanical),

	Mass:         8,
	FuelBonus:    500,
	HullSlotType: HullSlotTypeMechanical,
}
var ManeuveringJet = TechHullComponent{Tech: NewTech("Maneuvering Jet", NewCost(5, 0, 5, 10), TechRequirements{TechLevel: TechLevel{Energy: 2, Propulsion: 3}}, 60, TechCategoryMechanical),

	Mass:          5,
	MovementBonus: 1,
	HullSlotType:  HullSlotTypeMechanical,
}
var Overthruster = TechHullComponent{Tech: NewTech("Overthruster", NewCost(10, 0, 8, 20), TechRequirements{TechLevel: TechLevel{Energy: 5, Propulsion: 12}}, 70, TechCategoryMechanical),

	Mass:          5,
	MovementBonus: 2,
	HullSlotType:  HullSlotTypeMechanical,
}
var BeamDeflector = TechHullComponent{Tech: NewTech("Beam Deflector", NewCost(0, 0, 10, 8), TechRequirements{TechLevel: TechLevel{Energy: 6, Weapons: 6, Construction: 6, Electronics: 6}}, 80, TechCategoryMechanical),

	Mass:         1,
	HullSlotType: HullSlotTypeMechanical,
	BeamDefense:  1,
}
var Laser = TechHullComponent{Tech: NewTech("Laser", NewCost(0, 6, 0, 5), TechRequirements{TechLevel: TechLevel{}}, 0, TechCategoryBeamWeapon),

	Mass:         1,
	Initiative:   9,
	Power:        10,
	HullSlotType: HullSlotTypeWeapon,

	Range: 1,
}
var XRayLaser = TechHullComponent{Tech: NewTech("X-Ray Laser", NewCost(0, 6, 0, 6), TechRequirements{TechLevel: TechLevel{Weapons: 3}}, 10, TechCategoryBeamWeapon),

	Mass:         1,
	Initiative:   9,
	Power:        16,
	HullSlotType: HullSlotTypeWeapon,

	Range: 1,
}
var MiniGun = TechHullComponent{Tech: NewTech("Mini Gun", NewCost(0, 6, 0, 6), TechRequirements{TechLevel: TechLevel{Weapons: 5}, PRTRequired: IS}, 20, TechCategoryBeamWeapon),

	Mass:           3,
	Initiative:     12,
	Gattling:       true,
	Power:          16,
	HitsAllTargets: true,
	HullSlotType:   HullSlotTypeWeapon,

	Range: 2,
}
var YakimoraLightPhaser = TechHullComponent{Tech: NewTech("Yakimora Light Phaser", NewCost(0, 8, 0, 7), TechRequirements{TechLevel: TechLevel{Weapons: 6}}, 30, TechCategoryBeamWeapon),

	Mass:         1,
	Initiative:   9,
	Power:        26,
	HullSlotType: HullSlotTypeWeapon,

	Range: 1,
}
var Blackjack = TechHullComponent{Tech: NewTech("Blackjack", NewCost(0, 16, 0, 7), TechRequirements{TechLevel: TechLevel{Weapons: 7}}, 40, TechCategoryBeamWeapon),

	Mass:         10,
	Initiative:   10,
	Power:        90,
	HullSlotType: HullSlotTypeWeapon,

	Range: 0,
}
var PhaserBazooka = TechHullComponent{Tech: NewTech("Phaser Bazooka", NewCost(0, 8, 0, 11), TechRequirements{TechLevel: TechLevel{Weapons: 8}}, 50, TechCategoryBeamWeapon),

	Mass:         2,
	Initiative:   7,
	Power:        26,
	HullSlotType: HullSlotTypeWeapon,

	Range: 2,
}
var PulsedSapper = TechHullComponent{Tech: NewTech("Pulsed Sapper", NewCost(0, 0, 4, 12), TechRequirements{TechLevel: TechLevel{Energy: 5, Weapons: 9}}, 60, TechCategoryBeamWeapon),

	Mass:              1,
	Initiative:        14,
	DamageShieldsOnly: true,
	Power:             82,
	HullSlotType:      HullSlotTypeWeapon,

	Range: 3,
}
var ColloidalPhaser = TechHullComponent{Tech: NewTech("Colloidal Phaser", NewCost(0, 14, 0, 18), TechRequirements{TechLevel: TechLevel{Weapons: 10}}, 70, TechCategoryBeamWeapon),

	Mass:         2,
	Initiative:   5,
	Power:        26,
	HullSlotType: HullSlotTypeWeapon,

	Range: 3,
}
var GatlingGun = TechHullComponent{Tech: NewTech("Gatling Gun", NewCost(0, 20, 0, 13), TechRequirements{TechLevel: TechLevel{Weapons: 11}}, 80, TechCategoryBeamWeapon),

	Mass:           3,
	Initiative:     12,
	Gattling:       true,
	Power:          31,
	HitsAllTargets: true,
	HullSlotType:   HullSlotTypeWeapon,
	Range:          2,
}
var MiniBlaster = TechHullComponent{Tech: NewTech("Mini Blaster", NewCost(0, 10, 0, 9), TechRequirements{TechLevel: TechLevel{Weapons: 12}}, 90, TechCategoryBeamWeapon),

	Mass:         1,
	Initiative:   9,
	Power:        66,
	HullSlotType: HullSlotTypeWeapon,
	Range:        1,
}
var Bludgeon = TechHullComponent{Tech: NewTech("Bludgeon", NewCost(0, 22, 0, 9), TechRequirements{TechLevel: TechLevel{Weapons: 13}}, 100, TechCategoryBeamWeapon),

	Mass:         10,
	Initiative:   10,
	Power:        231,
	HullSlotType: HullSlotTypeWeapon,
	Range:        0,
}
var MarkIVBlaster = TechHullComponent{Tech: NewTech("Mark IV Blaster", NewCost(0, 12, 0, 15), TechRequirements{TechLevel: TechLevel{Weapons: 14}}, 110, TechCategoryBeamWeapon),

	Mass:         2,
	Initiative:   7,
	Power:        66,
	HullSlotType: HullSlotTypeWeapon,
	Range:        2,
}
var PhasedSapper = TechHullComponent{Tech: NewTech("Phased Sapper", NewCost(0, 0, 6, 16), TechRequirements{TechLevel: TechLevel{Energy: 8, Weapons: 15}}, 120, TechCategoryBeamWeapon),

	Mass:              1,
	Initiative:        14,
	DamageShieldsOnly: true,
	Power:             211,
	HullSlotType:      HullSlotTypeWeapon,
	Range:             3,
}
var HeavyBlaster = TechHullComponent{Tech: NewTech("Heavy Blaster", NewCost(0, 20, 0, 25), TechRequirements{TechLevel: TechLevel{Weapons: 16}}, 130, TechCategoryBeamWeapon),

	Mass:         2,
	Initiative:   5,
	Power:        66,
	HullSlotType: HullSlotTypeWeapon,
	Range:        3,
}
var GatlingNeutrinoCannon = TechHullComponent{Tech: NewTech("Gatling Neutrino Cannon", NewCost(0, 28, 0, 17), TechRequirements{TechLevel: TechLevel{Weapons: 17}, PRTRequired: WM}, 140, TechCategoryBeamWeapon),

	Mass:           3,
	Initiative:     13,
	Gattling:       true,
	Power:          80,
	HitsAllTargets: true,
	HullSlotType:   HullSlotTypeWeapon,
	Range:          2,
}
var MyopicDisruptor = TechHullComponent{Tech: NewTech("Myopic Disruptor", NewCost(0, 14, 0, 12), TechRequirements{TechLevel: TechLevel{Weapons: 18}}, 150, TechCategoryBeamWeapon),

	Mass:         1,
	Initiative:   9,
	Power:        169,
	HullSlotType: HullSlotTypeWeapon,
	Range:        1,
}
var Blunderbuss = TechHullComponent{Tech: NewTech("Blunderbuss", NewCost(0, 30, 0, 13), TechRequirements{TechLevel: TechLevel{Weapons: 19}, PRTRequired: WM}, 160, TechCategoryBeamWeapon),

	Mass:         10,
	Initiative:   11,
	Power:        592,
	HullSlotType: HullSlotTypeWeapon,
	Range:        0,
}
var Disruptor = TechHullComponent{Tech: NewTech("Disruptor", NewCost(0, 16, 0, 20), TechRequirements{TechLevel: TechLevel{Weapons: 20}}, 170, TechCategoryBeamWeapon),

	Mass:         2,
	Initiative:   8,
	Power:        169,
	HullSlotType: HullSlotTypeWeapon,
	Range:        2,
}
var SyncroSapper = TechHullComponent{Tech: NewTech("Syncro Sapper", NewCost(0, 0, 8, 21), TechRequirements{TechLevel: TechLevel{Energy: 11, Weapons: 21}}, 180, TechCategoryBeamWeapon),

	Mass:              1,
	Initiative:        14,
	DamageShieldsOnly: true,
	Power:             541,
	HullSlotType:      HullSlotTypeWeapon,
	Range:             3,
}
var MegaDisruptor = TechHullComponent{Tech: NewTech("Mega Disruptor", NewCost(0, 30, 0, 33), TechRequirements{TechLevel: TechLevel{Weapons: 22}}, 190, TechCategoryBeamWeapon),

	Mass:         2,
	Initiative:   6,
	Power:        169,
	HullSlotType: HullSlotTypeWeapon,
	Range:        3,
}
var BigMuthaCannon = TechHullComponent{Tech: NewTech("Big Mutha Cannon", NewCost(0, 36, 0, 23), TechRequirements{TechLevel: TechLevel{Weapons: 23}}, 200, TechCategoryBeamWeapon),

	Mass:           3,
	Initiative:     13,
	Gattling:       true,
	Power:          204,
	HitsAllTargets: true,
	HullSlotType:   HullSlotTypeWeapon,
	Range:          2,
}
var StreamingPulverizer = TechHullComponent{Tech: NewTech("Streaming Pulverizer", NewCost(0, 20, 0, 16), TechRequirements{TechLevel: TechLevel{Weapons: 24}}, 210, TechCategoryBeamWeapon),

	Mass:         1,
	Initiative:   9,
	Power:        433,
	HullSlotType: HullSlotTypeWeapon,
	Range:        1,
}
var AntiMatterPulverizer = TechHullComponent{Tech: NewTech("Anti-Matter Pulverizer", NewCost(0, 22, 0, 27), TechRequirements{TechLevel: TechLevel{Weapons: 26}}, 220, TechCategoryBeamWeapon),

	Mass:         1,
	Initiative:   8,
	Power:        433,
	HullSlotType: HullSlotTypeWeapon,
	Range:        2,
}
var AlphaTorpedo = TechHullComponent{Tech: NewTech("Alpha Torpedo", NewCost(8, 3, 3, 4), TechRequirements{TechLevel: TechLevel{}}, 0, TechCategoryTorpedo),

	Mass:         25,
	Initiative:   0,
	Accuracy:     35,
	Power:        5,
	HullSlotType: HullSlotTypeWeapon,
	Range:        4,
}
var BetaTorpedo = TechHullComponent{Tech: NewTech("Beta Torpedo", NewCost(18, 6, 4, 6), TechRequirements{TechLevel: TechLevel{Weapons: 5, Propulsion: 1}}, 10, TechCategoryTorpedo),

	Mass:         25,
	Initiative:   1,
	Accuracy:     45,
	Power:        12,
	HullSlotType: HullSlotTypeWeapon,
	Range:        4,
}
var DeltaTorpedo = TechHullComponent{Tech: NewTech("Delta Torpedo", NewCost(22, 8, 5, 8), TechRequirements{TechLevel: TechLevel{Weapons: 10, Propulsion: 2}}, 20, TechCategoryTorpedo),

	Mass:         25,
	Initiative:   1,
	Accuracy:     60,
	Power:        26,
	HullSlotType: HullSlotTypeWeapon,
	Range:        4,
}

var EpsilonTorpedo = TechHullComponent{Tech: NewTech("Epsilon Torpedo", NewCost(30, 10, 6, 10), TechRequirements{TechLevel: TechLevel{Weapons: 14, Propulsion: 3}}, 30, TechCategoryTorpedo),

	Mass:         25,
	Initiative:   2,
	Accuracy:     65,
	Power:        48,
	HullSlotType: HullSlotTypeWeapon,
	Range:        5,
}

var RhoTorpedo = TechHullComponent{Tech: NewTech("Rho Torpedo", NewCost(34, 12, 8, 12), TechRequirements{TechLevel: TechLevel{Weapons: 18, Propulsion: 4}}, 40, TechCategoryTorpedo),

	Mass:         25,
	Initiative:   2,
	Accuracy:     75,
	Power:        90,
	HullSlotType: HullSlotTypeWeapon,
	Range:        5,
}

var UpsilonTorpedo = TechHullComponent{Tech: NewTech("Upsilon Torpedo", NewCost(40, 14, 9, 15), TechRequirements{TechLevel: TechLevel{Weapons: 22, Propulsion: 5}}, 50, TechCategoryTorpedo),

	Mass:         25,
	Initiative:   3,
	Accuracy:     75,
	Power:        169,
	HullSlotType: HullSlotTypeWeapon,
	Range:        5,
}

var OmegaTorpedo = TechHullComponent{Tech: NewTech("Omega Torpedo", NewCost(52, 18, 12, 18), TechRequirements{TechLevel: TechLevel{Weapons: 26, Propulsion: 6}}, 60, TechCategoryTorpedo),

	Mass:         25,
	Initiative:   4,
	Accuracy:     80,
	Power:        316,
	HullSlotType: HullSlotTypeWeapon,
	Range:        5,
}
var JihadMissile = TechHullComponent{Tech: NewTech("Jihad Missile", NewCost(37, 13, 9, 13), TechRequirements{TechLevel: TechLevel{Weapons: 12, Propulsion: 6}}, 70, TechCategoryTorpedo),

	Mass:               35,
	Accuracy:           20,
	CapitalShipMissile: true,
	Power:              85,
	HullSlotType:       HullSlotTypeWeapon,
	Range:              5,
}
var JuggernautMissile = TechHullComponent{Tech: NewTech("Juggernaut Missile", NewCost(48, 16, 11, 16), TechRequirements{TechLevel: TechLevel{Weapons: 16, Propulsion: 8}}, 80, TechCategoryTorpedo),

	Mass:               35,
	Initiative:         1,
	Accuracy:           20,
	CapitalShipMissile: true,
	Power:              150,
	HullSlotType:       HullSlotTypeWeapon,
	Range:              5,
}

var DoomsdayMissile = TechHullComponent{Tech: NewTech("Doomsday Missile", NewCost(60, 20, 13, 20), TechRequirements{TechLevel: TechLevel{Weapons: 20, Propulsion: 10}}, 90, TechCategoryTorpedo),

	Mass:               35,
	Initiative:         2,
	Accuracy:           25,
	CapitalShipMissile: true,
	Power:              280,
	HullSlotType:       HullSlotTypeWeapon,
	Range:              6,
}

var ArmageddonMissile = TechHullComponent{Tech: NewTech("Armageddon Missile", NewCost(67, 23, 16, 24), TechRequirements{TechLevel: TechLevel{Weapons: 24, Propulsion: 10}}, 100, TechCategoryTorpedo),

	Mass:               35,
	Initiative:         3,
	Accuracy:           30,
	CapitalShipMissile: true,
	Power:              525,
	HullSlotType:       HullSlotTypeWeapon,
	Range:              6,
}
var MoleSkinShield = TechHullComponent{Tech: NewTech("Mole-skin Shield", NewCost(1, 0, 1, 4), TechRequirements{TechLevel: TechLevel{}}, 10, TechCategoryShield),

	Mass:         1,
	Shield:       25,
	HullSlotType: HullSlotTypeShield,
}
var CowHideShield = TechHullComponent{Tech: NewTech("Cow-hide Shield", NewCost(2, 0, 2, 5), TechRequirements{TechLevel: TechLevel{Energy: 3}}, 20, TechCategoryShield),

	Mass:         1,
	Shield:       40,
	HullSlotType: HullSlotTypeShield,
}
var WolverineDiffuseShield = TechHullComponent{Tech: NewTech("Wolverine Diffuse Shield", NewCost(3, 0, 3, 6), TechRequirements{TechLevel: TechLevel{Energy: 6}}, 30, TechCategoryShield),

	Mass:         1,
	Shield:       60,
	HullSlotType: HullSlotTypeShield,
}
var CrobySharmor = TechHullComponent{Tech: NewTech("Croby Sharmor", NewCost(7, 0, 4, 15), TechRequirements{TechLevel: TechLevel{Energy: 7, Construction: 4}, PRTRequired: IS}, 40, TechCategoryShield),

	Mass:         10,
	Shield:       60,
	Armor:        65,
	HullSlotType: HullSlotTypeShield,
}
var ShadowShield = TechHullComponent{Tech: NewTech("Shadow Shield", NewCost(3, 0, 3, 7), TechRequirements{TechLevel: TechLevel{Energy: 7, Electronics: 3}, PRTRequired: SS}, 50, TechCategoryShield),

	Mass:         2,
	Shield:       75,
	CloakUnits:   70,
	HullSlotType: HullSlotTypeShield,
}
var BearNeutrinoBarrier = TechHullComponent{Tech: NewTech("Bear Neutrino Barrier", NewCost(4, 0, 4, 8), TechRequirements{TechLevel: TechLevel{Energy: 10}}, 60, TechCategoryShield),

	Mass:         1,
	Shield:       100,
	HullSlotType: HullSlotTypeShield,
}
var GorillaDelagator = TechHullComponent{Tech: NewTech("Gorilla Delagator", NewCost(5, 0, 6, 11), TechRequirements{TechLevel: TechLevel{Energy: 14}}, 70, TechCategoryShield),

	Mass:         1,
	Shield:       175,
	HullSlotType: HullSlotTypeShield,
}
var ElephantHideFortress = TechHullComponent{Tech: NewTech("Elephant Hide Fortress", NewCost(8, 0, 10, 15), TechRequirements{TechLevel: TechLevel{Energy: 18}}, 80, TechCategoryShield),

	Mass:         1,
	Shield:       300,
	HullSlotType: HullSlotTypeShield,
}
var CompletePhaseShield = TechHullComponent{Tech: NewTech("Complete Phase Shield", NewCost(12, 0, 15, 20), TechRequirements{TechLevel: TechLevel{Energy: 22}}, 90, TechCategoryShield),

	Mass:         1,
	Shield:       500,
	HullSlotType: HullSlotTypeShield,
}
var SmallFreighter = TechHull{Tech: NewTech("Small Freighter", NewCost(12, 0, 17, 20), TechRequirements{TechLevel: TechLevel{}}, 10, TechCategoryShipHull),
	Type:          TechHullTypeFreighter,
	Mass:          25,
	Armor:         25,
	FuelCapacity:  130,
	CargoCapacity: 70,
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 1, Required: true},
		{Type: HullSlotTypeScannerElectricalMechanical, Capacity: 1},
		{Type: HullSlotTypeShieldArmor, Capacity: 1},
	},
}
var MediumFreighter = TechHull{Tech: NewTech("Medium Freighter", NewCost(20, 0, 19, 40), TechRequirements{TechLevel: TechLevel{Construction: 3}}, 20, TechCategoryShipHull),
	Type:          TechHullTypeFreighter,
	Mass:          60,
	Armor:         50,
	FuelCapacity:  450,
	CargoCapacity: 210,
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 1, Required: true},
		{Type: HullSlotTypeScannerElectricalMechanical, Capacity: 2},
		{Type: HullSlotTypeShieldArmor, Capacity: 2},
	},
}
var LargeFreighter = TechHull{Tech: NewTech("Large Freighter", NewCost(35, 0, 21, 100), TechRequirements{TechLevel: TechLevel{Construction: 8}}, 30, TechCategoryShipHull),
	Type:          TechHullTypeFreighter,
	Mass:          125,
	Armor:         150,
	FuelCapacity:  2600,
	CargoCapacity: 1200,
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 2, Required: true},
		{Type: HullSlotTypeScannerElectricalMechanical, Capacity: 2},
		{Type: HullSlotTypeShieldArmor, Capacity: 2},
	},
}
var SuperFreighter = TechHull{Tech: NewTech("Super Freighter", NewCost(35, 0, 21, 100), TechRequirements{TechLevel: TechLevel{Construction: 13}, PRTRequired: IS}, 40, TechCategoryShipHull),
	Type:          TechHullTypeFreighter,
	Mass:          175,
	Armor:         400,
	FuelCapacity:  8000,
	CargoCapacity: 3000,
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 3, Required: true},
		{Type: HullSlotTypeScannerElectricalMechanical, Capacity: 3},
		{Type: HullSlotTypeShieldArmor, Capacity: 5},
		{Type: HullSlotTypeElectrical, Capacity: 2},
	},
}
var Scout = TechHull{Tech: NewTech("Scout", NewCost(4, 2, 4, 10), TechRequirements{TechLevel: TechLevel{}}, 50, TechCategoryShipHull),
	Type:           TechHullTypeScout,
	Mass:           8,
	BuiltInScanner: true,
	Armor:          20,
	Initiative:     1,
	FuelCapacity:   50,
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 1, Required: true},
		{Type: HullSlotTypeScanner, Capacity: 1},
		{Type: HullSlotTypeGeneral, Capacity: 1},
	},
}
var Frigate = TechHull{Tech: NewTech("Frigate", NewCost(4, 2, 5, 12), TechRequirements{TechLevel: TechLevel{Construction: 6}}, 60, TechCategoryShipHull),
	Type:           TechHullTypeFighter,
	Mass:           8,
	BuiltInScanner: true,
	Armor:          45,
	Initiative:     4,
	FuelCapacity:   125,
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 1, Required: true},
		{Type: HullSlotTypeScanner, Capacity: 1},
		{Type: HullSlotTypeGeneral, Capacity: 3},
		{Type: HullSlotTypeShieldArmor, Capacity: 2},
	},
}
var Destroyer = TechHull{Tech: NewTech("Destroyer", NewCost(15, 3, 5, 35), TechRequirements{TechLevel: TechLevel{Construction: 3}}, 70, TechCategoryShipHull),
	Type:           TechHullTypeFighter,
	Mass:           30,
	BuiltInScanner: true,
	Armor:          200,
	Initiative:     3,
	FuelCapacity:   280,
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 1, Required: true},
		{Type: HullSlotTypeWeapon, Capacity: 1},
		{Type: HullSlotTypeWeapon, Capacity: 1},
		{Type: HullSlotTypeGeneral, Capacity: 1},
		{Type: HullSlotTypeArmor, Capacity: 2},
		{Type: HullSlotTypeMechanical, Capacity: 1},
		{Type: HullSlotTypeElectrical, Capacity: 1},
	},
}
var Cruiser = TechHull{Tech: NewTech("Cruiser", NewCost(40, 5, 8, 85), TechRequirements{TechLevel: TechLevel{Construction: 9}}, 80, TechCategoryShipHull),
	Type:         TechHullTypeFighter,
	Mass:         90,
	Armor:        700,
	Initiative:   5,
	FuelCapacity: 600,
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 2, Required: true},
		{Type: HullSlotTypeScannerElectricalMechanical, Capacity: 1},
		{Type: HullSlotTypeScannerElectricalMechanical, Capacity: 1},
		{Type: HullSlotTypeWeapon, Capacity: 2},
		{Type: HullSlotTypeWeapon, Capacity: 2},
		{Type: HullSlotTypeGeneral, Capacity: 2},
		{Type: HullSlotTypeShieldArmor, Capacity: 2},
	},
}
var BattleCruiser = TechHull{Tech: NewTech("Battle Cruiser", NewCost(55, 8, 12, 120), TechRequirements{TechLevel: TechLevel{Construction: 9}, PRTRequired: WM}, 90, TechCategoryShipHull),
	Type:         TechHullTypeFighter,
	Mass:         120,
	Armor:        1000,
	Initiative:   5,
	FuelCapacity: 1400,
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 2, Required: true},
		{Type: HullSlotTypeScannerElectricalMechanical, Capacity: 2},
		{Type: HullSlotTypeScannerElectricalMechanical, Capacity: 2},
		{Type: HullSlotTypeWeapon, Capacity: 3},
		{Type: HullSlotTypeWeapon, Capacity: 3},
		{Type: HullSlotTypeGeneral, Capacity: 3},
		{Type: HullSlotTypeShieldArmor, Capacity: 4},
	},
}
var Battleship = TechHull{Tech: NewTech("Battleship", NewCost(120, 25, 20, 225), TechRequirements{TechLevel: TechLevel{Construction: 13}}, 100, TechCategoryShipHull),
	Type:         TechHullTypeFighter,
	Mass:         222,
	Armor:        2000,
	Initiative:   10,
	FuelCapacity: 2800,
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 4, Required: true},
		{Type: HullSlotTypeScannerElectricalMechanical, Capacity: 1},
		{Type: HullSlotTypeShield, Capacity: 8},
		{Type: HullSlotTypeWeapon, Capacity: 6},
		{Type: HullSlotTypeWeapon, Capacity: 6},
		{Type: HullSlotTypeWeapon, Capacity: 2},
		{Type: HullSlotTypeWeapon, Capacity: 2},
		{Type: HullSlotTypeWeapon, Capacity: 4},
		{Type: HullSlotTypeArmor, Capacity: 6},
		{Type: HullSlotTypeElectrical, Capacity: 3},
		{Type: HullSlotTypeElectrical, Capacity: 3},
	},
}
var Dreadnought = TechHull{Tech: NewTech("Dreadnought", NewCost(140, 30, 25, 275), TechRequirements{TechLevel: TechLevel{Construction: 16}, PRTRequired: WM}, 110, TechCategoryShipHull),
	Type:         TechHullTypeFighter,
	Mass:         250,
	Armor:        4500,
	Initiative:   10,
	FuelCapacity: 4500,
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 5, Required: true},
		{Type: HullSlotTypeShieldArmor, Capacity: 4},
		{Type: HullSlotTypeShieldArmor, Capacity: 4},
		{Type: HullSlotTypeWeapon, Capacity: 6},
		{Type: HullSlotTypeWeapon, Capacity: 6},
		{Type: HullSlotTypeElectrical, Capacity: 4},
		{Type: HullSlotTypeElectrical, Capacity: 4},
		{Type: HullSlotTypeWeapon, Capacity: 8},
		{Type: HullSlotTypeWeapon, Capacity: 8},
		{Type: HullSlotTypeArmor, Capacity: 8},
		{Type: HullSlotTypeWeaponShield, Capacity: 5},
		{Type: HullSlotTypeWeaponShield, Capacity: 5},
		{Type: HullSlotTypeGeneral, Capacity: 2},
	},
}
var Privateer = TechHull{Tech: NewTech("Privateer", NewCost(50, 3, 3, 50), TechRequirements{TechLevel: TechLevel{Construction: 4}}, 120, TechCategoryShipHull),
	Type:          TechHullTypeArmedFreighter,
	Mass:          65,
	Armor:         150,
	Initiative:    3,
	FuelCapacity:  650,
	CargoCapacity: 250,
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 1, Required: true},
		{Type: HullSlotTypeShieldArmor, Capacity: 2},
		{Type: HullSlotTypeScannerElectricalMechanical, Capacity: 1},
		{Type: HullSlotTypeGeneral, Capacity: 1},
		{Type: HullSlotTypeGeneral, Capacity: 1},
	},
}
var Rogue = TechHull{Tech: NewTech("Rogue", NewCost(80, 5, 5, 60), TechRequirements{TechLevel: TechLevel{Construction: 8}, PRTRequired: SS}, 130, TechCategoryShipHull),
	Type:          TechHullTypeArmedFreighter,
	Mass:          75,
	Armor:         450,
	Initiative:    4,
	FuelCapacity:  2250,
	CargoCapacity: 500,
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 2, Required: true},
		{Type: HullSlotTypeShieldArmor, Capacity: 3},
		{Type: HullSlotTypeMineElectricalMechanical, Capacity: 2},
		{Type: HullSlotTypeScanner, Capacity: 1},
		{Type: HullSlotTypeGeneral, Capacity: 2},
		{Type: HullSlotTypeGeneral, Capacity: 2},
		{Type: HullSlotTypeMineElectricalMechanical, Capacity: 2},
		{Type: HullSlotTypeElectrical, Capacity: 1},
		{Type: HullSlotTypeElectrical, Capacity: 1},
	},
}
var Galleon = TechHull{Tech: NewTech("Galleon", NewCost(70, 5, 5, 105), TechRequirements{TechLevel: TechLevel{Construction: 11}}, 140, TechCategoryShipHull),
	Type:          TechHullTypeArmedFreighter,
	Mass:          125,
	Armor:         900,
	Initiative:    4,
	FuelCapacity:  2500,
	CargoCapacity: 1000,
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 4, Required: true},
		{Type: HullSlotTypeShieldArmor, Capacity: 2},
		{Type: HullSlotTypeShieldArmor, Capacity: 2},
		{Type: HullSlotTypeGeneral, Capacity: 3},
		{Type: HullSlotTypeGeneral, Capacity: 3},
		{Type: HullSlotTypeMineElectricalMechanical, Capacity: 2},
		{Type: HullSlotTypeMineElectricalMechanical, Capacity: 2},
		{Type: HullSlotTypeScanner, Capacity: 2},
	},
}
var MiniColonyShip = TechHull{Tech: NewTech("Mini-Colony Ship", NewCost(2, 0, 2, 3), TechRequirements{TechLevel: TechLevel{}, PRTRequired: HE}, 150, TechCategoryShipHull),
	Type:          TechHullTypeColonizer,
	Mass:          8,
	Armor:         10,
	FuelCapacity:  150,
	CargoCapacity: 10,
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 1, Required: true},
		{Type: HullSlotTypeMechanical, Capacity: 1},
	},
}
var ColonyShip = TechHull{Tech: NewTech("Colony Ship", NewCost(9, 0, 13, 18), TechRequirements{TechLevel: TechLevel{}}, 160, TechCategoryShipHull),
	Type:          TechHullTypeColonizer,
	Mass:          20,
	Armor:         20,
	FuelCapacity:  200,
	CargoCapacity: 25,
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 1, Required: true},
		{Type: HullSlotTypeMechanical, Capacity: 1},
	},
}
var MiniBomber = TechHull{Tech: NewTech("Mini Bomber", NewCost(18, 5, 9, 32), TechRequirements{TechLevel: TechLevel{Construction: 1}}, 170, TechCategoryShipHull),
	Type:         TechHullTypeBomber,
	Mass:         28,
	Armor:        50,
	FuelCapacity: 120,
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 1, Required: true},
		{Type: HullSlotTypeBomb, Capacity: 2},
	},
}
var B17Bomber = TechHull{Tech: NewTech("B-17 Bomber", NewCost(55, 10, 10, 150), TechRequirements{TechLevel: TechLevel{Construction: 6}}, 180, TechCategoryShipHull),
	Type:         TechHullTypeBomber,
	Mass:         69,
	Armor:        175,
	FuelCapacity: 400,
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 2, Required: true},
		{Type: HullSlotTypeBomb, Capacity: 4},
		{Type: HullSlotTypeBomb, Capacity: 4},
		{Type: HullSlotTypeScannerElectricalMechanical, Capacity: 1},
	},
}
var StealthBomber = TechHull{Tech: NewTech("Stealth Bomber", NewCost(55, 10, 15, 175), TechRequirements{TechLevel: TechLevel{Construction: 8}, PRTRequired: SS}, 190, TechCategoryShipHull),
	Type:         TechHullTypeBomber,
	FuelCapacity: 750,
	Armor:        225,
	Mass:         70,
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 2, Required: true},
		{Type: HullSlotTypeBomb, Capacity: 4},
		{Type: HullSlotTypeBomb, Capacity: 4},
		{Type: HullSlotTypeScannerElectricalMechanical, Capacity: 1},
		{Type: HullSlotTypeElectrical, Capacity: 3},
	},
}
var B52Bomber = TechHull{Tech: NewTech("B-52 Bomber", NewCost(90, 15, 10, 280), TechRequirements{TechLevel: TechLevel{Construction: 15}}, 200, TechCategoryShipHull),
	Type:         TechHullTypeBomber,
	FuelCapacity: 750,
	Armor:        450,
	Mass:         110,
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 2, Required: true},
		{Type: HullSlotTypeBomb, Capacity: 4},
		{Type: HullSlotTypeBomb, Capacity: 4},
		{Type: HullSlotTypeBomb, Capacity: 4},
		{Type: HullSlotTypeBomb, Capacity: 4},
		{Type: HullSlotTypeScannerElectricalMechanical, Capacity: 2},
		{Type: HullSlotTypeShield, Capacity: 2},
	},
}
var MidgetMiner = TechHull{Tech: NewTech("Midget Miner", NewCost(10, 0, 3, 20), TechRequirements{TechLevel: TechLevel{}, LRTsRequired: ARM}, 210, TechCategoryShipHull),
	Type:         TechHullTypeMiner,
	FuelCapacity: 210,
	Armor:        100,
	Mass:         10,
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 1, Required: true},
		{Type: HullSlotTypeMining, Capacity: 2},
	},
}
var MiniMiner = TechHull{Tech: NewTech("Mini-Miner", NewCost(25, 0, 6, 50), TechRequirements{TechLevel: TechLevel{Construction: 2}}, 220, TechCategoryShipHull),
	Type:         TechHullTypeMiner,
	Mass:         80,
	Armor:        130,
	FuelCapacity: 210,
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 1, Required: true},
		{Type: HullSlotTypeScannerElectricalMechanical},
		{Type: HullSlotTypeMining},
		{Type: HullSlotTypeMining},
	},
}
var Miner = TechHull{Tech: NewTech("Miner", NewCost(32, 0, 6, 110), TechRequirements{TechLevel: TechLevel{Construction: 6}, LRTsRequired: ARM}, 230, TechCategoryShipHull),
	Type:         TechHullTypeMiner,
	FuelCapacity: 500,
	Armor:        475,
	Mass:         110,
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 2, Required: true},
		{Type: HullSlotTypeArmorScannerElectricalMechanical, Capacity: 2},
		{Type: HullSlotTypeMining, Capacity: 2},
		{Type: HullSlotTypeMining, Capacity: 1},
		{Type: HullSlotTypeMining, Capacity: 2},
		{Type: HullSlotTypeMining, Capacity: 1},
	},
}
var MaxiMiner = TechHull{Tech: NewTech("Maxi-Miner", NewCost(32, 0, 6, 140), TechRequirements{TechLevel: TechLevel{Construction: 11}, LRTsDenied: OBRM}, 240, TechCategoryShipHull),
	Type:         TechHullTypeMiner,
	FuelCapacity: 850,
	Armor:        1400,
	Mass:         110,
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 3, Required: true},
		{Type: HullSlotTypeArmorScannerElectricalMechanical, Capacity: 2},
		{Type: HullSlotTypeMining, Capacity: 4},
		{Type: HullSlotTypeMining, Capacity: 1},
		{Type: HullSlotTypeMining, Capacity: 4},
		{Type: HullSlotTypeMining, Capacity: 1},
	},
}
var UltraMiner = TechHull{Tech: NewTech("Ultra-Miner", NewCost(30, 0, 6, 130), TechRequirements{TechLevel: TechLevel{Construction: 14}, LRTsRequired: ARM}, 250, TechCategoryShipHull),
	Type:         TechHullTypeMiner,
	FuelCapacity: 1300,
	Armor:        1500,
	Mass:         100,
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 2, Required: true},
		{Type: HullSlotTypeArmorScannerElectricalMechanical, Capacity: 3},
		{Type: HullSlotTypeMining, Capacity: 4},
		{Type: HullSlotTypeMining, Capacity: 2},
		{Type: HullSlotTypeMining, Capacity: 4},
		{Type: HullSlotTypeMining, Capacity: 2},
	},
}
var FuelTransport = TechHull{Tech: NewTech("Fuel Transport", NewCost(10, 0, 5, 50), TechRequirements{TechLevel: TechLevel{Construction: 4}, PRTRequired: IS}, 260, TechCategoryShipHull),
	Type:         TechHullTypeFuelTransport,
	Mass:         12,
	Armor:        5,
	FuelCapacity: 750,
	RepairBonus:  .05, // +5% repair
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 1, Required: true},
		{Type: HullSlotTypeShield},
	},
}
var SuperFuelXport = TechHull{Tech: NewTech("Super Fuel Xport", NewCost(20, 0, 8, 70), TechRequirements{TechLevel: TechLevel{Construction: 7}, PRTRequired: IS}, 270, TechCategoryShipHull),
	Type:         TechHullTypeFuelTransport,
	Mass:         111,
	Armor:        12,
	FuelCapacity: 2250,
	RepairBonus:  .1, // +10% repair
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 2, Required: true},
		{Type: HullSlotTypeShield, Capacity: 2},
		{Type: HullSlotTypeScanner},
	},
}
var MiniMineLayer = TechHull{Tech: NewTech("Mini Mine Layer", NewCost(8, 2, 5, 20), TechRequirements{TechLevel: TechLevel{}, PRTRequired: SD}, 280, TechCategoryShipHull),
	Type:                  TechHullTypeMineLayer,
	Mass:                  10,
	Armor:                 60,
	FuelCapacity:          400,
	MineLayingFactor:      2,
	ImmuneToOwnDetonation: true,
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 1, Required: true},
		{Type: HullSlotTypeMineLayer, Capacity: 2},
		{Type: HullSlotTypeMineLayer, Capacity: 2},
		{Type: HullSlotTypeScannerElectricalMechanical, Capacity: 1},
	},
}
var SuperMineLayer = TechHull{Tech: NewTech("Super Mine Layer", NewCost(20, 3, 9, 30), TechRequirements{TechLevel: TechLevel{Construction: 15}, PRTRequired: SD}, 290, TechCategoryShipHull),
	Type:                  TechHullTypeMineLayer,
	FuelCapacity:          2200,
	Armor:                 1200,
	Mass:                  30,
	MineLayingFactor:      2,
	ImmuneToOwnDetonation: true,
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 1, Required: true},
		{Type: HullSlotTypeMineLayer, Capacity: 8},
		{Type: HullSlotTypeMineLayer, Capacity: 8},
		{Type: HullSlotTypeShieldArmor, Capacity: 4},
		{Type: HullSlotTypeScannerElectricalMechanical, Capacity: 3},
		{Type: HullSlotTypeMineElectricalMechanical, Capacity: 3},
	},
}
var Nubian = TechHull{Tech: NewTech("Nubian", NewCost(75, 12, 12, 150), TechRequirements{TechLevel: TechLevel{Construction: 26}}, 300, TechCategoryShipHull),
	Type:         TechHullTypeFighter,
	FuelCapacity: 5000,
	Armor:        5000,
	Initiative:   2,
	Mass:         100,
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 3, Required: true},
		{Type: HullSlotTypeGeneral, Capacity: 3},
		{Type: HullSlotTypeGeneral, Capacity: 3},
		{Type: HullSlotTypeGeneral, Capacity: 3},
		{Type: HullSlotTypeGeneral, Capacity: 3},
		{Type: HullSlotTypeGeneral, Capacity: 3},
		{Type: HullSlotTypeGeneral, Capacity: 3},
		{Type: HullSlotTypeGeneral, Capacity: 3},
		{Type: HullSlotTypeGeneral, Capacity: 3},
		{Type: HullSlotTypeGeneral, Capacity: 3},
		{Type: HullSlotTypeGeneral, Capacity: 3},
		{Type: HullSlotTypeGeneral, Capacity: 3},
		{Type: HullSlotTypeGeneral, Capacity: 3},
	},
}
var MetaMorph = TechHull{Tech: NewTech("Meta Morph", NewCost(50, 12, 12, 120), TechRequirements{TechLevel: TechLevel{Construction: 10}, PRTRequired: HE}, 310, TechCategoryShipHull),
	Type:          TechHullTypeArmedFreighter,
	Mass:          85,
	Armor:         500,
	Initiative:    2,
	FuelCapacity:  700,
	CargoCapacity: 500,
	Slots: []TechHullSlot{
		{Type: HullSlotTypeEngine, Capacity: 3, Required: true},
		{Type: HullSlotTypeGeneral, Capacity: 8},
		{Type: HullSlotTypeGeneral, Capacity: 2},
		{Type: HullSlotTypeGeneral, Capacity: 2},
		{Type: HullSlotTypeGeneral, Capacity: 1},
		{Type: HullSlotTypeGeneral, Capacity: 2},
		{Type: HullSlotTypeGeneral, Capacity: 2},
	},
}
var OrbitalFort = TechHull{Tech: NewTech("Orbital Fort", NewCost(24, 0, 34, 80), TechRequirements{TechLevel: TechLevel{}}, 10, TechCategoryStarbaseHull),
	Type:                    TechHullTypeStarbase,
	SpaceDock:               0,
	Armor:                   100,
	Initiative:              10,
	RangeBonus:              1,
	Starbase:                true,
	OrbitalConstructionHull: true,
	RepairBonus:             .03, // 8% total repair rate
	Slots: []TechHullSlot{
		{Type: HullSlotTypeOrbitalElectrical, Capacity: 1},
		{Type: HullSlotTypeWeapon, Capacity: 12},
		{Type: HullSlotTypeShieldArmor, Capacity: 12},
		{Type: HullSlotTypeWeapon, Capacity: 12},
		{Type: HullSlotTypeShieldArmor, Capacity: 12},
	},
}

var SpaceDock = TechHull{Tech: NewTech("Space Dock", NewCost(40, 10, 50, 200), TechRequirements{TechLevel: TechLevel{Construction: 4}, LRTsRequired: ISB}, 20, TechCategoryStarbaseHull),
	Type:        TechHullTypeStarbase,
	SpaceDock:   200,
	Armor:       250,
	Initiative:  12,
	RangeBonus:  1,
	Starbase:    true,
	RepairBonus: .03, // 8% total repair rate
	Slots: []TechHullSlot{
		{Type: HullSlotTypeOrbitalElectrical, Capacity: 1},
		{Type: HullSlotTypeWeapon, Capacity: 16},
		{Type: HullSlotTypeShieldArmor, Capacity: 24},
		{Type: HullSlotTypeWeapon, Capacity: 16},
		{Type: HullSlotTypeShield, Capacity: 24},
		{Type: HullSlotTypeElectrical, Capacity: 2},
		{Type: HullSlotTypeElectrical, Capacity: 2},
		{Type: HullSlotTypeWeapon, Capacity: 16},
	},
}
var SpaceStation = TechHull{Tech: NewTech("Space Station", NewCost(120, 80, 250, 600), TechRequirements{TechLevel: TechLevel{}}, 30, TechCategoryStarbaseHull),
	Type:        TechHullTypeStarbase,
	SpaceDock:   UnlimitedSpaceDock,
	Armor:       500,
	Initiative:  14,
	RangeBonus:  1,
	Starbase:    true,
	RepairBonus: .15, // 20% total repair rate
	Slots: []TechHullSlot{
		{Type: HullSlotTypeOrbitalElectrical, Capacity: 1},
		{Type: HullSlotTypeWeapon, Capacity: 16},
		{Type: HullSlotTypeShield, Capacity: 16},
		{Type: HullSlotTypeWeapon, Capacity: 16},
		{Type: HullSlotTypeShieldArmor, Capacity: 16},
		{Type: HullSlotTypeShield, Capacity: 16},
		{Type: HullSlotTypeElectrical, Capacity: 3},
		{Type: HullSlotTypeWeapon, Capacity: 16},
		{Type: HullSlotTypeElectrical, Capacity: 3},
		{Type: HullSlotTypeWeapon, Capacity: 16},
		{Type: HullSlotTypeOrbitalElectrical, Capacity: 1},
		{Type: HullSlotTypeShieldArmor, Capacity: 16},
	},
}
var UltraStation = TechHull{Tech: NewTech("Ultra Station", NewCost(120, 80, 300, 600), TechRequirements{TechLevel: TechLevel{Construction: 12}, LRTsRequired: ISB}, 40, TechCategoryStarbaseHull),
	Type:                     TechHullTypeStarbase,
	SpaceDock:                UnlimitedSpaceDock,
	Armor:                    1000,
	Initiative:               16,
	Starbase:                 true,
	RepairBonus:              .15, // 20% total repair rate
	InnateScanRangePenFactor: .5,  // AR races get half innate scanning range for pen scanning
	Slots: []TechHullSlot{
		{Type: HullSlotTypeOrbitalElectrical, Capacity: 1},
		{Type: HullSlotTypeWeapon, Capacity: 16},
		{Type: HullSlotTypeElectrical, Capacity: 3},
		{Type: HullSlotTypeWeapon, Capacity: 6},
		{Type: HullSlotTypeShield, Capacity: 20},
		{Type: HullSlotTypeShield, Capacity: 20},
		{Type: HullSlotTypeElectrical, Capacity: 3},
		{Type: HullSlotTypeWeapon, Capacity: 16},
		{Type: HullSlotTypeElectrical, Capacity: 3},
		{Type: HullSlotTypeWeapon, Capacity: 16},
		{Type: HullSlotTypeOrbitalElectrical, Capacity: 1},
		{Type: HullSlotTypeShieldArmor, Capacity: 20},
		{Type: HullSlotTypeWeapon, Capacity: 16},
		{Type: HullSlotTypeShieldArmor, Capacity: 20},
		{Type: HullSlotTypeElectrical, Capacity: 3},
		{Type: HullSlotTypeWeapon, Capacity: 16},
	},
}
var DeathStar = TechHull{Tech: NewTech("Death Star", NewCost(120, 80, 350, 750), TechRequirements{TechLevel: TechLevel{Construction: 17}, PRTRequired: AR}, 50, TechCategoryStarbaseHull),
	Type:                     TechHullTypeStarbase,
	SpaceDock:                UnlimitedSpaceDock,
	Armor:                    1500,
	Initiative:               18,
	Starbase:                 true,
	RepairBonus:              .15, // 20% total repair rate
	InnateScanRangePenFactor: .5,  // AR races get half innate scanning range for pen scanning
	Slots: []TechHullSlot{
		{Type: HullSlotTypeOrbitalElectrical, Capacity: 1},
		{Type: HullSlotTypeWeapon, Capacity: 32},
		{Type: HullSlotTypeElectrical, Capacity: 4},
		{Type: HullSlotTypeElectrical, Capacity: 4},
		{Type: HullSlotTypeShield, Capacity: 20},
		{Type: HullSlotTypeShield, Capacity: 20},
		{Type: HullSlotTypeElectrical, Capacity: 4},
		{Type: HullSlotTypeWeapon, Capacity: 32},
		{Type: HullSlotTypeElectrical, Capacity: 4},
		{Type: HullSlotTypeWeapon, Capacity: 32},
		{Type: HullSlotTypeOrbitalElectrical, Capacity: 1},
		{Type: HullSlotTypeShieldArmor, Capacity: 20},
		{Type: HullSlotTypeElectrical, Capacity: 4},
		{Type: HullSlotTypeShieldArmor, Capacity: 20},
		{Type: HullSlotTypeElectrical, Capacity: 4},
		{Type: HullSlotTypeWeapon, Capacity: 32},
	},
}

func TechEngines() []TechEngine {
	return []TechEngine{
		SettlersDelight,
		QuickJump5,
		FuelMizer,
		LongHump6,
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

func TechTerraforms() []TechTerraform {
	return []TechTerraform{
		TotalTerraform3,
		TotalTerraform5,
		TotalTerraform7,
		TotalTerraform10,
		TotalTerraform15,
		TotalTerraform20,
		TotalTerraform25,
		TotalTerraform30,
		GravityTerraform3,
		GravityTerraform7,
		GravityTerraform11,
		GravityTerraform15,
		TempTerraform3,
		TempTerraform7,
		TempTerraform11,
		TempTerraform15,
		RadiationTerraform3,
		RadiationTerraform7,
		RadiationTerraform11,
		RadiationTerraform15,
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

func TechDefenses() []TechDefense {

	return []TechDefense{
		SDI,
		MissileBattery,
		LaserBattery,
		PlanetaryShield,
		NeutronShield,
	}
}

func TechHulls() []TechHull {
	return []TechHull{
		SmallFreighter,
		MediumFreighter,
		LargeFreighter,
		SuperFreighter,
		Scout,
		Frigate,
		Destroyer,
		Cruiser,
		BattleCruiser,
		Battleship,
		Dreadnought,
		Privateer,
		Rogue,
		Galleon,
		MiniColonyShip,
		ColonyShip,
		MiniBomber,
		B17Bomber,
		StealthBomber,
		B52Bomber,
		MidgetMiner,
		MiniMiner,
		Miner,
		MaxiMiner,
		UltraMiner,
		FuelTransport,
		SuperFuelXport,
		MiniMineLayer,
		SuperMineLayer,
		Nubian,
		MetaMorph,
		OrbitalFort,
		SpaceDock,
		SpaceStation,
		UltraStation,
		DeathStar,
	}
}

func TechHullComponents() []TechHullComponent {

	return []TechHullComponent{
		Stargate100_250,
		StargateAny_300,
		Stargate150_600,
		Stargate300_500,
		Stargate100_Any,
		StargateAny_800,
		StargateAny_Any,
		MassDriver5,
		MassDriver6,
		MassDriver7,
		SuperDriver8,
		SuperDriver9,
		UltraDriver10,
		UltraDriver11,
		UltraDriver12,
		UltraDriver13,
		RoboMiner,
		RoboMaxiMiner,
		RoboMidgetMiner,
		RoboMiniMiner,
		RoboSuperMiner,
		RoboUltraMiner,
		OrbitalAdjuster,
		LadyFingerBomb,
		BlackCatBomb,
		M70Bomb,
		M80Bomb,
		CherryBomb,
		LBU17Bomb,
		LBU32Bomb,
		LBU74Bomb,
		RetroBomb,
		SmartBomb,
		NeutronBomb,
		EnrichedNeutronBomb,
		PeerlessBomb,
		AnnihilatorBomb,
		BatScanner,
		RhinoScanner,
		MoleScanner,
		DNAScanner,
		PossumScanner,
		PickPocketScanner,
		ChameleonScanner,
		FerretScanner,
		DolphinScanner,
		GazelleScanner,
		RNAScanner,
		CheetahScanner,
		ElephantScanner,
		EagleEyeScanner,
		RobberBaronScanner,
		PeerlessScanner,
		Tritanium,
		Crobmnium,
		Carbonic,
		Strobnium,
		Organic,
		Kelarium,
		FieldedKelarium,
		DepletedNeutronium,
		Neutronium,
		Valanium,
		Superlatanium,
		TransportCloaking,
		StealthCloak,
		SuperStealthCloak,
		UltraStealthCloak,
		BattleComputer,
		BattleSuperComputer,
		BattleNexus,
		Jammer10,
		Jammer20,
		Jammer30,
		Jammer50,
		EnergyCapacitor,
		FluxCapacitor,
		EnergyDampener,
		TachyonDetector,
		AntiMatterGenerator,
		MineDispenser40,
		MineDispenser50,
		MineDispenser80,
		MineDispenser130,
		HeavyDispenser50,
		HeavyDispenser110,
		HeavyDispenser200,
		SpeedTrap20,
		SpeedTrap30,
		SpeedTrap50,
		ColonizationModule,
		OrbitalConstructionModule,
		CargoPod,
		SuperCargoPod,
		FuelTank,
		SuperFuelTank,
		ManeuveringJet,
		Overthruster,
		BeamDeflector,
		Laser,
		XRayLaser,
		MiniGun,
		YakimoraLightPhaser,
		Blackjack,
		PhaserBazooka,
		PulsedSapper,
		ColloidalPhaser,
		GatlingGun,
		MiniBlaster,
		Bludgeon,
		MarkIVBlaster,
		PhasedSapper,
		HeavyBlaster,
		GatlingNeutrinoCannon,
		MyopicDisruptor,
		Blunderbuss,
		Disruptor,
		SyncroSapper,
		MegaDisruptor,
		BigMuthaCannon,
		StreamingPulverizer,
		AntiMatterPulverizer,
		AlphaTorpedo,
		BetaTorpedo,
		DeltaTorpedo,
		EpsilonTorpedo,
		RhoTorpedo,
		UpsilonTorpedo,
		OmegaTorpedo,
		JihadMissile,
		JuggernautMissile,
		DoomsdayMissile,
		ArmageddonMissile,
		MoleSkinShield,
		CowHideShield,
		WolverineDiffuseShield,
		CrobySharmor,
		ShadowShield,
		BearNeutrinoBarrier,
		GorillaDelagator,
		ElephantHideFortress,
		CompletePhaseShield,
	}
}
