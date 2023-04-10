package cs

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
	PlayerOrders
	PlayerIntels
	PlayerPlans
	ID                    int64                `json:"id,omitempty"`
	CreatedAt             time.Time            `json:"createdAt,omitempty"`
	UpdatedAt             time.Time            `json:"updatedAt,omitempty"`
	GameID                int64                `json:"gameId,omitempty"`
	UserID                int64                `json:"userId,omitempty"`
	Name                  string               `json:"name,omitempty"`
	Num                   int                  `json:"num,omitempty"`
	Ready                 bool                 `json:"ready,omitempty"`
	AIControlled          bool                 `json:"aIControlled,omitempty"`
	SubmittedTurn         bool                 `json:"submittedTurn,omitempty"`
	Color                 string               `json:"color,omitempty"`
	DefaultHullSet        int                  `json:"defaultHullSet,omitempty"`
	Race                  Race                 `json:"race,omitempty"`
	TechLevels            TechLevel            `json:"techLevels,omitempty"`
	TechLevelsSpent       TechLevel            `json:"techLevelsSpent,omitempty"`
	ResearchSpentLastYear int                  `json:"researchSpentLastYear,omitempty"`
	Relations             []PlayerRelationship `json:"relations,omitempty"`
	Messages              []PlayerMessage      `json:"messages,omitempty"`
	Designs               []*ShipDesign        `json:"designs,omitempty"`
	Spec                  PlayerSpec           `json:"spec,omitempty"`
	Stats                 *PlayerStats         `json:"stats,omitempty"`
	leftoverResources     int
}

type PlayerIntels struct {
	PlayerIntels        []PlayerIntel        `json:"playerIntels,omitempty"`
	PlanetIntels        []PlanetIntel        `json:"planetIntels,omitempty"`
	FleetIntels         []FleetIntel         `json:"fleetIntels,omitempty"`
	ShipDesignIntels    []ShipDesignIntel    `json:"shipDesignIntels,omitempty"`
	MineralPacketIntels []MineralPacketIntel `json:"mineralPacketIntels,omitempty"`
	MineFieldIntels     []MineFieldIntel     `json:"mineFieldIntels,omitempty"`
	WormholeIntels      []WormholeIntel      `json:"wormholeIntels,omitempty"`
}

type PlayerPlans struct {
	ProductionPlans []ProductionPlan `json:"productionPlans,omitempty"`
	BattlePlans     []BattlePlan     `json:"battlePlans,omitempty"`
	TransportPlans  []TransportPlan  `json:"transportPlans,omitempty"`
}

type PlayerOrders struct {
	Researching       TechField         `json:"researching,omitempty"`
	NextResearchField NextResearchField `json:"nextResearchField,omitempty"`
	ResearchAmount    int               `json:"researchAmount,omitempty"`
}

type PlayerStats struct {
	FleetsBuilt      int `json:"fleetsBuilt,omitempty"`
	TokensBuilt      int `json:"tokensBuilt,omitempty"`
	PlanetsColonized int `json:"planetsColonized,omitempty"`
}

type PlayerRelationship struct {
	Relation PlayerRelation `json:"relation,omitempty"`
	ShareMap bool           `json:"shareMap,omitempty"`
}

type PlayerRelation string

const (
	PlayerRelationNeutral PlayerRelation = "Neutral"
	PlayerRelationFriend  PlayerRelation = "Friend"
	PlayerRelationEnemy   PlayerRelation = "Enemy"
)

type PlayerSpec struct {
	PlanetaryScanner         *TechPlanetaryScanner               `json:"planetaryScanner,omitempty"`
	Defense                  *TechDefense                        `json:"defense,omitempty"`
	Terraform                map[TerraformHabType]*TechTerraform `json:"terraform,omitempty"`
	ResourcesPerYear         int                                 `json:"resourcesPerYear,omitempty"`
	ResourcesPerYearResearch int                                 `json:"resourcesPerYearResearch,omitempty"`
	CurrentResearchCost      int                                 `json:"currentResearchCost,omitempty"`
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
	Tasks WaypointTransportTasks `json:"tasks,omitempty"`
}

