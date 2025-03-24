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
}
