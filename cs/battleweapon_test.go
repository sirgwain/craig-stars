//go:build !wasi && !wasm

package cs

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
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
		{"1 laser, 0 range", args{damage: 10, weaponRange: 1, dist: 0, beamDefense: 1}, 10},
		{"1 laser, 1 range", args{damage: 10, weaponRange: 1, dist: 1, beamDefense: 1}, 9},
		{"2 colloidal phasers, 3 range", args{damage: 52, weaponRange: 3, dist: 3, beamDefense: 1}, 47}, // real Stars! is 48...
		{"1 laser, 0 range, 1 deflector", args{damage: 10, weaponRange: 1, dist: 0, beamDefense: .9}, 9},
		{"1 laser, 1 range, 1 deflector", args{damage: 10, weaponRange: 1, dist: 1, beamDefense: .9}, 8},
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
		position Vector
	}
	tests := []struct {
		name   string
		weapon battleWeaponSlot
		args   args
		want   bool
	}{
		{"no distance, in range", battleWeaponSlot{token: &battleToken{BattleRecordToken: BattleRecordToken{Position: Vector{0, 0}}}}, args{Vector{0, 0}}, true},
		{"distance 1, in range", battleWeaponSlot{token: &battleToken{BattleRecordToken: BattleRecordToken{Position: Vector{0, 0}}}, weaponRange: 1}, args{Vector{1, 1}}, true},
		{"distance 2, out of range", battleWeaponSlot{token: &battleToken{BattleRecordToken: BattleRecordToken{Position: Vector{0, 0}}}, weaponRange: 1}, args{Vector{1, 2}}, false},
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
			want: 2.0 / 3.0,
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
			want: .325,
		},
		{
			name:   "torpedo, attractiveness 1",
			fields: fields{weaponType: battleWeaponTypeTorpedo, accuracy: .45},
			args: args{
				cost:    Cost{Boranium: 1, Resources: 1},
				armor:   1,
				shields: 1,
			},
			want: .51,
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
				ShipToken:      &ShipToken{Quantity: 1},
				cost:           tt.args.cost,
				armor:          tt.args.armor,
				shields:        tt.args.shields,
				stackShields:   tt.args.shields,
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
		{"jihad missile, 3 BC 20/30, 3 jammer 10/20", fields{accuracy: .20, torpedoBonus: .8244}, args{torpedoJamming: .6268}, .35808},
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

func Test_battleWeaponSlot_getBeamDamageToTarget(t *testing.T) {
	type fields struct {
		position           Vector
		shipQuantity       int
		slotQuantity       int
		weaponRange        int
		damagesShieldsOnly bool
		hitsAllTargets     bool
	}
	type args struct {
		damage               int
		position             Vector
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
				position:      Vector{1, 0},
				damage:        10,
				tokenQuantity: 1,
				armor:         20,
				shields:       0,
			},
			want: battleWeaponDamage{armorDamage: 9, damage: 9, quantityDamaged: 1},
		},
		{
			name:   "1 laser, 1 range away, 1 deflector, 20dp target, 8 damage done",
			fields: fields{shipQuantity: 1, slotQuantity: 1, weaponRange: 1},
			args: args{
				position:      Vector{1, 0},
				damage:        10,
				tokenQuantity: 1,
				armor:         20,
				shields:       0,
				beamDefense:   .9,
			},
			want: battleWeaponDamage{armorDamage: 8, damage: 8, quantityDamaged: 1},
		},
		{
			name:   "1 gattling, 1 range away, 1 deflector, 20dp target, 9 damage done",
			fields: fields{shipQuantity: 1, slotQuantity: 1, weaponRange: 1, hitsAllTargets: true},
			args: args{
				position:      Vector{1, 0},
				damage:        10,
				tokenQuantity: 1,
				armor:         20,
				shields:       0,
				beamDefense:   .9,
			},
			want: battleWeaponDamage{armorDamage: 9, damage: 9, quantityDamaged: 1},
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
				position:      Vector{1, 0}, // 1 away
				damage:        60,           // 3 lasers * 2 ships * 10 damage each
				tokenQuantity: 2,
				armor:         20,
				shields:       40,
			},
			want: battleWeaponDamage{shieldDamage: 40, armorDamage: 14, numDestroyed: 0, leftover: 0, damage: 7, quantityDamaged: 2},
		},
		{
			name:   "1 beam 75 damage, 2 targets with 2@50% damage",
			fields: fields{shipQuantity: 1, slotQuantity: 1, weaponRange: 1},
			args: args{
				damage:               75, // 1 big beam
				tokenQuantity:        2,
				armor:                100,
				tokenDamage:          50,
				tokenQuantityDamaged: 2,
			},
			want: battleWeaponDamage{shieldDamage: 0, armorDamage: 75, numDestroyed: 1, leftover: 0, damage: 75, quantityDamaged: 1},
		},
		{
			name:   "Shield sapper breaks shields with leftover",
			fields: fields{shipQuantity: 1, slotQuantity: 1, weaponRange: 1, damagesShieldsOnly: true},
			args: args{
				damage:        100,
				tokenQuantity: 1,
				armor:         100,
				shields:       50,
			},
			want: battleWeaponDamage{shieldDamage: 50, leftover: 50},
		},
		{
			name:   "Sappper no break shields",
			fields: fields{shipQuantity: 1, slotQuantity: 1, weaponRange: 1, damagesShieldsOnly: true},
			args: args{
				damage:        100,
				tokenQuantity: 1,
				armor:         100,
				shields:       500,
			},
			want: battleWeaponDamage{shieldDamage: 100},
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
				hitsAllTargets:     tt.fields.hitsAllTargets,
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
				t.Errorf("battleWeaponSlot.getTargetBeamDamage() = \n%#v\n, want \n%#v", got, tt.want)
			}
		})
	}
}

