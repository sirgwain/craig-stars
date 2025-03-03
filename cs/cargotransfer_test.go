package cs

import (
	"reflect"
	"testing"
)

func TestCargoTransfers_splitFleetCargoTransfers(t *testing.T) {
	type args struct {
		source *Fleet
		dest   *Fleet
	}
	tests := []struct {
		name           string
		cargoTransfers CargoTransfers
		args           args
		want           []ImmediateCargoTransfer
		wantErr        bool
	}{
		{
			name: "simple split",
			cargoTransfers: CargoTransfers{
				Vector{}.String(): []ImmediateCargoTransfer{
					{
						SourceFleetNum: 1,
						Cargo:          Cargo{Ironium: 2},
					},
				},
			},
			args: args{
				source: &Fleet{MapObject: MapObject{Num: 1}, Spec: FleetSpec{ShipDesignSpec: ShipDesignSpec{CargoCapacity: 1}}},
				dest:   &Fleet{MapObject: MapObject{Num: 2}, Spec: FleetSpec{ShipDesignSpec: ShipDesignSpec{CargoCapacity: 1}}},
			},
			want: []ImmediateCargoTransfer{
				{
					SourceFleetNum: 1,
					Cargo:          Cargo{Ironium: 1},
				},
				{
					SourceFleetNum: 2,
					Cargo:          Cargo{Ironium: 1},
				},
			},
		},
		{
			name: "split no cargo in dest",
			cargoTransfers: CargoTransfers{
				Vector{}.String(): []ImmediateCargoTransfer{
					{
						SourceFleetNum: 1,
						Cargo:          Cargo{Ironium: 2},
					},
				},
			},
			args: args{
				source: &Fleet{MapObject: MapObject{Num: 1}, Spec: FleetSpec{ShipDesignSpec: ShipDesignSpec{CargoCapacity: 2}}},
				dest:   &Fleet{MapObject: MapObject{Num: 2}},
			},
			want: []ImmediateCargoTransfer{
				{
					SourceFleetNum: 1,
					Cargo:          Cargo{Ironium: 2},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cargoTransfers.splitFleetCargoTransfers(tt.args.source, tt.args.dest); (err != nil) != tt.wantErr {
				t.Errorf("CargoTransfers.splitFleetCargoTransfers() error = %v, wantErr %v", err, tt.wantErr)
			}

			got := tt.cargoTransfers.getTransfers(tt.args.source.Position)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CargoTransfers.mergeFleetCargoTransfers() \ngot: \n%v\nwant: \n%v", got, tt.want)
			}
		})
	}
}

func TestCargoTransfers_mergeFleetCargoTransfers(t *testing.T) {
	type args struct {
		fleet         *Fleet
		mergingFleets []*Fleet
	}
	tests := []struct {
		name           string
		cargoTransfers CargoTransfers
		args           args
		want           []ImmediateCargoTransfer
	}{
		{
			name: "simple merge",
			cargoTransfers: CargoTransfers{
				Vector{}.String(): []ImmediateCargoTransfer{
					{
						SourceFleetNum: 1,
						Cargo:          Cargo{Ironium: 1},
					},
					{
						SourceFleetNum: 2,
						Cargo:          Cargo{Ironium: 1},
					},
				},
			},
			args: args{
				fleet:         &Fleet{MapObject: MapObject{Num: 1}},
				mergingFleets: []*Fleet{{MapObject: MapObject{Num: 2}}},
			},
			want: []ImmediateCargoTransfer{
				{
					SourceFleetNum: 1,
					Cargo:          Cargo{Ironium: 2}, // should merge ironium in
				},
			},
		},
		{
			name: "merge with different targets",
			cargoTransfers: CargoTransfers{
				Vector{}.String(): []ImmediateCargoTransfer{
					{
						SourceFleetNum:  1,
						Cargo:           Cargo{Ironium: 1},
						MapObjectTarget: MapObjectTarget{TargetType: MapObjectTypePlanet, TargetNum: 1},
					},
					{
						SourceFleetNum: 2,
						Cargo:          Cargo{Ironium: 1},
					},
					{
						SourceFleetNum:  1,
						Cargo:           Cargo{Ironium: 1},
						MapObjectTarget: MapObjectTarget{TargetType: MapObjectTypeFleet, TargetNum: 1, TargetPlayerNum: 2},
					},
				},
			},
			args: args{
				fleet:         &Fleet{MapObject: MapObject{Num: 1}},
				mergingFleets: []*Fleet{{MapObject: MapObject{Num: 2}}},
			},
			want: []ImmediateCargoTransfer{
				{
					SourceFleetNum:  1,
					Cargo:           Cargo{Ironium: 1},
					MapObjectTarget: MapObjectTarget{TargetType: MapObjectTypePlanet, TargetNum: 1},
				},
				{
					SourceFleetNum: 1,
					Cargo:          Cargo{Ironium: 1},
				},
				{
					SourceFleetNum:  1,
					Cargo:           Cargo{Ironium: 1},
					MapObjectTarget: MapObjectTarget{TargetType: MapObjectTypeFleet, TargetNum: 1, TargetPlayerNum: 2},
				},
			},
		},
		{
			name: "merge with different fleets",
			cargoTransfers: CargoTransfers{
				Vector{}.String(): []ImmediateCargoTransfer{
					{
						SourceFleetNum: 1,
						Cargo:          Cargo{Ironium: 1},
					},
					{
						SourceFleetNum: 2,
						Cargo:          Cargo{Ironium: 1},
					},
					{
						SourceFleetNum: 3,
						Cargo:          Cargo{Ironium: 1},
					},
				},
			},
			args: args{
				fleet:         &Fleet{MapObject: MapObject{Num: 1}},
				mergingFleets: []*Fleet{{MapObject: MapObject{Num: 2}}},
			},
			want: []ImmediateCargoTransfer{
				{
					SourceFleetNum: 1,
					Cargo:          Cargo{Ironium: 2},
				},
				{
					SourceFleetNum: 3,
					Cargo:          Cargo{Ironium: 1},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cargoTransfers.mergeFleetCargoTransfers(tt.args.fleet, tt.args.mergingFleets)
			got := tt.cargoTransfers.getTransfers(tt.args.fleet.Position)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CargoTransfers.mergeFleetCargoTransfers() \ngot: \n%v\nwant: \n%v", got, tt.want)
			}
		})
	}
}
