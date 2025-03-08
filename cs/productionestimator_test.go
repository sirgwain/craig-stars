package cs

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/test"
)

func Test_completionEstimate_GetYearsToBuildOne(t *testing.T) {
	type args struct {
		item                   ProductionQueueItem
		cost                   Cost
		mineralsOnHand         Mineral
		yearlyAvailableToSpend Cost
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "one item, no resources, indefinite build time",
			args: args{
				item: ProductionQueueItem{
					Type:     QueueItemTypeFactory,
					Quantity: 1,
				},
				cost:                   Cost{Ironium: 1},
				yearlyAvailableToSpend: Cost{},
			},
			want: Infinite,
		},
		{
			name: "two items, can build one a year",
			args: args{
				item: ProductionQueueItem{
					Quantity: 2,
				},
				cost:                   Cost{Ironium: 1},
				yearlyAvailableToSpend: Cost{Ironium: 1},
			},
			want: 1,
		},
		{
			name: "two items, two years to build each, four years total",
			args: args{
				item: ProductionQueueItem{
					Quantity: 2,
				},
				cost:                   Cost{Ironium: 2},
				yearlyAvailableToSpend: Cost{Ironium: 1},
			},
			want: 2,
		},
		{
			name: "one item half done, two more years to build",
			args: args{
				item: ProductionQueueItem{
					Quantity:  1,
					Allocated: Cost{Ironium: 2},
				},
				cost:                   Cost{Ironium: 4},
				yearlyAvailableToSpend: Cost{Ironium: 1},
			},
			want: 2,
		},
		{
			name: "one item, some minerals on hand, should complete in one year",
			args: args{
				item: ProductionQueueItem{
					Quantity: 1,
				},
				cost:                   Cost{Ironium: 4},
				mineralsOnHand:         Mineral{Ironium: 3},
				yearlyAvailableToSpend: Cost{Ironium: 1},
			},
			want: 1,
		},
		{
			name: "two items, lots of minerals on hand, should complete in two years for needing resources",
			args: args{
				item: ProductionQueueItem{
					Quantity: 2,
				},
				cost:                   Cost{Ironium: 4, Resources: 1},
				mineralsOnHand:         Mineral{Ironium: 100},
				yearlyAvailableToSpend: Cost{Ironium: 1, Resources: 1},
			},
			want: 1,
		},
		{
			name: "mine, 2 resources per year",
			args: args{
				item: ProductionQueueItem{
					Type:     QueueItemTypeAutoMines,
					Quantity: 5,
				},
				cost:                   Cost{Resources: 5},
				mineralsOnHand:         Mineral{},
				yearlyAvailableToSpend: Cost{Resources: 2},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &completionEstimate{}
			if got := e.GetYearsToBuildOne(tt.args.item, tt.args.cost, tt.args.mineralsOnHand, tt.args.yearlyAvailableToSpend); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("completionEstimate.GetCompletionEstimate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_completionEstimate_GetProductionWithEstimates(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
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
				name: "Invalid Design",
				args: args{
					items: []ProductionQueueItem{
						{
							Type: QueueItemTypeShipToken,
							design: NewShipDesign(1, 1).WithSlots([]ShipDesignSlot{
								{HullComponent: "not a component", HullSlotIndex: 20, Quantity: -1},
							}),
							Quantity: 1,
						},
					},
					planet: NewPlanet().WithCargo(Cargo{0, 0, 0, 1000}).
						WithContributesOnlyLeftoverToResearch(true),
				},
				want:                  []ProductionQueueItem{},
				wantLeftoverResources: 0,
				wantErr:               true,
			},
			{
				name: "3 half done ships",
				args: args{
					items: []ProductionQueueItem{
						{
							Type:      QueueItemTypeShipToken,
							design:    testLongRangeScoutDesign(1),
							Quantity:  3,
							Allocated: Cost{8, 1, 3, 4},
							// 40 iron left to build
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
							YearsToBuildOne: 1,
							YearsToBuildAll: 4,
							YearsToSkipAuto: Infinite,
						},
						Type:     QueueItemTypeShipToken,
						design:   testLongRangeScoutDesign(1),
						Quantity: 3,
					},
				},
				wantLeftoverResources: 22,
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
							YearsToBuildOne: Infinite,
							YearsToBuildAll: Infinite,
							YearsToSkipAuto: Infinite,
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
							YearsToBuildOne: 2,
							YearsToBuildAll: 2,
							YearsToSkipAuto: Infinite,
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
							YearsToBuildOne: 1,
							YearsToBuildAll: 1,
							YearsToSkipAuto: Infinite,
						},
						Type:     QueueItemTypeMine,
						Quantity: 1,
					},
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne: 3,
							YearsToBuildAll: 5,
							YearsToSkipAuto: Infinite,
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
							Quantity: 1,
						},
						{
							Type:     QueueItemTypeMine,
							Quantity: 2,
						},
					},
					planet: NewPlanet().WithMines(1).WithCargo(Cargo{100, 100, 100, 10}),
				},
				want: []ProductionQueueItem{
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne: 8,
							YearsToBuildAll: 8,
							YearsToSkipAuto: Infinite,
						},
						Type:     QueueItemTypeFactory,
						Quantity: 1,
					},
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne: 10,
							YearsToBuildAll: 11,
							YearsToSkipAuto: Infinite,
						},
						Type:     QueueItemTypeMine,
						Quantity: 2,
					},
				},
				wantErr: false,
			},
			{
				name: "5 auto factories, then 5 auto mines",
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
							YearsToBuildOne: 1,
							YearsToBuildAll: 6,
							YearsToSkipAuto: 2,
						},
						Type:     QueueItemTypeAutoFactories,
						Quantity: 5,
					},
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne: 1, // we build some mines in the first year
							YearsToBuildAll: 9, // we finish them the next year
							YearsToSkipAuto: Infinite,
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
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							Skipped:         true, // this is skipped as it's not needed
							YearsToBuildOne: Infinite,
							YearsToBuildAll: Infinite,
							YearsToSkipAuto: 1,
						},
						Type:     QueueItemTypeAutoMinTerraform,
						Quantity: 1,
					},
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne: 1, // we easily build all auto buildable factories in one turn
							YearsToBuildAll: 1,
							YearsToSkipAuto: Infinite,
						},
						Type:     QueueItemTypeAutoFactories,
						Quantity: 100,
					},
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne: 1, // we easily build all auto buildable mines in one turn
							YearsToBuildAll: 1,
							YearsToSkipAuto: Infinite,
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
							Skipped:         true,
							YearsToBuildOne: Infinite,
							YearsToBuildAll: Infinite,
							YearsToSkipAuto: 1,
						},
						Type:     QueueItemTypeAutoFactories,
						Quantity: 100,
					},
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							Skipped:         true,
							YearsToBuildOne: Infinite,
							YearsToBuildAll: Infinite,
							YearsToSkipAuto: 1,
						},
						Type:     QueueItemTypeAutoMines,
						Quantity: 100,
					},
				},
				wantLeftoverResources: 170,
				wantErr:               false,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				e := NewCompletionEstimator()

				player := NewPlayer(1, NewRace().WithSpec(&rules)).withSpec(&rules)
				planet := tt.args.planet
				planet.Hab = Hab{50, 50, 50}                         // perfect hab
				planet.MineralConcentration = Mineral{100, 100, 100} // perfect concentration for 1kT per mine output
				planet.PlayerNum = 1
				planet.Spec = computePlanetSpec(&rules, player, planet)

				// compute specs for designs in queue
				for _, item := range tt.args.items {
					if item.design == nil {
						continue
					}
					// discard error for non-failing cases
					var err error
					if item.design.Spec, err = ComputeShipDesignSpec(&rules, player.TechLevels, player.Race.Spec, item.design); err != nil && !tt.wantErr {
						t.Fatalf("ComputeShipDesignSpec() returned error: \n%v", err)
					}
				}
				for _, item := range tt.want {
					if item.design == nil {
						continue
					}
					var err error
					if item.design.Spec, err = ComputeShipDesignSpec(&rules, player.TechLevels, player.Race.Spec, item.design); err != nil && !tt.wantErr {
						t.Fatalf("ComputeShipDesignSpec() returned error: \n%v", err)
					}
				}
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
		santaMaria := NewShipDesign(1, 1).WithHull(ColonyShip.Name).WithSlots([]ShipDesignSlot{
			{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
			{HullComponent: OrbitalConstructionModule.Name, HullSlotIndex: 2, Quantity: 1},
		})
		potatoBug := NewShipDesign(1, 2).WithHull(MidgetMiner.Name).WithSlots([]ShipDesignSlot{
			{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
			{HullComponent: RoboMidgetMiner.Name, HullSlotIndex: 2, Quantity: 2},
		})

		type args struct {
			items  []ProductionQueueItem
			planet *Planet
		}
		tests := []struct {
			name                  string
			args                  args
			want                  []ProductionQueueItem
			wantLeftoverResources int
		}{
			{
				name: "santa maria",
				args: args{
					items: []ProductionQueueItem{
						{
							Type:     QueueItemTypeShipToken,
							Quantity: 1,
							design:   santaMaria,
							// {33, 15, 31, 43}
						},
					},
					planet: NewPlanet().WithCargo(Cargo{100, 100, 100, 250}), // 32 res per year
				},
				want: []ProductionQueueItem{
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne: 2,
							YearsToBuildAll: 2,
							YearsToSkipAuto: -1,
						},
						Type:     QueueItemTypeShipToken,
						Quantity: 1,
						design:   santaMaria,
					},
				},
				wantLeftoverResources: 1,
			},
			{
				name: "potato bug",
				args: args{
					items: []ProductionQueueItem{
						{
							Type:     QueueItemTypeShipToken,
							Quantity: 2,
							design:   potatoBug,
							// {41, 0, 12, 123}
						},
					},
					planet: NewPlanet().WithCargo(Cargo{100, 100, 100, 100}),
					// Res: 32 --> 34 --> 37 --> 40 --> 42 --> 45 --> 48
				},
				want: []ProductionQueueItem{
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne: 4,
							YearsToBuildAll: 7,
							YearsToSkipAuto: -1,
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
				args: args{
					items: []ProductionQueueItem{
						{
							Type:     QueueItemTypeMineralAlchemy,
							Quantity: 1,
						},
					},
					planet: NewPlanet().WithCargo(Cargo{100, 100, 100, 250}),
				},
				want: []ProductionQueueItem{
					{
						QueueItemCompletionEstimate: QueueItemCompletionEstimate{
							YearsToBuildOne: 3,
							YearsToBuildAll: 3,
							YearsToSkipAuto: -1,
						},
						Type:     QueueItemTypeMineralAlchemy,
						Quantity: 1,
					},
				},
				wantLeftoverResources: 1,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {

				e := NewCompletionEstimator()
				player := NewPlayer(1, NewRace().WithPRT(AR).WithSpec(&rules)).
					WithTechLevels(TechLevel{Energy: 1}).
					WithNum(1).withSpec(&rules)

				planet := tt.args.planet.WithPlayerNum(1).WithContributesOnlyLeftoverToResearch(true)
				planet.Hab = Hab{50, 50, 50}                         // perfect hab
				planet.MineralConcentration = Mineral{100, 100, 100} // perfect concentration for 1kT per mine output
				planet.Spec = computePlanetSpec(&rules, player, planet)

				// compute specs for designs in queue
				for _, item := range tt.args.items {
					if item.design != nil {
						item.design = item.design.WithSpec(&rules, player)
					}
				}
				for _, item := range tt.want {
					if item.design != nil {
						item.design = item.design.WithSpec(&rules, player)
					}
				}
				planet.ProductionQueue = tt.args.items

				got, gotLeftover, err := e.GetProductionWithEstimates(&rules, player, *planet)
				if err != nil {
					t.Fatalf("CompletionEstimator.GetProductionWithEstimates() errored unexpectedly; err = \n%v", err)
				}

				if gotLeftover != tt.wantLeftoverResources {
					t.Errorf("GetProductionWithEstimates() leftover = \n%v, wantLeftover \n%v", gotLeftover, tt.wantLeftoverResources)
				}

				test.CompareAsJSON(t, got, tt.want)

			})
		}
	})
}
