//go:build !wasi && !wasm

package cs

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

// Test_battleDamageScore verifies damage scoring for each battle tactic.
func Test_battleDamageScore(t *testing.T) {
	tests := []struct {
		name               string
		tactic             BattleTactic
		given, taken, want float64
	}{
		{"maximize damage", BattleTacticMaximizeDamage, 100, 9, -100},
		{"unchallenged pursuit", BattleTacticDisengageIfChallenged, 100, 9, -100},
		{"net damage uses ratio", BattleTacticMaximizeNetDamage, 100, 9, -10},
		{"damage ratio", BattleTacticMaximizeDamageRatio, 100, 9, -10},
		{"no outgoing damage minimizes incoming", BattleTacticMaximizeDamageRatio, 0, 9, 9},
		{"minimize incoming", BattleTacticMinimizeDamageToSelf, 100, 9, 9},
		{"retreat", BattleTacticDisengage, 100, 9, 9},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := battleDamageScore(tt.given, tt.taken, tt.tactic)
			assert.True(t, test.WithinTolerance(got, tt.want, .001), "got %v, want %v", got, tt.want)
		})
	}
}

// Test_battle_movementProjection verifies how enemy movement and tactics affect predicted damage.
func Test_battle_movementProjection(t *testing.T) {
	tests := []struct {
		name        string
		enemyMoves  int
		enemyTactic BattleTactic
		want        float64
	}{
		{"stationary enemy remains at current distance", 0, BattleTacticMaximizeDamage, -9},
		{"enemy with a move can close range", 1, BattleTacticMaximizeDamage, -10},
		{"enemy avoiding damage can leave range", 1, BattleTacticMinimizeDamageToSelf, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := testBattleToken(1, Vector{}, BattleTacticMaximizeDamage, 2).withTestBeam(10, 1)
			source.movesLeft = 1
			enemy := testBattleToken(2, Vector{1, 0}, tt.enemyTactic, 1).withTestBeam(20, 1)
			enemy.movesLeft = tt.enemyMoves
			b := &battle{rules: &rules, tokens: []*battleToken{source, enemy}}
			got := b.movementScore(source, Vector{}, BattleTargetAny)
			assert.True(t, test.WithinTolerance(got, tt.want, .001), "got %v, want %v", got, tt.want)
		})
	}
}

// Test_battle_movementSearch verifies planning distance and pursuit of targets beyond reach.
func Test_battle_movementSearch(t *testing.T) {
	tests := []struct {
		name        string
		enemyMoves  int
		tactic      BattleTactic
		wantRadius  int
		wantClosest bool
	}{
		{"plan over three remaining moves", 0, BattleTacticMaximizeDamage, 3, false},
		{"moving enemy remains beyond reach", 3, BattleTacticMaximizeDamage, 1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := testBattleToken(1, Vector{1, 4}, tt.tactic, 2).withTestBeam(10, 1)
			source.movesLeft = 3
			enemy := testBattleToken(2, Vector{5, 4}, BattleTacticMaximizeDamage)
			enemy.movesLeft = tt.enemyMoves
			b := &battle{rules: &rules, tokens: []*battleToken{source, enemy}}
			radius, closest := b.movementSearch(source, BattleTargetAny)
			assert.Equal(t, tt.wantRadius, radius)
			assert.Equal(t, tt.wantClosest, closest == enemy)
		})
	}
}

