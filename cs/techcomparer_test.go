//go:build !wasi && !wasm

package cs

import (
	"reflect"
	"testing"
)

func TestTechComparer_GetBestComponentWithTag(t *testing.T) {
	type fields struct {
		techLevels    TechLevel
		race          *Race
		acquiredParts []string
		beamShip      bool
	}
	type args struct {
		hullSlotType HullSlotType
		qty          int
		tag          TechTag
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *TechHullComponent
	}{
		{
			name: "Best beam with max techs",
			fields: fields{
				techLevels:    TechLevel{26, 26, 26, 26, 26, 26},
				race:          NewRace().WithPRT(JoaT),
				acquiredParts: []string{},
				beamShip:      false,
			},
			args: args{
				hullSlotType: HullSlotTypeWeapon,
				qty:          99,
				tag:          TechTagBeamWeapon,
			}, want: &AntiMatterPulverizer,
		},
		{
			name: "Best torpedo gun with max techs",
			fields: fields{
				techLevels:    TechLevel{26, 26, 26, 26, 26, 26},
				race:          NewRace().WithPRT(JoaT),
				acquiredParts: []string{},
				beamShip:      false,
			},
			args: args{
				hullSlotType: HullSlotTypeWeapon,
				qty:          99,
				tag:          TechTagTorpedo,
			}, want: &ArmageddonMissile,
		},
		{
			name: "Best NAS SS Scanner",
			fields: fields{
				techLevels:    TechLevel{15, 15, 15, 15, 15, 15},
				race:          NewRace().WithPRT(SS).WithLRT(NAS),
				acquiredParts: []string{"Mega Poly Shell", "Multi Cargo Pod"},
				beamShip:      false,
			},
			args: args{
				hullSlotType: HullSlotTypeArmorScannerElectricalMechanical,
				qty:          1,
				tag:          TechTagScanner,
			}, want: &RobberBaronScanner,
		},
		{
			name: "Best mining bot",
			fields: fields{
				techLevels:    TechLevel{14, 14, 14, 14, 14, 14},
				race:          NewRace().WithPRT(AR).WithLRT(ARM),
				acquiredParts: []string{AlienMiner.Name},
				beamShip:      false,
			},
			args: args{
				hullSlotType: HullSlotTypeMining,
				qty:          1,
				tag:          TechTagMiningRobot,
			}, want: &AlienMiner,
		},
		{
			name: "Best IT stargate at prop 6/con 10",
			fields: fields{
				techLevels:    TechLevel{0, 0, 6, 10, 0, 0},
				race:          NewRace().WithPRT(IT).WithLRT(ISB),
				acquiredParts: []string{},
				beamShip:      false,
			},
			args: args{
				hullSlotType: HullSlotTypeOrbitalElectrical,
				qty:          1,
				tag:          TechTagStargate,
			}, want: &StargateAny_300,
		},
		{
			name: "no matching part",
			fields: fields{
				techLevels:    TechLevel{0, 0, 0, 0, 0, 0},
				race:          NewRace().WithPRT(AR),
				acquiredParts: []string{},
				beamShip:      false,
			},
			args: args{
				hullSlotType: HullSlotTypeMineLayer,
				qty:          99,
				tag:          TechTagMineLayer,
			}, want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(1, tt.fields.race.WithSpec(&rules)).WithTechLevels(tt.fields.techLevels)
			if len(tt.fields.acquiredParts) > 0 {
				for _, tech := range tt.fields.acquiredParts {
					player = player.WithAcquiredTech(tech)
				}
			}
			tc := NewTechComparer(&rules, player)
			design := NewShipDesign(player.Num, 1).WithHull("Nubian").WithPurpose(ShipDesignPurposeTorpedoFighter).WithSpec(&rules, player)
			if tt.fields.beamShip {
				design.Purpose = ShipDesignPurposeBeamFighter
			}
			got := tc.GetBestComponentWithTag(design, tt.args.hullSlotType, 1, tt.args.tag)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TechComparer.GetBestComponentWithTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTechComparer_compareFieldsByTag(t *testing.T) {
	type args struct {
		hc    *TechHullComponent
		other *TechHullComponent
		tag   TechTag
		RS    bool
		light bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Croby Sharmor vs Bear Neutrino Barrier",
			args: args{
				hc:    &CrobySharmor,
				other: &BearNeutrinoBarrier,
				tag:   TechTagShield,
			},
			want: false,
		},
		{
			name: "Armageddon Missile vs Upsilon Torpedo",
			args: args{
				hc:    &ArmageddonMissile,
				other: &UpsilonTorpedo,
				tag:   TechTagTorpedo,
				RS:    false,
				light: false,
			},
			want: false,
		},
		{
			name: "Stargate 300-500 vs Stargate 100-Any",
			args: args{
				hc:    &Stargate300_500,
				other: &Stargate100_Any,
				tag:   TechTagStargate,
				RS:    false,
				light: false,
			},
			want: false,
		},
		{
			name: "Laser vs X Ray Laser",
			args: args{
				hc:    &Laser,
				other: &XRayLaser,
				tag:   TechTagBeamWeapon,
				RS:    false,
				light: false,
			},
			want: true,
		},
		{
			name: "Jammer 50 vs Jammer 20",
			args: args{
				hc:    &Jammer50,
				other: &Jammer20,
				tag:   TechTagTorpedoJammer,
				RS:    false,
				light: false,
			},
			want: false,
		},
		{
			name: "Mega Poly vs Superlatanium",
			args: args{
				hc:    &MegaPolyShell,
				other: &Superlatanium,
				tag:   TechTagArmor,
				RS:    false,
				light: false,
			},
			want: true,
		},
		{
			name: "Battle Nexus vs Battle Super Computer",
			args: args{
				hc:    &BattleNexus,
				other: &BattleSuperComputer,
				tag:   TechTagTorpedoBonus,
				RS:    false,
				light: false,
			},
			want: false,
		},
		{
			name: "Flux Capacitor vs Energy Capacitor",
			args: args{
				hc:    &FluxCapacitor,
				other: &EnergyCapacitor,
				tag:   TechTagBeamCapacitor,
				RS:    false,
				light: false,
			},
			want: false,
		},
		{
			name: "Neutronium versus organic armor with weight penalty",
			args: args{
				hc:    &Neutronium, // 275/1+(45-30)/10 = 275/2.5 = 110 effective dp
				other: &Organic,    // 175 dp
				tag:   TechTagArmor,
				RS:    false,
				light: true,
			},
			want: true,
		},
		{
			name: "Croby versus Neutronium with RS",
			args: args{
				hc:    &Neutronium,   // 275/2 (from RS) = 137.5 total dp
				other: &CrobySharmor, // (60*1.4) + 65 = 149 total dp
				tag:   TechTagArmor,
				RS:    true,
				light: false,
			},
			want: true,
		},
		{
			name: "Colloidal Phaser vs Heavy Blaster",
			args: args{
				hc:    &ColloidalPhaser,
				other: &HeavyBlaster,
				tag:   TechTagBeamWeapon,
				RS:    false,
				light: false,
			},
			want: true,
		},
		{
			name: "AMG vs Super Fuel Tank",
			args: args{
				hc:    &AntiMatterGenerator, // (200 mg + 50*5 gen) / 24
				other: &SuperFuelTank,       // 500 cap / 16
				tag:   TechTagFuelTank,
				RS:    false,
				light: false,
			},
			want: true,
		},
		{
			name: "Ultra Miner vs Alien Miner",
			args: args{
				hc:    &RoboUltraMiner,
				other: &AlienMiner,
				tag:   TechTagMiningRobot,
				RS:    false,
				light: false,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		race := NewRace()
		if tt.args.RS {
			race = race.WithLRT(RS)
		}
		player := NewPlayer(1, race.WithSpec(&rules)).WithTechLevels(TechLevel{26, 26, 26, 26, 26, 26})
		tc := techCompare{&rules, player}
		design := NewShipDesign(player.Num, 1).WithHull("Nubian").WithPurpose(ShipDesignPurposeTorpedoFighter).WithSpec(&rules, player)
		if tt.args.light {
			design.Purpose = ShipDesignPurposeFreighter
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := tc.compareFieldsByTag(design, tt.args.hc, tt.args.other, 1, tt.args.tag); got != tt.want {
				t.Errorf("TechComparer.compareFieldsByTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_techCompare_getMostNeededComponent_ArmorChecks(t *testing.T) {
	type fields struct {
		techLevels    TechLevel
		race          *Race
		acquiredParts []string
		beamShip      bool
		hull          string
	}
	type args struct {
		hullSlotType HullSlotType
		qty          int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TechHullComponent
		wantErr bool
	}{
		{
			name: "Should not use armor on BB due to being inefficient",
			fields: fields{
				techLevels:    TechLevel{14, 14, 14, 14, 14, 14},
				race:          NewRace().WithPRT(IS),
				acquiredParts: []string{},
				beamShip:      false, hull: Battleship.Name,
			},
			args: args{
				hullSlotType: HullSlotTypeArmor,
				qty:          1,
			}, want: nil, wantErr: false,
		},
		{
			name: "Best light armor",
			fields: fields{
				techLevels:    TechLevel{12, 12, 12, 12, 12, 14},
				race:          NewRace().WithPRT(WM),
				acquiredParts: []string{"Multi Cargo Pod"},
				beamShip:      true, hull: Scout.Name,
			},
			args: args{
				hullSlotType: HullSlotTypeArmor,
				qty:          1,
			}, want: &Organic, wantErr: false,
		},
		{
			name: "Best shield/armor item with tech 14 and RS",
			fields: fields{
				techLevels:    TechLevel{14, 14, 14, 14, 14, 14},
				race:          NewRace().WithPRT(IS).WithLRT(RS),
				acquiredParts: []string{"Mega Poly Shell", "Langston Shell", "Multi Cargo Pod"},
				beamShip:      false, hull: Scout.Name,
			},
			args: args{
				hullSlotType: HullSlotTypeShieldArmor,
				qty:          1,
			}, want: &MegaPolyShell, wantErr: false,
		},
		{
			name: "incorrect hull",
			fields: fields{
				techLevels:    TechLevel{14, 14, 14, 14, 14, 14},
				race:          NewRace().WithPRT(IS).WithLRT(RS),
				acquiredParts: []string{"Mega Poly Shell", "Langston Shell", "Multi Cargo Pod"},
				beamShip:      false, hull: "BANANA BOAT",
			},
			args: args{
				hullSlotType: HullSlotTypeShieldArmor,
				qty:          1,
			}, want: nil, wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(1, tt.fields.race.WithSpec(&rules)).WithTechLevels(tt.fields.techLevels).withSpec(&rules)
			for _, part := range tt.fields.acquiredParts {
				player.AcquiredTechs[part] = true
			}
			tc := NewTechComparer(&rules, player)
			design := NewShipDesign(player.Num, 1).WithName(tt.name).WithHull(tt.fields.hull).WithPurpose(ShipDesignPurposeTorpedoFighter)
			if tt.fields.beamShip {
				design.Purpose = ShipDesignPurposeBeamFighter
			}
			got, err := tc.GetMostNeededComponent(design, tt.args.hullSlotType, tt.args.qty)
			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Fatalf("techCompare.getMostNeededComponent() did not return error when expected")
				} else {
					t.Fatalf("techCompare.getMostNeededComponent() errored unexpectedly; err = \n%v", err)
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("techCompare.getMostNeededComponent() = %v, want %v", got, tt.want)
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
			name: "Useless part",
			args: args{
				armorMulti: 1, shieldMulti: 1,
				hc:         &AnnihilatorBomb,
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
			name: "Beam Deflector",
			args: args{
				armorMulti: 1, shieldMulti: 1,
				hc:         &BeamDeflector,
				qty:        1,
				shield:     1,
				armor:      1,
				beamBonus:  1,
				jamming:    0,
				computing:  0,
				deflecting: 0,
				starbase:   false,
			},
			want: 1.1, // 1.1 / 1
		},
		{
			name: "10 Flux Caps on 21% capped ship",
			args: args{
				armorMulti: 1, shieldMulti: 1,
				hc:         &FluxCapacitor,
				qty:        10,
				shield:     1,
				armor:      1,
				beamBonus:  1.21,
				jamming:    0,
				computing:  0,
				deflecting: 0,
				starbase:   false,
			},
			want: 2.11, // 2.55 / 1.21
		},
		{
			name: "1 Mega Poly on armored ship with RS",
			args: args{
				armorMulti: 0.5, shieldMulti: 1.4,
				hc:         &MegaPolyShell,
				qty:        1,
				shield:     460,
				armor:      1000,
				beamBonus:  1,
				jamming:    0,
				computing:  0,
				deflecting: 0,
				starbase:   false,
			},
			want: 1.41, // (1460+140+(200/1.7))/1460 * 1.2 = 1.17 * 1.2 = 1.41
		},
		{
			name: "3 Mega Polys on armored starbase with 10% jam",
			args: args{
				armorMulti: 1, shieldMulti: 1,
				hc:         &MegaPolyShell,
				qty:        3,
				shield:     400,
				armor:      1000,
				beamBonus:  1,
				jamming:    0.1,
				computing:  0,
				deflecting: 0,
				starbase:   false,
			},
			want: 2.12,
		},
	}
	for _, tt := range tests {
		player := NewPlayer(1, NewRace())
		player.Race.Spec.ArmorStrengthFactor = tt.args.armorMulti
		player.Race.Spec.ShieldStrengthFactor = tt.args.shieldMulti
		tc := techCompare{&rules, player}
		design := NewShipDesign(player.Num, 1).WithHull("Battleship").WithSpec(&rules, player)
		design.Spec.Shields = tt.args.shield
		design.Spec.Armor = tt.args.armor
		design.Spec.TorpedoBonus = tt.args.computing
		design.Spec.TorpedoJamming = tt.args.jamming
		design.Spec.BeamBonus = tt.args.beamBonus
		design.Spec.BeamDefense = tt.args.deflecting
		design.Spec.Starbase = tt.args.starbase
		t.Run(tt.name, func(t *testing.T) {
			// round the result to 2 decimal places for easier testing
			if got := roundFloat(tc.getWarshipPartBonus(design, tt.args.hc, tt.args.qty), 2); got != tt.want {
				t.Errorf("ShipDesign.getWarshipPartBonus() = %v, want %v", got, tt.want)
			}
		})
	}
}
