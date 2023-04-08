package cs

import (
	"fmt"
	"math"

	"github.com/rs/zerolog/log"
)

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
	PacketSpeed          int        `json:"packetSpeed,omitempty"`
	BonusResources       int        `json:"-"`
	Spec                 PlanetSpec `json:"spec,omitempty"`
	starbase             *Fleet
}

type PlanetOrders struct {
	ContributesOnlyLeftoverToResearch bool                  `json:"contributesOnlyLeftoverToResearch,omitempty"`
	ProductionQueue                   []ProductionQueueItem `json:"productionQueue,omitempty"`
	TargetType                        MapObjectType         `json:"targetType,omitempty"`
	TargetNum                         int                   `json:"targetNum,omitempty"`
	TargetPlayerNum                   int                   `json:"targetPlayerNum,omitempty"`
}



type PlanetSpec struct {
	Habitability              int     `json:"habitability,omitempty"`
	TerraformedHabitability   int     `json:"terraformedHabitability,omitempty"`
	MaxMines                  int     `json:"maxMines,omitempty"`
	MaxPossibleMines          int     `json:"maxPossibleMines,omitempty"`
	MaxFactories              int     `json:"maxFactories,omitempty"`
	MaxPossibleFactories      int     `json:"maxPossibleFactories,omitempty"`
	MaxDefenses               int     `json:"maxDefenses,omitempty"`
	PopulationDensity         float64 `json:"populationDensity,omitempty"`
	MaxPopulation             int     `json:"maxPopulation,omitempty"`
	GrowthAmount              int     `json:"growthAmount,omitempty"`
	MineralOutput             Mineral `json:"mineralOutput,omitempty"`
	ResourcesPerYear          int     `json:"resourcesPerYear,omitempty"`
	ResourcesPerYearAvailable int     `json:"resourcesPerYearAvailable,omitempty"`
	ResourcesPerYearResearch  int     `json:"resourcesPerYearResearch,omitempty"`
	Defense                   string  `json:"defense,omitempty"`
	DefenseCoverage           float64 `json:"defenseCoverage,omitempty"`
	DefenseCoverageSmart      float64 `json:"defenseCoverageSmart,omitempty"`
	Scanner                   string  `json:"scanner,omitempty"`
	ScanRange                 int     `json:"scanRange,omitempty"`
	ScanRangePen              int     `json:"scanRangePen,omitempty"`
	CanTerraform              bool    `json:"canTerraform,omitempty"`
	TerraformAmount           Hab     `json:"terraformAmount,omitempty"`
	HasStarbase               bool    `json:"hasStarbase,omitempty"`
	HasStargate               bool    `json:"hasStargate,omitempty"`
	DockCapacity              int     `json:"dockCapacity,omitempty"`
	SafeHullMass              int     `json:"safeHullMass,omitempty"`
	SafeRange                 int     `json:"safeRange,omitempty"`
	MaxHullMass               int     `json:"maxHullMass,omitempty"`
	MaxRange                  int     `json:"maxRange,omitempty"`
}

func (item *ProductionQueueItem) String() string {
	return fmt.Sprintf("ProductionQueueItem %d %s (%s)", item.Quantity, item.Type, item.DesignName)
}

