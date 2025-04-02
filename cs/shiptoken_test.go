package cs

import (
	"reflect"
	"testing"
)

func TestShipToken_applyMineDamage(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules))
	design := NewShipDesign(player.Num, 1)

	// set some spec values we care about
	design.Spec.Mass = 100
	design.Spec.Armor = 100

	designShielded := NewShipDesign(player.Num, 1)
	designShielded.Spec.Mass = 100
	designShielded.Spec.Armor = 150
	designShielded.Spec.Shields = 50

	type fields struct {
		quantity        int
		damage          float64
		quantityDamaged int
		design          *ShipDesign
	}
	tests := []struct {
		name                string
		fields              fields
		damage              int
		want                tokenDamage
		wantQuantity        int
		wantDamage          float64
		wantQuantityDamaged int
	}{
		{
			name: "ship remains intact",
			fields: fields{
				design:   design,
				quantity: 1,
			},
			damage: 50,
			want: tokenDamage{
				damage:         50,
				shipsDestroyed: 0,
			},
			wantQuantity:        1,
			wantQuantityDamaged: 1,
			wantDamage:          50,
		},
		{
			name: "most ships destroyed",
			fields: fields{
				design:   design,
				quantity: 10,
			},
			damage: 850,
			want: tokenDamage{
				damage:         850,
				shipsDestroyed: 8,
			},
			wantQuantity:        2,
			wantQuantityDamaged: 2,
			wantDamage:          25,
		},
		{
			name: "destroy partially damaged ships",
			fields: fields{
				design:          design,
				quantity:        2,
				quantityDamaged: 1,
				damage:          50,
			},
			damage: 75,
			want: tokenDamage{
				damage:         75,
				shipsDestroyed: 1,
			},
			wantQuantity:        1,
			wantQuantityDamaged: 1,
			wantDamage:          25,
		},
		{
			name: "shields & armor split damage",
			fields: fields{
				design:   designShielded,
				quantity: 1,
			},
			damage: 50,
			want: tokenDamage{
				damage:         25,
				shipsDestroyed: 0,
			},
			wantQuantity:        1,
			wantQuantityDamaged: 1,
			wantDamage:          25,
		},
		{
			name: "not enough shields; overflows into armor",
			fields: fields{
				design:   designShielded,
				quantity: 1,
			},
			// 150 damage; 50 hits shields and remainder hits hull
			damage: 150,
			want: tokenDamage{
				damage:         100,
				shipsDestroyed: 0,
			},
			wantQuantity:        1,
			wantQuantityDamaged: 1,
			wantDamage:          100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &ShipToken{
				Quantity:        tt.fields.quantity,
				Damage:          tt.fields.damage,
				QuantityDamaged: tt.fields.quantityDamaged,
				design:          tt.fields.design,
			}
			if got := st.applyMineDamage(tt.damage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ShipToken.ApplyMineDamage() = %v, want %v", got, tt.want)
			}

			if st.Quantity != tt.wantQuantity {
				t.Errorf("ShipToken.ApplyMineDamage() token.Quantity = %v, want %v", st.Quantity, tt.wantQuantity)
			}

			if st.QuantityDamaged != tt.wantQuantityDamaged {
				t.Errorf("ShipToken.ApplyMineDamage() token.QuantityDamaged = %v, want %v", st.QuantityDamaged, tt.wantQuantityDamaged)
			}

			if st.Damage != tt.wantDamage {
				t.Errorf("ShipToken.ApplyMineDamage() token.Damage = %v, want %v", st.Damage, tt.wantDamage)
			}
		})
	}
}

