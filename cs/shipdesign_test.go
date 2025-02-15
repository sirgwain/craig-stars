package cs

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/test"
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
			name: "invalid HullSlotIndex - negative",
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
			name: "invalid HullSlotIndex - out of bounds",
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
		{
			name: "invalid component - player banned",
			fields: fields{
				Name: "Santa Maria",
				Hull: ColonyShip.Name,
				Slots: []ShipDesignSlot{
					{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
					{HullComponent: OrbitalConstructionModule.Name, HullSlotIndex: 2, Quantity: 1},
				},
			},
			args: args{
				player: NewPlayer(1, NewRace().WithSpec(&rules)),
			},
			wantErr: true,
		},
		{
			name: "valid component - player is AR",
			fields: fields{
				Name: "Santa Maria",
				Hull: ColonyShip.Name,
				Slots: []ShipDesignSlot{
					{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
					{HullComponent: OrbitalConstructionModule.Name, HullSlotIndex: 2, Quantity: 1},
				},
			},
			args: args{
				player: NewPlayer(1, NewRace().WithPRT(AR).WithSpec(&rules)),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sd := &ShipDesign{
				Name:  tt.fields.Name,
				Hull:  tt.fields.Hull,
				Slots: tt.fields.Slots,
			}
			err := sd.Validate(&rules, tt.args.player)
			test.CheckUnexpectedError(t, err, tt.wantErr)
		})
	}
}

