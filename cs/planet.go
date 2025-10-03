package cs

import (
	"fmt"
	"math"
)

// Planets are the only static and constant MapObject. They don't move and they can't be destroyed.
// Players also start the game knowing all planet names and locations.
// I suppose these should have been named Stars, since they represent a star system, ah well..
type Planet struct {
	GameDBObject
	MapObject
	PlanetOrders
	Hab                  Hab        `json:"hab"`
	BaseHab              Hab        `json:"baseHab"`
	TerraformedAmount    Hab        `json:"terraformedAmount,omitzero"`
	MineralConcentration Mineral    `json:"mineralConcentration"`
	MineYears            Mineral    `json:"mineYears,omitzero"`
	Cargo                Cargo      `json:"cargo,omitzero"`
	PartialPopulation    int        `json:"partialPopulation"` // population not in a multiple of 100
	Mines                int        `json:"mines"`
	Factories            int        `json:"factories"`
	Defenses             int        `json:"defenses"`
	Homeworld            bool       `json:"homeworld,omitempty"`
	Scanner              bool       `json:"scanner,omitempty"`
	Spec                 PlanetSpec `json:"spec,omitzero"`
	RandomArtifact       bool       `json:"-"`
	Starbase             *Fleet     `json:"-"`
	Dirty                bool       `json:"-"`
	bonusResources       int
}

type PlanetOrders struct {
	ContributesOnlyLeftoverToResearch bool                  `json:"contributesOnlyLeftoverToResearch,omitempty"`
	ProductionQueue                   []ProductionQueueItem `json:"productionQueue,omitempty"`
	RouteTargetType                   MapObjectType         `json:"routeTargetType,omitempty"`
	RouteTargetNum                    int                   `json:"routeTargetNum,omitempty"`
	RouteTargetPlayerNum              int                   `json:"routeTargetPlayerNum,omitempty"`
	PacketTargetNum                   int                   `json:"packetTargetNum,omitempty"`
	PacketSpeed                       int                   `json:"packetSpeed,omitempty"`
}

type PlanetSpec struct {
	PlanetStarbaseSpec
	CanTerraform                              bool    `json:"canTerraform,omitempty"`
	Defense                                   string  `json:"defense,omitempty"`
	DefenseCoverage                           float64 `json:"defenseCoverage,omitempty"`
	DefenseCoverageSmart                      float64 `json:"defenseCoverageSmart,omitempty"`
	GrowthAmount                              int     `json:"growthAmount,omitempty"`
	Habitability                              int     `json:"habitability,omitempty"`
	MaxDefenses                               int     `json:"maxDefenses,omitempty"`
	MaxFactories                              int     `json:"maxFactories,omitempty"`
	MaxMines                                  int     `json:"maxMines,omitempty"`
	MaxPopulation                             int     `json:"maxPopulation,omitempty"`
	MaxPossibleFactories                      int     `json:"maxPossibleFactories,omitempty"`
	MaxPossibleMines                          int     `json:"maxPossibleMines,omitempty"`
	MiningOutput                              Mineral `json:"miningOutput,omitzero"`
	PopulationDensity                         float64 `json:"populationDensity,omitempty"`
	ResourcesPerYear                          int     `json:"resourcesPerYear,omitempty"`
	ResourcesPerYearAvailable                 int     `json:"resourcesPerYearAvailable,omitempty"`
	ResourcesPerYearResearch                  int     `json:"resourcesPerYearResearch,omitempty"`
	ResourcesPerYearResearchEstimatedLeftover int     `json:"resourcesPerYearResearchEstimatedLeftover,omitempty"`
	Scanner                                   string  `json:"scanner,omitempty"`
	ScanRange                                 int     `json:"scanRange,omitempty"`
	ScanRangePen                              int     `json:"scanRangePen,omitempty"`
	TerraformAmount                           Hab     `json:"terraformAmount,omitzero"`
	MinTerraformAmount                        Hab     `json:"minTerraformAmount,omitzero"`
	TerraformedHabitability                   int     `json:"terraformedHabitability,omitempty"`
}

