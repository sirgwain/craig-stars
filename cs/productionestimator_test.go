package cs

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/test"
)

func Test_completionEstimate_GetCompletionEstimate(t *testing.T) {
	type args struct {
		item                   ProductionQueueItem
		mineralsOnHand         Cost
		yearlyAvailableToSpend Cost
	}
	tests := []struct {
		name string
		args args
		want QueueItemCompletionEstimate
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
			want: QueueItemCompletionEstimate{
				YearsToBuildOne: Infinite,
				YearsToBuildAll: Infinite,
				PercentComplete: 0,
			},
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
			want: QueueItemCompletionEstimate{
				YearsToBuildOne: 1,
				YearsToBuildAll: 2,
				PercentComplete: 0,
			},
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
			want: QueueItemCompletionEstimate{
				YearsToBuildOne: 2,
				YearsToBuildAll: 4,
				PercentComplete: 0,
			},
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
			want: QueueItemCompletionEstimate{
				YearsToBuildOne: 2,
				YearsToBuildAll: 2,
				PercentComplete: .5,
			},
		},
		{
			name: "one item, some minerals on hand, should complete in one year",
			args: args{
				item: ProductionQueueItem{
					Quantity:  1,
					CostOfOne: Cost{Ironium: 4},
				},
				mineralsOnHand:         Cost{Ironium: 3},
				yearlyAvailableToSpend: Cost{Ironium: 1},
			},
			want: QueueItemCompletionEstimate{
				YearsToBuildOne: 1,
				YearsToBuildAll: 1,
				PercentComplete: 0,
			},
		},
		{
			name: "two items, lots of minerals on hand, should complete in two years for needing resources",
			args: args{
				item: ProductionQueueItem{
					Quantity:  2,
					CostOfOne: Cost{Ironium: 4, Resources: 1},
				},
				mineralsOnHand:         Cost{Ironium: 100},
				yearlyAvailableToSpend: Cost{Ironium: 1, Resources: 1},
			},
			want: QueueItemCompletionEstimate{
				YearsToBuildOne: 1,
				YearsToBuildAll: 2,
				PercentComplete: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &completionEstimate{}
			if got := e.GetCompletionEstimate(tt.args.item, tt.args.mineralsOnHand, tt.args.yearlyAvailableToSpend); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("completionEstimate.GetCompletionEstimate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_completionEstimate_GetProductionWithEstimates(t *testing.T) {
	type args struct {
		items                  []ProductionQueueItem
		mineralsOnHand         Cost
		yearlyAvailableToSpend Cost
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
						Quantity:  1,
						CostOfOne: Cost{Ironium: 1},
					},
				},
				yearlyAvailableToSpend: Cost{},
			},
			want: []ProductionQueueItem{
				{
					QueueItemCompletionEstimate: QueueItemCompletionEstimate{
						YearsToBuildOne: Infinite,
						YearsToBuildAll: Infinite,
						PercentComplete: 0,
					},
					Quantity:  1,
					CostOfOne: Cost{Ironium: 1},
				},
			},
		},
		{
			name: "one item, halfway done",
			args: args{
				items: []ProductionQueueItem{
					{
						Quantity:  1,
						CostOfOne: Cost{Ironium: 4},
						Allocated: Cost{Ironium: 2},
					},
				},
				yearlyAvailableToSpend: Cost{Ironium: 1},
			},
			want: []ProductionQueueItem{
				{
					QueueItemCompletionEstimate: QueueItemCompletionEstimate{
						YearsToBuildOne: 2,
						YearsToBuildAll: 2,
						PercentComplete: .5,
					},
					Quantity:  1,
					CostOfOne: Cost{Ironium: 4},
					Allocated: Cost{Ironium: 2},
				},
			},
		},
		{
			name: "two item, first completes this year, second takes 2 years",
			args: args{
				items: []ProductionQueueItem{
					{
						Quantity:  1,
						CostOfOne: Cost{Ironium: 1},
					},
					{
						Quantity:  2,
						CostOfOne: Cost{Boranium: 2},
					},
				},
				yearlyAvailableToSpend: Cost{Ironium: 1, Boranium: 1},
			},
			want: []ProductionQueueItem{
				{
					QueueItemCompletionEstimate: QueueItemCompletionEstimate{
						YearsToBuildOne: 1,
						YearsToBuildAll: 1,
						PercentComplete: 0,
					},
					CostOfOne: Cost{Ironium: 1},
					Quantity:  1,
				},
				{
					QueueItemCompletionEstimate: QueueItemCompletionEstimate{
						YearsToBuildOne: 2,
						YearsToBuildAll: 4,
						PercentComplete: 0,
					},
					CostOfOne: Cost{Boranium: 2},
					Quantity:  2,
				},
			},
		},
		{
			name: "two item, lots of minerals on hand, 1 resource each",
			args: args{
				items: []ProductionQueueItem{
					{
						Quantity:  1,
						CostOfOne: Cost{Ironium: 1, Resources: 1},
					},
					{
						Quantity:  2,
						CostOfOne: Cost{Boranium: 2, Resources: 1},
					},
				},
				mineralsOnHand:         Cost{Ironium: 100, Boranium: 100, Germanium: 100},
				yearlyAvailableToSpend: Cost{Ironium: 1, Boranium: 1, Resources: 1},
			},
			want: []ProductionQueueItem{
				{
					QueueItemCompletionEstimate: QueueItemCompletionEstimate{
						YearsToBuildOne: 1,
						YearsToBuildAll: 1,
						PercentComplete: 0,
					},
					CostOfOne: Cost{Ironium: 1, Resources: 1},
					Quantity:  1,
				},
				{
					QueueItemCompletionEstimate: QueueItemCompletionEstimate{
						YearsToBuildOne: 2,
						YearsToBuildAll: 3,
						PercentComplete: 0,
					},
					CostOfOne: Cost{Boranium: 2, Resources: 1},
					Quantity:  2,
				},
			},
		},
		{
			name: "5 auto factories, then 5 auto mines",
			args: args{
				items: []ProductionQueueItem{
					{
						Type:         QueueItemTypeAutoFactories,
						Quantity:     5,
						CostOfOne:    Cost{Germanium: 4, Resources: 10},
						MaxBuildable: 100,
					},
					{
						Type:         QueueItemTypeAutoMines,
						Quantity:     5,
						CostOfOne:    Cost{Resources: 5},
						MaxBuildable: 100,
					},
				},
				mineralsOnHand:         Cost{Ironium: 0, Boranium: 0, Germanium: 8},
				yearlyAvailableToSpend: Cost{Resources: 35},
			},
			want: []ProductionQueueItem{
				{
					QueueItemCompletionEstimate: QueueItemCompletionEstimate{
						YearsToBuildOne: 1,
						YearsToBuildAll: Infinite,
						PercentComplete: 0,
					},
					Type:         QueueItemTypeAutoFactories,
					Quantity:     5,
					CostOfOne:    Cost{Germanium: 4, Resources: 10},
					MaxBuildable: 100,
				},
				{
					QueueItemCompletionEstimate: QueueItemCompletionEstimate{
						YearsToBuildOne: 1, // we build some mines in the first year
						YearsToBuildAll: 1, // TODO this should be 2...
						PercentComplete: 0,
					},
					Type:         QueueItemTypeAutoMines,
					Quantity:     5,
					CostOfOne:    Cost{Resources: 5},
					MaxBuildable: 100,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newCompletionEstimator()

			planet := NewPlanet()

			if got := e.GetProductionWithEstimates(planet, tt.args.items, tt.args.mineralsOnHand, tt.args.yearlyAvailableToSpend); !test.CompareAsJSON(t, got, tt.want) {
				t.Errorf("PopulateCompletionEstimates() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}
