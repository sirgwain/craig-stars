package cs

import (
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
		wanterr bool
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
				ScanRange:          66,
				ScanRangePen:       30,
				Initiative:         1,
				Movement:           4,
				MovementFull:       4,
				EstimatedRange:     2272,
				EstimatedRangeFull: 2272,
			}, wanterr: false,
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
				NumEngines:         1,
				Cost:               Cost{12, 8, 7, 24},
				TechLevel:          TechLevel{Weapons: 3, Propulsion: 3, Electronics: 1},
				Mass:               23,
				Armor:              20,
				FuelCapacity:       50,
				Initiative:         1,
				Movement:           4,
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
			}, wanterr: false,
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
				CargoCapacity:      210,
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
			}, wanterr: false,
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
			}, wanterr: false,
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
				HullType:       TechHullTypeFighter,
				Engine:         TransGalacticFuelScoop.Engine,
				NumEngines:     4,
				Cost:           Cost{98, 28, 76, 168},
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
			}, wanterr: false,
		},
		{name: "Battleship with multiple jammers",
			args: args{
				techLevels: TechLevel{26, 26, 26, 26, 26, 26},
				raceSpec:   NewRace().WithSpec(&rules).Spec,
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
				HullType:       TechHullTypeFighter,
				Engine:         TransGalacticFuelScoop.Engine,
				NumEngines:     4,
				Cost:           Cost{98, 28, 43, 171},
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
			}, wanterr: false,
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
				HullType:       TechHullTypeFighter,
				Engine:         TransGalacticFuelScoop.Engine,
				NumEngines:     4,
				Cost:           Cost{98, 28, 52, 156},
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
			}, wanterr: false,
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
				HullType:       TechHullTypeFighter,
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
			}, wanterr: false,
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
				HullType:       TechHullTypeFighter,
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
			}, wanterr: false,
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
			}, wanterr: false,
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
			}, wanterr: false,
		},
		{
			name: "Incorrect components on PP starbase",
			args: args{
				techLevels: TechLevel{4, 0, 0, 0, 0, 0},
				raceSpec:   pps.Spec,
				design: NewShipDesign(player, 1).
					WithHull(SpaceStation.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: "BANANA!!!!!!!!", HullSlotIndex: 10, Quantity: 8},
					}),
			},
			want:    ShipDesignSpec{}, // doesn't matter since want value ignored if error desired
			wanterr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ComputeShipDesignSpec(&rules, tt.args.techLevels, tt.args.raceSpec, tt.args.design)
			if tt.wanterr && err == nil {
				t.Errorf("ComputeShipDesignSpec() did not error when expected")
			} else if !test.CompareAsJSON(t, got, tt.want) {
				t.Errorf("ComputeShipDesignSpec() = %v, want %v", got, tt.want)
			}
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

func TestShipDesign_getWarshipPartBonus(t *testing.T) {
	type args struct {
		armorMulti  float64
		shieldMulti float64
		hc          *TechHullComponent
		qty         int
		shield      int
		armor       int
		beamBonus   float64
		jamming     float64
		computing   float64
		deflecting  float64
		starbase    bool
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "1 jammer 50 on blank ship",
			args: args{
				armorMulti: 1, shieldMulti: 1,
				hc:         &Jammer50,
				qty:        1,
				shield:     1,
				armor:      1,
				beamBonus:  1,
				jamming:    0,
				computing:  0,
				deflecting: 0,
				starbase:   false,
			},
			want: 1,
		},
		{
			name: "12 jammer 50s hitting starbase hardcap",
			args: args{
				armorMulti: 1, shieldMulti: 1,
				hc:         &Jammer50,
				qty:        12,
				shield:     1,
				armor:      1,
				beamBonus:  1,
				jamming:    0,
				computing:  0,
				deflecting: 0,
				starbase:   true,
			},
			want: 1.95,
		},
		{
			name: "2 battle super comps on 90% computed ship",
			args: args{
				armorMulti: 1, shieldMulti: 1,
				hc:         &Jammer30,
				qty:        2,
				shield:     1,
				armor:      1,
				beamBonus:  1,
				jamming:    0,
				computing:  0.9, // new computing: 1-(0.1*0.49) = 94.9% computing 
				deflecting: 0,
				starbase:   false,
			},
			want: 1.02579, // check graph for number lol
		},
	}
	for _, tt := range tests {
		player := NewPlayer(1, NewRace())
		player.Race.Spec.ArmorStrengthFactor = tt.args.armorMulti
		player.Race.Spec.ShieldStrengthFactor = tt.args.shieldMulti
		design := NewShipDesign(player, 1).WithHull("Battleship").WithSpec(&rules, player)
		design.Spec.Shields = tt.args.shield
		design.Spec.Armor = tt.args.armor
		design.Spec.TorpedoBonus = tt.args.computing
		design.Spec.TorpedoJamming = tt.args.jamming
		design.Spec.BeamBonus = tt.args.beamBonus
		design.Spec.BeamDefense = tt.args.deflecting
		design.Spec.Starbase = tt.args.starbase
		t.Run(tt.name, func(t *testing.T) {
			if got := design.getWarshipPartBonus(&rules, player, tt.args.hc, tt.args.qty); got != tt.want {
				t.Errorf("ShipDesign.getWarshipPartBonus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDesignShip(t *testing.T) {
	type args struct {
		hull         *TechHull
		techLevels   TechLevel
		race         *Race
		purpose      ShipDesignPurpose
		fleetPurpose FleetPurpose
	}
	tests := []struct {
		name    string
		args    args
		want    []ShipDesignSlot
		wanterr bool
	}{
		{
			name: "Humanoid Starter Stalwart Defender",
			args: args{
				hull:         &Destroyer,
				techLevels:   TechLevel{3, 3, 3, 3, 3, 3},
				race:         NewRace().WithPRT(JoaT).WithSpec(&rules),
				purpose:      ShipDesignPurposeStartingFighter,
				fleetPurpose: FleetPurposeScout,
			},
			want: []ShipDesignSlot{
				{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
				{HullComponent: AlphaTorpedo.Name, HullSlotIndex: 2, Quantity: 1},
				{HullComponent: XRayLaser.Name, HullSlotIndex: 3, Quantity: 1},
				{HullComponent: RhinoScanner.Name, HullSlotIndex: 4, Quantity: 1},
				{HullComponent: Crobmnium.Name, HullSlotIndex: 5, Quantity: 1},
				{HullComponent: FuelTank.Name, HullSlotIndex: 6, Quantity: 1},
				{HullComponent: BattleComputer.Name, HullSlotIndex: 7, Quantity: 1},
			},
			wanterr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(1, tt.args.race).WithTechLevels(tt.args.techLevels).WithNum(1)
			got, err := DesignShip(&rules, tt.args.hull, "Test Ship", player, 1, 1, tt.args.purpose, tt.args.fleetPurpose)
			if (err != nil) != tt.wanterr {
				if tt.wanterr {
					t.Errorf("DesignShip() failed to error when expected; returned slots %v instead", got.Slots)
				} else {
					t.Errorf("DesignShip() errored unexpectedly; returned error %v", err)
				}
			}
			if !CompareSlicesUnordered(got.Slots, tt.want, true) {
				t.Errorf("ShipDesign from DesignShip() had slots %v, want %v", got.Slots, tt.want)
			}
		})
	}
}