type PlanetStarbaseSpec struct {
	HasMassDriver      bool   `json:"hasMassDriver,omitempty"`
	HasStarbase        bool   `json:"hasStarbase,omitempty"`
	HasStargate        bool   `json:"hasStargate,omitempty"`
	StarbaseDesignName string `json:"starbaseDesignName,omitempty"`
	StarbaseDesignNum  int    `json:"starbaseDesignNum,omitempty"`
	DockCapacity       int    `json:"dockCapacity,omitempty"`
	BasePacketSpeed    int    `json:"basePacketSpeed,omitempty"`
	SafePacketSpeed    int    `json:"safePacketSpeed,omitempty"`
	SafeHullMass       int    `json:"safeHullMass,omitempty"`
	SafeRange          int    `json:"safeRange,omitempty"`
	MaxRange           int    `json:"maxRange,omitempty"`
	MaxHullMass        int    `json:"maxHullMass,omitempty"`
	Stargate           string `json:"stargate,omitempty"`
	MassDriver         string `json:"massDriver,omitempty"`
}

func (item *ProductionQueueItem) String() string {
	return fmt.Sprintf("ProductionQueueItem %d %s (%d)", item.Quantity, item.Type, item.DesignNum)
}

func NewPlanet() *Planet {
	return &Planet{MapObject: MapObject{Type: MapObjectTypePlanet, PlayerNum: Unowned}, Dirty: true}
}

func (p *Planet) MarkDirty() {
	p.Dirty = true
}

func (p *Planet) WithOrders(orders PlanetOrders) *Planet {
	p.PlanetOrders = orders
	return p
}

func (p *Planet) withPosition(position Vector) *Planet {
	p.Position = position
	return p
}

func (p *Planet) WithNum(num int) *Planet {
	p.Num = num
	return p
}

// Set a planet's colonists to the specified number of colonists and return the resulting struct.
// Multiples of 100 go into its Cargo struct, while leftovers are assigned to PartialPopulation.
func (p *Planet) WithPopulation(pop int) *Planet {
	p.Cargo.Colonists = pop / 100
	p.PartialPopulation = pop % 100
	return p
}

// Set a planet's Hab and BaseHab and return the resulting struct.
func (p *Planet) WithHab(hab Hab) *Planet {
	p.Hab = hab
	p.BaseHab = hab
	return p
}

func (p *Planet) WithCargo(cargo Cargo) *Planet {
	p.Cargo = cargo
	return p
}

func (p *Planet) WithPlayerNum(playerNum int) *Planet {
	p.PlayerNum = playerNum
	return p
}

func (p *Planet) WithFactories(factories int) *Planet {
	p.Factories = factories
	return p
}

func (p *Planet) WithMines(mines int) *Planet {
	p.Mines = mines
	return p
}

func (p *Planet) WithDefenses(defenses int) *Planet {
	p.Defenses = defenses
	return p
}

func (p *Planet) WithMineralConcentration(mineralConcentration Mineral) *Planet {
	p.MineralConcentration = mineralConcentration
	return p
}

func (p *Planet) WithMineYears(mineYears Mineral) *Planet {
	p.MineYears = mineYears
	return p
}

func (p *Planet) WithScanner(scanner bool) *Planet {
	p.Scanner = scanner
	return p
}

func (p *Planet) WithHomeworld(homeworld bool) *Planet {
	p.Homeworld = homeworld
	return p
}

func (p *Planet) WithContributesOnlyLeftoverToResearch(contributes bool) *Planet {
	p.ContributesOnlyLeftoverToResearch = contributes
	return p
}

func (p *Planet) String() string {
	return fmt.Sprintf("Planet %v", p.MapObject)
}

// return planetary population rounded to the nearest multiple of 100
func (p *Planet) GetPopulation() (wholePop int) {
	return p.Cargo.Colonists * 100
}

func (p *Planet) exactPopulation() (exactPop int) {
	return p.Cargo.Colonists*100 + p.PartialPopulation
}

// set pop to specified value.
// TODO: Remove this - it risks tampering with planet partial pop
func (p *Planet) setPopulation(pop int) {
	p.Cargo.Colonists = pop / 100
	p.PartialPopulation = pop % 100
}

// add specified amount of population to this Planet, incrementing or
// decrementing its Cargo and PartialPopulation values as appropriate.
// Planet pop/partial pop is floored to a minimum of 0.
func (p *Planet) addPopulation(pop int) {
	p.PartialPopulation += pop
	p.Cargo.Colonists += p.PartialPopulation / 100
	p.PartialPopulation = p.PartialPopulation % 100
	if p.PartialPopulation < 0 {
		// rollover partial pop due to trunc division annoyance
		p.Cargo.Colonists -= 1
		p.PartialPopulation += 100
	}
	if p.Cargo.Colonists < 0 {
		p.Cargo.Colonists = 0
	}
}

