package cs

import (
	"fmt"
	"math"

	"github.com/rs/zerolog/log"
)

// A user can have multiple races stored in the database. Each time a game is created, a Race is copied
// into the Player object and stored separately (so changes to the User's race don't impact running games)
type Race struct {
	DBObject
	UserID                int64                 `json:"userId,omitempty"`
	Name                  string                `json:"name,omitempty"`
	PluralName            string                `json:"pluralName,omitempty"`
	SpendLeftoverPointsOn SpendLeftoverPointsOn `json:"spendLeftoverPointsOn,omitempty"`
	PRT                   PRT                   `json:"prt,omitempty"`
	LRTs                  Bitmask               `json:"lrts,omitempty"`
	HabLow                Hab                   `json:"habLow,omitempty"`
	HabHigh               Hab                   `json:"habHigh,omitempty"`
	GrowthRate            int                   `json:"growthRate,omitempty"`
	PopEfficiency         int                   `json:"popEfficiency,omitempty"`
	FactoryOutput         int                   `json:"factoryOutput,omitempty"`
	FactoryCost           int                   `json:"factoryCost,omitempty"`
	NumFactories          int                   `json:"numFactories,omitempty"`
	FactoriesCostLess     bool                  `json:"factoriesCostLess,omitempty"`
	ImmuneGrav            bool                  `json:"immuneGrav,omitempty"`
	ImmuneTemp            bool                  `json:"immuneTemp,omitempty"`
	ImmuneRad             bool                  `json:"immuneRad,omitempty"`
	MineOutput            int                   `json:"mineOutput,omitempty"`
	MineCost              int                   `json:"mineCost,omitempty"`
	NumMines              int                   `json:"numMines,omitempty"`
	ResearchCost          ResearchCost          `json:"researchCost,omitempty"`
	TechsStartHigh        bool                  `json:"techsStartHigh,omitempty"`
	Spec                  RaceSpec              `json:"spec,omitempty"`
}

type ResearchCostLevel string

const (
	ResearchCostExtra    ResearchCostLevel = "Extra"
	ResearchCostStandard ResearchCostLevel = "Standard"
	ResearchCostLess     ResearchCostLevel = "Less"
)

type SpendLeftoverPointsOn string

const (
	SpendLeftoverPointsOnSurfaceMinerals       SpendLeftoverPointsOn = "SurfaceMinerals"
	SpendLeftoverPointsOnMineralConcentrations SpendLeftoverPointsOn = "MineralConcentrations"
	SpendLeftoverPointsOnMines                 SpendLeftoverPointsOn = "Mines"
	SpendLeftoverPointsOnFactories             SpendLeftoverPointsOn = "Factories"
	SpendLeftoverPointsOnDefenses              SpendLeftoverPointsOn = "Defenses"
)

type ResearchCost struct {
	Energy        ResearchCostLevel `json:"energy,omitempty"`
	Weapons       ResearchCostLevel `json:"weapons,omitempty"`
	Propulsion    ResearchCostLevel `json:"propulsion,omitempty"`
	Construction  ResearchCostLevel `json:"construction,omitempty"`
	Electronics   ResearchCostLevel `json:"electronics,omitempty"`
	Biotechnology ResearchCostLevel `json:"biotechnology,omitempty"`
}

func (rc ResearchCost) Get(field TechField) ResearchCostLevel {
	switch field {
	case Energy:
		return rc.Energy
	case Weapons:
		return rc.Weapons
	case Propulsion:
		return rc.Propulsion
	case Construction:
		return rc.Construction
	case Electronics:
		return rc.Electronics
	case Biotechnology:
		return rc.Biotechnology
	}

	log.Error().Msgf("invalid field %s to get ResearchCost", field)
	return ResearchCostStandard
}

