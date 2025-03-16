package cs

import (
	"fmt"
	"math"
)

// Planets are the only static and constant MapObject. They don't move and they can't be destroyed.
// Players also start the game knowing all planet names and locations.
// I suppose these should have been named Stars, since they represent a star system, ah well..
type Planet struct {
	GameDBObject         `tstype:",extends"`
	MapObject            `tstype:",extends"`
	PlanetOrders         `tstype:",extends"`
	Hab                  Hab        `json:"hab"`
	BaseHab              Hab        `json:"baseHab"`
	TerraformedAmount    Hab        `json:"terraformedAmount"`
	MineralConcentration Mineral    `json:"mineralConcentration"`
	MineYears            Mineral    `json:"mineYears"`
	Cargo                Cargo      `json:"cargo"`
	Mines                int        `json:"mines"`
	Factories            int        `json:"factories"`
	Defenses             int        `json:"defenses"`
	Homeworld            bool       `json:"homeworld,omitempty"`
	Scanner              bool       `json:"scanner,omitempty"`
	Spec                 PlanetSpec `json:"spec"`
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
	PlanetStarbaseSpec                        `tstype:",extends"`
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
	MiningOutput                              Mineral `json:"miningOutput,omitempty"`
	Population                                int     `json:"population,omitempty"`
	PopulationDensity                         float64 `json:"populationDensity,omitempty"`
	ResourcesPerYear                          int     `json:"resourcesPerYear,omitempty"`
	ResourcesPerYearAvailable                 int     `json:"resourcesPerYearAvailable,omitempty"`
	ResourcesPerYearResearch                  int     `json:"resourcesPerYearResearch,omitempty"`
	ResourcesPerYearResearchEstimatedLeftover int     `json:"resourcesPerYearResearchEstimatedLeftover,omitempty"`
	Scanner                                   string  `json:"scanner,omitempty"`
	ScanRange                                 int     `json:"scanRange,omitempty"`
	ScanRangePen                              int     `json:"scanRangePen,omitempty"`
	TerraformAmount                           Hab     `json:"terraformAmount,omitempty"`
	MinTerraformAmount                        Hab     `json:"minTerraformAmount,omitempty"`
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

func (p *Planet) String() string {
	return fmt.Sprintf("Planet %v", p.MapObject)
}

func (p *Planet) population() int {
	return p.Cargo.Colonists * 100
}

// get the amount population that is productive for producing resources.
// This takes into account overcrowding
func productivePopulation(pop, maxPop int, overcrowdPenalty, overcrowdResourceMax float64) int {
	popOverCap := float64(pop) + Max(0, float64(pop-maxPop)*overcrowdPenalty)
	// TODO: Refactor to use RoundTo100 function once it gets refactored to
	// take function argument
	return int(Min(
		float64(maxPop)*(1+overcrowdResourceMax), popOverCap)/100) * 100
}

// get the population that will operate installations
// (it just maxes at max pop)
func productiveInstallationPopulation(pop, maxPop int) int {
	return Min(pop, maxPop)
}

func (p *Planet) setPopulation(pop int) {
	p.Cargo.Colonists = pop / 100
}

// return true if this planet can build a ship with a given mass
func (p *Planet) CanBuild(mass int) bool {
	return p.Spec.HasStarbase && (p.Starbase.Spec.SpaceDock == UnlimitedSpaceDock || p.Starbase.Spec.SpaceDock >= mass)
}

// populate a starbase design for a planet
func (p *Planet) PopulateStarbaseDesign(player *Player) error {
	if p.Starbase != nil {
		if len(p.Starbase.Tokens) != 1 {
			return fmt.Errorf("planet %s starbase has no tokens", p.Name)
		}
		designNum := p.Starbase.Tokens[0].DesignNum
		design := player.GetDesign(designNum)

		if design == nil {
			return fmt.Errorf("player %v does not have design %d", player, designNum)
		}
		p.Starbase.Tokens[0].design = player.GetDesign(designNum)
	}
	return nil
}

// add designs to each production queue item
func (p *Planet) PopulateProductionQueueDesigns(player *Player) error {
	for i := range p.ProductionQueue {
		item := &p.ProductionQueue[i]
		if item.Type == QueueItemTypeStarbase || item.Type == QueueItemTypeShipToken {
			design := player.GetDesign(item.DesignNum)
			if design == nil {
				return fmt.Errorf("player %v does not have design %d", player, item.DesignNum)
			}
			item.design = design
		}
	}
	return nil
}

// populate the costs of each item in the planet production queue
func (p *Planet) PopulateProductionQueueEstimates(rules *Rules, player *Player) error {
	// populate completion estimates
	completionEstimator := NewCompletionEstimator()
	var err error
	p.ProductionQueue, p.Spec.ResourcesPerYearResearchEstimatedLeftover, err = completionEstimator.GetProductionWithEstimates(rules, player, *p)
	return err
}

// empty this planet of pop & owner, typically used when transferring or losing ownership
func (p *Planet) emptyPlanet() {
	p.PlayerNum = Unowned
	p.Starbase = nil
	// defenses & scanner disappear, other structures stay though
	p.Scanner = false
	// clear any production or other orders from the previous owner
	p.Defenses = 0
	p.PlanetOrders = PlanetOrders{}
	p.setPopulation(0)
	p.Spec = PlanetSpec{}
	// reset any instaforming
	p.Hab = p.BaseHab.Add(p.TerraformedAmount)
}

// randomize a planet with new hab range, minerals, etc
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
			concAmount := minConc.GetAmount(minType)
			if concAmount < 40 {
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
			Ironium:   minConc.Ironium + rules.random.Intn(99-Min(minConc.Ironium, 98))/2,
			Boranium:  minConc.Boranium + rules.random.Intn(99-Min(minConc.Boranium, 98))/2,
			Germanium: minConc.Germanium + rules.random.Intn(99-Min(minConc.Germanium, 98))/2,
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
		p.Hab.Grav = habCenter.Grav + int(float64((habWidth.Grav-rules.random.Intn(habWidth.Grav-1)))/2*startingPlanet.HabPenaltyFactor)
	}
	if !player.Race.ImmuneTemp {
		p.Hab.Temp = habCenter.Temp + int(float64((habWidth.Temp-rules.random.Intn(habWidth.Temp-1)))/2*startingPlanet.HabPenaltyFactor)
	}
	if !player.Race.ImmuneRad {
		p.Hab.Rad = habCenter.Rad + int(float64((habWidth.Rad-rules.random.Intn(habWidth.Rad-1)))/2*startingPlanet.HabPenaltyFactor)
	}
	// BaseHab is the same as Hab
	p.BaseHab = p.Hab

	p.MineralConcentration = concentration
	p.Cargo = surface.ToCargo()

	// empty queue, no terraform
	p.ProductionQueue = []ProductionQueueItem{}
	p.TerraformedAmount = Hab{}

	raceSpec := player.Race.Spec

	// set the homeworld pop to our starting planet pop
	p.setPopulation(int(float64(startingPlanet.Population) * raceSpec.StartingPopulationFactor))

	if raceSpec.InnateMining {
		p.Mines = innateMines(raceSpec.InnateMinesFactor, p.population())
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
		plan := player.ProductionPlans[0]
		plan.Apply(p)
	}

}