// Return the amount of population considered productive for resource production,
// taking into account overcrowding penalties.
func ProductivePopulation(pop, maxPop int, overcrowdPenalty, overcrowdResourceMax float64) int {
	popOverCap := float64(pop) + max(0, float64(pop-maxPop)*overcrowdPenalty)
	return roundTo100(min(
		float64(maxPop)*(1+overcrowdResourceMax), popOverCap), math.Floor)
}

// return true if this planet is able to build a ship with the given mass;
// cost of ship not considered
func (p *Planet) CanBuild(mass int) bool {
	return p.Spec.HasStarbase && (p.Starbase.Spec.SpaceDock == UnlimitedSpaceDock || p.Starbase.Spec.SpaceDock >= mass)
}

// populate a starbase design for a planet
func (p *Planet) PopulateStarbaseDesign(player *Player) error {
	if p.Starbase == nil {
		return nil
	}

	if len(p.Starbase.Tokens) != 1 {
		return fmt.Errorf("planet %s' starbase has no tokens", p.Name)
	}

	designNum := p.Starbase.Tokens[0].DesignNum
	design := player.GetDesign(designNum)
	if design == nil {
		return fmt.Errorf("player %v does not have design %d", player, designNum)
	}

	p.Starbase.Tokens[0].design = player.GetDesign(designNum)
	return nil
}

// add designs to each production queue item with designs
func (p *Planet) PopulateProductionQueueDesigns(player *Player) error {
	for i := range p.ProductionQueue {
		item := &p.ProductionQueue[i]
		if item.Type != QueueItemTypeStarbase && item.Type != QueueItemTypeShipToken {
			// not ship or starbase; move on
			continue
		}
		design := player.GetDesign(item.DesignNum)
		if design == nil {
			return fmt.Errorf("player %v does not have design %d", player, item.DesignNum)
		}
		item.design = design
	}
	return nil
}

// populate the costs of each item in the planet production queue
func (p *Planet) PopulateProductionQueueEstimates(rules *Rules, player *Player) error {
	// populate completion estimates
	completionEstimator := NewCompletionEstimator()
	var err error
	p.ProductionQueue, p.Spec.ResourcesPerYearResearchEstimatedLeftover, err =
		completionEstimator.GetProductionWithEstimates(rules, player, *p)
	return err
}

// empty this planet of pop & owner, typically used when transferring or losing ownership
func (p *Planet) emptyPlanet() {
	p.PlayerNum = Unowned
	p.Starbase = nil
	// defenses & scanner disappear, other structures stay though
	p.Scanner = false
	p.Defenses = 0
	// clear any production or other orders from the previous owner
	p.PlanetOrders = PlanetOrders{}
	p.setPopulation(0)
	p.Spec = PlanetSpec{}
	p.Hab = p.BaseHab.Add(p.TerraformedAmount) // reset any instaforming
}

// randomize a planet with new hab range, minerals, etc;
// Used in universe generation as well as for genesis device resets
func (p *Planet) randomize(rules *Rules, accBBS bool) {
	// From @SuicideJunkie's tests and @edmundmk's previous research,
	// grav and temp are weighted slightly towards the center while
	// rad is completely random (though all 3 are clamped between 1 and 99).

	// First, we handle the first block of the hab randomness disregarding dropoff
	p.Hab = Hab{
		Grav: rules.MinHab + rules.random.Intn(rules.MaxHab-rules.MinHab-rules.HabDropoffRange.Grav+1), // 1+randint(99-1-9+1)
		Temp: rules.MinHab + rules.random.Intn(rules.MaxHab-rules.MinHab-rules.HabDropoffRange.Temp+1),
		Rad:  rules.MinHab + rules.random.Intn(rules.MaxHab-rules.MinHab-rules.HabDropoffRange.Rad+1),
	}

	// add random amounts to simulate dropoff at the extremes ranges
	var randomG, randomT, randomR int
	if rules.HabDropoffRange.Grav > 0 {
		randomG = rules.random.Intn(rules.HabDropoffRange.Grav + 1)
	}
	if rules.HabDropoffRange.Temp > 0 {
		randomT = rules.random.Intn(rules.HabDropoffRange.Temp + 1)
	}
	if rules.HabDropoffRange.Rad > 0 {
		randomR = rules.random.Intn(rules.HabDropoffRange.Rad + 1)
	}

	p.Hab = p.Hab.Add(Hab{randomG, randomT, randomR})

	// reset the other stuff
	p.BaseHab = p.Hab
	p.TerraformedAmount = Hab{}
	p.MineralConcentration = randomizeMinerals(rules, p.Hab.Rad, accBBS)
	p.MineYears = Mineral{}

}

