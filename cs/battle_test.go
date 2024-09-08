package cs

import (
	"reflect"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func testStalwartDefender(player *Player) *Fleet {
	fleet := &Fleet{
		MapObject: MapObject{
			PlayerNum: player.Num,
		},
		BaseName: "Stalwart Defender",
		Tokens: []ShipToken{
			{
				DesignNum: 1,
				Quantity:  1,
				design: NewShipDesign(player, 1).
					WithHull(Destroyer.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: BetaTorpedo.Name, HullSlotIndex: 2, Quantity: 1},
						{HullComponent: XRayLaser.Name, HullSlotIndex: 3, Quantity: 1},
						{HullComponent: RhinoScanner.Name, HullSlotIndex: 4, Quantity: 1},
						{HullComponent: Crobmnium.Name, HullSlotIndex: 5, Quantity: 1},
						{HullComponent: Overthruster.Name, HullSlotIndex: 6, Quantity: 1},
						{HullComponent: BattleComputer.Name, HullSlotIndex: 7, Quantity: 1},
					}).
					WithSpec(&rules, player)},
		},
		battlePlan:        &player.BattlePlans[0],
		OrbitingPlanetNum: None,
		FleetOrders: FleetOrders{
			Waypoints: []Waypoint{
				NewPositionWaypoint(Vector{}, 5),
			},
		},
	}
	fleet.Spec = ComputeFleetSpec(&rules, player, fleet)
	fleet.Fuel = fleet.Spec.FuelCapacity
	return fleet
}

func testJihadCruiser(player *Player) *Fleet {
	fleet := &Fleet{
		MapObject: MapObject{
			PlayerNum: player.Num,
		},
		BaseName: "Jihad Cruiser",
		Tokens: []ShipToken{
			{
				DesignNum: 1,
				Quantity:  1,
				design: NewShipDesign(player, 1).
					WithHull(Cruiser.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: TransStar10.Name, HullSlotIndex: 1, Quantity: 2},
						{HullComponent: Overthruster.Name, HullSlotIndex: 2, Quantity: 1},
						{HullComponent: BattleNexus.Name, HullSlotIndex: 3, Quantity: 1},
						{HullComponent: JihadMissile.Name, HullSlotIndex: 4, Quantity: 2},
						{HullComponent: JihadMissile.Name, HullSlotIndex: 5, Quantity: 2},
						{HullComponent: ElephantScanner.Name, HullSlotIndex: 6, Quantity: 2},
						{HullComponent: Kelarium.Name, HullSlotIndex: 5, Quantity: 2},
					}).
					WithSpec(&rules, player)},
		},
		battlePlan:        &player.BattlePlans[0],
		OrbitingPlanetNum: None,
		FleetOrders: FleetOrders{
			Waypoints: []Waypoint{
				NewPositionWaypoint(Vector{}, 5),
			},
		},
	}
	fleet.Spec = ComputeFleetSpec(&rules, player, fleet)
	fleet.Fuel = fleet.Spec.FuelCapacity
	return fleet
}

// create a new small freighter (with cargo pod) fleet for testing
func testTeamster(player *Player) *Fleet {
	fleet := &Fleet{
		MapObject: MapObject{
			PlayerNum: player.Num,
		},
		BaseName: "Teamster",
		Tokens: []ShipToken{
			{
				Quantity:  1,
				DesignNum: 1,
				design: NewShipDesign(player, 1).
					WithHull(MediumFreighter.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: Crobmnium.Name, HullSlotIndex: 2, Quantity: 1},
						{HullComponent: RhinoScanner.Name, HullSlotIndex: 3, Quantity: 1},
					}).
					WithSpec(&rules, player)},
		},
		battlePlan:        &player.BattlePlans[0],
		OrbitingPlanetNum: None,
	}

	fleet.Spec = ComputeFleetSpec(&rules, player, fleet)
	fleet.Fuel = fleet.Spec.FuelCapacity
	return fleet

}

// create a new small freighter (with cargo pod) fleet for testing
func testPrivateer(player *Player, quantity int) *Fleet {
	fleet := &Fleet{
		MapObject: MapObject{
			PlayerNum: player.Num,
		},
		BaseName: "Privateer",
		Tokens: []ShipToken{
			{
				Quantity:  quantity,
				DesignNum: 1,
				design: NewShipDesign(player, 1).
					WithHull(Privateer.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: Crobmnium.Name, HullSlotIndex: 2, Quantity: 1},
						{HullComponent: CargoPod.Name, HullSlotIndex: 3, Quantity: 1},
						{HullComponent: CargoPod.Name, HullSlotIndex: 4, Quantity: 1},
						{HullComponent: CargoPod.Name, HullSlotIndex: 5, Quantity: 1},
					}).
					WithSpec(&rules, player)},
		},
		battlePlan:        &player.BattlePlans[0],
		OrbitingPlanetNum: None,
	}

	fleet.Spec = ComputeFleetSpec(&rules, player, fleet)
	fleet.Fuel = fleet.Spec.FuelCapacity
	return fleet

}

