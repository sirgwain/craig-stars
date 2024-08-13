package cs

import (
	"fmt"
	"slices"
	"sort"
	"strings"
	"time"
)

// A Player contains all intel, messages, tech levels, and research orders for a single empire in the game.
// It is tied to a single User (or no user, for AI)
type Player struct {
	GameDBObject
	PlayerOrders
	PlayerIntels
	PlayerPlans
	UserID                    int64                `json:"userId,omitempty"`
	Name                      string               `json:"name,omitempty"`
	Num                       int                  `json:"num,omitempty"`
	Ready                     bool                 `json:"ready"`
	AIControlled              bool                 `json:"aiControlled,omitempty"`
	AIDifficulty              AIDifficulty         `json:"aiDifficulty,omitempty"`
	Guest                     bool                 `json:"guest,omitempty"`
	SubmittedTurn             bool                 `json:"submittedTurn"`
	Color                     string               `json:"color,omitempty"`
	DefaultHullSet            int                  `json:"defaultHullSet,omitempty"`
	Race                      Race                 `json:"race,omitempty"`
	TechLevels                TechLevel            `json:"techLevels,omitempty"`
	TechLevelsSpent           TechLevel            `json:"techLevelsSpent,omitempty"`
	ResearchSpentLastYear     int                  `json:"researchSpentLastYear,omitempty"`
	Relations                 []PlayerRelationship `json:"relations,omitempty"`
	Messages                  []PlayerMessage      `json:"messages,omitempty"`
	Designs                   []*ShipDesign        `json:"designs,omitempty"`
	ScoreHistory              []PlayerScore        `json:"scoreHistory"`
	AchievedVictoryConditions Bitmask              `json:"achievedVictoryConditions,omitempty"`
	Victor                    bool                 `json:"victor"`
	Stats                     *PlayerStats         `json:"stats,omitempty"`
	Spec                      PlayerSpec           `json:"spec,omitempty"`
	leftoverResources         int
	techLevelGained           bool
	discoverer                discoverer
}

// a player and all mapobjects the player owns
// this is used by the UI when loading a player's game
type FullPlayer struct {
	Player
	PlayerMapObjects
}

type PlayerStatus struct {
	UpdatedAt     *time.Time `json:"updatedAt,omitempty"`
	UserID        int64      `json:"userId,omitempty"`
	Name          string     `json:"name,omitempty"`
	Num           int        `json:"num,omitempty"`
	Ready         bool       `json:"ready,omitempty"`
	AIControlled  bool       `json:"aiControlled,omitempty"`
	Guest         bool       `json:"guest,omitempty"`
	SubmittedTurn bool       `json:"submittedTurn,omitempty"`
	Color         string     `json:"color,omitempty"`
	Victor        bool       `json:"victor,omitempty"`
}

type PlayerIntels struct {
	BattleRecords       []BattleRecord       `json:"battleRecords,omitempty"`
	PlayerIntels        []PlayerIntel        `json:"playerIntels,omitempty"`
	ScoreIntels         []ScoreIntel         `json:"scoreIntels,omitempty"`
	PlanetIntels        []PlanetIntel        `json:"planetIntels,omitempty"`
	FleetIntels         []FleetIntel         `json:"fleetIntels,omitempty"`
	StarbaseIntels      []FleetIntel         `json:"starbaseIntels,omitempty"`
	ShipDesignIntels    []ShipDesignIntel    `json:"shipDesignIntels,omitempty"`
	MineralPacketIntels []MineralPacketIntel `json:"mineralPacketIntels,omitempty"`
	MineFieldIntels     []MineFieldIntel     `json:"mineFieldIntels,omitempty"`
	WormholeIntels      []WormholeIntel      `json:"wormholeIntels,omitempty"`
	MysteryTraderIntels []MysteryTraderIntel `json:"mysteryTraderIntels,omitempty"`
	SalvageIntels       []SalvageIntel       `json:"salvageIntels,omitempty"`
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
	StarbasesBuilt   int `json:"starbasesBuilt,omitempty"`
	TokensBuilt      int `json:"tokensBuilt,omitempty"`
	PlanetsColonized int `json:"planetsColonized,omitempty"`
}

type PlayerRelationship struct {
	Relation PlayerRelation `json:"relation"`
	ShareMap bool           `json:"shareMap,omitempty"`
}

