package cs

import (
	"reflect"
	"testing"
)

func TestShipToken_applyMineDamage(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules))
	design := NewShipDesign(player, 1)

	// set some spec values we care about
	design.Spec.Mass = 100
	design.Spec.Armor = 100

	designShielded := NewShipDesign(player, 1)
	designShielded.Spec.Mass = 100
	designShielded.Spec.Armor = 150
	designShielded.Spec.Shields = 50

	type fields struct {
		Quantity        int
		Damage          float64
		QuantityDamaged int
		design          *ShipDesign
	}
	type args struct {
		damage int
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
			name: "1 ship, do 50 damage, don't destroy ship",
			fields: fields{
				design:   design,
				Quantity: 1,
			},
			args: args{
				damage: 50,
			},
			want: tokenDamage{
				damage:         50,
				shipsDestroyed: 0,
			},
			wantQuantity:        1,
			wantQuantityDamaged: 1,
			wantDamage:          50,
		},
		{
			name: "10 ships, do 850 damage, destroy 8 ships, leave 2 damaged",
			fields: fields{
				design:   design,
				Quantity: 10,
			},
			args: args{
				damage: 850,
			},
			want: tokenDamage{
				damage:         850,
				shipsDestroyed: 8,
			},
			wantQuantity:        2,
			wantQuantityDamaged: 2,
			wantDamage:          25,
		},
		{
			name: "2 ships, with 50 damage already, do 75 more damage, destroy ship",
			fields: fields{
				design:          design,
				Quantity:        2,
				QuantityDamaged: 1,
				Damage:          50,
			},
			args: args{
				damage: 75,
			},
			want: tokenDamage{
				damage:         75,
				shipsDestroyed: 1,
			},
			wantQuantity:        1,
			wantQuantityDamaged: 1,
			wantDamage:          25,
		},
		{
			name: "1 shielded ship, do 50 damage, don't destroy ship",
			fields: fields{
				design:   designShielded,
				Quantity: 1,
			},
			args: args{
				damage: 50,
			},
			want: tokenDamage{
				damage:         25,
				shipsDestroyed: 0,
			},
			wantQuantity:        1,
			wantQuantityDamaged: 1,
			wantDamage:          25,
		},
		{
			name: "take 150 mine damage, our shields absorb 50 and our token  takes the rest",
			fields: fields{
				design:   designShielded,
				Quantity: 1,
			},
			args: args{
				damage: 150,
			},
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
				Quantity:        tt.fields.Quantity,
				Damage:          tt.fields.Damage,
				QuantityDamaged: tt.fields.QuantityDamaged,
				design:          tt.fields.design,
			}
			if got := st.applyMineDamage(tt.args.damage); !reflect.DeepEqual(got, tt.want) {
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
	design := NewShipDesign(player, 1)
	heavyDesign := NewShipDesign(player, 1)

	// 100kT ship with 100 armor
	mass := 100
	design.Spec.Mass = mass
	design.Spec.Armor = 100

	// 200kT ship with 100 armor
	heavyDesign.Spec.Mass = 200
	heavyDesign.Spec.Armor = 100

	type fields struct {
		Quantity        int
		Damage          float64
		QuantityDamaged int
		design          *ShipDesign
	}
	type args struct {
		dist           float64
		safeRange      int
		safeSourceMass int
		safeDestMass   int
		maxMassFactor  int
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
			name: "no damage",
			fields: fields{
				design:   design,
				Quantity: 1,
			},
			args: args{
				dist:           100,
				safeRange:      100,
				safeSourceMass: mass,
				safeDestMass:   mass,
				maxMassFactor:  5,
			},
			want:                tokenDamage{},
			wantQuantity:        1,
			wantQuantityDamaged: 0,
			wantDamage:          0,
		},
		{
			name: "one token overmass",
			fields: fields{
				design:   design,
				Quantity: 1,
			},
			args: args{
				dist:           100,
				safeRange:      100,
				safeSourceMass: mass / 2,
				safeDestMass:   mass / 2,
				maxMassFactor:  5,
			},
			want:                tokenDamage{damage: 44},
			wantQuantity:        1,
			wantQuantityDamaged: 1,
			wantDamage:          44,
		},
		{
			name: "two token overmass",
			fields: fields{
				design:   design,
				Quantity: 2,
			},
			args: args{
				dist:           100,
				safeRange:      100,
				safeSourceMass: mass / 2,
				safeDestMass:   mass / 2,
				maxMassFactor:  5,
			},
			want:                tokenDamage{damage: 44},
			wantQuantity:        2,
			wantQuantityDamaged: 2,
			wantDamage:          44,
		},
		{
			name: "going over range by a little will not cause damage due to rounding",
			fields: fields{
				design:   design,
				Quantity: 1,
			},
			args: args{
				dist:           305,
				safeRange:      300,
				safeSourceMass: mass,
				safeDestMass:   mass,
				maxMassFactor:  5,
			},
			want:                tokenDamage{damage: 0},
			wantQuantity:        1,
			wantQuantityDamaged: 0,
			wantDamage:          0,
		},
		{
			name: "going over range by 2x should give 50 total damage, damaging both ship",
			fields: fields{
				design:   design,
				Quantity: 2,
			},
			args: args{
				dist:           300,
				safeRange:      100,
				safeSourceMass: mass,
				safeDestMass:   mass,
				maxMassFactor:  5,
			},
			want:                tokenDamage{damage: 50},
			wantQuantity:        2,
			wantQuantityDamaged: 2,
			wantDamage:          50,
		},
		{
			name: "existing damage, destroy damaged token",
			fields: fields{
				design:          design,
				Quantity:        2,
				QuantityDamaged: 1,
				Damage:          50,
			},
			args: args{
				dist:           300,
				safeRange:      100,
				safeSourceMass: mass,
				safeDestMass:   mass,
				maxMassFactor:  5,
			},
			want:                tokenDamage{damage: 50, shipsDestroyed: 1},
			wantQuantity:        1,
			wantQuantityDamaged: 1,
			wantDamage:          50,
		},
		{
			name: "100% damage is 4x over safe range (maxes out at 1 armor point left)",
			fields: fields{
				design:   design,
				Quantity: 1,
			},
			args: args{
				dist:           500,
				safeRange:      100,
				safeSourceMass: mass,
				safeDestMass:   mass,
				maxMassFactor:  5,
			},
			want:                tokenDamage{damage: 98},
			wantQuantity:        1,
			wantQuantityDamaged: 1,
			wantDamage:          98,
		},
		{
			// i.e. sending a 200kT ship through a 100kT gate source gate with infinite dest gate
			name: "25% damage for doubling allowed mass",
			fields: fields{
				design:   heavyDesign,
				Quantity: 1,
			},
			args: args{
				dist:           100,
				safeRange:      100,
				safeSourceMass: 100,
				safeDestMass:   InfiniteGate,
				maxMassFactor:  5,
			},
			want:                tokenDamage{damage: 25},
			wantQuantity:        1,
			wantQuantityDamaged: 1,
			wantDamage:          25,
		},
		{
			// // i.e. sending a 200kT ship through a 100kT gate source and dest gate
			name: "25% damage on each side for sending a ship through two gates with double mass limits",
			fields: fields{
				design:   heavyDesign,
				Quantity: 1,
			},
			args: args{
				dist:           100,
				safeRange:      100,
				safeSourceMass: 100,
				safeDestMass:   100,
				maxMassFactor:  5,
			},
			want:                tokenDamage{damage: 44}, // armor * (1 - .75 * .75)
			wantQuantity:        1,
			wantQuantityDamaged: 1,
			wantDamage:          44,
		},
		{
			// i.e. sending a 200kT ship through a 100kT gate source gate with infinite dest gate
			name: "25% damage for doubling allowed mass, 50% damage for range",
			fields: fields{
				design:   heavyDesign,
				Quantity: 1,
			},
			args: args{
				dist:           200,
				safeRange:      100,
				safeSourceMass: 100,
				safeDestMass:   InfiniteGate,
				maxMassFactor:  5,
			},
			want:                tokenDamage{damage: 44}, // armor * (1 - .75 * .75)
			wantQuantity:        1,
			wantQuantityDamaged: 1,
			wantDamage:          44,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &ShipToken{
				Quantity:        tt.fields.Quantity,
				Damage:          tt.fields.Damage,
				QuantityDamaged: tt.fields.QuantityDamaged,
				design:          tt.fields.design,
			}
			if got := st.applyOvergateDamage(tt.args.dist, tt.args.safeRange, tt.args.safeSourceMass, tt.args.safeDestMass, tt.args.maxMassFactor); !reflect.DeepEqual(got, tt.want) {
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
}

func TestShipToken_getStargateMassVanishingChance(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules))

	type fields struct {
		mass int
	}
	type args struct {
		safeSourceMass int
		maxMassFactor  int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name:   "no vanishing chance",
			fields: fields{mass: 200},
			args:   args{safeSourceMass: 200, maxMassFactor: 5},
			want:   0,
		},
		{
			name:   "no vanishing chance for 318kT on a 300/500 gate, due to rounding",
			fields: fields{mass: 318},
			args:   args{safeSourceMass: 300, maxMassFactor: 5},
			want:   0,
		},
		{
			name:   "200kT ship in a 100kt gate has a 14% chance of vanishing",
			fields: fields{mass: 600},
			args:   args{safeSourceMass: 300, maxMassFactor: 5},
			want:   .14,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			design := NewShipDesign(player, 1)
			design.Spec.Mass = tt.fields.mass

			tr := &ShipToken{
				Quantity: 1,
				design:   design,
			}
			if got := tr.getStargateMassVanishingChance(tt.args.safeSourceMass, tt.args.maxMassFactor); got != tt.want {
				t.Errorf("ShipToken.getStargateMassVanishingChance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShipToken_getStargateRangeVanishingChance(t *testing.T) {
	type args struct {
		dist      float64
		safeRange int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "no vanishing chance",
			args: args{
				dist:      100,
				safeRange: 100,
			},
			want: 0,
		},
		{
			name: "20% vanishing chance for 3.4x range",
			args: args{
				dist:      340,
				safeRange: 100,
			},
			want: .2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(1, NewRace().WithSpec(&rules))
			design := NewShipDesign(player, 1)
			tr := &ShipToken{
				Quantity: 1,
				design:   design,
			}

			if got := tr.getStargateRangeVanishingChance(tt.args.dist, tt.args.safeRange); got != tt.want {
				t.Errorf("ShipToken.getStargateRangeVanishingChance() = %v, want %v", got, tt.want)
			}
		})
	}
}
