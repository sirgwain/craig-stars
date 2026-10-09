//go:build !wasi && !wasm

package cs

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
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
			name: "1 ship, do 50 damage, don't destroy ship",
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
			name: "10 ships, do 850 damage, spread over all of them",
			fields: fields{
				design:   design,
				quantity: 10,
			},
			damage: 850,
			want: tokenDamage{
				damage:         850,
				shipsDestroyed: 0,
			},
			wantQuantity:        10,
			wantQuantityDamaged: 10,
			wantDamage:          85,
		},
		{
			name: "10 ships, do more damage than their armor, destroy all of them",
			fields: fields{
				design:   design,
				quantity: 10,
			},
			damage: 1010,
			want: tokenDamage{
				damage:         1010,
				shipsDestroyed: 10,
			},
		},
		{
			name: "2 ships, with 50 damage already, do 75 more damage, spread over both",
			fields: fields{
				design:          design,
				quantity:        2,
				quantityDamaged: 1,
				damage:          50,
			},
			damage: 75,
			want: tokenDamage{
				damage:         75,
				shipsDestroyed: 0,
			},
			wantQuantity:        2,
			wantQuantityDamaged: 2,
			wantDamage:          62,
		},
		{
			name: "2 ships, with 50 damage already, do 160 more damage, destroy both",
			fields: fields{
				design:          design,
				quantity:        2,
				quantityDamaged: 1,
				damage:          50,
			},
			damage: 160,
			want: tokenDamage{
				damage:         160,
				shipsDestroyed: 2,
			},
		},
		{
			name: "1 shielded ship, do 50 damage, don't destroy ship",
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
			name: "take 150 mine damage, our shields absorb 50 and our token takes the rest",
			fields: fields{
				design:   designShielded,
				quantity: 1,
			},
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

func TestShipToken_overgateDamagePercent(t *testing.T) {
	gate := func(safeRange, safeMass int) PlanetStarbaseSpec {
		return PlanetStarbaseSpec{SafeRange: safeRange, SafeHullMass: safeMass}
	}
	tests := []struct {
		name   string
		mass   int
		dist   float64
		source PlanetStarbaseSpec
		dest   PlanetStarbaseSpec
		want   int
	}{
		{"within limits", 100, 100, gate(100, 100), gate(100, 100), 0},
		{"a little over range rounds down to no damage", 100, 305, gate(300, 100), gate(300, 100), 0},
		{"2x over range", 100, 300, gate(100, 100), gate(100, 100), 50},
		{"at the max range", 100, 500, gate(100, 100), gate(100, 100), 100},
		{"only the source gate's range counts", 100, 300, gate(InfiniteGate, 100), gate(100, 100), 0},
		{"double the source gate's mass", 200, 100, gate(100, 100), gate(100, InfiniteGate), 25},
		{"double both gates' mass", 200, 100, gate(100, 100), gate(100, 100), 43},               // 1 - .75 * .75
		{"double the mass and 2x range", 200, 200, gate(100, 100), gate(100, InfiniteGate), 43}, // 1 - .75 * .75
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			design := NewShipDesign(1, 1)
			design.Spec.Mass = tt.mass
			st := &ShipToken{Quantity: 1, design: design}
			assert.Equal(t, tt.want, st.overgateDamagePercent(&rules, tt.dist, tt.source, tt.dest))
		})
	}
}