type PlayerRelation string

const (
	PlayerRelationNeutral PlayerRelation = "Neutral"
	PlayerRelationFriend  PlayerRelation = "Friend"
	PlayerRelationEnemy   PlayerRelation = "Enemy"
)

type PlayerSpec struct {
	PlanetaryScanner                  TechPlanetaryScanner                `json:"planetaryScanner"`
	Defense                           TechDefense                         `json:"defense"`
	Terraform                         map[TerraformHabType]*TechTerraform `json:"terraform"`
	ResourcesPerYear                  int                                 `json:"resourcesPerYear"`
	ResourcesPerYearResearch          int                                 `json:"resourcesPerYearResearch"`
	ResourcesPerYearResearchEstimated int                                 `json:"resourcesPerYearResearchEstimated"`
	CurrentResearchCost               int                                 `json:"currentResearchCost"`
}

type PlayerScore struct {
	Planets                   int     `json:"planets"`
	Starbases                 int     `json:"starbases"`
	UnarmedShips              int     `json:"unarmedShips"`
	EscortShips               int     `json:"escortShips"`
	CapitalShips              int     `json:"capitalShips"`
	TechLevels                int     `json:"techLevels"`
	Resources                 int     `json:"resources"`
	Score                     int     `json:"score"`
	Rank                      int     `json:"rank"`
	AchievedVictoryConditions Bitmask `json:"achievedVictoryConditions"`
}

type BattlePlan struct {
	Num             int             `json:"num"`
	Name            string          `json:"name"`
	PrimaryTarget   BattleTarget    `json:"primaryTarget"`
	SecondaryTarget BattleTarget    `json:"secondaryTarget"`
	Tactic          BattleTactic    `json:"tactic"`
	AttackWho       BattleAttackWho `json:"attackWho"`
	DumpCargo       bool            `json:"dumpCargo"`
}

type BattleTarget string

const (
	BattleTargetNone              BattleTarget = ""
	BattleTargetAny               BattleTarget = "Any"
	BattleTargetStarbase          BattleTarget = "Starbase"
	BattleTargetArmedShips        BattleTarget = "ArmedShips"
	BattleTargetBombersFreighters BattleTarget = "BombersFreighters"
	BattleTargetUnarmedShips      BattleTarget = "UnarmedShips"
	BattleTargetFuelTransports    BattleTarget = "FuelTransports"
	BattleTargetFreighters        BattleTarget = "Freighters"
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
	Num   int                    `json:"num"`
	Name  string                 `json:"name"`
	Tasks WaypointTransportTasks `json:"tasks,omitempty"`
}

type ProductionPlan struct {
	Num                               int                  `json:"num"`
	Name                              string               `json:"name"`
	Items                             []ProductionPlanItem `json:"items"`
	ContributesOnlyLeftoverToResearch bool                 `json:"contributesOnlyLeftoverToResearch,omitempty"`
}

type ProductionPlanItem struct {
	Type      QueueItemType `json:"type"`
	DesignNum int           `json:"designNum"`
	Quantity  int           `json:"quantity"`
}

// Apply a production plan to a planet
func (plan *ProductionPlan) Apply(planet *Planet) {
	planet.ProductionQueue = make([]ProductionQueueItem, len(plan.Items))
	for i, item := range plan.Items {
		planet.ProductionQueue[i] = ProductionQueueItem{
			Type:      item.Type,
			Quantity:  item.Quantity,
			DesignNum: item.DesignNum,
		}
	}
	planet.ContributesOnlyLeftoverToResearch = plan.ContributesOnlyLeftoverToResearch
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

	player := &Player{
		UserID: userID,
		Race:   playerRace,
		Color:  "#0000FF", // default to blue
		Stats:  &PlayerStats{},
		PlayerOrders: PlayerOrders{
			Researching:       Energy,
			ResearchAmount:    15,
			NextResearchField: NextResearchFieldLowestField,
		},
	}

	// start with a base discoverer
	player.discoverer = newDiscoverer(player)
	player.PlayerPlans = player.defaultPlans()
	return player
}

func (p *Player) WithNum(num int) *Player {
	p.Num = num
	return p
}