// Randomize a planet's mineral concentration within bounds set in Rules
func randomizeMinerals(rules *Rules, rad int, accBBS bool) Mineral {

	// These two variables are the shape of the normal distribution
	// based on comparing it with Stars! output
	mean := 80.0
	variance := 20.0

	// min and max of the minerals to be returned,
	// clamping the results
	mMin := rules.MinStartingMineralConcentration
	mMax := rules.MaxStartingMineralConcentration

	// create a normalized mineral concentration
	minConc := Mineral{
		Ironium:   1 + NormalSample(rules.random, mean, variance, mMax),
		Boranium:  1 + NormalSample(rules.random, mean, variance, mMax),
		Germanium: 1 + NormalSample(rules.random, mean, variance, mMax),
	}

	// add a small amount of minerals for accBBS
	if accBBS {
		for _, minType := range MineralTypes {
			if concAmount := minConc.GetAmount(minType); concAmount < 40 {
				minConc.Set(minType, concAmount+5)
			}
		}
	}

	// limit at least one mineral
	// TODO: make this limiting configurable?
	// it follows the original algorithm, but maybe there is a way to explain
	// what it's doing and make it easy to update for a mod
	// it picks a random number from 0 to 27, if under 18, limit a mineral
	// if the number is 9 to 18, only limit 1
	// if the number is 0, limit up to 4 times (limiter starts at 1, doubles each loop 1, 2, 4, 8)
	limiter := rules.random.Intn(27)
	if limiter < 18 {
		if limiter >= 9 {
			mineralType := MineralTypes[rules.random.Intn(len(MineralTypes))]
			value := 1 + rules.random.Intn(rules.LimitMineralConcentration)
			minConc.Set(mineralType, value)
		} else {
			limiter++
			for limiter < 16 {
				mineralType := MineralTypes[rules.random.Intn(len(MineralTypes))]
				value := 1 + rules.random.Intn(rules.LimitMineralConcentration)
				minConc.Set(mineralType, value)

				limiter *= 2
			}
		}
	}

	// we have high rad, add some bonus minerals
	if rad >= rules.HighRadMineralConcentrationBonusThreshold {
		minConc = Mineral{
			Ironium:   minConc.Ironium + rules.random.Intn(99-min(minConc.Ironium, 98))/2,
			Boranium:  minConc.Boranium + rules.random.Intn(99-min(minConc.Boranium, 98))/2,
			Germanium: minConc.Germanium + rules.random.Intn(99-min(minConc.Germanium, 98))/2,
		}
	}

	minConc.Ironium = Clamp(minConc.Ironium, mMin, mMax)
	minConc.Boranium = Clamp(minConc.Boranium, mMin, mMax)
	minConc.Germanium = Clamp(minConc.Germanium, mMin, mMax)

	return minConc
}

// Initialize a planet to be a homeworld for a player with ideal hab, starting mineral concentration, etc
func (p *Planet) initStartingWorld(player *Player, rules *Rules, startingPlanet StartingPlanet, concentration Mineral, surface Mineral) {
	p.Homeworld = startingPlanet.Homeworld

	p.RandomArtifact = false // no random artifacts on the homeworld
	p.PlayerNum = player.Num

	habWidth := player.Race.HabWidth()
	habCenter := player.Race.HabCenter()

	if !player.Race.ImmuneGrav {
		p.Hab.Grav = habCenter.Grav + int(float64(habWidth.Grav-rules.random.Intn(habWidth.Grav-1))/2*startingPlanet.HabPenaltyFactor)
	}
	if !player.Race.ImmuneTemp {
		p.Hab.Temp = habCenter.Temp + int(float64(habWidth.Temp-rules.random.Intn(habWidth.Temp-1))/2*startingPlanet.HabPenaltyFactor)
	}
	if !player.Race.ImmuneRad {
		p.Hab.Rad = habCenter.Rad + int(float64(habWidth.Rad-rules.random.Intn(habWidth.Rad-1))/2*startingPlanet.HabPenaltyFactor)
	}
	// BaseHab is the same as Hab
	p.BaseHab = p.Hab

	raceSpec := player.Race.Spec

	p.MineralConcentration = concentration
	p.Cargo = NewCargoFromMineral(surface,
		int(float64(startingPlanet.Population/100)*raceSpec.StartingPopulationFactor))

	// empty queue, no terraform
	p.ProductionQueue = []ProductionQueueItem{}
	p.TerraformedAmount = Hab{}

	if raceSpec.InnateMining {
		p.Mines = innateMines(raceSpec.InnateMinesFactor, p.GetPopulation())
		p.Factories = 0
	} else {
		p.Mines = startingPlanet.Mines
		p.Factories = startingPlanet.Factories
	}

	if raceSpec.CanBuildDefenses {
		p.Defenses = startingPlanet.Defenses
	} else {
		p.Defenses = 0
	}

	p.ContributesOnlyLeftoverToResearch = false
	p.Scanner = true

	if len(player.ProductionPlans) > 0 {
		// apply default production plan
		plan := player.ProductionPlans[0]
		plan.Apply(p)
	}

}

