package game

import (
	"fmt"
	"math"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Planet struct {
	MapObject
	Hab                               Hab                   `json:"hab,omitempty" gorm:"embedded;embeddedPrefix:hab_"`
	BaseHab                           Hab                   `json:"baseHab,omitempty" gorm:"embedded;embeddedPrefix:base_hab_"`
	TerraformedAmount                 Hab                   `json:"terraformedAmount,omitempty" gorm:"embedded;embeddedPrefix:terraform_hab_"`
	MineralConcentration              Mineral               `json:"mineralConcentration,omitempty" gorm:"embedded;embeddedPrefix:mineral_conc_"`
	MineYears                         Mineral               `json:"mineYears,omitempty" gorm:"embedded;embeddedPrefix:mine_years_"`
	Starbase                          *Fleet                `json:"starbase,omitempty"`
	Cargo                             Cargo                 `json:"cargo,omitempty" gorm:"embedded;embeddedPrefix:cargo_"`
	ProductionQueue                   []ProductionQueueItem `json:"productionQueue,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	Mines                             int                   `json:"mines,omitempty"`
	Factories                         int                   `json:"factories,omitempty"`
	Defenses                          int                   `json:"defenses,omitempty"`
	Homeworld                         bool                  `json:"homeworld,omitempty"`
	ContributesOnlyLeftoverToResearch bool                  `json:"contributesOnlyLeftoverToResearch,omitempty"`
	Scanner                           bool                  `json:"scanner,omitempty"`
	Spec                              *PlanetSpec           `json:"spec,omitempty" gorm:"serializer:json"`
}

type ProductionQueueItem struct {
	ID        uint           `gorm:"primaryKey" json:"id" header:"Username"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedat"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	PlanetID  uint           `json:"-"`
	Type      QueueItemType  `json:"type"`
	Quantity  int            `json:"quantity"`
	SortOrder int            `json:"-"`
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
	dockCapacity              int     `json:"dockCapacity,omitempty"`
}

func NewPlanet(gameID uint) Planet {
	return Planet{MapObject: MapObject{GameID: gameID, Dirty: true}}
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

func (p *Planet) Empty() {
	p.Hab = Hab{}
	p.BaseHab = Hab{}
	p.TerraformedAmount = Hab{}
	p.MineralConcentration = Mineral{}
	p.Cargo = Cargo{}
	p.ProductionQueue = []ProductionQueueItem{}
	p.MineYears = Mineral{}
}

// Randomize a planet with new hab range, minerals, etc
func (p *Planet) Randomize(rules *Rules) {
	p.Empty()

	// From @SuicideJunkie's tests and @edmundmk's previous research, grav and temp are weighted slightly towards
	// the center, rad is completely random
	// @edmundmk:
	// "I'm certain gravity and temperature probability is constant between 10 and 90 inclusive, and falls off towards 0 and 100.
	// It never generates 0 or 100 so I have to change my random formula to (1 to 90)+(0 to 9)
	// damn you all for sucking me into stars! again lol"
	p.Hab = Hab{
		Grav: rules.Random.Intn(91) + rules.Random.Intn(10),
		Temp: rules.Random.Intn(91) + rules.Random.Intn(10),
		Rad:  1 + rules.Random.Intn(100),
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
		Ironium:   rules.MinStartingMineralConcentration + rules.Random.Intn(rules.MaxStartingMineralConcentration+1),
		Boranium:  rules.MinStartingMineralConcentration + rules.Random.Intn(rules.MaxStartingMineralConcentration+1),
		Germanium: rules.MinStartingMineralConcentration + rules.Random.Intn(rules.MaxStartingMineralConcentration+1),
	}

	if p.MineralConcentration.Ironium > 30 {
		p.MineralConcentration.Ironium = 30 + rules.Random.Intn(45) + rules.Random.Intn(45)
	}

	if p.MineralConcentration.Boranium > 30 {
		p.MineralConcentration.Boranium = 30 + rules.Random.Intn(45) + rules.Random.Intn(45)
	}

	if p.MineralConcentration.Germanium > 30 {
		p.MineralConcentration.Germanium = 30 + rules.Random.Intn(45) + germRadBonus + rules.Random.Intn(45)
	}
}

// Initialize a planet to be a homeworld for a payer with ideal hab, starting mineral concentration, etc
func (p *Planet) initHomeworld(player *Player, rules *Rules, concentration Mineral, surface Mineral) error {

	if player.Race.Spec == nil || len(player.Race.Spec.StartingPlanets) == 0 {
		return fmt.Errorf("no starting planets defined for player %v, race %v", player, player.Race)
	}

	log.Debug().Msgf("Assigning %s to %s as homeworld", p, player)

	startingPlanet := player.Race.Spec.StartingPlanets[0]

	p.Homeworld = true
	p.PlayerNum = &player.Num
	p.PlayerID = player.ID
	p.Hab = player.Race.HabCenter()
	p.MineralConcentration = concentration
	p.Cargo = surface.ToCargo()

	// reset some fields in case this is called on an existing planet for some reason
	p.ProductionQueue = []ProductionQueueItem{}
	p.BaseHab = Hab{}
	p.TerraformedAmount = Hab{}

	raceSpec := player.Race.Spec

	// set the homeworld pop to our starting planet pop
	p.SetPopulation(startingPlanet.Population)
	p.Population()

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

	// TODO
	// // the homeworld gets a starbase
	// var design = Game.DesignsByGuid[player.GetLatestDesign(ShipDesignPurpose.Starbase).Guid];
	// CreateStarbaseOnPlanet(player, planet, design);

	// // apply the default plan, but remove the terraforming item because our homeworld is perfect
	// planetService.ApplyProductionPlan(planet.ProductionQueue.Items, player, player.ProductionPlans[0]);
	// planet.ProductionQueue.Items = planet.ProductionQueue.Items.Where(item => !item.IsTerraform).ToList();

	p.ProductionQueue = append(p.ProductionQueue, ProductionQueueItem{Type: QueueItemTypeAutoMinTerraform, Quantity: 1})
	p.ProductionQueue = append(p.ProductionQueue, ProductionQueueItem{Type: QueueItemTypeAutoFactories, Quantity: 10})
	p.ProductionQueue = append(p.ProductionQueue, ProductionQueueItem{Type: QueueItemTypeAutoMines, Quantity: 10})
	for i := range p.ProductionQueue {
		p.ProductionQueue[i].SortOrder = i
	}

	// Message.HomePlanet(player, planet);

	return nil
}

// Get the number of innate mines this player would have on this planet
func (p *Planet) GetInnateMines(player *Player) int {
	if player.Race.Spec.InnateMining {
		return int(math.Sqrt(float64(p.Population()) * float64(.1)))
	}
	return 0
}

func (p *Planet) shortestDistanceToPlanets(otherPlanets *[]Planet) float64 {
	minDistanceSquared := math.MaxFloat64
	for _, planet := range *otherPlanets {
		distSquared := p.Position.DistanceSquaredTo(&planet.Position)
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

func computePlanetSpec(rules *Rules, planet *Planet, player *Player) *PlanetSpec {
	spec := PlanetSpec{}
	spec.MaxPopulation = getMaxPopulation(rules, planet, player)
	spec.GrowthAmount = planet.getGrowthAmount(player, spec.MaxPopulation)
	spec.MineralOutput = planet.getMineralOutput(planet.Mines, player.Race.MineOutput)

	if planet.Scanner {
		scanner := rules.Techs.GetBestPlanetaryScanner(player)
		spec.Scanner = scanner.Name
		spec.ScanRange = scanner.ScanRange
		spec.ScanRangePen = scanner.ScanRangePen
	}
	return &spec
}

func getMaxPopulation(rules *Rules, planet *Planet, player *Player) int {
	race := &player.Race
	maxPopulationFactor := 1 + race.Spec.MaxPopulationOffset

	// get this player's planet habitability
	hab := race.GetPlanetHabitability(planet.Hab)
	return int(float64(rules.MaxPopulation) * maxPopulationFactor * float64(hab) / 100.0)
}
