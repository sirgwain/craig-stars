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
	return &Planet{MapObject: MapObject{Type: MapObjectTypePlanet, Dirty: true, PlayerNum: Unowned}}
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
func (p *Planet) productivePopulation(maxPop int) int {

	return MinInt(p.population(), 3*maxPop)
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
	p.Cargo = Cargo{}
	p.ProductionQueue = []ProductionQueueItem{}
	p.MineYears = Mineral{}
}

// empty this planet of pop, owner
func (p *Planet) emptyPlanet() {
	p.PlayerNum = Unowned
	p.Starbase = nil
	p.Scanner = false
	p.Defenses = 0 // defenses are all gone, rest of the structures can stay
	p.ProductionQueue = []ProductionQueueItem{}
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
	p.Hab = Hab{
		Grav: rules.random.Intn(91) + rules.random.Intn(10),
		Temp: rules.random.Intn(91) + rules.random.Intn(10),
		Rad:  1 + rules.random.Intn(100),
	}
	p.BaseHab = p.Hab
	p.TerraformedAmount = Hab{}

	// from @edmundmk on the Stars! discord, this is
	//  Generate mineral concentration.  There is a 30% chance of a
	//  concentration between 1 and 30.  Higher concentrations have a
	//  distribution centred on 75, minimum 31 and maximum 199.
	//  x = random 1 to 100 inclusive
	//  if x > 30 then
	//     x = 30 + random 0 to 44 inclusive + random 0 to 44 inclusive
	//  end
	//  return x

	// also from @SuicideJunkie about a bonus to germ for high rad
	// Only the exact example given in the help file it seems... "extreme values" is exactly rads being above 85, giving a small bonus to germanium.

	germRadBonus := int(0)
	if p.Hab.Rad >= rules.HighRadGermaniumBonusThreshold {
		germRadBonus = rules.HighRadGermaniumBonus
	}

	p.MineralConcentration = Mineral{
		Ironium:   rules.MinStartingMineralConcentration + rules.random.Intn(rules.MaxStartingMineralConcentration+1),
		Boranium:  rules.MinStartingMineralConcentration + rules.random.Intn(rules.MaxStartingMineralConcentration+1),
		Germanium: rules.MinStartingMineralConcentration + rules.random.Intn(rules.MaxStartingMineralConcentration+1),
	}

	if p.MineralConcentration.Ironium > 30 {
		p.MineralConcentration.Ironium = 30 + rules.random.Intn(45) + rules.random.Intn(45)
	}

	if p.MineralConcentration.Boranium > 30 {
		p.MineralConcentration.Boranium = 30 + rules.random.Intn(45) + rules.random.Intn(45)
	}

	if p.MineralConcentration.Germanium > 30 {
		p.MineralConcentration.Germanium = 30 + rules.random.Intn(45) + germRadBonus + rules.random.Intn(45)
	}

	// check if this planet has a random artifact
	if rules.RandomEventChances[RandomEventAncientArtifact] > rules.random.Float64() {
		p.RandomArtifact = true
	}
}

// Initialize a planet to be a homeworld for a payer with ideal hab, starting mineral concentration, etc
func (p *Planet) initStartingWorld(player *Player, rules *Rules, startingPlanet StartingPlanet, concentration Mineral, surface Mineral) {

	log.Debug().Msgf("Assigning %s to %s as homeworld", p, player)

	p.Homeworld = true
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

	p.MineralConcentration = concentration
	p.Cargo = surface.ToCargo()

	// reset some fields in case this is called on an existing planet for some reason
	p.ProductionQueue = []ProductionQueueItem{}
	p.BaseHab = Hab{}
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

// Get the number of innate mines this player would have on this planet
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
	spec.MiningOutput = planet.getMineralOutput(planet.Mines, race.MineOutput)

	// terraforming
	terraformer := NewTerraformer()
	spec.TerraformAmount = terraformer.getTerraformAmount(planet.Hab, planet.BaseHab, player, player)
	spec.MinTerraformAmount = terraformer.getMinTerraformAmount(planet.Hab, planet.BaseHab, player, player)
	spec.CanTerraform = spec.TerraformAmount.absSum() > 0
	spec.TerraformedHabitability = race.GetPlanetHabitability(planet.Hab.Add(spec.TerraformAmount))

	productivePop := planet.productivePopulation(spec.MaxPopulation)

	if !race.Spec.InnateMining {
		spec.MaxMines = productivePop * race.NumMines / 10000
		spec.MaxPossibleMines = spec.MaxPopulation * race.NumMines / 10000
	}

	if race.Spec.InnateResources {
		spec.ResourcesPerYear = int(math.Sqrt(float64(productivePop)*float64(player.TechLevels.Energy)/float64(race.PopEfficiency)) + .5)
	} else {
		// compute resources from population
		resourcesFromPop := productivePop / (race.PopEfficiency * 100)

		// compute resources from factories
		resourcesFromFactories := planet.Factories * race.FactoryOutput / 10

		spec.ResourcesPerYear = resourcesFromPop + resourcesFromFactories
		spec.MaxFactories = planet.population() * race.NumFactories / 10000
		spec.MaxPossibleFactories = spec.MaxPopulation * race.NumFactories / 10000
	}

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
		maxPossiblePop = p.Starbase.Spec.MaxPopulation
	}
	return roundToNearest100f(math.Max(minMaxPop, float64(maxPossiblePop)*maxPopulationFactor*float64(habitability)/100.0))
}

func (planet *Planet) maxBuildable(t QueueItemType) int {
	switch t {
	case QueueItemTypeAutoMines:
		return MaxInt(0, planet.Spec.MaxMines-planet.Mines)
	case QueueItemTypeMine:
		return MaxInt(0, planet.Spec.MaxPossibleMines-planet.Mines)
	case QueueItemTypeAutoFactories:
		return MaxInt(0, planet.Spec.MaxFactories-planet.Factories)
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
	planet.setPopulation(planet.population() + planet.Spec.GrowthAmount)

	if player.Race.Spec.InnateMining {
		productivePop := planet.productivePopulation(planet.Spec.MaxPopulation)
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