// set this planet's starbase on this planet
func (p *Planet) setStarbase(starbase *Fleet) {
	p.Starbase = starbase
	p.PacketSpeed = starbase.Spec.SafePacketSpeed
}

// return the amount of population this planet will have next year,
// truncated to the nearest multiple of 100.
func (p *Planet) PopNextYear() int {
	pop := p.exactPopulation() + p.Spec.GrowthAmount
	return roundTo100(pop, math.Floor)
}

// Get the number of innate mines a player would have with the given amount of population
func innateMines(innateMinesFactor float64, population int) int {
	// Verified - floored to nearest integer
	return int(math.Sqrt(float64(population)) * innateMinesFactor)
}

// Get the innate scanning distance a player would have wih the given amount of population
func innateScanner(innateScannerFactor float64, population int) int {
	// Verified - floored to nearest integer
	return int(math.Sqrt(float64(population) * innateScannerFactor))
}

// Find the shortest distance from one planet to a list of other planets
func (p *Planet) shortestDistanceToPlanets(otherPlanets []*Planet) float64 {
	minDistanceSquared := math.MaxInt
	for _, planet := range otherPlanets {
		distSquared := p.Position.DistanceSquaredTo(planet.Position)
		minDistanceSquared = min(minDistanceSquared, distSquared)
	}
	return math.Sqrt(float64(minDistanceSquared))
}

// getMineralOutput returns the mineral output of this Planet
// were it to be mined with the given numMines and mineOutput.
//
// Takes into account HW conc flooring as appropriate.
func (p *Planet) getMineralOutput(rules *Rules, numMines int, mineOutput int) (output Mineral) {
	for _, minType := range MineralTypes {
		conc := p.MineralConcentration.GetAmount(minType)
		if p.Homeworld && p.Owned() {
			// only apply HW conc floor if planet is owned.
			// We don't need to worry about remote miners since only unowned
			// or self-owned planets (for ARs) can be mined remotely.
			conc = max(conc, rules.MinHomeworldMineralConcentration)
		}

		// TODO: Check how Stars! does fractional mineral concs
		// and make this return a fractional output if needed.
		output.Set(minType, conc*numMines*mineOutput/1000)
	}
	return output
}

// Get how much a player will grow on a planet, given the max population the player can have on the planet.
//
// Returns exact value to nearest colonist.
func (p *Planet) GetGrowthAmount(player *Player, maxPopulation int, populationOvercrowdDieoffRate, populationOvercrowdDieoffRateMax float64) int {
	race := &player.Race
	pop := p.GetPopulation()
	habValue := race.GetPlanetHabitability(p.Hab)

	// First, we deal with cases where the pop doesn't grow
	if habValue < 0 {
		// Red worlds kill off (habValue / 10)% colonists every year
		// (habValue of -4% kills off 0.4%/yr)
		return int(math.Round(float64(pop*habValue) / 1000))
	}

	capacity := float64(pop) / float64(maxPopulation)
	if capacity > 1 {
		// Overpopulation kills 0.04% population per 1% over cap.
		// A 200% capacity planet is 100% over cap and thus loses
		// (0.04 * 100 = 4%) population each year.
		// This maxes out at 400% capacity (300% extra) at 12% deaths/yr.
		dieoffPercent := Clamp((capacity-1)*populationOvercrowdDieoffRate, 0, populationOvercrowdDieoffRateMax)
		return int(math.Round(float64(pop) * -dieoffPercent))
	}

	// perform normal pop growth calcs, applying penalty for partially crowded planets
	popGrowth := math.Round(float64(pop*race.GrowthRate*habValue) * race.Spec.GrowthFactor / 10000) // divide by 10000 as growthRate and habValue are both percents

	if capacity > 0.25 {
		crowdingFactor := math.Pow(1-capacity, 2) * 16 / 9
		popGrowth *= crowdingFactor
	}

	return int(popGrowth)
}