func Test_updateBestMoves_center(t *testing.T) {
	type args struct {
		better      bool
		newPosition Vector
		bestMoves   []Vector
	}
	tests := []struct {
		name string
		args args
		want []Vector
	}{
		{
			name: "better move 1,0, pick it",
			args: args{better: true, newPosition: Vector{1, 0}, bestMoves: []Vector{{0, 0}}},
			want: []Vector{{1, 0}},
		},
		{
			name: "better move away from center, pick it",
			args: args{better: true, newPosition: Vector{3, 3}, bestMoves: []Vector{{4, 4}, {4, 5}}},
			want: []Vector{{3, 3}},
		},
		{
			name: "equivalent damage move, but newPosition is closer to center",
			args: args{better: false, newPosition: Vector{4, 4}, bestMoves: []Vector{{4, 3}}},
			want: []Vector{{4, 4}},
		},
		{
			name: "equivalent damage move, newPosition is closer to center",
			args: args{better: false, newPosition: Vector{4, 5}, bestMoves: []Vector{{4, 3}}},
			want: []Vector{{4, 5}},
		},
		{
			name: "equivalent damage move, newPosition is same distance to center",
			args: args{better: false, newPosition: Vector{4, 5}, bestMoves: []Vector{{4, 4}}},
			want: []Vector{{4, 4}, {4, 5}},
		},
		{
			name: "equivalent damage move, newPosition is same distance to center",
			args: args{better: false, newPosition: Vector{5, 5}, bestMoves: []Vector{{4, 4}, {4, 5}}},
			want: []Vector{{4, 4}, {4, 5}, {5, 5}},
		},
		{
			name: "equivalent damage move, newPosition is farther from center, discard it",
			args: args{better: false, newPosition: Vector{6, 5}, bestMoves: []Vector{{4, 4}, {4, 5}}},
			want: []Vector{{4, 4}, {4, 5}},
		},
		{
			name: "equivalent damage move, middle of an edge is closer to center than a corner",
			args: args{better: false, newPosition: Vector{0, 4}, bestMoves: []Vector{{0, 0}}},
			want: []Vector{{0, 4}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := updateBestMoves(tt.args.better, tt.args.newPosition, tt.args.bestMoves, battleCenterDistance); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("updateBestMoves() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_updateBestMoves_nearest verifies that retreating tokens break ties toward their current square.
func Test_updateBestMoves_nearest(t *testing.T) {
	current := Vector{4, 4}
	distance := func(position Vector) int { return current.chebyshevDistance(position) }
	tests := []struct {
		name        string
		newPosition Vector
		bestMoves   []Vector
		want        []Vector
	}{
		{"nearer square replaces farther ties", Vector{5, 5}, []Vector{{6, 6}}, []Vector{{5, 5}}},
		{"farther square is discarded", Vector{6, 6}, []Vector{{5, 5}}, []Vector{{5, 5}}},
		{"same distance is kept even if it is toward the center", Vector{4, 5}, []Vector{{3, 3}}, []Vector{{3, 3}, {4, 5}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, updateBestMoves(false, tt.newPosition, tt.bestMoves, distance))
		})
	}
}

func Test_battle_getBestMoves_flee(t *testing.T) {
	// generate a fleeing token at a position with 10dp
	fleeingToken := func(position Vector) *battleToken {
		token := testBattleToken(1, position, BattleTacticDisengage, 2)
		token.armor = 10
		return token
	}

	// generate an enemy token at a position with a single laser
	enemyToken := func(position Vector) *battleToken {
		return testBattleToken(2, position, BattleTacticMaximizeDamage, 1).withTestBeam(10, 1)
	}

	tests := []struct {
		name    string
		token   *battleToken
		enemies []*battleToken
		want    []Vector
	}{
		{
			name:  "token at 1,4 move randomly 1st option",
			token: fleeingToken(Vector{1, 4}),
			// Equally safe moves tie, so pick any neighboring square rather than staying put.
			want: []Vector{{0, 3}, {0, 4}, {0, 5}, {1, 3}, {1, 5}, {2, 3}, {2, 4}, {2, 5}},
		},
		{
			name:  "token at 1,4 move away from surrounding weapons",
			token: fleeingToken(Vector{1, 4}),
			// make three weapons adjacent so we have to move straight back
			enemies: []*battleToken{
				enemyToken(Vector{1, 5}),
				enemyToken(Vector{2, 4}),
				enemyToken(Vector{1, 3}),
			},
			want: []Vector{{0, 3}, {0, 5}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &battle{
				tokens: append([]*battleToken{tt.token}, tt.enemies...),
				rules:  &rules,
			}

			// run away and record the new position
			if got := b.getBestMoves(tt.token); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("battle.getBestMoves() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_battle_getBestAttackMoves(t *testing.T) {
	type args struct {
		token   *battleToken
		enemies []*battleToken
	}

	laser := battleWeaponSlot{
		weaponType:   battleWeaponTypeBeam,
		power:        10,
		beamBonus:    1,
		weaponRange:  1,
		slotQuantity: 1,
	}

	// generate a fleeing token at a position with 10dp
	attackingToken := func(position Vector, tactic BattleTactic, weapon battleWeaponSlot) *battleToken {
		player := testPlayer().WithNum(1)
		token := battleToken{
			BattleRecordToken: BattleRecordToken{Position: position, PlayerNum: player.Num, Tactic: tactic, PrimaryTarget: BattleTargetAny, Movement: 4},
			attackPlayers:     map[int]bool{2: true},
			ShipToken:         &ShipToken{Quantity: 1},
			player:            player,
			attributes:        battleTokenAttributeArmed,
			armor:             100,
			cost:              Cost{1, 1, 1, 1},
		}

		// put our token in the weapon passed in
		weapon.token = &token
		token.weaponSlots = append(token.weaponSlots, &weapon)

		return &token
	}

	// generate an enemy weapon at a position with a single laser
	enemyToken := func(position Vector, weapon *battleWeaponSlot) *battleToken {
		player := testPlayer().WithNum(2)
		token := battleToken{
			BattleRecordToken: BattleRecordToken{Position: position, PlayerNum: player.Num, Tactic: BattleTacticMaximizeDamage, PrimaryTarget: BattleTargetAny, Movement: 4},
			attackPlayers:     map[int]bool{1: true},
			ShipToken:         &ShipToken{Quantity: 1},
			player:            testPlayer().WithNum(2),
			attributes:        battleTokenAttributeArmed,
			armor:             100,
			cost:              Cost{1, 1, 1, 1},
		}

		// if this enemy has a weapon, assign it now
		if weapon != nil {
			tokenWeapon := *weapon
			tokenWeapon.token = &token
			token.weaponSlots = append(token.weaponSlots, &tokenWeapon)
		}

		return &token
	}

	tests := []struct {
		name string
		args args
		want []Vector
	}{
		{
			// attacker should move over, or over and up/down
			// * * * *
			// A * * T
			// * * * *
			name: "move towards enemy",
			args: args{
				token: attackingToken(Vector{0, 1}, BattleTacticMaximizeDamage, laser),
				// make three weapons adjacent so we have to move straight back
				enemies: []*battleToken{
					enemyToken(Vector{4, 1}, nil),
				},
			},
			want: []Vector{{1, 2}},
		},
		{
			// attacker should move over, or over and down
			// A * * *
			// * * * *
			// * * * T
			name: "move towards enemy right or right/down",
			args: args{
				token: attackingToken(Vector{0, 0}, BattleTacticMaximizeDamage, laser),
				// make three weapons adjacent so we have to move straight back
				enemies: []*battleToken{
					enemyToken(Vector{3, 2}, nil),
				},
			},
			want: []Vector{{1, 1}},
		},
		{
			// attacker should move on top to maximize damage
			// A T * *
			// * * * *
			// * * * *
			name: "maximize beam damage one target",
			args: args{
				token: attackingToken(Vector{0, 0}, BattleTacticMaximizeDamage, laser),
				// make three weapons adjacent so we have to move straight back
				enemies: []*battleToken{
					enemyToken(Vector{1, 0}, nil),
				},
			},
			want: []Vector{{1, 0}},
		},
		{
			// attacker should move to cause the most damage vs damage taken
			// the board will have two tokens, a strong and a weak one
			// * S * *
			// A W * *
			// * * * *
			name: "maximize damage ratio strong and weak target",
			args: args{
				token: attackingToken(Vector{0, 1}, BattleTacticMaximizeDamageRatio, laser),
				// make three weapons adjacent so we have to move straight back
				enemies: []*battleToken{
					// strong 100 power beamer
					enemyToken(Vector{1, 0}, &battleWeaponSlot{weaponType: battleWeaponTypeBeam, power: 100, beamBonus: 1, slotQuantity: 1, weaponRange: 1}),
					// weak 10 power beamer
					enemyToken(Vector{1, 1}, &laser),
				},
			},
			want: []Vector{{1, 2}}, // best damage ratio, and towards center
		},
		{
			// attacker has torpedoes and wants to stay out of range of those lasers
			// it should move back
			// * * 1 *
			// * A * *
			// * * 2 *
			name: "maximize damage ratio, prefer no damage",
			args: args{
				token: attackingToken(Vector{1, 1}, BattleTacticMaximizeDamageRatio, battleWeaponSlot{weaponType: battleWeaponTypeTorpedo, power: 10, accuracy: 1, slotQuantity: 3, weaponRange: 2}),
				// make three weapons adjacent so we have to move straight back
				enemies: []*battleToken{
					enemyToken(Vector{2, 0}, &laser),
					// put two tokens here
					enemyToken(Vector{2, 2}, &laser),
					enemyToken(Vector{2, 2}, &laser),
				},
			},
			want: []Vector{{0, 2}},
		},
		{
			// attacker has torpedoes and wants to stay out of range of those lasers
			// it should move back but stay near the center if possible
			// * * 1 *
			// * A * *
			// * * 2 *
			name: "maximize damage ratio, prefer no damage, start in center (5,5)",
			args: args{
				token: attackingToken(Vector{5, 5}, BattleTacticMaximizeDamageRatio, battleWeaponSlot{weaponType: battleWeaponTypeTorpedo, power: 10, accuracy: 1, slotQuantity: 3, weaponRange: 2}),
				// make three weapons adjacent so we have to move straight back
				enemies: []*battleToken{
					enemyToken(Vector{6, 4}, &laser),
					// put two tokens here
					enemyToken(Vector{6, 6}, &laser),
					enemyToken(Vector{6, 6}, &laser),
				},
			},
			want: []Vector{{4, 4}, {4, 5}}, // we pick the 0 damageTaken options that keep us near center
		},
		{
			// attacker wants to cause the largest difference in damage
			// the board will have three tokens, one up and right, and 2 down and right (one unarmed)
			// we will have a powerful beam that can burn through multiple tokens so we'll move to the
			// 2 token space to attack them both and only take damage from one
			// * 1 * *
			// A * * *
			// * 2 * *
			name: "maximize net damage",
			args: args{
				token: attackingToken(Vector{0, 1}, BattleTacticMaximizeNetDamage, battleWeaponSlot{weaponType: battleWeaponTypeBeam, power: 20, beamBonus: 1, slotQuantity: 1, weaponRange: 1}),
				// make three weapons adjacent so we have to move straight back
				enemies: []*battleToken{
					enemyToken(Vector{1, 0}, &laser),
					// tokens at this spot
					enemyToken(Vector{1, 2}, nil),
					enemyToken(Vector{1, 2}, &laser),
				},
			},
			want: []Vector{{1, 2}}, // best net damage, move to two token square
		},
		{
			// attacker has torpedoes and wants to stay out of range of those lasers
			// it should move back
			// * * L *
			// * A * *
			// * * L *
			name: "minimize damage to self",
			args: args{
				token: attackingToken(Vector{1, 1}, BattleTacticMinimizeDamageToSelf, battleWeaponSlot{weaponType: battleWeaponTypeTorpedo, power: 10, accuracy: 1, slotQuantity: 2, weaponRange: 2}),
				// make three weapons adjacent so we have to move straight back
				enemies: []*battleToken{
					enemyToken(Vector{2, 0}, &laser),
					enemyToken(Vector{2, 2}, &laser),
				},
			},
			want: []Vector{{0, 2}}, // min damage, move out of range, away from the corner
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &battle{
				tokens: append([]*battleToken{tt.args.token}, tt.args.enemies...),
				rules:  &rules,
			}

			if got := b.getBestMoves(tt.args.token); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("battle.getBestMove() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_battleStepsToward verifies the single squares a token can step to toward a destination.
func Test_battleStepsToward(t *testing.T) {
	tests := []struct {
		name                  string
		position, destination Vector
		want                  []Vector
	}{
		{"adjacent destination is the only step", Vector{4, 4}, Vector{5, 5}, []Vector{{5, 5}}},
		{"staying put", Vector{4, 4}, Vector{4, 4}, []Vector{{4, 4}}},
		{"straight along x can also step diagonally", Vector{4, 4}, Vector{7, 4}, []Vector{{5, 4}, {5, 3}, {5, 5}}},
		{"straight along y can also step diagonally", Vector{4, 4}, Vector{4, 1}, []Vector{{4, 3}, {3, 3}, {5, 3}}},
		{"mostly x steps diagonally or along x", Vector{4, 4}, Vector{8, 6}, []Vector{{5, 5}, {5, 4}}},
		{"mostly y steps diagonally or along y", Vector{4, 4}, Vector{2, 8}, []Vector{{3, 5}, {4, 5}}},
		{"exact diagonal steps diagonally", Vector{4, 4}, Vector{7, 7}, []Vector{{5, 5}}},
		{"steps off the board are dropped", Vector{0, 0}, Vector{0, 5}, []Vector{{0, 1}, {1, 1}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, battleStepsToward(tt.position, tt.destination))
		})
	}
}
