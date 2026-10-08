//go:build !wasi && !wasm

package cs

import (
	"log/slog"
	"reflect"
	"slices"
	"testing"

	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func testStalwartDefender(player *Player) *Fleet {
	return testStalwartDefenderWithQuantity(player, 1)
}

func testStalwartDefenderWithQuantity(player *Player, quantity int) *Fleet {
	fleet := &Fleet{
		MapObject: MapObject{
			Type:      MapObjectTypeFleet,
			PlayerNum: player.Num,
		},
		BaseName: "Stalwart Defender",
		Tokens: []ShipToken{
			{
				DesignNum: 1,
				Quantity:  quantity,
				design: NewShipDesign(player.Num, 1).
					WithHull(Destroyer.Name).
					WithSlots(slices.Clone(DesignStalwartDefender.Slots)).
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

func testPrivateer(player *Player, quantity int) *Fleet {
	fleet := &Fleet{
		MapObject: MapObject{
			Type:      MapObjectTypeFleet,
			PlayerNum: player.Num,
		},
		BaseName: "Privateer",
		Tokens: []ShipToken{
			{
				Quantity:  quantity,
				DesignNum: 1,
				design: NewShipDesign(player.Num, 1).
					WithHull(Privateer.Name).
					WithSlots(slices.Clone(DesignPrivateer.Slots)).
					WithSpec(&rules, player)},
		},
		battlePlan:        &player.BattlePlans[0],
		OrbitingPlanetNum: None,
	}

	fleet.Spec = ComputeFleetSpec(&rules, player, fleet)
	fleet.Fuel = fleet.Spec.FuelCapacity
	return fleet

}

func Test_getMovesForRound(t *testing.T) {
	tests := []struct {
		spd   int
		round int
		want  int
	}{
		// spd % 4 == 0
		{spd: 0, round: 0, want: 1}, // base=0, even -> +1
		{spd: 0, round: 1, want: 0}, // odd -> no +1
		{spd: 4, round: 2, want: 2}, // base=1, even
		{spd: 4, round: 3, want: 1}, // odd

		// spd % 4 == 1
		{spd: 1, round: 0, want: 1}, // base=0, r%4=0 !=2 -> +1
		{spd: 1, round: 1, want: 1}, // r%4=1 !=2 -> +1
		{spd: 1, round: 2, want: 0}, // r%4=2 -> no +1
		{spd: 1, round: 3, want: 1}, // r%4=3 !=2 -> +1
		{spd: 5, round: 2, want: 1}, // base=1, r%4=2 -> no +1

		// spd % 4 == 2
		{spd: 2, round: 0, want: 1}, // base=1, no +1
		{spd: 2, round: 1, want: 1},
		{spd: 6, round: 3, want: 2}, // base=2, no +1

		// spd % 4 == 3
		{spd: 3, round: 0, want: 2}, // base=1, r%4=0 -> +1
		{spd: 3, round: 1, want: 1}, // r%4=1 -> no +1
		{spd: 3, round: 2, want: 1}, // r%4=2 -> no +1
		{spd: 3, round: 3, want: 1}, // r%4=3 -> no +1
		{spd: 7, round: 4, want: 3}, // base=2, r%4=0 -> +1
	}

	for _, tt := range tests {
		if got := getMovesForRound(tt.spd, tt.round); got != tt.want {
			t.Errorf("BattleGetMovesForRound(spd=%d, round=%d) = %d, want %d",
				tt.spd, tt.round, got, tt.want)
		}
	}
}

func Test_battle_fireBeamWeapon(t *testing.T) {

	type weapon struct {
		weaponSlot   *battleWeaponSlot
		shipQuantity int
		position     Vector
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
						beamBonus:    1,
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
						beamBonus:    1,
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
						beamBonus:    1,
						weaponRange:  2,
					},
					shipQuantity: 1,
					position:     Vector{2, 0}, // 1 away from target
				},
				targets: []*battleToken{
					{
						ShipToken: &ShipToken{
							Quantity: 1,
							design:   &ShipDesign{Name: "defender"}, // for logging
						},
						BattleRecordToken: BattleRecordToken{
							Position: Vector{0, 0},
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
						beamBonus:    1,
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
						beamBonus:    1,
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
			want: []want{{damage: 10.02, quantityDamaged: 1, quantityRemaining: 1}},
		},
		{name: "two weapons, two stacks, do 20 damage total, kill both",
			args: args{
				weapon: weapon{
					weaponSlot: &battleWeaponSlot{
						slotQuantity: 2, // 2 beam weapons
						power:        10,
						beamBonus:    1,
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
						beamBonus:    1,
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
						beamBonus:    1,
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
						beamBonus:    1,
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
						beamBonus:      1,
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
			b := &battle{
				rules:  &rules,
				log:    testLogger,
				record: newBattleRecord(1, None, Vector{}, []BattleRecordToken{}),
			}
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
		position     Vector
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
			want: []want{{damage: 7.52, quantityDamaged: 2, quantityRemaining: 2}},
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
		{name: "two torpedoes, do 15 damage each, kill ship with first hit",
			args: args{
				weapon: weapon{
					weaponSlot: &battleWeaponSlot{
						slotQuantity: 2,  // 2 torpedoes
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
		{name: "1 ship with 2 jihads hitting unarmored target",
			args: args{
				weapon: weapon{
					weaponSlot: &battleWeaponSlot{
						slotQuantity:       2,
						power:              85,
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
						armor: 350,
					},
				},
			},
			want: []want{{damage: 340.2, quantityDamaged: 1, quantityRemaining: 1}},
		},
		{name: "two capital missiles, start shielded, retain normal volley power",
			args: args{
				weapon: weapon{
					weaponSlot: &battleWeaponSlot{
						slotQuantity:       2,  // 2 torpedoes
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
			want: []want{{damage: 15.05, quantityDamaged: 1, quantityRemaining: 1}},
		},
		{name: "two torpedoes, two attacker ships, 4x torpedoes do 40 damage total, one kill, one damaged",
			args: args{
				weapon: weapon{
					weaponSlot: &battleWeaponSlot{
						slotQuantity: 2, // 2 torpedoes
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
			want: []want{{damage: 10.02, quantityDamaged: 1, quantityRemaining: 1}},
		},
		{name: "from testbed, two omega torps w 300 power, 2 1700dp1300 damage",
			args: args{
				weapon: weapon{
					weaponSlot: &battleWeaponSlot{
						slotQuantity: 2,   // 2 torpedoes
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
			want: []want{{damage: 1400.8, quantityDamaged: 2, quantityRemaining: 2}},
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
			b := &battle{
				rules:  &rules,
				log:    testLogger,
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

	battle := newBattler(slog.Default(), &rules, 1, map[int]*Player{1: player1, 2: player2}, fleets, nil)

	record := battle.runBattle()

	// ran some number of turns
	assert.Greater(t, len(record.ActionsPerRound), 1)
	assert.Equal(t, 2, record.Stats.NumShipsByPlayer[player1.Num])
	assert.Equal(t, 1, record.Stats.NumShipsByPlayer[player2.Num])
}

func Test_getBattleSpeed(t *testing.T) {
	type args struct {
		idealEngineSpeed int
		mass             int
		numEngines       int
		movementBonus    float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"248 kT Destroyer + Trans Galactic Drive + thruster", args{idealEngineSpeed: 9, mass: 248, numEngines: 1, movementBonus: 1}, 5},
		{"69 kT Destroyer + 1 Enigma Pulsar", args{idealEngineSpeed: 10, mass: 69, numEngines: 1, movementBonus: 0.5}, 9},
		{"71 kT Destroyer + 1 Enigma Pulsar", args{idealEngineSpeed: 10, mass: 71, numEngines: 1, movementBonus: 0.5}, 8},
		{"71 kT Destroyer + 1 Enigma Pulsar + WM", args{idealEngineSpeed: 10, mass: 71, numEngines: 1, movementBonus: 2.5}, 10},
		{"71 kT Cruiser w/ 2 Enigma Pulsars", args{idealEngineSpeed: 10, mass: 71, numEngines: 2, movementBonus: 1}, 9},
		{"572 kT Miner w/ Radiating Hydro Ram Scoop", args{idealEngineSpeed: 6, mass: 572, numEngines: 1, movementBonus: 0}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getBattleSpeed(rules.MovementMin, rules.MovementMax, tt.args.idealEngineSpeed, tt.args.movementBonus, tt.args.mass, tt.args.numEngines); got != tt.want {
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
			b := &battle{
				rules: &rules,
				log:   testLogger,
			}
			for _, token := range tt.args.tokens {
				token.movementMass = float64(token.Mass)
			}
			for round, want := range tt.wantMoveOrder {
				assert.Equal(t, want, b.buildMovementOrder(tt.args.tokens, round), "round %d", round)
			}
		})
	}
}

// Test_battle_fireWeaponSlot verifies beam and torpedo damage, casualties, and challenged retreat.
func Test_battle_fireWeaponSlot(t *testing.T) {
	type args struct {
		weaponType                                              battleWeaponType
		power, count, quantity, armor, shields, quantityDamaged int
		damage, accuracy                                        float64
		sapper, missile                                         bool
	}
	type want struct {
		quantity, shields, quantityDamaged, destroyed int
		damage                                        float64
		challenged                                    bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{name: "shield hit preserves armor damage", args: args{power: 10, count: 1, quantity: 2, armor: 100, shields: 100, damage: 50, quantityDamaged: 2}, want: want{quantity: 2, shields: 90, damage: 50, quantityDamaged: 2}},
		{name: "sapper preserves armor damage", args: args{power: 10, count: 1, quantity: 2, armor: 100, shields: 100, damage: 50, quantityDamaged: 2, sapper: true}, want: want{quantity: 2, shields: 90, damage: 50, quantityDamaged: 2}},
		{name: "small beam does not pool prior damage into kills", args: args{power: 10, count: 1, quantity: 2, armor: 100, damage: 50, quantityDamaged: 2}, want: want{quantity: 2, damage: 55, quantityDamaged: 2, challenged: true}},
		{name: "beam kills a damaged ship first", args: args{power: 75, count: 1, quantity: 2, armor: 100, damage: 50, quantityDamaged: 2}, want: want{quantity: 1, damage: 75, quantityDamaged: 1, destroyed: 1, challenged: true}},
		{name: "torpedo miss does splash damage", args: args{weaponType: battleWeaponTypeTorpedo, power: 80, count: 1, quantity: 1, armor: 100, shields: 100}, want: want{quantity: 1, shields: 90}},
		{name: "torpedo residual spreads across survivors", args: args{weaponType: battleWeaponTypeTorpedo, power: 75, count: 2, accuracy: 1, quantity: 3, armor: 100}, want: want{quantity: 2, damage: 25, quantityDamaged: 2, destroyed: 1, challenged: true}},
		{name: "shield breaking missile retains normal power", args: args{weaponType: battleWeaponTypeTorpedo, power: 100, count: 1, accuracy: 1, missile: true, quantity: 1, armor: 1000, shields: 25}, want: want{quantity: 1, damage: 76, quantityDamaged: 1, challenged: true}},
		{name: "unshielded missile doubles power", args: args{weaponType: battleWeaponTypeTorpedo, power: 100, count: 1, accuracy: 1, missile: true, quantity: 1, armor: 1000}, want: want{quantity: 1, damage: 200, quantityDamaged: 1, challenged: true}},
		{name: "one torpedo cannot destroy multiple ships", args: args{weaponType: battleWeaponTypeTorpedo, power: 100, count: 1, accuracy: 1, quantity: 3, armor: 10}, want: want{quantity: 2, destroyed: 1, challenged: true}},
		{name: "ship losses remove their shields", args: args{weaponType: battleWeaponTypeTorpedo, power: 20, count: 1, accuracy: 1, quantity: 2, armor: 10, shields: 200}, want: want{quantity: 1, shields: 95, destroyed: 1, challenged: true}},
		{name: "shield damage does not challenge retreat orders", args: args{power: 10, count: 1, quantity: 1, armor: 100, shields: 100}, want: want{quantity: 1, shields: 90}},
		{name: "armor damage starts a fresh retreat countdown", args: args{power: 10, count: 1, quantity: 1, armor: 100}, want: want{quantity: 1, damage: 10, quantityDamaged: 1, challenged: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRulesWithSeed(1)
			b := &battle{rules: &r, log: testLogger, round: 1, record: newBattleRecord(1, None, Vector{}, nil)}
			b.record.recordNewRound()
			target := &battleToken{BattleRecordToken: BattleRecordToken{PlayerNum: 2, Tactic: BattleTacticDisengageIfChallenged}, ShipToken: &ShipToken{Quantity: tt.args.quantity, Damage: tt.args.damage, QuantityDamaged: tt.args.quantityDamaged}, armor: tt.args.armor, stackShields: tt.args.shields, totalStackShields: tt.args.shields, movesMade: 12}
			weapon := &battleWeaponSlot{token: &battleToken{ShipToken: &ShipToken{Quantity: 1}}, weaponType: tt.args.weaponType, power: tt.args.power, beamBonus: 1, slotQuantity: tt.args.count, accuracy: tt.args.accuracy, damagesShieldsOnly: tt.args.sapper, capitalShipMissile: tt.args.missile}
			b.fireWeaponSlot(weapon, []*battleToken{target})
			assert.Equal(t, tt.want.quantity, target.Quantity)
			assert.Equal(t, tt.want.shields, target.stackShields)
			assert.Equal(t, tt.want.quantityDamaged, target.QuantityDamaged)
			assert.True(t, test.WithinTolerance(target.Damage, tt.want.damage, .001), "damage %v, want %v", target.Damage, tt.want.damage)
			round := b.record.ActionsPerRound[len(b.record.ActionsPerRound)-1]
			assert.Equal(t, tt.want.destroyed, round[0].TokensDestroyed)
			if tt.want.challenged {
				assert.Equal(t, BattleTacticDisengage, target.Tactic)
				assert.Equal(t, 0, target.movesMade)
			} else {
				assert.Equal(t, BattleTacticDisengageIfChallenged, target.Tactic)
			}
		})
	}
}

// Test_battle_retreatCountdown verifies escape timing and a fresh countdown after pursuit.
func Test_battle_retreatCountdown(t *testing.T) {
	tests := []struct {
		name       string
		moves      int
		challenged bool
		wantAway   bool
		wantMoves  int
	}{
		{name: "seventh retreat step remains on board", moves: 6, wantMoves: 7},
		{name: "leave after seven retreat steps", moves: 7, wantAway: true, wantMoves: 7},
		{name: "prior pursuit does not shorten retreat", moves: 12, challenged: true, wantMoves: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRulesWithSeed(1)
			token := testBattleToken(1, Vector{4, 4}, BattleTacticDisengage)
			token.movesMade = tt.moves
			b := &battle{rules: &r, round: 1, tokens: []*battleToken{token}, record: newBattleRecord(1, None, Vector{}, nil)}
			b.record.recordNewRound()
			if tt.challenged {
				token.Tactic = BattleTacticDisengageIfChallenged
				token.attributes = battleTokenAttributeArmed
				b.applyWeaponDamage(token, battleWeaponDamage{armorDamage: 10, damage: 10, quantityDamaged: 1})
			}
			b.moveToken(token)
			assert.Equal(t, tt.wantAway, token.ranAway)
			assert.Equal(t, tt.wantMoves, token.movesMade)
		})
	}
}

// Test_battle_hasHostility verifies targeting and battle continuation for stationary and retreating attackers.
func Test_battle_hasHostility(t *testing.T) {
	tests := []struct {
		name     string
		movement int
		tactic   BattleTactic
		hostile  bool
		want     bool
	}{
		{"stationary weapons sustain battle", 0, BattleTacticMaximizeDamage, true, true},
		{"retreating weapons retain firing targets", 4, BattleTacticDisengage, true, true},
		{"no hostility ends battle", 4, BattleTacticMaximizeDamage, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := &battleToken{BattleRecordToken: BattleRecordToken{PlayerNum: 1, Movement: tt.movement, Tactic: tt.tactic, PrimaryTarget: BattleTargetAny}, ShipToken: &ShipToken{Quantity: 1}, armor: 100, attributes: battleTokenAttributeArmed, attackPlayers: map[int]bool{2: tt.hostile}}
			enemy := &battleToken{BattleRecordToken: BattleRecordToken{PlayerNum: 2}, ShipToken: &ShipToken{Quantity: 1}, armor: 100, attackPlayers: map[int]bool{1: tt.hostile}}
			weapon := &battleWeaponSlot{token: source, power: 10, slotQuantity: 1, weaponRange: 1}
			source.weaponSlots = []*battleWeaponSlot{weapon}
			b := &battle{tokens: []*battleToken{source, enemy}}
			assert.Equal(t, tt.want, b.hasHostility())
			if tt.want {
				assert.Equal(t, []*battleToken{enemy}, weapon.findTargets(b.tokens))
			}
		})
	}
}

// Test_getBattleStartingPosition verifies formations for different participant counts.
func Test_getBattleStartingPosition(t *testing.T) {
	tests := []struct {
		name           string
		players, index int
		want           Vector
	}{
		{"two players first", 2, 0, Vector{1, 4}},
		{"two players second", 2, 1, Vector{8, 5}},
		{"three players first", 3, 0, Vector{4, 1}},
		{"three players second", 3, 1, Vector{8, 8}},
		{"three players third", 3, 2, Vector{1, 8}},
		{"sixteen players first", 16, 0, Vector{1, 1}},
		{"sixteen players last", 16, 15, Vector{6, 6}},
		{"seventeen players wrap to the sixteen player formation", 17, 16, Vector{1, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { assert.Equal(t, tt.want, getBattleStartingPosition(tt.players, tt.index)) })
	}
}

// Test_battle_torpedoHits verifies deterministic hit counts for large volleys.
func Test_battle_torpedoHits(t *testing.T) {
	tests := []struct {
		name     string
		count    int
		accuracy float64
		want     int
	}{
		{"all hit", 3, 1, 3},
		{"all miss", 3, 0, 0},
		{"large volley uses deterministic accuracy", 201, .5, 100},
		{"large fractional hit count is truncated", 999, .45, 449},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRulesWithSeed(1)
			b := &battle{rules: &r}
			assert.Equal(t, tt.want, b.torpedoHits(tt.count, tt.accuracy))
		})
	}
}

// Test_battle_salvageRecovery verifies recovery fractions for each battle location.
func Test_battle_salvageRecovery(t *testing.T) {
	starbase := &battleToken{ShipToken: &ShipToken{Quantity: 1}, attributes: battleTokenAttributeStarbase}
	destroyedStarbase := &battleToken{ShipToken: &ShipToken{}, attributes: battleTokenAttributeStarbase}
	tests := []struct {
		name   string
		planet *Planet
		tokens []*battleToken
		want   float64
	}{
		{"deep space", nil, nil, 0.75},
		{"planet without a base", &Planet{}, nil, 0.5},
		{"planet whose base was destroyed", &Planet{}, []*battleToken{destroyedStarbase}, 0.5},
		{"planet with a base", &Planet{}, []*battleToken{starbase}, 0.8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &battle{planet: tt.planet, tokens: tt.tokens}
			assert.Equal(t, tt.want, b.salvageRecovery())
		})
	}
}

// Test_battle_applyWeaponDamageRegeneration verifies that casualties reduce shields and regeneration capacity.
func Test_battle_applyWeaponDamageRegeneration(t *testing.T) {
	tests := []struct {
		name        string
		regen       bool
		wantShields int
	}{
		{"survivor has one ship of shields", false, 95},
		{"regeneration caps at surviving shield capacity", true, 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRulesWithSeed(1)
			player := NewPlayer(1, NewRace().WithLRT(RS).WithSpec(&r))
			target := &battleToken{ShipToken: &ShipToken{Quantity: 2}, player: player, armor: 10, shields: 100, stackShields: 200, totalStackShields: 200}
			weapon := &battleWeaponSlot{power: 20}
			b := &battle{rules: &r, record: newBattleRecord(1, None, Vector{}, nil)}
			b.applyWeaponDamage(target, weapon.getTorpedoVolleyDamage(target, 1, 0, 1, r.TorpedoSplashDamage))
			if tt.regen {
				target.regenerateShields()
			}
			assert.Equal(t, 1, target.Quantity)
			assert.Equal(t, 100, target.totalStackShields)
			assert.Equal(t, tt.wantShields, target.stackShields)
		})
	}
}

// Test_battle_initiativeTies verifies stable firing order for initiative ties and capped initiative.
func Test_battle_initiativeTies(t *testing.T) {
	tests := []struct {
		name       string
		initiative int
		want       int
	}{
		{"initiative below cap", 10, 20},
		{"initiative is capped", 60, 63},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := []*battleToken{{BattleRecordToken: BattleRecordToken{Initiative: tt.initiative}}, {BattleRecordToken: BattleRecordToken{Initiative: tt.initiative}}}
			for _, token := range tokens {
				token.weaponSlots = []*battleWeaponSlot{newBattleWeaponSlot(token, ShipDesignSlot{Quantity: 1}, &TechHullComponent{Initiative: 10}, 0, 0, 1), newBattleWeaponSlot(token, ShipDesignSlot{Quantity: 1}, &TechHullComponent{Initiative: 10}, 0, 0, 1)}
			}
			b := &battle{}
			want := []*battleWeaponSlot{tokens[1].weaponSlots[0], tokens[1].weaponSlots[1], tokens[0].weaponSlots[0], tokens[0].weaponSlots[1]}
			assert.True(t, reflect.DeepEqual(want, b.getSortedWeaponSlots(tokens)))
			assert.Equal(t, tt.want, want[0].initiative)
		})
	}
}

// Test_getBattleHostility verifies attack orders, retaliation, and allied participation.
func Test_getBattleHostility(t *testing.T) {
	tests := []struct {
		name    string
		primary BattleTarget
		armed   bool
		friend  bool
		want    map[int]map[int]bool
	}{
		{"retaliation includes a neutral defender", BattleTargetAny, true, false, map[int]map[int]bool{1: {2: true}, 2: {1: true}, 3: {}}},
		{"friend can join against the attacker", BattleTargetAny, true, true, map[int]map[int]bool{1: {2: true}, 2: {1: true}, 3: {1: true}}},
		{"no primary target does not initiate combat", BattleTargetNone, true, false, map[int]map[int]bool{1: {}, 2: {}, 3: {}}},
		{"unarmed orders do not initiate combat", BattleTargetAny, false, false, map[int]map[int]bool{1: {}, 2: {}, 3: {}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := map[int]*Player{}
			for number := 1; number <= 3; number++ {
				players[number] = testPlayer().WithNum(number)
				players[number].Relations = []PlayerRelationship{{Relation: PlayerRelationNeutral}, {Relation: PlayerRelationNeutral}, {Relation: PlayerRelationNeutral}}
			}
			players[1].Relations[1].Relation = PlayerRelationEnemy
			if tt.friend {
				players[3].Relations[1].Relation = PlayerRelationFriend
			}
			fleet := &Fleet{MapObject: MapObject{PlayerNum: 1}, battlePlan: &BattlePlan{PrimaryTarget: tt.primary, AttackWho: BattleAttackWhoEnemies}, Tokens: []ShipToken{{Quantity: 1, design: &ShipDesign{Spec: ShipDesignSpec{HasWeapons: tt.armed}}}}}
			assert.Equal(t, tt.want, getBattleHostility(players, []*Fleet{fleet}))
		})
	}
}

// Test_getBattleHostility_participation verifies initiation requirements and conflicting allied support.
func Test_getBattleHostility_participation(t *testing.T) {
	tests := []struct {
		name                                    string
		starbase, initiator, conflictingFriends bool
		want                                    map[int]map[int]bool
	}{
		{"starbase orders alone do not start a battle", true, false, false, map[int]map[int]bool{1: {}, 2: {}, 3: {}}},
		{"starbase orders apply when an armed fleet can initiate", true, true, false, map[int]map[int]bool{1: {2: true}, 2: {1: true}, 3: {}}},
		{"a spectator cannot support both opposing friends", false, true, true, map[int]map[int]bool{1: {2: true}, 2: {1: true}, 3: {}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := map[int]*Player{}
			for number := 1; number <= 3; number++ {
				players[number] = testPlayer().WithNum(number)
				players[number].Relations = []PlayerRelationship{{Relation: PlayerRelationNeutral}, {Relation: PlayerRelationNeutral}, {Relation: PlayerRelationNeutral}}
			}
			players[1].Relations[1].Relation = PlayerRelationEnemy
			if tt.conflictingFriends {
				players[3].Relations[0].Relation = PlayerRelationFriend
				players[3].Relations[1].Relation = PlayerRelationFriend
			}
			fleets := []*Fleet{{MapObject: MapObject{PlayerNum: 1}, Starbase: tt.starbase, battlePlan: &BattlePlan{PrimaryTarget: BattleTargetAny, AttackWho: BattleAttackWhoEnemies}, Tokens: []ShipToken{{Quantity: 1, design: &ShipDesign{Spec: ShipDesignSpec{HasWeapons: true}}}}}}
			if tt.initiator {
				fleets = append(fleets, &Fleet{MapObject: MapObject{PlayerNum: 2}, battlePlan: &BattlePlan{PrimaryTarget: BattleTargetAny, AttackWho: BattleAttackWhoEnemies}, Tokens: []ShipToken{{Quantity: 1, design: &ShipDesign{Spec: ShipDesignSpec{HasWeapons: true}}}}})
			}
			assert.Equal(t, tt.want, getBattleHostility(players, fleets))
		})
	}
}

// Test_selectBattleFleets verifies that only participating players' fleets join a battle.
func Test_selectBattleFleets(t *testing.T) {
	tests := []struct {
		name         string
		participants []int
		want         []int
	}{
		{"all players participate", []int{1, 2}, []int{1, 2, 2}},
		{"uninvolved player is excluded", []int{2}, []int{2, 2}},
		{"no participants", nil, []int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fleets := []*Fleet{{MapObject: MapObject{PlayerNum: 1}}, {MapObject: MapObject{PlayerNum: 2}}, {MapObject: MapObject{PlayerNum: 2}}}
			players := map[int]*Player{}
			for _, number := range tt.participants {
				players[number] = testPlayer().WithNum(number)
			}
			got := []int{}
			for _, fleet := range selectBattleFleets(fleets, players) {
				got = append(got, fleet.PlayerNum)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

// Test_battle_prepareCargo verifies mineral dumping, retained colonists, and starting mass and speed.
func Test_battle_prepareCargo(t *testing.T) {
	tests := []struct {
		name                   string
		dump                   bool
		wantCargo              Cargo
		wantDump               Mineral
		wantMass, wantMovement int
	}{
		{"loaded cargo affects speed", false, Cargo{Ironium: 140, Colonists: 20}, Mineral{}, 190, 3},
		{"dump minerals but retain colonists", true, Cargo{Colonists: 20}, Mineral{Ironium: 140}, 50, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRulesWithSeed(1)
			player := testPlayer().WithNum(1)
			design := &ShipDesign{Hull: SmallFreighter.Name, Spec: ShipDesignSpec{HullType: TechHullTypeFreighter, CargoCapacity: 200, Mass: 30, NumEngines: 1, Engine: Engine{IdealSpeed: 7}}}
			fleet := &Fleet{MapObject: MapObject{PlayerNum: 1}, Tokens: []ShipToken{{Quantity: 1, design: design}}, Cargo: Cargo{Ironium: 140, Colonists: 20}, Spec: FleetSpec{CargoCapacity: 200}, battlePlan: &BattlePlan{DumpCargo: tt.dump}}
			token := newBattleToken(&r, 1, Vector{1, 4}, &fleet.Tokens[0], *fleet.battlePlan, player)
			token.fleet = fleet
			b := &battle{rules: &r, tokens: []*battleToken{token}, fleets: []*Fleet{fleet}, record: newBattleRecord(1, None, Vector{}, nil)}
			b.prepareBattle()
			assert.Equal(t, tt.wantCargo, fleet.Cargo)
			assert.Equal(t, tt.wantDump, b.record.dumpedMinerals)
			assert.Equal(t, tt.wantMass, token.Mass)
			assert.Equal(t, tt.wantMovement, token.Movement)
			assert.Equal(t, token.Movement, b.record.Tokens[0].Movement)
		})
	}
}

// Test_battle_destroyedCargo verifies casualty losses of cargo and fuel and recovered wreckage.
func Test_battle_destroyedCargo(t *testing.T) {
	tests := []struct {
		name        string
		killed      int
		wantCargo   Cargo
		wantFuel    int
		wantSalvage Mineral
	}{
		{"one carrier loses half the load and fuel", 1, Cargo{Ironium: 50, Colonists: 10}, 50, Mineral{Ironium: 45}},
		{"destroyed fleet loses all cargo and fuel", 2, Cargo{}, 0, Mineral{Ironium: 90}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRulesWithSeed(1)
			design := &ShipDesign{Spec: ShipDesignSpec{CargoCapacity: 100, FuelCapacity: 100, Cost: Cost{Ironium: 30}}}
			fleet := &Fleet{Cargo: Cargo{Ironium: 100, Colonists: 20}, Fuel: 100, Tokens: []ShipToken{{Quantity: 2, design: design}}}
			token := &battleToken{BattleRecordToken: BattleRecordToken{PlayerNum: 1}, ShipToken: &fleet.Tokens[0], fleet: fleet, armor: 100}
			b := &battle{rules: &r, tokens: []*battleToken{token}, record: newBattleRecord(1, None, Vector{}, nil)}
			b.applyWeaponDamage(token, battleWeaponDamage{numDestroyed: tt.killed})
			assert.Equal(t, tt.wantCargo, fleet.Cargo)
			assert.Equal(t, tt.wantFuel, fleet.Fuel)
			assert.Equal(t, tt.wantSalvage, b.record.salvageMinerals)
			assert.Equal(t, Cargo{Ironium: 100 - tt.wantCargo.Ironium, Colonists: 20 - tt.wantCargo.Colonists}, b.record.Stats.CargoLostByPlayer[1])
		})
	}
}