// compute a planet's PlanetSpec.
func ComputePlanetSpec(rules *Rules, player *Player, planet *Planet) PlanetSpec {
	spec := PlanetSpec{}
	race := &player.Race

	// the player spec is computed at turn generation time, but when updating planets we won't have it
	// so compute it on demand
	if player.Spec.Terraform == nil {
		player.Spec = ComputePlayerSpec(player, rules)
	}

	defense := player.Spec.Defense
	scanner := player.Spec.PlanetaryScanner

	// hab/pop
	spec.Habitability = race.GetPlanetHabitability(planet.Hab)
	spec.MaxPopulation = planet.getMaxPopulation(rules, player, spec.Habitability)
	if spec.MaxPopulation > 0 {
		spec.PopulationDensity = float64(planet.GetPopulation()) / float64(spec.MaxPopulation)
	}
	spec.GrowthAmount = planet.GetGrowthAmount(player, spec.MaxPopulation, rules.PopulationOvercrowdDieoffRate, rules.PopulationOvercrowdDieoffRateMax)

	// terraforming
	terraformer := NewTerraformer()
	spec.TerraformAmount = terraformer.GetTerraformAmount(planet.Hab, planet.BaseHab, player, player)
	spec.MinTerraformAmount = terraformer.GetMinTerraformAmount(planet.Hab, planet.BaseHab, player, player)
	spec.CanTerraform = spec.TerraformAmount.absSum() > 0
	spec.TerraformedHabitability = race.GetPlanetHabitability(planet.Hab.Add(spec.TerraformAmount))

	// population will generate resources up to 3x max pop, but they can only
	// operate structures up to max pop
	productivePop := ProductivePopulation(planet.GetPopulation(), spec.MaxPopulation, rules.PopulationOvercrowdResourcePenalty, rules.PopulationOvercrowdResourceMax)
	installationPop := min(planet.GetPopulation(), spec.MaxPopulation)

	if !race.Spec.InnateMining {
		spec.MaxMines = getMaxInstallations(player.Race.NumMines, installationPop)
		spec.MaxPossibleMines = spec.MaxPopulation * race.NumMines / 10000
	} else {
		spec.MaxMines = planet.Mines
	}

	// Compute resources per year and mining output
	spec.ComputeResourcesPerYear(player, planet.Factories, productivePop, installationPop)
	spec.MiningOutput = planet.getMineralOutput(rules, min(spec.MaxMines, planet.Mines), race.MineOutput)
	spec.computeResourcesPerYearAvailable(player, planet)

	if race.Spec.CanBuildDefenses {
		spec.MaxDefenses = 100
		spec.Defense = defense.Name
		spec.computeDefenseCoverage(rules, defense.DefenseCoverage, planet.Defenses)
	}

	if race.Spec.InnateScanner {
		// compute AR organic scan range
		// TODO: confirm rounding behavior with NAS
		spec.Scanner = "Organic"
		spec.ScanRange = int(float64(innateScanner(player.Race.Spec.InnateScannerFactor, productivePop)) * player.Race.Spec.ScanRangeFactor)
		if !player.Race.Spec.NoAdvancedScanners && planet.Starbase != nil {
			spec.ScanRangePen = int(float64(spec.ScanRange) * planet.Starbase.Spec.InnateScanRangePenFactor)
		}
	} else if planet.Scanner {
		// normal scanner ranges
		spec.Scanner = scanner.Name
		spec.ScanRange = int(float64(scanner.ScanRange) * player.Race.Spec.ScanRangeFactor)
		spec.ScanRangePen = scanner.ScanRangePen
	}

	spec.PlanetStarbaseSpec = computePlanetStarbaseSpec(planet)

	return spec
}

