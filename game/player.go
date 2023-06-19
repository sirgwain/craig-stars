package game

import (
	"fmt"
	"time"
)

// a player and all mapobjects the player owns
// this is used by the UI when loading a player's game
type FullPlayer struct {
	Player
	PlayerMapObjects
}

type Player struct {
	ID                    int64                `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt             time.Time            `json:"createdAt,omitempty"`
	UpdatedAt             time.Time            `json:"updatedat,omitempty"`
	GameID                int64                `json:"gameId,omitempty"`
	UserID                int64                `json:"userId,omitempty"`
	Name                  string               `json:"name,omitempty"`
	Num                   int                  `json:"num"`
	Ready                 bool                 `json:"ready,omitempty"`
	AIControlled          bool                 `json:"aIControlled,omitempty"`
	SubmittedTurn         bool                 `json:"submittedTurn,omitempty"`
	Color                 string               `json:"color,omitempty"`
	DefaultHullSet        int                  `json:"defaultHullSet,omitempty"`
	Race                  Race                 `json:"race,omitempty"  gorm:"serializer:json"`
	TechLevels            TechLevel            `json:"techLevels,omitempty" gorm:"embedded;embeddedPrefix:tech_levels_"`
	TechLevelsSpent       TechLevel            `json:"techLevelsSpent,omitempty" gorm:"embedded;embeddedPrefix:tech_levels_spent"`
	ResearchAmount        int                  `json:"researchAmount,omitempty"`
	ResearchSpentLastYear int                  `json:"researchSpentLastYear,omitempty"`
	NextResearchField     NextResearchField    `json:"nextResearchField,omitempty"`
	Researching           TechField            `json:"researching,omitempty"`
	BattlePlans           []BattlePlan         `json:"battlePlans,omitempty" gorm:"serializer:json"`
	ProductionPlans       []ProductionPlan     `json:"productionPlans,omitempty" gorm:"serializer:json"`
	TransportPlans        []TransportPlan      `json:"transportPlans,omitempty" gorm:"serializer:json"`
	Messages              []PlayerMessage      `json:"messages,omitempty"`
	Designs               []ShipDesign         `json:"designs" gorm:"foreignKey:PlayerID;references:ID"`
	PlanetIntels          []PlanetIntel        `json:"planetIntels,omitempty"`
	FleetIntels           []FleetIntel         `json:"fleetIntels,omitempty"`
	DesignIntels          []ShipDesignIntel    `json:"designIntels,omitempty"`
	MineralPacketIntels   []MineralPacketIntel `json:"mineralPacketIntels,omitempty"`
	MineFieldIntels       []MineFieldIntel     `json:"mineFieldIntels,omitempty"`
	Stats                 *PlayerStats         `json:"stats,omitempty" gorm:"serializer:json"`
	Spec                  *PlayerSpec          `json:"spec,omitempty" gorm:"serializer:json"`
	leftoverResources     int                  `json:"-" gorm:"-"`
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
	Name  string                 `json:"name"`
	Tasks WaypointTransportTasks `json:"tasks,omitempty" gorm:"serializer:json"`
}

type ProductionPlan struct {
	Name                              string               `json:"name"`
	Items                             []ProductionPlanItem `json:"items" gorm:"serializer:json"`
	ContributesOnlyLeftoverToResearch bool                 `json:"contributesOnlyLeftoverToResearch,omitempty"`
}

type ProductionPlanItem struct {
	Type       QueueItemType `json:"type"`
	DesignName string        `json:"designName"`
	Quantity   int           `json:"quantity"`
	Allocated  Cost          `json:"allocated" gorm:"embedded;embeddedPrefix:allocated_"`
}

// All mapobjects that a player can issue commands to
type PlayerMapObjects struct {
	Planets        []*Planet        `json:"planets" gorm:"foreignKey:PlayerID;references:ID"`
	Fleets         []*Fleet         `json:"fleets" gorm:"foreignKey:PlayerID;references:ID"`
	MineFields     []*MineField     `json:"mineFields" gorm:"foreignKey:PlayerID;references:ID"`
	MineralPackets []*MineralPacket `json:"mineralPackets" gorm:"foreignKey:PlayerID;references:ID"`
}

// create a new player with an existing race. The race
// will be copied for the player
func NewPlayer(userID int64, race *Race) *Player {

	// copy this race for the player
	playerRace := *race
	playerRace.ID = 0
	playerRace.CreatedAt = time.Time{}
	playerRace.UpdatedAt = time.Time{}

	return &Player{
		UserID:            userID,
		Race:              playerRace,
		Color:             "#0000FF", // default to blue
		Stats:             &PlayerStats{},
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

// Get a player ShipDesign by name, or nil if no design found
func (p *Player) GetDesign(name string) *ShipDesign {
	for i := range p.Designs {
		design := &p.Designs[i]
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
				latest = &design
			} else {
				if latest.Version < design.Version {
					latest = &design
				}
			}

		}
	}

	return latest
}

func computePlayerSpec(player *Player, rules *Rules) *PlayerSpec {
	techs := rules.techs
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
