package game

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Tags map[string]string

type NewGamePlayerType string

const (
	NewGamePlayerTypeHost   NewGamePlayerType = "Host"
	NewGamePlayerTypeInvite NewGamePlayerType = "Invite"
	NewGamePlayerTypeOpen   NewGamePlayerType = "Open"
	NewGamePlayerTypeAI     NewGamePlayerType = "AI"
)

type NewGamePlayer struct {
	Type   NewGamePlayerType `json:"type,omitempty"`
	RaceID uint              `json:"raceID,omitempty"`
}

type GameSettings struct {
	Name                         string            `json:"name"`
	QuickStartTurns              int               `json:"quickStartTurns"`
	Size                         Size              `json:"size"`
	Density                      Density           `json:"density"`
	PlayerPositions              PlayerPositions   `json:"playerPositions"`
	RandomEvents                 bool              `json:"randomEvents"`
	ComputerPlayersFormAlliances bool              `json:"computerPlayersFormAlliances"`
	PublicPlayerScores           bool              `json:"publicPlayerScores"`
	StartMode                    GameStartMode     `json:"startMode"`
	VictoryConditions            VictoryConditions `json:"victoryConditions"`
	Players                      []NewGamePlayer   `json:"players,omitempty"`
}

type Game struct {
	ID                           uint              `gorm:"primaryKey" json:"id" header:"Name"`
	CreatedAt                    time.Time         `json:"createdAt"`
	UpdatedAt                    time.Time         `json:"updatedat"`
	DeletedAt                    gorm.DeletedAt    `gorm:"index" json:"deletedAt"`
	Name                         string            `json:"name" header:"Name"`
	HostID                       uint              `json:"hostId"`
	QuickStartTurns              int               `json:"quickStartTurns"`
	Size                         Size              `json:"size"`
	Density                      Density           `json:"density"`
	PlayerPositions              PlayerPositions   `json:"playerPositions"`
	RandomEvents                 bool              `json:"randomEvents"`
	ComputerPlayersFormAlliances bool              `json:"computerPlayersFormAlliances"`
	PublicPlayerScores           bool              `json:"publicPlayerScores"`
	StartMode                    GameStartMode     `json:"startMode"`
	Year                         int               `json:"year"`
	State                        GameState         `json:"state"`
	VictoryConditions            VictoryConditions `json:"victoryConditions" gorm:"embedded;embeddedPrefix:victory_condition_"`
	VictorDeclared               bool              `json:"victorDeclared"`
	Area                         Vector            `json:"area,omitempty" gorm:"embedded;embeddedPrefix:area_"`
	Players                      []Player          `json:"players,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	Planets                      []Planet          `json:"planets,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	Fleets                       []Fleet           `json:"fleets,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	Rules                        Rules             `json:"rules" gorm:"constraint:OnDelete:CASCADE;"`
}

type VictoryConditions struct {
	Conditions               []VictoryCondition `json:"conditions" gorm:"serializer:json"`
	NumCriteriaRequired      int                `json:"numCriteriaRequired"`
	YearsPassed              int                `json:"yearsPassed"`
	OwnPlanets               int                `json:"ownPlanets"`
	AttainTechLevel          int                `json:"attainTechLevel"`
	AttainTechLevelNumFields int                `json:"attainTechLevelNumFields"`
	ExceedsScore             int                `json:"exceedsScore"`
	ExceedsSecondPlaceScore  int                `json:"exceedsSecondPlaceScore"`
	ProductionCapacity       int                `json:"productionCapacity"`
	OwnCapitalShips          int                `json:"ownCapitalShips"`
	HighestScoreAfterYears   int                `json:"highestScoreAfterYears"`
}

type Size string

const (
	SizeTiny       Size = "Tiny"
	SizeTinyWide   Size = "TinyWide"
	SizeSmall      Size = "Small"
	SizeSmallWide  Size = "SmallWide"
	SizeMedium     Size = "Medium"
	SizeMediumWide Size = "MediumWide"
	SizeLarge      Size = "Large"
	SizeLargeWide  Size = "LargeWide"
	SizeHuge       Size = "Huge"
	SizeHugeWide   Size = "HugeWide"
)

type Density string

const (
	DensitySparse Density = "Sparse"
	DensityNormal Density = "Normal"
	DensityDense  Density = "Dense"
	DensityPacked Density = "Packed"
)

type PlayerPositions string

const (
	PlayerPositionsClose    PlayerPositions = "Close"
	PlayerPositionsModerate PlayerPositions = "Moderate"
	PlayerPositionsFarther  PlayerPositions = "Farther"
	PlayerPositionsDistant  PlayerPositions = "Distant"
)

type GameStartMode string

const (
	GameStartModeNormal   GameStartMode = "Normal"
	GameStartModeMidGame  GameStartMode = "MidGame"
	GameStartModeLateGame GameStartMode = "LateGame"
	GameStartModeEndGame  GameStartMode = "EndGame"
)

