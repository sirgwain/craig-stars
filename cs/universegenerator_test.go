//go:build !wasi && !wasm

package cs

import (
	"testing"

	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func TestGenerateUniverse(t *testing.T) {
	t.Run("Medium Packed", func(t *testing.T) {
		client := NewGamer()
		game := client.CreateGame(1, *NewGameSettings())
		game.Size = SizeMedium
		game.Density = DensityPacked
		game.GalaxyClumping = true

		numPlanets, err := game.Rules.GetNumPlanets(game.Size, game.Density)
		if err != nil {
			t.Error(err)
		}
		player := client.NewPlayer(1, *NewRace(), &game.Rules)
		players := []*Player{player}
		player.AIControlled = true
		player.Num = 1
		universe, err := client.GenerateUniverse(game, players)
		assert.NoError(t, err)

		assert.Equal(t, numPlanets, len(universe.Planets))
		assert.Greater(t, len(universe.Fleets), 0)
		assert.Greater(t, len(player.Designs), 0)
		assert.Greater(t, len(universe.Wormholes), 0)

		pmo := universe.GetPlayerMapObjects(player.Num)
		assert.Equal(t, 1, len(pmo.Planets))
		homeworld := pmo.Planets[0]
		assert.Equal(t, 25_000, homeworld.GetPopulation())
		assert.True(t, homeworld.Spec.HasStarbase)
	})
	t.Run("Normal", func(t *testing.T) {
		client := NewGamer()
		game := client.CreateGame(1, *NewGameSettings())

		numPlanets, err := game.Rules.GetNumPlanets(game.Size, game.Density)
		if err != nil {
			t.Error(err)
		}
		player := client.NewPlayer(1, *NewRace(), &game.Rules)
		players := []*Player{player}
		player.AIControlled = true
		player.Num = 1
		universe, err := client.GenerateUniverse(game, players)
		assert.NoError(t, err)

		assert.Equal(t, len(universe.Planets), numPlanets)
		assert.Greater(t, len(universe.Fleets), 0)
		assert.Greater(t, len(player.Designs), 0)
		assert.Greater(t, len(universe.Wormholes), 0)

		pmo := universe.GetPlayerMapObjects(player.Num)
		assert.Equal(t, 1, len(pmo.Planets))
		homeworld := pmo.Planets[0]
		assert.Equal(t, 25_000, homeworld.GetPopulation())
		assert.True(t, homeworld.Spec.HasStarbase)
	})

	t.Run("Multiple Players, multiple planets", func(t *testing.T) {
		client := NewGamer()
		game := client.CreateGame(1, *NewGameSettings())

		player1 := client.NewPlayer(1, *NewRace().WithPRT(IT), &game.Rules).WithAIControlled(true)
		player1.Num = 1
		player2 := client.NewPlayer(1, *NewRace(), &game.Rules).WithAIControlled(true)
		player2.Num = 2
		player3 := client.NewPlayer(1, *NewRace().WithPRT(PP), &game.Rules).WithAIControlled(true)
		player3.Num = 2

		players := []*Player{player1, player2, player3}
		universe, err := client.GenerateUniverse(game, players)
		assert.NoError(t, err)

		assert.Greater(t, len(universe.Fleets), 0)
		assert.Greater(t, len(player1.Designs), 0)
		assert.Greater(t, len(universe.Wormholes), 0)

		pmo := universe.GetPlayerMapObjects(player1.Num)
		assert.Equal(t, 2, len(pmo.Planets))
		assert.True(t, pmo.Planets[0].Spec.HasStarbase)
		assert.True(t, pmo.Planets[1].Spec.HasStarbase)
	})

	t.Run("Acc BBS Pop Test", func(t *testing.T) {
		client := NewGamer()
		game := client.CreateGame(1, *NewGameSettings().WithGameStartMode(GameStartModeAccBBS))

		numPlanets, err := game.Rules.GetNumPlanets(game.Size, game.Density)
		assert.NoError(t, err)

		hePlayer := client.NewPlayer(1, *NewRace().WithPRT(HE).WithGrowthRate(20).WithSpec(&game.Rules), &game.Rules).WithNum(1)
		ssPlayer := client.NewPlayer(2, *NewRace().WithPRT(SS).WithLRT(LSP).WithSpec(&game.Rules), &game.Rules).WithNum(2)
		itPlayer := client.NewPlayer(3, *NewRace().WithPRT(IT).WithGrowthRate(10).WithSpec(&game.Rules), &game.Rules).WithNum(3)
		sdPlayer := client.NewPlayer(4, *NewRace().WithPRT(SD).WithGrowthRate(1).WithSpec(&game.Rules), &game.Rules).WithNum(4)
		players := []*Player{hePlayer, ssPlayer, itPlayer, sdPlayer}

		universe, err := client.GenerateUniverse(game, players)
		assert.NoError(t, err)

		assert.Equal(t, len(universe.Planets), numPlanets)

		popPerPlayerPerPlanet := [][]int{
			{225_000},        // 40% HE; 9x pop
			{70_000},         // 15% LSP SS; 2.8x pop
			{60_000, 30_000}, // 10% IT; 3x pop
			{30_000},         // 1% SD; 1.2x pop
		}

		for playerNum := 1; playerNum <= len(players); playerNum++ {
			planets := universe.getPlanets(playerNum)
			popPerPlanet := popPerPlayerPerPlanet[playerNum-1]
			assert.Equal(t, len(popPerPlanet), len(planets))
			for _, p := range planets {
				if p.Homeworld {
					assert.Equal(t, popPerPlanet[0], p.GetPopulation())
				} else {
					assert.Equal(t, popPerPlanet[1], p.GetPopulation())

				}
			}
		}
	})

	t.Run("Max Mode", func(t *testing.T) {
		client := NewGamer()
		game := client.CreateGame(1, *NewGameSettings().WithGameStartMode(GameStartModeMax))

		numPlanets, err := game.Rules.GetNumPlanets(game.Size, game.Density)
		if err != nil {
			t.Error(err)
		}
		player := client.NewPlayer(1, *NewRace(), &game.Rules)
		players := []*Player{player}
		player.AIControlled = true
		player.Num = 1
		universe, err := client.GenerateUniverse(game, players)
		assert.NoError(t, err)

		assert.Equal(t, len(universe.Planets), numPlanets)
		assert.Greater(t, len(universe.Fleets), 0)
		assert.Greater(t, len(player.Designs), 0)
		assert.Greater(t, len(universe.Wormholes), 0)

		assert.Equal(t, player.TechLevels, TechLevel{rules.MaxTechLevel, rules.MaxTechLevel, rules.MaxTechLevel, rules.MaxTechLevel, rules.MaxTechLevel, rules.MaxTechLevel})

		pmo := universe.GetPlayerMapObjects(player.Num)
		assert.Equal(t, 1, len(pmo.Planets))
		homeworld := pmo.Planets[0]
		assert.Equal(t, homeworld.GetPopulation(), homeworld.Spec.MaxPopulation)
		assert.True(t, homeworld.Spec.HasStarbase)
		assert.Equal(t, homeworld.Factories, homeworld.Spec.MaxPossibleFactories)
		assert.Equal(t, homeworld.Mines, homeworld.Spec.MaxPossibleMines)

		// make sure all the fleet specs are ok - costs/scanner ranges are all as they should be
		for _, fleet := range pmo.Fleets {
			assert.NotNil(t, fleet)
			c := NewCostCalculator()
			design := fleet.Tokens[0].design
			hull := rules.techs.GetHull(design.Hull)
			design.Spec.computeScanRanges(&rules, player.Race.Spec.ScannerSpec, player.TechLevels, design, hull) // updates design scanrange but not fleet scan range
			assert.Equal(t, design.Spec.ScanRange, fleet.Spec.ScanRange)
			assert.Equal(t, design.Spec.ScanRangePen, fleet.Spec.ScanRangePen)
			calcCost, err := c.GetDesignCost(&rules, player.TechLevels, player.Race.Spec, design)
			assert.NoError(t, err)
			assert.Equal(t, calcCost, design.Spec.Cost)
		}
	})
}

func Test_universeGenerator_assignRaceStartingPointBonuses(t *testing.T) {
	type args struct {
		race        *Race
		extraPoints int
		pointsType  SpendLeftoverPointsOn
		planet      *Planet
	}
	tests := []struct {
		name string
		args args
		want *Planet
	}{
		{
			name: "10 points into factories, 2 facts",
			args: args{
				race:        NewRace().WithSpec(&rules),
				extraPoints: 10,
				pointsType:  SpendLeftoverPointsOnFactories,
			},
			want: NewPlanet().WithFactories(2),
		},
		{
			name: "too few for factories; spills over",
			args: args{
				race:        NewRace().WithSpec(&rules),
				extraPoints: 2,
				pointsType:  SpendLeftoverPointsOnFactories,
			},
			want: NewPlanet().WithCargo(Cargo{12, 5, 5, 0}), // yes this is weird i know del with it
		},
		{
			name: "2 points into mines; can't use",
			args: args{
				race:        NewRace().WithPRT(AR).WithSpec(&rules),
				extraPoints: 2,
				pointsType:  SpendLeftoverPointsOnMines,
			},
			want: NewPlanet().WithCargo(Cargo{12, 5, 5, 0}),
		},
		{
			name: "45 points into defenses; rounds half up",
			args: args{
				race:        NewRace().WithSpec(&rules),
				extraPoints: 45,
				pointsType:  SpendLeftoverPointsOnDefenses,
			},
			// first costs half as much pts
			want: NewPlanet().WithDefenses(5),
		},
		{
			name: "33 points into concentration; some already",
			args: args{
				race:        NewRace().WithSpec(&rules),
				extraPoints: 33,
				pointsType:  SpendLeftoverPointsOnMineralConcentrations,
				planet:      NewPlanet().WithMineralConcentration(Mineral{35, 35, 37}),
			},
			want: NewPlanet().WithMineralConcentration(Mineral{59, 43, 45}),
			// 16 points (33/2) into lowest, then 8 (16/2) into all
		},
		{
			name: "36 points into surface minerals with some cargo",
			args: args{
				race:        NewRace().WithSpec(&rules),
				extraPoints: 36,
				pointsType:  SpendLeftoverPointsOnSurfaceMinerals,
				planet:      NewPlanet().WithCargo(Cargo{1, 2, 3, 220}),
			},
			want: NewPlanet().WithCargo(Cargo{181, 92, 93, 220}),
		},
		{
			name: "invalid starting point type; uses surface mins",
			args: args{
				race:        NewRace().WithSpec(&rules),
				extraPoints: 1,
				pointsType:  "BANANANANA",
			},
			want: NewPlanet().WithCargo(Cargo{5, 2, 2, 0}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ug := universeGenerator{FullGame: &FullGame{Game: &Game{Rules: rules}}}
			if tt.args.planet == nil {
				tt.args.planet = NewPlanet()
			}
			ug.assignRaceStartingPointBonuses(tt.args.race, tt.args.planet, tt.args.extraPoints, tt.args.pointsType)

			test.CompareAsJSON(t, tt.args.planet, tt.want)
		})
	}
}

func Test_universeGenerator_getStartingStarbaseDesigns(t *testing.T) {
	tests := []struct {
		name   string
		player *Player
		want   []ShipDesign
	}{
		{
			name:   "Test JoaT Starter Base",
			player: testPlayer().WithNum(1),
			want: []ShipDesign{
				*NewShipDesign(1, 1).
					WithName(joatSpec().StartingPlanets[0].StarbaseDesignName).
					WithHull(joatSpec().StartingPlanets[0].StarbaseHull).
					WithPurpose(ShipDesignPurposeStarbase).
					WithSlots([]ShipDesignSlot{
						{HullComponent: Laser.Name, HullSlotIndex: 2, Quantity: 8},
						{HullComponent: MoleSkinShield.Name, HullSlotIndex: 3, Quantity: 8},
						{HullComponent: Laser.Name, HullSlotIndex: 4, Quantity: 8},
						{HullComponent: MoleSkinShield.Name, HullSlotIndex: 5, Quantity: 8},
						{HullComponent: MoleSkinShield.Name, HullSlotIndex: 6, Quantity: 8},
						{HullComponent: Laser.Name, HullSlotIndex: 8, Quantity: 8},
						{HullComponent: Laser.Name, HullSlotIndex: 10, Quantity: 8},
						{HullComponent: MoleSkinShield.Name, HullSlotIndex: 12, Quantity: 8},
					}),
			},
		},
		{
			name:   "PP Starter Base designs",
			player: NewPlayer(2, NewRace().WithPRT(PP).WithSpec(&rules)).WithNum(2),
			want: []ShipDesign{
				*NewShipDesign(2, 1).
					WithName(ppSpec().StartingPlanets[0].StarbaseDesignName).
					WithHull(ppSpec().StartingPlanets[0].StarbaseHull).
					WithPurpose(ShipDesignPurposeStarbase).
					WithSlots([]ShipDesignSlot{
						{HullComponent: MassDriver5.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: Laser.Name, HullSlotIndex: 2, Quantity: 8},
						{HullComponent: MoleSkinShield.Name, HullSlotIndex: 3, Quantity: 8},
						{HullComponent: Laser.Name, HullSlotIndex: 4, Quantity: 8},
						{HullComponent: MoleSkinShield.Name, HullSlotIndex: 5, Quantity: 8},
						{HullComponent: MoleSkinShield.Name, HullSlotIndex: 6, Quantity: 8},
						{HullComponent: Laser.Name, HullSlotIndex: 8, Quantity: 8},
						{HullComponent: Laser.Name, HullSlotIndex: 10, Quantity: 8},
						{HullComponent: MoleSkinShield.Name, HullSlotIndex: 12, Quantity: 8},
					}),
				*NewShipDesign(2, 2).
					WithName(ppSpec().StartingPlanets[1].StarbaseDesignName).
					WithHull(ppSpec().StartingPlanets[1].StarbaseHull).
					WithPurpose(ShipDesignPurposePacketThrower).
					WithSlots([]ShipDesignSlot{
						{HullComponent: MassDriver5.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: Laser.Name, HullSlotIndex: 2, Quantity: 6},
						{HullComponent: MoleSkinShield.Name, HullSlotIndex: 3, Quantity: 6},
						{HullComponent: Laser.Name, HullSlotIndex: 4, Quantity: 6},
						{HullComponent: MoleSkinShield.Name, HullSlotIndex: 5, Quantity: 6},
					}),
			},
		},
		{
			name:   "IT Base Designs",
			player: NewPlayer(1, NewRace().WithPRT(IT).WithSpec(&rules)).WithNum(15),
			want: []ShipDesign{
				*NewShipDesign(15, 1).
					WithName(itSpec().StartingPlanets[0].StarbaseDesignName).
					WithHull(itSpec().StartingPlanets[0].StarbaseHull).
					WithPurpose(ShipDesignPurposeStarbase).
					WithSlots([]ShipDesignSlot{
						{HullComponent: Stargate100_250.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: Laser.Name, HullSlotIndex: 2, Quantity: 8},
						{HullComponent: MoleSkinShield.Name, HullSlotIndex: 3, Quantity: 8},
						{HullComponent: Laser.Name, HullSlotIndex: 4, Quantity: 8},
						{HullComponent: MoleSkinShield.Name, HullSlotIndex: 5, Quantity: 8},
						{HullComponent: MoleSkinShield.Name, HullSlotIndex: 6, Quantity: 8},
						{HullComponent: Laser.Name, HullSlotIndex: 8, Quantity: 8},
						{HullComponent: Laser.Name, HullSlotIndex: 10, Quantity: 8},
						{HullComponent: MoleSkinShield.Name, HullSlotIndex: 12, Quantity: 8},
					}),
				*NewShipDesign(15, 2).
					WithName(itSpec().StartingPlanets[1].StarbaseDesignName).
					WithHull(itSpec().StartingPlanets[1].StarbaseHull).
					WithPurpose(ShipDesignPurposeStargater).
					WithSlots([]ShipDesignSlot{
						{HullComponent: Stargate100_250.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: Laser.Name, HullSlotIndex: 2, Quantity: 6},
						{HullComponent: MoleSkinShield.Name, HullSlotIndex: 3, Quantity: 6},
						{HullComponent: Laser.Name, HullSlotIndex: 4, Quantity: 6},
						{HullComponent: MoleSkinShield.Name, HullSlotIndex: 5, Quantity: 6},
					}),
			},
		},
		{
			name:   "AR Base Designs",
			player: NewPlayer(1, NewRace().WithPRT(AR).WithSpec(&rules)).WithNum(4),
			want: []ShipDesign{
				*NewShipDesign(4, 1).
					WithName(arSpec().StartingPlanets[0].StarbaseDesignName).
					WithHull(arSpec().StartingPlanets[0].StarbaseHull).
					WithPurpose(ShipDesignPurposeStarbase).
					WithSlots([]ShipDesignSlot{
						{HullComponent: Laser.Name, HullSlotIndex: 2, Quantity: 8},
						{HullComponent: MoleSkinShield.Name, HullSlotIndex: 3, Quantity: 8},
						{HullComponent: Laser.Name, HullSlotIndex: 4, Quantity: 8},
						{HullComponent: MoleSkinShield.Name, HullSlotIndex: 5, Quantity: 8},
						{HullComponent: MoleSkinShield.Name, HullSlotIndex: 6, Quantity: 8},
						{HullComponent: Laser.Name, HullSlotIndex: 8, Quantity: 8},
						{HullComponent: Laser.Name, HullSlotIndex: 10, Quantity: 8},
						{HullComponent: MoleSkinShield.Name, HullSlotIndex: 12, Quantity: 8},
					}),
				*NewShipDesign(4, 2).
					WithName("Starter Colony").
					WithHull("Orbital Fort").
					WithPurpose(ShipDesignPurposeStarterColony).
					WithCannotDelete(true),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ug := universeGenerator{
				FullGame: &FullGame{
					Game: &Game{
						Size: SizeSmall,
					},
				},
			}

			tt.player.Name = tt.name
			got := ug.createStartingStarbaseDesigns(&StaticTechStore, tt.player, 1)

			test.CompareAsJSON(t, got, tt.want)
		})
	}
}
