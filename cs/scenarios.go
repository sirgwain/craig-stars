//go:build !wasi && !wasm

package cs

import "fmt"

// ScenarioScoutTest is shared by browser and game-logic tests.
func ScenarioScoutTest() TestScenario {
	return TestScenario{Name: "Scout Test",
		Players: []ScenarioPlayer{{Designs: Designs(DesignLongRangeScout),
			Fleets: []ScenarioFleet{{Design: "Long Range Scout",
				At:   "Planet 1",
				Fuel: 300}}}},
		Planets: []ScenarioPlanet{Homeworld("Planet 1", 1),
			{Name: "Planet 2",
				Position: Vector{X: 0, Y: 49},
				Hab:      ScenarioValue(Hab{Grav: 25, Temp: 25, Rad: 25}),
				Cargo:    Cargo{}}}}
}

// ScenarioColonizerTest is shared by browser and game-logic tests.
func ScenarioColonizerTest() TestScenario {
	return TestScenario{Name: "Colonizer Test",
		Players: []ScenarioPlayer{{Designs: Designs(DesignSantaMaria, DesignLongRangeScout),
			Fleets: []ScenarioFleet{{Design: "Santa Maria",
				At:   "Planet 1",
				Fuel: 200},
				{Design: "Long Range Scout",
					At:   "Planet 1",
					Fuel: 300}}}},
		Planets: []ScenarioPlanet{Homeworld("Planet 1", 1),
			{Name: "Planet 2",
				Position: Vector{X: 0, Y: 25},
				Cargo:    Cargo{}}}}
}

// ScenarioColonizerTestAR is shared by browser and game-logic tests.
func ScenarioColonizerTestAR() TestScenario {
	return TestScenario{Name: "Colonizer Test AR",
		Players: []ScenarioPlayer{{Player: NewPlayer(1, NewRace().WithPRT(AR)),
			Designs: Designs(DesignSantaMariaAR, DesignLongRangeScout, ShipDesign{Name: "Starter Colony",
				Hull:    OrbitalFort.Name,
				Purpose: ShipDesignPurposeStarterColony}),
			Fleets: []ScenarioFleet{{Design: "Santa Maria",
				At:   "Planet 1",
				Fuel: 200},
				{Design: "Long Range Scout",
					At:   "Planet 1",
					Fuel: 300}}}},
		Planets: []ScenarioPlanet{Homeworld("Planet 1", 1),
			{Name: "Planet 2",
				Position: Vector{X: 0, Y: 25},
				Cargo:    Cargo{}}}}
}

// ScenarioSDMinefield is shared by browser and game-logic tests.
// An SD player with a standard minefield it can detonate.
func ScenarioSDMinefield() TestScenario {
	return TestScenario{Name: "SD Minefield",
		Players: []ScenarioPlayer{{Player: NewPlayer(1, NewRace().WithPRT(SD)),
			Minefields: []Minefield{
				{
					MinefieldType: MinefieldTypeStandard,
					NumMines:      1000,
				},
			}}},
		Planets: []ScenarioPlanet{Homeworld("Planet 1", 1)}}
}

