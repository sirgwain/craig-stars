package cs

import (
	"fmt"
	"math"
	"time"
)

type Race struct {
	ID                int64        `json:"id,omitempty"`
	CreatedAt         time.Time    `json:"createdAt,omitempty"`
	UpdatedAt         time.Time    `json:"updatedAt,omitempty"`
	UserID            int64        `json:"userId,omitempty"`
	Name              string       `json:"name,omitempty"`
	PluralName        string       `json:"pluralName,omitempty"`
	PRT               PRT          `json:"prt,omitempty"`
	LRTs              Bitmask      `json:"lrts,omitempty"`
	HabLow            Hab          `json:"habLow,omitempty"`
	HabHigh           Hab          `json:"habHigh,omitempty"`
	GrowthRate        int          `json:"growthRate,omitempty"`
	PopEfficiency     int          `json:"popEfficiency,omitempty"`
	FactoryOutput     int          `json:"factoryOutput,omitempty"`
	FactoryCost       int          `json:"factoryCost,omitempty"`
	NumFactories      int          `json:"numFactories,omitempty"`
	FactoriesCostLess bool         `json:"factoriesCostLess,omitempty"`
	ImmuneGrav        bool         `json:"immuneGrav,omitempty"`
	ImmuneTemp        bool         `json:"immuneTemp,omitempty"`
	ImmuneRad         bool         `json:"immuneRad,omitempty"`
	MineOutput        int          `json:"mineOutput,omitempty"`
	MineCost          int          `json:"mineCost,omitempty"`
	NumMines          int          `json:"numMines,omitempty"`
	ResearchCost      ResearchCost `json:"researchCost,omitempty"`
	TechsStartHigh    bool         `json:"techsStartHigh,omitempty"`
	Spec              RaceSpec     `json:"spec,omitempty"`
}

type ResearchCostLevel string

const (
	ResearchCostExtra    ResearchCostLevel = "Extra"
	ResearchCostStandard ResearchCostLevel = "Standard"
	ResearchCostLess     ResearchCostLevel = "Less"
)

type ResearchCost struct {
	Energy        ResearchCostLevel `json:"energy,omitempty"`
	Weapons       ResearchCostLevel `json:"weapons,omitempty"`
	Propulsion    ResearchCostLevel `json:"propulsion,omitempty"`
	Construction  ResearchCostLevel `json:"construction,omitempty"`
	Electronics   ResearchCostLevel `json:"electronics,omitempty"`
	Biotechnology ResearchCostLevel `json:"biotechnology,omitempty"`
}

type RaceSpec struct {
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
	BuiltInScannerMultiplier         int                    `json:"builtInScannerMultiplier,omitempty"`
	TechsCostExtraLevel              int                    `json:"techsCostExtraLevel,omitempty"`
	FreighterGrowthFactor            float64                `json:"freighterGrowthFactor,omitempty"`
	GrowthFactor                     float64                `json:"growthFactor,omitempty"`
	MaxPopulationOffset              float64                `json:"maxPopulationOffset,omitempty"`
	BuiltInCloakUnits                int                    `json:"builtInCloakUnits,omitempty"`
	StealsResearch                   TechLevel              `json:"stealsResearch,omitempty"`
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
	NewTechCostFactor                float64                `json:"newTechCostFactor,omitempty"`
	MiniaturizationMax               float64                `json:"miniaturizationMax,omitempty"`
	MiniaturizationPerLevel          float64                `json:"miniaturizationPerLevel,omitempty"`
	NoAdvancedScanners               bool                   `json:"noAdvancedScanners,omitempty"`
	ScanRangeFactor                  float64                `json:"scanRangeFactor,omitempty"`
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
	ShieldStrengthFactor             float64                `json:"shieldStrengthFactor,omitempty"`
	ShieldRegenerationRate           float64                `json:"shieldRegenerationRate,omitempty"`
	EngineFailureRate                float64                `json:"engineFailureRate,omitempty"`
	EngineReliableSpeed              int                    `json:"engineReliableSpeed,omitempty"`
}

type StealsResearch struct {
	Energy        float64 `json:"energy,omitempty"`
	Weapons       float64 `json:"weapons,omitempty"`
	Propulsion    float64 `json:"propulsion,omitempty"`
	Construction  float64 `json:"construction,omitempty"`
	Electronics   float64 `json:"electronics,omitempty"`
	Biotechnology float64 `json:"biotechnology,omitempty"`
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

type Bitmask uint32

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

func (r *Race) WithName(name string) *Race {
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

func (r *Race) WithSpec(rules *Rules) *Race {
	r.Spec = computeRaceSpec(r, rules)
	return r
}

func Humanoids() Race {
	return *NewRace()
}

func PPs() Race {
	r := NewRace()
	return *r.WithPRT(PP).WithName("Thrower").WithPluralName("Throwers")
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
		StartingTechLevels:       prtSpec.StartingTechLevels,
		StartingPlanets:          prtSpec.StartingPlanets,
		TechCostOffset:           prtSpec.TechCostOffset,
		MaxPopulationOffset:      prtSpec.MaxPopulationOffset,
		NewTechCostFactor:        1,
		MiniaturizationMax:       .75,
		MiniaturizationPerLevel:  .04,
		ScanRangeFactor:          1,
		StartingPopulationFactor: 1,
		ResearchFactor:           1,
		ShieldStrengthFactor:     1,
		EngineReliableSpeed:      10,

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
		BuiltInScannerMultiplier: prtSpec.BuiltInScannerMultiplier,
		TechsCostExtraLevel:      prtSpec.TechsCostExtraLevel,

		// IS
		FreighterGrowthFactor: prtSpec.FreighterGrowthFactor,
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

		spec.NewTechCostFactor += lrtSpec.NewTechCostFactor
		spec.TerraformCostOffset.Add(lrtSpec.TerraformCostOffset)
		spec.MiniaturizationMax += lrtSpec.MiniaturizationMax
		spec.MiniaturizationPerLevel += lrtSpec.MiniaturizationPerLevel
		spec.ScanRangeFactor += lrtSpec.ScanRangeFactor
		spec.FuelEfficiencyOffset += lrtSpec.FuelEfficiencyOffset
		spec.MaxPopulationOffset += lrtSpec.MaxPopulationOffset
		spec.MineralAlchemyCostOffset += lrtSpec.MineralAlchemyCostOffset
		spec.ScrapMineralOffset += lrtSpec.ScrapMineralOffset
		spec.ScrapMineralOffsetStarbase += lrtSpec.ScrapMineralOffsetStarbase
		spec.ScrapResourcesOffset += lrtSpec.ScrapResourcesOffset
		spec.ScrapResourcesOffsetStarbase += lrtSpec.ScrapResourcesOffsetStarbase
		spec.StartingPopulationFactor += lrtSpec.StartingPopulationFactor
		spec.StarbaseBuiltInCloakUnits += lrtSpec.StarbaseBuiltInCloakUnits
		spec.StarbaseCostFactor += lrtSpec.StarbaseCostFactor
		spec.ResearchFactor += lrtSpec.ResearchFactor
		spec.ResearchSplashDamage += lrtSpec.ResearchSplashDamage
		spec.ShieldStrengthFactor += lrtSpec.ShieldStrengthFactor
		spec.ShieldRegenerationRate += lrtSpec.ShieldRegenerationRate
		spec.EngineFailureRate += lrtSpec.EngineFailureRate
		spec.EngineReliableSpeed += lrtSpec.EngineReliableSpeed

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
