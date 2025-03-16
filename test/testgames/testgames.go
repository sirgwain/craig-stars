package testgames

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
)

type TestGame struct {
	Name           string       `json:"name"`
	Players        []TestPlayer `json:"players"`
	Planets        []cs.Planet  `json:"planets"`
	Wormholes      []cs.Wormhole
	MysteryTraders []cs.MysteryTrader
}

type TestPlayer struct {
	// if nil, this will be created as a default humanoid player
	*cs.Player
	Designs        []cs.ShipDesign    `json:"designs,omitempty"`
	Fleets         []cs.Fleet         `json:"fleets,omitempty"`
	Salvages       []cs.Salvage       `json:"salvages,omitempty"`
	MineralPackets []cs.MineralPacket `json:"mineralPackets,omitempty"`
	MineFields     []cs.MineField     `json:"mineFields,omitempty"`
}

var TestGames = []TestGame{
	{
		Name: "Single Unit Game",
		Players: []TestPlayer{
			{
				Designs: []cs.ShipDesign{{
					Name:  "Long Range Scout",
					Hull:  cs.Scout.Name,
					Slots: longRangeScoutSlots,
				}},
				Fleets: []cs.Fleet{
					{
						BaseName:          "Long Range Scout",
						Tokens:            []cs.ShipToken{{DesignNum: 1, Quantity: 1}},
						OrbitingPlanetNum: 1,
					},
				},
			},
		},
		Planets: []cs.Planet{{
			MapObject: cs.MapObject{
				Name:      "Planet 1",
				PlayerNum: 1,
			},
			Hab:                  cs.Hab{Grav: 50, Temp: 50, Rad: 50},
			MineralConcentration: cs.NewMineral(100, 100, 100),
			Cargo:                cs.Cargo{Ironium: 1000, Boranium: 1000, Germanium: 1000, Colonists: 2500},
		}},
	},
	{
		Name: "Cargo Transfer Invasion",
		Players: []TestPlayer{
			{
				Designs: []cs.ShipDesign{
					{
						Name:  "Teamster",
						Hull:  cs.MediumFreighter.Name,
						Slots: teamsterSlots,
					},
				},
				Fleets: []cs.Fleet{
					{
						BaseName:          "Teamster",
						Tokens:            []cs.ShipToken{{DesignNum: 1, Quantity: 1}},
						Fuel:              450,
						Cargo:             cs.Cargo{Colonists: 210},
						OrbitingPlanetNum: 1,
					},
				},
			},
			{
				Player: &cs.Player{Name: "Player 2", AIControlled: true, Race: *cs.NewRace()},
			},
		},
		Planets: []cs.Planet{{
			MapObject: cs.MapObject{
				Name:      "Planet 1",
				PlayerNum: 2, // give player2 a planet
			},
			Hab:   cs.Hab{Grav: 50, Temp: 50, Rad: 50},
			Cargo: cs.Cargo{Ironium: 1000, Boranium: 1000, Germanium: 1000, Colonists: 25},
		}},
	},
	{
		Name: "Cargo Transfer Planet Owned",
		Players: []TestPlayer{
			{
				Designs: []cs.ShipDesign{
					{
						Name:  "Teamster",
						Hull:  cs.MediumFreighter.Name,
						Slots: teamsterSlots,
					},
				},
				Fleets: []cs.Fleet{
					{
						BaseName:          "Teamster",
						Tokens:            []cs.ShipToken{{DesignNum: 1, Quantity: 1}},
						Fuel:              450,
						OrbitingPlanetNum: 1,
					},
				},
			},
		},
		Planets: []cs.Planet{{
			MapObject: cs.MapObject{
				Name:      "Planet 1",
				PlayerNum: 1,
			},
			Hab:   cs.Hab{Grav: 50, Temp: 50, Rad: 50},
			Cargo: cs.Cargo{Ironium: 1000, Boranium: 1000, Germanium: 1000, Colonists: 2500},
		}},
	},
	{
		Name: "Cargo Transfer Planet Steal",
		Players: []TestPlayer{
			{
				Player: &cs.Player{Name: "Player 1", UserID: 1, Race: *cs.NewRace().WithPRT(cs.SS)},
				Designs: []cs.ShipDesign{
					{
						Name: "Thief",
						Hull: cs.MediumFreighter.Name,
						Slots: []cs.ShipDesignSlot{
							{HullComponent: cs.LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
							{HullComponent: cs.Crobmnium.Name, HullSlotIndex: 2, Quantity: 1},
							{HullComponent: cs.RobberBaronScanner.Name, HullSlotIndex: 3, Quantity: 1},
						},
					},
					{
						Name:  "Teamster",
						Hull:  cs.MediumFreighter.Name,
						Slots: teamsterSlots,
					},
				},
				Fleets: []cs.Fleet{
					{
						BaseName:          "Thief", // can steal
						Tokens:            []cs.ShipToken{{DesignNum: 1, Quantity: 1}},
						Fuel:              450,
						OrbitingPlanetNum: 1,
					},
					{
						BaseName:          "Teamster", // can't steal
						Tokens:            []cs.ShipToken{{DesignNum: 2, Quantity: 1}},
						Fuel:              450,
						OrbitingPlanetNum: 1,
					},
				},
			},
			{
				Player: &cs.Player{Name: "Player 2", AIControlled: true, Race: *cs.NewRace()},
			},
		},
		Planets: []cs.Planet{{
			MapObject: cs.MapObject{
				Name:      "Planet 1",
				PlayerNum: 2, // give player2 a planet
			},
			Hab:   cs.Hab{Grav: 50, Temp: 50, Rad: 50},
			Cargo: cs.Cargo{Ironium: 1000, Boranium: 1000, Germanium: 1000, Colonists: 2500},
		}},
	},
	{
		Name: "Cargo Transfer Fleet Steal",
		Players: []TestPlayer{
			{
				Player: &cs.Player{Name: "Player 1", UserID: 1, Race: *cs.NewRace().WithPRT(cs.SS)},
				Designs: []cs.ShipDesign{
					{
						Name: "Thief",
						Hull: cs.MediumFreighter.Name,
						Slots: []cs.ShipDesignSlot{
							{HullComponent: cs.LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
							{HullComponent: cs.Crobmnium.Name, HullSlotIndex: 2, Quantity: 1},
							{HullComponent: cs.RobberBaronScanner.Name, HullSlotIndex: 3, Quantity: 1},
						},
					},
					{
						Name:  "Teamster",
						Hull:  cs.MediumFreighter.Name,
						Slots: teamsterSlots,
					},
				},
				Fleets: []cs.Fleet{
					{
						BaseName: "Thief", // can steal
						Tokens:   []cs.ShipToken{{DesignNum: 1, Quantity: 1}},
						Fuel:     100,
					},
					{
						BaseName: "Teamster", // can't steal
						Tokens:   []cs.ShipToken{{DesignNum: 2, Quantity: 1}},
						Fuel:     100,
					},
				},
			},
			{
				Player: &cs.Player{Name: "Player 2", AIControlled: true, Race: *cs.NewRace()},
				Designs: []cs.ShipDesign{
					{
						Name:  "Teamster",
						Hull:  cs.MediumFreighter.Name,
						Slots: teamsterSlots,
					},
				},
				Fleets: []cs.Fleet{
					{
						BaseName: "Teamster",
						Tokens:   []cs.ShipToken{{DesignNum: 1, Quantity: 1}},
						Cargo:    cs.Cargo{Ironium: 10, Boranium: 10, Germanium: 10, Colonists: 10},
						Fuel:     10,
					},
				},
			},
		},
	},
	{
		Name: "Cargo Transfer Jettison",
		Players: []TestPlayer{
			{
				Designs: []cs.ShipDesign{
					{
						Name:  "Teamster",
						Hull:  cs.MediumFreighter.Name,
						Slots: teamsterSlots,
					},
				},
				Fleets: []cs.Fleet{
					{
						BaseName: "Teamster Jettison",
						Tokens:   []cs.ShipToken{{DesignNum: 1, Quantity: 1}},
						Cargo:    cs.Cargo{Ironium: 50, Boranium: 50, Germanium: 50},
						Fuel:     500,
					},
				},
			},
		},
	},
	{
		Name: "Cargo Transfer Salvage",
		Players: []TestPlayer{
			{
				Designs: []cs.ShipDesign{
					{
						Name:  "Teamster",
						Hull:  cs.MediumFreighter.Name,
						Slots: teamsterSlots,
					},
				},
				Fleets: []cs.Fleet{
					{
						BaseName: "Teamster Salvager",
						Tokens:   []cs.ShipToken{{DesignNum: 1, Quantity: 1}},
						Cargo:    cs.Cargo{Ironium: 10, Boranium: 10, Germanium: 10},
						Fuel:     500,
					},
				},
				Salvages: []cs.Salvage{
					{
						Cargo: cs.Cargo{Ironium: 50, Boranium: 50, Germanium: 50},
					},
				},
			},
		},
	},
	{
		Name: "Cargo Transfer Fleets",
		Players: []TestPlayer{
			{
				Designs: []cs.ShipDesign{
					{
						Name:  "Teamster",
						Hull:  cs.MediumFreighter.Name,
						Slots: teamsterSlots,
					},
					{
						Name:  "Santa Maria",
						Hull:  cs.ColonyShip.Name,
						Slots: santaMariaSlots,
					},
				},
				Fleets: []cs.Fleet{
					{
						BaseName: "Teamster",
						Tokens:   []cs.ShipToken{{DesignNum: 1, Quantity: 1}},
						Cargo:    cs.Cargo{Ironium: 10, Boranium: 10, Germanium: 10, Colonists: 10},
						Fuel:     100,
					},
					{
						BaseName: "Santa Maria",
						Tokens:   []cs.ShipToken{{DesignNum: 2, Quantity: 1}},
						Cargo:    cs.Cargo{Ironium: 5, Boranium: 5, Germanium: 5, Colonists: 5},
						Fuel:     10,
					},
				},
			},
		},
	},
	{
		Name: "Cargo Transfer Split",
		Players: []TestPlayer{
			{
				Designs: []cs.ShipDesign{
					{
						Name:  "Teamster",
						Hull:  cs.MediumFreighter.Name,
						Slots: teamsterSlots,
					},
					{
						Name:  "Santa Maria",
						Hull:  cs.ColonyShip.Name,
						Slots: santaMariaSlots,
					},
				},
				Fleets: []cs.Fleet{
					{
						BaseName: "Teamster Jettison",
						Tokens:   []cs.ShipToken{{DesignNum: 1, Quantity: 1}, {DesignNum: 2, Quantity: 1}},
						Cargo:    cs.Cargo{Ironium: 50, Boranium: 50, Germanium: 50},
						Fuel:     500,
					},
				},
			},
		},
	},
	{
		Name: "Cargo Transfer MineralPacket",
		Players: []TestPlayer{
			{
				Designs: []cs.ShipDesign{
					{
						Name:  "Teamster",
						Hull:  cs.MediumFreighter.Name,
						Slots: teamsterSlots,
					},
				},
				Fleets: []cs.Fleet{
					{
						BaseName: "Teamster Mineral Packeter",
						Tokens:   []cs.ShipToken{{DesignNum: 1, Quantity: 1}},
						Cargo:    cs.Cargo{Ironium: 10, Boranium: 10, Germanium: 10, Colonists: 10},
						Fuel:     500,
					},
				},
				MineralPackets: []cs.MineralPacket{
					{
						Cargo:           cs.Cargo{Ironium: 50, Boranium: 50, Germanium: 50},
						WarpSpeed:       5,
						TargetPlanetNum: 1,
						Heading:         cs.Vector{X: 0, Y: 1},
					},
				},
			},
		},
	},
	{
		Name: "Two Player Game",
		Players: []TestPlayer{
			{
				Designs: []cs.ShipDesign{
					{
						Name:  "Long Range Scout",
						Hull:  cs.Scout.Name,
						Slots: longRangeScoutSlots,
					},
				},
				Fleets: []cs.Fleet{
					{
						BaseName: "Long Range Scout",
						Tokens:   []cs.ShipToken{{DesignNum: 1, Quantity: 3}}, OrbitingPlanetNum: 1,
					},
				},
			},
			{
				Player: &cs.Player{Name: "Player 2", AIControlled: true, Race: *cs.NewRace()},
				Designs: []cs.ShipDesign{
					{
						Name:  "Super Scout",
						Hull:  cs.Scout.Name,
						Slots: longRangeScoutSlots,
					},
				},
				Fleets: []cs.Fleet{
					{
						MapObject:         cs.MapObject{Position: cs.Vector{X: 50}},
						BaseName:          "Long Range Scout",
						Tokens:            []cs.ShipToken{{DesignNum: 1, Quantity: 2}},
						OrbitingPlanetNum: 2,
					},
				},
			},
		},

		Planets: []cs.Planet{
			{
				MapObject: cs.MapObject{
					Name:      "Planet 1",
					PlayerNum: 1,
				},
				Hab:                  cs.Hab{Grav: 50, Temp: 50, Rad: 50},
				MineralConcentration: cs.NewMineral(100, 100, 100),
				Cargo:                cs.Cargo{Ironium: 1000, Boranium: 1000, Germanium: 1000, Colonists: 2500},
			},
			{
				MapObject: cs.MapObject{
					Name:      "Planet 2",
					PlayerNum: 2,
					Position:  cs.Vector{X: 50},
				},
				Hab:                  cs.Hab{Grav: 50, Temp: 50, Rad: 50},
				MineralConcentration: cs.NewMineral(100, 100, 100),
				Cargo:                cs.Cargo{Ironium: 1000, Boranium: 1000, Germanium: 1000, Colonists: 2500},
			}},
	},
}

func CreateTestGames(db db.Client) error {
	for _, testGame := range TestGames {
		game := createTestGame(testGame)
		if err := db.CreateGame(game.Game); err != nil {
			return err
		}
		for _, player := range game.Players {
			player.GameID = game.ID
			if err := db.CreatePlayer(player); err != nil {
				return err
			}
			for _, design := range player.Designs {
				design.GameID = game.ID
			}
		}

		// do all new game gen stuff required
		ug := cs.NewUniverseGenerator(game.Game, game.Players)
		if err := ug.GenerateWithUniverse(game.Universe); err != nil {
			return err
		}

		// save to db
		if err := db.UpdateFullGame(game); err != nil {
			return err
		}
	}

	return nil
}

var colors = []string{
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

func createTestGame(tg TestGame) *cs.FullGame {
	game := newGame(tg.Name)

	for _, p := range tg.Planets {
		addPlanet(game, &p)
	}

	for _, wh := range tg.Wormholes {
		addWormhole(game, &wh)
	}

	for _, mt := range tg.MysteryTraders {
		addMysteryTrader(game, &mt)
	}

	for i, p := range tg.Players {
		if p.Player == nil {
			p.Player = &cs.Player{
				UserID:     1,
				Race:       *cs.NewRace(),
				TechLevels: cs.TechLevel{Energy: 3, Weapons: 3, Propulsion: 3, Construction: 3, Electronics: 3, Biotechnology: 3},
			}
		}
		player := addPlayer(game, p.WithNum(i+1))
		if p.AIControlled {
			player.SubmittedTurn = true
		}

		for _, d := range p.Designs {
			addDesign(player, &cs.ShipDesign{
				Name:  d.Name,
				Hull:  d.Hull,
				Slots: d.Slots,
			})
		}

		for i, f := range p.Fleets {
			addFleet(game, player, &f, i+1)
		}

		for i, mp := range p.MineralPackets {
			addMineralPacket(game, player, &mp, i+1)
		}

		for i, mf := range p.MineFields {
			addMineField(game, player, &mf, i+1)
		}

		for _, s := range p.Salvages {
			addSalvage(game, player, &s)
		}

	}

	return game
}

func newGame(name string) *cs.FullGame {
	client := cs.NewGamer()
	game := client.CreateGame(1, *cs.NewGameSettings().WithName(name))
	game.RandomEvents = false // don't allow random events in tests unless configured
	game.Area = cs.Vector{X: 200, Y: 200}
	game.Rules.ResetSeed(0) // keep the same seed for tests
	game.State = cs.GameStateWaitingForPlayers
	universe := cs.NewUniverse(log.Logger, &game.Rules)

	return &cs.FullGame{
		Game:      game,
		Universe:  &universe,
		TechStore: &cs.StaticTechStore,
		Players:   []*cs.Player{},
	}
}

func addPlayer(game *cs.FullGame, player *cs.Player) *cs.Player {
	player.Num = len(game.Players) + 1
	player.Color = colors[player.Num-1]

	if player.Name == "" {
		player.Name = fmt.Sprintf("Player #%d", player.Num)
	}

	game.Players = append(game.Players, player)
	return player
}

func addPlanet(game *cs.FullGame, planet *cs.Planet) *cs.Planet {
	planet.Type = cs.MapObjectTypePlanet
	planet.Num = len(game.Planets) + 1
	game.Planets = append(game.Planets, planet)
	return planet
}

func addDesign(player *cs.Player, design *cs.ShipDesign) *cs.ShipDesign {
	design.Num = len(player.Designs) + 1
	design.PlayerNum = player.Num
	player.Designs = append(player.Designs, design)
	return design
}

func addFleet(game *cs.FullGame, player *cs.Player, fleet *cs.Fleet, num int) *cs.Fleet {
	fleet.Type = cs.MapObjectTypeFleet
	fleet.Num = num
	fleet.PlayerNum = player.Num
	fleet.Name = fmt.Sprintf("%s #%d", fleet.BaseName, fleet.Num)
	fleet.Waypoints = []cs.Waypoint{
		cs.NewPositionWaypoint(fleet.Position, 0),
	}

	if fleet.OrbitingPlanetNum != cs.None {
		planet := game.Planets[fleet.OrbitingPlanetNum-1]
		fleet.Position = planet.Position
		fleet.Waypoints = []cs.Waypoint{
			cs.NewPlanetWaypoint(planet.Position, planet.Num, planet.Name, 0),
		}
	}

	game.Fleets = append(game.Fleets, fleet)
	return fleet
}

func addMineralPacket(game *cs.FullGame, player *cs.Player, mineralPacket *cs.MineralPacket, num int) *cs.MineralPacket {
	mineralPacket.Type = cs.MapObjectTypeMineralPacket
	mineralPacket.PlayerNum = player.Num
	mineralPacket.Num = num
	mineralPacket.Name = fmt.Sprintf("%s Mineral Packet #%d", player.Race.PluralName, mineralPacket.Num)
	game.MineralPackets = append(game.MineralPackets, mineralPacket)
	return mineralPacket
}

func addMineField(game *cs.FullGame, player *cs.Player, mineField *cs.MineField, num int) *cs.MineField {
	mineField.Type = cs.MapObjectTypeMineField
	mineField.PlayerNum = player.Num
	mineField.Num = len(game.MineFields) + 1
	mineField.Name = fmt.Sprintf("%s %s Mine Field #%d", player.Race.PluralName, mineField.MineFieldType.String(), num)
	game.MineFields = append(game.MineFields, mineField)
	return mineField
}

func addSalvage(game *cs.FullGame, player *cs.Player, salvage *cs.Salvage) *cs.Salvage {
	salvage.Type = cs.MapObjectTypeSalvage
	salvage.PlayerNum = player.Num
	salvage.Num = len(game.Salvages) + 1
	salvage.Name = fmt.Sprintf("Salvage #%d", salvage.Num)
	game.Salvages = append(game.Salvages, salvage)
	return salvage
}

func addWormhole(game *cs.FullGame, wormhole *cs.Wormhole) *cs.Wormhole {
	wormhole.Type = cs.MapObjectTypeWormhole
	wormhole.Num = len(game.Wormholes) + 1
	game.Wormholes = append(game.Wormholes, wormhole)
	return wormhole
}

func addMysteryTrader(game *cs.FullGame, mysteryTrader *cs.MysteryTrader) *cs.MysteryTrader {
	mysteryTrader.Type = cs.MapObjectTypePlanet
	mysteryTrader.Num = len(game.MysteryTraders) + 1
	mysteryTrader.Name = fmt.Sprintf("Mystery Trader #%d", mysteryTrader.Num)
	game.MysteryTraders = append(game.MysteryTraders, mysteryTrader)
	return mysteryTrader
}