// ScenarioKitchenSink is shared by browser and game-logic tests.
func ScenarioKitchenSink() TestScenario {
	return TestScenario{Name: "Kitchen Sink",
		Players: []ScenarioPlayer{{Designs: Designs(DesignLongRangeScout, DesignTeamster),
			Fleets: []ScenarioFleet{{Design: "Long Range Scout",
				At:   "Planet 1",
				Fuel: 300},
				{Design: "Teamster",
					Quantity: 2,
					Position: Vector{X: 0, Y: 50},
					Fuel:     900},
				// two fleets together in deep space, for the scanner's shared token count
				{Design: "Long Range Scout",
					Position: Vector{X: 80, Y: 20},
					Fuel:     300},
				{Design: "Teamster",
					Quantity: 3,
					Position: Vector{X: 80, Y: 20},
					Fuel:     1350}},
			MineralPackets: []ScenarioMineralPacket{{Cargo: Cargo{Ironium: 50, Boranium: 50, Germanium: 50},
				WarpSpeed: 5,
				Position:  Vector{X: 50, Y: 0},
				To:        "Planet 1"}},
			Minefields: []Minefield{
				{
					MinefieldType: MinefieldTypeStandard,
					NumMines:      1000,
				},
			},
			Salvages: []Salvage{
				{
					MapObject: MapObject{Position: Vector{X: 50, Y: 50}},
					Cargo:     Cargo{Ironium: 50, Boranium: 50, Germanium: 50},
				},
			}},
			{Player: &Player{Name: "Player 2", AIControlled: true, Race: *NewRace().WithPluralName("Rabbitoids")},
				Designs: Designs(DesignTeamster),
				Fleets: []ScenarioFleet{{Design: "Teamster",
					Position: Vector{X: 30, Y: 0},
					Cargo:    Cargo{Ironium: 10, Boranium: 10, Germanium: 10, Colonists: 10},
					Fuel:     10}}}},
		Planets: []ScenarioPlanet{{Name: "Planet 1",
			Owner: 1,
			Cargo: Cargo{Ironium: 1000, Boranium: 1000, Germanium: 1000, Colonists: 2500}},
			{Name: "Planet 2",
				Owner:    2,
				Position: Vector{X: 0, Y: 30},
				Hab:      ScenarioValue(Hab{Grav: 25, Temp: 25, Rad: 25}),
				Cargo:    Cargo{Ironium: 1000, Boranium: 1000, Germanium: 1000, Colonists: 2500}},
			{Name: "Planet 3",
				Position: Vector{X: 30, Y: 30},
				Hab:      ScenarioValue(Hab{Grav: 75, Temp: 75, Rad: 75})}},
		Wormholes: []Wormhole{
			{
				MapObject:      MapObject{Position: Vector{X: 10, Y: 10}},
				Stability:      WormholeStabilityRockSolid,
				DestinationNum: 2,
			},
			{
				MapObject:      MapObject{Position: Vector{X: 60, Y: 60}},
				Stability:      WormholeStabilityMostlyStable,
				DestinationNum: 1,
			},
		},
		MysteryTraders: []MysteryTrader{
			{
				MapObject:     MapObject{Position: Vector{X: 100, Y: 100}},
				WarpSpeed:     7,
				Destination:   Vector{X: -10, Y: -10},
				RequestedBoon: 5000,
			},
		}}
}

// ScenarioCargoTransferInvasion is shared by browser and game-logic tests.
func ScenarioCargoTransferInvasion() TestScenario {
	return TestScenario{Name: "Cargo Transfer Invasion",
		Players: []ScenarioPlayer{{Designs: Designs(DesignTeamster),
			Fleets: []ScenarioFleet{{Design: "Teamster",
				At:    "Planet 1",
				Fuel:  450,
				Cargo: Cargo{Colonists: 210}}}},
			AIPlayer("Player 2")},
		Planets: []ScenarioPlanet{{Name: "Planet 1",
			Owner:         2,
			Concentration: ScenarioValue(Mineral{}),
			Cargo:         Cargo{Ironium: 1000, Boranium: 1000, Germanium: 1000, Colonists: 25}}}}
}

// ScenarioCargoTransferInvasionStarbase is shared by browser and game-logic tests.
func ScenarioCargoTransferInvasionStarbase() TestScenario {
	return TestScenario{Name: "Cargo Transfer Invasion Starbase",
		Players: []ScenarioPlayer{{Designs: Designs(DesignTeamster),
			Fleets: []ScenarioFleet{{Design: "Teamster",
				At:    "Planet 1",
				Fuel:  450,
				Cargo: Cargo{Ironium: 10, Colonists: 200}}}},
			{Player: &Player{Name: "Player 2", AIControlled: true, Race: *NewRace()},
				Designs: Designs(ShipDesign{Name: "Starbase",
					Hull: SpaceStation.Name})}},
		Planets: []ScenarioPlanet{{Name: "Planet 1",
			Owner:         2,
			Concentration: ScenarioValue(Mineral{}),
			Starbase:      "Starbase",
			Cargo:         Cargo{Ironium: 1000, Boranium: 1000, Germanium: 1000, Colonists: 25}}}}
}

