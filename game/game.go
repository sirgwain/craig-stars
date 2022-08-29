package game

import (
	"fmt"
	"time"
)

type Tags map[string]string

type NewGamePlayerType string

const (
	NewGamePlayerTypeHost   NewGamePlayerType = "Host"
	NewGamePlayerTypeInvite NewGamePlayerType = "Invite"
	NewGamePlayerTypeOpen   NewGamePlayerType = "Open"
	NewGamePlayerTypeAI     NewGamePlayerType = "AI"
)

type AIDifficulty string

const (
	AIDifficultyNone   AIDifficulty = ""
	AIDifficultyEasy   AIDifficulty = "Easy"
	AIDifficultyNormal AIDifficulty = "Normal"
	AIDifficultyHard   AIDifficulty = "Hard"
)

type NewGamePlayer struct {
	Type         NewGamePlayerType `json:"type,omitempty"`
	RaceID       uint              `json:"raceID,omitempty"`
	AIDifficulty AIDifficulty      `json:"aiDifficulty,omitempty"`
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
	ID                           uint                      `gorm:"primaryKey" json:"id" header:"ID"`
	CreatedAt                    time.Time                 `json:"createdAt"`
	UpdatedAt                    time.Time                 `json:"updatedAt"`
	Name                         string                    `json:"name" header:"Name"`
	HostID                       uint                      `json:"hostId"`
	QuickStartTurns              int                       `json:"quickStartTurns"`
	Size                         Size                      `json:"size"`
	Density                      Density                   `json:"density"`
	PlayerPositions              PlayerPositions           `json:"playerPositions"`
	RandomEvents                 bool                      `json:"randomEvents"`
	ComputerPlayersFormAlliances bool                      `json:"computerPlayersFormAlliances"`
	PublicPlayerScores           bool                      `json:"publicPlayerScores"`
	StartMode                    GameStartMode             `json:"startMode"`
	Year                         int                       `json:"year"`
	State                        GameState                 `json:"state"`
	OpenPlayerSlots              uint                      `json:"openPlayerSlots"`
	NumPlayers                   int                       `json:"numPlayers"`
	VictoryConditions            VictoryConditions         `json:"victoryConditions" gorm:"embedded;embeddedPrefix:victory_condition_"`
	VictorDeclared               bool                      `json:"victorDeclared"`
	Area                         Vector                    `json:"area,omitempty" gorm:"embedded;embeddedPrefix:area_"`
	Players                      []Player                  `json:"players,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	Planets                      []Planet                  `json:"planets,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	Fleets                       []Fleet                   `json:"fleets,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	Wormholes                    []Wormohole               `json:"wormholes,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	MineralPackets               []MineralPacket           `json:"mineralPackets,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	Salvage                      []Salvage                 `json:"salvage,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	FleetsByNum                  map[playerFleetNum]*Fleet `json:"-" gorm:"-"`
	Rules                        Rules                     `json:"rules" gorm:"constraint:OnDelete:CASCADE;"`
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

// add a host to this game
func (settings *GameSettings) WithHost(raceID uint) *GameSettings {
	settings.Players = append(settings.Players, NewGamePlayer{Type: NewGamePlayerTypeHost, RaceID: raceID})
	return settings
}

// Add a player slot open to any players
func (settings *GameSettings) WithOpenPlayerSlot() *GameSettings {
	settings.Players = append(settings.Players, NewGamePlayer{Type: NewGamePlayerTypeOpen})
	return settings
}

// Add an AI player
func (settings *GameSettings) WithAIPlayer(aiDifficulty AIDifficulty) *GameSettings {
	settings.Players = append(settings.Players, NewGamePlayer{Type: NewGamePlayerTypeAI, AIDifficulty: aiDifficulty})
	return settings
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
	p.GameID = g.ID
	g.Players = append(g.Players, *p)
	g.NumPlayers = len(g.Players)
	g.OpenPlayerSlots = uint(clamp(int(g.OpenPlayerSlots-1), 0, g.NumPlayers))

	return &g.Players[len(g.Players)-1]
}

// // remove a fleet from the fleets list
// func (g *Game) removeFleet(fleet *Fleet) {
// 	key := playerFleetNum{fleet.PlayerNum, fleet.Num}
// 	delete(g.FleetsByNum, key)
// 	fleets := make([]Fleet, 0, len(g.Fleets)-1)
// 	for i := range g.Fleets {
// 		if &g.Fleets[i] != fleet {
// 			fleets = append(g.Fleets, g.Fleets[i])
// 		}
// 	}
// 	g.Fleets = fleets
// }

// update a player's planets, fleets, mineralpackets, etc
func (g *Game) updatePlayerOwnedObjects(player *Player) {
	// update player arrays of fleets and planets before we scan
	player.Fleets = []*Fleet{}
	for i := range g.Fleets {
		gameFleet := &g.Fleets[i]
		if gameFleet.PlayerNum == player.Num {
			player.Fleets = append(player.Fleets, gameFleet)
		}
	}

	player.Planets = []*Planet{}
	for i := range g.Planets {
		gamePlanet := &g.Planets[i]
		if gamePlanet.PlayerNum == player.Num {
			player.Planets = append(player.Planets, gamePlanet)
		}
	}

	player.MineralPackets = []*MineralPacket{}
	for i := range g.MineralPackets {
		gameMineralPacket := &g.MineralPackets[i]
		if gameMineralPacket.PlayerNum == player.Num {
			player.MineralPackets = append(player.MineralPackets, gameMineralPacket)
		}
	}

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
			player := &g.Players[planet.PlayerNum]
			planet.Spec = ComputePlanetSpec(&g.Rules, planet, player)
		}
	}
}

