package cs

import (
	"testing"
)

func TestShipDesign_Validate(t *testing.T) {
	type fields struct {
		Name  string
		Hull  string
		Slots []ShipDesignSlot
	}
	type args struct {
		player *Player
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid design",
			fields: fields{
				Name: "Scout",
				Hull: "Scout",
				Slots: []ShipDesignSlot{
					{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
					{HullComponent: RhinoScanner.Name, HullSlotIndex: 2, Quantity: 1},
					{HullComponent: FuelTank.Name, HullSlotIndex: 3, Quantity: 1},
				},
			},
			args: args{
				player: NewPlayer(1, NewRace().WithSpec(&rules)).WithTechLevels(TechLevel{3, 3, 3, 3, 3, 3}),
			},
			wantErr: false,
		},
		{
			name: "no name",
			fields: fields{
				Name: "",
				Hull: "Scout",
				Slots: []ShipDesignSlot{
					{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
				},
			},
			args: args{
				player: NewPlayer(1, NewRace().WithSpec(&rules)),
			},
			wantErr: true,
		},
		{
			name: "invalid hull",
			fields: fields{
				Name: "Scout",
				Hull: "some unknown hull",
			},
			args: args{
				player: NewPlayer(1, NewRace().WithSpec(&rules)),
			},
			wantErr: true,
		},
		{
			name: "invalid HullSlotIndex",
			fields: fields{
				Name: "Scout",
				Hull: "Scout",
				Slots: []ShipDesignSlot{
					{HullComponent: QuickJump5.Name, HullSlotIndex: -1, Quantity: 1},
				},
			},
			args: args{
				player: NewPlayer(1, NewRace().WithSpec(&rules)),
			},
			wantErr: true,
		},
		{
			name: "invalid HullSlotIndex2",
			fields: fields{
				Name: "Scout",
				Hull: "Scout",
				Slots: []ShipDesignSlot{
					{HullComponent: QuickJump5.Name, HullSlotIndex: 10, Quantity: 1},
				},
			},
			args: args{
				player: NewPlayer(1, NewRace().WithSpec(&rules)),
			},
			wantErr: true,
		},
		{
			name: "invalid Quantity",
			fields: fields{
				Name: "Scout",
				Hull: "Scout",
				Slots: []ShipDesignSlot{
					{HullComponent: BatScanner.Name, HullSlotIndex: 1, Quantity: 2},
				},
			},
			args: args{
				player: NewPlayer(1, NewRace().WithSpec(&rules)),
			},
			wantErr: true,
		},
		{
			name: "invalid Required",
			fields: fields{
				Name:  "Scout",
				Hull:  "Scout",
				Slots: []ShipDesignSlot{},
			},
			args: args{
				player: NewPlayer(1, NewRace().WithSpec(&rules)),
			},
			wantErr: true,
		}, {
			name: "invalid Required Quantity",
			fields: fields{
				Name: "Scout",
				Hull: "Scout",
				Slots: []ShipDesignSlot{
					{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 0},
				},
			},
			args: args{
				player: NewPlayer(1, NewRace().WithSpec(&rules)),
			},
			wantErr: true,
		},
		{
			name: "invalid component",
			fields: fields{
				Name: "Scout",
				Hull: "Scout",
				Slots: []ShipDesignSlot{
					{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
					{HullComponent: "unknown", HullSlotIndex: 2, Quantity: 1},
				},
			},
			args: args{
				player: NewPlayer(1, NewRace().WithSpec(&rules)),
			},
			wantErr: true,
		},
		{
			name: "invalid component type - cargo pod in scanner",
			fields: fields{
				Name: "Scout",
				Hull: "Scout",
				Slots: []ShipDesignSlot{
					{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
					{HullComponent: CargoPod.Name, HullSlotIndex: 2, Quantity: 1},
				},
			},
			args: args{
				player: NewPlayer(1, NewRace().WithSpec(&rules)),
			},
			wantErr: true,
		},
		{
			name: "invalid component - player can't build",
			fields: fields{
				Name: "Scout",
				Hull: "Scout",
				Slots: []ShipDesignSlot{
					{HullComponent: GalaxyScoop.Name, HullSlotIndex: 1, Quantity: 1},
				},
			},
			args: args{
				player: NewPlayer(1, NewRace().WithSpec(&rules)),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sd := &ShipDesign{
				Name:  tt.fields.Name,
				Hull:  tt.fields.Hull,
				Slots: tt.fields.Slots,
			}
			if err := sd.Validate(&rules, tt.args.player); (err != nil) != tt.wantErr {
				t.Errorf("ShipDesign.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