// ScenarioCargoTransferPlanetOwned is shared by browser and game-logic tests.
func ScenarioCargoTransferPlanetOwned() TestScenario {
	return TestScenario{Name: "Cargo Transfer Planet Owned",
		Players: []ScenarioPlayer{{Designs: Designs(DesignTeamster),
			Fleets: []ScenarioFleet{{Design: "Teamster",
				At:   "Planet 1",
				Fuel: 450}}}},
		Planets: []ScenarioPlanet{{Name: "Planet 1",
			Owner:         1,
			Concentration: ScenarioValue(Mineral{}),
			Cargo:         Cargo{Ironium: 1000, Boranium: 1000, Germanium: 1000, Colonists: 2500}}}}
}

// ScenarioCargoTransferPlanetSteal is shared by browser and game-logic tests.
func ScenarioCargoTransferPlanetSteal() TestScenario {
	return TestScenario{Name: "Cargo Transfer Planet Steal",
		Players: []ScenarioPlayer{{Player: &Player{Name: "Player 1", UserID: 1, Race: *NewRace().WithPRT(SS)},
			Designs: Designs(DesignThief, DesignTeamster),
			Fleets: []ScenarioFleet{{Design: "Thief",
				At:   "Planet 1",
				Fuel: 450},
				{Design: "Teamster",
					At:   "Planet 1",
					Fuel: 450}}},
			AIPlayer("Player 2")},
		Planets: []ScenarioPlanet{{Name: "Planet 1",
			Owner:         2,
			Concentration: ScenarioValue(Mineral{}),
			Cargo:         Cargo{Ironium: 1000, Boranium: 1000, Germanium: 1000, Colonists: 2500}}}}
}

// ScenarioCargoTransferFleetSteal is shared by browser and game-logic tests.
func ScenarioCargoTransferFleetSteal() TestScenario {
	return TestScenario{Name: "Cargo Transfer Fleet Steal",
		Players: []ScenarioPlayer{{Player: &Player{Name: "Player 1", UserID: 1, Race: *NewRace().WithPRT(SS)},
			Designs: Designs(DesignThief, DesignTeamster),
			Fleets: []ScenarioFleet{{Design: "Thief",
				Fuel: 100},
				{Design: "Teamster",
					Fuel: 100}}},
			{Player: &Player{Name: "Player 2", AIControlled: true, Race: *NewRace()},
				Designs: Designs(DesignTeamster),
				Fleets: []ScenarioFleet{{Design: "Teamster",
					Cargo: Cargo{Ironium: 10, Boranium: 10, Germanium: 10, Colonists: 10},
					Fuel:  10}}}}}
}

// ScenarioCargoTransferJettison is shared by browser and game-logic tests.
func ScenarioCargoTransferJettison() TestScenario {
	return TestScenario{Name: "Cargo Transfer Jettison",
		Players: []ScenarioPlayer{{Designs: Designs(DesignTeamster),
			Fleets: []ScenarioFleet{{Design: "Teamster",
				Name:  "Teamster Jettison #1",
				Cargo: Cargo{Ironium: 50, Boranium: 50, Germanium: 50},
				Fuel:  500}}}}}
}

// ScenarioCargoTransferSalvage is shared by browser and game-logic tests.
func ScenarioCargoTransferSalvage() TestScenario {
	return TestScenario{Name: "Cargo Transfer Salvage",
		Players: []ScenarioPlayer{{Designs: Designs(DesignTeamster),
			Fleets: []ScenarioFleet{{Design: "Teamster",
				Name:  "Teamster Salvager #1",
				Cargo: Cargo{Ironium: 10, Boranium: 10, Germanium: 10, Colonists: 10},
				Fuel:  500}},
			Salvages: []Salvage{
				{
					Cargo: Cargo{Ironium: 50, Boranium: 50, Germanium: 50},
				},
			}}}}
}