func Test_getNewJamming(t *testing.T) {
	type args struct {
		prevBonus      float64
		componentBonus float64
		multi          float64
		qty            int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"2x 20%", args{0, 0.2, 1, 2}, 0.36},
		{"4x 30%; 0.75x multi", args{0, 0.3, 0.75, 4}, 0.5699},
		{"2x 20%; prev 10%", args{0.1, 0.2, 1, 2}, 0.424},
		{"10x 30%; prev 10%", args{0.1, 0.3, 1, 10}, 0.9746},
		{"2x 20%; prev 10.8%, 0.75x multi", args{0.10875, 0.2, 0.75, 2}, 0.3396},
		{"4x 30%; prev 57%, 0.75x multi", args{0.569925, 0.3, 0.75, 4}, 0.7068},
		{"4x 20%; prev 57%, 0.75x multi", args{0.569925, 0.2, 0.75, 4}, 0.6762},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// round it to 4 decimal places to preserve my sanity
			if got := roundFloat(getNewJamming(tt.args.prevBonus, tt.args.componentBonus, tt.args.multi, tt.args.qty), 4); got != tt.want {
				t.Errorf("getNewJamming() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComputeShipDesignSpec(t *testing.T) {
	humanoids := NewRace().WithSpec(&rules)
	pps := NewRace().WithPRT(PP).WithSpec(&rules)
	player := NewPlayer(1, humanoids)
	type args struct {
		techLevels TechLevel
		raceSpec   RaceSpec
		design     *ShipDesign
	}
	tests := []struct {
		name    string
		args    args
		want    ShipDesignSpec
		wantErr bool
	}{
		{name: "Humanoid Starter Long Range Scout",
			args: args{
				techLevels: TechLevel{3, 3, 3, 3, 3, 3},
				raceSpec:   humanoids.Spec,
				design: NewShipDesign(player, 1).
					WithHull(Scout.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: RhinoScanner.Name, HullSlotIndex: 2, Quantity: 1},
						{HullComponent: FuelTank.Name, HullSlotIndex: 3, Quantity: 1},
					}),
			},
			want: ShipDesignSpec{
				HullType:           TechHullTypeScout,
				Engine:             LongHump6.Engine,
				NumEngines:         1,
				Cost:               Cost{17, 2, 7, 22},
				TechLevel:          TechLevel{Propulsion: 3, Electronics: 1},
				Mass:               25,
				Armor:              20,
				FuelCapacity:       300,
				ReduceCloaking:     1,
				BeamBonus:          1,
				Scanner:            true,
				ScanRange:          666,
				ScanRangePen:       30,
				Initiative:         1,
				Movement:           4,
				MovementFull:       4,
				EstimatedRange:     2272,
				EstimatedRangeFull: 2272,
			}, wantErr: false,
		},
		{name: "Humanoid Starter Armed Probe",
			args: args{
				techLevels: TechLevel{3, 3, 3, 3, 3, 3},
				raceSpec:   humanoids.Spec,
				design: NewShipDesign(player, 1).
					WithHull(Scout.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: XRayLaser.Name, HullSlotIndex: 2, Quantity: 1},
						{HullComponent: RhinoScanner.Name, HullSlotIndex: 3, Quantity: 1},
					}),
			},
			want: ShipDesignSpec{
				HullType:           TechHullTypeScout,
				Engine:             LongHump6.Engine,
				NumEngines:         42069,
				Cost:               Cost{12, 8, 7, 24},
				TechLevel:          TechLevel{Weapons: 3, Propulsion: 3, Electronics: 1},
				Mass:               23,
				Armor:              20,
				FuelCapacity:       50,
				Initiative:         1,
				Movement:           9,
				MovementFull:       4,
				MineSweep:          16,
				HasWeapons:         true,
				PowerRating:        13,
				ReduceCloaking:     1,
				BeamBonus:          1,
				Scanner:            true,
				ScanRange:          66,
				ScanRangePen:       30,
				EstimatedRange:     413,
				EstimatedRangeFull: 413,
				WeaponSlots: []ShipDesignSlot{
					{
						HullComponent: XRayLaser.Name,
						HullSlotIndex: 2,
						Quantity:      1,
					},
				},
			}, wantErr: false,
		},
		{name: "Humanoid Starter Teamster",
			args: args{
				techLevels: TechLevel{3, 3, 3, 3, 3, 3},
				raceSpec:   humanoids.Spec,
				design: NewShipDesign(player, 1).
					WithHull(MediumFreighter.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: FuelTank.Name, HullSlotIndex: 2, Quantity: 1},
						{HullComponent: Crobmnium.Name, HullSlotIndex: 3, Quantity: 1},
					}),
			},
			want: ShipDesignSpec{
				HullType:           TechHullTypeFreighter,
				Engine:             LongHump6.Engine,
				NumEngines:         1,
				Cost:               Cost{36, 0, 20, 63},
				TechLevel:          TechLevel{Propulsion: 3, Construction: 3},
				Mass:               128,
				Armor:              125,
				FuelCapacity:       700,
				CargoCapacity:      222,
				ReduceCloaking:     1,
				BeamBonus:          1,
				Scanner:            true,
				ScanRange:          0,
				ScanRangePen:       NoScanner,
				Initiative:         0,
				Movement:           3,
				MovementFull:       2,
				EstimatedRange:     1041,
				EstimatedRangeFull: 394,
			}, wantErr: false,
		},
		{name: "RS Shielded Destroyer",
			args: args{
				techLevels: TechLevel{3, 3, 3, 3, 3, 3},
				raceSpec:   NewRace().WithLRT(RS).WithSpec(&rules).Spec,
				design: NewShipDesign(player, 1).
					WithHull(Destroyer.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: Laser.Name, HullSlotIndex: 2, Quantity: 1},
						{HullComponent: AlphaTorpedo.Name, HullSlotIndex: 3, Quantity: 1},
						{HullComponent: MoleSkinShield.Name, HullSlotIndex: 4, Quantity: 1},
						{HullComponent: Tritanium.Name, HullSlotIndex: 5, Quantity: 1},
						{HullComponent: BattleComputer.Name, HullSlotIndex: 6, Quantity: 1},
					}),
			},
			want: ShipDesignSpec{
				HullType:       TechHullTypeFighter,
				Engine:         LongHump6.Engine,
				NumEngines:     1,
				Cost:           Cost{33, 11, 23, 67},
				TechLevel:      TechLevel{Propulsion: 3, Construction: 3},
				Mass:           127,
				Armor:          225, // 200 + 50/2 for the RS armor negative
				FuelCapacity:   280,
				Scanner:        true,
				ScanRange:      60,
				ScanRangePen:   30,
				TorpedoBonus:   .2,
				Initiative:     4,
				Movement:       3,
				MovementFull:   3,
				PowerRating:    12,
				ReduceCloaking: 1,
				BeamBonus:      1,
				Shields:        35, // 25*1.4 for RS 40% better shields
				MineSweep:      10,
				HasWeapons:     true,
				WeaponSlots: []ShipDesignSlot{
					{HullComponent: Laser.Name, HullSlotIndex: 2, Quantity: 1},
					{HullComponent: AlphaTorpedo.Name, HullSlotIndex: 3, Quantity: 1},
				},
				EstimatedRange:     419,
				EstimatedRangeFull: 419,
			}, wantErr: false,
		},
		{name: "Battleship with multiple battle computers",
			args: args{
				techLevels: TechLevel{26, 26, 26, 26, 26, 26},
				raceSpec:   NewRace().WithSpec(&rules).Spec,
				design: NewShipDesign(player, 1).
					WithHull(Battleship.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: TransGalacticFuelScoop.Name, HullSlotIndex: 1, Quantity: 4},
						{HullComponent: JihadMissile.Name, HullSlotIndex: 2, Quantity: 2},
						{HullComponent: BattleComputer.Name, HullSlotIndex: 10, Quantity: 3},
						{HullComponent: BattleSuperComputer.Name, HullSlotIndex: 11, Quantity: 3},
					}),
			},
			want: ShipDesignSpec{
				HullType:       TechHullTypeCapitalShip,
				Engine:         TransGalacticFuelScoop.Engine,
				NumEngines:     4,
				Cost:           Cost{98, 28, 76, 165},
				TechLevel:      TechLevel{Energy: 5, Weapons: 12, Propulsion: 9, Construction: 13, Electronics: 11},
				Mass:           374,
				Armor:          2000,
				FuelCapacity:   2800,
				Scanner:        true,
				ScanRange:      0,
				ScanRangePen:   NoScanner,
				TorpedoBonus:   .8244,
				Initiative:     19,
				Movement:       5,
				MovementFull:   5,
				PowerRating:    255,
				ReduceCloaking: 1,
				BeamBonus:      1,
				HasWeapons:     true,
				WeaponSlots: []ShipDesignSlot{
					{HullComponent: JihadMissile.Name, HullSlotIndex: 2, Quantity: 2},
				},
				EstimatedRange:     1497,
				EstimatedRangeFull: 1497,
			}, wantErr: false,
		},
		{name: "IS Battleship with multiple jammers",
			args: args{
				techLevels: TechLevel{26, 26, 26, 26, 26, 26},
				raceSpec:   NewRace().WithPRT(IS).WithSpec(&rules).Spec,
				design: NewShipDesign(player, 1).
					WithHull(Battleship.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: TransGalacticFuelScoop.Name, HullSlotIndex: 1, Quantity: 4},
						{HullComponent: JihadMissile.Name, HullSlotIndex: 2, Quantity: 2},
						{HullComponent: Jammer10.Name, HullSlotIndex: 10, Quantity: 3},
						{HullComponent: Jammer20.Name, HullSlotIndex: 11, Quantity: 3},
					}),
			},
			want: ShipDesignSpec{
				HullType:       TechHullTypeCapitalShip,
				Engine:         TransGalacticFuelScoop.Engine,
				NumEngines:     4,
				Cost:           Cost{109, 30, 45, 170},
				TechLevel:      TechLevel{Energy: 4, Weapons: 12, Propulsion: 9, Construction: 13, Electronics: 10},
				Mass:           374,
				Armor:          2000,
				FuelCapacity:   2800,
				Scanner:        true,
				ScanRange:      0,
				ScanRangePen:   NoScanner,
				TorpedoJamming: .6268,
				Initiative:     10,
				Movement:       5,
				MovementFull:   5,
				PowerRating:    255,
				ReduceCloaking: 1,
				BeamBonus:      1,
				HasWeapons:     true,
				WeaponSlots: []ShipDesignSlot{
					{HullComponent: JihadMissile.Name, HullSlotIndex: 2, Quantity: 2},
				},
				EstimatedRange:     1497,
				EstimatedRangeFull: 1497,
			}, wantErr: false,
		},
		{name: "Battleship with multiple deflectors",
			args: args{
				techLevels: TechLevel{26, 26, 26, 26, 26, 26},
				raceSpec:   NewRace().WithSpec(&rules).Spec,
				design: NewShipDesign(player, 1).
					WithHull(Battleship.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: TransGalacticFuelScoop.Name, HullSlotIndex: 1, Quantity: 4},
						{HullComponent: JihadMissile.Name, HullSlotIndex: 2, Quantity: 2},
						{HullComponent: BeamDeflector.Name, HullSlotIndex: 10, Quantity: 3},
						{HullComponent: BeamDeflector.Name, HullSlotIndex: 11, Quantity: 3},
					}),
			},
			want: ShipDesignSpec{
				HullType:       TechHullTypeCapitalShip,
				Engine:         TransGalacticFuelScoop.Engine,
				NumEngines:     4,
				Cost:           Cost{98, 28, 46, 156},
				TechLevel:      TechLevel{Energy: 6, Weapons: 12, Propulsion: 9, Construction: 13, Electronics: 6},
				Mass:           374,
				Armor:          2000,
				FuelCapacity:   2800,
				Scanner:        true,
				ScanRange:      0,
				ScanRangePen:   NoScanner,
				BeamDefense:    .5314,
				Initiative:     10,
				Movement:       5,
				MovementFull:   5,
				PowerRating:    255,
				ReduceCloaking: 1,
				BeamBonus:      1,
				HasWeapons:     true,
				WeaponSlots: []ShipDesignSlot{
					{HullComponent: JihadMissile.Name, HullSlotIndex: 2, Quantity: 2},
				},
				EstimatedRange:     1497,
				EstimatedRangeFull: 1497,
			}, wantErr: false,
		},
		{name: "Battleship with multiple capacitors",
			args: args{
				techLevels: TechLevel{26, 26, 26, 26, 26, 26},
				raceSpec:   NewRace().WithSpec(&rules).Spec,
				design: NewShipDesign(player, 1).
					WithHull(Battleship.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: TransGalacticFuelScoop.Name, HullSlotIndex: 1, Quantity: 4},
						{HullComponent: JihadMissile.Name, HullSlotIndex: 2, Quantity: 2},
						{HullComponent: FluxCapacitor.Name, HullSlotIndex: 10, Quantity: 1},
						{HullComponent: EnergyCapacitor.Name, HullSlotIndex: 11, Quantity: 3},
					}),
			},
			want: ShipDesignSpec{
				HullType:       TechHullTypeCapitalShip,
				Engine:         TransGalacticFuelScoop.Engine,
				NumEngines:     4,
				Cost:           Cost{98, 28, 44, 150},
				TechLevel:      TechLevel{Energy: 14, Weapons: 12, Propulsion: 9, Construction: 13, Electronics: 8},
				Mass:           372,
				Armor:          2000,
				FuelCapacity:   2800,
				Scanner:        true,
				ScanRange:      0,
				ScanRangePen:   NoScanner,
				BeamBonus:      1.5972,
				Initiative:     10,
				Movement:       5,
				MovementFull:   5,
				PowerRating:    255,
				ReduceCloaking: 1,
				HasWeapons:     true,
				WeaponSlots: []ShipDesignSlot{
					{HullComponent: JihadMissile.Name, HullSlotIndex: 2, Quantity: 2},
				},
				EstimatedRange:     1505,
				EstimatedRangeFull: 1505,
			}, wantErr: false,
		},
		{name: "Battleship with max capacitors",
			args: args{
				techLevels: TechLevel{26, 26, 26, 26, 26, 26},
				raceSpec:   NewRace().WithSpec(&rules).Spec,
				design: NewShipDesign(player, 1).
					WithHull(Battleship.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: TransGalacticFuelScoop.Name, HullSlotIndex: 1, Quantity: 4},
						{HullComponent: JihadMissile.Name, HullSlotIndex: 2, Quantity: 2},
						{HullComponent: FluxCapacitor.Name, HullSlotIndex: 10, Quantity: 3},
						{HullComponent: FluxCapacitor.Name, HullSlotIndex: 11, Quantity: 3},
					}),
			},
			want: ShipDesignSpec{
				HullType:       TechHullTypeCapitalShip,
				Engine:         TransGalacticFuelScoop.Engine,
				NumEngines:     4,
				Cost:           Cost{98, 28, 58, 162},
				TechLevel:      TechLevel{Energy: 14, Weapons: 12, Propulsion: 9, Construction: 13, Electronics: 8},
				Mass:           374,
				Armor:          2000,
				FuelCapacity:   2800,
				Scanner:        true,
				ScanRange:      0,
				ScanRangePen:   NoScanner,
				BeamBonus:      2.55,
				Initiative:     10,
				Movement:       5,
				MovementFull:   5,
				PowerRating:    255,
				ReduceCloaking: 1,
				HasWeapons:     true,
				WeaponSlots: []ShipDesignSlot{
					{HullComponent: JihadMissile.Name, HullSlotIndex: 2, Quantity: 2},
				},
				EstimatedRange:     1497,
				EstimatedRangeFull: 1497,
			}, wantErr: false,
		},
		{name: "Mini Bomber",
			args: args{
				techLevels: TechLevel{3, 3, 3, 3, 3, 3},
				raceSpec:   humanoids.Spec,
				design: NewShipDesign(player, 1).
					WithHull(MiniBomber.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: AlphaDrive8.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: LadyFingerBomb.Name, HullSlotIndex: 2, Quantity: 1},
					}),
			},
			want: ShipDesignSpec{
				HullType:           TechHullTypeBomber,
				Engine:             AlphaDrive8.Engine,
				NumEngines:         1,
				Cost:               Cost{34, 24, 11, 62},
				TechLevel:          TechLevel{Weapons: 2, Propulsion: 7, Construction: 1},
				Mass:               85,
				Armor:              50,
				FuelCapacity:       120,
				Movement:           5,
				MovementFull:       5,
				Bomber:             true,
				PowerRating:        16,
				ReduceCloaking:     1,
				BeamBonus:          1,
				Scanner:            true,
				ScanRangePen:       NoScanner,
				EstimatedRange:     245,
				EstimatedRangeFull: 245,
				Bombs: []Bomb{
					{
						Quantity:             1,
						KillRate:             .6,
						MinKillRate:          300,
						StructureDestroyRate: 2,
					},
				},
			}, wantErr: false,
		},
		{name: "PP Starbase",
			args: args{
				techLevels: TechLevel{4, 0, 0, 0, 0, 0},
				raceSpec:   pps.Spec,
				design: NewShipDesign(player, 1).
					WithHull(SpaceStation.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: MassDriver5.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: Laser.Name, HullSlotIndex: 2, Quantity: 8},
						{HullComponent: MoleSkinShield.Name, HullSlotIndex: 3, Quantity: 8},
						{HullComponent: Laser.Name, HullSlotIndex: 4, Quantity: 8},
						{HullComponent: MoleSkinShield.Name, HullSlotIndex: 6, Quantity: 8},
						{HullComponent: Laser.Name, HullSlotIndex: 8, Quantity: 8},
						{HullComponent: Laser.Name, HullSlotIndex: 10, Quantity: 8},
					}),
			},
			want: ShipDesignSpec{
				HullType:        TechHullTypeStarbase,
				Cost:            Cost{152, 196, 278, 782},
				TechLevel:       TechLevel{Energy: 4},
				Engine:          Engine{},
				Mass:            48,
				Armor:           500,
				Shields:         400,
				MineSweep:       640,
				PowerRating:     192,
				HasWeapons:      true,
				Initiative:      14,
				BasePacketSpeed: 5,
				SafePacketSpeed: 5,
				ReduceCloaking:  1,
				BeamBonus:       1,
				RepairBonus:     .15,
				Scanner:         true,
				ScanRangePen:    NoScanner,
				SpaceDock:       UnlimitedSpaceDock,
				Starbase:        true,
				MassDriver:      MassDriver5.Name,
				MaxPopulation:   1_000_000,
				WeaponSlots: []ShipDesignSlot{
					{HullComponent: Laser.Name, HullSlotIndex: 2, Quantity: 8},
					{HullComponent: Laser.Name, HullSlotIndex: 4, Quantity: 8},
					{HullComponent: Laser.Name, HullSlotIndex: 8, Quantity: 8},
					{HullComponent: Laser.Name, HullSlotIndex: 10, Quantity: 8},
				},
			}, wantErr: false,
		},
		{
			name: "Incorrect components on PP starbase",
			args: args{
				techLevels: TechLevel{4, 0, 0, 0, 0, 0},
				raceSpec:   pps.Spec,
				design: NewShipDesign(player, 1).
					WithHull(SpaceStation.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: "ERROR 412: I'M A TEAPOT", HullSlotIndex: 420, Quantity: 69},
					}),
			},
			want:    ShipDesignSpec{}, // doesn't matter since want value ignored if error desired
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ComputeShipDesignSpec(&rules, tt.args.techLevels, tt.args.raceSpec, tt.args.design)
			test.CheckUnexpectedError(t, err, tt.wantErr)
			test.CompareAsJSON(t, got, tt.want)
		})
	}
}