// set this planet's starbase on this planet
func (p *Planet) setStarbase(starbase *Fleet) {
	p.Starbase = starbase
	p.PacketSpeed = starbase.Spec.SafePacketSpeed
}

// Get the number of innate mines this player would have on this planet
func innateMines(innateMinesFactor float64, population int) int {
	return int(math.Sqrt(float64(population)) * innateMinesFactor)
}

// Get the innate scanning distance this player would have on this planet
func innateScanner(innateScannerFactor float64, population int) int {
	return int(math.Sqrt(float64(population) * innateScannerFactor))
}

func (p *Planet) shortestDistanceToPlanets(otherPlanets *[]*Planet) float64 {
	minDistanceSquared := math.MaxFloat64
	for _, planet := range *otherPlanets {
		distSquared := p.Position.DistanceSquaredTo(planet.Position)
		minDistanceSquared = math.Min(minDistanceSquared, distSquared)
	}
	return math.Sqrt(minDistanceSquared)
}

// get the mineral output of a planet based on mineOutput (10 for remote mining)
// TODO: Add fractional mineral outputs (% chance for extra) for sub-integer amounts
func (p *Planet) getMineralOutput(numMines int, mineOutput int) Mineral {
	return Mineral{
		int(float64(p.MineralConcentration.Ironium*numMines*mineOutput) / 1000),
		int(float64(p.MineralConcentration.Boranium*numMines*mineOutput) / 1000),
		int(float64(p.MineralConcentration.Germanium*numMines*mineOutput) / 1000),
	}
}