// Test_battleWeaponSlot_getTorpedoVolleyDamage_expectedHits verifies volley damage
// when hits and misses are split by average accuracy.
func Test_battleWeaponSlot_getTorpedoVolleyDamage_expectedHits(t *testing.T) {
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
			name: "4 powerful torpedoes but only 2 ships to destroy",
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

			count := float64(tt.fields.slotQuantity * tt.fields.shipQuantity)
			hits := count * weapon.getAccuracy(target.torpedoJamming)
			result := weapon.getTorpedoVolleyDamage(target, hits, count-hits, int(count), rules.TorpedoSplashDamage)
			got := battleWeaponDamage{shieldDamage: result.shieldDamage, armorDamage: result.armorDamage}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("battleWeaponSlot.getTorpedoVolleyDamage() = \n%#v\nwant: \n%#v", got, tt.want)
			}
		})
	}
}

// Test_battleWeaponSlot_getTorpedoVolleyDamage_singleHit verifies a volley of one torpedo that hits.
func Test_battleWeaponSlot_getTorpedoVolleyDamage_singleHit(t *testing.T) {
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
			if got := weapon.getTorpedoVolleyDamage(target, 1, 0, 1, rules.TorpedoSplashDamage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("battleWeaponSlot.getTorpedoVolleyDamage() = \n%#v\nwant: \n%#v", got, tt.want)
			}
		})
	}
}

