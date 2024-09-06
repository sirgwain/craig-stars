package cs

import (
	"fmt"
	"math"

	"github.com/rs/zerolog/log"
)

// Planets are the only static and constant MapObject. They don't move and they can't be destroyed.
// Players also start the game knowing all planet names and locations.
// I suppose these should have been named Stars, since they represent a star system, ah well..
type Planet struct {
	MapObject
	PlanetOrders
	Hab                  Hab        `json:"hab,omitempty"`
	BaseHab              Hab        `json:"baseHab,omitempty"`
	TerraformedAmount    Hab        `json:"terraformedAmount,omitempty"`
	MineralConcentration Mineral    `json:"mineralConcentration,omitempty"`
	MineYears            Mineral    `json:"mineYears,omitempty"`
	Cargo                Cargo      `json:"cargo,omitempty"`
	Mines                int        `json:"mines,omitempty"`
	Factories            int        `json:"factories,omitempty"`
	Defenses             int        `json:"defenses,omitempty"`
	Homeworld            bool       `json:"homeworld,omitempty"`
	Scanner              bool       `json:"scanner,omitempty"`
	Spec                 PlanetSpec `json:"spec,omitempty"`
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
	Contested                                 bool    `json:"contested,omitempty"`
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

func (p *Planet) WithMines(mines int) *Planet {
	p.Mines = mines
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
	return fmt.Sprintf("Planet %s", &p.MapObject)
}

func (p *Planet) population() int {
	return p.Cargo.Colonists * 100
}

// get the population that is productive. This takes into account overcrowding
// anything over 3x is unproductive
func (p *Planet) productivePopulation(pop, maxPop int) int {
	return MinInt(pop, 3*maxPop)
}

func (p *Planet) setPopulation(pop int) {
	p.Cargo.Colonists = pop / 100
}

// true if this planet can build a ship with a given mass
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
func (p *Planet) PopulateProductionQueueEstimates(rules *Rules, player *Player) {
	// populate completion estimates
	completionEstimator := NewCompletionEstimator()
	p.ProductionQueue, p.Spec.ResourcesPerYearResearchEstimatedLeftover = completionEstimator.GetProductionWithEstimates(rules, player, *p)
}

func (p *Planet) reset() {
	p.Hab = Hab{}
	p.BaseHab = Hab{}
	p.TerraformedAmount = Hab{}
	p.MineralConcentration = Mineral{}
	p.ProductionQueue = []ProductionQueueItem{}
	p.MineYears = Mineral{}
}

// empty this planet of pop, owner
func (p *Planet) emptyPlanet() {
	p.PlayerNum = Unowned
	p.Starbase = nil
	p.Scanner = false
	p.Defenses = 0                  // defenses are all gone, rest of the structures can stay
	p.PlanetOrders = PlanetOrders{} // clear any orders from previous owner
	p.setPopulation(0)
	p.Spec = PlanetSpec{}
	// reset any instaforming
	p.Hab = p.BaseHab.Add(p.TerraformedAmount)
}

// randomize a planet with new hab range, minerals, etc
func (p *Planet) randomize(rules *Rules) {
	p.reset()

	// From @SuicideJunkie's tests and @edmundmk's previous research, grav and temp are weighted slightly towards
	// the center, rad is completely random
	// @edmundmk:
	// "I'm certain gravity and temperature probability is constant between 10 and 90 inclusive, and falls off towards 0 and 100.
	// It never generates 0 or 100 so I have to change my random formula to (1 to 90)+(0 to 9)
	// damn you all for sucking me into stars! again lol"
	//
	// update: hab is 1 to 99
	p.Hab = Hab{
		Grav: 1 + rules.random.Intn(90) + rules.random.Intn(10),
		Temp: 1 + rules.random.Intn(90) + rules.random.Intn(10),
		Rad:  1 + rules.random.Intn(99),
	}
	p.BaseHab = p.Hab
	p.TerraformedAmount = Hab{}
	p.MineralConcentration = randomizeMinerals(rules, p.Hab.Rad)

	// check if this planet has a random artifact
	if rules.RandomEventChances[RandomEventAncientArtifact] >= rules.random.Float64() {
		p.RandomArtifact = true
	}
}

// Randomize a planet's mineral concentration within bounds set in Rules
func randomizeMinerals(rules *Rules, rad int) Mineral {

	// These two variables are the shape of the normal distribution
	// based on comparing it with Stars! output
	mean := 80.0
	variance := 20.0

	// These two are the min and max of the minerals to be returned,
	// They clamp the results
	mMin := rules.MinStartingMineralConcentration
	mMax := rules.MaxStartingMineralConcentration

	// creates a mineral concentration
	minConc := Mineral{
		Ironium:   NormalSample(rules.random, mean, variance, mMax),
		Boranium:  NormalSample(rules.random, mean, variance, mMax),
		Germanium: NormalSample(rules.random, mean, variance, mMax),
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
			value := 1 + rules.random.Intn(rules.MinStartingMineralConcentration)
			minConc.Set(mineralType, value)
		} else {
			limiter++
			for limiter < 16 {
				mineralType := MineralTypes[rules.random.Intn(len(MineralTypes))]
				value := 1 + rules.random.Intn(rules.MinStartingMineralConcentration)
				minConc.Set(mineralType, value)

				limiter *= 2
			}
		}
	}

	// we have high rad, add some bonus minerals
	if rad >= rules.HighRadMineralConcentrationBonusThreshold {
		minConc = Mineral{
			Ironium:   minConc.Ironium + rules.random.Intn(99-MinInt(minConc.Ironium, 98))/2,
			Boranium:  minConc.Boranium + rules.random.Intn(99-MinInt(minConc.Boranium, 98))/2,
			Germanium: minConc.Germanium + rules.random.Intn(99-MinInt(minConc.Germanium, 98))/2,
		}
	}

	minConc.Ironium = Clamp(minConc.Ironium, mMin, mMax)
	minConc.Boranium = Clamp(minConc.Boranium, mMin, mMax)
	minConc.Germanium = Clamp(minConc.Germanium, mMin, mMax)

	return minConc
}