func Test_battle_getBestFleeMoves(t *testing.T) {
	type args struct {
		token   *battleToken
		weapons []*battleWeaponSlot
	}

	// generate a fleeing token at a position with 10dp
	fleeingToken := func(position BattleVector) *battleToken {
		return &battleToken{
			BattleRecordToken: BattleRecordToken{Position: position, PlayerNum: 1},
			ShipToken:         &ShipToken{Quantity: 1},
			armor:             10,
		}
	}

	// generate an enemy weapon at a position with a single laser
	enemyWeapon := func(position BattleVector) *battleWeaponSlot {
		return &battleWeaponSlot{
			token: &battleToken{
				BattleRecordToken: BattleRecordToken{Position: position, PlayerNum: 2, Tactic: BattleTacticMaximizeDamage, AttackWho: BattleAttackWhoEveryone, PrimaryTarget: BattleTargetAny},
				ShipToken:         &ShipToken{Quantity: 1},
				player:            testPlayer().WithNum(2),
				attributes:        battleTokenAttributeArmed,
			},
			weaponType:   battleWeaponTypeBeam,
			power:        10,
			weaponRange:  1,
			slotQuantity: 1,
		}
	}

	tests := []struct {
		name string
		args args
		want []BattleVector
	}{
		{
			name: "token at 1,4 move randomly 1st option",
			args: args{
				token: fleeingToken(BattleVector{1, 4}),
			},
			// randomly pick from all moves
			want: []BattleVector{{0, 3}, {0, 4}, {0, 5}, {1, 3}, {1, 4}, {1, 5}, {2, 3}, {2, 4}, {2, 5}},
		},
		{
			name: "token at 1,4 move away from surrounding weapons",
			args: args{
				token: fleeingToken(BattleVector{1, 4}),
				// make three weapons adjacent so we have to move straight back
				weapons: []*battleWeaponSlot{
					enemyWeapon(BattleVector{1, 5}),
					enemyWeapon(BattleVector{2, 4}),
					enemyWeapon(BattleVector{1, 3}),
				},
			},
			want: []BattleVector{{0, 3}, {0, 5}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &battle{
				tokens: []*battleToken{tt.args.token},
				rules:  &rules,
			}

			// run away and record the new position
			if got := b.getBestFleeMoves(tt.args.token, tt.args.weapons); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("battle.getBestMove() = %v, want %v", got, tt.want)
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
		weaponRange:  1,
		slotQuantity: 1,
	}

	// generate a fleeing token at a position with 10dp
	attackingToken := func(position BattleVector, tactic BattleTactic, weapon battleWeaponSlot) *battleToken {
		player := testPlayer().WithNum(1)
		token := battleToken{
			BattleRecordToken: BattleRecordToken{Position: position, PlayerNum: player.Num, Tactic: tactic, AttackWho: BattleAttackWhoEveryone, PrimaryTarget: BattleTargetAny, Movement: 4},
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
	enemyToken := func(position BattleVector, weapon *battleWeaponSlot) *battleToken {
		player := testPlayer().WithNum(2)
		token := battleToken{
			BattleRecordToken: BattleRecordToken{Position: position, PlayerNum: player.Num, Tactic: BattleTacticMaximizeDamage, AttackWho: BattleAttackWhoEveryone, PrimaryTarget: BattleTargetAny, Movement: 4},
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
		want []BattleVector
	}{
		{
			// attacker should move over, or over and up/down
			// * * * *
			// A * * T
			// * * * *
			name: "move towards enemy",
			args: args{
				token: attackingToken(BattleVector{0, 1}, BattleTacticMaximizeDamage, laser),
				// make three weapons adjacent so we have to move straight back
				enemies: []*battleToken{
					enemyToken(BattleVector{4, 1}, nil),
				},
			},
			want: []BattleVector{{1, 0}, {1, 1}, {1, 2}},
		},
		{
			// attacker should move over, or over and down
			// A * * *
			// * * * *
			// * * * T
			name: "move towards enemy right or right/down",
			args: args{
				token: attackingToken(BattleVector{0, 0}, BattleTacticMaximizeDamage, laser),
				// make three weapons adjacent so we have to move straight back
				enemies: []*battleToken{
					enemyToken(BattleVector{3, 2}, nil),
				},
			},
			want: []BattleVector{{1, 0}, {1, 1}},
		},
		{
			// attacker should move on top to maximize damage
			// A T * *
			// * * * *
			// * * * *
			name: "maximize beam damage one target",
			args: args{
				token: attackingToken(BattleVector{0, 0}, BattleTacticMaximizeDamage, laser),
				// make three weapons adjacent so we have to move straight back
				enemies: []*battleToken{
					enemyToken(BattleVector{1, 0}, nil),
				},
			},
			want: []BattleVector{{1, 0}},
		},
		{
			// attacker should move to cause the most damage vs damage taken
			// the board will have two tokens, a strong and a weak one
			// * S * *
			// A W * *
			// * * * *
			name: "maximize damage ratio strong and weak target",
			args: args{
				token: attackingToken(BattleVector{0, 1}, BattleTacticMaximizeDamageRatio, laser),
				// make three weapons adjacent so we have to move straight back
				enemies: []*battleToken{
					// strong 100 power beamer
					enemyToken(BattleVector{1, 0}, &battleWeaponSlot{weaponType: battleWeaponTypeBeam, power: 100, slotQuantity: 1, weaponRange: 1}),
					// weak 10 power beamer
					enemyToken(BattleVector{1, 1}, &laser),
				},
			},
			want: []BattleVector{{1, 2}}, // best damage ratio, and towards center
		},
		{
			// attacker has torpedos and wants to stay out of range of those lasers
			// it should move back
			// * * 1 *
			// * A * *
			// * * 2 *
			name: "maximize damage ratio, prefer no damage",
			args: args{
				token: attackingToken(BattleVector{1, 1}, BattleTacticMaximizeDamageRatio, battleWeaponSlot{weaponType: battleWeaponTypeTorpedo, power: 10, accuracy: 1, slotQuantity: 3, weaponRange: 2}),
				// make three weapons adjacent so we have to move straight back
				enemies: []*battleToken{
					enemyToken(BattleVector{2, 0}, &laser),
					// put two tokens here
					enemyToken(BattleVector{2, 2}, &laser),
					enemyToken(BattleVector{2, 2}, &laser),
				},
			},
			want: []BattleVector{{0, 0}, {0, 1}, {0, 2}},
		},
		{
			// attacker has torpedos and wants to stay out of range of those lasers
			// it should move back but stay near the center if possible
			// * * 1 *
			// * A * *
			// * * 2 *
			name: "maximize damage ratio, prefer no damage, start in center (5,5)",
			args: args{
				token: attackingToken(BattleVector{5, 5}, BattleTacticMaximizeDamageRatio, battleWeaponSlot{weaponType: battleWeaponTypeTorpedo, power: 10, accuracy: 1, slotQuantity: 3, weaponRange: 2}),
				// make three weapons adjacent so we have to move straight back
				enemies: []*battleToken{
					enemyToken(BattleVector{6, 4}, &laser),
					// put two tokens here
					enemyToken(BattleVector{6, 6}, &laser),
					enemyToken(BattleVector{6, 6}, &laser),
				},
			},
			want: []BattleVector{{4, 4}, {4, 5}}, // we pick the 0 damageTaken options that keep us near center
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
				token: attackingToken(BattleVector{0, 1}, BattleTacticMaximizeNetDamage, battleWeaponSlot{weaponType: battleWeaponTypeBeam, power: 20, slotQuantity: 1, weaponRange: 1}),
				// make three weapons adjacent so we have to move straight back
				enemies: []*battleToken{
					enemyToken(BattleVector{1, 0}, &laser),
					// tokens at this spot
					enemyToken(BattleVector{1, 2}, nil),
					enemyToken(BattleVector{1, 2}, &laser),
				},
			},
			want: []BattleVector{{1, 2}}, // best net damage, move to two token square
		},
		{
			// attacker has torpedos and wants to stay out of range of those lasers
			// it should move back
			// * * L *
			// * A * *
			// * * L *
			name: "minimize damage to self",
			args: args{
				token: attackingToken(BattleVector{1, 1}, BattleTacticMinimizeDamageToSelf, battleWeaponSlot{weaponType: battleWeaponTypeTorpedo, power: 10, accuracy: 1, slotQuantity: 2, weaponRange: 2}),
				// make three weapons adjacent so we have to move straight back
				enemies: []*battleToken{
					enemyToken(BattleVector{2, 0}, &laser),
					enemyToken(BattleVector{2, 2}, &laser),
				},
			},
			want: []BattleVector{{0, 0}, {0, 1}, {0, 2}}, // min damage, move out of range
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &battle{
				tokens: append([]*battleToken{tt.args.token}, tt.args.enemies...),
				rules:  &rules,
			}

			// build a list of weapons on the board
			weapons := append([]*battleWeaponSlot{}, tt.args.token.weaponSlots...)
			for _, token := range tt.args.enemies {
				weapons = append(weapons, token.weaponSlots...)
			}

			b.findTargets()
			// run away and record the new position
			if got := b.getBestAttackMoves(tt.args.token, weapons); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("battle.getBestMove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_battle_fireBeamWeapon(t *testing.T) {

	type weapon struct {
		weaponSlot   *battleWeaponSlot
		shipQuantity int
		position     BattleVector
	}
	type args struct {
		weapon  weapon
		targets []*battleToken
	}
	type want struct {
		damage            float64
		quantityDamaged   int
		quantityRemaining int
		stackShields      int
	}
	tests := []struct {
		name string
		args args
		want []want
	}{
		{name: "Single weapon, do 10 damage, no kills",
			args: args{
				weapon: weapon{
					weaponSlot: &battleWeaponSlot{
						slotQuantity: 1, // 1 beam weapon
						power:        10,
						weaponRange:  1,
					},
					shipQuantity: 1,
				},
				targets: []*battleToken{
					{
						ShipToken: &ShipToken{
							Quantity: 1,
							design:   &ShipDesign{Name: "defender"}, // for logging
						},
						armor: 20,
					},
				},
			},
			want: []want{{damage: 10, quantityDamaged: 1, quantityRemaining: 1}},
		},
		{name: "Single weapon, do 30 damage, to a ship stack with two ships, one damaged",
			args: args{
				weapon: weapon{
					weaponSlot: &battleWeaponSlot{
						slotQuantity: 1, // 1 beam weapon
						power:        30,
						weaponRange:  1,
					},
					shipQuantity: 1,
				},
				targets: []*battleToken{
					{
						ShipToken: &ShipToken{
							Quantity:        2,
							Damage:          5,
							QuantityDamaged: 1,
							design:          &ShipDesign{Name: "defender"}, // for logging
						},
						armor: 20,
					},
				},
			},
			want: []want{{damage: 15, quantityDamaged: 1, quantityRemaining: 1}},
		},
		{name: "Single weapon, do 10 damage reduced to 9 for range",
			args: args{
				weapon: weapon{
					weaponSlot: &battleWeaponSlot{
						slotQuantity: 1, // 1 beam weapon
						power:        10,
						weaponRange:  2,
					},
					shipQuantity: 1,
					position:     BattleVector{2, 0}, // 1 away from target
				},
				targets: []*battleToken{
					{
						ShipToken: &ShipToken{
							Quantity: 1,
							design:   &ShipDesign{Name: "defender"}, // for logging
						},
						BattleRecordToken: BattleRecordToken{
							Position: BattleVector{0, 0},
						},
						armor: 20,
					},
				},
			},
			want: []want{{damage: 9, quantityDamaged: 1, quantityRemaining: 1}},
		},
		{name: "two weapons, do 30 damage total, one (over)kill",
			args: args{
				weapon: weapon{
					weaponSlot: &battleWeaponSlot{
						slotQuantity: 2,  // 2 beam weapons
						power:        15, // 15 damage per beam
						weaponRange:  2,
					},
					shipQuantity: 1, // one ship in the attacker stack
				},
				targets: []*battleToken{
					{
						ShipToken: &ShipToken{
							Quantity: 1,
							design:   &ShipDesign{Name: "defender"},
						},
						armor: 20, // 20 armor, will be destroyed
					},
				},
			},
			want: []want{{damage: 0, quantityDamaged: 0, quantityRemaining: 0}},
		},
		{name: "two weapons, two ships, do 40 damage total, one kill, one damaged",
			args: args{
				weapon: weapon{
					weaponSlot: &battleWeaponSlot{
						slotQuantity: 2, // 2 beam weapons
						power:        10,
					},
					shipQuantity: 2, // 2 ships in attacker stack
				},
				targets: []*battleToken{
					{
						ShipToken: &ShipToken{
							Quantity: 2, // two ships in defender
							design:   &ShipDesign{Name: "defender"},
						},
						armor: 30,
					},
				},
			},
			want: []want{{damage: 10, quantityDamaged: 1, quantityRemaining: 1}},
		},
		{name: "two weapons, two stacks, do 20 damage total, kill both",
			args: args{
				weapon: weapon{
					weaponSlot: &battleWeaponSlot{
						slotQuantity: 2, // 2 beam weapons
						power:        10,
					},
					shipQuantity: 1, // 1 ships in attacker stack
				},
				targets: []*battleToken{
					{
						ShipToken: &ShipToken{
							Quantity: 1,
							design:   &ShipDesign{Name: "defender"},
						},
						armor: 10,
					},
					{
						ShipToken: &ShipToken{
							Quantity: 1,
							design:   &ShipDesign{Name: "defender"},
						},
						armor: 10,
					},
				},
			},
			// both stacks gone
			want: []want{
				{damage: 0, quantityDamaged: 0, quantityRemaining: 0},
				{damage: 0, quantityDamaged: 0, quantityRemaining: 0},
			},
		},
		{name: "two weapons, two stacks, do 20 damage total, don't get through shield of the first stack",
			args: args{
				weapon: weapon{
					weaponSlot: &battleWeaponSlot{
						slotQuantity: 2, // 2 beam weapons
						power:        10,
					},
					shipQuantity: 1, // 1 ships in attacker stack
				},
				targets: []*battleToken{
					{
						ShipToken: &ShipToken{
							Quantity: 3,
							design:   &ShipDesign{Name: "defender"},
						},
						armor:        10,
						shields:      10,
						stackShields: 30,
					},
					{
						ShipToken: &ShipToken{
							Quantity: 3,
							design:   &ShipDesign{Name: "defender"},
						},
						armor:        10,
						shields:      10,
						stackShields: 30,
					},
				},
			},
			// both stacks alive, but first stack with 20 less stackShields
			want: []want{
				{damage: 0, quantityDamaged: 0, quantityRemaining: 3, stackShields: 10},
				{damage: 0, quantityDamaged: 0, quantityRemaining: 3, stackShields: 30},
			},
		},
		{name: "one weapon, do 10 damage to shields, no damage",
			args: args{
				weapon: weapon{
					weaponSlot: &battleWeaponSlot{
						slotQuantity: 1,
						power:        10,
					},
					shipQuantity: 1,
				},
				targets: []*battleToken{
					{
						ShipToken: &ShipToken{
							Quantity: 1,
							design:   &ShipDesign{Name: "defender"},
						},
						armor:        30,
						shields:      20,
						stackShields: 20,
					},
				},
			},
			want: []want{{damage: 0, quantityDamaged: 0, quantityRemaining: 1, stackShields: 10}},
		},
		{name: "one super beam, do 100 damage destroy one stack and damage another",
			args: args{
				weapon: weapon{
					weaponSlot: &battleWeaponSlot{
						slotQuantity: 1,
						power:        100,
					},
					shipQuantity: 1,
				},
				targets: []*battleToken{
					{
						ShipToken: &ShipToken{
							Quantity: 1,
							design:   &ShipDesign{Name: "defender1"},
						},
						armor: 10,
					},
					{
						ShipToken: &ShipToken{
							Quantity: 1,
							design:   &ShipDesign{Name: "defender2"},
						},
						armor: 100,
					},
				},
			},
			want: []want{
				{damage: 0, quantityDamaged: 0, quantityRemaining: 0},
				{damage: 90, quantityDamaged: 1, quantityRemaining: 1},
			},
		},
		{name: "one minigun, do 10 damage to all targets",
			args: args{
				weapon: weapon{
					weaponSlot: &battleWeaponSlot{
						slotQuantity:   1,
						power:          10,
						hitsAllTargets: true,
					},
					shipQuantity: 1,
				},
				targets: []*battleToken{
					{
						ShipToken: &ShipToken{
							Quantity: 2,
							design:   &ShipDesign{Name: "defender1"},
						},
						armor: 10,
					},
					{
						ShipToken: &ShipToken{
							Quantity: 2,
							design:   &ShipDesign{Name: "defender2"},
						},
						armor: 100,
					},
				},
			},
			want: []want{
				{damage: 0, quantityDamaged: 0, quantityRemaining: 1}, // destroy one token
				{damage: 5, quantityDamaged: 2, quantityRemaining: 2}, // damage both tokens
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &battle{rules: &rules,
				record: newBattleRecord(1, None, Vector{}, []BattleRecordToken{})}
			b.record.recordNewRound()

			// setup this weapon's token based on shipQuantity and position
			tt.args.weapon.weaponSlot.token = &battleToken{
				ShipToken: &ShipToken{
					Quantity: tt.args.weapon.shipQuantity,
					design:   &ShipDesign{Name: "attacker"}, // for logging
				},
				BattleRecordToken: BattleRecordToken{
					Position: tt.args.weapon.position,
				},
			}

			// fire the beam weapon!
			b.fireBeamWeapon(tt.args.weapon.weaponSlot, tt.args.targets)

			for i, target := range tt.args.targets {
				if target.Quantity != tt.want[i].quantityRemaining {
					t.Errorf("battleWeaponSlot.fireBeamWeapon() target: %d quantityRemaining = %v, want %v", i, target.Quantity, tt.want[i].quantityRemaining)
				}
				if target.Damage != tt.want[i].damage {
					t.Errorf("battleWeaponSlot.fireBeamWeapon() target: %d damage = %v, want %v", i, target.Damage, tt.want[i].damage)
				}
				if target.QuantityDamaged != tt.want[i].quantityDamaged {
					t.Errorf("battleWeaponSlot.fireBeamWeapon() target: %d quantityDamaged = %v, want %v", i, target.QuantityDamaged, tt.want[i].quantityDamaged)
				}
				if target.stackShields != tt.want[i].stackShields {
					t.Errorf("battleWeaponSlot.fireBeamWeapon() target: %d stackShields = %v, want %v", i, target.stackShields, tt.want[i].stackShields)
				}
			}

		})
	}
}

func Test_battle_fireTorpedo(t *testing.T) {

	type weapon struct {
		weaponSlot   *battleWeaponSlot
		shipQuantity int
		position     BattleVector
	}
	type args struct {
		weapon  weapon
		targets []*battleToken
	}
	type want struct {
		damage            float64
		quantityDamaged   int
		quantityRemaining int
		stackShields      int
	}
	tests := []struct {
		name string
		args args
		want []want
	}{
		{name: "Single torpedo, do 10 damage, no kills",
			args: args{
				weapon: weapon{
					weaponSlot: &battleWeaponSlot{
						slotQuantity: 1, // 1 torpedo
						power:        10,
						accuracy:     1,
					},
					shipQuantity: 1,
				},
				targets: []*battleToken{
					{
						ShipToken: &ShipToken{
							Quantity: 1,
							design:   &ShipDesign{Name: "defender"}, // for logging
						},
						armor: 20,
					},
				},
			},
			want: []want{{damage: 10, quantityDamaged: 1, quantityRemaining: 1}},
		},
		{name: "Single torpedo, do 10 damage to a 2 ship stack with 1@5 damage",
			args: args{
				weapon: weapon{
					weaponSlot: &battleWeaponSlot{
						slotQuantity: 1, // 1 torpedo
						power:        10,
						accuracy:     1,
					},
					shipQuantity: 1,
				},
				targets: []*battleToken{
					{
						ShipToken: &ShipToken{
							Quantity:        2,
							QuantityDamaged: 1,
							Damage:          5,
							design:          &ShipDesign{Name: "defender"}, // for logging
						},
						armor: 20,
					},
				},
			},
			// TODO: not sure about this. It doesn't make sense for a torpedo to splash damage at the end...
			want: []want{{damage: 15 / 2., quantityDamaged: 2, quantityRemaining: 2}},
		},
		{name: "Single torpedo, do 30 damage to a stack with two ships, destroy one, other undamaged",
			args: args{
				weapon: weapon{
					weaponSlot: &battleWeaponSlot{
						slotQuantity: 1, // 1 torpedo
						power:        30,
						accuracy:     1,
					},
					shipQuantity: 1,
				},
				targets: []*battleToken{
					{
						ShipToken: &ShipToken{
							Quantity: 2,
							design:   &ShipDesign{Name: "defender"}, // for logging
						},
						armor: 20,
					},
				},
			},
			want: []want{{damage: 0, quantityDamaged: 0, quantityRemaining: 1}},
		},
		{name: "two torpedos, do 15 damage each, kill ship with first hit",
			args: args{
				weapon: weapon{
					weaponSlot: &battleWeaponSlot{
						slotQuantity: 2,  // 2 torpedos
						power:        15, // 15 damage per torpedo
						accuracy:     1,
					},
					shipQuantity: 1, // one ship in the attacker stack
				},
				targets: []*battleToken{
					{
						ShipToken: &ShipToken{
							Quantity: 1,
							design:   &ShipDesign{Name: "defender"},
						},
						armor: 10, // 10 armor, will be destroyed
					},
				},
			},
			want: []want{{damage: 0, quantityDamaged: 0, quantityRemaining: 0}},
		},
		{name: "two capital missiles, do 10 damage each, take down shields with first hit, double damage with second",
			args: args{
				weapon: weapon{
					weaponSlot: &battleWeaponSlot{
						slotQuantity:       2,  // 2 torpedos
						power:              10, // 15 damage per torpedo
						accuracy:           1,
						capitalShipMissile: true,
					},
					shipQuantity: 1, // one ship in the attacker stack
				},
				targets: []*battleToken{
					{
						ShipToken: &ShipToken{
							Quantity: 1,
							design:   &ShipDesign{Name: "defender"},
						},
						stackShields: 5,  // 5 shields will be gone first hit
						armor:        35, // 10 armor gone first hit, then 10x2=20 damage on second hit
					},
				},
			},
			want: []want{{damage: 30, quantityDamaged: 1, quantityRemaining: 1}},
		},
		{name: "two torpedos, two attacker ships, 4x torpedos do 40 damage total, one kill, one damaged",
			args: args{
				weapon: weapon{
					weaponSlot: &battleWeaponSlot{
						slotQuantity: 2, // 2 torpedos
						power:        10,
						accuracy:     1,
					},
					shipQuantity: 2, // 2 ships in attacker stack
				},
				targets: []*battleToken{
					{
						ShipToken: &ShipToken{
							Quantity: 2, // two ships in defender
							design:   &ShipDesign{Name: "defender"},
						},
						armor: 30,
					},
				},
			},
			want: []want{{damage: 10, quantityDamaged: 1, quantityRemaining: 1}},
		},
		{name: "from testbed, two omega torps w 300 power, 2 1700dp1300 damage",
			args: args{
				weapon: weapon{
					weaponSlot: &battleWeaponSlot{
						slotQuantity: 2,   // 2 torpedos
						power:        300, // 600 damage total
						accuracy:     1,
					},
					shipQuantity: 1,
				},
				targets: []*battleToken{
					{
						ShipToken: &ShipToken{
							Quantity:        3, // three defenders, 3@5 damaged
							QuantityDamaged: 3,
							Damage:          1300, // 400dp left
							design:          &ShipDesign{Name: "defender"},
						},
						armor: 1700,
					},
				},
			},
			// 600 damage total, first ship takes 400, 200 split between remaining ships
			want: []want{{damage: 1400, quantityDamaged: 2, quantityRemaining: 2}},
		},
		{name: "one torpedo, do 5 damage to shields, 5 damage to hull",
			args: args{
				weapon: weapon{
					weaponSlot: &battleWeaponSlot{
						slotQuantity: 1,
						power:        10,
						accuracy:     1,
					},
					shipQuantity: 1,
				},
				targets: []*battleToken{
					{
						ShipToken: &ShipToken{
							Quantity: 1,
							design:   &ShipDesign{Name: "defender"},
						},
						armor:        20,
						shields:      20,
						stackShields: 20,
					},
				},
			},
			want: []want{{damage: 5, quantityDamaged: 1, quantityRemaining: 1, stackShields: 15}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &battle{rules: &rules,
				record: newBattleRecord(1, None, Vector{}, []BattleRecordToken{})}
			b.record.recordNewRound()

			// setup this weapon's token based on shipQuantity and position
			tt.args.weapon.weaponSlot.token = &battleToken{
				ShipToken: &ShipToken{
					Quantity: tt.args.weapon.shipQuantity,
					design:   &ShipDesign{Name: "attacker"}, // for logging
				},
				BattleRecordToken: BattleRecordToken{
					Position: tt.args.weapon.position,
				},
			}

			// fire the beam weapon!
			b.fireTorpedo(tt.args.weapon.weaponSlot, tt.args.targets)

			for i, target := range tt.args.targets {
				if target.Quantity != tt.want[i].quantityRemaining {
					t.Errorf("battleWeaponSlot.fireTorpedo() target: %d quantityRemaining = %v, want %v", i, target.Quantity, tt.want[i].quantityRemaining)
				}
				if target.Damage != tt.want[i].damage {
					t.Errorf("battleWeaponSlot.fireTorpedo() target: %d damage = %v, want %v", i, target.Damage, tt.want[i].damage)
				}
				if target.QuantityDamaged != tt.want[i].quantityDamaged {
					t.Errorf("battleWeaponSlot.fireTorpedo() target: %d quantityDamaged = %v, want %v", i, target.QuantityDamaged, tt.want[i].quantityDamaged)
				}
				if target.stackShields != tt.want[i].stackShields {
					t.Errorf("battleWeaponSlot.fireTorpedo() target: %d stackShields = %v, want %v", i, target.stackShields, tt.want[i].stackShields)
				}
			}

		})
	}
}

func Test_battle_runBattle1(t *testing.T) {
	player1 := testPlayer().WithNum(1)
	player2 := testPlayer().WithNum(2)
	player1.Relations = []PlayerRelationship{{Relation: PlayerRelationFriend}, {Relation: PlayerRelationEnemy}}
	player2.Relations = []PlayerRelationship{{Relation: PlayerRelationEnemy}, {Relation: PlayerRelationFriend}}

	fleets := []*Fleet{
		testStalwartDefender(player1),
		testLongRangeScout(player1),
		testTeamster(player2),
	}

	designNum := 1
	for _, fleet := range fleets {
		for _, token := range fleet.Tokens {
			token.design.Num = designNum
			designNum += 1
		}
	}

	battle := newBattler(log.Logger, &rules, &StaticTechStore, 1, map[int]*Player{1: player1, 2: player2}, fleets, nil)

	record := battle.runBattle()

	// ran some number of turns
	assert.Greater(t, len(record.ActionsPerRound), 1)
	assert.Equal(t, 2, record.Stats.NumShipsByPlayer[player1.Num])
	assert.Equal(t, 1, record.Stats.NumShipsByPlayer[player2.Num])
}

func Test_battle_runBattle2(t *testing.T) {
	player1 := NewPlayer(0, NewRace()).WithNum(1)
	player2 := NewPlayer(0, NewRace()).WithNum(2)
	player1.Name = AINames[0][1]
	player2.Name = AINames[1][1]
	player1.Race.PluralName = AINames[0][1]
	player2.Race.PluralName = AINames[1][1]
	player1.Relations = []PlayerRelationship{{Relation: PlayerRelationFriend}, {Relation: PlayerRelationEnemy}}
	player2.Relations = []PlayerRelationship{{Relation: PlayerRelationEnemy}, {Relation: PlayerRelationFriend}}
	player1.PlayerIntels.PlayerIntels = []PlayerIntel{{Num: player1.Num}, {Num: player2.Num}}
	player2.PlayerIntels.PlayerIntels = []PlayerIntel{{Num: player1.Num}, {Num: player2.Num}}

	player1.Designs = append(player1.Designs,
		NewShipDesign(player1, 1).
			WithName("Battle Cruiser").
			WithHull(BattleCruiser.Name).
			WithSlots([]ShipDesignSlot{
				{HullComponent: TransStar10.Name, HullSlotIndex: 1, Quantity: 2},
				{HullComponent: Overthruster.Name, HullSlotIndex: 2, Quantity: 2},
				{HullComponent: BattleSuperComputer.Name, HullSlotIndex: 3, Quantity: 2},
				{HullComponent: ColloidalPhaser.Name, HullSlotIndex: 4, Quantity: 3},
				{HullComponent: DeltaTorpedo.Name, HullSlotIndex: 5, Quantity: 3},
				{HullComponent: Overthruster.Name, HullSlotIndex: 6, Quantity: 3},
				{HullComponent: GorillaDelagator.Name, HullSlotIndex: 7, Quantity: 4},
			}),
	)

	player2.Designs = append(player2.Designs,
		NewShipDesign(player2, 1).
			WithName("Teamster").
			WithHull(SmallFreighter.Name).
			WithSlots([]ShipDesignSlot{
				{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
				{HullComponent: Crobmnium.Name, HullSlotIndex: 2, Quantity: 1},
				{HullComponent: RhinoScanner.Name, HullSlotIndex: 3, Quantity: 1},
			}),
		NewShipDesign(player2, 2).
			WithName("Long Range Scout").
			WithHull(Scout.Name).
			WithSlots([]ShipDesignSlot{
				{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
				{HullComponent: RhinoScanner.Name, HullSlotIndex: 2, Quantity: 1},
				{HullComponent: CompletePhaseShield.Name, HullSlotIndex: 3, Quantity: 1},
			}),
		NewShipDesign(player2, 3).
			WithName("Jammed&Fluxed Defender").
			WithHull(Destroyer.Name).
			WithSlots([]ShipDesignSlot{
				{HullComponent: TransStar10.Name, HullSlotIndex: 1, Quantity: 1},
				{HullComponent: ColloidalPhaser.Name, HullSlotIndex: 2, Quantity: 1},
				{HullComponent: ColloidalPhaser.Name, HullSlotIndex: 3, Quantity: 1},
				{HullComponent: RhinoScanner.Name, HullSlotIndex: 4, Quantity: 1},
				{HullComponent: Superlatanium.Name, HullSlotIndex: 5, Quantity: 1},
				{HullComponent: Jammer30.Name, HullSlotIndex: 6, Quantity: 1},
				{HullComponent: FluxCapacitor.Name, HullSlotIndex: 7, Quantity: 1},
			}),
		NewShipDesign(player2, 4).
			WithName("Stalwart Sapper").
			WithHull(Destroyer.Name).
			WithSlots([]ShipDesignSlot{
				{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
				{HullComponent: PulsedSapper.Name, HullSlotIndex: 2, Quantity: 1},
				{HullComponent: PulsedSapper.Name, HullSlotIndex: 3, Quantity: 1},
				{HullComponent: RhinoScanner.Name, HullSlotIndex: 4, Quantity: 1},
				{HullComponent: Superlatanium.Name, HullSlotIndex: 5, Quantity: 1},
				{HullComponent: Overthruster.Name, HullSlotIndex: 6, Quantity: 1},
				{HullComponent: Overthruster.Name, HullSlotIndex: 7, Quantity: 1},
			}),
	)

	fleets := []*Fleet{
		{
			MapObject: MapObject{
				PlayerNum: player1.Num,
			},
			BaseName: "Battle Cruiser",
			Tokens: []ShipToken{
				{
					DesignNum: player1.Designs[0].Num,
					Quantity:  2,
				},
			},
		},
		// player2's teamster
		{
			MapObject: MapObject{
				PlayerNum: player2.Num,
			},
			BaseName: "Teamster+",
			Tokens: []ShipToken{
				{
					Quantity:  5,
					DesignNum: player2.Designs[0].Num,
				},
				{
					Quantity:  2,
					DesignNum: player2.Designs[1].Num,
				},
				{
					Quantity:  3,
					DesignNum: player2.Designs[2].Num,
				},
				{
					Quantity:  4,
					DesignNum: player2.Designs[3].Num,
				},
			},
		}}

	record, _ := RunTestBattle([]*Player{player1, player2}, fleets)
	// ran some number of turns
	assert.Less(t, 5, len(record.ActionsPerRound))
}

func Test_battle_runBattleError(t *testing.T) {
	player1 := NewPlayer(0, NewRace()).WithNum(1)
	player2 := NewPlayer(0, NewRace()).WithNum(2)
	player1.Name = AINames[0][1]
	player2.Name = AINames[1][1]
	player1.Race.PluralName = AINames[0][1]
	player2.Race.PluralName = AINames[1][1]
	player1.Relations = []PlayerRelationship{{Relation: PlayerRelationFriend}, {Relation: PlayerRelationEnemy}}
	player2.Relations = []PlayerRelationship{{Relation: PlayerRelationEnemy}, {Relation: PlayerRelationFriend}}
	player1.PlayerIntels.PlayerIntels = []PlayerIntel{{Num: player1.Num}, {Num: player2.Num}}
	player2.PlayerIntels.PlayerIntels = []PlayerIntel{{Num: player1.Num}, {Num: player2.Num}}

	player1.Designs = append(player1.Designs,
		NewShipDesign(player1, 1).
			WithName("Battle Cruiser").
			WithHull(BattleCruiser.Name).
			WithSlots([]ShipDesignSlot{
				{HullComponent: TransStar10.Name, HullSlotIndex: 1, Quantity: 2},
				{HullComponent: Overthruster.Name, HullSlotIndex: 2, Quantity: 2},
				{HullComponent: BattleSuperComputer.Name, HullSlotIndex: 3, Quantity: 2},
				{HullComponent: ColloidalPhaser.Name, HullSlotIndex: 4, Quantity: 3},
				{HullComponent: DeltaTorpedo.Name, HullSlotIndex: 5, Quantity: 3},
				{HullComponent: Overthruster.Name, HullSlotIndex: 6, Quantity: 3},
				{HullComponent: GorillaDelagator.Name, HullSlotIndex: 7, Quantity: 4},
			}),
	)

	player2.Designs = append(player2.Designs,
		NewShipDesign(player2, 1).
			WithName("BANANA BOAT").
			WithHull("Banana Ship").
			WithSlots([]ShipDesignSlot{
				{HullComponent: "Ice Cream", HullSlotIndex: 1, Quantity: 1},
				{HullComponent: "Chocolate", HullSlotIndex: 2, Quantity: 1},
				{HullComponent: "Hot Fudge", HullSlotIndex: 3, Quantity: 1},
			}),
		NewShipDesign(player2, 2).
			WithName("Jammed&Fluxed Defender").
			WithHull(Destroyer.Name).
			WithSlots([]ShipDesignSlot{
				{HullComponent: TransStar10.Name, HullSlotIndex: 1, Quantity: 1},
				{HullComponent: ColloidalPhaser.Name, HullSlotIndex: 2, Quantity: 1},
				{HullComponent: ColloidalPhaser.Name, HullSlotIndex: 3, Quantity: 1},
				{HullComponent: RhinoScanner.Name, HullSlotIndex: 4, Quantity: 1},
				{HullComponent: Superlatanium.Name, HullSlotIndex: 5, Quantity: 1},
				{HullComponent: Jammer30.Name, HullSlotIndex: 6, Quantity: 1},
				{HullComponent: FluxCapacitor.Name, HullSlotIndex: 7, Quantity: 1},
			}),
	)

	fleets := []*Fleet{
		{
			MapObject: MapObject{
				PlayerNum: player1.Num,
			},
			BaseName: "Battle Cruiser",
			Tokens: []ShipToken{
				{
					DesignNum: player1.Designs[0].Num,
					Quantity:  2,
				},
			},
		},
		// player2's ~~teamster~~ BANANA
		{
			MapObject: MapObject{
				PlayerNum: player2.Num,
			},
			BaseName: "Banana+",
			Tokens: []ShipToken{
				{
					Quantity:  5,
					DesignNum: player2.Designs[0].Num,
				},
				{
					Quantity:  2,
					DesignNum: player2.Designs[1].Num,
				},
			},
		},
	}

	_, err := RunTestBattle([]*Player{player1, player2}, fleets)
	// should return error due to incorrect spec on teamster from nonexistent hull/parts
	assert.Error(t, err)
}

func Test_updateMovesWithCenterPreference(t *testing.T) {
	type args struct {
		better      bool
		newPosition BattleVector
		bestMoves   []BattleVector
	}
	tests := []struct {
		name string
		args args
		want []BattleVector
	}{
		{
			name: "better move 1,0, pick it",
			args: args{better: true, newPosition: BattleVector{1, 0}, bestMoves: []BattleVector{{0, 0}}},
			want: []BattleVector{{1, 0}},
		},
		{
			name: "better move away from center, pick it",
			args: args{better: true, newPosition: BattleVector{3, 3}, bestMoves: []BattleVector{{4, 4}, {4, 5}}},
			want: []BattleVector{{3, 3}},
		},
		{
			name: "equivalent damage move, but newPosition is closer to center",
			args: args{better: false, newPosition: BattleVector{4, 4}, bestMoves: []BattleVector{{4, 3}}},
			want: []BattleVector{{4, 4}},
		},
		{
			name: "equivalent damage move, newPosition is closer to center",
			args: args{better: false, newPosition: BattleVector{4, 5}, bestMoves: []BattleVector{{4, 3}}},
			want: []BattleVector{{4, 5}},
		},
		{
			name: "equivalent damage move, newPosition is same distance to center",
			args: args{better: false, newPosition: BattleVector{4, 5}, bestMoves: []BattleVector{{4, 4}}},
			want: []BattleVector{{4, 4}, {4, 5}},
		},
		{
			name: "equivalent damage move, newPosition is same distance to center",
			args: args{better: false, newPosition: BattleVector{5, 5}, bestMoves: []BattleVector{{4, 4}, {4, 5}}},
			want: []BattleVector{{4, 4}, {4, 5}, {5, 5}},
		},
		{
			name: "equivalent damage move, newPosition is farther from center, discard it",
			args: args{better: false, newPosition: BattleVector{6, 5}, bestMoves: []BattleVector{{4, 4}, {4, 5}}},
			want: []BattleVector{{4, 4}, {4, 5}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := updateMovesWithCenterPreference(tt.args.better, tt.args.newPosition, tt.args.bestMoves); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("updateBestPositions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getBattleMovement(t *testing.T) {
	type args struct {
		idealEngineSpeed int
		mass             int
		numEngines       int
		movementBonus    int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Destroyer + Trans Galactic Drive + thruster", args{idealEngineSpeed: 9, mass: 244, numEngines: 1, movementBonus: 1}, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getBattleMovement(tt.args.idealEngineSpeed, tt.args.movementBonus, tt.args.mass, tt.args.numEngines); got != tt.want {
				t.Errorf("getBattleMovement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_battle_buildMovementOrder(t *testing.T) {
	tokenMovement2Mass1 := &battleToken{
		ShipToken: &ShipToken{},
		BattleRecordToken: BattleRecordToken{
			PlayerNum: 1,
			Num:       1,
			Movement:  2,
			Mass:      1,
		},
	}
	tokenMovement4Mass1 := &battleToken{
		ShipToken: &ShipToken{},
		BattleRecordToken: BattleRecordToken{
			PlayerNum: 2,
			Num:       2,
			Movement:  4,
			Mass:      1,
		},
	}
	tokenMovement4Mass2 := &battleToken{
		ShipToken: &ShipToken{},
		BattleRecordToken: BattleRecordToken{
			PlayerNum: 3,
			Num:       3,
			Movement:  4,
			Mass:      2,
		},
	}
	tokenMovement5Mass2 := &battleToken{
		ShipToken: &ShipToken{},
		BattleRecordToken: BattleRecordToken{
			PlayerNum: 4,
			Num:       4,
			Movement:  5,
			Mass:      2,
		},
	}
	tokenMovement3Mass1 := &battleToken{
		ShipToken: &ShipToken{},
		BattleRecordToken: BattleRecordToken{
			PlayerNum: 4,
			Num:       4,
			Movement:  3,
			Mass:      1,
		},
	}

	tokenMovement6Mass190 := &battleToken{
		ShipToken: &ShipToken{},
		BattleRecordToken: BattleRecordToken{
			PlayerNum: 5,
			Num:       5,
			Movement:  6,
			Mass:      190,
		},
	}
	tokenMovement4Mass22 := &battleToken{
		ShipToken: &ShipToken{},
		BattleRecordToken: BattleRecordToken{
			PlayerNum: 6,
			Num:       6,
			Movement:  4,
			Mass:      22,
		},
	}
	tokenMovement4Mass41 := &battleToken{
		ShipToken: &ShipToken{},
		BattleRecordToken: BattleRecordToken{
			PlayerNum: 6,
			Num:       7,
			Movement:  4,
			Mass:      41,
		},
	}

	type args struct {
		tokens []*battleToken
	}
	tests := []struct {
		name          string
		args          args
		wantMoveOrder [4][]*battleToken
	}{
		{name: "one token, move 2", args: args{[]*battleToken{tokenMovement2Mass1}},
			wantMoveOrder: [4][]*battleToken{
				{tokenMovement2Mass1},
				nil,
				{tokenMovement2Mass1},
				nil,
			},
		},
		{name: "two tokens, same mass, different movements", args: args{[]*battleToken{tokenMovement2Mass1, tokenMovement4Mass1}},
			wantMoveOrder: [4][]*battleToken{
				{tokenMovement2Mass1, tokenMovement4Mass1},
				{tokenMovement4Mass1},
				{tokenMovement2Mass1, tokenMovement4Mass1},
				{tokenMovement4Mass1},
			},
		},
		// higher mass moves first
		{name: "two tokens, diff mass, same movement", args: args{[]*battleToken{tokenMovement4Mass1, tokenMovement4Mass2}},
			wantMoveOrder: [4][]*battleToken{
				{tokenMovement4Mass2, tokenMovement4Mass1},
				{tokenMovement4Mass2, tokenMovement4Mass1},
				{tokenMovement4Mass2, tokenMovement4Mass1},
				{tokenMovement4Mass2, tokenMovement4Mass1},
			},
		},
		// higher move/mass moves twice
		{name: "two tokens, one higher move/mass", args: args{[]*battleToken{tokenMovement4Mass1, tokenMovement5Mass2}},
			wantMoveOrder: [4][]*battleToken{
				{tokenMovement5Mass2, tokenMovement5Mass2, tokenMovement4Mass1},
				{tokenMovement5Mass2, tokenMovement4Mass1},
				{tokenMovement5Mass2, tokenMovement4Mass1},
				{tokenMovement5Mass2, tokenMovement4Mass1},
			},
		},
		// higher move/mass moves first
		{name: "two tokens, move3-mass1 and move5-mass2", args: args{[]*battleToken{tokenMovement3Mass1, tokenMovement5Mass2}},
			wantMoveOrder: [4][]*battleToken{
				{tokenMovement5Mass2, tokenMovement5Mass2, tokenMovement3Mass1},
				{tokenMovement5Mass2, tokenMovement3Mass1},
				{tokenMovement5Mass2},
				{tokenMovement5Mass2, tokenMovement3Mass1},
			},
		},
		// higher mass should move twice, then other tokens move
		{name: "three tokens", args: args{[]*battleToken{tokenMovement6Mass190, tokenMovement4Mass22, tokenMovement4Mass41}},
			wantMoveOrder: [4][]*battleToken{
				{tokenMovement6Mass190, tokenMovement6Mass190, tokenMovement4Mass41, tokenMovement4Mass22},
				{tokenMovement6Mass190, tokenMovement4Mass41, tokenMovement4Mass22},
				{tokenMovement6Mass190, tokenMovement6Mass190, tokenMovement4Mass41, tokenMovement4Mass22},
				{tokenMovement6Mass190, tokenMovement4Mass41, tokenMovement4Mass22},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &battle{rules: &rules}
			if gotMoveOrder := b.buildMovementOrder(tt.args.tokens); !reflect.DeepEqual(gotMoveOrder, tt.wantMoveOrder) {
				for _, r := range gotMoveOrder {
					for _, t := range r {
						s := t.String()
						_ = s
					}
				}
				t.Errorf("battle.buildMovementOrder() = %v, want %v", gotMoveOrder, tt.wantMoveOrder)
			}
		})
	}
}
