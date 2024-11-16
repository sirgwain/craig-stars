package cs

import "testing"

func TestTechComparer_CompareFieldsByTag(t *testing.T) {
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
				RS:    false,
				light: false,
			},
			want: false,
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
			name: "Flux Capacitator vs Energy Capacitator",
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
				hc:    &Neutronium,   // 275/2 = 137.5 total dp
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
				hc:    &AntiMatterGenerator, // 450 mg / 24 
				other: &SuperFuelTank, // 500 mg / 16
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
		tc := NewTechComparer()
		t.Run(tt.name, func(t *testing.T) {
			if got := tc.CompareFieldsByTag(player, tt.args.hc, tt.args.other, tt.args.tag, tt.args.light); got != tt.want {
				t.Errorf("TechComparer.CompareFieldsByTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTechComparer_GetBestComponentWithTags(t *testing.T) { // TODO: Fix nil pointer error... again
	type fields struct {
		techLevels    TechLevel
		race          *Race
		acquiredParts []string
		beamShip      bool
	}
	type args struct {
		hullSlotType HullSlotType
		qty          int
		tags         []TechTag
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TechHullComponent
		wantErr bool
	}{
		{
			name: "Best beam with max techs",
			fields: fields{
				techLevels:    TechLevel{1, 1, 26, 1, 1, 1},
				race:          NewRace().WithPRT(JoaT),
				acquiredParts: []string{},
				beamShip:      false,
			},
			args: args{
				hullSlotType: HullSlotTypeWeaponShield,
				qty:          99,
				tags:         []TechTag{TechTagBeamWeapon},
			}, want: &AntiMatterPulverizer, wantErr: false,
		},
		{
			name: "Best shield/armor with tech 14",
			fields: fields{
				techLevels:    TechLevel{14, 14, 14, 14, 14, 14},
				race:          NewRace().WithPRT(IS),
				acquiredParts: []string{"Mega Poly Shell", "Langston Shell", "Multi Cargo Pod"},
				beamShip:      false,
			},
			args: args{
				hullSlotType: HullSlotTypeGeneral,
				qty:          1,
				tags:         []TechTag{TechTagShield, TechTagArmor},
			}, want: &MegaPolyShell, wantErr: false,
		},
		{
			name: "Best light armor",
			fields: fields{
				techLevels:    TechLevel{12, 12, 12, 12, 12, 14},
				race:          NewRace().WithPRT(WM),
				acquiredParts: []string{"Multi Cargo Pod"},
				beamShip:      true,
			},
			args: args{
				hullSlotType: HullSlotTypeGeneral,
				qty:          1,
				tags:         []TechTag{TechTagArmor},
			}, want: &Organic, wantErr: false,
		},
		{
			name: "Best shield/armor with tech 14 and RS",
			fields: fields{
				techLevels:    TechLevel{14, 14, 14, 14, 14, 14},
				race:          NewRace().WithPRT(IS).WithLRT(RS),
				acquiredParts: []string{"Mega Poly Shell", "Langston Shell", "Multi Cargo Pod"},
				beamShip:      false,
			},
			args: args{
				hullSlotType: HullSlotTypeShieldArmor,
				qty:          1,
				tags:         []TechTag{TechTagShield, TechTagArmor},
			}, want: &MegaPolyShell, wantErr: false,
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
				hullSlotType: HullSlotTypeGeneral,
				qty:          99,
				tags:         []TechTag{TechTagMineLayer},
			}, want: nil, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := NewTechComparer()
			player := testPlayer().WithTechLevels(tt.fields.techLevels)
			player.Race = *tt.fields.race
			if len(tt.fields.acquiredParts) > 0 {
				for _, tech := range tt.fields.acquiredParts {
					player = player.WithAcquiredTech(tech)
				}
			}
			design := NewShipDesign(player, 1).WithHull("Nubian").WithPurpose(ShipDesignPurposeTorpedoFighter)
			if tt.fields.beamShip {
				design.Purpose = ShipDesignPurposeBeamFighter
			}
			got, err := tc.GetBestComponentWithTags(&rules, player, design, tt.args.hullSlotType, tt.args.tags...)
			if (err != nil) != tt.wantErr {
				t.Errorf("TechComparer.GetBestComponentWithTags() errored unexpectedly; error = %v", err)
			}
			if got != tt.want {
				t.Errorf("TechComparer.GetBestComponentWithTags() = %v, want %v", got.Name, tt.want.Name)
			}
		})
	}
}