// get how much a player will grow on a planet, given a max population the player can have on the planet
func (p *Planet) getGrowthAmount(player *Player, maxPopulation int, populationOvercrowdDieoffRate, populationOvercrowdDieoffRateMax float64) int {
	race := &player.Race
	growthFactor := race.Spec.GrowthFactor
	capacity := float64(p.population()) / float64(maxPopulation)
	habValue := race.GetPlanetHabitability(p.Hab)
	if habValue > 0 {
		popGrowth := int(float64(p.population())*float64(race.GrowthRate)*growthFactor/100.0*float64(habValue)/100.0 + .5)

		if capacity > 1 {
			// overpopulation calcs: https://wiki.starsautohost.org/wiki/Overpopulation
			// Population Death from overcrowding is 0.04% per % over 100% cap.
			// Thus a 200% capacity planet is 100% over and thus has (0.04 * 100 = 4%) a 4% death rate. This maxes out at 400% capacity at 12%
			// Credit: Thomas Harley
			// In addition to deaths:
			// excess population on overcrowded planets cannot work factories or mines
			// the first 200% overpopulation (300% capacity) only produce half their normal production(for a net population production of 200%).
			// Population over 300% produce nothing.

			dieoffPercent := Clamp((1-capacity)*populationOvercrowdDieoffRate, -populationOvercrowdDieoffRateMax, 0)
			popGrowth = int(float64(p.population()) * float64(dieoffPercent))
		} else if capacity > .25 {
			crowdingFactor := 16.0 / 9.0 * (1.0 - capacity) * (1.0 - capacity)
			popGrowth = int(float64(popGrowth) * crowdingFactor)
		}

		// round to the nearest 100 colonists
		return roundToNearest100(popGrowth)
	} else {
		// kill off (habValue / 10)% colonists every year. I.e. a habValue of -4% kills off .4%
		deathAmount := int(float64(p.population()) * (float64(habValue) / 1000.0))
		return roundToNearest100(Clamp(deathAmount, deathAmount, -100))
	}
}