func TestShipToken_applyOvergateDamage(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules))

	const lightMass = 100

	lightDesign := NewShipDesign(player.Num, 1) // 100kT ship with 100 armor
	lightDesign.Spec.Mass = lightMass
	lightDesign.Spec.Armor = 100

	type fields struct {
		Quantity        int
		Damage          float64
		QuantityDamaged int
	}
	type args struct {
		dist           float64
		safeRange      int
		safeSourceMass int
		safeDestMass   int
	}
	tests := []struct {
		name                string
		fields              fields
		args                args
		want                tokenDamage
		wantQuantity        int
		wantDamage          float64
		wantQuantityDamaged int
	}{
		{
			name:   "no damage",
			fields: fields{Quantity: 1},
			args: args{
				dist:           100,
				safeRange:      100,
				safeSourceMass: lightMass,
				safeDestMass:   lightMass,
			},
			want:                tokenDamage{},
			wantQuantity:        1,
			wantQuantityDamaged: 0,
			wantDamage:          0,
		},
		{
			name:   "double both gates' mass",
			fields: fields{Quantity: 1},
			args: args{
				dist:           100,
				safeRange:      100,
				safeSourceMass: lightMass / 2,
				safeDestMass:   lightMass / 2,
			},
			want:                tokenDamage{damage: 44}, // 1-(0.75^2)
			wantQuantity:        1,
			wantQuantityDamaged: 1,
			wantDamage:          44,
		},
		{
			name:   "2 tokens, 2x source mass only",
			fields: fields{Quantity: 2},
			args: args{
				dist:           100,
				safeRange:      100,
				safeSourceMass: lightMass / 2,
				safeDestMass:   lightMass,
			},
			// only 25% damage due to only exceeding the source gate's capabilities
			want:                tokenDamage{damage: 50},
			wantQuantity:        2,
			wantQuantityDamaged: 2,
			wantDamage:          25,
		},
		{
			name:   "range damage rounds down",
			fields: fields{Quantity: 1},
			args: args{
				dist:           305,
				safeRange:      300,
				safeSourceMass: lightMass,
				safeDestMass:   lightMass,
			},
			want:                tokenDamage{damage: 0},
			wantQuantity:        1,
			wantQuantityDamaged: 0,
			wantDamage:          0,
		},
		{
			name:   "3x safe range; half damage",
			fields: fields{Quantity: 2},
			args: args{
				dist:           300,
				safeRange:      100,
				safeSourceMass: lightMass,
				safeDestMass:   lightMass,
			},
			want:                tokenDamage{damage: 100},
			wantQuantity:        2,
			wantQuantityDamaged: 2,
			wantDamage:          50,
		},
		{
			name: "existing damage, destroy damaged token",
			fields: fields{
				Quantity:        2,
				QuantityDamaged: 1,
				Damage:          50,
			},
			args: args{
				dist:           300,
				safeRange:      100,
				safeSourceMass: lightMass,
				safeDestMass:   lightMass,
			},
			want:                tokenDamage{damage: 50, shipsDestroyed: 1}, // teeechnically should be 100 but it gets unused
			wantQuantity:        1,
			wantQuantityDamaged: 1,
			wantDamage:          50,
		},
		{
			name:   "5x safe range, damage cap",
			fields: fields{Quantity: 1},
			args: args{
				dist:           500,
				safeRange:      100,
				safeSourceMass: lightMass,
				safeDestMass:   lightMass,
			},
			want:                tokenDamage{damage: 98},
			wantQuantity:        1,
			wantQuantityDamaged: 1,
			wantDamage:          98,
		},
		{
			name:   "range & mass stack multiplicatively",
			fields: fields{Quantity: 1, Damage: 5, QuantityDamaged: 1},
			args: args{
				dist:           300,
				safeRange:      100,
				safeSourceMass: 50,
				safeDestMass:   InfiniteGate,
			},
			want:                tokenDamage{damage: 63}, // 100 * (0.5 * 0.75) = 100*0.625
			wantQuantity:        1,
			wantQuantityDamaged: 1,
			wantDamage:          68,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &ShipToken{
				Quantity:        tt.fields.Quantity,
				Damage:          tt.fields.Damage,
				QuantityDamaged: tt.fields.QuantityDamaged,
				design:          lightDesign,
			}

			if got := st.applyOvergateDamage(tt.args.dist, tt.args.safeRange, tt.args.safeSourceMass, tt.args.safeDestMass, rules.StargateMaxHullMassFactor, rules.StargateMaxRangeFactor); got != tt.want {
				t.Errorf("ShipToken.applyOvergateDamage() = %v, want %v", got, tt.want)
			}
			if st.Quantity != tt.wantQuantity {
				t.Errorf("ShipToken.applyOvergateDamage() token.Quantity = %v, want %v", st.Quantity, tt.wantQuantity)
			}
			if st.QuantityDamaged != tt.wantQuantityDamaged {
				t.Errorf("ShipToken.applyOvergateDamage() token.QuantityDamaged = %v, want %v", st.QuantityDamaged, tt.wantQuantityDamaged)
			}
			if st.Damage != tt.wantDamage {
				t.Errorf("ShipToken.applyOvergateDamage() token.Damage = %v, want %v", st.Damage, tt.wantDamage)
			}

		})
	}

	t.Run("cannot kill weak ship", func(t *testing.T) {
		// make ship with next to no armor gating very very far
		design := NewShipDesign(player.Num, 1)
		design.Spec.Mass = 500
		design.Spec.Armor = 12
		qty := 999_999
		st := &ShipToken{design: design, Quantity: qty}
		want := tokenDamage{damage: 11 * qty} // 98% of 12 gets rounded down to 11

		if got := st.applyOvergateDamage(500, 100, 100, 100, rules.StargateMaxHullMassFactor, rules.StargateMaxRangeFactor); got != want {
			t.Errorf("ShipToken.applyOvergateDamage() = %v, want %v", got, want)
		}
		if dead := qty - st.Quantity; dead > 0 {
			t.Errorf("ShipToken.applyOvergateDamage() killed %d ships unexpectedly", dead)
		}

		if numUndamaged := qty - st.QuantityDamaged; numUndamaged > 0 {
			t.Errorf("ShipToken.applyOvergateDamage() token.QuantityDamaged left %d ships undamaged", numUndamaged)
		}

	})
}