// Generate a new universe
func (g *Game) GenerateUniverse() error {
	ug := NewUniverseGenerator(g.Size, g.Density, g.Players, &g.Rules)
	universe, err := ug.Generate()

	if err != nil {
		return err
	}

	// todo save universe separately
	g.Area = universe.Area
	g.Planets = universe.Planets
	g.Fleets = universe.Fleets
	g.Wormholes = universe.Wormholes
	g.MineralPackets = universe.MineralPackets
	g.Salvage = universe.Salvage

	return nil
}

// generate a new turn
func (g *Game) GenerateTurn() error {
	turnGenerator := NewTurnGenerator(g)
	return turnGenerator.generateTurn()
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

func (g *Game) getOwnedPlanets() []Planet {
	var ownedPlanets []Planet

	for _, p := range g.Planets {
		if p.Owned() {
			ownedPlanets = append(ownedPlanets, p)
		}
	}

	return ownedPlanets
}

// transfer cargo from one cargo holder to another
func (g *Game) Transfer(source *Fleet, dest CargoHolder, cargoType CargoType, transferAmount int) {
	source.TransferCargoItem(dest, cargoType, transferAmount)

	// if (cargoType == CargoType.Fuel)	{
	// 	cargoTransferer.Transfer(source, dest, Cargo.Empty, transferAmount);
	// }	else if (cargoType == CargoType.Colonists)	{
	// 	// invasion?
	// 	if (dest is Planet planet && planet.PlayerNum != source.PlayerNum)		{
	// 		if (transferAmount > 0)			{
	// 			invasions.Add(new PlanetInvasion()
	// 			{
	// 				Planet = planet,
	// 				Fleet = source,
	// 				ColonistsToDrop = transferAmount * 100
	// 			});
	// 			// remove colonists from our cargo
	// 			source.Cargo = source.Cargo - Cargo.OfAmount(cargoType, transferAmount);
	// 		}
	// 		else			{
	// 			// can't beam enemy colonists onto your ship...
	// 			// TODO: send a message
	// 			log.Warn($"{Game.Year}: {source.PlayerNum} {source.Name} tried to beam colonists up from: {dest}");
	// 		}
	// 	}		else if (dest is Fleet otherFleet && otherFleet.PlayerNum != source.PlayerNum)		{
	// 		// ignore this, but send a message
	// 		// TODO: send a message
	// 		log.Warn($"{Game.Year}: {source.PlayerNum} {source.Name} tried to transfer colonists to/from a fleet they don't own: {otherFleet}");
	// 	}		else		{
	// 		cargoTransferer.Transfer(source, dest, Cargo.OfAmount(cargoType, transferAmount), 0);
	// 		log.Debug($"{Game.Year}: {source.PlayerNum} {source.Name} transferred {transferAmount}kT of {cargoType} to {dest.Name}");
	// 	}
	// }	else	{
	// 	// if this is a planet that is owned by someone else and we don't have "steal cargo from planets" ability in this fleet, make sure we are only giving cargo, not taking
	// 	if (dest is Planet planet && !planet.OwnedBy(source.PlayerNum) && !cargoTransferer.GetCanStealPlanetCargo(source, Game.MapObjectsByLocation))		{
	// 		transferAmount = Math.Max(0, transferAmount);
	// 	}

	// 	// if this is a fleet that is owned by someone else and we don't have "steal cargo from fleets" ability in this fleet, make sre we are only giving cargo, not taking it
	// 	if (dest is Fleet fleet && !fleet.OwnedBy(source.PlayerNum) && !cargoTransferer.GetCanStealFleetCargo(source, Game.MapObjectsByLocation))		{
	// 		transferAmount = Math.Max(0, transferAmount);
	// 	}

	// 	if (transferAmount != 0)		{
	// 		cargoTransferer.Transfer(source, dest, Cargo.OfAmount(cargoType, transferAmount), 0);
	// 		log.Debug($"{Game.Year}: {source.PlayerNum} {source.Name} transferred {transferAmount}kT of {cargoType} to {dest.Name}");
	// 	}
	// }
}
