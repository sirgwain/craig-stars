package cs

import (
	"math"
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/test"
)

func Test_completionEstimate_GetYearsToBuild(t *testing.T) {
	type args struct {
		item                   ProductionQueueItem
		costPerItem            Cost
		mineralsOnHand         Mineral
		yearlyAvailableToSpend Cost
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "infinite build time",
			args: args{
				item: ProductionQueueItem{
					Type:     QueueItemTypeFactory,
					Quantity: 1,
				},
				costPerItem:            Cost{Resources: 1},
				yearlyAvailableToSpend: Cost{Germanium: 99},
			},
			want: math.MaxInt,
		},
		{
			name: "2 items, 1 year each",
			args: args{
				item: ProductionQueueItem{
					Quantity: 2,
				},
				costPerItem:            Cost{Ironium: 1},
				yearlyAvailableToSpend: Cost{Ironium: 1},
			},
			want: 2,
		},
		{
			name: "2 items, 2 years each",
			args: args{
				item: ProductionQueueItem{
					Quantity: 2,
				},
				costPerItem:            Cost{Ironium: 2},
				yearlyAvailableToSpend: Cost{Ironium: 1},
			},
			want: 4,
		},
		{
			name: "partially done; 2 years to finish",
			args: args{
				item: ProductionQueueItem{
					Quantity:  1,
					Allocated: Cost{Ironium: 2},
				},
				costPerItem:            Cost{Ironium: 4},
				yearlyAvailableToSpend: Cost{Ironium: 1},
			},
			want: 2,
		},
		{
			name: "one item fully built",
			args: args{
				item: ProductionQueueItem{
					Quantity:  1,
					Allocated: Cost{Ironium: 1},
				},
				costPerItem:            Cost{Ironium: 4},
				mineralsOnHand:         Mineral{Ironium: 3},
				yearlyAvailableToSpend: Cost{Ironium: 1},
			},
			want: 0,
		},
		{
			name: "2 items, resource shortage",
			args: args{
				item: ProductionQueueItem{
					Quantity: 2,
				},
				costPerItem:            Cost{Ironium: 4, Resources: 1},
				mineralsOnHand:         Mineral{Ironium: 100},
				yearlyAvailableToSpend: Cost{Ironium: 1, Resources: 1},
			},
			want: 2,
		},
		{
			name: "5 mines, 2 resources per year",
			args: args{
				item: ProductionQueueItem{
					Type:     QueueItemTypeAutoMines,
					Quantity: 5,
				},
				costPerItem:            Cost{Resources: 5},
				mineralsOnHand:         Mineral{},
				yearlyAvailableToSpend: Cost{Resources: 2},
			},
			want: 13, // 25/2
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &completionEstimate{}
			if got := e.GetYearsToBuild(tt.args.item, tt.args.costPerItem, tt.args.mineralsOnHand, tt.args.yearlyAvailableToSpend); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("completionEstimate.GetYearsToBuild() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_completionEstimate_GetProductionWithEstimates(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		player := NewPlayer(1, NewRace().WithSpec(&rules)).withSpec(&rules)
		type args struct {
			items  []ProductionQueueItem
			planet *Planet
		}
		tests := []struct {
			name                  string
			args                  args
			want                  []ProductionQueueItem
			wantLeftoverResources int
			wantErr               bool
		}{
			{
				name: "3 half done ships",
				args: args{
					items: []ProductionQueueItem{
						{
							Type:   QueueItemTypeShipToken,
							design: testLongRangeScoutDesign(1).WithSpec(&rules, player),
							// Total cost: {54, 6, 24, 69}
							Quantity:  3,
							Allocated: Cost{15, 1, 3, 4},
							// 39 iron left to build
						},
					},
					planet: NewPlanet().WithCargo(Cargo{6, 100, 100, 1000}).
						WithMines(10). // 10kT per year
						WithFactories(1000).
						WithContributesOnlyLeftoverToResearch(true),
				},
				want: []ProductionQueueItem{
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne:     1,
							YearsToBuildAll:     4,
							YearsToSkipOrCancel: Infinite,
						},
						Type:      QueueItemTypeShipToken,
						design:    testLongRangeScoutDesign(1).WithSpec(&rules, player),
						Quantity:  3,
						Allocated: Cost{15, 1, 3, 4},
					},
				},
				wantLeftoverResources: 165,
				wantErr:               false,
			},
			{
				name: "one item, never completes",
				args: args{
					items: []ProductionQueueItem{
						{
							Type:     QueueItemTypeFactory,
							Quantity: 1,
						},
					},
					planet: NewPlanet().WithCargo(Cargo{0, 0, 0, 10}),
				},
				want: []ProductionQueueItem{
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne:     Infinite,
							YearsToBuildAll:     Infinite,
							YearsToSkipOrCancel: Infinite,
						},
						Type:     QueueItemTypeFactory,
						Quantity: 1,
					},
				},
				wantLeftoverResources: 1,
				wantErr:               false,
			},
			{
				name: "one item, 2 years to go",
				args: args{
					items: []ProductionQueueItem{
						{
							Type:      QueueItemTypeFactory,
							Quantity:  1,
							Allocated: Cost{Resources: 5},
						},
					},
					planet: NewPlanet().WithCargo(Cargo{0, 0, 4, 30}), // enough minerals to finish
				},
				want: []ProductionQueueItem{
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne:     2,
							YearsToBuildAll:     2,
							YearsToSkipOrCancel: Infinite,
						},
						Type:      QueueItemTypeFactory,
						Quantity:  1,
						Allocated: Cost{Resources: 5},
					},
				},
				wantLeftoverResources: 0,
				wantErr:               false,
			},
			{
				name: "two items, first completes this year, second takes 2 years",
				args: args{
					items: []ProductionQueueItem{
						{
							Type:     QueueItemTypeMine,
							Quantity: 1,
						},
						{
							Type:     QueueItemTypeFactory,
							Quantity: 2,
						},
					},
					planet: NewPlanet().WithMines(2).WithCargo(Cargo{0, 0, 0, 50}),
					// 5 resources/2kT per year
				},
				want: []ProductionQueueItem{
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne:     1,
							YearsToBuildAll:     1,
							YearsToSkipOrCancel: Infinite,
						},
						Type:     QueueItemTypeMine,
						Quantity: 1,
					},
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne:     3,
							YearsToBuildAll:     5,
							YearsToSkipOrCancel: Infinite,
						},
						Type:     QueueItemTypeFactory,
						Quantity: 2,
					},
				},
				wantErr: false,
			},
			{
				name: "two items, lots of minerals on hand, low resources",
				args: args{
					items: []ProductionQueueItem{
						{
							Type:     QueueItemTypeFactory,
							Quantity: 1, // 10 res; +1 res/yr
						},
						{
							Type:     QueueItemTypeMine,
							Quantity: 2, // 10 res total
						},
					},
					planet: NewPlanet().WithCargo(Cargo{100, 100, 100, 10}),
					// Pop: 1000 → 1115 → 1280 → 1445 → 1655 → 1895 → 2165 → 2480 → 3260 → 3740
				},
				want: []ProductionQueueItem{
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne:     8,
							YearsToBuildAll:     8,
							YearsToSkipOrCancel: Infinite,
						},
						Type:     QueueItemTypeFactory,
						Quantity: 1,
					},
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne:     10,
							YearsToBuildAll:     11,
							YearsToSkipOrCancel: Infinite,
						},
						Type:     QueueItemTypeMine,
						Quantity: 2,
					},
				},
				wantErr: false,
			},
			{
				name: "5 auto factories, then 10 auto mines",
				args: args{
					items: []ProductionQueueItem{
						{
							Type:     QueueItemTypeAutoFactories,
							Quantity: 5,
						},
						{
							Type:     QueueItemTypeAutoMines,
							Quantity: 10,
						},
					},
					planet: NewPlanet().WithCargo(Cargo{0, 0, 8, 350}),
				},
				want: []ProductionQueueItem{
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne:     1,
							YearsToBuildAll:     6,
							YearsToSkipOrCancel: 2,
						},
						Type:     QueueItemTypeAutoFactories,
						Quantity: 5,
					},
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne:     1, // we build some mines in the first year due to lack of germ
							YearsToBuildAll:     9, // we finish them some time later
							YearsToSkipOrCancel: Infinite,
						},
						Type:     QueueItemTypeAutoMines,
						Quantity: 10,
					},
				},
				wantErr: false,
			},
			{
				name: "Test later year planet with high everything",
				args: args{
					items: []ProductionQueueItem{
						{
							Type:     QueueItemTypeAutoMinTerraform,
							Quantity: 1,
						},
						{
							Type:     QueueItemTypeTerraformEnvironment,
							Quantity: 1,
						},
						{
							Type:     QueueItemTypeAutoFactories,
							Quantity: 100,
						},
						{
							Type:     QueueItemTypeAutoMines,
							Quantity: 100,
						},
					},
					planet: NewPlanet().WithCargo(Cargo{2000, 2000, 2000, 7000}).
						WithMines(700).WithFactories(700),
				},
				want: []ProductionQueueItem{
					// we skip/cancel the terraforming and easily finish everything else
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne:     Infinite,
							YearsToBuildAll:     Infinite,
							YearsToSkipOrCancel: 1, // not canceled due to auto
						},
						Type:     QueueItemTypeAutoMinTerraform,
						Quantity: 1,
					},
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne:     Infinite,
							YearsToBuildAll:     Infinite,
							YearsToSkipOrCancel: 1,
						},
						Type:     QueueItemTypeTerraformEnvironment,
						Quantity: 1,
					},
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne:     1,
							YearsToBuildAll:     1,
							YearsToSkipOrCancel: Infinite,
						},
						Type:     QueueItemTypeAutoFactories,
						Quantity: 100,
					},
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne:     1,
							YearsToBuildAll:     1,
							YearsToSkipOrCancel: Infinite,
						},
						Type:     QueueItemTypeAutoMines,
						Quantity: 100,
					},
				},
				wantLeftoverResources: 710,
				wantErr:               false,
			},
			{
				name: "auto factories when have more than max",
				args: args{
					items: []ProductionQueueItem{
						{
							Type:     QueueItemTypeAutoFactories,
							Quantity: 100,
						},
						{
							Type:     QueueItemTypeAutoMines,
							Quantity: 100,
						},
					},
					planet: NewPlanet().WithCargo(Cargo{2000, 2000, 2000, 1000}).
						WithFactories(200).WithMines(200), // have 200, only operate 100ypeAutoFactories,
				},
				want: []ProductionQueueItem{
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne:     Infinite,
							YearsToBuildAll:     Infinite,
							YearsToSkipOrCancel: 1,
						},
						Type:     QueueItemTypeAutoFactories,
						Quantity: 100,
					},
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne:     Infinite,
							YearsToBuildAll:     Infinite,
							YearsToSkipOrCancel: 1,
						},
						Type:     QueueItemTypeAutoMines,
						Quantity: 100,
					},
				},
				wantLeftoverResources: 170,
				wantErr:               false,
			},
			{
				name: "new planet 2500 cols",
				args: args{
					items: []ProductionQueueItem{
						{
							Type:     QueueItemTypeAutoFactories,
							Quantity: 250,
						},
						{
							Type:     QueueItemTypeAutoMines,
							Quantity: 250,
						},
					},
					// planet from OG game
					planet: NewPlanet().WithCargo(Cargo{18, 6, 17, 25}).
						WithMineralConcentration(Mineral{95, 75, 79}).WithHab(Hab{51, 45, 46}),
				},
				want: []ProductionQueueItem{
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne:     13,
							YearsToBuildAll:     45,
							YearsToSkipOrCancel: 9,
						},
						Type:     QueueItemTypeAutoFactories,
						Quantity: 250,
					},
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne:     10,
							YearsToBuildAll:     47,
							YearsToSkipOrCancel: Infinite,
						},
						Type:     QueueItemTypeAutoMines,
						Quantity: 250,
					},
				},
				wantLeftoverResources: 0,
				wantErr:               false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				e := NewCompletionEstimator()

				planet := tt.args.planet
				planet.Hab = Hab{50, 50, 50}                         // perfect hab
				planet.MineralConcentration = Mineral{100, 100, 100} // perfect concentration for 1kT per mine output
				planet.PlayerNum = 1
				planet.Spec = computePlanetSpec(&rules, player, planet)
				planet.ProductionQueue = tt.args.items

				got, gotLeftover, err := e.GetProductionWithEstimates(&rules, player, *planet)
				if (err != nil) != tt.wantErr {
					if tt.wantErr {
						t.Fatalf("CompletionEstimator.GetProductionWithEstimates() did not return error when expected")
					} else {
						t.Fatalf("CompletionEstimator.GetProductionWithEstimates() errored unexpectedly; err = \n%v", err)
					}
				}

				if gotLeftover != tt.wantLeftoverResources {
					t.Errorf("GetProductionWithEstimates() leftover = %d, wantLeftover %d", gotLeftover, tt.wantLeftoverResources)
				}

				test.CompareAsJSON(t, got, tt.want)
			})
		}
	})

	t.Run("AR", func(t *testing.T) {
		player := NewPlayer(1, NewRace().WithPRT(AR).WithSpec(&rules)).
			WithTechLevels(TechLevel{Energy: 1}).
			WithNum(1).withSpec(&rules)

		santaMaria := NewShipDesign(1, 1).WithHull(ColonyShip.Name).WithSlots([]ShipDesignSlot{
			{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
			{HullComponent: OrbitalConstructionModule.Name, HullSlotIndex: 2, Quantity: 1},
		}).WithSpec(&rules, player)
		potatoBug := NewShipDesign(1, 2).WithHull(MidgetMiner.Name).WithSlots([]ShipDesignSlot{
			{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
			{HullComponent: RoboMidgetMiner.Name, HullSlotIndex: 2, Quantity: 2},
		}).WithSpec(&rules, player)

		tests := []struct {
			name                  string
			items                 []ProductionQueueItem
			cargo                 Cargo
			want                  []ProductionQueueItem
			wantLeftoverResources int
		}{
			{
				name: "santa maria",
				items: []ProductionQueueItem{
					{
						Type:     QueueItemTypeShipToken,
						Quantity: 2,
						design:   santaMaria,
						// {33, 15, 31, 43}
					},
				},
				cargo: Cargo{100, 100, 100, 250}, // 50 res per year
				want: []ProductionQueueItem{
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne:     1,
							YearsToBuildAll:     2,
							YearsToSkipOrCancel: -1,
						},
						Type:     QueueItemTypeShipToken,
						Quantity: 2,
						design:   santaMaria,
					},
				},
				wantLeftoverResources: 0,
			},
			{
				name: "potato bug",
				items: []ProductionQueueItem{
					{
						Type:     QueueItemTypeShipToken,
						Quantity: 2,
						design:   potatoBug,
						// {41, 0, 12, 123}
					},
				},
				cargo: Cargo{100, 100, 100, 100},
				want: []ProductionQueueItem{
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne:     4,
							YearsToBuildAll:     7,
							YearsToSkipOrCancel: -1,
						},
						Type:     QueueItemTypeShipToken,
						Quantity: 2,
						design:   potatoBug,
					},
				},
				wantLeftoverResources: 0,
			},
			{
				name: "alchemy",
				items: []ProductionQueueItem{
					{
						Type:     QueueItemTypeMineralAlchemy,
						Quantity: 1,
					},
				},
				cargo: Cargo{100, 100, 100, 100}, // 32 res/yr
				want: []ProductionQueueItem{
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne:     3,
							YearsToBuildAll:     3,
							YearsToSkipOrCancel: -1,
						},
						Type:     QueueItemTypeMineralAlchemy,
						Quantity: 1,
					},
				},
				wantLeftoverResources: 0,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				e := NewCompletionEstimator()

				planet := NewPlanet().WithCargo(tt.cargo).WithPlayerNum(1).
					WithContributesOnlyLeftoverToResearch(true)
				planet.Starbase = testSpaceStation(player, planet)
				planet.Hab = Hab{50, 50, 50}                         // perfect hab
				planet.MineralConcentration = Mineral{100, 100, 100} // perfect concentration for 1kT per mine output
				planet.Spec = computePlanetSpec(&rules, player, planet)

				// compute specs for designs in queue
				for _, item := range tt.items {
					if item.design != nil {
						item.design = item.design.WithSpec(&rules, player)
					}
				}
				for _, item := range tt.want {
					if item.design != nil {
						item.design = item.design.WithSpec(&rules, player)
					}
				}
				planet.ProductionQueue = tt.items

				got, gotLeftover, err := e.GetProductionWithEstimates(&rules, player, *planet)
				if err != nil {
					t.Fatalf("CompletionEstimator.GetProductionWithEstimates() errored unexpectedly; err = \n%v", err)
				}

				if gotLeftover != tt.wantLeftoverResources {
					t.Errorf("CompletionEstimator.GetProductionWithEstimates() leftover = %v, wantLeftover %v", gotLeftover, tt.wantLeftoverResources)
				}

				test.CompareAsJSON(t, got, tt.want)

			})
		}
	})

	t.Run("Packets", func(t *testing.T) {
		player := NewPlayer(1, NewRace().WithSpec(&rules)).withSpec(&rules)
		driverBase := NewShipDesign(1, 1).WithHull(SpaceStation.Name).WithName("Driver Base").
			WithSlots([]ShipDesignSlot{
				{HullComponent: MassDriver5.Name, HullSlotIndex: 1, Quantity: 1},
			}).WithSpec(&rules, player)
		noDriverBase := NewShipDesign(1, 1).WithHull(SpaceStation.Name).WithName("No Driver Base").
			WithSpec(&rules, player)

		tests := []struct {
			name     string
			starbase *ShipDesign
			items    []ProductionQueueItem
			cargo    Cargo
			want     []ProductionQueueItem
		}{
			{
				name: "Packet without driver",
				items: []ProductionQueueItem{
					{
						Type:     QueueItemTypeIroniumMineralPacket,
						Quantity: 1,
					},
				},
				cargo:    Cargo{1000, 1000, 1000, 10_000},
				starbase: noDriverBase,
				want: []ProductionQueueItem{
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne:     Infinite,
							YearsToBuildAll:     Infinite,
							YearsToSkipOrCancel: 1,
						},
						Type:     QueueItemTypeIroniumMineralPacket,
						Quantity: 1,
					},
				},
			},
			{
				name: "Built new base between packets",
				items: []ProductionQueueItem{
					{
						Type:     QueueItemTypeIroniumMineralPacket,
						Quantity: 1,
					},
					{
						Type:     QueueItemTypeStarbase,
						design:   driverBase, // 70 res to upgrade
						Quantity: 1,
					},
					{
						Type:     QueueItemTypeMixedMineralPacket,
						Quantity: 2,
					},
				},
				cargo:    Cargo{1000, 1000, 1000, 800}, // 80 resources/yr, enough for the first 2 items
				starbase: noDriverBase,
				// First packet canceled due to being built when we had no driver;
				// 2nd packet finished normally due to being built after the driver base finishes
				want: []ProductionQueueItem{
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne:     Infinite,
							YearsToBuildAll:     Infinite,
							YearsToSkipOrCancel: 1,
						},
						Type:     QueueItemTypeIroniumMineralPacket,
						Quantity: 1,
					},
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne:     1,
							YearsToBuildAll:     1,
							YearsToSkipOrCancel: Infinite,
						},
						Type:     QueueItemTypeStarbase,
						design:   noDriverBase,
						Quantity: 1,
					},
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne:     1,
							YearsToBuildAll:     2,
							YearsToSkipOrCancel: Infinite,
						},
						Type:     QueueItemTypeMixedMineralPacket,
						Quantity: 2,
					},
				},
			},
			{
				name: "2 years of resources",
				items: []ProductionQueueItem{
					{
						Type:     QueueItemTypeIroniumMineralPacket,
						Quantity: 1,
					},
				},
				cargo:    Cargo{1000, 1000, 1000, 50}, // 5 res/yr
				starbase: driverBase,
				want: []ProductionQueueItem{
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne:     2,
							YearsToBuildAll:     2,
							YearsToSkipOrCancel: Infinite,
						},
						Type:     QueueItemTypeIroniumMineralPacket,
						Quantity: 1,
					},
				},
			},
			{
				name: "Alchemy blocks queue",
				items: []ProductionQueueItem{
					{
						Type:     QueueItemTypeMineralAlchemy,
						Quantity: 1,
					},
					{
						Type:     QueueItemTypeBoraniumMineralPacket,
						Quantity: 2,
					},
				},
				cargo:    Cargo{1000, 1000, 1000, 320}, // 32 --> 36 --> 42 res/yr; 110 in first 3 yrs
				starbase: driverBase,
				want: []ProductionQueueItem{
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne:     3,
							YearsToBuildAll:     3,
							YearsToSkipOrCancel: Infinite,
						},
						Type:     QueueItemTypeMineralAlchemy,
						Quantity: 1,
					},
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne:     3,
							YearsToBuildAll:     4,
							YearsToSkipOrCancel: Infinite,
						},
						Type:     QueueItemTypeBoraniumMineralPacket,
						Quantity: 2,
					},
				},
			},
			{
				name: "Packets built this and next year",
				items: []ProductionQueueItem{
					{
						Type:     QueueItemTypeIroniumMineralPacket,
						Quantity: 9,
					},
					{
						Type:     QueueItemTypeBoraniumMineralPacket,
						Quantity: 4,
					},
				},
				cargo:    Cargo{2000, 2000, 2000, 500}, // 50 res/yr, enough for 5 packetd
				starbase: driverBase,
				want: []ProductionQueueItem{
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne:     1,
							YearsToBuildAll:     2,
							YearsToSkipOrCancel: Infinite,
						},
						Type:     QueueItemTypeIroniumMineralPacket,
						Quantity: 9,
					},
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne:     2,
							YearsToBuildAll:     3,
							YearsToSkipOrCancel: Infinite,
						},
						Type:     QueueItemTypeBoraniumMineralPacket,
						Quantity: 4,
					},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				e := NewCompletionEstimator()

				player.Name = tt.name
				planet := NewPlanet().WithCargo(tt.cargo).
					WithContributesOnlyLeftoverToResearch(true)
				planet.PacketTargetNum = 33  // give packet target
				planet.Hab = Hab{50, 50, 50} // perfect hab
				planet.MineralConcentration = Mineral{100, 100, 100}
				planet.PlayerNum = 1
				if tt.starbase != nil {
					player.Designs = append(player.Designs, tt.starbase)
					base := newStarbase(player, planet, tt.starbase, "Old Base")
					base.Spec = ComputeFleetSpec(&rules, player, base)
					planet.Starbase = base
				}

				planet.Spec = computePlanetSpec(&rules, player, planet)
				planet.ProductionQueue = tt.items

				got, _, err := e.GetProductionWithEstimates(&rules, player, *planet)
				if err != nil {
					t.Fatalf("CompletionEstimator.GetProductionWithEstimates() errored unexpectedly; err = \n%v", err)
				}

				test.CompareAsJSON(t, got, tt.want)
			})
		}
	})

	// TODO: Fix this to be consistent with base game...?
	t.Run("HE 6%", func(t *testing.T) {
		race := NewRace().WithPRT(HE).WithLRT(OBRM).
			withImmuneGrav(true).
			withImmuneRad(true).
			withImmuneTemp(true).
			WithGrowthRate(6)
		// 12/9/22/3 10/3/19
		race.MineCost = 3
		race.NumMines = 19
		race.FactoryOutput = 12
		race.FactoryCost = 9
		race.NumFactories = 22
		race.FactoriesCostLess = true
		race = race.WithSpec(&rules)

		player := NewPlayer(1, race).withSpec(&rules)
		planet := NewPlanet().WithCargo(Cargo{44, 24, 33, 30}).WithContributesOnlyLeftoverToResearch(true)
		planet.Spec = computePlanetSpec(&rules, player, planet)
		planet.ProductionQueue = []ProductionQueueItem{
			{
				Type:     QueueItemTypeAutoFactories,
				Quantity: 100,
			},
			{
				Type:     QueueItemTypeAutoMines,
				Quantity: 100,
			},
			{
				Type:     QueueItemTypeAutoDefenses,
				Quantity: 10,
			},
		}

		wantQueue := []ProductionQueueItem{
			{
				QueueItemCompletionEstimate: QueueItemCompletionEstimate{
					YearsToBuildOne:     10,
					YearsToBuildAll:     32,
					YearsToSkipOrCancel: Infinite,
				},
				Type:     QueueItemTypeAutoFactories,
				Quantity: 100,
			},
			{
				QueueItemCompletionEstimate: QueueItemCompletionEstimate{
					YearsToBuildOne:     12,
					YearsToBuildAll:     34,
					YearsToSkipOrCancel: Infinite,
				},
				Type:     QueueItemTypeAutoMines,
				Quantity: 100,
			},
			{
				QueueItemCompletionEstimate: QueueItemCompletionEstimate{
					YearsToBuildOne:     34,
					YearsToBuildAll:     37,
					YearsToSkipOrCancel: Infinite,
				},
				Type:     QueueItemTypeAutoDefenses,
				Quantity: 10,
			},
		}

		e := NewCompletionEstimator()

		result, _, err := e.GetProductionWithEstimates(&rules, player, *planet)
		if err != nil {
			t.Fatalf("CompletionEstimator.GetProductionWithEstimates() errored unexpectedly; err = \n%v", err)
		}

		test.CompareAsJSON(t, result, wantQueue)
	})
}