func TestShipToken_applyOvergateVanishing(t *testing.T) {
	type fields struct {
		quantity        int
		quantityDamaged int
		mass            int
	}
	type args struct {
		distance    float64
		sourceRange int
		sourceMass  int
	}
	tests := []struct {
		name                string
		fields              fields
		args                args
		rng                 rng
		wantQuantity        int
		wantQuantityDamaged int
	}{
		{
			name: "no vanish; within limits",
			fields: fields{
				quantity:        1,
				quantityDamaged: 0,
				mass:            1,
			},
			args: args{
				distance:    1,
				sourceRange: 1,
				sourceMass:  1,
			},
			rng:                 newFloat64Random(0),
			wantQuantity:        1,
			wantQuantityDamaged: 0,
		},
		{
			name: "no vanish; failed rng roll",
			fields: fields{
				quantity:        1,
				quantityDamaged: 0,
				mass:            490,
			},
			args: args{
				distance:    340,
				sourceRange: 100,
				sourceMass:  100,
			},
			// overall chance: 1-(0.67*0.8)=46.4% vanish chance
			rng:                 newFloat64Random(0.47),
			wantQuantity:        1,
			wantQuantityDamaged: 0,
		},
		{
			name: "all vanish; 20 low rolls",
			fields: fields{
				quantity:        20,
				quantityDamaged: 0,
				mass:            180, // 11% vanish chance
			},
			args: args{
				distance:    100,
				sourceRange: 100,
				sourceMass:  100,
			},
			rng:                 newFloat64Random(0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1),
			wantQuantity:        0,
			wantQuantityDamaged: 0,
		},
		{
			name: "partial vanish; damaged ships",
			fields: fields{
				quantity:        5,
				quantityDamaged: 2,
				mass:            490,
			},
			args: args{
				distance:    340,
				sourceRange: 100,
				sourceMass:  100,
			},
			// 1 vanishes; still has 1 damaged token left
			rng:                 newFloat64Random(0, 1, 1, 1, 1),
			wantQuantity:        4,
			wantQuantityDamaged: 1,
		},
		{
			name: "full vanish; damaged ships",
			fields: fields{
				quantity:        5,
				quantityDamaged: 2,
				mass:            2,
			},
			args: args{
				distance:    1,
				sourceRange: 1,
				sourceMass:  1,
			},
			rng:                 newFloat64Random(0, 0, 0, 0, 0),
			wantQuantity:        0,
			wantQuantityDamaged: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rCopy := rules
			rCopy.random = tt.rng
			d := NewShipDesign(1, 1)
			d.Spec.Mass = tt.fields.mass
			st := &ShipToken{
				Quantity:        tt.fields.quantity,
				QuantityDamaged: tt.fields.quantityDamaged,
				design:          d,
			}
			if tt.fields.quantityDamaged > 0 {
				// give token some damage if any tokens are hurt
				st.Damage = 1
			}
			st.applyOvergateVanishing(&rCopy, tt.args.distance, tt.args.sourceRange, tt.args.sourceMass)
			if st.Quantity != tt.wantQuantity {
				t.Errorf("ShipToken.applyOvergateVanishing() produced token quantity %v, want %v", st.Quantity, tt.wantQuantity)
			}

			if st.QuantityDamaged != tt.wantQuantityDamaged {
				t.Errorf("ShipToken.applyOvergateVanishing() produced token with %v damaged tokens, want %v", st.QuantityDamaged, tt.wantQuantityDamaged)
			}

			if (st.Damage == 0) != (tt.wantQuantityDamaged == 0) {
				if tt.wantQuantityDamaged == 0 {
					t.Errorf("ShipToken.applyOvergateDamage() produced tokens with %f damage; expected none", st.Damage)
				} else {
					t.Error("ShipToken.applyOvergateDamage() produced undamaged tokens; expected damage")
				}
			}
		})
	}
}

