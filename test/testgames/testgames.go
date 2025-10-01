//go:build !wasi && !wasm

package testgames

import (
	"context"
	"fmt"
	"log/slog"

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
	Minefields     []cs.Minefield     `json:"minefields,omitempty"`
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
						Fuel:              300,
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
		Name: "Scout Test",
		Players: []TestPlayer{
			{
				Designs: []cs.ShipDesign{{
					Name:  "Long Range Scout",
					Hull:  cs.Scout.Name,
					Slots: longRangeScoutSlots,
				}},
				Fleets: []cs.Fleet{
					{
						Fuel:              300,
						BaseName:          "Long Range Scout",
						Tokens:            []cs.ShipToken{{DesignNum: 1, Quantity: 1}},
						OrbitingPlanetNum: 1,
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
				Homeworld:            true,
			},
			{
				MapObject: cs.MapObject{
					Name:     "Planet 2",
					Position: cs.Vector{X: 0, Y: 49}, // warp 7, one year
				},
				Hab:                  cs.Hab{Grav: 25, Temp: 25, Rad: 25},
				MineralConcentration: cs.NewMineral(100, 100, 100),
				Cargo:                cs.Cargo{},
			},
		},
	},
	{
		Name: "Colonizer Test",
		Players: []TestPlayer{
			{
				Designs: []cs.ShipDesign{
					{
						Name:  "Santa Maria",
						Hull:  cs.ColonyShip.Name,
						Slots: santaMariaSlots,
					},
					{
						Name:  "Long Range Scout",
						Hull:  cs.Scout.Name,
						Slots: longRangeScoutSlots,
					},
				},
				Fleets: []cs.Fleet{
					{
						Fuel:              200,
						BaseName:          "Santa Maria",
						Tokens:            []cs.ShipToken{{DesignNum: 1, Quantity: 1}},
						OrbitingPlanetNum: 1,
					},
					// scout for scanning
					{
						Fuel:              300,
						BaseName:          "Long Range Scout",
						Tokens:            []cs.ShipToken{{DesignNum: 2, Quantity: 1}},
						OrbitingPlanetNum: 1,
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
				Homeworld:            true,
			},
			{
				MapObject: cs.MapObject{
					Name:     "Planet 2",
					Position: cs.Vector{X: 0, Y: 25},
				},
				Hab:                  cs.Hab{Grav: 50, Temp: 50, Rad: 50},
				MineralConcentration: cs.NewMineral(100, 100, 100),
				Cargo:                cs.Cargo{},
			},
		},
	},
	{
		Name: "Kitchen Sink",
		Players: []TestPlayer{
			{
				Designs: []cs.ShipDesign{
					{
						Name:  "Long Range Scout",
						Hull:  cs.Scout.Name,
						Slots: longRangeScoutSlots,
					},
					{
						Name:  "Teamster",
						Hull:  cs.MediumFreighter.Name,
						Slots: teamsterSlots,
					},
				},
				Fleets: []cs.Fleet{
					{
						BaseName:          "Long Range Scout",
						Tokens:            []cs.ShipToken{{DesignNum: 1, Quantity: 1}},
						Fuel:              300,
						OrbitingPlanetNum: 1,
					},
					{
						MapObject: cs.MapObject{Position: cs.Vector{X: 0, Y: 50}},
						BaseName:  "Teamster",
						Tokens:    []cs.ShipToken{{DesignNum: 2, Quantity: 2}},
						Fuel:      900,
					},
				},
				MineralPackets: []cs.MineralPacket{
					{
						MapObject:       cs.MapObject{Position: cs.Vector{X: 50, Y: 0}},
						Cargo:           cs.Cargo{Ironium: 50, Boranium: 50, Germanium: 50},
						WarpSpeed:       5,
						TargetPlanetNum: 1,
					},
				},
				Minefields: []cs.Minefield{
					{
						MinefieldType: cs.MinefieldTypeStandard,
						NumMines:      1000,
					},
				},
				Salvages: []cs.Salvage{
					{
						MapObject: cs.MapObject{Position: cs.Vector{X: 50, Y: 50}},
						Cargo:     cs.Cargo{Ironium: 50, Boranium: 50, Germanium: 50},
					},
				},
			},
			{
				Player: &cs.Player{Name: "Player 2", AIControlled: true, Race: *cs.NewRace().WithPluralName("Rabbitoids")},
				Designs: []cs.ShipDesign{
					{
						Name:  "Teamster",
						Hull:  cs.MediumFreighter.Name,
						Slots: teamsterSlots,
					},
				},
				Fleets: []cs.Fleet{
					{
						MapObject: cs.MapObject{Position: cs.Vector{X: 30, Y: 0}},
						BaseName:  "Teamster",
						Tokens:    []cs.ShipToken{{DesignNum: 1, Quantity: 1}},
						Cargo:     cs.Cargo{Ironium: 10, Boranium: 10, Germanium: 10, Colonists: 10},
						Fuel:      10,
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
					Position:  cs.Vector{X: 0, Y: 30},
				},
				Hab:                  cs.Hab{Grav: 25, Temp: 25, Rad: 25},
				MineralConcentration: cs.NewMineral(100, 100, 100),
				Cargo:                cs.Cargo{Ironium: 1000, Boranium: 1000, Germanium: 1000, Colonists: 2500},
			},
			{
				MapObject: cs.MapObject{
					Name:     "Planet 3",
					Position: cs.Vector{X: 30, Y: 30},
				},
				Hab:                  cs.Hab{Grav: 75, Temp: 75, Rad: 75},
				MineralConcentration: cs.NewMineral(100, 100, 100),
			},
		},
		Wormholes: []cs.Wormhole{
			{
				MapObject:      cs.MapObject{Position: cs.Vector{X: 10, Y: 10}},
				Stability:      cs.WormholeStabilityRockSolid,
				DestinationNum: 2,
			},
			{
				MapObject:      cs.MapObject{Position: cs.Vector{X: 60, Y: 60}},
				Stability:      cs.WormholeStabilityMostlyStable,
				DestinationNum: 1,
			},
		},
		MysteryTraders: []cs.MysteryTrader{
			{
				MapObject:     cs.MapObject{Position: cs.Vector{X: 100, Y: 100}},
				WarpSpeed:     7,
				Destination:   cs.Vector{X: -10, Y: -10},
				RequestedBoon: 5000,
			},
		},
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
		Name: "Cargo Transfer Invasion Starbase",
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
						Cargo:             cs.Cargo{Ironium: 10, Colonists: 200},
						OrbitingPlanetNum: 1,
					},
				},
			},
			{
				Player: &cs.Player{Name: "Player 2", AIControlled: true, Race: *cs.NewRace()},
				Designs: []cs.ShipDesign{
					{
						Name: "Starbase",
						Hull: cs.SpaceStation.Name,
					},
				},
			},
		},
		Planets: []cs.Planet{{
			MapObject: cs.MapObject{
				Name:      "Planet 1",
				PlayerNum: 2, // give player2 a planet
			},
			Hab:   cs.Hab{Grav: 50, Temp: 50, Rad: 50},
			Cargo: cs.Cargo{Ironium: 1000, Boranium: 1000, Germanium: 1000, Colonists: 25},
			Starbase: &cs.Fleet{

				BaseName:  "Starbase",
				Tokens:    []cs.ShipToken{{DesignNum: 1, Quantity: 1}},
				PlanetNum: 1,
			},
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
						Cargo:    cs.Cargo{Ironium: 10, Boranium: 10, Germanium: 10, Colonists: 10},
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
						Heading:         cs.Vector{X: 1},
					},
				},
			},
		},
		Planets: []cs.Planet{{
			MapObject: cs.MapObject{
				Position: cs.Vector{X: 50},
				Name:     "Planet 1",
			},
		}},
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
	{
		Name: "Battle 1",
		Players: []TestPlayer{
			{
				Designs: []cs.ShipDesign{
					{
						Name:  "Destroyer",
						Hull:  cs.Destroyer.Name,
						Slots: destroyerDeltaSlots,
					},
				},
				Fleets: []cs.Fleet{
					{
						BaseName:          "Destroyer Delta",
						Tokens:            []cs.ShipToken{{DesignNum: 1, Quantity: 1}},
						OrbitingPlanetNum: 1,
					},
				},
			},
			{
				Player: &cs.Player{Name: "Player 2", AIControlled: true, Race: *cs.NewRace()},
				Designs: []cs.ShipDesign{
					{
						Name:  "Santa Maria",
						Hull:  cs.ColonyShip.Name,
						Slots: santaMariaSlots,
					},
				},
				Fleets: []cs.Fleet{
					{
						BaseName:          "Santa Maria",
						Tokens:            []cs.ShipToken{{DesignNum: 1, Quantity: 2}},
						OrbitingPlanetNum: 1,
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
		},
	},
}

// CreateTestGames creates one of each test game for manual testing
// for automated testing they should be created new each time
func CreateTestGames(db db.Client) error {
	ctx := context.Background()
	for _, testGame := range TestGames {
		var err error
		game := CreateTestGame(testGame)
		err = db.SaveGame(ctx, game.Game)
		if err != nil {
			return err
		}
		for _, player := range game.Players {
			player.GameID = game.ID

			if err := db.SavePlayer(ctx, player); err != nil {
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
		if err := db.UpdateFullGame(ctx, game); err != nil {
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

func CreateTestGame(tg TestGame) *cs.FullGame {
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
		var testPlayerPlayer cs.Player
		if p.Player == nil {
			testPlayerPlayer = cs.Player{
				UserID:        1,
				Race:          *cs.NewRace(),
				TechLevels:    cs.TechLevel{Energy: 3, Weapons: 3, Propulsion: 3, Construction: 3, Electronics: 3, Biotechnology: 3},
				Stats:         &cs.PlayerStats{},
				AcquiredTechs: map[string]bool{},
				PlayerOrders: cs.PlayerOrders{
					Researching:       cs.Energy,
					ResearchAmount:    15,
					NextResearchField: cs.NextResearchFieldSameField,
					CargoTransfers:    cs.CargoTransfers{},
				},
			}
		} else {
			testPlayerPlayer = *p.Player
		}
		player := addPlayer(game, testPlayerPlayer.WithNum(i+1))
		if testPlayerPlayer.AIControlled {
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

		// add any starbases on this planet
		for i, planet := range tg.Planets {
			if planet.PlayerNum == player.Num && planet.Starbase != nil {
				addStarbase(game, player, planet.Starbase, game.Planets[i])
			}
		}

		for i, mp := range p.MineralPackets {
			addMineralPacket(game, player, &mp, i+1)
		}

		for i, mf := range p.Minefields {
			addMinefield(game, player, &mf, i+1)
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
	universe := cs.NewUniverse(slog.Default(), &game.Rules)

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
	if planet.BaseHab == (cs.Hab{}) {
		planet.BaseHab = planet.Hab
	}
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

func addStarbase(game *cs.FullGame, player *cs.Player, fleet *cs.Fleet, planet *cs.Planet) *cs.Fleet {
	fleet.Type = cs.MapObjectTypeFleet
	fleet.PlayerNum = player.Num
	fleet.Name = fleet.BaseName
	fleet.Waypoints = []cs.Waypoint{
		cs.NewPositionWaypoint(fleet.Position, 0),
	}

	fleet.PlanetNum = planet.Num
	fleet.Starbase = true

	game.Starbases = append(game.Starbases, fleet)
	return fleet
}

func addMineralPacket(game *cs.FullGame, player *cs.Player, mineralPacket *cs.MineralPacket, num int) *cs.MineralPacket {
	mineralPacket.Type = cs.MapObjectTypeMineralPacket
	mineralPacket.PlayerNum = player.Num
	mineralPacket.Num = num
	mineralPacket.Name = fmt.Sprintf("%s Mineral Packet #%d", player.Race.PluralName, mineralPacket.Num)
	mineralPacket.Heading = (game.Planets[mineralPacket.TargetPlanetNum-1].Position.Subtract(mineralPacket.Position)).Normalized()

	game.MineralPackets = append(game.MineralPackets, mineralPacket)
	return mineralPacket
}

func addMinefield(game *cs.FullGame, player *cs.Player, minefield *cs.Minefield, num int) *cs.Minefield {
	minefield.Type = cs.MapObjectTypeMinefield
	minefield.PlayerNum = player.Num
	minefield.Num = len(game.Minefields) + 1
	minefield.Name = fmt.Sprintf("%s %s Minefield #%d", player.Race.PluralName, minefield.MinefieldType.String(), num)
	game.Minefields = append(game.Minefields, minefield)
	return minefield
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
	mysteryTrader.Heading = (mysteryTrader.Destination.Subtract(mysteryTrader.Position)).Normalized()
	game.MysteryTraders = append(game.MysteryTraders, mysteryTrader)
	return mysteryTrader
}