type RaceSpec struct {
	MiniaturizationSpec
	ScannerSpec
	HabCenter                        Hab                    `json:"habCenter,omitempty"`
	Costs                            map[QueueItemType]Cost `json:"costs,omitempty"`
	StartingTechLevels               TechLevel              `json:"startingTechLevels,omitempty"`
	StartingPlanets                  []StartingPlanet       `json:"startingPlanets,omitempty"`
	TechCostOffset                   TechCostOffset         `json:"techCostOffset,omitempty"`
	MineralsPerSingleMineralPacket   int                    `json:"mineralsPerSingleMineralPacket,omitempty"`
	MineralsPerMixedMineralPacket    int                    `json:"mineralsPerMixedMineralPacket,omitempty"`
	PacketResourceCost               int                    `json:"packetResourceCost,omitempty"`
	PacketMineralCostFactor          float64                `json:"packetMineralCostFactor,omitempty"`
	PacketReceiverFactor             float64                `json:"packetReceiverFactor,omitempty"`
	PacketDecayFactor                float64                `json:"packetDecayFactor,omitempty"`
	PacketOverSafeWarpPenalty        int                    `json:"packetOverSafeWarpPenalty,omitempty"`
	PacketBuiltInScanner             bool                   `json:"packetBuiltInScanner,omitempty"`
	DetectPacketDestinationStarbases bool                   `json:"detectPacketDestinationStarbases,omitempty"`
	DetectAllPackets                 bool                   `json:"detectAllPackets,omitempty"`
	PacketTerraformChance            float64                `json:"packetTerraformChance,omitempty"`
	PacketPermaformChance            float64                `json:"packetPermaformChance,omitempty"`
	PacketPermaTerraformSizeUnit     int                    `json:"packetPermaTerraformSizeUnit,omitempty"`
	CanGateCargo                     bool                   `json:"canGateCargo,omitempty"`
	CanDetectStargatePlanets         bool                   `json:"canDetectStargatePlanets,omitempty"`
	ShipsVanishInVoid                bool                   `json:"shipsVanishInVoid,omitempty"`
	TechsCostExtraLevel              int                    `json:"techsCostExtraLevel,omitempty"`
	FreighterGrowthFactor            float64                `json:"freighterGrowthFactor,omitempty"`
	GrowthFactor                     float64                `json:"growthFactor,omitempty"`
	MaxPopulationOffset              float64                `json:"maxPopulationOffset,omitempty"`
	BuiltInCloakUnits                int                    `json:"builtInCloakUnits,omitempty"`
	StealsResearch                   StealsResearch         `json:"stealsResearch,omitempty"`
	FreeCargoCloaking                bool                   `json:"freeCargoCloaking,omitempty"`
	MineFieldsAreScanners            bool                   `json:"mineFieldsAreScanners,omitempty"`
	MineFieldRateMoveFactor          float64                `json:"mineFieldRateMoveFactor,omitempty"`
	MineFieldSafeWarpBonus           int                    `json:"mineFieldSafeWarpBonus,omitempty"`
	MineFieldMinDecayFactor          float64                `json:"mineFieldMinDecayFactor,omitempty"`
	MineFieldBaseDecayRate           float64                `json:"mineFieldBaseDecayRate,omitempty"`
	MineFieldPlanetDecayRate         float64                `json:"mineFieldPlanetDecayRate,omitempty"`
	MineFieldMaxDecayRate            float64                `json:"mineFieldMaxDecayRate,omitempty"`
	CanDetonateMineFields            bool                   `json:"canDetonateMineFields,omitempty"`
	MineFieldDetonateDecayRate       float64                `json:"mineFieldDetonateDecayRate,omitempty"`
	DiscoverDesignOnScan             bool                   `json:"discoverDesignOnScan,omitempty"`
	CanRemoteMineOwnPlanets          bool                   `json:"canRemoteMineOwnPlanets,omitempty"`
	InvasionAttackBonus              float64                `json:"invasionAttackBonus,omitempty"`
	InvasionDefendBonus              float64                `json:"invasionDefendBonus,omitempty"`
	MovementBonus                    int                    `json:"movementBonus,omitempty"`
	Instaforming                     bool                   `json:"instaforming,omitempty"`
	PermaformChance                  float64                `json:"permaformChance,omitempty"`
	PermaformPopulation              int                    `json:"permaformPopulation,omitempty"`
	RepairFactor                     float64                `json:"repairFactor,omitempty"`
	StarbaseRepairFactor             float64                `json:"starbaseRepairFactor,omitempty"`
	InnateMining                     bool                   `json:"innateMining,omitempty"`
	InnateResources                  bool                   `json:"innateResources,omitempty"`
	InnateScanner                    bool                   `json:"innateScanner,omitempty"`
	InnatePopulationFactor           float64                `json:"innatePopulationFactor,omitempty"`
	CanBuildDefenses                 bool                   `json:"canBuildDefenses,omitempty"`
	LivesOnStarbases                 bool                   `json:"livesOnStarbases,omitempty"`
	FuelEfficiencyOffset             float64                `json:"fuelEfficiencyOffset,omitempty"`
	TerraformCostOffset              Cost                   `json:"terraformCostOffset,omitempty"`
	MineralAlchemyCostOffset         int                    `json:"mineralAlchemyCostOffset,omitempty"`
	ScrapMineralOffset               float64                `json:"scrapMineralOffset,omitempty"`
	ScrapMineralOffsetStarbase       float64                `json:"scrapMineralOffsetStarbase,omitempty"`
	ScrapResourcesOffset             float64                `json:"scrapResourcesOffset,omitempty"`
	ScrapResourcesOffsetStarbase     float64                `json:"scrapResourcesOffsetStarbase,omitempty"`
	StartingPopulationFactor         float64                `json:"startingPopulationFactor,omitempty"`
	StarbaseBuiltInCloakUnits        int                    `json:"starbaseBuiltInCloakUnits,omitempty"`
	StarbaseCostFactor               float64                `json:"starbaseCostFactor,omitempty"`
	ResearchFactor                   float64                `json:"researchFactor,omitempty"`
	ResearchSplashDamage             float64                `json:"researchSplashDamage,omitempty"`
	ArmorStrengthFactor              float64                `json:"armorStrengthFactor,omitempty"`
	ShieldStrengthFactor             float64                `json:"shieldStrengthFactor,omitempty"`
	ShieldRegenerationRate           float64                `json:"shieldRegenerationRate,omitempty"`
	EngineFailureRate                float64                `json:"engineFailureRate,omitempty"`
	EngineReliableSpeed              int                    `json:"engineReliableSpeed,omitempty"`
}

type MiniaturizationSpec struct {
	NewTechCostFactor       float64 `json:"newTechCostFactor,omitempty"`
	MiniaturizationMax      float64 `json:"miniaturizationMax,omitempty"`
	MiniaturizationPerLevel float64 `json:"miniaturizationPerLevel,omitempty"`
}

type ScannerSpec struct {
	BuiltInScannerMultiplier int     `json:"builtInScannerMultiplier,omitempty"`
	NoAdvancedScanners       bool    `json:"noAdvancedScanners,omitempty"`
	ScanRangeFactor          float64 `json:"scanRangeFactor,omitempty"`
}

type PRT string

const (
	/// This is only for tech requirements
	PRTNone PRT = ""

	/// Hyper Expansion
	HE PRT = "HE"

	/// Super Stealth
	SS PRT = "SS"

	/// Warmonger
	WM PRT = "WM"

	/// Claim Adjuster
	CA PRT = "CA"

	/// Inner Strength
	IS PRT = "IS"

	/// Space Demolition
	SD PRT = "SD"

	/// Packet Physics
	PP PRT = "PP"

	/// Interstellar Traveler
	IT PRT = "IT"

	/// Alternate Reality
	AR PRT = "AR"

	/// Jack of All Trades
	JoaT PRT = "JoaT"
)

var PRTs = [10]PRT{
	HE,
	SS,
	WM,
	CA,
	IS,
	SD,
	PP,
	IT,
	AR,
	JoaT,
}

type Bitmask uint32

func (mask Bitmask) countBits() int {
	count := 0

	for mask > 0 {
		count += int(mask & 1)
		mask >>= 1
	}

	return count
}

type LRT Bitmask

