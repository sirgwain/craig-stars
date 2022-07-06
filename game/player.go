package game

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Player struct {
	ID                    uint              `gorm:"primaryKey" json:"id,omitempty" header:"Username"`
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
	ID              uint             `gorm:"primaryKey" json:"id" header:"Username"`
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