type GameState string

const (
	GameStateSetup             GameState = "Setup"
	GameStateWaitingForPlayers GameState = "WaitingForPlayers"
	GameStateGeneratingTurn    GameState = "GeneratingTurn"
)

type VictoryCondition string

const (
	VictoryConditionOwnPlanets              VictoryCondition = "OwnPlanets"
	VictoryConditionAttainTechLevels        VictoryCondition = "AttainTechLevels"
	VictoryConditionExceedsScore            VictoryCondition = "ExceedsScore"
	VictoryConditionExceedsSecondPlaceScore VictoryCondition = "ExceedsSecondPlaceScore"
	VictoryConditionProductionCapacity      VictoryCondition = "ProductionCapacity"
	VictoryConditionOwnCapitalShips         VictoryCondition = "OwnCapitalShips"
	VictoryConditionHighestScoreAfterYears  VictoryCondition = "HighestScoreAfterYears"
)

func NewGame() *Game {
	rules := NewRules()
	return &Game{
		Name:            "A Barefoot Jaywalk",
		Size:            SizeSmall,
		Density:         DensityNormal,
		PlayerPositions: PlayerPositionsModerate,
		RandomEvents:    true,
		StartMode:       GameStartModeNormal,
		Year:            2400,
		State:           GameStateSetup,
		VictoryConditions: VictoryConditions{
			Conditions: []VictoryCondition{
				VictoryConditionOwnPlanets,
				VictoryConditionAttainTechLevels,
				VictoryConditionExceedsSecondPlaceScore,
			},
			NumCriteriaRequired:      1,
			YearsPassed:              50,
			OwnPlanets:               60,
			AttainTechLevel:          22,
			AttainTechLevelNumFields: 4,
			ExceedsScore:             11000,
			ExceedsSecondPlaceScore:  100,
			ProductionCapacity:       100000,
			OwnCapitalShips:          100,
			HighestScoreAfterYears:   100,
		},
		Rules: rules,
	}
}

func (g *Game) String() string {
	return fmt.Sprintf("%s (%d)", g.Name, g.ID)
}

// update the game with new settings
func (g *Game) WithSettings(settings *GameSettings) *Game {
	g.Name = settings.Name
	g.QuickStartTurns = settings.QuickStartTurns
	g.Size = settings.Size
	g.Density = settings.Density
	g.PlayerPositions = settings.PlayerPositions
	g.RandomEvents = settings.RandomEvents
	g.ComputerPlayersFormAlliances = settings.ComputerPlayersFormAlliances
	g.PublicPlayerScores = settings.PublicPlayerScores
	g.StartMode = settings.StartMode
	g.VictoryConditions = settings.VictoryConditions

	return g
}

// Add the player to the game, and compute the player's spec based on the game rules
func (g *Game) AddPlayer(p *Player) *Player {
	p.Race.Spec = computeRaceSpec(&p.Race, &g.Rules)
	g.Players = append(g.Players, *p)
	return p
}

// compute the specs for a universe, i.e. planets, designs, fleets
func (g *Game) computeSpecs() {
	for i := range g.Players {
		player := &g.Players[i]
		player.Spec = computePlayerSpec(player, &g.Rules)
	}

	for i := range g.Planets {
		planet := &g.Planets[i]
		if planet.Owned() {
			player := &g.Players[*planet.PlayerNum]
			planet.Spec = computePlanetSpec(&g.Rules, planet, player)
		}
	}
}

// Generate a new universe
func (g *Game) GenerateUniverse() error {
	log.Debug().Msgf("%s: Generating universe", g)
	area, err := g.Rules.GetArea(g.Size)
	if err != nil {
		return err
	}

	if err := generatePlanets(g, area); err != nil {
		return err
	}

	// save our area
	g.Area = area

	generateWormholes(g)

	generatePlayerTechLevels(g)
	generatePlayerPlans(g)
	generatePlayerShipDesigns(g)

	if err := generatePlayerHomeworlds(g, area); err != nil {
		return err
	}

	if err := generatePlayerPlanetReports(g); err != nil {
		return err
	}

	// generatePlayerFleets(g)
	applyGameStartModeModifier(g)

	// setup all the specs for planets, fleets, etc
	g.computeSpecs()

	return nil
}

// generate a new turn
func (g *Game) GenerateTurn() error {
	generateTurn(g)
	return nil
}

// check if all players have submitted their turn
func (g *Game) CheckAllPlayersSubmitted() bool {
	for _, player := range g.Players {
		if !player.SubmittedTurn {
			return false
		}
	}
	return true
}

func (g *Game) GetOwnedPlanets() []Planet {
	var ownedPlanets []Planet

	for _, p := range g.Planets {
		if p.Owned() {
			ownedPlanets = append(ownedPlanets, p)
		}
	}

	return ownedPlanets
}
