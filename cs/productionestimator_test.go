package cs

import (
	"math"
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/test"
)

func Test_completionEstimate_getCompletionEstimate(t *testing.T) {
	type args struct {
		item                   ProductionQueueItem
		yearlyAvailableToSpend Cost
		previousItemsCost      Cost
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
				YearsToBuildOne: math.MaxInt,
				YearsToBuildAll: math.MaxInt,
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
			name: "one item, two years to build",
			args: args{
				item: ProductionQueueItem{
					Quantity:  1,
					CostOfOne: Cost{Ironium: 2},
				},
				yearlyAvailableToSpend: Cost{Ironium: 1},
			},
			want: QueueItemCompletionEstimate{
				YearsToBuildOne: 2,
				YearsToBuildAll: 2,
				PercentComplete: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &completionEstimate{}
			if got := e.getCompletionEstimate(tt.args.item, tt.args.yearlyAvailableToSpend, tt.args.previousItemsCost); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("completionEstimate.GetCompletionEstimate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_completionEstimate_PopulateCompletionEstimates(t *testing.T) {
	type args struct {
		items                  []ProductionQueueItem
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
						YearsToBuildOne: math.MaxInt,
						YearsToBuildAll: math.MaxInt,
						PercentComplete: 0,
					},
					Quantity:  1,
					CostOfOne: Cost{Ironium: 1},
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
						YearsToBuildOne: 1,
						YearsToBuildAll: 2,
						PercentComplete: 0,
					},
					CostOfOne: Cost{Boranium: 2},
					Quantity:  2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newCompletionEstimator()
			e.PopulateCompletionEstimates(tt.args.items, tt.args.yearlyAvailableToSpend)

			got := tt.args.items
			if !test.CompareAsJSON(t, got, tt.want) {
				t.Errorf("PopulateCompletionEstimates() = \n%v, want \n%v", got, tt.want)
			}

		})
	}
}