func NewPlanet() *Planet {
	return &Planet{MapObject: MapObject{Type: MapObjectTypePlanet, Dirty: true, PlayerNum: Unowned}}
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

func (p *Planet) setPopulation(pop int) {
	p.Cargo.Colonists = pop / 100
}

// true if this planet can build a ship with a given mass
func (p *Planet) CanBuild(mass int) bool {
	return p.Spec.HasStarbase && (p.starbase.Spec.SpaceDock == UnlimitedSpaceDock || p.starbase.Spec.SpaceDock >= mass)
}

func (p *Planet) empty() {
	p.Hab = Hab{}
	p.BaseHab = Hab{}
	p.TerraformedAmount = Hab{}
	p.MineralConcentration = Mineral{}
	p.Cargo = Cargo{}
	p.ProductionQueue = []ProductionQueueItem{}
	p.MineYears = Mineral{}
}

// randomize a planet with new hab range, minerals, etc
func (p *Planet) randomize(rules *Rules) {
	p.empty()

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
}

// Initialize a planet to be a homeworld for a payer with ideal hab, starting mineral concentration, etc
func (p *Planet) initStartingWorld(player *Player, rules *Rules, startingPlanet StartingPlanet, concentration Mineral, surface Mineral) error {

	if len(player.Race.Spec.StartingPlanets) == 0 {
		return fmt.Errorf("no starting planets defined for player %v, race %v", player, player.Race)
	}

	log.Debug().Msgf("Assigning %s to %s as homeworld", p, player)

	p.Homeworld = true
	p.PlayerNum = player.Num

	habWidth := player.Race.HabWidth()
	habCenter := player.Race.HabCenter()
	p.Hab = Hab{
		Grav: habCenter.Grav + int(float64((habWidth.Grav-rules.random.Intn(habWidth.Grav-1)))/2*startingPlanet.HabPenaltyFactor),
		Temp: habCenter.Temp + int(float64((habWidth.Temp-rules.random.Intn(habWidth.Temp-1)))/2*startingPlanet.HabPenaltyFactor),
		Rad:  habCenter.Rad + int(float64((habWidth.Rad-rules.random.Intn(habWidth.Rad-1)))/2*startingPlanet.HabPenaltyFactor),
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
		p.Mines = p.innateMines(player)
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

	// // the homeworld gets a starbase
	starbaseDesign := player.GetDesign(startingPlanet.StarbaseDesignName)
	starbase := newStarbase(player, p, starbaseDesign, starbaseDesign.Name)
	starbase.Spec = computeFleetSpec(rules, player, &starbase)
	p.starbase = &starbase

	// p.PacketSpeed = p.Starbase.Spec.SafePacketSpeed

	// // apply the default plan, but remove the terraforming item because our homeworld is perfect
	// planetService.ApplyProductionPlan(planet.ProductionQueue.Items, player, player.ProductionPlans[0]);
	// planet.ProductionQueue.Items = planet.ProductionQueue.Items.Where(item => !item.IsTerraform).ToList();

	// p.ProductionQueue = append(p.ProductionQueue, ProductionQueueItem{Type: QueueItemTypeAutoMinTerraform, Quantity: 1})
	p.ProductionQueue = append(p.ProductionQueue, ProductionQueueItem{Type: QueueItemTypeAutoFactories, Quantity: 10})
	p.ProductionQueue = append(p.ProductionQueue, ProductionQueueItem{Type: QueueItemTypeAutoMines, Quantity: 10})

	messager.homePlanet(player, p)

	return nil
}

// Get the number of innate mines this player would have on this planet
func (p *Planet) innateMines(player *Player) int {
	if player.Race.Spec.InnateMining {
		return int(math.Sqrt(float64(p.population()) * float64(.1)))
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
func (p *Planet) getGrowthAmount(player *Player, maxPopulation int) int {
	race := &player.Race
	growthFactor := race.Spec.GrowthFactor
	capacity := float64(p.population()) / float64(maxPopulation)
	habValue := race.GetPlanetHabitability(p.Hab)
	if habValue > 0 {
		popGrowth := int(float64(p.population())*float64(race.GrowthRate)*growthFactor/100.0*float64(habValue)/100.0 + .5)

		if capacity > .25 {
			crowdingFactor := 16.0 / 9.0 * (1.0 - capacity) * (1.0 - capacity)
			popGrowth = int(float64(popGrowth) * crowdingFactor)
		}

		// round to the nearest 100 colonists
		return roundToNearest100(popGrowth)
	} else {
		// kill off (habValue / 10)% colonists every year. I.e. a habValue of -4% kills off .4%
		deathAmount := int(float64(p.population()) * (float64(habValue) / 1000.0))
		return roundToNearest100(clamp(deathAmount, deathAmount, -100))
	}
}

func computePlanetSpec(rules *Rules, player *Player, planet *Planet) PlanetSpec {
	spec := PlanetSpec{}
	race := &player.Race
	spec.Habitability = race.GetPlanetHabitability(planet.Hab)
	spec.MaxPopulation = getMaxPopulation(rules, spec.Habitability, player)
	spec.PopulationDensity = float64(planet.population()) / float64(spec.MaxPopulation)
	spec.GrowthAmount = planet.getGrowthAmount(player, spec.MaxPopulation)
	spec.MineralOutput = planet.getMineralOutput(planet.Mines, race.MineOutput)

	if !race.Spec.InnateMining {
		spec.MaxMines = planet.population() * race.NumMines / 10000
		spec.MaxPossibleMines = spec.MaxPopulation * race.NumMines / 10000
	}

	if race.Spec.InnateResources {
		spec.ResourcesPerYear = int(math.Sqrt(float64(planet.population()) * float64(player.TechLevels.Energy) / float64(race.PopEfficiency)))
	} else {
		// compute resources from population
		resourcesFromPop := planet.population() / (race.PopEfficiency * 100)

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
		spec.DefenseCoverage = float64(1.0 - (math.Pow((1 - (player.Spec.Defense.DefenseCoverage / 100)), float64(clamp(planet.Defenses, 0, spec.MaxDefenses)))))
		spec.DefenseCoverageSmart = float64(1.0 - (math.Pow((1 - (player.Spec.Defense.DefenseCoverage / 100 * rules.SmartDefenseCoverageFactor)), float64(clamp(planet.Defenses, 0, spec.MaxDefenses)))))
	}

	if planet.Scanner {
		scanner := player.Spec.PlanetaryScanner
		spec.Scanner = scanner.Name
		spec.ScanRange = scanner.ScanRange
		spec.ScanRangePen = scanner.ScanRangePen
	}

	starbase := planet.starbase
	spec.HasStarbase = starbase != nil
	if starbase != nil && starbase.Spec.HasStargate {
		spec.HasStargate = starbase != nil && starbase.Spec.HasStargate
		spec.SafeHullMass = starbase.Spec.SafeHullMass
		spec.SafeRange = starbase.Spec.SafeRange
		spec.MaxHullMass = starbase.Spec.MaxHullMass
		spec.MaxRange = starbase.Spec.MaxRange

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

func getMaxPopulation(rules *Rules, hab int, player *Player) int {
	race := &player.Race
	maxPopulationFactor := 1 + race.Spec.MaxPopulationOffset

	return roundToNearest100f(float64(rules.MaxPopulation) * maxPopulationFactor * float64(hab) / 100.0)
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