const (
	// Only used for TechRequirements
	LRTNone = 0

	// Improved Fuel Efficiency
	IFE LRT = 1 << (iota - 1)

	// Total Terraforming
	TT

	// Advanced Remote Mining
	ARM

	// Improved Starbases
	ISB

	// Generalized Research
	GR

	// Ultimate Recycling
	UR

	// No Ramscoop Engines
	NRSE

	// Only Basic Remote Mining
	OBRM

	// No Advanced Scanners
	NAS

	// Low Starting Population
	LSP

	// Bleeding Edge Technology
	BET

	// Regenerating Shields
	RS

	// Mineral Alchemy
	MA

	// Cheap Engines
	CE
)

var LRTs = []LRT{
	IFE,
	TT,
	ARM,
	ISB,
	GR,
	UR,
	NRSE,
	OBRM,
	NAS,
	LSP,
	BET,
	RS,
	MA,
	CE,
}

func (r *Race) IsImmune(habType HabType) bool {
	switch habType {
	case Grav:
		return r.ImmuneGrav
	case Temp:
		return r.ImmuneTemp
	case Rad:
		return r.ImmuneRad
	}

	return false
}

func NewRace() *Race {
	return &Race{
		Name:       "Humanoid",
		PluralName: "Humanoids",
		PRT:        JoaT,
		LRTs:       LRTNone,
		HabLow: Hab{
			Grav: 15,
			Temp: 15,
			Rad:  15,
		},
		HabHigh: Hab{
			Grav: 85,
			Temp: 85,
			Rad:  85,
		},
		GrowthRate:    15,
		PopEfficiency: 10,
		FactoryOutput: 10,
		FactoryCost:   10,
		NumFactories:  10,
		MineOutput:    10,
		MineCost:      5,
		NumMines:      10,
		ResearchCost: ResearchCost{
			Energy:        ResearchCostStandard,
			Weapons:       ResearchCostStandard,
			Propulsion:    ResearchCostStandard,
			Construction:  ResearchCostStandard,
			Electronics:   ResearchCostStandard,
			Biotechnology: ResearchCostStandard,
		},
	}
}

func (r *Race) WithUserID(userID int64) *Race {
	r.UserID = userID
	return r
}

func (r *Race) WithLRT(lrt LRT) *Race {
	r.LRTs |= Bitmask(lrt)
	return r
}

func (r *Race) WithPRT(prt PRT) *Race {
	r.PRT = prt
	return r
}

func (r Race) WithName(name string) Race {
	r.Name = name
	return r
}

func (r *Race) WithPluralName(pluralname string) *Race {
	r.PluralName = pluralname
	return r
}

func (r *Race) WithGrowthRate(growthRate int) *Race {
	r.GrowthRate = growthRate
	return r
}

func (r *Race) withImmuneGrav(immune bool) *Race {
	r.ImmuneGrav = immune
	return r
}

func (r *Race) withImmuneTemp(immune bool) *Race {
	r.ImmuneTemp = immune
	return r
}

func (r *Race) withImmuneRad(immune bool) *Race {
	r.ImmuneRad = immune
	return r
}

func (r *Race) withResearchCost(researchCost ResearchCost) *Race {
	r.ResearchCost = researchCost
	return r
}

func (r *Race) WithSpec(rules *Rules) *Race {
	r.Spec = computeRaceSpec(r, rules)
	return r
}

func Humanoids() Race {
	return *NewRace()
}

func Rabbitoids() Race {
	return Race{
		Name:              "Rabbitoid",
		PluralName:        "Rabbitoids",
		PRT:               IT,
		LRTs:              Bitmask(IFE) | Bitmask(TT) | Bitmask(CE) | Bitmask(NAS),
		HabLow:            Hab{10, 35, 13},
		HabHigh:           Hab{56, 81, 53},
		GrowthRate:        20,
		PopEfficiency:     10,
		FactoryOutput:     10,
		FactoryCost:       9,
		NumFactories:      17,
		FactoriesCostLess: true,
		MineOutput:        10,
		MineCost:          9,
		NumMines:          10,
		ResearchCost: ResearchCost{
			Energy:        ResearchCostExtra,
			Weapons:       ResearchCostExtra,
			Propulsion:    ResearchCostLess,
			Construction:  ResearchCostStandard,
			Electronics:   ResearchCostStandard,
			Biotechnology: ResearchCostLess,
		},
	}
}

func Insectoids() Race {
	return Race{
		Name:              "Insectoid",
		PluralName:        "Insectoids",
		PRT:               WM,
		LRTs:              Bitmask(ISB) | Bitmask(RS) | Bitmask(CE),
		HabLow:            Hab{-1, 0, 70},
		HabHigh:           Hab{-1, 100, 100},
		ImmuneGrav:        true,
		GrowthRate:        10,
		PopEfficiency:     10,
		FactoryOutput:     10,
		FactoryCost:       10,
		NumFactories:      10,
		FactoriesCostLess: false,
		MineOutput:        9,
		MineCost:          10,
		NumMines:          6,
		ResearchCost: ResearchCost{
			Energy:        ResearchCostLess,
			Weapons:       ResearchCostLess,
			Propulsion:    ResearchCostLess,
			Construction:  ResearchCostLess,
			Electronics:   ResearchCostStandard,
			Biotechnology: ResearchCostExtra,
		},
	}
}

func HEs() Race {
	r := NewRace()
	return *r.WithPRT(HE)
}
func SSs() Race {
	r := NewRace()
	return *r.WithPRT(SS)
}
func WMs() Race {
	r := NewRace()
	return *r.WithPRT(WM)
}
func CAs() Race {
	r := NewRace()
	return *r.WithPRT(CA)
}
func ISs() Race {
	r := NewRace()
	return *r.WithPRT(IS)
}
func SDs() Race {
	r := NewRace()
	return *r.WithPRT(SD)
}
func PPs() Race {
	r := NewRace()
	return *r.WithPRT(PP)
}
func ITs() Race {
	r := NewRace()
	return *r.WithPRT(IT)
}
func ARs() Race {
	r := NewRace()
	return *r.WithPRT(AR)
}
func JoaTs() Race {
	r := NewRace()
	return *r.WithPRT(JoaT)
}

func (r *Race) String() string {
	return fmt.Sprintf("Race %s (%d)", r.PluralName, r.ID)
}

func (r *Race) HasLRT(lrt LRT) bool {
	return r.LRTs&Bitmask(lrt) != 0
}