func TestShipToken_applyOvergateDamage(t *testing.T) {
	design := NewShipDesign(1, 1)
	design.Spec.Armor = 100

	tests := []struct {
		name                string
		quantity            int
		quantityDamaged     int
		damage              float64
		damagePercent       int
		want                tokenDamage
		wantQuantity        int
		wantQuantityDamaged int
		wantDamage          float64
	}{
		{name: "no damage", quantity: 1, damagePercent: 0, wantQuantity: 1},
		{name: "damage every ship", quantity: 2, damagePercent: 50, want: tokenDamage{damage: 50}, wantQuantity: 2, wantQuantityDamaged: 2, wantDamage: 50},
		{name: "never destroys an undamaged ship", quantity: 1, damagePercent: 99, want: tokenDamage{damage: 99}, wantQuantity: 1, wantQuantityDamaged: 1, wantDamage: 99},
		{name: "100% destroys every ship", quantity: 2, damagePercent: 100, want: tokenDamage{damage: 100, shipsDestroyed: 2}},
		{
			name: "damaged ships that can't take it are destroyed", quantity: 2, quantityDamaged: 1, damage: 50, damagePercent: 50,
			want: tokenDamage{damage: 50, shipsDestroyed: 1}, wantQuantity: 1, wantQuantityDamaged: 1, wantDamage: 50,
		},
		{
			// 30 + 30 + 20 spread over both ships
			name: "existing damage stacks", quantity: 2, quantityDamaged: 1, damage: 20, damagePercent: 30,
			want: tokenDamage{damage: 30}, wantQuantity: 2, wantQuantityDamaged: 2, wantDamage: 40,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &ShipToken{Quantity: tt.quantity, QuantityDamaged: tt.quantityDamaged, Damage: tt.damage, design: design}
			assert.Equal(t, tt.want, st.applyOvergateDamage(tt.damagePercent))
			assert.Equal(t, tt.wantQuantity, st.Quantity)
			assert.Equal(t, tt.wantQuantityDamaged, st.QuantityDamaged)
			assert.Equal(t, tt.wantDamage, st.Damage)
		})
	}

	t.Run("overgating again and again", func(t *testing.T) {
		// overgating 12 scouts 479.5 ly with a 250ly gate damages them 20% (4 of 20 armor),
		// then 40%, 60%, 80%, then destroys them
		scout := NewShipDesign(1, 1)
		scout.Spec.Armor = 20
		scout.Spec.Mass = 25
		st := &ShipToken{Quantity: 12, design: scout}
		damagePercent := st.overgateDamagePercent(&rules, 479.5, PlanetStarbaseSpec{SafeRange: 250, SafeHullMass: 300}, PlanetStarbaseSpec{SafeRange: 250, SafeHullMass: 300})
		for _, wantDamage := range []float64{4, 8, 12, 16} {
			st.applyOvergateDamage(damagePercent)
			assert.Equal(t, 12, st.Quantity)
			assert.Equal(t, wantDamage, st.Damage)
		}
		assert.Equal(t, tokenDamage{damage: 4, shipsDestroyed: 12}, st.applyOvergateDamage(damagePercent))
		assert.Equal(t, 0, st.Quantity)
	})
}

func TestShipToken_applyOvergateVanishing(t *testing.T) {
	tests := []struct {
		name                string
		quantity            int
		quantityDamaged     int
		damagePercent       int
		rng                 rng
		wantQuantity        int
		wantQuantityDamaged int
	}{
		{name: "no damage, no vanishing", quantity: 1, damagePercent: 0, rng: newFloat64Random(0), wantQuantity: 1},
		{name: "60% damage is a 20% chance, high roll", quantity: 1, damagePercent: 60, rng: newFloat64Random(0.21), wantQuantity: 1},
		{name: "60% damage is a 20% chance, low roll", quantity: 1, damagePercent: 60, rng: newFloat64Random(0.2), wantQuantity: 0},
		{name: "lose damaged ships first", quantity: 3, quantityDamaged: 1, damagePercent: 60, rng: newFloat64Random(0.9, 0.9, 0.1), wantQuantity: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules := NewRulesWithSeed(0)
			rules.random = tt.rng
			st := &ShipToken{Quantity: tt.quantity, QuantityDamaged: tt.quantityDamaged, Damage: 10, design: NewShipDesign(1, 1)}
			st.applyOvergateVanishing(&rules, tt.damagePercent)
			assert.Equal(t, tt.wantQuantity, st.Quantity)
			assert.Equal(t, tt.wantQuantityDamaged, st.QuantityDamaged)
		})
	}
}
