package cs

import (
	"reflect"
	"testing"

	"github.com/rs/zerolog/log"
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
		want           []ByHandCargoTransfer
		wantErr        bool
	}{
		{
			name: "simple split",
			cargoTransfers: CargoTransfers{
				Vector{}.String(): []ByHandCargoTransfer{
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
			want: []ByHandCargoTransfer{
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
				Vector{}.String(): []ByHandCargoTransfer{
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
			want: []ByHandCargoTransfer{
				{
					SourceFleetNum: 1,
					Cargo:          Cargo{Ironium: 2},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cargoTransfers.splitByHandTransfers(tt.args.source, tt.args.dest); (err != nil) != tt.wantErr {
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
		want           []ByHandCargoTransfer
	}{
		{
			name: "simple merge",
			cargoTransfers: CargoTransfers{
				Vector{}.String(): []ByHandCargoTransfer{
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
			want: []ByHandCargoTransfer{
				{
					SourceFleetNum: 1,
					Cargo:          Cargo{Ironium: 2}, // should merge ironium in
				},
			},
		},
		{
			name: "merge with different targets",
			cargoTransfers: CargoTransfers{
				Vector{}.String(): []ByHandCargoTransfer{
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
			want: []ByHandCargoTransfer{
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
				Vector{}.String(): []ByHandCargoTransfer{
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
			want: []ByHandCargoTransfer{
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
			tt.cargoTransfers.mergeByHandTransfers(tt.args.fleet, tt.args.mergingFleets)
			got := tt.cargoTransfers.getTransfers(tt.args.fleet.Position)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CargoTransfers.mergeFleetCargoTransfers() \ngot: \n%v\nwant: \n%v", got, tt.want)
			}
		})
	}
}

func TestCargoTransferer_getCargoLoadAmount(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules))
	planet := NewPlanet().WithCargo(Cargo{Ironium: 1000, Boranium: 1000, Germanium: 1000, Colonists: 1000})

	type args struct {
		dest      CargoHolder
		cargoType CargoType
		task      WaypointTransportTask
	}
	tests := []struct {
		name               string
		fleet              *Fleet
		args               args
		wantTransferAmount int
		wantWantToTransfer int
		wantWaitAtWaypoint bool
	}{
		{
			name:               "load 1kt ironium",
			fleet:              testSmallFreighter(player),
			args:               args{dest: planet, cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionLoadAmount, Amount: 1}},
			wantWantToTransfer: 1,
			wantTransferAmount: 1,
		},
		{
			name:               "load 1mg fuel",
			fleet:              testLongRangeScout(player).withFuel(0),
			args:               args{dest: testLongRangeScout(player).withFuel(10), cargoType: Fuel, task: WaypointTransportTask{Action: TransportActionLoadAmount, Amount: 1}},
			wantWantToTransfer: 1,
			wantTransferAmount: 1,
		},
		{
			name:               "load all ironium we can fit",
			fleet:              testSmallFreighter(player),
			args:               args{dest: planet, cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionLoadAll}},
			wantWantToTransfer: 1000,
			wantTransferAmount: 120, // small freighter has 120kT cargo capacity
		},
		{
			name:               "load all fuel we can fit",
			fleet:              testSmallFreighter(player).withFuel(0),
			args:               args{dest: testSmallFreighter(player), cargoType: Fuel, task: WaypointTransportTask{Action: TransportActionLoadAll}},
			wantWantToTransfer: 130,
			wantTransferAmount: 130, // small freighter has 130mg fuel capacity
		},
		{
			name:               "load all ironium we can fit (we already loaded 20)",
			fleet:              testSmallFreighter(player).withCargo(Cargo{Boranium: 20}),
			args:               args{dest: planet, cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionLoadAll}},
			wantWantToTransfer: 1000,
			wantTransferAmount: 100,
		},
		{
			name:               "load all fuel we can fit (we already loaded 20)",
			fleet:              testSmallFreighter(player).withFuel(20),
			args:               args{dest: testSmallFreighter(player), cargoType: Fuel, task: WaypointTransportTask{Action: TransportActionLoadAll}},
			wantWantToTransfer: 130,
			wantTransferAmount: 110,
		},
		{
			name:               "load fill percent",
			fleet:              testSmallFreighter(player),
			args:               args{dest: planet, cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionFillPercent, Amount: 50}},
			wantWantToTransfer: 60,
			wantTransferAmount: 60, // 50% of 120kT capacity
		},
		{
			name:               "load fill percent fuel",
			fleet:              testSmallFreighter(player).withFuel(0),
			args:               args{dest: testSmallFreighter(player), cargoType: Fuel, task: WaypointTransportTask{Action: TransportActionFillPercent, Amount: 50}},
			wantWantToTransfer: 65,
			wantTransferAmount: 65, // 50% of 130mg capacity
		},
		{
			name:               "load fill percent but wait",
			fleet:              testSmallFreighter(player),
			args:               args{dest: NewPlanet().WithCargo(Cargo{Ironium: 50}), cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionWaitForPercent, Amount: 50}},
			wantWantToTransfer: 60,
			wantTransferAmount: 50, // load all 50, wait for the additional 10
			wantWaitAtWaypoint: true,
		},
		{
			name:               "set amount to 20kT when we have 10kT already",
			fleet:              testSmallFreighter(player).withCargo(Cargo{Ironium: 10}),
			args:               args{dest: planet, cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionSetAmountTo, Amount: 20}},
			wantWantToTransfer: 10,
			wantTransferAmount: 10, // load 10kT more
		},
		{
			name:               "set fuel amount to 20mg when we have 10mg already",
			fleet:              testSmallFreighter(player).withFuel(10),
			args:               args{dest: testSmallFreighter(player), cargoType: Fuel, task: WaypointTransportTask{Action: TransportActionSetAmountTo, Amount: 20}},
			wantWantToTransfer: 10,
			wantTransferAmount: 10, // load 10mg more
		},
		{
			name:               "set amount to 20kT when we have 10kT already, but planet only has 5k",
			fleet:              testSmallFreighter(player).withCargo(Cargo{Ironium: 10}),
			args:               args{dest: NewPlanet().WithCargo(Cargo{Ironium: 5}), cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionSetAmountTo, Amount: 20}},
			wantWantToTransfer: 10,
			wantTransferAmount: 5,    // load 5kT more
			wantWaitAtWaypoint: true, // wait for remaining 5kT we want
		},
		{
			name:               "set amount to 20mg fuel when we have 10mg already, but dest fleet only has 5mg",
			fleet:              testSmallFreighter(player).withFuel(10),
			args:               args{dest: testSmallFreighter(player).withFuel(5), cargoType: Fuel, task: WaypointTransportTask{Action: TransportActionSetAmountTo, Amount: 20}},
			wantWantToTransfer: 10,
			wantTransferAmount: 5,    // load 5mg more
			wantWaitAtWaypoint: true, // wait for remaining 5mg we want
		},
		{
			name:               "set amount to 20kT when we have 30kT already. We should do nothing as we already have > 20kT",
			fleet:              testSmallFreighter(player).withCargo(Cargo{Ironium: 30}),
			args:               args{dest: planet, cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionSetAmountTo, Amount: 20}},
			wantWantToTransfer: 0,
			wantTransferAmount: 0, // don't unload
		},
		{
			name:               "set amount to 20mg fuel when we have 30mg already. We should do nothing, we already have > 20mg",
			fleet:              testSmallFreighter(player).withFuel(30),
			args:               args{dest: testSmallFreighter(player).withFuel(0), cargoType: Fuel, task: WaypointTransportTask{Action: TransportActionSetAmountTo, Amount: 20}},
			wantWantToTransfer: 0,
			wantTransferAmount: 0, // don't unload
		},
		{
			name:               "set waypoint to 2500kT colonists",
			fleet:              testSmallFreighter(player),
			args:               args{dest: NewPlanet().WithCargo(Cargo{Colonists: 2600}), cargoType: Colonists, task: WaypointTransportTask{Action: TransportActionSetWaypointTo, Amount: 2500}},
			wantWantToTransfer: 100,
			wantTransferAmount: 100, // load 100kT
		},
		{
			name:               "set waypoint to 2500kT colonists, nothing to load",
			fleet:              testSmallFreighter(player),
			args:               args{dest: NewPlanet().WithCargo(Cargo{Colonists: 2400}), cargoType: Colonists, task: WaypointTransportTask{Action: TransportActionSetWaypointTo, Amount: 2500}},
			wantWantToTransfer: 0,
			wantTransferAmount: 0, // no load, colonists too low
		},
		{
			name:               "set waypoint to 2500kT colonists, fleet has colonists, don't unload during fleetLoad",
			fleet:              testSmallFreighter(player).withCargo(Cargo{Colonists: 200}),
			args:               args{dest: NewPlanet().WithCargo(Cargo{Colonists: 2400}), cargoType: Colonists, task: WaypointTransportTask{Action: TransportActionSetWaypointTo, Amount: 2500}},
			wantWantToTransfer: 0,
			wantTransferAmount: 0, // no load, colonists too low
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cargoTransferer := newCargoTransferer(log.Logger, &FullGame{})
			gotTransferAmount, gotWantToTransfer, gotWaitAtWaypoint := cargoTransferer.getCargoLoadAmount(tt.fleet, tt.args.dest, tt.args.cargoType, tt.args.task)
			if gotTransferAmount != tt.wantTransferAmount {
				t.Errorf("cargoTransfer.getCargoLoadAmount() gotTransferAmount = %v, want %v", gotTransferAmount, tt.wantTransferAmount)
			}
			if gotWantToTransfer != tt.wantWantToTransfer {
				t.Errorf("cargoTransfer.getCargoLoadAmount() gotWantToTransfer = %v, want %v", gotWantToTransfer, tt.wantWantToTransfer)
			}

			if gotWaitAtWaypoint != tt.wantWaitAtWaypoint {
				t.Errorf("FlecargoTransferet.getCargoLoadAmount() gotWaitAtWaypoint = %v, want %v", gotWaitAtWaypoint, tt.wantWaitAtWaypoint)
			}
		})
	}
}

func TestCargoTransferer_getCargoUnloadAmount(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules))
	planet := NewPlanet().WithCargo(Cargo{Ironium: 1000, Boranium: 1000, Germanium: 1000, Colonists: 1000})

	type args struct {
		dest      CargoHolder
		cargoType CargoType
		task      WaypointTransportTask
	}
	tests := []struct {
		name               string
		fleet              *Fleet
		args               args
		wantTransferAmount int
		wantWantToTransfer int
		wantWaitAtWaypoint bool
	}{
		{
			name:               "unload 1kt ironium",
			fleet:              testSmallFreighter(player).withCargo(Cargo{Ironium: 1}),
			args:               args{dest: planet, cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionUnloadAmount, Amount: 1}},
			wantTransferAmount: 1,
			wantWantToTransfer: 1,
		},
		{
			name:               "unload 1mg fuel",
			fleet:              testLongRangeScout(player).withFuel(10),
			args:               args{dest: testLongRangeScout(player).withFuel(0), cargoType: Fuel, task: WaypointTransportTask{Action: TransportActionUnloadAmount, Amount: 1}},
			wantTransferAmount: 1,
			wantWantToTransfer: 1,
		},
		{
			name:               "unload all ironium we have",
			fleet:              testSmallFreighter(player).withCargo(Cargo{Ironium: 120}),
			args:               args{dest: planet, cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionUnloadAll}},
			wantTransferAmount: 120,
			wantWantToTransfer: 120,
		},
		{
			name:               "unload all fuel the dest can fit",
			fleet:              testSmallFreighter(player),
			args:               args{dest: testSmallFreighter(player).withFuel(0), cargoType: Fuel, task: WaypointTransportTask{Action: TransportActionUnloadAll}},
			wantTransferAmount: 130, // small freighter has 130mg fuel capacity
			wantWantToTransfer: 130,
		},
		{
			name:               "unload all ironium we have",
			fleet:              testSmallFreighter(player).withCargo(Cargo{Ironium: 20}),
			args:               args{dest: planet, cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionUnloadAll}},
			wantTransferAmount: 20,
			wantWantToTransfer: 20,
		},
		{
			name:               "unload all fuel dest can fit (they already have 20)",
			fleet:              testSmallFreighter(player),
			args:               args{dest: testSmallFreighter(player).withFuel(20), cargoType: Fuel, task: WaypointTransportTask{Action: TransportActionUnloadAll}},
			wantTransferAmount: 110,
			wantWantToTransfer: 130,
		},
		{
			name:               "unload all fuel at planet, does nothing",
			fleet:              testSmallFreighter(player),
			args:               args{dest: planet, cargoType: Fuel, task: WaypointTransportTask{Action: TransportActionUnloadAll}},
			wantTransferAmount: 0,
			wantWantToTransfer: 130,
		},
		{
			name:               "set amount to 20kT when we have 30kT already",
			fleet:              testSmallFreighter(player).withCargo(Cargo{Ironium: 30}),
			args:               args{dest: planet, cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionSetAmountTo, Amount: 20}},
			wantWantToTransfer: 10,
			wantTransferAmount: 10, // unload 10kT onto planet
		},
		{
			name:               "set fuel amount to 20mg when we have 30mg already",
			fleet:              testSmallFreighter(player).withFuel(30),
			args:               args{dest: testSmallFreighter(player).withFuel(0), cargoType: Fuel, task: WaypointTransportTask{Action: TransportActionSetAmountTo, Amount: 20}},
			wantWantToTransfer: 10,
			wantTransferAmount: 10, // load 10mg more
		},
		{
			name:               "set amount to 20kT when we have 10kT, should unload nothing",
			fleet:              testSmallFreighter(player).withCargo(Cargo{Ironium: 10}),
			args:               args{dest: planet, cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionSetAmountTo, Amount: 20}},
			wantWantToTransfer: 0,
			wantTransferAmount: 0,
		},
		{
			name:               "set waypoint to 2500kT colonists, should unload 100kT",
			fleet:              testSmallFreighter(player).withCargo(Cargo{Colonists: 100}),
			args:               args{dest: NewPlanet().WithCargo(Cargo{Colonists: 2400}), cargoType: Colonists, task: WaypointTransportTask{Action: TransportActionSetWaypointTo, Amount: 2500}},
			wantWantToTransfer: 100,
			wantTransferAmount: 100, // unload 100kT
		},
		{
			name:               "set waypoint to 2500kT colonists, nothing to unload",
			fleet:              testSmallFreighter(player),
			args:               args{dest: NewPlanet().WithCargo(Cargo{Colonists: 2400}), cargoType: Colonists, task: WaypointTransportTask{Action: TransportActionSetWaypointTo, Amount: 2500}},
			wantWantToTransfer: 100,
			wantTransferAmount: 0, // no unload
		},
		{
			name:               "set waypoint to 2500kT colonists, fleet has colonists, planet has enough, don't unload",
			fleet:              testSmallFreighter(player).withCargo(Cargo{Colonists: 200}),
			args:               args{dest: NewPlanet().WithCargo(Cargo{Colonists: 2500}), cargoType: Colonists, task: WaypointTransportTask{Action: TransportActionSetWaypointTo, Amount: 2500}},
			wantWantToTransfer: 0,
			wantTransferAmount: 0, // no unload, planet has enough
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cargoTransferer := newCargoTransferer(log.Logger, &FullGame{})
			gotTransferAmount, gotWantToTransfer, gotWaitAtWaypoint := cargoTransferer.getCargoUnloadAmount(tt.fleet, tt.args.dest, tt.args.cargoType, tt.args.task)
			if gotTransferAmount != tt.wantTransferAmount {
				t.Errorf("cargoTransfer.getCargoUnloadAmount() gotTransferAmount = %v, want %v", gotTransferAmount, tt.wantTransferAmount)
			}
			if gotWantToTransfer != tt.wantWantToTransfer {
				t.Errorf("cargoTransfer.getCargoUnloadAmount() gotWantToTransfer = %v, want %v", gotWantToTransfer, tt.wantWantToTransfer)
			}
			if gotWaitAtWaypoint != tt.wantWaitAtWaypoint {
				t.Errorf("cargoTransfer.getCargoUnloadAmount() gotWaitAtWaypoint = %v, want %v", gotWaitAtWaypoint, tt.wantWaitAtWaypoint)
			}
		})
	}
}

