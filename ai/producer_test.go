package ai

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/cs"
)

func Test_aiPlayer_addShipToQueue(t *testing.T) {

	ai := testAIPlayer()
	design, err := ai.designShip("Scout", cs.ShipDesignPurposeScout, cs.FleetPurposeScout)
	if err != nil {
		t.Errorf("couldn't design ship %v", err)
		return
	}

	type args struct {
		queue    []cs.ProductionQueueItem
		purpose  cs.FleetPurpose
		design   *cs.ShipDesign
		quantity int
	}
	tests := []struct {
		name      string
		args      args
		wantQueue []cs.ProductionQueueItem
		wantIndex int
	}{
		{
			name:      "add to top of empty queue",
			args:      args{queue: []cs.ProductionQueueItem{}, purpose: cs.FleetPurposeScout, design: design, quantity: 1},
			wantQueue: []cs.ProductionQueueItem{{Type: cs.QueueItemTypeShipToken, Quantity: 1, DesignNum: design.Num, Tags: cs.Tags{cs.TagPurpose: string(cs.FleetPurposeScout)}}},
			wantIndex: 0,
		},
		{
			name: "add to top of queue with auto",
			args: args{queue: []cs.ProductionQueueItem{{Type: cs.QueueItemTypeAutoFactories, Quantity: 5}, {Type: cs.QueueItemTypeAutoMines, Quantity: 5}}, purpose: cs.FleetPurposeScout, design: design, quantity: 1},
			wantQueue: []cs.ProductionQueueItem{
				{Type: cs.QueueItemTypeShipToken, Quantity: 1, DesignNum: design.Num, Tags: cs.Tags{cs.TagPurpose: string(cs.FleetPurposeScout)}},
				{Type: cs.QueueItemTypeAutoFactories, Quantity: 5},
				{Type: cs.QueueItemTypeAutoMines, Quantity: 5},
			},
			wantIndex: 0,
		},
		{
			name: "add after concrete, before auto",
			args: args{queue: []cs.ProductionQueueItem{{Type: cs.QueueItemTypeFactory, Quantity: 1}, {Type: cs.QueueItemTypeAutoFactories, Quantity: 5}, {Type: cs.QueueItemTypeAutoMines, Quantity: 5}}, purpose: cs.FleetPurposeScout, design: design, quantity: 1},
			wantQueue: []cs.ProductionQueueItem{
				{Type: cs.QueueItemTypeFactory, Quantity: 1},
				{Type: cs.QueueItemTypeShipToken, Quantity: 1, DesignNum: design.Num, Tags: cs.Tags{cs.TagPurpose: string(cs.FleetPurposeScout)}},
				{Type: cs.QueueItemTypeAutoFactories, Quantity: 5},
				{Type: cs.QueueItemTypeAutoMines, Quantity: 5},
			},
			wantIndex: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, gotIndex := ai.addShipToQueue(tt.args.queue, tt.args.purpose, tt.args.design, tt.args.quantity, nil)
			if !reflect.DeepEqual(got, tt.wantQueue) {
				t.Errorf("aiPlayer.addShipToQueue() gotQueue = %v, want %v", got, tt.wantQueue)
			}
			if gotIndex != tt.wantIndex {
				t.Errorf("aiPlayer.addShipToQueue() gotIndex = %v, wantIndex %v", gotIndex, tt.wantIndex)
			}
		})
	}
}

func Test_aiPlayer_getYearsToBuildShips(t *testing.T) {
	ai := testAIPlayer()
	design, err := ai.designShip("Scout", cs.ShipDesignPurposeScout, cs.FleetPurposeScout)

	if err != nil {
		t.Errorf("couldn't design ship %v", err)
		return
	}

	type args struct {
		planet *cs.Planet
		ships  []fleetShip
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "empty queue, 1 year",
			args: args{
				planet: cs.NewPlanet().
					WithNum(1).
					WithPlayerNum(ai.Num).
					WithCargo(cs.Cargo{Ironium: 1000, Boranium: 1000, Germanium: 1000, Colonists: 350}).
					WithMines(10).
					WithFactories(10),
				ships: []fleetShip{
					{
						design:   design,
						quantity: 1,
					},
				},
			},
			want: 1,
		},
		{
			name: "queue 2 ships, 1 year",
			args: args{
				planet: cs.NewPlanet().
					WithNum(1).
					WithPlayerNum(ai.Num).
					WithCargo(cs.Cargo{Ironium: 1000, Boranium: 1000, Germanium: 1000, Colonists: 350}).
					WithMines(10).
					WithFactories(10),
				ships: []fleetShip{
					{
						design:   design,
						quantity: 1,
					},
					{
						design:   design,
						quantity: 2,
					},
				},
			},
			want: 1,
		},
		{
			name: "queue with fleets, more than 1 year",
			args: args{
				planet: cs.NewPlanet().
					WithNum(1).
					WithPlayerNum(ai.Num).
					WithCargo(cs.Cargo{Ironium: 1000, Boranium: 1000, Germanium: 1000, Colonists: 350}).
					WithMines(10).
					WithFactories(10).
					WithProductionQueue([]cs.ProductionQueueItem{
						{Type: cs.QueueItemTypeFactory, Quantity: 100},
						{Type: cs.QueueItemTypeMine, Quantity: 100},
					}),
				ships: []fleetShip{
					{
						design:   design,
						quantity: 1,
					},
				},
			},
			want: 22, // takes 22 years to build the ship
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ai.getYearsToBuildShips(tt.args.planet, fleet{ships: tt.args.ships})
			if (err != nil) != tt.wantErr {
				t.Errorf("aiPlayer.getYearsToBuild() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("aiPlayer.getYearsToBuild() = %v, want %v", got, tt.want)
			}
		})
	}
}