// ScenarioCargoTransferFleets is shared by browser and game-logic tests.
func ScenarioCargoTransferFleets() TestScenario {
	return TestScenario{Name: "Cargo Transfer Fleets",
		Players: []ScenarioPlayer{{Designs: Designs(DesignTeamster, DesignSantaMaria),
			Fleets: []ScenarioFleet{{Design: "Teamster",
				Cargo: Cargo{Ironium: 10, Boranium: 10, Germanium: 10, Colonists: 10},
				Fuel:  100},
				{Design: "Santa Maria",
					Cargo: Cargo{Ironium: 5, Boranium: 5, Germanium: 5, Colonists: 5},
					Fuel:  10}}}}}
}

// ScenarioCargoTransferSplit is shared by browser and game-logic tests.
func ScenarioCargoTransferSplit() TestScenario {
	return TestScenario{Name: "Cargo Transfer Split",
		Players: []ScenarioPlayer{{Designs: Designs(DesignTeamster, DesignSantaMaria),
			Fleets: []ScenarioFleet{{Tokens: []ScenarioShipToken{{Design: "Teamster",
				Quantity: 1},
				{Design: "Santa Maria",
					Quantity: 1}},
				Name:  "Teamster Jettison #1",
				Cargo: Cargo{Ironium: 50, Boranium: 50, Germanium: 50},
				Fuel:  500}}}}}
}

// ScenarioCargoTransferMineralPacket is shared by browser and game-logic tests.
func ScenarioCargoTransferMineralPacket() TestScenario {
	return TestScenario{Name: "Cargo Transfer MineralPacket",
		Players: []ScenarioPlayer{{Designs: Designs(DesignTeamster),
			Fleets: []ScenarioFleet{{Design: "Teamster",
				Name:  "Teamster Mineral Packeter #1",
				Cargo: Cargo{Ironium: 10, Boranium: 10, Germanium: 10, Colonists: 10},
				Fuel:  500}},
			MineralPackets: []ScenarioMineralPacket{{Cargo: Cargo{Ironium: 50, Boranium: 50, Germanium: 50},
				WarpSpeed: 5,
				To:        "Planet 1"}}}},
		Planets: []ScenarioPlanet{{Name: "Planet 1", Hab: ScenarioValue(Hab{}),
			Position:      Vector{X: 50},
			Concentration: ScenarioValue(Mineral{})}}}
}

