//go:build !wasi && !wasm

package cs

import (
	"reflect"
	"slices"
	"testing"

	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	type args struct {
		items           []ProductionQueueItem
		surfaceMinerals Mineral
		population      int
		mines           int
		factories       int
	}
	tests := []struct {
		name                  string
		args                  args
		want                  []ProductionQueueItem
		wantLeftoverResources int
	}{
		{
			name: "one item, never complete",
			args: args{
				items: []ProductionQueueItem{
					{
						Type:     QueueItemTypeFactory,
						Quantity: 1,
					},
				},
				population:      1000,
				surfaceMinerals: Mineral{},
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
		},
		{
			name: "one item, two years to go",
			args: args{
				items: []ProductionQueueItem{
					{
						Type:      QueueItemTypeFactory,
						Quantity:  1,
						Allocated: Cost{Resources: 5},
					},
				},
				population:      3_000,                 // 3 resources per turn
				surfaceMinerals: Mineral{Germanium: 4}, // enough minerals to finish
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
		},
		{
			name: "two item, first completes this year, second takes 2 years",
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
				mines:      2,    // mine 2 germ a year
				population: 5000, // 5 resources per year
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
		},
		{
			name: "two item, lots of minerals on hand, low resources",
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
				surfaceMinerals: Mineral{Ironium: 100, Boranium: 100, Germanium: 100},
				mines:           1,    // 1kT per year in cargo
				population:      1000, // 1 resource per year for pop
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
				surfaceMinerals: Mineral{Ironium: 0, Boranium: 0, Germanium: 8}, // start with enough germ for two factories
				population:      35_000,                                         // 35 resources per year
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
		},
		{
			name: "5 auto factories, then 5 auto mines, no minerals on hand, low resources",
			args: args{
				items: []ProductionQueueItem{
					{
						Type:     QueueItemTypeAutoFactories,
						Quantity: 5,
					},
					{
						Type:     QueueItemTypeAutoMines,
						Quantity: 5,
					},
				},
				surfaceMinerals: Mineral{},
				population:      2000,
			},
			want: []ProductionQueueItem{
				{
					QueueItemCompletionEstimate: QueueItemCompletionEstimate{
						YearsToBuildOne: 9, // a partial factory started by the auto factories finishes
						YearsToBuildAll: 24,
						YearsToSkipAuto: 1,
					},
					Type:     QueueItemTypeAutoFactories,
					Quantity: 5,
				},
				{
					QueueItemCompletionEstimate: QueueItemCompletionEstimate{
						YearsToBuildOne: 3,  // 2 resources a year finishes a partial mine in year 3
						YearsToBuildAll: 25, // it takes a while to build all these mines
						YearsToSkipAuto: Infinite,
					},
					Type:     QueueItemTypeAutoMines,
					Quantity: 5,
				},
			},
			wantLeftoverResources: 0,
		},
		{
			name: "5 auto factories, then 5 auto mines, mines finish first but factories still get estimates",
			args: args{
				items: []ProductionQueueItem{
					{
						Type:     QueueItemTypeAutoFactories,
						Quantity: 5,
					},
					{
						Type:     QueueItemTypeAutoMines,
						Quantity: 5,
					},
				},
				surfaceMinerals: Mineral{}, // no germanium for factories until the mines dig some up
				population:      35_000,
			},
			want: []ProductionQueueItem{
				{
					QueueItemCompletionEstimate: QueueItemCompletionEstimate{
						YearsToBuildOne: 2,
						YearsToBuildAll: 5,
						YearsToSkipAuto: 1,
					},
					Type:     QueueItemTypeAutoFactories,
					Quantity: 5,
				},
				{
					QueueItemCompletionEstimate: QueueItemCompletionEstimate{
						YearsToBuildOne: 1,
						YearsToBuildAll: 1,
						YearsToSkipAuto: Infinite,
					},
					Type:     QueueItemTypeAutoMines,
					Quantity: 5,
				},
			},
			wantLeftoverResources: 5,
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
				surfaceMinerals: Mineral{2000, 2000, 2000},
				population:      700_000,
				mines:           700,
				factories:       700,
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
				surfaceMinerals: Mineral{2000, 2000, 2000},
				population:      100_000,
				mines:           200, // we have 200, pop only supports 100
				factories:       200,
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewCompletionEstimator()

			player := NewPlayer(1, NewRace().WithSpec(&rules)).withSpec(&rules)
			planet := NewPlanet().WithPlayerNum(1)
			planet.Hab = Hab{50, 50, 50}                         // perfect hab
			planet.MineralConcentration = Mineral{100, 100, 100} // perfect concentration for a 1kT per mine output
			planet.Cargo = planet.Cargo.AddMineral(tt.args.surfaceMinerals)
			planet.setPopulation(tt.args.population)
			planet.Mines = tt.args.mines
			planet.Factories = tt.args.factories
			planet.Spec = ComputePlanetSpec(&rules, player, planet)
			planet.ProductionQueue = tt.args.items

			got, gotLeftover, err := e.GetProductionWithEstimates(&rules, player, *planet)
			if err != nil {
				t.Errorf("PopulateCompletionEstimates() returned error")
			}

			test.CompareAsJSON(t, got, tt.want)

			if gotLeftover != tt.wantLeftoverResources {
				t.Errorf("PopulateCompletionEstimates() leftover = \n%v, wantLeftover \n%v", gotLeftover, tt.wantLeftoverResources)
			}
		})
	}
}

func Test_completionEstimate_autoAlchemy(t *testing.T) {
	autoAlchemy := ProductionQueueItem{Type: QueueItemTypeAutoMineralAlchemy, Quantity: 1}
	factory := ProductionQueueItem{Type: QueueItemTypeFactory, Quantity: 1}

	tests := []struct {
		name  string
		queue []ProductionQueueItem
		want  []QueueItemCompletionEstimate
	}{
		{
			name:  "auto alchemy before a factory has no estimate, the factory builds this year",
			queue: []ProductionQueueItem{autoAlchemy, factory},
			want: []QueueItemCompletionEstimate{
				{YearsToBuildOne: Infinite, YearsToBuildAll: Infinite, YearsToSkipAuto: Infinite},
				{YearsToBuildOne: 1, YearsToBuildAll: 1, YearsToSkipAuto: Infinite},
			},
		},
		{
			name:  "auto alchemy at the end of the queue builds this year and never finishes",
			queue: []ProductionQueueItem{autoAlchemy},
			want: []QueueItemCompletionEstimate{
				{YearsToBuildOne: 1, YearsToBuildAll: Infinite, YearsToSkipAuto: Infinite},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player, planet := newProductionTestPlanet(410, Mineral{}) // enough for 4 alchemy and a factory
			planet.ProductionQueue = tt.queue

			items, leftover, err := NewCompletionEstimator().GetProductionWithEstimates(&rules, player, *planet)
			require.NoError(t, err)

			got := make([]QueueItemCompletionEstimate, len(items))
			for i, item := range items {
				got[i] = item.QueueItemCompletionEstimate
			}
			assert.Equal(t, tt.want, got)
			assert.Zero(t, leftover)
		})
	}
}

func Test_completionEstimate_doesNotChangeQueue(t *testing.T) {
	player, planet := newProductionTestPlanet(5, Mineral{}) // enough to start a factory, not finish it
	planet.ProductionQueue = []ProductionQueueItem{{Type: QueueItemTypeFactory, Quantity: 1}, {Type: QueueItemTypeMine, Quantity: 1}}
	queue := slices.Clone(planet.ProductionQueue)

	_, _, err := NewCompletionEstimator().GetProductionWithEstimates(&rules, player, *planet)
	require.NoError(t, err)
	assert.Equal(t, queue, planet.ProductionQueue)
}

func Test_completionEstimate_packetAfterMassDriverBase(t *testing.T) {
	player, planet := newProductionTestPlanet(10_000, Mineral{10_000, 10_000, 10_000})
	player.Race.PRT = PP
	player.Race.Spec = ComputeRaceSpec(&player.Race, &rules)
	player.TechLevels = TechLevel{Energy: 4}
	massDriverBase := NewShipDesign(player.Num, 1).WithHull(SpaceStation.Name).
		WithSlots([]ShipDesignSlot{{HullComponent: MassDriver5.Name, HullSlotIndex: 1, Quantity: 1}}).
		WithSpec(&rules, player)
	player.Designs = append(player.Designs, massDriverBase)
	planet.PacketTargetNum = 1
	planet.ProductionQueue = []ProductionQueueItem{
		{Type: QueueItemTypeStarbase, Quantity: 1, DesignNum: massDriverBase.Num, design: massDriverBase},
		{Type: QueueItemTypeMixedMineralPacket, Quantity: 1},
	}

	items, _, err := NewCompletionEstimator().GetProductionWithEstimates(&rules, player, *planet)
	require.NoError(t, err)
	assert.Equal(t, 1, items[0].YearsToBuildOne)
	assert.Equal(t, 1, items[1].YearsToBuildOne) // the new base launches the packet the same year
	assert.Nil(t, planet.Starbase)               // estimating doesn't change the real planet
}

// in the browser the planet's starbase only has its design, not its fleet spec
func Test_completionEstimate_starbaseWithoutSpec(t *testing.T) {
	player, planet := newProductionTestPlanet(10, Mineral{1000, 1000, 1000})
	scout := NewShipDesign(player.Num, 1).WithHull(Scout.Name).WithSpec(&rules, player)
	baseDesign := NewShipDesign(player.Num, 2).WithHull(SpaceStation.Name).WithSpec(&rules, player)
	player.Designs = append(player.Designs, scout, baseDesign)
	starbase := newStarbase(player, planet, baseDesign, "Starbase")
	starbase.Spec = ComputeFleetSpec(&rules, player, &starbase)
	planet.Starbase = &starbase
	planet.Spec = ComputePlanetSpec(&rules, player, planet)
	planet.Spec.ResourcesPerYearAvailable = 10 // one scout this year, the rest next year
	planet.Starbase = &Fleet{Tokens: []ShipToken{{Quantity: 1, DesignNum: baseDesign.Num, design: baseDesign}}}
	planet.ProductionQueue = []ProductionQueueItem{{Type: QueueItemTypeShipToken, Quantity: 3, DesignNum: scout.Num, design: scout}}

	items, _, err := NewCompletionEstimator().GetProductionWithEstimates(&rules, player, *planet)
	require.NoError(t, err)
	assert.Equal(t, 1, items[0].YearsToBuildOne)
	assert.Equal(t, 2, items[0].YearsToBuildAll)       // next year still has a dock to build them
	assert.Equal(t, FleetSpec{}, planet.Starbase.Spec) // the caller's starbase isn't changed
}
