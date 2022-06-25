package game

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Player struct {
	ID                    uint           `gorm:"primaryKey" json:"id,omitempty" header:"Username"`
	CreatedAt             time.Time      `json:"createdAt,omitempty"`
	UpdatedAt             time.Time      `json:"updatedat,omitempty"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	GameID                uint           `json:"gameId"`
	UserID                uint           `json:"userId"`
	Name                  string         `json:"name"`
	Num                   int            `json:"num"`
	Ready                 bool           `json:"ready,omitempty"`
	AIControlled          bool           `json:"aIControlled,omitempty"`
	SubmittedTurn         bool           `json:"submittedTurn,omitempty"`
	Color                 string         `json:"color,omitempty"`
	Race                  Race           `json:"race,omitempty"`
	TechLevels            TechLevel      `json:"techLevels,omitempty" gorm:"embedded;embeddedPrefix:tech_levels_"`
	TechLevelsSpent       TechLevel      `json:"techLevelsSpent,omitempty" gorm:"embedded;embeddedPrefix:tech_levels_spent"`
	ResearchAmount        int            `json:"researchAmount,omitempty"`
	ResearchSpentLastYear int            `json:"researchSpentLastYear,omitempty"`
	Researching           TechField      `json:"researching,omitempty"`
	Fleets                []Fleet        `json:"fleets,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	Planets               []Planet       `json:"planets,omitempty" gorm:"foreignKey:PlayerID;references:ID"`
	PlanetIntels          []PlanetIntel  `json:"planetIntels,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
}

type BattlePlan struct {
	ID              uint           `gorm:"primaryKey" json:"id" header:"Username"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedat"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	PlayerID        uint           `json:"playerId"`
	Name            string         `json:"name"`
	PrimaryTarget   string         `json:"primaryTarget"`
	SecondaryTarget string         `json:"secondaryTarget"`
	Tactic          string         `json:"tactic"`
	AttackWho       string         `json:"attackWho"`
}

func NewPlayer(userID uint, race *Race) *Player {
	return &Player{
		UserID: userID,
		Race:   *race,
	}
}

func (p *Player) WithTechLevels(tl TechLevel) *Player {
	p.TechLevels = tl
	return p
}

func (p *Player) String() string {
	return fmt.Sprintf("Player %d (%d) %s", p.Num, p.ID, p.Race.PluralName)
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
