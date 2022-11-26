package game

import (
	"fmt"
	"math"

	"github.com/rs/zerolog/log"
)

type Planet struct {
	MapObject
	Hab                               Hab                   `json:"hab,omitempty" gorm:"embedded;embeddedPrefix:hab_"`
	BaseHab                           Hab                   `json:"baseHab,omitempty" gorm:"embedded;embeddedPrefix:base_hab_"`
	TerraformedAmount                 Hab                   `json:"terraformedAmount,omitempty" gorm:"embedded;embeddedPrefix:terraform_hab_"`
	MineralConcentration              Mineral               `json:"mineralConcentration,omitempty" gorm:"embedded;embeddedPrefix:mineral_conc_"`
	MineYears                         Mineral               `json:"mineYears,omitempty" gorm:"embedded;embeddedPrefix:mine_years_"`
	Cargo                             Cargo                 `json:"cargo,omitempty" gorm:"embedded;embeddedPrefix:cargo_"`
	Mines                             int                   `json:"mines,omitempty"`
	Factories                         int                   `json:"factories,omitempty"`
	Defenses                          int                   `json:"defenses,omitempty"`
	Homeworld                         bool                  `json:"homeworld,omitempty"`
	ContributesOnlyLeftoverToResearch bool                  `json:"contributesOnlyLeftoverToResearch,omitempty"`
	Scanner                           bool                  `json:"scanner,omitempty"`
	PacketSpeed                       int                   `json:"packetSpeed,omitempty"`
	BonusResources                    int                   `json:"-" gorm:"-"`
	ProductionQueue                   []ProductionQueueItem `json:"productionQueue,omitempty" gorm:"serializer:json"`
	Spec                              *PlanetSpec           `json:"spec,omitempty" gorm:"serializer:json"`
	Starbase                          *Fleet                `json:"starbase,omitempty"`
}

type ProductionQueueItem struct {
	Type       QueueItemType `json:"type"`
	DesignName string        `json:"designName"`
	Quantity   int           `json:"quantity"`
	Allocated  Cost          `json:"allocated"`
}

type QueueItemType string

const (
	QueueItemTypeIroniumMineralPacket   QueueItemType = "IroniumMineralPacket"
	QueueItemTypeBoraniumMineralPacket  QueueItemType = "BoraniumMineralPacket"
	QueueItemTypeGermaniumMineralPacket QueueItemType = "GermaniumMineralPacket"
	QueueItemTypeMixedMineralPacket     QueueItemType = "MixedMineralPacket"
	QueueItemTypeFactory                QueueItemType = "Factory"
	QueueItemTypeMine                   QueueItemType = "Mine"
	QueueItemTypeDefenses               QueueItemType = "Defenses"
	QueueItemTypeMineralAlchemy         QueueItemType = "MineralAlchemy"
	QueueItemTypeTerraformEnvironment   QueueItemType = "TerraformEnvironment"
	QueueItemTypeAutoMines              QueueItemType = "AutoMines"
	QueueItemTypeAutoFactories          QueueItemType = "AutoFactories"
	QueueItemTypeAutoDefenses           QueueItemType = "AutoDefenses"
	QueueItemTypeAutoMineralAlchemy     QueueItemType = "AutoMineralAlchemy"
	QueueItemTypeAutoMinTerraform       QueueItemType = "AutoMinTerraform"
	QueueItemTypeAutoMaxTerraform       QueueItemType = "AutoMaxTerraform"
	QueueItemTypeAutoMineralPacket      QueueItemType = "AutoMineralPacket"
	QueueItemTypeShipToken              QueueItemType = "ShipToken"
	QueueItemTypeStarbase               QueueItemType = "Starbase"
)

type PlanetSpec struct {
	Habitability              int     `json:"habitability,omitempty"`
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
	DockCapacity              int     `json:"dockCapacity,omitempty"`
}

func (item *ProductionQueueItem) String() string {
	return fmt.Sprintf("ProductionQueueItem %d %s (%s)", item.Quantity, item.Type, item.DesignName)
}