// ScenarioMysteryTrader is shared by browser and game-logic tests.
func ScenarioMysteryTrader() TestScenario {
	return TestScenario{Name: "Mystery Trader",
		Players: []ScenarioPlayer{{Designs: Designs(ShipDesign{Name: "Large Freighter",
			Hull: LargeFreighter.Name,
			Slots: []ShipDesignSlot{
				{HullComponent: TransStar10.Name, HullSlotIndex: 1, Quantity: 2},
				{HullComponent: SuperCargoPod.Name, HullSlotIndex: 2, Quantity: 2},
			}}),
			Fleets: []ScenarioFleet{{Design: "Large Freighter",
				Quantity: 10,
				Name:     "MT Meet-er #1",
				Position: Vector{X: 100},
				Cargo:    Cargo{Ironium: 1400 * 10},
				Fuel:     2600 * 10},
				{Design: "Large Freighter",
					Quantity: 10,
					Name:     "MT Meet-er #2",
					Position: Vector{X: 100, Y: 25},
					Cargo:    Cargo{Ironium: 1400 * 10},
					Fuel:     2600 * 10},
				{Design: "Large Freighter",
					Quantity: 10,
					Name:     "MT Meet-er #3",
					Position: Vector{X: 100, Y: 50},
					Cargo:    Cargo{Ironium: 1400 * 10},
					Fuel:     2600 * 10},
				{Design: "Large Freighter",
					Quantity: 10,
					Name:     "MT Meet-er #4",
					Position: Vector{X: 100, Y: 75},
					Cargo:    Cargo{Ironium: 1400 * 10},
					Fuel:     2600 * 10},
				{Design: "Large Freighter",
					Quantity: 10,
					Name:     "MT Meet-er #5",
					Position: Vector{X: 100, Y: 100},
					Cargo:    Cargo{Ironium: 1400 * 10},
					Fuel:     2600 * 10},
				{Design: "Large Freighter",
					Quantity: 10,
					Name:     "MT Meet-er #6",
					Position: Vector{X: 100, Y: 125},
					Cargo:    Cargo{Ironium: 1400 * 10},
					Fuel:     2600 * 10}}}},
		Planets: []ScenarioPlanet{{Name: "Planet 1",
			Owner:     1,
			Mines:     1000,
			Factories: 1000,
			Scanner:   true,
			Cargo:     Cargo{Ironium: 10000, Boranium: 10000, Germanium: 10000, Colonists: 10000}}},
		MysteryTraders: []MysteryTrader{
			{
				MapObject:     MapObject{Position: Vector{X: 200}},
				WarpSpeed:     10,
				Destination:   Vector{X: 0},
				RequestedBoon: 5000,
				RewardType:    MysteryTraderRewardResearch,
			},
			{
				MapObject:     MapObject{Position: Vector{X: 200, Y: 25}},
				WarpSpeed:     10,
				Destination:   Vector{X: 0, Y: 25},
				RequestedBoon: 5000,
				RewardType:    MysteryTraderRewardLifeboat,
			},
			{
				MapObject:     MapObject{Position: Vector{X: 200, Y: 50}},
				WarpSpeed:     10,
				Destination:   Vector{X: 0, Y: 50},
				RequestedBoon: 5000,
				RewardType:    MysteryTraderRewardArmor,
			},
			{
				MapObject:     MapObject{Position: Vector{X: 200, Y: 75}},
				WarpSpeed:     10,
				Destination:   Vector{X: 0, Y: 75},
				RequestedBoon: 5000,
				RewardType:    MysteryTraderRewardGenesis,
			},
			{
				MapObject:     MapObject{Position: Vector{X: 200, Y: 100}},
				WarpSpeed:     10,
				Destination:   Vector{X: 0, Y: 100},
				RequestedBoon: 5000,
				RewardType:    MysteryTraderRewardJumpGate,
			},
			{
				MapObject:     MapObject{Position: Vector{X: 200, Y: 125}},
				WarpSpeed:     10,
				Destination:   Vector{X: 0, Y: 125},
				RequestedBoon: 5000,
				RewardType:    MysteryTraderRewardShipHull,
			},
		}}
}

// ScenarioBattle1 is shared by browser and game-logic tests.
func ScenarioBattle1() TestScenario {
	return TestScenario{Name: "Battle 1",
		Players: []ScenarioPlayer{{Designs: Designs(DesignDestroyerDelta),
			Fleets: []ScenarioFleet{{Design: "Destroyer",
				Name:      "Destroyer Delta #1",
				At:        "Planet 1",
				EmptyFuel: true}}},
			{Player: &Player{Name: "Player 2", AIControlled: true, Race: *NewRace()},
				Designs: Designs(DesignSantaMaria),
				Fleets: []ScenarioFleet{{Design: "Santa Maria",
					Quantity:  2,
					At:        "Planet 1",
					EmptyFuel: true}}}},
		Planets: []ScenarioPlanet{{Name: "Planet 1",
			Owner: 1,
			Cargo: Cargo{Ironium: 1000, Boranium: 1000, Germanium: 1000, Colonists: 2500}}}}
}

// scenarioBattle puts fleets together in deep space, with separate homeworlds
// so every player survives to view the battle report. Player 1 is human.
func scenarioBattle(name string, designs ...ShipDesign) TestScenario {
	s := TestScenario{Name: name}
	for i, design := range designs {
		player := defaultScenarioPlayer()
		player.Name = fmt.Sprintf("Player %d", i+1)
		player.AIControlled = i > 0
		player.Race.PluralName = fmt.Sprintf("Player %d ships", i+1)
		player.TechLevels = TechLevel{Energy: 26, Weapons: 26, Propulsion: 26, Construction: 26, Electronics: 26, Biotechnology: 26}
		s.Players = append(s.Players, ScenarioPlayer{Player: &player,
			Designs: Designs(design),
			Fleets:  []ScenarioFleet{{Design: design.Name, Quantity: 3, Position: Vector{100, 100}}}})
		planet := Homeworld(fmt.Sprintf("Homeworld %d", i+1), i+1)
		planet.Position = Vector{10 + (i%5)*35, 10 + (i/5)*20}
		s.Planets = append(s.Planets, planet)
	}
	return s
}

