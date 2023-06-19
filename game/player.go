package game

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Player struct {
	ID                    uint              `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt             time.Time         `json:"createdAt,omitempty"`
	UpdatedAt             time.Time         `json:"updatedat,omitempty"`
	DeletedAt             gorm.DeletedAt    `gorm:"index" json:"deletedAt,omitempty"`
	GameID                uint              `json:"gameId,omitempty"`
	UserID                uint              `json:"userId,omitempty"`
	Name                  string            `json:"name,omitempty"`
	Num                   int               `json:"num,omitempty"`
	Ready                 bool              `json:"ready,omitempty"`
	AIControlled          bool              `json:"aIControlled,omitempty"`
	SubmittedTurn         bool              `json:"submittedTurn,omitempty"`
	Color                 string            `json:"color,omitempty"`
	DefaultHullSet        int               `json:"defaultHullSet,omitempty"`
	Race                  Race              `json:"race,omitempty"`
	TechLevels            TechLevel         `json:"techLevels,omitempty" gorm:"embedded;embeddedPrefix:tech_levels_"`
	TechLevelsSpent       TechLevel         `json:"techLevelsSpent,omitempty" gorm:"embedded;embeddedPrefix:tech_levels_spent"`
	ResearchAmount        int               `json:"researchAmount,omitempty"`
	ResearchSpentLastYear int               `json:"researchSpentLastYear,omitempty"`
	NextResearchField     NextResearchField `json:"nextResearchField,omitempty"`
	Researching           TechField         `json:"researching,omitempty"`
	BattlePlans           []BattlePlan      `json:"battlePlans,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	ProductionPlans       []ProductionPlan  `json:"productionPlans,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	TransportPlans        []TransportPlan   `json:"transportPlans,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	Messages              []PlayerMessage   `json:"messages,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	Designs               []*ShipDesign     `json:"designs,omitempty" gorm:"foreignKey:PlayerID;references:ID"`
	Fleets                []*Fleet          `json:"fleets,omitempty" gorm:"foreignKey:PlayerID;references:ID"`
	Planets               []*Planet         `json:"planets,omitempty" gorm:"foreignKey:PlayerID;references:ID"`
	PlanetIntels          []PlanetIntel     `json:"planetIntels,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	Spec                  *PlayerSpec       `json:"spec,omitempty" gorm:"serializer:json"`
	LeftoverResources     int               `json:"-" gorm:"-"`
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
	ID              uint             `gorm:"primaryKey" json:"id"`
	CreatedAt       time.Time        `json:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt   `gorm:"index" json:"deletedAt"`
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
	ID        uint                   `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
	DeletedAt gorm.DeletedAt         `gorm:"index" json:"deletedAt"`
	PlayerID  uint                   `json:"playerId"`
	Name      string                 `json:"name"`
	Tasks     WaypointTransportTasks `json:"tasks,omitempty" gorm:"-"`
}

type WaypointTransportTasks struct {
	Fuel      WaypointTransportTask `json:"fuel,omitempty"`
	Ironium   WaypointTransportTask `json:"ironium,omitempty"`
	Boranium  WaypointTransportTask `json:"boranium,omitempty"`
	Germanium WaypointTransportTask `json:"germanium,omitempty"`
	Colonists WaypointTransportTask `json:"colonists,omitempty"`
}

type WaypointTransportTask struct {
	Amount int                         `json:"amount,omitempty"`
	Action WaypointTaskTransportAction `json:"action,omitempty"`
}

type WaypointTaskTransportAction string

const (
	// No transport task for the specified cargo.
	TransportActionNone WaypointTaskTransportAction = ""

	// (fuel only) Load or unload fuel until the fleet carries only the exact amount
	// needed to reach the next waypoint. You can use this task to send a fleet
	// loaded with fuel to rescue a stranded fleet. The rescue fleet will transfer
	// only the amount of fuel it can spare without stranding itself.
	TransportActionLoadOptimal WaypointTaskTransportAction = "LoadOptimal"

	// Load as much of the specified cargo as the fleet can hold.
	TransportActionLoadAll WaypointTaskTransportAction = "LoadAll"

	// Unload all the specified cargo at the waypoint.
	TransportActionUnloadAll WaypointTaskTransportAction = "UnloadAll"

	// Load the amount specified only if there is room in the hold.
	TransportActionLoadAmount WaypointTaskTransportAction = "LoadAmount"

	// Unload the amount specified only if the fleet is carrying that amount.
	TransportActionUnloadAmount WaypointTaskTransportAction = "UnloadAmount"

	// Loads up to the specified portion of the cargo hold subject to amount available at waypoint and room left in hold.
	TransportActionFillPercent WaypointTaskTransportAction = "FillPercent"

	// Remain at the waypoint until exactly X % of the hold is filled.
	TransportActionWaitForPercent WaypointTaskTransportAction = "WaitForPercent"

	// (minerals and colonists only) This command waits until all other loads and unloads are complete,
	// then loads as many colonists or amount of a mineral as will fit in the remaining space. For example,
	// setting Load All Germanium, Load Dunnage Ironium, will load all the Germanium that is available,
	// then as much Ironium as possible. If more than one dunnage cargo is specified, they are loaded in
	// the order of Ironium, Boranium, Germanium, and Colonists.
	TransportActionLoadDunnage WaypointTaskTransportAction = "LoadDunnage"

	// Load or unload the cargo until the amount on board is the amount specified.
	// If less than the specified cargo is available, the fleet will not move on.
	TransportActionSetAmountTo WaypointTaskTransportAction = "SetAmountTo"

	// Load or unload the cargo until the amount at the waypoint is the amount specified.
	// This order is always carried out to the best of the fleet’s ability that turn but does not prevent the fleet from moving on.
	TransportActionSetWaypointTo WaypointTaskTransportAction = "SetWaypointTo"
)

type ProductionPlan struct {
	ID        uint                 `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time            `json:"createdAt"`
	UpdatedAt time.Time            `json:"updatedAt"`
	DeletedAt gorm.DeletedAt       `gorm:"index" json:"deletedAt"`
	PlayerID  uint                 `json:"playerId"`
	Name      string               `json:"name"`
	Items     []ProductionPlanItem `json:"items" gorm:"constraint:OnDelete:CASCADE;"`
}

type ProductionPlanItem struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	ProductionPlanID uint           `json:"-"`
	Type             QueueItemType  `json:"type"`
	DesignName       string         `json:"designName"`
	Quantity         int            `json:"quantity"`
	Allocated        Cost           `json:"allocated" gorm:"embedded;embeddedPrefix:allocated_"`
	SortOrder        int            `json:"-"`
}

// create a new player with an existing race. The race
// will be copied for the player
func NewPlayer(userID uint, race *Race) *Player {

	// copy this race for the player
	playerRace := *race
	playerRace.ID = 0
	playerRace.CreatedAt = time.Time{}
	playerRace.UpdatedAt = time.Time{}
	playerRace.DeletedAt = gorm.DeletedAt{}

	return &Player{
		UserID:            userID,
		Race:              playerRace,
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

func (p *Player) String() string {
	return fmt.Sprintf("Player %d (%d) %s", p.Num, p.ID, p.Race.PluralName)
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