func TestShipDesign_SlotsEqual(t *testing.T) {
	type args struct {
		sourceSlots []ShipDesignSlot
		otherSlots  []ShipDesignSlot
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "equal",
			args: args{
				sourceSlots: []ShipDesignSlot{
					{HullComponent: Laser.Name, HullSlotIndex: 2, Quantity: 8},
					{HullComponent: MoleSkinShield.Name, HullSlotIndex: 3, Quantity: 8},
				},
				otherSlots: []ShipDesignSlot{
					{HullComponent: Laser.Name, HullSlotIndex: 2, Quantity: 8},
					{HullComponent: MoleSkinShield.Name, HullSlotIndex: 3, Quantity: 8},
				},
			},
			want: true,
		},
		{
			name: "wrong num",
			args: args{
				sourceSlots: []ShipDesignSlot{
					{HullComponent: Laser.Name, HullSlotIndex: 2, Quantity: 8},
				},
				otherSlots: []ShipDesignSlot{
					{HullComponent: Laser.Name, HullSlotIndex: 2, Quantity: 8},
					{HullComponent: MoleSkinShield.Name, HullSlotIndex: 3, Quantity: 8},
				},
			},
			want: false,
		},
		{
			name: "unequal quantity",
			args: args{
				sourceSlots: []ShipDesignSlot{
					{HullComponent: Laser.Name, HullSlotIndex: 2, Quantity: 4},
					{HullComponent: MoleSkinShield.Name, HullSlotIndex: 3, Quantity: 8},
				},
				otherSlots: []ShipDesignSlot{
					{HullComponent: Laser.Name, HullSlotIndex: 2, Quantity: 8},
					{HullComponent: MoleSkinShield.Name, HullSlotIndex: 3, Quantity: 8},
				},
			},
			want: false,
		},
		{
			name: "unequal type",
			args: args{
				sourceSlots: []ShipDesignSlot{
					{HullComponent: Laser.Name, HullSlotIndex: 2, Quantity: 8},
					{HullComponent: MoleSkinShield.Name, HullSlotIndex: 3, Quantity: 8},
				},
				otherSlots: []ShipDesignSlot{
					{HullComponent: XRayLaser.Name, HullSlotIndex: 2, Quantity: 8},
					{HullComponent: MoleSkinShield.Name, HullSlotIndex: 3, Quantity: 8},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := &ShipDesign{
				Slots: tt.args.sourceSlots,
			}
			if got := source.SlotsEqual(tt.args.otherSlots); got != tt.want {
				t.Errorf("ShipDesign.SlotsEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShipDesignSpec_getJamOrComputerBonus(t *testing.T) {
	type fields struct {
		prevBonus float64
		starbase  bool
	}
	type args struct {
		hc           *TechHullComponent
		qty          int
		fieldToCheck TechTag
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name: "wrong field being checked",
			fields: fields{
				prevBonus: 0,
				starbase:  false,
			},
			args: args{
				hc:           &Jammer20,
				qty:          1,
				fieldToCheck: TechTagTorpedoBonus,
			},
			want: 1,
		},
		{
			name: "already at hardcap",
			fields: fields{
				prevBonus: 0.75,
				starbase:  true,
			},
			args: args{
				hc:           &Jammer20,
				qty:          1,
				fieldToCheck: TechTagTorpedoJammer,
			},
			want: 1,
		},
		{
			name: "10 Jammer 50s on starbase",
			fields: fields{
				prevBonus: 0,
				starbase:  true,
			},
			args: args{
				hc:           &Jammer50,
				qty:          10,
				fieldToCheck: TechTagTorpedoJammer,
			},
			want: 1.7493, // ((1-((1-0.5^10)*0.75))+1)/1
		},
		{
			name: "99 jammer 50s hitting ship hardcap",
			fields: fields{
				prevBonus: 0,
				starbase:  false,
			},
			args: args{
				hc:           &Jammer50,
				qty:          99,
				fieldToCheck: TechTagTorpedoJammer,
			},
			want: 1.95,
		},
		{
			name: "2 battle super comps on 90% computed ship",
			fields: fields{
				prevBonus: 0.9, // new computing: 1-(0.1*0.7^2) = 95.1% computing
				starbase:  false,
			},
			args: args{
				hc:           &BattleSuperComputer,
				qty:          2,
				fieldToCheck: TechTagTorpedoBonus,
			},
			want: 1.0268, // 1.951 / 1.9
		},
		{
			name: "1 Mega poly shell on 20% jammed starbase",
			fields: fields{
				prevBonus: 0.2,
				starbase:  true,
			},
			args: args{
				hc:           &MegaPolyShell,
				qty:          1,
				fieldToCheck: TechTagTorpedoJammer,
			},
			want: 1.0917, // 1.31 / 1.2
		},
		{
			name: "3 jammer 20s on 10% jammed starbase",
			fields: fields{
				prevBonus: 0.1,
				starbase:  true,
			},
			args: args{
				hc:           &Jammer20,
				qty:          3,
				fieldToCheck: TechTagTorpedoJammer,
			},
			want: 1.2884, // 1.4172 / 1.1
		},
		{
			name: "3 beam deflectors on 19% deflected starbase",
			fields: fields{
				prevBonus: 0.81,
				starbase:  true,
			},
			args: args{
				hc:           &BeamDeflector,
				qty:          3,
				fieldToCheck: TechTagBeamDeflector,
			},
			want: 1.1845, // 1.40951 / 1.19
		},
		{
			name: "3 deflectors on heavily deflected ship",
			fields: fields{
				prevBonus: 0.28243,
				starbase:  false,
			},
			args: args{
				hc:           &BeamDeflector,
				qty:          3,
				fieldToCheck: TechTagBeamDeflector,
			},
			want: 1.0446, // 1.7941 / 1.7176
		},
	}
	for _, tt := range tests {
		design := NewShipDesign(testPlayer(), 1).WithHull("Nubian").WithSpec(&rules, testPlayer())
		design.Spec.TorpedoBonus = tt.fields.prevBonus
		design.Spec.TorpedoJamming = tt.fields.prevBonus
		design.Spec.BeamDefense = tt.fields.prevBonus
		design.Spec.Starbase = tt.fields.starbase
		t.Run(tt.name, func(t *testing.T) {
			if got := roundFloat(design.Spec.getJamOrComputerBonus(&rules, tt.args.hc, tt.args.qty, tt.args.fieldToCheck), 4); got != tt.want {
				t.Errorf("ShipDesign.getJamOrComputerBonus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDesignShip(t *testing.T) {
	type args struct {
		hull         *TechHull
		techLevels   TechLevel
		player       *Player
		purpose      ShipDesignPurpose
		fleetPurpose FleetPurpose
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]int
		wantErr bool
	}{
		{
			name: "Humanoid Starter Stalwart Defender",
			args: args{
				techLevels:   TechLevel{3, 3, 3, 3, 3, 3},
				hull:         &Destroyer,
				player:       NewPlayer(1, NewRace().WithPRT(JoaT).WithSpec(&rules)).WithNum(1),
				purpose:      ShipDesignPurposeStartingFighter,
				fleetPurpose: FleetPurposeScout,
			},
			want: map[string]int{
				LongHump6.Name:      1,
				AlphaTorpedo.Name:   1,
				XRayLaser.Name:      1,
				RhinoScanner.Name:   1,
				Crobmnium.Name:      2,
				FuelTank.Name:       1,
				BattleComputer.Name: 1,
			},
			wantErr: false,
		},
		{
			name: "Humanoid Starter Teamster w/ IFE",
			args: args{
				techLevels:   TechLevel{3, 3, 3, 3, 3, 3},
				hull:         &MediumFreighter,
				player:       NewPlayer(1, NewRace().WithPRT(JoaT).WithLRT(IFE).WithSpec(&rules)).WithNum(1),
				purpose:      ShipDesignPurposeStartingFighter,
				fleetPurpose: FleetPurposeScout,
			},
			want: map[string]int{
				FuelMizer.Name:    1,
				RhinoScanner.Name: 1,
				Crobmnium.Name:    1,
			},
			wantErr: false,
		},
		{
			name: "IT starting Swashbuckler w/ radram",
			args: args{
				techLevels:   TechLevel{3, 3, 6, 5, 3, 3},
				hull:         &Privateer,
				player:       NewPlayer(1, NewRace().WithPRT(IT).WithLRT(CE).WithSpec(&rules)).WithNum(1),
				purpose:      ShipDesignPurposeStartingFighter,
				fleetPurpose: FleetPurposeFreighter,
			},
			want: map[string]int{
				RadiatingHydroRamScoop.Name: 1,
				RhinoScanner.Name:           1,
				AlphaTorpedo.Name:           1,
				XRayLaser.Name:              1,
				Crobmnium.Name:              2,
			},
			wantErr: false,
		},
		{
			name: "Large Freighter - avoids radram",
			args: args{
				techLevels:   TechLevel{3, 3, 6, 8, 3, 3},
				hull:         &LargeFreighter,
				player:       NewPlayer(1, NewRace().WithPRT(IT).WithSpec(&rules)).WithNum(1),
				purpose:      ShipDesignPurposeFreighter,
				fleetPurpose: FleetPurposeColonistFreighter,
			},
			want: map[string]int{
				DaddyLongLegs7.Name: 2,
				FuelTank.Name:       2,
				CowHideShield.Name:  2,
			},
			wantErr: false,
		},
		{
			name: "IFE Cargo Privateer",
			args: args{
				hull:         &Privateer,
				techLevels:   TechLevel{0, 0, 2, 4, 0, 0},
				player:       NewPlayer(1, NewRace().WithPRT(JoaT).WithLRT(IFE).WithLRT(RS).WithLRT(NRSE).WithSpec(&rules)).WithNum(1),
				purpose:      ShipDesignPurposeFreighter,
				fleetPurpose: FleetPurposeFreighter,
			},
			want: map[string]int{
				FuelMizer.Name:      1,
				CargoPod.Name:       1,
				FuelTank.Name:       2,
				MoleSkinShield.Name: 2,
			},
			wantErr: false,
		},
		{
			name: "Remote Miner",
			args: args{
				hull:         &UltraMiner,
				techLevels:   TechLevel{0, 0, 2, 15, 8, 0},
				player:       NewPlayer(1, NewRace().WithPRT(AR).WithLRT(IFE).WithLRT(ARM).WithLRT(NRSE).WithSpec(&rules)).WithNum(1),
				purpose:      ShipDesignPurposeMiner,
				fleetPurpose: FleetPurposeMiner,
			},
			want: map[string]int{
				FuelMizer.Name:      2,
				FuelTank.Name:       3,
				RoboUltraMiner.Name: 12,
			},
			wantErr: false,
		},
		{
			name: "SD Minelayer",
			args: args{
				hull:         &SuperMineLayer,
				techLevels:   TechLevel{4, 0, 2, 15, 8, 7},
				player:       NewPlayer(1, NewRace().WithPRT(SD).WithLRT(IFE).WithLRT(RS).WithLRT(NRSE).WithSpec(&rules)).WithNum(1),
				purpose:      ShipDesignPurposeDamageMineLayer,
				fleetPurpose: FleetPurposeMineLayer,
			},
			want: map[string]int{
				FuelMizer.Name:       1,
				FuelTank.Name:        3,
				MineDispenser80.Name: 19,
				CowHideShield.Name:   4,
			},
			wantErr: false,
		},
		{
			name: "Hush-A-Boom B-52 Bomber",
			args: args{
				hull:       &B52Bomber,
				techLevels: TechLevel{12, 16, 12, 15, 12, 12},
				player: NewPlayer(1, NewRace().WithPRT(WM).WithSpec(&rules)).WithNum(1).
					WithAcquiredTech(HushABoom.Name).WithAcquiredTech(LangstonShell.Name),
				purpose:      ShipDesignPurposeBomber,
				fleetPurpose: FleetPurposeBomber,
			},
			want: map[string]int{
				TransGalacticSuperScoop.Name: 2,
				FuelTank.Name:                2,
				HushABoom.Name:               16,
				LangstonShell.Name:           2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.player.TechLevels = tt.args.techLevels
			got, err := DesignShip(&rules, tt.args.hull, tt.name, tt.args.player, 1, 1, tt.args.purpose, tt.args.fleetPurpose)
			test.CheckUnexpectedError(t, err, tt.wantErr)

			tallyMap := map[string]int{}
			for _, slot := range got.Slots {
				tallyMap[slot.HullComponent] += slot.Quantity
			}
			if !reflect.DeepEqual(tallyMap, tt.want) {
				t.Errorf("ShipDesign from DesignShip() had parts \n%+v, want \n%+v", tallyMap, tt.want)
			}
		})
	}
}

func Test_designWarship(t *testing.T) {
	type fields struct {
		techLevel     TechLevel
		acquiredParts []Tech
		race          *Race
	}
	type args struct {
		hull    *TechHull
		purpose ShipDesignPurpose
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]int
		wantErr bool
	}{
		{
			name: "WM Weps 10 Battlecruiser",
			fields: fields{
				techLevel:     TechLevel{7, 10, 3, 10, 4, 0},
				acquiredParts: []Tech{},
				race:          NewRace().WithPRT(WM).WithLRT(IFE).WithLRT(RS).WithLRT(NRSE).WithSpec(&rules),
			},
			args: args{
				hull:    &BattleCruiser,
				purpose: ShipDesignPurposeBeamFighter,
			},
			want: map[string]int{
				FuelMizer.Name:              2,
				ManeuveringJet.Name:         2,
				EnergyCapacitor.Name:        2,
				ColloidalPhaser.Name:        6,
				PulsedSapper.Name:           3,
				WolverineDiffuseShield.Name: 4,
			}, wantErr: false,
		},
		{
			name: "HE Jihad Metamorph",
			fields: fields{
				techLevel:     TechLevel{10, 12, 7, 10, 0, 0},
				acquiredParts: []Tech{},
				race:          NewRace().WithPRT(HE).WithLRT(IFE).WithLRT(RS).WithLRT(NRSE).WithSpec(&rules),
			},
			args: args{
				hull:    &MetaMorph,
				purpose: ShipDesignPurposeTorpedoFighter,
			},
			want: map[string]int{
				AlphaDrive8.Name:         3,
				JihadMissile.Name:        8,
				BearNeutrinoBarrier.Name: 4,
				BattleComputer.Name:      3,
				ManeuveringJet.Name:      2,
			}, wantErr: false,
		},
		{
			name: "IS Croby Frigate Armed Scout",
			fields: fields{
				techLevel:     TechLevel{7, 10, 2, 6, 1, 1},
				acquiredParts: []Tech{},
				race:          NewRace().WithPRT(IS).WithLRT(IFE).WithLRT(RS).WithLRT(NRSE).WithSpec(&rules),
			},
			args: args{
				hull:    &Frigate,
				purpose: ShipDesignPurposeFighterScout,
			},
			want: map[string]int{
				FuelMizer.Name:       1,
				ColloidalPhaser.Name: 3,
				CrobySharmor.Name:    2,
				RhinoScanner.Name:    1,
			}, wantErr: false,
		},
		{
			name: "Max Techs HE AMP Nubian",
			fields: fields{
				techLevel:     TechLevel{26, 26, 26, 26, 26, 26},
				acquiredParts: []Tech{EnigmaPulsar.Tech},
				race:          NewRace().WithPRT(HE).WithLRT(IFE).WithLRT(RS).WithLRT(NRSE).WithSpec(&rules),
			},
			args: args{
				hull:    &Nubian,
				purpose: ShipDesignPurposeBeamFighter,
			},
			want: map[string]int{
				EnigmaPulsar.Name:         3,
				AntiMatterPulverizer.Name: 6,
				CompletePhaseShield.Name:  6,
				Jammer30.Name:             6,
				FluxCapacitor.Name:        6,
				BeamDeflector.Name:        12,
			}, wantErr: false,
		},
		{
			name: "Max Techs WM Missile Nubian",
			fields: fields{
				techLevel:     TechLevel{26, 26, 26, 26, 26, 26},
				acquiredParts: []Tech{},
				race:          NewRace().WithPRT(WM).WithLRT(IFE).WithLRT(RS).WithSpec(&rules),
			},
			args: args{
				hull:    &Nubian,
				purpose: ShipDesignPurposeTorpedoFighter,
			},
			want: map[string]int{
				GalaxyScoop.Name:         3,
				ArmageddonMissile.Name:   6,
				CompletePhaseShield.Name: 6,
				Jammer30.Name:            6,
				BattleNexus.Name:         6,
				BeamDeflector.Name:       12, // no need for jets as we already have 2 1/4 move
			}, wantErr: false,
		},
		{
			name: "W20 WM Dreadnought with Mega Poly",
			fields: fields{
				techLevel:     TechLevel{14, 20, 12, 16, 14, 14},
				acquiredParts: []Tech{MegaPolyShell.Tech},
				race:          NewRace().WithPRT(WM).WithLRT(IFE).WithLRT(NRSE).WithSpec(&rules),
			},
			args: args{
				hull:    &Dreadnought,
				purpose: ShipDesignPurposeTorpedoFighter,
			},
			want: map[string]int{
				Interspace10.Name:        5,
				DoomsdayMissile.Name:     28,
				GorillaDelagator.Name:    10,
				BattleSuperComputer.Name: 8,
				MegaPolyShell.Name:       16, // 400 armor + 100 shield blows gorilla delags out of the water
				Overthruster.Name:        2,
			}, wantErr: false,
		},
		{
			name: "ARM BB with Multi Function Pod",
			fields: fields{
				techLevel:     TechLevel{11, 24, 11, 16, 11, 7},
				acquiredParts: []Tech{MultiFunctionPod.Tech},
				race:          NewRace().WithPRT(JoaT).WithLRT(IFE).WithLRT(NRSE).WithLRT(RS).WithSpec(&rules),
			},
			args: args{
				hull:    &Battleship,
				purpose: ShipDesignPurposeTorpedoFighter,
			},
			want: map[string]int{
				Interspace10.Name:        4,
				ArmageddonMissile.Name:   20,
				BearNeutrinoBarrier.Name: 8,
				BattleSuperComputer.Name: 3,
				Jammer20.Name:            1,
				MultiFunctionPod.Name:    3,
			}, wantErr: false,
		},
		{
			name: "Weps 24 AR Death Star",
			fields: fields{
				techLevel:     TechLevel{22, 24, 22, 22, 22, 22},
				acquiredParts: []Tech{},
				race:          NewRace().WithPRT(AR).WithSpec(&rules),
			},
			args: args{
				hull:    &DeathStar,
				purpose: ShipDesignPurposeStarbase,
			},
			want: map[string]int{
				ArmageddonMissile.Name:   128,
				CompletePhaseShield.Name: 40,
				Valanium.Name:            40,
				Jammer30.Name:            8,
				BattleNexus.Name:         4,
				SuperStealthCloak.Name:   12,
				Stargate300_500.Name:     1,
				UltraDriver10.Name:       1,
			}, wantErr: false,
		},
		{
			name: "Tech 26 half Ultra Station; all MT Techs",
			fields: fields{
				techLevel:     TechLevel{26, 26, 26, 26, 26, 26},
				acquiredParts: MysteryTraderTechs,
				race:          NewRace().WithPRT(WM).WithLRT(ISB).WithSpec(&rules),
			},
			args: args{
				hull:    &UltraStation,
				purpose: ShipDesignPurposeStarbaseHalf,
			},
			want: map[string]int{
				ArmageddonMissile.Name:   48,
				CompletePhaseShield.Name: 20,
				MegaPolyShell.Name:       20,
				BattleNexus.Name:         6,
				SuperStealthCloak.Name:   6,
				Stargate300_500.Name:     1,
				UltraDriver10.Name:       1,
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(1, tt.fields.race).WithTechLevels(tt.fields.techLevel)
			for _, part := range tt.fields.acquiredParts {
				player.AcquiredTechs[part.Name] = true
			}
			got, err := designWarship(&rules, tt.args.hull, tt.name, player, 1, 2, tt.args.purpose)
			tallyMap := map[string]int{}
			for _, slot := range got.Slots {
				tallyMap[slot.HullComponent] += slot.Quantity
			}
			if err != nil {
				t.Errorf("DesignWarship() errored unexpectedly, error = %v", err)
			}
			if !reflect.DeepEqual(tallyMap, tt.want) {
				t.Errorf("ShipDesign from DesignWarship() had incorrect parts; test returned slots \n%v, expected \n%v", tallyMap, tt.want)
			}
		})
	}
}

func BenchmarkDesignShip(b *testing.B) {
	b.Run("Large", func(b *testing.B) {
		purposes := []ShipDesignPurpose{
			ShipDesignPurposeFreighter,
			ShipDesignPurposeSpeedMineLayer,
			ShipDesignPurposeMiner,
		}
		player := NewPlayer(1, NewRace().WithPRT(AR).WithLRT(IFE).WithLRT(ISB).WithLRT(RS).WithLRT(ARM).WithSpec(&rules)).WithTechLevels(TechLevel{26, 26, 26, 26, 26, 26})
		for _, tech := range MysteryTraderTechs {
			player.AcquiredTechs[tech.Name] = true
		}
		for b.Loop() {
			b.StopTimer()
			num := rules.random.Intn(3)
			purpose := purposes[num]
			var hull *TechHull
			switch num {
			case 0, 1:
				hull = &Nubian
			case 2:
				hull = &UltraMiner
			}
			fp := FleetPurposeFromShipDesignPurpose(purpose)
			b.StartTimer()
			DesignShip(&rules, hull, "Benchmark Ship", player, 1, 2, purpose, fp)
		}
	})

	b.Run("Small", func(b *testing.B) {
		purposes := []ShipDesignPurpose{
			ShipDesignPurposeFreighter,
			ShipDesignPurposeSpeedMineLayer,
			ShipDesignPurposeMiner,
		}
		player := NewPlayer(1, NewRace().WithPRT(AR).WithLRT(IFE).WithLRT(ISB).WithLRT(RS).WithLRT(ARM).WithSpec(&rules)).WithTechLevels(TechLevel{26, 26, 26, 26, 26, 26})
		for _, tech := range MysteryTraderTechs {
			player.AcquiredTechs[tech.Name] = true
		}
		for b.Loop() {
			b.StopTimer()
			num := rules.random.Intn(3)
			purpose := purposes[num]
			var hull *TechHull
			switch num {
			case 0:
				hull = &LargeFreighter
			case 1:
				hull = &Frigate
			case 2:
				hull = &MidgetMiner
			}
			fp := FleetPurposeFromShipDesignPurpose(purpose)
			b.StartTimer()
			DesignShip(&rules, hull, "Benchmark Ship", player, 1, 2, purpose, fp)
		}
	})

	b.Run("Large Warship", func(b *testing.B) {
		purposes := []ShipDesignPurpose{
			ShipDesignPurposeFighterScout,
			ShipDesignPurposeBeamFighter,
			ShipDesignPurposeTorpedoFighter,
			ShipDesignPurposeStarbase,
			ShipDesignPurposeStarbaseHalf,
			ShipDesignPurposeStarbaseQuarter,
		}
		c := rules.random.Intn(3)
		p := AR
		if c == 1 {
			p = WM
		}
		player := NewPlayer(1, NewRace().WithPRT(p).WithLRT(IFE).WithLRT(ISB).WithLRT(RS).WithLRT(ARM).WithSpec(&rules)).WithTechLevels(TechLevel{26, 26, 26, 26, 26, 26})
		for _, tech := range MysteryTraderTechs {
			player.AcquiredTechs[tech.Name] = true
		}
		b.ResetTimer()
		for b.Loop() {
			b.StopTimer()
			purpose := purposes[rules.random.Intn(6)]
			num := rules.random.Intn(8)
			var hull *TechHull
			switch num {
			case 0, 1, 2:
				if c == 1 {
					hull = &Dreadnought
				} else {
					hull = &Nubian
				}
			case 3, 4, 5:
				if c != 1 {
					hull = &DeathStar
				} else {
					hull = &UltraStation
				}
			default:
				hull = &Battleship
			}
			b.StartTimer()
			designWarship(&rules, hull, "Benchmark Ship", player, 1, 2, purpose)
		}
	})

	b.Run("Small Warship", func(b *testing.B) {
		purposes := []ShipDesignPurpose{
			ShipDesignPurposeFighterScout,
			ShipDesignPurposeBeamFighter,
			ShipDesignPurposeTorpedoFighter,
			ShipDesignPurposeStarbase,
			ShipDesignPurposeStarbaseHalf,
			ShipDesignPurposeStarbaseQuarter,
		}
		player := NewPlayer(1, NewRace().WithPRT(AR).WithLRT(IFE).WithLRT(ISB).WithLRT(RS).WithLRT(ARM).WithSpec(&rules)).WithTechLevels(TechLevel{26, 26, 26, 26, 26, 26})
		for _, tech := range MysteryTraderTechs {
			player.AcquiredTechs[tech.Name] = true
		}
		b.ResetTimer()
		for b.Loop() {
			b.StopTimer()
			purpose := purposes[rules.random.Intn(6)]
			num := rules.random.Intn(6)
			var hull *TechHull
			switch num {
			case 0:
				hull = &Frigate
			case 1, 2:
				hull = &Cruiser
			case 3, 4, 5:
				hull = &SpaceStation
			}
			b.StartTimer()
			designWarship(&rules, hull, "Benchmark Ship", player, 1, 2, purpose)
		}
	})
}
