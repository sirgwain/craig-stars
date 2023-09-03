package cs

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

type Tags map[string]string

type NewGamePlayerType string

const (
	NewGamePlayerTypeHost  NewGamePlayerType = "Host"
	NewGamePlayerTypeGuest NewGamePlayerType = "Guest"
	NewGamePlayerTypeOpen  NewGamePlayerType = "Open"
	NewGamePlayerTypeAI    NewGamePlayerType = "AI"
)

type AIDifficulty string

const (
	AIDifficultyNone   AIDifficulty = ""
	AIDifficultyEasy   AIDifficulty = "Easy"
	AIDifficultyNormal AIDifficulty = "Normal"
	AIDifficultyHard   AIDifficulty = "Hard"
)

type NewGamePlayer struct {
	Type           NewGamePlayerType `json:"type,omitempty"`
	RaceID         int64             `json:"raceId,omitempty"`
	AIDifficulty   AIDifficulty      `json:"aiDifficulty,omitempty"`
	Color          string            `json:"color,omitempty"`
	DefaultHullSet int               `json:"hullSetNum,omitempty"`
	Race           Race              `json:"race,omitempty"`
}

type GameSettings struct {
	Name                         string            `json:"name"`
	Public                       bool              `json:"public"`
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
	Rules                        *Rules            `json:"rules,omitempty"`
	TechStore                    *TechStore        `json:"techStore,omitempty"`
}

type Game struct {
	DBObject
	HostID                       int64             `json:"hostId"`
	Name                         string            `json:"name" header:"Name"`
	State                        GameState         `json:"state"`
	Public                       bool              `json:"public,omitempty"`
	Hash                         string            `json:"hash"`
	Size                         Size              `json:"size"`
	Density                      Density           `json:"density"`
	PlayerPositions              PlayerPositions   `json:"playerPositions"`
	RandomEvents                 bool              `json:"randomEvents,omitempty"`
	ComputerPlayersFormAlliances bool              `json:"computerPlayersFormAlliances,omitempty"`
	PublicPlayerScores           bool              `json:"publicPlayerScores,omitempty"`
	StartMode                    GameStartMode     `json:"startMode,omitempty"`
	QuickStartTurns              int               `json:"quickStartTurns,omitempty"`
	OpenPlayerSlots              int              `json:"openPlayerSlots,omitempty"`
	NumPlayers                   int               `json:"numPlayers,omitempty"`
	VictoryConditions            VictoryConditions `json:"victoryConditions"`
	Seed                         int64             `json:"seed"`
	Rules                        Rules             `json:"rules"`
	Area                         Vector            `json:"area,omitempty"`
	Year                         int               `json:"year,omitempty"`
	VictorDeclared               bool              `json:"victorDeclared"`
}

// struct for holding a game with a list of player status
type GameWithPlayers struct {
	Game
	Players []PlayerStatus `json:"players,omitempty"`
}

// return true if this is a single player game
func (g *GameWithPlayers) IsSinglePlayer() bool {
	nonAiPlayers := 0
	for _, p := range g.Players {
		if p.AIControlled {
			nonAiPlayers++
		}
	}
	return nonAiPlayers > 1
}

// A game with players and a universe, used in universe and turn generation
type FullGame struct {
	*Game
	*Universe
	*TechStore
	Players []*Player `json:"players,omitempty"`
}

// return true if this is a single player game
func (g *FullGame) IsSinglePlayer() bool {
	nonAiPlayers := 0
	for _, p := range g.Players {
		if p.AIControlled {
			nonAiPlayers++
		}
	}
	return nonAiPlayers > 1
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
	GameStartModeNormal   GameStartMode = ""
	GameStartModeMidGame  GameStartMode = "MidGame"
	GameStartModeLateGame GameStartMode = "LateGame"
	GameStartModeEndGame  GameStartMode = "EndGame"
)

type GameState string