func computePlanetStarbaseSpec(planet *Planet) PlanetStarbaseSpec {
	spec := PlanetStarbaseSpec{}

	starbase := planet.Starbase
	if starbase == nil {
		return spec
	}

	spec.HasStarbase = true
	spec.StarbaseDesignNum = planet.Starbase.Tokens[0].DesignNum
	spec.StarbaseDesignName = planet.Starbase.Tokens[0].design.Name
	if starbase.Spec.HasStargate {
		spec.HasStargate = true
		spec.Stargate = starbase.Spec.Stargate
		spec.SafeHullMass = starbase.Spec.SafeHullMass
		spec.SafeRange = starbase.Spec.SafeRange
		spec.MaxHullMass = starbase.Spec.MaxHullMass
		spec.MaxRange = starbase.Spec.MaxRange
	}
	if starbase.Spec.HasMassDriver {
		spec.HasMassDriver = true
		spec.MassDriver = starbase.Spec.MassDriver
		spec.BasePacketSpeed = starbase.Spec.BasePacketSpeed
		spec.SafePacketSpeed = starbase.Spec.SafePacketSpeed
	}
	spec.DockCapacity = starbase.Spec.SpaceDock

	return spec
}

// Compute and update this planet's regular and smart defense coverage values
// TODO: Test this
func (spec *PlanetSpec) computeDefenseCoverage(rules *Rules, coverage float64, numDefenses int) {
	// coverage is a percentage, so divide by 100
	blocked := math.Pow(1-coverage/100, float64(Clamp(numDefenses, 0, spec.MaxDefenses)))
	spec.DefenseCoverage = 1 - blocked
	blockedSmart := math.Pow(1-(coverage/100)*rules.SmartDefenseCoverageFactor, float64(Clamp(numDefenses, 0, spec.MaxDefenses)))
	spec.DefenseCoverageSmart = 1 - blockedSmart
}

// Compute the amount of resources this planet will produce per year, as well as its
// MaxFactories and MaxPossibleFactories fields.
func (spec *PlanetSpec) ComputeResourcesPerYear(player *Player, numFacts, productivePop, installationPop int) {
	if player.Race.Spec.InnateResources {
		// Compute resources for AR
		habMulti := float64(max(spec.Habitability, player.Race.Spec.MinHabFloor)) / 100
		// Confirmed: AR resources round up in base game
		spec.ResourcesPerYear = int(math.Ceil(habMulti *
			math.Sqrt(float64(productivePop*player.TechLevels.Energy)/float64(player.Race.PopEfficiency))))
	} else {
		// compute resources from population & factories
		resourcesFromPop := productivePop / (player.Race.PopEfficiency * 100)

		spec.MaxFactories = getMaxInstallations(player.Race.NumFactories, installationPop)
		spec.MaxPossibleFactories = spec.MaxPopulation * player.Race.NumFactories / 10000 // factory count rounds down
		resourcesFromFactories := int(math.Ceil(float64(min(numFacts, spec.MaxFactories)*player.Race.FactoryOutput) / 10))

		// Add them together
		spec.ResourcesPerYear = resourcesFromPop + resourcesFromFactories
	}
}

// Update a planet spec's ResourcesPerYearAvailable and ResourcesPerYearResearch stats.
//
// This is called by the main ComputePlanetSpec function as well as anytime a player
// changes research contribution amounts
func (spec *PlanetSpec) computeResourcesPerYearAvailable(player *Player, planet *Planet) {
	if planet.ContributesOnlyLeftoverToResearch {
		spec.ResourcesPerYearAvailable = spec.ResourcesPerYear
		spec.ResourcesPerYearResearch = 0
	} else {
		spec.ResourcesPerYearResearch = spec.ResourcesPerYear * player.ResearchAmount / 100
		spec.ResourcesPerYearAvailable = spec.ResourcesPerYear - spec.ResourcesPerYearResearch
	}
}

// get the max population for this planet for a player with the given hab value
func (p *Planet) getMaxPopulation(rules *Rules, player *Player, habitability int) int {
	maxPopulationFactor := 1 + player.Race.Spec.MaxPopulationOffset
	if player.Race.Spec.LivesOnStarbases && p.PlayerNum == player.Num {
		// AR races' max pop are independent of habitability
		// TODO: How does this round again?
		return int(roundTo100(float64(p.Starbase.Spec.MaxPopulation)*maxPopulationFactor, math.Floor))
	}

	// Habitability is floored at 5% when determining max population
	// (or 25% for AR races)
	// We divide by 100 since we store the value as an int rather than a percentage
	habitability = max(habitability, player.Race.Spec.MinHabFloor)
	return roundTo100(float64(rules.MaxPopulation*habitability)*maxPopulationFactor/100.0, math.Floor)
}

