package cs

import "testing"

func TestTechHullComponent_CompareFieldsByTag(t *testing.T) {
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
			name: "Croby vs bear",
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
			name: "Battle Nexus vs Super computer",
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
			name: "Flux cap vs en cap",
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
			name: "Heavy versus light armor",
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
			name: "Colloidal vs Heavy Blaster",
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
		design := NewShipDesign(player, 1).WithPurpose(ShipDesignPurposeTorpedoFighter)
		if tt.args.light {
			design.Purpose = ShipDesignPurposeBeamFighter
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := design.CompareFieldsByTag(player, tt.args.hc, tt.args.other, tt.args.tag); got != tt.want {
				t.Errorf("TechHullComponent.CompareFieldsByTag() = %v, want %v", got, tt.want)
			}
		})
	}
}