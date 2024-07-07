package cs

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/test"
)

func Test_getBeamDamageAtDistance(t *testing.T) {
	type args struct {
		damage      int
		weaponRange int
		dist        int
		beamDefense float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"1 laser, 0 range", args{damage: 10, weaponRange: 1, dist: 0}, 10},
		{"1 laser, 1 range", args{damage: 10, weaponRange: 1, dist: 1}, 9},
		{"2 colloidal phasers, 3 range", args{damage: 52, weaponRange: 3, dist: 3}, 47}, // real Stars! is 48...
		{"1 laser, 0 range, 1 deflector", args{damage: 10, weaponRange: 1, dist: 0, beamDefense: .1}, 9},
		{"1 laser, 1 range, 1 deflector", args{damage: 10, weaponRange: 1, dist: 1, beamDefense: .1}, 8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getBeamDamageAtDistance(tt.args.damage, tt.args.weaponRange, tt.args.dist, tt.args.beamDefense, rules.BeamRangeDropoff); got != tt.want {
				t.Errorf("getBeamDamageAtRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_battleWeaponSlot_isInRangePosition(t *testing.T) {
	type args struct {
		position BattleVector
	}
	tests := []struct {
		name   string
		weapon battleWeaponSlot
		args   args
		want   bool
	}{
		{"no distance, in range", battleWeaponSlot{token: &battleToken{BattleRecordToken: BattleRecordToken{Position: BattleVector{0, 0}}}}, args{BattleVector{0, 0}}, true},
		{"distance 1, in range", battleWeaponSlot{token: &battleToken{BattleRecordToken: BattleRecordToken{Position: BattleVector{0, 0}}}, weaponRange: 1}, args{BattleVector{1, 1}}, true},
		{"distance 2, out of range", battleWeaponSlot{token: &battleToken{BattleRecordToken: BattleRecordToken{Position: BattleVector{0, 0}}}, weaponRange: 1}, args{BattleVector{1, 2}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.weapon.isInRangePosition(tt.args.position); got != tt.want {
				t.Errorf("battleWeaponSlot.isInRangePosition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_battleWeaponSlot_getAttractiveness(t *testing.T) {
	type fields struct {
		weaponType         battleWeaponType
		accuracy           float64
		damagesShieldsOnly bool
		capitalShipMissile bool
	}
	type args struct {
		cost           Cost
		armor          int
		shields        int
		beamDefense    float64
		torpedoJamming float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name:   "beam, attractiveness 1",
			fields: fields{weaponType: battleWeaponTypeBeam},
			args: args{
				cost:    Cost{Boranium: 1, Resources: 1},
				armor:   1,
				shields: 1,
			},
			want: 1,
		},
		{
			name:   "torpedo, more shields than armor",
			fields: fields{weaponType: battleWeaponTypeTorpedo, accuracy: .45},
			args: args{
				cost:    Cost{Boranium: 1, Resources: 1},
				armor:   1,
				shields: 2,
			},
			want: .45,
		},
		{
			name:   "torpedo, more armor than shields",
			fields: fields{weaponType: battleWeaponTypeTorpedo, accuracy: .45},
			args: args{
				cost:    Cost{Boranium: 1, Resources: 1},
				armor:   2,
				shields: 1,
			},
			want: .3,
		},
		{
			name:   "torpedo, attractiveness 1",
			fields: fields{weaponType: battleWeaponTypeTorpedo, accuracy: .45},
			args: args{
				cost:    Cost{Boranium: 1, Resources: 1},
				armor:   1,
				shields: 1,
			},
			want: .45,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weapon := &battleWeaponSlot{
				weaponType:         tt.fields.weaponType,
				damagesShieldsOnly: tt.fields.damagesShieldsOnly,
				accuracy:           tt.fields.accuracy,
				capitalShipMissile: tt.fields.capitalShipMissile,
			}
			target := &battleToken{
				cost:           tt.args.cost,
				armor:          tt.args.armor,
				shields:        tt.args.shields,
				beamDefense:    tt.args.beamDefense,
				torpedoJamming: tt.args.torpedoJamming,
			}
			if got := weapon.getAttractiveness(target); !test.WithinTolerance(got, tt.want, .01) {
				t.Errorf("battleWeaponSlot.getAttractiveness() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_battleWeaponSlot_getAccuracy(t *testing.T) {
	type fields struct {
		accuracy     float64
		torpedoBonus float64
	}
	type args struct {
		torpedoJamming float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{"beta torpedo", fields{accuracy: .45}, args{}, .45},
		{"beta torpedo, 1 BC", fields{accuracy: .45, torpedoBonus: .2}, args{}, .56},
		{"beta torpedo, 1 BC, 1 jammer 20", fields{accuracy: .45, torpedoBonus: .2}, args{torpedoJamming: .2}, .45},
		{"beta torpedo, 1 BC, 1 jammer 10", fields{accuracy: .45, torpedoBonus: .2}, args{torpedoJamming: .1}, .505},
		{"beta torpedo, 1 BC, 1 jammer 30", fields{accuracy: .45, torpedoBonus: .1}, args{torpedoJamming: .2}, .405},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weapon := &battleWeaponSlot{
				accuracy:     tt.fields.accuracy,
				torpedoBonus: tt.fields.torpedoBonus,
			}
			if got := weapon.getAccuracy(tt.args.torpedoJamming); got != tt.want {
				t.Errorf("battleWeaponSlot.getAccuracy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_battleWeaponSlot_getDamage(t *testing.T) {
	type fields struct {
		weaponType  battleWeaponType
		weaponRange int
		power       int
	}
	type args struct {
		dist        int
		beamDefense float64
		beamDropoff float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{"torpedo, base damage", fields{weaponType: battleWeaponTypeTorpedo, power: 10}, args{}, 10},
		{
			name:   "laser, 0 range",
			fields: fields{weaponType: battleWeaponTypeBeam, weaponRange: 1, power: 10},
			args:   args{dist: 0, beamDropoff: .1},
			want:   10,
		},
		{
			name:   "laser, 1 range",
			fields: fields{weaponType: battleWeaponTypeBeam, weaponRange: 1, power: 10},
			args:   args{dist: 1, beamDropoff: .1},
			want:   9,
		},
		{
			name:   "ColloidalPhaser, 3 range",
			fields: fields{weaponType: battleWeaponTypeBeam, weaponRange: 3, power: 26},
			args:   args{dist: 3, beamDropoff: .1},
			want:   24,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weapon := &battleWeaponSlot{
				weaponType:  tt.fields.weaponType,
				weaponRange: tt.fields.weaponRange,
				power:       tt.fields.power,
			}
			if got := weapon.getDamage(tt.args.dist, tt.args.beamDefense, tt.args.beamDropoff); got != tt.want {
				t.Errorf("battleWeaponSlot.getDamage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_battleWeaponSlot_getBeamDamageToTarget(t *testing.T) {
	type fields struct {
		position           BattleVector
		shipQuantity       int
		slotQuantity       int
		weaponRange        int
		damagesShieldsOnly bool
	}
	type args struct {
		damage               int
		position             BattleVector
		armor                int
		shields              int
		beamDefense          float64
		tokenQuantity        int
		tokenDamage          float64
		tokenQuantityDamaged int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   battleWeaponDamage
	}{
		{
			name:   "1 laser, 20dp target, 10 damage done",
			fields: fields{shipQuantity: 1, slotQuantity: 1, weaponRange: 1},
			args: args{
				damage:        10,
				tokenQuantity: 1,
				armor:         20,
				shields:       0,
			},
			want: battleWeaponDamage{armorDamage: 10, damage: 10, quantityDamaged: 1},
		},
		{
			name:   "1 laser, 1 range away, 20dp target, 9 damage done",
			fields: fields{shipQuantity: 1, slotQuantity: 1, weaponRange: 1},
			args: args{
				position:      BattleVector{1, 0},
				damage:        10,
				tokenQuantity: 1,
				armor:         20,
				shields:       0,
			},
			want: battleWeaponDamage{armorDamage: 9, damage: 9, quantityDamaged: 1},
		},
		{
			name:   "1 laser, 1 range away, 1 defelctor, 20dp target, 8 damage done",
			fields: fields{shipQuantity: 1, slotQuantity: 1, weaponRange: 1},
			args: args{
				position:      BattleVector{1, 0},
				damage:        10,
				tokenQuantity: 1,
				armor:         20,
				shields:       0,
				beamDefense:   .1,
			},
			want: battleWeaponDamage{armorDamage: 8, damage: 8, quantityDamaged: 1},
		},
		{
			name:   "1 laser, 20 shields 20dp target, 10 damage done",
			fields: fields{shipQuantity: 1, slotQuantity: 1, weaponRange: 1},
			args: args{
				damage:        10,
				tokenQuantity: 1,
				armor:         20,
				shields:       20,
			},
			want: battleWeaponDamage{shieldDamage: 10},
		},
		{
			name:   "3 lasers, 20 shields 20dp target, 20 shield, 10 armor damage done",
			fields: fields{shipQuantity: 1, slotQuantity: 3, weaponRange: 1},
			args: args{
				damage:        30, // 3 lasers * 10 damage each
				tokenQuantity: 1,
				armor:         20,
				shields:       20,
			},
			want: battleWeaponDamage{shieldDamage: 20, armorDamage: 10, damage: 10, quantityDamaged: 1},
		},
		{
			name:   "2 ships, 3 lasers, 20 shields 20dp target, destroyed",
			fields: fields{shipQuantity: 2, slotQuantity: 3, weaponRange: 1},
			args: args{
				damage:        60, // 3 lasers * 2 ships * 10 damage each
				tokenQuantity: 1,
				armor:         20,
				shields:       20,
			},
			want: battleWeaponDamage{shieldDamage: 20, armorDamage: 20, numDestroyed: 1, leftover: 20, damage: 0, quantityDamaged: 0},
		},
		{
			name:   "2 ships, 3 lasers, 2 targets with 20x2 shields 20dp, destroy one",
			fields: fields{shipQuantity: 2, slotQuantity: 3, weaponRange: 1},
			args: args{
				damage:        60, // 3 lasers * 2 ships * 10 damage each
				tokenQuantity: 2,
				armor:         20,
				shields:       40,
			},
			want: battleWeaponDamage{shieldDamage: 40, armorDamage: 20, numDestroyed: 1, leftover: 0, damage: 0, quantityDamaged: 0},
		},
		{
			name:   "2 ships, 3 lasers, 2 targets 1sq away with 20x2 shields 20dp, damage both",
			fields: fields{shipQuantity: 2, slotQuantity: 3, weaponRange: 1},
			args: args{
				position:      BattleVector{1, 0}, // 1 away
				damage:        60,                 // 3 lasers * 2 ships * 10 damage each
				tokenQuantity: 2,
				armor:         20,
				shields:       40,
			},
			want: battleWeaponDamage{shieldDamage: 40, armorDamage: 14, numDestroyed: 0, leftover: 0, damage: 7, quantityDamaged: 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weapon := &battleWeaponSlot{
				token: &battleToken{
					BattleRecordToken: BattleRecordToken{Position: tt.fields.position},
					ShipToken: &ShipToken{
						Quantity: tt.fields.shipQuantity,
					},
				},
				slotQuantity:       tt.fields.slotQuantity,
				weaponType:         battleWeaponTypeBeam,
				weaponRange:        tt.fields.weaponRange,
				damagesShieldsOnly: tt.fields.damagesShieldsOnly,
			}

			target := &battleToken{
				BattleRecordToken: BattleRecordToken{Position: tt.args.position},
				ShipToken: &ShipToken{
					Quantity:        tt.args.tokenQuantity,
					Damage:          tt.args.tokenDamage,
					QuantityDamaged: tt.args.tokenQuantityDamaged,
				},
				armor:        tt.args.armor,
				stackShields: tt.args.shields,
				beamDefense:  tt.args.beamDefense,
			}
			if got := weapon.getBeamDamageToTarget(tt.args.damage, target, rules.BeamRangeDropoff); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("battleWeaponSlot.getTargetBeamDamage() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_battleWeaponSlot_getEstimatedTorpedoDamageToTarget(t *testing.T) {
	type fields struct {
		shipQuantity int
		slotQuantity int
		weaponPower  int
		accuracy     float64
		torpedoBonus float64
	}
	type args struct {
		armor                int
		shields              int
		tokenQuantity        int
		tokenDamage          float64
		tokenQuantityDamaged int
		torpedoJamming       float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   battleWeaponDamage
	}{
		{
			name: "1 beta, 20dp target, 12 damage done",
			fields: fields{
				weaponPower:  12,
				shipQuantity: 1,
				slotQuantity: 1,
				accuracy:     1, // always hits
			},
			args: args{
				tokenQuantity: 1,
				armor:         20,
				shields:       0,
			},
			want: battleWeaponDamage{armorDamage: 12},
		},
		{
			name: "2 beta, 20dp target, 12 damage done, 50% accuracy",
			fields: fields{
				weaponPower:  12,
				shipQuantity: 1,
				slotQuantity: 2,
				accuracy:     .5,
			},
			args: args{
				tokenQuantity: 1,
				armor:         20,
				shields:       0,
			},
			want: battleWeaponDamage{armorDamage: 12},
		},
		{
			name: "2 beta, 20 shields, 20dp target, 12 damage done, 50% accuracy",
			fields: fields{
				weaponPower:  12,
				shipQuantity: 1,
				slotQuantity: 2,
				accuracy:     .5,
			},
			args: args{
				tokenQuantity: 1,
				armor:         20,
				shields:       20,
			},
			// half damage shields/armor + 1/8th damage to shields for the miss
			want: battleWeaponDamage{shieldDamage: 8, armorDamage: 6},
		},
		{
			name: "4 beta, 20 shields, 20dp target, 24 damage done, 50% accuracy",
			fields: fields{
				weaponPower:  12,
				shipQuantity: 2,
				slotQuantity: 2,
				accuracy:     .5,
			},
			args: args{
				tokenQuantity: 1,
				armor:         20,
				shields:       20,
			},
			// half damage shields/armor + 1/8th damage to shields for the two misses
			want: battleWeaponDamage{shieldDamage: 15, armorDamage: 12},
		},
		{
			name: "4 powerful torpedos but only 2 ships to destroy",
			fields: fields{
				weaponPower:  100,
				shipQuantity: 2,
				slotQuantity: 2,
				accuracy:     1,
			},
			args: args{
				tokenQuantity: 2,
				armor:         20,
				shields:       40,
			},
			// half damage shields/armor + 1/8th damage to shields for the two misses
			want: battleWeaponDamage{shieldDamage: 40, armorDamage: 40},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weapon := &battleWeaponSlot{
				token: &battleToken{
					BattleRecordToken: BattleRecordToken{},
					ShipToken: &ShipToken{
						Quantity: tt.fields.shipQuantity,
					},
				},
				slotQuantity: tt.fields.slotQuantity,
				weaponType:   battleWeaponTypeTorpedo,
				power:        tt.fields.weaponPower,
				accuracy:     tt.fields.accuracy,
				torpedoBonus: tt.fields.torpedoBonus,
			}

			target := &battleToken{
				BattleRecordToken: BattleRecordToken{},
				ShipToken: &ShipToken{
					Quantity:        tt.args.tokenQuantity,
					Damage:          tt.args.tokenDamage,
					QuantityDamaged: tt.args.tokenQuantityDamaged,
				},
				armor:          tt.args.armor,
				stackShields:   tt.args.shields,
				torpedoJamming: tt.args.torpedoJamming,
			}

			if got := weapon.getEstimatedTorpedoDamageToTarget(target); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("battleWeaponSlot.getTargetBeamDamage() = \n%#v\nwant: \n%#v", got, tt.want)
			}
		})
	}
}

func Test_battleWeaponSlot_getTorpedoDamageToTarget(t *testing.T) {
	type args struct {
		weaponPower          int
		armor                int
		shields              int
		tokenQuantity        int
		tokenDamage          float64
		tokenQuantityDamaged int
	}
	tests := []struct {
		name string
		args args
		want battleWeaponDamage
	}{
		{
			name: "1 beta, 20dp target, 12 damage done",
			args: args{
				weaponPower:   12,
				tokenQuantity: 1,
				armor:         20,
				shields:       0,
			},
			want: battleWeaponDamage{armorDamage: 12, damage: 12, quantityDamaged: 1},
		},
		{
			name: "1 beta, 20 shields 20dp target, 6 damage done to shields and armor",
			args: args{
				weaponPower:   12,
				tokenQuantity: 1,
				armor:         20,
				shields:       20,
			},
			want: battleWeaponDamage{shieldDamage: 6, armorDamage: 6, damage: 6, quantityDamaged: 1},
		},
		{
			name: "delta torpedo, 12 shields 20dp target, 13 armor damage done",
			args: args{
				weaponPower:   26,
				tokenQuantity: 1,
				armor:         20,
				shields:       12,
			},
			want: battleWeaponDamage{shieldDamage: 12, armorDamage: 14, damage: 14, quantityDamaged: 1},
		},
		{
			name: "big 60 power torpedo 20 shields 20dp target, destroyed",
			args: args{
				weaponPower:   60,
				tokenQuantity: 1,
				armor:         20,
				shields:       20,
			},
			want: battleWeaponDamage{shieldDamage: 20, armorDamage: 20, numDestroyed: 1, damage: 0, quantityDamaged: 0},
		},
		{
			name: "big torpedo, 2 targets with 20x2 shields 20dp, destroy one",
			args: args{
				weaponPower:   40,
				tokenQuantity: 2,
				armor:         20,
				shields:       40,
			},
			want: battleWeaponDamage{shieldDamage: 20, armorDamage: 20, numDestroyed: 1, leftover: 0, damage: 0, quantityDamaged: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weapon := &battleWeaponSlot{
				token: &battleToken{
					BattleRecordToken: BattleRecordToken{},
					ShipToken:         &ShipToken{},
				},
				power:      tt.args.weaponPower,
				weaponType: battleWeaponTypeTorpedo,
			}

			target := &battleToken{
				BattleRecordToken: BattleRecordToken{},
				ShipToken: &ShipToken{
					Quantity:        tt.args.tokenQuantity,
					Damage:          tt.args.tokenDamage,
					QuantityDamaged: tt.args.tokenQuantityDamaged,
				},
				armor:        tt.args.armor,
				stackShields: tt.args.shields,
			}
			if got := weapon.getTorpedoDamageToTarget(target); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("battleWeaponSlot.getTargetBeamDamage() = \n%#v\nwant: \n%#v", got, tt.want)
			}
		})
	}
}