func (r *Race) HabCenter() Hab {
	return Hab{
		(r.HabHigh.Grav-r.HabLow.Grav)/2 + r.HabLow.Grav,
		(r.HabHigh.Temp-r.HabLow.Temp)/2 + r.HabLow.Temp,
		(r.HabHigh.Rad-r.HabLow.Rad)/2 + r.HabLow.Rad,
	}
}

func (r *Race) HabWidth() Hab {
	return Hab{
		(r.HabHigh.Grav - r.HabLow.Grav),
		(r.HabHigh.Temp - r.HabLow.Temp),
		(r.HabHigh.Rad - r.HabLow.Rad),
	}
}

// get this planet's habitabiliity from -45 to 100
func (r *Race) GetPlanetHabitability(hab Hab) int {
	planetValuePoints, redValue, ideality := 0, 0, 10000

	habValues := [3]int{hab.Grav, hab.Temp, hab.Rad}
	habCenters := [3]int{r.HabCenter().Grav, r.HabCenter().Temp, r.HabCenter().Rad}
	habLows := [3]int{r.HabLow.Grav, r.HabLow.Temp, r.HabLow.Rad}
	habHighs := [3]int{r.HabHigh.Grav, r.HabHigh.Temp, r.HabHigh.Rad}
	immune := [3]bool{r.ImmuneGrav, r.ImmuneTemp, r.ImmuneRad}

	var fromIdeal, tmp, habRadius, poorPlanetMod, habRed int

	for i := range habValues {
		habValue, habLower, habUpper, habCenter := habValues[i], habLows[i], habHighs[i], habCenters[i]

		if immune[i] {
			planetValuePoints += 10000
		} else {
			if habLower <= habValue && habUpper >= habValue {
				// green planet
				fromIdeal = int(math.Abs(float64(habValue-habCenter)) * 100)
				if habCenter > habValue {
					habRadius = habCenter - habLower
					fromIdeal /= habRadius
					tmp = habCenter - habValue
				} else {
					habRadius = habUpper - habCenter
					fromIdeal /= habRadius
					tmp = habValue - habCenter
				}
				poorPlanetMod = ((tmp) * 2) - habRadius
				fromIdeal = 100 - fromIdeal
				planetValuePoints += fromIdeal * fromIdeal
				if poorPlanetMod > 0 {
					ideality *= habRadius*2 - poorPlanetMod
					ideality /= habRadius * 2
				}
			} else {
				// red planet
				if habLower <= habValue {
					habRed = habValue - habUpper
				} else {
					habRed = habLower - habValue
				}

				if habRed > 15 {
					habRed = 15
				}

				redValue += habRed
			}
		}
	}

	if redValue != 0 {
		return -redValue
	}

	planetValuePoints = int(math.Sqrt(float64(planetValuePoints)/3.0) + 0.9)
	planetValuePoints = planetValuePoints * ideality / 10000

	return planetValuePoints
}

