//go:build !wasi && !wasm

package testgames

import (
	"context"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"
)

var TestGames = []cs.TestScenario{
	cs.SingleUnitScenario(),
	cs.TwoPlayerScenario(),
	{
		Name: "Scout Test",
		Players: []cs.TestScenarioPlayer{
			{
				Designs: []cs.ShipDesign{{
					Name:  "Long Range Scout",
					Hull:  cs.Scout.Name,
					Slots: cs.LongRangeScoutTestSlots,
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
		Players: []cs.TestScenarioPlayer{
			{
				Designs: []cs.ShipDesign{
					{
						Name:  "Santa Maria",
						Hull:  cs.ColonyShip.Name,
						Slots: cs.SantaMariaTestSlots,
					},
					{
						Name:  "Long Range Scout",
						Hull:  cs.Scout.Name,
						Slots: cs.LongRangeScoutTestSlots,
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
		Name: "Colonizer Test AR",
		Players: []cs.TestScenarioPlayer{
			{
				Player: cs.NewPlayer(1, cs.NewRace().WithPRT(cs.AR)),
				Designs: []cs.ShipDesign{
					{
						Name:  "Santa Maria",
						Hull:  cs.ColonyShip.Name,
						Slots: cs.SantaMariaARTestSlots,
					},
					{
						Name:  "Long Range Scout",
						Hull:  cs.Scout.Name,
						Slots: cs.LongRangeScoutTestSlots,
					},
					{
						Name:    "Starter Colony",
						Hull:    cs.OrbitalFort.Name,
						Purpose: cs.ShipDesignPurposeStarterColony,
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
		Players: []cs.TestScenarioPlayer{
			{
				Designs: []cs.ShipDesign{
					{
						Name:  "Long Range Scout",
						Hull:  cs.Scout.Name,
						Slots: cs.LongRangeScoutTestSlots,
					},
					{
						Name:  "Teamster",
						Hull:  cs.MediumFreighter.Name,
						Slots: cs.TeamsterTestSlots,
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
						Slots: cs.TeamsterTestSlots,
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
		Players: []cs.TestScenarioPlayer{
			{
				Designs: []cs.ShipDesign{
					{
						Name:  "Teamster",
						Hull:  cs.MediumFreighter.Name,
						Slots: cs.TeamsterTestSlots,
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
		Players: []cs.TestScenarioPlayer{
			{
				Designs: []cs.ShipDesign{
					{
						Name:  "Teamster",
						Hull:  cs.MediumFreighter.Name,
						Slots: cs.TeamsterTestSlots,
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
		Players: []cs.TestScenarioPlayer{
			{
				Designs: []cs.ShipDesign{
					{
						Name:  "Teamster",
						Hull:  cs.MediumFreighter.Name,
						Slots: cs.TeamsterTestSlots,
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
		Players: []cs.TestScenarioPlayer{
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
						Slots: cs.TeamsterTestSlots,
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
		Players: []cs.TestScenarioPlayer{
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
						Slots: cs.TeamsterTestSlots,
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
						Slots: cs.TeamsterTestSlots,
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
		Players: []cs.TestScenarioPlayer{
			{
				Designs: []cs.ShipDesign{
					{
						Name:  "Teamster",
						Hull:  cs.MediumFreighter.Name,
						Slots: cs.TeamsterTestSlots,
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
		Players: []cs.TestScenarioPlayer{
			{
				Designs: []cs.ShipDesign{
					{
						Name:  "Teamster",
						Hull:  cs.MediumFreighter.Name,
						Slots: cs.TeamsterTestSlots,
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
		Players: []cs.TestScenarioPlayer{
			{
				Designs: []cs.ShipDesign{
					{
						Name:  "Teamster",
						Hull:  cs.MediumFreighter.Name,
						Slots: cs.TeamsterTestSlots,
					},
					{
						Name:  "Santa Maria",
						Hull:  cs.ColonyShip.Name,
						Slots: cs.SantaMariaTestSlots,
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
		Players: []cs.TestScenarioPlayer{
			{
				Designs: []cs.ShipDesign{
					{
						Name:  "Teamster",
						Hull:  cs.MediumFreighter.Name,
						Slots: cs.TeamsterTestSlots,
					},
					{
						Name:  "Santa Maria",
						Hull:  cs.ColonyShip.Name,
						Slots: cs.SantaMariaTestSlots,
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
		Players: []cs.TestScenarioPlayer{
			{
				Designs: []cs.ShipDesign{
					{
						Name:  "Teamster",
						Hull:  cs.MediumFreighter.Name,
						Slots: cs.TeamsterTestSlots,
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
						Heading:         cs.VectorFloat64{X: 1},
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
		Name: "Mystery Trader",
		Players: []cs.TestScenarioPlayer{
			{
				Designs: []cs.ShipDesign{
					{
						Name: "Large Freighter",
						Hull: cs.LargeFreighter.Name,
						Slots: []cs.ShipDesignSlot{
							{HullComponent: cs.TransStar10.Name, HullSlotIndex: 1, Quantity: 2},
							{HullComponent: cs.SuperCargoPod.Name, HullSlotIndex: 2, Quantity: 2},
						},
					},
				},
				Fleets: []cs.Fleet{
					{
						MapObject: cs.MapObject{
							Position: cs.Vector{X: 100},
						},
						BaseName: "MT Meet-er",
						Tokens:   []cs.ShipToken{{DesignNum: 1, Quantity: 10}},
						Cargo:    cs.Cargo{Ironium: 1400 * 10},
						Fuel:     2600 * 10,
					},
					{
						MapObject: cs.MapObject{
							Position: cs.Vector{X: 100, Y: 25},
						},
						BaseName: "MT Meet-er",
						Tokens:   []cs.ShipToken{{DesignNum: 1, Quantity: 10}},
						Cargo:    cs.Cargo{Ironium: 1400 * 10},
						Fuel:     2600 * 10,
					},
					{
						MapObject: cs.MapObject{
							Position: cs.Vector{X: 100, Y: 50},
						},
						BaseName: "MT Meet-er",
						Tokens:   []cs.ShipToken{{DesignNum: 1, Quantity: 10}},
						Cargo:    cs.Cargo{Ironium: 1400 * 10},
						Fuel:     2600 * 10,
					},
					{
						MapObject: cs.MapObject{
							Position: cs.Vector{X: 100, Y: 75},
						},
						BaseName: "MT Meet-er",
						Tokens:   []cs.ShipToken{{DesignNum: 1, Quantity: 10}},
						Cargo:    cs.Cargo{Ironium: 1400 * 10},
						Fuel:     2600 * 10,
					},
					{
						MapObject: cs.MapObject{
							Position: cs.Vector{X: 100, Y: 100},
						},
						BaseName: "MT Meet-er",
						Tokens:   []cs.ShipToken{{DesignNum: 1, Quantity: 10}},
						Cargo:    cs.Cargo{Ironium: 1400 * 10},
						Fuel:     2600 * 10,
					},
					{
						MapObject: cs.MapObject{
							Position: cs.Vector{X: 100, Y: 125},
						},
						BaseName: "MT Meet-er",
						Tokens:   []cs.ShipToken{{DesignNum: 1, Quantity: 10}},
						Cargo:    cs.Cargo{Ironium: 1400 * 10},
						Fuel:     2600 * 10,
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
			Mines:                1000,
			Factories:            1000,
			Scanner:              true,
			Cargo:                cs.Cargo{Ironium: 10000, Boranium: 10000, Germanium: 10000, Colonists: 10000},
		}},
		MysteryTraders: []cs.MysteryTrader{
			{
				MapObject:     cs.MapObject{Position: cs.Vector{X: 200}}, // meet our freighter on turn 2
				WarpSpeed:     10,
				Destination:   cs.Vector{X: 0},
				RequestedBoon: 5000,
				RewardType:    cs.MysteryTraderRewardResearch,
			},
			{
				MapObject:     cs.MapObject{Position: cs.Vector{X: 200, Y: 25}},
				WarpSpeed:     10,
				Destination:   cs.Vector{X: 0, Y: 25},
				RequestedBoon: 5000,
				RewardType:    cs.MysteryTraderRewardLifeboat,
			},
			{
				MapObject:     cs.MapObject{Position: cs.Vector{X: 200, Y: 50}},
				WarpSpeed:     10,
				Destination:   cs.Vector{X: 0, Y: 50},
				RequestedBoon: 5000,
				RewardType:    cs.MysteryTraderRewardArmor,
			},
			{
				MapObject:     cs.MapObject{Position: cs.Vector{X: 200, Y: 75}},
				WarpSpeed:     10,
				Destination:   cs.Vector{X: 0, Y: 75},
				RequestedBoon: 5000,
				RewardType:    cs.MysteryTraderRewardGenesis,
			},
			{
				MapObject:     cs.MapObject{Position: cs.Vector{X: 200, Y: 100}},
				WarpSpeed:     10,
				Destination:   cs.Vector{X: 0, Y: 100},
				RequestedBoon: 5000,
				RewardType:    cs.MysteryTraderRewardJumpGate,
			},
			{
				MapObject:     cs.MapObject{Position: cs.Vector{X: 200, Y: 125}},
				WarpSpeed:     10,
				Destination:   cs.Vector{X: 0, Y: 125},
				RequestedBoon: 5000,
				RewardType:    cs.MysteryTraderRewardShipHull,
			},
		},
	},
	{
		Name: "Battle 1",
		Players: []cs.TestScenarioPlayer{
			{
				Designs: []cs.ShipDesign{
					{
						Name:  "Destroyer",
						Hull:  cs.Destroyer.Name,
						Slots: cs.DestroyerDeltaTestSlots,
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
						Slots: cs.SantaMariaTestSlots,
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
		game := cs.BuildScenario(testGame)
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

		// save to db
		if err := db.UpdateFullGame(ctx, game); err != nil {
			return err
		}
	}

	return nil
}