// ScenarioBattleTorpedoes exercises volley hits, misses, shield loss, and partial kills.
func ScenarioBattleTorpedoes() TestScenario {
	s := scenarioBattle("Battle Torpedoes", DesignBattleTorpedo, DesignBattleTorpedo)
	// Jamming makes accuracy differ between the otherwise identical stacks.
	s.Players[1].Designs[0].Slots[6].HullComponent = Jammer30.Name
	return s
}

// ScenarioBattleBeams exercises shield depletion, beam attenuation, and spillover
// between separate stacks of the same design.
func ScenarioBattleBeams() TestScenario {
	s := scenarioBattle("Battle Beams", DesignBattleBeam, DesignBattleBeam)
	s.Players[0].Fleets[0].Quantity = 8
	s.Players[1].Fleets = []ScenarioFleet{
		{Design: DesignBattleBeam.Name, Quantity: 2, Position: Vector{100, 100}},
		{Design: DesignBattleBeam.Name, Quantity: 2, Position: Vector{100, 100}},
		{Design: DesignBattleBeam.Name, Quantity: 2, Position: Vector{100, 100}},
	}
	s.Players[1].Designs[0].Slots[6].HullComponent = BeamDeflector.Name
	return s
}

// ScenarioBattleStarbases exercises an armed, stationary base with its range
// bonus, jamming, and escorts against fleets ordered to attack the base first.
func ScenarioBattleStarbases() TestScenario {
	s := scenarioBattle("Battle Starbases", DesignBattleTorpedo, DesignBattleBeam)
	s.Players[0].Designs = Designs(DesignBattleTorpedo, DesignBattleBeam)
	s.Players[0].Fleets = []ScenarioFleet{
		{Design: DesignBattleTorpedo.Name, Quantity: 12, At: "Battle Planet", BattlePlanNum: 1},
		{Design: DesignBattleBeam.Name, Quantity: 12, At: "Battle Planet", BattlePlanNum: 1},
	}
	s.Players[1].Designs = Designs(DesignBattleBeam, DesignBattleStarbase)
	s.Players[1].Fleets[0].At = "Battle Planet"
	s.Planets = append(s.Planets, ScenarioPlanet{Name: "Battle Planet", Owner: 2,
		Position: Vector{100, 100}, Cargo: Cargo{Colonists: 2500}, Starbase: DesignBattleStarbase.Name})
	return s
}

// ScenarioBattleDumpCargo compares identical loaded freighters with and without
// mineral dumping. Colonists remain aboard; dumping changes mass and speed.
func ScenarioBattleDumpCargo() TestScenario {
	s := scenarioBattle("Battle Dump Cargo", DesignPrivateer, DesignBattleTorpedo)
	s.Players[0].BattlePlans = []BattlePlan{{Num: 5, Name: "Dump Cargo and Retreat",
		PrimaryTarget: BattleTargetAny, Tactic: BattleTacticDisengage,
		AttackWho: BattleAttackWhoEnemiesAndNeutrals, DumpCargo: true}}
	s.Players[0].Fleets = []ScenarioFleet{
		{Name: "Dump Cargo", Design: DesignPrivateer.Name, Quantity: 3, Position: Vector{100, 100},
			Cargo: Cargo{Ironium: 150, Boranium: 150, Germanium: 150, Colonists: 150}, BattlePlanNum: 5},
		{Name: "Keep Cargo", Design: DesignPrivateer.Name, Quantity: 3, Position: Vector{100, 100},
			Cargo: Cargo{Ironium: 150, Boranium: 150, Germanium: 150, Colonists: 150}, BattlePlanNum: 4},
	}
	return s
}

// ScenarioBattleRunAway gives an armed ship time to fire while retreating and
// leave the board after its retreat countdown.
func ScenarioBattleRunAway() TestScenario {
	s := scenarioBattle("Battle Run Away", DesignBattleTorpedo, DesignBattleTorpedo)
	s.Players[0].Fleets[0].BattlePlanNum = 4
	return s
}