// compute the spec for this race
func computeRaceSpec(race *Race, rules *Rules) RaceSpec {
	prtSpec := rules.PRTSpecs[PRT(race.PRT)]
	spec := RaceSpec{
		HabCenter:           race.HabCenter(),
		StartingTechLevels:  prtSpec.StartingTechLevels,
		StartingPlanets:     prtSpec.StartingPlanets,
		TechCostOffset:      prtSpec.TechCostOffset,
		MaxPopulationOffset: prtSpec.MaxPopulationOffset,
		ScannerSpec: ScannerSpec{
			ScanRangeFactor: 1,
			// JoaT
			BuiltInScannerMultiplier: prtSpec.BuiltInScannerMultiplier,
		},
		StartingPopulationFactor: 1,
		ResearchFactor:           1,
		ShieldStrengthFactor:     1,
		ArmorStrengthFactor:      1,
		EngineReliableSpeed:      10,
		MiniaturizationSpec: MiniaturizationSpec{
			NewTechCostFactor:       1,
			MiniaturizationMax:      .75,
			MiniaturizationPerLevel: .04,
		},

		// PP
		MineralsPerSingleMineralPacket:   prtSpec.MineralsPerSingleMineralPacket,
		MineralsPerMixedMineralPacket:    prtSpec.MineralsPerMixedMineralPacket,
		PacketResourceCost:               prtSpec.PacketResourceCost,
		PacketMineralCostFactor:          prtSpec.PacketMineralCostFactor,
		PacketReceiverFactor:             prtSpec.PacketReceiverFactor,
		PacketDecayFactor:                prtSpec.PacketDecayFactor,
		PacketBuiltInScanner:             prtSpec.PacketBuiltInScanner,
		DetectPacketDestinationStarbases: prtSpec.DetectPacketDestinationStarbases,
		DetectAllPackets:                 prtSpec.DetectAllPackets,
		PacketTerraformChance:            prtSpec.PacketTerraformChance,
		PacketPermaformChance:            prtSpec.PacketPermaformChance,
		PacketPermaTerraformSizeUnit:     prtSpec.PacketPermaTerraformSizeUnit,

		// IT
		PacketOverSafeWarpPenalty: prtSpec.PacketOverSafeWarpPenalty,
		CanGateCargo:              prtSpec.CanGateCargo,
		CanDetectStargatePlanets:  prtSpec.CanDetectStargatePlanets,
		ShipsVanishInVoid:         prtSpec.ShipsVanishInVoid,

		// JoaT
		TechsCostExtraLevel: prtSpec.TechsCostExtraLevel,

		// IS
		FreighterGrowthFactor: prtSpec.FreighterGrowthFactor, // AR sets this negative
		InvasionDefendBonus:   prtSpec.InvasionDefendBonus,
		RepairFactor:          prtSpec.RepairFactor,
		StarbaseRepairFactor:  prtSpec.StarbaseRepairFactor,

		// HE
		GrowthFactor: prtSpec.GrowthFactor,

		// SS
		BuiltInCloakUnits: prtSpec.BuiltInCloakUnits,
		StealsResearch:    prtSpec.StealsResearch,
		FreeCargoCloaking: prtSpec.FreeCargoCloaking,

		// SD
		MineFieldsAreScanners:      prtSpec.MineFieldsAreScanners,
		MineFieldRateMoveFactor:    prtSpec.MineFieldRateMoveFactor,
		MineFieldSafeWarpBonus:     prtSpec.MineFieldSafeWarpBonus,
		MineFieldMinDecayFactor:    prtSpec.MineFieldMinDecayFactor,
		MineFieldBaseDecayRate:     prtSpec.MineFieldBaseDecayRate,
		MineFieldPlanetDecayRate:   prtSpec.MineFieldPlanetDecayRate,
		MineFieldMaxDecayRate:      prtSpec.MineFieldMaxDecayRate,
		CanDetonateMineFields:      prtSpec.CanDetonateMineFields,
		MineFieldDetonateDecayRate: prtSpec.MineFieldDetonateDecayRate,

		// WM
		DiscoverDesignOnScan: prtSpec.DiscoverDesignOnScan,
		InvasionAttackBonus:  prtSpec.InvasionAttackBonus,

		// AR
		CanRemoteMineOwnPlanets: prtSpec.CanRemoteMineOwnPlanets,
		MovementBonus:           prtSpec.MovementBonus,
		StarbaseCostFactor:      prtSpec.StarbaseCostFactor,
		InnateMining:            prtSpec.InnateMining,
		InnateResources:         prtSpec.InnateResources,
		InnateScanner:           prtSpec.InnateScanner,
		InnatePopulationFactor:  prtSpec.InnatePopulationFactor,
		CanBuildDefenses:        prtSpec.CanBuildDefenses,
		LivesOnStarbases:        prtSpec.LivesOnStarbases,

		// CA
		Instaforming:        prtSpec.Instaforming,
		PermaformChance:     prtSpec.PermaformChance,
		PermaformPopulation: prtSpec.PermaformPopulation,
	}

	if race.TechsStartHigh {
		// Jack of All Trades start at 4
		costsExtraLevel := prtSpec.TechsCostExtraLevel

		// if the race is configured with Extra cost for a tech, and the tech levels are less than the CostsExtraLevel, set it
		// to the costs extra level
		if race.ResearchCost.Energy == ResearchCostExtra && spec.StartingTechLevels.Energy < costsExtraLevel {
			spec.StartingTechLevels.Energy = costsExtraLevel
		}
		if race.ResearchCost.Weapons == ResearchCostExtra && spec.StartingTechLevels.Weapons < costsExtraLevel {
			spec.StartingTechLevels.Weapons = costsExtraLevel
		}
		if race.ResearchCost.Propulsion == ResearchCostExtra && spec.StartingTechLevels.Propulsion < costsExtraLevel {
			spec.StartingTechLevels.Propulsion = costsExtraLevel
		}
		if race.ResearchCost.Construction == ResearchCostExtra && spec.StartingTechLevels.Construction < costsExtraLevel {
			spec.StartingTechLevels.Construction = costsExtraLevel
		}
		if race.ResearchCost.Electronics == ResearchCostExtra && spec.StartingTechLevels.Electronics < costsExtraLevel {
			spec.StartingTechLevels.Electronics = costsExtraLevel
		}
		if race.ResearchCost.Biotechnology == ResearchCostExtra && spec.StartingTechLevels.Biotechnology < costsExtraLevel {
			spec.StartingTechLevels.Biotechnology = costsExtraLevel
		}
	}

	for _, lrt := range LRTs {
		if !race.HasLRT(lrt) {
			continue
		}
		lrtSpec := rules.LRTSpecs[lrt]

		spec.StartingTechLevels.Energy += lrtSpec.StartingTechLevels.Energy
		spec.StartingTechLevels.Weapons += lrtSpec.StartingTechLevels.Weapons
		spec.StartingTechLevels.Propulsion += lrtSpec.StartingTechLevels.Propulsion
		spec.StartingTechLevels.Construction += lrtSpec.StartingTechLevels.Construction
		spec.StartingTechLevels.Electronics += lrtSpec.StartingTechLevels.Electronics
		spec.StartingTechLevels.Biotechnology += lrtSpec.StartingTechLevels.Biotechnology

		spec.NewTechCostFactor += lrtSpec.NewTechCostFactorOffset
		spec.TerraformCostOffset = spec.TerraformCostOffset.Add(lrtSpec.TerraformCostOffset)
		spec.MiniaturizationMax += lrtSpec.MiniaturizationMax
		spec.MiniaturizationPerLevel += lrtSpec.MiniaturizationPerLevel
		spec.ScanRangeFactor += lrtSpec.ScanRangeFactorOffset
		spec.FuelEfficiencyOffset += lrtSpec.FuelEfficiencyOffset
		spec.MaxPopulationOffset += lrtSpec.MaxPopulationOffset
		spec.MineralAlchemyCostOffset += lrtSpec.MineralAlchemyCostOffset
		spec.ScrapMineralOffset += lrtSpec.ScrapMineralOffset
		spec.ScrapMineralOffsetStarbase += lrtSpec.ScrapMineralOffsetStarbase
		spec.ScrapResourcesOffset += lrtSpec.ScrapResourcesOffset
		spec.ScrapResourcesOffsetStarbase += lrtSpec.ScrapResourcesOffsetStarbase
		spec.StartingPopulationFactor += lrtSpec.StartingPopulationFactorDelta
		spec.StarbaseBuiltInCloakUnits += lrtSpec.StarbaseBuiltInCloakUnits
		spec.StarbaseCostFactor = math.Max(spec.StarbaseCostFactor, lrtSpec.StarbaseCostFactorOffset) // this isn't cumulative
		spec.ResearchFactor += lrtSpec.ResearchFactorOffset
		spec.ResearchSplashDamage += lrtSpec.ResearchSplashDamage
		spec.ShieldStrengthFactor += lrtSpec.ShieldStrengthFactorOffset
		spec.ShieldRegenerationRate += lrtSpec.ShieldRegenerationRateOffset
		spec.ArmorStrengthFactor += lrtSpec.ArmorStrengthFactorOffset
		spec.EngineFailureRate += lrtSpec.EngineFailureRateOffset
		spec.EngineReliableSpeed += lrtSpec.EngineReliableSpeed

		spec.StartingPlanets[0].StartingFleets = append(spec.StartingPlanets[0].StartingFleets, lrtSpec.StartingFleets...)

		if lrtSpec.NoAdvancedScanners {
			spec.NoAdvancedScanners = true
		}
	}

	factoryGermaniumOffset := 0
	if race.FactoriesCostLess {
		factoryGermaniumOffset = -1
	}

	spec.Costs = map[QueueItemType]Cost{

		QueueItemTypeMine:                   {Resources: race.MineCost},
		QueueItemTypeAutoMines:              {Resources: race.MineCost},
		QueueItemTypeFactory:                {Germanium: rules.FactoryCostGermanium + factoryGermaniumOffset, Resources: race.FactoryCost},
		QueueItemTypePlanetaryScanner:       rules.PlanetaryScannerCost,
		QueueItemTypeAutoFactories:          {Germanium: rules.FactoryCostGermanium + factoryGermaniumOffset, Resources: race.FactoryCost},
		QueueItemTypeMineralAlchemy:         {Resources: rules.MineralAlchemyCost + spec.MineralAlchemyCostOffset},
		QueueItemTypeAutoMineralAlchemy:     {Resources: rules.MineralAlchemyCost + spec.MineralAlchemyCostOffset},
		QueueItemTypeDefenses:               rules.DefenseCost,
		QueueItemTypeAutoDefenses:           rules.DefenseCost,
		QueueItemTypeTerraformEnvironment:   rules.TerraformCost.Add(spec.TerraformCostOffset),
		QueueItemTypeAutoMaxTerraform:       rules.TerraformCost.Add(spec.TerraformCostOffset),
		QueueItemTypeAutoMinTerraform:       rules.TerraformCost.Add(spec.TerraformCostOffset),
		QueueItemTypeIroniumMineralPacket:   {Resources: spec.PacketResourceCost, Ironium: int(float64(spec.MineralsPerSingleMineralPacket) * spec.PacketMineralCostFactor)},
		QueueItemTypeBoraniumMineralPacket:  {Resources: spec.PacketResourceCost, Boranium: int(float64(spec.MineralsPerSingleMineralPacket) * spec.PacketMineralCostFactor)},
		QueueItemTypeGermaniumMineralPacket: {Resources: spec.PacketResourceCost, Germanium: int(float64(spec.MineralsPerSingleMineralPacket) * spec.PacketMineralCostFactor)},
		QueueItemTypeMixedMineralPacket: {
			Resources: spec.PacketResourceCost,
			Ironium:   int(float64(spec.MineralsPerMixedMineralPacket) * spec.PacketMineralCostFactor),
			Boranium:  int(float64(spec.MineralsPerMixedMineralPacket) * spec.PacketMineralCostFactor),
			Germanium: int(float64(spec.MineralsPerMixedMineralPacket) * spec.PacketMineralCostFactor),
		},
		QueueItemTypeAutoMineralPacket: {
			Resources: spec.PacketResourceCost,
			Ironium:   int(float64(spec.MineralsPerMixedMineralPacket) * spec.PacketMineralCostFactor),
			Boranium:  int(float64(spec.MineralsPerMixedMineralPacket) * spec.PacketMineralCostFactor),
			Germanium: int(float64(spec.MineralsPerMixedMineralPacket) * spec.PacketMineralCostFactor),
		},
	}

	return spec
}