func (p *Player) WithAIControlled(aiControlled bool) *Player {
	p.AIControlled = aiControlled
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

// return the most recent PlayerScore or an empty score if there is no score history
func (p *Player) GetScore() PlayerScore {
	if len(p.ScoreHistory) > 0 {
		return p.ScoreHistory[len(p.ScoreHistory)-1]
	}
	return PlayerScore{}
}

// Get a player ShipDesign, or nil if no design found
func (p *Player) GetDesign(num int) *ShipDesign {
	for _, design := range p.Designs {
		if design.Num == num {
			return design
		}
	}
	return nil
}

// Get a ShipDesignIntel, or nil if no design found
func (p *Player) GetForeignDesign(playerNum int, num int) *ShipDesignIntel {
	for i := range p.ShipDesignIntels {
		design := &p.ShipDesignIntels[i]
		if design.PlayerNum == playerNum && design.Num == num {
			return design
		}
	}
	return nil
}

// Get a player ShipDesign, or nil if no design found
func (p *Player) GetDesignByName(name string) *ShipDesign {
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
func (p *Player) GetNextDesignNum(designs []*ShipDesign) int {
	num := 0
	for _, design := range designs {
		num = MaxInt(num, design.Num)
	}
	return num + 1
}

// get the next BattlePlan number to use
func (p *Player) GetNextBattlePlanNum() int {
	num := 0
	for _, plan := range p.BattlePlans {
		num = MaxInt(num, plan.Num)
	}
	return num + 1
}

// get the next ProductionPlan number to use
func (p *Player) GetNextProductionPlanNum() int {
	num := 0
	for _, plan := range p.ProductionPlans {
		num = MaxInt(num, plan.Num)
	}
	return num + 1
}

// get the next TransportPlan number to use
func (p *Player) GetNextTransportPlanNum() int {
	num := 0
	for _, plan := range p.TransportPlans {
		num = MaxInt(num, plan.Num)
	}
	return num + 1
}

// clear this player's transient intel
func (p *Player) clearTransientIntel() {
	p.FleetIntels = []FleetIntel{}
	p.MineFieldIntels = []MineFieldIntel{}
	p.SalvageIntels = []SalvageIntel{}
	p.MineralPacketIntels = []MineralPacketIntel{}
}

// for reports that stick around, increment the report age
func (p *Player) incrementReportAge() {
	for i := range p.PlanetIntels {
		planet := &p.PlanetIntels[i]
		if planet.ReportAge != ReportAgeUnexplored {
			planet.ReportAge++
		}
	}

	for i := range p.WormholeIntels {
		wormhole := &p.WormholeIntels[i]
		if wormhole.ReportAge != ReportAgeUnexplored {
			wormhole.ReportAge++
		}
	}
}

func (p *Player) getPlanetIntel(num int) *PlanetIntel {
	return &p.PlanetIntels[num-1]
}

func (p *Player) getWormholeIntel(num int) *WormholeIntel {
	for i := range p.WormholeIntels {
		intel := &p.WormholeIntels[i]
		if intel.Num == num {
			return intel
		}
	}
	return nil
}

func (p *Player) getMysteryTraderIntel(num int) *MysteryTraderIntel {
	for i := range p.MysteryTraderIntels {
		intel := &p.MysteryTraderIntels[i]
		if intel.Num == num {
			return intel
		}
	}
	return nil
}

func (p *Player) getMineFieldIntel(playerNum, num int) *MineFieldIntel {
	for i := range p.MineFieldIntels {
		intel := &p.MineFieldIntels[i]
		if intel.PlayerNum == playerNum && intel.Num == num {
			return intel
		}
	}

	return nil
}

func (p *Player) getMineralPacketIntel(playerNum, num int) *MineralPacketIntel {
	for i := range p.MineralPacketIntels {
		intel := &p.MineralPacketIntels[i]
		if intel.PlayerNum == playerNum && intel.Num == num {
			return intel
		}
	}

	return nil
}

func (p *Player) getFleetIntel(playerNum, num int) *FleetIntel {
	for i := range p.FleetIntels {
		intel := &p.FleetIntels[i]
		if intel.PlayerNum == playerNum && intel.Num == num {
			return intel
		}
	}

	return nil
}

func (p *Player) getShipDesignIntel(playerNum, num int) *ShipDesignIntel {
	for i := range p.ShipDesignIntels {
		intel := &p.ShipDesignIntels[i]
		if intel.PlayerNum == playerNum && intel.Num == num {
			return intel
		}
	}

	return nil
}

func (p *Player) getSalvageIntel(num int) *SalvageIntel {
	for i := range p.SalvageIntels {
		intel := &p.SalvageIntels[i]
		if intel.Num == num {
			return intel
		}
	}
	return nil
}

func computePlayerSpec(player *Player, rules *Rules, planets []*Planet) PlayerSpec {
	researcher := NewResearcher(rules)
	techs := rules.techs
	spec := PlayerSpec{
		PlanetaryScanner: *techs.GetBestPlanetaryScanner(player),
		Defense:          *techs.GetBestDefense(player),
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
			spec.ResourcesPerYearResearchEstimated += planet.Spec.ResourcesPerYearResearch + planet.Spec.ResourcesPerYearResearchEstimatedLeftover
		}
	}

	spec.CurrentResearchCost = researcher.getTotalCost(player.TechLevels, player.Researching, player.Race.ResearchCost.Get(player.Researching), player.TechLevels.Get(player.Researching))
	return spec
}

// return true if the player currently has this tech
func (p *Player) HasTech(tech *Tech) bool {
	return p.CanLearnTech(tech) && p.TechLevels.HasRequiredLevels(tech.Requirements.TechLevel)
}

func (p *Player) CanLearnTech(tech *Tech) bool {
	requirements := tech.Requirements
	if len(requirements.PRTsRequired) != 0 && !slices.Contains(requirements.PRTsRequired, p.Race.PRT) {
		return false
	}
	if len(requirements.PRTsDenied) != 0 && slices.Contains(requirements.PRTsDenied, p.Race.PRT) {
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

// return the cost, in resources, it will take to reach this tech level
func (p *Player) GetResearchCost(rules *Rules, techLevel TechLevel) int {
	researcher := NewResearcher(rules)

	techLevels := p.TechLevels
	resources := 0
	for _, field := range TechFields {
		requiredLevel := techLevel.Get(field)
		playerLevel := techLevels.Get(field)
		spent := p.TechLevelsSpent.Get(field)
		researchCostLevel := p.Race.ResearchCost.Get(field)

		// we require this field
		if requiredLevel > playerLevel {
			resources -= spent
			for i := 0; i < requiredLevel-playerLevel; i++ {
				// find the cost for the next level in this field
				resources += researcher.getTotalCost(techLevels, field, researchCostLevel, playerLevel+i)
				// assume we gained this for future calcs
				techLevels.Set(field, playerLevel+i+1)
			}
		}
	}

	return resources
}

// get the default relationships for a player with other players
func (p *Player) defaultRelationships(players []*Player, aiFormsAlliances bool) []PlayerRelationship {
	relations := make([]PlayerRelationship, len(players))

	for i, otherPlayer := range players {
		relationship := &relations[i]
		if otherPlayer.Num == p.Num {
			// we're friends with ourselves
			relationship.Relation = PlayerRelationFriend
		} else if aiFormsAlliances && p.AIControlled && otherPlayer.AIControlled {
			// team up! destroy all humans!
			relationship.Relation = PlayerRelationFriend
		} else if otherPlayer.AIControlled || p.AIControlled {
			// AI is always the enemy
			relationship.Relation = PlayerRelationEnemy
		} else {
			relationship.Relation = PlayerRelationNeutral
		}

	}
	return relations
}

// get the default relationships for a player with other players
func (p *Player) defaultPlans() PlayerPlans {

	// AR races don't build factories or mines
	// CA & tri-immune races don't do terraforming
	defaultProductionPlan := ProductionPlan{
		Num:  0,
		Name: "Default",
		Items: []
	}

	if !p.Race.Spec.Instaforming && !(p.Race.ImmuneGrav && p.Race.ImmuneTemp && p.Race.ImmuneRad) {
		defaultProductionPlan.Items = append(defaultProductionPlan.Items,
			ProductionPlanItem{Type: QueueItemTypeAutoMinTerraform, Quantity: 1},
		)
	}

	if !p.Race.Spec.InnateResources {
		defaultProductionPlan.Items = append(defaultProductionPlan.Items,
			ProductionPlanItem{Type: QueueItemTypeAutoFactories, Quantity: 250},
		)
	}

	if !p.Race.Spec.InnateMining {
		defaultProductionPlan.Items = append(defaultProductionPlan.Items,
			ProductionPlanItem{Type: QueueItemTypeAutoMines, Quantity: 250},
		)
	}

	if !p.Race.Spec.Instaforming && !(p.Race.ImmuneGrav && p.Race.ImmuneTemp && p.Race.ImmuneRad) {
		defaultProductionPlan.Items = append(defaultProductionPlan.Items,
			ProductionPlanItem{Type: QueueItemTypeAutoMaxTerraform, Quantity: 1},
		)
	}

	return PlayerPlans{ProductionPlans: []ProductionPlan{
		defaultProductionPlan,
	},
		BattlePlans: []BattlePlan{
			{
				Num:             0,
				Name:            "Default",
				PrimaryTarget:   BattleTargetArmedShips,
				SecondaryTarget: BattleTargetAny,
				Tactic:          BattleTacticMaximizeDamageRatio,
				AttackWho:       BattleAttackWhoEnemiesAndNeutrals,
			},
			{
				Num:		 1,
				Name: 		 "KillStarbase",
				PrimaryTarget:	 BattleTargetStarbase,
				SecondaryTarget: BattleTargetArmedShips,
				Tactic:		 BattleTacticMaximizeDamageRatio,
				AttackWho:	 BattleAttackWhoEnemiesAndNeutrals,
			},
			{
				Num:		 2,
				Name: 		 "Max Defense",
				PrimaryTarget:	 BattleTargetArmedShips,
				SecondaryTarget: BattleTargetBombersFreighters,
				Tactic:		 BattleTacticMaximizeNetDamage,
				AttackWho:	 BattleAttackWhoEnemiesAndNeutrals,
			},
			{
				Num:		 3,
				Name: 		 "Sniper",
				PrimaryTarget:	 BattleTargetUnarmedShips,
				SecondaryTarget: BattleTargetNone,
				Tactic:		 BattleTacticDisengageIfChallenged,
				AttackWho:	 BattleAttackWhoEnemiesAndNeutrals
			},
			{
				Num:		 4,
				Name:		 "Chicken",
				PrimaryTarget:	 BattleTargetAny,
				SecondaryTarget: BattleTargetNone,
				Tactic:		 BattleTacticDisengage,
				AttackWho:	 BattleAttackWhoEnemiesAndNeutrals,	
			},
		},
		TransportPlans: []TransportPlan{
			{
				Num:  0,
				Name: "Clear All",
			},
			{
				Num:  1,
				Name: "Quick Load",
				Tasks: WaypointTransportTasks{
					Fuel:      WaypointTransportTask{Action: TransportActionLoadOptimal},
					Ironium:   WaypointTransportTask{Action: TransportActionLoadAll},
					Boranium:  WaypointTransportTask{Action: TransportActionLoadAll},
					Germanium: WaypointTransportTask{Action: TransportActionLoadAll},
				},
			},
			{
				Num:  2,
				Name: "Quick Drop",
				Tasks: WaypointTransportTasks{
					Fuel:      WaypointTransportTask{Action: TransportActionLoadOptimal},
					Ironium:   WaypointTransportTask{Action: TransportActionUnloadAll},
					Boranium:  WaypointTransportTask{Action: TransportActionUnloadAll},
					Germanium: WaypointTransportTask{Action: TransportActionUnloadAll},
					Colonists: WaypointTransportTask{Action: TransportActionUnloadAll},
				},
			},
			{
				Num:  2,
				Name: "Wait Load",
				Tasks: WaypointTransportTasks{
					Fuel:      WaypointTransportTask{Action: TransportActionLoadOptimal},
					Ironium:   WaypointTransportTask{Action: TransportActionWaitForPercent, Amount: 100},
					Boranium:  WaypointTransportTask{Action: TransportActionWaitForPercent, Amount: 100},
					Germanium: WaypointTransportTask{Action: TransportActionWaitForPercent, Amount: 100},
				},
			},
			{
				Num:  3,
				Name: "Load Colonists",
				Tasks: WaypointTransportTasks{
					Colonists: WaypointTransportTask{Action: TransportActionLoadAll},
				},
			},
			{
				Num:  4,
				Name: "Unload Colonists",
				Tasks: WaypointTransportTasks{
					Colonists: WaypointTransportTask{Action: TransportActionLoadDunnage},
				},
			},
		},
	}
}

// get the default intels for a player for other players
func (p *Player) defaultPlayerIntels(players []*Player) []PlayerIntel {
	playerIntels := make([]PlayerIntel, len(players))
	for i, otherPlayer := range players {
		playerIntel := &playerIntels[i]
		playerIntel.Color = otherPlayer.Color
		playerIntel.Name = otherPlayer.Name
		playerIntel.Num = otherPlayer.Num

		// we know about ourselves
		if otherPlayer.Num == p.Num {
			playerIntel.Seen = true
			playerIntel.RaceName = p.Race.Name
			playerIntel.RacePluralName = p.Race.PluralName
		}
	}

	return playerIntels
}

// get the default intels for a player for other players
func (player *Player) initDefaultPlanetIntels(planets []*Planet) error {
	player.PlanetIntels = make([]PlanetIntel, len(planets))
	for j := range planets {
		// start with some defaults
		intel := &player.PlanetIntels[j]
		intel.ReportAge = ReportAgeUnexplored
		intel.Type = MapObjectTypePlanet
		intel.PlayerNum = Unowned

		// everyone knows about this
		planet := planets[j]
		intel.Position = planet.Position
		intel.Name = planet.Name
		intel.Num = planet.Num
	}

	return nil
}

func (p *Player) IsFriend(playerNum int) bool {
	return playerNum > 0 && playerNum <= len(p.Relations) && p.Relations[playerNum-1].Relation == PlayerRelationFriend
}

func (p *Player) IsEnemy(playerNum int) bool {
	return playerNum > 0 && playerNum <= len(p.Relations) && p.Relations[playerNum-1].Relation == PlayerRelationEnemy
}

func (p *Player) IsNeutral(playerNum int) bool {
	return playerNum > 0 && playerNum <= len(p.Relations) && p.Relations[playerNum-1].Relation == PlayerRelationNeutral
}

func (p *Player) IsSharingMap(playerNum int) bool {
	return playerNum != p.Num && playerNum > 0 && playerNum <= len(p.Relations) && p.Relations[playerNum-1].Relation == PlayerRelationFriend && p.Relations[playerNum-1].ShareMap
}

func (p *Player) getNextFleetNum(playerFleets []*Fleet) int {
	num := 1

	orderedFleets := make([]*Fleet, len(playerFleets))
	copy(orderedFleets, playerFleets)
	sort.Slice(orderedFleets, func(i, j int) bool { return orderedFleets[i].Num < orderedFleets[j].Num })

	for i := 0; i < len(orderedFleets); i++ {
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

// get the next mineral packet number to use
func (p *Player) getNextMineralPacketNum(packets []*MineralPacket) int {
	num := 0
	for _, packet := range packets {
		num = MaxInt(num, packet.Num)
	}
	return num + 1
}

// inject player designs into tokens for a slice of fleets
func (p *Player) InjectDesigns(fleets []*Fleet) error {

	designsByNum := make(map[int]*ShipDesign, len(p.Designs))
	for i := range p.Designs {
		design := p.Designs[i]
		designsByNum[design.Num] = design
	}

	// inject the design into this
	for _, f := range fleets {
		for i := range f.Tokens {
			token := &f.Tokens[i]
			token.design = designsByNum[token.DesignNum]
			if token.design == nil {
				return fmt.Errorf("unable to find design %d for fleet %s", token.DesignNum, f.Name)
			}
		}
	}

	return nil
}

// validate this battle plan
func (plan *BattlePlan) Validate(player *Player) error {
	if strings.TrimSpace(plan.Name) == "" {
		return fmt.Errorf("plan has no name")
	}
	return nil
}

// validate this production plan
func (plan *ProductionPlan) Validate(player *Player) error {
	if strings.TrimSpace(plan.Name) == "" {
		return fmt.Errorf("plan has no name")
	}
	return nil
}

// validate this transport plan
func (plan *TransportPlan) Validate(player *Player) error {
	if strings.TrimSpace(plan.Name) == "" {
		return fmt.Errorf("plan has no name")
	}
	return nil
}