func NewPlanet() *Planet {
	return &Planet{MapObject: MapObject{Type: MapObjectTypePlanet, Dirty: true, PlayerNum: Unowned}}
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

func (p *Planet) Population() int {
	return p.Cargo.Colonists * 100
}

func (p *Planet) SetPopulation(pop int) {
	p.Cargo.Colonists = pop / 100
}

// true if this planet can build a ship with a given mass
func (p *Planet) CanBuild(mass int) bool {
	return p.Spec.HasStarbase && (p.Starbase.Spec.SpaceDock == UnlimitedSpaceDock || p.Starbase.Spec.SpaceDock >= mass)
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

	if player.Race.Spec == nil || len(player.Race.Spec.StartingPlanets) == 0 {
		return fmt.Errorf("no starting planets defined for player %v, race %v", player, player.Race)
	}

	log.Debug().Msgf("Assigning %s to %s as homeworld", p, player)

	p.Homeworld = true
	p.PlayerNum = player.Num
	p.PlayerID = player.ID

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
	p.SetPopulation(int(float64(startingPlanet.Population) * raceSpec.StartingPopulationFactor))

	if raceSpec.InnateMining {
		p.Mines = p.GetInnateMines(player)
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

	p.ContributesOnlyLeftoverToResearch = true
	p.Scanner = true

	// // the homeworld gets a starbase
	starbaseDesign := player.GetDesign(startingPlanet.StarbaseDesignName)
	starbase := NewStarbase(player, p, starbaseDesign, starbaseDesign.Name)
	starbase.Spec = ComputeFleetSpec(rules, player, &starbase)
	p.Starbase = &starbase

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
func (p *Planet) GetInnateMines(player *Player) int {
	if player.Race.Spec.InnateMining {
		return int(math.Sqrt(float64(p.Population()) * float64(.1)))
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
	capacity := float64(p.Population()) / float64(maxPopulation)
	habValue := race.GetPlanetHabitability(p.Hab)
	if habValue > 0 {
		popGrowth := int(float64(p.Population())*float64(race.GrowthRate)*growthFactor/100.0*float64(habValue)/100.0 + .5)

		if capacity > .25 {
			crowdingFactor := 16.0 / 9.0 * (1.0 - capacity) * (1.0 - capacity)
			popGrowth = int(float64(popGrowth) * crowdingFactor)
		}

		// round to the nearest 100 colonists
		return roundToNearest100(popGrowth)
	} else {
		// kill off (habValue / 10)% colonists every year. I.e. a habValue of -4% kills off .4%
		deathAmount := int(float64(p.Population()) * (float64(habValue) / 1000.0))
		return roundToNearest100(clamp(deathAmount, deathAmount, -100))
	}
}

func ComputePlanetSpec(rules *Rules, planet *Planet, player *Player) *PlanetSpec {
	spec := PlanetSpec{}
	race := &player.Race
	spec.Habitability = race.GetPlanetHabitability(planet.Hab)
	spec.MaxPopulation = getMaxPopulation(rules, spec.Habitability, player)
	spec.PopulationDensity = float64(planet.Population()) / float64(spec.MaxPopulation)
	spec.GrowthAmount = planet.getGrowthAmount(player, spec.MaxPopulation)
	spec.MineralOutput = planet.getMineralOutput(planet.Mines, race.MineOutput)

	if !race.Spec.InnateMining {
		spec.MaxMines = planet.Population() * race.NumMines / 10000
		spec.MaxPossibleMines = spec.MaxPopulation * race.NumMines / 10000
	}

	if race.Spec.InnateResources {
		spec.ResourcesPerYear = int(math.Sqrt(float64(planet.Population()) * float64(player.TechLevels.Energy) / float64(race.PopEfficiency)))
	} else {
		// compute resources from population
		resourcesFromPop := planet.Population() / (race.PopEfficiency * 100)

		// compute resources from factories
		resourcesFromFactories := planet.Factories * race.FactoryOutput / 10

		spec.ResourcesPerYear = resourcesFromPop + resourcesFromFactories
		spec.MaxFactories = planet.Population() * race.NumFactories / 10000
		spec.MaxPossibleFactories = spec.MaxPopulation * race.NumFactories / 10000
	}

	if planet.ContributesOnlyLeftoverToResearch {
		spec.ResourcesPerYearAvailable = spec.ResourcesPerYear
	} else {
		spec.ResourcesPerYearResearch = int(float64(spec.ResourcesPerYear) * float64(player.ResearchAmount) / 100.0)
		spec.ResourcesPerYearAvailable = spec.ResourcesPerYear - spec.ResourcesPerYearResearch
	}

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

	spec.HasStarbase = planet.Starbase != nil

	return &spec
}

func getMaxPopulation(rules *Rules, hab int, player *Player) int {
	race := &player.Race
	maxPopulationFactor := 1 + race.Spec.MaxPopulationOffset

	return roundToNearest100f(float64(rules.MaxPopulation) * maxPopulationFactor * float64(hab) / 100.0)
}

// true if this is an auto type
func (t QueueItemType) IsAuto() bool {
	return t == QueueItemTypeAutoMines ||
		t == QueueItemTypeAutoFactories ||
		t == QueueItemTypeAutoDefenses ||
		t == QueueItemTypeAutoMineralAlchemy ||
		t == QueueItemTypeAutoMinTerraform ||
		t == QueueItemTypeAutoMaxTerraform ||
		t == QueueItemTypeAutoMineralPacket
}

// return the concrete version of this auto type
func (t QueueItemType) ConcreteType() QueueItemType {
	switch t {
	case QueueItemTypeAutoMines:
		return QueueItemTypeMine
	case QueueItemTypeAutoFactories:
		return QueueItemTypeFactory
	case QueueItemTypeAutoDefenses:
		return QueueItemTypeDefenses
	case QueueItemTypeAutoMaxTerraform:
		return QueueItemTypeTerraformEnvironment
	case QueueItemTypeAutoMinTerraform:
		return QueueItemTypeTerraformEnvironment
	case QueueItemTypeAutoMineralAlchemy:
		return QueueItemTypeMineralAlchemy
	case QueueItemTypeAutoMineralPacket:
		return QueueItemTypeMixedMineralPacket
	}
	return t
}

type ProductionResult struct {
	tokens   []ShipToken
	packets  []MineralPacket
	starbase *Fleet
}

// produce one turns worth of items from the production queue
func (planet *Planet) produce(player *Player) ProductionResult {
	result := ProductionResult{}
	available := Cost{Resources: planet.Spec.ResourcesPerYearAvailable}.AddCargoMinerals(planet.Cargo)
	newQueue := []ProductionQueueItem{}
	for itemIndex, item := range planet.ProductionQueue {

		// add in anything allocated in previous turns
		available = available.Add(item.Allocated)
		item.Allocated = Cost{}

		// get the cost of the current item
		cost := player.Race.Spec.Costs[item.Type]
		if item.Type == QueueItemTypeStarbase || item.Type == QueueItemTypeShipToken {
			design := player.GetDesign(item.DesignName)
			if design != nil {
				cost = design.Spec.Cost
			} else {
				log.Error().Msgf("player %s has no design named: %s", player, item.DesignName)
			}
		}

		if (cost != Cost{}) {
			// figure out how many we can build
			// and make sure we only build up to the quantity, and we don't build more than the planet supports
			numBuilt := MaxInt(0, available.NumBuildable(cost))
			numBuilt = MinInt(numBuilt, item.Quantity)
			numBuilt = MinInt(numBuilt, planet.maxBuildable(item.Type))

			if numBuilt > 0 {
				// build the items on the planet and remove from our available
				planet.buildItems(player, item, numBuilt, &result)
				available = available.Minus(cost.MultiplyInt(numBuilt))
			}

			if numBuilt < item.Quantity {
				// allocate to this item
				item.Allocated = planet.allocatePartialBuild(cost, available)
				available = available.Minus(item.Allocated)
			}

			if item.Type.IsAuto() {
				if available.Resources == 0 {
					// we are out of resources, create a partial item end production
					if (item.Allocated != Cost{}) && numBuilt < item.Quantity {
						// we partially built an auto items, create a partial concrete item
						// we have some leftover to allocate so create a concrete item
						concreteItem := ProductionQueueItem{Type: item.Type.ConcreteType(), Quantity: 1, Allocated: item.Allocated}
						item.Allocated = Cost{}

						// add the concreate item to the top of the queue
						newQueue = append([]ProductionQueueItem{concreteItem}, newQueue...)
					}
					// auto items stay in the list
					newQueue = append(newQueue, item)

					if available.Resources == 0 {
						// we are out of resources, so we are done building
						if itemIndex < len(planet.ProductionQueue)-1 {
							// append the unfinished queue back to the end of our remaining items
							newQueue = append(newQueue, planet.ProductionQueue[itemIndex+1:]...)
						}
						break
					}
				} else {
					// auto items stay in the list
					// and we have resources leftover so move on
					newQueue = append(newQueue, item)
				}
			} else {
				item.Quantity -= numBuilt
				if item.Quantity != 0 {
					// we didn't finish, add the item back onto the remaining list
					newQueue = append(newQueue, item)
					if itemIndex < len(planet.ProductionQueue)-1 {
						// append the unfinished queue back to the end of our remaining items
						newQueue = append(newQueue, planet.ProductionQueue[itemIndex+1:]...)
					}
					// we finished, break out
					break
				}
			}
		}
	}
	// replace the queue with what's leftover
	planet.ProductionQueue = newQueue
	player.leftoverResources += available.Resources
	planet.Cargo = Cargo{available.Ironium, available.Boranium, available.Germanium, planet.Cargo.Colonists}

	return result
}

// add built items to planet, build fleets, update player messages, etc
func (planet *Planet) buildItems(player *Player, item ProductionQueueItem, numBuilt int, result *ProductionResult) {

	switch item.Type {
	case QueueItemTypeAutoMines:
		fallthrough
	case QueueItemTypeMine:
		planet.Mines += numBuilt
		messager.minesBuilt(player, planet, numBuilt)
	case QueueItemTypeAutoFactories:
		fallthrough
	case QueueItemTypeFactory:
		planet.Factories += numBuilt
		messager.factoriesBuilt(player, planet, numBuilt)
	case QueueItemTypeAutoDefenses:
		fallthrough
	case QueueItemTypeDefenses:
		planet.Defenses += numBuilt
		messager.defensesBuilt(player, planet, numBuilt)
	case QueueItemTypeShipToken:
		design := player.GetDesign(item.DesignName)
		result.tokens = append(result.tokens, ShipToken{Quantity: numBuilt, design: design, DesignUUID: design.UUID})
	}

	log.Debug().
		Int64("PlayerID", player.ID).
		Int("Player", player.Num).
		Str("Planet", planet.Name).
		Str("Item", string(item.Type)).
		Str("DesignName", item.DesignName).
		Int("NumBuilt", numBuilt).
		Msgf("built item")

}

// Allocate resources to the top item on this production queue
// and return the leftover resources
//
// Costs are allocated by lowest percentage, i.e. if (we require
// Cost(10, 10, 10, 100) and we only have Cost(1, 10, 10, 100)
// we allocate Cost(1, 1, 1, 10)
//
// The min amount we have is 10 percent of the ironium, so we
// apply 10 percent to each cost amount
func (planet *Planet) allocatePartialBuild(costPerItem Cost, allocated Cost) Cost {
	ironiumPerc := 100.0
	if costPerItem.Ironium > 0 {
		ironiumPerc = float64(allocated.Ironium) / float64(costPerItem.Ironium)
	}
	boraniumPerc := 100.0
	if costPerItem.Boranium > 0 {
		boraniumPerc = float64(allocated.Boranium) / float64(costPerItem.Boranium)
	}
	germaniumPerc := 100.0
	if costPerItem.Germanium > 0 {
		germaniumPerc = float64(allocated.Germanium) / float64(costPerItem.Germanium)
	}
	resourcesPerc := 100.0
	if costPerItem.Resources > 0 {
		resourcesPerc = float64(allocated.Resources) / float64(costPerItem.Resources)
	}

	// figure out the lowest percentage
	minPerc := MinFloat64(ironiumPerc, boraniumPerc, germaniumPerc, resourcesPerc)

	// allocate the lowest percentage of each cost
	newAllocated := Cost{
		int(float64(costPerItem.Ironium) * minPerc),
		int(float64(costPerItem.Boranium) * minPerc),
		int(float64(costPerItem.Germanium) * minPerc),
		int(float64(costPerItem.Resources) * minPerc),
	}

	// return the amount we allocate to the top queued item
	return newAllocated
}

// get the maximum buildable amount of a queue item
func (planet *Planet) maxBuildable(t QueueItemType) int {
	switch t {
	case QueueItemTypeAutoMines:
		return planet.Spec.MaxMines - planet.Mines
	case QueueItemTypeMine:
		return planet.Spec.MaxPossibleMines - planet.Mines
	case QueueItemTypeAutoFactories:
		return planet.Spec.MaxFactories - planet.Factories
	case QueueItemTypeFactory:
		return planet.Spec.MaxPossibleFactories - planet.Factories
	case QueueItemTypeAutoDefenses:
		fallthrough
	case QueueItemTypeDefenses:
		return planet.Spec.MaxDefenses - planet.Defenses
	case QueueItemTypeTerraformEnvironment:
	case QueueItemTypeAutoMaxTerraform:
		// return planet.GetTerraformAmount(planet, player).AbsSum());
		break
	case QueueItemTypeAutoMinTerraform:
		// return GetMinTerraformAmount(planet, player).AbsSum());
		break

	case QueueItemTypeStarbase:
		return 1
	}
	// default to infinite
	return math.MaxInt
}

// reduce the mineral concentrations of a planet after mining.
func (planet *Planet) reduceMineralConcentration(rules *Rules) {
	mineralDecayFactor := rules.MineralDecayFactor
	minMineralConcentration := rules.MinMineralConcentration
	if planet.Homeworld {
		minMineralConcentration = rules.MinHomeworldMineralConcentration
	}

	planetMineYears := planet.MineYears.ToSplice()
	planetMineralConcentration := planet.MineralConcentration.ToSplice()
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
