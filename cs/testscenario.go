//go:build !wasi && !wasm

package cs

import (
	"fmt"
	"log/slog"
	"slices"
)

// TestScenario describes a small, deterministic game fixture for tests.
//
// BuildScenario expands the declarative objects into a FullGame and runs the
// normal universe-generation finalizer so specs, intel, maps, and related
// derived state are initialized the same way as a generated game.
type TestScenario struct {
	Name           string               `json:"name"`
	Players        []TestScenarioPlayer `json:"players"`
	Planets        []Planet             `json:"planets"`
	Wormholes      []Wormhole
	MysteryTraders []MysteryTrader
}

// TestScenarioPlayer describes the player-owned objects in a TestScenario.
//
// If Player is nil, BuildScenario creates a default humanoid player with enough
// starting tech and orders for turn-generation and e2e fixtures.
type TestScenarioPlayer struct {
	*Player
	Designs        []ShipDesign    `json:"designs,omitempty"`
	Fleets         []Fleet         `json:"fleets,omitempty"`
	Salvages       []Salvage       `json:"salvages,omitempty"`
	MineralPackets []MineralPacket `json:"mineralPackets,omitempty"`
	Minefields     []Minefield     `json:"minefields,omitempty"`
}

// LongRangeScoutTestSlots is the standard scout design used by small test scenarios.
var LongRangeScoutTestSlots = []ShipDesignSlot{
	{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
	{HullComponent: RhinoScanner.Name, HullSlotIndex: 2, Quantity: 1},
	{HullComponent: FuelTank.Name, HullSlotIndex: 3, Quantity: 1},
}

// SantaMariaTestSlots is the standard colony-ship design used by test scenarios.
var SantaMariaTestSlots = []ShipDesignSlot{
	{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
	{HullComponent: ColonizationModule.Name, HullSlotIndex: 2, Quantity: 1},
}

// SantaMariaARTestSlots is the Alternate Reality colony-ship design used by test scenarios.
var SantaMariaARTestSlots = []ShipDesignSlot{
	{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
	{HullComponent: OrbitalConstructionModule.Name, HullSlotIndex: 2, Quantity: 1},
}

// TeamsterTestSlots is the standard medium-freighter design used by cargo test scenarios.
var TeamsterTestSlots = []ShipDesignSlot{
	{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
	{HullComponent: Crobmnium.Name, HullSlotIndex: 2, Quantity: 1},
	{HullComponent: RhinoScanner.Name, HullSlotIndex: 3, Quantity: 1},
}

// DestroyerDeltaTestSlots is the standard destroyer design used by battle test scenarios.
var DestroyerDeltaTestSlots = []ShipDesignSlot{
	{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
	{HullComponent: DeltaTorpedo.Name, HullSlotIndex: 2, Quantity: 1},
	{HullComponent: DeltaTorpedo.Name, HullSlotIndex: 3, Quantity: 1},
	{HullComponent: DeltaTorpedo.Name, HullSlotIndex: 4, Quantity: 1},
	{HullComponent: Crobmnium.Name, HullSlotIndex: 5, Quantity: 1},
	{HullComponent: ManeuveringJet.Name, HullSlotIndex: 6, Quantity: 1},
	{HullComponent: BattleComputer.Name, HullSlotIndex: 7, Quantity: 1},
}

var scenarioColors = []string{
	"#0000FF",
	"#C33232",
	"#1F8BA7",
	"#43A43E",
	"#8D29CB",
	"#B88628",
	"#FF4500",
	"#FF8C00",
	"#008000",
	"#00FA9A",
	"#7FFFD4",
	"#8A2BE2",
	"#FF1493",
	"#D2691E",
	"#F0FFF0",
}

// SingleUnitScenario returns a one-player scenario with one planet and one scout.
func SingleUnitScenario() TestScenario {
	return TestScenario{
		Name: "Single Unit Game",
		Players: []TestScenarioPlayer{
			{
				Designs: []ShipDesign{{
					Name:  "Long Range Scout",
					Hull:  Scout.Name,
					Slots: LongRangeScoutTestSlots,
				}},
				Fleets: []Fleet{
					{
						Fuel:              300,
						BaseName:          "Long Range Scout",
						Tokens:            []ShipToken{{DesignNum: 1, Quantity: 1}},
						OrbitingPlanetNum: 1,
					},
				},
			},
		},
		Planets: []Planet{{
			MapObject: MapObject{
				Name:      "Planet 1",
				PlayerNum: 1,
			},
			Hab:                  Hab{Grav: 50, Temp: 50, Rad: 50},
			MineralConcentration: NewMineral(100, 100, 100),
			Cargo:                Cargo{Colonists: 2500},
		}},
	}
}

// TwoPlayerScenario returns a two-player scenario with one planet and one scout per player.
func TwoPlayerScenario() TestScenario {
	return TestScenario{
		Name: "Two Player Game",
		Players: []TestScenarioPlayer{
			{
				Designs: []ShipDesign{
					{
						Name:  "Long Range Scout",
						Hull:  Scout.Name,
						Slots: LongRangeScoutTestSlots,
					},
				},
				Fleets: []Fleet{
					{
						BaseName:          "Long Range Scout",
						Tokens:            []ShipToken{{DesignNum: 1, Quantity: 1}},
						OrbitingPlanetNum: 1,
					},
				},
			},
			{
				Player: &Player{Name: "Player 2", AIControlled: true, Race: *NewRace()},
				Designs: []ShipDesign{
					{
						Name:  "Super Scout",
						Hull:  Scout.Name,
						Slots: LongRangeScoutTestSlots,
					},
				},
				Fleets: []Fleet{
					{
						BaseName:          "Long Range Scout",
						Tokens:            []ShipToken{{DesignNum: 1, Quantity: 1}},
						OrbitingPlanetNum: 2,
					},
				},
			},
		},

		Planets: []Planet{
			{
				MapObject: MapObject{
					Name:      "Planet 1",
					PlayerNum: 1,
				},
				Hab:                  Hab{Grav: 50, Temp: 50, Rad: 50},
				MineralConcentration: NewMineral(100, 100, 100),
				Cargo:                Cargo{Colonists: 2500},
			},
			{
				MapObject: MapObject{
					Name:      "Planet 2",
					PlayerNum: 2,
					Position:  Vector{X: 100},
				},
				Hab:                  Hab{Grav: 50, Temp: 50, Rad: 50},
				MineralConcentration: NewMineral(100, 100, 100),
				Cargo:                Cargo{Colonists: 2500},
			}},
	}
}

// BuildScenario creates a fully initialized FullGame from a test scenario.
//
// BuildScenario panics if the scenario cannot be finalized by the universe
// generator, which keeps test setup call sites compact while still surfacing
// malformed fixtures immediately.
func BuildScenario(scenario TestScenario) *FullGame {
	game := newScenarioGame(scenario.Name)

	for _, planet := range scenario.Planets {
		addScenarioPlanet(game, &planet)
	}

	for _, wormhole := range scenario.Wormholes {
		addScenarioWormhole(game, &wormhole)
	}

	for _, mysteryTrader := range scenario.MysteryTraders {
		addScenarioMysteryTrader(game, &mysteryTrader)
	}

	for playerIndex, scenarioPlayer := range scenario.Players {
		playerTemplate := defaultScenarioPlayer()
		if scenarioPlayer.Player != nil {
			playerTemplate = *scenarioPlayer.Player
		}

		player := addScenarioPlayer(game, playerTemplate.WithNum(playerIndex+1))
		if player.AIControlled {
			player.SubmittedTurn = true
		}

		for _, design := range scenarioPlayer.Designs {
			addScenarioDesign(player, &ShipDesign{
				Name:    design.Name,
				Hull:    design.Hull,
				Purpose: design.Purpose,
				Slots:   slices.Clone(design.Slots),
			})
		}

		for fleetIndex, fleet := range scenarioPlayer.Fleets {
			addScenarioFleet(game, player, &fleet, fleetIndex+1)
		}

		for planetIndex, planet := range scenario.Planets {
			if planet.PlayerNum == player.Num && planet.Starbase != nil {
				addScenarioStarbase(game, player, planet.Starbase, game.Planets[planetIndex])
			}
		}

		for packetIndex, mineralPacket := range scenarioPlayer.MineralPackets {
			addScenarioMineralPacket(game, player, &mineralPacket, packetIndex+1)
		}

		for minefieldIndex, minefield := range scenarioPlayer.Minefields {
			addScenarioMinefield(game, player, &minefield, minefieldIndex+1)
		}

		for _, salvage := range scenarioPlayer.Salvages {
			addScenarioSalvage(game, player, &salvage)
		}
	}

	ug := NewUniverseGenerator(game.Game, game.Players)
	if err := ug.GenerateWithUniverse(game.Universe); err != nil {
		panic(fmt.Errorf("failed to generate universe for test scenario %v", err))
	}

	return game
}

func defaultScenarioPlayer() Player {
	return Player{
		UserID:        1,
		Race:          *NewRace(),
		TechLevels:    TechLevel{Energy: 3, Weapons: 3, Propulsion: 3, Construction: 3, Electronics: 3, Biotechnology: 3},
		Stats:         &PlayerStats{},
		AcquiredTechs: map[string]bool{},
		PlayerOrders: PlayerOrders{
			Researching:       Energy,
			ResearchAmount:    15,
			NextResearchField: NextResearchFieldSameField,
			CargoTransfers:    CargoTransfers{},
		},
	}
}

func newScenarioGame(name string) *FullGame {
	client := NewGamer()
	game := client.CreateGame(1, *NewGameSettings().WithName(name))
	game.RandomEvents = false
	game.Area = Vector{X: 200, Y: 200}
	game.Rules.ResetSeed(0)
	game.State = GameStateWaitingForPlayers
	universe := NewUniverse(slog.Default(), &game.Rules)

	return &FullGame{
		Game:      game,
		Universe:  &universe,
		TechStore: &StaticTechStore,
		Players:   []*Player{},
	}
}

func addScenarioPlayer(game *FullGame, player *Player) *Player {
	player.Num = len(game.Players) + 1
	player.Color = scenarioColors[player.Num-1]

	if player.Name == "" {
		player.Name = fmt.Sprintf("Player #%d", player.Num)
	}

	game.Players = append(game.Players, player)
	return player
}

func addScenarioPlanet(game *FullGame, planet *Planet) *Planet {
	planet.Type = MapObjectTypePlanet
	planet.Num = len(game.Planets) + 1
	if planet.BaseHab == (Hab{}) {
		planet.BaseHab = planet.Hab
	}
	game.Planets = append(game.Planets, planet)

	return planet
}

func addScenarioDesign(player *Player, design *ShipDesign) *ShipDesign {
	design.Num = len(player.Designs) + 1
	design.PlayerNum = player.Num
	player.Designs = append(player.Designs, design)
	return design
}

func addScenarioFleet(game *FullGame, player *Player, fleet *Fleet, num int) *Fleet {
	fleet.Type = MapObjectTypeFleet
	fleet.Num = num
	fleet.PlayerNum = player.Num
	fleet.Name = fmt.Sprintf("%s #%d", fleet.BaseName, fleet.Num)
	fleet.Waypoints = []Waypoint{
		NewPositionWaypoint(fleet.Position, 0),
	}

	if fleet.OrbitingPlanetNum != None {
		planet := game.Planets[fleet.OrbitingPlanetNum-1]
		fleet.Position = planet.Position
		fleet.Waypoints = []Waypoint{
			NewPlanetWaypoint(planet.Position, planet.Num, planet.Name, 0),
		}
	}

	game.Fleets = append(game.Fleets, fleet)
	return fleet
}

func addScenarioStarbase(game *FullGame, player *Player, fleet *Fleet, planet *Planet) *Fleet {
	fleet.Type = MapObjectTypeFleet
	fleet.PlayerNum = player.Num
	fleet.Name = fleet.BaseName
	fleet.Waypoints = []Waypoint{
		NewPositionWaypoint(fleet.Position, 0),
	}

	fleet.PlanetNum = planet.Num
	fleet.Starbase = true

	game.Starbases = append(game.Starbases, fleet)
	return fleet
}

func addScenarioMineralPacket(game *FullGame, player *Player, mineralPacket *MineralPacket, num int) *MineralPacket {
	mineralPacket.Type = MapObjectTypeMineralPacket
	mineralPacket.PlayerNum = player.Num
	mineralPacket.Num = num
	mineralPacket.Name = fmt.Sprintf("%s Mineral Packet #%d", player.Race.PluralName, mineralPacket.Num)
	mineralPacket.Heading = (game.Planets[mineralPacket.TargetPlanetNum-1].Position.Subtract(mineralPacket.Position)).Normalized()

	game.MineralPackets = append(game.MineralPackets, mineralPacket)
	return mineralPacket
}

func addScenarioMinefield(game *FullGame, player *Player, minefield *Minefield, num int) *Minefield {
	minefield.Type = MapObjectTypeMinefield
	minefield.PlayerNum = player.Num
	minefield.Num = len(game.Minefields) + 1
	minefield.Name = fmt.Sprintf("%s %s Minefield #%d", player.Race.PluralName, minefield.MinefieldType.String(), num)
	game.Minefields = append(game.Minefields, minefield)
	return minefield
}

func addScenarioSalvage(game *FullGame, player *Player, salvage *Salvage) *Salvage {
	salvage.Type = MapObjectTypeSalvage
	salvage.PlayerNum = player.Num
	salvage.Num = len(game.Salvages) + 1
	salvage.Name = fmt.Sprintf("Salvage #%d", salvage.Num)
	game.Salvages = append(game.Salvages, salvage)
	return salvage
}

func addScenarioWormhole(game *FullGame, wormhole *Wormhole) *Wormhole {
	wormhole.Type = MapObjectTypeWormhole
	wormhole.Num = len(game.Wormholes) + 1
	game.Wormholes = append(game.Wormholes, wormhole)
	return wormhole
}

func addScenarioMysteryTrader(game *FullGame, mysteryTrader *MysteryTrader) *MysteryTrader {
	mysteryTrader.Type = MapObjectTypePlanet
	mysteryTrader.Num = len(game.MysteryTraders) + 1
	mysteryTrader.Name = fmt.Sprintf("Mystery Trader #%d", mysteryTrader.Num)
	mysteryTrader.Heading = (mysteryTrader.Destination.Subtract(mysteryTrader.Position)).Normalized()
	mysteryTrader.PlayersRewarded = map[int]bool{}
	game.MysteryTraders = append(game.MysteryTraders, mysteryTrader)
	return mysteryTrader
}