// Initialize a planet to be a homeworld for a payer with ideal hab, starting mineral concentration, etc
func (p *Planet) initStartingWorld(player *Player, rules *Rules, startingPlanet StartingPlanet, concentration Mineral, surface Mineral) {

	log.Debug().Msgf("Assigning %s to %s as homeworld", p, player)

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
		p.Mines = p.innateMines(player, p.population())
		p.Factories = 0
	} else {
		p.Mines = rules.StartingMines
		p.Factories = rules.StartingFactories
	}

	if raceSpec.CanBuildDefenses {
		p.Defenses = rules.StartingDefenses
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
func (p *Planet) setStarbase(rules *Rules, player *Player, starbase *Fleet) {
	p.Starbase = starbase
	p.PacketSpeed = starbase.Spec.SafePacketSpeed
}

// Get the number of innate mines this player would have on this planet
func (p *Planet) innateMines(player *Player, population int) int {
	if player.Race.Spec.InnateMining {
		return int(math.Sqrt(float64(population)) * float64(player.Race.Spec.InnatePopulationFactor))
	}
	return 0
}

// Get the innate scanning distance this player would have on this planet
func (p *Planet) innateScanner(player *Player, population int) int {
	if player.Race.Spec.InnateScanner {
		return int(math.Sqrt(float64(population) * float64(player.Race.Spec.InnatePopulationFactor)))
	}
	return 0
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
func (p *Planet) getMineralOutput(numMines int, mineOutput int) Mineral {
	return Mineral{
		int(float64(p.MineralConcentration.Ironium) / 100 * float64(numMines) / 10 * float64(mineOutput)),
		int(float64(p.MineralConcentration.Boranium) / 100 * float64(numMines) / 10 * float64(mineOutput)),
		int(float64(p.MineralConcentration.Germanium) / 100 * float64(numMines) / 10 * float64(mineOutput)),
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

			dieoffPercent := ClampFloat64((1-capacity)*populationOvercrowdDieoffRate, -populationOvercrowdDieoffRateMax, 0)
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
	spec.TerraformAmount = terraformer.getTerraformAmount(planet.Hab, planet.BaseHab, player, player)
	spec.MinTerraformAmount = terraformer.getMinTerraformAmount(planet.Hab, planet.BaseHab, player, player)
	spec.CanTerraform = spec.TerraformAmount.absSum() > 0
	spec.TerraformedHabitability = race.GetPlanetHabitability(planet.Hab.Add(spec.TerraformAmount))

	productivePop := planet.productivePopulation(spec.Population, spec.MaxPopulation)

	if !race.Spec.InnateMining {
		spec.MaxMines = planet.getMaxMines(player, productivePop)
		spec.MaxPossibleMines = spec.MaxPopulation * race.NumMines / 10000
	} else {
		spec.MaxMines = planet.Mines
	}

	if race.Spec.InnateResources {
		spec.ResourcesPerYear = int(math.Sqrt(float64(productivePop)*float64(player.TechLevels.Energy)/float64(race.PopEfficiency)) + .5)
	} else {
		// compute resources from population
		resourcesFromPop := productivePop / (race.PopEfficiency * 100)

		spec.MaxFactories = planet.getMaxFactories(player, productivePop)
		spec.MaxPossibleFactories = spec.MaxPopulation * race.NumFactories / 10000

		// compute resources from factories
		resourcesFromFactories := MinInt(planet.Factories, spec.MaxFactories) * race.FactoryOutput / 10
		spec.ResourcesPerYear = resourcesFromPop + resourcesFromFactories
	}

	spec.MiningOutput = planet.getMineralOutput(MinInt(spec.MaxMines, planet.Mines), race.MineOutput)
	spec.computeResourcesPerYearAvailable(player, planet)

	if race.Spec.CanBuildDefenses {
		spec.MaxDefenses = 100
		spec.Defense = player.Spec.Defense.Name
		spec.DefenseCoverage = float64(1.0 - (math.Pow((1 - (player.Spec.Defense.DefenseCoverage / 100)), float64(Clamp(planet.Defenses, 0, spec.MaxDefenses)))))
		spec.DefenseCoverageSmart = float64(1.0 - (math.Pow((1 - (player.Spec.Defense.DefenseCoverage / 100 * rules.SmartDefenseCoverageFactor)), float64(Clamp(planet.Defenses, 0, spec.MaxDefenses)))))
	}

	if race.Spec.InnateScanner {
		spec.Scanner = "Organic"
		spec.ScanRange = int(float64(planet.innateScanner(player, productivePop)) * player.Race.Spec.ScanRangeFactor)
		if !player.Race.Spec.NoAdvancedScanners && planet.Starbase != nil {
			spec.ScanRangePen = int(float64(spec.ScanRange) * planet.Starbase.Spec.InnateScanRangePenFactor)
		}
	} else if planet.Scanner {
		scanner := player.Spec.PlanetaryScanner
		spec.Scanner = scanner.Name
		spec.ScanRange = int(float64(scanner.ScanRange) * player.Race.Spec.ScanRangeFactor)
		spec.ScanRangePen = scanner.ScanRangePen
	}

	spec.PlanetStarbaseSpec = computePlanetStarbaseSpec(rules, player, planet)

	return spec
}

func computePlanetStarbaseSpec(rules *Rules, player *Player, planet *Planet) PlanetStarbaseSpec {
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

// update a planet spec's resources per year
// this is called by the main ComputePlanetSpec as well as anytime a player
// updates a planet's ContributesOnlyLeftoverToResearch field
func (spec *PlanetSpec) computeResourcesPerYearAvailable(player *Player, planet *Planet) {
	if planet.ContributesOnlyLeftoverToResearch {
		spec.ResourcesPerYearAvailable = spec.ResourcesPerYear
		spec.ResourcesPerYearResearch = 0
	} else {
		spec.ResourcesPerYearResearch = int(float64(spec.ResourcesPerYear) * float64(player.ResearchAmount) / 100.0)
		spec.ResourcesPerYearAvailable = spec.ResourcesPerYear - spec.ResourcesPerYearResearch
	}
}

// get the max population for this planet for a player with a hab rating
func (p *Planet) getMaxPopulation(rules *Rules, player *Player, habitability int) int {
	maxPopulationFactor := 1 + player.Race.Spec.MaxPopulationOffset
	maxPossiblePop := rules.MaxPopulation

	// a planet's max pop can't go lower than 5% of a race's max, i.e.
	// for a regular race with 1 million max pop, the minimum max population is 50,000
	minMaxPop := float64(maxPossiblePop) * maxPopulationFactor * rules.MinMaxPopulationPercent

	if player.Race.Spec.LivesOnStarbases && p.PlayerNum == player.Num {
		return roundToNearest100f(float64(p.Starbase.Spec.MaxPopulation) * maxPopulationFactor)
	}
	return roundToNearest100f(math.Max(minMaxPop, float64(maxPossiblePop)*maxPopulationFactor*float64(habitability)/100.0))
}

// get max factories for a population
func (p *Planet) getMaxFactories(player *Player, population int) int {
	if player.Race.Spec.InnateResources {
		return 0
	} else {
		return population * player.Race.NumFactories / 10000
	}
}

// get max mines for a population
func (p *Planet) getMaxMines(player *Player, population int) int {
	if player.Race.Spec.InnateResources {
		return 0
	} else {
		return population * player.Race.NumMines / 10000
	}
}

func (planet *Planet) maxBuildable(player *Player, t QueueItemType) int {
	switch t {
	case QueueItemTypeAutoMines:
		// for autobuild purposes, the maxFactories is next year's pop
		futurePop := planet.productivePopulation(planet.population()+planet.Spec.GrowthAmount, planet.Spec.MaxPopulation)
		maxMines := planet.getMaxMines(player, futurePop)
		return MaxInt(0, maxMines-planet.Mines)
	case QueueItemTypeMine:
		return MaxInt(0, planet.Spec.MaxPossibleMines-planet.Mines)
	case QueueItemTypeAutoFactories:
		// for autobuild purposes, the maxFactories is next year's pop
		futurePop := planet.productivePopulation(planet.population()+planet.Spec.GrowthAmount, planet.Spec.MaxPopulation)
		maxFactories := planet.getMaxFactories(player, futurePop)
		return MaxInt(0, maxFactories-planet.Factories)
	case QueueItemTypeFactory:
		return MaxInt(0, planet.Spec.MaxPossibleFactories-planet.Factories)
	case QueueItemTypeAutoDefenses:
		fallthrough
	case QueueItemTypeDefenses:
		return MaxInt(0, planet.Spec.MaxDefenses-planet.Defenses)
	case QueueItemTypeTerraformEnvironment:
	case QueueItemTypeAutoMaxTerraform:
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
	planet.MineYears = planet.MineYears.AddInt(planet.Mines)
	planet.reduceMineralConcentration(rules)
}

// grow pop on this planet (or starbase)
func (planet *Planet) grow(player *Player) {
	planet.setPopulation(MaxInt(100, planet.population()+planet.Spec.GrowthAmount))

	if player.Race.Spec.InnateMining {
		productivePop := planet.productivePopulation(planet.population(), planet.Spec.MaxPopulation)
		planet.Mines = planet.innateMines(player, productivePop)
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
	for i := 0; i < 3; i++ {
		conc := planetMineralConcentration[i]
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
	planet.MineYears = NewMineral(planetMineYears)
	planet.MineralConcentration = NewMineral(planetMineralConcentration)
}