// Test_battleWeaponSlot_getAttractiveness_remainingDefenses verifies target ranking from remaining defenses and deflectors.
func Test_battleWeaponSlot_getAttractiveness_remainingDefenses(t *testing.T) {
	tests := []struct {
		name    string
		shields int
		damage  float64
		sapper  bool
		defense float64
		want    float64
	}{
		{"remaining armor raises attractiveness", 0, 50, false, 1, 100.0 / 51},
		{"remaining shields affect attractiveness", 10, 0, false, 1, 100.0 / 111},
		{"deflectors lower attractiveness", 0, 0, false, .5, 50.0 / 101},
		{"sapper ranks shield defense only", 10, 0, true, 1, 10},
		{"sapper cannot target an empty shield", 0, 0, true, 1, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := &battleToken{ShipToken: &ShipToken{Quantity: 1, Damage: tt.damage, QuantityDamaged: 1}, armor: 100, stackShields: tt.shields, cost: Cost{Resources: 100}, beamDefense: tt.defense}
			weapon := &battleWeaponSlot{damagesShieldsOnly: tt.sapper}
			if got := weapon.getAttractiveness(target); !test.WithinTolerance(got, tt.want, .001) {
				t.Errorf("battleWeaponSlot.getAttractiveness() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_getBattleArmorDamage verifies fractional armor damage and the volley kill limit.
func Test_getBattleArmorDamage(t *testing.T) {
	tests := []struct {
		name                     string
		damage, existing         float64
		damaged, quantity, limit int
		want                     battleWeaponDamage
	}{
		{name: "leftover damage does at least one point", damage: 1, existing: 99.5, damaged: 1, quantity: 2, limit: 2, want: battleWeaponDamage{armorDamage: 1, numDestroyed: 1, damage: 1, quantityDamaged: 1}},
		{name: "kill limit discards excess without harming survivors", damage: 200, quantity: 3, limit: 1, want: battleWeaponDamage{armorDamage: 100, numDestroyed: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := &battleToken{ShipToken: &ShipToken{Quantity: tt.quantity, Damage: tt.existing, QuantityDamaged: tt.damaged}, armor: 100}
			assert.Equal(t, tt.want, getBattleArmorDamage(target, tt.damage, tt.limit))
		})
	}
}

// Test_battleWeaponSlot_getAccuracy_bounds verifies accuracy limits after computers and jamming.
func Test_battleWeaponSlot_getAccuracy_bounds(t *testing.T) {
	tests := []struct {
		name                       string
		base, bonus, jamming, want float64
	}{
		{"zero base remains zero", 0, 1, 0, 0},
		{"positive accuracy has a minimum", .01, 0, .95, .01},
		{"accuracy cannot exceed one", .5, 2, 0, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weapon := &battleWeaponSlot{accuracy: tt.base, torpedoBonus: tt.bonus}
			if got := weapon.getAccuracy(tt.jamming); !test.WithinTolerance(got, tt.want, .001) {
				t.Errorf("battleWeaponSlot.getAccuracy() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_battleWeaponSlot_beamSpillover verifies that each target applies its own range attenuation and deflectors.
func Test_battleWeaponSlot_beamSpillover(t *testing.T) {
	tests := []struct {
		name                  string
		baseRange, rangeBonus int
		firstDefense          float64
		wantDamage            float64
	}{
		{"each target applies its own deflectors", 1, 0, .5, 34},
		{"starbase bonus does not dilute range attenuation", 1, 1, 1, 64},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRulesWithSeed(1)
			attacker := &battleToken{ShipToken: &ShipToken{Quantity: 1}}
			first := &battleToken{BattleRecordToken: BattleRecordToken{Position: Vector{1, 0}}, ShipToken: &ShipToken{Quantity: 1}, armor: 30, beamDefense: tt.firstDefense}
			if tt.rangeBonus > 0 {
				first.Position = Vector{2, 0}
			}
			second := &battleToken{ShipToken: &ShipToken{Quantity: 1}, armor: 1000}
			weapon := &battleWeaponSlot{token: attacker, power: 100, slotQuantity: 1, beamBonus: 1, weaponRange: tt.baseRange + tt.rangeBonus, rangeBonus: tt.rangeBonus}
			b := &battle{rules: &r, round: 1, record: newBattleRecord(1, None, Vector{}, nil)}
			b.record.recordNewRound()
			b.fireBeamWeapon(weapon, []*battleToken{first, second})
			assert.Equal(t, 0, first.Quantity)
			if !test.WithinTolerance(second.Damage, tt.wantDamage, .001) {
				t.Errorf("second target damage = %v, want %v", second.Damage, tt.wantDamage)
			}
		})
	}
}

// Test_roundBattleArmorDamage verifies that stored ship damage rounds up to 1/500th of armor.
func Test_roundBattleArmorDamage(t *testing.T) {
	tests := []struct {
		name          string
		damage, armor float64
		want          float64
	}{
		{"damage below 0.2% rounds up to 0.2%", 1, 1000, 2},
		{"damage below one point rounds up to one point", .01, 100, 1},
		{"damage on a step is unchanged", 4, 1000, 4},
		{"damage between steps rounds up", 340, 350, 340.2},
		{"surviving ships keep at least one step of armor", 99.99, 100, 99.8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := roundBattleArmorDamage(tt.damage, tt.armor); !test.WithinTolerance(got, tt.want, 1e-6) {
				t.Errorf("roundBattleArmorDamage() = %v, want %v", got, tt.want)
			}
		})
	}
}