func TestShipToken_getOvergateMassVanishingChance(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules))
	design := NewShipDesign(player.Num, 1)
	tests := []struct {
		name           string
		mass           int
		safeSourceMass int
		want           float64
	}{
		{
			name:           "infinite gate",
			mass:           2,
			safeSourceMass: InfiniteGate,
			want:           0,
		},
		{
			name:           "at limit",
			mass:           200,
			safeSourceMass: 200,
			want:           0,
		},
		{
			name:           "rounding check; no vanishing",
			mass:           318,
			safeSourceMass: 300,
			want:           0,
		},
		{
			name:           "200kT ship in a 100kT gate",
			mass:           200,
			safeSourceMass: 100,
			want:           0.14, // floor(33.333*(1-9/16))/100 = floor(14.58)/100 = 14%
		},
		{
			name:           "4.9x weight limit",
			mass:           490,
			safeSourceMass: 100,
			want:           0.33, // floor(33.333*(1-1/16000))/100 = floor(33.33)/100 = 33%
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			design.Spec.Mass = tt.mass
			st := &ShipToken{design: design}
			if got := st.getOvergateMassVanishingChance(tt.safeSourceMass, rules.StargateMaxHullMassFactor); got != tt.want {
				t.Errorf("ShipToken.getOvergateMassVanishingChance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShipToken_getOvergateRangeVanishingChance(t *testing.T) {
	tests := []struct {
		name      string
		dist      float64
		safeRange int
		want      float64
	}{
		{
			name:      "within limits; no vanish",
			dist:      50,
			safeRange: 100,
			want:      0,
		},
		{
			name:      "infinite gate, no vanish",
			dist:      200_000,
			safeRange: InfiniteGate,
			want:      0,
		},
		{
			name:      "3.4x range; 20%",
			dist:      340,
			safeRange: 100,
			want:      0.2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(1, NewRace().WithSpec(&rules))
			tr := &ShipToken{design: NewShipDesign(player.Num, 1)}

			if got := tr.getOvergateRangeVanishingChance(tt.dist, tt.safeRange, rules.StargateMaxRangeFactor); got != tt.want {
				t.Errorf("ShipToken.getOvergateRangeVanishingChance() = %v, want %v", got, tt.want)
			}
		})
	}
}