func TestCargoTransferer_transferToDest(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules))

	type args struct {
		dest           CargoHolder
		cargoType      CargoType
		transferAmount int
	}
	tests := []struct {
		name           string
		fleet          *Fleet
		args           args
		wantFleetCargo Cargo
		wantDestCargo  Cargo
		wantInvalid    CargoTransferStatus
	}{
		{
			name:           "transfer 10kT to planet",
			fleet:          testSmallFreighter(player).withCargo(Cargo{Ironium: 10}),
			args:           args{dest: NewPlanet().WithCargo(Cargo{Ironium: 1000}), cargoType: Ironium, transferAmount: 10},
			wantFleetCargo: Cargo{},
			wantDestCargo:  Cargo{Ironium: 1010},
		},
		{
			name:           "transfer 10kT to another fleet",
			fleet:          testSmallFreighter(player).withCargo(Cargo{Ironium: 120}),
			args:           args{dest: testSmallFreighter(player).withCargo(Cargo{Ironium: 100}), cargoType: Ironium, transferAmount: 10},
			wantFleetCargo: Cargo{Ironium: 110},
			wantDestCargo:  Cargo{Ironium: 110},
		},

		{
			name:           "transfer 10kT from planet",
			fleet:          testSmallFreighter(player),
			args:           args{dest: NewPlanet().WithCargo(Cargo{Ironium: 1000}), cargoType: Ironium, transferAmount: -10},
			wantFleetCargo: Cargo{Ironium: 10},
			wantDestCargo:  Cargo{Ironium: 990},
		},
		{
			name:          "transfer 1000kT from planet, error",
			fleet:         testSmallFreighter(player),
			args:          args{dest: NewPlanet().WithCargo(Cargo{Ironium: 1000}), cargoType: Ironium, transferAmount: -1000},
			wantDestCargo: Cargo{Ironium: 1000},
			wantInvalid:   CargoTransferStatusCargoCapacity,
		},
		{
			name:           "transfer 1000kT to planet, error",
			fleet:          testSmallFreighter(player).withCargo(Cargo{Ironium: 10}),
			args:           args{dest: NewPlanet().WithCargo(Cargo{Ironium: 1000}), cargoType: Ironium, transferAmount: 1000},
			wantFleetCargo: Cargo{Ironium: 10},
			wantDestCargo:  Cargo{Ironium: 1000},
			wantInvalid:    CargoTransferStatusCargo,
		},
		{
			name:           "transfer 120kT to another fleet with cargo, error",
			fleet:          testSmallFreighter(player).withCargo(Cargo{Ironium: 120}),
			args:           args{dest: testSmallFreighter(player).withCargo(Cargo{Ironium: 100}), cargoType: Ironium, transferAmount: 120},
			wantFleetCargo: Cargo{Ironium: 120},
			wantDestCargo:  Cargo{Ironium: 100},
			wantInvalid:    CargoTransferStatusDestCargoCapacity,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cargoTransferer := newCargoTransferer(log.Logger, &FullGame{})
			invalid := cargoTransferer.transferToDest(tt.fleet, tt.args.dest, tt.args.cargoType, tt.args.transferAmount)
			if invalid != tt.wantInvalid {
				t.Errorf("cargoTransferer.transferToDest() got %v, want %v", invalid, tt.wantInvalid)
			}

			if tt.args.dest.GetCargo() != tt.wantDestCargo {
				t.Errorf("cargoTransferer.transferToDest() gave destination cargo \n%v, wanted \n%v", tt.args.dest.GetCargo(), tt.wantDestCargo)
			}

			if tt.fleet.Cargo != tt.wantFleetCargo {
				t.Errorf("cargoTransferer.transferToDest() fleet.Cargo = %v, wantFleetCargo %v", tt.fleet.Cargo, tt.wantFleetCargo)
			}
		})
	}
}