// ScenarioBattleDisengageIfChallenged tests a fresh retreat countdown after armor
// damage; the high shield capacity also allows shield-only hits beforehand.
func ScenarioBattleDisengageIfChallenged() TestScenario {
	s := scenarioBattle("Battle Disengage if Challenged", DesignBattleBeam, DesignBattleTorpedo)
	s.Players[0].BattlePlans = []BattlePlan{{Num: 5, Name: "Disengage if Challenged",
		PrimaryTarget: BattleTargetAny, Tactic: BattleTacticDisengageIfChallenged,
		AttackWho: BattleAttackWhoEnemiesAndNeutrals}}
	s.Players[0].Fleets[0].BattlePlanNum = 5
	return s
}

// ScenarioBattle3Players exercises the triangular formation in a free-for-all.
func ScenarioBattle3Players() TestScenario {
	return scenarioBattle("Battle 3 Players", DesignBattleBeam, DesignBattleTorpedo, DesignBattleBeam)
}

// ScenarioBattle16Players exercises every position in the largest formation.
func ScenarioBattle16Players() TestScenario { return scenarioBattleManyPlayers(16) }

// ScenarioBattle20Players exercises reuse of formation positions above 16 players.
func ScenarioBattle20Players() TestScenario { return scenarioBattleManyPlayers(20) }

func scenarioBattleManyPlayers(count int) TestScenario {
	designs := make([]ShipDesign, count)
	for i := range designs {
		designs[i] = DesignBattleBeam
	}
	return scenarioBattle(fmt.Sprintf("Battle %d Players", count), designs...)
}

// ScenarioBattleMissilesAndChaff exercises the one-missile-one-kill limit and
// capital missiles against both shielded warships and unshielded chaff.
func ScenarioBattleMissilesAndChaff() TestScenario {
	s := scenarioBattle("Battle Missiles and Chaff", DesignJihadCruiser, DesignBattleBeam)
	s.Players[1].Designs = Designs(DesignBattleBeam, DesignBattleChaff)
	s.Players[1].Fleets[0].Tokens = []ScenarioShipToken{
		{Design: DesignBattleBeam.Name, Quantity: 3},
		{Design: DesignBattleChaff.Name, Quantity: 40},
	}
	s.Players[1].Fleets[0].Design = ""
	return s
}

// ScenarioBattleShieldSappers combines shield-only weapons with ordinary beams
// against regenerating shields (the default race has RS).
func ScenarioBattleShieldSappers() TestScenario {
	s := scenarioBattle("Battle Shield Sappers", DesignBattleBeam, DesignBattleBeam)
	s.Players[0].Designs[0].Slots[1].HullComponent = PulsedSapper.Name
	return s
}

// ScenarioBattleDamagedStacks starts with existing armor damage and unequal
// stacks, exercising damage preservation and casualty records.
func ScenarioBattleDamagedStacks() TestScenario {
	s := scenarioBattle("Battle Damaged Stacks", DesignBattleBeam, DesignBattleTorpedo)
	s.Players[0].Fleets[0].Design = ""
	s.Players[0].Fleets[0].Tokens = []ScenarioShipToken{{Design: DesignBattleBeam.Name,
		Quantity: 5, QuantityDamaged: 3, Damage: 100}}
	return s
}

// ScenarioBattleAlliedSupport has an aggressor, a retaliating defender, an ally
// with no attack orders, and an uninvolved observer at the same location.
func ScenarioBattleAlliedSupport() TestScenario {
	s := scenarioBattle("Battle Allied Support", DesignBattleBeam, DesignBattleTorpedo, DesignBattleBeam, DesignBattleBeam)
	for i := range s.Players {
		s.Players[i].Relations = make([]PlayerRelationship, len(s.Players))
		for j := range s.Players {
			s.Players[i].Relations[j].Relation = PlayerRelationNeutral
		}
		s.Players[i].Relations[i].Relation = PlayerRelationFriend
		s.Players[i].BattlePlans = []BattlePlan{{Num: 0, Name: "Enemies Only",
			PrimaryTarget: BattleTargetArmedShips, SecondaryTarget: BattleTargetAny,
			Tactic: BattleTacticMaximizeDamageRatio, AttackWho: BattleAttackWhoEnemies}}
	}
	s.Players[0].Relations[1].Relation = PlayerRelationEnemy
	s.Players[2].Relations[0].Relation = PlayerRelationFriend
	return s
}