const (
	GameStateSetup               GameState = "Setup"
	GameStateGeneratingUniverse  GameState = "GeneratingUniverse"
	GameStateWaitingForPlayers   GameState = "WaitingForPlayers"
	GameStateGeneratingTurn      GameState = "GeneratingTurn"
	GameStateGeneratingTurnError GameState = "GeneratingTurnError"
)

type playerGetter interface {
	getPlayer(playerNum int) *Player
}

func NewGame() *Game {
	seed := time.Now().UnixNano()
	rules := NewRulesWithSeed(seed)
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
			Conditions:               Bitmask(VictoryConditionOwnPlanets) | Bitmask(VictoryConditionAttainTechLevels) | Bitmask(VictoryConditionExceedsSecondPlaceScore),
			NumCriteriaRequired:      1,
			YearsPassed:              50,
			OwnPlanets:               60,
			AttainTechLevel:          22,
			AttainTechLevelNumFields: 4,
			ExceedsScore:             11000,
			ExceedsSecondPlaceScore:  100,
			ProductionCapacity:       100,
			OwnCapitalShips:          100,
			HighestScoreAfterYears:   100,
		},
		Seed:  seed,
		Rules: rules,
	}
}

// create a new GameSettings object for the default game
func NewGameSettings() *GameSettings {
	return &GameSettings{
		Name:            "A Barefoot Jaywalk",
		Size:            SizeSmall,
		Density:         DensityNormal,
		PlayerPositions: PlayerPositionsModerate,
		RandomEvents:    true,
		StartMode:       GameStartModeNormal,
		VictoryConditions: VictoryConditions{
			Conditions:               Bitmask(VictoryConditionOwnPlanets) | Bitmask(VictoryConditionAttainTechLevels) | Bitmask(VictoryConditionExceedsSecondPlaceScore),
			NumCriteriaRequired:      1,
			YearsPassed:              50,
			OwnPlanets:               60,
			AttainTechLevel:          22,
			AttainTechLevelNumFields: 4,
			ExceedsScore:             11000,
			ExceedsSecondPlaceScore:  100,
			ProductionCapacity:       100,
			OwnCapitalShips:          100,
			HighestScoreAfterYears:   100,
		},
	}
}

// set the name to this game
func (settings *GameSettings) WithName(name string) *GameSettings {
	settings.Name = name
	return settings
}

func (settings *GameSettings) WithSize(size Size) *GameSettings {
	settings.Size = size
	return settings
}

func (settings *GameSettings) WithDensity(density Density) *GameSettings {
	settings.Density = density
	return settings
}

func (settings *GameSettings) WithPublicPlayerScores(publicPlayerScores bool) *GameSettings {
	settings.PublicPlayerScores = publicPlayerScores
	return settings
}

// add a host to this game
func (settings *GameSettings) WithHost(raceID int64) *GameSettings {
	settings.Players = append(settings.Players, NewGamePlayer{Type: NewGamePlayerTypeHost, RaceID: raceID, Color: "#0000FF"})
	return settings
}

// Add a player slot open to any players
func (settings *GameSettings) WithOpenPlayerSlot() *GameSettings {
	settings.Players = append(settings.Players, NewGamePlayer{Type: NewGamePlayerTypeOpen})
	return settings
}

// Add an AI player
func (settings *GameSettings) WithAIPlayer(aiDifficulty AIDifficulty, defaultHullSet int) *GameSettings {
	settings.Players = append(settings.Players, NewGamePlayer{Type: NewGamePlayerTypeAI, AIDifficulty: aiDifficulty, DefaultHullSet: defaultHullSet})
	return settings
}

func (settings *GameSettings) WithAIPlayerRace(race Race, aiDifficulty AIDifficulty, defaultHullSet int) *GameSettings {
	settings.Players = append(settings.Players, NewGamePlayer{Type: NewGamePlayerTypeAI, AIDifficulty: aiDifficulty, DefaultHullSet: defaultHullSet, Race: race})
	return settings
}

func (g *Game) String() string {
	return fmt.Sprintf("%s (%d)", g.Name, g.ID)
}

