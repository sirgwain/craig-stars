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
		name string
		args args
		want ShipDesignSpec
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
				Cost:               Cost{18, 2, 7, 22},
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
			},
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
			},
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
				Cost:               Cost{37, 0, 20, 63},
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
			},
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
				Cost:           Cost{34, 12, 24, 70},
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
			},
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
				Cost:           Cost{100, 32, 76, 168},
				TechLevel:      TechLevel{Energy: 5, Weapons: 12, Propulsion: 9, Construction: 13, Electronics: 11},
				Mass:           374,
				Armor:          2000,
				FuelCapacity:   2800,
				Scanner:        true,
				ScanRange:      0,
				ScanRangePen:   NoScanner,
				TorpedoBonus:   .8244,
				Initiative:     13,
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
			},
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
				Cost:           Cost{103, 32, 43, 174},
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
			},
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
				Cost:           Cost{100, 32, 52, 156},
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
			},
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
				Cost:           Cost{100, 32, 45, 153},
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
			},
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
				Cost:           Cost{100, 32, 64, 162},
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
			},
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
				Cost:               Cost{34, 25, 12, 63},
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
						StructureDestroyRate: .2,
					},
				},
			},
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
				Cost:            Cost{160, 292, 286, 894},
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
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ComputeShipDesignSpec(&rules, tt.args.techLevels, tt.args.raceSpec, tt.args.design); !test.CompareAsJSON(t, got, tt.want) {
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
			other := &ShipDesign{
				Slots: tt.args.otherSlots,
			}
			if got := source.SlotsEqual(other); got != tt.want {
				t.Errorf("ShipDesign.SlotsEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}