// ScenarioBattleMixedFleet pits two battle cruisers against one fleet mixing
// freighters, shielded scouts, jammed beam destroyers, and shield sappers.
func ScenarioBattleMixedFleet() TestScenario {
	s := scenarioBattle("Battle Mixed Fleet", DesignBattleCruiser, DesignArmoredFreighter)
	s.Players[0].Fleets[0].Quantity = 2
	s.Players[1].Designs = Designs(DesignArmoredFreighter, DesignShieldedScout, DesignJammedDefender, DesignStalwartSapper)
	s.Players[1].Fleets = []ScenarioFleet{{Name: "Mixed Fleet", Position: Vector{100, 100}, Tokens: []ScenarioShipToken{
		{Design: DesignArmoredFreighter.Name, Quantity: 5},
		{Design: DesignShieldedScout.Name, Quantity: 2},
		{Design: DesignJammedDefender.Name, Quantity: 3},
		{Design: DesignStalwartSapper.Name, Quantity: 4},
	}}}
	return s
}

// ScenarioProductionStarbases is shared by browser and game-logic tests.
func ScenarioProductionStarbases() TestScenario {
	return TestScenario{Name: "Production Starbases",
		Players: []ScenarioPlayer{{Player: &Player{
			Name:       "Player 1",
			UserID:     1,
			Race:       *NewRace().WithPRT(PP),
			TechLevels: TechLevel{Energy: 4, Weapons: 3, Propulsion: 3, Construction: 3, Electronics: 3, Biotechnology: 3},
		},
			Designs: Designs(ShipDesign{Name: "Fort",
				Hull: OrbitalFort.Name}, ShipDesign{Name: "Station",
				Hull: SpaceStation.Name}, ShipDesign{Name: "Flinger",
				Hull:  SpaceStation.Name,
				Slots: []ShipDesignSlot{{HullComponent: MassDriver5.Name, HullSlotIndex: 1, Quantity: 1}}}, DesignLongRangeScout)}},
		Planets: []ScenarioPlanet{{Name: "Planet 1",
			Owner:     1,
			Starbase:  "Fort",
			Cargo:     Cargo{Ironium: 1000, Boranium: 1000, Germanium: 1000, Colonists: 2500},
			Homeworld: true}}}
}

// SingleUnitScenario returns one planet and one scout.
func SingleUnitScenario() TestScenario {
	return TestScenario{
		Name: "Single Unit Game",
		Players: []ScenarioPlayer{
			{
				Designs: Designs(DesignLongRangeScout),
				Fleets:  []ScenarioFleet{{Design: "Long Range Scout", At: "Planet 1", Fuel: 300}},
			},
		},
		Planets: []ScenarioPlanet{{Name: "Planet 1", Owner: 1, Cargo: Cargo{Colonists: 2500}}},
	}
}

// TwoPlayerScenario returns one planet and one scout per player.
func TwoPlayerScenario() TestScenario {
	second := AIPlayer("Player 2")
	second.Designs = Designs(ShipDesign{Name: "Super Scout", Hull: Scout.Name, Slots: LongRangeScoutTestSlots})
	second.Fleets = []ScenarioFleet{{Name: "Long Range Scout #1", Design: "Super Scout", At: "Planet 2", EmptyFuel: true}}
	return TestScenario{
		Name: "Two Player Game",
		Players: []ScenarioPlayer{
			{
				Designs: Designs(DesignLongRangeScout),
				Fleets:  []ScenarioFleet{{Design: "Long Range Scout", At: "Planet 1", EmptyFuel: true}},
			},
			second,
		},
		Planets: []ScenarioPlanet{
			{Name: "Planet 1", Owner: 1, Cargo: Cargo{Colonists: 2500}},
			{Name: "Planet 2", Owner: 2, Position: Vector{X: 100}, Cargo: Cargo{Colonists: 2500}},
		},
	}
}