func (race *Race) ComputeRacePoints(startingPoints int) int {

	prtPointCost := map[PRT]int{
		HE:   -40,
		SS:   -95,
		WM:   -45,
		CA:   -10,
		IS:   100,
		SD:   150,
		PP:   -120,
		IT:   -180,
		AR:   -90,
		JoaT: 66,
	}

	lrtPointCost := map[LRT]int{
		IFE:  -235,
		TT:   -25,
		ARM:  -159,
		ISB:  -201,
		GR:   40,
		UR:   -240,
		MA:   -155,
		NRSE: 160,
		CE:   240,
		OBRM: 255,
		NAS:  325,
		LSP:  180,
		BET:  70,
		RS:   30,
	}

	points := startingPoints

	// get points for hab ranges
	habPoints := int(race.getHabRangePoints() / 2000)

	growthRateFactor := race.GrowthRate // use raw growth rate, otherwise
	// HEs pay for GR at 2x
	grRate := float64(growthRateFactor)

	// update the points based on growth rate
	if growthRateFactor <= 5 {
		points += (6 - growthRateFactor) * 4200
	} else if growthRateFactor <= 13 {
		switch growthRateFactor {
		case 6:
			points += 3600
		case 7:
			points += 2250
		case 8:
			points += 600
		case 9:
			points += 225
		}
		growthRateFactor = growthRateFactor*2 - 5
	} else if growthRateFactor < 20 {
		growthRateFactor = (growthRateFactor - 6) * 3
	} else {
		growthRateFactor = 45
	}

	points -= int(habPoints*growthRateFactor) / 24

	// give points for off center habs
	numImmunities := 0
	hc := race.HabCenter()
	habCenter := [3]int{
		hc.Grav,
		hc.Temp,
		hc.Rad,
	}
	for habType := 0; habType < 3; habType++ {
		if race.IsImmune(HabType(habType)) {
			numImmunities++
		} else {
			points += int(math.Abs(float64(habCenter[habType]-50)) * 4)
		}
	}

	// multiple immunities are penalized extra
	if numImmunities > 1 {
		points -= 150
	}

	// determine factory costs
	operationPoints := race.NumFactories
	productionPoints := race.FactoryOutput

	if operationPoints > 10 || productionPoints > 10 {
		operationPoints -= 9
		if operationPoints < 1 {
			operationPoints = 1
		}
		productionPoints -= 9
		if productionPoints < 1 {
			productionPoints = 1
		}

		// HE penalty, 2 for all PRTs execpt 3 for HE
		factoryProductionCost := 2
		if race.PRT == HE {
			factoryProductionCost = 3
		}

		productionPoints *= factoryProductionCost

		// additional penalty for two- and three-immune
		if numImmunities >= 2 {
			points -= int(float64(productionPoints*operationPoints)*grRate) / 2
		} else {
			points -= int(float64(productionPoints*operationPoints)*grRate) / 9
		}
	}

	// pop efficiency
	popEfficiency := race.PopEfficiency
	if popEfficiency > 25 {
		popEfficiency = 25
	}
	if popEfficiency <= 7 {
		points -= 2400
	} else if popEfficiency == 8 {
		points -= 1260
	} else if popEfficiency == 9 {
		points -= 600
	} else if popEfficiency > 10 {
		points += (popEfficiency - 10) * 120
	}

	// factory points (AR races have very simple points)
	if race.PRT == AR {
		points += 210
	} else {
		productionPoints = 10 - race.FactoryOutput
		costPoints := 10 - race.FactoryCost
		operationPoints = 10 - race.NumFactories
		tmpPoints := 0

		if productionPoints > 0 {
			tmpPoints = productionPoints * 100
		} else {
			tmpPoints = productionPoints * 121
		}

		if costPoints > 0 {
			tmpPoints += costPoints * costPoints * -60
		} else {
			tmpPoints += costPoints * -55
		}

		if operationPoints > 0 {
			tmpPoints += operationPoints * 40
		} else {
			tmpPoints += operationPoints * 35
		}

		// limit low factory points
		llfp := 700
		if tmpPoints > llfp {
			tmpPoints = (tmpPoints-llfp)/3 + llfp
		}

		if operationPoints <= -7 {
			if operationPoints < -11 {
				if operationPoints < -14 {
					tmpPoints -= 360
				} else {
					tmpPoints += (operationPoints + 7) * 45
				}
			} else {
				tmpPoints += (operationPoints + 6) * 30
			}
		}

		if productionPoints <= -3 {
			tmpPoints += (productionPoints + 2) * 60
		}

		points += tmpPoints

		if race.FactoriesCostLess {
			points -= 175
		}

		// mines
		productionPoints = 10 - race.MineOutput
		costPoints = 3 - race.MineCost
		operationPoints = 10 - race.NumMines
		tmpPoints = 0

		if productionPoints > 0 {
			tmpPoints = productionPoints * 100
		} else {
			tmpPoints = productionPoints * 169
		}

		if costPoints > 0 {
			tmpPoints -= 360
		} else {
			tmpPoints += costPoints*(-65) + 80
		}

		if operationPoints > 0 {
			tmpPoints += operationPoints * 40
		} else {
			tmpPoints += operationPoints * 35
		}

		points += tmpPoints
	}

	// prt and lrt point costs
	points += prtPointCost[race.PRT]

	// too many lrts
	badLRTs := 0
	goodLRTs := 0

	// figure out how many bad vs good lrts we have.
	for _, lrt := range LRTs {
		if race.HasLRT(lrt) {

			if lrtPointCost[lrt] >= 0 {
				badLRTs++
			} else {
				goodLRTs++
			}
			points += lrtPointCost[lrt]
		}
	}

	if goodLRTs+badLRTs > 4 {
		points -= (goodLRTs + badLRTs) * (goodLRTs + badLRTs - 4) * 10
	}
	if badLRTs-goodLRTs > 3 {
		points -= (badLRTs - goodLRTs - 3) * 60
	}
	if goodLRTs-badLRTs > 3 {
		points -= (goodLRTs - badLRTs - 3) * 40
	}

	// No Advanced scanners is penalized in some races
	if race.HasLRT(NAS) {
		if race.PRT == PP {
			points -= 280
		} else if race.PRT == SS {
			points -= 200
		} else if race.PRT == JoaT {
			points -= 40
		}
	}

	// Techs
	//
	// Figure out the total number of Extra's, offset by the number of Less's
	techcosts := 0
	researchCost := [6]ResearchCostLevel{
		race.ResearchCost.Energy,
		race.ResearchCost.Weapons,
		race.ResearchCost.Propulsion,
		race.ResearchCost.Construction,
		race.ResearchCost.Electronics,
		race.ResearchCost.Biotechnology,
	}
	for i := 0; i < 6; i++ {
		rc := researchCost[i]
		if rc == ResearchCostExtra {
			techcosts--
		} else if rc == ResearchCostLess {
			techcosts++
		}
	}

	// if we have more less's then extra's, penalize the race
	if techcosts > 0 {
		points -= (techcosts * techcosts) * 130
		if techcosts >= 6 {
			points += 1430 // already paid 4680 so true cost is 3250
		} else if techcosts == 5 {
			points += 520 // already paid 3250 so true cost is 2730
		}
	} else if techcosts < 0 {
		// if we have more extra's, give the race a bonus that increases as
		// we have more extra's
		scienceCost := []int{150, 330, 540, 780, 1050, 1380}
		points += scienceCost[(-techcosts)-1]
		if techcosts < -4 && (popEfficiency < 10) {
			points -= 190
		}
	}
	if race.TechsStartHigh {
		points -= 180
	}

	// ART races get penalized extra for have cheap energy because it gives them such a boost
	if race.PRT == AR && race.ResearchCost.Energy == ResearchCostLess {
		points -= 100
	}

	return points / 3

}