func computePlanetSpec(rules *Rules, player *Player, planet *Planet) PlanetSpec {
	spec := PlanetSpec{}
	race := &player.Race

	// hab/pop
	spec.Habitability = race.GetPlanetHabitability(planet.Hab)
	spec.MaxPopulation = planet.getMaxPopulation(rules, player, spec.Habitability)
	spec.Population = planet.population()
	if spec.MaxPopulation > 0 {
		spec.PopulationDensity = float64(planet.population()) / float64(spec.MaxPopulation)
	}
	spec.GrowthAmount = planet.getGrowthAmount(player, spec.MaxPopulation, rules.PopulationOvercrowdDieoffRate, rules.PopulationOvercrowdDieoffRateMax)

	// terraforming
	terraformer := NewTerraformer()
	spec.TerraformAmount = terraformer.GetTerraformAmount(planet.Hab, planet.BaseHab, player, player)
	spec.MinTerraformAmount = terraformer.GetMinTerraformAmount(planet.Hab, planet.BaseHab, player, player)
	spec.CanTerraform = spec.TerraformAmount.absSum() > 0
	spec.TerraformedHabitability = race.GetPlanetHabitability(planet.Hab.Add(spec.TerraformAmount))

	// population will generate resources up to 3x max pop, but they can only
	// operate structures up to max pop
	productivePop := productivePopulation(spec.Population, spec.MaxPopulation, rules.PopulationOvercrowdResourcePenalty, rules.PopulationOvercrowdResourceMax)
	installationPop := productiveInstallationPopulation(spec.Population, spec.MaxPopulation)

	if !race.Spec.InnateMining {
		spec.MaxMines = getMaxInstallations(player.Race.NumMines, installationPop)
		spec.MaxPossibleMines = spec.MaxPopulation * race.NumMines / 10000
	} else {
		spec.MaxMines = planet.Mines
	}

	// Compute resources per year
	spec.computeResourcesPerYear(player, planet.Factories, productivePop, installationPop)
	spec.MiningOutput = planet.getMineralOutput(Min(spec.MaxMines, planet.Mines), race.MineOutput)
	spec.computeResourcesPerYearAvailable(player, planet)

	if race.Spec.CanBuildDefenses {
		spec.MaxDefenses = 100
		spec.Defense = player.Spec.Defense.Name
		spec.computeDefenseCoverage(rules, player.Spec.Defense.DefenseCoverage, planet.Defenses)
	}

	if race.Spec.InnateScanner {
		// calculate AR organic scan ranes
		spec.Scanner = "Organic"
		spec.ScanRange = int(float64(innateScanner(player.Race.Spec.InnateScannerFactor, productivePop)) * player.Race.Spec.ScanRangeFactor)
		if !player.Race.Spec.NoAdvancedScanners && planet.Starbase != nil {
			spec.ScanRangePen = int(float64(spec.ScanRange) * planet.Starbase.Spec.InnateScanRangePenFactor)
		}
	} else if planet.Scanner {
		scanner := player.Spec.PlanetaryScanner
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
	spec.HasStarbase = starbase != nil
	if starbase != nil {
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
	}

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
func (spec *PlanetSpec) computeResourcesPerYear(player *Player, numFacts, productivePop, installationPop int) {
	if player.Race.Spec.InnateResources {
		// Compute resources for AR
		spec.ResourcesPerYear = int(math.Ceil(float64(spec.Habitability) / 100 * // Confirmed: AR resources round up in base game
			math.Sqrt(float64(productivePop*player.TechLevels.Energy)/float64(player.Race.PopEfficiency))))
	} else {
		// compute resources from population & factories
		resourcesFromPop := productivePop / (player.Race.PopEfficiency * 100)

		spec.MaxFactories = getMaxInstallations(player.Race.NumFactories, installationPop)
		spec.MaxPossibleFactories = spec.MaxPopulation * player.Race.NumFactories / 10000 // factory count rounds down
		resourcesFromFactories := int(math.Ceil(float64(Min(numFacts, spec.MaxFactories)*player.Race.FactoryOutput) / 10))

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
		return roundToNearest100(float64(p.Starbase.Spec.MaxPopulation) * maxPopulationFactor)
	}

	// Habitability is floored at 5% when determining max population
	// (or 25% for AR races)
	if habitability < player.Race.Spec.MinHabFloor {
		habitability = player.Race.Spec.MinHabFloor
	}

	// TODO: Refactor to make this FLOOR to 100
	return roundToNearest100(float64(rules.MaxPopulation*habitability) * maxPopulationFactor / 100.0)
}

// return the maximum number count operable by the given population
func getMaxInstallations(installationsPer10K, population int) int {
	return population * installationsPer10K / 10000
}

func (planet *Planet) maxBuildable(player *Player, t QueueItemType) int {
	switch t {
	case QueueItemTypeAutoMines:
		// for autobuild purposes, the maxFactories is next year's pop
		futurePop := productiveInstallationPopulation(planet.population()+planet.Spec.GrowthAmount, planet.Spec.MaxPopulation)
		maxMines := getMaxInstallations(player.Race.NumMines, futurePop)
		return Max(0, maxMines-planet.Mines)
	case QueueItemTypeMine:
		return Max(0, planet.Spec.MaxPossibleMines-planet.Mines)
	case QueueItemTypeAutoFactories:
		// for autobuild purposes, the maxFactories is next year's pop
		futurePop := productiveInstallationPopulation(planet.population()+planet.Spec.GrowthAmount, planet.Spec.MaxPopulation)
		maxFactories := getMaxInstallations(player.Race.NumFactories, futurePop)
		return Max(0, maxFactories-planet.Factories)
	case QueueItemTypeFactory:
		return Max(0, planet.Spec.MaxPossibleFactories-planet.Factories)
	case QueueItemTypeAutoDefenses, QueueItemTypeDefenses:
		return Max(0, planet.Spec.MaxDefenses-planet.Defenses)
	case QueueItemTypeTerraformEnvironment, QueueItemTypeAutoMaxTerraform:
		return planet.Spec.TerraformAmount.absSum()
	case QueueItemTypeAutoMinTerraform:
		return planet.Spec.MinTerraformAmount.absSum()
	case QueueItemTypeStarbase:
		return 1
	case QueueItemTypePlanetaryScanner:
		if planet.Scanner {
			return 0
		}
		return 1
	case QueueItemTypeGenesisDevice:
		return 1
	}
	// default to infinite
	return Infinite
}

// mine minerals on this planet
func (planet *Planet) mine(rules *Rules) {
	planet.Cargo = planet.Cargo.AddMineral(planet.Spec.MiningOutput)
	planet.MineYears = planet.MineYears.AddToAll(planet.Mines)
	planet.reduceMineralConcentration(rules)
}

// grow pop on this planet (or starbase)
func (planet *Planet) grow(player *Player) {
	if planet.population() == 0 {
		// don't grow or reduce if at zero pop, planet is gone
		return
	}
	planet.setPopulation(Max(100, planet.population()+planet.Spec.GrowthAmount))

	if player.Race.Spec.InnateMining {
		planet.Mines = innateMines(player.Race.Spec.InnateMinesFactor, planet.population())
	}
}

// reduce the mineral concentrations of a planet after mining.
func (planet *Planet) reduceMineralConcentration(rules *Rules) {
	mineralDecayFactor := rules.MineralDecayFactor
	minMineralConcentration := rules.MinMineralConcentration
	if planet.Homeworld {
		minMineralConcentration = rules.MinHomeworldMineralConcentration
	}

	planetMineYears := planet.MineYears.ToSlice()
	planetMineralConcentration := planet.MineralConcentration.ToSlice()
	for i, conc := range planetMineralConcentration {
		if conc < minMineralConcentration {
			// can't have less than min, make sure we have that at least
			conc = minMineralConcentration
			planetMineralConcentration[i] = conc
		}

		minesPer := mineralDecayFactor / conc / conc
		mineYears := planetMineYears[i]
		if mineYears > minesPer {
			conc -= mineYears / minesPer
			if conc < minMineralConcentration {
				conc = minMineralConcentration
			}
			mineYears %= minesPer

			planetMineYears[i] = mineYears
			planetMineralConcentration[i] = conc
		}
	}
	planet.MineYears = NewMineral(planetMineYears[0], planetMineYears[1], planetMineYears[2])
	planet.MineralConcentration = NewMineral(planetMineralConcentration[0], planetMineralConcentration[1], planetMineralConcentration[2])
}
