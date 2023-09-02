package cs

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/test"
)

func Test_completionEstimate_GetYearsToBuildOne(t *testing.T) {
	type args struct {
		item                   ProductionQueueItem
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
					Type:      QueueItemTypeFactory,
					Quantity:  1,
					CostOfOne: Cost{Ironium: 1},
				},
				yearlyAvailableToSpend: Cost{},
			},
			want: Infinite,
		},
		{
			name: "two items, can build one a year",
			args: args{
				item: ProductionQueueItem{
					Quantity:  2,
					CostOfOne: Cost{Ironium: 1},
				},
				yearlyAvailableToSpend: Cost{Ironium: 1},
			},
			want: 1,
		},
		{
			name: "two items, two years to build each, four years total",
			args: args{
				item: ProductionQueueItem{
					Quantity:  2,
					CostOfOne: Cost{Ironium: 2},
				},
				yearlyAvailableToSpend: Cost{Ironium: 1},
			},
			want: 2,
		},
		{
			name: "one item half done, two more years to build",
			args: args{
				item: ProductionQueueItem{
					Quantity:  1,
					CostOfOne: Cost{Ironium: 4},
					Allocated: Cost{Ironium: 2},
				},
				yearlyAvailableToSpend: Cost{Ironium: 1},
			},
			want: 2,
		},
		{
			name: "one item, some minerals on hand, should complete in one year",
			args: args{
				item: ProductionQueueItem{
					Quantity:  1,
					CostOfOne: Cost{Ironium: 4},
				},
				mineralsOnHand:         Mineral{Ironium: 3},
				yearlyAvailableToSpend: Cost{Ironium: 1},
			},
			want: 1,
		},
		{
			name: "two items, lots of minerals on hand, should complete in two years for needing resources",
			args: args{
				item: ProductionQueueItem{
					Quantity:  2,
					CostOfOne: Cost{Ironium: 4, Resources: 1},
				},
				mineralsOnHand:         Mineral{Ironium: 100},
				yearlyAvailableToSpend: Cost{Ironium: 1, Resources: 1},
			},
			want: 1,
		},
		{
			name: "mine, 2 resources per year",
			args: args{
				item: ProductionQueueItem{
					Type:      QueueItemTypeAutoMines,
					Quantity:  5,
					CostOfOne: Cost{Resources: 5},
				},
				mineralsOnHand:         Mineral{},
				yearlyAvailableToSpend: Cost{Resources: 2},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &completionEstimate{}
			if got := e.GetYearsToBuildOne(tt.args.item, tt.args.mineralsOnHand, tt.args.yearlyAvailableToSpend); !reflect.DeepEqual(got, tt.want) {
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
		name string
		args args
		want []ProductionQueueItem
	}{
		{
			name: "one item, never complete",
			args: args{
				items: []ProductionQueueItem{
					{
						Type:      QueueItemTypeFactory,
						Quantity:  1,
						CostOfOne: Cost{Germanium: 4, Resources: 10},
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
						PercentComplete: 0,
					},
					Type:      QueueItemTypeFactory,
					Quantity:  1,
					CostOfOne: Cost{Germanium: 4, Resources: 10},
				},
			},
		},
		{
			name: "one item, two years to go",
			args: args{
				items: []ProductionQueueItem{
					{
						Type:      QueueItemTypeFactory,
						Quantity:  1,
						CostOfOne: Cost{Germanium: 4, Resources: 10},
						Allocated: Cost{Germanium: 2, Resources: 5},
					},
				},
				population:      3_000,                 // 3 resources per turn
				surfaceMinerals: Mineral{Germanium: 2}, // enough minerals to finish
			},
			want: []ProductionQueueItem{
				{
					QueueItemCompletionEstimate: QueueItemCompletionEstimate{
						YearsToBuildOne: 2,
						YearsToBuildAll: 2,
						PercentComplete: .5,
					},
					Type:      QueueItemTypeFactory,
					Quantity:  1,
					CostOfOne: Cost{Germanium: 4, Resources: 10},
					Allocated: Cost{Germanium: 2, Resources: 5},
				},
			},
		},
		{
			name: "two item, first completes this year, second takes 2 years",
			args: args{
				items: []ProductionQueueItem{
					{
						Type:      QueueItemTypeMine,
						Quantity:  1,
						CostOfOne: Cost{Resources: 2},
					},
					{
						Type:      QueueItemTypeFactory,
						Quantity:  2,
						CostOfOne: Cost{Germanium: 1, Resources: 2},
					},
				},
				mines:      1,    // mine 1 germ a year
				population: 2000, // 2 resources per year
			},
			want: []ProductionQueueItem{
				{
					QueueItemCompletionEstimate: QueueItemCompletionEstimate{
						YearsToBuildOne: 1,
						YearsToBuildAll: 1,
						PercentComplete: 0,
					},
					Type:      QueueItemTypeMine,
					CostOfOne: Cost{Resources: 2},
					Quantity:  1,
				},
				{
					QueueItemCompletionEstimate: QueueItemCompletionEstimate{
						YearsToBuildOne: 2,
						YearsToBuildAll: 3,
						PercentComplete: 0,
					},
					Type:      QueueItemTypeFactory,
					CostOfOne: Cost{Germanium: 1, Resources: 2},
					Quantity:  2,
				},
			},
		},
		{
			name: "two item, lots of minerals on hand, 1 resource each",
			args: args{
				items: []ProductionQueueItem{
					{
						Type:      QueueItemTypeFactory,
						Quantity:  1,
						CostOfOne: Cost{Ironium: 1, Resources: 1},
					},
					{
						Type:      QueueItemTypeMine,
						Quantity:  2,
						CostOfOne: Cost{Boranium: 2, Resources: 1},
					},
				},
				surfaceMinerals: Mineral{Ironium: 100, Boranium: 100, Germanium: 100},
				mines:           1,    // 1kT per year in cargo
				population:      1000, // 1 resource per year for pop
			},
			want: []ProductionQueueItem{
				{
					QueueItemCompletionEstimate: QueueItemCompletionEstimate{
						YearsToBuildOne: 1,
						YearsToBuildAll: 1,
						PercentComplete: 0,
					},
					Type:      QueueItemTypeFactory,
					Quantity:  1,
					CostOfOne: Cost{Ironium: 1, Resources: 1},
				},
				{
					QueueItemCompletionEstimate: QueueItemCompletionEstimate{
						YearsToBuildOne: 2,
						YearsToBuildAll: 2,
						PercentComplete: 0,
					},
					Type:      QueueItemTypeMine,
					Quantity:  2,
					CostOfOne: Cost{Boranium: 2, Resources: 1},
				},
			},
		},
		{
			name: "5 auto factories, then 5 auto mines",
			args: args{
				items: []ProductionQueueItem{
					{
						Type:      QueueItemTypeAutoFactories,
						Quantity:  5,
						CostOfOne: Cost{Germanium: 4, Resources: 10},
					},
					{
						Type:      QueueItemTypeAutoMines,
						Quantity:  10,
						CostOfOne: Cost{Resources: 5},
					},
				},
				surfaceMinerals: Mineral{Ironium: 0, Boranium: 0, Germanium: 8}, // start with enough germ for two factories
				population:      35_000,                                         // 35 resources per year
			},
			want: []ProductionQueueItem{
				{
					QueueItemCompletionEstimate: QueueItemCompletionEstimate{
						YearsToBuildOne: 1,
						YearsToBuildAll: 3,
					},
					Type:      QueueItemTypeAutoFactories,
					Quantity:  5,
					CostOfOne: Cost{Germanium: 4, Resources: 10},
				},
				{
					QueueItemCompletionEstimate: QueueItemCompletionEstimate{
						YearsToBuildOne: 1, // we build some mines in the first year
						YearsToBuildAll: 3, // we finish them the next year
					},
					Type:      QueueItemTypeAutoMines,
					Quantity:  10,
					CostOfOne: Cost{Resources: 5},
				},
			},
		},
		{
			name: "5 auto factories, then 5 auto mines, no minerals on hand, low resources",
			args: args{
				items: []ProductionQueueItem{
					{
						Type:      QueueItemTypeAutoFactories,
						Quantity:  5,
						CostOfOne: Cost{Germanium: 4, Resources: 10},
					},
					{
						Type:      QueueItemTypeAutoMines,
						Quantity:  5,
						CostOfOne: Cost{Resources: 5},
					},
				},
				surfaceMinerals: Mineral{},
				population:      2000,
			},
			want: []ProductionQueueItem{
				{
					QueueItemCompletionEstimate: QueueItemCompletionEstimate{
						YearsToBuildOne: 7,
						YearsToBuildAll: 13,
					},
					Type:      QueueItemTypeAutoFactories,
					Quantity:  5,
					CostOfOne: Cost{Germanium: 4, Resources: 10},
				},
				{
					QueueItemCompletionEstimate: QueueItemCompletionEstimate{
						YearsToBuildOne: 3,
						YearsToBuildAll: 17, // it takes a while to build all these mines
					},
					Type:      QueueItemTypeAutoMines,
					Quantity:  5,
					CostOfOne: Cost{Resources: 5},
				},
			},
		},
		{
			name: "Test later year planet with high everything",
			args: args{
				items: []ProductionQueueItem{
					{
						Type:      QueueItemTypeAutoMinTerraform,
						Quantity:  1,
						CostOfOne: Cost{Resources: 100},
					},
					{
						Type:      QueueItemTypeAutoFactories,
						Quantity:  100,
						CostOfOne: Cost{Germanium: 4, Resources: 10},
					},
					{
						Type:      QueueItemTypeAutoMines,
						Quantity:  100,
						CostOfOne: Cost{Resources: 5},
					},
				},
				surfaceMinerals: Mineral{2000, 2000, 2000},
				population:      700_000,
				mines:           634,
				factories:       664,
			},
			want: []ProductionQueueItem{
				{
					QueueItemCompletionEstimate: QueueItemCompletionEstimate{
						Skipped: true, // this is skipped as it's not needed
					},
					Type:      QueueItemTypeAutoMinTerraform,
					Quantity:  1,
					CostOfOne: Cost{Resources: 100},
				},
				{
					QueueItemCompletionEstimate: QueueItemCompletionEstimate{
						YearsToBuildOne: 1, // we easily build all auto buildable factories in one turn
						YearsToBuildAll: 1,
					},
					Type:      QueueItemTypeAutoFactories,
					Quantity:  100,
					CostOfOne: Cost{Germanium: 4, Resources: 10},
				},
				{
					QueueItemCompletionEstimate: QueueItemCompletionEstimate{
						YearsToBuildOne: 1, // we easily build all auto buildable mines in one turn
						YearsToBuildAll: 1,
					},
					Type:      QueueItemTypeAutoMines,
					Quantity:  100,
					CostOfOne: Cost{Resources: 5},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newCompletionEstimator()

			player := NewPlayer(1, NewRace().WithSpec(&rules)).withSpec(&rules)
			planet := NewPlanet().WithPlayerNum(1)
			planet.Hab = Hab{50, 50, 50}                         // perfect hab
			planet.MineralConcentration = Mineral{100, 100, 100} // perfect concentration for a 1kT per mine output
			planet.Cargo = planet.Cargo.AddMineral(tt.args.surfaceMinerals)
			planet.setPopulation(tt.args.population)
			planet.Mines = tt.args.mines
			planet.Factories = tt.args.factories
			planet.Spec = computePlanetSpec(&rules, player, planet)
			planet.ProductionQueue = tt.args.items

			if got := e.GetProductionWithEstimates(&rules, player, *planet); !test.CompareAsJSON(t, got, tt.want) {
				t.Errorf("PopulateCompletionEstimates() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}