func Test_cargoTransferer_loadByHands(t *testing.T) {
	player := NewPlayer(0, NewRace().WithSpec(&rules)).WithNum(1)

	type fields struct {
		fleets  []*Fleet
		targets []CargoHolder
	}
	tests := []struct {
		name            string
		fields          fields
		transfers       []ByHandCargoTransfer
		want            []cargoTransferResult
		wantSourceCargo []Cargo
		wantTargetCargo []Cargo
	}{
		{
			name: "load 10kT ironium from a planet",
			fields: fields{
				fleets:  []*Fleet{testSmallFreighter(player).withNum(1).withCargo(Cargo{Ironium: 10})},
				targets: []CargoHolder{NewPlanet().WithNum(1).WithCargo(Cargo{Ironium: 20})},
			},
			transfers: []ByHandCargoTransfer{
				{
					MapObjectTarget: MapObjectTarget{TargetType: MapObjectTypePlanet, TargetNum: 1},
					SourceFleetNum:  1,
					Cargo:           Cargo{Ironium: -10},
				},
			},
			want: []cargoTransferResult{
				{
					cargoType:   Ironium,
					transferred: -10,
					wanted:      -10,
				},
			},
			wantSourceCargo: []Cargo{
				{Ironium: 10}, // the fleet cargo should stay the same, it just transfers for real this time
			},
			wantTargetCargo: []Cargo{
				{Ironium: 10}, // should end up with 10 ironium left on the planet
			},
		},
		{
			name: "load 20kT ironium from a planet after another fleet unloads 10kT",
			fields: fields{
				fleets: []*Fleet{
					testSmallFreighter(player).withNum(1).withCargo(Cargo{Ironium: 20}),
					testSmallFreighter(player).withNum(2).withCargo(Cargo{Ironium: 0}),
				},
				targets: []CargoHolder{NewPlanet().WithNum(1).WithCargo(Cargo{Ironium: 20})},
			},
			transfers: []ByHandCargoTransfer{
				{
					MapObjectTarget: MapObjectTarget{TargetType: MapObjectTypePlanet, TargetNum: 1},
					SourceFleetNum:  2,
					Cargo:           Cargo{Ironium: 10}, // some other fleet dumps 10kT onto the planet by hand
				},
				{
					MapObjectTarget: MapObjectTarget{TargetType: MapObjectTypePlanet, TargetNum: 1},
					SourceFleetNum:  1,
					Cargo:           Cargo{Ironium: -20}, // our fleet loads 20kT (10kT from the other fleet's dump, 10kT from the planet)
				},
			},
			want: []cargoTransferResult{
				{
					cargoType:   Ironium,
					transferred: -10,
					wanted:      -10,
				},
			},
			wantSourceCargo: []Cargo{
				{Ironium: 20},
				{Ironium: 0},
			},
			wantTargetCargo: []Cargo{
				{Ironium: 10}, // should end up with 10 ironium left on the planet
			},
		},
		{
			name: "load 10kT ironium from a planet but someone else got it first",
			fields: fields{
				fleets:  []*Fleet{testSmallFreighter(player).withNum(1).withCargo(Cargo{Ironium: 10})},
				targets: []CargoHolder{NewPlanet().WithNum(1).WithCargo(Cargo{Ironium: 0})},
			},
			transfers: []ByHandCargoTransfer{
				{
					MapObjectTarget: MapObjectTarget{TargetType: MapObjectTypePlanet, TargetNum: 1},
					SourceFleetNum:  1,
					Cargo:           Cargo{Ironium: -10},
				},
			},
			want: []cargoTransferResult{
				{
					cargoType:   Ironium,
					transferred: 0,
					wanted:      -10,
				},
			},
			wantSourceCargo: []Cargo{
				{Ironium: 0}, // someone else got the cargo, we end up with nothing instead of 10
			},
			wantTargetCargo: []Cargo{
				{Ironium: 0}, // planet cargo stays the same
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := &FullGame{
				Game:      &Game{},
				Universe:  &Universe{},
				TechStore: &StaticTechStore,
				Players:   []*Player{player},
			}

			game.Fleets = tt.fields.fleets
			for _, target := range tt.fields.targets {
				switch t := target.(type) {
				case *Planet:
					game.Planets = append(game.Planets, t)
				case *Salvage:
					game.Salvages = append(game.Salvages, t)
				case *MineralPacket:
					game.MineralPackets = append(game.MineralPackets, t)
				case *Fleet:
					game.Fleets = append(game.Fleets, t)
				}
			}

			if err := game.Universe.buildMaps(game.Players); err != nil {
				t.Error(err)
				return
			}

			tr := newCargoTransferer(testLogger, game)
			got := tr.loadByHands(player, tt.transfers)
			// these are passed in as args, don't compare them
			for i := range got {
				tt.want[i].fleet = nil
				got[i].fleet = nil
				tt.want[i].dest = nil
				got[i].dest = nil
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cargoTransferer.loadByHands() = %v, want %v", got, tt.want)
			}

			for i, dest := range tt.fields.targets {
				if dest.GetCargo() != tt.wantTargetCargo[i] {
					t.Errorf("cargoTransferer.loadByHands() got dest cargo = %v, want %v", dest.GetCargo(), tt.wantTargetCargo[i])
				}
			}
			for i, fleet := range tt.fields.fleets {
				if fleet.Cargo != tt.wantSourceCargo[i] {
					t.Errorf("cargoTransferer.loadByHands() got source cargo = %v, want %v", fleet.Cargo, tt.wantSourceCargo[i])
				}
			}
		})
	}
}