type ProductionPlan struct {
	Name                              string               `json:"name"`
	Items                             []ProductionPlanItem `json:"items"`
	ContributesOnlyLeftoverToResearch bool                 `json:"contributesOnlyLeftoverToResearch,omitempty"`
}

type ProductionPlanItem struct {
	Type       QueueItemType `json:"type"`
	DesignName string        `json:"designName"`
	Quantity   int           `json:"quantity"`
	Allocated  Cost          `json:"allocated"`
}

// All mapobjects that a player can issue commands to
type PlayerMapObjects struct {
	Planets        []*Planet        `json:"planets"`
	Fleets         []*Fleet         `json:"fleets"`
	Starbases      []*Fleet         `json:"starbases"`
	MineFields     []*MineField     `json:"mineFields"`
	MineralPackets []*MineralPacket `json:"mineralPackets"`
}

// create a new player with an existing race. The race
// will be copied for the player
func NewPlayer(userID int64, race *Race) *Player {

	// copy this race for the player
	playerRace := *race
	playerRace.ID = 0
	playerRace.UpdatedAt = time.Time{}
	playerRace.CreatedAt = time.Time{}

	return &Player{
		UserID: userID,
		Race:   playerRace,
		Color:  "#0000FF", // default to blue
		Stats:  &PlayerStats{},
		PlayerOrders: PlayerOrders{
			Researching:       Energy,
			ResearchAmount:    15,
			NextResearchField: NextResearchFieldLowestField,
		},
		PlayerPlans: PlayerPlans{
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
		},
	}
}

func (p *Player) WithNum(num int) *Player {
	p.Num = num
	return p
}

func (p *Player) WithTechLevels(tl TechLevel) *Player {
	p.TechLevels = tl
	return p
}

func (p *Player) WithTechLevelsSpent(tl TechLevel) *Player {
	p.TechLevelsSpent = tl
	return p
}

func (p *Player) WithResearching(field TechField) *Player {
	p.Researching = field
	return p
}

func (p *Player) WithNextResearchField(field NextResearchField) *Player {
	p.NextResearchField = field
	return p
}

func (p *Player) withSpec(rules *Rules) *Player {
	p.Spec = computePlayerSpec(p, rules, []*Planet{})
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

// get the next design number to use
func (p *Player) GetNextDesignNum() int {
	num := 0
	for _, design := range p.Designs {
		num = maxInt(num, design.Num)
	}
	return num + 1
}

func computePlayerSpec(player *Player, rules *Rules, planets []*Planet) PlayerSpec {
	researcher := NewResearcher(rules)
	techs := rules.techs
	spec := PlayerSpec{
		PlanetaryScanner: techs.GetBestPlanetaryScanner(player),
		Defense:          techs.GetBestDefense(player),
		Terraform: map[TerraformHabType]*TechTerraform{
			TerraformHabTypeAll:  techs.GetBestTerraform(player, TerraformHabTypeAll),
			TerraformHabTypeGrav: techs.GetBestTerraform(player, TerraformHabTypeGrav),
			TerraformHabTypeTemp: techs.GetBestTerraform(player, TerraformHabTypeTemp),
			TerraformHabTypeRad:  techs.GetBestTerraform(player, TerraformHabTypeRad),
		},
	}

	for _, planet := range planets {
		if planet.OwnedBy(player.Num) {
			spec.ResourcesPerYear += planet.Spec.ResourcesPerYear
			spec.ResourcesPerYearResearch += planet.Spec.ResourcesPerYearResearch
		}
	}

	spec.CurrentResearchCost = researcher.getTotalCost(player, player.Researching, player.TechLevels.Get(player.Researching))
	return spec
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

func (p *Player) IsFriend(playerNum int) bool {
	return playerNum > 0 && playerNum <= len(p.Relations) && p.Relations[playerNum-1].Relation == PlayerRelationFriend
}

func (p *Player) IsEnemy(playerNum int) bool {
	return playerNum > 0 && playerNum <= len(p.Relations) && p.Relations[playerNum-1].Relation == PlayerRelationEnemy
}