// return the maximum number count operable by the given population
func getMaxInstallations(installationsPer10K, population int) int {
	return population * installationsPer10K / 10000
}

func (planet *Planet) MaxBuildable(player *Player, itemType QueueItemType) int {
	switch itemType {
	case QueueItemTypeAutoMines:
		// for autobuild purposes, the maxFactories is next year's pop
		// don't want to floor to 100
		futurePop := min(planet.PopNextYear(), planet.Spec.MaxPopulation)
		maxMines := getMaxInstallations(player.Race.NumMines, futurePop)
		return max(0, maxMines-planet.Mines)
	case QueueItemTypeAutoFactories:
		// for autobuild purposes, the maxFactories is next year's pop
		futurePop := min(planet.PopNextYear(), planet.Spec.MaxPopulation)
		maxFactories := getMaxInstallations(player.Race.NumFactories, futurePop)
		return max(0, maxFactories-planet.Factories)
	case QueueItemTypeMine:
		return max(0, planet.Spec.MaxPossibleMines-planet.Mines)
	case QueueItemTypeFactory:
		return max(0, planet.Spec.MaxPossibleFactories-planet.Factories)
	case QueueItemTypeAutoDefenses, QueueItemTypeDefenses:
		return max(0, planet.Spec.MaxDefenses-planet.Defenses)
	case QueueItemTypeTerraformEnvironment, QueueItemTypeAutoMaxTerraform:
		return planet.Spec.TerraformAmount.absSum()
	case QueueItemTypeAutoMinTerraform:
		return planet.Spec.MinTerraformAmount.absSum()
	case QueueItemTypeStarbase, QueueItemTypeGenesisDevice:
		return 1
	case QueueItemTypePlanetaryScanner:
		if planet.Scanner {
			return 0
		}
		return 1
	// TODO: Enable once auto alchemy gets fixed
	/* case QueueItemTypeAutoMineralAlchemy:
	return 1 */
	default:
		return Infinite
	}
}

// mine this planet using the given miningOutput and numMines
// TODO: add fractional mining support
func (planet *Planet) mine(rules *Rules, miningOutput Mineral, numMines int) {
	planet.Cargo = planet.Cargo.AddMineral(miningOutput)
	planet.MineYears = planet.MineYears.AddToAll(numMines)
	planet.reduceMineralConcentration(rules)
}

// grow pop on this planet (or starbase)
func (planet *Planet) grow(player *Player) {
	if planet.Cargo.Colonists == 0 {
		// don't grow or reduce if already at zero pop, planet is ded
		return
	}
	planet.addPopulation(planet.Spec.GrowthAmount)
	if planet.Cargo.Colonists == 0 {
		planet.Cargo.Colonists = 1
		planet.PartialPopulation = 0
	}

	if player.Race.Spec.InnateMining {
		planet.Mines = innateMines(player.Race.Spec.InnateMinesFactor, planet.GetPopulation())
	}
}

// reduce the mineral concentrations of a planet after mining based on MineYears.
func (planet *Planet) reduceMineralConcentration(rules *Rules) {
	mineralDecayFactor := rules.MineralDecayFactor // 1.5M by default
	minMineralConcentration := rules.MinMineralConcentration

	// In essence, mine years are like a weird odometer
	// where the amount needed to roll over and decrease
	// mineral conc increases the less minerals remain.

	// Check each mineral type separately
	for _, minType := range MineralTypes {
		conc := max(planet.MineralConcentration.GetAmount(minType), minMineralConcentration) // prevents division by 0

		mineYears := planet.MineYears.GetAmount(minType)
		mineYearsToRollover := mineralDecayFactor / (conc * conc)
		if mineYears <= mineYearsToRollover {
			// mine years below rollover amount; move on
			continue
		}

		newConc := max(conc-(mineYears/mineYearsToRollover), minMineralConcentration)
		planet.MineralConcentration.Set(minType, newConc)
		if newConc == minMineralConcentration {
			// If we're at the minimum, reset mine years to 0
			planet.MineYears.Set(minType, 0)
		} else {
			planet.MineYears.Set(minType, mineYears%mineYearsToRollover)
		}
	}
}