// update the game with new settings
func (g *Game) WithSettings(settings GameSettings) *Game {
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

	// use custom rules for this game
	if settings.Rules != nil {
		g.Rules = *settings.Rules
	}

	// use custome tech store for this game
	if settings.TechStore != nil {
		g.Rules.WithTechStore(settings.TechStore)
	}

	return g
}

// generate an invite hash for this game
func (g *Game) GenerateHash(salt string) string {
	hasher := sha1.New()
	hasher.Write([]byte(fmt.Sprintf("%d-%s", g.ID, salt)))
	sha := hex.EncodeToString(hasher.Sum(nil))
	return sha[10:]
}

func (g *Game) YearsPassed() int {
	return g.Year - g.Rules.StartingYear
}

func (fg *FullGame) getPlayer(playerNum int) *Player {

	if playerNum < 1 || playerNum > len(fg.Players)+1 {
		return nil
	}
	return fg.Players[playerNum-1]
}

func (g *FullGame) computeSpecs() error {

	g.buildMaps(g.Players)

	rules := &g.Rules
	for _, player := range g.Players {
		player.Race.Spec = computeRaceSpec(&player.Race, rules)
		player.Spec = computePlayerSpec(player, rules, g.Planets)
		log.Debug().Msgf("computed race and player spec for %v %s", player, player.Race.PluralName)

		for _, design := range player.Designs {
			numBuilt := design.Spec.NumBuilt
			design.Spec = ComputeShipDesignSpec(rules, player.TechLevels, player.Race.Spec, design)
			design.Spec.NumBuilt = numBuilt
			design.MarkDirty()
			log.Debug().Msgf("computed design spec for player %d, design %s", player.Num, design.Name)
		}
	}

	for _, starbase := range g.Starbases {
		player := g.getPlayer(starbase.PlayerNum)
		starbase.Spec = ComputeFleetSpec(rules, player, starbase)

		for _, token := range starbase.Tokens {
			design := g.designsByNum[playerObjectKey(starbase.PlayerNum, token.DesignNum)]
			design.Spec.NumInstances += token.Quantity
		}
		starbase.MarkDirty()
		log.Debug().Msgf("computed starbase spec for player %d, fleet %s", player.Num, starbase.Name)
	}

	for _, planet := range g.Planets {
		if planet.Owned() {
			player := g.getPlayer(planet.PlayerNum)
			planet.Spec = computePlanetSpec(rules, player, planet)
			if err := planet.PopulateProductionQueueDesigns(player); err != nil {
				return fmt.Errorf("planet %s unable to populate queue designs %w", planet.Name, err)
			}
			if err := planet.PopulateProductionQueueCosts(player); err != nil {
				return fmt.Errorf("planet %s unable to populate queue costs %w", planet.Name, err)
			}
			planet.PopulateProductionQueueEstimates(rules, player)

			planet.MarkDirty()
			log.Debug().Msgf("computed planet spec for player %d, planet %s", player.Num, planet.Name)
		}
	}

	for _, fleet := range g.Fleets {
		player := g.getPlayer(fleet.PlayerNum)
		fleet.Spec = ComputeFleetSpec(rules, player, fleet)

		for _, token := range fleet.Tokens {
			design := g.designsByNum[playerObjectKey(fleet.PlayerNum, token.DesignNum)]
			design.Spec.NumInstances += token.Quantity
		}
		fleet.MarkDirty()
		log.Debug().Msgf("computed fleet spec for player %d, fleet %s", player.Num, fleet.Name)
	}

	for _, mineField := range g.MineFields {
		player := g.getPlayer(mineField.PlayerNum)
		mineField.Spec = computeMinefieldSpec(rules, player, mineField, g.numPlanetsWithin(mineField.Position, mineField.Radius()))
		mineField.MarkDirty()
		log.Debug().Msgf("computed mineField spec for player %d, mineField %s", player.Num, mineField.Name)
	}

	for _, wormhole := range g.Wormholes {
		wormhole.Spec = computeWormholeSpec(wormhole, rules)
		wormhole.MarkDirty()
		log.Debug().Msgf("computed wormhole spec %s", wormhole.Name)
	}

	return nil

}
