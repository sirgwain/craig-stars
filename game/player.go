package game

import (
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
)

type Player struct {
	ID                    uint                           `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt             time.Time                      `json:"createdAt,omitempty"`
	UpdatedAt             time.Time                      `json:"updatedat,omitempty"`
	GameID                uint                           `json:"gameId,omitempty"`
	UserID                uint                           `json:"userId,omitempty"`
	Name                  string                         `json:"name,omitempty"`
	Num                   int                            `json:"num"`
	Ready                 bool                           `json:"ready,omitempty"`
	AIControlled          bool                           `json:"aIControlled,omitempty"`
	SubmittedTurn         bool                           `json:"submittedTurn,omitempty"`
	Color                 string                         `json:"color,omitempty"`
	DefaultHullSet        int                            `json:"defaultHullSet,omitempty"`
	Race                  Race                           `json:"race,omitempty"`
	TechLevels            TechLevel                      `json:"techLevels,omitempty" gorm:"embedded;embeddedPrefix:tech_levels_"`
	TechLevelsSpent       TechLevel                      `json:"techLevelsSpent,omitempty" gorm:"embedded;embeddedPrefix:tech_levels_spent"`
	ResearchAmount        int                            `json:"researchAmount,omitempty"`
	ResearchSpentLastYear int                            `json:"researchSpentLastYear,omitempty"`
	NextResearchField     NextResearchField              `json:"nextResearchField,omitempty"`
	Researching           TechField                      `json:"researching,omitempty"`
	BattlePlans           []BattlePlan                   `json:"battlePlans,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	ProductionPlans       []ProductionPlan               `json:"productionPlans,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	TransportPlans        []TransportPlan                `json:"transportPlans,omitempty" gorm:"serializer:json"`
	Messages              []PlayerMessage                `json:"messages,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	Designs               []*ShipDesign                  `json:"designs" gorm:"foreignKey:PlayerID;references:ID"`
	Fleets                []*Fleet                       `json:"fleets" gorm:"foreignKey:PlayerID;references:ID"`
	Planets               []*Planet                      `json:"planets" gorm:"foreignKey:PlayerID;references:ID"`
	MineralPackets        []*MineralPacket               `json:"mineralPackets" gorm:"foreignKey:PlayerID;references:ID"`
	PlanetIntels          []PlanetIntel                  `json:"planetIntels" gorm:"constraint:OnDelete:CASCADE;"`
	FleetIntels           []FleetIntel                   `json:"fleetIntels" gorm:"constraint:OnDelete:CASCADE;"`
	DesignIntels          []ShipDesignIntel              `json:"designIntels" gorm:"constraint:OnDelete:CASCADE;"`
	MineralPacketIntels   []MineralPacketIntel           `json:"mineralPacketIntels" gorm:"constraint:OnDelete:CASCADE;"`
	Spec                  *PlayerSpec                    `json:"spec,omitempty" gorm:"serializer:json"`
	Stats                 *PlayerStats                   `json:"stats,omitempty" gorm:"serializer:json"`
	PlanetsByNum          map[int]*Planet                `json:"-" gorm:"-"`
	FleetsByNum           map[int]*Fleet                 `json:"-" gorm:"-"`
	FleetIntelsByKey      map[string]*FleetIntel         `json:"-" gorm:"-"`
	DesignIntelsByKey     map[uuid.UUID]*ShipDesignIntel `json:"-" gorm:"-"`
	LeftoverResources     int                            `json:"-" gorm:"-"`
}

type PlayerStats struct {
	FleetsBuilt      int `json:"fleetsBuilt,omitempty"`
	TokensBuilt      int `json:"tokensBuilt,omitempty"`
	PlanetsColonized int `json:"planetsColonized,omitempty"`
}

type NextResearchField string

const (
	NextResearchFieldSameField     NextResearchField = "SameField"
	NextResearchFieldEnergy        NextResearchField = "Energy"
	NextResearchFieldWeapons       NextResearchField = "Weapons"
	NextResearchFieldPropulsion    NextResearchField = "Propulsion"
	NextResearchFieldConstruction  NextResearchField = "Construction"
	NextResearchFieldElectronics   NextResearchField = "Electronics"
	NextResearchFieldBiotechnology NextResearchField = "Biotechnology"
	NextResearchFieldLowestField   NextResearchField = "LowestField"
)

type PlayerSpec struct {
	PlanetaryScanner  *TechPlanetaryScanner `json:"planetaryScanner"`
	Defense           *TechDefense          `json:"defense"`
	ResourcesLeftover int                   `json:"resourcesAvailable,omitempty"`
}

type BattlePlan struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	PlayerID        uint             `json:"playerId"`
	Name            string           `json:"name"`
	PrimaryTarget   BattleTargetType `json:"primaryTarget"`
	SecondaryTarget BattleTargetType `json:"secondaryTarget"`
	Tactic          BattleTactic     `json:"tactic"`
	AttackWho       BattleAttackWho  `json:"attackWho"`
}

type BattleTargetType string

const (
	BattleTargetNone              BattleTargetType = "None"
	BattleTargetAny               BattleTargetType = "Any"
	BattleTargetStarbase          BattleTargetType = "Starbase"
	BattleTargetArmedShips        BattleTargetType = "ArmedShips"
	BattleTargetBombersFreighters BattleTargetType = "BombersFreighters"
	BattleTargetUnarmedShips      BattleTargetType = "UnarmedShips"
	BattleTargetFuelTransports    BattleTargetType = "FuelTransports"
	BattleTargetFreighters        BattleTargetType = "Freighters"
)

type BattleTactic string

const (
	// RUN AWAY!
	BattleTacticDisengage BattleTactic = "Disengage"
	// MaximizeDamage until we are damaged, then disengage
	BattleTacticDisengageIfChallenged BattleTactic = "DisengageIfChallenged"
	// If in range of enemy weapons, move away. Only fire if cornered or if from a safe range
	BattleTacticMinimizeDamageToSelf BattleTactic = "MinimizeDamageToSelf"
	BattleTacticMaximizeNetDamage    BattleTactic = "MaximizeNetDamage"
	BattleTacticMaximizeDamageRatio  BattleTactic = "MaximizeDamageRatio"
	BattleTacticMaximizeDamage       BattleTactic = "MaximizeDamage"
)

type BattleAttackWho string

const (
	BattleAttackWhoEnemies            BattleAttackWho = "Enemies"
	BattleAttackWhoEnemiesAndNeutrals BattleAttackWho = "EnemiesAndNeutrals"
	BattleAttackWhoEveryone           BattleAttackWho = "Everyone"
)

type TransportPlan struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	PlayerID uint                   `json:"playerId"`
	Name     string                 `json:"name"`
	Tasks    WaypointTransportTasks `json:"tasks,omitempty" gorm:"serializer:json"`
}

type ProductionPlan struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	PlayerID                          uint                 `json:"playerId"`
	Name                              string               `json:"name"`
	Items                             []ProductionPlanItem `json:"items" gorm:"serializer:json"`
	ContributesOnlyLeftoverToResearch bool                 `json:"contributesOnlyLeftoverToResearch,omitempty"`
}

type ProductionPlanItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	ProductionPlanID uint          `json:"-"`
	Type             QueueItemType `json:"type"`
	DesignName       string        `json:"designName"`
	Quantity         int           `json:"quantity"`
	Allocated        Cost          `json:"allocated" gorm:"embedded;embeddedPrefix:allocated_"`
	SortOrder        int           `json:"-"`
}

// create a new player with an existing race. The race
// will be copied for the player
func NewPlayer(userID uint, race *Race) *Player {

	// copy this race for the player
	playerRace := *race
	playerRace.ID = 0
	playerRace.CreatedAt = time.Time{}
	playerRace.UpdatedAt = time.Time{}

	return &Player{
		UserID:            userID,
		Stats:             &PlayerStats{},
		Race:              playerRace,
		Color:             "#0000FF", // default to blue
		ResearchAmount:    15,
		NextResearchField: NextResearchFieldLowestField,
		BattlePlans: []BattlePlan{
			{
				Name:            "Default",
				PrimaryTarget:   BattleTargetArmedShips,
				SecondaryTarget: BattleTargetAny,
				Tactic:          BattleTacticMaximizeDamageRatio,
				AttackWho:       BattleAttackWhoEnemiesAndNeutrals,
			},
		},
		ProductionPlans: []ProductionPlan{
			{Name: "Default", Items: []ProductionPlanItem{
				{Type: QueueItemTypeAutoMinTerraform, Quantity: 1},
				{Type: QueueItemTypeAutoFactories, Quantity: 10},
				{Type: QueueItemTypeAutoMines, Quantity: 10},
			}},
		},
		TransportPlans: []TransportPlan{
			{Name: "Default"},
			{Name: "Quick Load", Tasks: WaypointTransportTasks{
				Fuel:      WaypointTransportTask{Action: TransportActionLoadOptimal},
				Ironium:   WaypointTransportTask{Action: TransportActionLoadAll},
				Boranium:  WaypointTransportTask{Action: TransportActionLoadAll},
				Germanium: WaypointTransportTask{Action: TransportActionLoadAll},
			}},
			{Name: "Quick Drop", Tasks: WaypointTransportTasks{
				Fuel:      WaypointTransportTask{Action: TransportActionLoadOptimal},
				Ironium:   WaypointTransportTask{Action: TransportActionUnloadAll},
				Boranium:  WaypointTransportTask{Action: TransportActionUnloadAll},
				Germanium: WaypointTransportTask{Action: TransportActionUnloadAll},
			}},
			{Name: "Load Colonists", Tasks: WaypointTransportTasks{
				Colonists: WaypointTransportTask{Action: TransportActionLoadAll},
			}},
			{Name: "Unload Colonists", Tasks: WaypointTransportTasks{
				Colonists: WaypointTransportTask{Action: TransportActionUnloadAll},
			}},
		},
	}
}

func (p *Player) WithTechLevels(tl TechLevel) *Player {
	p.TechLevels = tl
	return p
}

func (p *Player) WithSpec(rules *Rules) *Player {
	p.Spec = computePlayerSpec(p, rules)
	return p
}

func (p *Player) String() string {
	return fmt.Sprintf("Player %d (%d) %s", p.Num, p.ID, p.Race.PluralName)
}

func (p *Player) getNextFleetNum() int {
	num := 1

	orderedFleets := make([]*Fleet, len(p.Fleets))
	copy(orderedFleets, p.Fleets)
	sort.Slice(orderedFleets, func(i, j int) bool { return orderedFleets[i].Num < orderedFleets[j].Num })

	for i := 0; i < len(orderedFleets); i++ {
		// todo figure out starbasees
		fleet := orderedFleets[i]
		if i > 0 {
			// if we are past fleet #1 and we skipped a number, used the skipped number
			if fleet.Num > 1 && fleet.Num != orderedFleets[i-1].Num+1 {
				return orderedFleets[i-1].Num + 1
			}
		}
		// we are the next num...
		num = fleet.Num + 1
	}

	return num
}

// get a player owned planet by num, or nil if it doesn't exist
func (p *Player) GetPlanet(num int) *Planet {
	return p.PlanetsByNum[num]
}

func (p *Player) AddPlanet(planet *Planet) {
	p.Planets = append(p.Planets, planet)
	p.PlanetsByNum[planet.Num] = planet
}

// get a player owned planet by num, or nil if it doesn't exist
func (p *Player) GetFleet(num int) *Fleet {
	return p.FleetsByNum[num]
}

func (p *Player) AddFleet(fleet *Fleet) {
	p.Fleets = append(p.Fleets, fleet)
	p.FleetsByNum[fleet.Num] = fleet
}

func (p *Player) RemoveFleet(fleet *Fleet) {
	delete(p.FleetsByNum, fleet.Num)

	fleets := make([]*Fleet, 0, len(p.Fleets)-1)
	for _, f := range p.Fleets {
		if f.Num != fleet.Num {
			fleets = append(fleets, f)
		}
	}
	p.Fleets = fleets
}

// Get a player ShipDesign by name, or nil if no design found
func (p *Player) GetDesign(name string) *ShipDesign {
	for i := range p.Designs {
		design := p.Designs[i]
		if design.Name == name {
			return design
		}
	}
	return nil
}

// get the latest ship design by purpose
func (p *Player) GetLatestDesign(purpose ShipDesignPurpose) *ShipDesign {
	var latest *ShipDesign = nil
	for i := range p.Designs {
		design := p.Designs[i]
		if design.Purpose == purpose {
			if latest == nil {
				latest = design
			} else {
				if latest.Version < design.Version {
					latest = design
				}
			}

		}
	}

	return latest
}

// get all planets the player owns that can build ships of mass mass
func (p *Player) GetBuildablePlanets(mass int) []*Planet {
	planets := []*Planet{}
	for _, planet := range p.Planets {
		if planet.CanBuild(mass) {
			planets = append(planets, planet)
		}
	}
	return planets
}

func computePlayerSpec(player *Player, rules *Rules) *PlayerSpec {
	techs := rules.Techs
	spec := PlayerSpec{
		PlanetaryScanner:  techs.GetBestPlanetaryScanner(player),
		Defense:           techs.GetBestDefense(player),
		ResourcesLeftover: 0,
	}

	return &spec
}

// return true if the player currently has this tech
func (p *Player) HasTech(tech *Tech) bool {
	return p.CanLearnTech(tech) && p.TechLevels.HasRequiredLevels(tech.Requirements.TechLevel)
}

func (p *Player) CanLearnTech(tech *Tech) bool {
	requirements := tech.Requirements
	if requirements.PRTRequired != PRTNone && requirements.PRTRequired != p.Race.PRT {
		return false
	}
	if requirements.PRTDenied != PRTNone && p.Race.PRT == requirements.PRTDenied {
		return false
	}

	if requirements.LRTsRequired != 0 && p.Race.LRTs&Bitmask(requirements.LRTsRequired) != Bitmask(requirements.LRTsRequired) {
		return false
	}

	if requirements.LRTsDenied != 0 && p.Race.LRTs&Bitmask(requirements.LRTsDenied) != 0 {
		return false
	}

	return true
}

func (p *Player) clearTransientReports() {
	p.FleetIntels = []FleetIntel{}
	p.FleetIntelsByKey = map[string]*FleetIntel{}
}

// build maps used for quick lookups for various player objects
func (p *Player) BuildMaps() {
	p.FleetIntelsByKey = map[string]*FleetIntel{}
	for i := range p.FleetIntels {
		intel := &p.FleetIntels[i]
		p.FleetIntelsByKey[intel.String()] = intel
	}

	p.DesignIntelsByKey = map[uuid.UUID]*ShipDesignIntel{}
	for i := range p.DesignIntels {
		intel := &p.DesignIntels[i]
		p.DesignIntelsByKey[intel.UUID] = intel
	}

	p.PlanetsByNum = make(map[int]*Planet, len(p.Planets))
	for _, planet := range p.Planets {
		p.PlanetsByNum[planet.Num] = planet
	}

	p.FleetsByNum = make(map[int]*Fleet, len(p.Fleets))
	for _, fleet := range p.Fleets {
		p.FleetsByNum[fleet.Num] = fleet
	}

}