// get points for a race's hab range
func (race *Race) getHabRangePoints() int64 {
	totalTerraforming := race.LRTs&Bitmask(TT) > 0

	// setup the starting values for each hab type, and the widths
	// for those
	terraformOffset := [3]int{}
	testHabStart := [3]int{}
	testHabWidth := [3]int{}

	points := 0.0

	numIterationsGrav := 0
	numIterationsRad := 0
	numIterationsTemp := 0

	// set the number of iterations for each hab type.  If we're immune it's just
	// 1 because all the planets in that range will be the same.  Otherwise we loop
	// over the entire hab range in 11 equal divisions (i.e. for Humanoids grav would be 15, 22, 29, etc. all the way to 85)
	if race.ImmuneGrav {
		numIterationsGrav = 1
	} else {
		numIterationsGrav = 11
	}
	if race.ImmuneTemp {
		numIterationsTemp = 1
	} else {
		numIterationsTemp = 11
	}
	if race.ImmuneRad {
		numIterationsRad = 1
	} else {
		numIterationsRad = 11
	}

	// We go through 3 main iterations.  During each the habitability of the test planet
	// varies between the low and high of the hab range for each hab type.  So for a humanoid
	// it goes (15, 15, 15), (15, 15, 22), (15, 15, 29), etc.   Until it's (85, 85, 85)
	// During the various loops the TTCorrectionFactor changes to account for the race's ability
	// to terrform.
	for loopIndex := 0; loopIndex < 3; loopIndex++ {

		// each main loop gets a different TTCorrectionFactor
		ttCorrectionFactor := 0
		if loopIndex == 0 {
			ttCorrectionFactor = 0
		} else if loopIndex == 1 {
			if totalTerraforming {
				ttCorrectionFactor = 8
			} else {
				ttCorrectionFactor = 5
			}
		} else {
			if totalTerraforming {
				ttCorrectionFactor = 17
			} else {
				ttCorrectionFactor = 15
			}
		}

		// for each hab type, set up the starts and widths
		// for this outer loop
		for _, habType := range HabTypes {
			// if we're immune, just make the hab values some middle value
			if race.IsImmune(habType) {
				testHabStart[habType] = 50
				testHabWidth[habType] = 11
			} else {
				// start at the minimum hab range
				testHabStart[habType] = Clamp(race.HabLow.Get(habType)-ttCorrectionFactor, 0, 100)

				// get the high range for this hab type
				habHigh := Clamp(race.HabHigh.Get(habType)+ttCorrectionFactor, 0, 100)

				// figure out the width for this hab type's starting range
				testHabWidth[habType] = habHigh - testHabStart[habType]
			}
		}

		// 3 nested for loops, one for each hab type.  The number of iterations is 11 for non immune habs, or 1 for immune habs
		// this starts iterations for the first hab (gravity)
		gravitySum := 0.0
		testPlanetHab := Hab{}
		for iterationGrav := 0; iterationGrav < numIterationsGrav; iterationGrav++ {
			testPlanetHab.Grav, terraformOffset[0] = race.getPlanetHabForHabIndex(iterationGrav, 0, loopIndex, numIterationsGrav, testHabStart[0], testHabWidth[0], ttCorrectionFactor)

			// go through iterations for temperature
			temperatureSum := 0.0
			for iterationTemp := 0; iterationTemp < numIterationsTemp; iterationTemp++ {
				testPlanetHab.Temp, terraformOffset[1] = race.getPlanetHabForHabIndex(iterationTemp, 1, loopIndex, numIterationsTemp, testHabStart[1], testHabWidth[1], ttCorrectionFactor)

				// go through iterations for radiation
				var radiationSum int64 = 0
				for iterationRad := 0; iterationRad < numIterationsRad; iterationRad++ {
					testPlanetHab.Rad, terraformOffset[2] = race.getPlanetHabForHabIndex(iterationRad, 2, loopIndex, numIterationsRad, testHabStart[2], testHabWidth[2], ttCorrectionFactor)

					planetDesirability := int64(race.GetPlanetHabitability(testPlanetHab))

					terraformOffsetSum := terraformOffset[0] + terraformOffset[1] + terraformOffset[2]
					if terraformOffsetSum > ttCorrectionFactor {
						// bring the planet desirability down by the difference between the terraformOffsetSum and the TTCorrectionFactor
						planetDesirability -= int64(terraformOffsetSum - ttCorrectionFactor)
						// make sure the planet isn't negative in desirability
						if planetDesirability < 0 {
							planetDesirability = 0
						}
					}
					planetDesirability *= planetDesirability

					// modify the planetDesirability by some factor based on which main loop we're going through
					switch loopIndex {
					case 0:
						planetDesirability *= 7
					case 1:
						planetDesirability *= 5
					default:
						planetDesirability *= 6
					}

					radiationSum += planetDesirability

					// log.Debug($"Hab: {testPlanetHab}, desirability: {planetDesirability}, sums: ({gravitySum}, {temperatureSum}, {radiationSum}) ");
				}

				// The radiationSum is the sum of the planetDesirability for each iteration in numIterationsRad
				// if we're immune to radiation it'll be the same very loop, so *= by 11
				if !race.ImmuneRad {
					radiationSum = radiationSum * int64(testHabWidth[2]) / 100
				} else {
					radiationSum *= 11
				}

				temperatureSum += float64(radiationSum)
			}

			// The tempSum is the sum of the radSums
			// if we're immune to radiation it'll be the same very loop, so *= by 11
			if !race.ImmuneTemp {
				temperatureSum = (temperatureSum * float64(testHabWidth[1])) / 100
			} else {
				temperatureSum *= 11
			}

			gravitySum += temperatureSum
		}
		if !race.ImmuneGrav {
			gravitySum = (gravitySum * float64(testHabWidth[0])) / 100
		} else {
			gravitySum *= 11
		}

		points += gravitySum
	}

	return int64(points/10.0 + 0.5)
}

// used by race point calculator to get points for a single iteration of a hab loop
func (race *Race) getPlanetHabForHabIndex(iterIndex int, habType HabType, loopIndex int, numIterations int, testHabStart int, testHabWidth int, ttCorrectionFactor int) (planetHab int, terraformOffset int) {
	// on the first iteration just use the testHabStart we already defined
	// if we're on a subsequent loop move the hab value along the habitable range of this race
	if iterIndex == 0 || numIterations <= 1 {
		planetHab = testHabStart
	} else {
		planetHab = (testHabWidth*iterIndex)/(numIterations-1) + testHabStart
	}

	habCenter := race.HabCenter().Get(habType)

	// if we on a main loop other than the first one, do some
	// stuff with the terraforming correction factor
	if loopIndex != 0 && !race.IsImmune(HabType(habType)) {
		offset := habCenter - planetHab
		if int(math.Abs(float64(offset))) <= ttCorrectionFactor {
			offset = 0
		} else if offset < 0 {
			offset += ttCorrectionFactor
		} else {
			offset -= ttCorrectionFactor
		}

		// we return this terraformOffset value for later use
		// when we do the summing
		terraformOffset = offset
		planetHab = habCenter - offset
	}

	return planetHab, terraformOffset
}
