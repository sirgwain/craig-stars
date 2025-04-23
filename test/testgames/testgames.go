package testgames

import "github.com/sirgwain/craig-stars/cs"

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
						Fuel:              300,
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
						OrbitingPlanetNum: 1,
					},
					{
						MapObject: cs.MapObject{Position: cs.Vector{X: 0, Y: 50}},
						BaseName:  "Teamster",
						Tokens:    []cs.ShipToken{{DesignNum: 2, Quantity: 2}},
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
				MineFields: []cs.MineField{
					{
						MineFieldType: cs.MineFieldTypeStandard,
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
				MapObject:      cs.MapObject{Position: cs.Vector{X: 20, Y: 0}},
				Stability:      cs.WormholeStabilityRockSolid,
				DestinationNum: 2,
			},
			{
				MapObject:      cs.MapObject{Position: cs.Vector{X: 0, Y: 20}},
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
	// TODO: Make a 12 player AI controlled game
}