func Test_cargoTransferer_unloadByHands(t *testing.T) {
	player := NewPlayer(0, NewRace().WithSpec(&rules)).WithNum(1)

	type fields struct {
		fleets  []*Fleet
		targets []CargoHolder
	}
	tests := []struct {
		name            string
		fields          fields
		transfers       []ByHandCargoTransfer
		want            []cargoTransferResult
		wantSourceCargo []Cargo
		wantTargetCargo []Cargo
	}{
		{
			name: "unload 10kT ironium to a planet",
			fields: fields{
				// the fleet has "0" ironium because it already did the transfer on the front end
				fleets:  []*Fleet{testSmallFreighter(player).withNum(1).withCargo(Cargo{Ironium: 0})},
				targets: []CargoHolder{NewPlanet().WithNum(1).WithCargo(Cargo{Ironium: 10})},
			},
			transfers: []ByHandCargoTransfer{
				{
					MapObjectTarget: MapObjectTarget{TargetType: MapObjectTypePlanet, TargetNum: 1},
					SourceFleetNum:  1,
					Cargo:           Cargo{Ironium: 10}, // unload 10kT ironium
				},
			},
			want: []cargoTransferResult{
				{
					cargoType:   Ironium,
					transferred: 10,
					wanted:      10,
				},
			},
			wantSourceCargo: []Cargo{
				{Ironium: 0}, // the fleet cargo should stay the same, it just transfers for real this time
			},
			wantTargetCargo: []Cargo{
				{Ironium: 20}, // should end up with 20 ironium total on the planet
			},
		},
		{
			name: "unload 10kT ironium, another fleet unloads 10kT germ, loads 10kT ironium",
			fields: fields{
				fleets: []*Fleet{
					testSmallFreighter(player).withNum(1).withCargo(Cargo{Ironium: 0}),
					testSmallFreighter(player).withNum(2).withCargo(Cargo{Ironium: 0}),
				},
				targets: []CargoHolder{NewPlanet().WithNum(1).WithCargo(Cargo{})},
			},
			transfers: []ByHandCargoTransfer{
				{
					MapObjectTarget: MapObjectTarget{TargetType: MapObjectTypePlanet, TargetNum: 1},
					SourceFleetNum:  2,
					Cargo:           Cargo{Ironium: 10}, // dump 10kT onto the planet by hand
				},
				{
					MapObjectTarget: MapObjectTarget{TargetType: MapObjectTypePlanet, TargetNum: 1},
					SourceFleetNum:  1,
					Cargo:           Cargo{Ironium: -10, Germanium: 10}, // load 10kT ironium, dump 10kT Germanium
				},
			},
			want: []cargoTransferResult{
				{
					cargoType:   Germanium,
					transferred: 10,
					wanted:      10,
				},
			},
			wantSourceCargo: []Cargo{
				{Ironium: 0},
				{Ironium: 0},
			},
			wantTargetCargo: []Cargo{
				{Ironium: 0, Germanium: 10}, // should end up with 0 ironium, 10 germ
			},
		},
		{
			name: "unload 50kTi ironium, another fleet loads 10kT, should end up with 40kT unloaded",
			fields: fields{
				fleets: []*Fleet{
					testSmallFreighter(player).withNum(1).withCargo(Cargo{Ironium: 0}),
					testSmallFreighter(player).withNum(2).withCargo(Cargo{Ironium: 10}),
				},
				targets: []CargoHolder{NewPlanet().WithNum(1).WithCargo(Cargo{})},
			},
			transfers: []ByHandCargoTransfer{
				{
					MapObjectTarget: MapObjectTarget{TargetType: MapObjectTypePlanet, TargetNum: 1},
					SourceFleetNum:  1,
					Cargo:           Cargo{Ironium: 50}, // dump 50kT onto the salvage
				},
				{
					MapObjectTarget: MapObjectTarget{TargetType: MapObjectTypePlanet, TargetNum: 1},
					SourceFleetNum:  2,
					Cargo:           Cargo{Ironium: -10}, // load 10kT ironium
				},
			},
			want: []cargoTransferResult{
				{
					cargoType:   Ironium,
					transferred: 40,
					wanted:      40,
				},
			},
			wantSourceCargo: []Cargo{
				{Ironium: 0},
				{Ironium: 10},
			},
			wantTargetCargo: []Cargo{
				{Ironium: 40}, // should end up with 40 ironium
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := &FullGame{
				Game:      &Game{},
				Universe:  &Universe{},
				TechStore: &StaticTechStore,
				Players:   []*Player{player},
			}

			game.Fleets = tt.fields.fleets

			for _, target := range tt.fields.targets {
				switch t := target.(type) {
				case *Planet:
					game.Planets = append(game.Planets, t)
				case *Salvage:
					game.Salvages = append(game.Salvages, t)
				case *MineralPacket:
					game.MineralPackets = append(game.MineralPackets, t)
				case *Fleet:
					game.Fleets = append(game.Fleets, t)
				}
			}

			if err := game.Universe.buildMaps(game.Players); err != nil {
				t.Error(err)
				return
			}

			tr := newCargoTransferer(testLogger, game)
			got := tr.unloadByHands(player, tt.transfers)
			// these are passed in as args, don't compare them
			for i := range got {
				tt.want[i].fleet = nil
				got[i].fleet = nil
				tt.want[i].dest = nil
				got[i].dest = nil
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cargoTransferer.unloadByHands() = %v, want %v", got, tt.want)
			}

			for i, dest := range tt.fields.targets {
				if dest.GetCargo() != tt.wantTargetCargo[i] {
					t.Errorf("cargoTransferer.unloadByHands() got dest cargo = %v, want %v", dest.GetCargo(), tt.wantTargetCargo[i])
				}
			}
			for i, fleet := range tt.fields.fleets {
				if fleet.Cargo != tt.wantSourceCargo[i] {
					t.Errorf("cargoTransferer.unloadByHands() got source cargo = %v, want %v", fleet.Cargo, tt.wantSourceCargo[i])
				}
			}
		})
	}
}
